
/* Copyright (C) 2014 CompleteDB LLC.
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

import pubsubsql.Client;
import java.util.*;

public class PubSubSQLTest {

	private int failCount = 0;
	private String currentFunction = "";
	private static final String ADDRESS = "localhost:7777";
	private String TABLE = "T" + System.currentTimeMillis();
	private int ROWS = 10;
	private int COLUMNS = 4; // including id

	//

	public static void main(String[] args) {
		PubSubSQLTest test = new PubSubSQLTest();
		test.TestNetHeader();		
		test.TestClient();
		
		if (test.failCount > 0) {
			System.out.println("Failed " + test.failCount + " tests.");
		} else {
			System.out.println("Passed all tests.");
		}
	}	

	// NetHeader
	private void TestNetHeader() {
		TestWriteRead();	
		TestGetBytes();
	}	

	private void TestWriteRead() {
		register("TestWriteRead");
		pubsubsql.NetHeader header1 = new pubsubsql.NetHeader(32567, 9875235);
		pubsubsql.NetHeader header2 = new pubsubsql.NetHeader(0, 0);
		byte[] bytes = new byte[100];
		header1.WriteTo(bytes);
		header2.ReadFrom(bytes);
		ASSERT_TRUE(header1.MessageSize == header2.MessageSize, "MessageSize do not match");
		ASSERT_TRUE(header1.RequestId == header2.RequestId, "RequestId do not match");
	}

	private void TestGetBytes() {
		register("TestGetBytes");
		pubsubsql.NetHeader header1 = new pubsubsql.NetHeader(32567, 9875235);
		pubsubsql.NetHeader header2 = new pubsubsql.NetHeader(0, 0);
		byte[] bytes = header1.GetBytes();
		header2.ReadFrom(bytes);
		ASSERT_TRUE(header1.MessageSize == header2.MessageSize, "MessageSize do not match");
		ASSERT_TRUE(header1.RequestId == header2.RequestId, "RequestId do not match");
	}

	// Client
	private void TestClient() {
		TestConnectDisconnect();						
		TestExecuteStatus();
		TestExecuteInvalidCommand();
		TestInsertOneRow();
		TestInsertManyRows();
		TestSelectOneRow();
		TestSelectManyRows();
		TestUpdateOneRow();
		TestUpdateManyRows();
		TestDeleteOneRow();
		TestDeleteManyRows();
		TestKey();
		TestTag();
		TestSubscribeUnsubscribe();
		TestSubscribeUnsubscribeByPubSubId();
		TestPubSubTimeout();
		TestSubscribeSkip();
		TestPubSubAddOnSubscribe();
		TestPubSubInsert();
		TestPubSubUpdate();
		TestPubSubDelete();
		TestPubSubRemove();
	}

	private void TestConnectDisconnect() {
		register("TestConnectDisconnect");
		Client client = pubsubsql.Factory.NewClient();
		ASSERT_CONNECT(client, ADDRESS, true);
		ASSERT_DISCONNECT(client);	
		ASSERT_CONNECT(client, "addresswithnoport", false);
		ASSERT_DISCONNECT(client);	
		ASSERT_CONNECT(client, "addresswithnoport:", false);
		ASSERT_DISCONNECT(client);	
		ASSERT_CONNECT(client, "localhost:7778", false);
		ASSERT_DISCONNECT(client);	
		//
		ASSERT_ACTION(client, "");
		ASSERT_ROW_COUNT(client, 0);
		ASSERT_NEXT_ROW(client, false);
	}

	private void TestExecuteStatus() {
		register("TestExecute");
		Client client = pubsubsql.Factory.NewClient();
		ASSERT_CONNECT(client, ADDRESS, true);
		ASSERT_EXECUTE(client, "status", true);
		ASSERT_ACTION(client, "status");
		//
		ASSERT_ROW_COUNT(client, 0);
		ASSERT_NEXT_ROW(client, false);
		ASSERT_DISCONNECT(client);	
	}

	private void TestExecuteInvalidCommand() {
		register("TestExecuteInvalidCommand");
		Client client = pubsubsql.Factory.NewClient();
		ASSERT_CONNECT(client, ADDRESS, true);
		ASSERT_EXECUTE(client, "blablabla", false);
		//
		ASSERT_ACTION(client, "");
		ASSERT_ROW_COUNT(client, 0);
		ASSERT_NEXT_ROW(client, false);
		ASSERT_DISCONNECT(client);	
	}

	private void TestInsertOneRow() {
		register("TestInsertOneRow");
		newtable();
		Client client = pubsubsql.Factory.NewClient();
		ASSERT_CONNECT(client, ADDRESS, true);
		String command = String.format("insert into %s (col1, col2, col3) values (1:col1, 1:col2, 1:col3)", TABLE);
		ASSERT_EXECUTE(client, command, true);
		ASSERT_ACTION(client, "insert");
		ASSERT_ROW_COUNT(client, 1);
		ASSERT_NEXT_ROW(client, true);
		//
		ASSERT_ID(client);
		ASSERT_VALUE(client, "col1", "1:col1", true);
		ASSERT_VALUE(client, "col2", "1:col2", true);
		ASSERT_VALUE(client, "col3", "1:col3", true);
		ASSERT_HAS_COLUMN(client, "col1", true);
		ASSERT_HAS_COLUMN(client, "col2", true);
		ASSERT_HAS_COLUMN(client, "col3", true);
		ASSERT_COLUMN_COUNT(client, 4); // including id
		//
		ASSERT_NEXT_ROW(client, false);
		ASSERT_DISCONNECT(client);
	}

	private void TestInsertManyRows() {
		register("TestInsertManyRows");
		newtable();
		Client client = pubsubsql.Factory.NewClient();
		ASSERT_CONNECT(client, ADDRESS, true);
		String command = String.format("insert into %s (col1, col2, col3) values (1:col1, 1:col2, 1:col3)", TABLE);
		for (int r = 0; r < ROWS; r++) {
			ASSERT_EXECUTE(client, command, true);
			ASSERT_ACTION(client, "insert");
			ASSERT_ROW_COUNT(client, 1);
			ASSERT_NEXT_ROW(client, true);
			//
			ASSERT_ID(client);
			ASSERT_VALUE(client, "col1", "1:col1", true);
			ASSERT_VALUE(client, "col2", "1:col2", true);
			ASSERT_VALUE(client, "col3", "1:col3", true);
			ASSERT_HAS_COLUMN(client, "col1", true);
			ASSERT_HAS_COLUMN(client, "col2", true);
			ASSERT_HAS_COLUMN(client, "col3", true);
			ASSERT_COLUMN_COUNT(client, 4); // including id
			//
			ASSERT_NEXT_ROW(client, false);
		}
		ASSERT_DISCONNECT(client);
	}
	
	private void TestSelectOneRow() {
		register("TestSelectOneRow");
		newtable();
		insertRow();
		Client client = pubsubsql.Factory.NewClient();
		ASSERT_CONNECT(client, ADDRESS, true);
		// select one row
		String command = String.format("select * from %s", TABLE);
		ASSERT_EXECUTE(client, command, true);
		ASSERT_ACTION(client, "select");
		ASSERT_ROW_COUNT(client, 1);
		ASSERT_NEXT_ROW(client, true);
		//
		ASSERT_ID(client);
		ASSERT_VALUE(client, "col1", "1:col1", true);
		ASSERT_VALUE(client, "col2", "1:col2", true);
		ASSERT_VALUE(client, "col3", "1:col3", true);
		ASSERT_HAS_COLUMN(client, "col1", true);
		ASSERT_HAS_COLUMN(client, "col2", true);
		ASSERT_HAS_COLUMN(client, "col3", true);
		ASSERT_COLUMN_COUNT(client, 4); // including id
		//
		ASSERT_NEXT_ROW(client, false);
		ASSERT_DISCONNECT(client);
	}

	private void TestSelectManyRows() {
		register("TestSelectRow");
		newtable();
		insertRows();
		Client client = pubsubsql.Factory.NewClient();
		ASSERT_CONNECT(client, ADDRESS, true);
		String command = String.format("select * from %s", TABLE);
		ASSERT_EXECUTE(client, command, true);
		ASSERT_ACTION(client, "select");
		ASSERT_ROW_COUNT(client, ROWS);
		for (int row = 0; row < ROWS; row++) {
			ASSERT_NEXT_ROW(client, true);
			ASSERT_ID(client);
			ASSERT_VALUE(client, "col1", row + ":col1", true);
			ASSERT_VALUE(client, "col2", row + ":col2", true);
			ASSERT_VALUE(client, "col3", row + ":col3", true);
			ASSERT_HAS_COLUMN(client, "col1", true);
			ASSERT_HAS_COLUMN(client, "col2", true);
			ASSERT_HAS_COLUMN(client, "col3", true);
			ASSERT_COLUMN_COUNT(client, 4); // including id
		}
		//
		ASSERT_NEXT_ROW(client, false);
		ASSERT_DISCONNECT(client);
	}

	private void TestUpdateOneRow() {
		register("TestUpdateOneRow");
		newtable();
		insertRow();
		Client client = pubsubsql.Factory.NewClient();
		ASSERT_CONNECT(client, ADDRESS, true);
		String command = String.format("update %s set col1 = newvalue", TABLE);
		ASSERT_EXECUTE(client, command, true);
		ASSERT_ACTION(client, "update");
		ASSERT_ROW_COUNT(client, 1);
		ASSERT_DISCONNECT(client);
	}

	private void TestUpdateManyRows() {
		register("TestUpdateManyRow");
		newtable();
		insertRows();
		Client client = pubsubsql.Factory.NewClient();
		ASSERT_CONNECT(client, ADDRESS, true);
		String command = String.format("update %s set col1 = newvalue", TABLE);
		ASSERT_EXECUTE(client, command, true);
		ASSERT_ACTION(client, "update");
		ASSERT_ROW_COUNT(client, ROWS);
		ASSERT_DISCONNECT(client);
	}

	private void TestDeleteOneRow() {
		register("TestDeleteOneRow");
		newtable();
		insertRow();
		Client client = pubsubsql.Factory.NewClient();
		ASSERT_CONNECT(client, ADDRESS, true);
		String command = String.format("delete from %s", TABLE);
		ASSERT_EXECUTE(client, command, true);
		ASSERT_ACTION(client, "delete");
		ASSERT_ROW_COUNT(client, 1);
		ASSERT_DISCONNECT(client);
	}

	private void TestDeleteManyRows() {
		register("TestDeleteManyRow");
		newtable();
		insertRows();
		Client client = pubsubsql.Factory.NewClient();
		ASSERT_CONNECT(client, ADDRESS, true);
		String command = String.format("delete from %s ", TABLE);
		ASSERT_EXECUTE(client, command, true);
		ASSERT_ACTION(client, "delete");
		ASSERT_ROW_COUNT(client, ROWS);
		ASSERT_DISCONNECT(client);
	}

	private void TestKey() {
		register("TestKey");
		newtable();
		insertRows();
		Client client = pubsubsql.Factory.NewClient();
		ASSERT_CONNECT(client, ADDRESS, true);
		String command = String.format("key %s col1", TABLE);
		ASSERT_EXECUTE(client, command, true);
		ASSERT_ACTION(client, "key");
		ASSERT_DISCONNECT(client);
	}

	private void TestTag() {
		register("TestTag");
		newtable();
		insertRows();
		Client client = pubsubsql.Factory.NewClient();
		ASSERT_CONNECT(client, ADDRESS, true);
		String command = String.format("tag %s col1", TABLE);
		ASSERT_EXECUTE(client, command, true);
		ASSERT_ACTION(client, "tag");
		ASSERT_DISCONNECT(client);
	}

	private void TestSubscribeUnsubscribe() {
		register("TestSubscribeUnsubscribe");
		newtable();
		Client client = pubsubsql.Factory.NewClient();
		ASSERT_CONNECT(client, ADDRESS, true);
		String command = String.format("subscribe * from %s", TABLE);
		// subscribe
		ASSERT_EXECUTE(client, command, true);
		ASSERT_ACTION(client, "subscribe");
		ASSERT_PUBSUBID(client);
		// unsubscribe
		command = String.format("unsubscribe from %s", TABLE);		
		ASSERT_EXECUTE(client, command, true);
		ASSERT_ACTION(client, "unsubscribe");
		//
		ASSERT_DISCONNECT(client);
	}

	private void TestSubscribeUnsubscribeByPubSubId() {
		register("TestSubscribeUnsubscribeByPubSubId");
		newtable();
		Client client = pubsubsql.Factory.NewClient();
		ASSERT_CONNECT(client, ADDRESS, true);
		String command = String.format("subscribe * from %s", TABLE);
		// subscribe
		ASSERT_EXECUTE(client, command, true);
		ASSERT_ACTION(client, "subscribe");
		ASSERT_PUBSUBID(client);
		// unsubscribe
		command = String.format("unsubscribe from %s where pubsubid = %s", TABLE, client.PubSubId());		
		ASSERT_EXECUTE(client, command, true);
		ASSERT_ACTION(client, "unsubscribe");
		//
		ASSERT_DISCONNECT(client);
	} 

	private void TestPubSubTimeout() {
		register("TestPubSubTimeout");
		newtable();
		Client client = pubsubsql.Factory.NewClient();
		ASSERT_CONNECT(client, ADDRESS, true);
		ASSERT_WAIT_FOR_PUBSUB(client, 10, false);	
	}

	private void TestSubscribeSkip() {
		register("TestSubscribeSkip");
		newtable();
		insertRows();
		Client client = pubsubsql.Factory.NewClient();
		ASSERT_CONNECT(client, ADDRESS, true);
		String command = String.format("subscribe skip * from %s", TABLE);
		ASSERT_EXECUTE(client, command, true);
		ASSERT_ACTION(client, "subscribe");
		ASSERT_PUBSUBID(client);
		ASSERT_WAIT_FOR_PUBSUB(client, 10, false);	
		ASSERT_DISCONNECT(client);
	}


	private void TestPubSubAddOnSubscribe() {
		register("TestPubSubAddOnSubscribe");
		newtable();
		insertRows();
		Client client = pubsubsql.Factory.NewClient();
		ASSERT_CONNECT(client, ADDRESS, true);
		String command = String.format("subscribe * from %s", TABLE);
		// subscribe
		ASSERT_EXECUTE(client, command, true);
		ASSERT_ACTION(client, "subscribe");
		ASSERT_PUBSUBID(client);
		// pubsub add
		String pubsubid = client.PubSubId();
		ASSERT_WAIT_FOR_PUBSUB(client, 10, true);	
		ASSERT_PUBSUBID_VALUE(client, pubsubid);
		ASSERT_ACTION(client, "add");
		ASSERT_RESULT_SET(client, ROWS, COLUMNS); 
		ASSERT_DISCONNECT(client);
	}

	private void TestPubSubInsert() {
		register("TestPubSubInsert");
		newtable();
		Client client = pubsubsql.Factory.NewClient();
		ASSERT_CONNECT(client, ADDRESS, true);
		String command = String.format("subscribe * from %s", TABLE);
		// subscribe
		ASSERT_EXECUTE(client, command, true);
		ASSERT_ACTION(client, "subscribe");
		ASSERT_PUBSUBID(client);
		// generate insert event
		insertRows();
		// pubsub insert
		ASSERT_PUBSUB_RESULT_SET(client, client.PubSubId(), "insert", ROWS, COLUMNS);
		ASSERT_DISCONNECT(client);
	}

	private void TestPubSubUpdate() {
		register("TestPubSubUpdate");
		newtable();
		insertRows();
		Client client = pubsubsql.Factory.NewClient();
		ASSERT_CONNECT(client, ADDRESS, true);
		String command = String.format("subscribe skip * from %s", TABLE);
		// subscribe
		ASSERT_EXECUTE(client, command, true);
		ASSERT_ACTION(client, "subscribe");
		ASSERT_PUBSUBID(client);
		String pubsubid = client.PubSubId();
		// generate update event
		command = String.format("update %s set col1 = newvalue", TABLE);	
		ASSERT_EXECUTE(client, command, true);
		// expected id and updated column (col1)
		ASSERT_PUBSUB_RESULT_SET(client, pubsubid, "update", ROWS, 2);
		ASSERT_DISCONNECT(client);
	}

	private void TestPubSubDelete() {
		register("TestPubSubDelete");
		newtable();
		insertRows();
		Client client = pubsubsql.Factory.NewClient();
		ASSERT_CONNECT(client, ADDRESS, true);
		String command = String.format("subscribe skip * from %s", TABLE);
		// subscribe
		ASSERT_EXECUTE(client, command, true);
		ASSERT_ACTION(client, "subscribe");
		ASSERT_PUBSUBID(client);
		String pubsubid = client.PubSubId();
		// generate update event
		command = String.format("delete from %s", TABLE);	
		ASSERT_EXECUTE(client, command, true);
		// expected id and updated column (col1)
		ASSERT_PUBSUB_RESULT_SET(client, pubsubid, "delete", ROWS, COLUMNS);
		ASSERT_DISCONNECT(client);
	}

	private void TestPubSubRemove() {
		register("TestPubSubRemove");
	}

	// helper functions

	private String generateTableName() {
		return "T" + System.currentTimeMillis();
	}

	private void newtable() {
		TABLE = generateTableName();
	}

	private void fail(String message) {
		System.out.println(currentFunction + " " + message);
		failCount++;	
	}

	private void print(String message) {
		System.out.println(message);
	}

	private void iferror(Client client, boolean expected, boolean got) {
		if (expected && !got) {
			print(String.format("Error: %s", client.Error()));
		}	
	}

	private void register(String function) {
		currentFunction = function;
	}

	private void insertRow() {
		Client client = pubsubsql.Factory.NewClient();
		ASSERT_CONNECT(client, ADDRESS, true);
		String command = String.format("insert into %s (col1, col2, col3) values (1:col1, 1:col2, 1:col3)", TABLE);
		ASSERT_EXECUTE(client, command, true);
		ASSERT_DISCONNECT(client);	
	}

	private void insertRows() {
		Client client = pubsubsql.Factory.NewClient();
		ASSERT_CONNECT(client, ADDRESS, true);
		for (int row = 0; row < ROWS; row++) {	
			String command = String.format("insert into %s (col1, col2, col3) values (%s:col1, %s:col2, %s:col3)", TABLE, row, row, row);
			ASSERT_EXECUTE(client, command, true);
		}
		ASSERT_DISCONNECT(client);
	}

	private void key(String column) {
		Client client = pubsubsql.Factory.NewClient();
		ASSERT_CONNECT(client, ADDRESS, true);
		String command = String.format("key %s %s", TABLE, column);
		ASSERT_EXECUTE(client, command, true);
		ASSERT_DISCONNECT(client);
	}

	private void tag(String column) {
		Client client = pubsubsql.Factory.NewClient();
		ASSERT_CONNECT(client, ADDRESS, true);
		String command = String.format("tag %s %s", TABLE, column);
		ASSERT_EXECUTE(client, command, true);
		ASSERT_DISCONNECT(client);
	}
	
	public void ASSERT_TRUE(boolean val, String message) {
		if (!val) {
			fail("ASSERT_TRUE failed: " + message);
		}	
	}

	public void ASSERT_FALSE(boolean val, String message) {
		if (val) {
			fail("ASSERT_FALSE failed: " + message);
		}
	}

	public void VALIDATE_RESULT(Client client, boolean result) {
		if (result && !client.Ok() ) fail("VALIDATE_RESULT failed: expected Ok");
		if (!result && !client.Failed() ) fail("VALIDATE_RESULT failed: expected Failed");
	}

	public void ASSERT_CONNECT(Client client, String address, boolean expected) {
		boolean got = client.Connect(address);
		if (expected != got) {
			fail(String.format("ASSERT_CONNECT failed: expected %s got %s ", expected, got));
		}	
		iferror(client, expected, got);
		VALIDATE_RESULT(client, got);
	}

	public void ASSERT_DISCONNECT(Client client) {
		client.Disconnect();
		if (client.Failed()) {
			fail("ASSERT_DISCONNECT failed: expected Ok() not Failed() after Disconnect()");
		}
	}

	public void ASSERT_CONNECTED(Client client, boolean expected) {
		boolean got = client.Connected();
		if (expected != got) {
			fail(String.format("ASSERT_CONNECTED failed: expected %s got %s", expected, got));
		}
	}

	public void ASSERT_EXECUTE(Client client, String command, boolean expected) {
		boolean got = client.Execute(command);
		if (expected != got) {
			fail(String.format("ASSERT_EXECUTE failed: expected %s got %s", expected, got));	
		}
		iferror(client, expected, got);
		VALIDATE_RESULT(client, got);
	}

	public void ASSERT_ACTION(Client client, String expected) {
		String got = client.Action();
		if (!expected.equals(got)) {
			fail(String.format("ASSERT_ACTION failed: expected %s got %s", expected, got));
		}
	}

	public void ASSERT_ROW_COUNT(Client client, int expected) {
		int got = client.RowCount();
		if (expected != got) {
			fail(String.format("ASSERT_ROW_COUNT failed: expected %s but got %s", expected, got));
		}
	}

	public void ASSERT_NEXT_ROW(Client client, boolean expected) {
		boolean got = client.NextRow();
		if (expected != got) {
			fail(String.format("ASSERT_NEXT_ROW failed: expected %s but got %s", expected, got));
		}
	}

	public void ASSERT_ID(Client client) {
		String id = client.Value("id");
		if (id.length() == 0) {
			fail("ASSERT_ID failed: expected non empty string");
		}
	}

	public void ASSERT_PUBSUBID(Client client) {
		String pubsubid = client.PubSubId();	
		if (pubsubid.length() == 0) {
			fail("ASSERT_PUBSUBID failed: expected non empty string");
		}	
	}

	public void ASSERT_PUBSUBID_VALUE(Client client, String expected) {
		String got = client.PubSubId(); 		
		if (!expected.equals(got)) {
			fail(String.format("ASSERT_PUBSUBID_VALUE failed: expected %s but got %s", expected, got));
		}
	}

	public void ASSERT_VALUE(Client client, String column, String value, boolean match) {
		String got = client.Value(column);	
		if (match && !value.equals(got)) {
			fail(String.format("ASSERT_VALUE failed: expected %s but got %s", value, got));
		}
		else if (!match && value.equals(got)) {
			fail(String.format("ASSERT_VALUE failed: not expected %s", value));
		}
	}

	public void ASSERT_COLUMN_COUNT(Client client, int expected) {
		int got = client.ColumnCount();
		if (expected != got) {
			fail(String.format("ASSERT_COLUMN_COUNT failed: expected %s but got %s", expected, got));
		}
	}

	public void ASSERT_HAS_COLUMN(Client client, String column, boolean expected) {
		boolean got = client.HasColumn(column);
		if (expected != got) {
			fail(String.format("ASSERT_HAS_COLUMN failed: expected %s but got %s", expected, got));
		}
	}	

	public void ASSERT_WAIT_FOR_PUBSUB(Client client, int timeout, boolean expected) {
		boolean got = client.WaitForPubSub(timeout);
		if (client.Failed()) {
			fail(String.format("ASSERT_WAIT_FOR_PUBSUB failed: %s", client.Error()));
		}
		else if (expected != got) {
			fail(String.format("ASSERT_WAIT_FOR_PUBSUB failed: expected %s but got %s", expected, got));
		}	
	}

	public void ASSERT_NON_EMPTY_VALUE(Client client, int ordinal) {
		if (client.ValueByOrdinal(ordinal).length() == 0) {
			fail(String.format("ASSERT_NON_EMPTY_VALUE failed: expected non empty string for ordinal %s", ordinal));
		}
	}

	public void ASSERT_RESULT_SET(Client client, int rows, int columns) {
		ASSERT_ROW_COUNT(client, rows);
		for (int row = 0; row < rows; row++) {
			ASSERT_NEXT_ROW(client, true);
			ASSERT_COLUMN_COUNT(client, columns); 
			for (int col = 0; col < columns; col++) {
				ASSERT_NON_EMPTY_VALUE(client, col);
			} 
		}
		ASSERT_NEXT_ROW(client, false);
	} 

	public void ASSERT_PUBSUB_RESULT_SET(Client client, String pubsubid, String action, int rows, int columns) {
		int readRows = 0;		
		while (readRows < rows) {
			if (!client.WaitForPubSub(100)) {
				fail(String.format("ASSERT_PUBSUB_RESULT_SET failed expected %s rows but got %s", rows, readRows));
				return;
			}
			ASSERT_PUBSUBID_VALUE(client, pubsubid);
			ASSERT_ACTION(client, action);
			while (client.NextRow()) {
				readRows++;
				ASSERT_COLUMN_COUNT(client, columns); 
				for (int col = 0; col < columns; col++) {
					ASSERT_NON_EMPTY_VALUE(client, col);
				} 
			}
		}
	}
}

