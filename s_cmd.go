//
//  s_cmd.go
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

const (
	global_replace = -1;
)

type s_cmd struct {
	addr	*address;
	regex	string;
	replace	[]byte;
	count	int;
	re	*regexp.Regexp;
}

func (c *s_cmd) match(line []byte, lineNumber, totalNumberOfLines int) bool {
	if c.addr != nil {
		if c.addr.rangeEnd == 0 {
			if lineNumber >= c.addr.rangeStart {
				return true
			}
		} else if lineNumber >= c.addr.rangeStart && lineNumber <= c.addr.rangeEnd {
			return true
		}
		if c.addr.lastLine && lineNumber == totalNumberOfLines {
			return true
		}
		if c.addr.regex != nil {
			return c.addr.regex.Match(line)
		}
		return false;
	}
	return true;
}

func (c *s_cmd) String() string {
	if c != nil {
		if c.addr != nil {
			return fmt.Sprintf("{Substitue Cmd regex:%s replace:%s count:%d addr:%v}", c.regex, c.replace, c.count, c.addr)
		}
		return fmt.Sprintf("{Substitue Cmd regex:%s replace:%s count:%d}", c.regex, c.replace, c.count);
	}
	return "{Substitue Cmd}";
}

func NewSCmd(pieces [][]byte, addr *address) (c *s_cmd, err os.Error) {
	if len(pieces) != 4 {
		return nil, WrongNumberOfCommandParameters
	}

	err = nil;
	c = new(s_cmd);
	c.addr = addr;

	c.regex = string(pieces[1]);
	if len(c.regex) == 0 {
		return nil, RegularExpressionExpected
	}
	c.re, err = regexp.Compile(string(c.regex));
	if err != nil {
		return nil, err
	}

	c.replace = pieces[2];

	flag := string(pieces[3]);
	if flag != "g" {
		c.count = 1;
		if len(flag) > 0 {
			c.count, err = strconv.Atoi(flag);
			if err != nil {
				return nil, InvalidSCommandFlag
			}
		}
	} else {
		c.count = global_replace
	}

	return c, err;
}

func (c *s_cmd) processLine(s *Sed) (stop bool, err os.Error) {
	stop, err = false, nil;

	switch c.count {
	case global_replace:
		s.patternSpace = c.re.ReplaceAll(s.patternSpace, c.replace)
	default:
		// a numeric flag command
		count := c.count;
		line := s.patternSpace;
		s.patternSpace = make([]byte, 0);
		for {
			if count <= 0 {
				s.patternSpace = bytes.Add(s.patternSpace, line);
				return;
			}
			matches := c.re.Execute(line);
			if len(matches) == 0 {
				s.patternSpace = bytes.Add(s.patternSpace, line);
				return;
			}
			s.patternSpace = bytes.Add(s.patternSpace, line[0:matches[0]]);
			s.patternSpace = bytes.Add(s.patternSpace, c.replace);
			line = line[matches[1]:];
			count--;
		}
	}
	return stop, err;
}
