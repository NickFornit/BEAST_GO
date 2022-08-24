/* Для Пульта выдать текущее состояние Самоощещения.

*/

package psychic

import (
	"BOT/brain/action_sensor"
	"BOT/brain/gomeostas"
	"BOT/brain/words_sensor"
	"strconv"
)
///////////////////////////////////////


func GetSelfPerceptionInfo()(string){
	ie:=CurrentInformationEnvironment
var out="Общее базовое состояние: <b>"
if gomeostas.CommonBadNormalWell==1{
	out+="ПЛОХО"
}
	if gomeostas.CommonBadNormalWell==2{
		out+="НОРМА"
	}
	if gomeostas.CommonBadNormalWell==3{
		out+="ХОРОШО"
	}
	out+="</b>"
	////////////////////////////
	out+="<br>Ощущаемое настроение: <b>"
	if ie.PsyMood==-1{
		out+="Плохое ("+strconv.Itoa(ie.Mood)+")"
	}
	if ie.PsyMood==0{
		out+="Нормальное ("+strconv.Itoa(ie.Mood)+")"
	}
	if ie.PsyMood==1{
		out+="Хорошее ("+strconv.Itoa(ie.Mood)+")"
	}
	out+="</b>"
	////////////////////////////////
	out+="<br>Текущая эмоция: "+getEmotonsComponentStr(ie.PsyEmotionImg)
	///////////////////////////////
	out+="<br>Опасность состояния: <b>"
	if ie.danger{out+="Опасное"}else{out+="Неопасное"}
	out+="</b>"
	////////////////////////////////
	if ie.PsyActionImg !=nil{
		out+="<br>Текущий образ сочетания действий с Пульта: <b>"
		aID:=ie.PsyActionImg.ActID
		for i := 0; i < len(aID); i++ {
			if i>0{out+=", "}
			out+=action_sensor.GetActionNameFromID(aID[i])
		}
		out+="</b>"
	}
	////////////////////////////////
	if ie.PsyVerbImg !=nil{
		out+="<br>Текущий образ фразы с Пульта: <b>"
		pID:=ie.PsyVerbImg.PhraseID
		for i := 0; i < len(pID); i++ {
			if i>0{out+=", "}
			out+=word_sensor.GetPhraseStringsFromPhraseID(pID[i])
		}
		out+="</b>"
	}
	///////////////////////////////


	// TODO остальное


return out
}