package main

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"tissueabuser"
	"tissueabuser/config"
)

const (
	ChdirFailedErr = -30
)

func main() {
	cwd, _ := os.Getwd()

	confdir := strings.Join([]string{cwd, "..", "config"}, string(os.PathSeparator))
	fmt.Println("temporarily change working dir to:", confdir)
	if err := os.Chdir(confdir); err != nil {
		fmt.Println("chdir failed:", err)
		os.Exit(ChdirFailedErr)
	}

	config.Init()
	conf := config.Conf

	fmt.Println("change working dir back:", cwd)
	if err := os.Chdir(cwd); err != nil {
		fmt.Println("chdir failed:", err)
		os.Exit(ChdirFailedErr)
	}

	client := tissueabsuer.NewClient(conf.HttpProxy)
	//client.UseProxy()

	var wg sync.WaitGroup
	wg.Add(1)

	mzitu := tissueabsuer.NewMzituCrawler(conf.MziTuStartUrl,
		conf.MziTuMode, conf.MzituTop, conf.DestDir, client, &wg)

	mzitu.Crawl()
	wg.Wait()
}
