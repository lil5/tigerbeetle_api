package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadFlag(t *testing.T) {
	tests := []struct {
		Fl       string
		Expected []string
	}{
		{"./config.yml", []string{".", "config", "yml"}},
		{"config.yml", []string{".", "config", "yml"}},
		{"/config.yml", []string{"/", "config", "yml"}},
		{"/and/config.yml", []string{"/and", "config", "yml"}},
		{"and/config.yml", []string{"and", "config", "yml"}},
	}

	for i, test := range tests {
		fpath, name, ext := ReadFlag(test.Fl)
		assert.Equal(t, test.Expected[0], fpath, "test number (%d)", i)
		assert.Equal(t, test.Expected[1], name, "test number (%d)", i)
		assert.Equal(t, test.Expected[2], ext, "test number (%d)", i)
	}
}
