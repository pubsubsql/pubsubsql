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

// tokenProducer produces tokens for the parser.
type tokenProducer interface {
	Produce() *token
}

// parser
type parser struct {
	tokens tokenProducer
}

// Indicates that error happened during parse phase and returns errorRequest
func (p *parser) parseError(s string) *errorRequest {
	e := errorRequest{
		err: s,
	}
	return &e
}

// Helper functions

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

func (p *parser) parseEOF(req request) request {
	t := p.tokens.Produce()
	if t.typ == tokenTypeEOF {
		return req
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

// INSERT sql statement

// Parses sql insert statement and returns sqlInsertRequest on success.
func (p *parser) parseSqlInsert() request {
	// into
	t := p.tokens.Produce()
	if t.typ != tokenTypeSqlInto {
		return p.parseError("expected into")
	}
	req := &sqlInsertRequest{
		colVals: make([]*columnValue, 0, PARSER_SQL_INSERT_REQUEST_COLUMN_CAPACITY),
	}
	// table name
	if errreq := p.parseTableName(&req.table); errreq != nil {
		return errreq
	}
	// (
	t = p.tokens.Produce()
	if t.typ != tokenTypeSqlLeftParenthesis {
		return p.parseError("expected ( ")
	}
	// columns
	columns := 0
	expectedType := tokenTypeSqlColumn
	var errreq request
	var str string
	for expectedType == tokenTypeSqlColumn {
		errreq, expectedType, str = p.parseSqlInsertColumn()
		if errreq != nil {
			return errreq
		}
		req.addColumn(str)
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
		errreq, expectedType, str = p.parseSqlInsertValue()
		if errreq != nil {
			return errreq
		}
		if values < columns {
			req.setValueAt(values, str)
		}
		values++
	}
	if columns != values {
		s := fmt.Sprintf("number of columns:%d and values:%d do not match", columns, values)
		return p.parseError(s)
	}
	// done
	return req
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

// SELECT sql statement

// Parses sql select statement and returns sqlSelectRequest on success.
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
	req := new(sqlSelectRequest)
	// table name
	if errreq := p.parseTableName(&req.table); errreq != nil {
		return errreq
	}
	// possible eof
	t = p.tokens.Produce()
	if t.typ == tokenTypeEOF {
		return req
	}
	// where
	if errreq := p.parseSqlWhere(&(req.filter), t); errreq != nil {
		return errreq
	}
	// we are good
	return req
}

// UPDATE sql statement

// Parses sql update statement and returns sqlUpdateRequest on success.
func (p *parser) parseSqlUpdate() request {
	req := &sqlUpdateRequest{
		colVals: make([]*columnValue, 0, PARSER_SQL_UPDATE_REQUEST_COLUMN_CAPACITY),
	}
	// table name
	if errreq := p.parseTableName(&req.table); errreq != nil {
		return errreq
	}
	// set
	t := p.tokens.Produce()
	if t.typ == tokenTypeSqlSet {
		return p.parseSqlUpdateColVals(req)
	}
	return p.parseError("expected set keyword")
}

func (p *parser) parseSqlUpdateColVals(req *sqlUpdateRequest) request {
	count := 0
loop:
	for t := p.tokens.Produce(); ; t = p.tokens.Produce() {
		switch t.typ {
		case tokenTypeSqlColumn:
			colval := new(columnValue)
			req.colVals = append(req.colVals, colval)
			if errreq := p.parseSqlEqualVal(colval, t); errreq != nil {
				return errreq
			}
			count++

		case tokenTypeSqlWhere:
			if errreq := p.parseSqlWhere(&(req.filter), t); errreq != nil {
				return errreq
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
	return req
}

// DELETE sql statement

// Parses sql delete statement and returns sqlDeleteRequest on success.
func (p *parser) parseSqlDelete() request {
	// from
	t := p.tokens.Produce()
	if t.typ != tokenTypeSqlFrom {
		return p.parseError("expected from")
	}
	req := new(sqlDeleteRequest)
	// table name
	if errreq := p.parseTableName(&req.table); errreq != nil {
		return errreq
	}
	// possible eof
	t = p.tokens.Produce()
	if t.typ == tokenTypeEOF {
		return req
	}
	// than it must be where
	if errreq := p.parseSqlWhere(&(req.filter), t); errreq != nil {
		return errreq
	}
	// we are good
	return req
}

// KEY sql statement

// Parses sql key statement and returns sqlKeyRequest on success.
func (p *parser) parseSqlKey() request {
	req := new(sqlKeyRequest)
	// table name
	if errreq := p.parseTableName(&req.table); errreq != nil {
		return errreq
	}
	// column name
	if errreq := p.parseColumnName(&req.column); errreq != nil {
		return errreq
	}
	return p.parseEOF(req)
}

// TAG sql statement

// Parses sql tag statement and returns sqlRequest on success.
func (p *parser) parseSqlTag() request {
	req := new(sqlTagRequest)
	// table name
	if errreq := p.parseTableName(&req.table); errreq != nil {
		return errreq
	}
	// column name
	if errreq := p.parseColumnName(&req.column); errreq != nil {
		return errreq
	}
	return p.parseEOF(req)
}

// SUBSCRIBE sql statement

// Parses sql subscribe statement and returns sqlSubscribeRequest on success.
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
	req := new(sqlSubscribeRequest)
	// table name
	if errreq := p.parseTableName(&req.table); errreq != nil {
		return errreq
	}
	// possible eof
	t = p.tokens.Produce()
	if t.typ == tokenTypeEOF {
		return req
	}
	// where
	if errreq := p.parseSqlWhere(&(req.filter), t); errreq != nil {
		return errreq
	}
	// we are good
	return req
}

// UNSUBSCRIBE sql statement

// Parses sql unsubscribe statement and returns sqlUnsubscribeRequest on success.
func (p *parser) parseSqlUnsubscribe() request {
	// from
	t := p.tokens.Produce()
	if t.typ != tokenTypeSqlFrom {
		return p.parseError("expected from")
	}
	req := new(sqlUnsubscribeRequest)
	// table name
	if errreq := p.parseTableName(&req.table); errreq != nil {
		return errreq
	}
	// possible eof
	t = p.tokens.Produce()
	if t.typ == tokenTypeEOF {
		return req
	}
	// than it must be where
	if errreq := p.parseSqlWhere(&(req.filter), t); errreq != nil {
		return errreq
	}
	// we are good
	return req
}

// Runs the parser.
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

// Parses tokens and returns an request.
func parse(tokens tokenProducer) request {
	p := &parser{
		tokens: tokens,
	}
	return p.run()
}
