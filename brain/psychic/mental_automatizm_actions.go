/*   дейсвтвия ментального автоматизма

*/

package psychic

import (
	"strconv"
)
////////////////////////////////////////////////

// момент запуска автоматизма в числе пульсов
var LastRunMentalAutomatizmPulsCount =0

/////////////////////////////////////////
/* запуск автоматизма на выполнение
возвращает true при успехе
 */
func RunMentalAutomatizmsID(id int)(bool){
	a:=MentalAutomatizmsFromID[id]
	if a==nil{
		return false
	}
	return RunMentalAutomatizm(a)
}
////////////////////
// Запуск
func RunMentalAutomatizm(am *MentalAutomatizm)(bool){
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

/* Блокировка выполнения плохого мент.автоматизма, если только не применена "СИЛА ВОЛИ"
	выполняется на уровне ментальных Правил, в функции consciousness
if am.Usefulness<0{
	return false Здесь не должно быть блокировки!
}
*/
	ai:=MentalActionsImagesArr[am.ActionsImageID]

	switch ai.typeID{
	case 1: //активация настроения Mood
		mentalMoodVolitionID=ai.valID
		mentalMoodVolitionPulsCount=PulsCount
			//gomeostas.CommonBadNormalWell=ai.valID  не трогать гомеостаз!
			understandingSituation(2)
	case 2: // активация эмоции EmotionID
			em:=EmotionFromIdArr[ai.valID]
			if em==nil{return false}
		mentalEmotionVolitionID=ai.valID
		mentalEmotionVolitionPulsCount=PulsCount
		/*    не трогать гомеостаз!
			lev2:=em.BaseIDarr
			gomeostas.SetCurContextActiveIDarr(lev2)
		 */
			understandingSituation(2)
	case 3:	//активация PurposeImage
			pp:=PurposeImageFromID[ai.valID]
			if pp==nil{return false}
		mentalPurposeImageID=ai.valID
		mentalPurposeImagePulsCount=PulsCount
			understandingSituation(2)
	case 4: runMentalFunctionID(ai.valID) //запуск инфо-функции
	case 5: 
	saveFromNextIDAnswerCicle=nil // обнуление памяти между циклами
	RumAutomatizmID(ai.valID) //запуск моторного автоматизма
	//case 6: //запуск Доминанты
	//case 7: // создание новой Доминанты

	default: return false
	}

	am.Count++
	currentMentalAutomatizmID=am.ID
	LastRunMentalAutomatizmPulsCount =PulsCount // активность мот.автоматизма в чисде пульсов
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

	switch ai.typeID{
	case 1: //активация настроения Mood
		val:=""
		switch ai.valID{
		case 1: val="Плохо"
		case 2: val="Норма"
		case 3: val="Хорошо"
		}
		if len(out)>0{out += "<hr>"}
		out += "<span style='font-size:19px;'>Переактивация Базового состояния Дерева понимания на :</span> " + val
	case 2: // активация эмоции EmotionID
		if len(out)>0{out += "<hr>"}
		em:=EmotionFromIdArr[ai.valID]
		emS:=getEmotonsComponentStr(em)
		out += "<span style='font-size:19px;'>Переактивация Эмоции на :</span> " + strconv.Itoa(ai.valID)+"<br>"+emS
	case 3:	//активация PurposeImage
		pp:=PurposeImageFromID[ai.valID]
		if pp==nil{return "Нет ментального образа действия с ID = "+strconv.Itoa(ai.valID)}
		out += "<span style='font-size:19px;'>Переактивация ментального образа действия:</span><br> " + getPurposeDetaileString(ai.valID)
	case 4: //запуск инфо-функции
		if len(out)>0{out += "<hr>"}
		out+="Запуск информационной функции № "+strconv.Itoa(ai.valID)
	case 5: //запуск моторного автоматизма
		motA:=AutomatizmFromIdArr[ai.valID]
		if motA != nil {
			if len(out)>0{out += "<hr>"}
			out += "<span style='font-size:19px;'>Запуск моторного автоматизма:</span><br> " + GetAutomotizmActionsString(motA, writeLog)
		}
	case 6: //запуск Доминанты
		out+="Запуск  Доминанты № "+strconv.Itoa(ai.valID)
	case 7: // создание новой Доминанты
		out+="создание новой Доминанты"

	default: out+="Нет ментального действия с ID = "+strconv.Itoa(ai.typeID)
	}

return out
}
/////////////////////////////////////////






