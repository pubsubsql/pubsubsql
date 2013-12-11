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

const (
	quitBySender int = iota
	quitByNetReader
	quitByNetWriter
	quitByNetwork
)

// responseSender is responsible for channeling reponses to client connection
type responseSender struct {
	sender       chan response // channel to publish responses to
	connectionId uint64
	quiter       *QuitChan
}

// factory
func newResponseSenderStub(connectionId uint64) *responseSender {
	return &responseSender{
		sender:       make(chan response, 1000),
		connectionId: connectionId,
		quiter:       NewQuitChan(),
	}
}

func (s *responseSender) send(r response) bool {
	select {
	case s.sender <- r:
		if !s.quiter.IsQuit() {
			return true
		}
		debug("sender IsQuit")
	case <-s.quiter.Chan():
		debug("sender quit")
	default:
		debug("sender queue is full ")
		// notify client connection that it needs to close due to inability to 
		// recv responses in a timely manner
		s.quiter.Quit(quitBySender)
	}
	return false
}

func (s *responseSender) tryRecv() response {
	select {
	case r := <-s.sender:
		return r
	default:
		return nil
	}
	return nil
}

func (s *responseSender) recv() response {
	return <-s.sender
}
