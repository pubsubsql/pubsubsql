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

package server

func lexMysqlConnectAddress(this *lexer) stateFn {
	this.skipWhiteSpaces()
	return this.lexSqlValue(nil)
}

// Helper function to process subscribe, status.
func lexMysqlCommandS(this *lexer) stateFn {
	switch this.next() {
	case 'u':
		return this.lexMatch(tokenTypeSqlSubscribe, "subscribe", 2, lexSqlSubscribe)
	case 't':
		return this.lexMatch(tokenTypeCmdStatus, "status", 2, nil)
	}
	return this.errorToken("Invalid command:" + this.current())
}

// Helper function to process mysql subscribe unsubscribe connect disconnect status tables commands.
func lexCmdMysql(this *lexer) stateFn {
	this.skipWhiteSpaces()
	switch this.next() {
	case 's': // subscribe, status
		return lexMysqlCommandS(this)
	case 'u': // unsubscribe
		return this.lexMatch(tokenTypeSqlUnsubscribe, "unsubscribe", 1, lexSqlUnsubscribeFrom)
	case 'c': // connect
		return this.lexMatch(tokenTypeCmdConnect, "connect", 1, lexMysqlConnectAddress)
	case 'd': // disconnect
		return this.lexMatch(tokenTypeCmdDisconnect, "disconnect", 1, nil)
	case 't': // tables
		return this.lexMatch(tokenTypeCmdTables, "tables", 1, nil)
	}
	return this.errorToken("Invalid command:" + this.current())
}
