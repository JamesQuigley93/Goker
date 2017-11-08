package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type card struct {
	suit   string
	value  string
	weight int
}

type hand struct {
	flush      bool
	straight   bool
	highCard   int
	duplicates []int
}

type deck []card

func newDeck() deck {
	cards := deck{}

	cardSuits := []string{"Spades", "Diamonds", "Hearts", "Clubs"}
	cardValues := []string{"Ace", "Two", "Three", "Four", "Five", "Six", "Seven", "Eight", "Nine", "Ten", "Jack", "Queen", "King"}
	cardWeights := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13}

	for _, suit := range cardSuits {
		for i, value := range cardValues {
			newCard := card{
				suit:   suit,
				value:  value,
				weight: cardWeights[i],
			}
			cards = append(cards, newCard)
		}
	}

	return cards
}

func (d deck) print() {
	for i, card := range d {
		fmt.Println(card.value+" of "+card.suit, i)
	}
	fmt.Printf("\n")
}

func (d deck) shuffle() {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	for i := range d {
		newPosition := r.Intn(len(d) - 1)

		d[i], d[newPosition] = d[newPosition], d[i]
	}
}

func deal(d deck, handSize int) (deck, deck) {
	return d[:handSize], d[handSize:]
}

func userDeal(h deck, d deck) (deck, deck) {
	d.shuffle()
	h, d = deal(d, 5)
	fmt.Printf("Your hand is: \n")
	h.print()
	return h, d
}

func replace(h deck, d deck) (deck, deck) {
	newHand := deck{}
	keep := getKeepInd()
	for _, index := range keep {
		newHand = append(newHand, h[index])
	}
	fmt.Printf("\n")
	newCards, d := deal(d, 5-len(keep))
	newHand = append(newHand, newCards...)
	fmt.Println("Your hand after replacement is: ")
	newHand.print()
	return newHand, d
}

func getKeepInd() []int64 {
	replaceInd := []int64{}
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Which cards do you want to keep?: ")
	text, _ := reader.ReadString('\n')
	text = strings.TrimSpace(text)
	if text == "" {
		return replaceInd
	}
	numbers := strings.Split(text, " ")
	for _, number := range numbers {
		i, err := strconv.ParseInt(number, 10, 64)
		if err != nil {
			fmt.Println(err)
		}
		replaceInd = append(replaceInd, i)
	}
	return replaceInd
}

func checkWin(d deck) hand {
	finalHand := hand{
		flush:      false,
		straight:   false,
		highCard:   0,
		duplicates: []int{},
	}
	if len(d) != 5 {
		return finalHand
	}
	weights := []int{}
	flush := true
	suit := d[0].suit
	for _, card := range d {
		weights = append(weights, card.weight)
		flush = (suit == card.suit && flush)
		suit = card.suit
		if finalHand.highCard < card.weight {
			finalHand.highCard = card.weight
		}
	}

	sort.Ints(weights)
	weight := weights[0] - 1
	straight := true
	for i, value := range weights {
		straight = (value == weight+1 && straight)
		if i > 0 && value == weights[i-1] {
			finalHand.duplicates = append(finalHand.duplicates, value)
		}
		weight = value
	}

	finalHand.flush = flush
	finalHand.straight = straight

	return finalHand
}

func findWinMultiplyer(h hand) int {
	if h.flush || h.straight {
		switch {
		case h.flush && h.straight && h.highCard == 13:
			return 800
		case h.flush && h.straight:
			return 50
		case h.flush:
			return 6
		default:
			return 4
		}
	}

	if len(h.duplicates) != 0 {
		switch {
		case len(h.duplicates) == 1 && h.duplicates[0] > 9:
			return 1
		case len(h.duplicates) == 2 && h.duplicates[0] != h.duplicates[1]:
			return 2
		case len(h.duplicates) == 2 && h.duplicates[0] == h.duplicates[1]:
			return 3
		case len(h.duplicates) == 3:
			return 25
		case len(h.duplicates) == 4:
			return 9

		default:
			return 0
		}
	}

	return 0
}

func playHand(bet int) int {
	cards := newDeck()
	hand := deck{}
	hand, cards = userDeal(hand, cards)
	hand, cards = replace(hand, cards)
	finalHand := checkWin(hand)
	return bet * findWinMultiplyer(finalHand)

}

func getBet(credits int) int {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("You have %v Credits. How much do you want to bet? ", credits)

	text, _ := reader.ReadString('\n')
	text = strings.TrimSpace(text)
	num, err := strconv.ParseInt(text, 10, 0)
	bet := int(num)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("\n")
	return bet
}
