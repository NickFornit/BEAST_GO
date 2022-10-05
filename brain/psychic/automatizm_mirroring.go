/* Формирование зеркальных автоматизмов

*/

package psychic

import (
	wordSensor "BOT/brain/words_sensor"
	"BOT/lib"
	"strconv"
	"strings"
)

/////////////////////////////////
// Формирование зеркальных автоматизмов на основе списка ответов
// тестирование - запуск из psychic.go
func FormingMirrorAutomatizmFromList(file string)(string){
	path:=lib.GetMainPathExeFile()
	strArr,_:=lib.ReadLines(path+file)
	//	triggPhrase|baseID|ContID_list|answerPhrase|Ton,Mood|actions1,...
	if len(strArr)<2{
		return "Пустой файл"
	}
/* Эти автоматизмы привязываются к baseID|ContID_list|0|, т.е. к нулевому образу пусковых
и к нулевому тону-настроению 90.
Но TODO: сделать более мягкую активацию автоматизмов дерева:
   ID|ParentNode|BaseID|EmotionID|ActivityID|ToneMoodID|SimbolID|VerbalID
   Если нет автоматизма для данного узла > ActivityID то смотреть для других узлов, начиная с данного уровня.
Т.е. если автоматизм привязан к ToneMoodID==90 а активировалась ветка с ToneMoodID==12 где нет автоматизма,
то пусть бы срабатывал привязанный к ToneMoodID==90 !!!
   НО ОСТОРОЖНО (понизить силу?)
 */


//	strArr:=strings.Split(list, "\r\n")
	// первую строку пропускаем из-за #utf8 bom
	for n := 1; n < len(strArr); n++ {
		if len(strArr[n])<10{
			continue
		}
		p := strings.Split(strArr[n], "|")
///////// УСЛОВИЯ ДЕРЕВА
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
		//emotionID,_:=createNewBaseStyle(0,0,lev2)

		// образ отсуствия пусковых
		//triggID,_:=createNewlastActivityID(0,nil)
		// образ отсуствия тона и настроения
		tm:=90
		//фраза
// засунуть фразу в дерево слов и дерево фраз
		prase:=p[0]
		wordSensor.VerbalDetection(prase, 1, 0, 0)
		PhraseID := wordSensor.CurrentPhrasesIDarr

		//первый симвод ответной фразы
		FirstSimbolID:=wordSensor.GetFirstSymbolFromPraseID(PhraseID)
		// создать образ Брока
		CreateVerbalImage(FirstSimbolID, PhraseID, 0, 0)

		nodeID := FindConditionsNode(baseID, lev2, nil,tm,PhraseID[0],FirstSimbolID)
		if nodeID>0{
			/*
если есть привязанный к узлу автоматизм, то он просто перестанет быть штатным,
т.к. авторитерный (зеркальный) автоматизм важнее
			exists:=ExistsAutomatizmForThisNodeID(nodeID)
			if exists {
				continue
			}
*/
			//  создать автоматизм и привязать его к nodeID
		var sequence="Snn:" // ответная фраза
			// засунуть фразу в дерево слов и дерево фраз
			wordSensor.VerbalDetection(p[3], 1, 0, 0)
			answerID := wordSensor.CurrentPhrasesIDarr
			sequence+=strconv.Itoa(answerID[0])

		sequence+="|Тnn:" // тон и настроение
		tnArr := strings.Split(p[4], ",")
		t,_:=strconv.Atoi(tnArr[0])
		m,_:=strconv.Atoi(tnArr[1])
		tonMoodID:=GetToneMoodID(t,m)
		sequence+=strconv.Itoa(tonMoodID)

		sequence+="|Dnn:" // перечень ответных действий
			aD := strings.Split(p[5], ",")
			for i := 0; i < len(aD); i++ {
				if i > 0 {
					sequence += ","
				}
				sequence += aD[i]
			}
			NoWarningCreateShow=true
			_,autmzm:=CreateAutomatizm(nodeID,sequence)
			NoWarningCreateShow=false
			if autmzm!=nil{
				SetAutomatizmBelief(autmzm, 2)// сделать автоматизм штатным
				// ?? autmzm.GomeoIdSuccesArr какие ID гомео-параметров улучшает это действие
				autmzm.Usefulness=1 // полезность
			}
		}
	}

	SaveAllPsihicMemory()
	return "OK"
}
//////////////////////////////////

/* на основе общего шаблона ответов
Без привязки к конкретным условиям первых трех уровней: все три уровня ставятся == 0
 */
func FormingMirrorAutomatizmFromTempList(file string)(string){
	path:=lib.GetMainPathExeFile()
	strArr,_:=lib.ReadLines(path+file)
	//	triggPhrase|baseID|ContID_list|answerPhrase|Ton,Mood|actions1,...
	if len(strArr)<2{
		return "Пустой файл"
	}
	/* Эти автоматизмы привязываются к baseID|ContID_list|0|, т.е. к нулевому образу пусковых
	   и к нулевому тону-настроению 90.
	   Но TODO: сделать более мягкую активацию автоматизмов дерева:
	      ID|ParentNode|BaseID|EmotionID|ActivityID|ToneMoodID|SimbolID|VerbalID
	      Если нет автоматизма для данного узла > ActivityID то смотреть для других узлов, начиная с данного уровня.
	   Т.е. если автоматизм привязан к ToneMoodID==90 а активировалась ветка с ToneMoodID==12 где нет автоматизма,
	   то пусть бы срабатывал привязанный к ToneMoodID==90 !!!
	      НО ОСТОРОЖНО (понизить силу?)
	*/


	//	strArr:=strings.Split(list, "\r\n")
	// первую строку пропускаем из-за #utf8 bom
	for n := 1; n < len(strArr); n++ {
		if len(strArr[n])<10{
			continue
		}
		p := strings.Split(strArr[n], "|")
		///////// УСЛОВИЯ ДЕРЕВА
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
		//emotionID,_:=createNewBaseStyle(0,0,lev2)

		// образ отсуствия пусковых
		//triggID,_:=createNewlastActivityID(0,nil)
		// образ отсуствия тона и настроения
		tm:=90
		//фраза
		// засунуть фразу в дерево слов и дерево фраз
		prase:=p[0]
		wordSensor.VerbalDetection(prase, 1, 0, 0)
		PhraseID := wordSensor.CurrentPhrasesIDarr

		//первый симвод ответной фразы
		FirstSimbolID:=wordSensor.GetFirstSymbolFromPraseID(PhraseID)
		// создать образ Брока
		CreateVerbalImage(FirstSimbolID, PhraseID, 0, 0)

		nodeID := FindConditionsNode(baseID, lev2, nil,tm,PhraseID[0],FirstSimbolID)
		if nodeID>0{
			/*
			если есть привязанный к узлу автоматизм, то он просто перестанет быть штатным,
			т.к. авторитерный (зеркальный) автоматизм важнее
						exists:=ExistsAutomatizmForThisNodeID(nodeID)
						if exists {
							continue
						}
			*/
			//  создать автоматизм и привязать его к nodeID
			var sequence="Snn:" // ответная фраза
			// засунуть фразу в дерево слов и дерево фраз
			wordSensor.VerbalDetection(p[3], 1, 0, 0)
			answerID := wordSensor.CurrentPhrasesIDarr
			sequence+=strconv.Itoa(answerID[0])

			sequence+="|Тnn:" // тон и настроение
			tnArr := strings.Split(p[4], ",")
			t,_:=strconv.Atoi(tnArr[0])
			m,_:=strconv.Atoi(tnArr[1])
			tonMoodID:=GetToneMoodID(t,m)
			sequence+=strconv.Itoa(tonMoodID)

			sequence+="|Dnn:" // перечень ответных действий
			aD := strings.Split(p[5], ",")
			for i := 0; i < len(aD); i++ {
				if i > 0 {
					sequence += ","
				}
				sequence += aD[i]
			}
			NoWarningCreateShow=true
			_,autmzm:=CreateAutomatizm(nodeID,sequence)
			NoWarningCreateShow=false
			if autmzm!=nil{
				SetAutomatizmBelief(autmzm, 2)// сделать автоматизм штатным
				// ?? autmzm.GomeoIdSuccesArr какие ID гомео-параметров улучшает это действие
				autmzm.Usefulness=1 // полезность
			}
		}
	}

	SaveAllPsihicMemory()
	return "OK"
}
//////////////////////////////////
