/*  Ориентировочные рефлексы
Привлечение осознанного внимания выявляет конечную цель - найти автоматизм или ничего не делать.
Начинается с Определение Цели в данной ситуации - на уровне наследственных функций
и активных доминант нерешенной проблемы.

Новизна при активации Дерева автоматизмов CurrentAutomatizTreeEnd будет сохранена при выполнении автоматизма в
*/

package psychic

import (
	"BOT/lib"
)

////////////////////////////////////////
var NoveltySituation []int // НОВИЗНА СИТУАЦИИ сохраняет значение CurrentAutomatizTreeEnd[] для решений

/////////////////////////////////////////
/* получить инфу после активации дерева рефлексов
Активация Дерева рефлексов всегда оказывается раньше, чем активации дерева понимания
и здесь получаем информацию о результате.
*/
var actualRelextActon []int
var veryActualSituation=false
var curTargetArrID []int
func GetReflexInformation(veryActual bool,targetArrID  []int,acrArr []int){
	//! получить при активации древа!!!! veryActualSituation=veryActual
	actualRelextActon = acrArr
	//! получить при активации древа!!!!curTargetArrID=targetArrID
}
////////////////////////////////////////////


// пульс PulsCount
func orientarionPuls(){

	/*  если еще не запущен автоматизм  НЕ НУЖНО ВЫЗЫВАТЬ ВСЕ ВРЕМЯ!!!
	if LastRunAutomatizmPulsCount==0{//20 сек ожидания (if LastRunAutomatizmPulsCount+20 < PulsCount {)
		orientation(saveAutomatizmID)
		saveAutomatizmID=0
	}
	 */
}
////////////////////////////////////////////////////


/////////////////////////////////////////////////////
/*  Выполнение ориентировочного рефлекса из активной ветки Дерева автоматизмов.
automatizmID: 0 - в активной ветке нет автоматзма, >0 - есть автоматизм
 */
//var saveAutomatizmID=0

// вызывается из func afterTreeActivation()
func orientation(automatizmID int)(int){

	lib.WritePultConsol("Ориентировчный рефлекс Дерева моторных автоматизмов.")

	notAllowScanInTreeThisTime = true
//	saveAutomatizmID=automatizmID
		var atmtzm *Automatizm
		if automatizmID == 0 {
//автоматизма нет, если нужно действовать, то какой-то предположить и сразу проверить
			atmtzm = orientation_1()
		}
		if automatizmID > 0 {
//проверить подходит ли автоматизм defAutomatizmID к текущим условиям
			atmtzm = orientation_2(automatizmID)
		}
		if atmtzm != nil {
			atmtzm.BranchID = detectedActiveLastNodID
			notAllowScanInTreeThisTime = false
			return atmtzm.ID
		}
notAllowScanInTreeThisTime = false
return 0
}
///////////////////////////////////////////////////////



/* автоматизма нет, если нужно действовать, то какой-то предположить и сразу проверить
Стадия отсуствия опыта в данных условиях.
 */
func orientation_1()(*Automatizm){

	lib.WritePultConsol("Простейший ориентировочный рефлекс полного непонимания (1 типа)")

	// новизна ситуации
	NoveltySituation=CurrentAutomatizTreeEnd // значение сохраняется в savedNoveltySituation

//  получение текущего состояния информационной среды: отражение Базового состояния и Активных Базовых контекстов
	GetCurrentInformationEnvironment()

	// оценка опасности ситуации, необходиомсть срочных действий
	veryActualSituation=CurrentInformationEnvironment.veryActualSituation
	// выявить ID парамктров гомеостаза как цели для улучшения в данных условиях
	curTargetArrID=CurrentInformationEnvironment.curTargetArrID


	if EvolushnStage < 3 {// до формирования зеркальных !!!!!
		/* Определение Цели в данной ситуации - ну уровне наследственных функций
		Здесь выбирается действие пробного автоматизма из выполнившегося рефлекса actualRelextActon
		и запускается автоматизм
		*/
		atmzm:=getPurposeGeneticAndRunAutomatizm()// в purpose_genetic.go
//  ЗДЕСЬ активировать Дерево Понимания НЕ НУЖНО, если действие уже запущено, омысление будет по результату.
		return atmzm
	}

	if EvolushnStage == 3 { // формирование зеркальных
/*если нет автоматизма - просто повторить фразу оператора как попугай.
		Он на нее ответит - и сработает функция отзеркаливания,
		сформируется автоматизм, как отвечать на фразу.
		Думается, это вполне могло оказаться эволюционной находкой:
		не знаешь как реагировать, подай сигнал и повтори действие родителя -
		он покажет, что надо сделать. И тупо рефлекторно запишется.
 */
		if WasOperatorActiveted { // оператор отреагировал
			purpose := getPurposeGenetic()
			// повторить действия оператора
			purpose.actionID=curActiveActions
			atmzm := createAndRunAutomatizmFromPurpose(purpose)
			if doWritingFile {SaveAutomatizm() }
			return atmzm
		}
	}

	// else НИЧЕГО НЕ ДЕЛАТЬ: при высокой актуальности - растерянность, при низкой - лень

	return nil
}
///////////////////////////////////////



////////////////// ОРИЕНТИРОВОКА, если есть автоматизм - ВЫЗЫВАТЕСЯ ВСЕГДА, не только при новых условиях
/*проверить подходит ли автоматизм defAutomatizmID к текущим условиям, если нет,
- по опыту того, к чему приводят новые условия - режим нахождения альтернативы
Или если автоматизма пока не имеет Belief==2, т.е. еще непроверненный

! важно: если вернуло автоматизм, значит хочет попробовать
 */
func orientation_2(nodeAutomatizmID int)(*Automatizm){
	lib.WritePultConsol("Простейший ориентировочный рефлекс частичного непонимания (2 типа)")

//  получение текущего состояния информационной среды: отражение Базового состояния и Активных Базовых контекстов
	GetCurrentInformationEnvironment()

	// оценка опасности ситуации, необходиомсть срочных действий
	veryActualSituation=CurrentInformationEnvironment.veryActualSituation
	// выявить ID парамктров гомеостаза как цели для улучшения в данных условиях
	curTargetArrID=CurrentInformationEnvironment.curTargetArrID

	// обработка автоматизма, рвущегося на выполнение. Есть ли опасная новизна?
	atmzm:=getPurposeGenetic2AndRunAutomatizm(nodeAutomatizmID)// в purpose_genetic.go
	if atmzm != nil {
		return atmzm
	}

return nil
}
//////////////////////////////////////////////











