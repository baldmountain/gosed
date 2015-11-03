//
//  p_cmd.go
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
)

type p_cmd struct {
	addr        *address
	upToNewLine bool
}

func (c *p_cmd) match(line []byte, lineNumber int) bool {
	return c.addr.match(line, lineNumber)
}

func (c *p_cmd) String() string {
	if c != nil && c.addr != nil {
		return fmt.Sprintf("{p command addr:%s}", c.addr.String())
	}
	return fmt.Sprint("{p command}")
}

func (c *p_cmd) processLine(s *Sed) (bool, error) {
	// print output space
	if c.upToNewLine {
		firstLine := bytes.SplitN(s.patternSpace, []byte{'\n'}, 1)[0]
		fmt.Fprintln(s.outputFile, string(firstLine))
	} else {
		fmt.Fprintln(s.outputFile, string(s.patternSpace))
	}
	return false, nil
}

func NewPCmd(pieces [][]byte, addr *address) (*p_cmd, error) {
	if len(pieces) > 1 {
		return nil, WrongNumberOfCommandParameters
	}
	cmd := new(p_cmd)
	cmd.addr = addr
	cmd.upToNewLine = pieces[0][0] == 'P'
	return cmd, nil
}
