//
//  b_cmd.go
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

type b_cmd struct {
	addr  *address
	label string
}

func (c *b_cmd) match(line []byte, lineNumber int) bool {
	return c.addr.match(line, lineNumber)
}

func (c *b_cmd) String() string {
	if c != nil {
		if c.addr != nil {
			return fmt.Sprintf("{b command label: %s Cmd addr:%s}", c.label, c.addr.String())
		}
		return fmt.Sprintf("{b command label: %s Cmd}", c.label)
	}
	return fmt.Sprintf("{b command}")
}

func (c *b_cmd) processLine(s *Sed) (bool, error) {
	return true, NotImplemented
}

func NewBCmd(pieces [][]byte, addr *address) (*b_cmd, error) {
	if len(pieces) != 1 {
		return nil, WrongNumberOfCommandParameters
	}
	cmd := new(b_cmd)
	cmd.addr = addr
	cmd.label = string(bytes.TrimSpace(pieces[0][1:]))
	return cmd, nil
}
