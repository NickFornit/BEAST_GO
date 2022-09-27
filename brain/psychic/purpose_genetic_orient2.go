/*   
Для ориентировочного рефлекса типа 2
функции для определения Цели в данной ситуации - на уровне наследственных функций
исходя из текущей информационной среды CurrentInformationEnvironment:

*/

package psychic



///////////////////////////////////////////////

// обработка автоматизма, рвущегося на выполнение, но в условиях есть новизна news
func getPurposeGenetic2AndRunAutomatizm(atmtzmID int)(*Automatizm){

	atmzm:=AutomatizmFromIdArr[atmtzmID]

	// Определение Цели в данной ситуации - на уровне наследственных функций
	purpose:=getPurposeGenetic()
	// мозжечковые рефлексы - самый первый уровень осознания - подгонка действий под заданную Цель.

	// есть ли очень значимые новые признаки?
	newsRes:=getImportantSigns()

	if purpose.veryActual{// нужно ли вообще шевелиться?

if newsRes{// повышенная опасность от оператора
	if purpose.veryActual{

	}else{
		
	}

}



		// срочность и важность ситуации: если очень срочно и важно - просто оставить имеющийся автоматизм
		if RumAutomatizm(atmzm) {// запустить автоматизм
			// отслеживать последствия в automatizm_result.go ИГНОРИРОВАТЬ ДРУГИЕ ПОКА НЕ БУДЕТ РЕЗУЛЬТАТ
			setAutomatizmRunning(atmzm,purpose)
		}
		// в automatizm_result.go после оценки результата будет осмысление с активацией Дерева Понимания
		return atmzm
	}
	////////////////////////////

	// если непроверенный автоматизм
	if  AutomatizmFromIdArr[atmtzmID].Belief!=2{

		// TODO

	}

	// ПЛОХОЙ АВТОМАТИЗМ ?!
	if atmzm.Usefulness < 0{ // плохой автоматизм, т.е. в прошлый раз был плохой результат
		// м.б. проэкспериментировать с мозжечковыми рефлексами - подгонка действий под заданную Цель.

		// если найден лучше автоматизм, то подставить его: atmzm=xxxx и запустить

		//return automatizm
	}else {// приемлемый по последствиям автоматизм


	}





	// НЕТ СРОЧНОСТИ, МОЖНО СПОКОЙНО ПОДУМАТЬ или тупо поэкспериментировать
	/* найти важные (по опыту) признаки в новизне news
	   Это - чисто рефлексторный процесс поиска в опыте
	*/
	if newsRes {
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


	/// !!!! если не найдено и нельзя выполнять то return nil
	if false {
		isReflexesActionBloking = false
		return nil
	}

	////////////  если нет результата - выполнить этот автоматизм
	return atmzm
}
/////////////////////////////////////////////////////////////////