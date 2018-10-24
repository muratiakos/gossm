package printer

import (
	"bytes"
	"github.com/fatih/color"
	"github.com/glassechidna/gossm/pkg/gossm"
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestPrintWrapped(t *testing.T) {
	buf := &bytes.Buffer{}
	p := New()
	c := color.New(color.FgBlue)

	p.print(buf, c, gossm.SsmMessage{}, "hello world\n")
	str := buf.String()
	expected := "\x1b[34m[] \x1b[0mhello world\n\x1b[34m[] \x1b[0m\n"
	assert.Equal(t, str, expected)
}
