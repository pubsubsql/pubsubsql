using System;
using Microsoft.VisualStudio.TestTools.UnitTesting;

namespace PubSubSQLTest
{
    [TestClass]
    public class NetworkHeaderTest
    {
        [TestMethod]
        public void Test()
        {
            UInt32 messageSize = 2389;
            UInt32 requestId = 2;
            PubSubSQL.NetworkHeader header1 = new PubSubSQL.NetworkHeader(messageSize, requestId);
            byte[] bytes = header1.GetBytes();
            PubSubSQL.NetworkHeader header2 = new PubSubSQL.NetworkHeader();
            header2.ReadFrom(bytes);
            Assert.AreEqual(header1, header2);
        }
    }
}
