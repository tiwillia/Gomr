package gomr

import (
	"errors"
	"github.com/golang/glog"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

// GET provided url over tcp. Returns a string with the response body.
func HttpGet(url string) (response []byte, err error) {
	client := &http.Client{}

	var req *http.Request
	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}

	var resp *http.Response
	resp, err = client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	response, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		return response, errors.New("HTTP status code: " + strconv.Itoa(resp.StatusCode))
	}

	return
}

// Test if 'match' matches 'line' and, if it does, return the substring indicated by regex 'pull'
func MatchAndPull(line, regex, pull string) (substr string) {
	matchRegex, err := regexp.Compile(regex)
	if err != nil {
		glog.Infoln("Unable to parse regex:", regex)
		return ""
	}

	if matchRegex.MatchString(line) {
		pullRegex := regexp.MustCompile(pull)
		substr := pullRegex.FindStringSubmatch(line)
		if len(substr) == 0 {
			return ""
		}
		return substr[1]
	} else {
		return ""
	}
}

// Test if line matches the provided regex string
func Match(line, regex string) bool {
	matchRegex, err := regexp.Compile(regex)
	if err != nil {
		glog.Infoln("Unable to parse regex:", regex)
		return false
	}

	return matchRegex.MatchString(line)
}

// CanonicalizeIrcNick converts a string to its downcased version according to
// the IRC protocol, with control characters removed.
func CanonicalizeIrcNick(nick string) string {
	// Strip control characters and convert to lowercase.
	return strings.Map(func(r rune) rune {
		switch {
		case r == '[':
			// '[' is uppercase '{'.
			return '{'
		case r == ']':
			// ']' is uppercase '}'.
			return '}'
		case r == '\\':
			// '\' is uppercase '|'.
			return '|'
		case r > 'A' && r <= 'Z':
			// Make uppercase letters lowercase.
			return r + 32
		case r > 32 && r < 127:
			// Non-control characters are fine.
			return r
		}
		return -1
	}, nick)
}

// Convert an int to a string and add the appropriate suffix
func addSuffix(num int) (fullNum string) {
	n := strconv.Itoa(num)
	suffix := ""
	switch {
	case Match(n, `^1.$`):
		suffix = "th"
	case Match(n, `.*[456789]+$`):
		suffix = "th"
	case Match(n, `[\d]+0$`):
		suffix = "th"
	case Match(n, `.*3$`):
		suffix = "rd"
	case Match(n, `.*2$`):
		suffix = "nd"
	case Match(n, `.*1$`):
		suffix = "st"
	}
	fullNum = n + suffix
	return
}
