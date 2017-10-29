package sofa

import (
	"net/url"
	"strings"
)

// FutonURL attempts to correctly convert any CouchDB path into a URL to be used to
// access the same documents through the Futon web GUI.
func (con *Connection) FutonURL(path string) url.URL {
	patharr := strings.Split(strings.Trim(path, "/"), "/")

	furl := con.URL("/")
	furl.Path = urlConcat(furl.Path, "_utils/")

	if len(patharr) == 0 || patharr[0] == "" {
		return furl
	}

	isDatabaseURL := false
	if len(patharr) == 1 {
		isDatabaseURL = true
	} else {
		switch patharr[1] {
		case "_design", "_all_docs":
			isDatabaseURL = true
		}
	}

	furl.RawQuery = strings.TrimLeft(path, "/")

	if isDatabaseURL {
		furl.Path = urlConcat(furl.Path, "database.html")
		return furl
	}

	furl.Path = urlConcat(furl.Path, "document.html")
	return furl
}

func FromBoolean(b bool) BooleanParameter {
	if b {
		return True
	}

	return False
}

func ToBoolean(b BooleanParameter) bool {
	return b == True
}
