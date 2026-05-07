package flags

import (
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/digitalghost-dev/poke-cli/cmd/utils"
	"github.com/digitalghost-dev/poke-cli/styling"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func latestReleaseTestClient(statusCode int, body string) *http.Client {
	return &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: statusCode,
				Header:     make(http.Header),
				Body:       io.NopCloser(strings.NewReader(body)),
				Request:    req,
			}, nil
		}),
	}
}

func TestLatestVersionFlag(t *testing.T) {
	err := os.Setenv("GO_TESTING", "1")
	if err != nil {
		t.Fatalf("Failed to set GO_TESTING env var: %v", err)
	}

	defer func() {
		err := os.Unsetenv("GO_TESTING")
		if err != nil {
			t.Logf("Warning: failed to unset GO_TESTING: %v", err)
		}
	}()

	originalClient := latestReleaseHTTPClient
	latestReleaseHTTPClient = latestReleaseTestClient(http.StatusOK, `{"tag_name":"v1.10.1"}`)
	defer func() { latestReleaseHTTPClient = originalClient }()

	tests := []struct {
		name           string
		args           []string
		expectedOutput string
		expectedError  bool
	}{
		{
			name:           "Get latest version with short flag",
			args:           []string{"-l"},
			expectedOutput: utils.LoadGolden(t, "main_latest_flag.golden"),
			expectedError:  false,
		},
		{
			name:           "Get latest version with long flag",
			args:           []string{"--latest"},
			expectedOutput: utils.LoadGolden(t, "main_latest_flag.golden"),
			expectedError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalArgs := os.Args
			os.Args = append([]string{"poke-cli"}, tt.args...)
			defer func() { os.Args = originalArgs }()

			output, err := LatestFlag()
			cleanOutput := styling.StripANSI(output)

			if tt.expectedError {
				require.Error(t, err, "Expected an error")
			} else {
				require.NoError(t, err, "Expected no error")
			}

			assert.Equal(t, tt.expectedOutput, cleanOutput, "Output should match expected")
		})
	}
}

func TestLatestReleaseFromURL(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		body       string
		wantErr    bool
		contains   string
		outputHas  string
	}{
		{
			name:       "valid latest release",
			statusCode: http.StatusOK,
			body:       `{"tag_name":"v1.2.3"}`,
			outputHas:  "v1.2.3",
		},
		{
			name:       "non-200 response",
			statusCode: http.StatusForbidden,
			body:       "rate limit",
			wantErr:    true,
			contains:   "unexpected GitHub response status: 403",
		},
		{
			name:       "empty release tag",
			statusCode: http.StatusOK,
			body:       `{"tag_name":""}`,
			wantErr:    true,
			contains:   "did not include a tag name",
		},
		{
			name:       "response body too large",
			statusCode: http.StatusOK,
			body:       strings.Repeat("x", maxLatestReleaseBytes+1),
			wantErr:    true,
			contains:   "response body exceeds",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var output strings.Builder
			err := latestReleaseFromURL(
				&output,
				"https://api.github.com/repos/digitalghost-dev/poke-cli/releases/latest",
				latestReleaseTestClient(tt.statusCode, tt.body),
			)
			cleanOutput := styling.StripANSI(output.String())

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.contains)
				assert.Contains(t, cleanOutput, tt.contains)
				return
			}

			require.NoError(t, err)
			assert.Contains(t, cleanOutput, tt.outputHas)
		})
	}
}
