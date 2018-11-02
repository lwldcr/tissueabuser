package main

import (
	"flag"
	"log"
	"os"
	"sync"

	"github.com/lwldcr/tissueabuser/config"
	"github.com/lwldcr/tissueabuser/crawler"
)

const (
	ChdirFailedErr = -30
	Prefix = "[cmd]"
)

var (
	confDir string
	logger *log.Logger
)

func init() {
	logger = log.New(os.Stdout, Prefix, log.Ldate|log.Ltime|log.Llongfile)
	logger.Println("init app")
	flag.StringVar(&confDir, "config_dir", "./config", "configuration file directory")

	flag.Parse()
	logger.Printf("init done,with given config_dir: %s", confDir)
}

func main() {
	config.Init(confDir, logger)
	conf := config.Conf

	crawler.SetLogger(logger)
	client := crawler.NewClient(conf.HttpProxy)
	//client.UseProxy()

	var wg sync.WaitGroup
	wg.Add(1)

	mzitu := crawler.NewMzituCrawler(conf.MziTuStartUrl,
		conf.MziTuMode, conf.MzituTop, conf.DestDir, client, &wg)

	mzitu.Crawl()
	wg.Wait()
}
