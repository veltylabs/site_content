package tests

import (
	"testing"

	sitecontent "github.com/veltylabs/site_content"
)

// The sentinel's text must not change when it moves to the domainError type.
func TestSentinelErrors(t *testing.T) {
	if sitecontent.ErrNotFound.Error() != "site_content not found" {
		t.Errorf("ErrNotFound text changed: got %q", sitecontent.ErrNotFound.Error())
	}
}
