package product

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/matchseller/jd-spider/util"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

type pageInfo struct {
	url  string `json:"url"`
	page int    `json:"page"`
}

type Product struct {
	SkuId       string `json:"sku_id"`
	ProductName string `json:"product_name"`
	ImgPath     string `json:"img_path"`
}

var retryMax = 30
var pInfoChan = make(chan pageInfo)
var productChan = make(chan []Product)

//start crawling
func Crawl(categoryUrlList []string) (products []Product, zeroPages []string) {
	maxTotalPage, cTotalPage, zeroPages := initPage(&categoryUrlList)
	for page := 1; page <= maxTotalPage; page++ {
		routineCount := 0
		for _, url := range categoryUrlList {
			if cTotalPage[url] <= 0 {
				continue
			}
			if page <= cTotalPage[url] {
				routineCount++
				go run(fmt.Sprintf("%s&page=%d", url, page))
			}
		}
		for i := 0; i < routineCount; i++ {
			products = append(products, <-productChan...)
		}
	}
	return
}

func SetRetryMax(number int) {
	retryMax = number
}

func initPage(categoryUrlList *[]string) (maxTotalPage int, cTotalPage map[string]int, zeroPages []string) {
	cTotalPage = make(map[string]int)
	for _, categoryUrl := range *categoryUrlList {
		go getTotalPage(categoryUrl)
	}
	var accept pageInfo
	for i := 0; i < len(*categoryUrlList); i++ {
		accept = <-pInfoChan
		if accept.page == 0 {
			zeroPages = append(zeroPages, accept.url)
		} else {
			cTotalPage[accept.url] = accept.page
			if accept.page > maxTotalPage {
				maxTotalPage = accept.page
			}
		}
	}
	return
}

func getTotalPage(categoryUrl string) {
	var pInfo pageInfo
	var totalPage int
	retStr := do(categoryUrl)
	if retStr != "" {
		dom, err := goquery.NewDocumentFromReader(strings.NewReader(retStr))
		if err == nil {
			totalPage, _ = strconv.Atoi(dom.Find("span.fp-text").First().Find("i").First().Text())
		}
	}
	pInfo.url = categoryUrl
	pInfo.page = totalPage
	defer func(pInfo pageInfo) {
		pInfoChan <- pInfo
	}(pInfo)
}

func run(url string) {
	var products []Product
	retStr := do(url)
	if retStr != "" {
		dom, err := goquery.NewDocumentFromReader(strings.NewReader(retStr))
		if err == nil {
			dom.Find("div.j-sku-item").Each(func(i int, selection *goquery.Selection) {
				imgPath, _ := selection.Find("div.p-img>a").First().Find("img").First().Attr("src")
				if imgPath == "" {
					imgPath, _ = selection.Find("div.p-img>a").First().Find("img").First().Attr("data-lazy-img")
				}
				productUrl, _ := selection.Find("div.p-name>a").First().Attr("href")
				if productUrl != "" && imgPath != "" {
					reg := regexp.MustCompile(`\/([0-9]+).html`)
					regRes := reg.FindStringSubmatch(productUrl)
					if len(regRes) > 1 {
						productName := util.TrimHtml(selection.Find("div.p-name>a").First().Find("em").First().Text())
						product := Product{
							SkuId:       strings.Trim(regRes[1], " "),
							ProductName: strings.Trim(productName, " "),
							ImgPath:     "https:" + strings.Trim(imgPath, " "),
						}
						if product.ProductName != "" && product.SkuId != "" {
							if utf8.RuneCountInString(product.ProductName) <= 500 {
								products = append(products, product)
							}
						}
					}
				}
			})
		}
	}
	defer func(products []Product) {
		productChan <- products
	}(products)
}

func do(url string) (retStr string) {
	reTryStart := 0
	for {
		ret, _ := util.GetUrl(url)
		retStr = string(ret)
		if (retStr != "" && !strings.Contains(retStr, "JengineD/1.7.2.1")) || reTryStart < retryMax {
			break
		}
		reTryStart++
	}
	if strings.Contains(retStr, "JengineD/1.7.2.1") {
		retStr = ""
		return
	}
	return
}
