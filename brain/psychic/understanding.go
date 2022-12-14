/* Базовая система самоощущения (актуализация доступной информации),
функция поддержвания информационной среды в зависимости от текущй ситуации -
обдщая картина понимания ситуации,
с функцией лени - думать или не думать: мотивированность думать зависит от ситуации, ее опасности.
Решение что-делать или просто игнорировать ситуацию - важнейший параметр индивидуальной адаптивности.
Не просто искать способ выйти из опасного состояния, а искать как улучшить состояние,
каким бы оно ни было вплоть до полной неудовлетворенности существующим, когда поис ведется в любой ситуации.
Нужна функция, определяющая лень с индивидуализированными параметрами
(м.б. зависящими от более базовых индивид.параментров).
Функция поддерживается системой самоощущения,
которая работает независимо от активации дерева понимания по каждому пульсу.

Активация func consciousness - после включения или пробуждения - начало цепочки мыслей,
а так же - по активности дерева (ор.рефлекс) - срабатывает функция делать-неделать.

После объективной активации (activationType==1) начинается рекурсивный субъективный вызов (activationType==2)
- цикл обдумывания (субъективный ориентировочный рефлекс), каждой шаг которого основывается на информации, даваемой предыдущим шагом
с целью найти подходящие действия для данной ситуации, что дает возможность снова сориентироваться.

Кроме гомеостатического инфо.окрыжения psi_information_environment.go
есть psi_information_environment_mental.go куда помещаются результаты информационных функций.
Использует субъективную часть эпизодов памяти - субъективный тип (для записи эпизодов цепочки мыслей).

Функция consciousness() проходит через 4 уровня решений - выделенных комментариями.
Если решение не найдено на данном уровне, то отрабатывает следующий уровень сложности обработки (что характеризует эволюционную последовательность их появления).

РАБОТА consciousness() каждого уровня описана по месту.

*/


package psychic

import "BOT/lib"

//////////////////////////////
var AllowConsciousnessProcess=false // при включении и просыпании - 1 раз
var isFirstConsciousnessActivation=false // первый вызов consciousness при включении или просыпании

// сохранение значения уровня осмысления == стадии развития при произвольном изменении уровня
var saveEvolushnStage=0

// true - текущий режим активации consciousness - субъектиынй (activationType=2)
var isActivationType2=false

var currrentFromNextID=0 // текущий fromNextID в текущем запуске consciousness


// сколько раз был стимул от оператора после последнего запуска Ответа
var stimulCount = 0
// не было моторного ответа на прошлый стимул, а уже последовавл новый
var isСonfusion=false
var timeOfLastStimul=0 //время с прошлого Стимула от Оператора timeOfLastStimul=PulsCount - curActiveActionsPulsCount

var newMentCickle=true // начало прохода ментального цикла, чтобы не повторять сообщения

var functionsInThisCickles []int//ID запускаемых в текущем цикле инфо-функций
func setCurIfoFuncID(infofID int){
	functionsInThisCickles=append(functionsInThisCickles,infofID)
}

// в прошлом объективной активации был запущен автоматизм, так что не нужно прерывать размышление
var isNormalRunAutomatimsInPreActive=false

///////////////////////////////////////////////////////////





////////////////////////////////////////////////////////////
/* Главная, активная с каждым ориентировочным рефлексов функция циклов осмысления
для поддержания информационной среды и произвольности.
Изолированная от непосредственных воздействий и поэтому самостоятельная система оценки и корректировки состояния,
происходит с формированием опыта прозвольности выполнения Правил и их выбора для данных условий,
в частности, Правил относительно самого себя (самосознание).

Начинает работать с if EvolushnStage > 3

По каждой активации дерева автоматизмов и дерева понимания запускается рекурсивный цикл осознания: consciousness.

Вид активации - при вызове функции осознания.
activationType == 0 - не бывает такого значея
activationType == 1 - активация ориентировочным рефлексом новой ситуации
activationType == 2 - активация "внутренним" (произвольным) ориентировочным рефлексом

В принципе здесь должны исправляться все лажи ответов предыдущих периодов...

fromNextID - ID MentalNext на котором была запущена переактивация consciousness при обдумывании
*/
func consciousness(activationType int,fromNextID int)(bool) {   //  return false

	if !AllowConsciousnessProcess {// при включении и просыпании - 1 раз AllowConsciousnessProcess=true
		isFirstConsciousnessActivation=false
		return false
	}
	var isFirstActivation=false
	if !isFirstConsciousnessActivation{// пробуждение
		isFirstConsciousnessActivation=true
		isFirstActivation=true
		initMentalMemories()
	}

	if activationType == 1 {
// ТЕСТИРОВАНИЕ РАЗНЫХ ФУНКЦИЙ
//		if infoMirroringStimul() {	return true	}

		isNormalRunAutomatimsInPreActive=false


		isActivationType2 = false
		//
		if stimulCount > 1{//не было моторного ответа на прошлый стимул, а уже последовавл новый
			isСonfusion=true
		}
		extremImportanceObject=nil
		extremImportanceMentalObject=nil

		// посмотреть, есть ли прерванные цепочки осмысления, и если есть, выбрать, с какой продолжить осмысление.
	}
	if activationType == 2 {
		isActivationType2 = true

	}

	if currrentFromNextID != fromNextID{//сохранеиние цикла осмысления между двумя объективными вызовами consciousness
		saveFromNextIDcurretCicle=append(saveFromNextIDcurretCicle,fromNextID)
		saveFromNextIDAnswerCicle=append(saveFromNextIDAnswerCicle,fromNextID)
	}
	currrentFromNextID=fromNextID


	var stopMentalWork=false
	if activationType == 1 && isActivationType2{// объективная активация
		if !isNormalRunAutomatimsInPreActive {
			// нужно прервать выполнение циклов субъективныъ активаций
			stopMentalWork = true
		}
	}



	// вернуть стадию развития
	if saveEvolushnStage >0 {// иначе оно обнуляет EvolushnStage
		EvolushnStage = saveEvolushnStage // возвращаем уровень осмысления
	}
	///////////////////////////////////////////////
	if activationType == 1 {
		refreshCurrentInformationEnvironment()
		saveFromNextIDcurretCicle = nil
	}
	//////////////////////////////////////////////////////////

	var limitCickleCountForEvolushnStage4 = 4 // ограничить число циклов для 4-й стадии
	var limitCickleCountForEvolushnStage5 = 20 // ограничить число циклов для 5-й стадии

	/////////////////////////////////////////////////////////////////////////////////////////////


	/////////////////////////////////
	if isFirstActivation{// проснулся, получил InterruptMemory, первые мысли
		oldRandActionsIDarr =nil // выполненные случайно действияID, чтобы не повторяться
		/* уже есть конец активной цепочки fromNextID, полученный выше при if fromNextID==0{
		 */
		infoFunc8() //infoFunc8() -> getMentalPurpose()

		// перезапуск осмысления после просыпания, но можно и не перезапускать, а ждать Стимула
		if EvolushnStage > 4 {// инициативность
			return reloadConsciousness(fromNextID)
		}
	}else{// при кажом вызове кроме первого
		/* определить текущие объекты восприятия и выделить один из них - самые важные НЕГАТИВНЫЕ
		по всем категориям importanceType
		улучшение которого становится текущей целью.
		*/
		if activationType==1 {// объективное восприятие
			// выделить наиболее нзачимое в восприятии
			curImportanceObjectArr=getGreatestImportance(curActiveActions)	//curImportanceObjectArr []extremImportance - здесь теперь - текущие цели внимания к наиболее важному
			//выбрать ОДИН, самый актуальный объект из curImportanceJbjectArr []extremImportance
			indexA:=getTopAttentionObject(curImportanceObjectArr)// - индекс объекта внимания
			if indexA>0{//выбран curImportanceObjectArr[indexA]
				extremImportanceObject = &curImportanceObjectArr[indexA]
			}
		}
		if activationType==2 {// субъективный цикл
			/* ЕЩЕ НЕТ ПОДДЕРЖКИ объектов собственных мыслей по аналогии с importance.go
			indexA:=getGreatestImportanceMental()// выделить наиболее нзачимое в мыслях
			if indexA>0 {
				extremImportanceMentalObject = &curMentalImportanceObjectArr[indexA]
			}*/
		}
	}
	//////////////////////////////////////
/* // тестирование Правил
	if !isFirstActivation { //это - не пробуждение
		rules := getSuitableRules() //тестирвоание
		if rules > 0 {
			rules = 0
		}
	}*/


	if fromNextID == 0 { // еще нет цикла осмысления, 1 и 2 уровни - не требуют осмысления

		//////////////////////////////  ПЕРВЫЙ УРОВЕНЬ  /////////////////////////////////////////////
		/* Период ожидания ответа LastRunAutomatizmPulsCount при поочередном Стимуле-Ответе есть всегда.
		   А здсь - поиск Ответа именно после каждого Стимула. Так что LastRunAutomatizmPulsCount в функции не учитываем.
		*/
		if true && !isFirstActivation { //это - не пробуждение
			////////////////////////////// 1 уровень ////////////////////
			// ПЕРВЫЙ УРОВЕНЬ, самый примитивный уровень:
			// есть ли штатный мот.автоматизм и нужно ли его менять или задумываться
				mentalInfoStruct.motorAtmzmID = 0 // сброс прежнего значения
				//НЕТ!!! atmtzm := GetBelief2AutomatizmListFromTreeId(detectedActiveLastNodID)
				/* учитывается именно тот автоматизм, что рвется на выполнение при активации дерева автоматизмов,
				   даже если он подобран "мягким алгоритмом" в getAutomatizmFromNodeID
				   При атасе он выполняется не раздумявая, иначе подвергается сомнению в infoFunc6()
				*/
				//atmtzm := AutomatizmFromIdArr[currentAutomatizmAfterTreeActivatedID] НЕ ВСЕГДА  ВЫЗЫВАЕТСЯ КОГДА ЕСТЬ currentAutomatizmAfterTreeActivatedID
				atmtzmID := getAutomatizmFromNodeID(detectedActiveLastNodID)
				atmtzm := AutomatizmFromIdArr[atmtzmID]
				if atmtzm != nil { // есть незаблокированный автоматизм
					if atmtzm.Usefulness >= 0{//незаблокированный
						// Если Период преступной инициативы, если важная ситуация, но нет опасности, то - ПОДВЕРГНУТЬ СОМНЕНИЮ автоматизм.
						if EvolushnStage == 4 || !CurrentInformationEnvironment.veryActualSituation || CurrentInformationEnvironment.danger {
							//ПОДВЕРГНУТЬ СОМНЕНИЮ автоматизм, если нет опасности (не нужно реагировать аффектно) и ситуация важна
							infoFunc6()
							//mentalInfoStruct.notOldAutomatizm true - НЕ позволить запустить рвущийся на выполнение старый автоматизм
							if !mentalInfoStruct.notOldAutomatizm {
								isNormalRunAutomatimsInPreActive=true
								//Мышление не прерывается если нормально выполняются штатные автоматизмы
								reloadConsciousness(fromNextID)// перезапуск осмысления

								//можно без опаски выполнять штатный автоматизм
								return false //Эпиз.память не пишется. При опасности - состояние аффекта.
							} // если нет - далее искать альтернативу
						}
						mentalInfoStruct.motorAtmzmID = atmtzm.ID // для последующего использования с инфо-фукнциях
						// нужно ПОДВЕРГНУТЬ СОМНЕНИЮ автоматизм
					}else{// если автоматизмв заблокирован
						/*удалить авторитарное Правило с таким действием.
						Иначе никогда не сраотает checkForUnbolokingAutomatizm, см. ниже об этом
						 */
						tryRemoveRules(atmtzm)
						/* НО при создании авторитарного правила (func fixNewRules(lastCommonDiffValue int)) определяется,
						есть ли автоматизм с действием оператора curActiveActionsID, и если у него atmtzm.Usefulness<0 -
						снять блокировку и сделать штатным (checkForUnbolokingAutomatizm(curActiveActionsID))
						*/
					}
				}
		}

		/*  так нельзя, пусть пишется в automatizm_result.go там где сформировалось Правило!
			// нет штатного автоматизма, ситуация осмысливается, эпиз.память пишется в момент объективного вызова:
			if activationType==1 {
				// новый кадр эпизодической памяти, тип - ОБЪЕКТИВНЫЙ
				// В эпиз.память пишется только если не вызвало автоматических (неосознанных) действий,
				//   а было привлечено осознанное внимание consciousness(2
				newEpisodeMemory(currentRulesID, 0) // запись ОБЪЕКТИВНОЙ эпизодической памяти saveEpisodicMenory()
		// МЕНТАЛЬНЫЙ кадр newEpisodeMemory(currentRulesID, 1) пишется в func afterWaitingPeriod(
			}
		*/
		/////////////////////////////////////////////////////////

		//return false //для тестирования

		//////////////////////////////// 2 уровень ////////////////////////////
		if true && !isFirstActivation { //это - не пробуждение false для тестирования
			// ВТОРОЙ УРОВЕНЬ - попытка использования примитивных Правил
				rules := getSuitableRules()// выбрать подходящее Правило
				if rules > 0 { // по правилу найти автоматизм и запустить его
					ta := TriggerAndActionArr[rules]
					purpose := getPurposeGenetic()
					// выбираем Ответное действие из Правила чтобы повторить его
					ai := ActionsImageArr[ta.Action]
					if ai != nil {
						purpose.actionID = ai
						atmzm := createAndRunAutomatizmFromPurpose(purpose)
						if atmzm != nil {
							isNormalRunAutomatimsInPreActive=true
							//Мышление не прерывается если нормально выполняются штатные автоматизмы
							reloadConsciousness(fromNextID)// перезапуск осмысления
							return true // заблокирвать более низкоуровневое
						}
					}
			}
			/////////////////////////////////////////////////////////
		} // конец блокирования для тестирования

		//	return false //для тестирования

		// НАЧАТЬ ЦИКЛ ОСМЫСЛЕНИЯ
		/* НАЙТИ или создать Базовое звено цепи fromNextID для данной активности деревьев
		и пройти цепочку до конца, чтобы продолжить цикл от него.
		*/
		fromNextID = createBasicLink()
		// перезапуск осмысления
		// !!! не нужен перезапуск  return reloadConsciousness(stopMentalWork, fromNextID)
		newMentCickle = true         // начало прохода ментального цикла
		functionsInThisCickles = nil //ID запускаемых в текущем цикле инфо-функций

	}//if fromNextID == 0 { // еще нет цикла осмысления, 1 и 2 уровни - не требуют осмысления
	/////////////////////////////////////////////////


	//////////////////////////////  ТРЕТИЙ УРОВЕНЬ  /////////////////////////////////////////////
	if !isFirstActivation { // это - не пробуждение

		//////// ТРЕТИЙ УРОВЕНЬ - попытка найти решение, используя всю текущую инфрмацию с учетом срочности.
		/* Если не было цикла осмысления, а проходилиь только уровни до 3-го,
		то и нет обработки, нет записи в Эпизод.память ментальныз кадров. Что определяет ощцщение субъективного времени.

		//if EvolushnStage > 3 { - уже обеспечено

		      Ментальные автоматизмы нужны только если нет мот.автоматизма или его нужно переделать.
		   Т.е. привязанный к ветке дерева понимания мент автоматизм должен срабатывать ЗДЕСЬ.
		   Он продолжается по NextID или даже ветвится в зависимости от ситуации.
		   А если его еще нет, то создать БАЗОВЫЙ.
		Работа третьего уровня.
		   Запуск ментального автоматизма сопровождается перезапуском consciousness() кроме запуска моторного автоматизма.

		У ветки UnderstandingNode всегда должен быть Базовый ментальный автоматизм, с которого начинается
		просмотр в функции consciousness() на ее Третьем уровне. От него может идти цепочка дочерних.

		Каждый запускаемый мент. автоматизм (кроме MentalActionsImages.activateMotorID) после отработки
		вызывает consciousness() прямо или косвенно. В течение одного пульса может быть множество
		перезапусков consciousness() с продолжением процесса мышления и добавления в цепь (.NextID) новых
		автоматизмов.

		Базовый автоматизм должен прикинуть, какой будет следующий – путем выбора MentalActionsImages,
		сделать его, запустить, а в следующем цикле consciousness() использовать инфу (и все окружение) для
		формирования моторного автоматизма (MentalActionsImages.activateMotorID)  и тогда запустить его с
		периодом ожидания.

		По результату записывается Правило. Если хорошо, то данная цепочка так и заканчивается запуском
		моторного автоматизма (MentalActionsImages.activateMotorID), если плохо – формируется следующая
		цепочка (.NextID) с выбором другого MentalActionsImages и т.д. Т.е. формирование следующего звена
		цепочки идет С УЧЕТОМ ОПЫТА (MentalActionsImages) ПРЕДЫДУЩИХ.

		Структура мент Правила MentalTriggerAndAction начинается или с мент.действия MentalActionsImage
		или с моторного ActionsImage, потом - Ответ MentalActionsImage и обычный Эффект.

		 После срабатывания инфо-функции (.activateInfoFunc) информация добавляется к текущему информационному окружения
		в виде глобальной структурц и задается значение глобальной переменной currentInfoStructId == ID инфо-функции),
		которые могут использоваться при запуске consciousness().
		*/
		///////////////////////////////////////////

		// детекция ленивого состояния
		if isIdleness() { // ЛЕНЬ
			// пофиг все, можно лениво обрабатывать накопившиеся структуры, реагирование - на уровне - до EvolushnStage < 4
			saveEvolushnStage = EvolushnStage
			EvolushnStage = 3 // нагло и просто :) - произвольный откат уровня осознания

			processingFreeState(stopMentalWork) // обработка структур в свободном состоянии может быть долгой -

			EvolushnStage = saveEvolushnStage // возвращаем уровень осмысления, иначе зависнет на этой стадии

			return false                      // пусть выполняется все менее высокоуровневое
		} //if isIdleness()
		/////////////////////////
		/////////////////////////  НЕТ ЛЕНИ

		if stopMentalWork {// прерывание цикла
			//lib.WritePultConsol("Прерывание осмысления")
			// запомнить текущую работу в момент ее прерывания, чтобы можно было вернуться к обдумыванию
			addInterruptMemory()
			fromNextID = 0
			return false
		}
		///////////////////////////////

		if isСonfusion {
			/*
			if timeOfLastStimul < 3 { // 3 пульса установлено для запуска периода ожидания
				lib.SentСonfusion("Beast не успел обдумать прошлый ответ, а уже есть новый.")
			}*/
			if newMentCickle { // начался новый цикл
				lib.SentСonfusion("Beast задумался... Beast не успел обдумать прошлый ответ, а уже есть новый.")
				newMentCickle = false
			}
		}
		///////////////////////////////////

		// ограничить число циклов для 4-й стадии (кстати, - средство против зацикливания)
		if EvolushnStage == 4 && len(saveFromNextIDcurretCicle) > limitCickleCountForEvolushnStage4 {
			fromNextID = 0 // перестать думать,  не вызывать больше consciousness на 4-й ступени
			return false
		}
		// ограничить число циклов для 5-й стадии (кстати, - средство против зацикливания)
		if EvolushnStage == 5 && len(saveFromNextIDcurretCicle) > limitCickleCountForEvolushnStage5 {
			fromNextID = 0 // перестать думать
			return false
		}

		// запустить базовую функцию infoFunc2() поиска по имеющейся информации и прекратить цикл, т.е. нет ментальных Правил на стадии 4.
		if !lib.ExistsValInArr(functionsInThisCickles, 2) {// еще не использовалась infoFunc2()
			// запустить infoFunc2():
			if mentalSimpleRelexSolution() { //есть какое-то решение при запуске infoFunc2(), запущен моторный автоматизм
				fromNextID = 0 // раз есть моторный автоматизм, не нужно больше думать (М,Б, ТОЛЬКО ДЛЯ 4-йСТАДИИ)
				return false
			}
		}

		// НАЙДЕНО fromNextID
		if fromNextID > 0 { // продолжение осмысления по цепочке goNext.ID == fromNextID

			// добавить в кратковреемнную память
			addShortTermMemory(fromNextID)

			// если нужно, учесть: ID последний запуск инфо-функции: switch currentInfoStructId{

			/*			// тестирование случая, когда создан мент.автоматизм, запускающий мот.автоматизм
						if true && PulsCount > 6 {
							mentalInfoStruct.ActionsImageID, _ = CreateNewlastActionsImageID(0, []int{111}, []int{124}, 4, 5,true)
							infoFunc7()
							goNextFromIDArr[fromNextID].AutomatizmID = mentalInfoStruct.mentalAtmzmID
						}
			*/

			// fromNextID определяет звено цепи (другого не может быть после перезапуска с fromNextID)
			mImgID := goNextFromIDArr[fromNextID].AutomatizmID // id ментального автоматизма
			if mImgID > 0 {
				// ментальный автоматизм звена
				matmzm := MentalAutomatizmsFromID[mImgID]
				if matmzm != nil {
					// если акция - моторный автоматизм или переактивация состояния, то запустить
					mentAct := MentalActionsImagesArr[matmzm.ActionsImageID]
					if mentAct.typeID == 5 {
						// запуск моторного автоматизма  или объектиный перезапуск (через переактивированные деревья)
						RunMentalAutomatizm(matmzm)
						//ментально запущен мот.автоматизм, но можно продолжать размышления

						/* НЕ ТРЕБУЕТСЯ выход из рекурсивного цикла в understanding_tree.go в блок if detectedActiveLastUnderstandingNodID>0{
						   Если моторный окажется успешным, то он будет записан штатным для ветки detectedActiveLastUnderstandingNodID.
						   В func afterWaitingPeriod( будут записаны ментальные Правила.
						*/
						if EvolushnStage == 4 {
							return true
						}
						// для 5-й ступени хотя уже ментально запущен мот.автоматизм, но можно продолжать размышления,
						//пусть продолжает думать, если есть Доминанта или еще корректировать Ответ.

					} // не моторный автоматизм и не объяктивный перезапуск, а ИНФО_ФУНКЦИЯ
					// запустить инфо-функцию
					RunMentalAutomatizm(matmzm)
					// обработать результат

					// получить следующее звено с учетом результата
					fromNextID = calcNexusFromNextID(fromNextID)
					if fromNextID > 0 {
						// перезапуск осмысления
						return reloadConsciousness(fromNextID)
					}
				}
			}
			//  не случилось ничего полезного и это конец цепочки - нужно доращивать: createNexusFromNextID(

			if mentalInfoStruct.motorAtmzmID > 0 { // нужно ПОДВЕРГНУТЬ СОМНЕНИЮ автоматизм
				/* сопоставить действия actImgID с другими из Правил и т.п. инфы и если компоненты найденных действий
				дают более желательный результат, то Создать новый образ действий из таких компонентов.
				Или же, если действия actImgID в данных условиях прогнозируют плохой эффект, то
				попытаться сгенерировать новое действие по имеющейся информации.
				*/
				fromNextID = createNexusFromNextID(fromNextID, 6) // создается мент.авт-м запуска infoFunc6()
				if fromNextID > 0 {
					// перезапуск осмысления с обновлением currrentFromNextID
					return reloadConsciousness(fromNextID)
				}
			}
			/////////////////////////////////

			// еще нет решения
			infoFunc1(fromNextID)// - находить другое действие, вернуть новый fromNextID

			///////////////////

		} // конец фазы осмысления по цепоске goNext.ID == fromNextID
		//if !isFirstActivation{
	}else{// начать думать при пробуждении
		// TODO
	}
//////////////////////////////////////





//////////////////////////////// 4 уровень ////////////////////////////
if fromNextID==0 {
if EvolushnStage > 4 {
// ЧЕТВЕРТЫЙ УРОВЕНЬ - доминанта нерешенной проблемы - только если нет срочности
// и тут - Стек отложенных дел.
if CurrentInformationEnvironment.veryActualSituation || CurrentInformationEnvironment.danger {
// TODO
} else { // нет решения, паника, откатиться на прежний уровень регирования
// TODO аварийное решение проблемы
return false // пусть выполняется все менее высокоуровневое
}

} //if EvolushnStage > 4
}
/////////////////////////////////////////////////////////







/////////////////////////////////////////////////////////



return false // не блокировать последующий код ориентировочного рефлекса.
}
////////////////////////////////////////////









