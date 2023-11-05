package viewer

import (
	"context"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	log "github.com/sirupsen/logrus"
)

func playVideo(ctx *context.Context) error {
	return chromedp.Run(*ctx,
		chromedp.Click(`.ytp-play-btn`, chromedp.NodeVisible),
	)
}

func getPlayerClasses(ctx *context.Context, classNames *string, ok *bool) error {
	return chromedp.Run(*ctx,
		chromedp.WaitVisible(`#container > .html5-video-player`),
		chromedp.AttributeValue(`#container > .html5-video-player`, "class", classNames, ok),
	)
}

func View(url *string) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
	)
	allocCtx, cancel := chromedp.NewExecAllocator(
		context.Background(),
		opts...,
	)

	ctx, cancel := chromedp.NewContext(
		allocCtx,
		chromedp.WithErrorf(logError),
		// chromedp.WithDebugf(logDebug),
	)
	defer cancel()

	var classNames string
	var ok bool
	err := chromedp.Run(ctx,
		chromedp.Navigate(*url),
	)
	if err != nil {
		log.WithField("reason", err).Error("Could not play video!")
		return
	} else {
		getPlayerClasses(&ctx, &classNames, &ok)
	}

	if ok {
		pausedMode := strings.Contains(classNames, "paused-mode")
		adCreated := strings.Contains(classNames, "ad-created")
		adShowing := strings.Contains(classNames, "ad-showing")

		if pausedMode {
			err := playVideo(&ctx)
			if err != nil {
				log.WithField("reason", err).Error("Could not play video!")
				return
			}
		}

		if adCreated {
			log.Info("Handling ads...")
			for adShowing {
				handleAds(&ctx, adCreated, adShowing)
				if err := getPlayerClasses(&ctx, &classNames, &ok); err != nil {
					log.Error(err)
					return
				}
				adShowing = strings.Contains(classNames, "ad-showing")
			}
			log.Info("Ad(s) has been skipped!")
		}

		err := chromedp.Run(ctx, chromedp.Sleep(50*time.Second))
		if err != nil {
			log.Error(err)
			return
		}
	} else {
		log.Error("Could not fetch video play-state")
		return
	}
}
