/* Значимость элементов восприятия - как объекта произвольного внимания:
того из всего воспринимаемого, что имеет наибольшую значимость
т.к. именно наибольшая значимость должна осмысливаться.

Сохраняется в файле importance.txt в формате ID|NodeAID|NodePID|Type|Value

Значимость - величина от -10 0 до 10, приобретаемая объектов внимания в данной ситуации
- берется из результата пробных действий и связывается ос всеми компонентами воспринимаемого в этих условиях
функцией setImportance().

Оценке значимости подлежат элементы действия оператора:
кнопки воздействия, фразы и отдельные слова, принимающие значимость фразы.

Значимость всегда определяется в контексте всех предшествующих условий,
т.е. специфична для активностей деревьев автоматизмов и понимания.

При каждом вызове consciousness определяется текущий объект наибольшой значимости в воспринимаемом -
в функции определения текущей Цели getMentalPurpose().

!!!Если уже была определена иная значимость ранее для объекта внимания в условиях "NodeAID_NodePID", то создается новый образ
так что можно получить обобщенное значение значимости методом простого усреднения или учитывать возможные граничные случаи:
func getAllObjectsImportance
*/

package psychic

import (
	"BOT/lib"
	"strconv"
	"strings"
)

/////////////////////////////////////////////////

// типы объектов значимости
var importanceType=[]int{
	1, // ID ActionsImage
	2, // ID MentalActionsImages
	3, // ID несловестного действия ActionsImage.ActID[n]
	4, // ID Verbal - при активации дерева автоматизмов
	5, // ID отдельной фразы Verbal.PhraseID[n]
	6,// ID отдельного слова  из Verbal.PhraseID[n]
	7,// ID тон сообщения с Пульта  Verbal.ToneID
	8,// ID настроение оператора  Verbal.MoodID
}


// Тот объект, применение которого привело к данной значимости  ID|NodeAID|NodePID|Type|Value
type importance struct {
	ID int
	// условия точного использования Правила:
	NodeAID int // конечный узел активной ветки дерева моторных автоматизмов detectedActiveLastNodID
	NodePID int // конечный узел активной ветки дерева ментальных автоматизмов detectedActiveLastUnderstandingNodID

	Type int // тип объекта importanceType
	Value int // значимость от -10 0 до 10
}
var importanceFromID=make(map[int]*importance)
// выборка по условиям: "NodeAID_NodePID"
//sinex:=strconv.Itoa(NodeAID)+"_"+strconv.Itoa(NodePID); importanceConditinArr[sinex]
var importanceConditinArr=make(map[string] []*importance)// Массив значимостей при данном условии


/////////////////////////////////////////////////
///////////////////////////////////////////
// вызывается из psychic.go
func ImportanceInit(){
	loadImportance()
}


////////////////////////////////////////////////
// создать новое сочетание ответных действий если такого еще нет
var lastImportanceID=0
var isNotImpLoading=true
func createNewlastImportanceID(id int,NodeAID int,NodePID int,Type int,Value int,CheckUnicum bool)(int,*importance){

	if CheckUnicum {
		oldID,oldVal:=checkUnicumImportance(NodeAID,NodePID,Type,Value)
		if oldVal!=nil{
			return oldID,oldVal
		}
	}

	if id==0{
		lastImportanceID++
		id=lastImportanceID
	}else{
		//		newW.ID=id
		if lastImportanceID<id{
			lastImportanceID=id
		}
	}

	var node importance
	node.ID = id
	node.NodeAID=NodeAID
	node.NodePID=NodePID
	node.Type=Type
	node.Value=Value

	importanceFromID[id]=&node
	sinex:=strconv.Itoa(NodeAID)+"_"+strconv.Itoa(NodePID)
	importanceConditinArr[sinex]=append(importanceConditinArr[sinex],&node)

	if doWritingFile{
		Saveimportance()
	}

	return id,&node
}
func checkUnicumImportance(NodeAID int,NodePID int,Type int,Value int)(int,*importance){
	for id, v := range importanceFromID {
		if NodeAID != v.NodeAID || NodePID != v.NodePID || Type != v.Type || Value != v.Value  {
			continue
		}
		return id,v
	}

	return 0,nil
}
/////////////////////////////////////////






//////////////////// сохранить Образы Importance
// ID|NodeAID|NodePID|Type|Value
func Saveimportance(){
	var out=""
	for k, v := range importanceFromID {
		out+=strconv.Itoa(k)+"|"
		out+=strconv.Itoa(v.NodeAID)+"|"
		out+=strconv.Itoa(v.NodePID)+"|"
		out+=strconv.Itoa(v.Value)+"|"

		out+="\r\n"
	}
	lib.WriteFileContent(lib.GetMainPathExeFile()+"/memory_psy/importance.txt",out)

}
////////////////////  загрузить образы Importance
// ID|NodeAID|NodePID|Type|Value
func loadImportance(){
	importanceFromID=make(map[int]*importance)
	importanceConditinArr=make(map[string] []*importance)
	strArr,_:=lib.ReadLines(lib.GetMainPathExeFile()+"/memory_psy/Iimportance.txt")
	cunt:=len(strArr)
	for n := 0; n < cunt; n++ {
		if len(strArr[n])==0{
			continue
		}
		p:=strings.Split(strArr[n], "|")
		id,_:=strconv.Atoi(p[0])
		NodeAID,_:=strconv.Atoi(p[1])
		NodePID,_:=strconv.Atoi(p[2])
		Type,_:=strconv.Atoi(p[3])
		Value,_:=strconv.Atoi(p[4])

		var saveDoWritingFile= doWritingFile; doWritingFile =false
		isNotImpLoading=false
		createNewlastImportanceID(id,NodeAID,NodePID,Type,Value,false)
		isNotImpLoading=true
		doWritingFile =saveDoWritingFile
	}
	return

}
/////////////////////////////////////////


// получить все значимости объекта внимания в контексте текущих условий
func getAllObjectsImportance(importanceTypeID int, NodeAID int,NodePID int) []*importance {
	sinex:=strconv.Itoa(NodeAID)+"_"+strconv.Itoa(NodePID)
var imp []*importance

	for _, v := range importanceConditinArr[sinex] {
		if v.Type==importanceTypeID{
			imp=append(imp,v)
		}
	}
	if imp!=nil{
		return imp
	}

	return imp
}
///////////////////////////////////////////////////////

/* средняя значимость объекта внимания
Возвращает:
1 - среднее значение значимости
2 - число значений
 */
func getObjectsImportanceValue(importanceTypeID int, NodeAID int,NodePID int)(int,int){
	sinex:=strconv.Itoa(NodeAID)+"_"+strconv.Itoa(NodePID)
	var value=0
	var count=0

	for _, v := range importanceConditinArr[sinex] {
		if v.Type==importanceTypeID{
			value+=v.Value
			count++
		}
	}
	if count>0{
		return value,count
	}

	return 0,0
}
///////////////////////////////////////////////////////////












