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
	"strconv"
)

/////////////////////////////////////////

/* Это используется для определения момента реакция оператора Пульта на действия автоматизма.
За 20 сек г.параметры могли бы просто натечь и вызывать сработавание при ожидании ответной реакции.
Флаг сбрасывается через пульс после запуска автоматизма.
*/
var WasOperatorActiveted=false

var WasConditionsActiveted=false

// период ожидания реакции оператора на действие автоматизма
const WaitingPeriodForActionsVal=20


var savePsyBaseMood=0 // -1 Плохое настроение, 0 Нормальное, 1 - хорошее настроение
// для более точной оценки
var savePsyMood=0//сила Плохо -10 ... 0 ...+10 Хорошо
// НОВИЗНА СИТУАЦИИ сохраняется значение CurrentAutomatizTreeEnd[] для решений
var savedNoveltySituation []int

// отслеживать Правила из Пульта в http://go/pages/rulles.php
var RullesOutputProcess=false // режим отслеживания
var RullesOutputStr="" // текущее состояние последний 10 правил

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
	lib.WritePultConsol("<span style='color:blue'>Ожидание ответа оператора.</span>")

	// при срабатывании автоматизма - блокируются все рефлексторные действия
	//MotorTerminalBlocking=true
	notAllowReflexRuning=false //уже есть, но на всякий случай :)

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
	//MotorTerminalBlocking=false
	notAllowReflexRuning=false

	LastAutomatizmWeiting=nil
	LastRunAutomatizmPulsCount=0
	WasOperatorActiveted=false
	onliOnceWasConditionsActiveted=false
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
		clinerAutomatizmRunning()
		return true
	}

	// не опасная ситуация, можно поэкспериментировать
	if EvolushnStage == 3 && !CurrentPurposeGenetic.veryActual{
	/* в случае отсуствия автоматизма в данных условиях - послать оператору те же стимулы, чтобы посмотреть его реакцию.
		   Создание автоматизма, повторяющего действия оператора в данных условиях
		НО если уже помылался provokatorMirrorAutomatizm то больше не делать этого (бесконечный цикл)
	*/
		if oldProvokatorAutomatizm != LastAutomatizmWeiting {// не повторять, если только что был такой ответ
			provokatorMirrorAutomatizm(LastAutomatizmWeiting, &CurrentPurposeGenetic)
			clinerAutomatizmRunning()
			return true
		}
	}

	// реакция была, но но оператор не обратил на нее внимания, нужно усилить силу действия мозжечковым рефлексом
	if cerebellumCoordination(LastAutomatizmWeiting,1){
		// и тут же снова запустить реакцию!
		if oldProvokatorAutomatizm != LastAutomatizmWeiting {// не повторять, если только что был такой ответ
			setAutomatizmRunning(LastAutomatizmWeiting, &CurrentPurposeGenetic)
			clinerAutomatizmRunning()
			return true
		}
	}
	clinerAutomatizmRunning()
	return false
}
/////////////////////////////////////////////////////////////////////








/* ПОСЛЕ ОРИЕНТИРОВОЧНОГО РЕФЛЕКСА оценивать действие запущенного автоматизма

 */
func calcAutomatizmResult(lastCommonDiffValue int,lastBetterOrWorse int,wellIDarr []int){

lib.WritePultConsol("<span style='color:blue;background-color:#FFD0FF;'>Был ОТВЕТ ОПЕРАТОРА. До ответа оператора сосотояние: <b>"+strconv.Itoa(lastBetterOrWorse)+"</b>, изменение на: <b>"+strconv.Itoa(lastCommonDiffValue)+"</b></span>")

	// lastBetterOrWorse больше не применяется т.к. lastCommonDiffValue более точен и информативен

	automatizmCorrection(lastCommonDiffValue,wellIDarr)

	// >3 потому, что раньше не пишется эпизодическая память и формируются более примитивные механизмы.
	if EvolushnStage > 3 {

		/* При каждом ответе на действия оператора - прописывать текущее правило rules
		   		и делать новый кадр эпизодической памяти
		      А так же просматривать эпизод память взад макчимум на EpisodeMemoryPause шагов или до паузы в общении > 30 шагов,
		   		фиксируя цепочку правил.
		*/
		ai1, _ := СreateNewlastActionsImageID(0, curActiveActions.ActID, curActiveActions.PhraseID, curActiveActions.ToneID, curActiveActions.MoodID)
		fixNewRules(lastCommonDiffValue,ai1)
	}

return
}
///////////////////////////////////////////////////////

//корректируется успешность автоматизма - реакция на результат lastCommonDiffValue
func automatizmCorrection(lastCommonDiffValue int,wellIDarr []int){
	/// если числа имеют разные знаки (одно положительное, другое отрицательное)
	if lib.IsDiffersOfSign(LastAutomatizmWeiting.Usefulness,lastCommonDiffValue){
		LastAutomatizmWeiting.Count=0 // сбрасываем  надежность
	} else {
		LastAutomatizmWeiting.Count++
	}

	// изменять полезность по 1 шагу!
	if lastCommonDiffValue>0 && LastAutomatizmWeiting.Usefulness<10 {
		LastAutomatizmWeiting.Usefulness++ // lastBetterOrWorse
	}
	if lastCommonDiffValue<0 && LastAutomatizmWeiting.Usefulness>-10 {
		LastAutomatizmWeiting.Usefulness-- // lastBetterOrWorse
	}

	if LastAutomatizmWeiting.Usefulness>0 {
		// задать тип автоматизма, 2 - проверенный
		SetAutomatizmBelief(LastAutomatizmWeiting, 2) // ТАК ПРОСТО НЕЛЬЗЯ ЗАДАВАТЬ Belief=2: LastAutomatizmWeiting.Belief=2
	}

	if lastCommonDiffValue>0{// стало лучше
		PsyBaseMood=1
		// список гомео параметро, которые улучшило это действие
		if wellIDarr != nil {
			LastAutomatizmWeiting.GomeoIdSuccesArr = wellIDarr // м.б. nil !!!! если нет таких явных действий
		}
		// пополняется список полезных автоматизмов
		AutomatizmSuccessFromIdArr[LastAutomatizmWeiting.ID] = LastAutomatizmWeiting
	}

	if lastCommonDiffValue<0{// стало хуже
		PsyBaseMood=-1
		// очистить списки улучшения
		LastAutomatizmWeiting.GomeoIdSuccesArr=nil
		if AutomatizmSuccessFromIdArr[LastAutomatizmWeiting.ID] !=nil {
			AutomatizmSuccessFromIdArr[LastAutomatizmWeiting.ID] =nil
		}
	}
}
//////////////////////////////////////////////////////



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


/* // текущий ID пускового стимула типов curActiveActions или curBaseStateImage
при активации дерева автоматизмов. Если тип curBaseStateImage, то ID отрицательное (ID<0)!
 */
var currentTriggerID=0



////////////////////////////////////////////////////////////////////////
/* на стадии >3 при каждом ответе на действия оператора - прописывать текущее правило rules.
   А так же просматривать эпизод память взад макчимум на 6 шагов или до паузы в общении > EpisodeMemoryPause шагов,
		фиксируя цепочку правил.
*/
func fixNewRules(lastCommonDiffValue int,ai1 int) int {
	currentTriggerID=ai1
	if LastAutomatizmWeiting == nil{
		return 0
	}

	// образ действий оператора
	if ai1 == 0  || LastAutomatizmWeiting == nil {
		return 0
	}
	// ответный образ действий Beast
	ai2:=LastAutomatizmWeiting.ActionsImageID
	if ai2 == 0{return 0}
	TriggerAndAction,_:=createNewlastTriggerAndActionID(0,ai1,ai2,lastCommonDiffValue)
	if TriggerAndAction == 0{return 0}
	rulesID, _ := createNewlastrulesID(0, []int{TriggerAndAction})
	if rulesID == 0{return 0}

	lib.WritePultConsol("<span style='color:green'>Записано <b>ПРАВИЛО № "+strconv.Itoa(rulesID)+"</b></span>")

	// новый кадр эпизодической памяти, сохраняющий
	newEpisodeMemory(rulesID,0) // запись эпизодической памяти saveEpisodicMenory()

	// теперь обрабатываем прошлую эпизодическую память
	GetRulesFromEpisodeMemory(0)

if RullesOutputProcess{// отслеживать Правила из Пульта в http://go/pages/rulles.php
	RullesOutputStr=getCur10lastRules()
	RullesOutputProcess=false
}

	return rulesID
}
///////////////////////////////////////////////////////////////////////


// обработать изменение состояния - записать Правило типа BaseStateImage
func fixRulesBaseStateImage(lastCommonDiffValue int){
	//корректируется успешность автоматизма - как в calcAutomatizmResult
	automatizmCorrection(lastCommonDiffValue,nil)
	/////////////////////// ПРАВИЛО:
	ai1, _ := СreateNewlastBaseStateImageID(0, curBaseStateImage.Mood, curBaseStateImage.EmotionID, curBaseStateImage.SituationID)
	ai1*=-1 // отрицательное значение идентифицирует образ - как текущего сосотояния!!!
	currentTriggerID=ai1
	fixNewRules(lastCommonDiffValue,ai1)
}
/////////////////////////////////////////////////////////////////////