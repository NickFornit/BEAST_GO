/* функции ментальных Правил


Функция мягкого поиска по деталям образов действий getSuitableMentalRulesСarefully()
который можно использовать в инфо-функциях.
*/


package psychic

import (
	"BOT/lib"
	"strconv"
	"strings"
)

///////////////////////////////////////////////
/* Создать групповое правило (более одного в цепочке EpisodeMemory.TAid)
из последнего участка эпизодической памяти объектиынх (EpisodeMemory.Type==0) элеметов.
limit 5 ограничивает выборку из эпиз.памяти, но она может получться и меньше.
*/
func GetMentalRulesFromEpisodeMemory(){
	rImg:=getLastRulesSequenceFromEpisodeMemory(5)
	if rImg!=nil {
		createNewlastrulesMentalID(0, detectedActiveLastNodID,detectedActiveLastUnderstandingNodID,rImg,true)//записать (если еще нет такого) групповое правило
	}
}
//////////////////////////////////////////////////

// вывести 10 последних Правил на Пульт в http://go/pages/rulles.php
func GetCur10lastMentalRules()string{
	rCount:=lastrulesMentalID
	if rCount >10{
		rCount=10
	}
	var out=""
	for i := 0; i < rCount; i++ {
		r:=rulesMentalArr[lastrulesMentalID-i]
		out+="ID="+strconv.Itoa(r.ID)+" для <span title='ID дерева автоматизмов'>"+strconv.Itoa(r.NodeAID)+"</span> и <span title='ID дерева понимания'>"+strconv.Itoa(r.NodePID)+"</span> :"
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
			if taa.ShortTermMemoryID!=nil  {
				out += "<b>Цепочка осмысления:</b>"
				CickleStr :=""
				for c := 0; c < len(taa.ShortTermMemoryID); c++ {
					if c>0{CickleStr+=", "}
					CickleStr+=strconv.Itoa(taa.ShortTermMemoryID[c])
				}
				out += "<span style='background-color:#FFECEB;cursor:pointer' onClick='show_cickles_strings(`"+CickleStr+"`)'>" + CickleStr + "</span> "
			}
			out+="=> <b>Ответ:</b> <span style='background-color:#E8E8FF;'>"+GetMentalActionsString(taa.Action)+"</span> "
			out+="<b>Эффект: "+strconv.Itoa(taa.Effect)+"</b>"
			out+="<br>"
		}
		out+="<hr>"
	}
	return out
}
/* для http://go/pages/mental_rules.php расшифровать цикл ментиального Правила типа "2, 4, 0, 0, 0, 0"
По строк перечисления ID goNext
 */
func GetMentalRulesCickleInfo(cickle string)string{
	if len(cickle)==0{
		return "Пустой цикл."
	}
	сArr:=strings.Split(cickle, ",")
	if len(сArr)==0{
		return "Пустой цикл."
	}
	out:="<br><b>Информация о звене цикла:</b><br>"
	out+="<table cellpadding=0 cellspacing=0 border=1 class='main_table'>"
	out+="<tr><th class='table_header'>goNext ID</th>"
	out+="<th class='table_header'>ID дерева автоматизмов</th>"
	out+="<th class='table_header'>ID ментального автоматизма</th></tr>"
	for i := 0; i < len(сArr); i++ {
		sID:=strings.Trim(сArr[i]," ")
		id,_:=strconv.Atoi(sID)
		if id == 0{
			continue
		}
		out+=GetGoNextInfo(id)
	}
	out+="</table>"
return out
}
///////////////////////////////////////////


///////////////////////////////////////////
/*  выбрать наилучшее Правило rulesID по действию с Пульта или измееннию состояния
Текущая ситуация - массив самых последних кадров эпизодической памяти и
активный пусковой стимул currentTriggerID типов curActiveActions или curBaseStateImage.
*/
func getSuitableMentalRules()(int){
	rID,doubt := getMentalRulesArrFromTrigger(saveFromNextIDAnswerCicle)
// попытка срочно найти действие, в опасной ситуации
	if CurrentInformationEnvironment.veryActualSituation || CurrentInformationEnvironment.danger{
		// на смотрим на сомнительность doubt
		return rID
	}else{
/* попытка более обстоятельно найти в эпиз.памяти подходящий фрагмент
	Чем больше limit тем маловероятнее найти совпадения сочетания Правил в ранней эпизодюпамяти,
   так что можно вызывать getRulesFromEpisodicsSlice постепенно уменьшая limit
Чем больше 	limit тем точнее результат обобщения, но меньше вероятность нахождения данного сочетания Правил
 */
		if doubt<3{// пойдет...
			return rID
		}
		/*
		maxSteps:=1000
		for limit:=5; limit > 1; limit-- {
			rID,_=getMentalRulesFromEpisodicsSlice(limit,maxSteps)
			if rID>0{
				return rID
			}
			maxSteps = maxSteps/2
		}*/
	}

	return 0
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

Возвращает:
1 - ID Правила
2 - index эпиз.памяти с таким Правилом
 */
func getMentalRulesFromEpisodicsSlice(limit int,maxSteps int)(int,int){

	//Вытащить из эпизод.памяти посленюю цепочку кадров
	rImg:=getLastMentalRulesSequenceFromEpisodeMemory(limit)

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
	if startF <0{//  а нет еще достаточной длины еп.памяти
		return 0,0
	}
	//Поиск - в контексте активных деревьев detectedActiveLastNodID и detectedActiveLastUnderstandingNodID
	var rеserve=0 // резервные Правила, если не найдено точно в контексте
	var index=0
		// идем назад по кускам lenFrag
/*TODO хорошо бы оптимизировать функцию так, чтобы можно было листать назад по -=lenFrag
если только это принципиально возможно
 */
		for i := startF; i >= 0; i -- { // =lenFrag - пролетает мимо
			if steps>maxSteps{
				return 0,0
			}
			var isConc=true
			var lastEM *EpisodeMemory
			for j := 0; j < lenFrag; j++ {
				lastEM =EpisodeMemoryObjects[i+j]
				if lastEM.TriggerAndActionID != rImg[j] {
					isConc=false
					break
				}
			}
			if lastEM == nil{
				return 0,0
			}
			if isConc{// нашли такой же фрагмент! но в нем нет пускового curActiveActions (раньше уже смотрели)
				// выдать конечное праило, если оно с хорошим эффектом
				ta:=MentalTriggerAndActionArr[lastEM.TriggerAndActionID]
				if ta !=nil {
						if ta.Effect >0{// с хорошим эффектом
							rеserve=lastEM.TriggerAndActionID
							index=i
							if lastEM.NodeAID == detectedActiveLastNodID{
								rеserve=lastEM.TriggerAndActionID
								index=i
								if lastEM.NodePID == detectedActiveLastUnderstandingNodID{
									return lastEM.TriggerAndActionID,i
								}
							}
							return rеserve,index
						}//else - продолжает искать хороший конец далее назад
				}
			}
			steps++
		}

	return 0,0
}
///////////////////////////////////////////////////

/* Вытащить из эпизод.памяти посленюю цепочку кадров
 */
func getLastMentalRulesSequenceFromEpisodeMemory(limit int)([]int){
	if EpisodeMemoryLastIDFrameID==0{
		return nil
	}
	var kind=1 // ментальный тип эпизод.памяти

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
	for i := EpisodeMemoryLastIDFrameID - beginID+1; i <= EpisodeMemoryLastIDFrameID; i++ {
		em := EpisodeMemoryObjects[i]
		rImg = append(rImg, em.TriggerAndActionID)
	}
	if len(rImg)>1{
		createNewlastrulesMentalID(0, detectedActiveLastNodID,detectedActiveLastUnderstandingNodID,rImg,true)// записать (если еще нет такого) групповое правило

		return rImg
	}

	return nil
}
/////////////////////////////////////////////////



///////////////////////////////////////////////////
/*  быстро выбрать самое лучшее правило из rulesArr по пусковому стимулу
используя шаблоном последнюю цепочку кадров эпизод. памяти.
mentalID - saveFromNextIDAnswerCicle []int

Возвращает:
1 - ID Правила
2 - неуверенность нахождения: 0 - макисмальная уваеренность, чем ниже, тем неопределеннее
*/
func getMentalRulesArrFromTrigger(mentalID []int)(int,int) {
	// сначала попробовать найти Правило с учетом тематического контекста
	for limit := 5; limit > 1; limit-- {
		//Вытащить из эпизод.памяти посленюю цепочку кадров
		rImg := getLastMentalRulesSequenceFromEpisodeMemory(limit)
		doubt:=0 // сомнение
		sinex := strconv.Itoa(detectedActiveLastNodID) + "_" + strconv.Itoa(detectedActiveLastUnderstandingNodID);
		rArr := rulesMentalArrConditinArr[sinex] // все правила для данного индекса
		rules := 0
		for _, v := range rArr {
			if len(v.TAid) != limit {
				continue
			}
			if lib.EqualArrs(rImg, v.TAid) {
				lastTa := v.TAid[len(v.TAid)-1:]
				ta := MentalTriggerAndActionArr[lastTa[0]]
				if ta != nil {
					if ta.Effect > 0 { // с хорошим эффектом
						rules = lastTa[0]
					} //else - продолжает искать хороший конец далее назад
				}
			}
		}
		if rules > 0 {
			return rules, doubt
		}
		doubt++
	}

	return 0, 10
}
	/*
			// искать эту цепочку в групповых Правилах
			rules := getMentalRulesFromTemp(rImg, limit)
			if rules>0{
				return rules
			}
// раз не нашли, то смотрим одиночные правила
	for k, v := range rulesMentalArr {
		for i := 0; i < len(v.TAid); i++ {
			if len(mentalID) != len(v.TAid){
				continue
			}
			ta:=MentalTriggerAndActionArr[k]
			if ta !=nil {
				if lib.EqualArrs(mentalID,v.TAid) {// есть такой пусковой
					if ta.Effect >0{// первый попавшийся с хорошим эффектом
						return k
					}
				}
			}
		}
	}

	return 0
}*/
///////////////////////////////////////////////
/*
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
}*/
////////////////////////////////////////////////
