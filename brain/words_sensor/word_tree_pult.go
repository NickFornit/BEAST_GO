/* выдать дерево слов на Пульт по GET-запросу http://go:8181?get_word_tree=1
Запросы даются при запуске и при изменении размера файла дерева /memory_reflex/word_tree.txt
 */

package word_sensor

import (
	"sort"
	"strconv"
)

func initWordPult(){
	//str:=GetPhraseTreeForPult()
	//if len(str)>0{}
}

// образ дерева фраз для вывода на Пульт
var wordTreeModel = ""

// проход дерева фраз
func GetWordTreeForPult() string {
	if notAllowScanInThisTime { return "!!!" }
	wordTreeModel = ""
	scanWordNodes(-1, &VernikeWordTree)
	return wordTreeModel
}

// сканировать узел дерева слов
func scanWordNodes(level int, wt *WordTree) {
	if wt.ID > 0 {
		wordTreeModel += setWordShift(level)
		wordTreeModel += wt.Symbol + "(" + strconv.Itoa(wt.ID) + ")<br>\r\n"
	}
	level++
	for n := 0; n < len(wt.Children); n++ {
		scanWordNodes(level,&wt.Children[n])
	}
}

// отступ
func setWordShift(level int) string {
	var sh = ""
	for n := 0; n < level; n++ {
		sh += "&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;"
	}
	return sh
}

// выдать на пульт список слов в алфавитном порядке
func GetWordsListForPult() string {
	if len(WordIdFormWord) == 0 {
		return "Еще не готовы данные, обновите чуть позже."
	}

	wArr := make([]string, 0, len(WordIdFormWord))
	for w, _ := range WordIdFormWord {
		wArr = append(wArr, w)
	}
	sort.Strings(wArr)
	var out = "<table class='main_table'  cellpadding=0 cellspacing=0 border=1 width='100%' style='font-size:14px;'>"
	out += "<tr><td class='table_header'>Слово</td>"
	out += "<td class='table_header' width=100>ID</td>"
	out += "<td class='table_header'>Слово</td>"
	out += "<td class='table_header' width=100>ID</td>"
	out += "<td class='table_header'>Слово</td>"
	out += "<td class='table_header' width=100>ID</td></tr><tr>"
	var col = 0
	for n := 0; n < len(wArr); n++ {
		if len(wArr[n]) == 0 { continue	}
		if col >= 3 {
			out += "</tr><tr>"
			col = 0
		}
		id := strconv.Itoa(WordIdFormWord[wArr[n]])
		out += "<td class='table_cell'>" + wArr[n] + "</td>"
		out += "<td class='table_cell'>" + id + "<img src='/img/delete.gif' class='select_control' onClick='delete_word(" + id + ")'></td>"
		col++
	}
	count := strconv.Itoa(len(wArr))
	out += "</tr></table><b>Всего: " + count + " слов</b>"
return out
}