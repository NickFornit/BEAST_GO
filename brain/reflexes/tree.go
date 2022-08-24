/*
дерево рефлексов, безусловных и условных

Узлы дерева создаются только рефлексами:
безусловными сразу при загрузке,
условными по мере возникновения новых - в новых уловиях.
У дерева - только 3 уровня для безусловных рефлексов, все узла которых - в виде ID:
1 - базовое состояние - ID Плохо, Норма, хорошо
2 - сочетаний активных Базовых контекстов - ID BaseStyleArr
3 - сочетаний пусковых стимулов - ID TriggerStimulsArr
При возникновении условных рефлексов просто добавляется новый узел
- образ новых условий (2 или 3-го уровней), запускающих усл.рефлекс.

Формат записи безусловного рефлекса: ID|baseID|styleID...|actionID...
*/

package reflexes

import (
	"BOT/lib"
	"sort"
	"strconv"
	"strings"
)

//////////////////////////////////////
func initReflexTree(){ // после инициализации loadGeneticReflexes()
	loadReflexTree()
	addGeneticReflexesToTree()
	//SaveReflexesAttributes()

	readyForRecognitionRflexes() // ини для дерева распознавания рефлексов


}
// ID|parentID|baseID|styleID|actionID|GeneticReflexID|ConditionedReflex
type ReflexNode struct { // узел дерева рефлексов
	ID int
	baseID int // базовое состояние
	StyleID int // стиль поведения - сочетание активностей Базовых контекстов  - ID BaseStyleArr
	ActionID int // сочетание пусковых стимулов  - ID TriggerStimulsArr

	GeneticReflexID int // безусловный рефлекс
	ConditionedReflex int // условный рефлекс, если есть, блокирует GeneticReflexID узла

	Children []ReflexNode // дочерние узлы (ветвление) НЕ АДРЕСА, А РЕАЛЬНЫЕ ОБЪЕКТЫ
	ParentID int     // ID родителя
	ParentNode *ReflexNode  // адрес родителя
}

var ReflexTree ReflexNode
var ReflexTreeFromID=make(map[int]*ReflexNode)

/*запрет показа карты при обновлении против паники типа "одновременная запись и считывание карты"
Использовать для всех операций записи узлов дерева
 */
var notAllowScanInReflexesThisTime=false

///////////////////////////////////////
var lastReflexNodeID=0
func createNewReflexNode(parent *ReflexNode,id int,baseID int,StyleID int,
				ActionID int,GeneticReflexID int,ConditionedReflex int)(int,*ReflexNode){
// если есть такой узел, то не создавать
	idOld,nodeOld:=FindReflexTreeNodeFromCondition(baseID,StyleID,ActionID)
	if idOld>0{
		return idOld,nodeOld
	}

	if id==0{
		lastReflexNodeID++
		id=lastReflexNodeID
	}else{
		//		newW.ID=id
		if lastReflexNodeID<id{
			lastReflexNodeID=id
		}
	}

	var node ReflexNode
	node.ID = id
	node.ParentNode=parent
	node.ParentID=parent.ID
	node.baseID=baseID
	node.StyleID=StyleID
	node.ActionID=ActionID
	node.GeneticReflexID=GeneticReflexID
	node.ConditionedReflex=ConditionedReflex

	parent.Children = append(parent.Children, node)
	// четко находим новый вставленный член (а не &parent.Children[count-1])
	var newN *ReflexNode
	for i := 0; i < len(parent.Children); i++ {
		if parent.Children[i].ID == node.ID {
			newN = &parent.Children[i]
		}
	}
// т.к. append меняет длину массива, перетусовывая адреса, то нужно обновить адреса в ReflexTreeFromID:
	updateReflexTreeFromID(parent)// здесь потому, что при загрузке из файла нужно на лету получать адреса

	return id,newN
}
// корректируем адреса всех узлов
func updateReflexTreeFromID(parent *ReflexNode){
	//updatingPhraseTreeFromID(&VernikePhraseTree)
	updatingPhraseTreeFromID(parent)
}
// проход всего дерева
func updatingPhraseTreeFromID(rt *ReflexNode){
	if rt.ID>0 {
		rt.ParentNode=ReflexTreeFromID[rt.ParentID] // wr.ParentNode адрес меняется из=за corretsParent(,
		ReflexTreeFromID[rt.ID] = rt
	}
	if rt.Children == nil{// конец ветки
		return
	}
	for i := 0; i < len(rt.Children); i++ {
		updatingPhraseTreeFromID(&rt.Children[i])
	}
}
///////////////////////////////////////////////////////////
// загрузить записанное дерево
// ID|parentID|baseID|styleID|actionID|geneticReflexID|conditionedReflex|
func loadReflexTree(){
	createNulLevelReflexTree(&ReflexTree)
	strArr,_:=lib.ReadLines(lib.GetMainPathExeFile()+"/memory_reflex/reflex_tree.txt")
	cunt:=len(strArr)
	//просто проход по всем строкам файла подряд так что сначала идут дочки, потом - их родители
	for n := 0; n < cunt; n++ {
		if len(strArr[n])==0{
			continue
		}
		if len(strArr[n])<2{
			panic("Сбой загрузки дерева рефлексов: ["+strconv.Itoa(n) + "] " + strArr[n])
			return
		}
		p:=strings.Split(strArr[n], "|")
		id,_:=strconv.Atoi(p[0])
		parentID,_:=strconv.Atoi(p[1])
		baseID,_:=strconv.Atoi(p[2])
		styleID,_:=strconv.Atoi(p[3])
		actionID,_:=strconv.Atoi(p[4])
		geneticReflexID,_:=strconv.Atoi(p[5])
		conditionedReflex,_:=strconv.Atoi(p[6])
		// новый узел с каждой строкой из файла
		createNewReflexNode(ReflexTreeFromID[parentID],id,baseID,styleID,
						actionID,geneticReflexID,conditionedReflex)
	}
	return
}
// создать первый, нулевой уровень дерева
func createNulLevelReflexTree(rt *ReflexNode){
	rt.ID=0
	ReflexTreeFromID[rt.ID]=rt
	return
}
/////////////////////////////////////
func saveReflexTree(){
	var out=""
	cnt:=len(ReflexTree.Children)
	for n := 0; n < cnt; n++ {
		out+=getReflexNode(&ReflexTree.Children[n])
	}
	lib.WriteFileContent(lib.GetMainPathExeFile()+"/memory_reflex/reflex_tree.txt",out)
	return
}
func getReflexNode(wt *ReflexNode)(string){
	var out=""
	//	if wt.ParentID>0 {
	out += strconv.Itoa(wt.ID) + "|"
	out += strconv.Itoa(wt.ParentID) + "|"
	out += strconv.Itoa(wt.baseID) + "|"
	out += strconv.Itoa(wt.StyleID) + "|"
	out += strconv.Itoa(wt.ActionID) + "|"
	out += strconv.Itoa(wt.GeneticReflexID) + "|"
	out += strconv.Itoa(wt.ConditionedReflex)
	out +="\r\n"
	//	}
	if(wt.Children==nil){// конец
		return out
	}
	for n := 0; n < len(wt.Children); n++ {
		out+=getReflexNode(&wt.Children[n])
	}
	return out
}
/////////////////////////////////////







////////////////////////////////////////
/*  распознавание условий в дереве рефлексов, нахождение ветки с данными условиями данного рефлекса
condArr получать с помощью func getConditionsArr(lev1ID int, lev2 []int, lev3 []int, PhraseID []int,ToneID int,MoodID int )([3]int){
 */
var detectedLastNodID=0// текущий последний распознанный узел дерева - результат распознавания
func ConditionsDetection(condArr []int){
	detectedLastNodID=0
	// основа дерева
	cnt := len(ReflexTree.Children)
	for n := 0; n < cnt; n++ {
		node := ReflexTree.Children[n]
		lev1 := node.baseID
		if condArr[0] == lev1 {
			detectedLastNodID=node.ID
			ost:=condArr[1:]
			getReflexTreeNode(1,ost, &node)
			break // только один из Базовых состояний
		}
	}
	// результат распознавания обрабатывается в func addGeneticReflexesToTree()
	return
}
/////////////////////
func getReflexTreeNode(level int,cond []int,node *ReflexNode){
	if len(cond)==0{
		return
	}
	ost:=cond[1:]

	for n := 0; n < len(node.Children); n++ {
		cld := node.Children[n]
	var levID=0
	switch level{
	case 1: levID=cld.StyleID
	case 2: levID=cld.ActionID
	}
//cond[0] потому, что на следующем уровне cond уже подрезана
	if cond[0] != levID {// пошло не туда
		return
	}

	detectedLastNodID=node.ID
		level++
		getReflexTreeNode(level,ost, &cld)
	}
	return
}
/////////////////////////////////////////////////
// создание новой ветки с новым рефлексом типа GeneticReflex, начиная с заданного узла
func createNewReflexToTreeFromNodes(level int,cond []int,node *ReflexNode)(int){
	if node==nil {
		return 0
	}
	if level>=len(cond) {
		return node.ID
	}
var id=0

/*
	switch(level){
	case 0:
		id,_=createNewReflexNode(node,0,cond[0],0,0,0,0)
	case 1:
		id,_=createNewReflexNode(node,0,cond[0],cond[1],0,0,0)
	case 2:
		id,_=createNewReflexNode(node,0,cond[0],cond[1],cond[2],0,0)
	}
	*/
	id,_=createNewReflexNode(node,0,cond[0],cond[1],cond[2],0,0)

	level++
	 id=createNewReflexToTreeFromNodes(level,cond, ReflexTreeFromID[id])
return id
}
/////////////////////////////////////


/////////////////////////////////////////////////////////////////////
/*  создать и добавить безусловные рефлексы в дерево, если таких узлов еще нет
т.к. безусловные рефлексы уже прописаны заранее, то нужно их всех прогнать для вставки в дерево
Условные рефлексы будут добавляться по мере возникновения.
Формат записи безусловного рефлекса: ID|baseID|styleID...|actionID...
Если у рефлекса пропущены условия, то этот рефлекс нужно привязать ко всем узлам пропущенного уровня.
 */
func addGeneticReflexesToTree(){
	notAllowScanInReflexesThisTime=true // запрет показа карты при обновлении
	keys := make([]int, 0, len(GeneticReflexes))
	for k := range GeneticReflexes {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, id := range keys {
		v:=GeneticReflexes[id]
// вытащить 3 уровня условий в виде ID их образов
		condArr:=getConditionsArr(v.lev1, v.lev2, v.lev3, nil,0,0)
		if v.lev3==nil{// если в рефлексе нет пускового стимула
			// иначе прописывает образ ID=2
			condArr[2]=0
		}


// поиск в дереве такого сочетания условий
		ConditionsDetection(condArr)

		if id==48{
			id=48
		}

// результат поиска:
		if detectedLastNodID > 0 {
			lastNode:=ReflexTreeFromID[detectedLastNodID]
			if lastNode!=nil {
				//насколько найденное соотвествует condArr?
					if lastNode.StyleID==0{
						// если уже есть такой узел, то ничего не делать с ним
						idOld,_:=FindReflexTreeNodeFromCondition(condArr[0], condArr[1], condArr[2])
						if idOld==0 {
							formingBranch(id, detectedLastNodID, 1, condArr)
						}
					}else{
						if condArr[2]>0 && lastNode.ActionID==0{
							// если уже есть такой узел, то ничего не делать с ним
							idOld,_:=FindReflexTreeNodeFromCondition(condArr[0], condArr[1], condArr[2])
							if idOld==0 {
								formingBranch(id, detectedLastNodID, 2, condArr)
							}
						}
					}
			}
		}else{// вообще нет, нарастить все с нуля
			formingBranch(id,ReflexTree.ID,0,condArr)
		}
	}
	// сохранение
	SaveReflexesAttributes()
//	SaveReflexesAttributes()

	notAllowScanInReflexesThisTime=false
}
/////////////////////////////////////



func formingBranch(reflexID int,fromID int,lastLevel int,condArr []int){
	// нарастить ветку недостающим
	lastNode:=ReflexTreeFromID[fromID]

	lastNodeID:=createNewReflexToTreeFromNodes(lastLevel,condArr,lastNode)
	// родителем должен быть последний найденный узел, а не тот, что будет создан первым
	//!!! НЕТ !!! ReflexTreeFromID[lastNodeID].ParentID=lastNode.ID
	//привязать рефлекс
	ReflexTreeFromID[lastNodeID].GeneticReflexID=reflexID
}




// найти КОНЕЧНЫЙ узел по условиям
func FindReflexTreeNodeFromCondition(baseID int,StyleID int,ActionID int)(int,*ReflexNode){
	for k, v := range ReflexTreeFromID {
		if v.baseID==baseID && v.StyleID==StyleID && v.ActionID==ActionID{
			return k, v
		}
	}
	return 0,nil
}
//////////////////////////////////////





// сохранение при выходе reflexes.SaveReflexesAttributes() и когда нужно
// !!! но только после того, как все данные будут загружены:
func SaveReflexesAttributes(){
	if ReflexPulsCount <5{
		return
	}
	// сохранить образы восприятия и пусковых стимулов после прохода всех безусловных рефлексов
	SaveBaseStyleArr()
	saveTriggerStimulsArr()

	saveReflexTree()
	SaveConditionReflex()
}