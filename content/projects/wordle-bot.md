---
title: "Wordle Bot"
date: "2024-03-01"
tags: ["python", "information-theory", "solver", "pandas"]
image: "assets/wordle-bot.png"
link: "https://github.com/iashyam/wordle-bot"
description: "A Python bot that employs information theory (entropy reduction) to solve Wordle puzzles efficiently."
---

A Python bot that employs information theory to solve Wordle puzzles efficiently.

### Key Features
- **Information Theory Solver**: Calculates the word that limits the space of possible solutions the most using entropy reduction (inspired by 3Blue1Brown).
- **Interactive Console Game**: Enter guesses and patterns (e.g. `bbgbb` where `b` is black/grey, `y` is yellow, and `g` is green) to narrow down the target word.
- **Efficiency Metrics**: Shows the next five most efficient guesses and how many possible solutions are remaining.
