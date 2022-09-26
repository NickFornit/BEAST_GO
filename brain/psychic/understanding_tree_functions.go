/* Функции для Дерева понимания (дерева ментальных автоматизмов)

*/

package psychic

import (
	"BOT/lib"
	"strconv"
	"strings"
)

///////////////////////////////

////////////////////////////////




/* Создать новый узел дерева понимания.
Формат записи:
ID|ParentNode|Mood|EmotionID|SituationID|PurposeID
*/
var lastUnderstandingNodeID=0
func createNewUnderstandingNode(parent *UnderstandingNode,id int,Mood int,EmotionID int,
	SituationID int,PurposeID int)(int,*UnderstandingNode){
	// если есть такой узел, то не создавать
	idOld,nodeOld:=FindUnderstandingTreeNodeFromCondition(Mood,EmotionID,SituationID,PurposeID)
	if idOld>0{
		return idOld,nodeOld
	}

	if id==0{
		lastUnderstandingNodeID++
		id=lastUnderstandingNodeID
	}else{
		//		newW.ID=id
		if lastUnderstandingNodeID<id{
			lastUnderstandingNodeID=id
		}
	}

	var node UnderstandingNode
	node.ID = id
	node.ParentNode=parent
	node.ParentID=parent.ID
	node.Mood=Mood
	node.EmotionID=EmotionID
	node.SituationID=SituationID
	node.PurposeID=PurposeID

	parent.Children = append(parent.Children, node)
	// четко находим новый вставленный член (а не &parent.Children[count-1])
	var newN *UnderstandingNode
	for i := 0; i < len(parent.Children); i++ {
		if parent.Children[i].ID == node.ID {
			newN = &parent.Children[i]
		}
	}
	// т.к. append меняет длину массива, перетусовывая адреса, то нужно обновить адреса в UnderstandingNodeFromID:
	updateUnderstandingNodeFromID(parent)// здесь потому, что при загрузке из файла нужно на лету получать адреса

	return id,newN
}
// корректируем адреса всех узлов
func updateUnderstandingNodeFromID(parent *UnderstandingNode){
	//updatingUnderstandingNodeFromID(&VernikePhraseTree)
	updatingUnderstandingNodeFromID(parent)
}
// проход всего дерева
func updatingUnderstandingNodeFromID(rt *UnderstandingNode){
	if rt.ID>0 {
		rt.ParentNode=UnderstandingNodeFromID[rt.ParentID] // wr.ParentNode адрес меняется из=за corretsParent(,
		UnderstandingNodeFromID[rt.ID] = rt
	}
	if rt.Children == nil{// конец ветки
		return
	}
	for i := 0; i < len(rt.Children); i++ {
		updatingUnderstandingNodeFromID(&rt.Children[i])
	}
}
///////////////////////////////////////////////////////////
/* загрузить записанное дерево
Формат записи:
ID|ParentNode|Mood|EmotionID|ActivityID|ToneMoodID|SimbolID|VerbalID
*/
func loadUnderstandingTree(){
	createNulLevelUnderstandingTree(&UnderstandingTree)
	strArr,_:=lib.ReadLines(lib.GetMainPathExeFile()+"/memory_psy/understanding_tree.txt")
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
		Mood,_:=strconv.Atoi(p[2])
		EmotionID,_:=strconv.Atoi(p[3])
		SituationID,_:=strconv.Atoi(p[4])
		PurposeID,_:=strconv.Atoi(p[5])

		// новый узел с каждой строкой из файла
		createNewUnderstandingNode(UnderstandingNodeFromID[parentID],id,Mood,EmotionID,
			SituationID,PurposeID)
	}
	return
}
// создать первый, нулевой уровень дерева
func createNulLevelUnderstandingTree(rt *UnderstandingNode){
	rt.ID=0
	UnderstandingNodeFromID[rt.ID]=rt
	return
}
///////////////////////////////////////////
// создать первые три ветки базовых состояний
func createBasicUnderstandingTree(){
	notAllowScanInTreeThisTime=true // запрет показа карты при обновлении
	createNewUnderstandingNode(&UnderstandingTree,0,1,0,0,0)
	createNewUnderstandingNode(&UnderstandingTree,0,2,0,0,0)
	createNewUnderstandingNode(&UnderstandingTree,0,3,0,0,0)
	SaveUnderstandingTree()
	notAllowScanInTreeThisTime = false // запрет показа карты при обновлении
}
/////////////////////////////////////
func SaveUnderstandingTree(){
	if EvolushnStage < 4 { // только со стадии развития 4
		return
	}
	notAllowScanInTreeThisTime=true
	var out=""
	cnt:=len(UnderstandingTree.Children)
	for n := 0; n < cnt; n++ {
		out+=getUnderstandingNode(&UnderstandingTree.Children[n])
	}
	lib.WriteFileContent(lib.GetMainPathExeFile()+"/memory_psy/understanding_tree.txt",out)
	notAllowScanInTreeThisTime = false
	return
}
func getUnderstandingNode(wt *UnderstandingNode)(string){
	var out=""
	//	if wt.ParentID>0 {
	out += strconv.Itoa(wt.ID) + "|"
	out += strconv.Itoa(wt.ParentID) + "|"
	out += strconv.Itoa(wt.Mood) + "|"
	out += strconv.Itoa(wt.EmotionID) + "|"
	out += strconv.Itoa(wt.SituationID) + "|"
	out += strconv.Itoa(wt.PurposeID)
	out +="\r\n"
	//	}
	if(wt.Children==nil){// конец
		return out
	}
	for n := 0; n < len(wt.Children); n++ {
		out+=getUnderstandingNode(&wt.Children[n])
	}
	return out
}
/////////////////////////////////////
// найти КОНЕЧНЫЙ узел по условиям
func FindUnderstandingTreeNodeFromCondition(Mood int,EmotionID int,
	SituationID int,PurposeID int)(int,*UnderstandingNode){
	for k, v := range UnderstandingNodeFromID {
		if v.Mood==Mood && v.EmotionID==EmotionID &&
			v.SituationID==SituationID && v.PurposeID==PurposeID{
			return k, v
		}
	}
	return 0,nil
}
//////////////////////////////////////





//////////////////////////////////////////////////////////////
/* создание иерархии АКТИВНЫХ образов контекстов условий и пусковых стимулов в виде ID образов в [4]int
создать последовательность уровней условий в виде массива  ID последовательности ID уровней
*/
func getUnderstandingActiveConditionsArr(lev1 int, lev2 int, lev3 int, lev4 int)([]int){
	arr:=make([]int, 4)
	arr[0]=lev1
	arr[1]=lev2
	arr[2]=lev3
	arr[3]=lev4
	return arr
}
/////////////////////////////////////////////////////////////////////


//////////////////////////////////////
// создание новой ветки с новым автоматизмом, начиная с заданного узла
func addNewUnderstandingBranchFromNodes(level int,cond []int,node *UnderstandingNode)(int){
	if node==nil {
		return 0
	}
	if level>=len(cond) {
		return node.ID
	}
	var id=0
	switch(level){
	case 0:
		id,_=createNewUnderstandingNode(node,0,cond[0],0,0,0)
	case 1:
		id,_=createNewUnderstandingNode(node,0,cond[0],cond[1],0,0)
	case 2:
		id,_=createNewUnderstandingNode(node,0,cond[0],cond[1],cond[2],0)
	case 3:
		id,_=createNewUnderstandingNode(node,0,cond[0],cond[1],cond[2],cond[3])
	}
	level++
	id=addNewUnderstandingBranchFromNodes(level,cond, UnderstandingNodeFromID[id])
	return id
}
/////////////////////////////////////

// создание ветки, начиная с заданного узла fromID
func formingUnderstandingBranch(fromID int,lastLevel int,condArr []int)(int){
	// нарастить ветку недостающим
	lastNode:=UnderstandingNodeFromID[fromID]

	lastNodeID:=addNewUnderstandingBranchFromNodes(lastLevel,condArr,lastNode)
	if lastNodeID>0{
		SaveUnderstandingTree()
	}
	return lastNodeID
}
/////////////////////////////////////////////////////



// выдать массив узлов ветки по заданному ID узла
func getBrangeUnderstandingNodeArr(lastNodeId int)([]*UnderstandingNode){
	var nArr []*UnderstandingNode
	node:=UnderstandingNodeFromID[lastNodeId]
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
// выдать массив ID узлов ветки по заданному ID узла
func getBrangeUnderstandingNodeIdArr(lastNodeId int)([]int){
	var nArr []int
	node:=UnderstandingNodeFromID[lastNodeId]
	if node==nil { return nil}
	for {
		if node==nil {
			break
		}
		nArr = append(nArr, node.ID)
		node=node.ParentNode
	}
	return nArr
}
//////////////////////////////////////