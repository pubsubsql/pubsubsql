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

import "sync/atomic"
import "time"

var debug_flag bool = true

func debug(str string) {
	if debug_flag {
		println("debug: " + str)
	}
}

// Stoper implements shutdown protocol to make sure that all active goroutines exit gracefully.
type Stoper struct {
	counter int64
	channel chan int
	stoping bool
}

// Stoper factory.
func NewStoper() *Stoper {
	return &Stoper{
		counter: 0,
		channel: make(chan int),
		stoping: false,
	}
}

// Enter starts goroutine participation in shutdown protocol.
func (s *Stoper) Enter() {
	atomic.AddInt64(&s.counter, 1)
}

// Leave notifies that participating goroutine gracesfully exited.
// Should be called with defer symantics.
func (s *Stoper) Leave() {
	atomic.AddInt64(&s.counter, -1)
}

// GetChan returns channel to be used in select {} in order to react to Stop event.
func (s *Stoper) GetChan() chan int {
	return s.channel
}

// Stop notifies all participating goroutines that shutdown protocol is in progress
// and waits for all go routines to exit until timeouti.
// Returns false when timed out.
func (s *Stoper) Stop(timeout time.Duration) bool {
	s.stoping = true
	close(s.channel)
	return s.Wait(timeout)
}

func (s *Stoper) Wait(timeout time.Duration) bool {
	t := time.Now()
	for atomic.LoadInt64(&s.counter) > 0 {
		time.Sleep(time.Millisecond * 10)
		if time.Since(t) > timeout {
			return false
		}
	}
	return true
}

func (s *Stoper) IsStoping() bool {
	return s.stoping
}

// Counter returns number of of participating goroutines.
func (s *Stoper) Counter() int64 {
	return atomic.LoadInt64(&s.counter)
}
