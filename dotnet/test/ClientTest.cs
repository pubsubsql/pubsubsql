using System;
using Microsoft.VisualStudio.TestTools.UnitTesting;
using PubSubSQL;

namespace PubSubSQLTest
{
     [TestClass]
     public class ClientTest {

        private String ADDRESS = "localhost:7777";
        private String TABLE = "T" + DateTime.Now.Ticks; 
        private int ROWS = 300;
        private int COLUMNS = 4; // including id

        [TestMethod]
        public void TestConnectDisconnect() {
            Client client = PubSubSQL.Factory.NewClient();
            ASSERT_CONNECT(client, ADDRESS, true);
            ASSERT_CONNECTED(client, true);
            ASSERT_DISCONNECT(client);	
            ASSERT_CONNECT(client, "addresswithnoport", false);
            ASSERT_CONNECTED(client, false);
            ASSERT_DISCONNECT(client);	
            ASSERT_CONNECT(client, "addresswithnoport:", false);
            ASSERT_CONNECTED(client, false);
            ASSERT_DISCONNECT(client);	
            ASSERT_CONNECT(client, "localhost:7778", false);
            ASSERT_CONNECTED(client, false);
            ASSERT_DISCONNECT(client);	
        }

        [TestMethod]
        public void TestExecuteStatus() {
            Client client = PubSubSQL.Factory.NewClient();
            ASSERT_CONNECT(client, ADDRESS, true);
            ASSERT_EXECUTE(client, "status", true);
            ASSERT_ACTION(client, "status");
            ASSERT_DISCONNECT(client);	
        }

        [TestMethod]
        public void TestExecuteInvalidCommand() {
            Client client = PubSubSQL.Factory.NewClient();
            ASSERT_CONNECT(client, ADDRESS, true);
            ASSERT_EXECUTE(client, "blablabla", false);
        }

        [TestMethod]
        public void TestInsertOneRow() {
            newtable();
            Client client = PubSubSQL.Factory.NewClient();
            ASSERT_CONNECT(client, ADDRESS, true);
            String command = string.Format("insert into {0} (col1, col2, col3) values (1:col1, 1:col2, 1:col3)", TABLE);
            ASSERT_EXECUTE(client, command, true);
            ASSERT_ACTION(client, "insert");
            ASSERT_ROW_COUNT(client, 1);
            ASSERT_NEXT_ROW(client, true);
            ASSERT_ID(client);
            ASSERT_VALUE(client, "col1", "1:col1", true);
            ASSERT_VALUE(client, "col2", "1:col2", true);
            ASSERT_VALUE(client, "col3", "1:col3", true);
            ASSERT_HAS_COLUMN(client, "col1", true);
            ASSERT_HAS_COLUMN(client, "col2", true);
            ASSERT_HAS_COLUMN(client, "col3", true);
            ASSERT_COLUMN_COUNT(client, 4); // including id
            ASSERT_NEXT_ROW(client, false);
            ASSERT_DISCONNECT(client);
        }

        [TestMethod]
        public void TestInsertManyRows() {
            newtable();
            Client client = PubSubSQL.Factory.NewClient();
            ASSERT_CONNECT(client, ADDRESS, true);
            String command = string.Format("insert into {0} (col1, col2, col3) values (1:col1, 1:col2, 1:col3)", TABLE);
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
        
        [TestMethod]
        public void TestSelectOneRow() {
            newtable();
            insertRow();
            Client client = PubSubSQL.Factory.NewClient();
            ASSERT_CONNECT(client, ADDRESS, true);
            // select one row
            String command = string.Format("select * from {0}", TABLE);
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

        [TestMethod]
        public void TestSelectManyRows() {
            newtable();
            insertRows();
            Client client = PubSubSQL.Factory.NewClient();
            ASSERT_CONNECT(client, ADDRESS, true);
            String command = string.Format("select * from {0}", TABLE);
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

        [TestMethod]
        public void TestUpdateOneRow() {
            newtable();
            insertRow();
            Client client = PubSubSQL.Factory.NewClient();
            ASSERT_CONNECT(client, ADDRESS, true);
            String command = string.Format("update {0} set col1 = newvalue", TABLE);
            ASSERT_EXECUTE(client, command, true);
            ASSERT_ACTION(client, "update");
            ASSERT_ROW_COUNT(client, 1);
            ASSERT_DISCONNECT(client);
        }

        [TestMethod]
        public void TestUpdateManyRows() {
            newtable();
            insertRows();
            Client client = PubSubSQL.Factory.NewClient();
            ASSERT_CONNECT(client, ADDRESS, true);
            String command = string.Format("update {0} set col1 = newvalue", TABLE);
            ASSERT_EXECUTE(client, command, true);
            ASSERT_ACTION(client, "update");
            ASSERT_ROW_COUNT(client, ROWS);
            ASSERT_DISCONNECT(client);
        }

        [TestMethod]
        public void TestDeleteOneRow() {
            newtable();
            insertRow();
            Client client = PubSubSQL.Factory.NewClient();
            ASSERT_CONNECT(client, ADDRESS, true);
            String command = string.Format("delete from {0}", TABLE);
            ASSERT_EXECUTE(client, command, true);
            ASSERT_ACTION(client, "delete");
            ASSERT_ROW_COUNT(client, 1);
            ASSERT_DISCONNECT(client);
        }

        [TestMethod]
        public void TestDeleteManyRows() {
            newtable();
            insertRows();
            Client client = PubSubSQL.Factory.NewClient();
            ASSERT_CONNECT(client, ADDRESS, true);
            String command = string.Format("delete from {0} ", TABLE);
            ASSERT_EXECUTE(client, command, true);
            ASSERT_ACTION(client, "delete");
            ASSERT_ROW_COUNT(client, ROWS);
            ASSERT_DISCONNECT(client);
        }

        [TestMethod]
        public void TestKey() {
            newtable();
            insertRows();
            Client client = PubSubSQL.Factory.NewClient();
            ASSERT_CONNECT(client, ADDRESS, true);
            String command = string.Format("key {0} col1", TABLE);
            ASSERT_EXECUTE(client, command, true);
            ASSERT_ACTION(client, "key");
            ASSERT_DISCONNECT(client);
        }

        [TestMethod]
        public void TestTag() {
            newtable();
            insertRows();
            Client client = PubSubSQL.Factory.NewClient();
            ASSERT_CONNECT(client, ADDRESS, true);
            String command = string.Format("tag {0} col1", TABLE);
            ASSERT_EXECUTE(client, command, true);
            ASSERT_ACTION(client, "tag");
            ASSERT_DISCONNECT(client);
        }

        [TestMethod]
        public void TestSubscribeUnsubscribe() {
            newtable();
            Client client = PubSubSQL.Factory.NewClient();
            ASSERT_CONNECT(client, ADDRESS, true);
            String command = string.Format("subscribe * from {0}", TABLE);
            // subscribe
            ASSERT_EXECUTE(client, command, true);
            ASSERT_ACTION(client, "subscribe");
            ASSERT_PUBSUBID(client);
            // unsubscribe
            command = string.Format("unsubscribe from {0}", TABLE);		
            ASSERT_EXECUTE(client, command, true);
            ASSERT_ACTION(client, "unsubscribe");
            //
            ASSERT_DISCONNECT(client);
        }

        [TestMethod]
        public void TestSubscribeUnsubscribeByPubSubId() {
            newtable();
            Client client = PubSubSQL.Factory.NewClient();
            ASSERT_CONNECT(client, ADDRESS, true);
            String command = string.Format("subscribe * from {0}", TABLE);
            // subscribe
            ASSERT_EXECUTE(client, command, true);
            ASSERT_ACTION(client, "subscribe");
            ASSERT_PUBSUBID(client);
            // unsubscribe
            command = string.Format("unsubscribe from {0} where pubsubid = {1}", TABLE, client.PubSubId());		
            ASSERT_EXECUTE(client, command, true);
            ASSERT_ACTION(client, "unsubscribe");
            //
            ASSERT_DISCONNECT(client);
        } 

        [TestMethod]
        public void TestPubSubTimeout() {
            newtable();
            Client client = PubSubSQL.Factory.NewClient();
            ASSERT_CONNECT(client, ADDRESS, true);
            ASSERT_WAIT_FOR_PUBSUB(client, 10, false);	
        }

        [TestMethod]
        public void TestSubscribeSkip() {
            newtable();
            insertRows();
            Client client = PubSubSQL.Factory.NewClient();
            ASSERT_CONNECT(client, ADDRESS, true);
            String command = string.Format("subscribe skip * from {0}", TABLE);
            ASSERT_EXECUTE(client, command, true);
            ASSERT_ACTION(client, "subscribe");
            ASSERT_PUBSUBID(client);
            ASSERT_WAIT_FOR_PUBSUB(client, 10, false);	
            ASSERT_DISCONNECT(client);
        }

        [TestMethod]
        public void TestPubSubAddOnSubscribe() {
            newtable();
            insertRows();
            Client client = PubSubSQL.Factory.NewClient();
            ASSERT_CONNECT(client, ADDRESS, true);
            String command = string.Format("subscribe * from {0}", TABLE);
            // subscribe
            ASSERT_EXECUTE(client, command, true);
            ASSERT_ACTION(client, "subscribe");
            ASSERT_PUBSUBID(client);
            // pubsub add
            String pubsubid = client.PubSubId();
            ASSERT_PUBSUB_RESULT_SET(client, pubsubid, "add", ROWS, COLUMNS);
            ASSERT_DISCONNECT(client);
        }

        [TestMethod]
        public void TestPubSubInsert() {
            newtable();
            Client client = PubSubSQL.Factory.NewClient();
            ASSERT_CONNECT(client, ADDRESS, true);
            String command = string.Format("subscribe * from {0}", TABLE);
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

        [TestMethod]
        public void TestPubSubUpdate() {
            newtable();
            insertRows();
            Client client = PubSubSQL.Factory.NewClient();
            ASSERT_CONNECT(client, ADDRESS, true);
            String command = string.Format("subscribe skip * from {0}", TABLE);
            // subscribe
            ASSERT_EXECUTE(client, command, true);
            ASSERT_ACTION(client, "subscribe");
            ASSERT_PUBSUBID(client);
            String pubsubid = client.PubSubId();
            // generate update event
            command = string.Format("update {0} set col1 = newvalue", TABLE);	
            ASSERT_EXECUTE(client, command, true);
            ASSERT_ROW_COUNT(client, ROWS);
            // expected id and updated column (col1)
            ASSERT_PUBSUB_RESULT_SET(client, pubsubid, "update", ROWS, 2);
            ASSERT_DISCONNECT(client);
        }

        [TestMethod]
        public void TestPubSubDelete() {
            newtable();
            insertRows();
            Client client = PubSubSQL.Factory.NewClient();
            ASSERT_CONNECT(client, ADDRESS, true);
            String command = string.Format("subscribe skip * from {0}", TABLE);
            // subscribe
            ASSERT_EXECUTE(client, command, true);
            ASSERT_ACTION(client, "subscribe");
            ASSERT_PUBSUBID(client);
            String pubsubid = client.PubSubId();
            // generate delete event
            command = string.Format("delete from {0}", TABLE);	
            ASSERT_EXECUTE(client, command, true);
            ASSERT_ROW_COUNT(client, ROWS);
            // expected id and updated column (col1)
            ASSERT_PUBSUB_RESULT_SET(client, pubsubid, "delete", ROWS, COLUMNS);
            ASSERT_DISCONNECT(client);
        }

        [TestMethod]
        public void TestPubSubRemove() {
            newtable();
            insertRows();
            Client client = PubSubSQL.Factory.NewClient();
            ASSERT_CONNECT(client, ADDRESS, true);
            // key col1
            String command = string.Format("key {0} col1", TABLE);
            ASSERT_EXECUTE(client, command, true);
            command = string.Format("subscribe skip * from {0} where col1 = 1:col1", TABLE);
            ASSERT_EXECUTE(client, command, true);
            String pubsubid = client.PubSubId();
            // generate remove event
            command = string.Format("update {0} set col1 = newvalue where col1 = 1:col1", TABLE);	
            ASSERT_EXECUTE(client, command, true);
            ASSERT_PUBSUB_RESULT_SET(client, pubsubid, "remove", 1, COLUMNS);
            ASSERT_DISCONNECT(client);
        }

        // helper functions

        public String generateTableName() {
            return "T" + DateTime.Now.Ticks; 
        }

        public void newtable() {
            TABLE = generateTableName();
        }

        public void fail(String message) {
            Assert.Fail(message);
        }

        public void iferror(Client client, bool expected, bool got) {
            if (expected && !got) {
                fail("Error " + client.Error());
            }	
        }

        public void insertRow() {
            Client client = PubSubSQL.Factory.NewClient();
            ASSERT_CONNECT(client, ADDRESS, true);
            String command = string.Format("insert into {0} (col1, col2, col3) values (1:col1, 1:col2, 1:col3)", TABLE);
            ASSERT_EXECUTE(client, command, true);
            ASSERT_DISCONNECT(client);	
        }

        public void insertRows() {
            Client client = PubSubSQL.Factory.NewClient();
            ASSERT_CONNECT(client, ADDRESS, true);
            for (int row = 0; row < ROWS; row++) {	
                String command = string.Format("insert into {0} (col1, col2, col3) values ({1}:col1, {2}:col2, {3}:col3)", TABLE, row, row, row);
                ASSERT_EXECUTE(client, command, true);
            }
            ASSERT_DISCONNECT(client);
        }

        public void key(String column) {
            Client client = PubSubSQL.Factory.NewClient();
            ASSERT_CONNECT(client, ADDRESS, true);
            String command = string.Format("key {0} {1}", TABLE, column);
            ASSERT_EXECUTE(client, command, true);
            ASSERT_DISCONNECT(client);
        }

        public void tag(String column) {
            Client client = PubSubSQL.Factory.NewClient();
            ASSERT_CONNECT(client, ADDRESS, true);
            String command = string.Format("tag {0} {1}", TABLE, column);
            ASSERT_EXECUTE(client, command, true);
            ASSERT_DISCONNECT(client);
        }
        
        public void ASSERT_TRUE(bool val, String message) {
            if (!val) {
                fail("ASSERT_TRUE failed: " + message);
            }	
        }

        public void ASSERT_FALSE(bool val, String message) {
            if (val) {
                fail("ASSERT_FALSE failed: " + message);
            }
        }

        public void VALIDATE_RESULT(Client client, bool result) {
            if (result && !client.Ok() ) fail("VALIDATE_RESULT failed: expected Ok");
            if (!result && !client.Failed() ) fail("VALIDATE_RESULT failed: expected Failed");
        }

        public void ASSERT_CONNECT(Client client, String address, bool expected) {
            bool got = client.Connect(address);
            if (expected != got) {
                fail(string.Format("ASSERT_CONNECT failed: expected {0} got {1} ", expected, got));
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

        public void ASSERT_CONNECTED(Client client, bool expected) {
            bool got = client.Connected();
            if (expected != got) {
                fail(string.Format("ASSERT_CONNECTED failed: expected {0} got {1}", expected, got));
            }
        }

        public void ASSERT_EXECUTE(Client client, String command, bool expected) {
            bool got = client.Execute(command);
            if (expected != got) {
                fail(string.Format("ASSERT_EXECUTE failed: expected {0} got {1}", expected, got));	
            }
            iferror(client, expected, got);
            VALIDATE_RESULT(client, got);
        }

        public void ASSERT_ACTION(Client client, String expected) {
            String got = client.Action();
            if (expected != got) {
                fail(string.Format("ASSERT_ACTION failed: expected {0} got {1}", expected, got));
            }
        }

        public void ASSERT_ROW_COUNT(Client client, int expected) {
            int got = client.RowCount();
            if (expected != got) {
                fail(string.Format("ASSERT_ROW_COUNT failed: expected {0} but got {1}", expected, got));
            }
        }

        public void ASSERT_NEXT_ROW(Client client, bool expected) {
            bool got = client.NextRow();
            if (expected != got) {
                fail(string.Format("ASSERT_NEXT_ROW failed: expected {0} but got {1}", expected, got));
            }
        }

        public void ASSERT_ID(Client client) {
            String id = client.Value("id");
            if (string.IsNullOrEmpty(id)) {
                fail("ASSERT_ID failed: expected non empty string");
            }
        }

        public void ASSERT_PUBSUBID(Client client) {
            String pubsubid = client.PubSubId();	
            if (string.IsNullOrEmpty(pubsubid)) {
                fail("ASSERT_PUBSUBID failed: expected non empty string");
            }	
        }

        public void ASSERT_PUBSUBID_VALUE(Client client, String expected) {
            String got = client.PubSubId(); 		
            if (expected != got) {
                fail(string.Format("ASSERT_PUBSUBID_VALUE failed: expected {0} but got {1}", expected, got));
            }
        }

        public void ASSERT_VALUE(Client client, String column, String value, bool match) {
            String got = client.Value(column);	
            if (match && value != got) {
                fail(string.Format("ASSERT_VALUE failed: expected {0} but got {1}", value, got));
            }
            else if (!match && value == got) {
                fail(string.Format("ASSERT_VALUE failed: not expected {0}", value));
            }
        }

        public void ASSERT_COLUMN_COUNT(Client client, int expected) {
            int got = client.ColumnCount();
            if (expected != got) {
                fail(string.Format("ASSERT_COLUMN_COUNT failed: expected {0} but got {1}", expected, got));
            }
        }

        public void ASSERT_HAS_COLUMN(Client client, String column, bool expected) {
            bool got = client.HasColumn(column);
            if (expected != got) {
                fail(string.Format("ASSERT_HAS_COLUMN failed: expected {0} but got {1}", expected, got));
            }
        }	

        public void ASSERT_WAIT_FOR_PUBSUB(Client client, int timeout, bool expected) {
            bool got = client.WaitForPubSub(timeout);
            if (client.Failed()) {
                fail(string.Format("ASSERT_WAIT_FOR_PUBSUB failed: {0}", client.Error()));
            }
            else if (expected != got) {
                fail(string.Format("ASSERT_WAIT_FOR_PUBSUB failed: expected {0} but got {1}", expected, got));
            }	
        }

        public void ASSERT_NON_EMPTY_VALUE(Client client, int ordinal) {
            if (string.IsNullOrEmpty(client.ValueByOrdinal(ordinal))) {
                fail(string.Format("ASSERT_NON_EMPTY_VALUE failed: expected non empty string for ordinal {0}", ordinal));
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
                    fail(string.Format("ASSERT_PUBSUB_RESULT_SET failed expected {0} rows but got {1}", rows, readRows));
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
}
