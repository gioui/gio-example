/*
Package playing provides types for modeling a deck of conventional
playing cards.
*/
package playing

type (
	Suit  uint8
	Rank  uint8
	Color bool
)

const (
	Spades Suit = iota
	Clubs
	Hearts
	Diamonds
	UnknownSuit
)

const (
	Ace Rank = iota
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
	UnknownRank
)

const (
	Red   Color = true
	Black Color = false
)

type Card struct {
	Suit
	Rank
}

func Deck() []Card {
	d := make([]Card, 0, 52)
	for i := range 4 {
		for k := range 13 {
			d = append(d, Card{
				Suit: Suit(i),
				Rank: Rank(k),
			})
		}
	}
	return d
}

func (r Rank) String() string {
	switch r {
	case Ace:
		return "A"
	case Two:
		return "2"
	case Three:
		return "3"
	case Four:
		return "4"
	case Five:
		return "5"
	case Six:
		return "6"
	case Seven:
		return "7"
	case Eight:
		return "8"
	case Nine:
		return "9"
	case Ten:
		return "10"
	case Jack:
		return "J"
	case Queen:
		return "Q"
	case King:
		return "K"
	default:
		return "?"
	}
}

func (s Suit) String() string {
	switch s {
	case Spades:
		return "♠"
	case Hearts:
		return "♥"
	case Diamonds:
		return "♦"
	case Clubs:
		return "♣"
	default:
		return "?"
	}
}

func (s Suit) Color() Color {
	switch s {
	case Spades, Clubs:
		return Black
	case Hearts, Diamonds:
		return Red
	default:
		return Black
	}
}
