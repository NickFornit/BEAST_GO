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
				createNewlastImportanceID(0, detectedActiveLastNodID, detectedActiveLastUnderstandingNodID, 6, curActiveActions.PhraseID[i],effect, true)
			}
		}
	}

	// ID тон сообщения с Пульта  Verbal.ToneID
	if curActiveActions.ToneID>0{
		createNewlastImportanceID(0, detectedActiveLastNodID, detectedActiveLastUnderstandingNodID, 7, curActiveActions.MoodID, effect, true)
	}

	// ID настроение оператора  Verbal.MoodID
	if curActiveActions.MoodID>0{
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
type extremImportance struct {
	objID int//  объект значимости
	kind int // тип объекта
	extremVal int// экстремальная значимость
}
var curImportanceObjectArr []extremImportance //- здесь сохраняются текущие цели внимания к наиболее важному


func getGreatestImportance(curActions *ActionsImage)[]extremImportance{
	var importanceObjectArr []extremImportance
	if curActions == nil{
		return nil
	}
	var id,v = 0,0
	// целостный образ действий
	id,v=getObjectsImportanceValue(1,curActions.ID, detectedActiveLastNodID, detectedActiveLastUnderstandingNodID)
	if v <0{
		importanceObjectArr =append(importanceObjectArr,extremImportance{id,1,v})
	}

	// ID несловестного действия ActionsImage.ActID[n]
	if len(curActions.ActID) > 0{
		min:=0
		minID:=0
		for i := 0; i < len(curActions.ActID); i++ {
			id,v=getObjectsImportanceValue(3,curActions.ActID[i], detectedActiveLastNodID, detectedActiveLastUnderstandingNodID)
			if v <0{
				if min > v{
					min = v
					minID=id
				}
			}
		}
		if minID>0{
			importanceObjectArr =append(importanceObjectArr,extremImportance{minID,3,min})
		}
	}

	// ID Verbal - при активации дерева автоматизмов
	if curActiveVerbalID >0 {
		id,v=getObjectsImportanceValue(4,curActiveVerbalID, detectedActiveLastNodID, detectedActiveLastUnderstandingNodID)
		if v <0{
			importanceObjectArr =append(importanceObjectArr,extremImportance{id,4,v})
		}
	}

	//ID отдельной фразы Verbal.PhraseID[n]
	if len(curActions.PhraseID) >0 {
		min:=0
		minID:=0
		for i := 0; i < len(curActions.PhraseID); i++ {
			id,v=getObjectsImportanceValue(5,curActions.PhraseID[i], detectedActiveLastNodID, detectedActiveLastUnderstandingNodID)
			if v <0{
				if min > v{
					min = v
					minID=id
				}
			}

			//ID отдельного слова  из Verbal.PhraseID[n]
			min2:=0
			minID2:=0
			wRr:=word_sensor.WordsArrFromPhraseID[curActions.PhraseID[i]]
			for j := 0; j < len(wRr); j++ {
				id,v=getObjectsImportanceValue(6,wRr[j], detectedActiveLastNodID, detectedActiveLastUnderstandingNodID)
				if v <0{
					if min2 > v{
						min2 = v
						minID2=id
					}
				}
			}
			if minID2>0{
				importanceObjectArr =append(importanceObjectArr,extremImportance{minID2,6,min2})
			}

		}
		if minID>0{
			importanceObjectArr =append(importanceObjectArr,extremImportance{minID,5,min})
		}
	}

	// ID тон сообщения с Пульта  Verbal.ToneID
	if curActions.ToneID>0{
		id,v=getObjectsImportanceValue(7,curActions.ToneID, detectedActiveLastNodID, detectedActiveLastUnderstandingNodID)
		if v <0{
			importanceObjectArr =append(importanceObjectArr,extremImportance{id,7,v})
		}
	}

	// ID настроение оператора  Verbal.MoodID
	if curActions.MoodID>0{
		id,v=getObjectsImportanceValue(8,curActions.MoodID, detectedActiveLastNodID, detectedActiveLastUnderstandingNodID)
		if v <0{
			importanceObjectArr =append(importanceObjectArr,extremImportance{id,8,v})
		}
	}

	return importanceObjectArr
}
/////////////////////////////////////////////////

/* выбрать один, самый актуальный объект из curImportanceJbjectArr []extremImportance
Возвращает индекс массива curImportanceObjectArr или -1
 */
func getTopAttentionObject(eij []extremImportance)int{
	if eij==nil{
		return -1
	}
	// тупо выбираем самый негативный
	min:=0
	index:=-1
	for i := 0; i < len(eij); i++ {
		if min > eij[i].extremVal{
			min = eij[i].extremVal
			index=i
		}
	}
	return index
}
/////////////////////////////////////////////////






////////////////////////////////////////////////////
/* выделить наиболее значимое в мыслях в текущем цикле - привлечение внимания к собственным мыслям.
При каждом СУБъективном вызове consciousness
 */
type extremImportanceMental struct {
	objID int//  объект значимости
	kind int // тип объекта
	extremVal int// экстремальная значимость
}
var curImportanceObjectMentalArr []extremImportanceMental // сюда складываются текущие цели внимания к наиболее важному
func getGreatestImportanceMental(){


}
/////////////////////////////////////////////////









