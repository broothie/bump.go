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
	version = "0.0.2"

	// SemVer regexp
	semVerMatcher = regexp.MustCompile(`\d+\.\d+\.\d+`)

	// Args
	command    = kingpin.New("bump", "SemVer bumping made easy!")
	filename   = command.Arg("filename", "File containing SemVer pattern to bump.").Required().ExistingFile()
	segment    = command.Flag("segment", "SemVer segment to bump (major, minor, or patch).").Short('s').Default("patch").String()
	lineNumber = command.Flag("line", "Line number to look for SemVer pattern.").Short('l').Int()
	occurrence = command.Flag("occurrence", "If multiple SemVer patterns can be found, use this to indicate which one to bump.").Short('o').Default("1").Int()
	dryRun     = command.Flag("dry-run", "Don't rewrite file, just print output").Short('d').Bool()
	quiet      = command.Flag("quiet", "Suppress output.").Short('q').Bool()
)

func init() {
	command.Author("Andrew Booth")
	command.Version(version)
	command.VersionFlag.Short('v')
	command.HelpFlag.Short('h')
	if _, err := command.Parse(os.Args[1:]); err != nil {
		fmt.Println(err.Error())
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
	offset := 0
	if *lineNumber > 0 {
		// Line number is specified
		// Split file into lines
		lines := strings.Split(fileContents, "\n")
		if *lineNumber > len(lines) {
			fmt.Printf("line %d of %s doesn't exist (only %d lines)\n", *lineNumber, *filename, len(lines))
			os.Exit(1)
			return
		}

		// Count chars prior to match and set search space
		lineNumberIndex := *lineNumber - 1
		searchSpace = lines[lineNumberIndex]
		for i := 0; i < lineNumberIndex; i++ {
			offset += len(lines[i]) + 1
		}
	}

	// Check occurrences vs matches
	occurrenceIndex := *occurrence - 1
	matches := semVerMatcher.FindAllString(searchSpace, *occurrence)
	if len(matches) == 0 {
		fmt.Printf("no SemVer pattern found in %s\n", *filename)
		os.Exit(1)
		return
	} else if occurrenceIndex > len(matches) {
		fmt.Printf("occurrence %d doesn't exist (only %d SemVer matches found)\n", *occurrence, len(matches))
		os.Exit(1)
		return
	}

	// Bump
	semVer := matches[occurrenceIndex]
	bumpedSemver, err := bump(semVer, *segment)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
		return
	}

	// Print SemVer bump
	if *dryRun || !*quiet {
		fmt.Printf("%s -> %s\n", semVer, bumpedSemver)
	}

	// Don't write if dry run
	if *dryRun {
		return
	}

	// Replace and write
	position := offset + semVerMatcher.FindAllStringIndex(fileContents, *occurrence)[occurrenceIndex][0]
	newFileContents := fileContents[:position] + bumpedSemver + fileContents[position+len(semVer):]
	if err := ioutil.WriteFile(*filename, []byte(newFileContents), os.ModePerm); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func bump(semVer string, segment string) (string, error) {
	segments := strings.Split(semVer, ".")
	majorStr, minorStr, patchStr := segments[0], segments[1], segments[2]

	switch strings.ToLower(segment) {
	case "major":
		majorInt, _ := strconv.Atoi(majorStr)
		majorInt++
		majorStr = strconv.Itoa(majorInt)
		minorStr = "0"
		patchStr = "0"

	case "minor":
		minorInt, _ := strconv.Atoi(minorStr)
		minorInt++
		minorStr = strconv.Itoa(minorInt)
		patchStr = "0"

	case "patch":
		patchInt, _ := strconv.Atoi(patchStr)
		patchInt++
		patchStr = strconv.Itoa(patchInt)

	default:
		return "", fmt.Errorf("'%s' is not a valid SemVer segment", segment)
	}

	return fmt.Sprintf("%s.%s.%s", majorStr, minorStr, patchStr), nil
}
