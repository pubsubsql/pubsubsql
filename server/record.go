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

// link
type link struct {
	pubsub *pubSub
	tg     *tag
}

func (l *link) clear() {
	l.pubsub = nil
	l.tg = nil
}

// record  
type record struct {
	values []string
	links  []link
}

// record factory
func newRecord(columns int, idx int) *record {
	r := record{
		values: make([]string, columns, columns),
	}
	r.setValue(0, strconv.Itoa(idx))
	return &r
}

// Returns record index in a table.
func (r *record) idx() int {
	i, err := strconv.Atoi(r.values[0])
	if err != nil {
		return -1
	}
	return i
}

// Returns record index in a table as string.
func (r *record) idAsString() string {
	return r.values[0]
}

// Returns value based on column ordinal.
// Empty string is returned for invalid ordinal. 
func (r *record) getValue(ordinal int) string {
	if len(r.values) > ordinal {
		return r.values[ordinal]
	}
	return ""
}

// Sets value based on column ordinal. 
// Automatically adjusts the record if ordinal is invalid.
func (r *record) setValue(ordinal int, val string) {
	l := len(r.values)
	if l <= ordinal {
		delta := ordinal - l + 1
		temp := make([]string, delta)
		r.values = append(r.values, temp...)
	}
	r.values[ordinal] = val
}
