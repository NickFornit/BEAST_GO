/*   дейсвтвия ментального автоматизма

*/

package psychic

import (
	"BOT/brain/gomeostas"
	"strconv"
)
////////////////////////////////////////////////



/////////////////////////////////////////
/* запуск автоматизма на выполнение
возвращает true при успехе
 */
func RunMentalMentalAutomatizmsID(id int)(bool){
	a:=MentalAutomatizmsFromID[id]
	if a==nil{
		return false
	}
	return RunMentalMentalAutomatizm(a)
}
////////////////////
// todo = true - выполнить полюбому,
func RunMentalMentalAutomatizm(am *MentalAutomatizm)(bool){
	if am==nil{
		return false
	}

// NotAllowAnyActions ставится тогда, когда сохранение памяти должно выполняться в тишине, в бездействии
	if  NotAllowAnyActions{
		return false
	}
	if am.ActionsImageID==0{
		return false
	}

// блокировка выполнения плохого мент.автоматизма, если только не применена "СИЛА ВОЛИ"
if am.ActionsImageID>0{

	return false
}
	ai:=MentalActionsImagesArr[am.ActionsImageID]

	if ai.activateMotorID >0 {
		// здесь начинается период ожидания: LastRunAutomatizmPulsCount =PulsCount
		RumAutomatizmID(ai.activateMotorID)
	}

	if ai.activateInfoFunc >0 {
		runMenyalFunctionID(ai.activateInfoFunc)
	}

	if ai.activateBaseID >0 {// на один текущий пульс, во время которого происходит обдумывание
		gomeostas.CommonBadNormalWell=ai.activateBaseID
		understandingSituation()
	}
	if ai.activateEmotion >0 {// на один текущий пульс, во время которого происходит обдумывание

		// найти эмоцию по ее ID
		lev2:=EmotionFromIdArr[ai.activateEmotion].BaseIDarr
		gomeostas.SetCurContextActiveIDarr(lev2)
		understandingSituation()
	}

	currentMentalAutomatizmID=am.ID
	//выполнить мозжечковый рефлекс сразу после выполняющегося автоматизма
	//runCerebellumAdditionalMentalAutomatizm(0,am.ID)
//	notAllowReflexRuning=true // блокировка рефлексов
	return true
}
//////////////////////////////////////////


/* Действия автоматизма MentalActionsImages
только для показа на Пульте т.к. мент.автоматизм не имеет видимых действий
 */
func GetMentalAutomotizmActionsString(id int,writeLog bool)(string){
	if id==0{
		return ""
	}
	out:=""
	am:=MentalAutomatizmsFromID[id]
	if am==nil{
		return ""
	}

	ai:=MentalActionsImagesArr[am.ActionsImageID]

	if ai.activateMotorID >0 {
		motA:=AutomatizmFromIdArr[ai.activateMotorID]
		if motA != nil {
			if len(out)>0{out += "<hr>"}
			out += "<span style='font-size:19px;'>Запуск моторного автоматизма:</span><br> " + GetAutomotizmActionsString(motA, writeLog)
		}
	}

	if ai.activateInfoFunc >0 {
		if len(out)>0{out += "<hr>"}
		out+="Запуск информационной функции № "+strconv.Itoa(ai.activateInfoFunc)
	}

	if ai.activateBaseID >0 {
		val:=""
		switch ai.activateBaseID{
		case 1: val="Плохо"
		case 2: val="Норма"
		case 3: val="Хорошо"
		}
		if len(out)>0{out += "<hr>"}
		out += "<span style='font-size:19px;'>Переактивация Базового состояния Дерева понимания на :</span> " + val

	}
	if ai.activateEmotion >0 {
		if len(out)>0{out += "<hr>"}
		em:=EmotionFromIdArr[ai.activateEmotion]
		emS:=getEmotonsComponentStr(em)
		out += "<span style='font-size:19px;'>Переактивация Эмоции на :</span> " + strconv.Itoa(ai.activateEmotion)+"<br>"+emS
	}

return out
}
/////////////////////////////////////////






