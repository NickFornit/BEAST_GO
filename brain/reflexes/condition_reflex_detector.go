/* распознаватель условного рефлекса

*/

package reflexes

import (
	wordSensor "BOT/brain/words_sensor"
	"BOT/tools"
	"strings"
)

//////////////////////////////////////////

func conditionRexlexFound(cond []int)(bool){
	if cond==nil || len(cond)==0{
		return false
	}
// var ConditionReflexesFrom3=make(map[int]*ConditionReflex)
	reflex:=ConditionReflexesFrom3[cond[0]]
	if reflex==nil {
// попробовать найти схожие по образу рефлексы чтобы не так жестко привязываться к точности фразы
		reflex=findConditionsReflesFromImgID(cond[0])
	}
	if reflex==nil {
		return false
	}
	///////////////////////

		// время жизни рефлекса минус 10 дней
		//life:=LifeTime - reflex.activationTime -10
		life:= reflex.lastActivation - reflex.activationTime -3600*24*10
		if life<0{life=0}
		// коэффициент влияния времени жизни: каждые 10 дней укрепляют рефлекс в 2 раза
		k:=1+ (2*life)/(3600*24*10)
// определить время просрочки рефлекса при неиспользовании > 10 дней с учетом времни жизни
		delay:=(3600*24*10)*k
//Рефлкс НЕ активен, если его lastActivation меньше, чем activationTime-delay
		if reflex.lastActivation > reflex.activationTime-delay {
			// обновить актуальность рефлекса
			reflex.lastActivation=reflex.activationTime
			conditionReflexesIdArr = append(conditionReflexesIdArr, reflex.ID)
			return true
		}else{// рефлекс просрочен и должен быть дезактивирован
			reflex.lastActivation=0
		}

	return false
}
//////////////////////////////////////////////



//////////////////////////////////////////////////
/* попробовать найти другие образы типа TriggerStimuls,
	упрощая фразу из массива фраз TriggerStimulsArr[cond[0]].PhraseID []int
перебором массива var TriggerStimulsArr = make(map[int]*TriggerStimuls)
*/
func findConditionsReflesFromImgID(ImgID int)(*ConditionReflex){
	var reflex *ConditionReflex
// выделить текущую фразу
img:=TriggerStimulsArr[ImgID]
if img!=nil{
prase:=wordSensor.GetPhraseStringsFromPhraseID(img.PhraseID[0])
if len(prase)>0{
//если есть не буквенные символы, то убрать их
prase=wordSensor.ClinerNotAlphavit(prase)
// есть ли такой образ?
reflex=findConditionsReflesFromPrase(prase)
if reflex==nil {// все еще не нашли подходящий рефлекс
// если во фразе несколько слов, то попробовать все сочетания слов фразы по порядку (а не перемещивая)
wArr:=strings.Split(prase, " ")
if len(wArr)>1 {
	max:=len(wArr)
	if max > 5 {max=5} // не более 5 слов во фразе для подбора условного рефлекса
	limit:=len(wArr)-1 //максимальное число элементов в сочетании
		if limit>3{limit=3}
// найти все сочетания ряда чисел от 0 до максимального подряд, без перемешивания, не менее чем по 2 слова
	combArr := tools.GetAllCombinationsNumbers(len(wArr),limit)
	// передор сочетаний слов combArr
	for i := 0; i < len(combArr); i++ {
		var words=""
		for n := 0; n < len(combArr[i]); n++ {
			words +=wArr[combArr[i][n]]
		}
		// есть ли такой образ?
		reflex=findConditionsReflesFromPrase(words)
		if reflex!=nil {
			return reflex
		}
	}
}
// напоследок посмотреть по одному длинному слову, > 5 символов (у русских 5*2)
	for i := 0; i < len(wArr); i++ {
		if len(wArr[i]) < 10{
			continue
		}
		// есть ли такой образ?
		reflex=findConditionsReflesFromPrase(wArr[i])
		if reflex!=nil {
			return reflex
		}
	}
}
}
}
return reflex
}
////////////////////////////////////

func findConditionsReflesFromPrase(prase string)(*ConditionReflex){
// есть ли такая фраза в Дереве фраз?
id:=wordSensor.GetExistsPraseID(prase)
if id>0{// id фразы есть, найти ее образ по TriggerStimulsArr
for k, v := range TriggerStimulsArr {
if v.PhraseID[0]==id {// есть образ с такой фразой
reflex:=ConditionReflexesFrom3[k]
if reflex!=nil {// есть рефлекс с таким образом
 return reflex
}
}
}
}
return nil
}