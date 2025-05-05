package cribbage

import (
	"fmt"
	"math/rand"
	"slices"

	"gioui.org/example/outlay/fan/playing"
)

type Phase uint8

const (
	BetweenHands Phase = iota
	Dealing
	Sacrifice
	Cut
	CircularCount
	CountHands
	CountCrib
)

func (p Phase) String() string {
	switch p {
	case BetweenHands:
		return "between"
	case Dealing:
		return "dealing"
	case Sacrifice:
		return "sacrifice"
	case Cut:
		return "cut"
	case CircularCount:
		return "circular count"
	case CountHands:
		return "count hands"
	case CountCrib:
		return "count crib"
	default:
		return "unknown"
	}
}

type Game struct {
	Phase
	Deck    []playing.Card
	CutCard *playing.Card
	Dealer  int
	Crib    []playing.Card
	Players []Player
}

type Player struct {
	Hand, Table []playing.Card
}

func (p Player) String() string {
	return fmt.Sprintf("[Hand: %s, Table: %s]", p.Hand, p.Table)
}

func (g Game) String() string {
	return fmt.Sprintf("[Phase: %v\nDealer: %v\nCrib: %v\nCut: %v\nPlayers: %v\nDeck: %v]\n", g.Phase, g.Dealer, g.Crib, g.CutCard, g.Players, g.Deck)
}

const MinHand = 4

func NewGame(players int) Game {
	var g Game
	g.Players = make([]Player, players)
	g.Dealer = g.NumPlayers() - 1
	for i := range 4 {
		for j := range 13 {
			g.Deck = append(g.Deck, playing.Card{
				Suit: playing.Suit(i),
				Rank: playing.Rank(j),
			})
		}
	}
	g.Phase = Dealing
	return g
}

func (g Game) NumPlayers() int {
	return len(g.Players)
}

func (g Game) Right(player int) int {
	return (player + g.NumPlayers() - 1) % g.NumPlayers()
}

func (g Game) Left(player int) int {
	return (player + 1) % g.NumPlayers()
}

func (g *Game) CutAt(depth int) {
	g.CutCard = &g.Deck[depth]
	g.Phase = CircularCount
}

func DrainInto(src, dest *[]playing.Card) {
	for _, c := range *src {
		*dest = append(*dest, c)
	}
	*src = (*src)[:0]
}

func (g *Game) Reset() {
	for i := range g.Players {
		DrainInto(&(g.Players[i].Hand), &g.Deck)
		DrainInto(&(g.Players[i].Table), &g.Deck)
	}
	DrainInto(&(g.Crib), &g.Deck)
	g.Phase = Dealing
	g.CutCard = nil
}

func (g *Game) DealCardTo(dest *[]playing.Card) {
	card := g.Deck[0]
	g.Deck = g.Deck[1:]
	*dest = append(*dest, card)
}

func (g *Game) DealRound() {
	g.Dealer = g.Left(g.Dealer)
	g.Reset()
	g.Shuffle()
	for range g.CardsToDealPerPlayer() {
		for i := range g.Players {
			g.DealCardTo(&(g.Players[i].Hand))
		}
	}
	for range g.CardsDealtToCrib() {
		g.DealCardTo(&g.Crib)
	}
	g.Phase = Sacrifice
}

func (g Game) CardsToDealPerPlayer() int {
	switch g.NumPlayers() {
	case 2:
		return 6
	case 3:
		return 5
	case 4:
		return 5
	default:
		return 0
	}
}

func (g Game) CardsDealtToCrib() int {
	if g.NumPlayers() == 3 {
		return 1
	}
	return 0
}

func (g *Game) Shuffle() {
	rand.Shuffle(len(g.Deck), func(i, j int) {
		g.Deck[i], g.Deck[j] = g.Deck[j], g.Deck[i]
	})
}

func (g *Game) Sacrifice(player, card int) {
	hand := g.Players[player].Hand
	if len(hand) <= MinHand {
		return
	}
	c := hand[card]
	g.Players[player].Hand = slices.Delete(hand, card, card+1)
	g.Crib = append(g.Crib, c)
}
