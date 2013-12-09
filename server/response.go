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

import "strconv"

type responseStatusType int8

const (
	responseStatusOk  responseStatusType = iota // ok.
	responseStatusErr                           // error.
)

// response
type response interface {
	getResponseStatus() responseStatusType
	String() string
	toNetworkReadyJSON() []byte
}

// json helper functions
func ok(builder *JSONBuilder) {
	builder.nameValue("status", "ok")
}

func id(builder *JSONBuilder, id string) {
	builder.nameValue("id", id)
}

func action(builder *JSONBuilder, action string) {
	builder.nameValue("action", action)
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

func (r *errorResponse) toNetworkReadyJSON() []byte {
	builder := networkReadyJSONBuilder()
	builder.beginObject()
	builder.nameValue("status", "err")
	builder.valueSeparator()
	builder.nameValue("msg", r.msg)
	builder.endObject()
	return builder.getNetworkBytes()
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

func (r *okResponse) toNetworkReadyJSON() []byte {
	builder := networkReadyJSONBuilder()
	builder.beginObject()
	ok(builder)
	builder.endObject()
	return builder.getNetworkBytes()
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

func (r *sqlInsertResponse) toNetworkReadyJSON() []byte {
	builder := networkReadyJSONBuilder()
	builder.beginObject()
	ok(builder)
	builder.valueSeparator()
	action(builder, "insert")
	builder.valueSeparator()
	id(builder, r.id)
	builder.endObject()
	return builder.getNetworkBytes()
}

// sqlSelectResponse is a response for sql select statement
type sqlSelectResponse struct {
	response
	columns []*column
	records []*record
}

func (r *sqlSelectResponse) data(builder *JSONBuilder) {
	builder.nameIntValue("rows", len(r.records)) 
	builder.valueSeparator()
	builder.string("data")
	builder.nameSeparator()
	builder.beginArray()
	for recIndex, rec := range r.records {
		// another row
		if recIndex != 0 {
			builder.valueSeparator()
		}
		builder.beginObject()
		// columns and values
		for colIndex, col := range r.columns {
			if colIndex != 0 {
				builder.valueSeparator()
			}
			builder.nameValue(col.name, rec.getValue(colIndex))
		}
		builder.endObject()
	}
	builder.endArray()
}

func (r *sqlSelectResponse) toNetworkReadyJSON() []byte {
	builder := networkReadyJSONBuilder()
	builder.beginObject()
	ok(builder)
	builder.valueSeparator()
	action(builder, "select")
	builder.valueSeparator()
	r.data(builder)
	builder.endObject()
	return builder.getNetworkBytes()
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

func (r *sqlDeleteResponse) toNetworkReadyJSON() []byte {
	builder := networkReadyJSONBuilder()
	builder.beginObject()
	ok(builder)
	builder.valueSeparator()
	action(builder, "delete")
	builder.valueSeparator()
	builder.nameIntValue("rows", r.deleted)
	builder.endObject()
	return builder.getNetworkBytes()
}

// sqlUpdateResponse
type sqlUpdateResponse struct {
	response
	updated int
}

func (r *sqlUpdateResponse) toNetworkReadyJSON() []byte {
	builder := networkReadyJSONBuilder()
	builder.beginObject()
	ok(builder)
	builder.valueSeparator()
	action(builder, "update")
	builder.valueSeparator()
	builder.nameIntValue("rows", r.updated)
	builder.endObject()
	return builder.getNetworkBytes()
}

// sqlSubscribeResponse
type sqlSubscribeResponse struct {
	response
	pubsubid uint64
}

func (r *sqlSubscribeResponse) toNetworkReadyJSON() []byte {
	builder := networkReadyJSONBuilder()
	builder.beginObject()
	ok(builder)
	builder.valueSeparator()
	action(builder, "subscribe")
	builder.valueSeparator()
	builder.nameValue("pubsubid", strconv.FormatUint(r.pubsubid, 10))
	builder.endObject()
	return builder.getNetworkBytes()
}

func newSubscribeResponse(sub *subscription) response {
	return &sqlSubscribeResponse{
		pubsubid: sub.id,
	}
}

// sqlActionDataResponse
type sqlActionDataResponse struct {
	sqlSelectResponse
	pubsubid uint64
}

// sqlActionAddResponse
type sqlActionAddResponse struct {
	sqlActionDataResponse
}

// sqlActionInsertResponse
type sqlActionInsertResponse struct {
	sqlActionDataResponse
}

// sqlActonDeleteResponse
type sqlActionDeleteResponse struct {
	response
	id       string
	pubsubid uint64
}

// sqlActionRemoveResponse
type sqlActionRemoveResponse struct {
	response
	id       string
	pubsubid uint64
}

// sqlActionUpdateResponse
type sqlActionUpdateResponse struct {
	response
	pubsubid uint64
	cols     []*column
	rec      *record
}

func newSqlActionUpdateResponse(pubsubid uint64, cols []*column, rec *record) *sqlActionUpdateResponse {
	res := sqlActionUpdateResponse{
		pubsubid: pubsubid,
		cols:     cols,
	}
	// copy updated data
	l := len(cols)
	res.rec = &record{
		values: make([]string, l, l),
	}
	for idx, col := range cols {
		res.rec.setValue(idx, rec.getValue(col.ordinal))
	}
	return &res
}

// sqlUnsubscribeResponse
type sqlUnsubscribeResponse struct {
	response
	unsubscribed int
}
