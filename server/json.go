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
 * You should have idxeived a copy of the GNU Affero General Public License
 * along with PubSubSQL.  If not, see <http://www.gnu.org/licenses/>.
 */

package pubsubsql

import "bytes"
import "unicode/utf8"

type JSONBuilder struct {
	bytes.Buffer
	err bool //
}

// implementation for string function was copied from go source code
var hex = "0123456789abcdef"
func (j *JSONBuilder) string(s string) (int) {
    len0 := j.Len() 
    j.WriteByte('"')
    start := 0
    for i := 0; i < len(s); {
        if b := s[i]; b < utf8.RuneSelf {
            if 0x20 <= b && b != '\\' && b != '"' && b != '<' && b != '>' {
                i++     
                continue
            }       
            if start < i {
                j.WriteString(s[start:i])
            }       
            switch b {
            case '\\', '"':
                j.WriteByte('\\')
                j.WriteByte(b)
            case '\n':
                j.WriteByte('\\')
                j.WriteByte('n')
            case '\r':
                j.WriteByte('\\')
                j.WriteByte('r')
            default:
                // This encodes bytes < 0x20 except for \n and \r, 
                // as well as < and >. The latter are escaped because they
                // can lead to security holes when user-controlled strings 
                // are rendered into JSON and served to some browsers.
                j.WriteString(`\u00`)
                j.WriteByte(hex[b>>4])
                j.WriteByte(hex[b&0xF])
            }       
            i++     
            start = i
            continue
        }       
        c, size := utf8.DecodeRuneInString(s[i:])
        if c == utf8.RuneError && size == 1 {
			j.err = true
        }       
        i += size 
    }
    if start < len(s) {
        j.WriteString(s[start:])
    }
    j.WriteByte('"')
    return j.Len() - len0
}

func (j *JSONBuilder) beginArray() {
	j.WriteByte('[')	
}

func (j *JSONBuilder) beginObject() {
	j.WriteByte('{')	
}
  
func (j *JSONBuilder) endArray() {
	j.WriteByte(']')	
}

func (j *JSONBuilder) endObject() {
	j.WriteByte('}')	
}

func (j *JSONBuilder) nameSeparator() {
	j.WriteByte(':')
}

func (j *JSONBuilder) valueSeparator() {
	j.WriteByte(',')
}

var errorString = `{ "status":"error" "msg":"Failed to build json document due to invalid utf8 string."`
func (j *JSONBuilder) Bytes() ([]byte) {
	if j.err {
		return []byte(errorString)
	}				
	return j.Bytes()	
}

