package nyaa

import (
	"bufio"
	"context"
	"errors"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/anacrolix/torrent"
	"github.com/anicine/anicine-torrent/client"
	"github.com/anicine/anicine-torrent/internal/analyze"
	"github.com/anicine/anicine-torrent/models"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrMissingSeeders = errors.New("missing seeders")
)

func Load(path string) ([]string, error) {
	_, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var (
		trackers []string
		scanner  = bufio.NewScanner(file)
	)

	for scanner.Scan() {
		txt := strings.TrimSpace(scanner.Text())
		if txt == "" || txt == "\n" || txt == "\t" {
			continue
		}

		uri, err := url.Parse(txt)
		if err != nil {
			return nil, err
		}

		trackers = append(trackers, uri.String())
	}

	return trackers, nil
}

func (s *Nyaa) extract(ctx context.Context, uri string) ([]*torrent.File, error) {
	magnet, err := s.client.AddMagnet(uri)
	if err != nil {
		return nil, err
	}

	magnet.AddTrackers([][]string{s.trackers})
	magnet.DisallowDataUpload()
	magnet.DisallowDataDownload()

	now := time.Now()
	for {
		select {
		case <-ctx.Done():
			return nil, context.Canceled
		case <-magnet.GotInfo():
			files := magnet.Files()
			return files, err
		default:
			if time.Since(now) > time.Second*90 {
				return nil, ErrMissingSeeders
			}
		}
	}
}

type page struct {
	Total      int
	Torrents   []*models.TorrentData
	PageNumber int
	NextPage   bool
}

func (s *Nyaa) browse(ctx context.Context, args *client.Args) (*page, error) {
	if args == nil {
		return nil, errors.New("missing args")
	}

	args.Headers = map[string]string{
		"accept-language": "en-US,en;q=0.9",
		"sec-fetch-site":  "same-origin",
		"referer":         args.Endpoint.String(),
	}

	resp, err := client.Do(ctx, *args)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	page := new(page)
	if num := analyze.ExtractParentheses(doc.Find("div.row").Find("h3").Text()); len(num) > 0 {
		page.Total, err = strconv.Atoi(analyze.CleanSymbols(num[len(num)-1]))
		if err != nil {
			s.log.Error("cannot get the total numbers of torrents", "error", err)
		}
	}
	doc.Find("tbody").Find("tr").Each(func(i int, t *goquery.Selection) {
		torrent := new(models.TorrentData)
		torrent.TorrentInfo = new(models.TorrentInfo)
		t.Find("td").Each(func(i int, s *goquery.Selection) {
			switch i {
			case 1:
				s.Find("a").Each(func(_ int, v *goquery.Selection) {
					if href, ok := v.Attr("href"); ok {
						if _, ok = v.Attr("class"); !ok {
							torrent.Name = strings.TrimSpace(v.Text())
							torrent.NyaaID = analyze.ExtractNyaaID(href + "/")
						}
					}
				})
			case 2:
				s.Find("a").Each(func(_ int, v *goquery.Selection) {
					if href, ok := v.Attr("href"); ok {
						if !strings.Contains(href, "/download/") {
							torrent.Hash = analyze.ExtractMagnetHash(href)
						}
					}
				})
			case 3:
				size := strings.ToLower(strings.ReplaceAll(s.Text(), " ", ""))
				if strings.HasSuffix(size, "gib") {
					torrent.Size, _ = analyze.ExtractMiB(size[:3])
				} else if strings.HasSuffix(size, "mib") {
					torrent.Size, _ = strconv.ParseFloat(size[:3], 64)
				}
			case 4:
				if timestamp, ok := s.Attr("data-timestamp"); ok {
					torrent.Date, _ = strconv.ParseInt(timestamp, 10, 64)
				}
			case 5:
				torrent.Seeders, _ = strconv.Atoi(s.Text())
			case 6:
				torrent.Leechers, _ = strconv.Atoi(s.Text())
			}
		})
		page.Torrents = append(page.Torrents, torrent)
	})

	var active bool
	doc.Find("ul.pagination").Find("li").Each(func(i int, s *goquery.Selection) {
		if class, ok := s.Attr("class"); ok {
			if strings.Contains(class, "active") {
				active = true
				txt := s.Find("a").Text()
				page.PageNumber = analyze.ExtractNum(txt)
			}
			return
		}
		txt := s.Find("a").Text()
		num := analyze.ExtractNum(txt)
		if active && num > page.PageNumber {
			page.NextPage = true
		}
	})

	return page, nil
}

func (s *Nyaa) scrape(ctx context.Context, collection *mongo.Collection, cfg *config, path string) error {

	query := "s=id&o=asc&p="

	args := &client.Args{
		Proxy:  true,
		Method: http.MethodGet,
		Endpoint: &url.URL{
			Scheme:   "http",
			Host:     "nyaa.si",
			Path:     "/user/" + path,
			RawQuery: query + "0",
		},
	}
	var (
		err   error
		page  *page
		skip  bool
		state uint8
	)

action:
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			time.Sleep(3 * time.Second)

			page, err = s.browse(ctx, args)
			if err != nil {
				return err
			}

			s.log.Info("check the page ["+path+"] ...", "number", page.PageNumber)
			skip = false
			for _, v1 := range page.Torrents {
				if v1 == nil {
					continue
				}

				if _, err = collection.InsertOne(ctx, v1); err != nil {
					if mongo.IsDuplicateKeyError(err) {
						if cfg.ByID && cfg.BySeed {
							s.log.Info("pause on page ["+path+"] !!", "number", page.PageNumber)
							return nil
						}
						if cfg.ByID && !cfg.BySeed && state == 0 {
							page.NextPage = false
						}
						skip = true
						s.log.Info("skip page ["+path+"] !!", "number", page.PageNumber)
						break
					} else {
						return err
					}
				}

				if v1.NyaaID > cfg.NyaaID {
					cfg.NyaaID = v1.NyaaID
				}
			}
			if !skip {
				s.log.Info("scrape the page ["+path+"] ...", "number", page.PageNumber)
			}

			if !page.NextPage && state != 0 {
				if state == 1 {
					cfg.ByID = true
				} else {
					cfg.BySeed = true
				}
			}

			switch {
			case page.NextPage:
				args.Endpoint.RawQuery = query + strconv.FormatInt(int64(page.PageNumber+1), 10)
			case !cfg.ByID && page.Total > 7500:
				state = 1
				query = "s=id&o=desc&p="
				args.Endpoint.RawQuery = query + "0"
				goto action
			case !cfg.BySeed && page.Total > 15000:
				state = 2
				query = "s=seeders&o=desc&p="
				args.Endpoint.RawQuery = query + "0"
				goto action
			default:
				if page.Total < 7501 {
					cfg.ByID = true
				}
				if cfg.ByID && page.Total < 15001 {
					cfg.BySeed = true
				}
				return nil
			}
		}
	}
}
