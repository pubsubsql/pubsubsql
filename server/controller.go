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

package pubsubsql

import "time"
import "os"

type Controller struct {
	stoper *Stoper	
	network *network	
}

func (c *Controller) Run() {
	if !config.processCommandLine(os.Args[1:]) {
		return
	}
	// stoper
	c.stoper = NewStoper()
	// commands
	switch config.COMMAND {
	case "start":
		c.runAsServer()	
	case "connect":
		c.runAsClient()
	case "help":
		println("help")
	}
}

func (c *Controller) runAsClient() {
	client := newCli()
	client.run()
}

func (c *Controller) runAsServer() {
	// dataservice
	datasrv := newDataService(config.CHAN_TABLE_REQUESTS_BUFFER_SIZE, c.stoper)
	go datasrv.run()
	// router 
	router := newRequestRouter(datasrv)
	// network context
	context := new(networkContext)
	context.stoper = c.stoper
	context.router = router
	// network	
	c.network = newNetwork(context)
	c.network.start(config.netAddress())
	println("started")						
	// wait for quit input
	rd := newLineReader("q")
	for rd.readLine() {
	}
	// shutdown
	c.network.stop()
	c.stoper.Stop(0)
	c.stoper.Wait(time.Millisecond * 3000)
	println("stoped")	
}

