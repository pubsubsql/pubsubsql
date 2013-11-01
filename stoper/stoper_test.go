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
import "runtime"

func Test(t *testing.T) {
	stoper := NewStoper()
	if !stoper.Stop(0) {
		t.Errorf("stoper.Stop() expected true but got false");
	}
}


func testStoper (stoper *Stoper, level int, perLevel int) {
	stoper.Enter()
	defer stoper.Leave()
	level--
	if (level < 0) {
		return
	}
	//start other go routines
	for i := 0; i < perLevel; i++ {
		go testStoper(stoper, level, perLevel)
	}
	//wait for stop event
	c := stoper.GetChan()
	<-c;
}

func TestMultiGoroutines(t *testing.T) {
	stoper := NewStoper()
	levels := 5
	perLevel := 10
	go testStoper(stoper, levels, perLevel)
	time.Sleep(time.Millisecond * 500)
	if !stoper.Stop(time.Millisecond * 1000) {
		t.Errorf("stoper.Stop() expected true but got false");
	}
}

func TestMultiGoroutinesMultiCores(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU() - 1)	
	stoper := NewStoper()
	levels := 5
	perLevel := 10
	go testStoper(stoper, levels, perLevel)
	time.Sleep(time.Millisecond * 500)
	if !stoper.Stop(time.Millisecond * 1000) {
		t.Errorf("stoper.Stop() expected true but got false");
	}
}

