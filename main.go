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
	var cards []Card

	for i := uint(0); i < num; i++ {
		cards = append(cards, d.cards[0])
		d.cards = d.cards[1:]
	}

	return cards
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

	fmt.Printf("Player has been dealt: %s\n\n", game.generateDealtString(game.playerCards))

	dealerScore := game.calculateDealtCards(game.dealerCards)
	if dealerScore == 21 {
		fmt.Printf("Dealer shows: %s\n\n", game.generateDealtString(game.dealerCards))
	} else {
		fmt.Printf("Dealer shows: %s\n\n", game.generateDealtString(game.dealerCards[:1]))
	}
}

func (game *Game) calculateDealtCards(cards []Card) int {
	sum := 0
	aces := 0
	for _, card := range cards {
		if card.value >= 10 {
			sum += 10
			continue
		}

		if card.value == 1 {
			aces++
			continue
		}

		sum += card.value
	}

	if aces == 0 {
		return sum
	}

	if sum < (11 + aces - 1) {
		return sum + 11 + aces - 1
	} else {
		return sum + aces
	}
}

func (game *Game) generateDealtString(cards []Card) string {
	var dealtStr string
	dealtSum := game.calculateDealtCards(cards)

	for _, card := range cards {
		dealtStr += card.getString() + " "
	}

	if len(cards) != 1 {
		dealtStr += "= " + strconv.Itoa(dealtSum)
	}

	return dealtStr
}

func (game *Game) play(bet float64) float64 {
	defer fmt.Printf("----------------------------------\n\n")
	fmt.Printf("\n----------------------------------\n\n")

	game.deck.create()
	game.deck.shuffle()
	game.dealStartingCards()

	playerBlackjack := game.calculateDealtCards(game.playerCards) == 21
	dealerBlackjack := game.calculateDealtCards(game.dealerCards) == 21

	if playerBlackjack && dealerBlackjack {
		fmt.Printf("Player and Dealer both have Blackjack.\n\n")
		return 0
	}

	if playerBlackjack && !dealerBlackjack {
		fmt.Printf("Blackjack!\n\n")
		return bet * 1.5
	}

	if !playerBlackjack && dealerBlackjack {
		fmt.Printf("Dealer has Blackjack: %s\n\n", game.generateDealtString(game.dealerCards))
		return -bet
	}

	game.playerTurn()
	playerScore := game.calculateDealtCards(game.playerCards)
	if playerScore > 21 {
		fmt.Printf("Player bust: %s\n\n", game.generateDealtString(game.playerCards))
		return -bet
	}
	fmt.Printf("\n----------------------------------\n\n")

	game.dealerTurn()
	dealerScore := game.calculateDealtCards(game.dealerCards)
	if dealerScore > 21 {
		fmt.Printf("Dealer bust!\n\n")
		return bet
	}

	if dealerScore == playerScore {
		fmt.Println("Dealer and player tie")
		return 0
	} else if dealerScore > playerScore {
		fmt.Println("Dealer wins")
		return -bet
	}

	fmt.Println("Player wins!")
	return bet
}

func (game *Game) playerTurn() {
	playerHit := "h"
	fmt.Printf("Would you like to hit or stay (H/S)? ")
	playerHit = enterString()

	for playerHit == "h" || playerHit == "H" {
		newCard := game.deck.deal(1)
		game.playerCards = append(game.playerCards, newCard...)

		fmt.Printf("\nYou are dealt: %s\n", game.generateDealtString(newCard))

		playerValue := game.calculateDealtCards(game.playerCards)
		if playerValue > 21 {
			return
		}

		if playerValue == 21 {
			fmt.Printf("Player has 21: %s\n\n", game.generateDealtString(game.playerCards))
			return
		}

		fmt.Printf("Player now has: %s\n\n", game.generateDealtString(game.playerCards))

		fmt.Printf("Would you like to hit or stay (H/S)? ")
		playerHit = enterString()
	}
}

func (game *Game) dealerTurn() {
	dealerScore := game.calculateDealtCards(game.dealerCards)
	for dealerScore < 17 {
		newCard := game.deck.deal(1)
		game.dealerCards = append(game.dealerCards, newCard...)
		dealerScore = game.calculateDealtCards(game.dealerCards)
	}

	fmt.Printf("Dealer has: %s\n\n", game.generateDealtString(game.dealerCards))
}

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
