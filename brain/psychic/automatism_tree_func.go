/*  функции Дерева автоматизмов


*/

package psychic

import (
	"BOT/lib"
	"strconv"
	"strings"
)

/////////////////////////////////////////////////////////////


/* Создать новый узел дерева автоматизма.
Формат записи:
ID|ParentNode|BaseID|EmotionID|ActivityID|ToneMoodID|SimbolID|VerbalID

 */
var lastAutomatizmNodeID=0
var noRunThisOperation=false // не проверять на дубли
func createNewAutomatizmNode(parent *AutomatizmNode,id int,baseID int,EmotionID int,
	ActivityID int,ToneMoodID int,SimbolID int,VerbalID int,CheckUnicum bool)(int,*AutomatizmNode){

	if parent == nil{
		return 0,nil
	}

	//if !noRunThisOperation { НЕЛЬЗЯ ИГНОРИРОВАТЬ ИНАЧЕ СОЗДАЕТ ЛИШНЕЕ
	// если есть такой узел, то не создавать
	if CheckUnicum {
		idOld, nodeOld := FindAutomatizmTreeNodeFromCondition(baseID, EmotionID, ActivityID, ToneMoodID, SimbolID, VerbalID)
		if idOld > 0 {
			return idOld, nodeOld
		}
	}

	if id==0{
		lastAutomatizmNodeID++
		id=lastAutomatizmNodeID
	}else{
		//		newW.ID=id
		if lastAutomatizmNodeID<id{
			lastAutomatizmNodeID=id
		}
	}

	var node AutomatizmNode
	node.ID = id
	node.ParentNode=parent
	node.ParentID=parent.ID
	node.BaseID=baseID
	node.EmotionID=EmotionID
	node.ActivityID=ActivityID
	if ToneMoodID==0{
		//  !!!!! ToneMoodID=90
	}
	node.ToneMoodID=ToneMoodID
	node.SimbolID=SimbolID
	node.VerbalID=VerbalID

	parent.Children = append(parent.Children, node)
	// четко находим новый вставленный член (а не &parent.Children[count-1])
	var newN *AutomatizmNode
	for i := 0; i < len(parent.Children); i++ {
		if parent.Children[i].ID == node.ID {
			newN = &parent.Children[i]
		}
	}
	AutomatizmTreeFromID[node.ID]=&node
	// т.к. append меняет длину массива, перетусовывая адреса, то нужно обновить адреса в AutomatizmTreeFromID:
	updateAutomatizmTreeFromID(parent)// здесь потому, что при загрузке из файла нужно на лету получать адреса

	return id,newN
}
// корректируем адреса всех узлов
func updateAutomatizmTreeFromID(parent *AutomatizmNode){
	//updatingPhraseTreeFromID(&VernikePhraseTree)
	updatingPhraseTreeFromID(parent)
}
// проход всего дерева
func updatingPhraseTreeFromID(rt *AutomatizmNode){
	if rt.ID>0 {
		rt.ParentNode=AutomatizmTreeFromID[rt.ParentID] // wr.ParentNode адрес меняется из=за corretsParent(,
		AutomatizmTreeFromID[rt.ID] = rt
	}
	if rt.Children == nil{// конец ветки
		return
	}
	for i := 0; i < len(rt.Children); i++ {
		updatingPhraseTreeFromID(&rt.Children[i])
	}
}
///////////////////////////////////////////////////////////
/* загрузить записанное дерево
Формат записи:
ID|ParentNode|BaseID|EmotionID|ActivityID|ToneMoodID|SimbolID|VerbalID
 */
func loadAutomatizmTree(){
	//нулевой узел
	AutomatizmTreeFromID[0]=&AutomatizmTree // все по нулям по умолчанию

//	createNulLevelAutomatizmTree(&AutomatizmTree)
	strArr,_:=lib.ReadLines(lib.GetMainPathExeFile()+"/memory_psy/automatizm_tree.txt")
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
		/*
		if AutomatizmTreeFromID[parentID]==nil{
			// сбой загрузки рефлексов
			return
		}
		*/
		baseID,_:=strconv.Atoi(p[2])
		EmotionID,_:=strconv.Atoi(p[3])
		ActivityID,_:=strconv.Atoi(p[4])
		ToneMoodID,_:=strconv.Atoi(p[5])
		SimbolID,_:=strconv.Atoi(p[6])
		VerbalID,_:=strconv.Atoi(p[7])
		// новый узел с каждой строкой из файла
var saveDoWritingFile= doWritingFile; doWritingFile =false
		createNewAutomatizmNode(AutomatizmTreeFromID[parentID],id,baseID,EmotionID,
			ActivityID,ToneMoodID,SimbolID,VerbalID,false)
doWritingFile =saveDoWritingFile
	}
	return
}
/*
// создать первый, нулевой уровень дерева
func createNulLevelAutomatizmTree(rt *AutomatizmNode){
	rt.ID=0
	AutomatizmTreeFromID[rt.ID]=rt
	return
}
 */
/////////////////////////////////////
func SaveAutomatizmTree(){
	notAllowScanInTreeThisTime = true
	var out=""
	cnt:=len(AutomatizmTree.Children)
	for n := 0; n < cnt; n++ { // чтобы записывалось по порядку родителей
		out+=getAutomatizmNode(&AutomatizmTree.Children[n])
	}
	lib.WriteFileContent(lib.GetMainPathExeFile()+"/memory_psy/automatizm_tree.txt",out)
	notAllowScanInTreeThisTime = false
	return
}
func getAutomatizmNode(wt *AutomatizmNode)(string){
	var out=""
	//	if wt.ParentID>0 {
	out += strconv.Itoa(wt.ID) + "|"
	out += strconv.Itoa(wt.ParentID) + "|"
	out += strconv.Itoa(wt.BaseID) + "|"
	out += strconv.Itoa(wt.EmotionID) + "|"
	out += strconv.Itoa(wt.ActivityID) + "|"
	out += strconv.Itoa(wt.ToneMoodID) + "|"
	out += strconv.Itoa(wt.SimbolID) + "|"
	out += strconv.Itoa(wt.VerbalID)
	out +="\r\n"
	//	}
	if wt.Children==nil{// конец
		return out
	}
	for n := 0; n < len(wt.Children); n++ {
		out+=getAutomatizmNode(&wt.Children[n])
	}
	return out
}
/////////////////////////////////////
// найти КОНЕЧНЫЙ узел по условиям
func FindAutomatizmTreeNodeFromCondition(baseID int,EmotionID int,
	ActivityID int,ToneMoodID int,SimbolID int,VerbalID int)(int,*AutomatizmNode){
	/*
//	 поиск по дереву, что и эффективнее.
	for k, v := range AutomatizmTreeFromID {
		if v.BaseID==baseID && v.EmotionID==EmotionID &&
			v.ActivityID==ActivityID && ToneMoodID==v.ToneMoodID && v.SimbolID==SimbolID && v.VerbalID==VerbalID{
			return k, v
		}
	}
	return 0,nil
	*/
	var id=0
	var aut *AutomatizmNode
	cnt:=len(AutomatizmTree.Children)
	for n := 0; n < cnt; n++ {
		id,aut=checkAutomatizmTree(&AutomatizmTree.Children[n],baseID, EmotionID, ActivityID, ToneMoodID, SimbolID, VerbalID)
		if id >0{
			return id,aut
		}
	}
	return 0,nil
}
////////////////
func checkAutomatizmTree(v *AutomatizmNode,baseID int,EmotionID int,
ActivityID int,ToneMoodID int,SimbolID int,VerbalID int)(int,*AutomatizmNode){
	var id=v.ID
	var aut =v

	if v.ID==13{
	//	v.ID=13
	}
// как только наткнется в предыдущих на такое услове - выдаст ID этой ветки
	if v.BaseID==baseID && v.EmotionID==EmotionID &&
		v.ActivityID==ActivityID && ToneMoodID==v.ToneMoodID && v.SimbolID==SimbolID && v.VerbalID==VerbalID {
		return v.ID,v
	}

	if v.Children==nil {// конец
		return 0,nil
	}
	for n := 0; n < len(v.Children); n++ {
		id,aut=checkAutomatizmTree(&v.Children[n],baseID, EmotionID, ActivityID, ToneMoodID, SimbolID, VerbalID)
		if id>0{
			return id,aut
		}
	}
	return 0,nil //v.ID

}
//////////////////////////////////////




//////////////////////////////////////////////////////////////
/* создание иерархии АКТИВНЫХ образов контекстов условий и пусковых стимулов в виде ID образов в [5]int
создать последовательность уровней условий в виде массива  ID последовательности ID уровней
*/
func getActiveConditionsArr(lev1 int, lev2 int, lev3 int, lev4 int, lev5 int, lev6 int)([]int){
	arr:=make([]int, 6)
	arr[0]=lev1
	arr[1]=lev2
	arr[2]=lev3
	arr[3]=lev4
	arr[4]=lev5
	arr[5]=lev6
	return arr
}
func getConditionsCount(condArr []int)(int){
	var count=0
	for i := 0; i < len(condArr); i++ {
		if condArr[i]>0{
			count++
		}
	}
	return count
}
////////////////////////////////////////////////////


// выдать массив узлов ветки по заданному ID узла
func getBrangeNodeArr(lastNodeId int)([]*AutomatizmNode){
	var nArr []*AutomatizmNode
	node:=AutomatizmTreeFromID[lastNodeId]
	if node==nil { return nil}
	for {
		if node==nil {
			break
		}
		nArr = append(nArr, node)
		node=node.ParentNode
	}
	return nArr
}
//////////////////////////////////////


//////////////////////////////////////////////
/* поиск узла дерева автоматизмов по условиям у.рефлекса для automatizms_from_reflexes.go
Если нет такого узла - дорастить ветку.
Выдать  ID узла

 */
func FindConditionsNode(lev1 int,lev2 []int,actArr []int, tonMood int,fistSymb int,verbalID int)(int){
	// образ эмоции
	eID,_:=createNewBaseStyle(0,lev2,true)
	// образ действий дерева: из TriggerStimuls -> Activity
	aID,_:=createNewlastActivityID(0,actArr,true)// конвертировать образ типа reflexes.TriggerStimuls в psychic.Activity

	// проход дераева автоматизмов:
	detectedActiveLastNodID=0
	ActiveBranchNodeArr=nil
	CurrentAutomatizTreeEnd=nil
	currentStepCount=0
	var lastNodeID=0
	condArr:=getActiveConditionsArr(lev1,eID,aID,tonMood,fistSymb,verbalID)
	if verbalID==15{
	//	verbalID=15
	}
	// основа дерева
	cnt := len(AutomatizmTree.Children)
	for n := 0; n < cnt; n++ {
		node := AutomatizmTree.Children[n]
		lev1 := node.BaseID
		if condArr[0] == lev1 {
			detectedActiveLastNodID=node.ID
			ost:=condArr[1:]
			if len(ost)==0{

			}
			//conditionAutomatizmFound2(1,ost, &node)
			conditionAutomatizmFound(1,ost, &node)
			break // другие ветки не смотреть
		}
	}

	// результат активации Дерева:
	if detectedActiveLastNodID>0{
		// есть ли еще неучтенные, нулевые условия? т.е. просто показаь число ненулевых значений condArr
		conditionsCount:=getConditionsCount(condArr)
		if currentStepCount<conditionsCount { // не пройдено до конца имеющихся условий
			noRunThisOperation=true // не проверять на дубли, раз уже проверено
			lastNodeID = formingBranch(detectedActiveLastNodID, currentStepCount+1, condArr)
			noRunThisOperation=false
		}else{// все condArr пойдены, ветка существует,
			// НО если последнее условие не совпадает, нужно создать ветку на основе предпоследнего узла
			//lenNcount:=len(ActiveBranchNodeArr)
			//if lenNcount<8{// остановилось на предпоследнем узле detectedActiveLastNodID

			//}

return detectedActiveLastNodID
		}
	}

/*
	// проверка
	n:=AutomatizmTreeFromID[lastNodeID];
	if n.VerbalID == 0{
		return 0
	}
*/
	return lastNodeID
}
//.......
/*
func conditionAutomatizmFound2(level int,cond []int,node *AutomatizmNode){
	if cond==nil || len(cond)==0{
		return
	}

	ost:=cond[1:]

	if node.ID==6{
		node.ID=6
	}

	if level==4 && cond[0]==15{
		cond[0]=15
	}

	for n := 0; n < len(node.Children); n++ {
		cld := node.Children[n]
		var val = 0
		switch level {
		case 0:
			val = cld.BaseID
		case 1:
			val = cld.EmotionID
		case 2:
			val = cld.ActivityID
		case 3:
			val = cld.ToneMoodID
		case 4:
			val = cld.SimbolID
		case 5:
			val = cld.VerbalID
		}
		if cond[0] == val {
			detectedActiveLastNodID=cld.ID
			ActiveBranchNodeArr=append(ActiveBranchNodeArr,cld.ID)
		}else {
			currentStepCount=level-1
			return
			if level == 5 { // последний не совпадает


			}
		}

		level++
		currentStepCount=level
		conditionAutomatizmFound2(level,ost, &node.Children[n])
		return // раз совпало, то другие ветки не смотреть
	}

	return
}
*/
////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////
// создание ветки, начиная с заданного узла fromID
func formingBranch(fromID int,lastLevel int,condArr []int)(int){
	// нарастить ветку недостающим
	lastNode:=AutomatizmTreeFromID[fromID]

	lastNodeID:=addNewBranchFromNodes(lastLevel,condArr,lastNode)
	if lastNodeID>0{
		if !noRunThisOperation {
			if doWritingFile {SaveAutomatizmTree() }
		}
	}
	return lastNodeID
}
/////////////////////////////////////////////////////
// создание новой ветки с новым автоматизмом, начиная с заданного узла при проходе дерева
func addNewBranchFromNodes(level int,cond []int,node *AutomatizmNode)(int){
	if node==nil {
		return 0
	}
	if level>=len(cond) {
		return node.ID
	}
	var id=0
	switch(level){
	case 0:
		id,_=createNewAutomatizmNode(node,0,cond[0],0,0,0,0,0,true)
	case 1:
		id,_=createNewAutomatizmNode(node,0,cond[0],cond[1],0,0,0,0,true)
	case 2:
		id,_=createNewAutomatizmNode(node,0,cond[0],cond[1],cond[2],0,0,0,true)
	case 3:
		id,_=createNewAutomatizmNode(node,0,cond[0],cond[1],cond[2],cond[3],0,0,true)
	//case 4:
	//	id,_=createNewAutomatizmNode(node,0,cond[0],cond[1],cond[2],cond[3],cond[4],0)
	case 4, 5:
		// нельзя отдельно указывать символ фразы и фразу - они всегда в паре
		id,_=createNewAutomatizmNode(node,0,cond[0],cond[1],cond[2],cond[3],cond[4],cond[5],true)
	}
	level++
	id=addNewBranchFromNodes(level,cond, AutomatizmTreeFromID[id])
	return id
}
/////////////////////////////////////


