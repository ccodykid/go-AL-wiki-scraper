package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/gocolly/colly/v2"
)

var c = colly.NewCollector(
	colly.AllowedDomains("azurlane.koumakan.jp", "azurlane.netojuu.com"),
)

// TODO: func get_link()
// TODO: implement get class
// TODO: implement download image

func get_img(img_src string) {
	img_coll := c.Clone()

	img_coll.OnRequest(func(r *colly.Request) {
		log.Println("Coll-3 Visiting: ", r.URL)
	})

	img_coll.OnResponse(func(r *colly.Response) {
		log.Printf("%d", r.StatusCode)
		r.Save(r.FileName())
		log.Printf("Coll-3 Successfully saved %s", r.FileName())
	})

	img_coll.Visit(img_src)
}

func get_img_url(url string, link string) {
	icon_coll := c.Clone()

	icon_coll.OnHTML("div[class*=shipgirl-image]", func(h *colly.HTMLElement) {
		img_src := h.ChildAttr("img:first-child", "src")
		log.Printf(img_src)

		get_img(img_src)
	})

	icon_coll.OnRequest(func(r *colly.Request) {
		log.Println("Coll-2 Visiting: ", r.URL)
	})

	icon_coll.Visit(url)
}

func create_filename(link string, counter int) string {
	ship_name := strings.Split(link, "/")[2]
	filename := fmt.Sprintf("%d_%s.png", counter, ship_name)
	return filename
}

func main() {

	counter := 0

	link_coll := c.Clone()

	link_coll.OnHTML("table tbody tr", func(h *colly.HTMLElement) {

		link := h.ChildAttr("td:nth-child(2) > a", "href")
		url := h.Request.AbsoluteURL(link)

		get_img_url(url, link)
	})

	link_coll.OnHTML("title", func(h *colly.HTMLElement) {
		fmt.Println(h.Text)
	})

	link_coll.OnRequest(func(r *colly.Request) {
		fmt.Println("Coll-1 Visiting", r.URL)
	})

	link_coll.Visit("https://azurlane.koumakan.jp/wiki/List_of_Ships")

	link_coll.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong:", err)
	})

	fmt.Println(counter)
}
