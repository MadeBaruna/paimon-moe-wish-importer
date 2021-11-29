package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App application struct
type App struct {
	ctx context.Context
	ch  chan bool
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		ch: make(chan bool),
	}
}

// startup is called at application startup
func (b *App) startup(ctx context.Context) {
	// Perform your setup here
	b.ctx = ctx
}

// domReady is called after the front-end dom has been loaded
func (b *App) domReady(ctx context.Context) {
	// runtime.WindowSetSize(b.ctx, 500, 500)
}

// shutdown is called at application termination
func (b *App) shutdown(ctx context.Context) {
	// Perform your teardown here
}

func (b *App) GetLog(server string) {
	var loc string
	if server == "global" {
		loc = "$userprofile\\AppData\\LocalLow\\miHoYo\\Genshin Impact\\output_log.txt"
	} else {
		loc = "$userprofile\\AppData\\LocalLow\\miHoYo\\原神\\output_log.txt"
	}

	fmt.Println(loc)

	path := os.ExpandEnv(loc)
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		//runtime.LogInfo(b.ctx, "Log file does not exists")
		runtime.EventsEmit(b.ctx, "error", map[string]string{"name": "LOGNOTFOUND"})
		return
	}

	f, err := os.OpenFile(path, os.O_RDONLY, os.ModePerm)
	if err != nil {
		//runtime.LogError(b.ctx, err.Error())
		runtime.EventsEmit(b.ctx, "error", map[string]string{"name": "ERRORREADLOG", "message": err.Error()})
		return
	}
	defer f.Close()
	sc := bufio.NewScanner(f)

	r := regexp.MustCompile(`^OnGetWebViewPageFinish.*log$`)
	var links []string
	for sc.Scan() {
		if str := sc.Text(); r.MatchString(str) {
			links = append(links, str)
		}
	}

	if err := sc.Err(); err != nil {
		runtime.EventsEmit(b.ctx, "error", map[string]string{"name": "ERRORREADLOG", "message": err.Error()})
		return
	}

	if len(links) == 0 {
		//runtime.LogError(b.ctx, "Cannot found wish history link")
		runtime.EventsEmit(b.ctx, "error", map[string]string{"name": "LINKNOTFOUND"})
		return
	}

	selected := strings.Replace(links[len(links)-1], "OnGetWebViewPageFinish:", "", 1)
	u, err := GetUrl(selected, server)
	if err != nil {
		//runtime.LogError(b.ctx, "Invalid link")
		runtime.EventsEmit(b.ctx, "error", map[string]string{"name": "LINKINVALID"})
		return
	}

	//runtime.LogInfo(b.ctx, "Wish History link found")
	//runtime.LogInfo(b.ctx, u.String())

	b.GetBannerLog(u)
}

func (b *App) Start(server string) {
	go b.GetLog(server)
}

func (b *App) ResizeWindow(width, height int) {
	runtime.WindowSetSize(b.ctx, width, height)
}
