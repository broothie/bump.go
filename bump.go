package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/alecthomas/kingpin"
)

var (
	// Version
	version = "0.0.1"

	// SemVer regexp
	semVerMatcher = regexp.MustCompile(`\d+\.\d+\.\d+`)

	// Args
	command    = kingpin.New("bump", "SemVer bumping made easy.")
	filename   = command.Arg("filename", "File containing SemVer pattern to bump.").Required().ExistingFile()
	segment    = command.Flag("segment", "SemVer segment to bump (major, minor, or patch).").Short('s').Default("patch").String()
	lineNumber = command.Flag("line", "Line number to look for SemVer pattern.").Short('l').Int()
	occurrence = command.Flag("occurrence", "If multiple SemVer patterns are found, use this to indicate which one to bump.").Short('o').Default("1").Int()
	debump     = command.Flag("debump", "Decrement instead of increment.").Short('d').Bool()
	quiet      = command.Flag("quiet", "Suppress output.").Short('q').Bool()
)

func init() {
	command.Version(version)
	command.VersionFlag.Short('v')
	command.HelpFlag.Short('h')
	if _, err := command.Parse(os.Args); err != nil {
		fmt.Printf(err.Error())
		os.Exit(1)
		return
	}
}

func main() {
	// Read file
	fileContentsBytes, err := ioutil.ReadFile(*filename)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
		return
	}

	fileContents := string(fileContentsBytes)
	searchSpace := fileContents

	// Find SemVer pattern
	occurrenceIndex := *occurrence - 1
	var position int
	switch {
	case *lineNumber > 0: // Line number is specified
		lines := strings.Split(fileContents, "\n")
		if *lineNumber > len(lines) {
			fmt.Printf("line %d of %s doesn't exist (%s only has %d lines)\n",
				*lineNumber,
				*filename,
				*filename,
				len(lines),
			)

			os.Exit(1)
			return
		}

		lineNumberIndex := *lineNumber - 1
		priorChars := 0
		for i := 0; i < lineNumberIndex; i++ {
			priorChars += len(lines[i]) + 1
		}

		line := lines[lineNumberIndex]
		searchSpace = line
		position = priorChars + semVerMatcher.FindAllStringIndex(line, -1)[occurrenceIndex][0]

	default: // Replace first occurrence in file
		position = semVerMatcher.FindAllStringIndex(fileContents, -1)[occurrenceIndex][0]
	}

	// Bump
	semVer := semVerMatcher.FindAllString(searchSpace, -1)[occurrenceIndex]
	bumpedSemver, err := bump(semVer, *segment, *debump)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
		return
	}

	if !*quiet {
		fmt.Printf("%s: %s -> %s\n", *filename, semVer, bumpedSemver)
	}

	// Replace and write
	newFileContents := fileContents[:position] + bumpedSemver + fileContents[position+len(semVer):]
	if err := ioutil.WriteFile(*filename, []byte(newFileContents), os.ModePerm); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func bump(semVer string, segment string, debump bool) (string, error) {
	segments := strings.Split(semVer, ".")
	majorStr, minorStr, patchStr := segments[0], segments[1], segments[2]

	addee := 1
	if debump {
		addee = -1
	}

	switch strings.ToLower(segment) {
	case "major":
		majorInt, _ := strconv.Atoi(majorStr)
		majorInt += addee
		majorStr = strconv.Itoa(majorInt)

	case "minor":
		minorInt, _ := strconv.Atoi(minorStr)
		minorInt += addee
		minorStr = strconv.Itoa(minorInt)

	case "patch":
		patchInt, _ := strconv.Atoi(patchStr)
		patchInt += addee
		patchStr = strconv.Itoa(patchInt)

	default:
		return "", fmt.Errorf("'%s' is not a valid SemVer segment", segment)
	}

	return fmt.Sprintf("%s.%s.%s", majorStr, minorStr, patchStr), nil
}
