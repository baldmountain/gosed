//
//  cmd.go
//  sed
//
// Copyright (c) 2009 Geoffrey Clements
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.
//

package sed

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

var (
	WrongNumberOfCommandParameters	os.Error	= os.ErrorString("Wrong number of parameters for command")
	UnknownScriptCommand		os.Error	= os.ErrorString("Unknown script command")
	InvalidSCommandFlag		os.Error	= os.ErrorString("Invalid flag for s command")
	RegularExpressionExpected	os.Error	= os.ErrorString("Expected a regular expression, got zero length string")
	UnterminatedRegularExpression	os.Error	= os.ErrorString("Unterminated regular expression")
	NoSupportForTwoAddress		os.Error	= os.ErrorString("This command doesn't support an address range or to end of file")
	NotImplemented			os.Error	= os.ErrorString("This command command hasn't been implemented yet")
)

type Cmd interface {
	fmt.Stringer
	processLine(s *Sed) (stop bool, err os.Error)
}

type Address interface {
	match(line []byte, lineNumber int) bool
}

const (
	ADDRESS_LINE	= iota
	ADDRESS_RANGE
	ADDRESS_TO_END_OF_FILE
	ADDRESS_LAST_LINE
	ADDRESS_REGEX
)

type address struct {
	not		bool
	address_type	int
	rangeStart	int
	rangeEnd	int
	regex		*regexp.Regexp
}

func (a *address) getTypeAsString() string {
	if a != nil {
		switch a.address_type {
		case ADDRESS_LINE:
			return "ADDRESS_LINE"
		case ADDRESS_RANGE:
			return "ADDRESS_RANGE"
		case ADDRESS_TO_END_OF_FILE:
			return "ADDRESS_TO_END_OF_FILE"
		case ADDRESS_LAST_LINE:
			return "ADDRESS_LAST_LINE"
		case ADDRESS_REGEX:
			return "ADDRESS_REGEX"
		default:
			return "ADDRESS_UNKNOWN"
		}
	}
	return "nil"
}

func (a *address) String() string {
	return fmt.Sprintf("address{type: %s rangeStart:%d rangeEnd:%d regex:%v}", a.getTypeAsString(), a.rangeStart, a.rangeEnd, a.regex)
}

func (a *address) match(line []byte, lineNumber int) bool {
	val := true
	if a != nil {
		switch a.address_type {
		case ADDRESS_LINE:
			val = lineNumber == a.rangeStart
		case ADDRESS_RANGE:
			val = lineNumber >= a.rangeStart && lineNumber <= a.rangeEnd
		case ADDRESS_TO_END_OF_FILE:
			val = lineNumber >= a.rangeStart
		case ADDRESS_LAST_LINE:
			val = false	// this is wrong!
		case ADDRESS_REGEX:
			val = a.regex.Match(line)
		default:
			val = false
		}
		if a.not {
			val = !val
		}
	}
	return val
}

func getNumberFromLine(s []byte) ([]byte, int, os.Error) {
	idx := 0
	for {
		if s[idx] < '0' || s[idx] > '9' {
			break
		}
		idx++
	}
	i, err := strconv.Atoi(string(s[0:idx]))
	if err != nil {
		return s, -1, err
	}
	return s[idx:], i, nil
}

// A nil address means match any line
func checkForAddress(s []byte) ([]byte, *address, os.Error) {
	var err os.Error
	if s[0] == '/' {
		// regular expression address
		s = s[1:]
		idx := bytes.IndexByte(s, '/')
		if idx < 0 {
			return s, nil, UnterminatedRegularExpression
		}
		r := s[0:idx]
		if len(r) == 0 {
			return s, nil, RegularExpressionExpected
		}
		// s is now just the command
		s = s[idx+1:]
		addr := new(address)
		addr.address_type = ADDRESS_REGEX
		addr.regex, err = regexp.Compile(string(r))
		if err != nil {
			return s, nil, err
		}
		return s, addr, nil
	} else if s[0] == '$' {
		// end of file
		addr := new(address)
		addr.address_type = ADDRESS_LAST_LINE
		// s is now just the command
		s = s[1:]
		return s, addr, nil
	} else if s[0] >= '0' && s[0] <= '9' {
		// numeric line address
		addr := new(address)
		addr.address_type = ADDRESS_LINE
		s, addr.rangeStart, err = getNumberFromLine(s)
		if err != nil {
			return s, nil, err
		}
		addr.rangeEnd = addr.rangeStart
		if s[0] == ',' {
			s = s[1:]
			if len(s) > 0 && s[0] >= '0' && s[0] <= '9' {
				addr.address_type = ADDRESS_RANGE
				s, addr.rangeEnd, err = getNumberFromLine(s)
				if err != nil {
					return s, nil, err
				}
				// if end range is less than start only match single line
				if addr.rangeEnd < addr.rangeStart {
					addr.address_type = ADDRESS_LINE
					addr.rangeEnd = 0
				}
			} else {
				addr.address_type = ADDRESS_TO_END_OF_FILE
			}
		}
		if s[0] == '!' {
			addr.not = true
			s = s[1:]
		}
		return s, addr, nil
	}
	return s, nil, nil
}

func NewCmd(s *Sed, line []byte) (Cmd, os.Error) {

	var err os.Error
	var addr *address
	line, addr, err = checkForAddress(line)
	if err != nil {
		return nil, err
	}

	if len(line) > 0 {
		switch line[0] {
		case 'a':
			return NewACmd(s, line, addr)
		case 'b':
			return NewBCmd(bytes.Split(line, []byte{'/'}, 0), addr)
		case 'c':
			return NewCCmd(s, line, addr)
		case 'd', 'D':
			return NewDCmd(bytes.Split(line, []byte{'/'}, 0), addr)
		case 'g', 'G':
			return NewGCmd(bytes.Split(line, []byte{'/'}, 0), addr)
		case 'h', 'H':
			return NewHCmd(bytes.Split(line, []byte{'/'}, 0), addr)
		case 'i':
			return NewICmd(s, line, addr)
		case 'n':
			return NewNCmd(bytes.Split(line, []byte{'/'}, 0), addr)
		case 'P', 'p':
			return NewPCmd(bytes.Split(line, []byte{'/'}, 0), addr)
		case 'q':
			return NewQCmd(bytes.Split(line, []byte{'/'}, 0), addr)
		case 'r':
			return NewQCmd(bytes.Split(line, []byte{'/'}, 0), addr)
		case 's':
			return NewSCmd(bytes.Split(line, []byte{'/'}, 0), addr)
		case '=':
			return NewEqlCmd(bytes.Split(line, []byte{'/'}, 0), addr)
		}
	}

	return nil, UnknownScriptCommand
}
