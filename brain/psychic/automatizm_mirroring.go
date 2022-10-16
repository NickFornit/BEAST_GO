/* Формирование зеркальных автоматизмов */

package psychic

import (
	wordSensor "BOT/brain/words_sensor"
	"BOT/lib"
	"strconv"
	"strings"
)

// Формирование зеркальных автоматизмов на основе списка ответов lib/mirror_reflexes_basic_phrases/...
// тестирование - запуск из psychic.go
func FormingMirrorAutomatizmFromList(file string) string {
	path := lib.GetMainPathExeFile()
	strArr, _ := lib.ReadLines(path + file)
	// triggPhrase|baseID|ContID_list|answerPhrase|Ton,Mood|actions1,...
	if len(strArr) < 2 { return "Пустой файл"	}
	/* Эти автоматизмы привязываются к baseID|ContID_list|0|, т.е. к нулевому образу пусковых
	и к нулевому тону-настроению 90.
	Но TODO: сделать более мягкую активацию автоматизмов дерева:
		 ID|ParentNode|BaseID|EmotionID|ActivityID|ToneMoodID|SimbolID|VerbalID
		 Если нет автоматизма для данного узла > ActivityID то смотреть для других узлов, начиная с данного уровня.
	Т.е. если автоматизм привязан к ToneMoodID==90 а активировалась ветка с ToneMoodID==12 где нет автоматизма,
	то пусть бы срабатывал привязанный к ToneMoodID==90 !!!
		 НО ОСТОРОЖНО (понизить силу?)
	 */

	// первую строку пропускаем из-за #utf8 bom
	for n := 1; n < len(strArr); n++ {
		if len(strArr[n]) < 10 { continue	}
		p := strings.Split(strArr[n], "|")
		// УСЛОВИЯ ДЕРЕВА
		// базовое состояние
		baseID,_ := strconv.Atoi(p[1])

		// базовые контексты
		pn := strings.Split(p[2], ",")
		var lev2 []int
		for i := 0; i < len(pn); i++ {
			b, _ := strconv.Atoi(pn[i])
			if b > 0 {
				lev2 = append(lev2, b)
			}
		}

		// образ отсуствия тона и настроения
		tm := 90

		// засунуть фразу в дерево слов и дерево фраз
		prase := p[0]
		wordSensor.VerbalDetection(prase, 1, 0, 0)
		PhraseID := wordSensor.CurrentPhrasesIDarr

		// первый символ ответной фразы
		FirstSimbolID := wordSensor.GetFirstSymbolFromPraseID(PhraseID)
		// создать образ Брока
		CreateVerbalImage(FirstSimbolID, PhraseID, 0, 0)

		nodeID := FindConditionsNode(baseID, lev2, nil, tm, PhraseID[0], FirstSimbolID)
		/* если есть привязанный к узлу автоматизм, то он просто перестанет быть штатным,
		т.к. авторитерный (зеркальный) автоматизм важнее
		exists:=ExistsAutomatizmForThisNodeID(nodeID)
		if exists {
			continue
		}	*/
		if nodeID > 0 {
			// создать автоматизм и привязать его к nodeID
			var sequence = "Snn:" // ответная фраза
			// засунуть фразу в дерево слов и дерево фраз
			wordSensor.VerbalDetection(p[3], 1, 0, 0)
			answerID := wordSensor.CurrentPhrasesIDarr
			//sequence += strconv.Itoa(answerID[0])
			for i := 0; i < len(answerID); i++ {
				if i > 0 { sequence += ","}
				sequence += strconv.Itoa(answerID[i]) // ответная фраза
			}

			sequence += "|Тnn:" // тон и настроение
			tnArr := strings.Split(p[4], ",")
			t, _ := strconv.Atoi(tnArr[0])
			m, _ := strconv.Atoi(tnArr[1])
			tonMoodID := GetToneMoodID(t, m)
			sequence += strconv.Itoa(tonMoodID)

			sequence += "|Dnn:" // перечень ответных действий
			aD := strings.Split(p[5], ",")
			for i := 0; i < len(aD); i++ {
				if i > 0 { sequence += "," }
				sequence += aD[i]
			}
			NoWarningCreateShow = true
			_, autmzm := CreateAutomatizm(nodeID, sequence,0)
			NoWarningCreateShow = false
			if autmzm != nil {
				SetAutomatizmBelief(autmzm, 2) // сделать автоматизм штатным
				// ?? autmzm.GomeoIdSuccesArr какие ID гомео-параметров улучшает это действие
				autmzm.Usefulness = 1 // полезность
			}
		}
	}
	SaveAllPsihicMemory()
	return "OK"
}

/* на основе общего шаблона ответов lib/mirror_basic_phrases_common.txt
Создаются автоматизмы, привязанные к пусковой фразе, а не к узлу дерева,
с BranchID > 2000000.
var AutomatizmIdFromPhraseId=make(map[int] []*Automatizm)
// тестирование - запуск из psychic.go
 */
func FormingMirrorAutomatizmFromTempList(file string) string {
	path := lib.GetMainPathExeFile()
	strArr, _ := lib.ReadLines(path + file)
	// triggPhrase|baseID|ContID_list|answerPhrase|Ton,Mood|actions1,...
	if len(strArr) == 0 { return "Пустой файл" }
	/* Эти автоматизмы привязываются к baseID|ContID_list|0|, т.е. к нулевому образу пусковых
	   и к нулевому тону-настроению 90.
	   Но TODO: сделать более мягкую активацию автоматизмов дерева:
	      ID|ParentNode|BaseID|EmotionID|ActivityID|ToneMoodID|SimbolID|VerbalID
	      Если нет автоматизма для данного узла > ActivityID то смотреть для других узлов, начиная с данного уровня.
	   Т.е. если автоматизм привязан к ToneMoodID==90 а активировалась ветка с ToneMoodID==12 где нет автоматизма,
	   то пусть бы срабатывал привязанный к ToneMoodID==90 !!!
	      НО ОСТОРОЖНО (понизить силу?)
	*/

	for n := 0; n < len(strArr); n++ {
		if len(strArr[n]) < 10 { continue	}
		p := strings.Split(strArr[n], "|")
		// УСЛОВИЯ ДЕРЕВА
		// пусковая фраза
		triggerPrase := p[0]

		// ответ
		answerPrase := p[1]

		// тон, настроение
		pt := strings.Split(p[2], ",")
		t,_ := strconv.Atoi(pt[0])
		m,_ := strconv.Atoi(pt[1])
		tm := GetToneMoodID(t, m + 19)

		// засунуть фразу в дерево слов и дерево фраз
		wordSensor.VerbalDetection(triggerPrase, 1, 0, 0)
		triggerPraseID := wordSensor.CurrentPhrasesIDarr

		wordSensor.VerbalDetection(answerPrase, 1, 0, 0)
		answerPraseID := wordSensor.CurrentPhrasesIDarr

		// создать автоматизм и привязать его к объекту
		var sequence = "Snn:" // ответная фраза
		for i := 0; i < len(answerPraseID); i++ {
			if i > 0 { sequence += ","}
			sequence += strconv.Itoa(answerPraseID[i]) // ответная фраза
		}
		sequence += "|Tnn:" + strconv.Itoa(tm) // тон и настроение
		sequence += "|Dnn:" // перечень ответных действий
		aD := strings.Split(p[3], ",")
		for i := 0; i < len(aD); i++ {
			if i > 0 { sequence += "," }
			sequence += aD[i]
		}

		NoWarningCreateShow=true
		// для фразы triggerPraseID создаем привязанный к ней автоматизм
		_, autmzm := CreateAutomatizm(2000000 + triggerPraseID[0], sequence,0)
		NoWarningCreateShow = false
		if autmzm != nil {
			SetAutomatizmBelief(autmzm, 2) // сделать автоматизм штатным
			// ?? autmzm.GomeoIdSuccesArr какие ID гомео-параметров улучшает это действие
			autmzm.Usefulness = 1 // полезность
		}
	}
	SaveAllPsihicMemory()
	return "OK"
}

/* создание зеркального автоматизма, повторяющего действия оператора в данных условиях
в ответ на действия sourceAtmzm - причина ответа оператора
Только что действиями оператора была активирована ветка detectedActiveLastNodID дерева и
есть информация об этих действиях в curActiveActions
Автоматизм прикрепляется к ветке предыдущей активации дерева LastDetectedActiveLastNodID (причине) -
которая становится пусковым стимулом отзеркаливания.
*/
func createNewMirrorAutomatizm(sourceAtmzm *Automatizm) {
	if sourceAtmzm == nil { return }
	var sequence = ""
/* вытащить действия исходного автоматизма чтобы найти или сделать узел дерева с таким пускателем
и существубшими BaseID и EmotionID
 */
	curNode:=AutomatizmTreeFromID[detectedActiveLastNodID]
	targetNodeID:=findTreeNodrFromAutomatizmSequence(curNode.BaseID, curNode.EmotionID, sourceAtmzm.Sequence)
	if targetNodeID==0{
		return
	}
	SaveAutomatizmTree()
// найти узел, который может реагировать на данные действия и если нет - создать его чтобы привязать зеркальный автоматизм

	// создать автоматизм и привязать его к объекту
	if len(curActiveActions.phraseID) > 0 {
		sequence += "Snn:"
		for i := 0; i < len(curActiveActions.phraseID); i++ {
			if i > 0 { sequence += ","}
			sequence += strconv.Itoa(curActiveActions.phraseID[i]) // ответная фраза
		}
		// тон, настроение
		tm := GetToneMoodID(curActiveActions.toneID, curActiveActions.moodID + 19)
		sequence += "|Tnn:" + strconv.Itoa(tm) // тон и настроение
	}

	sequence += "|Dnn:" // перечень ответных действий
	for i := 0; i < len(curActiveActions.actID); i++ {
		if i > 0 { sequence += "," }
		sequence += strconv.Itoa(curActiveActions.actID[i])
	}

	// NoWarningCreateShow=true
	// для фразы triggerPraseID создаем привязанный к ней автоматизм
	//_, autmzm := CreateAutomatizm(LastDetectedActiveLastNodID, sequence,1)
	//_, autmzm := CreateAutomatizm(detectedActiveLastNodID, sequence,1)
	_, autmzm := CreateAutomatizm(targetNodeID, sequence,1)
	//	NoWarningCreateShow=false
	if autmzm != nil {
		SetAutomatizmBelief(autmzm, 2) // сделать автоматизм штатным, т.к. действия авторитарно верные
		autmzm.Usefulness = 1 // авторитарная полезность

		// SaveAutomatizm()
	}
}
/* вытащить действия исходного автоматизма чтобы найти или сделать узел дерева с таким пускателем
и существубшими BaseID и EmotionID
*/
func findTreeNodrFromAutomatizmSequence(baseID int, EmotionID int, AtmzmS string)(int){

	lev2:=EmotionFromIdArr[EmotionID].BaseIDarr
	var activityID []int
	var toneMoodID=90
	var simbolID=0
	var verbalID []int
	sArr:=strings.Split(AtmzmS, "|")
	for i := 0; i < len(sArr); i++ {
		if len(sArr[i]) == 0 {
			continue
		}
		pArr := strings.Split(sArr[i], ":")
		switch pArr[0] {
		case "Snn": // есть ли такой у второго
			v, _ := strconv.Atoi(pArr[1])
			verbalID=append(verbalID,v)
			// первый символ ответной фразы
			simbolID = wordSensor.GetFirstSymbolFromPraseID(verbalID)
		case "Dnn":
			aArr := strings.Split(pArr[1], ",")
			for n := 0; n < len(aArr); n++ {
				a, _ := strconv.Atoi(aArr[n])
				activityID=append(activityID,a)
			}

		/* последовательный запуск автоматизмов НЕ ПРОВЕРЯЕМ ЭКЗОТИКУ...
		case "Ann":
		*/
		case "Tnn":
			toneMoodID, _ = strconv.Atoi(pArr[1])
		}
	}

	nodeID := FindConditionsNode(baseID, lev2, activityID, toneMoodID, simbolID, verbalID[0])
	if nodeID > 0 {
return nodeID
	}

	return 0
}
////////////////////////////////////////////////////////




/* в случае отсуствия автоматизма в данных условиях - послать оператору те же стимулы, чтобы посмотреть его реакцию.
Создание автоматизма, повторяющего действия оператора в данных условиях */
func provokatorMirrorAutomatizm(sourceAtmzm *Automatizm, purposeGenetic *PurposeGenetic) {
	if sourceAtmzm == nil || purposeGenetic == nil { return	}
	var linkID = 0
	var sequence = ""

	// создать автоматизм и привязать его к объекту
	sequence += "|Dnn:" // перечень ответных действий
	for i := 0; i < len(curActiveActions.actID); i++ {
		if i > 0 { sequence += "," }
		sequence += strconv.Itoa(curActiveActions.actID[i])
		// образ пусковых
		linkID = 1000000 + curActiveActions.triggID
	}

	if len(curActiveActions.phraseID) > 0 {
		sequence += "Snn:"  // ответная фраза
		for i := 0; i < len(curActiveActions.phraseID); i++ {
			if i > 0 { sequence += ","}
			sequence += strconv.Itoa(curActiveActions.phraseID[i]) // ответная фраза
		}
		// тон, настроение
		tm := GetToneMoodID(curActiveActions.toneID, curActiveActions.moodID + 19)
		sequence += "|Tnn:" + strconv.Itoa(tm) // тон и настроение
		linkID = 2000000 + curActiveActions.phraseID[0] // только по первому слову?
	}

	// NoWarningCreateShow=true
	// для фразы triggerPraseID создаем привязанный к ней автоматизм
	_, autmzm := CreateAutomatizm(detectedActiveLastNodID, sequence,1)
	// NoWarningCreateShow=false
	if autmzm != nil {
		autmzm.BranchID += linkID // не привязывать к узлу
		SetAutomatizmBelief(autmzm, 2) // сделать автоматизм штатным, т.к. действия авторитарно верные
		autmzm.Usefulness = 1 // авторитарная полезность
	}
	// и тут же запустить реакцию с ожиданием ответа
	setAutomatizmRunning(autmzm, purposeGenetic)
}
//////////////////////////////////
