/*   дейсвтвия ментального автоматизма

*/

package psychic

import (

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
	return RumMentalMentalAutomatizm(a)
}
////////////////////
// todo = true - выполнить полюбому,
func RumMentalMentalAutomatizm(am *MentalAutomatizm)(bool){
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
if am.Usefulness<0{

	return false
}

	GetMentalAutomotizmActionsString(am,true)

	currentMentalAutomatizmID=am.ID

	//выполнить мозжечковый рефлекс сразу после выполняющегося автоматизма
	//runCerebellumAdditionalMentalAutomatizm(0,am.ID)

//	notAllowReflexRuning=true // блокировка рефлексов

	return true
}
//////////////////////////////////////////



func GetMentalAutomotizmActionsString(am *MentalAutomatizm,writeLog bool){


	// TODO по MentalActionsImages автоматизма



}
/////////////////////////////////////////






