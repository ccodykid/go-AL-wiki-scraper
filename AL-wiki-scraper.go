package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/gocolly/colly/v2"
)

func dl_ship_img(collector *colly.Collector, root_url string, dir string) {
	base_url_collector := collector.Clone()

	base_url_collector.OnHTML("table tbody tr", func(h *colly.HTMLElement) {

		href_attr := h.ChildAttr("td:nth-child(2) > a", "href")
		url := h.Request.AbsoluteURL(href_attr)

		img_url_collector := collector.Clone()

		img_url_collector.OnHTML("div[class*=shipgirl-image]", func(h *colly.HTMLElement) {

			img_src := h.ChildAttr("img:first-child", "src")
			log.Println("[Coll-3] Found: ", img_src)

			img_src_collector := collector.Clone()

			img_src_collector.OnRequest(func(r *colly.Request) {
				log.Printf("[Coll-3] Visiting: %s", r.URL)
			})

			img_src_collector.OnResponse(func(r *colly.Response) {
				r.Save(filepath.Join(dir, r.FileName()))
				log.Printf("[Coll-3] Successfully saved %s", r.FileName())
			})

			img_src_collector.Visit(img_src)
		})

		img_url_collector.OnRequest(func(r *colly.Request) {
			log.Printf("[Coll-2] Visiting: %s", r.URL)
		})

		img_url_collector.Visit(url)
	})

	base_url_collector.OnHTML("title", func(h *colly.HTMLElement) {
		fmt.Println(h.Text)
	})

	base_url_collector.OnRequest(func(r *colly.Request) {
		log.Printf("[Coll-1] Visiting: %s", r.URL)
	})

	base_url_collector.Visit(root_url)

	base_url_collector.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})
}

func create_dir(dir string) string {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}

	return dir
}

func main() {

	var root_url = "https://azurlane.koumakan.jp/wiki/List_of_Ships"

	var parent_collector = colly.NewCollector(
		colly.AllowedDomains("azurlane.koumakan.jp", "azurlane.netojuu.com"),
	)

	var output_dir = create_dir("images")

	dl_ship_img(parent_collector, root_url, output_dir)
}
