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
            for (int i = 0; i < TestUtils.ROWS; i++)
            {
                string command = "select * from " + TestUtils.TABLE;
                TestUtils.ASSERT_EXECUTE(client, command, "select failed");
                TestUtils.ASSERT_ACTION(client, "select");
                // extra column for id
                TestUtils.ASSERT_COLUMNS(client, 4);
            }
            client.Disconnect();  
        }
    }
}
