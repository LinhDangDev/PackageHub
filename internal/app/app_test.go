package app

import (
	"testing"

	"github.com/jchv/go-webview2"
)

func TestWebView2Init(t *testing.T) {
	w := webview2.NewWithOptions(webview2.WebViewOptions{
		Debug: false,
		WindowOptions: webview2.WindowOptions{
			Title:  "Test",
			Width:  400,
			Height: 300,
		},
	})
	if w == nil {
		t.Skip("WebView2 not available in this headless test environment")
		return
	}
	defer w.Destroy()
}
