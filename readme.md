# `bump`
SemVer bumping made easy!

## Installation
### Brew
```shell script
$ brew tap broothie/bump
$ brew install bump
```

### Go
```shell script
$ go install github.com/broothie/bump
```

### Source
```shell script
$ git checkout https://github.com/broothie/bump.git
$ go install
```

## Usage
Let's assume you have a file `version.txt`:
```shell script
$ cat version.txt
version = "5.0.2"
```

`bump` will increment the patch segment by default:
```
$ bump version.txt
5.0.2 -> 5.0.3
$ cat version.txt
version = "5.0.3"
```

Specify the minor segment with `-s`:
```
$ bump version.txt -s minor
5.0.2 -> 5.1.0
$ cat version.txt
version = "5.1.0"
```

or the major segment: 
```
$ bump version.txt -s major
5.0.2 -> 6.0.0
$ cat version.txt
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
  -l, --line=LINE        Line number to look for SemVer pattern.
  -o, --occurrence=1     If multiple SemVer patterns can be found, use this to indicate which one to bump.
  -q, --quiet            Suppress output.
  -v, --version          Show application version.

Args:
  <filename>  File containing SemVer pattern to bump.

```
