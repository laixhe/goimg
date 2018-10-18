package config

import (
	"github.com/go-ini/ini"
	"log"
	"strings"
)

var conf *ini.File

func init() {

	var err error
	conf, err = ini.Load("config.ini")
	if err != nil {
		log.Fatalln(err)
	}

}

func Get(str string) string {

	strArr := strings.Split(str, ".")

	if len(strArr) == 2 {
		return conf.Section(strArr[0]).Key(strArr[1]).String()
	}

	return conf.Section("").Key(strArr[0]).String()
}

func HttpAddr() string {

	addr := Get("http.addr")
	if addr == "" {
		return ":8080"
	}

	return addr
}

func PathImg() string {

	path := Get("img.path")
	if path == "" {
		return "img/"
	}

	return path
}
