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

func validateTokens(t *testing.T, expected []token, tokens chan token) {
	breakLoop := false
	for _, e := range expected {
		g := <-tokens
		if e.typ != g.typ {
			t.Errorf("expected type " + e.typ.String() + " but got " + g.typ.String() + " value: " + g.val)
			breakLoop = true
		}
		if e.val != g.val {
			t.Errorf("expected value " + e.val + " but got " + g.val)
			breakLoop = true
		}
		if breakLoop {
			break
		}
	}
}

// INSERT 
func TestSqlInsertStatement(t *testing.T) {
	consumer := chanTokenConsumer{channel: make(chan token)}
	go lex("insert into stocks (	ticker,bid, ask		 ) values (IBM, '34.43', 465.123)", &consumer)
	expected := []token{
		{tokenTypeSqlInsert, "insert"},
		{tokenTypeSqlInto, "into"},
		{tokenTypeSqlTable, "stocks"},
		{tokenTypeSqlLeftParenthesis, "("},
		{tokenTypeSqlColumn, "ticker"},
		{tokenTypeSqlComma, ","},
		{tokenTypeSqlColumn, "bid"},
		{tokenTypeSqlComma, ","},
		{tokenTypeSqlColumn, "ask"},
		{tokenTypeSqlRightParenthesis, ")"},
		{tokenTypeSqlValues, "values"},
		{tokenTypeSqlLeftParenthesis, "("},
		{tokenTypeSqlValue, "IBM"},
		{tokenTypeSqlComma, ","},
		{tokenTypeSqlValue, "34.43"},
		{tokenTypeSqlComma, ","},
		{tokenTypeSqlValue, "465.123"},
		{tokenTypeSqlRightParenthesis, ")"},
		{tokenTypeEOF, ""}}

	validateTokens(t, expected, consumer.channel)
}

// DELETE 
func TestSqlDeleteStatement1(t *testing.T) {
	consumer := chanTokenConsumer{channel: make(chan token)}
	go lex(" delete	 	 from stocks", &consumer)
	expected := []token{
		{tokenTypeSqlDelete, "delete"},
		{tokenTypeSqlFrom, "from"},
		{tokenTypeSqlTable, "stocks"},
		{tokenTypeEOF, ""}}

	validateTokens(t, expected, consumer.channel)
}

func TestSqlDeleteStatement2(t *testing.T) {
	consumer := chanTokenConsumer{channel: make(chan token)}
	go lex(" delete	 	 from stocks where ticker = 'IBM'  ", &consumer)
	expected := []token{
		{tokenTypeSqlDelete, "delete"},
		{tokenTypeSqlFrom, "from"},
		{tokenTypeSqlTable, "stocks"},
		{tokenTypeSqlWhere, "where"},
		{tokenTypeSqlColumn, "ticker"},
		{tokenTypeSqlEqual, "="},
		{tokenTypeSqlValue, "IBM"},
		{tokenTypeEOF, ""}}

	validateTokens(t, expected, consumer.channel)
}

// SELECT
func TestSqlSelectStatement1(t *testing.T) {
	consumer := chanTokenConsumer{channel: make(chan token)}
	go lex(" select * 	from stocks", &consumer)
	expected := []token{
		{tokenTypeSqlSelect, "select"},
		{tokenTypeSqlStar, "*"},
		{tokenTypeSqlFrom, "from"},
		{tokenTypeSqlTable, "stocks"},
		{tokenTypeEOF, ""}}

	validateTokens(t, expected, consumer.channel)
}

func TestSqlSelectStatement2(t *testing.T) {
	consumer := chanTokenConsumer{channel: make(chan token)}
	go lex(" select	* 	 from stocks where ticker = IBM", &consumer)
	expected := []token{
		{tokenTypeSqlSelect, "select"},
		{tokenTypeSqlStar, "*"},
		{tokenTypeSqlFrom, "from"},
		{tokenTypeSqlTable, "stocks"},
		{tokenTypeSqlWhere, "where"},
		{tokenTypeSqlColumn, "ticker"},
		{tokenTypeSqlEqual, "="},
		{tokenTypeSqlValue, "IBM"},
		{tokenTypeEOF, ""}}

	validateTokens(t, expected, consumer.channel)
}

// UPDATE
func TestSqlUpdateStatement1(t *testing.T) {
	consumer := chanTokenConsumer{channel: make(chan token)}
	go lex(" update stocks set bid = 140.45, ask = 142.01 ", &consumer)
	expected := []token{
		{tokenTypeSqlUpdate, "update"},
		{tokenTypeSqlTable, "stocks"},
		{tokenTypeSqlSet, "set"},
		{tokenTypeSqlColumn, "bid"},
		{tokenTypeSqlEqual, "="},
		{tokenTypeSqlValue, "140.45"},
		{tokenTypeSqlComma, ","},
		{tokenTypeSqlColumn, "ask"},
		{tokenTypeSqlEqual, "="},
		{tokenTypeSqlValue, "142.01"},
		{tokenTypeEOF, ""}}

	validateTokens(t, expected, consumer.channel)
}

func TestSqlUpdateStatement2(t *testing.T) {
	consumer := chanTokenConsumer{channel: make(chan token)}
	go lex(" update stocks set bid = 140.45, ask = '142.01' where ticker = 'GOOG'", &consumer)
	expected := []token{
		{tokenTypeSqlUpdate, "update"},
		{tokenTypeSqlTable, "stocks"},
		{tokenTypeSqlSet, "set"},
		{tokenTypeSqlColumn, "bid"},
		{tokenTypeSqlEqual, "="},
		{tokenTypeSqlValue, "140.45"},
		{tokenTypeSqlComma, ","},
		{tokenTypeSqlColumn, "ask"},
		{tokenTypeSqlEqual, "="},
		{tokenTypeSqlValue, "142.01"},
		{tokenTypeSqlWhere, "where"},
		{tokenTypeSqlColumn, "ticker"},
		{tokenTypeSqlEqual, "="},
		{tokenTypeSqlValue, "GOOG"},
		{tokenTypeEOF, ""}}

	validateTokens(t, expected, consumer.channel)
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
