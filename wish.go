package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type WishItem struct {
	UID       string `json:"uid"`
	GachaType string `json:"gacha_type"`
	ItemID    string `json:"item_id"`
	Count     string `json:"count"`
	Time      string `json:"time"`
	Name      string `json:"name"`
	Lang      string `json:"lang"`
	ItemType  string `json:"item_type"`
	RankType  string `json:"rank_type"`
	ID        string `json:"id"`
}

type WishResponse struct {
	Retcode int    `json:"retcode"`
	Message string `json:"message"`
	Data    struct {
		Page   string     `json:"page"`
		Size   string     `json:"size"`
		Total  string     `json:"total"`
		List   []WishItem `json:"list"`
		Region string     `json:"region"`
	} `json:"data"`
}

func (b *App) GetJson(url string) (*WishResponse, error) {
	var client = &http.Client{Timeout: 30 * time.Second}

	retries := 3
	var temperror error
	for retries > 0 {
		//runtime.LogInfo(b.ctx, fmt.Sprintf("fetching retry: %d", retries))
		res, err := client.Get(url)
		if err != nil {
			temperror = err
			retries -= 1
			continue
		}

		defer res.Body.Close()

		//runtime.LogInfo(b.ctx, fmt.Sprintf("fetching status: %d", res.StatusCode))

		if res.StatusCode == http.StatusOK {
			data := &WishResponse{}
			json.NewDecoder(res.Body).Decode(data)

			//runtime.LogInfo(b.ctx, fmt.Sprintf("fetching retcode: %d", data.Retcode))

			if data.Retcode == 0 {
				return data, nil
			} else {
				if data.Message == "authkey error" || data.Message == "authkey timeout" {
					return nil, errors.New("authkey error " + strconv.Itoa(data.Retcode))
				}
			}
		} else {
			temperror = errors.New("error code " + strconv.Itoa(res.StatusCode))
		}

		retries -= 1
	}

	return nil, temperror
}

func (b *App) GetBannerLog(u *url.URL) {
	//runtime.LogError(b.ctx, "Start Get Banner Log")

	banners := []int{100, 200, 301, 302}

	var wishes []string
	for _, e := range banners {
		runtime.EventsEmit(b.ctx, "banner", e)

		res, err := b.GetWishes(u, e)
		if err != nil {
			return
		}

		time.Sleep(2 * time.Second)
		wishes = append(wishes, res...)
	}

	result := "paimonmoe,importer,version,1,0\n"
	result += strings.Join(wishes, "\n")
	err := b.Copy(result)
	if err == nil {
		//runtime.LogInfo(b.ctx, "Copied to clipboard")
		runtime.EventsEmit(b.ctx, "copied")
	}

	err = b.SaveToDownload(result)
	if err == nil {
		//runtime.LogInfo(b.ctx, "Saved to Downloads folder")
		runtime.EventsEmit(b.ctx, "saved")
	}

	runtime.EventsEmit(b.ctx, "result", result)
}

func (b *App) GetWishes(u *url.URL, code int) ([]string, error) {
	var wishes []string

	q := u.Query()
	q.Set("lang", "en")
	q.Set("gacha_type", strconv.Itoa(code))
	q.Set("size", "20")
	q.Set("lang", "en-us")

	page := 1
	lastId := "0"
	var temp []string

	for loop := true; loop; loop = (len(temp) > 0) {
		select {
		case <-b.ch:
			//runtime.LogInfo(b.ctx, "cancelled")
			return nil, errors.New("cancelled")
		default:
		}

		q.Set("page", strconv.Itoa(page))
		q.Set("end_id", lastId)
		u.RawQuery = q.Encode()

		runtime.EventsEmit(b.ctx, "page", page)

		//runtime.LogInfo(b.ctx, fmt.Sprintf("wish banner: %d page: %d\n", code, page))

		res, err := b.GetJson(u.String())
		//runtime.LogInfo(b.ctx, res.Message)

		if err != nil {
			if err.Error() == "authkey error" {
				runtime.EventsEmit(b.ctx, "error", map[string]string{"name": "AUTHKEYERROR", "message": err.Error()})
			} else if err.Error() == "error code" {
				runtime.EventsEmit(b.ctx, "error", map[string]string{"name": "ERRORCODE", "message": err.Error()})
			} else {
				runtime.EventsEmit(b.ctx, "error", map[string]string{"name": "CONNECTIONERROR", "message": err.Error()})
			}
			return nil, err
		}

		//runtime.LogInfo(b.ctx, strconv.Itoa(len(res.Data.List)))

		temp = nil
		if len(res.Data.List) > 0 {
			runtime.EventsEmit(b.ctx, "uid", res.Data.List[0].UID)
			lastId = res.Data.List[len(res.Data.List)-1].ID

			for _, item := range res.Data.List {
				str := fmt.Sprintf("%s,%s,%s,%s,%s", item.GachaType, item.Time, item.Name, item.ItemType, item.RankType)
				temp = append(temp, str)
			}

			wishes = append(wishes, temp...)
			runtime.EventsEmit(b.ctx, "total", map[string]int{"code": code, "total": len(wishes)})
		}

		page += 1
		time.Sleep(time.Second)
	}

	return wishes, nil
}

func (b *App) Cancel() {
	b.ch <- true
}

func (b *App) Copy(text string) error {
	return clipboard.WriteAll(text)
}

func (b *App) SaveToDownload(text string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	path := home + "\\Downloads"
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return err
	}

	t := time.Now()
	f, err := os.Create(path + "\\paimon-moe-import-" + t.Format("20060102150405") + ".csv")
	if err != nil {
		return err
	}

	_, err = f.WriteString(text)
	if err != nil {
		return err
	}

	f.Sync()
	return nil
}
