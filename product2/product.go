package product2

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/matchseller/jd-spider/util"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

type Product struct {
	SkuId       string `json:"sku_id"`
	ProductName string `json:"product_name"`
	ImgPath     string `json:"img_path"`
}

var retryMax = 30

//start crawling
func Crawl(categoryUrl string) (products []Product) {
	totalPage := getTotalPage(categoryUrl)
	for page := 1; page <= totalPage; page++ {
		products = append(products, run(fmt.Sprintf("%s&page=%d", categoryUrl, page))...)
	}
	return
}

func SetRetryMax(number int) {
	retryMax = number
}

func getTotalPage(categoryUrl string) (totalPage int) {
	retStr := do(categoryUrl)
	if retStr != "" {
		dom, err := goquery.NewDocumentFromReader(strings.NewReader(retStr))
		if err == nil {
			totalPage, _ = strconv.Atoi(dom.Find("span.fp-text>i").First().Text())
		}
	}
	return totalPage
}

func run(url string) (products []Product) {
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
						productName = strings.ReplaceAll(productName, "\n", " ")
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
	return
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
