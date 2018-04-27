package imghand

import (
	"fmt"
	"github.com/laixhe/goimg/config"
	"strconv"
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