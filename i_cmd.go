//
//  i_cmd.go
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

type i_cmd struct {
	command;
	text	[]byte;
}

func (c *i_cmd) String() string {
	if c != nil {
		if c.addr != nil {
			return fmt.Sprintf("{Insert Cmd addr:%v text:%s}", c.addr, c.text)
		}
		return fmt.Sprintf("{Insert Cmd text:%s}", c.text);
	}
	return fmt.Sprintf("{Insert Cmd}");
}

func (c *i_cmd) processLine(s *Sed) (bool, os.Error) {
	b := bytes.NewBuffer(nil);
	b.Write(c.text);
	b.Write(s.patternSpace);
	s.patternSpace = b.Bytes();
	return false, nil;
}

func (c *i_cmd) getAddress() *address	{ return c.addr }

func NewICmd(s *Sed, pieces [][]byte, addr *address) (*i_cmd, os.Error) {
	if len(pieces) != 2 {
		return nil, WrongNumberOfCommandParameters
	}
	cmd := new(i_cmd);
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
