package viewer

import (
	"context"

	"github.com/chromedp/chromedp"
	log "github.com/sirupsen/logrus"
)

func skipAd(ctx *context.Context, text *string) (bool, error) {
	log.Trace("Trying to skip ad!")
	if *text == "1" || *text == "0" {
		log.Trace("Skipping ad...")
		err := chromedp.Run(*ctx,
			chromedp.WaitVisible(`.ytp-ad-skip-button`),
			chromedp.Click(`.ytp-ad-skip-button`, chromedp.NodeEnabled),
		)
		if err != nil {
			return false, err
		}
		return true, nil
	}
	return false, nil
}

func handleAds(ctx *context.Context, adCreated, adShowing bool) {
	if !adShowing {
		log.Trace("No ads to skip!")
		return
	}
	for {
		var previewText string
		_ = chromedp.Run(*ctx,
			chromedp.WaitVisible(`.ytp-ad-preview-text`),
			chromedp.Text(`.ytp-ad-preview-text`, &previewText),
		)
		if len(previewText) > 1 {
			log.Info("Ad is not skippable")
			_ = chromedp.Run(*ctx,
				chromedp.WaitNotPresent(`.ad-showing`),
			)
			break
		}
		skipped, err := skipAd(ctx, &previewText)
		if err != nil {
			log.Error(err)
			log.Warn("Could not skip ad")
			break
		} else if skipped {
			break
		}
	}
}
