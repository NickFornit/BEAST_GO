/* Образ обощенного (без условия применения) Правила: 
Cтимул (действий оператора) - ответа Beast - Эффекта

Правило приобретает определнную значимость (смысл) в контексте опредеелнных условий.
Поэтому нет никакого смысла пытаться использовать обобщенное правило само по себе 
т.к. одно и то же сочетания Стиму и Ответ может иметь разный эФфект.
НО если для данного сочетания Стимул-Ответ есть только один вод эффекта, 
то это - уже Информация, и чем больше опыт (количество обобщенных правил), тем такая информация полезнее.
*/

package psychic

import (
	"BOT/lib"
	"strconv"
	"strings"
)

//////////////////////////////////

type TriggerAndAction struct {
	ID int
/*образ пусковых стимулов:
  >0 - образы действий оператора с Пульта ActionsImage,
  <0 - образ текущего состояния BaseStateImage
*/
	Trigger int
	Action int // образ ответных действий - всегда ActionsImage
	Effect int // эффект от действий: -1 или 0 или 1
}
////////////////////////

var TriggerAndActionArr=make(map[int]*TriggerAndAction)

//////////////////////////////////////////

// вызывается из psychic.go
func TriggerAndActionInit(){
	loadTriggerAndActionArr()
}


////////////////////////////////////////////////
// создать новое правило если такого еще нет
var lastTriggerAndActionID=0
func createNewlastTriggerAndActionID(id int,Trigger int,Action int,Effect int,CheckUnicum bool)(int,*TriggerAndAction){
	if Effect<0{Effect=-1}
	if Effect>0{Effect=1}

	if CheckUnicum {
		oldID,oldVal:=checkUnicumTriggerAndAction(Trigger,Action,Effect)
		if oldVal!=nil{
			return oldID,oldVal
		}
	}

	if id==0{
		lastTriggerAndActionID++
		id=lastTriggerAndActionID
	}else{
		//		newW.ID=id
		if lastTriggerAndActionID<id{
			lastTriggerAndActionID=id
		}
	}

	var node TriggerAndAction
	node.ID = id
	node.Trigger = Trigger
	node.Action=Action
	node.Effect=Effect

	TriggerAndActionArr[id]=&node

	if doWritingFile{SaveTriggerAndActionArr() }

	return id,&node
}
func checkUnicumTriggerAndAction(Trigger int,Action int,Effect int)(int,*TriggerAndAction){
	for id, v := range TriggerAndActionArr {
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
func SaveTriggerAndActionArr(){
	var out=""
	for k, v := range TriggerAndActionArr {
		out+=strconv.Itoa(k)+"|"
		out+=strconv.Itoa(v.Trigger)+"|"
		out+=strconv.Itoa(v.Action)+"|"
		out+=strconv.Itoa(v.Effect)
		out+="\r\n"
	}
	lib.WriteFileContent(lib.GetMainPathExeFile()+"/memory_psy/trigger_and_actions.txt",out)

}
////////////////////  загрузить образы стимула (действий оператора) - ответа Beast
// ID|ActID через ,|PhraseID через ,|ToneID|MoodID|
func loadTriggerAndActionArr(){
	TriggerAndActionArr=make(map[int]*TriggerAndAction)
	strArr,_:=lib.ReadLines(lib.GetMainPathExeFile()+"/memory_psy/trigger_and_actions.txt")
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
		createNewlastTriggerAndActionID(id,Trigger,Action,Effect,false)
doWritingFile =saveDoWritingFile
	}
	return

}
///////////////////////////////////////////
