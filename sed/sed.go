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
	"bufio"
	"bytes"
	"container/list"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"unicode"
	"unicode/utf8"
)

const (
	versionMajor = 0
	versionMinor = 2
	versionPoint = 1
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
var line_wrap = flag.Uint("l", 0, "Specify the default line-wrap length for the l command. A length of 0 (zero) means to never wrap long lines. If not specified, it is taken to be 70.")
var unbuffered = flag.Bool("u", false, "Buffer both input and output as minimally as practical. (ignored)")
var treat_files_as_seperate = flag.Bool("s", false, "Treat files as searate entites. Line numbers reset to 1 for each file")

var usageShown bool = false

var newLine = []byte{'\n'}

type Sed struct {
	inputFile               *os.File
	input                   *bufio.Reader
	lineNumber              int
	currentLine             string
	beforeCommands          *list.List
	commands                *list.List
	afterCommands           *list.List
	outputFile              *os.File
	patternSpace, holdSpace []byte
	scriptLines             [][]byte
	scriptLineNumber        int
}

func (s *Sed) Init() {
	s.beforeCommands = new(list.List)
	s.commands = new(list.List)
	s.afterCommands = new(list.List)
	s.outputFile = os.Stdout
	s.patternSpace = make([]byte, 0)
	s.holdSpace = make([]byte, 0)
}

func copyByteSlice(a []byte) []byte {
	newSlice := make([]byte, len(a))
	copy(newSlice, a)
	return newSlice
}

func usage() {
	// only show usage once.
	if !usageShown {
		usageShown = true
		fmt.Fprint(os.Stdout, "sed [options] [script] input_file\n\n")
		flag.PrintDefaults()
	}
}

var inputFilename string

func (s *Sed) getNextScriptLine() ([]byte, error) {
	if s.scriptLineNumber < len(s.scriptLines) {
		val := s.scriptLines[s.scriptLineNumber]
		s.scriptLineNumber++
		return val, nil
	}
	return nil, io.EOF
}

func trimSpaceFromBeginning(s []byte) []byte {
	start, end := 0, len(s)
	for start < end {
		r, wid := utf8.DecodeRune(s[start:end])
		if !unicode.IsSpace(r) {
			break
		}
		start += wid
	}
	return s[start:end]
}

func (s *Sed) parseScript(scriptBuffer []byte) (err error) {
	// a script may be a single command or it may be several
	s.scriptLines = bytes.Split(scriptBuffer, newLine)
	s.scriptLineNumber = 0
	var line []byte
	var serr error
	for line, serr = s.getNextScriptLine(); serr == nil; line, serr = s.getNextScriptLine() {
		// line = bytes.TrimSpace(line);
		line = trimSpaceFromBeginning(line)
		if len(line) == 0 {
			// zero length line
			continue
		}
		if line[0] == '#' {
			if s.scriptLineNumber == 1 && len(line) > 1 && line[1] == 'n' {
				// spcial case where the first 2 characters of the file are #n which is
				// equivalent to passing -n on the command line
				*quiet = true
			}
			continue
		}
		c, err := NewCmd(s, line)
		if err != nil {
			fmt.Printf("Script error: %s -> %d: %s\n", err.Error(), s.scriptLineNumber, line)
			os.Exit(-1)
		}
		if _, ok := c.(*i_cmd); ok {
			s.beforeCommands.PushBack(c)
		} else if _, ok := c.(*a_cmd); ok {
			s.afterCommands.PushBack(c)
		} else {
			s.commands.PushBack(c)
		}
	}
	return nil
}

func (s *Sed) printLine(line []byte) {
	l := len(line)
	if *line_wrap <= 0 || l < int(*line_wrap) {
		fmt.Fprintf(s.outputFile, "%s\n", line)
	} else {
		// print the line in segments
		for i := 0; i < l; i += int(*line_wrap) {
			endOfLine := i + int(*line_wrap)
			if endOfLine > l {
				endOfLine = l
			}
			fmt.Fprintf(s.outputFile, "%s\n", line[i:endOfLine])
		}
	}
}

func (s *Sed) printPatternSpace() {
	lines := bytes.Split(s.patternSpace, newLine)
	for _, line := range lines {
		s.printLine(line)
	}
}

func (s *Sed) process() {
	if *treat_files_as_seperate || *edit_inplace {
		s.lineNumber = 0
	}
	var err error
	s.patternSpace, err = s.input.ReadSlice('\n')
	for err != io.EOF {
		lineLength := len(s.patternSpace)
		if lineLength > 0 {
			s.patternSpace = s.patternSpace[0 : lineLength-1]
		}
		s.currentLine = string(s.patternSpace)
		// track line number starting with line 1
		s.lineNumber++
		stop := false
		// process i commands
		for c := s.beforeCommands.Front(); c != nil; c = c.Next() {
			// ask the sed if we should process this command, based on address
			if cmd, ok := c.Value.(*i_cmd); ok {
				if c.Value.(Address).match(s.patternSpace, s.lineNumber) {
					s.outputFile.Write(cmd.text)
				}
			}
		}
		for c := s.commands.Front(); c != nil; c = c.Next() {
			// ask the sed if we should process this command, based on address
			if c.Value.(Address).match(s.patternSpace, s.lineNumber) {
				var err error
				stop, err = c.Value.(Cmd).processLine(s)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
					fmt.Fprintf(os.Stderr, "Line: %d:%s\n", s.lineNumber, s.currentLine)
					fmt.Fprintf(os.Stderr, "Command: %s\n", c.Value.(Cmd).String())
					os.Exit(-1)
				}
				if stop {
					break
				}
			}
		}
		if !*quiet && !stop {
			s.printPatternSpace()
		}
		// process a commands
		for c := s.afterCommands.Front(); c != nil; c = c.Next() {
			// ask the sed if we should process this command, based on address
			if cmd, ok := c.Value.(*a_cmd); ok {
				if c.Value.(Address).match(s.patternSpace, s.lineNumber) {
					fmt.Fprintf(s.outputFile, "%s\n", cmd.text)
				}
			}
		}
		s.patternSpace, err = s.input.ReadSlice('\n')
	}
}

func Main() {
	var err error
	s := new(Sed)
	s.Init()
	flag.Parse()
	if *show_version {
		fmt.Fprintf(os.Stdout, "Version: %s (c)2009-2010 Geoffrey Clements All Rights Reserved\n\n", versionString)
	}
	if *show_help {
		usage()
		return
	}

	// the first parameter may be a script or an input file. This helps us track which
	currentFileParameter := 0
	var scriptBuffer []byte

	// we need a script
	if len(*script) == 0 {
		// no -e so try -f
		if len(*script_file) > 0 {
			sb, err := os.ReadFile(*script_file)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading script file %s\n", *script_file)
				os.Exit(-1)
			}
			scriptBuffer = sb
		} else if flag.NArg() > 1 {
			scriptBuffer = []byte(flag.Arg(0))

			// change semicoluns to newlines for scripts on command line
			idx := bytes.IndexByte(scriptBuffer, ';')
			for idx >= 0 {
				scriptBuffer[idx] = '\n'
				s := scriptBuffer[idx+1:]
				idx = bytes.IndexByte(s, ';')
			}
			// first parameter was the script so move to second parameter
			currentFileParameter++
		}
	} else {
		scriptBuffer = []byte(*script)
		// change semicoluns to newlines for scripts on command line
		idx := bytes.IndexByte(scriptBuffer, ';')
		for idx >= 0 {
			scriptBuffer[idx] = '\n'
			s := scriptBuffer[idx+1:]
			idx = bytes.IndexByte(s, ';')
		}
	}

	// if script still isn't set we are screwed, exit.
	if len(scriptBuffer) == 0 {
		fmt.Fprint(os.Stderr, "No script found.\n\n")
		usage()
		os.Exit(-1)
	}

	// parse script
	s.parseScript(scriptBuffer)

	if currentFileParameter >= flag.NArg() {
		if *edit_inplace {
			fmt.Fprintf(os.Stderr, "Warning: Option -i ignored\n")
		}
		s.input = bufio.NewReader(os.Stdin)
		s.process()
	} else {
		for ; currentFileParameter < flag.NArg(); currentFileParameter++ {
			inputFilename = flag.Arg(currentFileParameter)
			// actually do the processing
			s.inputFile, err = os.Open(inputFilename)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error openint input file: %s.\n\n", inputFilename)
				usage()
				os.Exit(-1)
			}
			s.input = bufio.NewReader(s.inputFile)
			var tempFilename string
			if *edit_inplace {
				tempFilename = inputFilename + ".tmp"
				tmpc := 0
				dir, _ := os.Stat(tempFilename)
				for dir != nil {
					tmpc++
					tempFilename = inputFilename + "-" + strconv.Itoa(tmpc) + ".tmp"
					dir, _ = os.Stat(tempFilename)
				}
				f, err := os.Create(tempFilename)
				if err != nil {
					s.inputFile.Close()
					fmt.Fprintf(os.Stderr, "Error opening temp file file for inplace editing: %s\n", err.Error())
					os.Exit(-1)
				}
				s.outputFile = f
			}
			s.process()
			// done processing, close input file
			s.inputFile.Close()
			s.input = nil
			if *edit_inplace {
				s.outputFile.Seek(0, 0)
				// find out about
				dir, err := os.Stat(inputFilename)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error getting information about input file: %s %v\n", err)
					// os.Remove(tempFilename);
					os.Exit(-1)
				}
				// reopen input file
				s.inputFile, err = os.OpenFile(inputFilename, os.O_WRONLY|os.O_TRUNC, dir.Mode())
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error opening input file for in place editing: %s\n", err.Error())
					// os.Remove(tempFilename);
					os.Exit(-1)
				}

				_, e := io.Copy(s.inputFile, s.outputFile)
				s.outputFile.Close()
				s.inputFile.Close()
				if e != nil {
					fmt.Fprintf(os.Stderr, "Error copying temp file back to input file: %s\nFull output is in %s", err.Error(), tempFilename)
				} else {
					os.Remove(tempFilename)
				}
			}
		}
	}
}
