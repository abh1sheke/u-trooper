package viewer

import (
	"context"

	"github.com/chromedp/chromedp"
)

type retrier struct {
	count int
}

func (r *retrier) inc() {
	r.count += 1
}

func (r *retrier) reset() {
	r.count = 0
}

func handleConsentDialogue(ctx *context.Context) error {
	var dialogueExists bool
	err := chromedp.Run(*ctx,
		chromedp.WaitVisible("#container > .html5-video-player"),
		chromedp.EvaluateAsDevTools(
			`document.querySelectorAll(".eom-v1-dialog").length >= 1`,
			&dialogueExists,
		),
	)
	if err != nil {
		return err
	}
	if dialogueExists {
		err = chromedp.Run(*ctx,
			chromedp.WaitVisible(".yt-spec-button-shape-next--filled"),
			chromedp.EvaluateAsDevTools(
				`document.querySelectorAll(".yt-spec-button-shape-next--filled")[1].click()`,
				nil,
			),
		)
		if err != nil {
			return err
		}
	}
	return err
}

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

func newChromeCtx(proxy *string) (context.Context, context.CancelFunc) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.ProxyServer(*proxy),
		chromedp.Flag("proxy-bypass-list", "<-loopback>"),
	)
	alloc, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	ctx, cancel := chromedp.NewContext(
		alloc,
		chromedp.WithErrorf(logError),
		// chromedp.WithDebugf(logDebug),
	)
	return ctx, cancel
}
