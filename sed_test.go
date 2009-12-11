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

func TestNewCmd(t *testing.T) {
	pieces := []byte{'4', 'r', '5', 'o', '/', '0', '/', 'g'};
	c, err := NewCmd(nil, pieces);
	if c != nil {
		t.Error("2: Got a command when we shouldn't have " + c.String())
	}
	if err == nil {
		t.Error("Didn't get an error we expected")
	} else {
		checkString(t, "Expected Unknown script command", "Unknown script command", err.String())
	}

	// s
	pieces = []byte{'s', '/', 'o', '/', '0', '/', 'g'};
	c, err = NewCmd(nil, pieces);
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
	pieces := []byte{'d', '/', 'o', '/', '0', '/', 'g'};
	c, err := NewCmd(nil, pieces);
	dc := c.(*d_cmd);
	if dc != nil {
		t.Error("1: Got a command when we shouldn't have " + c.String())
	}
	if err == nil {
		t.Error("Didn't get an error we expected")
	} else {
		checkString(t, "Expected: Wrong number of parameters for command", "Wrong number of parameters for command", err.String())
	}

	pieces = []byte{'d', '/', 'd'};
	c, err = NewCmd(nil, pieces);
	dc = c.(*d_cmd);
	if dc != nil {
		t.Error("2: Got a command when we shouldn't have " + c.String())
	}
	if err == nil {
		t.Error("Didn't get an error we expected")
	} else {
		checkString(t, "Expected: Wrong number of parameters for command", "Wrong number of parameters for command", err.String())
	}

	pieces = []byte{'d'};
	c, err = NewCmd(nil, pieces);
	dc = c.(*d_cmd);
	if dc == nil {
		t.Error("Didn't get a d command that we expected")
	} else if err != nil {
		t.Error("Got an error we didn't expect: " + err.String())
	}

	pieces = []byte{'$', 'd'};
	c, err = NewCmd(nil, pieces);
	dc = c.(*d_cmd);
	if dc == nil {
		t.Error("Didn't get a d command that we expected")
	} else if err != nil {
		t.Error("Got an error we didn't expect: " + err.String())
	}

	pieces = []byte{'4', '5', '7', 'd'};
	c, err = NewCmd(nil, pieces);
	dc = c.(*d_cmd);
	if dc == nil {
		t.Error("Didn't get a d command that we expected")
	} else if err != nil {
		t.Error("Got an error we didn't expect: " + err.String())
	}
}

func TestNewNCmd(t *testing.T) {
	pieces := []byte{'n', '/', 'o', '/', '0', '/', 'g'};
	c, err := NewCmd(nil, pieces);
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

	pieces = []byte{'n', '/', 'd'};
	c, err = NewCmd(nil, pieces);
	nc = c.(*n_cmd);
	if nc != nil {
		t.Error("2: Got a command when we shouldn't have " + c.String())
	}
	if err == nil {
		t.Error("Didn't get an error we expected")
	} else {
		checkString(t, "Expected: Wrong number of parameters for command", "Wrong number of parameters for command", err.String())
	}

	pieces = []byte{'n'};
	c, err = NewCmd(nil, pieces);
	nc = c.(*n_cmd);
	if nc == nil {
		t.Error("Didn't get a n command that we expected")
	} else if err != nil {
		t.Error("Got an error we didn't expect: " + err.String())
	}

	pieces = []byte{'$', 'n'};
	c, err = NewCmd(nil, pieces);
	nc = c.(*n_cmd);
	if nc == nil {
		t.Error("Didn't get a d command that we expected")
	} else if err != nil {
		t.Error("Got an error we didn't expect: " + err.String())
	}

	pieces = []byte{'4', '5', '7', 'n'};
	c, err = NewCmd(nil, pieces);
	nc = c.(*n_cmd);
	if nc == nil {
		t.Error("Didn't get a n command that we expected")
	} else if err != nil {
		t.Error("Got an error we didn't expect: " + err.String())
	}
}

func TestNewPCmd(t *testing.T) {
	pieces := []byte{'P', '/', 'o', '/', '0', '/', 'g'};
	c, err := NewCmd(nil, pieces);
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

	pieces = []byte{'P', '/', 'd'};
	c, err = NewCmd(nil, pieces);
	pc = c.(*p_cmd);
	if pc != nil {
		t.Error("2: Got a command when we shouldn't have " + c.String())
	}
	if err == nil {
		t.Error("Didn't get an error we expected")
	} else {
		checkString(t, "Expected: Wrong number of parameters for command", "Wrong number of parameters for command", err.String())
	}

	pieces = []byte{'P'};
	c, err = NewCmd(nil, pieces);
	pc = c.(*p_cmd);
	if pc == nil {
		t.Error("Didn't get a p command that we expected")
	} else if err != nil {
		t.Error("Got an error we didn't expect: " + err.String())
	}

	pieces = []byte{'$', 'P'};
	c, err = NewCmd(nil, pieces);
	pc = c.(*p_cmd);
	if pc == nil {
		t.Error("Didn't get a p command that we expected")
	} else if err != nil {
		t.Error("Got an error we didn't expect: " + err.String())
	}

	pieces = []byte{'4', '5', '7', 'P'};
	c, err = NewCmd(nil, pieces);
	pc = c.(*p_cmd);
	if pc == nil {
		t.Error("Didn't get a p command that we expected")
	} else if err != nil {
		t.Error("Got an error we didn't expect: " + err.String())
	}
}

func TestNewQCmd(t *testing.T) {
	pieces := []byte{'q', '/', 'o', '/', '0', '/', 'g'};
	c, err := NewCmd(nil, pieces);
	qc := c.(*q_cmd);
	if qc != nil {
		t.Error("1: Got a command when we shouldn't have " + c.String())
	}
	if err == nil {
		t.Error("Didn't get an error we expected")
	} else {
		checkString(t, "Expected: Wrong number of parameters for command", "Wrong number of parameters for command", err.String())
	}

	pieces = []byte{'q', '/', 'q'};
	c, err = NewCmd(nil, pieces);
	qc = c.(*q_cmd);
	if qc != nil {
		t.Error("2: Got a command when we shouldn't have " + c.String())
	}
	if err == nil {
		t.Error("Didn't get an error we expected")
	} else {
		checkString(t, "Expected: parsing q: invalid argument", "parsing q: invalid argument", err.String())
	}

	pieces = []byte{'q'};
	c, err = NewCmd(nil, pieces);
	qc = c.(*q_cmd);
	if qc == nil {
		t.Error("Didn't get a q command that we expected")
	} else if err != nil {
		t.Error("Got an error we didn't expect: " + err.String())
	}

	pieces = []byte{'q', '/', '1'};
	c, err = NewCmd(nil, pieces);
	qc = c.(*q_cmd);
	if qc == nil {
		t.Error("Didn't get a q command that we expected")
	} else if err != nil {
		t.Error("Got an error we didn't expect: " + err.String())
	}

	pieces = []byte{'$', 'q'};
	c, err = NewCmd(nil, pieces);
	qc = c.(*q_cmd);
	if qc == nil {
		t.Error("Didn't get a q command that we expected")
	} else if err != nil {
		t.Error("Got an error we didn't expect: " + err.String())
	}

	pieces = []byte{'4', '5', '7', 'q'};
	c, err = NewCmd(nil, pieces);
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
	pieces := []byte{'s', '/', 'o', '/', '0', '/', 'g'};
	c, _ := NewCmd(nil, pieces);
	_s.patternSpace = []byte{'g', 'o', 'o', 'd'};
	stop, err := c.(Cmd).processLine(_s);
	if stop {
		t.Error("Got stop when we shouldn't have")
	}
	if err != nil {
		t.Errorf("Got and error when we shouldn't have %v", err)
	}
	checkString(t, "bad global s command", "g00d", string(_s.patternSpace));

	pieces = []byte{'s', '/', 'o', '/', '0', '/', '1'};
	c, _ = NewCmd(nil, pieces);
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
