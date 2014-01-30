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

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestConnectDisconnect(t *testing.T) {
	register("TestConnectDisconnect", t)
	client := NewClient()
	ASSERT_CONNECT(client, ADDRESS, true)
	ASSERT_CONNECTED(client, true)
	ASSERT_DISCONNECT(client)
	ASSERT_CONNECT(client, "addresswithnoport", false)
	ASSERT_CONNECTED(client, false)
	ASSERT_DISCONNECT(client)
	ASSERT_CONNECT(client, "addresswithnoport:", false)
	ASSERT_CONNECTED(client, false)
	ASSERT_DISCONNECT(client)
	ASSERT_CONNECT(client, "localhost:7778", false)
	ASSERT_CONNECTED(client, false)
	ASSERT_DISCONNECT(client)
}

func TestExecuteStatus(t *testing.T) {
	register("TestExecuteStatus", t)
	client := NewClient()
	ASSERT_CONNECT(client, ADDRESS, true)
	ASSERT_EXECUTE(client, "status", true)
	ASSERT_ACTION(client, "status")
	ASSERT_DISCONNECT(client)
}

func TestExecuteInvalidCommand(t *testing.T) {
	register("TestExecuteInvalidCommand", t)
	client := NewClient()
	ASSERT_CONNECT(client, ADDRESS, true)
	ASSERT_EXECUTE(client, "blablabla", false)
	ASSERT_DISCONNECT(client)
}

func TestInsertOneRow(t *testing.T) {
	register("TestInsertOneRow", t)
	newtable()
	client := NewClient()
	ASSERT_CONNECT(client, ADDRESS, true)
	command := fmt.Sprintf("insert into %v (col1, col2, col3) values (1:col1, 1:col2, 1:col3)", TABLE)
	ASSERT_EXECUTE(client, command, true)
	ASSERT_ACTION(client, "insert")
	ASSERT_ROW_COUNT(client, 1)
	ASSERT_NEXT_ROW(client, true)
	ASSERT_ID(client)
	ASSERT_VALUE(client, "col1", "1:col1", true)
	ASSERT_VALUE(client, "col2", "1:col2", true)
	ASSERT_VALUE(client, "col3", "1:col3", true)
	ASSERT_HAS_COLUMN(client, "col1", true)
	ASSERT_HAS_COLUMN(client, "col2", true)
	ASSERT_HAS_COLUMN(client, "col3", true)
	ASSERT_COLUMN_COUNT(client, 4) // including id
	ASSERT_NEXT_ROW(client, false)
	ASSERT_DISCONNECT(client)
}

func TestInsertManyRows(t *testing.T) {
	register("TestInsertManyRows", t)
	newtable()
	client := NewClient()
	ASSERT_CONNECT(client, ADDRESS, true)
	command := fmt.Sprintf("insert into %v (col1, col2, col3) values (1:col1, 1:col2, 1:col3)", TABLE)
	for r := 0; r < ROWS; r++ {
		ASSERT_EXECUTE(client, command, true)
		ASSERT_ACTION(client, "insert")
		ASSERT_ROW_COUNT(client, 1)
		ASSERT_NEXT_ROW(client, true)
		//
		ASSERT_ID(client)
		ASSERT_VALUE(client, "col1", "1:col1", true)
		ASSERT_VALUE(client, "col2", "1:col2", true)
		ASSERT_VALUE(client, "col3", "1:col3", true)
		ASSERT_HAS_COLUMN(client, "col1", true)
		ASSERT_HAS_COLUMN(client, "col2", true)
		ASSERT_HAS_COLUMN(client, "col3", true)
		ASSERT_COLUMN_COUNT(client, 4) // including id
		//
		ASSERT_NEXT_ROW(client, false)
	}
	ASSERT_DISCONNECT(client)
}

func TestSelectOneRow(t *testing.T) {
	register("TestSelectOneRow", t)
	newtable()
	insertRow()
	client := NewClient()
	ASSERT_CONNECT(client, ADDRESS, true)
	// select one row
	command := fmt.Sprintf("select * from %v", TABLE)
	ASSERT_EXECUTE(client, command, true)
	ASSERT_ACTION(client, "select")
	ASSERT_ROW_COUNT(client, 1)
	ASSERT_NEXT_ROW(client, true)
	//
	ASSERT_ID(client)
	ASSERT_VALUE(client, "col1", "1:col1", true)
	ASSERT_VALUE(client, "col2", "1:col2", true)
	ASSERT_VALUE(client, "col3", "1:col3", true)
	ASSERT_HAS_COLUMN(client, "col1", true)
	ASSERT_HAS_COLUMN(client, "col2", true)
	ASSERT_HAS_COLUMN(client, "col3", true)
	ASSERT_COLUMN_COUNT(client, 4) // including id
	//
	ASSERT_NEXT_ROW(client, false)
	ASSERT_DISCONNECT(client)
}

func TestSelectManyRows(t *testing.T) {
	register("TestSelectRow", t)
	newtable()
	insertRows()
	client := NewClient()
	ASSERT_CONNECT(client, ADDRESS, true)
	command := fmt.Sprintf("select * from %v", TABLE)
	ASSERT_EXECUTE(client, command, true)
	ASSERT_ACTION(client, "select")
	ASSERT_ROW_COUNT(client, ROWS)
	for row := 0; row < ROWS; row++ {
		ASSERT_NEXT_ROW(client, true)
		ASSERT_ID(client)
		ASSERT_VALUE(client, "col1", fmt.Sprintf("%v:col1", row), true)
		ASSERT_VALUE(client, "col2", fmt.Sprintf("%v:col2", row), true)
		ASSERT_VALUE(client, "col3", fmt.Sprintf("%v:col3", row), true)
		ASSERT_HAS_COLUMN(client, "col1", true)
		ASSERT_HAS_COLUMN(client, "col2", true)
		ASSERT_HAS_COLUMN(client, "col3", true)
		ASSERT_COLUMN_COUNT(client, 4) // including id
	}
	//
	ASSERT_NEXT_ROW(client, false)
	ASSERT_DISCONNECT(client)
}

func TestUpdateOneRow(t *testing.T) {
	register("TestUpdateOneRow", t)
	newtable()
	insertRow()
	client := NewClient()
	ASSERT_CONNECT(client, ADDRESS, true)
	command := fmt.Sprintf("update %v set col1 = newvalue", TABLE)
	ASSERT_EXECUTE(client, command, true)
	ASSERT_ACTION(client, "update")
	ASSERT_ROW_COUNT(client, 1)
	ASSERT_DISCONNECT(client)
}

func TestUpdateManyRows(t *testing.T) {
	register("TestUpdateManyRow", t)
	newtable()
	insertRows()
	client := NewClient()
	ASSERT_CONNECT(client, ADDRESS, true)
	command := fmt.Sprintf("update %v set col1 = newvalue", TABLE)
	ASSERT_EXECUTE(client, command, true)
	ASSERT_ACTION(client, "update")
	ASSERT_ROW_COUNT(client, ROWS)
	ASSERT_DISCONNECT(client)
}

func TestDeleteOneRow(t *testing.T) {
	register("TestDeleteOneRow", t)
	newtable()
	insertRow()
	client := NewClient()
	ASSERT_CONNECT(client, ADDRESS, true)
	command := fmt.Sprintf("delete from %v", TABLE)
	ASSERT_EXECUTE(client, command, true)
	ASSERT_ACTION(client, "delete")
	ASSERT_ROW_COUNT(client, 1)
	ASSERT_DISCONNECT(client)
}

func TestDeleteManyRows(t *testing.T) {
	register("TestDeleteManyRow", t)
	newtable()
	insertRows()
	client := NewClient()
	ASSERT_CONNECT(client, ADDRESS, true)
	command := fmt.Sprintf("delete from %v ", TABLE)
	ASSERT_EXECUTE(client, command, true)
	ASSERT_ACTION(client, "delete")
	ASSERT_ROW_COUNT(client, ROWS)
	ASSERT_DISCONNECT(client)
}

func TestKey(t *testing.T) {
	register("TestKey", t)
	newtable()
	insertRows()
	client := NewClient()
	ASSERT_CONNECT(client, ADDRESS, true)
	command := fmt.Sprintf("key %v col1", TABLE)
	ASSERT_EXECUTE(client, command, true)
	ASSERT_ACTION(client, "key")
	ASSERT_DISCONNECT(client)
}

func TestTag(t *testing.T) {
	register("TestKey", t)
	newtable()
	insertRows()
	client := NewClient()
	ASSERT_CONNECT(client, ADDRESS, true)
	command := fmt.Sprintf("tag %v col1", TABLE)
	ASSERT_EXECUTE(client, command, true)
	ASSERT_ACTION(client, "tag")
	ASSERT_DISCONNECT(client)
}

func TestSubscribeUnsubscribe(t *testing.T) {
	register("TestSubscribeUnsubscribe", t)
	newtable()
	client := NewClient()
	ASSERT_CONNECT(client, ADDRESS, true)
	command := fmt.Sprintf("subscribe * from %v", TABLE)
	// subscribe
	ASSERT_EXECUTE(client, command, true)
	ASSERT_ACTION(client, "subscribe")
	ASSERT_PUBSUBID(client)
	// unsubscribe
	command = fmt.Sprintf("unsubscribe from %v", TABLE)
	ASSERT_EXECUTE(client, command, true)
	ASSERT_ACTION(client, "unsubscribe")
	//
	ASSERT_DISCONNECT(client)
}

func TestSubscribeUnsubscribeByPubSubId(t *testing.T) {
	register("TestSubscribeUnsubscribeByPubSubId", t)
	newtable()
	client := NewClient()
	ASSERT_CONNECT(client, ADDRESS, true)
	command := fmt.Sprintf("subscribe * from %v", TABLE)
	// subscribe
	ASSERT_EXECUTE(client, command, true)
	ASSERT_ACTION(client, "subscribe")
	ASSERT_PUBSUBID(client)
	// unsubscribe
	command = fmt.Sprintf("unsubscribe from %v where pubsubid = %v", TABLE, client.PubSubId())
	ASSERT_EXECUTE(client, command, true)
	ASSERT_ACTION(client, "unsubscribe")
	//
	ASSERT_DISCONNECT(client)
}

func TestPubSubTimeout(t *testing.T) {
	register("TestPubSubTimeout", t)
	newtable()
	client := NewClient()
	ASSERT_CONNECT(client, ADDRESS, true)
	ASSERT_WAIT_FOR_PUBSUB(client, 10, false)
}

func TestSubscribeSkip(t *testing.T) {
	register("TestSubscribeSkip", t)
	newtable()
	insertRows()
	client := NewClient()
	ASSERT_CONNECT(client, ADDRESS, true)
	command := fmt.Sprintf("subscribe skip * from %v", TABLE)
	ASSERT_EXECUTE(client, command, true)
	ASSERT_ACTION(client, "subscribe")
	ASSERT_PUBSUBID(client)
	ASSERT_WAIT_FOR_PUBSUB(client, 10, false)
	ASSERT_DISCONNECT(client)
}

func TestPubSubAddOnSubscribe(t *testing.T) {
	register("TestPubSubAddOnSubscribe", t)
	newtable()
	insertRows()
	client := NewClient()
	ASSERT_CONNECT(client, ADDRESS, true)
	command := fmt.Sprintf("subscribe * from %v", TABLE)
	// subscribe
	ASSERT_EXECUTE(client, command, true)
	ASSERT_ACTION(client, "subscribe")
	ASSERT_PUBSUBID(client)
	// pubsub add
	pubsubid := client.PubSubId()
	ASSERT_PUBSUB_RESULT_SET(client, pubsubid, "add", ROWS, COLUMNS)
	ASSERT_DISCONNECT(client)
}

func TestPubSubInsert(t *testing.T) {
	register("TestPubSubInsert", t)
	newtable()
	client := NewClient()
	ASSERT_CONNECT(client, ADDRESS, true)
	command := fmt.Sprintf("subscribe * from %v", TABLE)
	// subscribe
	ASSERT_EXECUTE(client, command, true)
	ASSERT_ACTION(client, "subscribe")
	ASSERT_PUBSUBID(client)
	// generate insert event
	insertRows()
	// pubsub insert
	ASSERT_PUBSUB_RESULT_SET(client, client.PubSubId(), "insert", ROWS, COLUMNS)
	ASSERT_DISCONNECT(client)
}

func TestPubSubUpdate(t *testing.T) {
	register("TestPubSubUpdate", t)
	newtable()
	insertRows()
	client := NewClient()
	ASSERT_CONNECT(client, ADDRESS, true)
	command := fmt.Sprintf("subscribe skip * from %v", TABLE)
	// subscribe
	ASSERT_EXECUTE(client, command, true)
	ASSERT_ACTION(client, "subscribe")
	ASSERT_PUBSUBID(client)
	pubsubid := client.PubSubId()
	// generate update event
	command = fmt.Sprintf("update %v set col1 = newvalue", TABLE)
	ASSERT_EXECUTE(client, command, true)
	ASSERT_ROW_COUNT(client, ROWS)
	// expected id and updated column (col1)
	ASSERT_PUBSUB_RESULT_SET(client, pubsubid, "update", ROWS, 2)
	ASSERT_DISCONNECT(client)
}

func TestPubSubDelete(t *testing.T) {
	register("TestPubSubDelete", t)
	newtable()
	insertRows()
	client := NewClient()
	ASSERT_CONNECT(client, ADDRESS, true)
	command := fmt.Sprintf("subscribe skip * from %v", TABLE)
	// subscribe
	ASSERT_EXECUTE(client, command, true)
	ASSERT_ACTION(client, "subscribe")
	ASSERT_PUBSUBID(client)
	pubsubid := client.PubSubId()
	// generate update event
	command = fmt.Sprintf("delete from %v", TABLE)
	ASSERT_EXECUTE(client, command, true)
	ASSERT_ROW_COUNT(client, ROWS)
	// expected id and updated column (col1)
	ASSERT_PUBSUB_RESULT_SET(client, pubsubid, "delete", ROWS, COLUMNS)
	ASSERT_DISCONNECT(client)
}

func TestPubSubRemove(t *testing.T) {
	register("TestPubSubRemove", t)
	newtable()
	insertRows()
	client := NewClient()
	ASSERT_CONNECT(client, ADDRESS, true)
	// key col1
	command := fmt.Sprintf("key %v col1", TABLE)
	ASSERT_EXECUTE(client, command, true)
	command = fmt.Sprintf("subscribe skip * from %v where col1 = 1:col1", TABLE)
	ASSERT_EXECUTE(client, command, true)
	pubsubid := client.PubSubId()
	// generate remove
	command = fmt.Sprintf("update %v set col1 = newvalue where col1 = 1:col1", TABLE)
	ASSERT_EXECUTE(client, command, true)
	ASSERT_PUBSUB_RESULT_SET(client, pubsubid, "remove", 1, COLUMNS)
	ASSERT_DISCONNECT(client)
}

// helpers

var ADDRESS = "localhost:7777"
var T *testing.T = nil
var F = ""
var TABLE = generateTableName()
var ROWS = 300
var COLUMNS = 4

func generateTableName() string {
	return "T" + strconv.FormatInt(time.Now().UnixNano()/100000, 10)
}

func register(f string, t *testing.T) {
	F = f
	T = t
	println(F)
}

func fail(msg string) {
	T.Errorf("%v %v", F, msg)
}

func iferror(client Client, expected bool, got bool) {
	if expected && !got {
		println(fmt.Sprintf("Error: %v", client.Error()))
	}
}

func newtable() {
	TABLE = generateTableName()
}

func insertRow() {
	client := NewClient()
	ASSERT_CONNECT(client, ADDRESS, true)
	command := fmt.Sprintf("insert into %v (col1, col2, col3) values (1:col1, 1:col2, 1:col3)", TABLE)
	ASSERT_EXECUTE(client, command, true)
	ASSERT_DISCONNECT(client)
}

func insertRows() {
	client := NewClient()
	ASSERT_CONNECT(client, ADDRESS, true)
	ASSERT_CONNECT(client, ADDRESS, true)
	for row := 0; row < ROWS; row++ {
		command := fmt.Sprintf("insert into %v (col1, col2, col3) values (%v:col1, %v:col2, %v:col3)", TABLE, row, row, row)
		ASSERT_EXECUTE(client, command, true)
	}
	ASSERT_DISCONNECT(client)
}

func ASSERT_CONNECT(client Client, address string, expected bool) {
	got := client.Connect(address)
	if expected != got {
		fail(fmt.Sprintf("ASSERT_CONNECT failed: expected %v got %v ", expected, got))
	}
}

func VALIDATE_RESULT(client Client, result bool) {
	if result && !client.Ok() {
		fail("VALIDATE_RESULT failed: expected Ok")
	}
	if !result && !client.Failed() {
		fail("VALIDATE_RESULT failed: expected Failed")
	}
}

func ASSERT_DISCONNECT(client Client) {
	client.Disconnect()
	if client.Failed() {
		fail("ASSERT_DISCONNECT failed: expected Ok() not Failed() after Disconnect()")
	}
}

func ASSERT_CONNECTED(client Client, expected bool) {
	got := client.Connected()
	if expected != got {
		fail(fmt.Sprintf("ASSERT_CONNECTED failed: expected %v got %v", expected, got))
	}
}

func ASSERT_EXECUTE(client Client, command string, expected bool) {
	got := client.Execute(command)
	if expected != got {
		fail(fmt.Sprintf("ASSERT_EXECUTE failed: expected %v got %v", expected, got))
	}
	iferror(client, expected, got)
	VALIDATE_RESULT(client, got)
}

func ASSERT_ACTION(client Client, expected string) {
	got := client.Action()
	if expected != got {
		fail(fmt.Sprintf("ASSERT_ACTION failed: expected %v got %v", expected, got))
	}
}

func ASSERT_ROW_COUNT(client Client, expected int) {
	got := client.RowCount()
	if expected != got {
		fail(fmt.Sprintf("ASSERT_ROW_COUNT failed: expected %v but got %v", expected, got))
	}
}

func ASSERT_NEXT_ROW(client Client, expected bool) {
	got := client.NextRow()
	if expected != got {
		fail(fmt.Sprintf("ASSERT_NEXT_ROW failed: expected %v but got %v", expected, got))
	}
}

func ASSERT_ID(client Client) {
	id := client.Value("id")
	if id == "" {
		fail("ASSERT_ID failed: expected non empty string")
	}
}

func ASSERT_VALUE(client Client, column string, value string, match bool) {
	got := client.Value(column)
	if match && value != got {
		fail(fmt.Sprintf("ASSERT_VALUE failed: expected %v but got %v", value, got))
	} else if !match && value != got {
		fail(fmt.Sprintf("ASSERT_VALUE failed: not expected %v", value))
	}
}

func ASSERT_HAS_COLUMN(client Client, column string, expected bool) {
	got := client.HasColumn(column)
	if expected != got {
		fail(fmt.Sprintf("ASSERT_HAS_COLUMN failed: expected %v but got %v", expected, got))
	}
}

func ASSERT_COLUMN_COUNT(client Client, expected int) {
	got := client.ColumnCount()
	if expected != got {
		fail(fmt.Sprintf("ASSERT_COLUMN_COUNT failed: expected %v but got %v", expected, got))
	}
}
func ASSERT_PUBSUBID(client Client) {
	pubsubid := client.PubSubId()
	if pubsubid == "" {
		fail("ASSERT_PUBSUBID failed: expected non empty string")
	}
}

func ASSERT_WAIT_FOR_PUBSUB(client Client, timeout int, expected bool) {
	got := client.WaitForPubSub(timeout)
	if client.Failed() {
		fail(fmt.Sprintf("ASSERT_WAIT_FOR_PUBSUB failed: %v", client.Error()))
	} else if expected != got {
		fail(fmt.Sprintf("ASSERT_WAIT_FOR_PUBSUB failed: expected %v but got %v", expected, got))
	}
}

func ASSERT_PUBSUBID_VALUE(client Client, expected string) {
	got := client.PubSubId()
	if expected != got {
		fail(fmt.Sprintf("ASSERT_PUBSUBID_VALUE failed: expected %v but got %v", expected, got))
	}
}

func ASSERT_NON_EMPTY_VALUE(client Client, ordinal int) {
	if client.ValueByOrdinal(ordinal) == "" {
		fail(fmt.Sprintf("ASSERT_NON_EMPTY_VALUE failed: expected non empty string for ordinal %v", ordinal))
	}
}

func ASSERT_RESULT_SET(client Client, rows int, columns int) {
	ASSERT_ROW_COUNT(client, rows)
	for row := 0; row < rows; row++ {
		ASSERT_NEXT_ROW(client, true)
		ASSERT_COLUMN_COUNT(client, columns)
		for col := 0; col < columns; col++ {
			ASSERT_NON_EMPTY_VALUE(client, col)
		}
	}
	ASSERT_NEXT_ROW(client, false)
}

func ASSERT_PUBSUB_RESULT_SET(client Client, pubsubid string, action string, rows int, columns int) {
	readRows := 0
	for readRows < rows {
		if !client.WaitForPubSub(100) {
			fail(fmt.Sprintf("ASSERT_PUBSUB_RESULT_SET failed expected %v rows but got %v error %v", rows, readRows, client.Error()))
			return
		}
		ASSERT_PUBSUBID_VALUE(client, pubsubid)
		ASSERT_ACTION(client, action)
		for client.NextRow() {
			readRows++
			ASSERT_COLUMN_COUNT(client, columns)
			for col := 0; col < columns; col++ {
				ASSERT_NON_EMPTY_VALUE(client, col)
			}
		}
	}
}
