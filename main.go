package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

type Card struct {
	value int
	suit  int // 0 - spades, 1 - hearts, 2 - diamonds, 3 - clubs
}

func (card Card) getString() string {
	var suit string
	var value string

	switch card.suit {
	case 0:
		suit = "♠"
	case 1:
		suit = "♥"
	case 2:
		suit = "♦"
	case 3:
		suit = "♣"
	}

	switch card.value {
	case 11:
		value = "J"
	case 12:
		value = "Q"
	case 13:
		value = "K"
	case 1:
		value = "A"
	default:
		value = fmt.Sprintf("%d", card.value)
	}

	return value + suit
}

type Deck struct {
	cards []Card
}

func (d *Deck) deal(num uint) []Card {
	var card []Card

	for i := uint(0); i < num; i++ {
		len := len(d.cards)
		card = append(card, d.cards[len-1])
		d.cards = d.cards[:len-1]
	}

	return card
}

func (d *Deck) create() {
	for suit := 0; suit <= 3; suit++ {
		for num := 1; num <= 13; num++ {
			card := Card{suit: suit, value: num}
			d.cards = append(d.cards, card)
		}
	}
}

func (d *Deck) shuffle() {
	rand.Shuffle(len(d.cards), func(i, j int) { d.cards[i], d.cards[j] = d.cards[j], d.cards[i] })
}

type Game struct {
	deck        Deck
	playerCards []Card
	dealerCards []Card
}

func (game *Game) dealStartingCards() {
	pCards := game.deck.deal(2)
	dCards := game.deck.deal(2)

	game.playerCards = append(game.playerCards, pCards...)
	game.dealerCards = append(game.dealerCards, dCards...)
}

func (game *Game) play(bet float64) float64 {
	var playerValue int
	var playerCardsStr string
	// var dealerValue int

	game.deck.create()
	game.deck.shuffle()
	game.dealStartingCards()

	fmt.Printf("\n----------------------------------\n\n")

	playersTurn := true
	for playersTurn {
		for _, card := range game.playerCards {
			playerValue += card.value
			playerCardsStr = card.getString()
		}
		fmt.Printf("Player has been dealt: %s\n", playerCardsStr)
		playersTurn = game.playerTurn()
	}

	fmt.Printf("\n----------------------------------\n\n")

	return bet
}

func (game *Game) playerTurn() bool {
	return false
}

// func (game *Game) dealerTurn() {

// }

func enterString() string {
	reader := bufio.NewReader(os.Stdin)
	// ReadString will block until the delimiter is entered
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("An error occurred while reading input. Please try again", err)
		return ""
	}

	// remove the delimiter from the string
	input = strings.TrimSuffix(input, "\n")
	return input
}

func main() {
	balance := float64(100)

	for balance > 0 {
		fmt.Printf("Your balance is: $%.2f\n", balance)
		fmt.Printf("Enter your bet (q to quit): ")
		bet, err := strconv.ParseFloat(enterString(), 64)
		if err != nil {
			break
		}
		if bet > balance || bet <= 0 {
			fmt.Println("Invalid bet.")
			continue
		}

		game := Game{}
		balance += game.play(bet)
	}

	fmt.Printf("You left with: $%2.f\n", balance)
}
