/* Copyright (C) 2013 CompleteDB LLC.
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with PubSubSQL.  If not, see <http://www.gnu.org/licenses/>.
 */

package server

import (
	"runtime"
	"fmt"
	"os"
	"time"
)

// Controller is a container that initializes, binds and controls server components.
type Controller struct {
	network			*network
	requests chan	*requestItem
	quit			*Quitter
}

// Run is a main server entry function. It processes command line options and runs the server in the appropriate mode.
func (this *Controller) Run() {
	runtime.GOMAXPROCS(runtime.NumCPU() - 1)
	if !config.processCommandLine(os.Args[1:]) {
		return
	}
	this.quit = NewQuitter()
	// process commands
	switch config.COMMAND {
	case "help":
		this.displayHelp()
	case "cli":
		this.runAsClient()
	case "start":
		this.runAsServer()
	case "stop":
		this.runOnce("stop")
	}
}

// displayHelp displays help to the cli user.
func (this *Controller) displayHelp() {
	fmt.Println("")
	fmt.Println("commands:")
	fmt.Println(validCommandsUsageString())
	fmt.Println("")
	fmt.Println("options:")
	config.flags.PrintDefaults()
}

// runAsClient runs the program in cli mode.
func (this *Controller) runAsClient() {
	client := newCli()
	// start cli event loop
	client.run()
}

// run command once
func (this *Controller) runOnce(command string) {
	client := newCli()
	client.runOnce(command)
}

// runAsServer runs the program in server mode.
func (this *Controller) runAsServer() {
	// initialize server components
	// requests
	this.requests = make(chan *requestItem)
	// data service
	dataService := newDataService(this.quit)
	go dataService.run()
	// router
	router := newRequestRouter(dataService)
	router.controllerRequests = this.requests
	// network context
	context := new(networkContext)
	context.quit = this.quit
	context.router = router
	// network
	this.network = newNetwork(context)
	if !this.network.start(config.netAddress()) {
		this.quit.Quit(0)
		return
	}
	info("started")
	// watch for quit (q) input
	go this.readInput()
	// wait for command to process or stop event
LOOP:
	for {
		select {
		case <-this.quit.GetChan():
			break LOOP
		case item := <-this.requests:
			this.onCommandRequest(item)
		}
	}
	// shutdown
	this.network.stop()
	this.quit.Quit(0)
	this.quit.Wait(time.Millisecond * config.WAIT_MILLISECOND_SERVER_SHUTDOWN)
	info("stopped")
}

// readInput reads a command line input from the standard until quit (q) input.
func (this *Controller) readInput() {
	cin := newLineReader("q")
	for cin.readLine() {
	}
	this.quit.Quit(0)
	debug("controller done readInput")
}

// onCommandRequest processes request from a connected client, sending respond back to the client.
func (this *Controller) onCommandRequest(item *requestItem) {
	switch item.req.(type) {
	case *cmdStatusRequest:
		logInfo("client connection:", item.sender.connectionId, "requested server status")
		if item.req.isStreaming() {
			return
		}
		res := newCmdStatusResponse(this.network.connectionCount())
		res.requestId = item.getRequestId()
		item.sender.send(res)
	case *cmdStopRequest:
		logInfo("client connection:", item.sender.connectionId, "requested to stop the server")
		this.quit.Quit(0)
	case *mysqlConnectRequest:
		logInfo("client connection:", item.sender.connectionId, "requested mysql connect")
		if item.req.isStreaming() {
			return
		}
		request := item.req.(*mysqlConnectRequest)
		response := newCmdMysqlConnectResponse(request)
		response.requestId = item.getRequestId()
		//
		item.dbConn.connect(request.address)
		if item.dbConn.hasError() {
			response.error = item.dbConn.getLastError()
		}
		//
		item.sender.send(response)
	case *mysqlDisconnectRequest:
		logInfo("client connection:", item.sender.connectionId, "requested mysql diconnect")
		if item.req.isStreaming() {
			return
		}
		request := item.req.(*mysqlDisconnectRequest)
		response := newCmdMysqlDisconnectResponse(request)
		response.requestId = item.getRequestId()
		//
		item.dbConn.disconnect()
		if item.dbConn.hasError() {
			response.error = item.dbConn.getLastError()
		}
		//
		item.sender.send(response)
	case *mysqlStatusRequest:
		logInfo("client connection:", item.sender.connectionId, "requested mysql status")
		if item.req.isStreaming() {
			return
		}
		request := item.req.(*mysqlStatusRequest)
		response := newCmdMysqlStatusResponse(request)
		response.requestId = item.getRequestId()
		//
		connected := item.dbConn.isConnected()
		response.setOnline(connected)
		if item.dbConn.hasError() {
			response.error = item.dbConn.getLastError()
		}
		//
		item.sender.send(response)
	case *mysqlTablesRequest:
		logInfo("client connection:", item.sender.connectionId, "requested mysql tables")
		if item.req.isStreaming() {
			return
		}
		request := item.req.(*mysqlTablesRequest)
		response := newCmdMysqlTablesResponse(request)
		response.requestId = item.getRequestId()
		//
		tables := item.dbConn.findTables()
		if item.dbConn.hasError() {
			response.error = item.dbConn.getLastError()
		} else {
			response.tables = tables
		}
		//
		item.sender.send(response)
	}
}
