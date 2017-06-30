package site

import (
	"errors"
	"fmt"

	"github.com/PuerkitoBio/goquery"
)

type Ciba struct {
}

func (c Ciba) GetTranslate(text string) (error, string, string,[]string) {
	doc, err := goquery.NewDocument("http://www.iciba.com/" + text)
	if err != nil {
		fmt.Println(err)
		return errors.New("获取翻译失败"), "", "",nil
	}
	props := doc.Find(`span.prop`)
	translates := doc.Find(`li.clearfix>p`)
	if props.Length() == 0 {
		fmt.Println("err:", err)
		return errors.New("该词暂无释义"), "", "",nil
	}
	var translatesArr string
	//好变态!!我能怎么办
	props.Each(func(i int, prop *goquery.Selection) {
		result := prop.Text()
		translates.Each(func(j int, ss *goquery.Selection) {
			if i == j {
				spans := ss.Find(`span`)
				var spanText string
				spans.Each(func(k int, span *goquery.Selection) {
					spanText += span.Text()
				})
				result += spanText
				fmt.Println(result)
				translatesArr +=  result+"\n"
			}
		})
	})
	//s := doc.Find(`.base-speak>span>span`)
	//if s.Length() == 0 {
	//	fmt.Println("err:", err)
	//	return errors.New("该词暂无音标"), "", "",nil
	//}
	//var phonetics string
	//s.Each(func(i int, ss *goquery.Selection) {
	//
	//	fmt.Printf("%d:%s", i, ss.Text())
	//	phonetics +=  ss.Text()+"\n"
	//
	//})
	return nil, translatesArr, "",nil
}
func (c Ciba) GetName()string{
	return "ciba"
}