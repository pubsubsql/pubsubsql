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

func newErrorResponse(msg string) *errorResponse {
	return &errorResponse{msg: msg}
}

func (r *errorResponse) getResponsStatus() responseStatusType {
	return responseStatusErr
}

func (r *errorResponse) String() string {
	return `{"status":"err" "msg":"` + r.msg + `"}`
}

// okResponse
type okResponse struct {
	response
}

func newOkResponse() *okResponse {
	return &okResponse{}
}

func (r *okResponse) getResponsStatus() responseStatusType {
	return responseStatusOk
}

func (r *okResponse) String() string {
	return `{"status":"ok"}`
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

// sqlSelectResponse is a response for sql select statement
type sqlSelectResponse struct {
	response
	columns []*column
	records []*record
}

func (r *sqlSelectResponse) copyRecordData(source *record) {
	l := len(r.columns)
	dest := &record{
		values: make([]string, l, l),
	}
	for idx, col := range r.columns {
		dest.setValue(idx, source.getValue(col.ordinal))
	}
	addRecordToSlice(&r.records, dest)
}

// sqlDeleteResponse
type sqlDeleteResponse struct {
	response
	deleted int
}

// sqlUpdateResponse
type sqlUpdateResponse struct {
	response
	updated int
}

// sqlSubscribeResponse
type sqlSubscribeResponse struct {
	response
	pubsubid uint64
}

func newSubscribeResponse(sub *subscription) response {
	return &sqlSubscribeResponse{
		pubsubid: sub.id,
	}
}

// sqlActionAddResponse
type sqlActionAddResponse struct {
	sqlSelectResponse
	pubsubid uint64
}

// sqlUnsubscribeResponse
type sqlUnsubscribeResponse struct {
	response
	unsubscribed int
}
