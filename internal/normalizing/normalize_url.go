package normalizing

import (
	"net/url"
)

func NormalizeURL(url_ string) (string, error) {
	parsed, err := url.Parse(url_)
	if err != nil {
		return "", err
	}
	return parsed.Host + parsed.Path, nil
}
