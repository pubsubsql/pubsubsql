using System;
using Microsoft.VisualStudio.TestTools.UnitTesting;

namespace PubSubSQLTest
{
    [TestClass]
    public class NetHeaderTest
    {
        [TestMethod]
        public void TestWriteRead()
        {
            PubSubSQL.NetHeader header1 = new PubSubSQL.NetHeader(32567, 9875235);
            PubSubSQL.NetHeader header2 = new PubSubSQL.NetHeader();
            byte[] bytes = new byte[100];
            header1.WriteTo(bytes);
            header2.ReadFrom(bytes);
            Assert.AreEqual(header1, header2);
        }

        [TestMethod]
        public void TestGetBytes()
        {
            PubSubSQL.NetHeader header1 = new PubSubSQL.NetHeader(32567, 9875235);
            PubSubSQL.NetHeader header2 = new PubSubSQL.NetHeader();
            byte[] bytes = header1.GetBytes();
            header1.WriteTo(bytes);
            header2.ReadFrom(bytes);
            Assert.AreEqual(header1, header2);
        }
    }
}
