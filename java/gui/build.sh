javac -d . -cp .:../lib/* PubSubSQLGUI.java MainForm.java TableDataset.java TableView.java
echo compiled
jar cvf pubsubsqlgui.jar *.class  images/*.png   
jar ufe pubsubsqlgui.jar PubSubSQLGUI PubSubSQLGUI.class
jar ufm pubsubsqlgui.jar manifest.txt 
rm *.class
mv pubsubsqlgui.jar ../pubsubsqlgui.jar
echo running...
java -jar ../pubsubsqlgui.jar

