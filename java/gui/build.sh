#javac -d . gui/PubSubSQLGUI.java gui/MainForm.java
#jar cvmf gui/main_manifest pubsubsqlgui.jar PubSubSQLGUI.class PubSubSQLGUI1$.class MainForm.class gui/images/New.png
#java -jar pubsubsqlgui.jar
 
rm *.jar
javac -d . PubSubSQLGUI.java MainForm.java
jar cvf pubsubsqlgui.jar *.class  images/*.png   
jar ufe pubsubsqlgui.jar PubSubSQLGUI PubSubSQLGUI.class
rm *.class
java -jar pubsubsqlgui.jar

