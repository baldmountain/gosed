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
	"os";
	"fmt";
	"strconv";
)

type eql_cmd struct {
	command;
}

func (c *eql_cmd) String() string {
	if c != nil && c.addr != nil {
		return fmt.Sprint("{Output current line number}", c.addr)
	}
	return fmt.Sprint("{Output current line number Cmd}");
}

func (c *eql_cmd) processLine(s *Sed) (bool, os.Error) {
	fmt.Fprint(s.outputFile, strconv.Itoa(s.lineNumber));
	return false, nil;
}

func (c *eql_cmd) getAddress() *address	{ return c.addr }

func NewEqlCmd(pieces [][]byte, addr *address) (*eql_cmd, os.Error) {
	if len(pieces) > 1 {
		return nil, WrongNumberOfCommandParameters
	}
	cmd := new(eql_cmd);
	cmd.addr = addr;
	return cmd, nil;
}
