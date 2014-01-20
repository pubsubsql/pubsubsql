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
    public partial class ConnectForm : Form
    {
        public string TableName = string.Format("T{0:HHmmss}", DateTime.Now);
        public string Address = string.Empty;
        public ConnectForm()
        {
            InitializeComponent();
        }

        private void SimulatorForm_Load(object sender, EventArgs e)
        {
            CenterToParent();
        }

        private void okButton_Click_1(object sender, EventArgs e)
        {
            Address = string.Format("{0}:{1}", hostText.Text, Convert.ToInt32(portUpDown.Value));
        }
    }
}
