package viewer

import (
	"context"
	"fmt"
	"time"

	"github.com/chromedp/chromedp"
	log "github.com/sirupsen/logrus"
)

const SkipButton string = `ytp-ad-skip-button`

func handleAds(ctx *context.Context) {
	start := time.Now()
	for {
		if time.Since(start) >= (2 * time.Minute) {
			log.Fatal("Ad timeout reached. Quitting.")
		}
		var adShowing, skipExists bool
		err := chromedp.Run(*ctx,
			chromedp.EvaluateAsDevTools(
				`document.getElementsByClassName("ad-showing").length >= 1`,
				&adShowing,
			),
			chromedp.EvaluateAsDevTools(
				fmt.Sprintf(
					`document.getElementsByClassName("%s").length >= 1`,
					SkipButton,
				), &skipExists,
			),
		)
		if err != nil {
			log.Error(err)
			continue
		}
		if !adShowing {
			break
		}
		if skipExists {
			log.Trace("Skipping ad...")
			err := chromedp.Run(*ctx,
				chromedp.WaitVisible("."+SkipButton),
				chromedp.Click("."+SkipButton, chromedp.NodeEnabled),
			)
			if err != nil {
				log.Error(err)
				continue
			}
			break
		}
	}
	log.Info("Ads have been skipped.")
}
