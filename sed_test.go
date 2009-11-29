//
//  sed_test.go
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
  "testing";
)

func TestProcessLine(t *testing.T) {
  _s := new(Sed);
  _s.Init();
  pieces := [...]string{"s", "o", "0", "g"};
  c, _ := NewCmd(pieces[0:len(pieces)]);
  _s.patternSpace = "good";
  stop, err := c.(Cmd).processLine(_s);
  if stop {
    t.Error("Got stop when we shouldn't have")
  }
  if err != nil {
    t.Errorf("Got and error when we shouldn't have %v", err)
  }
  checkString(t, "bad global s command", "g00d", _s.patternSpace);

  pieces = [...]string{"s", "o", "0", "1"};
  c, _ = NewCmd(pieces[0:len(pieces)]);
  _s.patternSpace = "good";
  stop, err = c.(Cmd).processLine(_s);
  if stop {
    t.Error("Got stop when we shouldn't have")
  }
  if err != nil {
    t.Errorf("Got and error when we shouldn't have %v", err)
  }
  checkString(t, "bad global s command", "g0od", _s.patternSpace);
}

func checkString(t *testing.T, message, expected, actual string) {
  if expected != actual {
    t.Errorf("%s: '%s' != '%s'", message, expected, actual)
  }
}
