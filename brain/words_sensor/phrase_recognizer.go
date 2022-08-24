/*  распознавание фразы
вербальная иерархия распознавателей

Распознавание фраз начинается в main.go с word_sensor.VerbalDetection(text_dlg, is_input_rejim, moodID)
Память о воспринятых фразах в текущем активном контексте (Vernike_detector.go): var MemoryDetectedArr []MemoryDetected
*/

package word_sensor

import (
	_ "strconv"
	_ "strings"
)

///////////////////////////////////////////////////

func iniPraseRecognising(){
	/*


	//%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%

	*/
	//%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%

}
//////////////////////////////////////////////////////

// текущий уникальный ID последней активной ветки дерева - результат детекции фразы - для дальнейшего использования
var DetectedUnicumPhraseID=0
// нераспознанный остаток
var CurrentVerbalPhraseEnd []int
// текущий найденный ID последней активной ветки дерева слов
var DetectedCurrentPhraseID=0
// текущий номер слова распознаваемой фразы
var currentStepPhraseCount=0

/////////////////////////////////////////////////////////////////////////

/* переносим в дерево слов достаточно повторяющиеся из tempArr для trees_former.go

Распознать и вставить новое слово-фразу в дерево:
найти подходящий узел и если еще нет - вставить новый.
 */


/////////////////////////////////////////////////////////
// проход одной фразы - распознавание ID слов фразы
func PhraseDetection(words []int)(string) {

	if len(words) == 0 {
		return ""
	}
	if len(words) == 1 && words[0]==0 {// пустые строки не писать
		return ""
	}
	CurrentVerbalPhraseEnd = nil
	DetectedUnicumPhraseID = 0
	//var pultOut=""
	DetectedCurrentPhraseID = 0
	currentStepPhraseCount = len(words)

	r := words
	// основа дерева
	cnt := len(VernikePhraseTree.Children)
	for n := 0; n < cnt; n++ {
		phraseNode := VernikePhraseTree.Children[n]
		rt := phraseNode.WordID
		if r[0] == rt {
			cldrn := VernikePhraseTree.Children[n] //.Children
			getPhraseTreeNode(r, &cldrn)
		}

		if currentStepPhraseCount==0 { // распознанно точно, не смотреть другие
			break
		}
	}

	//////////////// результат распознавания
	if DetectedCurrentPhraseID > 0 {
		if currentStepPhraseCount==0 { // полностью распознан
			DetectedUnicumPhraseID=DetectedCurrentPhraseID
		}else{
			var nr=len(r)-currentStepPhraseCount
			CurrentVerbalPhraseEnd=r[nr:]
		}
	}
	/////////////////////////////////
	var needSave=false
	if DetectedUnicumPhraseID ==0 {
		// нераспознанный остаток
		if len(CurrentVerbalPhraseEnd) > 0 {
			r := CurrentVerbalPhraseEnd
			var tree *PhraseTree
			if DetectedCurrentPhraseID > 0 {
				tree = PhraseTreeFromID[DetectedCurrentPhraseID]
			} else {
				tree = &VernikePhraseTree
			}
			// просто добавить новую ветку - из диалога это стоит делать за 1 раз т.к. слова уже известны
			node := createNewNodePhraseTree(tree, 0, r[0])
			tree = node
			id := createPhraseTreeNodes(r, PhraseTreeFromID[tree.ID])
			DetectedUnicumPhraseID = id
			needSave=true
		}
	}

	// нет вообще такого, добавить все слово
	if DetectedUnicumPhraseID ==0{
		tree := PhraseTreeFromID[0]
		// сразу создать первый узел
		if len(r) > 0 {
			node := createNewNodePhraseTree(tree, 0, r[0])
			tree = node
			if tree !=nil {
				id := createPhraseTreeNodes(r, PhraseTreeFromID[tree.ID])
				DetectedUnicumPhraseID = id
				needSave = true
			}
		}
	}
	if needSave{
		SavePhraseTree()
	}


	out:=GetPhraseStringsFromPhraseID(DetectedUnicumPhraseID)

return out//pultOut+"{"+strconv.Itoa(DetectedUnicumPhraseID)+")"
}
/////////////////////////////////////////////


func getPhraseTreeNode(words []int,wt *PhraseTree){

	if len(words)==0{
		return
	}
	ost:=words[1:]

	if words[0] != wt.WordID {// пошло не туда
		return
	}

	DetectedCurrentPhraseID=wt.ID
	currentStepPhraseCount=len(ost)

	for n := 0; n < len(wt.Children); n++ {
			getPhraseTreeNode(ost, &wt.Children[n])
	}
	return
}
//////////////////////////////////////////////////////////////////////


func getNodeCountFromLastID(lastID int)(int){
	if lastID==0{
		return 0
	}
	var count=0
	for {
		node:=PhraseTreeFromID[lastID]
		if node==nil || node.WordID==0 {
			break
		}
		count++
		lastID=node.ParentID
	}
	return count
}








///////////////////////////////////////
// выдать строку из массива wordsArr[]int
func GetStrFromArrID(wArr []int)(string){
var out=""
	for i := 0; i < len(wArr); i++ {
		out += GetWordFromWordID(wArr[i])+" "
	}
	return out
}


