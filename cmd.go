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
  "strconv";
  "regexp";
  "fmt";
)

type Cmd interface {
  fmt.Stringer;
  getAddress() *address;
  processLine(s *Sed) (stop bool, err os.Error);
}

type address struct {
  lineNumber int;
  lastLine   bool;
  regex      *regexp.Regexp;
}

func (a *address) String() string {
  return fmt.Sprintf("address{lineNumber:%d lastLine:%t regex:%v}", a.lineNumber, a.lastLine, a.regex)
}


type command struct {
  addr *address;
}

// A nil address means match any line
func checkForAddress(s string) *address {
  if s == "$" {
    return &address{-1, true, nil}
  }
  if ln, ok := strconv.Atoi(s); ok == nil {
    return &address{ln, false, nil}
  }
  return nil;
}

func (s *Sed) lineMatchesAddress(addr *address) bool {
  if addr != nil {
    if s.lineNumber == addr.lineNumber {
      return true
    }
    if addr.lastLine && s.lineNumber == len(s.inputLines) {
      return true
    }
    if addr.regex != nil {
      return addr.regex.MatchString(s.patternSpace)
    }
    return false;
  }
  return true;
}

func NewCmd(pieces []string) (Cmd, os.Error) {
  retryOnce := true;

  addr := checkForAddress(pieces[0]);
  if addr != nil {
    pieces = pieces[1:]
  }
retry:
  if retryOnce {
    switch pieces[0] {
    case "s":
      return NewSCmd(pieces, addr)
    case "q":
      return NewQCmd(pieces, addr)
    case "d":
      return NewDCmd(pieces, addr)
    case "P":
      return NewPCmd(pieces, addr)
    case "n":
      return NewNCmd(pieces, addr)
    }
    if re, ok := regexp.Compile(pieces[0]); ok == nil {
      pieces = pieces[1:];
      addr = &address{-1, false, re};
      retryOnce = false;
      goto retry;
    }
  }

  return nil, os.ErrorString("unknown script command");
}
