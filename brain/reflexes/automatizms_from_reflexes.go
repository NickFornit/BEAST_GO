/* Создать автоматизмы на основе существующих рефлексов

Запускается по ссылке не ранее 4-го пульса, так что должны быть готовы все массивы.

Для тестирования возможно избежать долгий период воспитания с формированием автоматизмов и просто сгенерировать автоматизмы на основе существующих рефлексов (с приоритетом условных рефлексов).
При этом у автоматизмов будут установлены опции уже проверенного автоматизма с полезностью, равной 1 (вполне полезно). Это правомерно потому, что рефлексы создавались уже полезными для своих условий.
В дальнейшем такие автоматизмы будут проверяться в зависимости от реакции оператора и изменения состояния Beast, корректируясь настолько эффективно, насколько позволяет текущая стадия развития.
*/

package reflexes

import (
	"BOT/brain/psychic"
	wordSensor "BOT/brain/words_sensor"
	"sort"
	"strconv"
)

func testingRunMakeAutomatizmsFromReflexes(){
	//  RunMakeAutomatizmsFromReflexes()
}



////////////////////////////////////////////
/* сканировать для всех условий у.рефлексы (и если нет  - смотреть б.рефлексы),
создавать ветку дерева автоматизма если такой еще нет,
создавать автоматизм, прикрепляя его к нужно ветке.
 */
func RunMakeAutomatizmsFromReflexes()(string){
	// проверить готовность рабочих массивов и сообщить если нет
	if ConditionReflexes == nil || len(ConditionReflexes)==0 ||
		psychic.AutomatizmTreeFromID == nil || len(psychic.AutomatizmTreeFromID)==0{
		return "Еще не сформировалась оперативная память, пожалуйста перезапустите процесс через пару секунд."
	}
	var newCount=0
// сортировка по ID чтобы тестировалось воспроизводимо
	keys := make([]int, 0, len(ConditionReflexes))
	for k := range ConditionReflexes {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _,id:= range keys { v:=ConditionReflexes[id]
	//for _, v := range ConditionReflexes {

		//      v=ConditionReflexes[3673]
		if newCount>20{
			continue // чтобы записала файлы      return ""
		}

		/* поиск узла дерева автоматизмов по условиям у.рефлекса
		Если нет такого узла - дорастить ветку.
		Выдать  ID узла
		*/
		actID:=TriggerStimulsArr[v.lev3]
		tm:=psychic.GetToneMoodID(actID.ToneID,actID.MoodID)// тон-настроение
		verbalID:=actID.PhraseID// фраза VerbalID

	//	sss:=wordSensor.GetWordFromPraseNodeID(verbalID[0]); if len(sss)>0{}

		FirstSimbolID:=wordSensor.GetFirstSymbolFromPraseID(verbalID[0])
		// создать образ Брока
		psychic.CreateVerbalImage(FirstSimbolID, verbalID, actID.ToneID, actID.MoodID)

		nodeID := psychic.FindConditionsNode(v.lev1, v.lev2, actID.RSarr,tm,FirstSimbolID,verbalID[0])
		if nodeID>0{
// если есть привязанный к узлу автоматизм, то не привязывать еще
if psychic.ExistsAutomatizmForThisNodeID(nodeID) {
	continue
}

			//  создать автоматизм и привязать его к nodeID
			var sequence="Dnn:"
			aArr:=v.ActionIDarr
				for i := 0; i < len(aArr); i++ {
					if i > 0 {
						sequence += ","
					}
					sequence += strconv.Itoa(aArr[i])
			}
			_,autmzm:=psychic.CreateAutomatizm(nodeID,sequence)
			if autmzm!=nil{
				autmzm.Belief=2 // сделать автоматизм штатным
				// ?? autmzm.GomeoIdSuccesArr какие ID гомео-параметров улучшает это действие
				autmzm.Usefulness=1 // полезность

				newCount++
			}
		}
	}

	psychic.SaveAllPsihicMemory()

return "Процесс нормально завершен, создано "+strconv.Itoa(newCount)+" новых автоматизмов."
}
/////////////////////////////////////////


//





