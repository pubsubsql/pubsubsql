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

type configuration struct {
	// logger
	LOG_DEBUG bool
 	LOG_INFO bool
	LOG_WARN bool
	LOG_ERR bool

	// resources
	CHAN_RESPONSE_SENDER_BUFFER_SIZE int
	CHAN_TABLE_REQUESTS_BUFFER_SIZE int
	PARSER_SQL_INSERT_REQUEST_COLUMN_CAPACITY int
	PARSER_SQL_UPDATE_REQUEST_COLUMN_CAPACITY int
	TABLE_COLUMNS_CAPACITY int 
	TABLE_RECORDS_CAPACITY int 
	TABLE_GET_RECORDS_BY_TAG_CAPACITY int
}

func defaultConfig() configuration {
	return configuration {
		// logger
		LOG_DEBUG: true,
		LOG_INFO: true,
		LOG_WARN: true, 
		LOG_ERR: true,
		// resources
		CHAN_RESPONSE_SENDER_BUFFER_SIZE: 10000, 
		CHAN_TABLE_REQUESTS_BUFFER_SIZE: 1000,
		PARSER_SQL_INSERT_REQUEST_COLUMN_CAPACITY: 10,
		PARSER_SQL_UPDATE_REQUEST_COLUMN_CAPACITY: 10,
		TABLE_COLUMNS_CAPACITY: 10, 
		TABLE_RECORDS_CAPACITY: 1000, 
		TABLE_GET_RECORDS_BY_TAG_CAPACITY: 20,
	}
}

var config = defaultConfig()

