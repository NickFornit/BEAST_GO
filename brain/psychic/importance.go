/* Значимость элементов восприятия - как объекта произвольного внимания:
того из всего воспринимаемого, что имеет наибольшую значимость
т.к. именно наибольшая значимость должна осмысливаться.

Кроме того, значимости объектов - это и есть модель понимания данного объекта внимания -
его значимость в разных условиях и то, какие действия могут быть совершены при этом.

Сохраняется в файле importance.txt в формате ID|NodeAID|NodePID|Type|ObjectID|Value

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
// названия типов объектов значимости
var importanceTypeName=[]string{
	1:"Образ действия",
	2:"образ мысли",
	3:"Моторное действие",
	4:"Вербальный образ",
	5:"Фраза",
	6:"Слово",
	7:"Тон",
	8:"Настроение",
}

/* Для определения текущих объектов восприятия и выделения одного из них - самого важного НЕГАТИВНОГО
		по всем категориям importanceType
При каждом ОБъективном вызове consciousness определяется текущий объект наибольшой значимости в воспринимаемом -
в функции определения текущей Цели getMentalPurpose()
*/
// текущий объект внимания
var extremImportanceObject *extremImportance
// текущий субъект внимания (объект внимания к собственным мыслям)    ЕЩЕ НЕТ ТАКОГО по аналогии с importance.go
var extremImportanceMentalObject *extremImportance

type extremImportance struct {
	objID int//  объект значимости
	kind int // тип объекта
	extremVal int// экстремальная значимость
}
var curImportanceObjectArr []extremImportance //- здесь сохраняются текущие цели внимания к наиболее важному


// Тот объект, применение которого привело к данной значимости  ID|NodeAID|NodePID|Type|ObjectID|Value
type importance struct {
	ID int
	// условия точного использования Правила:
	NodeAID int // конечный узел активной ветки дерева моторных автоматизмов detectedActiveLastNodID
	NodePID int // конечный узел активной ветки дерева ментальных автоматизмов detectedActiveLastUnderstandingNodID

	Type int // тип объекта importanceType
	ObjectID int // ID объекта восприянтия
	Value int // значимость от -10 0 до 10 Чудовищно большие будут говорить о сверх значимости.
}
var importanceFromID=make(map[int]*importance)
// выборка по условиям: "NodeAID_NodePID"
//sinex:=strconv.Itoa(NodeAID)+"_"+strconv.Itoa(NodePID); importanceConditinArr[sinex]
var importanceConditinArr=make(map[string] []*importance)// Массив значимостей при данном условии


//объекты внимания, имеющую высокую значимость, не сохраняются в памяти и гасятки во сне
var importanceObjectID []*extremImportance
///////////////////////////////////////////
// вызывается из psychic.go
func ImportanceInit(){
	loadImportance()
}


////////////////////////////////////////////////
// создать новый образ значимости объекта восприятия если такого еще нет
var lastImportanceID=0
var isNotImpLoading=true
func createNewlastImportanceID(id int,NodeAID int,NodePID int,Type int,ObjectID int,Value int,CheckUnicum bool)(int,*importance){

	if CheckUnicum {
		oldID,oldVal:=checkUnicumImportance(NodeAID,NodePID,Type,ObjectID)
		if oldVal!=nil{
			oldVal.Value+=Value
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
	node.ObjectID=ObjectID
	node.Value=Value

	importanceFromID[id]=&node
	sinex:=strconv.Itoa(NodeAID)+"_"+strconv.Itoa(NodePID)
	importanceConditinArr[sinex]=append(importanceConditinArr[sinex],&node)

	if doWritingFile{
		Saveimportance()
	}

	return id,&node
}
func checkUnicumImportance(NodeAID int,NodePID int,Type int,ObjectID int)(int,*importance){
	for id, v := range importanceFromID {
		if NodeAID != v.NodeAID || NodePID != v.NodePID || Type != v.Type || ObjectID != v.ObjectID {
			continue
		}
		return id,v
	}

	return 0,nil
}
/////////////////////////////////////////






//////////////////// сохранить Образы Importance
// ID|NodeAID|NodePID|Type|ObjectID|Value
func Saveimportance(){
	var out=""
	for k, v := range importanceFromID {
		out+=strconv.Itoa(k)+"|"
		out+=strconv.Itoa(v.NodeAID)+"|"
		out+=strconv.Itoa(v.NodePID)+"|"
		out+=strconv.Itoa(v.Type)+"|"
		out+=strconv.Itoa(v.ObjectID)+"|"
		out+=strconv.Itoa(v.Value)

		out+="\r\n"
	}
	lib.WriteFileContent(lib.GetMainPathExeFile()+"/memory_psy/importance.txt",out)

}
////////////////////  загрузить образы Importance
// ID|NodeAID|NodePID|Type|Value
func loadImportance(){
	importanceFromID=make(map[int]*importance)
	importanceConditinArr=make(map[string] []*importance)
	strArr,_:=lib.ReadLines(lib.GetMainPathExeFile()+"/memory_psy/importance.txt")
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
		ObjectID,_:=strconv.Atoi(p[4])
		Value,_:=strconv.Atoi(p[5])

		var saveDoWritingFile= doWritingFile; doWritingFile =false
		isNotImpLoading=false
		createNewlastImportanceID(id,NodeAID,NodePID,Type,ObjectID,Value,false)
		isNotImpLoading=true
		doWritingFile =saveDoWritingFile
	}
	return

}
/////////////////////////////////////////





///////////////////////////////////////////////////////
/* объект значимости (kind int,ObjectID int) - точно по условиям
Возвращает:
ссылку на объект значимости
*/
func getObjectImportance(kind int,ObjectID int, NodeAID int,NodePID int)(*importance){
	sinex:=strconv.Itoa(NodeAID)+"_"+strconv.Itoa(NodePID)
	for _, v := range importanceConditinArr[sinex] {
		if v.Type==kind && v.ObjectID ==ObjectID{
			return v
		}
	}

	return nil
}
///////////////////////////////////////////////////////////
/* объект значимости (kind int,ObjectID int) - игнорируя условия detectedActiveLastNodID detectedActiveLastUnderstandingNodID
Возвращает массив найденных объектов
*/
func getObjectImportanceЦithoutСonditions(kind int,ObjectID int)([]*importance){
	var outArr []*importance
	for _, v := range importanceFromID {
		if v.Type==kind && v.ObjectID ==ObjectID{
			outArr=append(outArr,v)
		}
	}

	return outArr
}
///////////////////////////////////////////////////////////


///////////////////////////////////////////////////////
/* значимость ID объекта внимания (kind int,ObjectID int)
в условиях detectedActiveLastNodID, detectedActiveLastUnderstandingNodID
Возвращает:
1 - ID объекта значимости
2 - значимость
*/
func getObjectsImportanceValue(kind int,ObjectID int, NodeAID int,NodePID int)(int,int){
	sinex:=strconv.Itoa(NodeAID)+"_"+strconv.Itoa(NodePID)
	for _, v := range importanceConditinArr[sinex] {
		if v.Type==kind && v.ObjectID ==ObjectID{
			return v.ID,v.Value
		}
	}

	return 0,0
}
///////////////////////////////////////////////////////////















