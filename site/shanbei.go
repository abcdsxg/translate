package site

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Shanbei struct {
	Msg  string
	Data struct {
		Pronunciations struct {
			Uk string
			Us string
		}
		Audio_addresses struct {
			Uk []string
			Us []string
		}
		Definitions struct {
			Cn []struct {
				Pos  string
				Defn string
			}
		}
		Id int
	}
}

func (s Shanbei) GetTranslate(text string) (error, string, string, []string) {
	response, err := http.Get("https://www.shanbay.com/api/v1/bdc/search/?version=2&word=" + text + "&_=" + time.Now().String())
	if err != nil {
		fmt.Println(err)
		return errors.New("获取翻译失败"), "", "",nil
	}
	body := response.Body
	defer body.Close()
	bodys, err := ioutil.ReadAll(body)
	if err != nil {
		fmt.Println(err)
		return errors.New("获取翻译失败"), "", "",nil
	}
	fmt.Print()
	var shanbei Shanbei
	err = json.Unmarshal(bodys, &shanbei)
	if err != nil {
		fmt.Println(err)
		return errors.New("解析失败"), "", "",nil
	}
	var translates string
	for _, v := range shanbei.Data.Definitions.Cn {
		translates +=  v.Pos+v.Defn+"\n"
	}
	var audios []string
	audios = append(audios, shanbei.Data.Audio_addresses.Uk...)
	audios = append(audios, shanbei.Data.Audio_addresses.Us...)

	//var phonetics string
	//phonetics +=  shanbei.Data.Pronunciations.Uk+"\n"
	//phonetics += shanbei.Data.Pronunciations.Us+"\n"
	//
	//return nil, translates, phonetics, audios
	return nil, translates, "", audios
}
func (c Shanbei) GetName()string{
	return "shanbei"
}