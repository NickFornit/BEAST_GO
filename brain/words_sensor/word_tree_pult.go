/*   выдать дерево слов на Пульт по GET-запросу http://go:8181?get_word_tree=1
Запросы даются при запуске и при изменении размера файла дерева /memory_reflex/word_tree.txt
 */

package word_sensor

import "strconv"

func initWordPult(){
	//str:=GetPhraseTreeForPult()
	//if len(str)>0{}
}
/////////////////////////

// образ дерева фраз для вывода
var wordTreeModel=""


// проход дерева фраз

func GetWordTreeForPult()(string){
	if notAllowScanInThisTime{
		return "!!!"
	}
	wordTreeModel=""
	scanWordNodes(-1,&VernikeWordTree)

	return wordTreeModel
}
//////////////////////

func scanWordNodes(level int,wt *WordTree){

	if wt.ID>0 {
		wordTreeModel += setWordShift(level)
		wordTreeModel += wt.Symbol + "(" + strconv.Itoa(wt.ID) + ")<br>\r\n"
	}
	level++
	for n := 0; n < len(wt.Children); n++ {
		scanWordNodes(level,&wt.Children[n])
	}
}
// отступ
func setWordShift(level int)(string){
var sh=""
	for n := 0; n < level; n++ {
		sh+="&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;"
	}
	return sh
}