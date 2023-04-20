package rss

import (
	"APIGateway/pkg/storage"
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"time"

	strip "github.com/grokify/html-strip-tags-go"
)

const DATE_TIME_TEMPL = "Mon, 2 Jan 2006 15:04:05 -0700"

type RssParserInput struct {
	UpdateInterval int // in minutes
	Url            string
}

type IrssFeed interface {
	News() []storage.News
}

type Rss struct {
	Channel ChannelStruct `xml:"channel"`
}

type ChannelStruct struct {
	Items []Item `xml:"item"`
}

type Item struct {
	Guid        string `xml:"guid"`
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

// Возвращает список подгтовленный новостей
func (rss Rss) News(guidPrefix string) []storage.News {
	return prepareNews(rss.Channel.Items, guidPrefix)
}

func prepareNews(items []Item, guidPrefix string) []storage.News {
	var newsList []storage.News
	for _, item := range items {
		news := storage.NewNews(
			guidPrefix+"_"+item.Guid,
			item.Title,
			strip.StripTags(item.Description),
			uint64(parseDateTime(item.PubDate).Unix()),
			item.Link,
		)
		newsList = append(newsList, news)
	}

	return newsList
}

func parseDateTime(dateTimeRaw string) time.Time {
	pubTime, err := time.Parse(time.RFC1123, dateTimeRaw)
	if err == nil {
		return pubTime
	}

	pubTime, err = time.Parse(time.RFC1123Z, dateTimeRaw)
	if err == nil {
		return pubTime
	}

	pubTime, err = time.Parse(DATE_TIME_TEMPL, dateTimeRaw)
	if err == nil {
		return pubTime
	}

	panic(err)
}

func getFromRss(url string, guidPrefix string) []storage.News {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	rss := new(Rss)
	err = xml.Unmarshal([]byte(body), rss)
	if err != nil {
		log.Fatal(err)
	}

	return rss.News(guidPrefix)
}

func extractGuidPrefix(url string) string {
	re := regexp.MustCompile(`(?:https:\/\/){0,1}(?:www\.){0,1}([-a-zA-Z0-9]{1,256}\.[a-zA-Z0-9]{1,6}){1}`)
	submatches := re.FindAllStringSubmatch(url, 1)

	return submatches[0][1]
}

// Запускает обновление новостей по URL RSS
func UpdateNews(url string, refrPeriod int, ch chan storage.News) {
	refrPerDur := time.Duration(refrPeriod) * time.Minute
	guidPrefix := extractGuidPrefix(url)

	for {
		newsList := getFromRss(url, guidPrefix)
		for _, news := range newsList {
			ch <- news
		}

		time.Sleep(refrPerDur)
	}
}
