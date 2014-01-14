using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;

namespace PubSubSQLGUI
{
    class ListViewDoubleBuffered : System.Windows.Forms.ListView
    {
        public ListViewDoubleBuffered()
            : base()
        {
            DoubleBuffered = true;
        }
    }
}
