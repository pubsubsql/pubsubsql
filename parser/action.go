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

type actionType uint8

const (
	actionTypeError actionType = iota // error action indicates scan or parse error.
	actionTypeCmd                     // cmd actions help, status etc.
	actionTypeSql                     // sql actins insert, update etc.
)

// action
type action interface {
	getActionType() actionType
}

// errorAction is an error action
type errorAction struct {
	action
	err string
}

func (act errorAction) getActionType() actionType {
	return actionTypeError
}

// sqlAction is a generic sql action
type sqlAction struct {
	action
}

func (act sqlAction) getActionType() actionType {
	return actionTypeSql
}

// cmdAction is a generic command action 
type cmdAction struct {
	action
}

func (act cmdAction) getActionType() actionType {
	return actionTypeCmd
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

// sqlInsertAction is an action for sql insert statement
type sqlInsertAction struct {
	sqlAction
	table   string
	colVals []*columnValue
}

func (a *sqlInsertAction) addColumn(col string) {
	a.colVals = append(a.colVals, &columnValue{col: col})
}

func (a *sqlInsertAction) addColVal(col string, val string) {
	a.colVals = append(a.colVals, &columnValue{col: col, val: val})
}

func (a *sqlInsertAction) setValueAt(idx int, val string) {
	a.colVals[idx].val = val
}

// sqlSelectAction is an action for sql select statement
type sqlSelectAction struct {
	sqlAction
	table  string
	filter sqlFilter
}

// sqlUpdateAction is an action for sql update statement
type sqlUpdateAction struct {
	sqlAction
	table   string
	colVals []*columnValue
	filter  sqlFilter
}

func (a *sqlUpdateAction) addColVal(col string, val string) {
	a.colVals = append(a.colVals, &columnValue{col: col, val: val})
}

// sqlDeleteAction is an action for sql delete statement
type sqlDeleteAction struct {
	sqlAction
	table  string
	filter sqlFilter
}

// sqlSubscribeAction is an action for sql subscribe statement
type sqlSubscribeAction struct {
	sqlAction
	table  string
	filter sqlFilter
}

// sqlUnsubscribeAction is an action for sql unsubscribe statement
type sqlUnsubscribeAction struct {
	sqlAction
	table string
}

// sqlKeyAction is an action for key statement 
// key defines unique index
type sqlKeyAction struct {
	sqlAction
	table  string
	column string
}

// sqlTagAction is an action for tag statement 
// tag defines non-unique index
type sqlTagAction struct {
	sqlAction
	table  string
	column string
}
