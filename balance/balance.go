package balance

import (
	"sync"
)

type Weighted struct {
	host      string
	weight    int
	curWeight int
}

var Weighteds []*Weighted
var locker = &sync.Mutex{}

func Init(cfg []string) {
	Weighteds = make([]*Weighted, len(cfg))
	for i := 0; i < len(cfg); i++ {
		Weighteds[i] = &Weighted{cfg[i], 1, 0}
	}
}

func Next() *Weighted {
	index, total := -1, 0
	locker.Lock()
	for i := 0; i < len(Weighteds); i++ {
		Weighteds[i].curWeight += Weighteds[i].weight
		total += Weighteds[i].weight

		if index == -1 || Weighteds[index].curWeight < Weighteds[i].curWeight {
			index = i
		}
	}

	Weighteds[index].curWeight -= total
	locker.Unlock()

	return Weighteds[index]
}
