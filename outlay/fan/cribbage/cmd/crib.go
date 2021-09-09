package main

import (
	"fmt"

	"gioui.org/example/outlay/fan/cribbage"
)

func main() {
	g := cribbage.NewGame(2)
	fmt.Println(g)
	g.DealRound()
	fmt.Println(g)
	g.Sacrifice(0, 0)
	g.Sacrifice(0, 4)
	g.Sacrifice(1, 0)
	g.Sacrifice(1, 4)
	fmt.Println(g)
	g.CutAt(10)
	fmt.Println(g)
	g.Reset()
	fmt.Println(g)
	g.DealRound()
	fmt.Println(g)
}
