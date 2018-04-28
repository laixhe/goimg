package imghand

import (
	"fmt"
	"github.com/laixhe/goimg/config"
	"strconv"
	"net/url"
)

// 匹配是否是 md5 的长度
func IsMD5Path(str string) bool {

	return regexpUrlParse.MatchString(str)

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

// 组合文件目录路径
func JoinPath(md5_str string) string {

	sortPath := SortPath([]byte(md5_str[:5]))
	return config.PathImg() + sortPath + "/" + md5_str

}

// 进行 url 部分解析 - md5，并组合文件完整路径
func UrlParse(md5_url string) string {

	if md5_url == "" {
		return ""
	}

	if len(md5_url) < 32 {
		return ""
	}

	// 进行 url 解析
	parse, err := url.Parse(md5_url)
	if err != nil {
		return ""
	}

	parsePath := parse.Path

	if len(parse.Path) != 32 {
		return ""
	}

	// 匹配是否是 md5 的长度
	if !IsMD5Path(parsePath) {
		return ""
	}

	// 组合文件完整路径
	return JoinPath(parsePath) + "/" + parsePath

}

func StringToInt(str string) int {
	if str == "" {
		return 0
	}

	toint, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}

	if toint < 0 {
		return 0
	}

	return toint
}