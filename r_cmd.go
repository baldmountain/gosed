//
//  r_cmd.go
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
	"fmt"
	"os"
)

type r_cmd struct {
	addr *address
	text []byte
}

func (c *r_cmd) match(line []byte, lineNumber int) bool {
	return c.addr.match(line, lineNumber)
}

func (c *r_cmd) String() string {
	if c != nil && c.addr != nil {
		return fmt.Sprintf("{r command addr:%s}", c.addr.String())
	}
	return fmt.Sprint("{r command}")
}

func (c *r_cmd) processLine(s *Sed) (bool, os.Error) {
	// print output space
	if c.text != nil {
		s.outputFile.Write(c.text)
	}
	return false, nil
}

func NewRCmd(line []byte, addr *address) (*r_cmd, os.Error) {
	line = line[1:]
	cmd := new(r_cmd)
	cmd.addr = addr
	if len(line) > 0 {
		cmd.text = line
	} else {
		cmd.text = nil
	}
	return cmd, nil
}
