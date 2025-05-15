# Letter unBoxed

A computed archive of every accepted 2-word solution for the New York Times' daily [Letter Boxed](https://www.nytimes.com/puzzles/letter-boxed) puzzles.

## How it works

The solver program (`api/`) uses a brute-force algorithm to find solutions for the current Letter Boxed puzzle. The main site (`website/`) displays the archive of past solutions.

A new Letter Boxed puzzle is released daily at 7AM UTC, so at 6AM UTC every day, the main site calls the solver API and stores the solutions on a [Turso](https://turso.tech/) database.

Since it takes too long to run than most serverless providers allow, the solver API is hosted on [fly.io](https://fly.io/) (with the cheapest machine they offer), while the website is hosted on [Vercel](https://vercel.com/home) (free tier) to minimise costs.

## Usage

To start the webapp:

```
$ cd website
$ pnpm run dev
```

To start the API server:

```
$ cd api
$ go run cmd/serve/main.go
```

You can also run the solver without starting the server:

```
$ cd api
$ go run cmd/solve/main.go --max-words 2
```
