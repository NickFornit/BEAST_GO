/* выдать дерево фраз на Пульт по GET-запросу http://go:8181?get_phrase_tree=1
Запросы даются при запуске и при изменении размера файла дерева /memory_reflex/word_tree.txt
 */

package word_sensor

import "strconv"

func initPrasePult() {
	//str:=GetPhraseTreeForPult()
	//if len(str)>0{}
}

// образ дерева фраз для вывода на Пульт
var praseTreeModel = ""

// проход дерева фраз
func GetPhraseTreeForPult() string {
	if notAllowScanInThisTime { return "!!!" }
	praseTreeModel = ""
	scanPraseNodes(-1, &VernikePhraseTree)

	return praseTreeModel
}

// сканировать узел дерева фраз
func scanPraseNodes(level int,wt *PhraseTree) {
	if wt.ID > 0 {
		praseTreeModel += setShift(level)
		// получить слово из ID конечного узла дерева слов
		praseTreeModel += GetWordFromWordID(wt.WordID) + "(" + strconv.Itoa(wt.ID) + " WordID="+strconv.Itoa(wt.WordID)+")<br>\r\n"
	}
	level++
	for n := 0; n < len(wt.Children); n++ {
		scanPraseNodes(level,&wt.Children[n])
	}
}

// отступ
func setShift(level int) string {
	var sh = ""
	for n := 0; n < level; n++ {
		sh += "&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;"
	}
	return sh
}