/*
дерево рефлексов, безусловных и условных

Узлы дерева создаются только рефлексами:
безусловными сразу при загрузке, условными по мере возникновения новых - в новых уловиях.
У дерева - только 3 уровня для безусловных рефлексов, все узла которых - в виде ID:
1 - базовое состояние - ID Плохо, Норма, хорошо
2 - сочетаний активных Базовых контекстов - ID BaseStyleArr
3 - сочетаний пусковых стимулов - ID TriggerStimulsArr
При возникновении условных рефлексов просто добавляется новый узел
- образ новых условий (2 или 3-го уровней), запускающих усл.рефлекс.

Дерево обязательно должно иметь три базовых состояния вначале - без рефлексов:
1|0|1|0|0|0|0
2|0|2|0|0|0|0
3|0|3|0|0|0|0

Формат записи безусловного рефлекса: ID|baseID|styleID...|actionID...
*/

package reflexes

import (
	"BOT/lib"
	"strconv"
	"strings"
)

func initReflexTree() { // после инициализации loadGeneticReflexes()
	// tools.GetAllCombinationsOfSeriesNumbers(5,3)
	loadReflexTree()
	if len(ReflexTree.Children) == 0 { // еще нет никаких веток
		// создать первые три ветки базовых состояний
		createBasicReflexTree()
	}
	// SaveReflexesAttributes()
	readyForRecognitionRflexes() // ини для дерева распознавания рефлексов
}

/* Структура дерева рефлексов
ID|parentID|baseID|styleID|actionID|GeneticReflexID|ConditionedReflex */
type ReflexNode struct { // узел дерева рефлексов
	ID int
	baseID int 		// базовое состояние
	StyleID int 	// стиль поведения - сочетание активностей Базовых контекстов  - ID BaseStyleArr
	ActionID int 	// сочетание пусковых стимулов  - ID TriggerStimulsArr

	GeneticReflexID int 		// безусловный рефлекс
	ConditionedReflex int 	// условный рефлекс, если есть, блокирует GeneticReflexID узла

	Children []ReflexNode 	// дочерние узлы (ветвление) НЕ АДРЕСА, А РЕАЛЬНЫЕ ОБЪЕКТЫ
	ParentID int     				// ID родителя
	ParentNode *ReflexNode  // адрес родителя
}

var ReflexTree ReflexNode // дерево рефлексов
var ReflexTreeFromID = make(map[int]*ReflexNode) // карта для поиска рефлексов

/* запрет показа карты на пульте (func GetReflexesTreeForPult()) при обновлении
против паники типа "одновременная запись и считывание карты"
Использовать для всех операций записи узлов дерева
 */
var notAllowScanInReflexesThisTime = false

var lastReflexNodeID = 0 // последний узео в дереве автоматизмов

// Создать новый узел дерева рефлексов
func createNewReflexNode(parent *ReflexNode, id int, baseID int, StyleID int,
				ActionID int, GeneticReflexID int, ConditionedReflex int,CheckUnicum bool)(int, *ReflexNode) {
	// если есть такой узел, то не создавать
	if CheckUnicum {
		idOld, nodeOld := FindReflexTreeNodeFromCondition(baseID, StyleID, ActionID)
		if idOld > 0 { return idOld,nodeOld }
	}

	if id == 0 {
		lastReflexNodeID++
		id = lastReflexNodeID
	} else {
		if lastReflexNodeID < id {
			lastReflexNodeID = id
		}
	}

	var node ReflexNode
	node.ID = id
	node.ParentNode = parent
	node.ParentID = parent.ID
	node.baseID = baseID
	node.StyleID = StyleID
	node.ActionID = ActionID
	node.GeneticReflexID = GeneticReflexID
	node.ConditionedReflex = ConditionedReflex

	parent.Children = append(parent.Children, node)
	// четко находим новый вставленный член (а не &parent.Children[count-1])
	var newN *ReflexNode
	for i := 0; i < len(parent.Children); i++ {
		if parent.Children[i].ID == node.ID {
			newN = &parent.Children[i]
		}
	}
	// т.к. append меняет длину массива, перетусовывая адреса, то нужно обновить адреса в ReflexTreeFromID:
	updateReflexTreeFromID(parent) // здесь потому, что при загрузке из файла нужно на лету получать адреса

	return id, newN
}

// корректируем адреса всех узлов
func updateReflexTreeFromID(parent *ReflexNode) {
	// updatingPhraseTreeFromID(&VernikePhraseTree)
	updatingPhraseTreeFromID(parent)
}

// проход всего дерева
func updatingPhraseTreeFromID(rt *ReflexNode) {
	if rt.ID > 0 {
		rt.ParentNode = ReflexTreeFromID[rt.ParentID] // wr.ParentNode адрес меняется из=за corretsParent(,
		ReflexTreeFromID[rt.ID] = rt
	}
	// конец ветки
	if rt.Children == nil { return }
	for i := 0; i < len(rt.Children); i++ {
		updatingPhraseTreeFromID(&rt.Children[i])
	}
}

// загрузить записанное дерево
// ID|parentID|baseID|styleID|actionID|geneticReflexID|conditionedReflex|
func loadReflexTree(){
	createNulLevelReflexTree(&ReflexTree)
	strArr, _ := lib.ReadLines(lib.GetMainPathExeFile() + "/memory_reflex/reflex_tree.txt")
	cunt := len(strArr)
	// просто проход по всем строкам файла подряд так что сначала идут дочки, потом - их родители
	for n := 0; n < cunt; n++ {
		if len(strArr[n]) == 0 { continue	}
		if len(strArr[n]) < 2 {
			panic("Сбой загрузки дерева рефлексов: ["+strconv.Itoa(n) + "] " + strArr[n])
			return
		}
		p := strings.Split(strArr[n], "|")
		id, _ := strconv.Atoi(p[0])
		parentID, _ := strconv.Atoi(p[1])
		baseID, _ := strconv.Atoi(p[2])
		styleID, _ := strconv.Atoi(p[3])
		actionID, _ := strconv.Atoi(p[4])
		geneticReflexID, _ := strconv.Atoi(p[5])
		conditionedReflex, _ := strconv.Atoi(p[6])
		// новый узел с каждой строкой из файла
		createNewReflexNode(ReflexTreeFromID[parentID], id, baseID, styleID,
			actionID,geneticReflexID,conditionedReflex,false)
	}
	return
}

// создать первый, нулевой уровень дерева
func createNulLevelReflexTree(rt *ReflexNode) {
	rt.ID = 0
	ReflexTreeFromID[rt.ID] = rt
	return
}

// создать первые три ветки базовых состояний
func createBasicReflexTree() {
	notAllowScanInReflexesThisTime = true // запрет показа карты при обновлении
	createNewReflexNode(&ReflexTree,0,1,0,0,0,0,false)
	createNewReflexNode(&ReflexTree,0,2,0,0,0,0,false)
	createNewReflexNode(&ReflexTree,0,3,0,0,0,0,false)
	saveReflexTree()
	notAllowScanInReflexesThisTime = false // запрет показа карты при обновлении
	return
}

// сохранить дерево рефлексов
func saveReflexTree() {
	var out = ""
	cnt := len(ReflexTree.Children)
	for n := 0; n < cnt; n++ {
		out += getReflexNode(&ReflexTree.Children[n])
	}
	lib.WriteFileContent(lib.GetMainPathExeFile() + "/memory_reflex/reflex_tree.txt", out)
	return
}

// получить ветку дерева в виде строки
func getReflexNode(wt *ReflexNode) string {
	var out = ""

	out += strconv.Itoa(wt.ID) + "|"
	out += strconv.Itoa(wt.ParentID) + "|"
	out += strconv.Itoa(wt.baseID) + "|"
	out += strconv.Itoa(wt.StyleID) + "|"
	out += strconv.Itoa(wt.ActionID) + "|"
	out += strconv.Itoa(wt.GeneticReflexID) + "|"
	out += strconv.Itoa(wt.ConditionedReflex)
	out += "\r\n"

	if wt.Children == nil { return out } // конец
	for n := 0; n < len(wt.Children); n++ {
		out += getReflexNode(&wt.Children[n])
	}
	return out
}

var detectedLastNodID = 0 // текущий последний распознанный узел дерева - результат распознавания

func ConditionsDetection(condArr []int) {
	detectedLastNodID = 0 /* распознавание условий в дереве рефлексов, нахождение ветки с данными условиями данного рефлекса
	condArr получать с помощью func getConditionsArr(lev1ID int, lev2 []int, lev3 []int, PhraseID []int,ToneID int,MoodID int )([3]int){
	*/
	// основа дерева
	cnt := len(ReflexTree.Children)
	for n := 0; n < cnt; n++ {
		node := ReflexTree.Children[n]
		lev1 := node.baseID
		if condArr[0] == lev1 {
			detectedLastNodID = node.ID
			ost := condArr[1:]
			getReflexTreeNode(1, ost, &node)
			break // только один из Базовых состояний
		}
	}
	return
}

// получить последний узел ветки - результат распознавания
func getReflexTreeNode(level int, cond []int, node *ReflexNode) {
	if len(cond) == 0 { return }

	ost := cond[1:]
	for n := 0; n < len(node.Children); n++ {
		cld := node.Children[n]
		var levID = 0
		switch level {
		case 1: levID = cld.StyleID
		case 2: levID = cld.ActionID
		}
		// cond[0] потому, что на следующем уровне cond уже подрезана
		if cond[0] != levID {	return } // пошло не туда

		detectedLastNodID = node.ID
		level++
		getReflexTreeNode(level, ost, &cld)
	}
	return
}

// создание новой ветки с новым рефлексом типа GeneticReflex, начиная с заданного узла
func createNewReflexToTreeFromNodes(level int,cond []int,node *ReflexNode) int {
	if node == nil { return 0	}
	if level >= len(cond) {	return node.ID }
	var id=0
	/* switch(level){
	case 0:
		id,_=createNewReflexNode(node,0,cond[0],0,0,0,0)
	case 1:
		id,_=createNewReflexNode(node,0,cond[0],cond[1],0,0,0)
	case 2:
		id,_=createNewReflexNode(node,0,cond[0],cond[1],cond[2],0,0)
	}
	*/
	id, _ = createNewReflexNode(node,0, cond[0], cond[1], cond[2],0,0,true)
	level++
	id = createNewReflexToTreeFromNodes(level, cond, ReflexTreeFromID[id])
	return id
}

/* Сразу создать и добавить безусловный рефлекс в дерево, если таких узлов еще нет.
Формат записи безусловного рефлекса: ID|baseID|styleID...|actionID...
Если у рефлекса пропущены условия, то этот рефлекс нужно привязать ко всем узлам пропущенного уровня.
 */
func addGeneticReflexesToTree(detectedActiveLastNodID int, condArr []int) {
	notAllowScanInReflexesThisTime = true // запрет показа карты при обновлении

	// найти ID GeneticReflexes (список всех dnk_reflexes.txt) по условиям
	reflexID := findGeneticReflexFromCondinion(strconv.Itoa(condArr[0]), condArr[1], condArr[2])
	if reflexID > 0 {
		// v := GeneticReflexes[reflexID]
		// trigger:=v.ActionIDarr
		level := getLevelFromNodeID(detectedActiveLastNodID)
		lastNodeID := formingBranch(reflexID, detectedActiveLastNodID, level, condArr)
		detectedActiveLastNodID = lastNodeID
		if ReflexTreeFromID[detectedActiveLastNodID].GeneticReflexID > 0 {
			if condArr[2] == 0 { // древний рефлекс
				oldReflexesIdArr = append(oldReflexesIdArr, ReflexTreeFromID[detectedActiveLastNodID].GeneticReflexID)
			} else { // нормальный безусловный рефлекс (с пусковым стимулом)
				geneticReflexesIdArr = append(geneticReflexesIdArr, ReflexTreeFromID[detectedActiveLastNodID].GeneticReflexID)
			}
		}
	}
	// сохранение  - учти, что срабатывает - после пятого пульса
	SaveReflexesAttributes()
	notAllowScanInReflexesThisTime = false
}

// найти ID GeneticReflexes (список всех dnk_reflexes.txt) по условиям
func findGeneticReflexFromCondinion(basic string,img1id int,img2id int) int {
	img1 := BaseStyleArr[img1id]
	var img2 *TriggerStimuls

	if img2id > 0 {
		img2 = TriggerStimulsArr[img2id]
	}
	lev1str := ""
	for i := 0; i < len(img1.BSarr); i++ {
		if len(lev1str) > 0 { lev1str += "," }
		lev1str += strconv.Itoa(img1.BSarr[i])
	}
	lev2str := ""
	if img2 != nil {
		for i := 0; i < len(img2.RSarr); i++ {
			if len(lev2str) > 0 { lev2str += "," }
			lev2str += strconv.Itoa(img2.RSarr[i])
		}
	}

	for id, v:= range geneticReflexesStr {
		if v.lev1 == basic && v.lev2 == lev1str && v.lev3 == lev2str { return id }
	}
	return 0
}

// нарастить ветку недостающим узлом
func formingBranch(reflexID int,fromID int,lastLevel int,condArr []int) int {
	lastNode := ReflexTreeFromID[fromID]
	lastNodeID := createNewReflexToTreeFromNodes(lastLevel, condArr, lastNode)
	// родителем должен быть последний найденный узел, а не тот, что будет создан первым
	// !!! НЕТ !!! ReflexTreeFromID[lastNodeID].ParentID=lastNode.ID
	// привязать рефлекс
	ReflexTreeFromID[lastNodeID].GeneticReflexID = reflexID

	return lastNodeID
}

// найти уровень вложения данного узла в ветке
func getLevelFromNodeID(nodeID int) int {
	lastNode := ReflexTreeFromID[nodeID]
	var level = 0
	for lastNode.ParentNode != nil {
		level++
		lastNode = lastNode.ParentNode
	}
	return level
}

// найти КОНЕЧНЫЙ узел по условиям
func FindReflexTreeNodeFromCondition(baseID int, StyleID int, ActionID int)(int, *ReflexNode) {
	for k, v := range ReflexTreeFromID {
		if v.baseID == baseID && v.StyleID == StyleID && v.ActionID == ActionID { return k, v	}
	}
	return 0, nil
}

// сохранение при выходе reflexes.SaveReflexesAttributes() и когда нужно
// !!! но только после того, как все данные будут загружены:
func SaveReflexesAttributes() {
	if ReflexPulsCount < 5 { return	}
	// сохранить образы восприятия и пусковых стимулов после прохода всех безусловных рефлексов
	SaveBaseStyleArr()
	SaveTriggerStimulsArr()
	saveReflexTree()
	SaveConditionReflex()
}