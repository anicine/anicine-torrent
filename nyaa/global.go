package nyaa

import (
	"context"
	"errors"
	"net/http"
	"net/url"
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
	trackers          = []string{
		"http://125.227.35.196:6969/announce",
		"http://210.244.71.25:6969/announce",
		"http://210.244.71.26:6969/announce",
		"http://213.159.215.198:6970/announce",
		"http://37.19.5.139:6969/announce",
		"http://37.19.5.155:6881/announce",
		"http://46.4.109.148:6969/announce",
		"http://87.248.186.252:8080/announce",
		"http://anidex.moe:6969/announce",
		"http://asmlocator.ru:34000/1hfZS1k4jh/announce",
		"http://bt.evrl.to/announce",
		"http://bt.rutracker.org/ann",
		"http://mgtracker.org:6969/announce",
		"http://nyaa.tracker.wf:7777/announce",
		"http://pubt.net:2710/announce",
		"http://tracker.acgnx.se/announce",
		"http://tracker.baravik.org:6970/announce",
		"http://tracker.dler.org:6969/announce",
		"http://tracker.filetracker.pl:8089/announce",
		"http://tracker.grepler.com:6969/announce",
		"http://tracker.mg64.net:6881/announce",
		"http://tracker.openbittorrent.com:80/announce",
		"http://tracker.tiny-vps.com:6969/announce",
		"http://tracker.torrentyorg.pl/announce",
		"https://computer1.sitelio.me/",
		"https://internet.sitelio.me/",
		"https://opentracker.i2p.rocks:443/announce",
		"https://www.artikelplanet.nl",
		"udp://168.235.67.63:6969",
		"udp://182.176.139.129:6969",
		"udp://37.19.5.155:2710",
		"udp://46.148.18.250:2710",
		"udp://46.4.109.148:6969",
		"udp://allerhandelenlaag.nl",
		"udp://c3t.org",
		"udp://computerbedrijven.bestelinks.nl/",
		"udp://computerbedrijven.startsuper.nl/",
		"udp://computershop.goedbegin.nl/",
		"udp://coppersurfer.tk:6969/announce",
		"udp://exodus.desync.com:6969",
		"udp://exodus.desync.com:6969/announce",
		"udp://ipv4.tracker.harry.lu:80/announce",
		"udp://open.demonii.com:1337/announce",
		"udp://open.stealth.si:80/announce",
		"udp://open.tracker.cl:1337/announce",
		"udp://opentracker.i2p.rocks:6969/announce",
		"udp://p4p.arenabg.com:1337/announce",
		"udp://public.popcorn-tracker.org:6969/announce",
		"udp://tracker.bittor.pw:1337/announce",
		"udp://tracker.dler.org:6969/announce",
		"udp://tracker.internetwarriors.net:1337/announce",
		"udp://tracker.leechers-paradise.org:6969/announce",
		"udp://tracker.openbittorrent.com:80",
		"udp://tracker.opentrackr.org:1337",
		"udp://tracker.opentrackr.org:1337/announce",
		"udp://tracker.publicbt.com:80",
		"udp://tracker.tiny-vps.com:6969",
		"udp://tracker.torrent.eu.org:451/announce",
		"udp://tracker.zer0day.to:1337/announc",
	}
)

func (s *Nyaa) extract(ctx context.Context, magnet string) ([]*torrent.File, error) {
	torrent, err := s.client.AddMagnet(magnet)
	if err != nil {
		return nil, err
	}

	torrent.AddTrackers([][]string{trackers})
	torrent.DisallowDataUpload()
	torrent.DisallowDataDownload()

	now := time.Now()
	for {
		select {
		case <-ctx.Done():
			return nil, context.Canceled
		case <-torrent.GotInfo():
			files := torrent.Files()
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
