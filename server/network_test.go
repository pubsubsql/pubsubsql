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

func TestNetworkWriteRead(t *testing.T) {
	context := newNetworkContextStub()
	s := context.stoper
	n := newNetwork(context)
	n.start("localhost:54321")
	c, err := net.Dial("tcp", "localhost:54321")
	if err != nil {
		t.Error(err)
	}
	// send valid message get result
	rw := newNetMessageReaderWriter(c, nil)
	message := []byte("key stocks ticker")
	err = rw.writeHeaderAndMessage(message)
	if err != nil {
		t.Error(err)
	}
	message, err = rw.readMessage()
	if err != nil {
		t.Error(err)
	}
	debug(string(message))
	if n.connectionCount() != 1 {
		t.Error("Expected 1 network connection")
	}
	// send invalid message get result
	message = []byte("bla bla bla")
	err = rw.writeHeaderAndMessage(message)
	if err != nil {
		t.Error(err)
	}
	message, err = rw.readMessage()
	if err != nil {
		t.Error(err)
	}
	debug(string(message))
	if n.connectionCount() != 1 {
		t.Error("Expected 1 network connection")
	}
	// close connection
	c.Close()
	time.Sleep(time.Millisecond * 500)
	if n.connectionCount() > 0 {
		t.Error("Expected 0 network connection")
	}
	// shutdown
	s.Stop(0)
	n.stop()
	s.Wait(time.Millisecond * 500)
}
