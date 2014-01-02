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
	toNetworkReadyJSON() ([]byte, bool)
}

type requestIdResponse struct {
	response
	requestId uint32
}

func getRequestIdResponse(requestId uint32) requestIdResponse {
	return requestIdResponse {requestId: requestId,}
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

func (this *errorResponse) getResponsStatus() responseStatusType {
	return responseStatusErr
}

func (this *errorResponse) String() string {
	return `{"status":"err" "msg":"` + this.msg + `"}`
}

func (this *errorResponse) toNetworkReadyJSON() ([]byte, bool) {
	builder := networkReadyJSONBuilder()
	builder.beginObject()
	builder.nameValue("status", "err")
	builder.valueSeparator()
	builder.nameValue("msg", this.msg)
	builder.endObject()
	return builder.getNetworkBytes(), false
}

// okResponse
type okResponse struct {
	response
}

func newOkResponse() *okResponse {
	return &okResponse{}
}

func (this *okResponse) getResponsStatus() responseStatusType {
	return responseStatusOk
}

func (this *okResponse) String() string {
	return `{"status":"ok"}`
}

func (this *okResponse) toNetworkReadyJSON() ([]byte, bool) {
	builder := networkReadyJSONBuilder()
	builder.beginObject()
	ok(builder)
	builder.endObject()
	return builder.getNetworkBytes(), false
}

// cmdStatusResponse
type cmdStatusResponse struct {
	response
	connections int
}

func (this *cmdStatusResponse) toNetworkReadyJSON() ([]byte, bool) {
	builder := networkReadyJSONBuilder()
	builder.beginObject()
	ok(builder)
	builder.nameIntValue("connections", this.connections)
	builder.endObject()
	return builder.getNetworkBytes(), false
}

// sqlInsertResponse is a response for sql insert statement
type sqlInsertResponse struct {
	requestIdResponse
	id string
}

func newSqlInsertResponse(id string, requestId uint32) *sqlInsertResponse  {
	return &sqlInsertResponse {
		id: id,
		requestIdResponse: getRequestIdResponse(requestId),
	} 
}

func (this *sqlInsertResponse) getResponsStatus() responseStatusType {
	return responseStatusOk
}

func (this *sqlInsertResponse) String() string {
	return `{"response":"insert" "status":"ok" "id":"` + this.id + `"}`
}

func (this *sqlInsertResponse) toNetworkReadyJSON() ([]byte, bool) {
	builder := networkReadyJSONBuilder()
	builder.beginObject()
	ok(builder)
	builder.valueSeparator()
	action(builder, "insert")
	builder.valueSeparator()
	id(builder, this.id)
	builder.endObject()
	return builder.getNetworkBytes(), false
}

// sqlSelectResponse is a response for sql select statement
type sqlSelectResponse struct {
	response
	columns []*column
	records []*record
	batch   int
}

func row(builder *JSONBuilder, columns []*column, rec *record) {
	builder.beginObject()
	// columns and values
	for colIndex, col := range columns {
		if colIndex != 0 {
			builder.valueSeparator()
		}
		builder.nameValue(col.name, rec.getValue(colIndex))
	}
	builder.endObject()
}

func (this *sqlSelectResponse) data(builder *JSONBuilder) bool {
	more := len(this.records) > config.DATA_BATCH_SIZE
	records := this.records
	if more {
		records = this.records[0:config.DATA_BATCH_SIZE]
		this.records = this.records[config.DATA_BATCH_SIZE:]
		//
		this.batch++
		builder.nameIntValue("batch", this.batch)
		builder.valueSeparator()
	}
	builder.nameIntValue("rows", len(this.records))
	builder.valueSeparator()
	builder.string("data")
	builder.nameSeparator()
	builder.beginArray()
	for recIndex, rec := range records {
		// another row
		if recIndex != 0 {
			builder.objectSeparator()
		}
		row(builder, this.columns, rec)
	}
	builder.endArray()
	return more
}

func (this *sqlSelectResponse) toNetworkReadyJSON() ([]byte, bool) {
	builder := networkReadyJSONBuilder()
	builder.beginObject()
	ok(builder)
	builder.valueSeparator()
	action(builder, "select")
	builder.valueSeparator()
	more := this.data(builder)
	builder.endObject()
	return builder.getNetworkBytes(), more
}

func (this *sqlSelectResponse) copyRecordData(source *record) {
	l := len(this.columns)
	dest := &record{
		values: make([]string, l, l),
	}
	for idx, col := range this.columns {
		dest.setValue(idx, source.getValue(col.ordinal))
	}
	addRecordToSlice(&this.records, dest)
}

// sqlDeleteResponse
type sqlDeleteResponse struct {
	response
	deleted int
}

func (this *sqlDeleteResponse) toNetworkReadyJSON() ([]byte, bool) {
	builder := networkReadyJSONBuilder()
	builder.beginObject()
	ok(builder)
	builder.valueSeparator()
	action(builder, "delete")
	builder.valueSeparator()
	builder.nameIntValue("rows", this.deleted)
	builder.endObject()
	return builder.getNetworkBytes(), false
}

// sqlUpdateResponse
type sqlUpdateResponse struct {
	response
	updated int
}

func (this *sqlUpdateResponse) toNetworkReadyJSON() ([]byte, bool) {
	builder := networkReadyJSONBuilder()
	builder.beginObject()
	ok(builder)
	builder.valueSeparator()
	action(builder, "update")
	builder.valueSeparator()
	builder.nameIntValue("rows", this.updated)
	builder.endObject()
	return builder.getNetworkBytes(), false
}

// sqlSubscribeResponse
type sqlSubscribeResponse struct {
	response
	pubsubid uint64
}

func (this *sqlSubscribeResponse) toNetworkReadyJSON() ([]byte, bool) {
	builder := networkReadyJSONBuilder()
	builder.beginObject()
	ok(builder)
	builder.valueSeparator()
	action(builder, "subscribe")
	builder.valueSeparator()
	builder.nameValue("pubsubid", strconv.FormatUint(this.pubsubid, 10))
	builder.endObject()
	return builder.getNetworkBytes(), false
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

func (this *sqlActionDataResponse) toNetworkReadyJSONHelper(act string) ([]byte, bool) {
	builder := networkReadyJSONBuilder()
	builder.beginObject()
	ok(builder)
	builder.valueSeparator()
	action(builder, act)
	builder.valueSeparator()
	builder.nameValue("pubsubid", strconv.FormatUint(this.pubsubid, 10))
	builder.valueSeparator()
	more := this.data(builder)
	builder.endObject()
	return builder.getNetworkBytes(), more
}

// sqlActionAddResponse
type sqlActionAddResponse struct {
	sqlActionDataResponse
}

func (this *sqlActionAddResponse) toNetworkReadyJSON() ([]byte, bool) {
	return this.toNetworkReadyJSONHelper("add")
}

// sqlActionInsertResponse
type sqlActionInsertResponse struct {
	sqlActionDataResponse
}

func (this *sqlActionInsertResponse) toNetworkReadyJSON() ([]byte, bool) {
	return this.toNetworkReadyJSONHelper("insert")
}

// sqlActonDeleteResponse
type sqlActionDeleteResponse struct {
	response
	id       string
	pubsubid uint64
}

func (this *sqlActionDeleteResponse) toNetworkReadyJSON() ([]byte, bool) {
	builder := networkReadyJSONBuilder()
	builder.beginObject()
	ok(builder)
	builder.valueSeparator()
	action(builder, "delete")
	builder.valueSeparator()
	builder.nameValue("pubsubid", strconv.FormatUint(this.pubsubid, 10))
	builder.valueSeparator()
	builder.nameValue("id", this.id)
	builder.endObject()
	return builder.getNetworkBytes(), false
}

// sqlActionRemoveResponse
type sqlActionRemoveResponse struct {
	response
	id       string
	pubsubid uint64
}

func (this *sqlActionRemoveResponse) toNetworkReadyJSON() ([]byte, bool) {
	builder := networkReadyJSONBuilder()
	builder.beginObject()
	ok(builder)
	builder.valueSeparator()
	action(builder, "remove")
	builder.valueSeparator()
	builder.nameValue("pubsubid", strconv.FormatUint(this.pubsubid, 10))
	builder.valueSeparator()
	builder.nameValue("id", this.id)
	builder.endObject()
	return builder.getNetworkBytes(), false
}

// sqlActionUpdateResponse
type sqlActionUpdateResponse struct {
	response
	pubsubid uint64
	cols     []*column
	rec      *record
}

func (this *sqlActionUpdateResponse) toNetworkReadyJSON() ([]byte, bool) {
	builder := networkReadyJSONBuilder()
	builder.beginObject()
	ok(builder)
	builder.valueSeparator()
	action(builder, "update")
	builder.valueSeparator()
	builder.nameValue("pubsubid", strconv.FormatUint(this.pubsubid, 10))
	builder.valueSeparator()

	builder.string("data")
	builder.nameSeparator()
	builder.beginArray()
	row(builder, this.cols, this.rec)
	builder.endArray()

	builder.endObject()
	return builder.getNetworkBytes(), false
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

func (this *sqlUnsubscribeResponse) toNetworkReadyJSON() ([]byte, bool) {
	builder := networkReadyJSONBuilder()
	builder.beginObject()
	ok(builder)
	builder.valueSeparator()
	action(builder, "unsubscribe")
	builder.valueSeparator()
	builder.nameIntValue("subscriptions", this.unsubscribed)
	builder.endObject()
	return builder.getNetworkBytes(), false
}
