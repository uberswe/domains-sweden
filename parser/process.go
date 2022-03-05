package parser

import (
	"bytes"
	"context"
	"fmt"
	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/esimov/stackblur-go"
	"github.com/nfnt/resize"
	"image"
	"image/jpeg"
	"log"
	"sync"
	"time"
)

func (p *Parser) process(url string) (content string, requestSize float64, responseTime float64, screenshot []byte, blurredScreenshot []byte, events []Event, requested time.Time, err error) {
	// create chrome instance
	ctx, cancel := chromedp.NewContext(
		context.Background(),
		chromedp.WithLogf(log.Printf),
	)
	// TODO we should maybe ignore ssl errors
	defer cancel()
	size := int64(0)

	var ws sync.WaitGroup

	chromedp.ListenTarget(ctx, func(v interface{}) {
		switch v.(type) {
		case *network.EventRequestWillBeSent:
			ws.Add(1)
			go func(r *network.EventRequestWillBeSent) {
				events = append(events, Event{
					URL:  r.Request.URL,
					Type: "request",
					Time: time.Now(),
				})
				ws.Done()
			}(v.(*network.EventRequestWillBeSent))
			break

		case *network.EventResponseReceived:
			ws.Add(1)
			go func(r *network.EventResponseReceived) {
				events = append(events, Event{
					URL:  r.Response.URL,
					Type: "response",
					Time: time.Now(),
				})
				ws.Done()
			}(v.(*network.EventResponseReceived))
			break
		case *network.EventDataReceived:
			// Fired when data chunk was received over the network.
			ws.Add(1)
			go func() {
				edr := v.(*network.EventDataReceived)
				//log.Printf("Data Received : %d\n", edr.DataLength)
				size += edr.DataLength
				ws.Done()
			}()
		default:
			//log.Println(reflect.TypeOf(v).Elem())

			// case *network.EventLoadingFinished:
			// 	go func() {
			// 		lf := v.(*network.EventLoadingFinished)
			// 		log.Printf("Loading finished : %f\n", lf.EncodedDataLength)
			// 	}()
			// case *network.EventLoadingFailed:
			// 	// Fired when HTTP request has failed to load.
			// 	go func() {
			// 		lf := v.(*network.EventLoadingFailed)
			// 		log.Printf("Loading finished : %s\n", lf.ErrorText)
			// 	}()
		}
	})

	// create a timeout
	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	start := time.Now()
	// navigate to a page, wait for an element, click
	err = chromedp.Run(ctx,
		chromedp.EmulateViewport(1920, 1080),
		emulation.SetUserAgentOverride(fmt.Sprintf("Domaner.xyz Spider Bot - Please contact support@domaner.xyz regarding any abuse or problem. Visit https://www.domaner.xyz for more information.")),
		chromedp.Navigate(url),
		chromedp.Sleep(10*time.Second),
		// wait for footer element is visible (ie, page is loaded)
		chromedp.FullScreenshot(&screenshot, 5),
		chromedp.ActionFunc(func(ctx context.Context) error {
			node, err := dom.GetDocument().Do(ctx)
			if err != nil {
				return err
			}
			content, err = dom.GetOuterHTML().WithNodeID(node.NodeID).Do(ctx)
			return err
		}),
	)
	if err != nil {
		return "", 0, 0, nil, nil, nil, start, err
	}

	img, _, err2 := image.Decode(bytes.NewReader(screenshot))

	if err2 != nil {
		return "", 0, 0, nil, nil, nil, start, err2
	}

	tmpImg := img.(interface {
		SubImage(r image.Rectangle) image.Image
	}).SubImage(image.Rect(0, 0, 1920, 1080))

	thumbnail := resize.Resize(960, 540, tmpImg, resize.Lanczos3)

	blurredThumbnail, err3 := stackblur.Run(thumbnail, 5)
	if err3 != nil {
		return "", 0, 0, nil, nil, nil, start, err3
	}
	//for _, e := range events {
	//	log.Printf("Event %s at %s: %s\n", e.Type, e.Time.Format("04:05.999999999"), e.URL)
	//}

	buf := new(bytes.Buffer)
	err4 := jpeg.Encode(buf, thumbnail, nil)
	if err4 != nil {
		return "", 0, 0, nil, nil, nil, start, err4
	}

	blurBuf := new(bytes.Buffer)
	err5 := jpeg.Encode(blurBuf, blurredThumbnail, nil)
	if err5 != nil {
		return "", 0, 0, nil, nil, nil, start, err5
	}

	requestSize = float64(size) / 1024.0 / 1024.0
	responseTime = time.Since(start.Add(10 * time.Second)).Seconds()
	return content, requestSize, responseTime, buf.Bytes(), blurBuf.Bytes(), events, start, nil
}
