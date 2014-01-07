using System;
using Microsoft.VisualStudio.TestTools.UnitTesting;
using System.Net.Sockets;
using PubSubSQL;

namespace PubSubSQLTest
{
    [TestClass]
    public class NetHelperTest
    {
        int port = 7777;

        [TestMethod]
        public void TestWriteRead()
        {
            Socket socket = new Socket(AddressFamily.InterNetwork, SocketType.Stream, ProtocolType.Tcp);
            socket.Connect("localhost", port);
            NetHelper helper = new NetHelper();
            helper.Set(socket, 100);
            helper.WriteWithHeader(1, NetHelper.ToUTF8("status"));
            NetHeader header = new NetHeader();
            byte[] bytes;
            helper.Read(ref header, out bytes);
            string beginExpected = "\n{\"status\":\"ok\"";
            string beginRead = NetHelper.FromUTF8(bytes).Substring(0, beginExpected.Length);
            Assert.AreEqual(beginExpected, beginRead);
            helper.Close();
        }

        public void TestReadTimeout()
        {
            Socket socket = new Socket(AddressFamily.InterNetwork, SocketType.Stream, ProtocolType.Tcp);
            socket.Connect("localhost", port);
            NetHelper helper = new NetHelper();
            helper.Set(socket, 100);
            NetHeader header = new NetHeader();
            byte[] bytes;
            bool ok = helper.ReadTimeout(200, ref header, out bytes);
            Assert.AreEqual(ok, false);
            helper.Close();
        }
    }
}
