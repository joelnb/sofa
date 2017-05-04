package sofa

import (
	"errors"
	"testing"
)

func TestErrorStatus(t *testing.T) {
	tests := []struct {
		err            error
		statusCode     int
		expectedResult bool
	}{
		{
			err: ResponseError{
				Method:     "PUT",
				StatusCode: 412,
				URL:        "http://couch.db:5984/somedb",

				Err:    "file_exists",
				Reason: "The database could not be created, the file already exists.",
			},
			statusCode:     412,
			expectedResult: true,
		},
		{
			err: ResponseError{
				Method:     "PUT",
				StatusCode: 412,
				URL:        "http://couch.db:5984/somedb",

				Err:    "file_exists",
				Reason: "The database could not be created, the file already exists.",
			},
			statusCode:     409,
			expectedResult: false,
		},
		{
			err:            errors.New("not a ResponseError"),
			statusCode:     401,
			expectedResult: false,
		},
		{
			err: ResponseError{
				Method:     "PUT",
				StatusCode: 401,
				URL:        "http://couch.db:5984/somedb/_design/someDesign/view",
			},
			statusCode:     401,
			expectedResult: true,
		},
	}

	for i, test := range tests {
		if ErrorStatus(test.err, test.statusCode) != test.expectedResult {
			t.Errorf("%d: expected %t got %t", i, test.expectedResult, !test.expectedResult)
		}
	}
}
