package crawler

import (
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

var logger *log.Logger

func SetLogger(l *log.Logger) {
	logger = l
}

type Crawler interface {
	Crawl()
}

type Album struct { // album struct definition
	title string
	url   string
}

type MzituCrawler struct {
	StartUrl   string          // start url
	Action     string          // actions: full/daily/category
	DestDir    string          // where to store files
	Mode       string          // mode: "daily" or "full"
	Topn       int             // topn
	Client     *Client         // client
	Wg         *sync.WaitGroup // sync
	Albums     chan *Album     // Albums
	MainImgPat *regexp.Regexp  // main image reg-exp pattern
	CurAlbum   *Album          // current album
}

func NewMzituCrawler(url string, mode string, top int, dest string,
	client *Client, wg *sync.WaitGroup) *MzituCrawler {

	albums := make(chan *Album)

	mainImgPattern, err := regexp.Compile("<div class=\"main-image\"><p><a href=\"(.*?)\" ><img src=\"(.*?)\"")

	if err != nil {
		logger.Println("cannot compile regexp for main image")
		os.Exit(REGEXPERR)
	}

	mzitu := MzituCrawler{
		StartUrl:   url,
		Mode:       mode,
		DestDir:    dest,
		Client:     client,
		Wg:         wg,
		Albums:     albums,
		MainImgPat: mainImgPattern,
		Topn:       top,
	}
	return &mzitu
}

func (m *MzituCrawler) Crawl() {
	go m.GetLinks()
	m.CrawlAlbums()

	close(m.Albums)
}

func (m *MzituCrawler) GetLinks() {
	defer m.Wg.Done()
	logger.Println("start crawling mzitu:", m.StartUrl)

	r, _ := http.NewRequest("GET", m.StartUrl, nil)
	resp, err := m.Client.DoRequest(r, "")
	if err != nil {
		logger.Println("request failed:", err)
		os.Exit(REQUESTERR)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Println("read response body failed:", err)
		os.Exit(READDATAERR)
	}

	content := string(body)

	//_, _, day := time.Now().Date()
	//dateStr := strconv.Itoa(day) + "日: "
	albumPattern, _ := regexp.Compile("<a href=\"(http://.*?)\" target=\"_blank\">(.*?)</a>")

	var topn int = -1
	if m.Mode != "full" {
		topn = m.Topn
	}
	matcher := albumPattern.FindAllStringSubmatch(content, topn)
	//matcher := albumPattern.FindAllStringSubmatch(content, -1)

	albums := make(map[string]string)
	for _, a := range matcher {
		albums[a[2]] = a[1] // album: { "title": link }
	}

	for title, link := range albums {
		album := Album{title, link}
		m.Albums <- &album
	}

	m.Albums <- &Album{url: ""}
}

func (m *MzituCrawler) CrawlAlbums() {
	logger.Println("listening for albums...")
	for {
		select {
		case album := <-m.Albums:
			if album.url == "" {
				logger.Println("got blank album, stop crawling")
				return
			}

			logger.Println("got album:", album.title, album.url)
			albumDir := strings.Join([]string{m.DestDir, album.title}, string(os.PathSeparator))

			if stat, err := os.Stat(albumDir); os.IsNotExist(err) {
				logger.Println("trying to build album dir")
				if err := os.Mkdir(albumDir, os.ModePerm); err != nil {
					logger.Println("failed:", err)
					continue
				}
				logger.Println("dir created successfully")
			} else if !stat.IsDir() {
				logger.Println("path already exists but not a valid directory:", albumDir)
				continue
			}

			m.CurAlbum = album
			m.CrawlPage(album.url)

			time.Sleep(time.Duration(rand.Intn(5)) * time.Second)
		}
	}
}

// crawl single page
func (m *MzituCrawler) CrawlPage(url string) {

	r, _ := http.NewRequest("GET", url, nil)
	resp, err := m.Client.DoRequest(r, "")
	if err != nil {
		logger.Println("fetching album failed:", err)
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Println("read data failed:", err)
		return
	}

	content := string(body)
	matcher := m.MainImgPat.FindAllStringSubmatch(content, -1)

	if len(matcher) > 0 {
		matched := matcher[0]
		nextPage := matched[1]
		imgUrl := matched[2]
		time.Sleep(time.Duration(rand.Intn(5)) * time.Second)
		m.CrawlNext(url, imgUrl, nextPage)
	}
}

func (m *MzituCrawler) CrawlNext(referer string, imgUrl string, nextPage string) {
	logger.Println("got image link:", imgUrl)
	r, _ := http.NewRequest("GET", imgUrl, nil)
	resp, err := m.Client.DoRequest(r, referer)
	if err != nil {
		logger.Println("cannot get img:", err)
		return
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Println("read data failed:", err)
		return
	}

	urlParts := strings.Split(imgUrl, "/")
	fileName := urlParts[len(urlParts)-1]
	fullName := strings.Join([]string{m.DestDir, m.CurAlbum.title, fileName}, string(os.PathSeparator))

	if _, err := os.Stat(fullName); os.IsExist(err) {
		logger.Println("warning: already exists:", fullName)
		goto RECURSIVE
	}
	logger.Println("storing file:", fullName)
	err = ioutil.WriteFile(fullName, data, os.ModePerm)
	if err != nil {
		logger.Println("store file failed:", err)
	}

RECURSIVE:
	if !isNewAlbum(nextPage) {
		m.CrawlPage(nextPage)
	}
}

// judge if next page link is a new album
func isNewAlbum(url string) bool {
	match, err := regexp.Match("http://www.mzitu.com/\\d+/\\d+", []byte(url))
	if err != nil {
		logger.Println("matching failed:", err)
		return true
	}
	return !match
}
