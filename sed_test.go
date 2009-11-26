package sed

import (
	"testing";
)

func TestProcessLine(t *testing.T) {
  pieces := [...]string{"s", "o", "0", "g"};
  c, _ := NewCmd(pieces[0:len(pieces)]);
	s,stop,err := c.processLine("good");
	if stop { t.Error("Got stop when we shouldn't have") }
	if err != nil { t.Errorf("Got and error when we shouldn't have %v", err) }
	checkString(t, "bad global s command", "g00d", s);

  pieces = [...]string{"s", "o", "0", "1"};
  c, _ = NewCmd(pieces[0:len(pieces)]);
	s,stop,err = c.processLine("good");
	if stop { t.Error("Got stop when we shouldn't have") }
	if err != nil { t.Errorf("Got and error when we shouldn't have %v", err) }
	checkString(t, "bad global s command", "g0od", s);
}

func checkString(t *testing.T, message, expected, actual string) {
	if expected != actual {
		t.Errorf("%s: '%s' != '%s'", message, expected, actual)
	}
}
