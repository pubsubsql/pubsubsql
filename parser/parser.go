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

// tokenProducer produces tokens for the parser 
type tokenProducer interface {
	Produce() *token
}

// parser 
type parser struct {
	tokens tokenProducer
}

func (p *parser) parseError(s string) *errorAction {
	e := errorAction{
		err: s,
	}
	return &e
}

// UPDATE
func (p *parser) parseSqlUpdate() action {
	t := p.tokens.Produce()
	if t.typ != tokenTypeSqlTable {
		return p.parseError("expected table name")

	}
	act := &sqlUpdateAction{
		table:   t.val,
		colVals: make([]*columnValue, 0, 10),
	}
	t = p.tokens.Produce()
	if t.typ == tokenTypeSqlSet {
		return p.parseSqlUpdateColVals(act)
	}
	return p.parseError("expected set keyword")
}

func (p *parser) parseSqlUpdateColVals(act *sqlUpdateAction) action {
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
			if erract := p.parseSqlEqualVal(&(act.filter), nil); erract != nil {
				return erract
			}
			// we must be at the end
			break loop 

		case tokenTypeEOF:
			break loop

		case tokenTypeSqlComma:
			continue

		default:
			return p.parseError("expected column or where keyword")
		}
	}
	if count == 0 {
		return p.parseError("expected at least on column value pair")
	}
	return act
}

func (p *parser) parseSqlEqualVal(colval *columnValue, t *token) action {
	// column
	if t == nil {
		t = p.tokens.Produce()
	}
	if t.typ != tokenTypeSqlColumn {
		return p.parseError("expected column name")
	}
	colval.column = t.val
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

// INSERT
func (p *parser) parseSqlInsert() action {
	return nil
}

// SELECT
func (p *parser) parseSqlSelect() action {
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
	// table name
	t = p.tokens.Produce()
	if t.typ != tokenTypeSqlTable {
		return p.parseError("expected table name")
	}
	act := &sqlSelectAction{
		table:   t.val,
	}
	// possible eof
	t = p.tokens.Produce()
	if t.typ == tokenTypeEOF {
		return act	
	}	
	// than it must be where
	if t.typ != tokenTypeSqlWhere {
		return p.parseError("expected where clause")

	}
	if erract := p.parseSqlEqualVal(&(act.filter), nil); erract != nil {
		return erract
	}
	// we are good
	return act
}

// DELETE
func (p *parser) parseSqlDelete() action {
	return nil
}

// run runs the parser
func (p *parser) run() action {
	t := p.tokens.Produce()
	switch t.typ {
	case tokenTypeSqlUpdate:
		return p.parseSqlUpdate()

	case tokenTypeSqlInsert:
		return p.parseSqlInsert()

	case tokenTypeSqlSelect:
		return p.parseSqlSelect()

	case tokenTypeSqlDelete:
		return p.parseSqlDelete()

	}

	return p.parseError("invalid action")
}

// parse parses tokens and returns an action 
func parse(tokens tokenProducer) action {
	p := &parser{
		tokens: tokens,
	}
	return p.run()
}
