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

// column  
type column struct {
	name string	
	ordinal int
}

// table  
type table struct {
	name string	
	colMap map[string]*column 
	colSlice []*column 
}

// table factory 
func newTable(name string) *table {
	t := &table {
		name: name, 
		colMap: make(map[string]*column),
		colSlice: make([]*column, 0, 10),
	}
	t.addColumn("id")
	return t
} 

// addColumn adds column and returns column ordinal
func (t *table) addColumn(name string) int {
	ordinal := len(t.colSlice)
	col := &column {
		name: name, 
		ordinal: ordinal,
	}
	t.colMap[name] = col
	t.colSlice = append(t.colSlice, col)
	return ordinal	
}

// getAddColumn tries to retrieve existing column  or adds it if does not exist
func (t *table) getAddColumn(name string) int {
	col, ok := t.colMap[name]
	if ok {
		return col.ordinal
	} 
	return t.addColumn(name)
} 


