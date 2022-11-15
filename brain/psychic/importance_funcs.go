/* Функции Значимости объектов восприятия

Значимость всегда определяется в контексте всех предшествующих условий,
т.е. специфична для активностей деревьев автоматизмов и понимания.

При каждом вызове consciousness определяется текущий объект наибольшой значимости в воспринимаемом -
в функции определения текущей Цели getMentalPurpose().
*/

package psychic

import word_sensor "BOT/brain/words_sensor"

/////////////////////////////////////////////////
/* Фиксация значимости объекта ОБъективного восприятия.
Значимость - величина от 0 до 10, приобретаемая объектов внимания в данной ситуации
- берется из результата пробных действий и связывается ос всеми компонентами воспринимаемого в этих условиях
Тот объект, применение которого привело к данной значимости.
Вызывается в момент ответных действий оператора в период ожидания.
Используется curActiveActions - как ответ типа ActionsImage в период ожидания и получение эффекта effect
*/
func setImportance(effect int){
	// целостный образ действий ID ActionsImage
	createNewlastImportanceID(0, detectedActiveLastNodID, detectedActiveLastUnderstandingNodID, 1, curActiveActions.ID,effect, true)

	// ID несловестного действия ActionsImage.ActID[n]
	if len(curActiveActions.ActID) > 0{
		for i := 0; i < len(curActiveActions.ActID); i++ {
			createNewlastImportanceID(0, detectedActiveLastNodID, detectedActiveLastUnderstandingNodID, 3, curActiveActions.ActID[i],effect, true)
		}
	}

	// ID Verbal - при активации дерева автоматизмов
	if curActiveVerbalID >0 {
		createNewlastImportanceID(0, detectedActiveLastNodID, detectedActiveLastUnderstandingNodID, 4, curActiveVerbalID, effect, true)
	}

	// ID отдельной фразы Verbal.PhraseID[n]
	if len(curActiveActions.PhraseID) >0{
		for i := 0; i < len(curActiveActions.PhraseID); i++ {
			createNewlastImportanceID(0, detectedActiveLastNodID, detectedActiveLastUnderstandingNodID, 5, curActiveActions.PhraseID[i],effect, true)
			// по каждому слову фразы
			wRr:=word_sensor.WordsArrFromPhraseID[curActiveActions.PhraseID[i]]
			for j := 0; j < len(wRr); j++ {
				createNewlastImportanceID(0, detectedActiveLastNodID, detectedActiveLastUnderstandingNodID, 6, curActiveActions.PhraseID[j],effect, true)
			}
		}
	}

	// ID тон сообщения с Пульта  Verbal.ToneID
	if curActiveActions.MoodID>0{
		createNewlastImportanceID(0, detectedActiveLastNodID, detectedActiveLastUnderstandingNodID, 7, curActiveActions.MoodID, effect, true)
	}

	// ID настроение оператора  Verbal.MoodID
	if curActiveActions.ToneID>0{
		createNewlastImportanceID(0, detectedActiveLastNodID, detectedActiveLastUnderstandingNodID, 8, curActiveActions.ToneID, effect, true)
	}
}
////////////////////////////////////////////////

/* Фиксация значимости объекта СУБъективного восприятия.

заполняется из afterWaitingPeriod(
*/
func setImportanceMental(effect int){

}
////////////////////////////////////////////////////



/////////////////////////////////////////////////////
/* определить текущие объекты восприятия и выделить один из них - самые важные НЕГАТИВНЫЕ
		по всем категориям importanceType
При каждом ОБъективном вызове consciousness определяется текущий объект наибольшой значимости в воспринимаемом -
в функции определения текущей Цели getMentalPurpose()
*/
func getGreatestImportance(){
	curImportanceJbjectArr = nil
	var id,s,c = 0,0,0
	// целостный образ действий
	id,s,c=getObjectsImportanceValue(1,curActiveActions.ID, detectedActiveLastNodID, detectedActiveLastUnderstandingNodID)
	if s/c <0{
		curImportanceJbjectArr=append(curImportanceJbjectArr,id)
	}

	// ID несловестного действия ActionsImage.ActID[n]
	if len(curActiveActions.ActID) > 0{
		min:=0
		minID:=0
		for i := 0; i < len(curActiveActions.ActID); i++ {
			id,s,c=getObjectsImportanceValue(3,curActiveActions.ActID[i], detectedActiveLastNodID, detectedActiveLastUnderstandingNodID)
			if s/c <0{
				if min > s/c{
					min = s/c
					id=minID
				}
			}
		}
		if minID>0{
			curImportanceJbjectArr=append(curImportanceJbjectArr,minID)
		}
	}

	// ID Verbal - при активации дерева автоматизмов
	if curActiveVerbalID >0 {
		id,s,c=getObjectsImportanceValue(4,curActiveVerbalID, detectedActiveLastNodID, detectedActiveLastUnderstandingNodID)
		if s/c <0{
			curImportanceJbjectArr=append(curImportanceJbjectArr,id)
		}
	}

	//ID отдельной фразы Verbal.PhraseID[n]
	if len(curActiveActions.PhraseID) >0 {
		min:=0
		minID:=0
		for i := 0; i < len(curActiveActions.PhraseID); i++ {
			id,s,c=getObjectsImportanceValue(5,curActiveActions.PhraseID[i], detectedActiveLastNodID, detectedActiveLastUnderstandingNodID)
			if s/c <0{
				if min > s/c{
					min = s/c
					id=minID
				}
			}

			//ID отдельного слова  из Verbal.PhraseID[n]
			min2:=0
			minID2:=0
			wRr:=word_sensor.WordsArrFromPhraseID[curActiveActions.PhraseID[i]]
			for j := 0; j < len(wRr); j++ {
				id,s,c=getObjectsImportanceValue(6,wRr[j], detectedActiveLastNodID, detectedActiveLastUnderstandingNodID)
				if s/c <0{
					if min2 > s/c{
						min2 = s/c
						id=minID2
					}
				}
			}
			if minID2>0{
				curImportanceJbjectArr=append(curImportanceJbjectArr,minID2)
			}

		}
		if minID>0{
			curImportanceJbjectArr=append(curImportanceJbjectArr,minID)
		}
	}

	// ID тон сообщения с Пульта  Verbal.ToneID
	if curActiveActions.MoodID>0{
		id,s,c=getObjectsImportanceValue(7,curActiveActions.MoodID, detectedActiveLastNodID, detectedActiveLastUnderstandingNodID)
		if s/c <0{
			curImportanceJbjectArr=append(curImportanceJbjectArr,id)
		}
	}

	// ID настроение оператора  Verbal.MoodID
	if curActiveActions.ToneID>0{
		id,s,c=getObjectsImportanceValue(8,curActiveActions.ToneID, detectedActiveLastNodID, detectedActiveLastUnderstandingNodID)
		if s/c <0{
			curImportanceJbjectArr=append(curImportanceJbjectArr,id)
		}
	}
}
/////////////////////////////////////////////////




/* выделить наиболее нзачимое в мыслях в текущем цикле
При каждом СУБъективном вызове consciousness
 */
var curImportanceJbjectArr []int // сюда складываются текущие цели внимания к наиболее важному
func getGreatestImportanceMental(){


}
/////////////////////////////////////////////////









