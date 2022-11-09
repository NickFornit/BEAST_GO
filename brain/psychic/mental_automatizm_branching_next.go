/* Переходы на следующий автоматизм от предыдущего в цепи 
с ветвлением в зависимости от BranchID узла дерева моторных автоматизмов AutomatizmNode

Это позволит:
1. убрать NextID из структуры ментальных автоматизмов и сделать автоматизмы универсальными
2. делать ветвление следующего звена в цепочке автоматизмв в зависимости от активной BranchID узла дерева моторных автоматизмов AutomatizmNode

К узлу UnderstandingNode цепляется не автоматизм, а goNext
по goNext.NextID получаем первый ID мент.автоматизма
и создается соотвествующий goNextFromUnderstandingNodeIDArr

При добавлении звена в цепочку создается новая goNext с goNext.NextID следующего ID мент.автоматизма.
В каждом звене цепочки - свой goNext - направляет цепочку с ветвелениями по активному MotorBranchID
т.е. проход цепи идет по последовательности goNext.

Формат в файле ID|MotorBranchID|FromID|NextID|AutomatizmID|
*/


package psychic

import (
	"BOT/lib"
	"strconv"
	"strings"
)

////////////////////////////////

// вызывается из psychic.go
func goNextInit(){
	loadgoNext()
}
//////////////////////////////////////////////


/* разветвитель следующих по цепочке ментальных автоматизмов
*/
type goNext struct {
	ID int
	MotorBranchID int // узел AutomatizmNode
	FromID int // ID родительского goNext
	NextID int // ID дочернего goNext
	AutomatizmID int // id ментального автоматизма
}
// по ID goNext -> *goNext
var goNextFromIDArr=make(map[int]*goNext)
/* по узлу UnderstandingNode.ID -> goNext
Может быть множество goNext, привязанных к узлу UnderstandingNode.ID
Карта сохраняет первые вхождения всех цепочек от самого их начала и создается
при формировании первого автоматизма от узла дерева понимания.
ПРИМЕР uID - это ID узла дерева понимания:
uID=15 -> goNextFromIDArr[goNext.ID=1] -> goNextFromIDArr[goNext.ID=2] - цепочка для MotorBranchID==1
		-> goNextFromIDArr[goNext.ID=10] -> -> goNextFromIDArr[goNext.ID=11] - цепочка для MotorBranchID==2
Для каждого MotorBranchID создается своя линейная цепочка!

К узлу uID дерева понимания пристегнуты 2 goNext с разными MotorBranchID,
а к тем пристегнуты по 2 других goNext с разными MotorBranchID
goNextFromUnderstandingNodeIDArr[1]=goNext.ID=1
goNextFromUnderstandingNodeIDArr[2]=goNext.ID=10

goNextFromUnderstandingNodeIDArr представляет собой память о том, как шли размышения до конечного звена:
создания успешного моторного автомтаизма.
При этом ранее созданные успешные могут стать неуспешными к каких-то условиях и построение цепи продолжится.
Но последовательности Правил в эпиз.памяти хранят случаи успешного действия и они могут использоваться.
Так что необходим отслеживающий режим, который выявляет, какое именно Правило нужно применить в данном случае.
*/
var goNextFromUnderstandingNodeIDArr=make(map[int][]*goNext)

//////////////////////////////////////////




////////////////////////////////////////////////
/* начать новую цепочку
Возвращает новый *goNext
 */
func createNewNextFromUnderstandingNodeID(UnderstandingNodeID int,
	MotorBranchID int,
	FromID int, // записывается текущее currrentFromNextID
	AutomatizmID int)(int,*goNext){

	// если ли уже такое начало цепочки?
	firstArr:=goNextFromUnderstandingNodeIDArr[UnderstandingNodeID]
	if firstArr != nil {
		for i := 0; i < len(firstArr); i++ {
			if firstArr[i].MotorBranchID == MotorBranchID{
				return firstArr[i].ID, goNextFromIDArr[firstArr[i].ID]
			}
		}
	}
	id,newID:=createNewlastgoNextID(0,MotorBranchID,FromID,0,AutomatizmID)
	goNextFromUnderstandingNodeIDArr[UnderstandingNodeID]=append(goNextFromUnderstandingNodeIDArr[UnderstandingNodeID],newID)
// прицепить мент.автоматизм
createMentalAutomatizmID(0, mentalInfoStruct.mImgID, 1)

	return id,newID
}
///////////////////////////////////////////////

////////////////////////////////////////////////
// создать новую goNext
var lastgoNextID=0
func createNewlastgoNextID(id int,MotorBranchID int,FromID int,NextID int,AutomatizmID int)(int,*goNext){

//	Практически все goNext должны быть уникальными и не требуют проверки на уникальность.
	if id==0{
		lastgoNextID++
		id=lastgoNextID
	}else{
		//		newW.ID=id
		if lastgoNextID<id{
			lastgoNextID=id
		}
	}

	var node goNext
	node.ID = id
	node.MotorBranchID=MotorBranchID
	node.FromID=FromID
	node.NextID=NextID
	node.AutomatizmID=AutomatizmID

	goNextFromIDArr[id]=&node

	if doWritingFile{SavegoNext() }

	return id,&node
}
/////////////////////////////////////////






//////////////////// сохранить Образы goNext
// ID|MotorBranchID|FromID|NextID|AutomatizmID|
func SavegoNext(){
	var out=""
	for k, v := range goNextFromIDArr {
		out+=strconv.Itoa(k)+"|"
		out+=strconv.Itoa(v.MotorBranchID)+"|"
		out+=strconv.Itoa(v.FromID)+"|"
		out+=strconv.Itoa(v.NextID)+"|"
		out+=strconv.Itoa(v.AutomatizmID)
		out+="\r\n"
	}
	lib.WriteFileContent(lib.GetMainPathExeFile()+"/memory_psy/goNext.txt",out)

}
////////////////////  загрузить образы goNext
// ID|ActID через ,|PhraseID через ,|ToneID|MoodID|
func loadgoNext(){
	goNextFromIDArr=make(map[int]*goNext)
	strArr,_:=lib.ReadLines(lib.GetMainPathExeFile()+"/memory_psy/goNext.txt")
	cunt:=len(strArr)
	for n := 0; n < cunt; n++ {
		if len(strArr[n])==0{
			continue
		}
		p:=strings.Split(strArr[n], "|")
		id,_:=strconv.Atoi(p[0])
		MotorBranchID,_:=strconv.Atoi(p[1])
		FromID,_:=strconv.Atoi(p[2])
		NextID,_:=strconv.Atoi(p[3])
		AutomatizmID,_:=strconv.Atoi(p[4])

		var saveDoWritingFile= doWritingFile; doWritingFile =false
		createNewlastgoNextID(id,MotorBranchID,FromID,NextID,AutomatizmID)
		doWritingFile =saveDoWritingFile
	}
	return
}
///////////////////////////////////////////




