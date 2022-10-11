/* сенсор символов, слов и фраз */
package word_sensor

import (
	"BOT/lib"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// вспомогательный шаблон воспринятых слов с числом повторений восприятия
var tempArr=make(map[string]int)

func init(){
	// залить сохранненный tempArr
	loadTempArr()
	afterLoadTempArr()
	//GetExistsPraseID("дурак")
}

// сохранить шаблон воспринятых слов
func SaveTempArr() {
	var str[]string
	for k, v := range tempArr {
		str = append(str, k + "|#|" + strconv.Itoa(v))
	}
	sort.Strings(str)
	var out = ""
	for i := 0; i < len(str); i++ {
		p:=strings.Split(str[i], "|#|")
		out+=p[1] + "|#|" + p[0] + "\r\n"
	}
	lib.WriteFileContent(lib.GetMainPathExeFile() + "/memory_reflex/words_temp_arr.txt", out)
}

// загрузить шаблон воспринятых слов
func loadTempArr() {
	tempArr = make(map[string]int)
	wArr, _ := lib.ReadLines(lib.GetMainPathExeFile() + "/memory_reflex/words_temp_arr.txt")
	for n := 0; n < len(wArr); n++ {
		if len(wArr[n]) < 4 { continue }
		p := strings.Split(wArr[n], "|#|")
		v, _ := strconv.Atoi(p[0])
		tempArr[p[1]] = v
	}
}

// Новая порция текста для формирования дерева слов
// Эта функция работает при накачке текстов из main.go
func SetNewTextBlock(txt string) string {
	/*
	txt, err := url.QueryUnescape(txt)
	if err != nil {
		log.Fatal(err)
		return "ОШИБКА раскодировки"
	}
 */
	txt = strings.Replace(txt, "{#1}", "%", -1)
	// txt= strings.Replace(txt, "{#2}", "", -1)// кавычки просто очищены (пусть будет афазия :)
	var res = ""
	// разделяем на фразы
	strArr := strings.Split(txt, "|#") // а не |#| - чтобы оставлять разделитель "|"
	for i := 0; i < len(strArr); i++ {
		addNewtempArr(strArr[i] + "|")
	}
	updateWordTreeFromTempArr(4,6)

	return res
}

/* добавляются как целиком фраза, так и все слова во фразе
Тут же дозаполняетс дерево слов уже многократно провторяющимися элементами.

!?. в конце фразы не записывать, они распознаются отдельно

ret bool если true - вернуть массив разделенных строк текста (режим диалога с Пультом, а не режим накачки текстами)
 */
func addNewtempArr(str string) []string {
	if len(str) < 2 { return []string{str} }
	var out_str []string
	// разделить слова
	// разбить по разделителям
	sArr := strings.Split(str, "|")
	var wordCount = 0
	// фраза без разделителей и отдельные слова (более 1 символа) - в шаблон:
	var out = ""
	var isBegin = true
	cnt := len(sArr)
	for i := 0; i < cnt; i++ {
		if len(sArr[i]) == 0 { continue	}
		// это - слово(более 1 символа), а символ
		r := []rune(sArr[i])
		if isBegin && len(r) > 1 {
			tempArr[sArr[i]]++
		}
		wordCount++
		// во фразе оставляем и слова и символы, т.е. восстанавливаем фразу как она была
		out += sArr[i]

		if sArr[i] == " " { isBegin = true} else { isBegin = false }
	}

	if wordCount < 6 { // не запоминать фразы, длинее 6 слов
		if wordCount > 1 { tempArr[out]++	} // это фраза (> 1 слова)
			out_str = append(out_str, out) // короткая фраза
		} else {
			// разбить по знакам препинания
			zp := regexp.MustCompile(`[,:.+\(\)]`)
			sPath := zp.Split(out, -1)
			for i := 0; i < len(sPath); i++ {
				s := strings.Trim(sPath[i], " ")
				sSps := strings.Split(s, " ")
				if len(sSps) > 6 { // не запоминать фразы, длинее 6 слов, обрезать их до 6 слов
					if wordCount > 1 { tempArr[s]++	} // это фраза (> 1 слова)
					out_str = append(out_str, s)
					continue
				}
				r := []rune(s)
				if len(r) < 2 { continue }
				if wordCount > 1 { // это фраза (> 1 слова)
					tempArr[s]++
				}
				out_str = append(out_str, s)
			}
		}
	SaveTempArr()

	return out_str
}