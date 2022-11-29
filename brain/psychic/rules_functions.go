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
	rImg:=getLastRulesSequenceFromEpisodeMemory(0,5)
	if rImg!=nil {
		createNewRules(0, detectedActiveLastNodPrevID,detectedActiveLastUnderstandingNodPrevID,rImg,true) //записать (если еще нет такого) групповое правило
	}
}
//////////////////////////////////////////////////

////////////////////////////////////////////////////////////////////////
/* на стадии >3 при каждом ответе на действия оператора - прописывать текущее ПРАВИЛО rules.
   А так же просматривать эпизод память взад макcимум на 6 шагов или до паузы в общении > EpisodeMemoryPause шагов,
		фиксируя цепочку правил.
*/
func fixNewRules(lastCommonDiffValue int) int {
	if LastAutomatizmWeiting == nil{
		return 0
	}

	// curStimulImageID - образ действий оператора перед Ответом для записи Правила
	if curStimulImageID == 0  || LastAutomatizmWeiting == nil {
		return 0
	}
	if curStimulImage.ActID == nil && curStimulImage.PhraseID ==nil{ // не писать Правила с пустым Стимулом
		return 0
	}

	answer:=LastAutomatizmWeiting.ActionsImageID// ответный образ действий Beast
	// не записывать Правило с curActiveActionsID==answer
	if answer == 0  || ActionsImageArr[answer]==nil  || curStimulImageID==answer{
		return 0}

	TriggerAndAction,_:=createNewlastTriggerAndActionID(0,curStimulImageID,answer,lastCommonDiffValue,true)
	if TriggerAndAction == 0{return 0}
	// Запись Правила (именно detectedActiveLastNodPrevID а не detectedActiveLastNodID)
	currentRulesID, _ = createNewRules(0, detectedActiveLastNodPrevID,detectedActiveLastUnderstandingNodPrevID,[]int{TriggerAndAction},true)
	if currentRulesID == 0{return 0}

	//lib.WritePultConsol("<span style='color:green'>Записано <b>ПРАВИЛО № "+strconv.Itoa(currentRulesID)+"</b></span>") // уже есть сообщение в createNewRules()

	// новый кадр эпизодической памяти, сохраняющий Правил
	newEpisodeMemory(currentRulesID,0) // запись эпизодической памяти saveEpisodicMenory()

	// теперь обрабатываем прошлую эпизодическую память (необязательно, т.к. при каждом поиске в эп.памяти это происходит)
	GetRulesFromEpisodeMemory()

	return currentRulesID
}
///////////////////////////////////////////////////////////////////////

/* Не записывать Правила по изменению состояния, а только - по стимулу от Оператора!
// записать Правило типа BaseStateImage Стимул - не от оператора, а при активации изменением состояния
func fixRulesBaseStateImage(lastCommonDiffValue int){
	//корректируется успешность автоматизма - как в calcAutomatizmResult
	automatizmCorrection(lastCommonDiffValue,nil)
	/////////////////////// ПРАВИЛО:
	stimul, _ := CreateNewStatImageID(0, curStimulImage.Mood, curStimulImage.EmotionID, curBaseStateImage.SituationID,true)
	stimul*=-1 // отрицательное значение идентифицирует образ - как текущего сосотояния!!!
	fixNewRules(lastCommonDiffValue,stimul)
}*/
/////////////////////////////////////////////////////////////////////

/*Отзеркаливает авторитерный ответ Оператора на совершенное действие.
Записать одиночное Правило как Оператор отвечает на действиЯ Beast -
	авторитетное Правило всегда имеет эффект +1
Такое Правило используется в случае отсуствия решения как отвечать,
т.к. не пришется групповое Правило дял точного бездумного реагирования
(хотя на уровне эпиз.памяти и можно вычленять такие групповые Правила,
выделяя Стимул следующего Правила как ответ на действия Beast).
 */
func fixNewTeachRules()int{
	if LastAutomatizmWeiting == nil{
		return 0
	}

	// curActiveActionsID - образ действий оператора после Ответа для записи Правила
	if curActiveActionsID == 0  || LastAutomatizmWeiting == nil {
		return 0
	}
	if curActiveActions.ActID == nil && curActiveActions.PhraseID ==nil{ // не писать Правила с пустым Стимулом
		return 0
	}

	/* Если есть ли автоматизм с действием оператора curActiveActionsID, и если у него atmtzm.Usefulness<0 - снять блокировку
	потому как это - новое авторитарное подтвержение полезности.
	*/
	checkForUnbolokingAutomatizm(curActiveActionsID)

	answer:=LastAutomatizmWeiting.ActionsImageID// образ действий Beast перел ответом Оператора
	// не записывать Правило с curActiveActionsID==answer
	if answer == 0  || ActionsImageArr[answer]==nil || curActiveActionsID==answer{
		return 0}
	TriggerAndAction,_:=createNewlastTriggerAndActionID(0,answer,curActiveActionsID,1,true)
	if TriggerAndAction == 0{return 0}
	// Запись Правила
	currentRulesID, _ = createNewRules(0, detectedActiveLastNodPrevID,detectedActiveLastUnderstandingNodPrevID,[]int{TriggerAndAction},true)
	if currentRulesID == 0{return 0}

	return currentRulesID
}
/////////////////////////////////////////////////////////////////////


// удалить авторитарное Правило
func tryRemoveRules(atmtzm *Automatizm){
	if atmtzm==nil{
		return
	}
	treeNodeID:=atmtzm.BranchID
	actID:=atmtzm.ActionsImageID
	for k, v := range rulesArr { // пропускаем одиночные правила len(v.TAid) == 1
		if len(v.TAid) == 1 && v.NodeAID != treeNodeID {
			ta:=TriggerAndActionArr[v.TAid[0]]
			if ta.Trigger==actID && ta.Effect==1{
				//Обобщенное Правило пусть остается
				//delete(TriggerAndActionArr, v.TAid[0])
				delete(rulesArr, k)
				return
			}
		}
	}
}
///////////////////////////////////////////////////////////////////




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
/*  выбрать наилучшее Правило rulesID для текущего действию Оператора с Пульта (curActiveActionsID)
Текущая ситуация - массив самых последних кадров эпизодической памяти и
активный пусковой стимул curActiveActionsID.
*/
func getSuitableRules()(int){


// попытка срочно найти действие, в опасной ситуации
	if CurrentInformationEnvironment.veryActualSituation || CurrentInformationEnvironment.danger{
		rID,exact := getRulesArrFromTrigger(curActiveActionsID,true)
		if exact>0{// пойдет хоть что-то...
			return rID
		}
	}else{ //
		rID,exact := getRulesArrFromTrigger(curActiveActionsID,false)
		if exact>5{// пойдет...
			return rID
		}
	}

	return 0
}
//////////////////////////////////////////////////////////
/////////////////////////////////////////////////
// для выбора в func getRulesArrFromTrigger()

///////////////////////////////////////////////////
/*  быстро выбрать ранее успешное правило из rulesArr для данных условий
и заданного Стимула trigID типа ActionsImage
Алгоритм:
.....................

Возвращает наиболее совпадающее ПОЛЕЗНОЕ (Effect>0) Правило:
1 - ID Правила, совпадающие по цепочке эп.памяти с уверенностью exact
2 - точность совпадения: если 10...15 - для полны0 условий, 5-9 - только для условий дерева автоматизмов
*/
func getRulesArrFromTrigger(trigger int,veryActualSituation bool)(int,int){

	/* В конце эпиз.памяти еще нет Правила с новым Стимулом curActiveActionsID,
	но его последние Правила нужны чтобы по ним находить в групповых Правилах
	предшествовашие события: при совпадении цепочки последнего куска эпиз.памяти и групп.Правила
	велика вероятность верности такого реагирования.
	Поэтому сначала выделяем последнюю цепочку эпиз.памяти.
	*/
	//Вытащить из эпизод.памяти посленюю цепочку кадров, максимум в 5 кадров.
	rImg := getLastRulesSequenceFromEpisodeMemory(0,5)

	// полный образ текущих условий
	rulesID,exact:=searchingRules(trigger,rImg,0)
	if rulesID>0 { // не найдено для точного совпадения условий
		return rulesID,exact
	}

	//// не найдено для точного совпадения условий
	// смотрим тоолько для условий дерева автоматизмов
	rulesID,exact=searchingRules(trigger,rImg,1)
	if rulesID>0 { // не найдено для точного совпадения условий
		return rulesID,exact
	}

	if veryActualSituation {
		// смотрим безусловно (самый неувернный вариант)
		rulesID, exact = searchingRules(trigger, rImg, 2)
		if rulesID > 0 { // не найдено для точного совпадения условий
			return rulesID, exact
		}

		// для условия дерева автоматизмов в одиночных Правилах выбираем наилучшее
		rulesID:=getBeastIDRulesFromCondA(detectedActiveLastNodID)
		if rulesID>0 { // не найдено для точного совпадения условий
			return rulesID,1
		}
	}


	return 0,0
}
//
func searchingRules(trigger int,rImg []int,condType int )(int,int){
	type гUsefool struct {
		rID int
		exact int
	}
	var гUsefoolArr []гUsefool
	// текущие значения
	exact:=0 // точность совпадения
	rules:=0 //
	for _, v := range rulesArr {
		switch condType{
		// с учетом обоих деревьев
		case 0: if v.NodeAID!=detectedActiveLastNodID || v.NodePID!=detectedActiveLastUnderstandingNodID{
			continue
		}
		// с учетом только дерева автоматизмов
		case 1: if v.NodeAID!=detectedActiveLastNodID {
			continue
		}
		}
		exact=0 // точность совпадения
		rules=0 //
		for i := 0; i < len(v.TAid); i++ {
			rul:=TriggerAndActionArr[v.TAid[i]]
			if rul==nil{ lib.WritePultConsol("Нет карты TriggerAndActionArr для iD="+strconv.Itoa(v.TAid[i]));return 0,0}
			if rul.Trigger == trigger && rul.Effect>0 { // есть такое эффективное Правило
				// уже есть Ответ
				rules = rul.ID
				exact = 1 // предварительное начальное значение
				//смотрим совпадения предыдущих звеньев Правила и rImg
				var eIndex = len(rImg) - 1 // последний кадр эпиз.памяти
				// уходим назад, начиная с пердыдущего звена от найденного
				for r := i - 1; r >= 0; r-- { // смотрим совпадения предыдущих звеньев rImg
					if eIndex < 0 {
						break
					}
					eR := rImg[eIndex]
					rulR := TriggerAndActionArr[eR]
					rulE := TriggerAndActionArr[v.TAid[r]]
					if rulR == nil {
						lib.WritePultConsol("Нет карты TriggerAndActionArr для iD=" + strconv.Itoa(eR));
						return 0, 0
					}
					if rulE == nil {
						lib.WritePultConsol("Нет карты TriggerAndActionArr для iD=" + strconv.Itoa(eR));
						return 0, 0
					}
					if rulR.ID == rulE.ID { // совпадает
						exact++ // более 5 не бывает
					}else{
						break
					}
					eIndex--
				}
			}
			// запоминаем лучший результат текущего группового правла
			if rules>0 {
				гUsefoolArr = append(гUsefoolArr, гUsefool{rules, exact})
			}
		}
	}
	if гUsefoolArr!=nil {
		maxExact:=0
		curR:=0
		// выбираем самое правильное Правило
		for i := 0; i < len(гUsefoolArr); i++ {
			if гUsefoolArr[i].exact>maxExact{
				curR=гUsefoolArr[i].rID
				maxExact=гUsefoolArr[i].exact
			}
		}
		switch condType{
		// с учетом обоих деревьев
		case 0: maxExact += 10
		// с учетом только дерева автоматизмов
		case 1: maxExact += 5
			//case 2:
		}
		return curR,maxExact
	}

	return 0,0
}
///////////////////////////////////////////////






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

	rImg:=getLastRulesSequenceFromEpisodeMemory(0,limit)

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


///////////////////////////////////////////////////
/*  Выбрать ранее успешное правило из rulesArr для данных условий
КОНКРЕТНО ДЛЯ определенного объекта внимания extremImportanceObject (objID int, kind int) в структуре ActionsImage,
используя шаблоном последнюю цепочку кадров эпизод. памяти.
Возвращает - ID одиночного Правила
Чем более частный вид объекта внимания extremImportanceObject, тем более недостовеный резаультат.
*/
func getRulesArrFromAttentionObject(objID int, kind int)(int){
	if extremImportanceObject == nil{
		return 0
	}
	// пустое ответное действие
	var act ActionsImage
// м.б. лучше взять за основу прошлый ответный образ действий Beast? хотя он уже не должен быть ответом на новый Стимул...
	//act:=*ActionsImageArr[LastAutomatizmWeiting.ActionsImageID]
// После Стимула ищется Ответ и еще нет действия дла него, нужно сформировать по образу extremImportanceObject
switch extremImportanceObject.kind{
case 1: // ID ActionsImage
	act = *ActionsImageArr[extremImportanceObject.objID]
case 2: // ID MentalActionsImages

case 3: // ID несловестного действия ActionsImage.ActID[n]
	act.ActID=nil
	act.ActID =append(act.ActID,extremImportanceObject.objID)
case 4: // ID Verbal - при активации дерева автоматизмов
verb:=VerbalFromIdArr[extremImportanceObject.objID]
if verb==nil{return 0}
	act.PhraseID=verb.PhraseID
	act.MoodID=verb.MoodID
	act.ToneID=verb.ToneID
case 5: // ID отдельной фразы Verbal.PhraseID[n]
	act.PhraseID=nil
	act.PhraseID =append(act.PhraseID,extremImportanceObject.objID)
case 6:// ID отдельного слова  из Verbal.PhraseID[n]
	verb:=VerbalFromIdArr[extremImportanceObject.objID]
	if verb==nil{return 0}// м.б. применялось такое слово в Правилах... ??
	// сделать фразу из слова
	PhraseID:=AddwordIDToPhraseTree([]int{extremImportanceObject.objID})
	act.PhraseID=PhraseID
case 7:// ID тон сообщения с Пульта  Verbal.ToneID
	act.ToneID = extremImportanceObject.objID
case 8:// ID настроение оператора  Verbal.MoodID
	act.MoodID = extremImportanceObject.objID

}
	// образ ActionsImage
	ActionsImageID,_:=CreateNewlastActionsImageID(0,act.ActID,act.PhraseID,act.ToneID,act.MoodID,true)

	rules:=0
	// найти одиночное (а не групповое) Правило с учетом тематического контекста (групповые Правила)
		sinex:=strconv.Itoa(detectedActiveLastNodID)+"_"+strconv.Itoa(detectedActiveLastUnderstandingNodID);
		rArr:=rulesArrConditinArr[sinex] // все правила для данного индекса
		for _, v := range rArr {
			if len(v.TAid)>1{// найти одиночное (а не групповое)
				continue
			}
			lastTa:=v.TAid[len(v.TAid)-1:]
				ta:=TriggerAndActionArr[lastTa[0]]
				if ta !=nil && ta.Trigger==ActionsImageID{
					act:=ActionsImageArr[ta.Action]
					if act==nil{continue}
					if ta.Effect >0{// с хорошим эффектом
						rules = lastTa[0]
					}//else - продолжает искать хороший конец далее назад
				}
		}

	return rules
}
///////////////////////////////////////////////


/*Eсли для данного сочетания Стимул-Ответ есть только один вид эффекта,
то это - уже Информация, и чем больше опыт (количество обобщенных правил), тем такая информация полезнее.
Найти доминирующий эффект для сочетания Стимул-Ответ
 */
func getDominantEffect(triggerID int, actionID int)int{
var effect=0
	for _, v := range TriggerAndActionArr {
		if v.Trigger==triggerID && v.Action==actionID{
			effect+=v.Effect
		}
	}
	if effect > 1{effect=1}
	if effect < -1{effect=-1}
	return effect
}
////////////////////////////////////////////////

// для условия дерева автоматизмов (NodeAID) в одиночных Правилах выбираем наилучшее
func getBeastIDRulesFromCondA(NodeAID int)int {
	type res struct {
		rulesID int
		sumEffect int
	}
	var rules []res
	var oldRes *res

	for k, v := range rulesArr {
		if len(v.TAid) > 1 || v.NodeAID != NodeAID {
			continue
		}
		if TriggerAndActionArr[k] == nil{
			continue
		}
		r:=TriggerAndActionArr[k]
		effect:=getDominantEffect(r.Trigger, r.Action)
		oldRes=nil
		for i := 0; i < len(rules); i++ {
			if rules[i].rulesID==r.ID{
				oldRes=&rules[i]
				rules[i].rulesID=r.ID
				rules[i].sumEffect+=effect
				break
			}
		}
		if oldRes == nil {
			rules = append(rules, res{r.ID,effect})
		}
	}
	maxE:=0
	rulesID:=0
	if rules!=nil && len(rules)>0 {
		for i := 0; i < len(rules); i++ {
			if rules[i].sumEffect > maxE {
				maxE = rules[i].sumEffect
				rulesID = rules[i].rulesID
			}
		}
	}
	if rulesID>0{
		return rulesID
	}
return 0
}
////////////////////////////////////////////////////////////////////



/* есть ли положительный эффект у Правила, следующего за действием автоматизма,
чтобы если он хороший, посчитать такое действие приемлемым и запустить автоматизм.
Для текущего detectedActiveLastNodID.
т.е. смотрим цепочку Правил,
в которой есть Правило с действием==actionsImageID,
а последующее Правило имеет эффект >0
 */
func isNextWellEffectFromActonRules(actionsImageID int)bool{
	// Правила, у которых TriggerAndAction.Action == actionsImageID
	var taIdArr []int
	for k, v := range TriggerAndActionArr {
		if v.Action == actionsImageID{
			taIdArr=append(taIdArr,k)
		}
	}
	if rulesArr == nil{
		return false
	}
	// смотрим групповые Правила, у которых есть такой TriggerAndAction при NodeAID==detectedActiveLastNodID
	for _, v := range rulesArr {// пропускаем одиночные правила len(v.TAid) == 1
		if len(v.TAid) == 1 || v.NodeAID != detectedActiveLastNodID {
			continue
		}
		for i := 0; i < len(v.TAid); i++ {
			ta:=TriggerAndActionArr[v.TAid[i]]
			if ta.Action==actionsImageID{
				taNext:=TriggerAndActionArr[v.TAid[i]+1]
				if taNext != nil && taNext.Action == actionsImageID && taNext.Effect>0{
					return true
				}
			}
		}
	}
	return false
}
////////////////////////////////////////////////