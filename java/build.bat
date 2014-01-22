javac -d . api/NetHeader.java api/Client.java
jar cfv pubsubsql.jar pubsubsql/NetHeader.class pubsubsql/Client.class
javac -d . -cp pubsubsql.jar test/PubSubSQLTest.java
