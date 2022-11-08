/* Процессы осмысления: создание и использование ментальных автоматизмов
для Дерева понимания (или дерева ментальных автоматизмов)

*/

package psychic

import (
	"BOT/lib"
)

///////////////////////////////


//////////////////////////////////////////////////
/* Определение Цели в данной ситуации - на уровне дерева понимания
	Здесь выбирается действие пробного автоматизма из выполнившегося рефлекса actualRelextActon
	и запускается МЕНТАЛЬНЫЙ автоматизм.
	На стадии 4 - провоцировать оператора на ответы (почему, зачем, что такое?)
*/
func getPurposeUndestandingAndRunAutomatizm()(bool) {


	// мозжечковые рефлексы - самый первый уровень осознания - подгонка действий под заданную Цель.
	if EvolushnStage == 4 {
		/*  на стадии 4 - провоцировать оператора на ответы (почему, зачем, что такое?)


		 */
	}
/*
	// переосмыслить ситуацию - Активировать Дерево Понимания
	//understandingSituation()
	и затем создать новую цель understanding_purpose_image.go
*/


	return false
}
////////////////////////////////////////////////




//////////////////////////////////////////////////////////////////
// перезапустить осмысление или остановить цикл
func reloadConsciousness(stop bool,fromNextID int)(bool){
	if fromNextID == currrentFromNextID{//тихо (без стека прерываний) предотвратить зацикливание
		return false
	}
	if stop{// было прервано объективной активацией, запомнить остановленное и прекратить ментальный цикл
		// не запускать consciousness(2,fromNextID)
		saveInterruptMemory()

		return false
	}else {
		// перезапуск осмысления
		consciousness(2, fromNextID)
	}
	return true
}
///////////////////////////////////////////////////


///////////////////////////////
// обновление состояния информационной среды
func refreshCurrentInformationEnvironment(){
	///////// Информационная среда осознания ситуации
	// Нужно собрать всю информацию, которая может повлиять на решение.
	//  получение текущего состояния информационной среды: отражение Базового состояния и Активных Базовых контекстов
	GetCurrentInformationEnvironment()

	// оценка опасности ситуации, необходиомсть срочных действий
	veryActualSituation=CurrentInformationEnvironment.veryActualSituation
	// выявить ID парамктров гомеостаза как цели для улучшения в данных условиях
	curTargetArrID=CurrentInformationEnvironment.curTargetArrID

	/* Еще информация:
	жизненный опыт  psy_Experience.go
	доминанта psy_problem_dominanta.go
	субъектиная оценка ситуации для применения произвольности
	*/

	// актуальной инфой являются узлы активной ветки дерева понимания, особенно контекст SituationID
}
///////////////////////////////////////////////


// детекция ленивого состояния
func isIdleness()(bool){
	if CurrentInformationEnvironment.veryActualSituation || CurrentInformationEnvironment.danger {
		return false
	}

	if isCurrentProblemDominanta != nil{
		return false
	}

	return true
}
///////////////////////////////////////////////////////
// обработка структур в свободном состоянии, в первую очередь - эпизодической памяти
func processingFreeState(){
	// TODO переработка происходившего, в первую очередь - эпизодической памяти
	//EpisodeMemoryLastCalcID - последний эпизод, который был осмыслен в лени или во сне
}
//////////////////////////////////////////////////////




///////////////////////////////////////////////////////
// создание или использование ментального автоматизма инфо-функции c ID= infoID
func createNexusFromNextID(fromNextID int,infoID int)(int){

	imgID,_:=CreateNewlastMentalActionsImagesID(0,0,0,infoID,0)
	if imgID>0 {
		aID, _ := createMentalAutomatizmID(0, imgID, 1)
		if aID > 0 {
			// создание звена цепочки (всегда уникальное звено)
			fromNextID0 := fromNextID
			// создание нового элемента цепочки
			fromNextID, _ = createNewNextFromUnderstandingNodeID(
				detectedActiveLastUnderstandingNodID,
				detectedActiveLastNodID,
				fromNextID0,
				aID)

			//addShortTermMemory(fromNextID)

			return fromNextID
		}
	}

	return 0
}
//////////////////////////////////////////////////


/* обработать результат инфо-фукнции и прикинуть, что дальше с fromNextID
 */
func calcNexusFromNextID(fromNextID int)(int){
	switch currentInfoStructId{
	// если это - самая первая ф-ция, то вызвать 2-ю ф-цию поиска продолжения
	case 1: return createNexusFromNextID(fromNextID,2)
	// анализ инфо стркутуры по currentInfoStructId и выдача решения
	case 2: return analisAndSintez(fromNextID)
	}

	return 0
}
//////////////////////////////////////////////////////



//////////////////////////////////////////////////////
/* после периода ожидания, если мот.автоматизма был запущен ментально:
в конце saveFromNextIDAnswerCicle есть звено с таким запуском.
Учесть последствия ментального запуска мот.автоматизма,
из saveFromNextIDAnswerCicle выявить Правила и записать в эпизод.пямять.
Если нужно - обдумать результат в новом ментальном consciousness(2,currrentFromNextID)
Сразу после обработки периода ожидания запускается дерево понимания и объявный запуск consciousness(1,0)
так что никаких действий совершать в afterWaitingPeriod() или в самой afterWaitingPeriod() не следует.

effect =lastCommonDiffValue
*/
func afterWaitingPeriod(effect int){
	if saveFromNextIDAnswerCicle==nil || len(saveFromNextIDAnswerCicle)==0{
		return
	}
	//вытащить последний член, который ДОЛЖЕН БЫТЬ с ментальным автоматизмом запуска моторного автоматизма - ответного действия
	last:=saveFromNextIDAnswerCicle[len(saveFromNextIDAnswerCicle)-1]// всегда есть
	lastNextFrom:=goNextFromIDArr[last]
	if lastNextFrom==nil{// не должно быть такого
		return
	}
	mentAtmzm:=MentalAutomatizmsFromID[lastNextFrom.AutomatizmID]
	if mentAtmzm==nil{// не должно быть такого, поэтому выдадим панику
		lib.TodoPanic("Нет автоматизма для func afterWaitingPeriod()")
		return
	}
	mact:= MentalActionsImagesArr[mentAtmzm.ActionsImageID]
	if mact==nil{// не должно быть такого, поэтому выдадим панику
		lib.TodoPanic("Нет действия автоматизма для func afterWaitingPeriod()")
		return
	}
	if mact.activateMotorID ==0{// конечный член saveFromNextIDAnswerCicle не содержит мент.автоматизм моторного запуска
/* в конце saveFromNextIDAnswerCicle должен быть мент.автоматизм моторного запуска,
   если только Стимул не возникнет, не дожидаясь ответа на предыдущий.
   В таком случае ментальное Правило не формируется.
 */
		return
	}

	mRules,_:=createNewlastMentalTriggerAndActionID(0,saveFromNextIDAnswerCicle,mact.activateMotorID,effect)

	rID,_:=createNewlastrulesMentalID(0,[]int{mRules})
	if rID>0 {
		// записать в эпизод.пямять ментальный кадр - всегда если дошло до этой строки
		newEpisodeMemory(rID,1)
	}

	// нужно ли осмыслить это?
	//??  consciousness(2,currrentFromNextID)
}
///////////////////////////////////////////////////////





