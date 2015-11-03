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
	"bytes"
	"fmt"
)

type i_cmd struct {
	addr *address
	text []byte
}

func (c *i_cmd) match(line []byte, lineNumber int) bool {
	return c.addr.match(line, lineNumber)
}

func (c *i_cmd) String() string {
	if c != nil {
		if c.addr != nil {
			return fmt.Sprintf("{i command addr:%s text:%s}", c.addr.String(), string(c.text))
		}
		return fmt.Sprintf("{i command text:%s}", string(c.text))
	}
	return fmt.Sprintf("{i command}")
}

func (c *i_cmd) processLine(s *Sed) (bool, error) {
	return false, nil
}

func NewICmd(s *Sed, line []byte, addr *address) (*i_cmd, error) {
	cmd := new(i_cmd)
	cmd.addr = addr
	cmd.text = line[1:]
	for bytes.HasSuffix(cmd.text, []byte{'\\'}) {
		cmd.text = cmd.text[0 : len(cmd.text)-1]
		line, err := s.getNextScriptLine()
		if err != nil {
			break
		}
		// cmd.text = bytes.AddByte(cmd.text, '\n')
		buf := bytes.NewBuffer(cmd.text)
		buf.Write(line)
		s.patternSpace = buf.Bytes()
	}
	return cmd, nil
}
