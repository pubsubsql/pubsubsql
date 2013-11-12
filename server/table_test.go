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

import "testing"

func validateTableRecordsCount(t *testing.T, tbl *table, expected int) {
	val := tbl.getRecordCount()
	if val != expected {
		t.Errorf("records count do not match expected:%d but got:%d", expected, val)
	}
}

func TestTable1(t *testing.T) {
	tbl := newTable("table1")
	tbl.getAddColumn("col1")
	r := tbl.addRecord()
	validateTableRecordsCount(t, tbl, 1)
	validateRecordValuesCount(t, r, 2)
	validateRecordValue(t, r, 0, "0")
	//
	r = tbl.getRecord(0)
	validateTableRecordsCount(t, tbl, 1)
	validateRecordValuesCount(t, r, 2)
	validateRecordValue(t, r, 0, "0")
}

func TestTable2(t *testing.T) {
	tbl := newTable("table1")
	tbl.getAddColumn("col1")
	tbl.getAddColumn("col2")
	tbl.getAddColumn("col3")
	tbl.getAddColumn("col3")
	col1 := tbl.getColumn("col1").ordinal
	col2 := tbl.getColumn("col2").ordinal
	col3 := tbl.getColumn("col3").ordinal
	//
	r := tbl.addRecord()
	validateTableRecordsCount(t, tbl, 1)
	validateRecordValuesCount(t, r, 4)
	validateRecordValue(t, r, 0, "0")
	r = tbl.getRecord(0)
	validateTableRecordsCount(t, tbl, 1)
	validateRecordValuesCount(t, r, 4)
	validateRecordValue(t, r, 0, "0")
	//	
	r = tbl.addRecord()
	validateTableRecordsCount(t, tbl, 2)
	validateRecordValuesCount(t, r, 4)
	validateRecordValue(t, r, 0, "1")
	validateRecordValue(t, r, 1, "")
	validateRecordValue(t, r, 2, "")
	validateRecordValue(t, r, 3, "")
	r.setValue(col1, "val1")
	r.setValue(col2, "val2")
	r.setValue(col3, "val3")
	validateRecordValue(t, r, 0, "1")
	validateRecordValue(t, r, 1, "val1")
	validateRecordValue(t, r, 2, "val2")
	validateRecordValue(t, r, 3, "val3")
	r = tbl.getRecord(1)
	validateTableRecordsCount(t, tbl, 2)
	validateRecordValuesCount(t, r, 4)
	validateRecordValue(t, r, 0, "1")
	validateRecordValue(t, r, 1, "val1")
	validateRecordValue(t, r, 2, "val2")
	validateRecordValue(t, r, 3, "val3")
}
