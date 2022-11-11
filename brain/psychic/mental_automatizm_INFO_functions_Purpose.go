/* Инфофункция поиска Цели

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

	// определить эмоциональный контекст по newEmotionID или прямо по CurrentEmotionReception
	contextArr:=CurrentEmotionReception.BaseIDarr

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

	/* ищем PurposeImage.actionID с учетом текущего эмоционального контекста
	найти - с учетом Правил!!!!
	На стадии 4 - провоцировать оператора на ответы (почему, зачем, что такое?)
	*/
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


	// мозжечковые рефлексы - самый первый уровень осознания - подгонка действий под заданную Цель.
	if EvolushnStage == 4 {
		/*


		 */
	}

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

