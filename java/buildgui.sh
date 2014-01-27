javac -d . gui/PubSubSQLGUI.java gui/MainForm.java
jar cvmf gui/main_manifest pubsubsqlgui.jar PubSubSQLGUI.class MainForm.class
java -cp .:lib/* PubSubSQLGUI
