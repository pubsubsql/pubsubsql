javac -d . -cp .:lib/gson-2.2.4.jar api/NetHeader.java api/Client.java api/ClientImpl.java api/Factory.java api/NetHelper.java api/ResponseData.java 
jar cfv lib/pubsubsql.jar pubsubsql/NetHeader.class pubsubsql/Client.class pubsubsql/ClientImpl.class pubsubsql/Factory.class pubsubsql/NetHelper.class pubsubsql/ResponseData.class
javac -d . -cp .:lib/* test/PubSubSQLTest.java
java -cp .:lib/* PubSubSQLTest
