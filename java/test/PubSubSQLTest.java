
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

public class PubSubSQLTest {

	private int failCount = 0;
	private String currentFunction = "";
	private static final String ADDRESS = "localhost:7777";

	// simple test framework
	private void fail(String message) {
		System.out.println(currentFunction + " " + message);
		failCount++;	
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

	public void ASSERT_CONNECT(Client client) {
		if (!client.Connect(ADDRESS)) {
			fail("ASSERT_CONNECT failed " + client.Error());
		}	
		ASSERT_TRUE(client.Ok(), "client.Ok");
		ASSERT_FALSE(client.Failed(), "client.Failed");
	}

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
	}

	private void TestConnectDisconnect() {
		register("TestConnectDisconnect");
		Client client = pubsubsql.Factory.NewClient();
		ASSERT_CONNECT(client);
		client.Disconnect();	
		ASSERT_FALSE(client.Connect("addresswithnoport"), "address with no port");
		ASSERT_FALSE(client.Connect("addresswithnoport:"), "address with separator no port");
		ASSERT_FALSE(client.Connect("localhost:7778"), "invalid address");
	}

}
