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
                string command = "insert into " + TestUtils.TABLE + " (col1, col2, col3) values ( val1, val2, val3 ) ";
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
            string command = "select * from " + TestUtils.TABLE;
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
            string command = "update " + TestUtils.TABLE + " set col1 = val1updated, col2 = val2updated, col3 = val3updated ";
            TestUtils.ASSERT_EXECUTE(client, command, "update failed");
            TestUtils.ASSERT_RECORD_COUNT(client, TestUtils.ROWS);
            client.Disconnect();  
        }

        [TestMethod]
        public void TestDelete()
        {
            Client client = Factory.NewClient();
            TestUtils.ASSERT_CONNECT(client);
            string command = "delete from " + TestUtils.TABLE;
            TestUtils.ASSERT_EXECUTE(client, command, "delete failed");
            TestUtils.ASSERT_RECORD_COUNT(client, TestUtils.ROWS);
            client.Disconnect();  
        }

        [TestMethod]
        public void TestKey()
        {
            Client client = Factory.NewClient();
            TestUtils.ASSERT_CONNECT(client);
            string command = "key " + TestUtils.TABLE + " col1 ";
            TestUtils.ASSERT_EXECUTE(client, command, "key failed");
            client.Disconnect();  
        }

        [TestMethod]
        public void TestTag()
        {
            Client client = Factory.NewClient();
            TestUtils.ASSERT_CONNECT(client);
            string command = "tag " + TestUtils.TABLE + " col2 ";
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
            string command = "subscribe * from " + TestUtils.TABLE;
            TestUtils.ASSERT_EXECUTE(client, command, "subscribe failed");
            TestUtils.ASSERT_PUBSUBID(client);
            command = string.Format("unsubscribe from {0} where pubsubid = {1}", TestUtils.TABLE, client.PubSubId());
            TestUtils.ASSERT_EXECUTE(client, command, "unsubscribe failed");
            client.Disconnect();
        }

        [TestMethod]
        public void TestWaitPubSubAdd()
        {
            
        }

    }
}
