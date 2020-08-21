# `bump`
SemVer bumping made easy!

Bump is a tool for bumping SemVer patterns. Just run `bump <filename>` and bump will increment the first SemVer
pattern it finds.

## Installation
### Brew
```
$ brew tap broothie/bump
$ brew install bump
```

### Go
```
$ go install github.com/broothie/bump
```

### Source
```
$ git checkout https://github.com/broothie/bump.git
$ go install
```

## Usage
Let's assume you have a file `version.txt`:
```
$ cat version.txt
blahblahblah
version = "5.0.2"
```

`bump` will increment the patch segment by default:
```
$ bump version.txt
5.0.2 -> 5.0.3
$ cat version.txt
blahblahblah
version = "5.0.3"
```

You can bump a specific segment with `-s` e.g.:
```
$ bump version.txt -s minor
5.0.3 -> 5.1.0
$ cat version.txt
blahblahblah
version = "5.1.0"
$ bump version.txt -s major
5.1.0 -> 6.0.0
$ cat version.txt
blahblahblah
version = "6.0.0"
```

## Options
```
$ bump -h
usage: bump [<flags>] <filename>

SemVer bumping made easy!

Flags:
  -h, --help             Show context-sensitive help (also try --help-long and --help-man).
  -s, --segment="patch"  SemVer segment to bump (major, minor, or patch).
  -l, --line=LINE        Line number on which to look for SemVer pattern.
  -o, --occurrence=1     If multiple SemVer patterns can be found, use this to indicate which one to bump.
  -d, --dry-run          Don't rewrite file, just print output (overrides '-q').
  -q, --quiet            Suppress output.
  -v, --version          Show application version.

Args:
  <filename>  File containing SemVer pattern to bump.

```
