# Letter Boxed Solver

This program solves a [Letter Boxed](https://www.nytimes.com/puzzles/letter-boxed) puzzle with a brute-force algorithm. The Go version (`go-src/`) is a re-write of the initial, slower Python version (`py-src/`). The dictionary used (`dictionary.txt`) is taken from [sindresorhus/word-list](https://github.com/sindresorhus/word-list), but Letter Boxed actually uses a much smaller dictionary.

## Usage

For the Python version:

```
$ python3 py-src/main.py
```

For the Go version:

```
$ cd go-src
$ go run cmd/solve/main.go
```

