//
//  g_cmd.go
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

type g_cmd struct {
	addr	*address;
	replace	bool;
}

func (c *g_cmd) match(line []byte, lineNumber, totalNumberOfLines int) bool {
	return c.addr.match(line, lineNumber, totalNumberOfLines)
}

func (c *g_cmd) String() string {
	if c != nil && c.addr != nil {
		if c.replace {
			return fmt.Sprint("{Replace pattern space with contents of hold space Cmd addr:%v}", c.addr)
		} else {
			return fmt.Sprint("{Append a newline and the hold space to the pattern space Cmd addr:%v}", c.addr)
		}
	}
	return fmt.Sprint("{Append/Replace pattern space with contents of hold space}");
}

func (c *g_cmd) processLine(s *Sed) (bool, os.Error) {
	if c.replace {
		s.patternSpace = copyByteSlice(s.holdSpace)
	} else {
		s.patternSpace = bytes.AddByte(s.patternSpace, '\n');
		s.patternSpace = bytes.Add(s.patternSpace, s.holdSpace);
	}
	return false, nil;
}

func NewGCmd(pieces [][]byte, addr *address) (*g_cmd, os.Error) {
	if len(pieces) > 1 {
		return nil, WrongNumberOfCommandParameters
	}
	cmd := new(g_cmd);
	if pieces[0][0] == 'g' {
		cmd.replace = true
	}
	cmd.addr = addr;
	return cmd, nil;
}
