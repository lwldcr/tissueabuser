// config module for crawler
// read configurations from local ini file

package config

import (
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/widuu/goini"
)

var logger *log.Logger

const (
	FileNotFoundErr    = -10
	OptionsNotValidErr = -20
	MakeDirectoryErr   = -30
	FilePath           = "conf.ini"
)

type conf struct {
	UseProxy  bool   // use http proxy or not
	HttpProxy string // http proxy
	DestDir   string // destination of local path for storing files

	MziTuStartUrl string // start url of mzitu.com
	MziTuMode     string // crawl mode: daily or full
	MzituTop      int    // top n albums
}

var Conf *conf

// read configuration

func Init(dir string, l *log.Logger) {
    logger = l

	var myconf conf

	configPath := filepath.Join(dir, FilePath)
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		logger.Println("cannot access configuration file:", err)
		os.Exit(FileNotFoundErr)
	}

	config := goini.SetConfig(configPath)
	config.ReadList()

	// read dest config and check validation
	dest := config.GetValue("local", "dest")
	if dest == "" {
		dest = "."
	}

	stat, err := os.Stat(dest)
	if os.IsNotExist(err) {
		logger.Println("dest dir not exist, will try to make dir")
		if err1 := os.Mkdir(dest, os.ModePerm); err1 != nil {
			logger.Println("making dir failed:", err)
			os.Exit(MakeDirectoryErr)
		}
		logger.Println("new directory created:", dest)
	} else if !stat.IsDir() {
		logger.Println("dest dir is not a valid directory")
		os.Exit(OptionsNotValidErr)
	}
	myconf.DestDir = dest

	// read http settings
	proxy := config.GetValue("http", "proxy")
	if proxy == "" {
		myconf.UseProxy = false
	} else {
		myconf.UseProxy = true
		myconf.HttpProxy = proxy
	}

	// read target sites
	target := config.GetValue("mzitu", "start")
	if target == "" {
		logger.Println("No target url given, exiting!")
		os.Exit(OptionsNotValidErr)
	}
	myconf.MziTuStartUrl = target

	mode := config.GetValue("mzitu", "mode")
	if mode == "" {
		mode = "daily" // default "daily" mode
	}
	myconf.MziTuMode = mode

	topn := config.GetValue("mzitu", "top")
	topnInt, err := strconv.Atoi(topn)
	if err != nil {
		topnInt = 100
	}
	myconf.MzituTop = topnInt

	Conf = &myconf
}
