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
	logger.Println("start crawling")
	go m.GetLinks()
	m.CrawlAlbums()
}

func (m *MzituCrawler) GetLinks() {
	defer m.Wg.Done()
	logger.Println("start crawling mzitu album links:", m.StartUrl)

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

	topN := -1
	if m.Mode != "full" {
		topN = m.Topn
	}
	logger.Printf("find matching albums of top: %d", topN)
	matcher := albumPattern.FindAllStringSubmatch(content, topN)

	logger.Printf("got matched albums: %d", len(matcher))
	albums := make(map[string]string)
	for _, a := range matcher {
		albums[a[2]] = a[1] // album: { "title": link }
	}

	for title, link := range albums {
		album := Album{title, link}
		m.Albums <- &album
	}

	// close albums channel
	logger.Printf("all %d albums sent, will close albums channel", len(albums))
	close(m.Albums)
}

func (m *MzituCrawler) CrawlAlbums() {
	logger.Println("listening for albums...")
	var wg sync.WaitGroup
	FOR:for {
		select {
		case album, ok := <-m.Albums:
			if !ok {
				logger.Println("albums channel closed, will exit crawler after unfinished jobs done")
				break FOR
			}
			wg.Add(1)
			go m.DoCrawl(album, &wg)
		}
	}

	wg.Wait()
	logger.Println("all crawling jobs done, returning")
}

func (m *MzituCrawler) DoCrawl(album *Album, wg *sync.WaitGroup) {
	defer wg.Done()
	logger.Println("got album:", album.title, album.url)
	albumDir := strings.Join([]string{m.DestDir, album.title}, string(os.PathSeparator))

	if stat, err := os.Stat(albumDir); os.IsNotExist(err) {
		logger.Println("trying to build album dir")
		if err := os.Mkdir(albumDir, os.ModePerm); err != nil {
			logger.Println("failed:", err)
			return
		}
		logger.Println("dir created successfully")
	} else if !stat.IsDir() {
		logger.Println("path already exists but not a valid directory:", albumDir)
		return
	}

	//m.CurAlbum = album
	m.CrawlPage(album.url, album)
	time.Sleep(time.Duration(rand.Intn(3)) * time.Millisecond)
}

// crawl single page
func (m *MzituCrawler) CrawlPage(url string, album *Album) {

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
		m.CrawlNext(url, imgUrl, nextPage, album)
	}
}

func (m *MzituCrawler) CrawlNext(referer string, imgUrl string, nextPage string, album *Album) {
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
	fullName := strings.Join([]string{m.DestDir, album.title, fileName}, string(os.PathSeparator))
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
		m.CrawlPage(nextPage, album)
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
