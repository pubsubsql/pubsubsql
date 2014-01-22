javac -d . api/NetHeader.java api/Client.java api/client.java
jar cfv pubsubsql.jar pubsubsql/NetHeader.class pubsubsql/Client.class pubsubsql/client.class
javac -d . -cp pubsubsql.jar test/PubSubSQLTest.java
