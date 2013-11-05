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
	// valid 
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
	// invalid 
	{
		consumer := chanTokenConsumer{channel: make(chan token)}
		go lex("insert1234", &consumer)
		c := consumer.channel
		// insert
		tk := <-c
		if tk.typ != tokenTypeError {
			t.Errorf("expected error token")
		}
		// eof	
		tk = <-c
		if tk.typ != tokenTypeEOF {
			t.Errorf("expected eof token")
		}
		//
	}
}

// Tests delete command
func TestDeleteCommand(t *testing.T) {
	// valid
	{
		consumer := chanTokenConsumer{channel: make(chan token)}
		go lex("delete from table1 where ticker = 'IBM'", &consumer)
		c := consumer.channel
		// delete
		tk := <-c
		if tk.typ != tokenTypeSqlDelete {
			t.Errorf("expected delete token")
		}
		// eof	
		tk = <-c
		if tk.typ != tokenTypeEOF {
			t.Errorf("expected eof token")
		}
		//
	}
	// invalid
	{
		consumer := chanTokenConsumer{channel: make(chan token)}
		go lex("delete1234", &consumer)
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

// Tests help command
func TestHelpCommand(t *testing.T) {
	// valid
	{
		consumer := chanTokenConsumer{channel: make(chan token)}
		go lex("help", &consumer)
		c := consumer.channel
		// help
		tk := <-c
		if tk.typ != tokenTypeCmdHelp {
			t.Errorf("expected help token")
		}
		// eof	
		tk = <-c
		if tk.typ != tokenTypeEOF {
			t.Errorf("expected eof token")
		}
		//
	}
	// invalid
	{
		consumer := chanTokenConsumer{channel: make(chan token)}
		go lex("help1234", &consumer)
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

// Tests update command
func TestUpdateCommand(t *testing.T) {
	// valid 
	{
		consumer := chanTokenConsumer{channel: make(chan token)}
		go lex("update table1 set bid = 123.12 where ticker = 'IBM'", &consumer)
		c := consumer.channel
		// help
		tk := <-c
		if tk.typ != tokenTypeSqlUpdate {
			t.Errorf("expected update token")
		}
		// eof	
		tk = <-c
		if tk.typ != tokenTypeEOF {
			t.Errorf("expected eof token")
		}
		//
	}
	// invalid
	{
		consumer := chanTokenConsumer{channel: make(chan token)}
		go lex("update123", &consumer)
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

// Tests unsubscribe command
func TestUnsubscribeCommand(t *testing.T) {
	// valid 
	{
		consumer := chanTokenConsumer{channel: make(chan token)}
		go lex("unsubscribe from table1", &consumer)
		c := consumer.channel
		// help
		tk := <-c
		if tk.typ != tokenTypeSqlUnsubscribe {
			t.Errorf("expected unsubscribe token")
		}
		// eof	
		tk = <-c
		if tk.typ != tokenTypeEOF {
			t.Errorf("expected eof token")
		}
		//
	}
	// invalid
	{
		consumer := chanTokenConsumer{channel: make(chan token)}
		go lex("unsubscribe123", &consumer)
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

// Tests subscribe command
func TestSubscribeCommand(t *testing.T) {
	// valid 
	{
		consumer := chanTokenConsumer{channel: make(chan token)}
		go lex("subscribe * from table1", &consumer)
		c := consumer.channel
		// help
		tk := <-c
		if tk.typ != tokenTypeSqlSubscribe {
			t.Errorf("expected subscribe token")
		}
		// eof	
		tk = <-c
		if tk.typ != tokenTypeEOF {
			t.Errorf("expected eof token")
		}
		//
	}
	// invalid
	{
		consumer := chanTokenConsumer{channel: make(chan token)}
		go lex("subscribe123", &consumer)
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

// Tests select command
func TestSelectCommand(t *testing.T) {
	// valid 
	{
		consumer := chanTokenConsumer{channel: make(chan token)}
		go lex("select * from table1", &consumer)
		c := consumer.channel
		// help
		tk := <-c
		if tk.typ != tokenTypeSqlSelect {
			t.Errorf("expected select token")
		}
		// eof	
		tk = <-c
		if tk.typ != tokenTypeEOF {
			t.Errorf("expected eof token")
		}
		//
	}
	// invalid
	{
		consumer := chanTokenConsumer{channel: make(chan token)}
		go lex("select123", &consumer)
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

// Tests start command
func TestStartCommand(t *testing.T) {
	// valid 
	{
		consumer := chanTokenConsumer{channel: make(chan token)}
		go lex("start", &consumer)
		c := consumer.channel
		// help
		tk := <-c
		if tk.typ != tokenTypeCmdStart {
			t.Errorf("expected start token")
		}
		// eof	
		tk = <-c
		if tk.typ != tokenTypeEOF {
			t.Errorf("expected eof token")
		}
		//
	}
	// invalid
	{
		consumer := chanTokenConsumer{channel: make(chan token)}
		go lex(" start12344", &consumer)
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

// Tests stop command
func TestStopCommand(t *testing.T) {
	// valid 
	{
		consumer := chanTokenConsumer{channel: make(chan token)}
		go lex(" 	stop", &consumer)
		c := consumer.channel
		// help
		tk := <-c
		if tk.typ != tokenTypeCmdStop {
			t.Errorf("expected stop token")
		}
		// eof	
		tk = <-c
		if tk.typ != tokenTypeEOF {
			t.Errorf("expected eof token")
		}
		//
	}
	// invalid
	{
		consumer := chanTokenConsumer{channel: make(chan token)}
		go lex("stop1234", &consumer)
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

// Tests status command
func TestStatusCommand(t *testing.T) {
	// valid 
	{
		consumer := chanTokenConsumer{channel: make(chan token)}
		go lex(" 	status   ", &consumer)
		c := consumer.channel
		// help
		tk := <-c
		if tk.typ != tokenTypeCmdStatus {
			t.Errorf("expected stop token")
		}
		// eof	
		tk = <-c
		if tk.typ != tokenTypeEOF {
			t.Errorf("expected eof token")
		}
		//
	}
	// invalid
	{
		consumer := chanTokenConsumer{channel: make(chan token)}
		go lex("stop1234", &consumer)
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
