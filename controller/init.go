package controller

import (
	"regexp"
	"log"
)

var regexpUrlParse *regexp.Regexp

func init(){

	var err error
	// 初始化正则表达式
	regexpUrlParse, err = regexp.Compile("[a-z0-9]{32}")
	if err != nil {
		log.Println("regexpUrlParse:", err)
	}

}
