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
	"net"
	"strconv"
	"sync"
)

// networkContext
type networkContext struct {
	quit   *Quitter
	router *requestRouter
}

func newNetworkContextStub() *networkContext {
	quit := NewQuitter()
	//
	datasrv := newDataService(quit)
	go datasrv.run()
	//
	router := newRequestRouter(datasrv)
	//
	context := new(networkContext)
	context.quit = quit
	context.router = router
	//
	return context
}

// network
type networkConnectionContainer interface {
	removeConnection(*networkConnection)
}

type network struct {
	networkConnectionContainer
	mutex       sync.Mutex
	connections map[uint64]*networkConnection
	listener    net.Listener
	context     *networkContext
}

func (this *network) addConnection(netConn *networkConnection) {
	if this.context.quit.Done() {
		return
	}
	this.mutex.Lock()
	defer this.mutex.Unlock()
	if this.connections == nil {
		this.connections = make(map[uint64]*networkConnection)
	}
	this.connections[netConn.getConnectionId()] = netConn
	logInfo("new client connection id:", strconv.FormatUint(netConn.getConnectionId(), 10))
}

func (this *network) removeConnection(netConn *networkConnection) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	if this.connections != nil {
		delete(this.connections, netConn.getConnectionId())
	}
}

func (this *network) connectionCount() int {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	count := len(this.connections)
	return count
}

func (this *network) closeConnections() {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	for _, c := range this.connections {
		c.close()
	}
	this.connections = nil
}

func newNetwork(context *networkContext) *network {
	return &network{
		listener: nil,
		context:  context,
	}
}

func (this *network) start(address string) bool {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		logError("Failed to listen for incoming connections ", err.Error())
		return false
	}
	// host, port := net.SplitHostPort(address)
	logInfo("listening for incoming connections on ", address)
	this.listener = listener
	var connectionId uint64 = 0
	// accept connections
	acceptor := func() {
		quit := this.context.quit
		quit.Join()
		defer quit.Leave()
		for {
			conn, err := this.listener.Accept()
			// stop was called
			if quit.Done() {
				return
			}
			if err == nil {
				connectionId++
				netConn := newNetworkConnection(conn, this.context, connectionId, this)
				this.addConnection(netConn)
				go netConn.run()
			} else {
				logError("failed to accept client connection", err.Error())
			}
		}
	}
	go acceptor()
	return true
}

func (this *network) stop() {
	if this.listener != nil {
		this.listener.Close()
	}
	this.closeConnections()
}
