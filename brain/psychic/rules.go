/* обобщение примитивных правил на основе эпизодической памяти

Заполняется при активации дерева (один эпизод)
и при обобщении эпизодической памяти (последовательность эпизодов).

Бывает Правило со стимулом от Оператора и НЕ бывает со стимулом - по изменению состояния:
//записать ПРАВИЛО, Стимул - ОТ ОПЕРАТОРА
fixNewRules(lastCommonDiffValue,curActiveActionsID)

Цепочки Правил в Эпиз.памяти создабт карту решений в контексте одной темы:
карты местности - куда идти после очередного шага,
карту игры в шахматы: как ходить в данной позиции и на сколько шагов вперед обдумывать решения.

Поиск по Павилам ведется не по отдельным действиям (а это - не только слова, но и тон и/или кнопки действий),
а тупо по ID образов действий. Но, если крепко заматься и вспоминать детально,
то можно делать поиск по всем компонентам действий и делать его мягким.
Заготовлен прототип getSuitableMentalRulesСarefully() для использования в ментальных инфо-функциях.
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
type rules struct {
	ID int
	// условия точного использования Правила:
	NodeAID int // конечный узел активной ветки дерева моторных автоматизмов detectedActiveLastNodID
	NodePID int // конечный узел активной ветки дерева ментальных автоматизмов detectedActiveLastUnderstandingNodID

	TAid []int // цепочка стимул-ответов ID TriggerAndAction - последовательность из эпизодов памяти подряд, сохраняющая последовательность общения ( дурак -> сам дурак!, маме скажу -> ябеда, щас в морду дам -> ну попробуй)
}
var rulesArr=make(map[int]*rules)
// выборка по условиям Правила: "NodeAID_NodePID"
//sinex:=strconv.Itoa(NodeAID)+"_"+strconv.Itoa(NodePID); rulesArrConditinArr[sinex]
var rulesArrConditinArr=make(map[string] []*rules)// Массив Правил

//////////////////////////////////////////


///////////////////////////////////////////
// вызывается из psychic.go
func rulesInit(){
	loadrulesArr()

//	getCur10lastRules()
}


////////////////////////////////////////////////
// создать новое сочетание ответных действий если такого еще нет
var lastrulesID=0
var isNotLoading=true
func createNewRules(id int,NodeAID int,NodePID int,TAid []int,CheckUnicum bool)(int,*rules){

	if TAid == nil{
		return 0,nil
	}
	if CheckUnicum {
		oldID,oldVal:=checkUnicumrules(NodeAID,NodePID,TAid)
		if oldVal!=nil{
			return oldID,oldVal
		}
	}

	if id==0{
		lastrulesID++
		id=lastrulesID
	}else{
		//		newW.ID=id
		if lastrulesID<id{
			lastrulesID=id
		}
	}

	var node rules
	node.ID = id
	node.TAid=TAid
	node.NodeAID = NodeAID
	node.NodePID = NodePID

	rulesArr[id]=&node
	sinex:=strconv.Itoa(NodeAID)+"_"+strconv.Itoa(NodePID)
	rulesArrConditinArr[sinex]=append(rulesArrConditinArr[sinex],&node)

	if doWritingFile{
		SaveRulesArr()
	}
	if isNotLoading {
		if len(TAid) > 1 {
			lib.WritePultConsol("<span style='color:green'>Записано групповое <b>ПРАВИЛО № " + strconv.Itoa(id) + "</b></span>")
		} else {
			lib.WritePultConsol("<span style='color:green'>Записано <b>ПРАВИЛО № " + strconv.Itoa(id) + "</b></span>")
		}
	}

	return id,&node
}
func checkUnicumrules(NodeAID int,NodePID int,TAid []int)(int,*rules){
	for id, v := range rulesArr {
		if NodeAID != v.NodeAID || NodePID != v.NodePID || !lib.EqualArrs(TAid,v.TAid) {
			continue
		}
		return id,v
	}

	return 0,nil
}
/////////////////////////////////////////






//////////////////// сохранить Образы rules
// ID|NodeAID|NodePID|TAid через ,
func SaveRulesArr(){
	var out=""
	for k, v := range rulesArr {
		out+=strconv.Itoa(k)+"|"
		out+=strconv.Itoa(v.NodeAID)+"|"
		out+=strconv.Itoa(v.NodePID)+"|"
		for i := 0; i < len(v.TAid); i++ {
			out+=strconv.Itoa(v.TAid[i])+","
		}
		out+="\r\n"
	}
	lib.WriteFileContent(lib.GetMainPathExeFile()+"/memory_psy/rules.txt",out)

}
////////////////////  загрузить образы rules
// ID|NodeAID|NodePID|TAid через ,
func loadrulesArr(){
	rulesArr=make(map[int]*rules)
	rulesArrConditinArr=make(map[string] []*rules)
	strArr,_:=lib.ReadLines(lib.GetMainPathExeFile()+"/memory_psy/rules.txt")
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
		isNotLoading=false
		createNewRules(id,NodeAID,NodePID,TAid,false)
		isNotLoading=true
doWritingFile =saveDoWritingFile
	}
	return

}
///////////////////////////////////////////


