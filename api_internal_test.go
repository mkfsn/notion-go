package notion

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithUserAgent(t *testing.T) {
	var settings apiSettings

	userAgent := "test-user-agent"

	WithUserAgent("test-user-agent")(&settings)

	assert.Equal(t, settings.userAgent, userAgent)
}

func TestWithBaseURL(t *testing.T) {
	var settings apiSettings

	baseURL := "https://example.com"

	WithBaseURL(baseURL)(&settings)

	assert.Equal(t, settings.baseURL, baseURL)
}

func TestNotionVersion(t *testing.T) {
	var settings apiSettings

	notionVersion := "2021-05-19"

	WithNotionVersion(notionVersion)(&settings)

	assert.Equal(t, settings.notionVersion, notionVersion)
}

func TestHTTPClient(t *testing.T) {
	var settings apiSettings

	httpClient := &http.Client{}

	WithHTTPClient(httpClient)(&settings)

	assert.Same(t, settings.httpClient, httpClient)
}
