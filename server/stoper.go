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

//Stoper implements shutdown protocol to make sure that all avtive goroutines exit gracefully.
type Stoper struct {
	counter int64
	channel chan int
}

//Stoper factory.
func NewStoper() *Stoper {
	stoper := new(Stoper)
	stoper.counter = 0
	stoper.channel = make(chan int)
	return stoper
}

//Enter starts participation in shutdown protocol. 
func (this *Stoper) Enter() {
	atomic.AddInt64(&this.counter, 1)
}

//Leave notifies that participating goroutine gracesfully exited.
//should be called with defer symantics
func (this *Stoper) Leave() {
	atomic.AddInt64(&this.counter, -1)
}

//GetChan returns channel to be used in select {} in order to react to Stop event.
func (this *Stoper) GetChan() chan int {
	return this.channel
}

//Stop notifies all participating goroutines that shutdown protocol is in progress
//and waits for all go routines to exit until timeout
//Returns false when timed out
func (this *Stoper) Stop(timeout time.Duration) bool {
	close(this.channel)
	t := time.Now()
	for atomic.LoadInt64(&this.counter) > 0 {
		time.Sleep(time.Millisecond * 10)
		if time.Since(t) > timeout {
			return false
		}
	}
	return true
}

//Counter returns number of of participating goroutines
func (this *Stoper) Counter() int64 {
	return atomic.LoadInt64(&this.counter)
}
