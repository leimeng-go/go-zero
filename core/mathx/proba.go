package mathx

import (
	"math/rand"
	"sync"
	"time"
)

// A Proba is used to test if true on given probability.
type Proba struct {
	// rand.New(...) returns a non thread safe object
	r    *rand.Rand
	lock sync.Mutex
}

// NewProba returns a Proba.
func NewProba() *Proba {
	return &Proba{
		r: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}
// 用于在给定的概率下返回`true`或`false`，并发安全的
// TrueOnProba checks if true on given probability.
func (p *Proba) TrueOnProba(proba float64) (truth bool) {
	p.lock.Lock()
	truth = p.r.Float64() < proba
	p.lock.Unlock()
	return
}
