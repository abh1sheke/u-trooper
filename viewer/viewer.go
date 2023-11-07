package viewer

import (
	"context"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	log "github.com/sirupsen/logrus"
)

func playVideo() chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.WaitVisible(".ytp-play-button"),
		chromedp.Click(".ytp-play-button"),
	}
}

func getPlayerClasses(classNames *string, ok *bool) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.WaitVisible(`#container > .html5-video-player`),
		chromedp.AttributeValue(`#container > .html5-video-player`, "class", classNames, ok),
	}
}

func newChromeCtx() (context.Context, context.CancelFunc) {
	ctx, cancel := chromedp.NewContext(
		context.Background(),
		chromedp.WithErrorf(logError),
		// chromedp.WithDebugf(logDebug),
	)
	return ctx, cancel
}

func view(args *viewerArgs) {
	defer args.wg.Done()
	for {
		if count := args.viewCount.Load(); count == uint64(args.target) {
			return
		}
		ctx, cancel := newChromeCtx()
		var classNames string
		var ok bool
		err := chromedp.Run(ctx, chromedp.Navigate(args.url))
		if err != nil {
			log.WithField("reason", err).Error("Could not play video.")
			return
		} else {
			err := chromedp.Run(ctx, getPlayerClasses(&classNames, &ok))
			if err != nil {
				log.WithField("reason", err).Error("Could not get player classes.")
				return
			}
		}
		if !ok {
			log.Error("Could not fetch video play-state.")
			return
		}
		pausedMode := strings.Contains(classNames, "paused-mode")
		adCreated := strings.Contains(classNames, "ad-created")
		unstarted := strings.Contains(classNames, "unstarted-mode")
		if pausedMode || unstarted {
			err := chromedp.Run(ctx, playVideo())
			if err != nil {
				log.WithField("reason", err).Error("Could not play video.")
				return
			}
		}
		if adCreated {
			log.Info("Handling ads...")
			handleAds(&ctx)
		}
		err = chromedp.Run(ctx, chromedp.Sleep(args.duration))
		if err != nil {
			log.WithField("reason", err).Errorf("Could not watch video for %v.", args.duration)
			return
		} else {
			args.viewCount.Add(1)
			log.WithFields(
				log.Fields{
					"views": args.viewCount.Load(),
					"tss":   time.Since(args.start),
				},
			).Info("View action completed.")
			cancel()
		}
	}
}
