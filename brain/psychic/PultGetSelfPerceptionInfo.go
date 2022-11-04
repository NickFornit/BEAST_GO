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

	refreshCurrentInformationEnvironment()
	// опасность
	danger:=GetAttentionDanger()
	// актуальность ситуации
	veryActualSituation,_:=gomeostas.FindTargetGomeostazID()

	// против паники типа "одновременная запись и считывание карты"
	if notAllowScanInTreeThisTime{
		return "!!!"
	}
	ie:=CurrentInformationEnvironment
var out="Общее состояние жизненных параметров: <b>"
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
	out+="<br>Текущая эмоция: "+GetCurrentEmotionReception() //getEmotonsComponentStr(ie.PsyEmotionImg)
	///////////////////////////////
	out+="<br>Опасность состояния: <b>"
	if danger {out+="Опасное"}else{out+="Неопасное"}
	out+="</b>"
	out+="<br>Важность состояния: <b>"
	if veryActualSituation {out+="Очень важное состояние"}else{out+="Спокойное состояние."}
	out+="</b>"
	////////////////////////////////
	if ie.PsyActionImg !=nil {
		str := ""
		aID := ie.PsyActionImg.ActID
		for i := 0; i < len(aID); i++ {
			if i > 0 {
				out += ", "
			}
			str += action_sensor.GetActionNameFromID(aID[i])
		}
		if len(str) > 0 {
			out += "<br>Текущий образ сочетания действий с Пульта: <b>" + str + "</b>"
		}
	}
	////////////////////////////////
	if ie.PsyVerbImg !=nil{
		str:=""
		pID:=ie.PsyVerbImg.PhraseID
		for i := 0; i < len(pID); i++ {
			if i>0{out+=", "}
			str+=word_sensor.GetPhraseStringsFromPhraseID(pID[i])
		}
if len(str) >0 {
	out += "<br>Текущий образ фразы с Пульта: <b>"+str+"</b>"
}
	}
	///////////////////////////////


	// TODO остальное


return out
}