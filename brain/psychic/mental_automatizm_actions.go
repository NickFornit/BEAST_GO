/*   дейсвтвия ментального автоматизма

*/

package psychic

import (
	"BOT/brain/gomeostas"
)
////////////////////////////////////////////////



/* для разделения строки Sequence автоматизма на составляющие
типы действий:
1 Mnn - выполнить ментальную функцию с ID
2 Ann - выполнить моторный автоматизм с ID
*/
type ActsMentalAutomatizm struct {
	Type int  	// тип совершаемого действия
	Acts string // само действие
}


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

// блокировка выполнения плохого автоматизма, если только не применена "СИЛА ВОЛИ"
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
		// обязательная переактивация
		consciousness(2,3)
	}

	if ai.activateBaseID >0 {// на один текущий пульс, во время которого происходит обдумывание
		gomeostas.CommonBadNormalWell=ai.activateBaseID
		automatizmTreeActivation()
		// и сама последует consciousness(1,0)
	}
	if ai.activateEmotion >0 {// на один текущий пульс, во время которого происходит обдумывание

		// найти эмоцию по ее ID
		lev2:=EmotionFromIdArr[ai.activateEmotion].BaseIDarr
		gomeostas.SetCurContextActiveIDarr(lev2)
		automatizmTreeActivation()
		// и сама последует consciousness(1,0)
	}

	currentMentalAutomatizmID=am.ID
	//выполнить мозжечковый рефлекс сразу после выполняющегося автоматизма
	//runCerebellumAdditionalMentalAutomatizm(0,am.ID)
//	notAllowReflexRuning=true // блокировка рефлексов
	return true
}
//////////////////////////////////////////


// только для показа на Пульте т.к. мент.автоматизм не имеет видимых действий 
func GetMentalAutomotizmActionsString(am *MentalAutomatizm,writeLog bool){


	// TODO по MentalActionsImages автоматизма



}
/////////////////////////////////////////






