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
)

type Cmd interface {
  processLine(s *Sed) (stop bool, err os.Error);
}

func NewCmd(pieces []string) (Cmd, err os.Error) {
  err = nil;
  switch pieces[0] {
  case "s":
    return NewSCmd(pieces)
  case "q":
    return NewQCmd(pieces)
  case "d":
    return NewDCmd(pieces)
  case "P":
    return NewPCmd(pieces)
  case "n":
    return NewNCmd(pieces)
  }
  return nil, os.ErrorString("unknown script command");
}
