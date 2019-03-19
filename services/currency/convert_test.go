package currency

import "testing"

func TestConvert(t *testing.T) {
	s := currencyLayerServer(List{})

	testCases := []struct {
		URL       string
		accessKey string
		from      string
		to        string
		value     float32
		want      float32
	}{
		{s.URL, *accessKey, "USD", "GBP", 10, 6.58443},
	}

	for _, tc := range testCases {
		c := New(Credential{tc.URL, tc.accessKey})
		res, err := c.Convert(tc.from, tc.to, tc.value)
		if err != nil {
			t.Fatalf("Error while requesting currency list: %s", err.Error())
		}

		if res != tc.want {
			t.Errorf("got: %f; want: %f", res, tc.want)
		}
	}
}
