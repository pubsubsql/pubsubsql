/* Copyright (C) 2013 CompleteD LLC.
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

import "fmt"

// tokenProducer produces tokens for the parser 
type tokenProducer interface {
	Produce() *token
}

// parser 
type parser struct {
	tokens tokenProducer
}

func (p *parser) parseError(s string) *errorRequest {
	e := errorRequest{
		err: s,
	}
	return &e
}

func (p *parser) parseSqlEqualVal(colval *columnValue, t *token) request {
	//col
	if t == nil {
		t = p.tokens.Produce()
	}
	if t.typ != tokenTypeSqlColumn {
		return p.parseError("expected.col name")
	}
	colval.col = t.val
	// =
	t = p.tokens.Produce()
	if t.typ != tokenTypeSqlEqual {
		return p.parseError("expected = sign")
	}
	// value
	t = p.tokens.Produce()
	if t.typ != tokenTypeSqlValue {
		return p.parseError("expected valid value")
	}
	colval.val = t.val
	return nil
}

func (p *parser) parseTableName(table *string) request {
	t := p.tokens.Produce()
	if t.typ != tokenTypeSqlTable {
		return p.parseError("expected table name")
	}
	*table = t.val
	return nil
}

func (p *parser) parseColumnName(column *string) request {
	t := p.tokens.Produce()
	if t.typ != tokenTypeSqlColumn {
		return p.parseError("expected column name")
	}
	*column = t.val
	return nil
}

func (p *parser) parseEOF(act request) request {
	t := p.tokens.Produce()
	if t.typ == tokenTypeEOF {
		return act
	}
	return p.parseError("expected EOF")
}

func (p *parser) parseSqlWhere(filter *sqlFilter, t *token) request {
	//must be where
	if t != nil && t.typ != tokenTypeSqlWhere {
		return p.parseError("expected where clause")
	}
	return p.parseSqlEqualVal(&(filter.columnValue), nil)
}

// INSERT
func (p *parser) parseSqlInsert() request {
	// into
	t := p.tokens.Produce()
	if t.typ != tokenTypeSqlInto {
		return p.parseError("expected into")
	}
	act := &sqlInsertRequest{
		colVals: make([]*columnValue, 0, 10),
	}
	// table name
	if erract := p.parseTableName(&act.table); erract != nil {
		return erract
	}
	// (
	t = p.tokens.Produce()
	if t.typ != tokenTypeSqlLeftParenthesis {
		return p.parseError("expected ( ")
	}
	// columns
	columns := 0
	expectedType := tokenTypeSqlColumn
	var erract request
	var str string
	for expectedType == tokenTypeSqlColumn {
		erract, expectedType, str = p.parseSqlInsertColumn()
		if erract != nil {
			return erract
		}
		act.addColumn(str)
		columns++
	}
	// values
	t = p.tokens.Produce()
	if t.typ != tokenTypeSqlValues {
		return p.parseError("expected values keyword")
	}
	// (
	t = p.tokens.Produce()
	if t.typ != tokenTypeSqlLeftParenthesis {
		return p.parseError("expected values ( ")
	}
	//
	expectedType = tokenTypeSqlValue
	values := 0
	for expectedType == tokenTypeSqlValue {
		erract, expectedType, str = p.parseSqlInsertValue()
		if erract != nil {
			return erract
		}
		if values < columns {
			act.setValueAt(values, str)
		}
		values++
	}
	if columns != values {
		s := fmt.Sprintf("number of columns:%d and values:%d do not match", columns, values)
		return p.parseError(s)
	}
	// done
	return act
}

func (p *parser) parseSqlInsertColumn() (request, tokenType, string) {
	t := p.tokens.Produce()
	if t.typ != tokenTypeSqlColumn {
		return p.parseError("expected column name"), tokenTypeError, ""
	}
	str := t.val
	t = p.tokens.Produce()
	if t.typ == tokenTypeSqlComma {
		return nil, tokenTypeSqlColumn, str
	}
	if t.typ == tokenTypeSqlRightParenthesis {
		return nil, tokenTypeSqlValues, str
	}
	return p.parseError("expected , or ) "), tokenTypeError, ""
}

func (p *parser) parseSqlInsertValue() (request, tokenType, string) {
	t := p.tokens.Produce()
	if t.typ != tokenTypeSqlValue {
		return p.parseError("expected value"), tokenTypeError, ""
	}
	str := t.val
	t = p.tokens.Produce()
	if t.typ == tokenTypeSqlComma {
		return nil, tokenTypeSqlValue, str
	}
	if t.typ == tokenTypeSqlRightParenthesis {
		return nil, tokenTypeEOF, str
	}
	return p.parseError("expected , or ) "), tokenTypeError, ""
}

// SELECT
func (p *parser) parseSqlSelect() request {
	// *
	t := p.tokens.Produce()
	if t.typ != tokenTypeSqlStar {
		return p.parseError("expected * symbol")
	}
	// from
	t = p.tokens.Produce()
	if t.typ != tokenTypeSqlFrom {
		return p.parseError("expected from")
	}
	act := new(sqlSelectRequest)
	// table name
	if erract := p.parseTableName(&act.table); erract != nil {
		return erract
	}
	// possible eof
	t = p.tokens.Produce()
	if t.typ == tokenTypeEOF {
		return act
	}
	// where
	if erract := p.parseSqlWhere(&(act.filter), t); erract != nil {
		return erract
	}
	// we are good
	return act
}

// UPDATE
func (p *parser) parseSqlUpdate() request {
	act := &sqlUpdateRequest{
		colVals: make([]*columnValue, 0, 10),
	}
	// table name
	if erract := p.parseTableName(&act.table); erract != nil {
		return erract
	}
	// set
	t := p.tokens.Produce()
	if t.typ == tokenTypeSqlSet {
		return p.parseSqlUpdateColVals(act)
	}
	return p.parseError("expected set keyword")
}

func (p *parser) parseSqlUpdateColVals(act *sqlUpdateRequest) request {
	count := 0
loop:
	for t := p.tokens.Produce(); ; t = p.tokens.Produce() {
		switch t.typ {
		case tokenTypeSqlColumn:
			colval := new(columnValue)
			act.colVals = append(act.colVals, colval)
			if erract := p.parseSqlEqualVal(colval, t); erract != nil {
				return erract
			}
			count++

		case tokenTypeSqlWhere:
			if erract := p.parseSqlWhere(&(act.filter), t); erract != nil {
				return erract
			}
			// we must be at the end
			break loop

		case tokenTypeEOF:
			break loop

		case tokenTypeSqlComma:
			continue

		default:
			return p.parseError("expected.col or where keyword")
		}
	}
	if count == 0 {
		return p.parseError("expected at least on.col value pair")
	}
	return act
}

// DELETE
func (p *parser) parseSqlDelete() request {
	// from
	t := p.tokens.Produce()
	if t.typ != tokenTypeSqlFrom {
		return p.parseError("expected from")
	}
	act := new(sqlDeleteRequest)
	// table name
	if erract := p.parseTableName(&act.table); erract != nil {
		return erract
	}
	// possible eof
	t = p.tokens.Produce()
	if t.typ == tokenTypeEOF {
		return act
	}
	// than it must be where
	if erract := p.parseSqlWhere(&(act.filter), t); erract != nil {
		return erract
	}
	// we are good
	return act
}

// SUBSCRIBE 
func (p *parser) parseSqlSubscribe() request {
	// *
	t := p.tokens.Produce()
	if t.typ != tokenTypeSqlStar {
		return p.parseError("expected * symbol")
	}
	// from
	t = p.tokens.Produce()
	if t.typ != tokenTypeSqlFrom {
		return p.parseError("expected from")
	}
	act := new(sqlSubscribeRequest)
	// table name
	if erract := p.parseTableName(&act.table); erract != nil {
		return erract
	}
	// possible eof
	t = p.tokens.Produce()
	if t.typ == tokenTypeEOF {
		return act
	}
	// where
	if erract := p.parseSqlWhere(&(act.filter), t); erract != nil {
		return erract
	}
	// we are good
	return act
}

// UNSUBSCRIBE
func (p *parser) parseSqlUnsubscribe() request {
	// from
	t := p.tokens.Produce()
	if t.typ != tokenTypeSqlFrom {
		return p.parseError("expected from")
	}
	act := new(sqlUnsubscribeRequest)
	// table name
	if erract := p.parseTableName(&act.table); erract != nil {
		return erract
	}
	return p.parseEOF(act)
}

// KEY
func (p *parser) parseSqlKey() request {
	act := new(sqlKeyRequest)
	// table name
	if erract := p.parseTableName(&act.table); erract != nil {
		return erract
	}
	// column name
	if erract := p.parseColumnName(&act.column); erract != nil {
		return erract
	}
	return p.parseEOF(act)
}

// TAG
func (p *parser) parseSqlTag() request {
	act := new(sqlTagRequest)
	// table name
	if erract := p.parseTableName(&act.table); erract != nil {
		return erract
	}
	// column name
	if erract := p.parseColumnName(&act.column); erract != nil {
		return erract
	}
	return p.parseEOF(act)
}

// run runs the parser
func (p *parser) run() request {
	t := p.tokens.Produce()
	switch t.typ {
	case tokenTypeSqlInsert:
		return p.parseSqlInsert()

	case tokenTypeSqlSelect:
		return p.parseSqlSelect()

	case tokenTypeSqlUpdate:
		return p.parseSqlUpdate()

	case tokenTypeSqlDelete:
		return p.parseSqlDelete()

	case tokenTypeSqlSubscribe:
		return p.parseSqlSubscribe()

	case tokenTypeSqlUnsubscribe:
		return p.parseSqlUnsubscribe()

	case tokenTypeSqlKey:
		return p.parseSqlKey()

	case tokenTypeSqlTag:
		return p.parseSqlTag()

	}

	return p.parseError("invalid request")
}

// parse parses tokens and returns an request 
func parse(tokens tokenProducer) request {
	p := &parser{
		tokens: tokens,
	}
	return p.run()
}
