
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

	// simple test framework

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
		register("TestInsert");
		newtable();
		Client client = pubsubsql.Factory.NewClient();
		ASSERT_CONNECT(client, ADDRESS, true);
		String command = String.format("insert into %s (col1, col2, col3) values (1:col1, 1:col2, 1:col3)", TABLE);
		ASSERT_EXECUTE(client, command, true);
		ASSERT_ACTION(client, "insert");
		ASSERT_ROW_COUNT(client, 1);
		ASSERT_NEXT_ROW(client, true);
		ASSERT_NEXT_ROW(client, false);
		ASSERT_DISCONNECT(client);
	}

	private void TestInsertManyRows() {
		register("TestInsert");
		newtable();
		int ROWS = 10;
		Client client = pubsubsql.Factory.NewClient();
		ASSERT_CONNECT(client, ADDRESS, true);
		String command = String.format("insert into %s (col1, col2, col3) values (1:col1, 1:col2, 1:col3)", TABLE);
		for (int r = 0; r < ROWS; r++) {
			ASSERT_EXECUTE(client, command, true);
			ASSERT_ACTION(client, "insert");
			ASSERT_ROW_COUNT(client, 1);
			ASSERT_NEXT_ROW(client, true);
			ASSERT_NEXT_ROW(client, false);
		}
		ASSERT_DISCONNECT(client);
	}
}

