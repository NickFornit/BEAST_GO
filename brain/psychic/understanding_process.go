/* Процессы осмысления: создание и использование ментальных автоматизмов
для Дерева понимания (или дерева ментальных автоматизмов)

*/

package psychic

import (
	"BOT/lib"
)

///////////////////////////////


//////////////////////////////////////////////////
/* Определение ближайшей Цели в данной ситуации
!!!!это - не PurposeImage (understanding_purpose_image.go)!!!
Это - постановка цели для текущего цикла размышления, чтобы оценить эффект для Правила.
Дополнение стека savePorposeIDcurrentCicle[]  addMewPorposeMemory(porposeID)
*/
func getPurposeUndestandingAndRunAutomatizm()(bool) {

	/* TODO ДОДЕЛАТЬ:
	func getPurposeUndestandingAndRunAutomatizm()(bool)
	func valuationPorposeIDcurrentCicle()
	СДЕЛАТЬ:
	ментальную функцию поставновки целей !!!! альтернативную getPurposeUndestandingAndRunAutomatizm()
	ментальную функцию оценки эффекта альтернативную getMentalEffect
	Проставить lib.WritePultConsol() для главных событий цикла
	 */


	// мозжечковые рефлексы - самый первый уровень осознания - подгонка действий под заданную Цель.
	if EvolushnStage == 4 {
		/*  на стадии 4 - провоцировать оператора на ответы (почему, зачем, что такое?)


		 */
	}


	//lib.WritePultConsol()
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
		addInterruptMemory()

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



/* получить конечный goNext цепочки
*/
func getLastAutomatizmFrom(fromNextID int)(int){
	nextID:=goNextFromIDArr[fromNextID].ID // если нет продолжения, то остается данный fromNextID
	for goNextFromIDArr[nextID].NextID>0 {
		nextID=goNextFromIDArr[nextID].NextID
	}
	return nextID
}
///////////////////////////////////////



///////////////////////////////////////////////////////
/* НАЙТИ или создать Базовое звено цепи fromNextID для данной активности деревьев
и пройти цепочку до конца, чтобы продолжить цикл от него.
 */
func createBasicLink()(int){
	fromNextID:=0
	// если такое Базовое звено?
	firstArr:=goNextFromUnderstandingNodeIDArr[detectedActiveLastUnderstandingNodID]
	if firstArr == nil { // нет веток у данного узла дерева - создать базовое звено цепочки
		// создание нового элемента цепочки
		fromNextID, _ = createNewNextFromUnderstandingNodeID(
			detectedActiveLastUnderstandingNodID,
			detectedActiveLastNodID,
			0,
			0)
	}else{// уже есть базовое звено
		fromNextID=firstArr[0].ID
	}
	//стек для обобщений: 7 Базовых fromNextID
	addMewBaseLinksMemory(fromNextID)

	// найти конец цепочки чтобы продолжить цикл от него
	fromNextID = getLastAutomatizmFrom(fromNextID)

	return fromNextID
}
//////////////////////////////////////////////////



///////////////////////////////////////////////////////
// создание или использование ментального автоматизма инфо-функции c ID= infoID
func createNexusFromNextID(fromNextID int,infoID int)(int){

	imgID,_:=CreateNewlastMentalActionsImagesID(0,0,0,0,infoID,0)
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

	// оценить совокупный эффект
	effectValuation:=getMentalEffect(effect)

	mRules,_:=createNewlastMentalTriggerAndActionID(0,saveFromNextIDAnswerCicle,mact.activateMotorID,effectValuation)

	rID,_:=createNewlastrulesMentalID(0,[]int{mRules})
	if rID>0 {
		// записать в эпизод.пямять ментальный кадр - всегда если дошло до этой строки
		newEpisodeMemory(rID,1)
	}

	// нужно ли осмыслить это?
	//??  consciousness(2,currrentFromNextID)
}
///////////////////////////////////////////////////////

func getMentalEffect(effect0 int)int{
	/* улучшилось ли положение с учетом текущего PurposeImage 4-го узла ветки понимания?
	currentUnderstandingActivedNodes[]*UnderstandingNode // начиная с конечного к первому
	*/
	effect4:=valuationPurpose()

	/* оценить эффект с учетом текущих целей savePorposeIDcurrentCicle []
	И эта цель имеет значимое преимущество для EvolushnStage > 4, где произвольная цель - главное,
	особенно в случае Доминанты.
	 */
	maneEffect:=valuationPorposeIDcurrentCicle()

	// Коэффициенты эффектов разного вида должны сильно влиять не ментальность твари, что для нее важнее.
	effectValuation:=effect0*3 + effect4*1 + maneEffect*4
	if effectValuation>1{effectValuation=1}
	if effectValuation<1{effectValuation=-1}

	return effectValuation
}
///// оценить эффект по цели 4-го узла
func valuationPurpose()(int){
	effect:=0
	node4:=currentUnderstandingActivedNodes[1:]
	purpose4id:=node4[0].PurposeID
	purpose4:=PurposeImageFromID[purpose4id]
	if purpose4!=nil && prePurpose4 !=nil {
		/* использовать для сравнения с предыдущим образом var prePurpose4 *PurposeImage*/
		if !purpose4.veryActual && prePurpose4.veryActual{
			effect+=2
		}
		if purpose4.veryActual && !prePurpose4.veryActual{
			effect-=2
		}
		//purpose4.targetID Уменьшилось ли число критических параметров с прошлого раза
		if len(purpose4.targetID) < len(prePurpose4.targetID){
			effect++
		}
		if len(purpose4.targetID) > len(prePurpose4.targetID){
			effect--
		}
		//purpose4.actionID - непонятно как использовать, разве что по Правилам или эпиз.памяти

		prePurpose4 = purpose4// сохранить для следующего
	}
	if effect>0 {
		return 1
	}
	if effect<0 {
		return -1
	}
	return 0
}
////////////////////////////////////

/*оценить эффект с учетом текущих целей savePorposeIDcurrentCicle []
И эта цель имеет значимое преимущество для EvolushnStage > 4, где произвольная цель - главное,
особенно в случае Доминанты.
 */
func valuationPorposeIDcurrentCicle()int{



	return 0
}
/////////////////////////////////////////////////////////




