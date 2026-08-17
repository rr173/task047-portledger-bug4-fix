package registry

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestProbeConcurrentSnapshotSafe(t *testing.T) {
	r := New()
	const workers = 32
	const rounds = 20
	var wg sync.WaitGroup
	start := make(chan struct{})
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			<-start
			for j := 0; j < rounds; j++ {
				ts := time.Date(2026, 8, 16, 10, j+1, 0, 0, time.UTC)
				_, _ = r.Submit(fmt.Sprintf("host-%d", i), ts, []int{80 + i})
			}
		}(i)
	}
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			for j := 0; j < rounds*2; j++ {
				r.Snapshot(5)
			}
		}()
	}
	close(start)
	wg.Wait()
}
