/*  функции для фраз

*/

package word_sensor

import (
	"BOT/lib"
	"regexp"
	"strconv"
	"strings"
)

//////////////////////////////////////
// слово из ID узла дерева фраз
func GetWordFromPraseNodeID(nodeID int)(string){
	if nodeID==0{
		return ""
	}
	ph:=PhraseTreeFromID[nodeID]
	word := GetWordFromWordID(ph.WordID)
	return word
}

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
		if len(str)>0{
			str+=" "
		}
		str+=idArr[i]
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
///////////////////////////////////////////////

// удалить слово все всех упоминаниях в Дереве фраз
func deleteWordFromPhrase(wordID int){
	strArr,_:=lib.ReadLines(lib.GetMainPathExeFile()+"/memory_reflex/phrase_tree.txt")
	var out=""
	var parentNewID=0
	var parentOdID=0
	for n := 0; n < len(strArr); n++ {
		p:=strings.Split(strArr[n], "|")
		id,_:=strconv.Atoi(p[0])
		pID,_:=strconv.Atoi(p[1])
		node:=PhraseTreeFromID[id]
		if node==nil{
			return
		}
		p=strings.Split(strArr[n], "|#|")
		wID,_:=strconv.Atoi(p[1])
		if wID==wordID{
			if len(node.Children) >0{// всем дочкам переписать родителей - node.ParentID
				parentNewID=node.ParentID
				parentOdID=node.ID
			}// если нет родителя, то можно просто удалить
			continue // не писать удаляемую строку
		}

			if pID>0 && pID == parentOdID { // заменить родителя
				out += strconv.Itoa(id) + "|" + strconv.Itoa(parentNewID) + "|#|" + strconv.Itoa(wID) + "\r\n"
			} else {
				out += strArr[n] + "\r\n"
			}
	}
	lib.WriteFileContent(lib.GetMainPathExeFile()+"/memory_reflex/phrase_tree.txt",out)
}
//////////////////////////////////////////////