/* Ожидание результата запущенного автоматизма и его обработка

В BAD_detector.go в самом низу есть func BetterOrWorseNow() с комментариями по делу. Я ее отрабатывал как раз для того, чтобы фиксировать любые улучшения или ухудшения для определения эффекта автоматизма.
Она вызывается (через трансформатор против цицличности wasChangingMoodCondition()) 2 раза: в момент запуска автоматизма и как только совершится любое действие оператора на пульте. Таким образом в automatizm_result.go получается дифферент:
oldCommonDiffValue,oldBetterOrWorse,oldParIdSuccesArr = wasChangingMoodCondition()
Т.е. если ты поставишь точку прерывания на
oldCommonDiffValue,oldBetterOrWorse,oldParIdSuccesArr = wasChangingMoodCondition()
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


//////////////////////// ПУЛЬС
// ПУЛЬС срабатывания по каждому Пульсу здесь для удобства
var oldBetterOrWorse=0 //- стали лучше или хуже: величина измнения от -10 через 0 до 10
var oldParIdSuccesArr []int//стали лучше следующие г.параметры []int гоменостаза
var oldCommonDiffValue=0 // насколько изменилось общее состояние, значение от  -10(максимально Плохо) через 0 до 10(максимально Хорошо)
func automatizmActionsPuls(){
	if LastRunAutomatizmPulsCount==0 {
		return
	}

	// если был запущен автоматизм, возможно, без ориентировчного рефлека, рефлексы блокируются на 2 пульса
	if MotorTerminalBlocking && LastRunAutomatizmPulsCount+2 > PulsCount {
		MotorTerminalBlocking=false // снять блокировку
	}

/* Ожидание результата автоматизма ПОСЛЕ ОРИЕНТИРОВОЧНОГО РЕФЛЕКСА.
Реакция ожидается на слелующем пульcе после срабатывания автоматизма	и в течение 20 пульсов
 за это время получим уверенное wasChangingMoodCondition() по значению gomeostas.BetterOrWorseNow()
 */
		if (LastRunAutomatizmPulsCount+1) == PulsCount {
			WasOperatorActiveted=false
			// зафиксировать текущее состояние на момент срабатывания автоматизма
			oldCommonDiffValue,oldBetterOrWorse,oldParIdSuccesArr = wasChangingMoodCondition(1)
			if oldCommonDiffValue>0{}
		}
		if (LastRunAutomatizmPulsCount+2) < PulsCount {// следить со следующего пульса
/* 	Контроль за изменением состояния, возвращает:
	lastCommonDiffValue - насколько изменилось общее состояние
   	lastBetterOrWorse - стали лучше или хуже: величина измнения от -10 через 0 до 10
   	gomeoParIdSuccesArr - стали лучше следующие г.параметры []int гоменостаза
 */
if WasOperatorActiveted { // оператор отреагировал
	lastCommonDiffValue,lastBetterOrWorse,gomeoParIdSuccesArr := wasChangingMoodCondition(2)
				// обработать изменение состояния
				calcAutomatizmResult(lastCommonDiffValue,lastBetterOrWorse, gomeoParIdSuccesArr)
				//  clinerAutomatizmRunning()  есть в calcAutomatizmResult
			}
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
	//////////////////////////////////////////////





/* ПОКА НИКАК НЕ ИСПОЛЬЗУЕТСЯ - после срабатывания актуального автоматизма ветки дерева.
т.е. активная ветка не сопровождается новизной, но м.б. есть технически невидимая новизна
   и нужно так же оценить последствия, и если они плохие, то задуматься.
   Ожидание результата автоматизма БЕЗ ОРИЕНТИРОВОЧНОГО РЕФЛЕКСА (автоматически запущенного из Дерева).
	Реакция ожидается на слелующем пульcе после срабатывания автоматизма	и в течение WaitingPeriodForActionsVal пульсов
	 за это время получим уверенное wasChangingMoodCondition() по значению gomeostas.BetterOrWorseNow()
*/
/*
	if (AutomatizmRunningPulsCountAut+1)<PulsCount{
		WasOperatorActiveted=false
		// зафиксировать текущее состояние на момент срабатывания автоматизма
		oldCommonDiffValue,oldBetterOrWorse,oldParIdSuccesArr = wasChangingMoodCondition()
	}


	if (AutomatizmRunningPulsCountAut+2)<PulsCount && AutomatizmRunningPulsCountAut+WaitingPeriodForActionsVal < PulsCount{

		// раскомментировать когда дойдешь до отладки этого
		//lastCommonDiffValue,lastBetterOrWorse,gomeoParIdSuccesArr := wasChangingMoodCondition()
		if WasOperatorActiveted { // оператор отреагировал
			// обработать изменение состояния
			//calcAutomatizmResultAut(lastCommonDiffValue,lastBetterOrWorse, gomeoParIdSuccesArr)
		}
	}
	if AutomatizmRunningPulsCountAut+WaitingPeriodForActionsVal < PulsCount {
		// Из МОЗЖУЧКА как-то отреагировать на отсуствие реакции - повторить автоматизм с большей силой Energy
		if noAutovatizmResult(){// была попытка отреагировать сильнее - в cerebellum.go
			return
		}
		//сбрасывать ожидание результата автоматизма если прошло WaitingPeriodForActionsVal пульсов
		clinerAutomatizmRunningAut()
	}
	*/



}
/////////////////////////////////////////////////////////////////////
// отреагировать на отсуствие реакции - повторить автоматизм с большей силой Energy
func noAutovatizmResult()(bool){

	// не опасная ситуация, можно поэкспериментировать
	if EvolushnStage == 3 && !CurrentPurposeGenetic.veryActual{
	/* в случае отсуствия автоматизма в данных условиях - послать оператору те же стимулы, чтобы посмотреть его реакцию.
		   Создание автоматизма, повторяющего действия оператора в данных условиях
	*/
		provokatorMirrorAutomatizm(AutomatizmRunning,&CurrentPurposeGenetic)
		return true
	}

	if EvolushnStage > 3 {
		// создать образ ситуации
		autmzmTreeNodeID := AutomatizmRunning.BranchID
		id, _ := createSituationImage(0, autmzmTreeNodeID, 6)
		// осмыслить ситуацию - Активировать Дерево Понимания
		understandingSituation(id, savePurposeGenetic)
		return true
	}

	// реакция была, но но оператор не обратил на нее внимания, нужно усилить силу действия
	if cerebellumCoordination(AutomatizmRunning,1){
		// и тут же снова запустить реакцию!
		setAutomatizmRunning(AutomatizmRunning, &CurrentPurposeGenetic)
		return true
	}

	return false
}
/////////////////////////////////////////////////////////////////////




/////////////////////////////////////////////////////////////////////
// отслеживание запущенных автоматизмов
var AutomatizmRunning *Automatizm // запущенный автоматизм
var AutomatizmRunningPulsCount=0 // время запуска автоматизма WaitingPeriodForActionsVal сек ожидания (if AutomatizmRunningPulsCount+WaitingPeriodForActionsVal < PulsCount {)
var savePurposeGenetic *PurposeGenetic // массив примитивных целей, создающих контекст ситуации

func setAutomatizmRunning(am *Automatizm,ps *PurposeGenetic){
	AutomatizmRunning=am
	savePsyBaseMood=PsyBaseMood
	savePsyMood=PsyMood
	savedNoveltySituation=NoveltySituation
	if ps != nil {
		savePurposeGenetic = ps
	}
}
///////////////////////////////////
func clinerAutomatizmRunning(){
	AutomatizmRunning=nil
	LastRunAutomatizmPulsCount=0
	isReflexesActionBloking=false
}
///////////////////////////////////



/* ПОСЛЕ ОРИЕНТИРОВОЧНОГО РЕФЛЕКСА оценивать действие запущенного автоматизма

 */
func calcAutomatizmResult(commonDiffValue int,diffPsyBaseMood int,wellIDarr []int){
	/*
	if AutomatizmRunningPulsCount==0 || AutomatizmRunning==nil{
		clinerAutomatizmRunning()
		return
	}
	 */
	// commonDiffValue - точно изменился, иначе бы не было вызова calcAutomatizmResult
	/// если числа имеют разные знаки (одно положительное, другое отрицательное)
	if lib.IsDiffersOfSign(AutomatizmRunning.Usefulness,commonDiffValue){
		AutomatizmRunning.Count=0 // сбрасываем  надежность
	} else {
		AutomatizmRunning.Count++
	}
	// задать тип автоматизма, 2 - проверенный
	SetAutomatizmBelief(AutomatizmRunning,2)// ТАК ПРОСТО НЕЛЬЗЯ ЗАДАВАТЬ Belief=2: AutomatizmRunning.Belief=2

	AutomatizmRunning.Usefulness =commonDiffValue // diffPsyBaseMood


	if commonDiffValue>0{// стало лучше
		PsyBaseMood=1
		// список гомео параметро, которые улучшило это действие
		AutomatizmRunning.GomeoIdSuccesArr=wellIDarr // м.б. nil !!!! если нет таких явных действий
		// пополняется список полезных автоматизмов
		if commonDiffValue>0 {
			AutomatizmSuccessFromIdArr[AutomatizmRunning.ID] = AutomatizmRunning
		}
	}
	if EvolushnStage == 3{
/* отзеркаливание ответа оператора не зависимо от того, стало хуже или лучше
потому, что это был ответ оператора на действия автоматизма, значит - авторитетный ответ
   Создание автоматизма, повторяющего действия оператора в данных условиях
 */
		createNewMirrorAutomatizm(AutomatizmRunning)
	}

	if commonDiffValue<0{// стало хуже
		PsyBaseMood=-1
		// очистить списки улучшения
		AutomatizmRunning.GomeoIdSuccesArr=nil
		AutomatizmSuccessFromIdArr[AutomatizmRunning.ID] =nil

	}

	// только если серьезно изменилась ситуация
	if diffPsyBaseMood!=0{// изменилась ситуация
		// обновить информационное окружение
		GetCurrentInformationEnvironment()
		// переактивировать дерево рефлексов
		automatizmTreeActivation()//и возникает новый цикл активации Дерева, уже по внутренним причинам
	}
	if EvolushnStage > 3 {
		// создать образ ситуации
		autmzmTreeNodeID := AutomatizmRunning.BranchID
		id, _ := createSituationImage(0, autmzmTreeNodeID, 1)
		// осмыслить ситуацию - Активировать Дерево Понимания
		understandingSituation(id, savePurposeGenetic)
	}

// оценить значимость поизнесенной фразы в VerbalFromIdArr структурах Дерева Понимания??

/* !!!!!! допонение cerebellumReflexFromID[LastAutomatizmWeiting.ID] другими корректирующими действиеями
если это еще получается, но при отсуствии эффекта нужно создавать новый автоматизм.
Это - только на уровне осмысления в Дереве Понимания:
   cerebellumCoordination(AutomatizmRunning.ID)
Должна быть осознание цели и перебеора-недобора!!!!!!
   В каждом автоматизме есть параметр силы: Automatizm.Energy вот он и корректируется.
 */

	clinerAutomatizmRunning()
return
}
///////////////////////////////////////////////////////
















// начало неработающего кода
//#############################################################################
/*  ПОКА НИКАК НЕ ИСПОЛЬЗУЕТСЯ ?? TODO
Отслеживание запущенных автоматизмов ВНЕ ОРИЕНТИРОВОЧНОГО РЕФЛЕКСА
должно использовать всю ту же функцию
func automatizmActionsPuls() !!!!!

ПОКА ОСТАВЛЯЮ ЧТОБЫ СОХРАНИТЬ ИДЕИ....

ЕСЛИ срабатывает автоматизм без ориентировочного рефлекса, значит есть технически невидимая новизна
и нужно так же оценить последствия, и если они плохие, то задуматься.
*/
var savePsyBaseMoodAut=0 // -1 Плохое настроение, 0 Нормальное, 1 - хорошее настроение
// для более точной оценки
var savePsyMoodAut=0//сила Плохо -10 ... 0 ...+10 Хорошо

var AutomatizmRunningAut *Automatizm // запущенный автоматизм
var AutomatizmRunningPulsCountAut=0 // время запуска автоматизма
func setAutomatizmRunningAut(am *Automatizm){
	AutomatizmRunningAut=am
	AutomatizmRunningPulsCountAut=PulsCount
	savePsyBaseMoodAut=PsyBaseMood
	savePsyMoodAut=PsyMood
}
///////////////////////////////////
func clinerAutomatizmRunningAut(){
	AutomatizmRunningAut=nil
	AutomatizmRunningPulsCountAut=0
}
///////////////////////////////////////////////////////////////////
// При любых изменениях wasChangingMoodCondition() оценивать действие запущенного автоматизма
func calcAutomatizmResultAut(diffPsyBaseMood int,wellIDarr []int){
	if AutomatizmRunningPulsCountAut==0 || AutomatizmRunningAut==nil{
		clinerAutomatizmRunningAut()
		return
	}
	// diffPsyBaseMood - точно изменился, иначе бы не было вызова calcAutomatizmResult
	/// если числа имеют разные знаки (одно положительное, другое отрицательное)
	if lib.IsDiffersOfSign(AutomatizmRunningAut.Usefulness,diffPsyBaseMood)		{
		AutomatizmRunningAut.Count--
	} else {
		AutomatizmRunningAut.Count=0
	}
	SetAutomatizmBelief(AutomatizmRunning,2)
	// ТАК НЕЛЬЗЯ ЗАДАВАТЬ Belief=2: AutomatizmRunning.Belief=2
	AutomatizmRunningAut.Usefulness =diffPsyBaseMood

	if diffPsyBaseMood<0{// СТАЛО ХУЖЕ В ПРИВЫЧНОМ АВТОМАТИЗМЕ !!!! криминал
		PsyBaseMood=-1
		/* значит есть технически невидимая новизна и нужно ОСМЫСЛИТЬ ТАКУЮ СИТУАЦИЮ

		 */
	}
	if diffPsyBaseMood>0{// стало лучше
		PsyBaseMood=1
	}
	if EvolushnStage > 3 {
		// создать образ ситуации
		autmzmTreeNodeID := AutomatizmRunning.BranchID
		id, _ := createSituationImage(0, autmzmTreeNodeID, 2)
		// осмыслить ситуацию - Активировать Дерево Понимания
		understandingSituation(id, savePurposeGenetic)
	}

	// оценить значимость поизнесенной фразы в VerbalFromIdArr структурах Дерева Понимания??

	/* !!!!!! допонение cerebellumReflexFromID[LastAutomatizmWeiting.ID] другими корректирующими действиеями
	   если это еще получается, но при отсуствии эффекта нужно создавать новый автоматизм.
	   Это - только на уровне осмысления в Дереве Понимания:
	      cerebellumCoordination(AutomatizmRunning.ID)
	   Должна быть осознание цели и перебеора-недобора!!!!!!
	      В каждом автоматизме есть параметр силы: Automatizm.Energy вот он и корректируется.
	*/

	clinerAutomatizmRunningAut()
	return
}
//#############################################################################
// конец неработающего кода







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


//////////////////////////////////////////////////////////////////
// для индикации период ожидания реакции оператора на действие автоматизма
//   psychicWaitingPeriodForActions()
func WaitingPeriodForActions()(bool,int){

	if LastRunAutomatizmPulsCount>0{
		time:=WaitingPeriodForActionsVal - (PulsCount-LastRunAutomatizmPulsCount)
		return true,time
	}
	if AutomatizmRunningPulsCountAut>0{
		time:=WaitingPeriodForActionsVal - (PulsCount-AutomatizmRunningPulsCountAut)
		return true,time
	}
return false,0
}
