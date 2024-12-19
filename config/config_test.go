package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	t.Run("Empty TB_ADDRESSES", func(t *testing.T) {
		os.Unsetenv("TB_ADDRESSES")
		assert.False(t, NewConfig())
	})

	t.Run("Contains TB_ADDRESSES", func(t *testing.T) {
		os.Setenv("TB_ADDRESSES", "127.0.0.1:3033")
		assert.True(t, NewConfig())
	})

	t.Run("Buffered cluster", func(t *testing.T) {
		os.Setenv("TB_ADDRESSES", "127.0.0.1:3033")
		os.Setenv("IS_BUFFERED", "true")
		os.Setenv("BUFFER_SIZE", "100")
		os.Setenv("BUFFER_DELAY", "100ms")
		os.Setenv("BUFFER_CLUSTER", "1")
		assert.True(t, NewConfig())
		assert.True(t, Config.IsBuffered)
		assert.Equal(t, 100, Config.BufferSize)
		{
			d, _ := time.ParseDuration("100ms")
			assert.Equal(t, d, Config.BufferDelay)
		}
		assert.Equal(t, 1, Config.BufferCluster)
	})
}
