package main

import (
	"crypto/md5"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
)

func main()  {
	fmt.Println("main")
	res, err := http.Get("http://www.bikebound.com/tag/trackers/")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find("article.post").Each(func(i int, s *goquery.Selection) {
		text := s.Find("h3.entry-title").First().Text()
		link := s.Find("h3 a").First().AttrOr("href", "")
		fmt.Println(text, link)
		crawBike(text, link)
	})

}
func crawBike(name string, link string) {
	hash := fmt.Sprintf("%x", md5.Sum([]byte(name)))
	baseDir := "motocycles/" + hash
	err := os.Mkdir(baseDir, 0755)
	if err != nil {
		panic(err)
	}
	res, err := http.Get(link)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find("div.entry-content a").Each(func(i int, s *goquery.Selection) {
		image := s.AttrOr("href", "")
		if image != "" && strings.HasPrefix(image, "http://www.bikebound.com/wp-content/uploads"){
			filename := baseDir + "/" + strconv.Itoa(i) + path.Ext(image)
			download(image, filename)
		}
	})
}

func download(link string, filename string)  {
	log.Println("Downloading", link)
	res, err := http.Get(link)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	ioutil.WriteFile(filename, data, 0755)

}
