//
//  s_cmd.go
//  sed
//
// Copyright (c) 2009 Geoffrey Clements
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "nthOccuranceSoftware"), to deal
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

const (
	global_replace = -1
)

type s_cmd struct {
	addr         *address
	regex        string
	replace      []byte
	nthOccurance int
	re           *regexp.Regexp
}

func (c *s_cmd) match(line []byte, lineNumber int) bool {
	return c.addr.match(line, lineNumber)
}

func (c *s_cmd) String() string {
	if c != nil {
		if c.addr != nil {
			return fmt.Sprintf("{s command addr:%s regex:%v replace:%s nth occurance:%d}", c.addr, c.regex, c.replace, c.nthOccurance)
		}
		return fmt.Sprintf("{s command regex:%v replace:%s nth occurance:%d}", c.regex, c.replace, c.nthOccurance)
	}
	return "{s command}"
}

func NewSCmd(pieces [][]byte, addr *address) (c *s_cmd, err os.Error) {
	if len(pieces) != 4 {
		return nil, WrongNumberOfCommandParameters
	}

	err = nil
	c = new(s_cmd)
	c.addr = addr

	c.regex = string(pieces[1])
	if len(c.regex) == 0 {
		return nil, RegularExpressionExpected
	}
	c.re, err = regexp.Compile(string(c.regex))
	if err != nil {
		return nil, err
	}

	c.replace = pieces[2]

	flag := string(pieces[3])
	if flag != "g" {
		c.nthOccurance = 1
		if len(flag) > 0 {
			c.nthOccurance, err = strconv.Atoi(flag)
			if err != nil {
				return nil, InvalidSCommandFlag
			}
		}
	} else {
		c.nthOccurance = global_replace
	}

	return c, err
}

func (c *s_cmd) processLine(s *Sed) (stop bool, err os.Error) {
	stop, err = false, nil

	switch c.nthOccurance {
	case global_replace:
		s.patternSpace = c.re.ReplaceAll(s.patternSpace, c.replace)
	default:
		// a numeric flag command
		count := 0
		line := s.patternSpace
		s.patternSpace = make([]byte, 0)
		for {
			matches := c.re.FindIndex(line)
			if len(matches) > 0 {
				count++
				if count == c.nthOccurance {
					s.patternSpace = bytes.Add(s.patternSpace, line[0:matches[0]])
					s.patternSpace = bytes.Add(s.patternSpace, c.replace)
					s.patternSpace = bytes.Add(s.patternSpace, line[matches[1]:])
					break
				} else {
					s.patternSpace = bytes.Add(s.patternSpace, line[0:matches[0]+1])
				}
				line = line[matches[0]+1:]
			} else {
				s.patternSpace = bytes.Add(s.patternSpace, line)
				break
			}
		}
	}
	return stop, err
}
