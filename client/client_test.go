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

import  (
	"testing"
)

var address = "localhost:7777"
var T *testing.T = nil

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
	if !client.Execute(command)	{
		T.Error("Execute failed")
		T.Error(err)	
		ASSERT_FALSE(client.Ok())
		ASSERT_TRUE(client.Failed())
		return
	}
	ASSERT_TRUE(client.Ok())
	ASSERT_FALSE(client.Failed())
}

func ASSERT_CONNECT(client Client) {
	if !client.Connect(address) {
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

func TestConnectDisconnect(t *testing.T) {
	T  = t
	client := NewClient()
	ASSERT_CONNECT(client)	
	client.Disconnect()
}

func TestStatusCommand(t *testing.T) {
	T  = t
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
	ASSERT_ID(client)
	client.Disconnect()
}


