/*   дейсвтвия ментального автоматизма

*/

package psychic

import (
	"BOT/lib"
	"strconv"
)
////////////////////////////////////////////////


// момент запуска автоматизма в числе пульсов
var LastRunMentalMentalAutomatizmPulsCount =0 //сбрасывать ожидание результата автоматизма если прошло 20 пульсов
// ожидается результат запущенного MotMentalAutomatizm
var LastMentalMentalAutomatizmWeiting *MentalAutomatizm
// активный узел дерева в момент запуска автоматизма
var LastDetectedActiveLastMentalNodID=0

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
//	if MotorTerminalBlocking { //блокировка моторных терминалов во сне или произвольно
//		return false
//	}

// NotAllowAnyActions ставится тогда, когда сохранение памяти должно выполняться в тишине, в бездействии
	if  NotAllowAnyActions{
		return false
	}
	if len(am.Sequence)==0{
		return false
	}

// блокировка выполнения плохого автоматизма, если только не применена "ИЛА ВОЛИ"
if am.Usefulness<0{

	return false
}

	var out="3|"
	actArr:=ParceMentalAutomatizmSequence(am.Sequence)
	for i := 0; i < len(actArr); i++ {
		act,_:=strconv.Atoi(actArr[i].Acts)

		switch actArr[i].Type{
		case 1: // выполнить ментальную функцию с ID
			runMenyalFunctionID(act)

		case 2: //выполнить моторный автоматизм с ID
			RumAutomatizmID(act)

		///////////////////////////////////////
		}
	}
	lib.SentActionsForPult(out)

	//выполнить мозжечковый рефлекс сразу после выполняющегося автоматизма
	//runCerebellumAdditionalMentalAutomatizm(0,am.ID)

	notAllowReflexRuning=true // блокировка рефлексов
	LastRunMentalMentalAutomatizmPulsCount =PulsCount // активность мот.автоматизма в чисде пульсов
	LastMentalMentalAutomatizmWeiting=am
	LastDetectedActiveLastMentalNodID=detectedActiveLastNodID

	return true
}
//////////////////////////////////////////






