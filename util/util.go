package util

import (
	"errors"
	"github.com/axgle/mahonia"
	"io/ioutil"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

func GetUrl(url string) ([]byte, error) {
	client := &http.Client{}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.81 Safari/537.36")
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New("request err,http code:" + strconv.Itoa(resp.StatusCode))
	}
	ret, _ := ioutil.ReadAll(resp.Body)
	return ret, nil
}

func convertToString(src string, srcCode string, tagCode string) string {
	srcCoder := mahonia.NewDecoder(srcCode)
	srcResult := srcCoder.ConvertString(src)
	tagCoder := mahonia.NewDecoder(tagCode)
	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)
	result := string(cdata)
	return result
}

func GbkToUtf8(str string) string {
	return convertToString(str, "gbk", "utf-8")
}

//去除html标签
func TrimHtml(src string) string {
	//将HTML标签全转换成小写
	re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllStringFunc(src, strings.ToLower)
	//去除STYLE
	re, _ = regexp.Compile("\\<style[\\S\\s]+?\\</style\\>")
	src = re.ReplaceAllString(src, "")
	//去除SCRIPT
	re, _ = regexp.Compile("\\<script[\\S\\s]+?\\</script\\>")
	src = re.ReplaceAllString(src, "")
	//去除所有尖括号内的HTML代码，并换成空格
	re, _ = regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllString(src, " ")
	//去除连续的换行符
	re, _ = regexp.Compile("\\s{2,}")
	src = re.ReplaceAllString(src, " ")
	return strings.TrimSpace(src)
}

func RemoveDuplicatesAndEmpty(arr []string) (ret []string) {
	sort.Strings(arr)
	arrLen := len(arr)
	for i := 0; i < arrLen; i++ {
		if (i > 0 && arr[i-1] == arr[i]) || len(arr[i]) == 0 {
			continue
		}
		ret = append(ret, arr[i])
	}
	return
}

func RemoveEmptyString(arr *[]string) {
	arrLen := len(*arr)
	for i := 0; i < arrLen; i++ {
		(*arr)[i] = strings.ReplaceAll((*arr)[i], " ", "")
	}
}
