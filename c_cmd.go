//
//  c_cmd.go
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
)

type c_cmd struct {
	addr	*address
	text	[]byte
}

func (c *c_cmd) match(line []byte, lineNumber int) bool {
	return c.addr.match(line, lineNumber)
}

func (c *c_cmd) String() string {
	if c != nil {
		if c.addr != nil {
			return fmt.Sprintf("{Append Cmd addr:%s text:%s}", c.addr.String(), c.text)
		}
		return fmt.Sprintf("{Append Cmd text:%s}", c.text)
	}
	return fmt.Sprintf("{Append Cmd}")
}

func (c *c_cmd) printText(s *Sed) {
	// we are going to get the newline from the
	fmt.Fprint(s.outputFile, string(c.text))
}

func (c *c_cmd) processLine(s *Sed) (bool, os.Error) {
  s.patternSpace = s.patternSpace[0:0]
	if c.addr != nil {
	  switch c.addr.address_type {
	    case ADDRESS_RANGE:
  	    if s.lineNumber+1 == c.addr.rangeEnd {
					c.printText(s);
					return true, nil
  	    }
	    case ADDRESS_LINE, ADDRESS_REGEX, ADDRESS_LAST_LINE:
				c.printText(s);
				return true, nil
	    case ADDRESS_TO_END_OF_FILE:
	      // FIX need to output at end of file
    }
	} else {
		c.printText(s);
		return true, nil
	}
	return false, nil
}

func NewCCmd(s *Sed, line []byte, addr *address) (*c_cmd, os.Error) {
	cmd := new(c_cmd)
	cmd.addr = addr
	cmd.text = line[1:]
	for bytes.HasSuffix(cmd.text, []byte{'\\'}) {
		cmd.text = cmd.text[0 : len(cmd.text)-1]
		line, err := s.getNextScriptLine()
		if err != nil {
			break
		}
		cmd.text = bytes.AddByte(cmd.text, '\n')
		cmd.text = bytes.Add(cmd.text, line)
	}
	return cmd, nil
}
