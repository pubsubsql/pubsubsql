javac -d . -cp .;..\lib\* *.java
echo compiled
jar cvf pubsubsqlgui.jar *.class  images\*.png   
jar ufe pubsubsqlgui.jar PubSubSQLGUI PubSubSQLGUI.class
jar ufm pubsubsqlgui.jar manifest.txt 
del *.class
move pubsubsqlgui.jar ..\pubsubsqlgui.jar
echo running...
java -jar ..\pubsubsqlgui.jar

