package utils

import (
	"regexp"
)

func TidyString(s string) string {
	regex := regexp.MustCompile(`[\n\t]`)
	cleanedString := regex.ReplaceAllString(s, "") //ReplaceAllString 方法会将所有匹配正则表达式的子字符串替换为空字符串（即删除这些子字符串）。
	return cleanedString
} //删除输入字符串中的所有换行符和制表符，并返回清理后的字符串。
