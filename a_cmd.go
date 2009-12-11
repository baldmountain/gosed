//
//  a_cmd.go
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
)

type a_cmd struct {
	addr	*address;
	text	[]byte;
}

func (c *a_cmd) match(line []byte, lineNumber, totalNumberOfLines int) bool {
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

func (c *a_cmd) String() string {
	if c != nil {
		if c.addr != nil {
			return fmt.Sprintf("{Append Cmd addr:%v text:%s}", c.addr, c.text)
		}
		return fmt.Sprintf("{Append Cmd text:%s}", c.text);
	}
	return fmt.Sprintf("{Append Cmd}");
}

func (c *a_cmd) processLine(s *Sed) (bool, os.Error) {
	s.patternSpace = bytes.Add(s.patternSpace, c.text);
	return false, nil;
}

func NewACmd(s *Sed, pieces [][]byte, addr *address) (*a_cmd, os.Error) {
	if len(pieces) != 2 {
		return nil, WrongNumberOfCommandParameters
	}
	cmd := new(a_cmd);
	cmd.addr = addr;
	cmd.text = pieces[1];
	for bytes.HasSuffix(cmd.text, []byte{'\\'}) {
		cmd.text = cmd.text[0 : len(cmd.text)-1];
		line, err := s.getNextScriptLine();
		if err != nil {
			break
		}
		cmd.text = bytes.AddByte(cmd.text, '\n');
		cmd.text = bytes.Add(cmd.text, line);
	}
	return cmd, nil;
}
