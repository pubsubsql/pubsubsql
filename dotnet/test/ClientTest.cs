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
        
    }
}
