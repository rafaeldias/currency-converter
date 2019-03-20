package currency

import (
	"errors"
	"testing"
)

func TestSuccessful(t *testing.T) {
	var testCases = []struct {
		success bool
		want    error
	}{
		{true, nil},                          // successful
		{false, errors.New("Testing error")}, // error
	}

	for _, tc := range testCases {
		var ep errPayload

		if tc.want != nil {
			ep = errPayload{Info: tc.want.Error()}
		}

		var d = defaultPayload{Success: tc.success, Error: ep}
		err := d.Successful()
		if (tc.want == nil && err != tc.want) || (tc.want != nil && (err == nil || err.Error() != tc.want.Error())) {
			t.Errorf("got: %s, want: %s", err, tc.want)
		}
	}
}
