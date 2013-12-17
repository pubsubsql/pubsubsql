/* Copyright (C) 2013 CompleteDB LLC.
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with PubSubSQL.  If not, see <http://www.gnu.org/licenses/>.
 */

package pubsubsql

import "strconv"
import "bufio"
import "os"
import "strings"

/*
import "time"
*/

type lineReader struct {
	reader *bufio.Reader		
	quit string
	line string
}

func newLineReader(quit string) *lineReader {
	return &lineReader {
		reader: bufio.NewReader(os.Stdin),
		quit: quit,
	}
}

func (l *lineReader) readLine() bool {
	line, err := l.reader.ReadString('\n')
	l.line = strings.TrimSpace(line)	
	if err != nil {
		return false
	}
	return l.line != l.quit
}

type cli struct {
	prefix string		
	quit *QuitChan 
	input chan string	
}

func newCli() *cli {
	return &cli {
		quit: NewQuitChan(),
		input: make(chan string),
	}  
}

func (c *cli) readInput() {
	l := newLineReader("q")
	for l.readLine() {
		if len(l.line) > 0 {
			c.input <- l.line	
		}			
	}
	c.quit.Quit(0)
}

func (c *cli) run() {
	c.initPrefix()
	go c.readInput()
	
	cout := bufio.NewWriter(os.Stdout)
	for {
		cout.WriteString(c.prefix)	
		cout.Flush()
		select {
		case input := <-c.input:
			cout.WriteString(input)			
			cout.WriteString("\n")
			cout.Flush()
		case <-c.quit.GetChan():
			return
		}	
	}
}

func (c *cli) initPrefix() {
	def := defaultConfig()
	c.prefix = "pubsubsql"
	if def.IP != config.IP {
		c.prefix += " " + config.netAddress()
	} else if def.PORT != config.PORT {
		c.prefix += ":" + strconv.Itoa(int(config.PORT))
	}
	c.prefix += ">"
}

