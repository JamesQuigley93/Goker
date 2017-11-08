package main

import (
	"fmt"
)

func main() {
	score := poker()
	fmt.Printf("Game Over! Your score was %v\n", score)
}

func poker() int {
	credits := 100
	highScore := 100
	for credits > 0 {
		bet := getBet(credits)
		credits = credits - bet
		winnings := playHand(bet)
		if winnings == 0 {
			fmt.Println("You lose!")
		} else {
			fmt.Printf("You won: %v credits! \n", winnings)
			credits = credits + winnings
		}
		if credits > highScore {
			highScore = credits
		}
	}
	return highScore
}
