
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


public class PubSubSQLTest {

	private int failCount = 0;

	// simple test framework
	// ASSERT
	public void ASSERT_TRUE(boolean val, String message) {
		if (!val) {
			System.out.println(message);
			failCount++;
		}	
	}

	public static void main(String[] args) {
		PubSubSQLTest test = new PubSubSQLTest();
		test.TestNetHeader();		
		
		if (test.failCount > 0) {
			System.out.println("Failed " + test.failCount + " tests.");
		} else {
			System.out.println("Passed all tests.");
		}
	}	

	// NetHeader
	private void TestNetHeader() {
		TestWriteRead();	
	}	

	private void TestWriteRead() {
		pubsubsql.NetHeader header1 = new pubsubsql.NetHeader(32567, 9875235);
		pubsubsql.NetHeader header2 = new pubsubsql.NetHeader(0, 0);
		byte[] bytes = new byte[100];
		header1.WriteTo(bytes);
		header2.ReadFrom(bytes);
		ASSERT_TRUE(header1.MessageSize == header2.MessageSize, "TestWriteRead: MessageSize do not match");
		ASSERT_TRUE(header1.RequestId == header2.RequestId, "TestWriteRead: MessageSize do not match");
	}

	private void TestGetBytes() {
		pubsubsql.NetHeader header1 = new pubsubsql.NetHeader(32567, 9875235);
		pubsubsql.NetHeader header2 = new pubsubsql.NetHeader(0, 0);
		byte[] bytes = header1.GetBytes();
		header2.ReadFrom(bytes);
		ASSERT_TRUE(header1.MessageSize == header2.MessageSize, "TestWriteRead: MessageSize do not match");
		ASSERT_TRUE(header1.RequestId == header2.RequestId, "TestWriteRead: MessageSize do not match");
	}

	//
	

}
