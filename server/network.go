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

type network struct {
	listener net.Listener
	stopFlag bool
	quit     chan int
}

func newNetwork() *network {
	return &network{
		listener: nil,
		stopFlag: false,
		quit:     make(chan int),
	}
}

func (n *network) start(address string) bool {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Println("Error listening to incoming connections ", err.Error())
		return false
	}
	n.listener = listener
	// accept connections
	acceptor := func() {
		for {
			conn, err := n.listener.Accept()
			if n.stopFlag {
				close(n.quit)
				return
			}
			if err == nil {
				go n.handleConnection(conn)
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

func (n *network) handleConnection(conn net.Conn) {

}
