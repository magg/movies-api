package movies

import (
	"net/http"
	"strings"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestSearchMovies(t *testing.T) {
	cases := []struct {
		name                string
		mockResponseBody    string
		expectedMovies      []Movie
		expectedErrorString string
	}{
		{
			name:             "RegularCase",
			mockResponseBody: `{"Search":[{"Title":"Star Wars: A New Hope","Year":"1977"},{"Title":"Star Wars: The Empire Strikes Back","Year":"1980"}]}`,
			expectedMovies: []Movie{
				{Title: "Star Wars: A New Hope", Year: "1977"},
				{Title: "Star Wars: The Empire Strikes Back", Year: "1980"},
			},
			expectedErrorString: "",
		},
		{
			name:             "SortedCase",
			mockResponseBody: `{"Search":[{"Title":"Star Wars: The Empire Strikes Back","Year":"1980"}, {"Title":"Star Wars: A New Hope","Year":"1977"}]}`,
			expectedMovies: []Movie{
				{Title: "Star Wars: A New Hope", Year: "1977"},
				{Title: "Star Wars: The Empire Strikes Back", Year: "1980"},
			},
			expectedErrorString: "",
		},
		{
			name:             "SortedCaseAlphanumeric",
			mockResponseBody: `{"Search":[{"Title":"Star Wars: The Empire Strikes Back","Year":"1977"}, {"Title":"Star Wars: A New Hope","Year":"1977"}]}`,
			expectedMovies: []Movie{
				{Title: "Star Wars: A New Hope", Year: "1977"},
				{Title: "Star Wars: The Empire Strikes Back", Year: "1977"},
			},
			expectedErrorString: "",
		},

		{
			name:                "NotFoundError",
			mockResponseBody:    `{"Response":"False","Error":"Movie not found!"}`,
			expectedMovies:      nil,
			expectedErrorString: "Movie not found!",
		},
		{
			name:                "JSONParseError",
			mockResponseBody:    "",
			expectedMovies:      nil,
			expectedErrorString: "unexpected end of JSON input",
		},

		{
			name:                "JSONInvalidError",
			mockResponseBody:    `{Response":"False","Error":"Movie not found!"}`,
			expectedMovies:      nil,
			expectedErrorString: "invalid character 'R' looking for beginning of object key string",
		},
	}

	searcher := &APIMovieSearcher{
		URL:    "http://example.com/",
		APIKey: "mock-api-key",
	}

	for _, c := range cases {
		// register http mock
		httpmock.RegisterResponder(
			"GET",
			"http://example.com/",
			func(req *http.Request) (*http.Response, error) {
				return httpmock.NewStringResponse(200, c.mockResponseBody), nil
			},
		)
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		// run test
		t.Run(c.name, func(t *testing.T) {

			var sort bool = false
			if strings.Contains(c.name, "Sort") {
				sort = true
			}

			actualMovies, actualError := searcher.SearchMovies("star wars", sort)
			assert.EqualValues(t, c.expectedMovies, actualMovies)
			if c.expectedErrorString == "" {
				assert.NoError(t, actualError)
			} else {
				assert.EqualError(t, actualError, c.expectedErrorString)
			}
		})
	}
}
