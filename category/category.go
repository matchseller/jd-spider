package category

import (
	"errors"
	jsoniter "github.com/json-iterator/go"
	"github.com/matchseller/jd-spider/util"
	"regexp"
	"strings"
)

const url string = "https://dc.3.cn/category/get"

//start crawling
func Crawl() (catList []string, err error) {
	var ret []byte
	ret, err = util.GetUrl(url)
	if err != nil {
		return
	}
	content := util.GbkToUtf8(string(ret))
	reader := strings.NewReader(content)
	var params map[string]interface{}
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	decoder := json.NewDecoder(reader)
	err = decoder.Decode(&params)
	if err != nil {
		return
	}
	data, isOk := params["data"].([]interface{})
	if !isOk {
		return nil, errors.New("json parse error: params[\"data\"]")
	}
	catList = composeData(data)
	catList = removeDuplicates(catList)
	return
}

func removeDuplicates(arr []string) (ret []string) {
	for k, v := range arr {
		if index := strings.Index(v, "&"); index != -1 {
			arr[k] = v[:index]
		}
	}
	util.RemoveEmptyString(&arr)
	return util.RemoveDuplicatesAndEmpty(arr)
}

//compose data and return category list
func composeData(data []interface{}) (catList []string) {
	for _, v := range data {
		innerDataMap, isOk := v.(map[string]interface{})
		if !isOk {
			continue
		}
		_, isOk = innerDataMap["s"]
		if !isOk {
			continue
		}
		innerData, isOk := innerDataMap["s"].([]interface{})
		if !isOk {
			continue
		}
		_, isOk = innerDataMap["n"]
		if isOk {
			urlStr, isOk := innerDataMap["n"].(string)
			if !isOk {
				continue
			}
			urlArr := strings.Split(urlStr, "|")
			if len(urlArr) == 0 {
				continue
			}
			var categoryUrl string
			if strings.Contains(urlArr[0], "i-list.jd.com/list.html") || strings.Contains(urlArr[0], "list.jd.com/list.html") || strings.Contains(urlArr[0], "coll.jd.com/list.html") {
				categoryUrl = urlArr[0]
			} else if isMatch, _ := regexp.MatchString(`^[0-9]+-[0-9]+-[0-9]`, urlArr[0]); isMatch {
				categoryUrl = "list.jd.com/list.html?cat=" + strings.ReplaceAll(urlArr[0], "-", ",")
			}
			if categoryUrl != "" {
				catList = append(catList, "https://"+categoryUrl)
			}
		}
		catList = append(catList, composeData(innerData)...)
	}
	return catList
}
