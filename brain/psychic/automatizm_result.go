/* Ожидание результата запущенного автоматизма и его обработка

*/


package psychic

import (
	"BOT/brain/gomeostas"
	"BOT/lib"
)

/////////////////////////////////////////
var savePsyBaseMood=0 // -1 Плохое настроение, 0 Нормальное, 1 - хорошее настроение
// для более точной оценки
var savePsyMood=0//сила Плохо -10 ... 0 ...+10 Хорошо
// НОВИЗНА СИТУАЦИИ сохраняется значение CurrentAutomatizTreeEnd[] для решений
var savedNoveltySituation []int


//////////////////////// ПУЛЬС
// ПУЛЬС срабатывания по каждому Пульсу здесь для удобства
func automatizmActionsPuls(){

	// если был запущен автоматизм, возможно, без ориентировчного рефлека, рефлексы блокируются на 2 пульса
	if MotorTerminalBlocking && LastRunAutomatizmPulsCount+2 > PulsCount {
		MotorTerminalBlocking=false // снять блокировку
	}

/* Ожидание результата автоматизма ПОСЛЕ ОРИЕНТИРОВОЧНОГО РЕФЛЕКСА.
Реакция ожидается на слелующем пульcе после срабатывания автоматизма	и в течение 20 пульсов
 за это время получим уверенное wasChangingMoodCondition() по значению gomeostas.BetterOrWorseNow()
 */
	if isPeriodResultWaiting {
		if (AutomatizmRunningPulsCount+1) < PulsCount {// следить со следующего пульса
			// Из МОЗЖУЧКА как-то отреагировать на отсуствие реакции - повторить автоматизм с большей силой Energy
			/*if noAutovatizmResult() { // была попытка отреагировать сильнее - в cerebellum.go
				return
			}*/
			res, wellIDarr := wasChangingMoodCondition()
			if res != 0 { // ИЗМЕНИЛОСЬ СОСТОЯНИЕ
				// обработать изменение состояния
				calcAutomatizmResult(res, wellIDarr)
				//  clinerAutomatizmRunning()  есть в calcAutomatizmResult
			}
		}
		if AutomatizmRunningPulsCount+20 < PulsCount {
			//сбрасывать ожидание результата автоматизма если прошло 20 пульсов
			clinerAutomatizmRunning()
		}
	}
	//////////////////////////////////////////////


/* ПОКА НИКАК НЕ ИСПОЛЬЗУЕТСЯ
   Ожидание результата автоматизма БЕЗ ОРИЕНТИРОВОЧНОГО РЕФЛЕКСА (автоматически запущенного из Дерева).
	Реакция ожидается на слелующем пульcе после срабатывания автоматизма	и в течение 20 пульсов
	 за это время получим уверенное wasChangingMoodCondition() по значению gomeostas.BetterOrWorseNow()
*/
	if (AutomatizmRunningPulsCountAut+1)<PulsCount && AutomatizmRunningPulsCountAut+20 < PulsCount{
		// Из МОЗЖУЧКА как-то отреагировать на отсуствие реакции - повторить автоматизм с большей силой Energy
		if noAutovatizmResult(){// была попытка отреагировать сильнее - в cerebellum.go
			return
		}
		res,wellIDarr:=wasChangingMoodCondition()
		if res!=0{// ИЗМЕНИЛОСЬ СОСТОЯНИЕ
			// обработать изменение состояния
			calcAutomatizmResultAut(res,wellIDarr)
		}
	}
	if AutomatizmRunningPulsCountAut+20 > PulsCount {
		//сбрасывать ожидание результата автоматизма если прошло 20 пульсов
		clinerAutomatizmRunningAut()
	}
}
/////////////////////////////////////////////////////////////////////


/////////////////////////////////////////////////////////////////////
// отслеживание осознанно запущенных автоматизмов
var AutomatizmRunning *Automatizm // запущенный автоматизм
var AutomatizmRunningPulsCount=0 // время запуска автоматизма
var savePurposeGenetic *PurposeGenetic // массив примитивных целей

func setAutomatizmRunning(am *Automatizm,ps *PurposeGenetic){
	isPeriodResultWaiting=true // игнорировать новое внимание на время ожидания результата автоматизма
	AutomatizmRunning=am
	AutomatizmRunningPulsCount=PulsCount
	savePsyBaseMood=PsyBaseMood
	savePsyMood=PsyMood
	savedNoveltySituation=NoveltySituation
	savePurposeGenetic=ps
}
///////////////////////////////////
func clinerAutomatizmRunning(){
	AutomatizmRunning=nil
	AutomatizmRunningPulsCount=0
	isReflexesActionBloking=false
	isPeriodResultWaiting=false
}
///////////////////////////////////



/* ПОСЛЕ ОРИЕНТИРОВОЧНОГО РЕФЛЕКСА оценивать действие запущенного автоматизма

 */
func calcAutomatizmResult(diffPsyBaseMood int,wellIDarr []int){
	if AutomatizmRunningPulsCount==0 || AutomatizmRunning==nil{
		clinerAutomatizmRunning()
		return
	}
	// diffPsyBaseMood - точно изменился, иначе бы не было вызова calcAutomatizmResult
	/// если числа имеют разные знаки (одно положительное, другое отрицательное)
	if lib.IsDiffersOfSign(AutomatizmRunning.Usefulness,diffPsyBaseMood)		{
		AutomatizmRunning.Count--
	} else {
		AutomatizmRunning.Count=0
	}
	setAutomatizmBelief(AutomatizmRunning,2)
	// ТАК НЕЛЬЗЯ ЗАДАВАТЬ Belief=2: AutomatizmRunning.Belief=2
	AutomatizmRunning.Usefulness =diffPsyBaseMood

	if diffPsyBaseMood<0{// стало хуже
		PsyBaseMood=-1
	}
	if diffPsyBaseMood>0{// стало лучше
		PsyBaseMood=1
		// список гомео параметро, которые улучшило это действие
		AutomatizmRunning.GomeoIdSuccesArr=wellIDarr
		// пополняется список полезных автоматизмов
		AutomatizmSuccessFromIdArr[AutomatizmRunning.ID]=AutomatizmRunning
	}

	if diffPsyBaseMood!=0{// изменилась ситуация
		// обновить информационное окружение
		GetCurrentInformationEnvironment()
		// переактивировать дерево рефлексов
		automatizmTreeActivation()//и возникает новый цикл активации Дерева, уже по внутренним причинам
	}
// создать образ ситуации
autmzmTreeNodeID:=AutomatizmRunning.BranchID
	id,_:=createSituationImage(0,autmzmTreeNodeID,1)
	// осмыслить ситуацию - Активировать Дерево Понимания
	understandingSituation(id,savePurposeGenetic)

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









/////////////////////////////////////////////////////////////////////////
/* сканируется с каждым пульсом в func automatizmActionsPuls() во время ожидания
В  gomeostas.BetterOrWorseNow() учитывается CommonMoodAfterAction - Общее (де)мотивирующее действие с Пульта
 */
func wasChangingMoodCondition()(int,[]int){
	//стало хуже или лучше теперь, возвращает величину измнения от -10 через 9 до 10
	res,wellIDarr:=gomeostas.BetterOrWorseNow()

	return res,wellIDarr
}
/////////////////////////////////////////////////////////////////////////







/////////////////////////////////////////////////////////////////////
/*  ПОКА НИКАК НЕ ИСПОЛЬЗУЕТСЯ ?? TODO
Отслеживание запущенных автоматизмов ВНЕ ОРИЕНТИРОВОЧНОГО РЕФЛЕКСА.
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
/* При любых изменениях wasChangingMoodCondition() оценивать действие запущенного автоматизма

 */
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
	setAutomatizmBelief(AutomatizmRunning,2)
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
	// создать образ ситуации
	autmzmTreeNodeID:=AutomatizmRunning.BranchID
	id,_:=createSituationImage(0,autmzmTreeNodeID,2)
	// осмыслить ситуацию - Активировать Дерево Понимания
	understandingSituation(id,savePurposeGenetic)

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
/////////////////////////////////////////////////////////////////////////








