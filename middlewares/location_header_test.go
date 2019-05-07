package middlewares

// Middleware tests based on https://github.com/unrolled/secure

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/containous/traefik/log"
	"github.com/stretchr/testify/assert"
)

func TestLocationHeaderRegex(t *testing.T) {
	testCases := []struct {
		desc        string
		location    string
		regex       string
		replacement string
		expected    string
	}{
		{
			desc:        "no regex",
			location:    "http://example.com/foo",
			regex:       ``,
			replacement: "",
			expected:    "http://example.com/foo",
		},
		{
			desc:        "single replacement",
			location:    "http://example.com/foo",
			regex:       `(.*)/foo`,
			replacement: "$1/bar",
			expected:    "http://example.com/bar",
		},
		{
			desc:        "scheme replacement",
			location:    "http://example.com/foo",
			regex:       `http://(.*)`,
			replacement: "https://$1",
			expected:    "https://example.com/foo",
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			exp, err := regexp.Compile(strings.TrimSpace(test.regex))
			if err != nil {
				log.Errorf("Error compiling regular expression %s: %s", test.regex, err)
			}

			header := &LocationHeaderStruct{
				LocationRegex:       exp,
				LocationReplacement: test.replacement,
			}

			// Build the recorder
			res := httptest.NewRecorder()

			// Add the location header
			res.HeaderMap.Add("Location", test.location)

			// Generate the resulting response
			response := res.Result()

			// Modify the headers
			header.ModifyLocationHeader(response)

			assert.Equal(t, http.StatusOK, response.StatusCode, "Status not OK")
			assert.Equal(t, test.expected, response.Header.Get("Location"), "Did not get expected header")
		})
	}
}
