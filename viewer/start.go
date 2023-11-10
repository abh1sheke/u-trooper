package viewer

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/abh1sheke/utrooper/tor"
	log "github.com/sirupsen/logrus"
)

type viewerArgs struct {
	duration  time.Duration
	mu        *sync.Mutex
	start     time.Time
	target    int
	url       string
	viewCount *atomic.Uint64
	wg        *sync.WaitGroup
	proxy     string
}

func newViewerArgs(url, proxy string, target, dur int) *viewerArgs {
	mu := new(sync.Mutex)
	wg := new(sync.WaitGroup)
	viewCount := new(atomic.Uint64)
	start := time.Now()
	duration := time.Duration(dur) * time.Second

	return &viewerArgs{duration, mu, start, target, url, viewCount, wg, proxy}
}

func StartViewing(views, instances, duration int, url, proxy *string) {
	args := newViewerArgs(*url, *proxy, views, duration)
	err := tor.Start(args.mu)
	if err != nil {
		log.WithField("reason", err).Fatal("Could not start tor service")
	}
	defer func() {
		err := tor.Stop(args.mu)
		if err != nil {
			log.WithField("reason", err).Fatal("Could not stop tor service")
		}
	}()
	for i := 0; i < instances; i++ {
		args.wg.Add(1)
		go view(args)
	}
	args.wg.Wait()
	log.WithFields(log.Fields{
		"took":  fmt.Sprintf("%vms", time.Since(args.start)),
		"views": args.viewCount.Load(),
	}).Info("Finished operation.")
}
