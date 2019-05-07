package middlewares

// Middleware based on https://github.com/unrolled/secure

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/containous/traefik/log"
	"github.com/containous/traefik/types"
)

const (
	// LocationHeader is the key within the response context used to
	// access the location header
	LocationHeader = "Location"
	// ReplacedLocationHeader is the header to set the old Location to
	ReplacedLocationHeader = "X-Replaced-Location"
)

// LocationHeaderStruct is a middleware that helps setup a few basic security features. A single headerOptions struct can be
// provided to configure which features should be enabled, and the ability to override a few of the default values.
type LocationHeaderStruct struct {
	LocationRegex       *regexp.Regexp
	LocationReplacement string
}

// NewLocationHeaderFromStruct constructs a new header instance from supplied frontend header struct.
func NewLocationHeaderFromStruct(headers *types.Headers) *LocationHeaderStruct {
	if headers == nil || !headers.HasLocationHeaderDefined() {
		return nil
	}

	exp, err := regexp.Compile(strings.TrimSpace(headers.LocationRegex))
	if err != nil {
		log.Errorf("Error compiling regular expression %s: %s", headers.LocationRegex, err)
	}

	return &LocationHeaderStruct{
		LocationRegex:       exp,
		LocationReplacement: strings.TrimSpace(headers.LocationReplacement),
	}
}

// ModifyLocationHeader modifies the Location header
func (s *LocationHeaderStruct) ModifyLocationHeader(res *http.Response) error {
	// Handle the Location header regex
	if s.LocationRegex != nil && len(s.LocationReplacement) > 0 {
		locationValue := res.Header.Get(LocationHeader)
		if len(locationValue) > 0 {
			// If there is a location regex run it against the location
			res.Header.Set(LocationHeader, s.LocationRegex.ReplaceAllString(locationValue, s.LocationReplacement))

			// Store the original value in the replacement location
			res.Header.Set(ReplacedLocationHeader, locationValue)
		}
	}
	return nil
}
