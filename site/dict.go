package site

import (
	"errors"
	//"fmt"

	"encoding/json"
	"net/url"

	"strconv"

	"github.com/PuerkitoBio/goquery"
)

type Dict struct {
}

func (d Dict) GetTranslate(text string) (error, string, string, []string) {
	doc, err := goquery.NewDocument("http://dict.cn/" + text)
	if err != nil {
		return errors.New("获取翻译失败"), "", "", nil
	}
	props := doc.Find(`ul.dict-basic-ul>li`)
	if props.Length() == 0 {
		return errors.New("该词暂无释义"), "", "", nil
	}
	var translatesArr string

	percents, _ := doc.Find(`.dict-chart`).Attr(`data`)
	percents, err = url.QueryUnescape(percents)

	percents, err = parsePercent(translatesArr, percents)
	if err != nil {
		return errors.New("解析词频失败"), "", "", nil
	}
	translatesArr += percents
	props.Each(func(i int, prop *goquery.Selection) {
		result := prop.Find(`span`).Text()
		translates := prop.Find(`strong`)
		if translates.Length() != 0 {
			translates.Each(func(j int, translate *goquery.Selection) {
				result += translate.Text()
				//fmt.Println(result)
				translatesArr += result + "\n"
			})
		}
	})
	s := doc.Find(`.phonetic>span`)
	if s.Length() == 0 {
		//fmt.Println("err:", err)
		return errors.New("该词暂无音标"), "", "", nil
	}
	var phonetics string
	var audios []string
	s.Each(func(i int, ss *goquery.Selection) {
		phonetic := ss.Find(`bdo`)
		//fmt.Printf("%d:%s", i, phonetic.Text())
		if i == 0 {
			phonetics += "英：" + phonetic.Text() + "\t"
		} else {
			phonetics += "美：" + phonetic.Text() + "\n"
		}
		audio, _ := ss.Find(`i.sound`).Attr("naudio")
		audios = append(audios, audio)

	})
	return nil, translatesArr, phonetics, audios
}
func parsePercent(translatesArr string, percents string) (string, error) {
	dicts := make(map[int]struct {
		Percent int
		Sense   string
	}, 0)
	var percentCon string
	err := json.Unmarshal([]byte(percents), &dicts)
	if err != nil {
		return "", err
	}
	for i:=1;i<len(dicts);i++{
		lineStr:="\t\t"
		if i%3 == 0 {
			lineStr="\n"
		}
		percentCon += dicts[i].Sense + ":" + strconv.Itoa(dicts[i].Percent) + "%"+lineStr

	}
	percentCon += "\n\n"
	translatesArr += percentCon
	return translatesArr, nil
}
func (c Dict) GetName() string {
	return "dict"
}
