using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.Data;
using System.Drawing;
using System.Linq;
using System.Text;
using System.Windows.Forms;

namespace PubSubSQLGUI
{
    public partial class MainForm : Form
    {
        public MainForm()
        {
            InitializeComponent();
            // set up gui events
            newMenu.Click += new_;
            newButton.Click += new_;
            exitMenu.Click += exit;
            connectLocalMenu.Click += connectLocal;     
            connectLocalButton.Click += connectLocal;     
            connectMenu.Click += connect;
            connectButton.Click += connect;
            disconnectMenu.Click += disconnect;
            disconnectButton.Click += disconnect;
            executeMenu.Click += execute;
            executeButton.Click += execute;
            nextPaneMenu.Click += nextPane;
            aboutMenu.Click += about;
        }

        // gui events

        private void new_(object sender, EventArgs e)
        {
            
        }

        private void exit(object sender, EventArgs e)
        {
            this.Close();
        }

        private void connectLocal(object sender, EventArgs e)
        {

        }

        private void connect(object sender, EventArgs e)
        {

        }

        private void disconnect(object sender, EventArgs e)
        {

        }

        private void execute(object sender, EventArgs e)
        {

        }

        private void cancelExecute(object sender, EventArgs e)
        {

        }

        private void nextPane(object sender, EventArgs e)
        {

        }

        private void about(object sender, EventArgs e)
        {

        }
    }
}
