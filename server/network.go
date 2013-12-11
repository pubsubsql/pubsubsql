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

/*
func (b *networkBuffer) readHeader() {
} 
*/

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
	stopFlag bool
	quit     chan int
	context  *networkContext
}

func (n *network) addConnection(c *networkConnection) {
	n.mutex.Lock()
	if n.connections == nil {
		n.connections = make(map[uint64]*networkConnection)
	}
	n.connections[c.getConnectionId()] = c
	n.mutex.Unlock()
}

func (n *network) removeConnection(c *networkConnection) {
	n.mutex.Lock()
	delete(n.connections, c.getConnectionId())
	n.mutex.Unlock()
}

func (n *network) connectionCount() int {
	n.mutex.Lock()
	count := len(n.connections)
	n.mutex.Unlock()
	return count
}

func newNetwork(context *networkContext) *network {
	return &network{
		listener: nil,
		stopFlag: false,
		quit:     make(chan int),
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
		for {
			conn, err := n.listener.Accept()
			// stop was called
			if n.stopFlag {
				debug("stop was called")
				close(n.quit)
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
	n.stopFlag = true
	if n.listener != nil {
		n.listener.Close()
		<-n.quit
	}
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
		sender: &responseSender{sender: make(chan response, 1000), connectionId: connectionId},
	}
}

func (c *networkConnection) close() {
	if c.conn != nil {
		c.parent.removeConnection(c)
		c.conn.Close()
	}
}

func (c *networkConnection) getConnectionId() uint64 {
	return c.sender.connectionId
}

func (c *networkConnection) run() {
	go c.read()
	c.write()
}

func readHeader(conn net.Conn, bytes []byte) (uint32, error) {
	n, err := conn.Read(bytes[0:4])
	if err != nil {
		log.Println("Error reading from network connection: ", err.Error())
		return 0, err
	}
	if n < 4 {
		err = errors.New("Failed to read header.")
		log.Println("Error reading from network connection: ", err.Error())
		return 0, err
	}
	header := binary.LittleEndian.Uint32(bytes)
	return header, nil
}

func (c *networkConnection) readData(bytes []byte) error {
	left := len(bytes)
	n := 0
	var err error
	for left > 0 {
		if c.sender.quiter.IsQuit() {
			err = errors.New("Read was interupted by quit event.")
			return err
		}
		bytes = bytes[n:]
		n, err = c.conn.Read(bytes)
		if err != nil {
			log.Println("Error reading from network connection: ", err.Error())
			return err
		}
		left = left - n
	}
	return nil
}

func (c *networkConnection) read() {
	defer c.close()
	var max uint32 = 2048
	size := max + 4
	bytes := make([]byte, size, size)
	for {
		if c.sender.quiter.IsQuit() {
			return
		}
		// read header
		header, err := readHeader(c.conn, bytes)
		if err != nil {
			return
		}
		if header > max {
			log.Println("Error reading from network connection: message size exceedes max allowed value: ", header)
			return
		}
		// read data
		data := bytes[0:header]
		err = c.readData(data)
		if err != nil {
			return
		}
		// convert message to string
		message := string(data)
		debug(message)
		/*
			tokens := newTokens()
			if !lex(message, tokens) {
				// scan error
				// send error response to the client			
			}	
			// 
		*/

	}

}

func (c *networkConnection) write() {

}
