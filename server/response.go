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

type responseStatusType int8

const (
	responseStatusOk  responseStatusType = iota // ok.
	responseStatusErr                           // error.
)

// response 
type response interface {
	getResponseStatus() responseStatusType
	String() string
}

// errorResponse 
type errorResponse struct {
	response
	msg string
}

func (r *errorResponse) getResponsStatus() responseStatusType {
	return responseStatusErr
}

func (r *errorResponse) String() string {
	return `{"status":"err" "msg":"` + r.msg + `"}`
}

// sqlInsertResponse is a response for sql insert statement
type sqlInsertResponse struct {
	response
	id string
}

func (r *sqlInsertResponse) getResponsStatus() responseStatusType {
	return responseStatusOk
}

func (r *sqlInsertResponse) String() string {
	return `{"response":"insert" "status":"ok" "id":"` + r.id + `"}`
}

//
