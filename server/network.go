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

import (
	"encoding/binary"
	"errors"
	"net"
	"strconv"
	"sync"
)

// networkContext
type networkContext struct {
	stoper *Stoper
	router *requestRouter
}

func newNetworkContextStub() *networkContext {
	stoper := NewStoper()
	//
	datasrv := newDataService(stoper)
	go datasrv.run()
	//
	router := newRequestRouter(datasrv)
	//
	context := new(networkContext)
	context.stoper = stoper
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

func (this *network) addConnection(netconn *networkConnection) {
	if this.context.stoper.Stoped() {
		return
	}
	this.mutex.Lock()
	if this.connections == nil {
		this.connections = make(map[uint64]*networkConnection)
	}
	this.connections[netconn.getConnectionId()] = netconn
	this.mutex.Unlock()
	loginfo("new client connection id:", strconv.FormatUint(netconn.getConnectionId(), 10))
}

func (this *network) removeConnection(netconn *networkConnection) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	if this.connections != nil {
		delete(this.connections, netconn.getConnectionId())
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
		logerror("Failed to listen to incoming connections ", err.Error())
		return false
	}
	this.listener = listener
	var connectionId uint64 = 0
	// accept connections
	acceptor := func() {
		stoper := this.context.stoper
		stoper.Join()
		defer stoper.Leave()
		for {
			conn, err := this.listener.Accept()
			// stop was called
			if stoper.Stoped() {
				return
			}
			if err == nil {
				connectionId++
				netconn := newNetworkConnection(conn, this.context, connectionId, this)
				this.addConnection(netconn)
				go netconn.run()
			} else {
				logerror("failed to accept client connection", err.Error())
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

//

type networkConnection struct {
	parent networkConnectionContainer
	conn   net.Conn
	stoper *Stoper
	router *requestRouter
	sender *responseSender
}

func newNetworkConnection(conn net.Conn, context *networkContext, connectionId uint64, parent networkConnectionContainer) *networkConnection {
	return &networkConnection{
		parent: parent,
		conn:   conn,
		stoper: context.stoper,
		router: context.router,
		sender: newResponseSenderStub(connectionId),
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
	case <-this.sender.connectionStoper.GetChan():
	case <-this.stoper.GetChan():
	}
	this.conn.Close()
	this.parent.removeConnection(this)
}

func (this *networkConnection) close() {
	this.sender.connectionStoper.Stop(0)
}

func (this *networkConnection) run() {
	go this.watchForQuit()
	go this.read()
	this.write()
}

func (this *networkConnection) Stoped() bool {
	// connection can be stoped becuase of global shutdown sequence
	// or response sender is full
	// or socket error
	return this.sender.connectionStoper.Stoped() || this.stoper.Stoped()
}

// message reader
type netMessageReaderWriter struct {
	conn   net.Conn
	bytes  []byte
	stoper IStoper
}

func newNetMessageReaderWriter(conn net.Conn, stoper IStoper) *netMessageReaderWriter {
	return &netMessageReaderWriter{
		conn:   conn,
		bytes:  make([]byte, 2048, 2048),
		stoper: stoper,
	}
}

func (this *netMessageReaderWriter) Stoped() bool {
	return this.stoper != nil && this.stoper.Stoped()
}

func (this *netMessageReaderWriter) writeMessage(bytes []byte) error {
	leftToWrite := len(bytes)
	for {
		if this.Stoped() {
			err := errors.New("Write was interupted by quit event.")
			return err
		}
		written, err := this.conn.Write(bytes)
		if err != nil {
			return err
		}
		leftToWrite -= written
		if leftToWrite == 0 {
			break
		}
		bytes = bytes[written:]
	}
	return nil
}

// for cli
func (this *netMessageReaderWriter) writeHeaderAndMessage(bytes []byte) error {
	header := make([]byte, HEADER_SIZE, HEADER_SIZE)
	binary.LittleEndian.PutUint32(header, uint32(len(bytes)))
	err := this.writeMessage(header)
	if err != nil {
		return err
	}
	return this.writeMessage(bytes)
}

func (this *netMessageReaderWriter) readMessage() ([]byte, error) {
	// header
	read, err := this.conn.Read(this.bytes[0:HEADER_SIZE])
	if err != nil {
		return nil, err
	}
	if read < HEADER_SIZE {
		err = errors.New("Failed to read header.")
		return nil, err
	}
	header := binary.LittleEndian.Uint32(this.bytes)
	// prepare buffer
	if len(this.bytes) < int(header) {
		this.bytes = make([]byte, header, header)
	}
	// message
	bytes := this.bytes[:header]
	left := len(bytes)
	message := bytes
	read = 0
	for left > 0 {
		if this.Stoped() {
			err = errors.New("Read was interupted by quit event.")
			return nil, err
		}
		bytes = bytes[read:]
		read, err = this.conn.Read(bytes)
		if err != nil {
			return nil, err
		}
		left -= read
	}
	return message, nil
}

func (c *networkConnection) route(req request) {
	item := &requestItem{
		req:    req,
		sender: c.sender,
	}
	c.router.route(item)
}

func (this *networkConnection) read() {
	this.stoper.Join()
	defer this.stoper.Leave()
	reader := newNetMessageReaderWriter(this.conn, this)
	//
	var err error
	var message []byte
	for {
		err = nil
		if this.Stoped() {
			break
		}
		message, err = reader.readMessage()
		if err != nil {
			break
		}
		// parse and route the message
		tokens := newTokens()
		lex(string(message), tokens)
		req := parse(tokens)
		this.route(req)
	}
	if err != nil && !this.Stoped() {
		logerror("failed to read from client connection: ", this.sender.connectionId)
		logerror(err.Error())
		// notify writer and sender that we are done
		this.sender.connectionStoper.Stop(0)
	}
}

func (this *networkConnection) write() {
	this.stoper.Join()
	defer this.stoper.Leave()
	writer := newNetMessageReaderWriter(this.conn, this)
	for {
		select {
		case res := <-this.sender.sender:
			if this.Stoped() {
				return
			}
			err := writer.writeMessage(res.toNetworkReadyJSON())
			if err != nil {
				if !this.Stoped() {
					logerror("failed to write to client connection: ", this.sender.connectionId)
					logerror(err.Error())
					// notify reader and sender that we are done
					this.sender.connectionStoper.Stop(0)
				}
				return
			}
		case <-this.stoper.GetChan():
			debug("on write stop")
			return
		case <-this.sender.connectionStoper.GetChan():
			debug("on write connection stop")
			return
		}
	}
}
