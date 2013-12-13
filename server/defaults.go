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

var DEBUG bool = false
var CHAN_RESPONSE_SENDER_BUFFER_SIZE = 10000 
var CHAN_TABLE_REQUESTS_BUFFER_SIZE = 1000
var PARSER_SQL_INSERT_REQUEST_COLUMN_CAPACITY = 10
var PARSER_SQL_UPDATE_REQUEST_COLUMN_CAPACITY = 10
var TABLE_COLUMNS_CAPACITY = 10
var TABLE_RECORDS_CAPACITY = 1000
var TABLE_GET_RECORDS_BY_TAG_CAPACITY = 20

func debug(str string) {
	if DEBUG {
		println("debug: " + str)
	}
}
