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
limit 5 ограничивает выборку из эпиз.памяти, но она может получться и меньше.
*/
func GetRulesFromEpisodeMemory(){
	rImg:=getLastRulesSequenceFromEpisodeMemory(5)
	if rImg!=nil {
		createNewlastrulesID(0, detectedActiveLastNodID,detectedActiveLastUnderstandingNodID,rImg,true)//записать (если еще нет такого) групповое правило
	}
}
//////////////////////////////////////////////////




// вывести 10 последних Правил на Пульт в http://go/pages/rulles.php
func GetCur10lastRules()string{
	rCount:=lastrulesID
	if rCount >10{
		rCount=10
	}
	var out=""
	for i := 0; i < rCount; i++ {
		r:=rulesArr[lastrulesID-i]
		out+="ID="+strconv.Itoa(r.ID)+" для <span title='ID дерева автоматизмов'>"+strconv.Itoa(r.NodeAID)+"</span> и <span title='ID дерева понимания'>"+strconv.Itoa(r.NodePID)+"</span> :"
		for n := 0; n < len(r.TAid); n++ {
			taa:=TriggerAndActionArr[r.TAid[n]]
			if taa == nil{
				out+="<span style='padding:40px;'></span><span style='color:red;' title='Нет образа TriggerAndActionArr "+strconv.Itoa(r.TAid[n])+"'>нет "+strconv.Itoa(r.TAid[n])+"&nbsp;&nbsp;&nbsp;&nbsp;</span>"
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
func getSuitableRules()(int){

	rID,doubt := getRulesArrFromTrigger(currentTriggerID)
// попытка срочно найти действие, в опасной ситуации
	if CurrentInformationEnvironment.veryActualSituation || CurrentInformationEnvironment.danger{
// на смотрим на сомнительность doubt
		return rID
	}else{ //
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
			rID,_=getRulesFromEpisodicsSlice(limit,maxSteps)
			if rID>0{
				return rID
			}
			maxSteps = maxSteps/2
		}*/
	}

	return 0
}

/////////////////////////////////////////////////
/* Найти последнее известное Правило по цепочке последних limit кадров эпиз.памяти (шаблон решений)
с учетом шаблона уже реализованных правил.

Последовательность шаблона использует прежний опыт цепочек Стимул-Ответ-Эффект
с ожиданием очередного Стимула для последнего нахождения подходящего Ответа.

Каждый раз, находя последнее правило в данной ситуации,
оно может использоваться для того, чтобы перейти к следующему известному правилу
или, если такого нет, начать формировать новое решение
ментальным автоматизмом получения инфы с запуском ментального осмысления.

Получается ветвление дерева решений (основанного на эпиз.памяти) по каждому Стимулу

Такой поиск одинаково работает как для кадров объективной эпизод.памяти,
давая решение для запуска моторных автоматизмов,
так и для кадров субъективных - давая решения запуска ментальным автоматизмам.

Чем больше limit тем маловероятнее найти совпадения,
так что можно вызывать getRulesFromEpisodicsSlice постепенно уменьшая limit

 */
/*
смотреть собственную память имеет смысл только в виде цепочки эпизодов - для прогнозирования предположительных
последствий запуска такой цепочки. Отдельные эпизоды просто показывают срабатывания автоматизмов.
Но если смотреть память оператора, то можно продолжить отзеркаливать его реакции - это будет развитие простейшего отзеркаливания
из 3 стадии, но уже без попугайского провоцирования.
Для этого нужна функция, которая "сдвигает по фазе" звенья эпизода. Например:
	Стимул1 - Ответ1 - Эффект1 | Стимул2 - Ответ2 - Эффект2 | Стимул2 - Ответ2 - Эффект2 ...
преобразуются в цепочку:
Стимул10(Ответ1) - Ответ10(Стимул2) - Count1 | Стимул11(Ответ2) - Ответ11(Стимул2) - Count2
здесь эффекта нет потому, что мы не знаем, как отразилось на операторе его действие, а уверенно гадать на 4 стадии еще не умеем.
Поэтому просто ищем среди этих пар подходящий стимул, ориентируясь на Count как число совпадений такой пары при преобразовании.
	То есть сколько раз была зафкисирована в памяти такая реакци оператора - это будет фактором уверенности ее отзеркаливания, если
придется выбирать из нескольких вариантов. Поиск будет быстрее, если отсортировать по Count по убыванию.
	Отзеркаленный автоматизм тем не менее не может быть авторитарным несмотря ни на какую Count, только пробным.
	Такое отзеркаливание, с оценкой прошлого опыта, можно назвать более осмысленным, по сравнению с попугайским рефлекторным на 3 стадии.

Возвращает:
1 - ID Правила
2 - index эпиз.памяти с таким Правилом
*/

func getRulesFromEpisodicsSlice(limit int,maxSteps int)(int,int){

	rImg:=getLastRulesSequenceFromEpisodeMemory(limit)

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
				lastEM=EpisodeMemoryObjects[i+j]
				if lastEM.TriggerAndActionID != rImg[j] {
					isConc=false
					break
				}
			}
			if lastEM == nil{
				return 0,0
			}
			if isConc{// уж ты, нашли такой же фрагмент! но в нем нет пускового curActiveActions (раньше уже смотрели)
				// выдать конечное праило, если оно с хорошим эффектом
				ta:=TriggerAndActionArr[lastEM.TriggerAndActionID]
				if ta !=nil {
						if ta.Effect >0{// с хорошим эффектом
/* TODO тут можно посмотреть далее на сколько-то шагов вперед чтобы прикинуть, чем закончится комбинацияя Стимул-Ответ
это - как думать на сколько-то шагов врепед в шахматах. Можно запустить цикл обдумывания.
Найденный ta.Effect >0 - это и есть примитивная ЦЕЛЬ, в отличие от Доминанты нерешенной проблемы.
							!!!!!!!!!!!!!!
							Получая эффект последнего эпизода цепочки мы по сути просматриваем цепочку на предмет последствий ее выполнения.
							Стало быть эта функция должна только выдавать true/false как оценку, стоит ли начинать реагировать по этому правилу.
							А не запускать действия последнего эпизода цепочки (как это делается потом через createAndRunAutomatizmFromPurpose()) потому,
							что получится что то типа:
							найденный фрагмент эпиз. памяти:  привет - хай/как дела - так себе/ну и дурак - сам дурак
							если запустить реакцию последнего эпизода то получим: привет - сам дурак.
							*/
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
						}
//else - продолжает искать хороший конец далее назад, хотя это уже менее вероятно, но в прощлом при меньшей  длине шаблона можно найти.
				}
			}
			steps++
		}

	return 0,0
}
///////////////////////////////////////////////////

/* Вытащить из эпизод.памяти посленюю цепочку кадров
 */
func getLastRulesSequenceFromEpisodeMemory(limit int)([]int){
	if EpisodeMemoryLastIDFrameID==0{
		return nil
	}
	var kind=0 // здесь всегда - объективнй тип эпизод.памяти
	var beginID=0
	var preLifeTime=0
	for i := EpisodeMemoryLastIDFrameID; i >=0; i-- {
		em:=EpisodeMemoryObjects[i]
		// если самый последний эпизод уже является em.Type == kind
		if i==EpisodeMemoryLastIDFrameID && em.Type == kind || beginID > 5{
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
	//beginID+1 чтобы число проходов цикла было равно beginID и окончился на i <= EpisodeMemoryLastIDFrameID
	for i := EpisodeMemoryLastIDFrameID - beginID+1; i <= EpisodeMemoryLastIDFrameID; i++ {
		em := EpisodeMemoryObjects[i]
		rImg = append(rImg, em.TriggerAndActionID)
	}
	if len(rImg)>1{
		createNewlastrulesID(0, detectedActiveLastNodID,detectedActiveLastUnderstandingNodID,rImg,true)// записать (если еще нет такого) групповое правило

		return rImg
	}

	return nil
}
/////////////////////////////////////////////////



///////////////////////////////////////////////////
/*  быстро выбрать ранее успешное правило из rulesArr для данных условий
используя шаблоном последнюю цепочку кадров эпизод. памяти.

Возвращает:
1 - ID Правила
2 - неуверенность нахождения: 0 - макисмальная уваеренность, чем ниже, тем неопределеннее
 */
func getRulesArrFromTrigger(trigID int)(int,int){
	doubt:=0 // сомнение
	// сначала попробовать найти Правило с учетом тематического контекста (групповые Правила)
	for limit:=5; limit > 1; limit-- {
		//Вытащить из эпизод.памяти посленюю цепочку кадров
		rImg := getLastRulesSequenceFromEpisodeMemory(limit)
		sinex:=strconv.Itoa(detectedActiveLastNodID)+"_"+strconv.Itoa(detectedActiveLastUnderstandingNodID);
		rArr:=rulesArrConditinArr[sinex] // все правила для данного индекса
		rules:=0
		for _, v := range rArr {
			if len(v.TAid)!=limit {
				continue
			}
			if lib.EqualArrs(rImg, v.TAid){
				lastTa:=v.TAid[len(v.TAid)-1:]
				ta:=TriggerAndActionArr[lastTa[0]]
				if ta !=nil {
					if ta.Effect >0{// с хорошим эффектом
						rules = lastTa[0]
					}//else - продолжает искать хороший конец далее назад
				}
			}
		}
		if rules>0{
			return rules,doubt
		}
		doubt++
	}

	return 0,10
}
///////////////////////////////////////////////
/*
func getRulesFromTemp(rImg []int,limit int)(int){
	for _, v := range rulesArr {
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
