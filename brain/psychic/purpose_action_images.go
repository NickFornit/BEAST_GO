/*
Образы отвестных действий
*/

package psychic

import (
	"BOT/lib"
	"strconv"
	"strings"
)

///////////////////////////////////



/* Образ действия

 */
type ActionImage struct {
	ID int  // идентификатор данного сочетания действий
	ActID[] int // массив действий (http://go/pages/terminal_actions.php)
	// для текущего сообщения с Пусльта:
	PhraseID []int// массив фразID (DetectedUnicumPhraseID) слова каждой фразы вытаскиваются wordSensor.GetWordArrFromPhraseID(PhraseID[0])
	PhraseToneID int  // тон фразы
	PhraseMoodID int // настроение фразы
}
var ActionImageArr=make(map[int]*ActionImage)

var ActiveCurActionImageID=0 // ID последнего совершенного образа сдействия
//////////////////////////////////////////


////////////////////////////////////////////////
// создать новое сочетание ответных действий если такого еще нет с ЗАПОМИНАНИЕМ
func CreateNewActionImageImage(ActID []int,PhraseID []int,ToneID int,MoodID int)(int,*ActionImage){
	oldID,oldVal:=checkUnicumActionImage(ActID,PhraseID,ToneID,MoodID)
	if oldVal!=nil{
		return oldID,oldVal
	}
	aImgID,aImg:=createNewlastActionImageID(0,ActID,PhraseID,ToneID,MoodID)

	saveActionImageArr()

	return aImgID,aImg
}
/////////////////////////////////////////
// создать образ сочетаний ответных действий
var lastActionImageID=0
func createNewlastActionImageID(id int,ActID []int,PhraseID []int,ToneID int,MoodID int)(int,*ActionImage){

	if id==0{
		lastActionImageID++
		id=lastActionImageID
	}else{
		//		newW.ID=id
		if lastActionImageID<id{
			lastActionImageID=id
		}
	}

	var node ActionImage
	node.ID = id
	node.ActID = ActID
	node.PhraseID=PhraseID
	node.PhraseToneID=ToneID
	node.PhraseMoodID=MoodID

	ActionImageArr[id]=&node
	return id,&node
}
func checkUnicumActionImage(ActID []int,PhraseID []int,ToneID int,MoodID int)(int,*ActionImage){
	for id, v := range ActionImageArr {
		if !lib.EqualArrs(ActID,v.ActID) {
			continue
		}
		if !lib.EqualArrs(PhraseID,v.PhraseID) {
			continue
		}
		if ToneID!=v.PhraseToneID || MoodID!=v.PhraseMoodID {
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
func saveActionImageArr(){
	var out=""
	for k, v := range ActionImageArr {
		out+=strconv.Itoa(k)+"|"
		for i := 0; i < len(v.ActID); i++ {
			out+=strconv.Itoa(v.ActID[i])+","
		}
		out+="|"
		for i := 0; i < len(v.PhraseID); i++ {
			out+=strconv.Itoa(v.PhraseID[i])+","
		}
		out+="|"
		out+=strconv.Itoa(v.PhraseToneID)+"|"
		out+=strconv.Itoa(v.PhraseMoodID)+"|"
		out+="\r\n"
	}
	lib.WriteFileContent(lib.GetMainPathExeFile()+"/memory_reflex/action_images.txt",out)

}
////////////////////  загрузить образы сочетаний ответных действий
// ID|ActID через ,|PhraseID через ,|ToneID|MoodID|
func loadActionImageArr(){
	ActionImageArr=make(map[int]*ActionImage)
	strArr,_:=lib.ReadLines(lib.GetMainPathExeFile()+"/memory_reflex/action_images.txt")
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

		createNewlastActionImageID(id,ActID,PhraseID,ToneID,MoodID)
	}
	return

}
///////////////////////////////////////////



//////// из действий автоматизма создать новое сочетание пусковых стимулов TriggerStimuls если такого еще нет
func CreateNewActionImageFromAutomatizm(atmzmAct string)(*ActionImage){
	var aArr []int
	var sArr []int
	actArr:=ParceAutomatizmSequence(atmzmAct)
	for i := 0; i < len(actArr); i++ {
		if actArr[i].Type == 1 && len(actArr[i].Acts)>0 { // Snn- перечень ID сенсора слов через запятую
			p:=strings.Split(actArr[i].Acts, ",")
			for n := 0; n < len(p); n++ {
				aID, _ := strconv.Atoi(p[n])
				sArr = append(sArr, aID)
			}
		}
		if actArr[i].Type == 2  && len(actArr[i].Acts)>0 { //Dnn - ID прогрмаммы действий, через запятую
			p:=strings.Split(actArr[i].Acts, ",")
			for n := 0; n < len(p); n++ {
				aID, _ := strconv.Atoi(p[n])
				aArr = append(aArr, aID)
			}
		}
	}
		_,trigger:=CreateNewActionImageImage(aArr,sArr,0,0)

	return trigger
}
///////////////////////////////////////////////////////