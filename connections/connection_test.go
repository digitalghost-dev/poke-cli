package connections

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBaseApiCallSuccess(t *testing.T) {
	expectedData := map[string]string{"key": "value"}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(expectedData)
		assert.Nil(t, err)
	}))
	defer ts.Close()

	var target map[string]string

	baseApiCall(ts.URL, &target)

	assert.Equal(t, expectedData, target)
}
