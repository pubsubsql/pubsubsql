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

import "sync"

type QuitChan struct {
	quit		chan int
	isquit		bool
	mutex	    *sync.Mutex 	
}

// factory
func newQuitChan() *QuitChan {
	return &QuitChan{
		quit: make(chan int, 1),
		isquit: false,
	}
}

func (q *QuitChan) Quit(quitCookie int) {
	q.mutex.Lock()
	if !q.isquit 	{
		quit <- quitCookie
		q.isquit = true
		close(q.quit)
	}
	q.mutex.Unlock()
}

func (q *QuitChan) Chan() (chan int) {
	return q.quit 
}

func (q *QuitChan) IsQuit() bool {
	return q.isquit 
}

