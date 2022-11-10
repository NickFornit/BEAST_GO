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
	TAid []int // цепочка стимул-ответов ID MentalTriggerAndAction - последовательность из эпизодов памяти подряд, сохраняющая последовательность общения ( дурак -> сам дурак!, маме скажу -> ябеда, щас в морду дам -> ну попробуй)
}
var rulesMentalArr=make(map[int]*rulesMental)

//////////////////////////////////////////

// вызывается из psychic.go
func rulesMentalInit(){
	loadrulesMentalArr()
}


////////////////////////////////////////////////
// создать новое сочетание ответных действий если такого еще нет
var lastrulesMentalID=0
var isNotMentLoading=true
func createNewlastrulesMentalID(id int,TAid []int)(int,*rulesMental){
	if TAid == nil{
		return 0,nil
	}

	oldID,oldVal:=checkUnicumrulesMental(TAid)
	if oldVal!=nil{
		return oldID,oldVal
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
func checkUnicumrulesMental(TAid []int)(int,*rulesMental){
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
// ID|TAid через ,
func SaverulesMentalArr(){
	var out=""
	for k, v := range rulesMentalArr {
		out+=strconv.Itoa(k)+"|"
		for i := 0; i < len(v.TAid); i++ {
			out+=strconv.Itoa(v.TAid[i])+","
		}
		out+="\r\n"
	}
	lib.WriteFileContent(lib.GetMainPathExeFile()+"/memory_psy/rulesMental.txt",out)

}
////////////////////  загрузить образы rulesMental
// ID|TAid через ,
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

		s:=strings.Split(p[1], ",")
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
		createNewlastrulesMentalID(id,TAid)
		isNotMentLoading=true
doWritingFile =saveDoWritingFile
	}
	return

}
///////////////////////////////////////////


