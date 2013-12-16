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

import "flag"
import "strings"

type configuration struct {
	// logger
	LOG_DEBUG bool
 	LOG_INFO bool
	LOG_WARN bool
	LOG_ERROR bool

	// resources
	CHAN_RESPONSE_SENDER_BUFFER_SIZE int
	CHAN_TABLE_REQUESTS_BUFFER_SIZE int
	PARSER_SQL_INSERT_REQUEST_COLUMN_CAPACITY int
	PARSER_SQL_UPDATE_REQUEST_COLUMN_CAPACITY int
	TABLE_COLUMNS_CAPACITY int 
	TABLE_RECORDS_CAPACITY int 
	TABLE_GET_RECORDS_BY_TAG_CAPACITY int
	
	// command 
	COMMAND string

	// network
	IP string
	PORT uint	

	// run mode
	CLI bool	
	SERVER bool
}

func defaultConfig() configuration {
	return configuration {
		// logger
		LOG_DEBUG: false,
		LOG_INFO: true,
		LOG_WARN: true, 
		LOG_ERROR: true,

		// resources
		CHAN_RESPONSE_SENDER_BUFFER_SIZE: 10000, 
		CHAN_TABLE_REQUESTS_BUFFER_SIZE: 1000,
		PARSER_SQL_INSERT_REQUEST_COLUMN_CAPACITY: 10,
		PARSER_SQL_UPDATE_REQUEST_COLUMN_CAPACITY: 10,
		TABLE_COLUMNS_CAPACITY: 10, 
		TABLE_RECORDS_CAPACITY: 1000, 
		TABLE_GET_RECORDS_BY_TAG_CAPACITY: 20,

		// command 
		COMMAND: "start", 

		// network
		IP: "127.0.0.1",
		PORT: 7777,
		
		// run mode 
		CLI: true, 
		SERVER: true,
	}
}

var config = defaultConfig()

var validCommands = map[string]string {
	"start": "",
	"connect": "",
	"help": "", 
}

func validCommandsUsageString() string {
	str := "["
	for command, _ := range validCommands {
		str += " " + command 						
	}
	str += " ]"
	return str
}

func (c *configuration) setLogLevel(loglevel string) bool {
	c.LOG_DEBUG = false
	c.LOG_INFO = false
	c.LOG_WARN = false
	c.LOG_ERROR = false
	levels := strings.Split(loglevel, ",")		
	for _, s := range levels {
		switch s {
		case "debug":
			c.LOG_DEBUG = true	
		case "info":
			c.LOG_INFO = true
		case "warn":
			c.LOG_WARN = true
		case "error":
			c.LOG_ERROR = true
		default:
			return false
		}
	}
	return true
}

func (c *configuration) processCommandLine(args []string) bool {
	fset := flag.NewFlagSet("pubsubsql", flag.ExitOnError) 	
	// set up flags
	var loglevel string					
	fset.StringVar(&loglevel, "loglevel", "info", `logging level "debug,info,warn,error"`) 
	fset.StringVar(&c.IP, "ip", config.IP, "ip address")		
	fset.UintVar(&c.PORT, "port", config.PORT, "port number")
	fset.BoolVar(&c.CLI, "cli", config.CLI, "true indicates that server runs in interactive mode")

	// set command 
	if len(args) > 0 {	
		first := args[0]
		if first[0] != '-' {
			if len(args) > 1 {
				args = args[1:]
			} else { 
				args = nil
			}
			c.COMMAND = first	
		}
	}
	if _, contains := validCommands[c.COMMAND]; !contains {
		println("Command ", c.COMMAND, "is not valid. Valid commands ", validCommandsUsageString() )
		return false
	}

	// parse options
	if len(args) > 0 { 
		err := fset.Parse(args)
		if err != nil {
			println(err)
			fset.PrintDefaults()
			return false			
		}	
	}

	// log level
	// set loglevel
	if !c.setLogLevel(loglevel) {
		println("Invalid --loglevel option, usage: " + fset.Lookup("loglevel").Usage)	
		return false
	} 
		
	return true	
}

