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
	"runtime";
	"testing";
)

func TestNewCmd(t *testing.T) {
	pieces := [][]byte{[]byte{'b', 'a', 'd'}, []byte{'o'}, []byte{'0'}, []byte{'g'}};
	c, err := NewCmd(pieces);
	if c != nil {
		_, file, line, ok := runtime.Caller(0);
		if ok {
			t.Errorf("%s:%d: Got a command when we shouldn't have %s", file, line, c.String())
		} else {
			t.Error("1: Got a command when we shouldn't have " + c.String())
		}
	}
	if err == nil {
		t.Error("Didn't get an error we expected")
	} else {
		checkString(t, "Expected Unknown script command", "Unknown script command", err.String())
	}

	pieces = [][]byte{[]byte{'4', 'r', '5'}, []byte{'o'}, []byte{'0'}, []byte{'g'}};
	c, err = NewCmd(pieces);
	if c != nil {
		t.Error("2: Got a command when we shouldn't have " + c.String())
	}
	if err == nil {
		t.Error("Didn't get an error we expected")
	} else {
		checkString(t, "Expected Unknown script command", "Unknown script command", err.String())
	}

	// s
	pieces = [][]byte{[]byte{'s'}, []byte{'o'}, []byte{'0'}, []byte{'g'}};
	c, err = NewCmd(pieces);
	if c.getAddress() != nil {
		t.Error("Got an address when we shouldn't have " + c.String())
	}
	sc := c.(*s_cmd);
	if sc == nil {
		t.Error("Didn't get a command that we expected")
	} else if sc.regex != "o" && len(sc.replace) == 1 && sc.replace[0] == '0' && sc.count == -1 {
		t.Error("We didn't get the s command we expected")
	} else if err != nil {
		t.Error("Got an error we didn't expect: " + err.String())
	}
}

func TestNewDCmd(t *testing.T) {
	pieces := [][]byte{[]byte{'d'}, []byte{'o'}, []byte{'0'}, []byte{'g'}};
	c, err := NewCmd(pieces);
	dc := c.(*d_cmd);
	if dc != nil {
		t.Error("1: Got a command when we shouldn't have " + c.String())
	}
	if err == nil {
		t.Error("Didn't get an error we expected")
	} else {
		checkString(t, "Expected: Wrong number of parameters for command", "Wrong number of parameters for command", err.String())
	}

	pieces = [][]byte{[]byte{'d'}, []byte{'d'}};
	c, err = NewCmd(pieces);
	dc = c.(*d_cmd);
	if dc != nil {
		t.Error("2: Got a command when we shouldn't have " + c.String())
	}
	if err == nil {
		t.Error("Didn't get an error we expected")
	} else {
		checkString(t, "Expected: Wrong number of parameters for command", "Wrong number of parameters for command", err.String())
	}

	pieces = [][]byte{[]byte{'d'}};
	c, err = NewCmd(pieces);
	if c.getAddress() != nil {
		t.Error("Got an address when we shouldn't have " + c.String())
	}
	dc = c.(*d_cmd);
	if dc == nil {
		t.Error("Didn't get a d command that we expected")
	} else if err != nil {
		t.Error("Got an error we didn't expect: " + err.String())
	}

	pieces = [][]byte{[]byte{'$'}, []byte{'d'}};
	c, err = NewCmd(pieces);
	if c.getAddress() == nil {
		t.Error("Got an address when we shouldn't have " + c.String())
	} else if a := c.getAddress(); a.lineNumber != -1 || !a.lastLine || a.regex != nil {
		t.Error("Did not get the address we expected: " + a.String())
	}
	dc = c.(*d_cmd);
	if dc == nil {
		t.Error("Didn't get a d command that we expected")
	} else if err != nil {
		t.Error("Got an error we didn't expect: " + err.String())
	}

	pieces = [][]byte{[]byte{'4', '5', '7'}, []byte{'d'}};
	c, err = NewCmd(pieces);
	if c.getAddress() == nil {
		t.Error("Got an address when we shouldn't have " + c.String())
	} else if a := c.getAddress(); a.lineNumber != 457 || a.lastLine || a.regex != nil {
		t.Error("Did not get the address we expected: " + a.String())
	}
	dc = c.(*d_cmd);
	if dc == nil {
		t.Error("Didn't get a d command that we expected")
	} else if err != nil {
		t.Error("Got an error we didn't expect: " + err.String())
	}
}

func TestNewNCmd(t *testing.T) {
	pieces := [][]byte{[]byte{'n'}, []byte{'o'}, []byte{'0'}, []byte{'g'}};
	c, err := NewCmd(pieces);
	nc := c.(*n_cmd);
	if nc != nil {
		t.Error("1: Got a command when we shouldn't have " + c.String())
	}
	if err == nil {
		t.Error("Didn't get an error we expected")
	}
	if err == nil {
		t.Error("Didn't get an error we expected")
	} else {
		checkString(t, "Expected: Wrong number of parameters for command", "Wrong number of parameters for command", err.String())
	}

	pieces = [][]byte{[]byte{'n'}, []byte{'d'}};
	c, err = NewCmd(pieces);
	nc = c.(*n_cmd);
	if nc != nil {
		t.Error("2: Got a command when we shouldn't have " + c.String())
	}
	if err == nil {
		t.Error("Didn't get an error we expected")
	} else {
		checkString(t, "Expected: Wrong number of parameters for command", "Wrong number of parameters for command", err.String())
	}

	pieces = [][]byte{[]byte{'n'}};
	c, err = NewCmd(pieces);
	if c.getAddress() != nil {
		t.Error("Got an address when we shouldn't have " + c.String())
	}
	nc = c.(*n_cmd);
	if nc == nil {
		t.Error("Didn't get a n command that we expected")
	} else if err != nil {
		t.Error("Got an error we didn't expect: " + err.String())
	}

	pieces = [][]byte{[]byte{'$'}, []byte{'n'}};
	c, err = NewCmd(pieces);
	if c.getAddress() == nil {
		t.Error("Got an address when we shouldn't have " + c.String())
	} else if a := c.getAddress(); a.lineNumber != -1 || !a.lastLine || a.regex != nil {
		t.Error("Did not get the address we expected: " + a.String())
	}
	nc = c.(*n_cmd);
	if nc == nil {
		t.Error("Didn't get a d command that we expected")
	} else if err != nil {
		t.Error("Got an error we didn't expect: " + err.String())
	}

	pieces = [][]byte{[]byte{'4', '5', '7'}, []byte{'n'}};
	c, err = NewCmd(pieces);
	if c.getAddress() == nil {
		t.Error("Got an address when we shouldn't have " + c.String())
	} else if a := c.getAddress(); a.lineNumber != 457 || a.lastLine || a.regex != nil {
		t.Error("Did not get the address we expected: " + a.String())
	}
	nc = c.(*n_cmd);
	if nc == nil {
		t.Error("Didn't get a n command that we expected")
	} else if err != nil {
		t.Error("Got an error we didn't expect: " + err.String())
	}
}

func TestNewPCmd(t *testing.T) {
	pieces := [][]byte{[]byte{'P'}, []byte{'o'}, []byte{'0'}, []byte{'g'}};
	c, err := NewCmd(pieces);
	pc := c.(*p_cmd);
	if pc != nil {
		t.Error("1: Got a command when we shouldn't have " + c.String())
	}
	if err == nil {
		t.Error("Didn't get an error we expected")
	}
	if err == nil {
		t.Error("Didn't get an error we expected")
	} else {
		checkString(t, "Expected: Wrong number of parameters for command", "Wrong number of parameters for command", err.String())
	}

	pieces = [][]byte{[]byte{'P'}, []byte{'d'}};
	c, err = NewCmd(pieces);
	pc = c.(*p_cmd);
	if pc != nil {
		t.Error("2: Got a command when we shouldn't have " + c.String())
	}
	if err == nil {
		t.Error("Didn't get an error we expected")
	} else {
		checkString(t, "Expected: Wrong number of parameters for command", "Wrong number of parameters for command", err.String())
	}

	pieces = [][]byte{[]byte{'P'}};
	c, err = NewCmd(pieces);
	if c.getAddress() != nil {
		t.Error("Got an address when we shouldn't have " + c.String())
	}
	pc = c.(*p_cmd);
	if pc == nil {
		t.Error("Didn't get a p command that we expected")
	} else if err != nil {
		t.Error("Got an error we didn't expect: " + err.String())
	}

	pieces = [][]byte{[]byte{'$'}, []byte{'P'}};
	c, err = NewCmd(pieces);
	if c.getAddress() == nil {
		t.Error("Got an address when we shouldn't have " + c.String())
	} else if a := c.getAddress(); a.lineNumber != -1 || !a.lastLine || a.regex != nil {
		t.Error("Did not get the address we expected: " + a.String())
	}
	pc = c.(*p_cmd);
	if pc == nil {
		t.Error("Didn't get a p command that we expected")
	} else if err != nil {
		t.Error("Got an error we didn't expect: " + err.String())
	}

	pieces = [][]byte{[]byte{'4', '5', '7'}, []byte{'P'}};
	c, err = NewCmd(pieces);
	if c.getAddress() == nil {
		t.Error("Got an address when we shouldn't have " + c.String())
	} else if a := c.getAddress(); a.lineNumber != 457 || a.lastLine || a.regex != nil {
		t.Error("Did not get the address we expected: " + a.String())
	}
	pc = c.(*p_cmd);
	if pc == nil {
		t.Error("Didn't get a p command that we expected")
	} else if err != nil {
		t.Error("Got an error we didn't expect: " + err.String())
	}
}

func TestNewQCmd(t *testing.T) {
	pieces := [][]byte{[]byte{'q'}, []byte{'o'}, []byte{'0'}, []byte{'g'}};
	c, err := NewCmd(pieces);
	qc := c.(*q_cmd);
	if qc != nil {
		t.Error("1: Got a command when we shouldn't have " + c.String())
	}
	if err == nil {
		t.Error("Didn't get an error we expected")
	} else {
		checkString(t, "Expected: Wrong number of parameters for command", "Wrong number of parameters for command", err.String())
	}

	pieces = [][]byte{[]byte{'q'}, []byte{'q'}};
	c, err = NewCmd(pieces);
	qc = c.(*q_cmd);
	if qc != nil {
		t.Error("2: Got a command when we shouldn't have " + c.String())
	}
	if err == nil {
		t.Error("Didn't get an error we expected")
	} else {
		checkString(t, "Expected: parsing q: invalid argument", "parsing q: invalid argument", err.String())
	}

	pieces = [][]byte{[]byte{'q'}};
	c, err = NewCmd(pieces);
	if c.getAddress() != nil {
		t.Error("Got an address when we shouldn't have " + c.String())
	}
	qc = c.(*q_cmd);
	if qc == nil {
		t.Error("Didn't get a q command that we expected")
	} else if err != nil {
		t.Error("Got an error we didn't expect: " + err.String())
	}

	pieces = [][]byte{[]byte{'q'}, []byte{'1'}};
	c, err = NewCmd(pieces);
	if c.getAddress() != nil {
		t.Error("Got an address when we shouldn't have " + c.String())
	}
	qc = c.(*q_cmd);
	if qc == nil {
		t.Error("Didn't get a q command that we expected")
	} else if err != nil {
		t.Error("Got an error we didn't expect: " + err.String())
	}

	pieces = [][]byte{[]byte{'$'}, []byte{'q'}};
	c, err = NewCmd(pieces);
	if c.getAddress() == nil {
		t.Error("Got an address when we shouldn't have " + c.String())
	} else if a := c.getAddress(); a.lineNumber != -1 || !a.lastLine || a.regex != nil {
		t.Error("Did not get the address we expected: " + a.String())
	}
	qc = c.(*q_cmd);
	if qc == nil {
		t.Error("Didn't get a q command that we expected")
	} else if err != nil {
		t.Error("Got an error we didn't expect: " + err.String())
	}

	pieces = [][]byte{[]byte{'4', '5', '7'}, []byte{'q'}};
	c, err = NewCmd(pieces);
	if c.getAddress() == nil {
		t.Error("Got an address when we shouldn't have " + c.String())
	} else if a := c.getAddress(); a.lineNumber != 457 || a.lastLine || a.regex != nil {
		t.Error("Did not get the address we expected: " + a.String())
	}
	qc = c.(*q_cmd);
	if qc == nil {
		t.Error("Didn't get a d command that we expected")
	} else if err != nil {
		t.Error("Got an error we didn't expect: " + err.String())
	}
}

func TestProcessLine(t *testing.T) {
	_s := new(Sed);
	_s.Init();
	pieces := [][]byte{[]byte{'s'}, []byte{'o'}, []byte{'0'}, []byte{'g'}};
	c, _ := NewCmd(pieces);
	_s.patternSpace = []byte{'g', 'o', 'o', 'd'};
	stop, err := c.(Cmd).processLine(_s);
	if stop {
		t.Error("Got stop when we shouldn't have")
	}
	if err != nil {
		t.Errorf("Got and error when we shouldn't have %v", err)
	}
	checkString(t, "bad global s command", "g00d", string(_s.patternSpace));

	pieces = [][]byte{[]byte{'s'}, []byte{'o'}, []byte{'0'}, []byte{'1'}};
	c, _ = NewCmd(pieces);
	_s.patternSpace = []byte{'g', 'o', 'o', 'd'};
	stop, err = c.(Cmd).processLine(_s);
	if stop {
		t.Error("Got stop when we shouldn't have")
	}
	if err != nil {
		t.Errorf("Got and error when we shouldn't have %v", err)
	}
	checkString(t, "bad global s command", "g0od", string(_s.patternSpace));
}

func checkInt(t *testing.T, val, expected int, actual string) {
	if expected != val {
		t.Errorf("%s: '%d' != '%d'", val, expected, actual)
	}
}

func checkString(t *testing.T, message, expected, actual string) {
	if expected != actual {
		t.Errorf("%s: '%s' != '%s'", message, expected, actual)
	}
}
