# wordler

Play [WORDLE](https://www.powerlanguage.co.uk/wordle/) on the command line.

[![demo](https://asciinema.org/a/452632.svg)](https://asciinema.org/a/452632?autoplay=1)

## Installation

If you don't have go installed, you can install it with:
```
brew install go
```

Check your version:
```
go version
```

It should be at least 1.17.

And then:

```
go install github.com/bertrandom/wordler@latest
```

## Usage

```
wordler
```

If you want to play a different date, try:
```
wordler -date 2021-11-24
```

Wordler now uses the NYTimes word list because that's what everyone else is using, but if you want to use the original wordlist, use the `-legacy` flag, like this:
```
wordler -legacy
```