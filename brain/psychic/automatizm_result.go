/* Ожидание результата запущенного автоматизма и его обработка

В BAD_detector.go в самом низу есть func BetterOrWorseNow() с комментариями по делу. Я ее отрабатывал как раз для того, чтобы фиксировать любые улучшения или ухудшения для определения эффекта автоматизма.
Она вызывается (через трансформатор против цицличности wasChangingMoodCondition()) 2 раза: в момент запуска автоматизма и как только совершится любое действие оператора на пульте. Таким образом в automatizm_result.go получается дифферент:
oldlastBetterOrWorse,oldBetterOrWorse,oldParIdSuccesArr = wasChangingMoodCondition()
Т.е. если ты поставишь точку прерывания на
oldlastBetterOrWorse,oldBetterOrWorse,oldParIdSuccesArr = wasChangingMoodCondition()
то и получишь эффект автоматизма.
*/


package psychic

import (
	"BOT/brain/gomeostas"
	"BOT/lib"
)

/////////////////////////////////////////

/* Это используется для определения момента реакция оператора Пульта на действия автоматизма.
За 20 сек г.параметры могли бы просто натечь и вызывать сработавание при ожидании ответной реакции.
Флаг сбрасывается через пульс после запуска автоматизма.
*/
var WasOperatorActiveted=false

// период ожидания реакции оператора на действие автоматизма
const WaitingPeriodForActionsVal=20


var savePsyBaseMood=0 // -1 Плохое настроение, 0 Нормальное, 1 - хорошее настроение
// для более точной оценки
var savePsyMood=0//сила Плохо -10 ... 0 ...+10 Хорошо
// НОВИЗНА СИТУАЦИИ сохраняется значение CurrentAutomatizTreeEnd[] для решений
var savedNoveltySituation []int


/////////////////////////////////////////////////////////////////////
// отслеживание запущенных автоматизмов
// структура примитивных целей, создающих контекст ситуации НЕ СБРАСЫВАЕТСЯ после ожидания
var savePurposeGenetic *PurposeGenetic
/* При запуске автоматизма определяются:
// момент запуска автоматизма в числе пульсов
var LastRunAutomatizmPulsCount =0 //сбрасывать ожидание результата автоматизма если прошло 20 пульсов
// ожидается результат запущенного MotAutomatizm
var LastAutomatizmWeiting *Automatizm //сбрасывается указатель автоматизма
 */
//////////////////////////////////////////////////////////////////////

func setAutomatizmRunning(am *Automatizm,ps *PurposeGenetic){
	// при срабатывании автоматизма - блокируются все рефлексторные действия
	MotorTerminalBlocking=true // уже есть, но на всякий случай :)
	LastAutomatizmWeiting=am // уже есть, но для надежности :)
	LastDetectedActiveLastNodID=detectedActiveLastNodID // уже есть, но для надежности :)

	savePsyBaseMood=PsyBaseMood
	savePsyMood=PsyMood
	savedNoveltySituation=NoveltySituation
	if ps != nil {
		savePurposeGenetic = ps
	}
	WasOperatorActiveted=false // ждем ответа оператора
	// зафиксировать текущее состояние на момент срабатывания автоматизма
	oldlastBetterOrWorse,oldBetterOrWorse,oldParIdSuccesArr = wasChangingMoodCondition(1)
}
///////////////////////////////////
func clinerAutomatizmRunning(){
	MotorTerminalBlocking=false
	LastAutomatizmWeiting=nil
	LastRunAutomatizmPulsCount=0
// !!!! НЕ СБРАСЫВАТЬ savePurposeGenetic=nil - он может определяться независимо от запуска автоматизма
}
///////////////////////////////////


//////////////////////// ПУЛЬС
// ПУЛЬС срабатывания по каждому Пульсу здесь для удобства
var oldBetterOrWorse=0 //- стали лучше или хуже: величина измнения от -10 через 0 до 10
var oldParIdSuccesArr []int//стали лучше следующие г.параметры []int гоменостаза
var oldlastBetterOrWorse=0 // насколько изменилось общее состояние, значение от  -10(максимально Плохо) через 0 до 10(максимально Хорошо)
func automatizmActionsPuls(){

	if LastRunAutomatizmPulsCount==0 {
		return
	}
// вышло время ожидания реакции
		if (LastRunAutomatizmPulsCount+WaitingPeriodForActionsVal) < PulsCount {
			// отреагировать на отсуствие реакции - повторить автоматизм с большей силой Energy
			// Из МОЗЖУЧКА как-то отреагировать на отсуствие реакции - повторить автоматизм с большей силой Energy
			if noAutovatizmResult(){// была попытка отреагировать сильнее - в cerebellum.go
				return // чтобы не сбрасывать clinerAutomatizmRunning()
			}
			//сбрасывать ожидание результата автоматизма если прошло WaitingPeriodForActionsVal пульсов
			clinerAutomatizmRunning()
		}
}
/////////////////////////////////////////////////////////////////////
// отреагировать на отсуствие реакции - повторить автоматизм с большей силой Energy
func noAutovatizmResult()(bool){

	if EvolushnStage > 3 {
		// осмыслить ситуацию - Активировать Дерево Понимания
		understandingSituation()
		return true
	}

	// не опасная ситуация, можно поэкспериментировать
	if EvolushnStage == 3 && !CurrentPurposeGenetic.veryActual{
	/* в случае отсуствия автоматизма в данных условиях - послать оператору те же стимулы, чтобы посмотреть его реакцию.
		   Создание автоматизма, повторяющего действия оператора в данных условиях
	*/
		provokatorMirrorAutomatizm(LastAutomatizmWeiting,&CurrentPurposeGenetic)
		return true
	}

	// реакция была, но но оператор не обратил на нее внимания, нужно усилить силу действия
	if cerebellumCoordination(LastAutomatizmWeiting,1){
		// и тут же снова запустить реакцию!
		setAutomatizmRunning(LastAutomatizmWeiting, &CurrentPurposeGenetic)
		return true
	}

	return false
}
/////////////////////////////////////////////////////////////////////








/* ПОСЛЕ ОРИЕНТИРОВОЧНОГО РЕФЛЕКСА оценивать действие запущенного автоматизма

 */
func calcAutomatizmResult(lastCommonDiffValue int,lastBetterOrWorse int,wellIDarr []int){

	// lastBetterOrWorse - точно изменился, иначе бы не было вызова calcAutomatizmResult
	/// если числа имеют разные знаки (одно положительное, другое отрицательное)
	if lib.IsDiffersOfSign(LastAutomatizmWeiting.Usefulness,lastBetterOrWorse){
		LastAutomatizmWeiting.Count=0 // сбрасываем  надежность
	} else {
		LastAutomatizmWeiting.Count++
	}
	// задать тип автоматизма, 2 - проверенный
	SetAutomatizmBelief(LastAutomatizmWeiting,2)// ТАК ПРОСТО НЕЛЬЗЯ ЗАДАВАТЬ Belief=2: LastAutomatizmWeiting.Belief=2

	// изменять полезность по 1 шагу!
	if lastBetterOrWorse>0 && LastAutomatizmWeiting.Usefulness<10 {
		LastAutomatizmWeiting.Usefulness++ // lastBetterOrWorse
	}
	if lastBetterOrWorse<0 && LastAutomatizmWeiting.Usefulness>-10 {
		LastAutomatizmWeiting.Usefulness-- // lastBetterOrWorse
	}


	if lastBetterOrWorse>0{// стало лучше
		PsyBaseMood=1
		// список гомео параметро, которые улучшило это действие
		LastAutomatizmWeiting.GomeoIdSuccesArr=wellIDarr // м.б. nil !!!! если нет таких явных действий
		// пополняется список полезных автоматизмов
		if lastBetterOrWorse>0 {
			AutomatizmSuccessFromIdArr[LastAutomatizmWeiting.ID] = LastAutomatizmWeiting
		}
	}
	if EvolushnStage == 3{
/* отзеркаливание ответа оператора не зависимо от того, стало хуже или лучше
потому, что это был ответ оператора на действия автоматизма, значит - авторитетный ответ
   Создание автоматизма, повторяющего действия оператора в данных условиях
 */
		createNewMirrorAutomatizm(LastAutomatizmWeiting)

	}

	if lastBetterOrWorse<0{// стало хуже
		PsyBaseMood=-1
		// очистить списки улучшения
		LastAutomatizmWeiting.GomeoIdSuccesArr=nil
		AutomatizmSuccessFromIdArr[LastAutomatizmWeiting.ID] =nil
	}

return
}
///////////////////////////////////////////////////////



//////////////////////////////////////////////////////////////////
// для индикации период ожидания реакции оператора на действие автоматизма
//   psychicWaitingPeriodForActions()
func WaitingPeriodForActions()(bool,int){

	if LastRunAutomatizmPulsCount>0{
		time:=WaitingPeriodForActionsVal - (PulsCount-LastRunAutomatizmPulsCount)
		return true,time
	}

	return false,0
}
//////////////////////////////////////////



/////////////////////////////////////////////////////////////////////////
/* сканируется с каждым пульсом в func automatizmActionsPuls() во время ожидания
В  gomeostas.BetterOrWorseNow() учитывается CommonMoodAfterAction - Общее (де)мотивирующее действие с Пульта

res - стали лучше или хуже: величина измнения от -10 через 0 до 10
wellIDarr - стали лучше следующие г.параметры []int гоменостаза
*/
func wasChangingMoodCondition(kind int)(int,int,[]int){
	//стало хуже или лучше теперь, возвращает величину измнения от -10 через 9 до 10
	res0,res,wellIDarr:=gomeostas.BetterOrWorseNow(kind)

	return res0,res,wellIDarr
}
/////////////////////////////////////////////////////////////////////////