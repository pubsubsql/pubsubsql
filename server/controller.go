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
)

type Controller struct {
	stoper  *Stoper
	network *network
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
		println("help")
	case "status":
		println("status")
	}
}

func (this *Controller) runAsClient() {
	client := newCli()
	client.run()
}

func (this *Controller) runAsServer() {
	// data service
	datasrv := newDataService(this.stoper)
	go datasrv.run()
	// router 
	router := newRequestRouter(datasrv)
	// network context
	context := new(networkContext)
	context.stoper = this.stoper
	context.router = router
	// network	
	this.network = newNetwork(context)
	this.network.start(config.netAddress())
	info("started")
	// wait for quit input
	cin := newLineReader("q")
	for cin.readLine() {
	}
	// shutdown
	this.network.stop()
	this.stoper.Stop(0)
	this.stoper.Wait(time.Millisecond * 3000)
	info("stoped")
}
