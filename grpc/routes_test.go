package grpc

import (
	"sync"
	"testing"
	"time"

	"github.com/charithe/timedbuf/v2"
	"github.com/lil5/tigerbeetle_api/config"
	"github.com/stretchr/testify/assert"
)

func TestGetRandomBufferCluster(t *testing.T) {
	n := 100000
	cluster := 3
	flushFunc := func(payloads []TimedPayload) {}
	tbufs := make([]*timedbuf.TimedBuf[TimedPayload], cluster)
	for i := range cluster {
		tbufs[i] = timedbuf.New[TimedPayload](cluster, 2*time.Millisecond, flushFunc)
	}

	// set config
	config.Config.BufferCluster = cluster

	// set app
	a := App{
		TBuf:  tbufs[0],
		TBufs: tbufs,
	}

	// test
	t.Run("sync get random buffer", func(t *testing.T) {
		for range n {
			assert.NotPanics(t, func() {
				b := a.getRandomTBuf()
				assert.NotNil(t, b)
			})
		}
	})

	t.Run("async get random buffer", func(t *testing.T) {
		var wg sync.WaitGroup
		for range n {
			wg.Add(1)
			go func() {
				assert.NotPanics(t, func() {
					b := a.getRandomTBuf()
					assert.NotNil(t, b)
				})
				wg.Done()
			}()
		}
		wg.Wait()
	})
}
