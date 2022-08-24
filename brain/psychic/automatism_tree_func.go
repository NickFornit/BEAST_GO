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
func createNewAutomatizmNode(parent *AutomatizmNode,id int,baseID int,EmotionID int,
	ActivityID int,ToneMoodID int,SimbolID int,VerbalID int)(int,*AutomatizmNode){
	// если есть такой узел, то не создавать
	idOld,nodeOld:=FindAutomatizmTreeNodeFromCondition(baseID,EmotionID,ActivityID,ToneMoodID,SimbolID,VerbalID)
	if idOld>0{
		return idOld,nodeOld
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
	createNulLevelAutomatizmTree(&AutomatizmTree)
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
		baseID,_:=strconv.Atoi(p[2])
		EmotionID,_:=strconv.Atoi(p[3])
		ActivityID,_:=strconv.Atoi(p[4])
		ToneMoodID,_:=strconv.Atoi(p[5])
		SimbolID,_:=strconv.Atoi(p[6])
		VerbalID,_:=strconv.Atoi(p[7])
		// новый узел с каждой строкой из файла
		createNewAutomatizmNode(AutomatizmTreeFromID[parentID],id,baseID,EmotionID,
			ActivityID,ToneMoodID,SimbolID,VerbalID)
	}
	return
}
// создать первый, нулевой уровень дерева
func createNulLevelAutomatizmTree(rt *AutomatizmNode){
	rt.ID=0
	AutomatizmTreeFromID[rt.ID]=rt
	return
}
/////////////////////////////////////
func SaveAutomatizmTree(){
	var out=""
	cnt:=len(AutomatizmTree.Children)
	for n := 0; n < cnt; n++ {
		out+=getAutomatizmNode(&AutomatizmTree.Children[n])
	}
	lib.WriteFileContent(lib.GetMainPathExeFile()+"/memory_psy/automatizm_tree.txt",out)
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
	if(wt.Children==nil){// конец
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
	for k, v := range AutomatizmTreeFromID {
		if v.BaseID==baseID && v.EmotionID==EmotionID &&
			v.ActivityID==ActivityID && ToneMoodID==v.ToneMoodID && v.SimbolID==SimbolID && v.VerbalID==VerbalID{
			return k, v
		}
	}
	return 0,nil
}
//////////////////////////////////////
// создание новой ветки с новым автоматизмом, начиная с заданного узла
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
		id,_=createNewAutomatizmNode(node,0,cond[0],0,0,0,0,0)
	case 1:
		id,_=createNewAutomatizmNode(node,0,cond[0],cond[1],0,0,0,0)
	case 2:
		id,_=createNewAutomatizmNode(node,0,cond[0],cond[1],cond[2],0,0,0)
	case 3:
		id,_=createNewAutomatizmNode(node,0,cond[0],cond[1],cond[2],cond[3],0,0)
	case 4:
		id,_=createNewAutomatizmNode(node,0,cond[0],cond[1],cond[2],cond[3],cond[4],0)
	case 5:
		id,_=createNewAutomatizmNode(node,0,cond[0],cond[1],cond[2],cond[3],cond[4],cond[5])
	}
	level++
	id=addNewBranchFromNodes(level,cond, AutomatizmTreeFromID[id])
	return id
}
/////////////////////////////////////


// создание ветки, начиная с заданного узла fromID
func formingBranch(fromID int,lastLevel int,condArr []int)(int){
	// нарастить ветку недостающим
	lastNode:=AutomatizmTreeFromID[fromID]

	lastNodeID:=addNewBranchFromNodes(lastLevel,condArr,lastNode)
	if lastNodeID>0{
		SaveAutomatizmTree()
	}
	return lastNodeID
}
/////////////////////////////////////////////////////



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



