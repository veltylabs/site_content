package sitecontent

import "testing"

func TestSentinelErrors(t *testing.T) {
	if ErrNotFound.Error() != "site_content not found" {
		t.Errorf("ErrNotFound text changed: got %q", ErrNotFound.Error())
	}
}
