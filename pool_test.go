package pool_test

import (
	"fmt"
	"strings"
	"testing"

	pool "github.com/rocketlaunchr/go-pool"
)

func TestPool(t *testing.T) {
	type X struct {
		ID string
	}

	randoms := []string{"AAA", "BBB", "CCC", "DDD", "EEE", "FFF", "GGG", "HHH", "III", "JJJ", "KKK", "LLL", "MMM", "NNN"}
	ids := map[string]struct{}{}
	for _, v := range randoms {
		ids[v] = struct{}{}
	}

	i := 0
	p := pool.New(func() *X {
		// This function will panic if Max field is not able to restrict growth
		// of pool and slice goes "index out of range"
		defer func() {
			i++
		}()
		return &X{randoms[i]}
	}, pool.Options{Max: 1})

	for i := 0; i < 1_000_000; i++ {
		i := i
		t.Run(fmt.Sprintf("test %v", i), func(t *testing.T) {
			t.Parallel()
			borrowed := p.Borrow()
			defer borrowed.Return()

			id := borrowed.Item().ID

			if _, exists := ids[id]; !exists {
				t.Errorf("Result was incorrect, got: %s, want: %s.", id, strings.Join(randoms, ","))
			}
		})
	}
}
