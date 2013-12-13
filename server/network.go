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

import "net"
import "log"
import "sync"
import "encoding/binary"
import "errors"

//import "encoding/binary"

type tokensProducerConsumer struct {
	idx    int
	tokens []*token
}

func newTokens() *tokensProducerConsumer {
	return &tokensProducerConsumer{
		idx:    0,
		tokens: make([]*token, 0, 30),
	}
}

func reuseTokens(pc *tokensProducerConsumer) {
	pc.idx = 0
}

func (c *tokensProducerConsumer) Consume(t *token) {
	c.tokens = append(c.tokens, t)
}

func (p *tokensProducerConsumer) Produce() *token {
	if p.idx >= len(p.tokens) {
		return &token{
			typ: tokenTypeEOF,
		}
	}
	t := p.tokens[p.idx]
	p.idx++
	return t
}

// networkContext
type networkContext struct {
	stoper *Stoper
	router *requestRouter
}

func newNetworkContextStub() *networkContext {
	stoper := NewStoper()
	//
	datasrv := newDataService(1000, stoper)
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
	//
	listener net.Listener
	context  *networkContext
}

func (n *network) addConnection(c *networkConnection) {
	if n.context.stoper.IsStoping() {
		return
	}
	n.mutex.Lock()
	if n.connections == nil {
		n.connections = make(map[uint64]*networkConnection)
	}
	n.connections[c.getConnectionId()] = c
	n.mutex.Unlock()
}

func (n *network) removeConnection(c *networkConnection) {
	n.mutex.Lock()
	if n.connections != nil {
		delete(n.connections, c.getConnectionId())
	}
	n.mutex.Unlock()
}

func (n *network) connectionCount() int {
	n.mutex.Lock()
	count := len(n.connections)
	n.mutex.Unlock()
	return count
}

func (n *network) closeConnections() {
	n.mutex.Lock()
	for _, c := range n.connections {
		c.close()
	}
	n.connections = nil
	n.mutex.Unlock()
}

func newNetwork(context *networkContext) *network {
	return &network{
		listener: nil,
		context:  context,
	}
}

func (n *network) start(address string) bool {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Println("Error listening to incoming connections ", err.Error())
		return false
	}
	n.listener = listener
	var connectionId uint64 = 0
	// accept connections
	acceptor := func() {
		s := n.context.stoper
		s.Enter()
		defer s.Leave()
		for {
			conn, err := n.listener.Accept()
			// stop was called
			if s.IsStoping() {
				debug("stop was called")
				return
			}
			if err == nil {
				debug("new network connection")
				connectionId++
				c := newNetworkConnection(conn, n.context, connectionId, n)
				n.addConnection(c)
				go c.run()
			} else {
				log.Println("Error accepting client connection", err.Error())
			}
		}
	}
	go acceptor()
	//	
	return true
}

func (n *network) stop() {
	if n.listener != nil {
		n.listener.Close()
	}
	n.closeConnections()
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

func (c *networkConnection) remove() {
	c.parent.removeConnection(c)
}

func (c *networkConnection) getConnectionId() uint64 {
	return c.sender.connectionId
}

func (c *networkConnection) watchForQuit() {
	select {
		case <-c.sender.quiter.GetChan():
		case <-c.stoper.GetChan():	
	}
	c.conn.Close()
	c.parent.removeConnection(c)
}

func (c *networkConnection) close() {
	c.sender.quiter.Quit(0)
} 

func (c *networkConnection) run() {
	go c.watchForQuit()	
	go c.read()
	c.write()
}

func (c *networkConnection) shouldStop() bool {
	return c.sender.quiter.IsQuit() || c.stoper.IsStoping()
}

type IStoper interface {
	shouldStop() bool
}

// message reader
type netMessageReaderWriter struct {
	conn  net.Conn
	bytes []byte
	s     IStoper
}

func newNetMessageReaderWriter(conn net.Conn, s IStoper) *netMessageReaderWriter {
	return &netMessageReaderWriter{
		conn:  conn,
		bytes: make([]byte, 2048, 2048),
		s:     s,
	}
}

func (r *netMessageReaderWriter) shouldStop() bool {
	return r.s != nil && r.s.shouldStop()
}

func (r *netMessageReaderWriter) writeMessage(bytes []byte) error {
	l := len(bytes)
	for {
		if r.shouldStop() {
			err := errors.New("Write was interupted by quit event.")
			return err
		}
		n, err := r.conn.Write(bytes)
		if err != nil {
			return err
		}
		if l == 0 {
			break
		}
		l = l - n
		bytes = bytes[n:]
	}
	return nil
}

// for testing
func (r *netMessageReaderWriter) writeHeaderAndMessage(bytes []byte) error {
	header := make([]byte, 4, 4)
	binary.LittleEndian.PutUint32(header, uint32(len(bytes)))
	err := r.writeMessage(header)
	if err != nil {
		return err
	}
	return r.writeMessage(bytes)
}

func (r *netMessageReaderWriter) readMessage() ([]byte, error) {
	// header
	n, err := r.conn.Read(r.bytes[0:4])
	if err != nil {
		return nil, err
	}
	if n < 4 {
		err = errors.New("Failed to read header.")
		return nil, err
	}
	header := binary.LittleEndian.Uint32(r.bytes)
	// prepare buffer
	if len(r.bytes) < int(header) {
		r.bytes = make([]byte, header, header)
	}
	// message
	bytes := r.bytes[:header]
	left := len(bytes)
	message := bytes
	n = 0
	for left > 0 {
		if r.shouldStop() {
			err = errors.New("Read was interupted by quit event.")
			return nil, err
		}
		bytes = bytes[n:]
		n, err = r.conn.Read(bytes)
		if err != nil {
			return nil, err
		}
		left = left - n
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

func (c *networkConnection) read() {
	s := c.stoper
	s.Enter()
	defer s.Leave()
	reader := newNetMessageReaderWriter(c.conn, c)
	//
	var err error
	var message []byte
	for {
		err = nil
		if c.shouldStop() {
			break
		}
		message, err = reader.readMessage()
		if err != nil {
			break
		}
		// parse and route the message
		pc := newTokens()
		lex(string(message), pc)
		req := parse(pc)
		c.route(req)
	}
	if !c.shouldStop() {
		if err != nil {
			log.Println(err.Error())
			// notify writer and sender that we are done
			c.sender.quiter.Quit(quitByNetReader)
		}
	}
}

func (c *networkConnection) write() {
	s := c.stoper
	s.Enter()
	defer s.Leave()
	writer := newNetMessageReaderWriter(c.conn, c)
	for {
		select {
		case res := <-c.sender.sender:
			if c.shouldStop() {
				return
			}
			err := writer.writeMessage(res.toNetworkReadyJSON())
			if err != nil {
				if !c.shouldStop() {
					log.Println(err.Error())
					// notify reader and sender that we are done
					c.sender.quiter.Quit(quitByNetWriter)
				}
				return
			}
		case <-s.GetChan():
			debug("stop event on write")
			return
		case <-c.sender.quiter.GetChan():
			debug("stop event on write")
			return
		}
	}
}
