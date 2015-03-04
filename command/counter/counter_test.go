// Tests the counter package
// makes sure that the counts
// in a concurrent manner

package counter

import (
	"sync"
	"testing"
)

// Constants to let the compiler find typos.
const (
	A100s  = "100s"
	A200s  = "200s"
	A300s  = "300s"
	A400s  = "400s"
	A500s  = "500s"
	Errors = "Errors"
	Totals = "Totals"
)

// Creates 4 concurrent thread groups and wait for
// completion
func body(c *Counter, t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(4)
	go func() {
		c.Incr(A100s, 2)
		c.Incr(A200s, 1)
		c.Incr(A400s, 1)
		c.Incr(A300s, 1)
		c.Incr(A500s, 1)
		c.Incr(Totals, 100)
		wg.Done()
	}()
	go func() {
		c.Incr("200s", 21)
		c.Incr("100s", 6)
		c.Incr("300s", 1)
		c.Incr("400s", 4)
		c.Incr("Totals", 100)
		wg.Done()
	}()
	go func() {
		c.Incr("100s", 2)
		c.Incr("200s", 6)
		c.Incr("400s", 1)
		c.Incr("300s", 1)
		c.Incr("Totals", 100)
		wg.Done()
	}()
	go func() {
		c.Incr("400s", 2)
		c.Incr("200s", 1)
		c.Incr("100s", 3)
		c.Incr("300s", 1)
		c.Incr("Totals", 100)
		wg.Done()
	}()
	// sync.WaitGroups: wait until all 4 threads report Done.
	wg.Wait()

	expected := map[string]int{
		"100s":   13,
		"200s":   29,
		"300s":   4,
		"400s":   8,
		"500s":   1,
		"Errors": 0,
		"Totals": 400,
	}
	for k, v := range expected {
		if v != c.Get(k) {
			t.Errorf("counter %s: expected %d, got %d", k, v, c.Get(k))
		}
	}
}

// Run the test and create a new counter to test
func TestCounter(t *testing.T) {
	c := New()
	body(c, t)
}
