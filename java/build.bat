javac -d . api/NetHeader.java
jar cfv pubsubsql.jar pubsubsql/NetHeader.class
javac -d . -cp pubsubsql.jar test/PubSubSQLTest.java
