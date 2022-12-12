/*   
Для ориентировочного рефлекса типа 2
функции для определения Цели в данной ситуации - на уровне наследственных функций
исходя из текущей информационной среды CurrentInformationEnvironment:

*/

package psychic

import "BOT/brain/gomeostas"

///////////////////////////////////////////////
// обработка автоматизма, рвущегося на выполнение, но в условиях может быть новизна news
/* Здесь - очень органиченные возможности адаптации автоматизма:
плохой - не выполнять, хороший - выполнять
при опасной ситуации выполнять тот, какой есть,
при спокойной ситуации - пробовать рефлексы мозжечка.

Из-за столь скудных возможностей и разросся функционал мыслительных автоматизмов
с их произвольностью (- перекрытием имеющихся автоматизмов новыми).
 */
var oldNodeAutomatizm=0 // прошлы раз запускался такой штатный автоматизм
func getPurposeGenetic2AndRunAutomatizm(atmtzmID int)(*Automatizm){

	atmzm:=AutomatizmFromIdArr[atmtzmID]

	// Определение Цели в данной ситуации - на уровне наследственных функций
	purpose:=getPurposeGenetic()
	// мозжечковые рефлексы - самый первый уровень осознания - подгонка действий под заданную Цель.

	// есть ли очень значимые новые признаки?
	newsRes:=getImportantSigns()
	if newsRes { // повышенная опасность от оператора
		// срочность и важность ситуации: если очень срочно и важно - просто оставить имеющийся автоматизм
		runAutomatizmFromPurpose(atmzm, purpose)
		return atmzm
	}

	if purpose.veryActual{// нужно ли вообще шевелиться?
	// высокий purpose.veryActual, 	нужно выполнить штатный автоматизм, а не придуманный

	// если в прошлый раз уже был такой автоматизм, то ничего не делать, иначе начинает повторять одно и то же
if oldNodeAutomatizm != atmzm.ID {
	runAutomatizmFromPurpose(atmzm, purpose)
	oldNodeAutomatizm = atmzm.ID
	return atmzm
}

	// список всех автоматизмов для ID узла Дерева
	//aArr:=GetMotorsAutomatizmListFromTreeId(detectedActiveLastNodID)

	/*aID := getAutomatizmFromNodeID(detectedActiveLastNodID)
	atmzm=AutomatizmFromIdArr[aID]
	purpose.actionID=ActionsImageArr[atmzm.ActionsImageID]
	runAutomatizmFromPurpose(atmzm, purpose)*/
	return nil

//if purpose.veryActual
}else {              // нет опасности и нет опасной новизны

		// плохой автоматизм,
		if atmzm.Usefulness < 0 {
			// была ли уже оптимизация?
			if cerebellumCoordination(atmzm, 0) {
				runAutomatizmFromPurpose(atmzm, purpose)
				return atmzm
			} else {
				if gomeostas.BaseContextActive[2] || gomeostas.BaseContextActive[3] { // активен Поиск или Игра
					// тупо метод тыка
					// Тупо поэкспериментировать для пополнения опыта (не)удачных автоматизмов
					// TODO !не проверено!
					// в отличии от createAndRunAutomatizmFromPurpose(purpose) не использовать текущие рефлексы, а пробовать всякое
					// Выдавая это на стадии 3, тварь получает реакцию оператора, которую отзеркаливает
					atmzm := findAnySympleRandActions()
					return atmzm
				} else { // НЕ ИГРА  И НЕ ПОИСК, плохой автоматизм просто не выполнять
					return nil
				}
			}
		}

		//все нормально, просто выполнить автоматизм и отслеживать последствия
		runAutomatizmFromPurpose(atmzm, purpose)
		return atmzm
	}

	return nil
}
/////////////////////////////////////////////////////////////////