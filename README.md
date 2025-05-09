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

## Letter Boxed's dictionary

The dictionary that Letter Boxed uses is exposed on `window.gameData.dictionary` on [Letter Boxed](https://www.nytimes.com/puzzles/letter-boxed), though it's different every day because it only contains the words that can be formed with the letters for the current day. As of 3 May 2025, it can be retrieved like so:

```
$ curl -s https://www.nytimes.com/puzzles/letter-boxed | grep -oP '"dictionary":\K\[.*?\]'
```

As a note, the `window.gameData` object also contains the suggested solution for the current day.
