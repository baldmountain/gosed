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
	"fmt";
	"os";
)

type n_cmd struct {
	command;
}

func (c *n_cmd) String() string {
	if c != nil && c.addr != nil {
		return fmt.Sprint("{Output pattern space and get next line Cmd addr:%v}", c.addr)
	}
	return fmt.Sprint("{Output pattern space and get next line Cmd}");
}

func (c *n_cmd) processLine(s *Sed) (bool, os.Error) {
	if !*quiet {
		s.printPatternSpace()
	}
	return true, nil;
}

func (c *n_cmd) getAddress() *address	{ return c.addr }

func NewNCmd(pieces [][]byte, addr *address) (*n_cmd, os.Error) {
	if len(pieces) > 1 {
		return nil, WrongNumberOfCommandParameters
	}
	cmd := new(n_cmd);
	cmd.addr = addr;
	return cmd, nil;
}
