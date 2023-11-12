package viewer

import (
	"strings"
	"time"

	"github.com/abh1sheke/utrooper/tor"
	"github.com/chromedp/chromedp"
	log "github.com/sirupsen/logrus"
)

func view(args *viewerArgs) {
	defer args.wg.Done()
	retries := new(retrier)
	for {
		if retries.count > 5 {
			log.WithField("retries", retries.count).Fatalf("Too many retries.")
		}
		count := args.viewCount.Load()
		if count == uint64(args.target) {
			return
		}
		if args.proxy.isTor && count > 1 && count%5 == 0 {
			err := tor.Restart(args.mu)
			if err != nil {
				log.WithField("reason", err).
					Warn("Could not restart tor. IP has not been changed.")
			}
		}
		proxy := args.proxy.url.String()
		ctx, cancel := newChromeCtx(&proxy)
		var classNames string
		var ok bool
		err := chromedp.Run(ctx, chromedp.Navigate(args.url))
		if err != nil {
			log.WithField("reason", err).Error("Could not play video.")
			retries.inc()
			continue
		} else {
			err := chromedp.Run(ctx, getPlayerClasses(&classNames, &ok))
			if err != nil {
				log.WithField("reason", err).Error("Could not get player classes.")
				retries.inc()
				continue
			}
		}
		if !ok {
			log.Error("Could not fetch video play-state.")
			retries.inc()
			continue
		}
		pausedMode := strings.Contains(classNames, "paused-mode")
		adCreated := strings.Contains(classNames, "ad-created")
		unstarted := strings.Contains(classNames, "unstarted-mode")
		err = handleConsentDialogue(&ctx)
		if err != nil {
			log.WithField("reason", err).Error("Could not handle consent dialog")
			retries.inc()
			continue
		}
		if pausedMode || unstarted {
			err := chromedp.Run(ctx, playVideo())
			if err != nil {
				log.WithField("reason", err).Error("Could not play video.")
				retries.inc()
				continue
			}
		}
		if adCreated {
			log.Info("Handling ads...")
			handleAds(&ctx)
		}
		err = chromedp.Run(ctx, chromedp.Sleep(args.duration))
		if err != nil {
			log.WithField("reason", err).Errorf("Could not watch video for %v.", args.duration)
			retries.inc()
			continue
		} else {
			args.viewCount.Add(1)
			log.WithFields(
				log.Fields{
					"views": args.viewCount.Load(),
					"tss":   time.Since(args.start),
				},
			).Info("View action completed.")
			retries.reset()
			cancel()
		}
	}
}
