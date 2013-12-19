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
	"sync"
	"sync/atomic"
	"time"
)

// Stoper implements shutdown protocol to make sure that all active goroutines exit gracefully.
type Stoper struct {
	counter int64
	channel chan int
	stoped  bool
	mutex   sync.Mutex
}

// Stoper factory.
func NewStoper() *Stoper {
	return &Stoper{
		counter: 0,
		channel: make(chan int),
		stoped:  false,
	}
}

// Join starts goroutine participation in shutdown protocol.
func (this *Stoper) Join() {
	atomic.AddInt64(&this.counter, 1)
}

// Leave notifies that participating goroutine gracesfully exited.
// Should be called with defer symantics.
func (this *Stoper) Leave() {
	atomic.AddInt64(&this.counter, -1)
}

// GetChan returns channel to be used in select {} in order to react to Stop event.
func (this *Stoper) GetChan() chan int {
	return this.channel
}

// Stop notifies all participating goroutines that shutdown protocol is in progress
// and waits for all go routines to exit until timeouti.
// Returns false when timed out.
func (this *Stoper) Stop(timeout time.Duration) bool {
	this.stop()
	return this.Wait(timeout)
}

// stop helper
func (this *Stoper) stop() {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	if !this.stoped {
		this.stoped = true
		close(this.channel)
	}
}

func (this *Stoper) Wait(timeout time.Duration) bool {
	now := time.Now()
	for atomic.LoadInt64(&this.counter) > 0 {
		time.Sleep(time.Millisecond * 10)
		if time.Since(now) > timeout {
			return false
		}
	}
	return true
}

func (this *Stoper) Stoped() bool {
	return this.stoped
}

// Counter returns number of of participating goroutines.
func (this *Stoper) Counter() int64 {
	return atomic.LoadInt64(&this.counter)
}
