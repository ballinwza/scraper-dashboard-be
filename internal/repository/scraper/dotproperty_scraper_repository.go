package scraper

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"go.uber.org/zap"

	"github.com/ballinwza/scraper-dashboard-be/internal/domain"
	"github.com/ballinwza/scraper-dashboard-be/pkg/logger"
)

// ScraperRepositoryInterface
type IDotpropertyScraperRepository interface {
	ScrapeMainPage(targetURL string, currentPage, maxPages int) ([]domain.RentalEstate, error)
}

type DotpropertyScraperRepository struct {
	log *zap.Logger
}

func NewDotpropertyScraperRepository() IDotpropertyScraperRepository {
	return &DotpropertyScraperRepository{
		log: logger.Log,
	}
}

func (s *DotpropertyScraperRepository) ScrapeMainPage(targetURL string, startPage, maxPages int) ([]domain.RentalEstate, error) {
	var items []domain.RentalEstate
	var mu sync.Mutex

	c := colly.NewCollector(
		colly.AllowedDomains("www.dotproperty.co.th", "dotproperty.co.th"),
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"),
	)

	c.SetRequestTimeout(15 * time.Second)
	c.WithTransport(&http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	})

	c.Limit(&colly.LimitRule{
		DomainRegexp: `dotproperty\.co\.th`,
		Parallelism:  2,
		RandomDelay:  1 * time.Second,
		Delay:        2 * time.Second,
	})

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Accept-Language", "en-US,en;q=0.9")
	})

	c.OnError(func(r *colly.Response, err error) {
		logger.Log.Error("Failed to scrapping", zap.Error(err))
	})

	c.OnResponse(func(r *colly.Response) {
		if r.StatusCode != 200 {
			logger.Log.Error("Blocked or Error Page!", zap.Int("StatusCode", 500))
		}
	})

	c.OnHTML("script[type='application/ld+json']", func(e *colly.HTMLElement) {
		if !strings.Contains(e.Text, `"ItemList"`) {
			return
		}

		var itemList dotpropertyItemListEntity
		err := json.Unmarshal([]byte(e.Text), &itemList)
		if err != nil {
			logger.Log.Error("Error unmarshaling JSON-LD", zap.Error(err))
			return
		}

		doc := e.DOM.Parents().Last()
		baths := 0
		area := 0.0

		for _, el := range itemList.ItemListElement {
			listing := el.Item

			// Price to float64
			var priceVal float64
			switch v := listing.Offers.Price.(type) {
			case string:
				priceVal, _ = strconv.ParseFloat(v, 64)
			case float64:
				priceVal = v
			}

			// Location
			locationStr := fmt.Sprintf("%s, %s, %s",
				listing.About.Address.StreetAddress,
				listing.About.Address.AddressLocality,
				listing.About.Address.AddressRegion,
			)

			// Scrap from HTML
			listingID := extractListingID(listing.URL)
			card := doc.Find(fmt.Sprintf("a[href*='%s']", listingID))
			if card.Length() > 0 {
				card.Find("div[data-testid='unit-search-result-item-details']").Each(func(_ int, s *goquery.Selection) {
					text := strings.TrimSpace(s.Text())
					baths = extractBathrooms(text)
					area = extractAreaSqM(text)
				})
			}

			if listing.About.NumberOfBedrooms == 0 {
				listing.About.NumberOfBedrooms = 1
			}

			item := domain.RentalEstate{
				Title:        listing.Name,
				FormalName:   listing.About.ContainedInPlace.Name,
				Description:  listing.Description,
				DatePosted:   listing.DatePosted,
				PropertyType: listing.About.Category,
				Price:        priceVal,
				Location:     locationStr,
				Bedrooms:     listing.About.NumberOfBedrooms,
				Bathrooms:    baths,
				AreaSqM:      area,
				ImageURL:     listing.Image,
				SourceURL:    listing.URL,
				Latitude:     listing.About.Geo.Latitude,
				Longitude:    listing.About.Geo.Longitude,
			}

			mu.Lock()
			items = append(items, item)
			mu.Unlock()
		}
	})

	c.OnHTML("li.pagination-next a, a.next, a[rel='next']", func(e *colly.HTMLElement) {
		if startPage < maxPages {
			nextPage := e.Request.AbsoluteURL(e.Attr("href"))

			if nextPage != "" && nextPage != e.Request.URL.String() {
				startPage++
				logger.Log.Info("Visiting page",
					zap.Int("CurrentPage", startPage),
					zap.String("NextPage", nextPage),
				)
				e.Request.Visit(nextPage)
			}
		}
	})

	// Start visiting
	err := c.Visit(targetURL)
	if err != nil {
		log.Fatal(err)
	}

	c.Wait()

	s.log.Info("Scraper job completed", zap.Int("total scraped", len(items)))
	return items, nil
}

// Helper Functions
func extractBathrooms(desc string) int {
	if desc == "" {
		return 0
	}

	patterns := []*regexp.Regexp{
		regexp.MustCompile(`(\d+)\s*ห้องน้ำ`),
		regexp.MustCompile(`ห้องน้ำ\s*(\d+)`),
		regexp.MustCompile(`(?i)(\d+)\s*(?:Bathroom|Bath|baths|bathrooms)\b`),
	}

	for _, re := range patterns {
		matches := re.FindStringSubmatch(desc)
		if len(matches) > 1 {
			if val, err := strconv.Atoi(matches[1]); err == nil {
				return val
			}
		}
	}

	return 0
}

func extractAreaSqM(desc string) float64 {
	if desc == "" {
		return 0
	}

	patterns := []*regexp.Regexp{
		regexp.MustCompile(`([\d\.]+)\s*(?:ตร\.?ม\.?|ตารางเมตร)`),
		regexp.MustCompile(`(?i)([\d\.]+)\s*(?:Sq\.?M|Sqm|Sq\s*m|m²)\b`),
	}

	for _, re := range patterns {
		matches := re.FindStringSubmatch(desc)
		if len(matches) > 1 {
			if val, err := strconv.ParseFloat(matches[1], 64); err == nil {
				return val
			}
		}
	}

	return 0
}

func extractListingID(fullURL string) string {
	u, err := url.Parse(fullURL)
	if err != nil {
		return ""
	}

	parts := strings.Split(u.Path, "_")
	if len(parts) > 1 {
		return parts[len(parts)-1]
	}
	return u.Path
}
