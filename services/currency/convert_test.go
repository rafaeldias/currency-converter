package currency

import (
	"errors"
	"testing"
)

func TestConvert(t *testing.T) {
	testCases := []struct {
		rp       *reqParserTest
		from, to string
		value    float32
		err      error
	}{
		{&reqParserTest{rc: &readCloserTest{}}, "USD", "BRA", 6.68443, nil},
		{
			&reqParserTest{rc: &readCloserTest{}, reqErr: errors.New("Req Error")},
			"USD", "BRA", 6.68443,
			errors.New("Req Error"),
		},
		{
			&reqParserTest{rc: &readCloserTest{}, parseErr: errors.New("Parse Error")},

			"USD", "BRA", 6.68443,
			errors.New("Parse Error"),
		},
	}

	for _, tc := range testCases {
		c := &convert{tc.rp}
		_, err := c.Convert(tc.from, tc.to, tc.value)

		if tc.err == nil && err != nil {
			t.Fatalf("Error while testing currency list: %s", err.Error())
		}

		if tc.err != nil && (err == nil || err.Error() != tc.err.Error()) {
			t.Errorf("got: %s, want: %s", err, tc.err)
		}
	}
}
