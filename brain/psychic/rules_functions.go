/* функции Правил

*/


package psychic

import (
	"BOT/lib"
	"strconv"
)

///////////////////////////////////////////////
/* Создать групповое правило (более одного в цепочке EpisodeMemory.TAid)
из последнего участка эпизодической памяти объектиынх (EpisodeMemory.Type==0) элеметов.
*/
func GetRulesFromEpisodeMemory(kind int){
	if EpisodeMemoryLastIDFrameID==0{
		return
	}
	var beginID=0
	for i := EpisodeMemoryLastIDFrameID; i >=0; i-- {
		em:=EpisodeMemoryObjects[i]
		if em == nil || em.Type != kind || (LifeTime - em.LifeTime) >EpisodeMemoryPause{
			break // закончить выборку
		}
		beginID++
	}
	if beginID > 1 {
		var taID []int
		if (EpisodeMemoryLastIDFrameID - beginID) > 1 { // только групповые правила, более 1
			for i := EpisodeMemoryLastIDFrameID - beginID; i < EpisodeMemoryLastIDFrameID; i++ {
				em := EpisodeMemoryObjects[i]
				taID = append(taID, em.RulesID)
			}
			// создать новое, групповое правило
			rulesID,_:=createNewlastrulesID(0, taID)
			if rulesID>0 {
				lib.WritePultConsol("<span style='color:green'>Записано групповое <b>ПРАВИЛО № " + strconv.Itoa(rulesID)+"</b></span>")
			}
		}
	}
}
//////////////////////////////////////////////////




//отслеживать Правила из Пульта в http://go/pages/rulles.php
func getCur10lastRules()string{
	rCount:=lastrulesID
	if rCount >10{
		rCount=10
	}
	var out=""
	for i := 0; i < rCount; i++ {
		r:=rulesArr[lastrulesID-i]
		out+="ID="+strconv.Itoa(r.ID)+":"
		for n := 0; n < len(r.TAid); n++ {
			taa:=TriggerAndActionArr[r.TAid[n]]
			if taa == nil{
				continue
			}
			if n>0{
				out+="<span style='padding:40px;'></span>"
			}else{
				out+="<span style='padding:10px;'></span>"
			}
			if taa.Trigger >0 {
				out += "<b>Стимул:</b> <span style='background-color:#FFECEB;'>" + GetActionsString(taa.Trigger) + "</span> "
			}
			if taa.Trigger <0 {
				out += "<b>Стимул:</b> <span style='background-color:#FFECEB;'>" + GetBaseStateImageString(taa.Trigger) + "</span> "
			}
			out+="=> <b>Ответ:</b> <span style='background-color:#E8E8FF;'>"+GetActionsString(taa.Action)+"</span> "
			out+="<b>Эффект: "+strconv.Itoa(taa.Effect)+"</b>"
			out+="<br>"
		}
		out+="<hr>"
	}
	return out
}
///////////////////////////////////////////


///////////////////////////////////////////
/*  выбрать наилучшее Правило rulesID по действию с Пульта или измееннию состояния
Текущая ситуация - массив самых последних кадров эпизодической памяти и
активный пусковой стимул currentTriggerID типов curActiveActions или curBaseStateImage.
*/
func getSuitableRules(activation_type int)(int){
	var rID=0
/* Нет смысла разделять WasOperatorActiveted и WasConditionsActiveted т.к. currentTriggerID уже все дает :)
	// попытка тупо выбрать лучшее правило из rulesArr
	if WasOperatorActiveted { // оператор отреагировал, искать по curActiveActions
		// попытка найти в эпиз.памяти подходящий фрагмент
		rID=getRulesFromEpisodicsSlice(activation_type)
if rID==0 {
	rID=getRulesArrFromTrigger(currentTriggerID)
}
	}else{
		if WasConditionsActiveted {// изменение условий, искать по curBaseStateImage
			rID=getRulesArrFromTrigger(currentTriggerID)
		}
	}
 */
// попытка срочно найти действие, в опасной ситуации
if CurrentInformationEnvironment.veryActualSituation || CurrentInformationEnvironment.danger{
	rID = getRulesArrFromTrigger(currentTriggerID)
}

	if rID==0 {
/* попытка более обстоятельно найти в эпиз.памяти подходящий фрагмент
	Чем больше limit тем маловероятнее найти совпадения сочетания Правил в ранней эпизодюпамяти,
   так что можно вызывать getRulesFromEpisodicsSlice постепенно уменьшая limit
Чем больше 	limit тем точнее результат обобщения, но меньше вероятность нахождения данного сочетания Правил
 */
	maxSteps:=1000
	limit:=5
	rID=getRulesFromEpisodicsSlice(activation_type,limit,maxSteps)
	if rID==0{
		limit=4
		maxSteps=500
		rID=getRulesFromEpisodicsSlice(activation_type,limit,maxSteps)
	}
	if rID==0{
		limit=3
		maxSteps=400
		rID=getRulesFromEpisodicsSlice(activation_type,limit,maxSteps)
	}
	if rID==0{
		limit=2
		maxSteps=300
		rID=getRulesFromEpisodicsSlice(activation_type,limit,maxSteps)
	}
	}

	return rID
}
/////////////////////////////////////////////////
/* нйти последнее Правило по цепочке последних limit кадров эпиз.памяти
Чем больше limit тем маловероятнее найти совпадения,
так что можно вызывать getRulesFromEpisodicsSlice постепенно уменьшая limit
 */
func getRulesFromEpisodicsSlice(activation_type int,limit int,maxSteps int)(int){
	if EpisodeMemoryLastIDFrameID==0{
		return 0
	}
	var kind=0 // объективнй тип эпизод.памяти
	if activation_type==2{
		kind=1
	}

	var beginID=0
	var preLifeTime=0
	for i := EpisodeMemoryLastIDFrameID; i >=0; i-- {
		em:=EpisodeMemoryObjects[i]
// если самый последний эпизод уже является em.Type == kind
		if i==EpisodeMemoryLastIDFrameID && em.Type == kind{
			continue
		}
		if preLifeTime==0{
			preLifeTime=em.LifeTime
		}
		if em == nil || em.Type != kind || (em.LifeTime - preLifeTime) >EpisodeMemoryPause{
			break // закончить выборку
		}
		beginID++
	}
	if beginID == 0 {// скорее всего вышло время ожидания EpisodeMemoryPause
		return 0
	}
	var rID []int
	// перебор последнего фрагмента кадров эпиз.памяти
	for i := EpisodeMemoryLastIDFrameID - beginID; i <= EpisodeMemoryLastIDFrameID; i++ {
				em := EpisodeMemoryObjects[i]
		ta:=TriggerAndActionArr[em.RulesID]
		if ta !=nil {
			//r.Trigger по знаку всегда совпадает с currentTriggerID
			if ta.Trigger == currentTriggerID{// есть такой пусковой
				if ta.Effect >0{// с хорошим эффектом - раз недавно получилось хорошо, то и повторить
					return em.RulesID
				}
			}
		}
				rID = append(rID, em.RulesID)
	}
// найти такую последовательность в предыдущей эпизод.памяти, но не далее 1000 фрагментов
	// limit=2// для тестирования - оставить 2 последних
if len(rID)>limit{// limit последних
	rID=rID[len(rID)-limit:]
}

	lenFrag:=len(rID)
	steps:=0
	lenEp:=len(EpisodeMemoryObjects)
	var startF = lenEp - 2*lenFrag // отмотать на 2 длины, чтобы не проверять в rID саму себя
	if startF > lenEp{//  а нет еще достаточной длины еп.памяти
		return 0
	}
		// идем назад по кускам lenFrag
		for i := startF; i >= 0; i -- { // =lenFrag - пролетает мимо
			if steps>maxSteps{
				return 0
			}
			var isConc=true
			var lastEM *EpisodeMemory
			for j := 0; j < lenFrag; j++ {
				lastEM=EpisodeMemoryObjects[i+j]
				if lastEM.RulesID != rID[j] {
					isConc=false
					break
				}
			}
			if isConc{// уж ты, нашли такой же фрагмент! но в нем нет пускового curActiveActions (раньше уже смотрели)
				// выдать конечное праило, если оно с хорошим эффектом
				rArr:=rulesArr[lastEM.RulesID]
				lastTa:=rArr.TAid[len(rArr.TAid)-1:]
				ta:=TriggerAndActionArr[lastTa[0]]
				if ta !=nil {
					//r.Trigger по знаку всегда совпадает с currentTriggerID
					//if ta.Trigger == currentTriggerID{// есть такой пусковой
						if ta.Effect >0{// с хорошим эффектом
							return lastEM.RulesID
						}//else - продолжает искать хороший конец далее назад
					//}
				}
			}
			steps++
		}

	return 0
}
///////////////////////////////////////////////////


//  выбрать самое лучшее правило из rulesArr по пусковому стимулу
func getRulesArrFromTrigger(trigID int)(int){
	for k, v := range rulesArr {
		for i := 0; i < len(v.TAid); i++ {
			if trigID!=v.TAid[i]{
				continue
			}
			ta:=TriggerAndActionArr[k]
			if ta !=nil {
				//r.Trigger по знаку всегда совпадает с currentTriggerID
				if ta.Trigger == currentTriggerID{// есть такой пусковой
					if ta.Effect >0{// первый попавшийся с хорошим эффектом
						return k
					}
				}
			}
		}
	}

	return 0
}
///////////////////////////////////////////////
