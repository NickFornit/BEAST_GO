/*  Сформировать условные рефлексы на основе списка фраз-синонимов */

package reflexes

import (
	wordSensor "BOT/brain/words_sensor"
	"strconv"
	"strings"
)

/* Затрагивает файлы памяти:
condition_reflexes.txt
trigger_stimuls_images.txt
word_tree.txt
phrase_tree.txt

т.к. вставляются новые слова и фразы в вербальные деревья, формируются образы пусковых фраз
 */

/*  Сформировать условные рефлексы на основе списка фраз-синонимов */
func FormingConditionsRefleaxFromList(list string) string {
	// тестовый запуск в tree_activation.go readyForRecognitionRflexes()
 // list="4774|1|1,2,8||18,50,76,79|ненавижу\r\n4775|1|1,2,8|1|17,50,76,79,80|нет\r\n4777|1|1,2,8|1,4|14,17,50,76,79|ну\r\n4780|1|1,2,8|1,7|9,17,50,76,79,80|понятно\r\n4783|1|1,2,8|1,10|38,50,73,76,79,80|смеюсь"

	if len(list) < 5 { return "Пустой файл" }

	strArr := strings.Split(list, "\r\n")
	for n := 0; n < len(strArr); n++ {
		p:=strings.Split(strArr[n], "|")
		lev1, _ := strconv.Atoi(p[1])
		// второй уровень
		pn := strings.Split(p[2], ",")
		var lev2 []int
		for i := 0; i < len(pn); i++ {
			b, _ := strconv.Atoi(pn[i])
			if b > 0 {
				lev2 = append(lev2, b)
			}
		}
		// третий уровень - создать образ пускового стимула фразы типа TriggerStimulsID
		// засунуть фразу в дерево слов и дерево фраз
		prase := p[5]
		wordSensor.VerbalDetection(prase, 1, 0, 0)
		PhraseID := wordSensor.CurrentPhrasesIDarr
		tID, vt := CreateNewlastTriggerStimulsID(0, nil, PhraseID, 0, 0)
		if vt != nil{}
		SaveTriggerStimulsArr()
		lev3 := tID
		pn = strings.Split(p[4], ",")
		var ActionIDarr []int
		for i := 0; i < len(pn); i++ {
			b, _ := strconv.Atoi(pn[i])
			if b > 0 {
				ActionIDarr = append(ActionIDarr, b)
			}
		}
		CreateNewConditionReflex(0, lev1, lev2, lev3, ActionIDarr, 0)
	}
  // lib.WriteNewString(lib.GetMainPathExeFile()+"/memory_reflex/condition_reflexes.txt", out)
	SaveConditionReflex()
	return "OK"
}