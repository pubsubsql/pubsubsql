using System;
using Microsoft.VisualStudio.TestTools.UnitTesting;
using System.Net.Sockets;
using PubSubSQL;

namespace PubSubSQLTest
{
    class TestUtils
    {
        static readonly string ADDRESS = "localhost:7777";

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
                Assert.Fail("Expected action " + action + " but got " + client.Action() + "\n" + client.JSON());
            }
        }
    }
}
