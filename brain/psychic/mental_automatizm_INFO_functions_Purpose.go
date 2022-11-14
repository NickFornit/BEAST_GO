/* Инфо-функция поиска Цели - как объекта произвольного внимания:
того из всего воспринимаемого, что имеет наибольшую значимость
т.к. именно наибольшая значимость должна осмысливаться.

targetID [] можно пока вообще не трогать, так же как и veryActual -
оба просто констатируют что нужно бы улучшить.

Тут главное - определить PurposeImage.actionID -
выбранный образ действия бота для достижения данной цели (ActionsImage)

т.е. если не определно произвольное PurposeImage.actionID или вообще не задана Цель,
то целью автоматически становится улучшение состояния, см. func valuationPurpose().

На пятой ступени развития целью может стать решение Доминанты нерешенной проблемы, и только
если нет Доминант, то ищем PurposeImage.actionID с учетом текущего эмоционального контекста.
*/


package psychic

import "BOT/lib"

//////////////////////////////////////////

//////////////////////////////////////////////////////////
/* Ментальное определение ближайшей Цели в данной ситуации
Это - для переактивации с новым PurposeImage (understanding_purpose_image.go)
Постановка цели для текущего цикла размышления, чтобы оценить эффект для Правила.
Учесть текущую эмоцию, в том числе переактивированную!!!
*/
func infoFunc8() {
	if mentalInfoStruct.mentalPurposeID >0 {
		getMentalPurpose()
	}
	// получили mentalInfoStruct.mImgID
	currentInfoStructId=8 // определение актуального поля mentalInfo
}
//////////////////////////////////////
// найти подходящий PurposeImage.actionID и определить mentalInfoStruct.mentalPurposeID
func getMentalPurpose(){
	mentalInfoStruct.mentalPurposeID=0

	// определяется текущий объект наибольшой значимости в воспринимаемом
	getGreatestImportance()


	if EvolushnStage > 4 { // главное - Доминанта нерешенной пробелмы
		dominfntaID := getMainDominanta(CurrentEmotionReception.ID)
		if dominfntaID > 0 {

// TODO: использовать Цель Доминанты и определить по ней mentalInfoStruct.mentalPurposeID= !!!!
//TODO:  придумать какой должна быть цель Доминанты, чтобы по ней можно было определить mentalInfoStruct.mentalPurposeID

			if mentalInfoStruct.mentalPurposeID >0 {
				createAndRunPurposeAutomatizm()
				return
			}
		}
	}
	///////////////////////////
	// если нет Доминант, то ищем PurposeImage.actionID с учетом текущего эмоционального контекста

	/* ищем PurposeImage.actionID в контексте активных деревьев
	найти - с учетом Правил!!!!
	??На стадии 4 - провоцировать оператора на ответы (почему, зачем, что такое?)
Нужно искать не только в контексте эмоции, а активных веток деревьев detectedActiveLastNodID и detectedActiveLastUnderstandingNodID
	Алгоритм:
Ищем в ОБЪЕКТИВНЫХ Правилах подходящее действие (текущий saveFromNextIDAnswerCicle),
	смотрим в последних кадрах эпизод.памяти такое Правило и в его продолжении есть позитивный эффект,
	то берем оттуда действие оператора.
	*/
	if EpisodeMemoryObjects!=nil {// еще нет эпиз.памяти, так что и цели нет...
		return
	}
// тупо - для 4-й ступени (и 5-й, если нет Доминанты), когда набирается значимость и Правила
		index := 0 // индекс в массиве эпизод.памяти, чтобы по нему смотреть, что было дальне
		rID, doubt := getRulesArrFromTrigger(currentTriggerID)
		if doubt > 3 { // слишком сомнительно

			rulexID := 0
			maxSteps := 1000
			for limit := 5; limit > 1; limit-- {
				rulexID, index = getRulesFromEpisodicsSlice(limit, maxSteps)
				if rulexID > 0 {
					break
				}
				maxSteps = maxSteps / 2
			}
		} else { // вполне уверенно найдено Правило, ищем его индекс в эпиз.памчяти
			for i := len(EpisodeMemoryObjects); i >= 0; i-- {
				ep := EpisodeMemoryObjects[i]
				if ep.Type == 0 && ep.TriggerAndActionID == rID {
					index = i
					break
				}
			}
		}
		//////////////////////////
		if index == 0 {
			return
		}

		// есть ли последующий кадр эпизод.памяти
		lastEM := EpisodeMemoryObjects[index+1]
		if lastEM == nil {
			return // придется обойтись без PurposeImage.actionID
		}
		// хороший ли эфект
		rArr := rulesArr[lastEM.TriggerAndActionID]
		if rArr == nil {
			return // придется обойтись без PurposeImage.actionID
		}
		// выдать конечное праило, если оно с хорошим эффектом
		var rеserve = 0 // резервные Правила, если не найдено точно в контексте
		ta := TriggerAndActionArr[lastEM.TriggerAndActionID]
		if ta == nil {
			return // придется обойтись без PurposeImage.actionID
		}
		if ta.Effect > 0 { // с хорошим эффектом
			rеserve = ta.Action
			if lastEM.NodeAID == detectedActiveLastNodID {
				rеserve = ta.Action
				if lastEM.NodePID == detectedActiveLastUnderstandingNodID {
					rеserve = ta.Action
				}
			}
		}
		if rеserve > 0 { // нашли действие
			if savePurposeGenetic == nil {
				getPurposeGenetic()
			}
			// создать Цель
			purposeID, _ := createPurposeImageID(0, savePurposeGenetic.veryActual, savePurposeGenetic.targetID, rеserve, true)
			// передать результат
			mentalInfoStruct.mentalPurposeID = purposeID
	}

	/*   оставляю как пример расшифровки эмоции
	// определить эмоциональный контекст по newEmotionID или прямо по CurrentEmotionReception
	contextArr:=CurrentEmotionReception.BaseIDarr
	for i := 0; i < len(contextArr); i++ {
		switch contextArr[i]{
		case 1:	//Пищевой	- Пищевое поведение, восполнение энергии, на что тратится время и тормозятся антагонистические стили поведения.

		case 2:	//Поиск	- Поисковое поведение, любопытство. Обследование объекта внимания, поиск новых возможностей.

		case 3:	//Игра	- Игровое поведение - отработка опыта в облегченных ситуациях или при обучении.

		case 4:	//Гон	- Половое поведение. Тормозятся антагонистические стили

		case 5:	//Защита	- Оборонительные поведение для явных признаков угрозы или плохом состоянии.

		case 6:	//Лень	- Апатия в благополучном или безысходном состоянии.

		case 7:	//Ступор	- Оцепенелость при непреодолимой опастbase_context_activnostности или когда нет мотивации при благополучии или отсуствии любых возможностей для активного поведения.

		case 8:	//Страх	- Осторожность при признаках опасной ситуации.

		case 9:	//Агрессия	- Агрессивное поведение для признаков легкой добычи или защиты (иногда - при плохом состоянии).

		case 10: //Злость	- Безжалостность в случае низкой оценки .

		case 11: //Доброта	- Альтруистическое поведение.

		case 12: //Сон - Состояние сна. Освобождение стрессового состояния. Реконструкция необработанной информации.

		}
	}
	*/



	// мозжечковые рефлексы - самый первый уровень осознания - подгонка действий под заданную Цель.

		/*


		 */


	createAndRunPurposeAutomatizm()

	return
}
////////////////////////////////////
func createAndRunPurposeAutomatizm(){
	if mentalInfoStruct.mentalPurposeID==0{
		return
	}
	// создать мент.автоматизм приоизвольной активацции ментальной цели
	actImgID,_:=CreateNewlastMentalActionsImagesID(0,3,mentalInfoStruct.mentalPurposeID,true)
	id, matmzm := createMentalAutomatizmID(0, actImgID, 1)
	if id >0 {
		// запустить мент.автоматизм
		RunMentalMentalAutomatizm(matmzm)

		lib.WritePultConsol("Найдена и активирована Цель.")
	}
}
//////////////////////////////////////////////////////////

