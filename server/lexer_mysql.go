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

// Helper function to process mysql subscribe unsubscribe connect disconnect commands.
func lexSqlMysql(this *lexer) stateFn {
	this.skipWhiteSpaces()
	switch this.next() {
	case 's':
		return this.lexMatch(tokenTypeSqlSubscribe, "subscribe", 1, lexSqlSubscribe)
	case 'u':
		return this.lexMatch(tokenTypeSqlUnsubscribe, "unsubscribe", 1, lexSqlUnsubscribeFrom)
	case 'c':
		return this.lexMatch(tokenTypeSqlConnect, "connect", 1, lexSqlConnectValue)
	case 'd':
		return this.lexMatch(tokenTypeSqlDisconnect, "disconnect", 1, nil)
	}
	return this.errorToken("Invalid command:" + this.current())
}
