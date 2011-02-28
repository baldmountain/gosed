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
	"bytes"
	"fmt"
	"os"
)

type a_cmd struct {
	addr *address
	text []byte
}

func (c *a_cmd) match(line []byte, lineNumber int) bool {
	return c.addr.match(line, lineNumber)
}

func (c *a_cmd) String() string {
	if c != nil {
		if c.addr != nil {
			return fmt.Sprintf("{a command addr:%s text:%s}", c.addr.String(), c.text)
		}
		return fmt.Sprintf("{a command text:%s}", c.text)
	}
	return fmt.Sprintf("{a command}")
}

func (c *a_cmd) processLine(s *Sed) (bool, os.Error) {
	return false, nil
}

func NewACmd(s *Sed, line []byte, addr *address) (*a_cmd, os.Error) {
	cmd := new(a_cmd)
	cmd.addr = addr
	cmd.text = line[1:]
	for bytes.HasSuffix(cmd.text, []byte{'\\'}) {
		cmd.text = cmd.text[0 : len(cmd.text)-1]
		line, err := s.getNextScriptLine()
		if err != nil {
			break
		}
		buf := bytes.NewBuffer(cmd.text)
		buf.WriteRune('\n')
		buf.Write(line)
		cmd.text = buf.Bytes()
	}
	cmd.text = trimSpaceFromBeginning(cmd.text)
	return cmd, nil
}
