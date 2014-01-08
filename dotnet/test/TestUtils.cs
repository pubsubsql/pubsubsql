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
                ASSERT_FALSE(client.Ok());
                ASSERT_TRUE(client.Failed());
            }
            ASSERT_TRUE(client.Ok());
            ASSERT_FALSE(client.Failed());
        }
    }
}
