/* Процессы осмысления: создание и использование ментальных автоматизмов
для Дерева понимания (или дерева ментальных автоматизмов)

*/

package psychic

///////////////////////////////


//////////////////////////////////////////////////
/* Определение Цели в данной ситуации - ну уровне дерева понимания
Здесь выбирается действие пробного автоматизма из выполнившегося рефлекса actualRelextActon.
*/
func getPurposeUndestandingAndRunAutomatizm()(*Automatizm) {
var atmzm *Automatizm

	// мозжечковые рефлексы - самый первый уровень осознания - подгонка действий под заданную Цель.
	if EvolushnStage == 4 {
		/*  на стадии 4 - провоцировать оператора на ответы (почему, зачем, что такое?)


		 */
	}
/*
	// осмыслить ситуацию - Активировать Дерево Понимания
	autmzmTreeNodeID:=AutomatizmRunning.BranchID// создать образ ситуации
	id,_:=createSituationImage(0,autmzmTreeNodeID,4,true)
	// осмыслить ситуацию - Активировать Дерево Понимания
	understandingSituation(id,purpose)
	и затем создать новую цель understanding_purpose_image.go
*/


	return atmzm
}
////////////////////////////////////////////////



// обработка автоматизма, рвущегося на выполнение, но в условиях есть новизна news
func getPurposeUndestanding2AndRunAutomatizm(atmtzmID int)(*Automatizm){

	atmzm:=AutomatizmFromIdArr[atmtzmID]

	// есть ли очень значимые новые признаки?
	newsRes:=getImportantSigns()
	if newsRes{

		// создать образ ситуации
		autmzmTreeNodeID:=detectedActiveLastNodID
		id,_:=createSituationImage(0,autmzmTreeNodeID,5,true)
		// осмыслить ситуацию - Активировать Дерево Понимания
		understandingSituation(id,savePurposeGenetic)
		// м.б. были совершены действия, тогда isReflexesActionBloking был выставлен

	}else{ // нет новизны

		// TODO м.б. есть сила воли, блокирующая выполнение данного автоматизма и подставляющая свой?
		// просто выполнить автоматизм НО М.Б. задать purpose:=getPurposeGenetic() а не ставить nil ?
		runAutomatizmFromPurpose(atmzm, nil)
	}

	// мозжечковые рефлексы - самый первый уровень осознания - подгонка действий под заданную Цель.

	return atmzm // TODO с учетом активации дерева понимания !!!!



	/// !!!! если не найдено и нельзя выполнять то return nil
	if false {
		isReflexesActionBloking = false
		return nil
	}

	////////////  если нет результата - выполнить этот автоматизм
	return atmzm
}
//////////////////////////////////////////////////////