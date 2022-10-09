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
		// срочность и важность ситуации: если очень срочно и важно - просто оставить имеющийся автоматизм
		runAutomatizmFromPurpose(atmzm, purpose)
		return atmzm
	}else{//....................................... не очень актуально
		// проверенный и хороший автоматизм не трогать
		if  atmzm.Belief == 2 && atmzm.Usefulness>0{
			//просто выполнить автоматизм и отслеживать последствия
			runAutomatizmFromPurpose(atmzm, purpose)
			return atmzm
		}
	}
		// плохой автоматизм, попробовать оптимизировать с помощью рефлексов МОЗЖЕЧКА
		if atmzm.Usefulness < 0{
			// была ли уже оптимизация?
var cr *cerebellumReflex
			if atmzm.Belief == 3{// была неудачная попытка оптимизации (т.к. atmzm.Usefulness < 0)
cr=cerebellumReflexFromMotorsID[atmzm.ID]
if cr!=nil { // уже имеющийся мозжечковый рефлекс
	// в прошлый раз было добавлено энергии
	oldEnerg:=cr.addEnergy
	if oldEnerg+atmzm.Energy >9{// был добавлено лишее, нужно постепенно уменьшить, а не скакать от плюса к минусу
		cr.addEnergy-- // может быть отрицательное число, чтобы уменьшить энергию автоматизма
	}
	if oldEnerg+atmzm.Energy <1{// был добавлено лишее, нужно постепенно добавлять, а не скакать от плюса к минусу
		cr.addEnergy++
	}
	SaveCerebellumReflex()// может не стоит, пусть записывается вместе со всеми сохранениями
}
			}else{// еще не было мозжечкового рефлекса
// добавить
				_,cr=createNewCerebellumReflex(0,0,atmzm.ID)
				if cr !=nil {
					// добавить энергичность по максимум
					cr.addEnergy=10-atmzm.Energy // хотя при выполнении автоматизма будет подрезана лишняя энергия
					atmzm.Belief = 3
					SaveCerebellumReflex()// может не стоит, пусть записывается вместе со всеми сохранениями
				}
			}
if cr!=nil {
	runAutomatizmFromPurpose(atmzm, purpose)
	return atmzm
}else{
	// просто не выполнять плохой автоматизм, раз что-то не так с рефлексами мозжечка
	return nil
}
		}

//if newsRes{// повышенная опасность от оператора
}else{
	/* от оператора нет опасности, но высокий purpose.veryActual,
	нужно выполнить штатный автоматизм, а не придуманный
	 */
	runAutomatizmFromPurpose(atmzm, purpose)
	return atmzm
}
//if purpose.veryActual
}else{// нет опасности и нет опасной новизны
		// плохой автоматизм,
		if atmzm.Usefulness < 0 {
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

//просто выполнить автоматизм и отслеживать последствия
	runAutomatizmFromPurpose(atmzm, purpose)
	return atmzm
}

	return nil
}
/////////////////////////////////////////////////////////////////