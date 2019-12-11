##Installation
1.Install jd-spider
```
$ go get github.com/matchseller/jd-spider
```
2.Import it in your code:
```
import (
	"github.com/matchseller/jd-spider/category"
	"github.com/matchseller/jd-spider/price"
	"github.com/matchseller/jd-spider/product"
)
```
##Usage
1.Crawl categories
```
func main(){
    categoryUrls, err := category.Crawl()
}
```
2.Crawl products
```
func main(){
    categoryUrls, err := category.Crawl()
    if err == nil {
        products, _ := product.Crawl(categoryUrls[:10])
    }
}
```
3.Crawl price
```
func main(){
    categoryUrls, err := category.Crawl()
    if err == nil {
        products, _ := product.Crawl(categoryUrls[:10])
        var skuIds []string
        for _, v := range products{
            skuIds = append(skuIds, v.SkuId)
        }
        jPrices := price.Crawl(&skuIds)
    }
}
```