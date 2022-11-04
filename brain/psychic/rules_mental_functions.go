/* функции ментальных Правил


Функция мягкого поиска по деталям образов действий getSuitableMentalRulesСarefully()
который можно использовать в инфо-функциях.
*/


package psychic

import (
	"BOT/lib"
	"strconv"
)

///////////////////////////////////////////////
/* Создать групповое правило (более одного в цепочке EpisodeMemory.TAid)
из последнего участка эпизодической памяти объектиынх (EpisodeMemory.Type==0) элеметов.
limit 5 ограничивает выборку из эпиз.памяти, но она может получться и меньше.
*/
func GetMentalRulesFromEpisodeMemory(kind int){
	rImg:=getLastRulesSequenceFromEpisodeMemory(kind,5)
	if rImg!=nil {
		createNewlastrulesMentalID(0, rImg)//записать (если еще нет такого) групповое правило
	}
}
//////////////////////////////////////////////////

// тестовый вызов из main.go
func GetRRRRRRRRR(){
	RullesMentalOutputStr=getCur10lastMentalRules()
}

// вывести 10 последних Правил на Пульт в http://go/pages/rulles.php
func getCur10lastMentalRules()string{
	rCount:=lastrulesMentalID
	if rCount >10{
		rCount=10
	}
	var out=""
	for i := 0; i < rCount; i++ {
		r:=rulesMentalArr[lastrulesMentalID-i]
		out+="ID="+strconv.Itoa(r.ID)+":"
		for n := 0; n < len(r.TAid); n++ {
			taa:=MentalTriggerAndActionArr[r.TAid[n]]
			if taa == nil{
				continue
			}
			if n>0{
				out+="<span style='padding:40px;'></span>"
			}else{
				out+="<span style='padding:10px;'></span>"
			}
			if taa.Trigger >0 {
				out += "<b>Стимул (оператора):</b> <span style='background-color:#FFECEB;'>" + GetActionsString(taa.Trigger) + "</span> "
			}
			if taa.Trigger <0 {
				out += "<b>Стимул (ментальный):</b> <span style='background-color:#FFECEB;'>" + GetMentalActionsString(-taa.Trigger) + "</span> "
			}
			out+="=> <b>Ответ:</b> <span style='background-color:#E8E8FF;'>"+GetMentalActionsString(taa.Action)+"</span> "
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
func getSuitableMentalRules()(int){
	var rID=0
	var activationType=2
// попытка срочно найти действие, в опасной ситуации
	if CurrentInformationEnvironment.veryActualSituation || CurrentInformationEnvironment.danger{
		rID = getMentalRulesArrFromTrigger(currentTriggerID)

	}else{
/* попытка более обстоятельно найти в эпиз.памяти подходящий фрагмент
	Чем больше limit тем маловероятнее найти совпадения сочетания Правил в ранней эпизодюпамяти,
   так что можно вызывать getRulesFromEpisodicsSlice постепенно уменьшая limit
Чем больше 	limit тем точнее результат обобщения, но меньше вероятность нахождения данного сочетания Правил
 */
		maxSteps:=1000
		for limit:=5; limit > 1; limit-- {
			rID=getMentalRulesFromEpisodicsSlice(activationType,limit,maxSteps)
			if rID>0{
				return rID
			}
			maxSteps = maxSteps/2
		}
	}

	return rID
}
/////////////////////////////////////////////////////////////
/* Внимательно выбрать наилучшее Правило rulesID по действию с Пульта или измееннию состояния.
Мягкий поиск (с игнорирование второстепенного) - не по ID образов действий, а по детализованных действиям:
по словам фразы, тон игнорируется, и по действиям.
Это - очень ресурсоемко, так что пока не реализовано
 */
func getSuitableMentalRulesСarefully()(int){

	// TODO

	return 0
}
/////////////////////////////////////////////////
/* Найти последнее известное Правило по цепочке последних limit кадров эпиз.памяти (шаблон решений)
с учетом шаблона уже реализованных правил.

Последовательность шаблона использует прежний опыт цепочек Стимул-Ответ-Эффект
с ожиданием очередного Стимула для последуего нахождения подходящего Ответа.

Каждый раз, находя последнее правило в данной ситуации,
оно может использоваться для того, чтобы перейти к следующему известному правилу
или, если такого нет, начать фрмировать новое решение
ментальным автоматизмом получения инфы с запуском ментального осмысления.

Получается ветвление дерева решений (основанного на эпиз.памяти) по каждому Стимулу

Такой поиск одинаково работате как для кадров объективной эпизод.памяти,
давая решение для запуска моторных автоматизмов,
так и для кадров субъектиных - давая решения запуска метнальным автоматизмам.

Чем больше limit тем маловероятнее найти совпадения,
так что можно вызывать getRulesFromEpisodicsSlice постепенно уменьшая limit
 */
func getMentalRulesFromEpisodicsSlice(activationType int,limit int,maxSteps int)(int){

	rImg:=getLastMentalRulesSequenceFromEpisodeMemory(activationType,limit)

// найти такую последовательность в предыдущей эпизод.памяти, но не далее 1000 фрагментов
/* уже обеспечено
if len(rImg)>limit{// limit последних
	rImg=rImg[len(rImg)-limit:]
}
*/
	lenFrag:=len(rImg)
	steps:=0
	lenEp:=len(EpisodeMemoryObjects)
	var startF = lenEp - 2*lenFrag // отмотать на 2 длины, чтобы не проверять в rID саму себя
	if startF > lenEp{//  а нет еще достаточной длины еп.памяти
		return 0
	}
		// идем назад по кускам lenFrag
/*TODO хорошо бы оптимизировать функцию так, чтобы можно было листать назад по -=lenFrag
если только это принципиально возможно
 */
		for i := startF; i >= 0; i -- { // =lenFrag - пролетает мимо
			if steps>maxSteps{
				return 0
			}
			var isConc=true
			var lastEM EpisodeMemory
			for j := 0; j < lenFrag; j++ {
				lastEM =*EpisodeMemoryObjects[i+j]
				if lastEM.RulesID != rImg[j] {
					isConc=false
					break
				}
			}
			if isConc{// уж ты, нашли такой же фрагмент! но в нем нет пускового curActiveActions (раньше уже смотрели)
				// выдать конечное праило, если оно с хорошим эффектом
				rArr:=rulesMentalArr[lastEM.RulesID]
				lastTa:=rArr.TAid[len(rArr.TAid)-1:]
				ta:=TriggerAndActionArr[lastTa[0]]
				if ta !=nil {
						if ta.Effect >0{// с хорошим эффектом
							return lastEM.RulesID
						}//else - продолжает искать хороший конец далее назад
				}
			}
			steps++
		}

	return 0
}
///////////////////////////////////////////////////

/*

 */
func getLastMentalRulesSequenceFromEpisodeMemory(activationType int,limit int)([]int){
	if EpisodeMemoryLastIDFrameID==0{
		return nil
	}
	var kind=0 // объективнй тип эпизод.памяти
	if activationType==2{
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
		if em == nil || em.Type != kind ||
			(em.LifeTime - preLifeTime) >EpisodeMemoryPause ||
			beginID >=limit{
			break // закончить выборку
		}
		beginID++
	}
	if beginID == 0 {
		return nil
	}
	var rImg []int
	// перебор последнего фрагмента кадров эпиз.памяти
	for i := EpisodeMemoryLastIDFrameID - beginID; i <= EpisodeMemoryLastIDFrameID; i++ {
		em := EpisodeMemoryObjects[i]
		rImg = append(rImg, em.RulesID)
	}
	if len(rImg)>1{
		createNewlastrulesMentalID(0, rImg)// записать (если еще нет такого) групповое правило

		return rImg
	}

	return nil
}
/////////////////////////////////////////////////



///////////////////////////////////////////////////
/*  быстро выбрать самое лучшее правило из rulesArr по пусковому стимулу
используя шаблоном последнюю цепочку кадров эпизод. памяти.
 */
func getMentalRulesArrFromTrigger(trigID int)(int){
	// сначала попробовать найти Правило с учетом тематического контекста
	for limit:=5; limit > 1; limit-- {
		rImg := getLastRulesSequenceFromEpisodeMemory(1, limit)
		rules := getRulesFromTemp(rImg, limit)
		if rules>0{
			return rules
		}
	}
// раз не нашли, то смотрим одиночные правила
	for k, v := range rulesMentalArr {
		for i := 0; i < len(v.TAid); i++ {
			if trigID!=v.TAid[i] || len(v.TAid)>1{
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
func getMentalRulesFromTemp(rImg []int,limit int)(int){
	for _, v := range rulesMentalArr {
		if len(v.TAid)!=limit{
			continue
		}
		if lib.EqualArrs(rImg, v.TAid){
			lastTa:=v.TAid[len(v.TAid)-1:]
			ta:=TriggerAndActionArr[lastTa[0]]
			if ta !=nil {
				if ta.Effect >0{// с хорошим эффектом
					return lastTa[0]
				}//else - продолжает искать хороший конец далее назад
			}
		}
	}
	return 0
}
////////////////////////////////////////////////
