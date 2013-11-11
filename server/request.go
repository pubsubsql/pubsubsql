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

type requestType uint8

const (
	requestTypeError requestType = iota // error request indicates scan or parse error.
	requestTypeCmd                     // cmd requests help, status etc.
	requestTypeSql                     // sql actins insert, update etc.
)

// request
type request interface {
	getRequestType() requestType
}

// errorRequest is an error request
type errorRequest struct {
	request
	err string
}

func (act errorRequest) getRequestType() requestType {
	return requestTypeError
}

// sqlRequest is a generic sql request
type sqlRequest struct {
	request
}

func (act sqlRequest) getRequestType() requestType {
	return requestTypeSql
}

// cmdRequest is a generic command request 
type cmdRequest struct {
	request
}

func (act cmdRequest) getRequestType() requestType {
	return requestTypeCmd
}

// columnValue pair 
type columnValue struct {
	col string
	val string
}

// temporarely stub for sqlFilter type that will be more capble in future versions
type sqlFilter struct {
	columnValue
}

func (f *sqlFilter) addFilter(col string, val string) {
	f.col = col
	f.val = val
}

// sqlInsertRequest is an request for sql insert statement
type sqlInsertRequest struct {
	sqlRequest
	table   string
	colVals []*columnValue
}

func (a *sqlInsertRequest) addColumn(col string) {
	a.colVals = append(a.colVals, &columnValue{col: col})
}

func (a *sqlInsertRequest) addColVal(col string, val string) {
	a.colVals = append(a.colVals, &columnValue{col: col, val: val})
}

func (a *sqlInsertRequest) setValueAt(idx int, val string) {
	a.colVals[idx].val = val
}

// sqlSelectRequest is an request for sql select statement
type sqlSelectRequest struct {
	sqlRequest
	table  string
	filter sqlFilter
}

// sqlUpdateRequest is an request for sql update statement
type sqlUpdateRequest struct {
	sqlRequest
	table   string
	colVals []*columnValue
	filter  sqlFilter
}

func (a *sqlUpdateRequest) addColVal(col string, val string) {
	a.colVals = append(a.colVals, &columnValue{col: col, val: val})
}

// sqlDeleteRequest is an request for sql delete statement
type sqlDeleteRequest struct {
	sqlRequest
	table  string
	filter sqlFilter
}

// sqlSubscribeRequest is an request for sql subscribe statement
type sqlSubscribeRequest struct {
	sqlRequest
	table  string
	filter sqlFilter
}

// sqlUnsubscribeRequest is an request for sql unsubscribe statement
type sqlUnsubscribeRequest struct {
	sqlRequest
	table string
}

// sqlKeyRequest is an request for key statement 
// key defines unique index
type sqlKeyRequest struct {
	sqlRequest
	table  string
	column string
}

// sqlTagRequest is an request for tag statement 
// tag defines non-unique index
type sqlTagRequest struct {
	sqlRequest
	table  string
	column string
}
