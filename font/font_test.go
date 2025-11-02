package font

import (
	"testing"

	"github.com/dnsoa/go/assert"
)

func TestXxx(t *testing.T) {
	r := assert.New(t)
	f, err := GetFontFamily(OpenSansSans40TTF)
	r.NoError(err)
	t.Logf("%v", f)
}
