package currency

import "testing"

func TestNew(t *testing.T) {
	var curr = New("x", "y")

	if c, ok := curr.(currency); !ok {
		t.Errorf("got: %v, want: *currency", c)
	}
}
