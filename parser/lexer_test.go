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
import "fmt"

// prints consumed tokens on a separate line
type printlnTokenConsumer struct {
}

func (c *printlnTokenConsumer) Consume(t token) {
	fmt.Println(t)
}

// sends consumed tokens to the channel
type chanTokenConsumer struct {
	channel chan token
}

func (consumer *chanTokenConsumer) Consume(t token) {
	consumer.channel <- t
	if t.typ == tokenTypeEOF {
		close(consumer.channel)
	}
}

// Tests insert command
func TestInsertCommand(t *testing.T) {
	// valid insert
	{
		consumer := chanTokenConsumer{channel: make(chan token)}
		go lex("insert into table1 (ticker, bid, ask) values (IBM, 12.45, 34.67)", &consumer)
		c := consumer.channel
		//test error
		tk := <-c
		if tk.typ != tokenTypeSqlInsert {
			t.Errorf("expected insert token")
		}
		//test eof	
		tk = <-c
		if tk.typ != tokenTypeEOF {
			t.Errorf("expected eof token")
		}
		//
	}
	// invalid command
	{
		consumer := chanTokenConsumer{channel: make(chan token)}
		go lex("insert1234", &consumer)
		c := consumer.channel
		//test error
		tk := <-c
		if tk.typ != tokenTypeError {
			t.Errorf("expected error token")
		}
		//test eof	
		tk = <-c
		if tk.typ != tokenTypeEOF {
			t.Errorf("expected eof token")
		}
		//
	}
}


