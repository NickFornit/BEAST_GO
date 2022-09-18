/*  Ориентировочные рефлексы
Привлечение осознанного внимания выявляет конечную цель - найти автоматизм или ничего не делать.
Начинается с Определение Цели в данной ситуации - на уровне наследственных функций
и активных доминант нерешенной проблемы.

Новизна при активации Дерева автоматизмов CurrentAutomatizTreeEnd будет сохранена при выполнении автоматизма в
*/

package psychic

import (
	"BOT/lib"
	"BOT/brain/gomeostas"
)

////////////////////////////////////////
var NoveltySituation []int // НОВИЗНА СИТУАЦИИ сохраняет значение CurrentAutomatizTreeEnd[] для решений
var isPeriodResultWaiting=false // идет период ожидания резйльтата автоматизма

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
var needOrientationReflex=false
var wasRunOrientationReflex=false
func orientarionPuls(){
	// как только отпустит AutomatizmRunningPulsCount - сразу выполнить
	if wasRunOrientationReflex && AutomatizmRunningPulsCount==0{
		orientation(saveAutomatizmID)
		saveAutomatizmID=0
	}
}
////////////////////////////////////////////////////


/////////////////////////////////////////////////////
/*  Выполнение ориентировочного рефлекса из активной ветки Дерева автоматизмов.
automatizmID: 0 - в активной ветке нет автоматзма, >0 - есть автоматизм
 */
var saveAutomatizmID=0
func orientation(automatizmID int)(int){
	wasRunOrientationReflex=true
	saveAutomatizmID=automatizmID
	if AutomatizmRunningPulsCount >0{// идет время ожидания результата выполненного автоматизма
		needOrientationReflex=true
	}else {// Нет ожидания результата выполненного автоматизма
		wasRunOrientationReflex=false // покасить сразу
		var atmtzm *Automatizm
		if automatizmID == 0 {
//автоматизма нет, если нужно действовать, то какой-то предположить и сразу проверить
			atmtzm = orientation_1()
		}
		if saveAutomatizmID > 0 {
//проверить подходит ли автоматизм defAutomatizmID к текущим условиям
			atmtzm = orientation_2(saveAutomatizmID)
		}
		if atmtzm != nil {
			atmtzm.BranchID = detectedActiveLastNodID
			notAllowScanInTreeThisTime = false
			return atmtzm.ID
		}
	}

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

	// Определение Цели в данной ситуации - ну уровне наследственных функций
	purpose:=getPurposeGenetic()
	// мозжечковые рефлексы - самый первый уровень осознания - подгонка действий под заданную Цель.

	// нужно ли вообще шевелиться?
	// veryActualSituation: плохо для  1, 2, 7 и/или 8  параметров гомеостаза
	if purpose.veryActual{// нужно создавать автоматизм и тут же запускать его
		if purpose.actionID.ID>0 {
			/* сформировать пробный автоматизм моторного действия и сразу запустить его в действие
			   Зафиксироваь время действия
			   10 пульсов следить за измнением состояния и ответными действиями - считать следствием действия
			   оценить результат и скорректировать силу мозжечком в записи автоматизма.

			Выбрать любое действие, т.к. послед создания автоматизма в данной ветке detectedActiveLastNodID
			он уже не вызовет orientation_1(), а будет orientation_2()
			*/
			atmzm:=createAutomatizm(purpose)
			// запустить автоматизм
			if RumAutomatizm(atmzm) {
				// отслеживать последствия в automatizm_result.go
				setAutomatizmRunning(atmzm,purpose)
			}
// в automatizm_result.go после оценки результата будет осмысление с активацией Дерева Понимания

			return atmzm
		}
		///////////////////////////////////////
		// игнорировать новое внимание на время ожидания результата автоматизма
		if isPeriodResultWaiting {
			return AutomatizmRunning
		}

// нет действий, попробовать использовать AutomatizmSuccessFromIdArr.GomeoIdSuccesArr
// TODO .....................

// нет действий, попробовать бессмысленно выдать фразы имеющиеся Вериниковские раз нужно что-то срочно делать
			if purpose.veryActual {
				// TODO подобрать хоть как-то ассоциирующуюся фразу из имеющизся
				phraseID:=findSuitablePhrase()
				if len(phraseID)>0{
					purpose.actionID.PhraseID=phraseID
					atmzm:=createAutomatizm(purpose)
					// запустить автоматизм
					if RumAutomatizm(atmzm) {
						// отслеживать последствия в automatizm_result.go
						setAutomatizmRunning(atmzm,purpose)
					}
					// в automatizm_result.go после оценки результата будет осмысление с активацией Дерева Понимания
					return atmzm
				}


				// return autmzm
			}
		//  ЗДЕСЬ активировать Дерево Понимания НЕ НУЖНО, действие уже запущено, омысление будет по результату.

	}else{// нет атаса, можно спокойно поэкспериментивроать, если есть любопытсво
		if gomeostas.BaseContextActive[2] {// активен Поиск или Игра

			// нет действий, попробовать использовать AutomatizmSuccessFromIdArr.GomeoIdSuccesArr
			// TODO .....................

// Тупо поэкспериментировать в контексте поиска или игры для пополнения опыта (не)удачных автоматизмов
// TODO ................................

			/*
					// осмыслить ситуацию - Активировать Дерево Понимания
				autmzmTreeNodeID:=AutomatizmRunning.BranchID// создать образ ситуации
				id,_:=createSituationImage(0,autmzmTreeNodeID,4)
				// осмыслить ситуацию - Активировать Дерево Понимания
				understandingSituation(id,purpose)
*/
			//return automatizm
		}
	}
	isReflexesActionBloking=false
	return nil
}
///////////////////////////////////////



////////////////// ОРИЕНТИРОВОКА в условиях, когда есть автоматизм
/*проверить подходит ли автоматизм defAutomatizmID к текущим условиям, если нет,
- по опыту того, к чему приводят новые условия - режим нахождения альтернативы
Или если автоматизма пока не имеет Belief==2, т.е. еще непроверненный

! важно: если вернуло автоматизм, значит хочет попробовать
 */
func orientation_2(nodeAutomatizmID int)(*Automatizm){
	lib.WritePultConsol("Ориентировочный рефлекс частичного непонимания (2 типа)")

	NoveltySituation=CurrentAutomatizTreeEnd

//  получение текущего состояния информационной среды: отражение Базового состояния и Активных Базовых контекстов
// только при ориентировчном рефлексе - обновление самоощущения! и запись эпизодической памяти
	GetCurrentInformationEnvironment()

	// оценка опасности ситуации, необходиомсть срочных действий
	veryActualSituation=CurrentInformationEnvironment.veryActualSituation
	// выявить ID парамктров гомеостаза как цели для улучшения в данных условиях
	curTargetArrID=CurrentInformationEnvironment.curTargetArrID

	/* Блокировать выполнение рефлексов на время ожидания результата автоматизма
	вызывается из reflex_action.go рефлексов
	*/
	isReflexesActionBloking=true // отмена в automatizm_result.go или просто isReflexesActionBloking=false


	atmzm:=AutomatizmFromIdArr[nodeAutomatizmID]

	// Определение Цели в данной ситуации - на уровне наследственных функций
	purpose:=getPurposeGenetic()
// мозжечковые рефлексы - самый первый уровень осознания - подгонка действий под заданную Цель.

	// если непроверенный автоматизм
	if  AutomatizmFromIdArr[nodeAutomatizmID].Belief!=2{

		// TODO

	}

	// ПЛОХОЙ АВТОМАТИЗМ ?!
	if atmzm.Usefulness < 0{ // плохой автоматизм, т.е. в прошлый раз был плохой результат
		// осмыслить ситуацию - Активировать Дерево Понимания
		autmzmTreeNodeID:=AutomatizmRunning.BranchID// создать образ ситуации
		id,_:=createSituationImage(0,autmzmTreeNodeID,3)
		// осмыслить ситуацию - Активировать Дерево Понимания
		understandingSituation(id,purpose)

		// если найден лучше автоматизм, то подставить его: atmzm=xxxx и запустить

		//return automatizm
	}else {
		// АКТИВАЦИЯ ДЕРЕВА ПОНИМАНИЯ СРАЗУ (и потом при получении результата реагирования)
		// создать образ ситуации
		autmzmTreeNodeID := AutomatizmRunning.BranchID
		id, _ := createSituationImage(0, autmzmTreeNodeID, 1)
		// осмыслить ситуацию - Активировать Дерево Понимания
		understandingSituation(id, purpose)
	}


	if purpose.veryActual{// нужно ли вообще шевелиться?
		// срочность и важность ситуации: если очень срочно и важно - просто оставить имеющийся автоматизм
		atmzm:=createAutomatizm(purpose)
		// запустить автоматизм
		if RumAutomatizm(atmzm) {
			// отслеживать последствия в automatizm_result.go ИГНОРИРОВАТЬ ДРУГИЕ ПОКА НЕ БУДЕТ РЕЗУЛЬТАТ
			setAutomatizmRunning(atmzm,purpose)
		}

		// в automatizm_result.go после оценки результата будет осмысление с активацией Дерева Понимания
		return atmzm
	}
	////////////////////////////
	// игнорировать новое внимание на время ожидания результата автоматизма
	if isPeriodResultWaiting {
		return AutomatizmRunning
	}

	// НЕТ СРОЧНОСТИ, МОЖНО СПОКОЙНО ПОДУМАТЬ или тупо поэкспериментировать
	/* найти важные (по опыту) признаки в новизне NoveltySituation
	   Это - чисто рефлексторный процесс поиска в опыте
	*/
	is:=getImportantSigns()
	if is!=nil {
		/*
				// осмыслить ситуацию - Активировать Дерево Понимания
			autmzmTreeNodeID:=AutomatizmRunning.BranchID// создать образ ситуации
			id,_:=createSituationImage(0,autmzmTreeNodeID,5)
			// осмыслить ситуацию - Активировать Дерево Понимания
			understandingSituation(id,purpose)
		*/
	}

	// если не найдено решение
	/*
			// осмыслить ситуацию - Активировать Дерево Понимания
		autmzmTreeNodeID:=AutomatizmRunning.BranchID// создать образ ситуации
		id,_:=createSituationImage(0,autmzmTreeNodeID,1)
		// осмыслить ситуацию - Активировать Дерево Понимания
		understandingSituation(id,purpose)
	*/

	// Тупо поэкспериментировать в контексте поиска или игры для пополнения опыта (не)удачных автоматизмов

	// TODO ................................


	////////////  если нет результата
	isReflexesActionBloking=false
	return nil
}
//////////////////////////////////////////////






