// config module for crawler
// read configurations from local ini file

package config

import (
	"fmt"
	"github.com/widuu/goini"
	"os"
	"strconv"
)

const (
	FileNotFoundErr    = -10
	OptionsNotValidErr = -20
	MakeDirectoryErr   = -30
	FilePath           = "./conf.ini"
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

func Init() {
	var myconf conf

	if _, err := os.Stat(FilePath); os.IsNotExist(err) {
		fmt.Println("cannot access configuration file:", err)
		os.Exit(FileNotFoundErr)
	}

	config := goini.SetConfig("./conf.ini")
	config.ReadList()

	// read dest config and check validation
	dest := config.GetValue("local", "dest")
	if dest == "" {
		dest = "."
	}

	stat, err := os.Stat(dest)
	if os.IsNotExist(err) {
		fmt.Println("dest dir not exist, will try to make dir")
		if err1 := os.Mkdir(dest, os.ModePerm); err1 != nil {
			fmt.Println("making dir failed:", err)
			os.Exit(MakeDirectoryErr)
		}
		fmt.Println("new directory created:", dest)
	} else if !stat.IsDir() {
		fmt.Println("dest dir is not a valid directory")
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
		fmt.Println("No target url given, exiting!")
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
