package price

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/matchseller/jd-spider/util"
	"math"
	"strconv"
	"strings"
	"sync"
)

type catchInfo struct {
	Cfb string
	Id  string
	M   string
	Op  string
	P   string
}

const url string = "https://p.3.cn/prices/mgets?skuIds="

var maxCrawlAmount = 100
var retryMax = 30

var wg sync.WaitGroup
var mutex sync.RWMutex

var jPrice map[string]float64

//start crawling
func Crawl(skuIds *[]string) map[string]float64 {
	jPrice = make(map[string]float64)
	util.RemoveEmptyString(skuIds)
	*skuIds = util.RemoveDuplicatesAndEmpty(*skuIds)
	skuIdLen := len(*skuIds)
	maxRoutines := int(math.Ceil(float64(skuIdLen) / float64(maxCrawlAmount)))
	for i := 0; i < maxRoutines; i++ {
		sliceStart := i * maxCrawlAmount
		var sliceEnd int
		if i == (maxRoutines - 1) {
			sliceEnd = skuIdLen
		} else {
			sliceEnd = sliceStart + maxCrawlAmount
		}
		wg.Add(1)
		go run((*skuIds)[sliceStart:sliceEnd])
	}
	wg.Wait()
	return jPrice
}

func SetRetryMax(max int) {
	retryMax = max
}

func run(skuIds []string) {
	var catchList []catchInfo
	urlStr := url + strings.Join(skuIds, ",")
	retryStart := 1
	var content string
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	for {
		ret, _ := util.GetUrl(urlStr)
		content = string(ret)
		containsError := strings.Contains(content, "error")
		if (content == "" || containsError) && retryStart <= retryMax {
			retryStart++
			continue
		}
		reader := strings.NewReader(content)
		decoder := json.NewDecoder(reader)
		decoder.Decode(&catchList)
		break
	}
	if len(skuIds) == len(catchList) {
		for k, skuId := range skuIds {
			price, _ := strconv.ParseFloat(catchList[k].P, 64)
			mutex.Lock()
			jPrice[skuId] = price
			mutex.Unlock()
		}
	}
	defer wg.Add(-1)
}
