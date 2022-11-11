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

func testingRunMakeAutomatizmsFromReflexes() {
	// RunMakeAutomatizmsFromReflexes()
	// RunMakeAutomatizmsFromGeneticReflexes()
}

/* сканировать для всех условных рефлексов,
создавать ветку дерева автоматизма если такой еще нет,
создавать автоматизм, прикрепляя его к нужно ветке.
 */
func RunMakeAutomatizmsFromReflexes() string {
	// проверить готовность рабочих массивов и сообщить если нет
	if ConditionReflexes == nil || len(ConditionReflexes) == 0 ||
		psychic.AutomatizmTreeFromID == nil || len(psychic.AutomatizmTreeFromID) == 0 {
		return "Еще не сформировалась оперативная память, пожалуйста перезапустите процесс через пару секунд."
	}
	var newCount = 0
	var count = 0
	// сортировка по ID чтобы тестировалось воспроизводимо
	keys := make([]int, 0, len(ConditionReflexes))
	for k := range ConditionReflexes {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _,id:= range keys {
		v := ConditionReflexes[id]
		// for _, v := range ConditionReflexes {
		//      v=ConditionReflexes[3673]
		// для проверки
		//	if count>6{
		//	psychic.SaveAllPsihicMemory();
		//	return ""; }

		/* поиск узла дерева автоматизмов по условиям у.рефлекса
		Если нет такого узла - дорастить ветку.
		Выдать ID узла
		*/
		actID := TriggerStimulsArr[v.lev3]
		tm := psychic.GetToneMoodID(actID.ToneID, actID.MoodID) // тон-настроение
		verbalID := actID.PhraseID // фраза VerbalID
		// s:=wordSensor.GetPhraseStringsFromPhraseID(verbalID[0]);if len(s)>0{}
		FirstSimbolID := wordSensor.GetFirstSymbolFromPraseID(verbalID)
		// создать образ Брока
		psychic.CreateVerbalImage(FirstSimbolID, verbalID, actID.ToneID, actID.MoodID)
		nodeID := psychic.FindConditionsNode(v.lev1, v.lev2, actID.RSarr, tm, FirstSimbolID, verbalID[0])
		//,,,,,,,,,,,,,,,,, для проверки
		//		psychic.SaveAllPsihicMemory() // чтобы сразу видеть какой узел возник
		/*
				lastNode:=psychic.AutomatizmTreeFromID[nodeID]; if lastNode!=nil{}
				if lastNode.VerbalID == 0{
					continue
				}
		*/
		//,,,,,,,,,,,,,,,,,
		if nodeID > 0 {
			// если есть привязанный к узлу автоматизм, то не привязывать еще
			exists := psychic.ExistsAutomatizmForThisNodeID(nodeID)
			if exists {
				//,,,,,,,,,,,,,,,,, для проверки
				//	aArr:=psychic.AutomatizmBelief2FromTreeNodeId[nodeID];if aArr!=nil{}
				count++
				continue
			}

			psychic.NoWarningCreateShow = true
			ActionsImageID,_:=psychic.CreateNewlastActionsImageID(0,v.ActionIDarr,nil,0,0,true)
			_, autmzm := psychic.CreateAutomatizm(nodeID, ActionsImageID)
			psychic.NoWarningCreateShow = false
			if autmzm != nil {
				psychic.SetAutomatizmBelief(autmzm, 2) // сделать автоматизм штатным
				// ?? autmzm.GomeoIdSuccesArr какие ID гомео-параметров улучшает это действие
				autmzm.Usefulness = 1 // полезность
				count++
				newCount++
			}
		}
	}
	psychic.SaveAllPsihicMemory()
	return "Процесс нормально завершен, создано " + strconv.Itoa(newCount) + " новых автоматизмов."
}

/* сканировать для всех безусловных рефлексов,
создавать ветку дерева автоматизма если такой еще нет,
создавать автоматизм, прикрепляя его к нужно ветке.
*/
func RunMakeAutomatizmsFromGeneticReflexes() string {
	// проверить готовность рабочих массивов и сообщить если нет
	if GeneticReflexes == nil || len(GeneticReflexes) == 0 ||
		psychic.AutomatizmTreeFromID == nil || len(psychic.AutomatizmTreeFromID) == 0 {
		return "Еще не сформировалась оперативная память, пожалуйста перезапустите процесс через пару секунд."
	}
	var newCount = 0
	var count = 0
	// сортировка по ID чтобы тестировалось воспроизводимо
	keys := make([]int, 0, len(GeneticReflexes))
	for k := range GeneticReflexes {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, id := range keys {
		v := GeneticReflexes[id]
		//for _, v := range GeneticReflexes {

		//	v=GeneticReflexes[3673]
		// для проверки
		//	if count>6{
		//	psychic.SaveAllPsihicMemory();
		//	return ""; }

		/* поиск узла дерева автоматизмов по условиям у.рефлекса
		Если нет такого узла - дорастить ветку.
		Выдать  ID узла
		*/
		nodeID := psychic.FindConditionsNode(v.lev1, v.lev2, v.lev3,0,0,0)
		//,,,,,,,,,,,,,,,,, для проверки
		//		psychic.SaveAllPsihicMemory() // чтобы сразу видеть какой узел возник
		/*
			lastNode:=psychic.AutomatizmTreeFromID[nodeID]; if lastNode!=nil{}
			if lastNode.VerbalID == 0{
				continue
			}
		*/
		//,,,,,,,,,,,,,,,,,
		if nodeID > 0 {
			// если есть привязанный к узлу автоматизм, то не привязывать еще
			exists := psychic.ExistsAutomatizmForThisNodeID(nodeID)
			if exists {
				//,,,,,,,,,,,,,,,,, для проверки
				//	aArr:=psychic.AutomatizmBelief2FromTreeNodeId[nodeID];if aArr!=nil{}
				count++
				continue
			}
			//  создать автоматизм и привязать его к nodeID
			psychic.NoWarningCreateShow = true
			ActionsImageID,_:=psychic.CreateNewlastActionsImageID(0,v.ActionIDarr,nil,0,0,true)
			_, autmzm := psychic.CreateAutomatizm(nodeID, ActionsImageID)
			psychic.NoWarningCreateShow = false
			if autmzm != nil {
				psychic.SetAutomatizmBelief(autmzm, 2) // сделать автоматизм штатным
				// ?? autmzm.GomeoIdSuccesArr какие ID гомео-параметров улучшает это действие
				autmzm.Usefulness = 1 // полезность
				count++
				newCount++
			}
		}
	}
	psychic.SaveAllPsihicMemory()
	return "Процесс нормально завершен, создано " + strconv.Itoa(newCount) + " новых автоматизмов."
}