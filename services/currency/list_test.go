package currency

import (
	"errors"
	"io"
	"net/url"
	"testing"
)

type readCloserTest struct {
}

func (rc readCloserTest) Read(b []byte) (int, error) {
	return 0, nil
}

func (rc readCloserTest) Close() error {
	return nil
}

type reqParserTest struct {
	reqErr, parseErr error
	rc               io.ReadCloser
}

func (rp *reqParserTest) Request(s string, v url.Values) (io.ReadCloser, error) {
	return rp.rc, rp.reqErr
}

func (rp *reqParserTest) Parse(p payloader, r io.Reader) error {
	return rp.parseErr
}

func TestList(t *testing.T) {
	testCases := []struct {
		rp  *reqParserTest
		err error
	}{
		{
			&reqParserTest{rc: &readCloserTest{}},
			nil,
		},
		{
			&reqParserTest{rc: &readCloserTest{}, reqErr: errors.New("Req Error")},
			errors.New("Req Error"),
		},
		{
			&reqParserTest{rc: &readCloserTest{}, parseErr: errors.New("Parse Error")},
			errors.New("Parse Error"),
		},
	}

	for _, tc := range testCases {
		c := &list{tc.rp}
		_, err := c.List()

		if tc.err == nil && err != nil {
			t.Fatalf("Error while testing currency list: %s", err.Error())
		}

		if tc.err != nil && (err == nil || err.Error() != tc.err.Error()) {
			t.Errorf("got: %s, want: %s", err, tc.err)
		}
	}
}
