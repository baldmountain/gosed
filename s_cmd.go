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
  "regexp";
  "strconv";
)

type s_cmd struct {
  command;
  regex   string;
  replace string;
  flag    string;
  count   int;
  re      *regexp.Regexp;
}

func (c *s_cmd) String() string {
  if c.addr != nil {
    return fmt.Sprintf("{Substitue Cmd regex:%s replace:%s flag:%s addr:%v}", c.regex, c.replace, c.flag, c.addr)
  }
  return fmt.Sprintf("{Substitue Cmd regex:%s replace:%s flag:%s}", c.regex, c.replace, c.flag)
}

func NewSCmd(pieces []string, addr *address) (c *s_cmd, err os.Error) {
  if len(pieces) != 4 {
    return nil, os.ErrorString("invalid script line")
  }

  err = nil;
  c = new(s_cmd);
  c.addr = addr;
  
  c.regex = pieces[1];
  if len(c.regex) == 0 {
    return nil, os.ErrorString("Regular expression in s command can't be zero length.")
  }
  c.re, err = regexp.Compile(c.regex);
  if err != nil {
    return nil, err
  }

  c.replace = pieces[2];

  c.flag = pieces[3];
  if c.flag != "g" {
    c.count = 1;
    if len(c.flag) > 0 {
      c.count, err = strconv.Atoi(c.flag);
      if err != nil {
        return nil, os.ErrorString("Invalid flag for s command " + err.String())
      }
    }
  }

  return c, err;
}

func (c *s_cmd) getAddress()(*address) {
  return c.addr;
}

func (c *s_cmd) processLine(s *Sed) (stop bool, err os.Error) {
  stop, err = false, nil;

  switch c.flag {
  case "g":
    s.patternSpace = c.re.ReplaceAllString(s.patternSpace, c.replace)
  default:
    // a numeric flag command
    count := c.count;
    line := s.patternSpace;
    s.patternSpace = "";
    for {
      if count <= 0 {
        s.patternSpace += line;
        return;
      }
      matches := c.re.ExecuteString(line);
      if len(matches) == 0 {
        s.patternSpace += line;
        return;
      }
      s.patternSpace += line[0:matches[0]] + c.replace;
      line = line[matches[1]:];
      count--;
    }
  }
  return stop, err;
}
