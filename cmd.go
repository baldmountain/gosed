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
	"bytes";
	"fmt";
	"os";
	"regexp";
	"strconv";
)

var (
	WrongNumberOfCommandParameters	os.Error	= os.ErrorString("Wrong number of parameters for command");
	UnknownScriptCommand		os.Error	= os.ErrorString("Unknown script command");
	InvalidSCommandFlag		os.Error	= os.ErrorString("Invalid flag for s command");
	RegularExpressionExpected	os.Error	= os.ErrorString("Expected a regular expression, got zero length string");
	UnterminatedRegularExpression	os.Error	= os.ErrorString("Unterminated regular expression");
)

type Cmd interface {
	fmt.Stringer;
	processLine(s *Sed) (stop bool, err os.Error);
}

type Address interface {
	match(line []byte, lineNumber, totalNumberOfLines int) bool;
}

const (
	ADDRESS_LINE	= iota;
	ADDRESS_RANGE;
	ADDRESS_TO_END_OF_FILE;
	ADDRESS_LAST_LINE;
	ADDRESS_REGEX;
)

type address struct {
	address_type	int;
	rangeStart	int;
	rangeEnd	int;
	regex		*regexp.Regexp;
}

func (a *address)getTypeAsString() string {
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
	return "nil";
}

func (a *address) String() string {
	return fmt.Sprintf("address{type: %s rangeStart:%d rangeEnd:%d regex:%v}", a.getTypeAsString(), a.rangeStart, a.rangeEnd, a.regex)
}

func (a *address) match(line []byte, lineNumber, totalNumberOfLines int) bool {
	if a != nil {
		switch a.address_type {
		case ADDRESS_LINE:
			return lineNumber == a.rangeStart
		case ADDRESS_RANGE:
			return lineNumber >= a.rangeStart && lineNumber <= a.rangeEnd
		case ADDRESS_TO_END_OF_FILE:
			return lineNumber >= a.rangeStart
		case ADDRESS_LAST_LINE:
			return lineNumber == totalNumberOfLines
		case ADDRESS_REGEX:
			return a.regex.Match(line)
		default:
			return false
		}
	}
	return true;
}

func getNumberFromLine(s []byte) ([]byte, int, os.Error) {
	idx := 0;
	for {
		if s[idx] < '0' || s[idx] > '9' {
			break
		}
		idx++;
	}
	i, err := strconv.Atoi(string(s[0:idx]));
	if err != nil {
		return s, -1, err
	}
	return s[idx:], i, nil;
}

// A nil address means match any line
func checkForAddress(s []byte) ([]byte, *address, os.Error) {
	var err os.Error;
	if s[0] == '/' {
		// regular expression address
		s = s[1:];
		idx := bytes.IndexByte(s, '/');
		if idx < 0 {
		  return s, nil, UnterminatedRegularExpression;
		}
		r := s[0:idx];
		if len(r) == 0 {
		  return s, nil, RegularExpressionExpected;
		}
		// s is now just the command
		s = s[idx+1:];
		addr := new(address);
		addr.address_type = ADDRESS_REGEX;
		addr.regex, err = regexp.Compile(string(r));
		if err != nil {
			return s, nil, err
		}
		return s, addr, nil;
	} else if s[0] == '$' {
		// end of file
		addr := new(address);
		addr.address_type = ADDRESS_LAST_LINE;
		// s is now just the command
		s = s[1:];
		return s, addr, nil;
	} else if s[0] >= '0' && s[0] <= '9' {
		// numeric line address
		addr := new(address);
		addr.address_type = ADDRESS_LINE;
		s, addr.rangeStart, err = getNumberFromLine(s);
		if err != nil {
			return s, nil, err
		}
		addr.rangeEnd = addr.rangeStart;
		if s[0] == ',' {
			s = s[1:];
			if len(s) > 0 && s[0] >= '0' && s[0] <= '9' {
    		addr.address_type = ADDRESS_RANGE;
				s, addr.rangeEnd, err = getNumberFromLine(s);
				if err != nil {
					return s, nil, err
				}
			} else {
    		addr.address_type = ADDRESS_TO_END_OF_FILE;
			}
		}
		return s, addr, nil;
	}
	return s, nil, nil;
}

func NewCmd(s *Sed, line []byte) (Cmd, os.Error) {

	var err os.Error;
	var addr *address;
	line, addr, err = checkForAddress(line);
	if err != nil {
		return nil, err
	}

	pieces := bytes.Split(line, []byte{'/'}, 0);

	if len(pieces[0]) > 0 {
		switch pieces[0][0] {
		case 's':
			return NewSCmd(pieces, addr)
		case 'q':
			return NewQCmd(pieces, addr)
		case 'd', 'D':
			return NewDCmd(pieces, addr)
		case 'P', 'p':
			return NewPCmd(pieces, addr)
		case 'n':
			return NewNCmd(pieces, addr)
		case '=':
			return NewEqlCmd(pieces, addr)
		case 'a':
			return NewACmd(s, pieces, addr)
		case 'i':
			return NewICmd(s, pieces, addr)
		case 'g', 'G':
			return NewGCmd(pieces, addr)
		case 'h', 'H':
			return NewHCmd(pieces, addr)
		}
	}

	return nil, UnknownScriptCommand;
}

func copyByteSlice(a []byte) []byte {
	newSlice := make([]byte, len(a));
	copy(newSlice, a);
	return newSlice;
}
