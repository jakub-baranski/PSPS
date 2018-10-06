package main

import (
	"errors"
	"fmt"
	"github.com/gocolly/colly"
	"strconv"
	"strings"
)

type DiscountRecord struct {
	Title     string
	Link      string
	Discount  string
	Price     float64
	PlusPrice float64
	Currency  string
	Platform  string
}

type ScrapperOptions struct {
	baseUrl string
	region  Region
}

type Scrapper struct {
	collector *colly.Collector
	options   *ScrapperOptions
}

func (s *Scrapper) Init() {
	s.options = &ScrapperOptions{
		"psprices.com",
		US,
	}
	s.CreateCollector()
}

func (s *Scrapper) CreateCollector() {
	s.collector = colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36"), )
}

func (s *Scrapper) GetLink(platform Platform) string {
	b := strings.Builder{}
	b.WriteString("https://")
	b.WriteString(s.options.baseUrl)
	b.WriteRune('/')
	b.WriteString(string(s.options.region))
	b.WriteString("/discounts")
	b.WriteString("/?platform=")
	b.WriteString(string(platform))
	b.WriteString("&display=table")

	return b.String()
}

func (s *Scrapper) FetchNewDiscounts(platform Platform, ch chan DiscountRecord) {
	pageUrl := s.GetLink(platform) + "&sort=date"
	s.scrapPage(pageUrl, ch)
}

func (s *Scrapper) FetchBestDiscounts(platform Platform, ch chan DiscountRecord) {
	pageUrl := s.GetLink(platform) + "&sort=likes"
	s.scrapPage(pageUrl, ch)

}

func (s *Scrapper) scrapPage(pageUrl string, ch chan DiscountRecord) {
	s.collector.OnHTML("tr", func(e *colly.HTMLElement) {
		result, err := parseResultRow(e)
		if err == nil {
			ch <- result
		}
	})

	s.collector.OnScraped(func(r *colly.Response) {
		close(ch)
	})

	err := s.collector.Visit(pageUrl)
	if err != nil {
		fmt.Println(err)
	}
}

func parseResultRow(e *colly.HTMLElement) (record DiscountRecord, err error) {

	discount := e.ChildText(`td.text-right>strong`)
	title := e.ChildText(`a`)
	link := e.ChildAttr(`a`, `href`)
	platform := e.ChildText(`td>small`)

	price, _ := strconv.ParseFloat(strings.TrimPrefix(e.ChildText(`td.text-nowrap>strong:first-child`), "$"), 32)
	plusPrice, _ := strconv.ParseFloat(strings.TrimPrefix(e.ChildText(`td.text-nowrap>.content__game_card__price_plus_ico+strong`), "$"), 32)
	if discount != "" {
		record = DiscountRecord{title, link, discount, price, plusPrice, "$", platform}
		return record, nil
	} else {
		return DiscountRecord{}, errors.New("pupa")
	}

}

func NewScrapper(options ...func(*Scrapper)) *Scrapper {
	s := &Scrapper{}
	s.Init()

	for _, f := range options {
		f(s)
	}

	return s
}

// Methods for setting options
func SetBaseUrl(bu string) func(*Scrapper) {
	return func(s *Scrapper) {
		s.options.baseUrl = bu
	}
}

func SetRegion(r Region) func(*Scrapper) {
	return func(s *Scrapper) {
		s.options.region = r
	}
}
