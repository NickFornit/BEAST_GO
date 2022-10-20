/* Ориентировочный рефлекс осознания ситуации.

В любом случае orientationConsciousness(kind int) - осмысление ситуации SituationImage

Попытка подставить вместо текущего автоматизма более правильный в новых условиях - ор.рефлекс в дереве
Оценка результата автоматизма и корректировка мозжечковыми рефлексами
ПРИ ОСОЗНАНИИ ЦЕЛИ И ПЕРЕБОРА_НЕДОБРА, А ТАК ЖЕ НЕОБХОДИМОСТИ ДОПОЛНИТЕЛЬНЫХ (корректирующих)ДЕЙСТВИЙ.
*/

package psychic

import "BOT/lib"

///////////////////////////////////////



/* kind:
0 - вообще еще нет ветки с нуля - полная новизна условий
1 - частино есть ветка - частичная новизна условий
2 - ветка полностью есть, только новизна ситуации

Доступная Инфа:
curActiveActions - структура действий при активации дерева автоматизмов
savePurposeGenetic
saveSituationImageID
currentMentalAutomatizmID
LastRunAutomatizmPulsCount - идет период ожидания ответа
LastAutomatizmWeiting - ожидается результат запущенного MotAutomatizm

// нераспознанный остаток - НОВИЗНА: var CurrentUnderstandingTreeEnd []int

return true - были совершены моторные действия - заблокировать все более низкоуровневые действия
*/
func orientationConsciousness(kind int)(bool){
	lib.WritePultConsol("Ориентировчный рефлекс Дерева понимания.")

/* kind ==1: в активной ветке нет ментального автоматизма - нужна субъектиная оценка ситуации новизны
   НЕТ УСЛОВИЙ ДЛЯ ПРОЯВЛЕНИЯ ПРОИЗВОЛЬНОСТИ, найти решение,
особенно не заморачиваясь (первое приближение).
 */
	if kind ==1 {
		if consciousness(1){
			return true
		}
/* НОВИЗНА: CurrentUnderstandingTreeEnd=condArr[currentUnderstandingStepCount:] []int
- нет еще таких веток, начиная с currentUnderstandingStepCount

	// новый кадр эпизодической памяти, сохраняющий
   	newEpisodeMemory()  // запись эпизодической памяти saveEpisodicMenory()

   	// есть ли очень значимые новые признаки?
   	newsRes:=getImportantSigns()
 */

// Если для данного случая есть подходящая доминанта, то - перейти к ее решению
		// TODO

		/* Определение Цели в данной ситуации - на уровне дерева понимания
		Здесь выбирается действие пробного автоматизма из выполнившегося рефлекса actualRelextActon
		и запускается МЕНТАЛЬНЫЙ автоматизм.
		На стадии 4 - провоцировать оператора на ответы (почему, зачем, что такое?)
		*/
		res := getPurposeUndestandingAndRunAutomatizm()
		return res
	}
	///////////////////////////////////////////////
/* kind==2: в активной ветке ЕСТЬ ментальный автоматизм, который нужно проверить перед выполненением
	 штатный мент.автоматизм, привязанный к ветке: currentMentalAutomatizmID
   УСЛОВИЯ ДЛЯ ПРИМЕНЕНИЯ ПРОИЗВОЛЬНОСТИ: старый мент.автоматизм может быть заменен на новый.
 */
	if kind==2 {// условия применения ПРОИЗВОЛЬНОСТИ
		if consciousness(2){
			return true
		}
/* нужна субъектиная (отвлеченная от гомеостаза) оценка ситуации -
  с учетом новизны ситуации, которая может не детектироваться активными ветками деревьев,
  особенно это касается понимания ФРАЗЫ сообщения.
Потому что активность ветки дерева ментальных автоматизмов не связана
	с активностью ветки дерева моторных автоматизмов,
	начиная с самого первого уровня (Mood - Ощущаемое настроение, а не Базовое состояние).

   если есть привязаный автоматизм, то в ор.рефлексе второго типа смотрится:
   если в конце цепочки цель достигнута, то просто запускать такой автоматизм,
		если нет, то найти один из вариантов: заменить новым автоматизмом,
		дорастить следующим или если уже есть следующий - ответвить.
		Метод: проход всей цепочки автоматизмов по одной для прогноза
		с попытками модернизации на основе текущего опыта и возможностей.
 */


/*Если нет настораживающих признаков, и цели в результате цепочки реакции соотвествуют желаемому,
		то такая реакция запускается без осмысления (фактически это редуцированная ориентировочная реакция).
 */
		if true{// TODO не нужно осмысление
			RunMentalMentalAutomatizmsID(currentMentalAutomatizmID)
			return true
		}
// нужно привлечение внимания, осмысление
			// новый кадр эпизодической памяти, сохраняющий
// вызывается только в func calcAutomatizmResult:	newEpisodeMemory()  // запись эпизодической памяти saveEpisodicMenory()

      	// есть ли очень значимые новые признаки?
      //	newsRes:=getImportantSigns()


		// Если для данного случая есть доминанта, то - перейти к ее решению
		// TODO

		// запустить
		if true {// TODO условие запуска
			RunMentalMentalAutomatizmsID(currentMentalAutomatizmID)
			return true
		}
	}

return false // - ничего не произвошло
}
////////////////////////////////////


