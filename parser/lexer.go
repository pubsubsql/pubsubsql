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
	tokenTypeSqlCreate                            // create -> table
	tokenTypeSqlTable                             // table -> table name
	tokenTypeSqlInsert                            // insert -> into
	tokenTypeSqlInto                              // into -> table
	tokenTypeSqlUpdate                            // update -> table	
	tokenTypeSqlDelete                            // delete -> from
	tokenTypeSqlFrom                              // from -> table
	tokenTypeSqlSelect                            // select
	tokenTypeSqlSubscribe                         // subscribe
	tokenTypeSqlUnsubscribe                       // unsubscribe 
	tokenTypeSqlWhere                             // where
	tokenTypeSqlStar                              // *
	tokenTypeSqlEqual                             // =
	tokenTypeSqlLeftParenthesis                   // (
	tokenTypeSqlRightParenthesis                  // )
	tokenTypeSqlComma                             // ,
	tokenTypeSqlId                                // starts with alpha contains alnum 
	tokenTypeSqlValue                             // continous sequence of chars delimited by WHITE SPACE | ' | , | ( | ) 
	tokenTypeSqlAnsiQuote                         // '
	tokenTypeSqlString                            // ' + any character + '  '' becomes ' inside the string
	tokenTypeWhiteSpace                           // \n,\r,\t, space
)

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

func isWhiteSpace(rune int32) bool {
	return (unicode.IsSpace(rune) || rune == 0)
}

// read till first unicode White space character
func (l *lexer) scanTillWhiteSpace() {
	for rune := l.next(); !isWhiteSpace(rune); rune = l.next() {

	}
}

// match reads input and matches against the string
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

// lexMatchCommand matches expected command
func (l *lexer) lexMatchCommand(typ tokenType, command string, skip int, fn stateFn) stateFn {
	if l.match(command, skip) {
		l.emit(typ)
		return fn
	}
	l.errorToken("Invalid command:" + l.current())
	return nil
}

// lexCommandST helper function to process status stop start commands
func lexCommandST(l *lexer) stateFn {
	switch l.next() {
	case 'a':
		if l.next() == 'r' {
			return l.lexMatchCommand(tokenTypeCmdStart, "start", 4, nil)
		}
		return l.lexMatchCommand(tokenTypeCmdStatus, "status", 4, nil)
	default:
		return l.lexMatchCommand(tokenTypeCmdStop, "stop", 3, nil)
	}
	l.errorToken("Invalid command:" + l.current())
	return nil
}

// lexCommandS helper function to process select subscribe status stop start commands
func lexCommandS(l *lexer) stateFn {
	switch l.next() {
	case 'e':
		return l.lexMatchCommand(tokenTypeSqlSelect, "select", 2, nil)
	case 'u':
		return l.lexMatchCommand(tokenTypeSqlSubscribe, "subscribe", 2, nil)
	case 't':
		return lexCommandST(l)
	}
	l.errorToken("Invalid command:" + l.current())
	return nil
}

// skipWhiteSpaces skips white space characters
func lexSkipWhiteSpaces(l *lexer) {
	for rune := l.next(); unicode.IsSpace(rune); rune = l.next() {
	}
	l.backup()
	l.ignore()
}

// lexCommand is the initial state function
func lexCommand(l *lexer) stateFn {
	lexSkipWhiteSpaces(l)
	switch l.next() {
	case 'u': // update unsubscribe
		if l.next() == 'p' {
			return l.lexMatchCommand(tokenTypeSqlUpdate, "update", 2, nil)
		}
		return l.lexMatchCommand(tokenTypeSqlUnsubscribe, "unsubscribe", 2, nil)
	case 's': // select subscribe status stop start
		return lexCommandS(l)
	case 'i': // insert
		return l.lexMatchCommand(tokenTypeSqlInsert, "insert", 1, nil)
	case 'd': // delete
		return l.lexMatchCommand(tokenTypeSqlDelete, "delete", 1, nil)
	case 'h': // help
		return l.lexMatchCommand(tokenTypeCmdHelp, "help", 1, nil)
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
