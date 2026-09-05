package app

import (
	"fmt"
	"net"
	"net/http"
	"runtime"
	"syscall"
	"unsafe"

	"github.com/jchv/go-webview2"
	"packetinstall/internal/scanner"
	"packetinstall/internal/web"
)

// RunDesktopApp launches packetinstall as a native Windows desktop window using WebView2 with dark titlebar.
func RunDesktopApp(opts scanner.ScanOptions) error {
	// 1. Allocate an available local port and keep listener active
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return fmt.Errorf("failed to allocate local port: %w", err)
	}
	port := listener.Addr().(*net.TCPAddr).Port

	server := web.NewServer(opts)
	url := fmt.Sprintf("http://127.0.0.1:%d", port)

	// 2. Start HTTP server directly on the established listener
	go func() {
		_ = http.Serve(listener, server.Handler())
	}()

	// 3. Initialize native Windows WebView2 Desktop Window
	w := webview2.NewWithOptions(webview2.WebViewOptions{
		Debug:     false,
		AutoFocus: true,
		WindowOptions: webview2.WindowOptions{
			Title:  "packetinstall — Smart Tools & Agent Management",
			Width:  1360,
			Height: 880,
			Center: true,
		},
	})
	if w == nil {
		return fmt.Errorf("could not create native WebView2 window. Please ensure Microsoft Edge WebView2 runtime is installed")
	}
	defer w.Destroy()

	// 4. Transform native titlebar to dark theme matching app palette
	hwnd := w.Window()
	if hwnd != nil {
		setWindowDarkTitleBar(hwnd)
	}

	w.SetSize(1360, 880, webview2.HintNone)
	w.Navigate(url)
	w.Run()

	return nil
}

// setWindowDarkTitleBar enables DWM immersive dark mode and applies custom dark caption color on Windows.
func setWindowDarkTitleBar(hwnd unsafe.Pointer) {
	if hwnd == nil || runtime.GOOS != "windows" {
		return
	}

	dwmapi := syscall.NewLazyDLL("dwmapi.dll")
	setAttr := dwmapi.NewProc("DwmSetWindowAttribute")

	// 1. DWMWA_USE_IMMERSIVE_DARK_MODE (20 on Windows 11 & Windows 10 20H1+, 19 on older Win10)
	darkMode := int32(1)
	_, _, _ = setAttr.Call(uintptr(hwnd), uintptr(20), uintptr(unsafe.Pointer(&darkMode)), uintptr(unsafe.Sizeof(darkMode)))
	_, _, _ = setAttr.Call(uintptr(hwnd), uintptr(19), uintptr(unsafe.Pointer(&darkMode)), uintptr(unsafe.Sizeof(darkMode)))

	// 2. DWMWA_CAPTION_COLOR = 35 (Windows 11 build 22000+)
	// Color #0b1019 in COLORREF (0x00BBGGRR): R=0x0b, G=0x10, B=0x19 -> 0x0019100B
	captionColor := uint32(0x0019100B)
	_, _, _ = setAttr.Call(uintptr(hwnd), uintptr(35), uintptr(unsafe.Pointer(&captionColor)), uintptr(unsafe.Sizeof(captionColor)))

	// 3. DWMWA_TEXT_COLOR = 36 (Windows 11 build 22000+)
	// Soft white #E2E8F0 -> R=0xE2, G=0xE8, B=0xF0 -> 0x00F0E8E2
	textColor := uint32(0x00F0E8E2)
	_, _, _ = setAttr.Call(uintptr(hwnd), uintptr(36), uintptr(unsafe.Pointer(&textColor)), uintptr(unsafe.Sizeof(textColor)))

	// 4. DWMWA_BORDER_COLOR = 34 (Windows 11 build 22000+)
	borderColor := uint32(0x00162032)
	_, _, _ = setAttr.Call(uintptr(hwnd), uintptr(34), uintptr(unsafe.Pointer(&borderColor)), uintptr(unsafe.Sizeof(borderColor)))

	// 5. Force window frame redraw via SetWindowPos
	user32 := syscall.NewLazyDLL("user32.dll")
	setWindowPos := user32.NewProc("SetWindowPos")
	// SWP_FRAMECHANGED(0x0020) | SWP_NOMOVE(0x0002) | SWP_NOSIZE(0x0001) | SWP_NOZORDER(0x0004) = 0x0027
	_, _, _ = setWindowPos.Call(uintptr(hwnd), 0, 0, 0, 0, 0, uintptr(0x0027))
}
