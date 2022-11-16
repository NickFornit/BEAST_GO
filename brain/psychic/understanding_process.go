/* Процессы осмысления: создание и использование ментальных автоматизмов
для Дерева понимания (или дерева ментальных автоматизмов)

*/

package psychic

import (
	"BOT/brain/gomeostas"
	"BOT/lib"
)

//////////////////////////////////////////////////





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
func processingFreeState(stopMentalWork bool){

	// if stopMentalWork{ - прекратить обработку

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
/* Продолжить цепочку осмысления: найти мотивированное продолжение:
- создается мент.авт-м запуска infoFuncNNN() - infoID
 т.е. - создание ментального автоматизма инфо-функции c ID= infoID
Это не ментальная функция! а наследственная структура для мотивации действий и направления мышления.
 */

func createNexusFromNextID(fromNextID int,infoID int)(int){
// typeD==4 - запуск инфо-функции
	imgID,_:=CreateNewlastMentalActionsImagesID(0,4,infoID,true)
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
	/* Если не было цикла осмысления, а проходилиь только уровни до 3-го,
	то и нет обработки, нет записи в Эпизод.память ментальныз кадров. Что определяет ощцщение субъективного времени.
	 */
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
// ???		lib.TodoPanic("Нет автоматизма для func afterWaitingPeriod()")
		return
	}
	mact:= MentalActionsImagesArr[mentAtmzm.ActionsImageID]
	if mact==nil{// не должно быть такого, поэтому выдадим панику
		lib.TodoPanic("Нет действия автоматизма в func afterWaitingPeriod()")
		return
	}
	if mact.typeID !=5{// конечный член saveFromNextIDAnswerCicle не содержит мент.автоматизм моторного запуска
/* в конце saveFromNextIDAnswerCicle должен быть мент.автоматизм моторного запуска,
   если только Стимул не возникнет, не дожидаясь ответа на предыдущий.
   В таком случае ментальное Правило не формируется.
 */
		return
	}

	// оценить совокупный эффект
	effectValuation:=getMentalEffect(effect)
// mact.valID потому что это - точно моторный запуск, выше: if mact.typeID !=5{
	mRules,_:=createNewlastMentalTriggerAndActionID(0,saveFromNextIDAnswerCicle,mact.valID,effectValuation,true)

	rID,_:=createNewlastrulesMentalID(0,detectedActiveLastNodID,detectedActiveLastUnderstandingNodID,[]int{mRules},true)
	if rID>0 {
		// записать в эпизод.пямять ментальный кадр - всегда если дошло до этой строки
		newEpisodeMemory(rID,1)
	}

// образ значимости ментального действия
createNewlastImportanceID(0, detectedActiveLastNodID, detectedActiveLastUnderstandingNodID, 2, mentAtmzm.ActionsImageID, effectValuation, true)


	// нужно ли осмыслить это?
	//??  consciousness(2,currrentFromNextID)
}
///////////////////////////////////////////////////////

func getMentalEffect(effect0 int)int{
	/* улучшилось ли положение с учетом текущего PurposeImage 4-го узла ветки понимания?
	currentUnderstandingActivedNodes[]*UnderstandingNode // начиная с конечного к первому
	*/
	effect4:=valuationPurpose()

	// улучшилась ли значимость объекта внимания
	effectA:=0
	if extremImportanceObject!=nil {
		oldVal := extremImportanceObject.extremVal
		_, newVal := getObjectsImportanceValue(extremImportanceObject.kind, extremImportanceObject.objID, detectedActiveLastNodID, detectedActiveLastUnderstandingNodID)
		if newVal > oldVal { // улучшилось
			effectA = 1
		}
		if newVal < oldVal { // улучшилось
			effectA = -1
		}
	}

	// улучшилась ли значимость субъекта внимания
	effectMA:=0
	if extremImportanceMentalObject!=nil {
		oldVal := extremImportanceMentalObject.extremVal
		_, newVal := getObjectsImportanceValue(extremImportanceMentalObject.kind, extremImportanceMentalObject.objID, detectedActiveLastNodID, detectedActiveLastUnderstandingNodID)
		if newVal > oldVal { // улучшилось
			effectMA = 1
		}
		if newVal < oldVal { // улучшилось
			effectMA = -1
		}
	}

	// Коэффициенты эффектов разного вида должны сильно влиять не ментальность твари: что для нее важнее.
	effectValuation:=0
	if EvolushnStage < 5 {
		effectValuation = effect0*3 + effect4*1 + effectA*1 +effectMA*0
	}else{// повышенная роль заданной цели effect4 и объекта внимания
		effectValuation = effect0*2 + effect4*3 + effectA*3 +effectMA*2
	}
	if effectValuation>1{effectValuation=1}
	if effectValuation<1{effectValuation=-1}

	return effectValuation
}
///// оценить эффект по цели 4-го узла
func valuationPurpose()(int){
	effect:=0
	// т.к. в начале конечный узел, то берем первый (PurposeID)
	node4:=currentUnderstandingActivedNodes[0]
	purpose4id:=node4.PurposeID
	purpose4:=PurposeImageFromID[purpose4id]
	prePurpose4:=PurposeImageFromID[prePurposeID]
	if purpose4!=nil && prePurpose4 !=nil {
		/* использовать для сравнения с предыдущим образом var prePurpose4 *PurposeImage*/
		if !purpose4.veryActual && prePurpose4.veryActual{
			effect+=2
		}
		if purpose4.veryActual && !prePurpose4.veryActual{
			effect-=2
		}

		effect+=compareGomeoPars(purpose4.targetID,prePurpose4.targetID)

		// добились ли того, что оператор сделал желаемое?
		if purpose4.actionID>0{
			effect+=compareOperatorsAction(purpose4.actionID)
		}
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

func compareGomeoPars(gPar []int,gParPre []int)int{
	effect:=0
	//purpose4.targetID Уменьшилось ли число критических параметров с прошлого раза
	if len(gPar) < len(gParPre){
		effect++
	}
	if len(gPar) > len(gParPre){
		effect--
	}

	/* сравинить бодее тщательно, по значимости параметров: чем более значмый параметр обнулился, тем сильнее эффект
	 ID параметров гомеостаза как цели для улучшения в данных условиях
	уже отсортированы по убыванию значимости
	 */
	for i := 0; i < len(gParPre); i++ {
		if !lib.ExistsValInArr(gPar, gParPre[i]){
	// больше нет такого значения в новом массиве, эффект увеличивается со значимостью
			w:=gomeostas.GomeostazParamsWeight[gParPre[i]]
			if w>50{
				effect+=3
			}else{
				if w>10{
					effect+=2
				}else{
					effect++
				}
			}
		}
	}

	return effect
}

func compareOperatorsAction(purpose4actionID int)int{
	ppA:=ActionsImageArr[purpose4actionID]
	if ppA==nil{
		return 0
	}
	/* curActiveActions - структура действий оператора при активации дерева автоматизмов типа ActionsImage
		curActiveActions.ActID
		curActiveActions.PhraseID
		curActiveActions.ToneID=0
		curActiveActions.MoodID=0
		 */
if lib.EqualArrs(ppA.ActID,curActiveActions.ActID) && lib.EqualArrs(ppA.PhraseID, curActiveActions.PhraseID){
	// достаточно полное совпадение
	if mentalPurposeImageID>0{// при призвольно заданной цели
		return 4
	}else {
		return 2
	}
}
	if lib.EqualArrs(ppA.ActID,curActiveActions.ActID) || lib.EqualArrs(ppA.PhraseID, curActiveActions.PhraseID){
		// частичное совпадение
		if mentalPurposeImageID>0{// при призвольно заданной цели
			return 2
		}else {
			return 1
		}
	}
	return 0
}
/////////////////////////////////////////////////////////




