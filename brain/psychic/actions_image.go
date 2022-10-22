/* образы действий, оператора или действий Beast.
Фактически структура повторяет ActionsImage из рефлексов и позволяет сохранять
как образы действий в автоматизмах, так и образы действий оператора, отражаемые в дереве мот.автомтаизмов.
Используется для формирования пар стимул (действия оператора) - действия (ответ beast) 
для эпизодической памяти и структуры rules - Правил примитивного опыта.

Формат: ID|RSarr через ,|PhraseID через ,|ToneID|MoodID|
*/

package psychic

import (
	termineteAction "BOT/brain/terminete_action"
	word_sensor "BOT/brain/words_sensor"
	"BOT/lib"
	"strconv"
	"strings"
)

///////////////////////////////////
type ActionsImage struct {
	ID    int   // идентификатор данного сочетания пусковых стимулов
	ActID []int // массив действий с Пульта
	// для текущего сообщения с Пусльта:
	PhraseID []int // массив фразID (DetectedUnicumPhraseID) слова каждой фразы вытаскиваются wordSensor.GetWordArrFromPhraseID(PhraseID[0])
	ToneID   int   // тон сообщения с Пульта
	MoodID   int   // настроение оператора
}

var ActionsImageArr=make(map[int]*ActionsImage)

//////////////////////////////////////////

// вызывается из psychic.go
func ActionsImageInit(){
	loadActionsImageArr()
}


////////////////////////////////////////////////
// создать новое сочетание ответных действий если такого еще нет 
var lastActionsImageID=0
func СreateNewlastActionsImageID(id int,ActID []int,PhraseID []int,ToneID int,MoodID int)(int,*ActionsImage){

	if id==0{
		lastActionsImageID++
		id=lastActionsImageID
	}else{
		//		newW.ID=id
		if lastActionsImageID<id{
			lastActionsImageID=id
		}
	}

	var node ActionsImage
	node.ID = id
	node.ActID = ActID
	node.PhraseID=PhraseID
	node.ToneID=ToneID
	node.MoodID=MoodID

	ActionsImageArr[id]=&node

	if doWritingFile { SaveActionsImageArr() }

	return id,&node
}
func checkUnicumActionsImage(ActID []int,PhraseID []int,ToneID int,MoodID int)(int,*ActionsImage){
	for id, v := range ActionsImageArr {
		if !lib.EqualArrs(ActID,v.ActID) {
			continue
		}
		if !lib.EqualArrs(PhraseID,v.PhraseID) {
			continue
		}
		if ToneID!=v.ToneID || MoodID!=v.MoodID {
			continue
		}
		return id,v
	}

	return 0,nil
}
/////////////////////////////////////////






//////////////////// сохранить образы сочетаний ответных действий
//В случае отсуствия ответных действий создается ID такого отсутсвия, пример такой записи: 2|||0|0|
// ID|ActID через ,|PhraseID через ,|ToneID|MoodID|
func SaveActionsImageArr(){
	var out=""
	for k, v := range ActionsImageArr {
		out+=strconv.Itoa(k)+"|"
		for i := 0; i < len(v.ActID); i++ {
			out+=strconv.Itoa(v.ActID[i])+","
		}
		out+="|"
		for i := 0; i < len(v.PhraseID); i++ {
			out+=strconv.Itoa(v.PhraseID[i])+","
		}
		out+="|"
		out+=strconv.Itoa(v.ToneID)+"|"
		out+=strconv.Itoa(v.MoodID)+"|"
		out+="\r\n"
	}
	lib.WriteFileContent(lib.GetMainPathExeFile()+"/memory_psy/action_images.txt",out)

}
////////////////////  загрузить образы сочетаний ответных действий
// ID|ActID через ,|PhraseID через ,|ToneID|MoodID|
func loadActionsImageArr(){
	ActionsImageArr=make(map[int]*ActionsImage)
	strArr,_:=lib.ReadLines(lib.GetMainPathExeFile()+"/memory_psy/action_images.txt")
	cunt:=len(strArr)
	for n := 0; n < cunt; n++ {
		if len(strArr[n])==0{
			continue
		}
		p:=strings.Split(strArr[n], "|")
		id,_:=strconv.Atoi(p[0])

		s:=strings.Split(p[1], ",")
		var ActID[] int
		for i := 0; i < len(s); i++ {
			if len(s[i])==0{
				continue
			}
			si,_:=strconv.Atoi(s[i])
			ActID=append(ActID,si)
		}

		s=strings.Split(p[2], ",")
		var PhraseID[] int
		for i := 0; i < len(s); i++ {
			if len(s[i])==0{
				continue
			}
			si,_:=strconv.Atoi(s[i])
			PhraseID=append(PhraseID,si)
		}
		x,_:=strconv.Atoi(p[3])
		ToneID:=x
		x,_=strconv.Atoi(p[4])
		MoodID:=x
var saveDoWritingFile= doWritingFile; doWritingFile =false
		СreateNewlastActionsImageID(id,ActID,PhraseID,ToneID,MoodID)
doWritingFile =saveDoWritingFile
	}
	return

}
///////////////////////////////////////////

func GetActionsString(act int)(string){
	var out=""
	ai:=ActionsImageArr[act]
	if ai.ActID != nil {
		out += "Действие: "
		for i := 0; i < len(ai.ActID); i++ {
			if i > 0 {
				out += ", "
			}
			actName:= termineteAction.TerminalActonsNameFromID[ai.ActID[i]]
			out += "<b>" + actName + "</b>"
		}
		out += " "
	}

	if ai.PhraseID != nil {
		out += "Фраза: "
		for i := 0; i < len(ai.PhraseID); i++ {
			if i > 0 {
				out += " "
			}
			prase := word_sensor.GetPhraseStringsFromPhraseID(ai.PhraseID[i])
			out += "<b>\""+prase+"\"</b>"
		}
		out += " "
	}

	if ai.ToneID != 0 {
		out+=" Тон: "+getToneStrFromID(ai.ToneID)+" "
	}

	if ai.MoodID != 0 {
		out+=" Настрой: "+getMoodStrFromID(ai.MoodID)+"<br>"
	}

	return out
}
/////////////////////////////////////////







