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

import "net"

type networkConnection struct {
	parent networkConnectionContainer
	conn   net.Conn
	quit   *Quitter
	router *requestRouter
	sender *responseSender
	dbConn *mysqlConnection
}

func newNetworkConnection(conn net.Conn, context *networkContext, connectionId uint64, parent networkConnectionContainer) *networkConnection {
	return &networkConnection {
		parent: parent,
		conn:   conn,
		quit:   context.quit,
		router: context.router,
		sender: newResponseSenderStub(connectionId),
		dbConn: newMysqlConnection(),
	}
}

func (this *networkConnection) remove() {
	this.parent.removeConnection(this)
}

func (this *networkConnection) getConnectionId() uint64 {
	return this.sender.connectionId
}

func (this *networkConnection) watchForQuit() {
	select {
	case <-this.sender.quit.GetChan():
	case <-this.quit.GetChan():
	}
	this.conn.Close()
	this.parent.removeConnection(this)
}

func (this *networkConnection) close() {
	this.sender.quit.Quit(0)
}

func (this *networkConnection) run() {
	go this.watchForQuit()
	go this.read()
	defer this.dbConn.disconnect()
	this.write()
}

func (this *networkConnection) Done() bool {
	// connection can be stopped because of global shutdown sequence
	// or response sender is full
	// or socket error
	return this.sender.quit.Done() || this.quit.Done()
}

func (this *networkConnection) route(header *netHeader, req request) {
	item := &requestItem {
		header: header,
		req:    req,
		sender: this.sender,
		dbConn: this.dbConn,
	}
	this.router.route(item)
}

func (this *networkConnection) read() {
	this.quit.Join()
	defer this.quit.Leave()
	reader := newNetHelper(this.conn, config.NET_READWRITE_BUFFER_SIZE)
	//
	var err error
	var message []byte
	var header *netHeader
	tokens := newTokens()
	for {
		err = nil
		if this.Done() {
			break
		}
		header, message, err = reader.readMessage()
		if err != nil {
			break
		}
		tokens.reuse()
		// parse and route the message
		lex(string(message), tokens)
		req := parse(tokens)
		this.route(header, req)
	}
	if err != nil && !this.Done() {
		logWarn("failed to read from client connection:", this.sender.connectionId, err.Error())
		// notify writer and sender that we are done
		this.sender.quit.Quit(0)
	}
}

func (this *networkConnection) write() {
	this.quit.Join()
	defer this.quit.Leave()
	writer := newNetHelper(this.conn, config.NET_READWRITE_BUFFER_SIZE)
	var err error
	for {
		select {
		case res := <-this.sender.sender:
			debug("response is ready to be send over tcp")
			// merge responses if applicable
			nextRes := this.sender.tryRecv()
			for nextRes != nil && res.merge(nextRes) {
				nextRes = this.sender.tryRecv()
			}
			// write messages in batches if applicable
			var msg []byte
			more := true
			for err == nil && more {
				if this.Done() {
					return
				}
				msg, more = res.toNetworkReadyJSON()
				err = writer.writeMessage(msg)
				if err != nil {
					break
				}
				if !more && nextRes != nil {
					res = nextRes
					nextRes = nil
					more = true
				}
			}
			if err != nil && !this.Done() {
				logWarn("failed to write to client connection:", this.sender.connectionId, err.Error())
				// notify reader and sender that we are done
				this.sender.quit.Quit(0)
				return
			}
		case <-this.quit.GetChan():
			debug("on write stop")
			return
		case <-this.sender.quit.GetChan():
			debug("on write connection stop")
			return
		}
	}
}
