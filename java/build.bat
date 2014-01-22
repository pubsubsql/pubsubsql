javac -d . api/NetHeader.java
jar cvf pubsubsql.jar pubsubsql/NetHeader.class
javac -d . -cp pubsubsql.jar test/PubSubSQLTest.java
