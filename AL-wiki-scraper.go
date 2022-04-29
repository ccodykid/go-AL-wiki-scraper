package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gocolly/colly/v2"
)

func get_img(collector *colly.Collector, img_src_dump []string) {
	img_collector := collector.Clone()

	img_collector.OnRequest(func(r *colly.Request) {
		log.Println("Coll-3 Visiting: ", r.URL)
	})

	img_collector.OnResponse(func(r *colly.Response) {
		log.Printf("%d", r.StatusCode)
		r.Save(r.FileName())
		log.Printf("Coll-3 Successfully saved %s", r.FileName())
	})

	for _, link := range img_src_dump {
		img_collector.Visit(link)
	}

}

func get_img_url(collector *colly.Collector, link_dump []string) []string {
	icon_collector := collector.Clone()
	var img_src_dump []string

	icon_collector.OnHTML("div[class*=shipgirl-image]", func(h *colly.HTMLElement) {
		img_src := h.ChildAttr("img:first-child", "src")
		log.Println("Found: ", img_src)

		img_src_dump = append(img_src_dump, img_src)
	})

	icon_collector.OnRequest(func(r *colly.Request) {
		log.Println("Coll-2 Visiting: ", r.URL)
	})

	for _, link := range link_dump {
		icon_collector.Visit(link)
	}

	return img_src_dump
}

func get_ship_url(collector *colly.Collector, root_url string) []string {
	link_collector := collector.Clone()
	var link_dump []string

	link_collector.OnHTML("table tbody tr", func(h *colly.HTMLElement) {

		href_attr := h.ChildAttr("td:nth-child(2) > a", "href")
		url := h.Request.AbsoluteURL(href_attr)

		link_dump = append(link_dump, url)
	})

	link_collector.OnHTML("title", func(h *colly.HTMLElement) {
		fmt.Println(h.Text)
	})

	link_collector.OnRequest(func(r *colly.Request) {
		fmt.Println("Coll-1 Visiting", r.URL)
	})

	link_collector.Visit(root_url)

	link_collector.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong:", err)
	})

	return link_dump
}

func create_filename(link string, counter int) string {
	ship_name := strings.Split(link, "/")[2]
	filename := fmt.Sprintf("%d_%s.png", counter, ship_name)
	return filename
}

func create_dir(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func main() {

	var root_url = "https://azurlane.koumakan.jp/wiki/List_of_Ships"

	var parent_collector = colly.NewCollector(
		colly.AllowedDomains("azurlane.koumakan.jp", "azurlane.netojuu.com"),
	)

	ship_url_list := get_ship_url(parent_collector, root_url)
	img_url_list := get_img_url(parent_collector, ship_url_list)
	get_img(parent_collector, img_url_list)
}
