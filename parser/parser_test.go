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

type tokensProducerConsumer struct {
	idx    int
	tokens []*token
}

func newTokens() *tokensProducerConsumer {
	return &tokensProducerConsumer{
		idx:    0,
		tokens: make([]*token, 0, 30),
	}
}

func reuseTokens(pc *tokensProducerConsumer) {
	pc.idx = 0
}

func (c *tokensProducerConsumer) Consume(t *token) {
	c.tokens = append(c.tokens, t)
}

func (p *tokensProducerConsumer) Produce() *token {
	if p.idx >= len(p.tokens) {
		return &token{
			typ: tokenTypeEOF,
		}
	}
	t := p.tokens[p.idx]
	p.idx++
	return t
}

func expectedError(t *testing.T, a action) {
	switch a.(type) {
	case *errorAction:

	default:
		t.Errorf("parse error: expected error")
	}

}

func validateSelect(t *testing.T, a action, y *sqlSelectAction) {
	switch a.(type) {
	case *errorAction:
		e := a.(*errorAction)
		t.Errorf("parse error: " + e.err)

	case *sqlSelectAction:
		x := a.(*sqlSelectAction)
		// table name
		if x.table != y.table {
			t.Errorf("parse error: unexpected table name: " + x.table)
		}
		// filter
		if x.filter != y.filter {
			t.Errorf("parse error: filters do not match")
		}

	default:
		t.Errorf("parse error: invalid action type expected sqlSelectAction")
	}

}

func validateDelete(t *testing.T, a action, y *sqlDeleteAction) {
	switch a.(type) {
	case *errorAction:
		e := a.(*errorAction)
		t.Errorf("parse error: " + e.err)

	case *sqlDeleteAction:
		x := a.(*sqlDeleteAction)
		// table name
		if x.table != y.table {
			t.Errorf("parse error: unexpected table name: " + x.table)
		}
		// filter
		if x.filter != y.filter {
			t.Errorf("parse error: filters do not match")
		}

	default:
		t.Errorf("parse error: invalid action type expected sqlSelectAction")
	}
}

func validateUpdate(t *testing.T, a action, y *sqlUpdateAction) {
	switch a.(type) {
	case *errorAction:
		e := a.(*errorAction)
		t.Errorf("parse error: " + e.err)

	case *sqlUpdateAction:
		x := a.(*sqlUpdateAction)
		// table name
		if x.table != y.table {
			t.Errorf("parse error: table names do not match " + x.table)
		}
		// number of columns and values
		if len(x.colVals) != len(y.colVals) {
			t.Errorf("parse error: colVals lens do not match")
			break
		}
		// columns and values
		for i := 0; i < len(x.colVals); i++ {
			if *(y.colVals[i]) != *(x.colVals[i]) {
				t.Errorf("parse error: colVals do not match")
			}
		}
		// filter
		if x.filter != y.filter {
			t.Errorf("parse error: filters do not match")

		}

	default:
	}

}

func TestParseSqlUpdateStatement2(t *testing.T) {
	pc := newTokens()
	lex(" update stocks set bid = 140.45, ask = 142.01", pc)
	x := parse(pc)
	var y sqlUpdateAction
	y.table = "stocks"
	y.addColVal("bid", "140.45")
	y.addColVal("ask", "142.01")
	validateUpdate(t, x, &y)

}

func TestParseSqlUpdateStatement1(t *testing.T) {
	pc := newTokens()
	lex(" update stocks set bid = 140.45, ask = 142.01, sector = 'TECH' where ticker = IBM", pc)
	x := parse(pc)
	var y sqlUpdateAction
	y.table = "stocks"
	y.addColVal("bid", "140.45")
	y.addColVal("ask", "142.01")
	y.addColVal("sector", "TECH")
	y.filter.addFilter("ticker", "IBM")
	validateUpdate(t, x, &y)
}

func TestParseSqlUpdateStatement3(t *testing.T) {
	pc := newTokens()
	lex(" update stocks set bid = ", pc)
	x := parse(pc)
	expectedError(t, x)
	//
	pc = newTokens()
	lex(" update stocks ", pc)
	x = parse(pc)
	expectedError(t, x)
	//
	pc = newTokens()
	lex(" update stocks set ", pc)
	x = parse(pc)
	expectedError(t, x)
}

// SELECT
func TestParseSqlSelectStatement1(t *testing.T) {
	pc := newTokens()
	lex(" select *  from stocks ", pc)
	x := parse(pc)
	var y sqlSelectAction
	y.table = "stocks"
	validateSelect(t, x, &y)
}

func TestParseSqlSelectStatement2(t *testing.T) {
	pc := newTokens()
	lex(" select *  from stocks where  ticker = 'IBM'", pc)
	x := parse(pc)
	var y sqlSelectAction
	y.table = "stocks"
	y.filter.addFilter("ticker", "IBM")
	validateSelect(t, x, &y)
}

func TestParseSqlSelectStatement3(t *testing.T) {
	pc := newTokens()
	lex(" select ", pc)
	x := parse(pc)
	expectedError(t, x)
	//
	pc = newTokens()
	lex(" select *", pc)
	x = parse(pc)
	expectedError(t, x)
	//
	pc = newTokens()
	lex(" select * from ", pc)
	x = parse(pc)
	expectedError(t, x)
	//
	pc = newTokens()
	lex(" select * from stocks where", pc)
	x = parse(pc)
	expectedError(t, x)
	//
	pc = newTokens()
	lex(" select * from stocks where ticker ", pc)
	x = parse(pc)
	expectedError(t, x)
	//
	pc = newTokens()
	lex(" select * from stocks where ticker =", pc)
	x = parse(pc)
	expectedError(t, x)
}

// DELETE 
func TestParseSqlDeleteStatement1(t *testing.T) {
	pc := newTokens()
	lex(" delete  from stocks ", pc)
	x := parse(pc)
	var y sqlDeleteAction
	y.table = "stocks"
	validateDelete(t, x, &y)
}

func TestParseSqlDeleteStatement2(t *testing.T) {
	pc := newTokens()
	lex(" delete  from stocks where  ticker = 'IBM'", pc)
	x := parse(pc)
	var y sqlDeleteAction
	y.table = "stocks"
	y.filter.addFilter("ticker", "IBM")
	validateDelete(t, x, &y)
}

func TestParseSqlDeleteStatement3(t *testing.T) {
	pc := newTokens()
	lex(" delete ", pc)
	x := parse(pc)
	expectedError(t, x)
	//
	pc = newTokens()
	lex(" delete from", pc)
	x = parse(pc)
	expectedError(t, x)
	//
	pc = newTokens()
	lex(" delete from stocks where", pc)
	x = parse(pc)
	expectedError(t, x)
	//
	pc = newTokens()
	lex(" delete from stocks where ticker ", pc)
	x = parse(pc)
	expectedError(t, x)
	//
	pc = newTokens()
	lex(" delete from stocks where ticker =", pc)
	x = parse(pc)
	expectedError(t, x)
}
