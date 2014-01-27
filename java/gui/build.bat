javac -d . PubSubSQLGUI.java MainForm.java
jar cvf pubsubsqlgui.jar *.class images/*.png  
jar ufe pubsubsqlgui.jar PubSubSQLGUI PubSubSQLGUI.class
del *.class
java -jar pubsubsqlgui.jar

