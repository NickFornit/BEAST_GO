/*  функции для фраз

*/

package word_sensor

import (
	"regexp"
	"strings"
)

//////////////////////////////////////

//////////////////////////////////////////////////
/// строка из ID фразы дерева фраз
func GetPhraseStringsFromPhraseID(lastID int)(string){

	var idArr []string
	for {
		node:=PhraseTreeFromID[lastID]
		if node==nil {
			break
		}
		w:=GetWordFromWordID(node.WordID)
		idArr=append(idArr,w)
		lastID=node.ParentID
		if lastID==0{
			break
		}
	}

	var str=""
	for i := len(idArr)-1; i >=0; i-- {
		str+=idArr[i]+" "
	}

	return str
}
//////////////////////////////////////////


///////////////////////////////////////
// выдать строку из массива wordsArr[]int
// используется в update_genom.go
func GetStrFromArrID(wArr []int)(string){
	var out=""
	for i := 0; i < len(wArr); i++ {
		out += GetWordFromWordID(wArr[i])+" "
	}
	return out
}
////////////////////////////////////////


// очистить фразу от неалфавитных символов
func ClinerNotAlphavit(prase string)(string){

	var out=""
	reg := regexp.MustCompile(`[а-я ]`)
	res:=reg.FindAllString(prase,-1)
	for i := 0; i < len(res); i++ {
		out+=res[i]
	}

	return out
}
//////////////////////////////////////

// если есть такая фраза в Дереве, то выдать ее ID
func GetExistsPraseID(text string) int {
	var id = 0
	// чистим лишние пробелы
	rp := regexp.MustCompile("s+")
	text = rp.ReplaceAllString(text, " ")
	text = strings.TrimSpace(text)

	wordsArr := GetWordIDfromPhrase(text)
	str := PhraseDetection(wordsArr) // распознаватель фразы
	if len(str)>0 {
		id = DetectedUnicumPhraseID
	}
	return id
}