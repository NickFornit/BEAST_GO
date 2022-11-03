/* Образ ментального Правила: Cтимул (действий оператора или мент.автом=ма) - ответа Beast - Эффекта

*/
package psychic

import (
	"BOT/lib"
	"strconv"
	"strings"
)

/////////////////////////
type MentalTriggerAndAction struct {
	ID int
/*образ пусковых стимулов:
  >0 - образы действий оператора с Пульта ActionsImage,
  <0 - образ действий ментального автоматизма, активировавшего consciousness
 */
	Trigger int// м.б. два вида (см. сверху)
	Action int // образ ответных действий - всегда MentalActionsImages
	Effect int // эффект от действий: -1 или 0 или 1
}
////////////////////////

var MentalTriggerAndActionArr=make(map[int]*MentalTriggerAndAction)

//////////////////////////////////////////

// вызывается из psychic.go
func MentalTriggerAndActionInit(){
	loadMentalTriggerAndActionArr()
}


////////////////////////////////////////////////
// создать новое сочетание ответных действий если такого еще нет
var lastMentalTriggerAndActionID=0
func createNewlastMentalTriggerAndActionID(id int,Trigger int,Action int,Effect int)(int,*MentalTriggerAndAction){
	if Effect<0{Effect=-1}
	if Effect>0{Effect=1}

	oldID,oldVal:=checkUnicumMentalTriggerAndAction(Trigger,Action,Effect)
	if oldVal!=nil{
		return oldID,oldVal
	}
	if id==0{
		lastMentalTriggerAndActionID++
		id=lastMentalTriggerAndActionID
	}else{
		//		newW.ID=id
		if lastMentalTriggerAndActionID<id{
			lastMentalTriggerAndActionID=id
		}
	}

	var node MentalTriggerAndAction
	node.ID = id
	node.Trigger = Trigger
	node.Action=Action
	node.Effect=Effect

	MentalTriggerAndActionArr[id]=&node

	if doWritingFile{SaveMentalTriggerAndActionArr() }

	return id,&node
}
func checkUnicumMentalTriggerAndAction(Trigger int,Action int,Effect int)(int,*MentalTriggerAndAction){
	for id, v := range MentalTriggerAndActionArr {
		if Trigger != v.Trigger {
			continue
		}
		if Action != v.Action {
			continue
		}
		if Effect != v.Effect {
			continue
		}
		return id,v
	}

	return 0,nil
}
/////////////////////////////////////////






//////////////////// сохранить Образы стимула (действий оператора) - ответа Beast
//В случае отсуствия ответных действий создается ID такого отсутсвия, пример такой записи: 2|||0|0|
// ID|ActID через ,|PhraseID через ,|ToneID|MoodID|
func SaveMentalTriggerAndActionArr(){
	var out=""
	for k, v := range MentalTriggerAndActionArr {
		out+=strconv.Itoa(k)+"|"
		out+=strconv.Itoa(v.Trigger)+"|"
		out+=strconv.Itoa(v.Action)+"|"
		out+=strconv.Itoa(v.Effect)
		out+="\r\n"
	}
	lib.WriteFileContent(lib.GetMainPathExeFile()+"/memory_psy/trigger_and_actions_mental.txt",out)

}
////////////////////  загрузить образы стимула (действий оператора) - ответа Beast
// ID|ActID через ,|PhraseID через ,|ToneID|MoodID|
func loadMentalTriggerAndActionArr(){
	MentalTriggerAndActionArr=make(map[int]*MentalTriggerAndAction)
	strArr,_:=lib.ReadLines(lib.GetMainPathExeFile()+"/memory_psy/trigger_and_actions_mental.txt")
	cunt:=len(strArr)
	for n := 0; n < cunt; n++ {
		if len(strArr[n])==0{
			continue
		}
		p:=strings.Split(strArr[n], "|")
		id,_:=strconv.Atoi(p[0])

		Trigger,_:=strconv.Atoi(p[1])
		Action,_:=strconv.Atoi(p[2])
		Effect,_:=strconv.Atoi(p[3])
var saveDoWritingFile= doWritingFile; doWritingFile =false
		createNewlastMentalTriggerAndActionID(id,Trigger,Action,Effect)
doWritingFile =saveDoWritingFile
	}
	return

}
///////////////////////////////////////////

