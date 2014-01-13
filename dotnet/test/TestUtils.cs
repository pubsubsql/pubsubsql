using System;
using Microsoft.VisualStudio.TestTools.UnitTesting;
using System.Net.Sockets;
using PubSubSQL;

namespace PubSubSQLTest
{
    class TestUtils
    {
        static readonly string ADDRESS = "localhost:7777";
        public static int ROWS = 300;
        public static string TABLE = generateTableName(); 

        public static string generateTableName()
        {
            return "T" + DateTime.Now.Ticks.ToString();
        }

        public static void ASSERT_TRUE(bool b)
        {
            Assert.AreEqual(true, b);
        }

        public static void ASSERT_FALSE(bool b)
        {
            Assert.AreEqual(false, b);
        }

        public static void ASSERT_CONNECT(Client client)
        {
            if (!client.Connect(ADDRESS))
            {
                Assert.Fail("Connect failed " + client.Error());
            }
            ASSERT_TRUE(client.Ok());
            ASSERT_FALSE(client.Failed());
        }

        public static void ASSERT_EXECUTE(Client client, string command, string err)
        {
            if (!client.Execute(command))
            {
                Assert.Fail("Execute failed " + client.Error() + " " + err);
                Assert.Fail(client.JSON());
                ASSERT_FALSE(client.Ok());
                ASSERT_TRUE(client.Failed());
                return;
            }
            ASSERT_TRUE(client.Ok());
            ASSERT_FALSE(client.Failed());
        }

        public static void ASSERT_ACTION(Client client, string action)
        {
            if (client.Action() != action)
            {
                Assert.Fail("Expected action {0} but got {1} ", action ,client.Action());
            }
        }

        public static void ASSERT_ID(Client client)
        {
            if (client.Id() == string.Empty)
            {
                Assert.Fail("Expected id but got empty string");
            }
        }

        public static void ASSERT_COLUMNS(Client client, int count)
        {
            if (client.ColumnCount() != count)
            {
                Assert.Fail("Expected columns count {0} but got {1}", count, client.ColumnCount());
            }
        }

        public static void ASSERT_COLUMN(Client client, string column)
        {
            if (!client.HasColumn(column))
            {
                Assert.Fail("Expected column {0}", column);
            }
        }

        public static void ASSERT_RECORD_COUNT(Client client, int count)
        {
            if (client.RecordCount() != count)
            {
                Assert.Fail("Expected records count {0} but got {1}", count, client.RecordCount());
            }
        }

        public static void ASSERT_VALUE(Client client, string column, string value)
        {
            if (client.Value(column) != value)
            {
                Assert.Fail("Expected value {0} for column {1} but got {2}", value, column, client.Value(column));
            }
        }

        public static void ASSERT_PUBSUBID(Client client)
        {
            if (client.PubSubId()== string.Empty)
            {
                Assert.Fail("Expected pubsubid but got empty string");
            }
        }
    }
}
