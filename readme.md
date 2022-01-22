# Solving Wordle, the Hard Way

It turns out that, contrary to some, Wordly is solvable using minmax ... just barely.

We are doing the needful version - can't cheat by backtracking to words that are impossible, 1's and 2's have to be
preserved.

I did it in Go because Python was just slow enough to not work. Bummer.

## Dataset

Using
the [alphabetized Wordle list](https://gist.githubusercontent.com/cfreshman/a03ef2cba789d8cf00c08f767e0fad7b/raw/a9e55d7e0c08100ce62133a1fa0d9c4f0f542f2c/wordle-answers-alphabetical.txt)
to avoid spoiling myself.

## Solutions

It is clear the approach to the classical game 20 questions can work here, viva Information Theory. Actually even that
greedy version is very strong - at each step, simply figure out for each possible guest what the answer will tell you -
and make the maximum sized split as small as possible.

This is implemented as ```GreedyNeedfulPlayer```, fails on 13 cases, or 0.562% of the time, average guess length is
3.655. Not bad.

A better solution, naturally, is the exact solution, which is minimax. Wordle is barely solvable, because it has too
many cases to be solvable straight away, but because the first guess can be pre-computed optimally, it can be stored and
reused.

## Optimizations

We prune very aggressively. Once we have a guess with path of length t, any path that goes beyond t is assumed to be
infinite.

We prune both the min and the max step.

Words and answers are encoded in integers to save space and reduce the need for heap usage. It really speeds it up.

## Optimizations That Didn't Work Out

Alphabetized word list seems like it could be improved. It can't.

Tried it sorted so that more "splitty" words are tried first. The idea is to find the max path early. It doesn't really
help the hard case though.

It can be run on multiple threads. But it is fast enough as it is.

# And now, Ladies and Gentlemen...

I ran ```play(reducedWordlist, dm, MinMaxPlayer(true))``` to find the best first guess, and it is ... well, a lot of the words
when used as first have a max depth of 6. I chose ```learn```, because I really learned that sometimes things barely work
out.

With that start, simply run a minmax for each further set of allowed words as
```
play(reducedWordlist, dm, FastFirstHand(game.ToWord("learn"), MinMaxPlayer(false)))
```

The results show 0% failure rate with that dictionary. Beat that, ML.

![All solved plot](doc/best-algorithm.png?raw=true "Best algorithms always finish in 6")


