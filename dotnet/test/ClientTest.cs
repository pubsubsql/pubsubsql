using System;
using Microsoft.VisualStudio.TestTools.UnitTesting;
using PubSubSQL;

namespace PubSubSQLTest
{
    [TestClass]
    public class ClientTest
    {
        
        [TestMethod]
        public void TestConnectDisconnect()
        {
            Client client = Factory.NewClient();
            TestUtils.ASSERT_CONNECT(client);
            client.Disconnect();
        }

        [TestMethod]
        public void TestStatus()
        {
            Client client = Factory.NewClient();
            TestUtils.ASSERT_CONNECT(client);
            TestUtils.ASSERT_EXECUTE(client, "status", "status failed");
            TestUtils.ASSERT_ACTION(client, "status");
            client.Disconnect();
        }

        [TestMethod]
        public void TestInsert()
        {
            Client client = Factory.NewClient();
            TestUtils.ASSERT_CONNECT(client);
            for (int i = 0; i < TestUtils.ROWS; i++)
            {
                string command = string.Format("insert into {0} (col1, col2, col3) values ( val1, val2, val3 ) ", TestUtils.TABLE);
                TestUtils.ASSERT_EXECUTE(client, command, "insert failed");
                TestUtils.ASSERT_ACTION(client, "insert");
                TestUtils.ASSERT_ID(client);
            }
            client.Disconnect();  
        }

        [TestMethod]
        public void TestSelect()
        {
            Client client = Factory.NewClient();
            TestUtils.ASSERT_CONNECT(client);
            string command = string.Format("select * from {0} ", TestUtils.TABLE);
            TestUtils.ASSERT_EXECUTE(client, command, "select failed");
            TestUtils.ASSERT_ACTION(client, "select");
            TestUtils.ASSERT_COLUMNS(client, 4);
            TestUtils.ASSERT_RECORD_COUNT(client, TestUtils.ROWS);
            TestUtils.ASSERT_COLUMN(client, "col1");  
            TestUtils.ASSERT_COLUMN(client, "col2");  
            TestUtils.ASSERT_COLUMN(client, "col2");
            int rows = 0;
            while (rows < TestUtils.ROWS)
            {
                TestUtils.ASSERT_TRUE(client.NextRecord());
                TestUtils.ASSERT_VALUE(client, "col1", "val1");
                TestUtils.ASSERT_VALUE(client, "col2", "val2");
                TestUtils.ASSERT_VALUE(client, "col3", "val3");
                rows++;
            }
            client.Disconnect();  
        }

        [TestMethod]
        public void TestUpdate()
        {
            Client client = Factory.NewClient();
            TestUtils.ASSERT_CONNECT(client);
            string command = string.Format("update {0} set col1 = val1updated, col2 = val2updated, col3 = val3updated ", TestUtils.TABLE);
            TestUtils.ASSERT_EXECUTE(client, command, "update failed");
            TestUtils.ASSERT_RECORD_COUNT(client, TestUtils.ROWS);
            client.Disconnect();  
        }

        [TestMethod]
        public void TestDelete()
        {
            Client client = Factory.NewClient();
            TestUtils.ASSERT_CONNECT(client);
            string command = string.Format("delete from {0} ", TestUtils.TABLE);
            TestUtils.ASSERT_EXECUTE(client, command, "delete failed");
            TestUtils.ASSERT_RECORD_COUNT(client, TestUtils.ROWS);
            client.Disconnect();  
        }

        [TestMethod]
        public void TestKey()
        {
            Client client = Factory.NewClient();
            TestUtils.ASSERT_CONNECT(client);
            string command = string.Format("key {0} col1 ", TestUtils.TABLE);
            TestUtils.ASSERT_EXECUTE(client, command, "key failed");
            client.Disconnect();  
        }

        [TestMethod]
        public void TestTag()
        {
            Client client = Factory.NewClient();
            TestUtils.ASSERT_CONNECT(client);
            string command = string.Format("tag {0} col2 ", TestUtils.TABLE);
            TestUtils.ASSERT_EXECUTE(client, command, "tag failed");
            client.Disconnect();  
        }

        [TestMethod]
        public void TestInsertUniqueKey()
        {
            Client client = Factory.NewClient();
            TestUtils.ASSERT_CONNECT(client);
            for (int i = 0; i < TestUtils.ROWS; i++)
            {
                string command = string.Format("insert into {0} (col1, col2, col3) values ( {1}, {2}, {3} ) ", TestUtils.TABLE, i, i, i );
                TestUtils.ASSERT_EXECUTE(client, command, "insert failed");
                TestUtils.ASSERT_ACTION(client, "insert");
                TestUtils.ASSERT_ID(client);
            }
            client.Disconnect();  
        }

        [TestMethod]
        public void TestSubscribeUnsubscribe()
        {
            Client client = Factory.NewClient();
            TestUtils.ASSERT_CONNECT(client);
            string command = string.Format("subscribe * from {0}", TestUtils.TABLE);
            TestUtils.ASSERT_EXECUTE(client, command, "subscribe failed");
            TestUtils.ASSERT_PUBSUBID(client);
            command = string.Format("unsubscribe from {0} where pubsubid = {1}", TestUtils.TABLE, client.PubSubId());
            TestUtils.ASSERT_EXECUTE(client, command, "unsubscribe failed");
            client.Disconnect();
        }

        [TestMethod]
        public void TestWaitPubSubAdd()
        {
            Client client = Factory.NewClient();
            TestUtils.ASSERT_CONNECT(client);
            string command = string.Format("subscribe * from {0} ", TestUtils.TABLE);
            TestUtils.ASSERT_EXECUTE(client, command, "subscribe failed");
            TestUtils.ASSERT_PUBSUBID(client);
            string pubsubid = client.PubSubId();
            TestUtils.ASSERT_TRUE(client.WaitForPubSub(100));
            TestUtils.ASSERT_ACTION(client, "add");
            TestUtils.ASSERT_PUBSUBID(client);
            TestUtils.ASSERT_TRUE(pubsubid == client.PubSubId());
            TestUtils.ASSERT_COLUMNS(client, 4);
            int rows = 0;
            for (rows = 0; rows < TestUtils.ROWS; rows++ )
            {
                TestUtils.ASSERT_TRUE(client.NextRecord());
                string val = rows.ToString();
                TestUtils.ASSERT_COLUMN(client, "id");
                TestUtils.ASSERT_VALUE(client, "col1", val);
                TestUtils.ASSERT_VALUE(client, "col2", val);
                TestUtils.ASSERT_VALUE(client, "col3", val);
            }
            client.Disconnect();
        }

        [TestMethod]
        public void TestWaitPubSubUpdate()
        {
            Client client = Factory.NewClient();
            TestUtils.ASSERT_CONNECT(client);
            string command = string.Format("subscribe skip * from {0} ", TestUtils.TABLE);
            TestUtils.ASSERT_EXECUTE(client, command, "subscribe failed");
            TestUtils.ASSERT_ACTION(client, "subscribe");
            TestUtils.ASSERT_PUBSUBID(client);
            string pubsubid = client.PubSubId();
            // generate pubsub update action
            string val = "updatedvalue";
            command = string.Format("update {0} set col3 = {1}", TestUtils.TABLE, val); 
            TestUtils.ASSERT_EXECUTE(client, command, "update failed");
            TestUtils.ASSERT_ACTION(client, "update");
            //
            for (int rows = 0; rows < TestUtils.ROWS; rows++ )
            {
                TestUtils.ASSERT_TRUE(client.WaitForPubSub(10));
                TestUtils.ASSERT_ACTION(client, "update");
                TestUtils.ASSERT_PUBSUBID(client);
                TestUtils.ASSERT_TRUE(pubsubid == client.PubSubId());
                TestUtils.ASSERT_COLUMNS(client, 2);
                TestUtils.ASSERT_TRUE(client.NextRecord());
                TestUtils.ASSERT_COLUMN(client, "id");
                TestUtils.ASSERT_VALUE(client, "col3", val);
            }
            client.Disconnect();
        }

        [TestMethod]
        public void TestWaitPubSubRemove()
        {
            Client client = Factory.NewClient();
            TestUtils.ASSERT_CONNECT(client);
            string val = "subscribedvalue";
            // first update tag with the same value so that when we subscribe to the tag and update
            // the tag with different value we will get pubsub remove action
            string command = string.Format("update {0} set col2 = {1}", TestUtils.TABLE, val);
            TestUtils.ASSERT_EXECUTE(client, command, "update failed");
            command = string.Format("subscribe skip * from {0} where col2 = {1}", TestUtils.TABLE, val);
            TestUtils.ASSERT_EXECUTE(client, command, "subscribe failed");
            TestUtils.ASSERT_ACTION(client, "subscribe");
            TestUtils.ASSERT_PUBSUBID(client);
            string pubsubid = client.PubSubId();
            // now update the tag and generate pubsub remove    
            val = "newtagvalue";
            command = string.Format("update {0} set col2 = {1}", TestUtils.TABLE, val);
            TestUtils.ASSERT_EXECUTE(client, command, "update failed");
            TestUtils.ASSERT_ACTION(client, "update");
            for (int rows = 0; rows < TestUtils.ROWS; rows++ )
            {
                TestUtils.ASSERT_TRUE(client.WaitForPubSub(10));
                TestUtils.ASSERT_ACTION(client, "remove");
                TestUtils.ASSERT_PUBSUBID(client);
                TestUtils.ASSERT_TRUE(pubsubid == client.PubSubId());
                TestUtils.ASSERT_ID(client);
            }
            client.Disconnect();
        }

        [TestMethod]
        public void TestWaitPubSubDelete()
        {
            Client client = Factory.NewClient();
            TestUtils.ASSERT_CONNECT(client);
            string command = string.Format("subscribe skip * from {0}", TestUtils.TABLE);
            TestUtils.ASSERT_EXECUTE(client, command, "subscribe failed");
            TestUtils.ASSERT_ACTION(client, "subscribe");
            TestUtils.ASSERT_PUBSUBID(client);
            string pubsubid = client.PubSubId();
            command = string.Format("delete from {0} ", TestUtils.TABLE);
            TestUtils.ASSERT_EXECUTE(client, command, "delete failed");
            TestUtils.ASSERT_ACTION(client, "delete");
            for (int rows = 0; rows < TestUtils.ROWS; rows++)
            {
                TestUtils.ASSERT_TRUE(client.WaitForPubSub(10));
                TestUtils.ASSERT_ACTION(client, "delete");
                TestUtils.ASSERT_PUBSUBID(client);
                TestUtils.ASSERT_TRUE(pubsubid == client.PubSubId());
                TestUtils.ASSERT_ID(client);
            }
            client.Disconnect();
        }

        [TestMethod]
        public void TestWaitPubSubInsert()
        {
            Client client = Factory.NewClient();
            TestUtils.ASSERT_CONNECT(client);
            string command = string.Format("subscribe * from {0}", TestUtils.TABLE);
            TestUtils.ASSERT_EXECUTE(client, command, "subscribe failed");
            TestUtils.ASSERT_ACTION(client, "subscribe");
            TestUtils.ASSERT_PUBSUBID(client);
            string pubsubid = client.PubSubId();
            // insert record to generate pubsub insert
            for (int row = 0; row < TestUtils.ROWS; row++)
            {
                command = string.Format("insert into {0} (col1, col2, col3) values ({1}, {1}, {1})", TestUtils.TABLE, row);
                TestUtils.ASSERT_EXECUTE(client, command, "insert failed");
                TestUtils.ASSERT_ACTION(client, "insert");
                TestUtils.ASSERT_ID(client);
            }
            //
            for (int rows = 0; rows < TestUtils.ROWS; rows++)
            {
                string val = rows.ToString();
                TestUtils.ASSERT_TRUE(client.WaitForPubSub(10));
                TestUtils.ASSERT_ACTION(client, "insert");
                TestUtils.ASSERT_TRUE(pubsubid == client.PubSubId());
                TestUtils.ASSERT_TRUE(client.NextRecord());
                TestUtils.ASSERT_COLUMN(client, "id");
                TestUtils.ASSERT_VALUE(client, "col1", val);
                TestUtils.ASSERT_VALUE(client, "col2", val);
                TestUtils.ASSERT_VALUE(client, "col3", val);
            }
            client.Disconnect();
        }
        

    }
}
