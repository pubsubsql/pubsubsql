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
	column string
	val    string
}

// sqlInsertAction is an action for sql insert statement
type sqlInsertAction struct {
	table  string
	values []columnValue
}

// sqlSelectAction is an action for sql select statement
type sqlSelectAction struct {
	table  string
	filter columnValue
}

// sqlUpdateAction is an action for sql update statement
type sqlUpdateAction struct {
	sqlAction
	table   string
	colVals []*columnValue
	filter  columnValue
}

func (a *sqlUpdateAction) addColVal(col string, val string) {
	a.colVals = append(a.colVals, &columnValue{column: col, val: val})
}

func (a *sqlUpdateAction) addFilter(col string, val string) {
	a.filter.column = col
	a.filter.val = val
}

// sqlDeleteAction is an action for sql delete statement
type sqlDeleteAction struct {
	table  string
	filter columnValue
}
