/* распознаватель условного рефлекса

1. С помощью findConditionsReflesFromPrase( из всех у.рефлексов с данным ID образа пускового стимула (imgId3)
выбирается тот, что подходит к данным условиям 1 и 2 уровня.
2. Если на публьте была вбита фраза, для которой нет imgId3, то фраза очищается от неалфавитных символов
и снова пробуется найти подходящий imgId3
3. Если все еще нет подходящего imgId3 то фраза комбинируется:
перебираются все сочетания слов до максимального числа, без перемешивания, не менее чем по 2 слова
4. Если все еще нет подходящего imgId3 то пробуются все слова фразы, не менее 5 символов.
Это позволяет найти у.рефлекс среди длинной фразы, например,
во фразе "я боюсь тебя" будет найден рефлекс на слово "боюсь".
*/

package reflexes

import (
	"BOT/brain/gomeostas"
	wordSensor "BOT/brain/words_sensor"
	"BOT/lib"
	"BOT/tools"
	"strings"
)

// поиск условного рефлекса
func conditionRexlexFound(cond []int) bool {
	if cond == nil || len(cond) == 0 { return false	}
	reflex := getRightConditionReflexesFrom3(cond[0])
	if reflex == nil {
		// попробовать найти схожие по образу рефлексы чтобы не так жестко привязываться к точности фразы
		reflex = findConditionsReflesFromImgID(cond, cond[0])
	}
	if reflex == nil { return false }

	delay := 3600 * 24 * 50 // 50 дней в секундах
	// коэффициент влияния разницы по времени между текущей и последней активациями
	life := lib.Round(float64((LifeTime - reflex.lastActivation) / delay))
	if life < 0 {	life = 0 }
	// определить время просрочки рефлекса при неиспользовании > 50 дней с учетом времени последней активации
	delay = delay * (life + 1)
	// если рефлекс активировался менее чем через 50 дней от последней активации - сдвинуть проверку
	if reflex.lastActivation > LifeTime - delay {
		// обновить актуальность рефлекса
		reflex.lastActivation = LifeTime * 2
		conditionReflexesIdArr = append(conditionReflexesIdArr, reflex.ID)
		return true
	} else { // рефлекс просрочен и должен быть дезактивирован
		reflex.lastActivation = 0
	}

	return false
}

/* попробовать найти другие образы типа TriggerStimuls,
	упрощая фразу из массива фраз TriggerStimulsArr[cond[0]].PhraseID []int
перебором массива var TriggerStimulsArr = make(map[int]*TriggerStimuls)
*/
func findConditionsReflesFromImgID(cond []int,ImgID int) *ConditionReflex {
	var reflex *ConditionReflex
	var wArr []string
	var prase = ""
	var words = ""

	// выделить текущую фразу
	img := TriggerStimulsArr[ImgID]
	if img == nil || img.PhraseID == nil { return nil	}
	for i := 0; i < len(img.PhraseID); i++ {
		prase += wordSensor.GetPhraseStringsFromPhraseID(img.PhraseID[i])
	}
	prase = strings.Trim(prase,"")
	if len(prase) > 0 {
		// если есть не буквенные символы, то убрать их
		prase = wordSensor.ClinerNotAlphavit(prase)
		// есть ли такой образ?
		reflex = findConditionsReflesFromPrase(cond, prase)
		if reflex == nil { // все еще не нашли подходящий рефлекс
			// если во фразе несколько слов, то попробовать все сочетания слов фразы по порядку (а не перемещивая)
			pArr := strings.Split(prase, " ")
			for i := 0; i < len(pArr); i++ {
				if len(pArr[i]) == 0 { continue	}
				wArr = append(wArr, pArr[i])
			}
			if len(wArr) > 1 {
				max := len(wArr)
				if max > 5 { max = 5 } // не более 5 слов во фразе для подбора условного рефлекса
				limit := len(wArr) - 1 // максимальное число элементов в сочетании
				if limit > 3 { limit = 3 }
				// найти все сочетания ряда чисел от 0 до максимального подряд, без перемешивания, не менее чем по 2 слова
				combArr := tools.GetAllCombinationsOfSeriesNumbers(len(wArr), limit)
				// перебор сочетаний слов combArr
				for i := 0; i < len(combArr); i++ {
					for n := 0; n < len(combArr[i]); n++ {
						if n > 0 { words += " " }
						words += wArr[combArr[i][n]]
					}
					// есть ли такой образ?
					reflex = findConditionsReflesFromPrase(cond, words)
					if reflex != nil {
						return reflex
					}
				}
			}
			// напоследок посмотреть по одному длинному слову, > 5 символов (у русских 5*2)
			for i := 0; i < len(wArr); i++ {
				if len(wArr[i]) < 10 { continue	}
				// есть ли такой образ?
				reflex = findConditionsReflesFromPrase(cond, wArr[i])
				if reflex != nil {
					return reflex
				}
				// м.б. еще и первая-последняя буквы - точно, остальные впремешку?
				// TODO
			}
		}
	}
	return reflex
}

// поиск образа у-рефлекса
func findConditionsReflesFromPrase(cond []int, prase string) *ConditionReflex {
	if len(prase) == 0 { return nil	}
	// есть ли такая фраза в Дереве фраз?
	id := wordSensor.GetExistsPraseID(prase)
	if id > 0 { // id фразы есть, найти ее образ по TriggerStimulsArr
		for k, v := range TriggerStimulsArr {
			if v.PhraseID == nil { continue	}
			if v.PhraseID[0] == id { // есть образ с такой фразой
				reflex := getRightConditionReflexesFrom3(k)
				// есть рефлекс с таким образом
				if reflex != nil { return reflex }
			}
		}
	}
	return nil
}

// выбор наиболее близкого по условиям рефлекса из массива с данным пусковым стимулом
// var ConditionReflexesFrom3=make(map[int][]*ConditionReflex)
func getRightConditionReflexesFrom3(imgId3 int) *ConditionReflex {
	ActiveCurBaseID = gomeostas.CommonBadNormalWell
	bsIDarr := gomeostas.GetCurContextActiveIDarr()
	rArr := ConditionReflexesFrom3[imgId3]
	if rArr == nil { return nil	}
	for _, v := range rArr {
		// это - способ прохода дерева без рекурсии, т.к. строго заданы уровни веток:
		if v.lev1 == ActiveCurBaseID && lib.EqualArrs(v.lev2,bsIDarr) {
			return v
		}
	}
	return nil
}