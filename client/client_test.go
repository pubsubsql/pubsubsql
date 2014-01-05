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
)

var ADDRESS = "localhost:7777"
var T *testing.T = nil
var TABLE = "T" + strconv.FormatInt(time.Now().Unix(), 10)

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

func ASSERT_RECORD_COUNT(client Client, count int) {
	if client.RecordCount() != count {
		T.Error("Expected record count ", count, "but got", client.RecordCount())
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
	client.Disconnect()
}

func TestInsertCommand(t *testing.T) {
	T = t
	client := NewClient()
	ASSERT_CONNECT(client)
	ASSERT_EXECUTE(client, "insert into insertcommand (col1, col2) values ('HELLO', WORLD)", "insert failed")
	ASSERT_ACTION(client, "insert")
	ASSERT_ID(client)
	client.Disconnect()
}

func TestSelectCommand(t *testing.T) {
	println(TABLE)
	T = t
	client := NewClient()
	ASSERT_CONNECT(client)
	rows := 250
	command := "insert into " + TABLE + " (col1, col2, col3) values (col1, col2, col3) "
	for i := 0; i < rows; i++ {
		ASSERT_EXECUTE(client, command, "insert failed "+command)
		ASSERT_ACTION(client, "insert")
		ASSERT_ID(client)
	}
	//
	command = "select * from " + TABLE
	ASSERT_EXECUTE(client, command, "select failed")
	ASSERT_ACTION(client, "select")
	ASSERT_NOID(client)
	ASSERT_RECORD_COUNT(client, rows)
	client.Disconnect()
}

func TestUpdateCommand(t *testing.T) {
	T = t
	client := NewClient()
	ASSERT_CONNECT(client)
	rows := 250
	// previous test inserted 250 records
	command := "update " + TABLE + " set col1 = updated_value "
	ASSERT_EXECUTE(client, command, "update failed")
	ASSERT_ACTION(client, "update")
	ASSERT_NOID(client)
	ASSERT_RECORD_COUNT(client, rows)
	client.Disconnect()
}

func TestDeleteCommand(t *testing.T) {
	T = t
	client := NewClient()
	ASSERT_CONNECT(client)
	rows := 250
	// previous test inserted 250 records
	command := "delete from " + TABLE
	ASSERT_EXECUTE(client, command, "delete failed")
	ASSERT_ACTION(client, "delete")
	ASSERT_NOID(client)
	ASSERT_RECORD_COUNT(client, rows)
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
	client.Disconnect()
}
