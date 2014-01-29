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
    public partial class SimulatorForm : Form
    {
        public string TableName = string.Format("T{0:HHmmss}", DateTime.Now);
        public int Columns = 0;
        public int Rows = 0;
        public SimulatorForm()
        {
            InitializeComponent();
        }

        private void SimulatorForm_Load(object sender, EventArgs e)
        {
            CenterToParent();
        }

        private void okButton_Click_1(object sender, EventArgs e)
        {
            Columns = Convert.ToInt32(columnsUpDown.Value);
            Rows = Convert.ToInt32(rowsUpDown.Value);
        }
    }
}
