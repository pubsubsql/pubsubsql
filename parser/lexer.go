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

// tokenType identifies the type of lex tokens
type tokenType uint8

const (
	tokenTypeError               tokenType = iota // error occured 
	tokenTypeEOF                                  // last token
	tokenTypeCmdHelp                              // help
	tokenTypeCmdStatus                            // status
	tokenTypeCmdStop                              // stop
	tokenTypeCmdStart                             // start	
	tokenTypeSqlTable                             // table name
	tokenTypeSqlColumn                            // column name
	tokenTypeSqlInsert                            // insert
	tokenTypeSqlInto                              // into
	tokenTypeSqlUpdate                            // update	
	tokenTypeSqlSet                               // set
	tokenTypeSqlDelete                            // delete
	tokenTypeSqlFrom                              // from
	tokenTypeSqlSelect                            // select
	tokenTypeSqlSubscribe                         // subscribe
	tokenTypeSqlUnsubscribe                       // unsubscribe 
	tokenTypeSqlWhere                             // where
	tokenTypeSqlValues                            // values
	tokenTypeSqlStar                              // *
	tokenTypeSqlEqual                             // =
	tokenTypeSqlLeftParenthesis                   // (
	tokenTypeSqlRightParenthesis                  // )
	tokenTypeSqlComma                             // ,
	tokenTypeSqlValue                             // 'some string' string or continous sequence of chars delimited by WHITE SPACE | ' | , | ( | ) 
	tokenTypeSqlValueWithSingleQuote 			  // '' becomes ' inside the string, parser will need to replace the string
	tokenTypeSqlKey								  // key
	tokenTypeSqlTag								  // tag
)

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

// token is a a symbol representing lexical unit 
type token struct {
	typ tokenType
	// string identified by lexer as a token based on
	// the pattern rule for the tokenType
	val string
}

func (t token) String() string {
	if t.typ == tokenTypeEOF {
		return "EOF"
	}
	return t.val
}

// tokenConsumer consumes tokens emited by lexer
type tokenConsumer interface {
	Consume(t token)
}

// lexer holds the state of the scanner
type lexer struct {
	input  string        // the string being scanned
	start  int           // start position of this item
	pos    int           // currenty position in the input
	width  int           // width of last rune read from input
	tokens tokenConsumer // consumed tokens
}

// stateFn represents the state of the lexer
// as a function that returns the next state.
type stateFn func(*lexer) stateFn

// errorToken emits an error toekan and terminates the scan
// by passing back a nil ponter that will be the next statei,
// terminating l.run
func (l *lexer) errorToken(format string, args ...interface{}) stateFn {
	l.tokens.Consume(token{tokenTypeError, fmt.Sprintf(format, args...)})
	return nil
}

// emit passes a token to the token consumer 
func (l *lexer) emit(t tokenType) {
	l.tokens.Consume(token{t, l.current()})
}

// returns current lexeme string
func (l *lexer) current() string {
	str := l.input[l.start:l.pos]
	l.start = l.pos
	return str
}

// next returns the next rune in the input
func (l *lexer) next() (rune int32) {
	if l.pos >= len(l.input) {
		l.width = 0
		return 0
	}
	rune, l.width = utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += l.width
	return rune
}

// end returns whether end was reached in the input
func (l *lexer) end() bool {
	if l.pos >= len(l.input) {
		return true
	}
	return false
}

// ignore skips over the pending input before this point
func (l *lexer) ignore() {
	l.start = l.pos
}

// backup steps back one rune
func (l *lexer) backup() {
	l.pos -= l.width
}

// peek returns but does not consume the next rune in the input
func (l *lexer) peek() int32 {
	rune := l.next()
	l.backup()
	return rune
}

// isWhiteSpace determines if rune is valid unicode space character or 0
func isWhiteSpace(rune int32) bool {
	return (unicode.IsSpace(rune) || rune == 0)
}

// scanTillWhiteSpace reads till first unicode White space character
func (l *lexer) scanTillWhiteSpace() {
	for rune := l.next(); !isWhiteSpace(rune); rune = l.next() {

	}
}

// match scans input and matches against the string
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
	if false == isWhiteSpace(l.peek()) {
		done = false
		l.scanTillWhiteSpace()
	}
	return done
}

// tryMatch scans input and tries to match the string value
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

// lexMatch matches expected string value
func (l *lexer) lexMatch(typ tokenType, value string, skip int, fn stateFn) stateFn {
	if l.match(value, skip) {
		l.emit(typ)
		return fn
	}
	l.errorToken("Unexpected token:" + l.current())
	return nil
}

// lexSqlIndentifier scans input for valid sql identifier
func (l *lexer) lexSqlIdentifier(typ tokenType, fn stateFn) stateFn {
	l.lexSkipWhiteSpaces()
	// first rune has to be valid unicode letter	
	if !unicode.IsLetter(l.next()) {
		l.errorToken("identifier must begin with a letter " + l.current())
		return nil
	}
	for rune := l.next(); unicode.IsLetter(rune) || unicode.IsDigit(rune); rune = l.next() {

	}
	l.backup()
	l.emit(typ)
	return fn
}

// lexSqlLeftParenthesis scans input for '('
func (l *lexer) lexSqlLeftParenthesis(fn stateFn) stateFn {
	l.lexSkipWhiteSpaces()
	if l.next() != '(' {
		l.errorToken("expected ( ")
		return nil
	}
	l.emit(tokenTypeSqlLeftParenthesis)
	return fn
}

// lexSqlValue scans input for valid sql value
// single quoted string 'some string '' '
// unicode character sequence delimited by white space or ')' or ','
func (l *lexer) lexSqlValue(fn stateFn) stateFn {
	l.lexSkipWhiteSpaces()
	if l.end() {
		l.errorToken("expected value but go eof")
		return nil
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
				l.errorToken("string was not delimited")
				return nil
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

// lexTryMatch tries to match expected value returns next state function depending on the match
func (l *lexer) lexTryMatch(typ tokenType, val string, fnMatch stateFn, fnNoMatch stateFn) stateFn {
	l.lexSkipWhiteSpaces()
	if l.tryMatch(val) {
		l.emit(typ)
		return fnMatch
	}
	return fnNoMatch
}

// INSERT

func lexSqlInsertInto(l *lexer) stateFn {
	l.lexSkipWhiteSpaces()
	return l.lexMatch(tokenTypeSqlInto, "into", 0, lexSqlInsertIntoTable)
}

func lexSqlInsertIntoTable(l *lexer) stateFn {
	return l.lexSqlIdentifier(tokenTypeSqlTable, lexSqlInsertIntoTableLeftParenthesis)
}

func lexSqlInsertIntoTableLeftParenthesis(l *lexer) stateFn {
	return l.lexSqlLeftParenthesis(lexSqlInsertColumn)
}

func lexSqlInsertColumn(l *lexer) stateFn {
	l.lexSkipWhiteSpaces()
	return l.lexSqlIdentifier(tokenTypeSqlColumn, lexSqlInsertColumnCommaOrRightParenthesis)
}

func lexSqlInsertColumnCommaOrRightParenthesis(l *lexer) stateFn {
	l.lexSkipWhiteSpaces()
	switch l.next() {
	case ',':
		l.emit(tokenTypeSqlComma)
		return lexSqlInsertColumn
	case ')':
		l.emit(tokenTypeSqlRightParenthesis)
		return lexSqlInsertValues
	}
	l.errorToken("expected , or ) ")
	return nil
}

func lexSqlInsertValues(l *lexer) stateFn {
	l.lexSkipWhiteSpaces()
	return l.lexMatch(tokenTypeSqlValues, "values", 0, lexSqlInsertValuesLeftParenthesis)
}

func lexSqlInsertValuesLeftParenthesis(l *lexer) stateFn {
	return l.lexSqlLeftParenthesis(lexSqlInsertVal)
}

func lexSqlInsertVal(l *lexer) stateFn {
	return l.lexSqlValue(lexSqlInsertValueCommaOrRigthParenthesis)
}

func lexSqlInsertValueCommaOrRigthParenthesis(l *lexer) stateFn {
	l.lexSkipWhiteSpaces()
	switch l.next() {
	case ',':
		l.emit(tokenTypeSqlComma)
		return lexSqlInsertVal
	case ')':
		l.emit(tokenTypeSqlRightParenthesis)
		// we are done with insert
		return nil
	}
	l.errorToken("expected , or ) ")
	return nil
}

// DELETE

func lexSqlFrom(l *lexer) stateFn {
	l.lexSkipWhiteSpaces()
	return l.lexMatch(tokenTypeSqlFrom, "from", 0, lexSqlFromTable)
}

func lexSqlFromTable(l *lexer) stateFn {
	return l.lexSqlIdentifier(tokenTypeSqlTable, lexSqlWhere)
}

func lexSqlWhere(l *lexer) stateFn {
	return l.lexTryMatch(tokenTypeSqlWhere, "where", lexSqlWhereColumn, nil)
}

// WHERE
// TODO add more where complex expressions

func lexSqlWhereColumn(l *lexer) stateFn {
	return l.lexSqlIdentifier(tokenTypeSqlColumn, lexSqlWhereColumnEqual)
}

func lexSqlWhereColumnEqual(l *lexer) stateFn {
	l.lexSkipWhiteSpaces()
	if l.next() == '=' {
		l.emit(tokenTypeSqlEqual)
		return lexSqlWhereColumnEqualValue
	}
	l.errorToken("expected = ")
	return nil
}

func lexSqlWhereColumnEqualValue(l *lexer) stateFn {
	l.lexSkipWhiteSpaces()
	return l.lexSqlValue(lexEof)
}

func lexEof(l *lexer) stateFn {
	l.lexSkipWhiteSpaces()
	if l.end() {
		return nil
	}
	l.errorToken("unexpected token at the end of statement")
	return nil
}

// SELECT
// TODO add ability to select individual columns

func lexSqlSelectStar(l *lexer) stateFn {
	l.lexSkipWhiteSpaces()
	if l.next() == '*' {
		l.emit(tokenTypeSqlStar)
		return lexSqlFrom
	}
	l.errorToken("expected columns or *")
	return nil
}

// UPDATE

func lexSqlUpdateTable(l *lexer) stateFn {
	l.lexSkipWhiteSpaces()
	return l.lexSqlIdentifier(tokenTypeSqlTable, lexSqlUpdateTableSet)
}

func lexSqlUpdateTableSet(l *lexer) stateFn {
	l.lexSkipWhiteSpaces()
	return l.lexMatch(tokenTypeSqlSet, "set", 0, lexSqlColumn)
}

func lexSqlColumn(l *lexer) stateFn {
	l.lexSkipWhiteSpaces()
	if l.end() {
		return nil
	}
	return l.lexSqlIdentifier(tokenTypeSqlColumn, lexSqlColumnEqual)
}

func lexSqlColumnEqual(l *lexer) stateFn {
	l.lexSkipWhiteSpaces()
	if l.next() == '=' {
		l.emit(tokenTypeSqlEqual)
		return lexSqlColumnEqualValue
	}
	l.errorToken("expecgted = ")
	return nil
}

func lexSqlColumnEqualValue(l *lexer) stateFn {
	l.lexSkipWhiteSpaces()
	return l.lexSqlValue(lexSqlCommaOrWhere)
}

func lexSqlCommaOrWhere(l *lexer) stateFn {
	l.lexSkipWhiteSpaces()
	if l.next() == ',' {
		l.emit(tokenTypeSqlComma)
		return lexSqlColumn
	}
	l.backup()
	return lexSqlWhere
}

// UNSUBSCRIBE

func lexSqlUnsubscribeFrom(l *lexer) stateFn {
	l.lexSkipWhiteSpaces()
	return l.lexMatch(tokenTypeSqlFrom, "from", 0, lexSqlUnsubscribeFromTable)
}

func lexSqlUnsubscribeFromTable(l *lexer) stateFn {
	return l.lexSqlIdentifier(tokenTypeSqlTable, nil)
}

// KEY 
// TAG

func lexSqlKeyTable(l *lexer) stateFn {
	return l.lexSqlIdentifier(tokenTypeSqlTable, lexSqlKeyColumn)
}

func lexSqlKeyColumn(l *lexer) stateFn {
	return l.lexSqlIdentifier(tokenTypeSqlColumn, nil)
}

// END SQL

// lexCommandST helper function to process status stop start commands
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
	l.errorToken("Invalid command:" + l.current())
	return nil
}

// lexCommandS helper function to process select subscribe status stop start commands
func lexCommandS(l *lexer) stateFn {
	switch l.next() {
	case 'e':
		return l.lexMatch(tokenTypeSqlSelect, "select", 2, lexSqlSelectStar)
	case 'u':
		return l.lexMatch(tokenTypeSqlSubscribe, "subscribe", 2, lexSqlSelectStar)
	case 't':
		return lexCommandST(l)
	}
	l.errorToken("Invalid command:" + l.current())
	return nil
}

// skipWhiteSpaces skips white space characters
func (l *lexer) lexSkipWhiteSpaces() {
	for rune := l.next(); unicode.IsSpace(rune); rune = l.next() {
	}
	l.backup()
	l.ignore()
}

// lexCommand is the initial state function
func lexCommand(l *lexer) stateFn {
	l.lexSkipWhiteSpaces()
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
	l.errorToken("Invalid command:" + l.current())
	return nil
}

// run scans the input by executing state functon until 
// the state is nil
func (l *lexer) run() {
	for state := lexCommand; state != nil; {
		state = state(l)
	}
	l.emit(tokenTypeEOF)
}

// lex scans the input by running lexer 
func lex(input string, tokens tokenConsumer) {
	l := &lexer{
		input:  input,
		tokens: tokens,
	}
	l.run()
}
