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

import "testing"
import "time"
import "net"
import "strconv"

func TestNetworkStartStop(t *testing.T) {
	context := newNetworkContextStub()
	s := context.stoper
	n := newNetwork(context)
	if !n.start("localhost:54321") {
		t.Error(`network.start("localhost:54321") failed`)
	}
	// shutdown
	s.Stop(0)
	n.stop()
	s.Wait(time.Millisecond * 1000)
}

func TestNetworkConnections(t *testing.T) {
	context := newNetworkContextStub()
	s := context.stoper
	n := newNetwork(context)
	n.start("localhost:54321")
	c, err := net.Dial("tcp", "localhost:54321")
	if err != nil {
		t.Error(err)
	}
	time.Sleep(time.Millisecond * 1000)
	if n.connectionCount() != 1 {
		t.Error("Expected 1 network connection")
	}
	// shutdown
	s.Stop(0)
	n.stop()
	s.Wait(time.Millisecond * 1000)
	c.Close()
}

func validateWriteRead(t *testing.T, conn net.Conn, message string) {
	rw := newNetMessageReaderWriter(conn, nil)
	bytes := []byte(message)
	err := rw.writeHeaderAndMessage(bytes)
	if err != nil {
		t.Error(err)
	}
	bytes, err = rw.readMessage()
	if err != nil {
		t.Error(err)
	}
	debug(string(bytes))
}

func validateRead(t *testing.T, conn net.Conn) {
	rw := newNetMessageReaderWriter(conn, nil)
	bytes, err := rw.readMessage()
	if err != nil {
		t.Error(err)
	}
	debug(string(bytes))
}

func validateConnect(t *testing.T, address string) net.Conn {
	c, err := net.Dial("tcp", address)
	if err != nil {
		t.Error(err)
	}
	return c
}

func TestNetworkWriteRead(t *testing.T) {
	context := newNetworkContextStub()
	address := "localhost:54321"
	s := context.stoper
	n := newNetwork(context)
	n.start(address)
	c := validateConnect(t, address)
	// send valid message get result
	validateWriteRead(t, c, "key stocks ticker")
	validateWriteRead(t, c, "bla bla bla")
	validateWriteRead(t, c, "insert into stocks (ticker, bid, ask) values (IBM,123,124)")
	validateWriteRead(t, c, "insert into stocks (ticker, bid, ask) values (MSFT,37,38.45)")
	validateWriteRead(t, c, "select * from stocks")
	validateWriteRead(t, c, "key stocks ticker")
	// test pubsub
	c2 := validateConnect(t, address)
	validateWriteRead(t, c2, "subscribe * from stocks")
	// on add
	validateRead(t, c2)
	validateWriteRead(t, c, "insert into stocks (ticker, bid, ask) values (ORCL,37,38.45)")
	// on insert
	validateRead(t, c2)
	//
	if n.connectionCount() != 2 {
		t.Error("Expected 1 network connection")
	}
	// close connections
	c.Close()
	time.Sleep(time.Millisecond * 60)
	if n.connectionCount() != 1 {
		t.Error("Expected 1 network connection")
	}
	c2.Close()
	time.Sleep(time.Millisecond * 60)
	if n.connectionCount() != 0 {
		t.Error("Expected 0 network connection")
	}
	// shutdown
	s.Stop(0)
	n.stop()
	s.Wait(time.Millisecond * 500)
}

func TestNetworMultiInsert(t *testing.T) {
	context := newNetworkContextStub()
	address := "localhost:54321"
	s := context.stoper
	n := newNetwork(context)
	n.start(address)
	c := validateConnect(t, address)

	// subscribe
	validateWriteRead(t, c, "key stocks ticker")
	validateWriteRead(t, c, "subscribe * from stocks")

	// insert bunch of records
	idblock := 1
	insertsPerConnection := 10
	totalConnections := 10
	for i := 0; i < totalConnections; i++ {
		idblock += 100000
		go func(cid int, tickerid int) {
			c2 := validateConnect(t, address)
			for j := 0; j < insertsPerConnection; j++ {
				ticker := strconv.Itoa(tickerid)
				tickerid++
				validateWriteRead(t, c2, "insert into stocks (ticker, bid, ask) values ( "+ticker+",10,10)")
			}

		}(i, idblock)
	}
	// read inserted published records 
	for j := 0; j < (insertsPerConnection * totalConnections); j++ {
		validateRead(t, c)
	}
	c.Close()
	// shutdown
	s.Stop(0)
	n.stop()
	s.Wait(time.Millisecond * 500)
}
