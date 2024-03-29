package sofa

import (
	"net/url"
	"strings"
)

// FutonURL attempts to correctly convert any CouchDB path into a URL to be used to
// access the same documents through the Futon web GUI.
func (con *CouchDB1Connection) FutonURL(path string) url.URL {
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

// FauxtonURL attempts to correctly convert any CouchDB path into a URL to be used to
// access the same documents through the Fauxton web GUI.
func (con *CouchDB2Connection) FauxtonURL(path string) url.URL {
	patharr := strings.Split(strings.Trim(path, "/"), "/")

	furl := con.URL("/")
	furl.Path = urlConcat(furl.Path, "_utils/")

	if len(patharr) == 0 {
		return furl
	} else if len(patharr) == 1 {
		path = path + "/_all_docs"
	}

	furl.Fragment = urlConcat("database", path)

	return furl
}

// FromBoolean converts a standard bool value into a sofa.BooleanParameter for
// use in documents to allow not including unset booleans.
func FromBoolean(b bool) BooleanParameter {
	if b {
		return True
	}

	return False
}

// ToBoolean converts a sofa.BooleanParameter value into a standard bool.
func ToBoolean(b BooleanParameter) bool {
	return b == True
}
