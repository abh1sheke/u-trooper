package viewer

import (
	"context"
	"net/url"

	"github.com/chromedp/chromedp"
	log "github.com/sirupsen/logrus"
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
	uri := *proxy
	if _, err := url.ParseRequestURI(uri); err != nil {
		log.WithFields(log.Fields{"url": uri, "reason": err}).
			Warn("Could not parse given proxy url")
		uri = "socks5://127.0.0.1:9050"
	}
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
		chromedp.ProxyServer(uri),
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
