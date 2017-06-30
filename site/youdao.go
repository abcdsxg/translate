package site

import (
	"errors"
	"fmt"

	"github.com/PuerkitoBio/goquery"
)

type Youdao struct {
}

func (y Youdao) GetTranslate(text string) (error, string,string,[]string) {
	fmt.Println("url:https://www.youdao.com/w/eng/" + text + "/#keyfrom=dict2.index")
	doc, err := goquery.NewDocument("https://www.youdao.com/w/eng/" + text + "/#keyfrom=dict2.index")
	if err != nil {
		fmt.Println(err)
		return errors.New("获取翻译失败"), "","",nil
	}
	translates := doc.Find(`.trans-container>ul`).First().Find(`li`)
	if translates.Length() == 0 {
		fmt.Println("err:", err)
		return errors.New("该词暂无释义"), "","",nil
	}
	var results string
	translates.Each(func(i int, translate *goquery.Selection) {

		fmt.Printf("%d:%s", i, translate.Text())
		results +=  translate.Text()+"\n"

	})
	//s := doc.Find(`span.pronounce`).Find(`span.phonetic`)
	//if s.Length() == 0 {
	//	fmt.Println("err:", err)
	//	return errors.New("该词暂无音标"), "","",nil
	//}
	//var phonetics string
	//s.Each(func(i int, ss *goquery.Selection) {
	//
	//	fmt.Printf("%d:%s", i, ss.Text())
	//	phonetics +=  ss.Text()+"\n"
	//
	//})

	return nil, results,"",nil
}
func (c Youdao) GetName()string{
	return "youdao"
}