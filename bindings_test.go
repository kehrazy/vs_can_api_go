package vscanapi_test

import (
	"testing"

	vscanapi "github.com/kehrazy/vs_can_api_go"
)

func TestApiVersion(t *testing.T) {
	_, err := vscanapi.GetApiVersion()
	if err != nil {
		t.Fatal(err)
	}
}
