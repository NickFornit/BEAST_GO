/*  дерево фраз
вербальная иерархия распознавателей
Первый уровень дерева фраз может может заполняться любыми ID слов
Память о воспринятых фразах в текущем активном контексте (Vernike_detector.go): var MemoryDetectedArr []MemoryDetected
*/

package word_sensor

import (
	"BOT/lib"
	"regexp"
	"strconv"
	"strings"
)

// подошла очередь инициализации
func afterLoadPhraseArr(){
	loadPhraseTree()
	/*
	//%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
	SetNewPhraseTreeNode("повести и игра") //
	//%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
	SavePhraseTree()
	*/
	iniPraseRecognising()
	afetrInitPhraseTree()
	// для старых слов получить WordIdFormWord
	getWordIdFormWord()
	// deleteWordFromPhrase(11)
}

// дерево фраз, разбитых на слова, формат записи ID|ParentID|#|WordID
type PhraseTree struct {
	ID int 				// id узла слова
	ParentID int  // ID родителя
	WordID int 		// одно  слово, м.б. пробелорм или любым символом

	Children []PhraseTree 	// дочерние узлы (ветвление) НЕ АДРЕСА, А РЕАЛЬНЫЕ ОБЪЕКТЫ
	ParentNode *PhraseTree  // адрес родителя
}

// дерево фраз
var VernikePhraseTree PhraseTree
// карта поиска дерева фраз
var PhraseTreeFromID = make(map[int]*PhraseTree)

// Последовательность wordID в ветке дерева получается из GetWordArrFromPhraseID(PhraseID int)
var WordsArrFromPhraseID = make(map[int][]int)

// Все ID фраз по wordID: в каких ID фраз содержится данное слово
var PhraseTreeFromWordID = make(map[int][]int)

// для обеспечения уникальности узлов:
/*  лишнее
type PhraseUnicum struct {
	ID int
	wordID int
}
var PhraseUnicumIdStr=make(map[PhraseUnicum]int)// для каждого сочетания  выдается ID узла
*/

var lastPhraseTreeID = 0 // конечный узел дерева фраз

// создать новый узел дерева фраз
func createNewNodePhraseTree(parent *PhraseTree, id int, wordID int) *PhraseTree {
	if parent == nil {	return nil 	}
	//if wordID==0{ return nil }

	// после удаления слова - запрет на вставку новых слов до перезагрузки
	if blockingNewInsertWordAfterDeleted {
		lib.WritePultConsol("ПОСЛЕ УДАЛЕНИЯ СЛОВА - ЗАПРЕТ НА ВСТАВКУ НОВЫХ СЛОВ ДО ПЕРЕЗАГРУЗКИ")
		return nil
	}

	// notAllowScanInThisTime=true // запрет показа карты при обновлении
	if id == 0 {
		lastPhraseTreeID++
		id = lastPhraseTreeID
	} else {
		if lastPhraseTreeID < id {
			lastPhraseTreeID = id
		}
	}

	var newW PhraseTree
	newW.ID = id
	newW.ParentID = parent.ID
	newW.ParentNode = parent
	newW.WordID = wordID

	parent.Children = append(parent.Children, newW)
	// четко находим новый вставленный член (а не &parent.Children[count-1])
	var newTP *PhraseTree
	for i := 0; i < len(parent.Children); i++ {
		if parent.Children[i].ID == newW.ID {
			newTP = &parent.Children[i]
		}
	}
	//!!!! PhraseTreeFromID[new.ID]=new  т.к. append меняет длину массива, перетусовывая адреса, то нужно:
	scanAllTree(parent) // здесь потому, что при загрузке из файла нужно на лету получать адреса
	/*if newTP != nil {
		WordsArrFromPhraseID[newTP.ID] = append(WordsArrFromPhraseID[newTP.ID], newTP.WordID)
		PhraseTreeFromWordID[newTP.WordID] = append(PhraseTreeFromWordID[newTP.WordID], newTP)
	}*/
	// notAllowScanInThisTime=false

	return newTP
}

// корректируем адреса всех узлов
func scanAllTree(parent *PhraseTree){
	updatingPhraseTreeFromID(parent)
}
func updatingPhraseTreeFromID(wt *PhraseTree) {
	if wt.ID > 0 {
		wt.ParentNode = PhraseTreeFromID[wt.ParentID] // wt.ParentNode адрес меняется из=за corretsParent(,
		PhraseTreeFromID[wt.ID] = wt
	}
	if wt.Children == nil {	return } // конец ветки
	for i := 0; i < len(wt.Children); i++ {
		updatingPhraseTreeFromID(&wt.Children[i])
	}
}



// Загрузка дерева фраз
func loadPhraseTree() {
	initPhraseTree(&VernikePhraseTree)
	// initPhraseTree(&VernikePhraseTree)
	strArr, _ := lib.ReadLines(lib.GetMainPathExeFile() + "/memory_reflex/phrase_tree.txt")
	cunt := len(strArr)
	// просто проход по всем строкам файла подряд так что сначала идут дочки, потом - их родители
	for n := 0; n < cunt; n++ {
		if len(strArr[n]) < 2 {
			panic("Сбой загрузки дерева фраз: ["+strconv.Itoa(n) + "] " + strArr[n])
			return
		}
		p := strings.Split(strArr[n], "|#|")
		id, _ := strconv.Atoi(p[1])
		wordID := id
		if WordTreeFromID[wordID] == nil { continue } // нет такого узла дерева слов
		idP := strings.Split(p[0], "|")
		id, _ = strconv.Atoi(idP[0])
		parentID, _ := strconv.Atoi(idP[1])
		// новый узел с каждой строкой из файла
		createNewNodePhraseTree(PhraseTreeFromID[parentID], id, wordID)
	}

// заполнить PhraseTreeFromWordID
	finishScanAllTree()

	return
}
var curBrangeArr []int
func finishScanAllTree(){
	WordsArrFromPhraseID = make(map[int][]int)
	PhraseTreeFromWordID = make(map[int][]int)
	curBrangeArr=nil
	curScanAllTree(&VernikePhraseTree)
}
func curScanAllTree(wt *PhraseTree) {
	if wt.ID > 0 {
		wt.ParentNode = PhraseTreeFromID[wt.ParentID] // wt.ParentNode адрес меняется из=за corretsParent(,
		curBrangeArr=append(curBrangeArr,wt.WordID)
	}
	if wt.Children == nil {	// конец ветки
		WordsArrFromPhraseID[wt.ID]=curBrangeArr
		// перебрать все слова                       wordID - ЭТО НЕТ nodePhraseID !!!!!!!!!
		for i := 0; i < len(curBrangeArr); i++ {
			PhraseTreeFromWordID[curBrangeArr[i]]=append(PhraseTreeFromWordID[curBrangeArr[i]],wt.ID)
		}
		if len(curBrangeArr)>1{
			curBrangeArr=nil
		}
		if len(PhraseTreeFromWordID[wt.ID])>1{
			curBrangeArr=nil
		}
		curBrangeArr=nil
		return
	}
	for i := 0; i < len(wt.Children); i++ {
		//curBrangeArr=nil
		curScanAllTree(&wt.Children[i])
	}
}
////////////////////////////////////////////////////////////





/////////////////////////////////////////////////////////
// создать первый, нулевой уровень дерева
func initPhraseTree(vt *PhraseTree) {
	// createNewNodePhraseTree(vt,0,0)
	vt.ID = 0
	vt.WordID = 0
	PhraseTreeFromID[vt.ID] = vt
	//updateWordTreeFromID()
	return
}

// Сохранить дерево фраз
// ID|ParentID|#|WordID
func SavePhraseTree() {
	var out = ""
	cnt := len(VernikePhraseTree.Children)
	for n := 0; n < cnt; n++ {
		out += getPtreeNode(&VernikePhraseTree.Children[n])
	}
	lib.WriteFileContent(lib.GetMainPathExeFile() + "/memory_reflex/phrase_tree.txt", out)
	return
}

// получить ветку дерева фраз в виде строки
func getPtreeNode(wt *PhraseTree) string {
	var out = ""

	out += strconv.Itoa(wt.ID) + "|"
	out += strconv.Itoa(wt.ParentID) + "|#|"
	out += strconv.Itoa(wt.WordID) + "\r\n"

	if wt.Children == nil {	return out } // конец
	for n := 0; n < len(wt.Children); n++ {
		out += getPtreeNode(&wt.Children[n])
	}
	return out
}

/* вставка новой фразы со вставкой новых слов фразы,
так что фраза будет распознанна всегда.
 */
func SetNewPhraseTreeNode(word string) *WordTree {
	// чистим лишние пробелы
	rp := regexp.MustCompile("s+")
	word = rp.ReplaceAllString(word, " ")
	word = strings.TrimSpace(word)

	var wordsIDstr[] int // строка (не)распознанных слов

	/* сначала добавляем слова в дерево слов, потом - всю фразу в дерево фраз
	   Делим фразу на слова (в строке нет других разделительных символов,
	т.к. они уже сработали при разделении на фразы).
	*/
	wArr := strings.Split(word, " ")
	for n := 0; n < len(wArr); n++ { // перебор отдельных слов
		curWord := strings.TrimSpace(wArr[n])
		if len(curWord) == 0 { return nil	}

		id := SetNewWordTreeNode(curWord)
		// распознавание будет ВСЕГДА т.к. в случае новго слова оно вставляется в дерево слов тут же
		wordsIDstr = append(wordsIDstr, id)
	} //for n := 0; n < len(wArr); n++ { закончен проход отдельных слов
	// updateWordTreeFromID()// обновляем массив адресов узлов после всех append() родителей, меняющих адреса

	//  проход фразы
	// var needSave=false
	if len(wordsIDstr) > 0 {
		PhraseDetection(wordsIDstr)
		// Запись недостающего в дерево фраз происходит в PhraseDetection(wordsIDstr)
	}
	// if needSave {
	//	savePhraseTree()
	// }
	return nil
}

// создание ветки фраз, начиная с заданного узла
func createPhraseTreeNodes(word []int, wt *PhraseTree) int {
	ost := word[1:]
	if len(ost) == 0 { return wt.ID	}

	node:=createNewNodePhraseTree(PhraseTreeFromID[wt.ID], 0, ost[0])
	id:=createPhraseTreeNodes(ost, PhraseTreeFromID[node.ID])

	return id
}