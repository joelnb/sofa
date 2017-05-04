package sofa

import (
	"testing"

	"github.com/nbio/st"
)

func TestUrlConcat(t *testing.T) {
	res := "test/url"
	tests := []string{
		urlConcat("test/", "url"),
		urlConcat("test/", "/url"),
		urlConcat("test", "/url"),
		urlConcat("test", "url"),
	}

	for i, test := range tests {
		st.Expect(t, test, res, i)
	}
}

func TestEncodeValue(t *testing.T) {
	tests := []struct {
		value          interface{}
		expectedResult string
	}{
		{
			value:          "simple",
			expectedResult: "simple",
		},
		{
			value:          "es<@p3d",
			expectedResult: "es<@p3d",
		},
		{
			value:          map[string]string{"a": "map"},
			expectedResult: "{\"a\":\"map\"}",
		},
	}

	for n, tc := range tests {
		str, err := encodeValue(tc.value)
		if err != nil {
			t.Errorf("encodeValue error: %s", err.Error())
		}

		st.Expect(t, str, tc.expectedResult, n)
	}
}

func TestEncodeOptions(t *testing.T) {
	singleTests := []struct {
		name           string
		value          interface{}
		expectedResult string
	}{
		{
			name:           "booly",
			value:          true,
			expectedResult: "booly=true",
		},
		{
			name:           "string",
			value:          "something_stringy",
			expectedResult: "string=something_stringy",
		},
		{
			name:           "hash",
			value:          "ba9a9a00dece730a6c97cb2bb3527990  tails-i386-2.5.iso",
			expectedResult: "hash=ba9a9a00dece730a6c97cb2bb3527990++tails-i386-2.5.iso",
		},
		{
			name:           "countarr",
			value:          []string{"one", "two", "three"},
			expectedResult: "countarr=%5B%22one%22%2C%22two%22%2C%22three%22%5D",
		},
		{
			name:           "mapper",
			value:          map[string]interface{}{"name": "Harry", "age": 25},
			expectedResult: "mapper=%7B%22age%22%3A25%2C%22name%22%3A%22Harry%22%7D",
		},
	}

	for n, tc := range singleTests {
		opts := NewURLOptions()
		opts.Add(tc.name, tc.value)

		st.Expect(t, opts.Encode(), tc.expectedResult, n)
	}
}
