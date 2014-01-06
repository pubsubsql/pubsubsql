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
	"strconv"
	"testing"
	"time"
	"fmt"
)

func generateTableName() string {
	return "T" + strconv.FormatInt(time.Now().Unix(), 10)
}

var ADDRESS = "localhost:7777"
var T *testing.T = nil
var TABLE = generateTableName()
var ROWS = 300 

func ASSERT_TRUE(b bool) {
	if !b {
		T.Error("Expected true")
	}
}

func ASSERT_FALSE(b bool) {
	if b {
		T.Error("Expected false")
	}
}

func ASSERT_INT_EQ(val1 int, val2 int, err string) {
	if val1 != val2 {
		T.Error("int values do not match val1:", val1, "val2:", val2)
		T.Error(err)
	}
}

func ASSERT_EXECUTE(client Client, command string, err string) {
	if !client.Execute(command) {
		T.Error("Execute failed")
		T.Error(client.Error())
		T.Error(err)
		ASSERT_FALSE(client.Ok())
		ASSERT_TRUE(client.Failed())
		return
	}
	ASSERT_TRUE(client.Ok())
	ASSERT_FALSE(client.Failed())
}

func ASSERT_CONNECT(client Client) {
	if !client.Connect(ADDRESS) {
		T.Error("Connect failed.", client.Error())
		ASSERT_FALSE(client.Ok())
		ASSERT_TRUE(client.Failed())
		return
	}
	ASSERT_TRUE(client.Ok())
	ASSERT_FALSE(client.Failed())
}

func ASSERT_ACTION(client Client, action string) {
	if client.Action() != action {
		T.Error("Expected action", action, "but got", client.Action())
	}
}

func ASSERT_ID(client Client) {
	if client.Id() == "" {
		T.Error("Expected id but got empty string")
	}
}

func ASSERT_NOID(client Client) {
	if client.Id() != "" {
		T.Error("Expected no id but got", client.Id())
	}
}

func ASSERT_PUBSUBID(client Client) {
	if client.PubSubId() == "" {
		T.Error("Expected pubsubid but got empty string")
	}
}

func ASSERT_NOPUBSUBID(client Client) {
	if client.PubSubId() != "" {
		T.Error("Expected no pubsubid but got", client.PubSubId())
	}
}

func ASSERT_RECORD_COUNT(client Client, count int) {
	if client.RecordCount() != count {
		T.Error("Expected record count ", count, "but got", client.RecordCount())
	}
}

func ASSERT_VALUE(client Client, column string, value string) {
	if client.Value(column) != value {
		T.Error("Expected value ", value, "but got", client.Value(column))
	}
}

func ASSERT_COLUMN(client Client, column string) {
	if !client.HasColumn(column) {
		T.Error("Expexted column ", column)
	}
}

func TestConnectDisconnect(t *testing.T) {
	T = t
	client := NewClient()
	ASSERT_CONNECT(client)
	client.Disconnect()
}

func TestStatusCommand(t *testing.T) {
	T = t
	client := NewClient()
	ASSERT_CONNECT(client)
	ASSERT_EXECUTE(client, "status", "status failed")
	ASSERT_ACTION(client, "status")
	ASSERT_NOID(client)
	ASSERT_NOPUBSUBID(client)
	client.Disconnect()
}

func TestInsertCommand(t *testing.T) {
	T = t
	client := NewClient()
	ASSERT_CONNECT(client)
	ASSERT_EXECUTE(client, "insert into insertcommand (col1, col2) values ('HELLO', WORLD)", "insert failed")
	ASSERT_ACTION(client, "insert")
	ASSERT_ID(client)
	ASSERT_NOPUBSUBID(client)
	client.Disconnect()
}

func TestSelectCommand(t *testing.T) {
	println(TABLE)
	T = t
	client := NewClient()
	ASSERT_CONNECT(client)
	command := "insert into " + TABLE + " (col1, col2, col3) values (col1, col2, col3) "
	for i := 0; i < ROWS; i++ {
		ASSERT_EXECUTE(client, command, "insert failed "+command)
		ASSERT_ACTION(client, "insert")
		ASSERT_ID(client)
		ASSERT_NOPUBSUBID(client)
	}
	//
	command = "select * from " + TABLE
	ASSERT_EXECUTE(client, command, "select failed")
	ASSERT_ACTION(client, "select")
	ASSERT_NOID(client)
	ASSERT_RECORD_COUNT(client, ROWS)
	rowsread := 0
	for client.NextRecord() {
		rowsread++
		ASSERT_RECORD_COUNT(client, ROWS)
	}
	println("rowsread:", rowsread)
	ASSERT_INT_EQ(ROWS, rowsread, "NextRecord failed")
	ASSERT_NOPUBSUBID(client)
	client.Disconnect()
}

func TestUpdateCommand(t *testing.T) {
	T = t
	client := NewClient()
	ASSERT_CONNECT(client)
	command := "update " + TABLE + " set col1 = updated_value "
	ASSERT_EXECUTE(client, command, "update failed")
	ASSERT_ACTION(client, "update")
	ASSERT_NOID(client)
	ASSERT_RECORD_COUNT(client, ROWS)
	ASSERT_NOPUBSUBID(client)
	client.Disconnect()
}

func TestDeleteCommand(t *testing.T) {
	T = t
	client := NewClient()
	ASSERT_CONNECT(client)
	command := "delete from " + TABLE
	ASSERT_EXECUTE(client, command, "delete failed")
	ASSERT_ACTION(client, "delete")
	ASSERT_NOID(client)
	ASSERT_RECORD_COUNT(client, ROWS)
	ASSERT_NOPUBSUBID(client)
	client.Disconnect()
}

func TestKeyCommand(t *testing.T) {
	T = t
	client := NewClient()
	ASSERT_CONNECT(client)
	command := "key " + TABLE + " col1"
	ASSERT_EXECUTE(client, command, "key failed")
	ASSERT_ACTION(client, "key")
	ASSERT_NOID(client)
	ASSERT_RECORD_COUNT(client, 0)
	ASSERT_NOPUBSUBID(client)
	client.Disconnect()
}

func TestTagCommand(t *testing.T) {
	T = t
	client := NewClient()
	ASSERT_CONNECT(client)
	command := "tag " + TABLE + " col2"
	ASSERT_EXECUTE(client, command, "tag failed")
	ASSERT_ACTION(client, "tag")
	ASSERT_NOID(client)
	ASSERT_RECORD_COUNT(client, 0)
	ASSERT_NOPUBSUBID(client)
	client.Disconnect()
}

func TestSubscribeUnsubscribeCommand(t *testing.T) {
	T = t
	client := NewClient()
	ASSERT_CONNECT(client)
	// subscribe
	command := "subscribe skip * from " + TABLE 
	ASSERT_EXECUTE(client, command, "subscribe failed")
	ASSERT_ACTION(client, "subscribe")
	ASSERT_NOID(client)
	ASSERT_RECORD_COUNT(client, 0)
	ASSERT_PUBSUBID(client)
	// unsubscribe
	command = "unsubscribe from " + TABLE + " where pubsubid = " + client.PubSubId()
	ASSERT_EXECUTE(client, command, "subscribe failed")
	ASSERT_ACTION(client, "unsubscribe")
	ASSERT_NOID(client)
	ASSERT_RECORD_COUNT(client, 0)
	ASSERT_NOPUBSUBID(client)
	//
	client.Disconnect()
}

func TestValueAndColumns(t *testing.T) {
	T = t
	client := NewClient()
	ASSERT_CONNECT(client)
	// subscribe
	command := "delete from " + TABLE 
	// clear the table first
	ASSERT_EXECUTE(client, command, "delete failed")
	// insert values
	for i := 0; i < ROWS; i++ {
		val1 := "1:" + strconv.Itoa(i)
		val2 := "2:" + strconv.Itoa(i)
		val3 := "3:" + strconv.Itoa(i)
		command := fmt.Sprintf("insert into %s (col1, col2, col3) values (%s, %s, %s)", TABLE, val1, val2, val3)
		ASSERT_EXECUTE(client, command, "insert failed")
	}	
	// test value	
	command = "select * from " + TABLE	
	ASSERT_EXECUTE(client, command, "select failed")
	i := 0
	for client.NextRecord() {
		val1 := "1:" + strconv.Itoa(i)
		ASSERT_COLUMN(client, "col1")
		ASSERT_VALUE(client, "col1", val1)
		val2 := "2:" + strconv.Itoa(i)
		ASSERT_COLUMN(client, "col2")
		ASSERT_VALUE(client, "col2", val2)
		val3 := "3:" + strconv.Itoa(i)
		ASSERT_COLUMN(client, "col3")
		ASSERT_VALUE(client, "col3", val3)
		ASSERT_VALUE(client, "invalid_column", "")
				
		i++
	}
	ASSERT_INT_EQ(len(client.Columns()), 4, "Columns failed")
	client.Disconnect()
	
}
