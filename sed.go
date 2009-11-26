package sed

import (
	"flag";
	"fmt";
	"io";
	"os";
	"strings";
	"container/vector";
)

const (
	versionMajor	= 0;
	versionMinor	= 1;
	versionPoint	= 0;
)

var versionString string

func init() {
	versionString = fmt.Sprintf("%d.%d.%d", versionMajor, versionMinor, versionPoint)
}

//var show_version = flag.Bool("v", false, "Show version information.")
var show_help = flag.Bool("h", false, "Show help information.")
var quiet = flag.Bool("n", false, "Don't print the pattern space at the end of each script cycle.")
var script = flag.String("e", "", "The script used to process the input file.")
var script_file = flag.String("f", "", "Specify a file to read as the script. Ignored if -e present")
var inplace_suffix = flag.String("i", "", "This option specifies that files are to be edited in-place.")
var line_wrap = flag.Uint("l", 70, "Specify the default line-wrap length for the l command. A length of 0 (zero) means to never wrap long lines. If not specified, it is taken to be 70.")
var unbuffered = flag.Bool("u", false, "Buffer both input and output as minimally as practical.")

// the input file split into lines
var inputLines []string
var usageShown bool = false
var commands *vector.Vector
var outputFile = os.Stdout

func usage() {
	// only show usage once.
	if !usageShown {
		usageShown = true;
		fmt.Fprint(os.Stdout, "sed [options] [script] input_file\n\n");
		flag.PrintDefaults();
	}
}

func readInputFile() {
	var inputFilename string;
	if flag.NArg() > 1 {
		inputFilename = flag.Arg(1)
	} else if flag.NArg() == 1 {
		inputFilename = flag.Arg(0)
	}
	if len(inputFilename) == 0 {
		fmt.Fprint(os.Stderr, "No input file specified.\n\n");
		usage();
		os.Exit(-1);
	}
	f, err := os.Open(inputFilename, os.O_RDONLY, 0);
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening input file %s\n", inputFilename);
		os.Exit(-1);
	}
	b, _ := io.ReadAll(f);
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input file %s\n", inputFilename);
		os.Exit(-1);
	}
	inputLines = strings.Split(string(b), "\n", 0);
}

func parseScript() (err os.Error) {
	// assume only one
	commands = new(vector.Vector);
	// a script may be a single command or it may be several
	scriptLines := strings.Split(*script, "\n", 0);
	for idx, line := range scriptLines {
		line = strings.TrimSpace(line);
		if strings.HasPrefix(line, "#") || len(line) == 0 {
			// comment line
			continue
		}
		// this isn't really right. There may be slashes in the regular expression
		pieces := strings.Split(line, "/", 0);
		c, err := NewCmd(pieces);
		if err != nil {
			fmt.Printf("%v line %d: %s\n", err, idx+1, line);
			os.Exit(-1);
		}
		commands.Push(c);
	}
	return nil;
}

func printPatternSpace(s string) {
	l := len(s);
	if *line_wrap <= 0 || l < int(*line_wrap) {
		fmt.Fprintf(outputFile, "%s\n", s)
	} else {
		for i := 0; i < l; i += int(*line_wrap) {
			endOfLine := i + int(*line_wrap);
			if endOfLine > l {
				endOfLine = l
			}
			buf := s[i:endOfLine];
			fmt.Fprintf(outputFile, "%s\n", buf);
		}
	}
}

func process() {
	for _, patternSpace := range inputLines {
		for c := range commands.Iter() {
			out, stop, err := c.(*cmd).processLine(patternSpace);
			if err != nil {
				fmt.Printf("%v\n", err);
				os.Exit(-1);
			}
			if stop {
				break
			}
			patternSpace = out;
		}
		if !*quiet {
			printPatternSpace(patternSpace)
		}
	}
}

func Main() {
	flag.Parse();
  // if *show_version {
  //  fmt.Fprintf(os.Stdout, "Version: %s (c)2009 Geoffrey Clements All Rights Reserved\n\n", versionString)
  // }
	if *show_help {
		usage();
		return;
	}
	// we need a script
	if len(*script) == 0 {
		// no -e so try -f
		if len(*script_file) > 0 {
			f, err := os.Open(*script_file, os.O_RDONLY, 0);
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error opening file %s\n", *script_file);
				os.Exit(-1);
			}
			b, _ := io.ReadAll(f);
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading script file %s\n", *script_file);
				os.Exit(-1);
			}
			s := string(b);
			script = &s;
		} else if flag.NArg() > 1 {
			s := flag.Arg(0);
			script = &s;
		}
	}

	// if script still isn't set we are screwed, exit.
	if len(*script) == 0 {
		fmt.Fprint(os.Stderr, "No script found.\n\n");
		usage();
		os.Exit(-1);
	}

	// actually do the processing
	readInputFile();
	parseScript();
	process();
}
