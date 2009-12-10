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
)

type Cmd interface {
	fmt.Stringer;
	getAddress() *address;
	processLine(s *Sed) (stop bool, err os.Error);
}

type address struct {
	rangeStart	int;
	rangeEnd	int;
	lastLine	bool;
	regex		*regexp.Regexp;
}

type command struct {
	addr *address;
}

func (a *address) String() string {
	return fmt.Sprintf("address{rangeStart:%d rangeEnd:%d lastLine:%t regex:%v}", a.rangeStart, a.rangeEnd, a.lastLine, a.regex)
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
	if s[0] == '/' {
		// regular expression address
	} else if s[0] >= '0' && s[0] <= '9' {
		// numeric line address
		addr := new(address);
		var err os.Error;
		s, addr.rangeStart, err = getNumberFromLine(s);
		if err != nil {
			return s, nil, err
		}
		addr.rangeEnd = addr.rangeStart;
		if s[0] == ',' {
			s = s[1:];
			if len(s) > 0 && s[0] >= '0' && s[0] <= '9' {
				s, addr.rangeEnd, err = getNumberFromLine(s);
				if err != nil {
					return s, nil, err
				}
			} else {
				addr.rangeEnd = 0	// to end of file
			}
		}
		return s, addr, nil;
	}
	return s, nil, nil;
}

func (s *Sed) lineMatchesAddress(addr *address) bool {
	if addr != nil {
		if addr.rangeEnd == 0 {
			if s.lineNumber >= addr.rangeStart {
				return true
			}
		} else if s.lineNumber >= addr.rangeStart && s.lineNumber <= addr.rangeEnd {
			return true
		}
		if addr.lastLine && s.lineNumber == len(s.inputLines) {
			return true
		}
		if addr.regex != nil {
			return addr.regex.Match(s.patternSpace)
		}
		return false;
	}
	return true;
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
