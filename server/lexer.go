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

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

// tokenType identifies the type of lex tokens.
type tokenType uint8

const (
	tokenTypeError                   tokenType = iota // error occured 
	tokenTypeEOF                                      // last token
	tokenTypeCmdHelp                                  // help
	tokenTypeCmdStatus                                // status
	tokenTypeCmdStop                                  // stop
	tokenTypeCmdStart                                 // start	
	tokenTypeSqlTable                                 // table name
	tokenTypeSqlColumn                                // column name
	tokenTypeSqlInsert                                // insert
	tokenTypeSqlInto                                  // into
	tokenTypeSqlUpdate                                // update	
	tokenTypeSqlSet                                   // set
	tokenTypeSqlDelete                                // delete
	tokenTypeSqlFrom                                  // from
	tokenTypeSqlSelect                                // select
	tokenTypeSqlSubscribe                             // subscribe
	tokenTypeSqlUnsubscribe                           // unsubscribe 
	tokenTypeSqlWhere                                 // where
	tokenTypeSqlValues                                // values
	tokenTypeSqlStar                                  // *
	tokenTypeSqlEqual                                 // =
	tokenTypeSqlLeftParenthesis                       // (
	tokenTypeSqlRightParenthesis                      // )
	tokenTypeSqlComma                                 // ,
	tokenTypeSqlValue                                 // 'some string' string or continous sequence of chars delimited by WHITE SPACE | ' | , | ( | ) 
	tokenTypeSqlValueWithSingleQuote                  // '' becomes ' inside the string, parser will need to replace the string
	tokenTypeSqlKey                                   // key
	tokenTypeSqlTag                                   // tag
)

// String converts tokenType value to a string. 
func (typ tokenType) String() string {
	switch typ {
	case tokenTypeError:
		return "tokenTypeError"
	case tokenTypeEOF:
		return "tokenTypeEOF"
	case tokenTypeCmdHelp:
		return "tokenTypeCmdHelp"
	case tokenTypeCmdStatus:
		return "tokenTypeCmdStatus"
	case tokenTypeCmdStop:
		return "tokenTypeCmdStop"
	case tokenTypeCmdStart:
		return "tokenTypeCmdStart"
	case tokenTypeSqlTable:
		return "tokenTypeSqlTable"
	case tokenTypeSqlColumn:
		return "tokenTypeSqlColumn"
	case tokenTypeSqlInsert:
		return "tokenTypeSqlInsert"
	case tokenTypeSqlInto:
		return "tokenTypeSqlInto"
	case tokenTypeSqlUpdate:
		return "tokenTypeSqlUpdate"
	case tokenTypeSqlSet:
		return "tokenTypeSqlSet"
	case tokenTypeSqlDelete:
		return "tokenTypeSqlDelete"
	case tokenTypeSqlFrom:
		return "tokenTypeSqlFrom"
	case tokenTypeSqlSelect:
		return "tokenTypeSqlSelect"
	case tokenTypeSqlSubscribe:
		return "tokenTypeSqlSubscribe"
	case tokenTypeSqlUnsubscribe:
		return "tokenTypeSqlUnsubscribe"
	case tokenTypeSqlWhere:
		return "tokenTypeSqlWhere"
	case tokenTypeSqlValues:
		return "tokenTypeSqlValues"
	case tokenTypeSqlStar:
		return "tokenTypeSqlStar"
	case tokenTypeSqlEqual:
		return "tokenTypeSqlEqual"
	case tokenTypeSqlLeftParenthesis:
		return "tokenTypeSqlLeftParenthesis"
	case tokenTypeSqlRightParenthesis:
		return "tokenTypeSqlRightParenthesis"
	case tokenTypeSqlComma:
		return "tokenTypeSqlComma"
	case tokenTypeSqlValue:
		return "tokenTypeSqlValue"
	case tokenTypeSqlValueWithSingleQuote:
		return "tokenTypeSqlValueWithSingleQuote"
	case tokenTypeSqlKey:
		return "tokenTypeSqlKey"
	case tokenTypeSqlTag:
		return "tokenTypeSqlTag"
	}
	return "not implemented"
}

// token is a symbol representing lexical unit. 
type token struct {
	typ tokenType
	// string identified by lexer as a token based on
	// the pattern rule for the tokenType
	val string
}

// String converts token to a string. 
func (t token) String() string {
	if t.typ == tokenTypeEOF {
		return "EOF"
	}
	return t.val
}

// tokenConsumer consumes tokens emited by lexer.
type tokenConsumer interface {
	Consume(t *token)
}

// lexer holds the state of the scanner.
type lexer struct {
	input  string        // the string being scanned
	start  int           // start position of this item
	pos    int           // currenty position in the input
	width  int           // width of last rune read from input
	tokens tokenConsumer // consumed tokens
	err    string        // error message
}

// stateFn represents the state of the lexer
// as a function that returns the next state.
type stateFn func(*lexer) stateFn

// Emits an error token and terminates the scan
// by passing back a nil ponter that will be the next state
// terminating lexer.run function
func (l *lexer) errorToken(format string, args ...interface{}) stateFn {
	l.err = fmt.Sprintf(format, args...)
	l.tokens.Consume(&token{tokenTypeError, l.err})
	return nil
}

// Returns true if scan was a success. 
func (l *lexer) ok() bool {
	return len(l.err) > 0
}

// Passes a token to the token consumer. 
func (l *lexer) emit(t tokenType) {
	l.tokens.Consume(&token{t, l.current()})
}

// Returns current lexeme string.
func (l *lexer) current() string {
	str := l.input[l.start:l.pos]
	l.start = l.pos
	return str
}

// Returns the next rune in the input.
func (l *lexer) next() (rune int32) {
	if l.pos >= len(l.input) {
		l.width = 0
		return 0
	}
	rune, l.width = utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += l.width
	return rune
}

// Returns whether end was reached in the input.
func (l *lexer) end() bool {
	if l.pos >= len(l.input) {
		return true
	}
	return false
}

// Skips over the pending input before this point.
func (l *lexer) ignore() {
	l.start = l.pos
}

// Steps back one rune.
func (l *lexer) backup() {
	l.pos -= l.width
}

// Returns but does not consume the next rune in the input.
func (l *lexer) peek() int32 {
	rune := l.next()
	l.backup()
	return rune
}

// Determines if rune is valid unicode space character or 0.
func isWhiteSpace(rune int32) bool {
	return (unicode.IsSpace(rune) || rune == 0)
}

// Reads till first white space character
// as defined by isWhiteSpace function
func (l *lexer) scanTillWhiteSpace() {
	for rune := l.next(); !isWhiteSpace(rune); rune = l.next() {

	}
}

// Skips white space characters in the input.
func (l *lexer) skipWhiteSpaces() {
	for rune := l.next(); unicode.IsSpace(rune); rune = l.next() {
	}
	l.backup()
	l.ignore()
}

// Scans input and matches against the string.
// Returns true if the expected string was matched.
func (l *lexer) match(str string, skip int) bool {
	done := true
	for _, rune := range str {
		if skip > 0 {
			skip--
			continue
		}
		if rune != l.next() {
			done = false
		}
	}
	if !isWhiteSpace(l.peek()) {
		done = false
		l.scanTillWhiteSpace()
	}
	return done
}

// Scans input and tries to match the expected string.
// Returns true if the expected string was matched.
// Does not advance the input if the string was not matched.
func (l *lexer) tryMatch(val string) bool {
	i := 0
	for _, rune := range val {
		i++
		if rune != l.next() {
			for ; i > 0; i-- {
				l.backup()
			}
			return false
		}
	}
	return true
}

// lexMatch matches expected string value emiting the token on success
// and returning passed state function.
func (l *lexer) lexMatch(typ tokenType, value string, skip int, fn stateFn) stateFn {
	if l.match(value, skip) {
		l.emit(typ)
		return fn
	}
	return l.errorToken("Unexpected token:" + l.current())
}

// lexSqlIndentifier scans input for valid sql identifier emiting the token on success
// and returning passed state function.
func (l *lexer) lexSqlIdentifier(typ tokenType, fn stateFn) stateFn {
	l.skipWhiteSpaces()
	// first rune has to be valid unicode letter	
	if !unicode.IsLetter(l.next()) {
		return l.errorToken("identifier must begin with a letter " + l.current())
	}
	for rune := l.next(); unicode.IsLetter(rune) || unicode.IsDigit(rune); rune = l.next() {

	}
	l.backup()
	l.emit(typ)
	return fn
}

// lexSqlLeftParenthesis scans input for '(' emiting the token on success 
// and returning passed state function.
func (l *lexer) lexSqlLeftParenthesis(fn stateFn) stateFn {
	l.skipWhiteSpaces()
	if l.next() != '(' {
		return l.errorToken("expected ( ")
	}
	l.emit(tokenTypeSqlLeftParenthesis)
	return fn
}

// lexSqlValue scans input for valid sql value emiting the token on success
// and returing passed state function.
func (l *lexer) lexSqlValue(fn stateFn) stateFn {
	l.skipWhiteSpaces()
	if l.end() {
		return l.errorToken("expected value but go eof")
	}
	rune := l.next()
	typ := tokenTypeSqlValue
	// real string
	if rune == '\'' {
		l.ignore()
		for rune = l.next(); ; rune = l.next() {
			if rune == '\'' {
				if !l.end() {
					rune = l.next()
					// check for '''
					if rune == '\'' {
						typ = tokenTypeSqlValueWithSingleQuote
					} else {
						// since we read lookahead after single quote that ends the string 
						// for lookahead
						l.backup()
						// for single quote which is not part of the value
						l.backup()
						l.emit(typ)
						// now ignore that single quote 
						l.next()
						l.ignore()
						//
						return fn
					}
				} else {
					// at the very end
					l.backup()
					l.emit(typ)
					l.next()
					return fn
				}
			}
			if rune == 0 {
				return l.errorToken("string was not delimited")
			}
		}
		// value 
	} else {
		for rune = l.next(); !isWhiteSpace(rune) && rune != ',' && rune != ')'; rune = l.next() {
		}
		l.backup()
		l.emit(typ)
		return fn
	}
	return nil
}

// Tries to match expected value returns next state function depending on the match.
func (l *lexer) lexTryMatch(typ tokenType, val string, fnMatch stateFn, fnNoMatch stateFn) stateFn {
	l.skipWhiteSpaces()
	if l.tryMatch(val) {
		l.emit(typ)
		return fnMatch
	}
	return fnNoMatch
}

// WHERE sql where clause scan state functions.

func lexSqlWhereColumn(l *lexer) stateFn {
	return l.lexSqlIdentifier(tokenTypeSqlColumn, lexSqlWhereColumnEqual)
}

func lexSqlWhereColumnEqual(l *lexer) stateFn {
	l.skipWhiteSpaces()
	if l.next() == '=' {
		l.emit(tokenTypeSqlEqual)
		return lexSqlWhereColumnEqualValue
	}
	return l.errorToken("expected = ")
}

func lexSqlWhereColumnEqualValue(l *lexer) stateFn {
	l.skipWhiteSpaces()
	return l.lexSqlValue(lexEof)
}

func lexEof(l *lexer) stateFn {
	l.skipWhiteSpaces()
	if l.end() {
		return nil
	}
	return l.errorToken("unexpected token at the end of statement")
}

// INSERT sql statement scan state functions.

func lexSqlInsertInto(l *lexer) stateFn {
	l.skipWhiteSpaces()
	return l.lexMatch(tokenTypeSqlInto, "into", 0, lexSqlInsertIntoTable)
}

func lexSqlInsertIntoTable(l *lexer) stateFn {
	return l.lexSqlIdentifier(tokenTypeSqlTable, lexSqlInsertIntoTableLeftParenthesis)
}

func lexSqlInsertIntoTableLeftParenthesis(l *lexer) stateFn {
	return l.lexSqlLeftParenthesis(lexSqlInsertColumn)
}

func lexSqlInsertColumn(l *lexer) stateFn {
	l.skipWhiteSpaces()
	return l.lexSqlIdentifier(tokenTypeSqlColumn, lexSqlInsertColumnCommaOrRightParenthesis)
}

func lexSqlInsertColumnCommaOrRightParenthesis(l *lexer) stateFn {
	l.skipWhiteSpaces()
	switch l.next() {
	case ',':
		l.emit(tokenTypeSqlComma)
		return lexSqlInsertColumn
	case ')':
		l.emit(tokenTypeSqlRightParenthesis)
		return lexSqlInsertValues
	}
	return l.errorToken("expected , or ) ")
}

func lexSqlInsertValues(l *lexer) stateFn {
	l.skipWhiteSpaces()
	return l.lexMatch(tokenTypeSqlValues, "values", 0, lexSqlInsertValuesLeftParenthesis)
}

func lexSqlInsertValuesLeftParenthesis(l *lexer) stateFn {
	return l.lexSqlLeftParenthesis(lexSqlInsertVal)
}

func lexSqlInsertVal(l *lexer) stateFn {
	return l.lexSqlValue(lexSqlInsertValueCommaOrRigthParenthesis)
}

func lexSqlInsertValueCommaOrRigthParenthesis(l *lexer) stateFn {
	l.skipWhiteSpaces()
	switch l.next() {
	case ',':
		l.emit(tokenTypeSqlComma)
		return lexSqlInsertVal
	case ')':
		l.emit(tokenTypeSqlRightParenthesis)
		// we are done with insert
		return nil
	}
	return l.errorToken("expected , or ) ")
}

// SELECT sql statement scan state functions.

func lexSqlSelectStar(l *lexer) stateFn {
	l.skipWhiteSpaces()
	if l.next() == '*' {
		l.emit(tokenTypeSqlStar)
		return lexSqlFrom
	}
	return l.errorToken("expected columns or *")
}

// UPDATE sql statement scan state functions.

func lexSqlUpdateTable(l *lexer) stateFn {
	l.skipWhiteSpaces()
	return l.lexSqlIdentifier(tokenTypeSqlTable, lexSqlUpdateTableSet)
}

func lexSqlUpdateTableSet(l *lexer) stateFn {
	l.skipWhiteSpaces()
	return l.lexMatch(tokenTypeSqlSet, "set", 0, lexSqlColumn)
}

func lexSqlColumn(l *lexer) stateFn {
	l.skipWhiteSpaces()
	if l.end() {
		return nil
	}
	return l.lexSqlIdentifier(tokenTypeSqlColumn, lexSqlColumnEqual)
}

func lexSqlColumnEqual(l *lexer) stateFn {
	l.skipWhiteSpaces()
	if l.next() == '=' {
		l.emit(tokenTypeSqlEqual)
		return lexSqlColumnEqualValue
	}
	return l.errorToken("expecgted = ")
}

func lexSqlColumnEqualValue(l *lexer) stateFn {
	l.skipWhiteSpaces()
	return l.lexSqlValue(lexSqlCommaOrWhere)
}

func lexSqlCommaOrWhere(l *lexer) stateFn {
	l.skipWhiteSpaces()
	if l.next() == ',' {
		l.emit(tokenTypeSqlComma)
		return lexSqlColumn
	}
	l.backup()
	return lexSqlWhere
}

// DELETE sql statement scan state functions.

func lexSqlFrom(l *lexer) stateFn {
	l.skipWhiteSpaces()
	return l.lexMatch(tokenTypeSqlFrom, "from", 0, lexSqlFromTable)
}

func lexSqlFromTable(l *lexer) stateFn {
	return l.lexSqlIdentifier(tokenTypeSqlTable, lexSqlWhere)
}

func lexSqlWhere(l *lexer) stateFn {
	return l.lexTryMatch(tokenTypeSqlWhere, "where", lexSqlWhereColumn, nil)
}

// KEY and TAG sql statement scan state functions. 

func lexSqlKeyTable(l *lexer) stateFn {
	return l.lexSqlIdentifier(tokenTypeSqlTable, lexSqlKeyColumn)
}

func lexSqlKeyColumn(l *lexer) stateFn {
	return l.lexSqlIdentifier(tokenTypeSqlColumn, nil)
}

// UNSUBSCRIBE

func lexSqlUnsubscribeFrom(l *lexer) stateFn {
	l.skipWhiteSpaces()
	return l.lexMatch(tokenTypeSqlFrom, "from", 0, lexSqlUnsubscribeFromTable)
}

func lexSqlUnsubscribeFromTable(l *lexer) stateFn {
	return l.lexSqlIdentifier(tokenTypeSqlTable, nil)
}

// END SQL

// Helper function to process status stop start commands.
func lexCommandST(l *lexer) stateFn {
	switch l.next() {
	case 'a':
		if l.next() == 'r' {
			return l.lexMatch(tokenTypeCmdStart, "start", 4, nil)
		}
		return l.lexMatch(tokenTypeCmdStatus, "status", 4, nil)
	default:
		return l.lexMatch(tokenTypeCmdStop, "stop", 3, nil)
	}
	return l.errorToken("Invalid command:" + l.current())
}

// Helper function to process select subscribe status stop start commands.
func lexCommandS(l *lexer) stateFn {
	switch l.next() {
	case 'e':
		return l.lexMatch(tokenTypeSqlSelect, "select", 2, lexSqlSelectStar)
	case 'u':
		return l.lexMatch(tokenTypeSqlSubscribe, "subscribe", 2, lexSqlSelectStar)
	case 't':
		return lexCommandST(l)
	}
	return l.errorToken("Invalid command:" + l.current())
}

// Initial state function.
func lexCommand(l *lexer) stateFn {
	l.skipWhiteSpaces()
	switch l.next() {
	case 'u': // update unsubscribe
		if l.next() == 'p' {
			return l.lexMatch(tokenTypeSqlUpdate, "update", 2, lexSqlUpdateTable)
		}
		return l.lexMatch(tokenTypeSqlUnsubscribe, "unsubscribe", 2, lexSqlUnsubscribeFrom)
	case 's': // select subscribe status stop start
		return lexCommandS(l)
	case 'i': // insert
		return l.lexMatch(tokenTypeSqlInsert, "insert", 1, lexSqlInsertInto)
	case 'd': // delete
		return l.lexMatch(tokenTypeSqlDelete, "delete", 1, lexSqlFrom)
	case 'h': // help
		return l.lexMatch(tokenTypeCmdHelp, "help", 1, nil)
	case 'k': // key
		return l.lexMatch(tokenTypeSqlKey, "key", 1, lexSqlKeyTable)
	case 't': // tag
		return l.lexMatch(tokenTypeSqlTag, "tag", 1, lexSqlKeyTable)
	}
	return l.errorToken("Invalid command:" + l.current())
}

// Scans the input by executing state functon until. 
// the state is nil
func (l *lexer) run() {
	for state := lexCommand; state != nil; {
		state = state(l)
	}
	l.emit(tokenTypeEOF)
}

// Scans the input by running lexer.
func lex(input string, tokens tokenConsumer) bool {
	l := &lexer{
		input:  input,
		tokens: tokens,
	}
	l.run()
	return l.ok()
}
