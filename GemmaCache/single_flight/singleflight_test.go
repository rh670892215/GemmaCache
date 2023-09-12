package single_flight

import (
	"testing"
)

func TestDo(t *testing.T) {
	var g ConcurrencyLimiter
	v, err := g.Do("key", func() (interface{}, error) {
		return "bar", nil
	})

	if v != "bar" || err != nil {
		t.Errorf("Do v = %v, error = %v", v, err)
	}
}
