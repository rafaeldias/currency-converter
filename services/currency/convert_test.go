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
		c := convert{Credential{tc.URL, tc.accessKey}}
		res, err := c.Convert(tc.from, tc.to, tc.value)
		if err != nil {
			t.Fatalf("Error while requesting currency list: %s", err.Error())
		}

		if res != tc.want {
			t.Errorf("got: %f; want: %f", res, tc.want)
		}
	}
}

func TestConvertError(t *testing.T) {
	s := currencyLayerServer(List{})

	testCases := []struct {
		URL       string
		accessKey string
		from      string
		to        string
		value     float32
		want      string
	}{
		{s.URL, "", "USD", "GBP", 10, 6.58443, "User did not supply an access key or supplied an invalid access key."},
		{"", "", "Get /api/convert?access_key=: unsupported protocol scheme \"\""},
	}

	for _, tc := range testCases {
		c := convert{Credential{tc.URL, tc.accessKey}}
		_, err := c.Convert(tc.from, tc.to, tc.value)
		if err != nil && err.Error() != tc.want {
			t.Fatalf("got: %s; want: %s", err.Error(), tc.want)
		}
	}
}
