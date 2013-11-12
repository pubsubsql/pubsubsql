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

// record  
type record struct {
	values []string
}

// record factory
func newRecord(columns int) *record {
	r := record {
		values: make([]string, columns, columns),
	}
	return &r
}

// getValue retrieves value based on column ordinal
func (r *record) getValue(ordinal int) string {
	if len(r.values) > ordinal {
		return r.values[ordinal]
	} 
	return ""
}

// setValue sets value based on column ordinal 
func (r *record) setValue(ordinal int, val string) {
	l := len(r.values)
	if l <= ordinal {
		delta := ordinal - l + 1	
		temp := make([]string, delta)
		r.values = append(r.values, temp...) 	
	} 
	r.values[ordinal] = val	
}

