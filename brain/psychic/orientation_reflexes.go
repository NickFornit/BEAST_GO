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

// вызывается из func automatizmTreeActivation()
func orientation(automatizmID int)(int){
	if LastRunAutomatizmPulsCount>0{// не перебивать ожидание от запущенного автоматизма
/* НО! во время периода ожидания реакции оператора на действие автоматизма
	более важная ориентировочная реакция может превать ожидание ответа (стек прерываний)

	TODO запомнить какой автоматизма прерван и потом вспомнить чтобы снова запустить и как-то обработать.
		возможно, создать доминанту если это - на достаточной стадии развития
 */

		return 0
	}
	// Нет ожидания результата выполненного автоматизма
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

	lib.WritePultConsol("Ориентировочный рефлекс полного непонимания (1 типа)")

	NoveltySituation=CurrentAutomatizTreeEnd // значение сохраняется в savedNoveltySituation

//  получение текущего состояния информационной среды: отражение Базового состояния и Активных Базовых контекстов
// только при ориентировчном рефлексе - обновление самоощущения! и запись кадра эпизодической памяти
	GetCurrentInformationEnvironment()

	// оценка опасности ситуации, необходиомсть срочных действий
	veryActualSituation=CurrentInformationEnvironment.veryActualSituation
	// выявить ID парамктров гомеостаза как цели для улучшения в данных условиях
	curTargetArrID=CurrentInformationEnvironment.curTargetArrID

	/* Блокировать выполнение рефлексов на время ожидания результата автоматизма
	вызывается из reflex_action.go рефлексов
	*/
	isReflexesActionBloking=true // отмена в automatizm_result.go или просто isReflexesActionBloking=false

	if EvolushnStage < 4 {
		/* Определение Цели в данной ситуации - ну уровне наследственных функций
		Здесь выбирается действие пробного автоматизма из выполнившегося рефлекса actualRelextActon
		и запускается автоматизм
		*/
		atmzm:=getPurposeGeneticAndRunAutomatizm()// в purpose_genetic.go
//  ЗДЕСЬ активировать Дерево Понимания НЕ НУЖНО, если действие уже запущено, омысление будет по результату.
		return atmzm
	}

	if EvolushnStage > 3 {
		/* Определение Цели в данной ситуации - ну уровне дерева понимания
		Здесь выбирается действие пробного автоматизма из выполнившегося рефлекса actualRelextActon
		и запускается автоматизм.
		На стадии 4 - провоцировать оператора на ответы (почему, зачем, что такое?)
		*/
		atmzm:=getPurposeUndestandingAndRunAutomatizm() //в understanding_purpose_image.go
		return atmzm
	}
	// else НИЧЕГО НЕ ДЕЛАТЬ: при высокой актуальности - растерянность, при низкой - лень


	isReflexesActionBloking=false
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
	lib.WritePultConsol("Ориентировочный рефлекс частичного непонимания (2 типа)")

//  получение текущего состояния информационной среды: отражение Базового состояния и Активных Базовых контекстов
// только при ориентировчном рефлексе - обновление самоощущения! и запись эпизодической памяти
	GetCurrentInformationEnvironment()

	// оценка опасности ситуации, необходиомсть срочных действий
	veryActualSituation=CurrentInformationEnvironment.veryActualSituation
	// выявить ID парамктров гомеостаза как цели для улучшения в данных условиях
	curTargetArrID=CurrentInformationEnvironment.curTargetArrID

	if EvolushnStage < 4 {
		/* // обработка автоматизма, рвущегося на выполнение, но в условиях есть новизна news
		Если опасной новизны нет, то
		*/
		atmzm:=getPurposeGenetic2AndRunAutomatizm(nodeAutomatizmID)// в purpose_genetic.go
		//  ЗДЕСЬ активировать Дерево Понимания НЕ НУЖНО, если действие уже запущено, омысление будет по результату.
		return atmzm
	}
	if EvolushnStage > 3 {
		// TODO - осмысление
		atmzm:=getPurposeUndestanding2AndRunAutomatizm(nodeAutomatizmID) //в understanding_purpose_image.go
		return atmzm
	}

return nil
}
//////////////////////////////////////////////











