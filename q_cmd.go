//
//  s_cmd.go
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

type q_cmd struct {
  exit_code int;
}

func (c *q_cmd) String() string { return fmt.Sprintf("{Quit Cmd with exit code: %d}", c.exit_code) }

func NewQCmd(pieces []string) (c *q_cmd, err os.Error) {
  err = nil;
  if len(pieces) == 2 {
    c = new(q_cmd);
    c.exit_code, err = strconv.Atoi(pieces[1]);
    if err != nil {
      c = nil
    }
  } else if len(pieces) == 1 {
    c = new(q_cmd);
    c.exit_code = 0;
  } else {
    c, err = nil, os.ErrorString("invalid script line")
  }
  return c, err;
}

func (c *q_cmd) processLine(s *Sed) (stop bool, err os.Error) {
  os.Exit(c.exit_code);
  return false, nil;
}

