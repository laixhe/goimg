package controller

import (
	"fmt"
	"log"
	"regexp"
)

var regexpUrlParse *regexp.Regexp

func init() {

	var err error
	// 初始化正则表达式
	regexpUrlParse, err = regexp.Compile("[a-z0-9]{32}")
	if err != nil {
		log.Println("regexpUrlParse:", err)
	}

}

// 路径部分排序做目录
func SortPath(str []byte) string {

	strLen := len(str)
	for i := 0; i < strLen; i++ {
		for j := 1 + i; j < strLen; j++ {
			if str[i] > str[j] {
				str[i], str[j] = str[j], str[i]
			}
		}
	}

	ret := ""

	for i := 0; i < strLen; i++ {
		ret += fmt.Sprintf("%d", str[i])
	}

	return ret
}
