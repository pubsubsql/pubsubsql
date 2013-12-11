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

func TestQuitChan(t *testing.T) {
	quit := NewQuitChan()

	for i := 0; i < 10; i++ {
		go func() {
			<-quit.Chan()
			if quit.QuitCookie() != 12345 {
				t.Error("Invalid quit cookie")
			}
			if !quit.IsQuit() {
				t.Error("Invalid quit flag")
			}
		}()
	}

	quit.Quit(12345)
	time.Sleep(time.Millisecond * 100)
}
