/* обобщение ментальных правил на основе ментальных кадров эпизодической памяти

Цепочки Правил в Эпиз.памяти создабт карту решений в контексте одной темы:
карты местности - куда идти после очередного шага,
карту игры в шахматы: как ходить в данной позиции и на сколько шагов вперед обдумывать решения.

Правило "На что обращать внимание" - перечисляет автоматизмы инфо-функций выделения обпределенных объектов внимания
Шпионы развивают это правило, перечислением типов объектов внимания.
Цель - переактивация с помещением в эпиз.память всех этих объектов внимания.
*/


package psychic

import (
	"BOT/lib"
	"strconv"
	"strings"
)

///////////////////////////////////

/* Правила примитивного опыта, обобщающие стимулы->ответы->эффект для таких цепочек в диалогах
На основе этих правил становятся возможны более системные обобщения.
 */
type rulesMental struct {
	ID int
	// условия точного использования Правила:
	NodeAID int // конечный узел активной ветки дерева моторных автоматизмов detectedActiveLastNodID
	NodePID int // конечный узел активной ветки дерева ментальных автоматизмов detectedActiveLastUnderstandingNodID

	TAid []int // цепочка стимул-ответов ID MentalTriggerAndAction - последовательность из эпизодов памяти подряд, сохраняющая последовательность общения ( дурак -> сам дурак!, маме скажу -> ябеда, щас в морду дам -> ну попробуй)
}
var rulesMentalArr=make(map[int]*rulesMental)

// выборка по условиям Правила: "NodeAID_NodePID"
//sinex:=strconv.Itoa(NodeAID)+"_"+strconv.Itoa(NodePID); rulesArrConditinArr[sinex]
var rulesMentalArrConditinArr=make(map[string] []*rulesMental)// Массив Правил
//////////////////////////////////////////

// вызывается из psychic.go
func rulesMentalInit(){
	loadrulesMentalArr()
}


////////////////////////////////////////////////
// создать новое ментальное правило если такого еще нет
var lastrulesMentalID=0
var isNotMentLoading=true
func createNewRulesMentalID(id int,NodeAID int,NodePID int,TAid []int,CheckUnicum bool)(int,*rulesMental){
	if TAid == nil{
		return 0,nil
	}
	if CheckUnicum {
		oldID,oldVal:=checkUnicumrulesMental(NodeAID,NodePID,TAid)
		if oldVal!=nil{
			return oldID,oldVal
		}
	}

	if id==0{
		lastrulesMentalID++
		id=lastrulesMentalID
	}else{
		//		newW.ID=id
		if lastrulesMentalID<id{
			lastrulesMentalID=id
		}
	}

	var node rulesMental
	node.ID = id
	node.TAid=TAid

	rulesMentalArr[id]=&node
	sinex:=strconv.Itoa(NodeAID)+"_"+strconv.Itoa(NodePID)
	rulesMentalArrConditinArr[sinex]=append(rulesMentalArrConditinArr[sinex],&node)

	if doWritingFile{
		SaverulesMentalArr()

	}
	if isNotMentLoading {
		if len(TAid)>1{
			lib.WritePultConsol("<span style='color:green'>Записано групповое <b>ПРАВИЛО № " + strconv.Itoa(id) + "</b></span>")
		}else{
			lib.WritePultConsol("<span style='color:green'>Записано <b>ПРАВИЛО № " + strconv.Itoa(id) + "</b></span>")
		}
	}

	return id,&node
}
func checkUnicumrulesMental(NodeAID int,NodePID int,TAid []int)(int,*rulesMental){
	for id, v := range rulesMentalArr {
		if !lib.EqualArrs(TAid,v.TAid) {
			continue
		}
		return id,v
	}

	return 0,nil
}
/////////////////////////////////////////






//////////////////// сохранить Образы rulesMental
// ID|NodeAID|NodePID|TAid через ,
func SaverulesMentalArr(){
	var out=""
	for k, v := range rulesMentalArr {
		out+=strconv.Itoa(k)+"|"
		out+=strconv.Itoa(v.NodeAID)+"|"
		out+=strconv.Itoa(v.NodePID)+"|"
		for i := 0; i < len(v.TAid); i++ {
			out+=strconv.Itoa(v.TAid[i])+","
		}
		out+="\r\n"
	}
	lib.WriteFileContent(lib.GetMainPathExeFile()+"/memory_psy/rulesMental.txt",out)

}
////////////////////  загрузить образы rulesMental
// ID|NodeAID|NodePID|TAid через ,
func loadrulesMentalArr(){
	rulesMentalArr=make(map[int]*rulesMental)
	strArr,_:=lib.ReadLines(lib.GetMainPathExeFile()+"/memory_psy/rulesMental.txt")
	cunt:=len(strArr)
	for n := 0; n < cunt; n++ {
		if len(strArr[n])==0{
			continue
		}
		p:=strings.Split(strArr[n], "|")
		id,_:=strconv.Atoi(p[0])
		NodeAID,_:=strconv.Atoi(p[1])
		NodePID,_:=strconv.Atoi(p[2])

		s:=strings.Split(p[3], ",")
		var TAid[] int
		for i := 0; i < len(s); i++ {
			if len(s[i])==0{
				continue
			}
			si,_:=strconv.Atoi(s[i])
			TAid=append(TAid,si)
		}
var saveDoWritingFile= doWritingFile; doWritingFile =false
		isNotMentLoading=false
		createNewRulesMentalID(id,NodeAID,NodePID,TAid,false)
		isNotMentLoading=true
doWritingFile =saveDoWritingFile
	}
	return

}
///////////////////////////////////////////


