package currency

import "testing"

func TestList(t *testing.T) {
	l := List{"AED": "United Arab Mirates Dirham"}
	s := currencyLayerServer(l)

	testCases := []struct {
		URL       string
		accessKey string
		want      List
	}{
		{s.URL, *accessKey, l},
	}

	for _, tc := range testCases {
		c := New(Credential{tc.URL, tc.accessKey})
		list, err := c.List()
		if err != nil {
			t.Fatalf("Error while requesting currency list: %s", err.Error())
		}

		for k, v := range tc.want {
			if curr, ok := list[k]; !ok || curr != v {
				t.Fatalf("got: %s; want: %s", list, tc.want)
			}
		}
	}
}

func TestListError(t *testing.T) {
	l := List{"AED": "United Arab Mirates Dirham"}
	s := currencyLayerServer(l)

	testCases := []struct {
		URL       string
		accessKey string
		want      string
	}{
		{s.URL, "", "User did not supply an access key or supplied an invalid access key."},
		{"", "", "Get /api/list?access_key=: unsupported protocol scheme \"\""},
	}

	for _, tc := range testCases {
		c := New(Credential{tc.URL, tc.accessKey})
		_, err := c.List()
		if err != nil && err.Error() != tc.want {
			t.Fatalf("got: %s; want: %s", err.Error(), tc.want)
		}
	}
}
