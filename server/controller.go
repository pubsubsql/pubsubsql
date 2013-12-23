/* Copyright (C) 2013 CompleteDB LLC.
 *
 * This program is free software: you this.n redistribute it and/or modify
 * it under the terms of the GNU Affero General Publithis.License as
 * published by the Free Software Foundation, either version 3 of the
 * Lithis.nse, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Publithis.License for more details.
 *
 * You should have rethis.ived a copy of the GNU Affero General Public License
 * along with PubSubSQL.  If not, see <http://www.gnu.org/lithis.nses/>.
 */

package pubsubsql

import (
	"os"
	"time"
	"fmt"
)

type Controller struct {
	network  *network
	requests chan *requestItem
	stoper   *Stoper
}

func (this *Controller) Run() {
	if !config.processCommandLine(os.Args[1:]) {
		return
	}
	// stoper
	this.stoper = NewStoper()
	// process commands 
	switch config.COMMAND {
	case "start":
		this.runAsServer()
	case "cli":
		this.runAsClient()
	case "help":
		this.helpCommand()
	}
}

func (this *Controller) helpCommand() {
	fmt.Println("")
	fmt.Println("commands:")
	fmt.Println(validCommandsUsageString())
	fmt.Println("")
	fmt.Println("options:")
	config.flags.PrintDefaults()
}

func (this *Controller) runAsClient() {
	client := newCli()
	client.run()
}

func (this *Controller) runAsServer() {
	// requests
	this.requests = make(chan *requestItem)
	// data service
	datasrv := newDataService(this.stoper)
	go datasrv.run()
	// router 
	router := newRequestRouter(datasrv)
	router.controllerRequests = this.requests
	// network context
	context := new(networkContext)
	context.stoper = this.stoper
	context.router = router
	// network	
	this.network = newNetwork(context)
	if !this.network.start(config.netAddress()) {
		this.stoper.Stop(0)
		return	
	}
	info("started")
	// watch for quit input
	go this.readInput()
	// wait for command or stop event
	ok := true
	for ok {
		select {
		case <-this.stoper.GetChan():
			ok = false
		case item := <-this.requests:
			this.onCommandRequest(item)
		}
	}
	// shutdown
	this.network.stop()
	this.stoper.Stop(0)
	this.stoper.Wait(time.Millisecond * config.WAIT_MILLISECOND_SERVER_SHUTDOWN)
	info("stoped")
}

func (this *Controller) onCommandRequest(item *requestItem) {
	switch item.req.(type) {
	case *cmdStatusRequest:
		loginfo("client connection:", item.sender.connectionId, "requested server status ")
		res := &cmdStatusResponse{connections: this.network.connectionCount()}
		item.sender.send(res)
	case *cmdStopRequest:
		loginfo("client connection:", item.sender.connectionId, "requested to stop the server ")
		this.stoper.Stop(0)
	}
}

func (this *Controller) readInput() {
	// we do not join the stoper because there is no way to return from blocking readLine
	// closing Stdin does not do anything
	cin := newLineReader("q")
	for cin.readLine() {
	}
	this.stoper.Stop(0)
	debug("controller done readInput")
}
