//
//  sed.go
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
  "flag";
  "fmt";
  "io/ioutil";
  "os";
  "strings";
  "container/vector";
)

const (
  versionMajor = 0;
  versionMinor = 1;
  versionPoint = 0;
)

var versionString string

func init() {
  versionString = fmt.Sprintf("%d.%d.%d", versionMajor, versionMinor, versionPoint)
}

var show_version = flag.Bool("version", false, "Show version information.")
var show_help = flag.Bool("h", false, "Show help information.")
var quiet = flag.Bool("n", false, "Don't print the pattern space at the end of each script cycle.")
var script = flag.String("e", "", "The script used to process the input file.")
var script_file = flag.String("f", "", "Specify a file to read as the script. Ignored if -e present")
var edit_inplace = flag.Bool("i", false, "This option specifies that files are to be edited in-place. Otherwise output is printed to stdout.")
var line_wrap = flag.Uint("l", 70, "Specify the default line-wrap length for the l command. A length of 0 (zero) means to never wrap long lines. If not specified, it is taken to be 70.")
var unbuffered = flag.Bool("u", false, "Buffer both input and output as minimally as practical. (ignored)")
var treat_files_as_seperate = flag.Bool("s", false, "Treat files as searate entites. Line numbers reset to 1 for each file")

var usageShown bool = false

type Sed struct {
  inputLines              []string;
  commands                *vector.Vector;
  outputFile              *os.File;
  patternSpace, holdSpace string;
  lineNumber              int;
}

func (s *Sed) Init() {
  s.commands = new(vector.Vector);
  s.outputFile = os.Stdout;
  s.patternSpace = "";
  s.holdSpace = "";
}

func usage() {
  // only show usage once.
  if !usageShown {
    usageShown = true;
    fmt.Fprint(os.Stdout, "sed [options] [script] input_file\n\n");
    flag.PrintDefaults();
  }
}

var inputFilename string

func (s *Sed) readInputFile() {
  b, err := ioutil.ReadFile(inputFilename);
  if err != nil {
    fmt.Fprintf(os.Stderr, "Error reading input file %s\n", inputFilename);
    os.Exit(-1);
  }
  s.inputLines = strings.Split(string(b), "\n", 0);
}

func (s *Sed) parseScript() (err os.Error) {
  // a script may be a single command or it may be several
  scriptLines := strings.Split(*script, "\n", 0);
  for idx, line := range scriptLines {
    line = strings.TrimSpace(line);
    if strings.HasPrefix(line, "#") || len(line) == 0 {
      // comment
      continue
    }
    // this isn't really right. There may be slashes in the regular expression
    pieces := strings.Split(line, "/", 0);
    c, err := NewCmd(pieces);
    if err != nil {
      fmt.Printf("%v line %d: %s\n", err, idx+1, line);
      os.Exit(-1);
    }
    s.commands.Push(c);
  }
  return nil;
}

func (s *Sed) printPatternSpace() {
  l := len(s.patternSpace);
  if *line_wrap <= 0 || l < int(*line_wrap) {
    fmt.Fprintf(s.outputFile, "%s\n", s.patternSpace)
  } else {
    // print the line in segments
    for i := 0; i < l; i += int(*line_wrap) {
      endOfLine := i + int(*line_wrap);
      if endOfLine > l {
        endOfLine = l
      }
      fmt.Fprintf(s.outputFile, "%s\n", s.patternSpace[i:endOfLine]);
    }
  }
}

func (s *Sed) process() {
  if *treat_files_as_seperate || *edit_inplace {
    s.lineNumber = 0
  }
  for _, s.patternSpace = range s.inputLines {
    // track line number starting with line 1
    s.lineNumber++;
    for c := range s.commands.Iter() {
      // println("cmd: ", c.(fmt.Stringer).String());
      if s.lineMatchesAddress(c.(Cmd).getAddress()) {
        stop, err := c.(Cmd).processLine(s);
        if err != nil {
          fmt.Printf("%v\n", err);
          os.Exit(-1);
        }
        if stop {
          break
        }
      }
    }
    if !*quiet {
      s.printPatternSpace()
    }
  }
}

func Main() {
  s := new(Sed);
  s.Init();
  flag.Parse();
  if *show_version {
    fmt.Fprintf(os.Stdout, "Version: %s (c)2009 Geoffrey Clements All Rights Reserved\n\n", versionString)
  }
  if *show_help {
    usage();
    return;
  }

  // the first parameter may be a script or an input file. This helps us track which
  currentFileParameter := 0;

  // we need a script
  if len(*script) == 0 {
    // no -e so try -f
    if len(*script_file) > 0 {
      b, err := ioutil.ReadFile(*script_file);
      if err != nil {
        fmt.Fprintf(os.Stderr, "Error reading script file %s\n", *script_file);
        os.Exit(-1);
      }
      s := string(b);
      script = &s;
    } else if flag.NArg() > 1 {
      s := flag.Arg(0);
      script = &s;
      // first parameter was the script so move to second parameter
      currentFileParameter++;
    }
  }

  // if script still isn't set we are screwed, exit.
  if len(*script) == 0 {
    fmt.Fprint(os.Stderr, "No script found.\n\n");
    usage();
    os.Exit(-1);
  }

  // parse script
  s.parseScript();

  if currentFileParameter >= flag.NArg() {
    fmt.Fprint(os.Stderr, "No input file specified.\n\n");
    usage();
    os.Exit(-1);
  }

  for ; currentFileParameter < flag.NArg(); currentFileParameter++ {
    inputFilename = flag.Arg(currentFileParameter);
    // actually do the processing
    s.readInputFile();
    if *edit_inplace {
      dir, err := os.Stat(inputFilename);
      if err != nil {
        fmt.Fprintf(os.Stderr, "Error getting information about input file: %s %v\n", err);
        os.Exit(-1);
      }
      f, err := os.Open(inputFilename, os.O_WRONLY|os.O_TRUNC, int(dir.Mode));
      if err != nil {
        fmt.Fprint(os.Stderr, "Error opening input file for inplace editing: %s %v\n", err);
        os.Exit(-1);
      }
      s.outputFile = f;
    }
    s.process();
    if *edit_inplace {
      s.outputFile.Close()
    }
  }
}
