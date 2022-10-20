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
			//var sequence = "Snn:" // ответная фраза
			// засунуть фразу в дерево слов и дерево фраз
			wordSensor.VerbalDetection(p[3], 1, 0, 0)
			answerID := wordSensor.CurrentPhrasesIDarr

			//sequence += "|Тnn:" // тон и настроение
			tnArr := strings.Split(p[4], ",")
			t, _ := strconv.Atoi(tnArr[0])
			m, _ := strconv.Atoi(tnArr[1])

			var aArr []int
			aD := strings.Split(p[5], ",")
			for i := 0; i < len(aD); i++ {
				a,_:=strconv.Atoi(aD[i])
				aArr=append(aArr,a)
			}

			NoWarningCreateShow = true

			ActionsImageID,_:=CreateNewActionsImageImage(aArr,answerID,t,m)
			_, autmzm := CreateAutomatizm(nodeID, ActionsImageID)
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
		//tm := GetToneMoodID(t, m + 19)

		// засунуть фразу в дерево слов и дерево фраз
		wordSensor.VerbalDetection(triggerPrase, 1, 0, 0)
		triggerPraseID := wordSensor.CurrentPhrasesIDarr

		wordSensor.VerbalDetection(answerPrase, 1, 0, 0)
		answerPraseID := wordSensor.CurrentPhrasesIDarr

		// создать автоматизм и привязать его к объекту
		NoWarningCreateShow=true
		// для фразы triggerPraseID создаем привязанный к ней автоматизм
		ActionsImageID,_:=CreateNewActionsImageImage(nil,answerPraseID,t,m)
		_, autmzm := CreateAutomatizm(2000000 + triggerPraseID[0], ActionsImageID)
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
/* вытащить действия исходного автоматизма чтобы найти или сделать узел дерева с таким пускателем
и существубшими BaseID и EmotionID
 */
	curNode:=AutomatizmTreeFromID[detectedActiveLastNodID]
	targetNodeID:=findTreeNodeFromAutomatizmActionsImage(curNode.BaseID, curNode.EmotionID, sourceAtmzm.ActionsImageID)
	if targetNodeID==0{
		return
	}
//	SaveAutomatizmTree()
// найти узел, который может реагировать на данные действия и если нет - создать его чтобы привязать зеркальный автоматизм

	// создать автоматизм и привязать его к объекту
	// NoWarningCreateShow=true
	ActionsImageID,_:=CreateNewActionsImageImage(curActiveActions.ActID,curActiveActions.PhraseID,curActiveActions.ToneID,curActiveActions.MoodID)
	_, autmzm := CreateAutomatizm(targetNodeID, ActionsImageID)
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
func findTreeNodeFromAutomatizmActionsImage(baseID int, EmotionID int, ActionsImageID int)(int) {

	lev2:=EmotionFromIdArr[EmotionID].BaseIDarr

	ai:=ActionsImageArr[ActionsImageID]

// первый символ ответной фразы
	simbolID := wordSensor.GetFirstSymbolFromPraseID(ai.PhraseID)

	tm := GetToneMoodID(ai.ToneID, ai.MoodID + 19)

	nodeID := FindConditionsNode(baseID, lev2, ai.ActID, ai.PhraseID[0], simbolID, tm)
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

	ActionsImageID,_:=CreateNewActionsImageImage(curActiveActions.ActID,curActiveActions.PhraseID,curActiveActions.ToneID,curActiveActions.MoodID)
	// NoWarningCreateShow=true
	// для фразы triggerPraseID создаем привязанный к ней автоматизм
	_, autmzm := CreateAutomatizm(detectedActiveLastNodID, ActionsImageID)
	// NoWarningCreateShow=false
	if autmzm != nil {
		//autmzm.BranchID += linkID // не привязывать к узлу
		SetAutomatizmBelief(autmzm, 2) // сделать автоматизм штатным, т.к. действия авторитарно верные (копируем действия оператора)
		autmzm.Usefulness = 1 // авторитарная полезность
	}
	// и тут же запустить реакцию с ожиданием ответа
	setAutomatizmRunning(autmzm, purposeGenetic)
}
//////////////////////////////////
