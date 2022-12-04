/* Информационные функции, вызываемые как действия ментального автоматизма по их ID функции.

Инфо-функции - разные методы получения инфы, систематизации, поиска и т.п.
с целью найти верное действие для моторного автоматизма, а если нет,
то создания нового ментального автоматизма для продолжения итеации поиска.

У инфо-функций не должно быть вхлжного аргумента, иначе невозможно будет их вызывать из runMentalFunctionID(id int)
Поэтому в инфо-функции могут вызываться вспомогательные функции с аргументами, полученными в инфофункци
которые вызываются только если есть нужная инфа, например, сохраненная в mentalInfoStruct

Результат работы инфо-функции записывается в mentalInfoStruct
и определяется общая переменная currentInfoStructId == ID инфо-функции
*/

package psychic

import (
	wordSensor "BOT/brain/words_sensor"
	"BOT/lib"
	"strconv"
)

///////////////////////////////////////////////////////////////////////////

/* Дополнительное инфо-окружение.
Общая для всех информационных функций структура (типа информацонного окружения)
для сохранения найденной информации.
*/
type mentalInfo struct {
	mImgID int // ID MentalActionsImages найденного целевого действия
	ActionsImageID int//ID ActionsImage действия (стимула или ответа)
	motorAtmzmID int // ID моторного автоматизма
	mentalAtmzmID int // ID моторного автоматизма
	mentalPurposeID int // ID ментальной цели, альтернативной текущей  PurposeImage
	notOldAutomatizm bool// true - НЕ позволить запустить рвущийся на выполнение старый автоматизм
	runInfoFuncID int // запуск инфо-функции
	epizodFrameIndex int // индекс кадра КРАТКОВРЕМЕННОЙ памяти (shortTermMemory)
	moodeID int
	emotonID int
}
var mentalInfoStruct mentalInfo

func clinerMentalInfo(){
	mentalInfoStruct.mImgID=0
	mentalInfoStruct.ActionsImageID=0
		mentalInfoStruct.motorAtmzmID=0
		mentalInfoStruct.mentalAtmzmID=0
		mentalInfoStruct.mentalPurposeID=0
		mentalInfoStruct.notOldAutomatizm=false
	// ВСТАВЛЯТЬ ДРУГИЕ ЧЛЕНЫ ПО МЕРЕ ПОЯВЛЕНИЯ !!!!
}

/* произвольно активированные параметры, определяются при замуске ментального автоматизма.
Держатся на время, пока не изменятся генетически определенные соотвествующие параметры или
	если активация была в данном пульсе
 */
var mentalMoodVolitionID=0// призвольно активированное настроение
var mentalMoodVolitionPulsCount=0// призвольно активированное настроение

var mentalEmotionVolitionID=0// призвольно активированная эмоция
var mentalEmotionVolitionPulsCount=0// призвольно активированная эмоция

var mentalPurposeImageID=0// призвольно активированная цель
var mentalPurposeImagePulsCount=0// призвольно активированная цель
/////////////////////////////////////


/* Общая переменная currentInfoStructId == ID инфо-функции),
которые могут использоваться при запуске consciousness()
для определения какой параметр из mentalInfoStruct использовать как выходная инфа функции.

Можно использовать switch currentInfoStructId{ для выявления поля структуры mentalInfoStruct
 */
var currentInfoStructId=0

///////////////////////////////////
/* Функция вызова пронумерованной функции

 */
func runMentalFunctionID(id int){
	switch id {
	// НЕ ВЫЗЫВАЕТСЯ ментально case 1: infoFunc1()//Подобрать MentalActionsImages для начального звена цепочки
	case 2: infoFunc2()//Подобрать MentalActionsImages для последующего звена цепочки
	case 3: infoFunc3()//найти подходящий мент.автоматизм по опыту ментальных Правил
	case 4: infoFunc4()//Анализ инфо стркутуры и др. информации по currentInfoStructId и выдача решения
	case 5: infoFunc5()//создать и запустить ментальный автоматизм по ID акции
	case 6: infoFunc6()//ПОДВЕРГНУТЬ СОМНЕНИЮ автоматизм, если нет опасности (не нужно реагировать аффектно) и ситуация важна
	case 7: infoFunc7()//создать и запустить ментальный автоматизм запуска моторного автоматизма по действию ActionsImageID
	case 8: infoFunc8()//Ментальное определение ближайшей Цели в данной ситуации
	case 9: infoFunc9()//найти способ улучшения значимости объекта внимания extremImportanceObject
	case 10: infoFunc10()//найти способ улучшения значимости субъекта внимания extremImportanceMentalObject
	case 11: infoFunc11()//Ментальное отзеркаливание
	case 12: infoFunc12()//Cинтез новой фразы из компонентов, имеющих плюсы в Правилах
	case 122: infoFunc122()//выбрать ID действия имеющего плюсы в Правилах
	case 13: infoFunc13()//Отзеркалить Стимул от оператора
	case 14: infoFunc14()//ментально переактиваровать дерево понимания с mentalInfoStruct.moodeID и mentalInfoStruct.emotonID
	case 15: infoFunc15()//Для условия дерева автоматизмов (NodeAID) в одиночных Правилах выбираем наилучшее
	case 16: infoFunc16()//Случайно выдать любую известную фразу или действие и затем infoFunc7()

	default: return
	}

//	setCurIfoFuncID(id)
}
//////////////////////////////////////////////////

func getMentalFunctionString(id int)string{
	switch id {
	case 1: return "Подобрать MentalActionsImages для базового звена цепочки"
	case 2: return "Подобрать MentalActionsImages для последующего звена цепочки"
	case 3: return "Найти подходящий мент.автоматизм по опыту ментальных Правил"
	case 4: return "Анализ инфо стркутуры и др. информации по currentInfoStructId и выдача решения"
	case 5: return "Создать и запустить ментальный автоматизм по ID акции"
	case 6: return "Подвергнуть сомнению автоматизм, если нет опасности (не нужно реагировать аффектно) и ситуация важна"
	case 7: return "Создать и запустить ментальный автоматизм запуска моторного автоматизма по действию ActionsImageID"
	case 8: return "Ментальное определение ближайшей Цели в данной ситуации"
	case 9: return "Найти способ улучшения значимости объекта внимания extremImportanceObject"
	case 10: return "Найти способ улучшения значимости субъекта внимания extremImportanceMentalObject"
	case 11: return "Ментальное отзеркаливание"
	case 12: return "Cинтез новой фразы из компонентов, имеющих плюсы в Правилах"
	case 122: return "Выбрать ID действия имеющего плюсы в Правилах"
	case 13: return "Отзеркалить Стимул от оператора"
	case 14: return "Ментально переактиваровать дерево понимания с mentalInfoStruct.moodeID и mentalInfoStruct.emotonID"
	case 15: return "Для условия дерева автоматизмов (NodeAID) в одиночных Правилах выбираем наилучшее"
	case 16: return "Случайно выдать любую известную фразу и затем infoFunc7()"
	}
	return "Нет функции с ID = "+strconv.Itoa(id)
}
//////////////////////////////////////////////////////////


//////////////////////////////////////////////////////////
/* далее идут ПРОНУМЕРОВАННЫЕ ИНФОРМАЦИОННЫЕ ФУНКЦИИ,
для которых в mental_automatizm_INFO_structs.go определяются ИНФОРМАЦИОННЫЕ ГЛОБАЛЬНЫЕ СТРУКТУРЫ - для
передачи в них полученной информации.
Так же для передачи информации в инфо-функции (если это нужно, например, что найти) применяюися входне структуры.
 */
//////////////////////////////////////////////////////////




/* 1. Подобрать MentalActionsImages для продолжения цикла осмысления.
находить другое действие, вернуть новый fromNextID
Вызывается только из consciousness
*/
func infoFunc1(fromNextID int){
	clinerMentalInfo()
	setCurIfoFuncID(1)
	iID:=0// ID MentalActionsImages

	/* СНАЧАЛА поиск в ментальных Правилах готового решения
	      Найдя нужную комбинацию Правил для другой эмоции для успешного решения цепочки
	   воссоздать тот базовый контекст при той же ситуации (эмоция в найденном фрагменте эпиз.памяти)
	   и отработать цепочку «правильно».
	*/
	actionID:= infoFindRightMentalRules()
	if actionID>0{
		mAtmzmID,_:=createMentalAutomatizmID(0,actionID, 1)
		if mAtmzmID>0 {
			/* ментально переактивировать дерево понимания с той эмоцией, что была в КРАТКОВРЕМЕННОЙ памяти из Правила
			При этом тут же будет совершено моторное действие уже в новых условиях дерева понимания.
				 Это - непонятный пока момент...
			*/
			if mentalInfoStruct.epizodFrameIndex > 0 {
				tm := termMemory[mentalInfoStruct.epizodFrameIndex]
				if tm.uTreeNodID > 0 {
					uNode := UnderstandingNodeFromID[tm.uTreeNodID]
					if uNode != nil {
						mentalInfoStruct.moodeID = uNode.Mood
						mentalInfoStruct.emotonID = uNode.EmotionID
						reactivateEmotionUnderstandingЕree()// с перезапуском функции consciousness
					}
				}
			}
		}
		//id, newID := createNewlastgoNextID(0, detectedActiveLastUnderstandingNodID,detectedActiveLastNodID, 0, mAtmzmID)
		// создать новое звено типа 5(запуск моторного действия)
		iID, _ = CreateNewlastMentalActionsImagesID(0, 5, mAtmzmID, true)
		//mentalInfoStruct.mImgID = iID // передача инфы в структуру
		infoCreateAndRunMentMotorAtmzmFromAction(iID)// запустить моторный автоматизм
		reloadConsciousness(0)// конец цикла осмысления
		return
	}// конец поиска в правилах
	///////////////////////////////////////////////////

	// найти следующую инфо-функцию для вызова
	var nextIF=0

	//случайная инфо-фунция
	if EvolushnStage == 4 {
		// случайный выбор инфо-функции (хотел сделать Пультовй редактор последовательности выбора инфо-функции)
		nextIF= infoFindRundomMentalFunction()
	}
	if EvolushnStage == 5 {

		// TODO есть есть проблема - создать Доминанту

	}


	if nextIF>0 { // найдена инфо-функция, ну и просто запустить ее !!!!!!
		/* создание ментального автоматизма для запуска инфо-функции
		Можно было бы просто запускать инфо-функции, но в мент.Правилах - действие мент.автоматизма!
		*/
		fromNextID = createNexusFromNextID(fromNextID, nextIF) // создается мент.авт-м запуска infoFunc2()
		if fromNextID > 0 {
			// перезапуск осмысления с обновлением currrentFromNextID - с последующим запуском мент-го автоматизма.
			reloadConsciousness(fromNextID)
			return
		}
	}

	reloadConsciousness(0)// конец цикла осмысления
}
/////////////////////////////////////////////////////////


/* №2 Здесь выбор той или иной функции делается в контексте имеющейся инфо-среды,
а если этого не удается, то выбирается случайно одна из инфо-функций.

Не обязательно - ответ на Стимул, а может быть собственная ИНИЦИАТИВА (начинается со случайных действий infoFindRundomMentalFunction()).
*/
func infoFunc2(){
	currentInfoStructId=2 // определение актуального поля mentalInfo
	lib.WritePultConsol("Запущена функция:<br> "+getPurposeDetaileString(2))
	mentalSimpleRelexSolution()
	currentInfoStructId=2 // определение актуального поля mentalInfo
}
func mentalSimpleRelexSolution()bool{
	clinerMentalInfo()
	setCurIfoFuncID(2)
	iID:=0// ID MentalActionsImages

//	infoFunc13();return true // тестивароние попугайского отзеркаливания


	// если нет готового правила, ТО ИДЕМ ДАЛЬШЕ

	/* ГЛАВНОЕ - ПРИВЛЕЧЕНИЕ ОСОЗНАННОГО ВНИМАНИЯ к наиболее значимому,
	 есть ли актуальный объект внимания с отрицательной значимостью.
	Информационная функция "Понимание объекта восприятия" - выборка данного образа восприятия в дереве с прослеживанием,
	что оно означало в разных условиях. т.к. образ включает в себя все составляющие объекта восприятия, то он - обобщение,
	а его понимание - Вид всех последствий в разных условиях.
	Для образа фразы type Verbal struct - составляющие отдельные слова, которые могут быть в разных образах фраз,
	так что можно сделать функцию Вида - выборки данного слова в разных условиях с последствиями.
	*/
	if extremImportanceObject!=nil{
		// найти способ улучшения значимости объекта extremImportanceObject - запуск моторного автоматизма
		if infoFindAttentionObjImprovement(){// сделано, но еще не тестировалось
			return true// больше не искать, уже создан мент.автоматизм объективного действия
		}
	}

	/*  текущий субъект внимания (наиболее нзачимое в мыслях)
	TODO ЕЩЕ НЕТ ПОДДЕРЖКИ объектов собственных мыслей по аналогии с importance.go
	var extremImportanceMentalObject extremImportance

	Здесь должно находиться решение (выбор действия) на основе использования информации
	о наиболее важном в мыслях, скорее всего связанное с Доминантами.
	*/
	if infoFindAttentionObjMentalImprovement(){

		return true// больше не искать, уже создан мент.автоматизм объективного действия
	}
	///////////////////////////////////////

	/*если для данного сочетания Стимул-Ответ есть только один вид эффекта,
	то это - уже Информация, и чем больше опыт (количество обобщенных правил), тем такая информация полезнее.
	infoFunc15()
	*/
	//getDominantEffect(triggerID int, actionID int)



	/*Если оператор нажал кнопку Учитель, это - стимул привлечения внимания для наблюдения за ним,
	за достигаемой им целью и как он ее достигает. Это - начиная с 4-й стадии развития.
	Для понимания отзеркаленных действий оператора цель, которую ставил тварь перед своими действиями PurposeGenetic -
	фксируется в узлах дерева понимания PurposeImage -> SituationImage .
	Кнопки действий теперь активируют смысл (контекст) ситуации SituationImage
	*/
	// TODO

	/*Посмотреть, есть ли другие ветки у данного узла и использовать их начальную инфу.
	В ходе цикла осмысления могут быть переключения на другие условия goNext.MotorBranchID деревьев (перезапуск осмысления)
	при том же контексте Цели (сначала PurposeImage -> SituationImage). Попробовать решения других веток.
	М.б. применить goNextFromUnderstandingNodeIDArr
	В общем - анализ состояния веток циклов осмысления для данной Цели, т.к. кнопки действий активируют смысл (контекст) ситуации SituationImage.
	В том числе, в контексте, если оператор нажал кнопку Учитель, это - стимул привлечения внимания для наблюдения за ним, з
	а достигаемой им целью и как он ее достигает.
	*/
	// TODO







	funcID:=0
	if EvolushnStage == 4 {
		if CurrentInformationEnvironment.veryActualSituation || CurrentInformationEnvironment.danger {
			funcID = 13 // спросить, как правильно ответить, цикл осмысления прекращается
		}else{
			// случайный выбор инфо-функции - поддержка продолжения цикла осмысления
			funcID = infoFindRundomMentalFunction()
		}
	}
	if EvolushnStage > 4 {
		if isIdleness() { // ЛЕНЬ думать
			funcID = 13 // спросить, как правильно ответить, цикл осмысления прекращается
		} else {
			// случайный выбор инфо-функции - поддержка продолжения цикла осмысления
			funcID = infoFindRundomMentalFunction()
		}
	}


	if funcID>0 {// есть всегда
		// создать новое звено типа 4(запуск инфо-функции runMentalFunctionID(funcID))
		// мент.автоматизм нужен для Мент.Правила MentalTriggerAndAction
		iID,_=CreateNewlastMentalActionsImagesID(0,4,funcID,true)
		mentalInfoStruct.mImgID = iID // передача инфы в структуру
		mentalInfoStruct.runInfoFuncID=funcID

		runMentalFunctionID(funcID)// запуск инфо-функции
		return true
	}


	//TODO другие способы нахождения нового звена (пока не реализовано) это
	if EvolushnStage > 4 {
		// доминанты нерешенной проблемы НЕ ЗДЕСЬ :)

	}
	////////////////////////////////////////////////////////



	return false
}
//////////////////////////////////////////////////////////




/* №3 найти подходящий мент.автоматизм по опыту ментальных Правил
*/
func infoFunc3() {
	lib.WritePultConsol("Запущена функция:<br> "+getPurposeDetaileString(3))
	clinerMentalInfo()
	infoFindRightMentalRules()
	// получили mentalInfoStruct.mImgID
	currentInfoStructId=3 // определение актуального поля mentalInfo
}
func infoFindRightMentalRules()(int){
	setCurIfoFuncID(3)
	mrulesID:=getSuitableMentalRules()
	if mrulesID > 0 { // по правилу найти автоматизм и запустить его
		mta := MentalTriggerAndActionArr[mrulesID]
		// выбираем Ответное действие из Правила
		if mta != nil {
// последний кадр эпиз.памяти
			epizodFrameIndex :=0
			if len(mta.ShortTermMemoryID)>0 {
				epizodFrameIndex = mta.ShortTermMemoryID[len(mta.ShortTermMemoryID)-1]
			}
			mentalInfoStruct.epizodFrameIndex=epizodFrameIndex  // индекс кадра КРАТКОВРЕМЕННОЙ памяти
			mentalInfoStruct.mImgID=mta.Action
			return mentalInfoStruct.mImgID
		}
	}
	mentalInfoStruct.mImgID=0
	return 0
}
//////////////////////////////////////////////////////////



/* нализ инфо стркутуры и др. информации по currentInfoStructId и выдача решения
Нужна таблица, какие инфо-функции вызывать при данной ситуации

 */
func infoFunc4(){
	lib.WritePultConsol("Запущена функция:<br> "+getPurposeDetaileString(4))
	clinerMentalInfo()
	setCurIfoFuncID(4)
	if currrentFromNextID>0 {
		analisAndSintez(currrentFromNextID) // возвращает mentalInfoStruct.mImgID
	}else{
		mentalInfoStruct.mImgID=0
	}
	// получили mentalInfoStruct.mImgID
	currentInfoStructId=4 // определение актуального поля mentalInfo
}
/* анализ инфо стркутуры и др. информации по currentInfoStructId и выдача решения
Нужна таблица, какие инфо-функции вызывать при данной ситуации
Возвращает fromNextID следующего звена, даже если найден моторный автоматизм или задана объективная переактивация.
*/
func analisAndSintez(fromNextID int)(int){

	// сначала стандартно:
	/*
		1. Поиск MentalActionsImages для следующего .NextID начинается по ментальным Правилам.
		2. Если нет правил, посмотреть, есть ли дургие ветки у данного узла и использовать их начальную инфу.
		3. Если нет дргих веток, выбрать какой-то провоцирующий.
	*/
	// поиск в Правилах
	//action:=infoFindRightMentalRules()

	// сложные типы данных, полученные в инфо-функциях
	switch currentInfoStructId{

	}




	return mentalInfoStruct.mImgID // очищенный или заполненный
}
//////////////////////////////////////////////////////////


/* №5 создать и запустить ментальный автоматизм по действию mImgID -
ВСЕГДА ПОСЛЕ ПОЛУЧЕНИЯ ОБРАЗА ДЕЙСТВИЯ mentalInfoStruct.mImgID
 */
func infoFunc5() {
	if mentalInfoStruct.mImgID >0 {
		lib.WritePultConsol("Запущена функция:<br> "+getPurposeDetaileString(5))
		clinerMentalInfo()
		infoCreateAndRunMentAtmzmFromAction(mentalInfoStruct.mImgID)
	}
	// получили mentalInfoStruct.mImgID
	currentInfoStructId=5 // определение актуального поля mentalInfo
}
func infoCreateAndRunMentAtmzmFromAction(actImgID int){
	if actImgID ==0 {
		return
	}
	setCurIfoFuncID(5)
	id, matmzm := createMentalAutomatizmID(0, actImgID, 1)
	if id >0 {
		// запустить мент.автоматизм
		RunMentalAutomatizm(matmzm)
	}
}
//////////////////////////////////////////////////////////


/* №6 нужно ПОДВЕРГНУТЬ СОМНЕНИЮ автоматизм, если нет опасности (не нужно реагировать аффектно) и ситуация важна.
Создать новый образ действий, если он лучше текущего в моторном автоматизме, рвущесмя на выполнение
и запустить ментальный автоматизм по акции -
ВСЕГДА ПОСЛЕ ПОЛУЧЕНИЯ ОБРАЗА ДЕЙСТВИЯ mentalInfoStruct.mImgID
*/
func infoFunc6() {
	if mentalInfoStruct.motorAtmzmID >0 {
		lib.WritePultConsol("Запущена функция:<br> "+getPurposeDetaileString(6))
		clinerMentalInfo()
		infoCreateAndRunNewActionMentAtmzmFromAction(mentalInfoStruct.motorAtmzmID)
	}
	// получили mentalInfoStruct.mImgID
	currentInfoStructId=6 // определение актуального поля mentalInfo
}
func infoCreateAndRunNewActionMentAtmzmFromAction(actImgID int){
	if actImgID ==0 {// ID действия мот.автоматизма, рвущегося на выполнение
		return
	}
	setCurIfoFuncID(6)
/* учитывать ли инфу, что автоматизм - общий (не привязанный к конечному узлу ветки)?
var isCommon=false
if automatizm.BranchID>1000000{// это общий, привязанный не к ветке, а к действию или фразе
	isCommon=true
}
 */

// сопоставление, грозит ли опасностной значимостью запускаемый автоматизм
todoComparison(actImgID)
/*
	id, matmzm := createMentalAutomatizmID(0, actImgID, 1)
	if id >0 {
		// запустить мент.автоматизм
		RunMentalAutomatizm(matmzm)
	}
 */
}
/* сопоставление, грозит ли опасностной значимостью запускаемый автоматизм
*/
func todoComparison(actImgID int){
	// учесть значимости компонентов образа действия автоматизма
	actImg:=ActionsImageArr[actImgID]
	if actImg==nil{
		mentalInfoStruct.notOldAutomatizm=true
		return
	}
	objImportance:=getGreatestImportance(actImg)
	// оценить, насколько приемелемо запускать акции с такими значимостями
	for i := 0; i < len(objImportance); i++ {
		if objImportance[i].extremVal < 0{ // отрицательная значимость - не позволять запуск автоматизма
			mentalInfoStruct.notOldAutomatizm=true
			return
		}
		}

mentalInfoStruct.notOldAutomatizm=false
	return
}
//////////////////////////////////////////////////////////



/* №7 создать и запустить ментальный автоматизм запуска моторного автоматизма по действию ActionsImageID -
ВСЕГДА ПОСЛЕ ПОЛУЧЕНИЯ ОБРАЗА ДЕЙСТВИЯ mentalInfoStruct.ActionsImageID
Создается моторный автоматизм (если такого еще нет), привязанный к ветке текущей активности дерева автоматизмов.
*/
var prevMotorAtmzmID=0
func infoFunc7() {
	if mentalInfoStruct.ActionsImageID >0 {
		lib.WritePultConsol("Запущена функция:<br> "+getPurposeDetaileString(7))
		infoCreateAndRunMentMotorAtmzmFromAction(mentalInfoStruct.ActionsImageID)
	}else{
		lib.WritePultConsol("Для функции infoFunc7() НЕ ОПРЕДЕЛНО mentalInfoStruct.ActionsImageID!")
		return
	}
	// получили mentalInfoStruct.mImgID
	currentInfoStructId=7 // определение актуального поля mentalInfo
}
func infoCreateAndRunMentMotorAtmzmFromAction(ActionsImageID int){
	if ActionsImageID ==0 {
		return
	}
	setCurIfoFuncID(7)
	motorID,motorAtmzm:=createNewAutomatizmID(0,detectedActiveLastNodID,ActionsImageID,true)
	if motorID==0{
		return
	}
	if prevMotorAtmzmID==motorID{// раньше был запущен ментально такой мот.автоматизм
		// применить мозжечковый рефлекс
		cerebellumCoordination(motorAtmzm, 1)// 1 - усилить действие
	}
	prevMotorAtmzmID=motorID
	clinerMentalInfo()
	mentalInfoStruct.motorAtmzmID=motorID
if motorID==0{
	return
}
	actImgID,_:=CreateNewlastMentalActionsImagesID(0,5,motorID,true)
if actImgID==0{
	return
}
	id, matmzm := createMentalAutomatizmID(0, actImgID, 1)
	if id >0 {
		mentalInfoStruct.mentalAtmzmID=id
		// запустить мент.автоматизм
		RunMentalAutomatizm(matmzm)
	}
}
//////////////////////////////////////////////////////////



// найти способ улучшения значимости объекта extremImportanceObject
func infoFunc9() {
	lib.WritePultConsol("Запущена функция:<br> "+getPurposeDetaileString(9))
	clinerMentalInfo()
	infoFindAttentionObjImprovement()
	currentInfoStructId=9 // определение актуального поля mentalInfo
}
// улучшение значимости объекта внимания
func infoFindAttentionObjImprovement()bool {
	setCurIfoFuncID(9)
	if extremImportanceObject == nil {
		return false
	}
	// найти в Правилах ответное действие с объектом extremImportanceObject, приводящее к успеху
	rulesID:=getRulesArrFromAttentionObject(extremImportanceObject.objID, extremImportanceObject.kind)
	if rulesID>0 {// достаточная уверенность
		// действие из Правила
		rules:=rulesArr[rulesID]
		if rules==nil{return false}
		actID:=TriggerAndActionArr[rules.TAid[0]].Action
		// создание мент.авто-ма запуска действия, улучшающего значимость объекта внимания - return true
		infoCreateAndRunMentMotorAtmzmFromAction(actID)
		return true
	}
return false
}
////////////////////////////////////////////////////////////////////////
// найти способ улучшения значимости субъекта внимания extremImportanceMentalObject
func infoFunc10() {
	lib.WritePultConsol("Запущена функция:<br> "+getPurposeDetaileString(10))
	clinerMentalInfo()
	infoFindAttentionObjMentalImprovement()
	currentInfoStructId=10 // определение актуального поля mentalInfo
}
// улучшение объекта внимания
func infoFindAttentionObjMentalImprovement()bool {
	setCurIfoFuncID(10)
	if extremImportanceMentalObject == nil {
		return false
	}

	//ЕЩЕ НЕТ ПОДДЕРЖКИ объектов собственных мыслей по аналогии с importance.go
	// TODO по аналогии rulesID:=getRulesArrFromAttentionObject( - найти в Правилах ответное действие с объектом extremImportanceMentalObject, приводящее к успеху

	mentalInfoStruct.motorAtmzmID=0
	return false
}
////////////////////////////////////////////////////////////////////////



////////////////////////////////////////////////////////////////////////
// Ментальное отзеркаливание
func infoFunc11() {
	lib.WritePultConsol("Запущена функция:<br> "+getPurposeDetaileString(11))
	clinerMentalInfo()
	infoMentalMirriring()
	currentInfoStructId=11 // определение актуального поля mentalInfo
}
func infoMentalMirriring()bool {
	setCurIfoFuncID(11)
	//есть ли фраза в действиях оператора
	if curActiveActions.PhraseID==nil{
		return false
	}
	/* алгоритм:
	1. Найти такую фразу в Ответах Beast, в Правилах: rulesID
	2. Посмотреть какое последовало действие оператора на это - в эпиз.памяти после rulesID: answer
	3. Создать автоматизм на такое действие и выдать на Пульт.
	 */

	// Пункт 1: ищем в одиночных Правилах
	for _, v := range rulesArr {
		if len(v.TAid)==1{
			ta:=TriggerAndActionArr[v.TAid[0]]
			if ta==nil{
				lib.WritePultConsol("Нет карты TriggerAndActionArr для iD="+strconv.Itoa(v.TAid[0]))
				return false
			}
			// условия должны совпадать, чтобы фраза не была сказана невпопад
			if v.NodeAID!=detectedActiveLastNodID { // пока v.NodePID не учитываем, иначе фиг найдет...
				continue
			}
			act:=ActionsImageArr[ta.Action]
			if act==nil{
				lib.WritePultConsol("Нет карты ActionsImageArr для iD="+strconv.Itoa(ta.Action))
				return false
			}
			if lib.EqualArrs(act.PhraseID, curActiveActions.PhraseID){// нашли
				//rulesID=k // правило, где в Ответе применена такая же фраза
				// Пункт 2: ищем в эпизод памяти такое Правило, начиная с конца
				end:=10000 // не смотрим далее, чем на end кадров
				for i := len(EpisodeMemoryObjects)-1; i >= 0 || end==0; i -- {
					em:=EpisodeMemoryObjects[i]
					if em.TriggerAndActionID == ta.ID{//нашли кадр
						// вынуть действие из следующего кадра
						nextEM:=EpisodeMemoryObjects[i+1]
						nextTa:=TriggerAndActionArr[nextEM.TriggerAndActionID]
						if nextTa==nil{
							lib.WritePultConsol("Нет карты TriggerAndActionArr для iD="+strconv.Itoa(nextEM.TriggerAndActionID))
							return false
						}
						acting:=ActionsImageArr[nextTa.Trigger]
						if acting==nil{
							lib.WritePultConsol("Нет карты ActionsImageArr для iD="+strconv.Itoa(nextTa.Trigger))
							return false
						}
						// Пункт 3: здесь определяется mentalInfoStruct.motorAtmzmID
						infoCreateAndRunMentMotorAtmzmFromAction(acting.ID)
return true
					}
					end--
				}
			}
		}
	}

	mentalInfoStruct.motorAtmzmID=0
	return false
}
////////////////////////////////////////////////////////////////////////




/*Cинтез новой фразы из компонентов, имеющих плюсы в Правилах
Возможен при хорошем словарном запасе и хорошем опыте значимости (importance) фраз.

Пока что выдается первая подходящая фраза, имеющая +значимость в данных услових
 */
func infoFunc12() {
	lib.WritePultConsol("Запущена функция:<br> "+getPurposeDetaileString(12))
	clinerMentalInfo()
	infoSynthesisOwnPrase()
	currentInfoStructId=12 // определение актуального поля mentalInfo
}
func infoSynthesisOwnPrase()bool {
	setCurIfoFuncID(12)
/*TODO сделать фразу, состояющую из 2-3-х известных фраз, найденных в Правилах при данных условиях и выдать ее на Пульт
Чтобы фраза не была бессмысленным попугайством, нужно проверять ее смысл по importanceFromID
	importance.Type=5//фраза
  importance.ObjectID=praseID
  importance.Value>0
	для текущих условий
  importance.NodeAID
  importance.NodePID
 */
	var phraseArr []int // здесь собираем плюсовые фразы
	for _, v := range rulesArr {
		if len(v.TAid) == 1 {
			ta := TriggerAndActionArr[v.TAid[0]]
			if ta == nil {
				lib.WritePultConsol("Нет карты TriggerAndActionArr для iD=" + strconv.Itoa(v.TAid[0]))
				return false
			}
			if ta.Effect<=0{
				continue
			}
			// условия должны совпадать, чтобы фраза не была сказана невпопад
			if v.NodeAID != detectedActiveLastNodID { // пока v.NodePID не учитываем, иначе фиг найдет...
				continue
			}
			act := ActionsImageArr[ta.Action]// берем именно Ответ Beast
			if act == nil {
				lib.WritePultConsol("Нет карты ActionsImageArr для iD=" + strconv.Itoa(ta.Action))
				return false
			}
			if act.PhraseID!=nil{
				// посмотртеь значимость фразы, только для act.PhraseID[0]
				imp:=getObjectImportance(5,act.PhraseID[0], detectedActiveLastNodID,detectedActiveLastUnderstandingNodID)
				if imp!=nil && imp.Value>0{
					phraseArr=append(phraseArr,act.PhraseID[0])
// TODO более крутой способ подбора фраз
					if len(phraseArr)>1{// пока не смотреть больше. т.к. будет выбрано просто 2 первые.
break
					}
				}
			}
		}
	}
// суммарная фраза - на выход
	if len(phraseArr)>0{// выдать нормальным тоном с настроением wordSensor.CurPultMood
		actID,_:=CreateNewlastActionsImageID(0,nil,phraseArr,90,wordSensor.CurPultMood,true)
if actID>0{
	// здесь определяется mentalInfoStruct.motorAtmzmID
	infoCreateAndRunMentMotorAtmzmFromAction(actID)
	return true
}
	}
	return false
}
/////////////////////////////////////////////////////////////////////////

/*выбрать ID действия имеющего плюсы в Правилах
*/
func infoFunc122() {
	lib.WritePultConsol("Запущена функция:<br> "+getPurposeDetaileString(122))
	clinerMentalInfo()
	infoSynthesisOenAction()
	currentInfoStructId=12 // определение актуального поля mentalInfo
}
func infoSynthesisOenAction()bool {
	setCurIfoFuncID(122)
	for _, v := range rulesArr {
		if len(v.TAid) == 1 {
			ta := TriggerAndActionArr[v.TAid[0]]
			if ta == nil {
				lib.WritePultConsol("Нет карты TriggerAndActionArr для iD=" + strconv.Itoa(v.TAid[0]))
				return false
			}
			if ta.Effect<=0{
				continue
			}
			// условия должны совпадать, чтобы фраза не была сказана невпопад
			if v.NodeAID != detectedActiveLastNodID { // пока v.NodePID не учитываем, иначе фиг найдет...
				continue
			}
			act := ActionsImageArr[ta.Action]// берем именно Ответ Beast
			if act == nil {
				lib.WritePultConsol("Нет карты ActionsImageArr для iD=" + strconv.Itoa(ta.Action))
				return false
			}
			if act.ActID!=nil{
				aID,_:=CreateNewlastActionsImageID(0,act.ActID,nil,90,0,true)
				// посмотртеь значимость действия
				imp:=getObjectImportance(1,aID, detectedActiveLastNodID,detectedActiveLastUnderstandingNodID)
				if imp!=nil && imp.Value>0{
					infoCreateAndRunMentMotorAtmzmFromAction(aID)
					return true
				}
			}
		}
	}
	return false
}
/////////////////////////////////////////////////////////////////////////




/* Тупое повторение Стимула оператора. Попугайство.
Показать непонимания, растерянность с вопросом о том, как нужно реагировать:
"Ответь сам на "+спопугайничать оператора+" чтобы показать, как лучше ответить."

Отзеркалить последний Стимул от оператора и совершить такое же действие.
Нужно учесть, что отзеркаливание происходит и в orientation_reflexes.go func orientation_1()!!! м.б. стоит оттуда это убрать?
Отзеркаливает авторитерный ответ Оператора на совершенное действие с помощью func fixNewTeachRules()
- запись авторитарного Правила.
*/
func infoFunc13() {
	lib.WritePultConsol("Запущена функция:<br> "+getPurposeDetaileString(13))
	clinerMentalInfo()
	infoMirroringStimul()
	currentInfoStructId=13 // определение актуального поля mentalInfo
}
func infoMirroringStimul()bool {
	setCurIfoFuncID(13)
	// свежесть ответа оператора - не позже, чем 5 пульсов назад
	if curActiveActionsID>0 && (curActiveActionsPulsCount > PulsCount-5){
		// выполнить действие curActiveActionsID
		isTeachQuestion=true//Показать непонимание, растерянность с предложением научить
		infoCreateAndRunMentMotorAtmzmFromAction(curActiveActionsID)
		return true
	}
	return false
}
/////////////////////////////////////////////////////////////////////



/* Ментально переактиваровать дерево понимания с mentalInfoStruct.moodeID и mentalInfoStruct.emotonID
Хотя можно переактивировать четез действия ментального автоматизма 1 или 2 (или 3 для Цели), но только что-то одно.
*/
func infoFunc14() {
	lib.WritePultConsol("Запущена функция:<br> "+getPurposeDetaileString(14))
	// clinerMentalInfo() чтобы можно было задавать mentalInfoStruct.moodeID и mentalInfoStruct.emotonID
	reactivateEmotionUnderstandingЕree()
	currentInfoStructId=14 // определение актуального поля mentalInfo
}
func reactivateEmotionUnderstandingЕree()bool {
	setCurIfoFuncID(14)
if mentalInfoStruct.moodeID==0{// выбрать случайно - для infoFindRundomMentalFunction()
	var infoArr=[]int{1,2,3}
	mentalMoodVolitionID=lib.RandChooseIntArr(infoArr)
}else{
	mentalMoodVolitionID=mentalInfoStruct.moodeID
}
if mentalInfoStruct.emotonID==0{// выбрать случайно - для infoFindRundomMentalFunction()
	var infoArr []int
	for k,_ := range EmotionFromIdArr {
		infoArr=append(infoArr,k)
	}
	mentalEmotionVolitionID=lib.RandChooseIntArr(infoArr)
}else{
	mentalEmotionVolitionID=mentalInfoStruct.emotonID
}
	mentalMoodVolitionPulsCount=PulsCount
	mentalEmotionVolitionPulsCount=PulsCount
	understandingSituation(2)
	clinerMentalInfo()
	return false
}
/////////////////////////////////////////////////////////////////////



/* для условия дерева автоматизмов (NodeAID) в одиночных Правилах выбираем наилучшее
Eсли для данного сочетания Стимул-Ответ есть только один вид эффекта,
	то это - уже Информация, и чем больше опыт (количество обобщенных правил), тем такая информация полезнее.
*/
func infoFunc15() {
	lib.WritePultConsol("Запущена функция:<br> "+getPurposeDetaileString(15))
	clinerMentalInfo()
	beastIDRulesFromCondA()
	currentInfoStructId=15 // определение актуального поля mentalInfo
}
func beastIDRulesFromCondA()bool {
	setCurIfoFuncID(15)
	// для условия дерева автоматизмов в одиночных Правилах выбираем наилучшее
	rulesID:=getBeastIDRulesFromCondA(detectedActiveLastNodID)
	if rulesID>0 { // не найдено для точного совпадения условий
		rules:=rulesArr[rulesID]
		if rules==nil{
			return false
		}
		actID:=TriggerAndActionArr[rules.TAid[0]].Action
		if actID>0{
			// здесь определяется mentalInfoStruct.motorAtmzmID
			infoCreateAndRunMentMotorAtmzmFromAction(actID)
			return true
		}
		return true
	}
return false
}
///////////////////////////////////////////////////


/* По ЗНАЧИМОСТИ или Случайно выдать действие и затем infoFunc7()
и сразу запустить (Писать мент.Правило на случайным действиям бесполезно)
*/
func infoFunc16() {
	lib.WritePultConsol("Запущена функция:<br> "+getPurposeDetaileString(15))
	clinerMentalInfo()
	randomAction()
	currentInfoStructId=16 // определение актуального поля mentalInfo
}
var oldRandActionsIDarr []int // выполненные действияID, чтобы не повторяться. Сбрасывается при просыпании if isFirstActivation{
var wasRandPhrase=false // последний раз была не фраза
func randomAction()bool {
	setCurIfoFuncID(16)
	var rActID=0

	//значимый образ ActionsImage
	praseI:=getObjectImportance(1,curActiveActions.PhraseID[0], detectedActiveLastNodID, detectedActiveLastUnderstandingNodID)
	if praseI != nil{//
		rActID=ActionsImageArr[praseI.ObjectID].ID
	}
	if rActID==0{// ID несловестного действия ActionsImage.ActID[n]
		praseI:=getObjectImportance(3,curActiveActions.PhraseID[0], detectedActiveLastNodID, detectedActiveLastUnderstandingNodID)
		if praseI != nil{//
			rActID,_=CreateNewlastActionsImageID(0,[]int{praseI.ObjectID},nil,90,0,true)
		}
	}
	if rActID==0{// ID отдельной фразы Verbal.PhraseID[n]
		praseI:=getObjectImportance(4,curActiveActions.PhraseID[0], detectedActiveLastNodID, detectedActiveLastUnderstandingNodID)
		if praseI != nil{//
			rActID,_=CreateNewlastActionsImageID(0,nil,[]int{praseI.ObjectID},90,0,true)
		}
	}
/*
	// попробовать найти по значимостям - infoFunc9()
	if infoFindAttentionObjImprovement(){
		return true
	}

	if !wasRandPhrase { // выбрать случайную фразу
		//infoFunc12()//Cинтез новой фразы из компонентов, имеющих плюсы в Правилах
		if infoSynthesisOwnPrase() {// уже совершено найденное действие
			wasRandPhrase = true
			return true
		}
	}

		// infoFunc122():
if infoSynthesisOenAction() {// уже совершено найденное действие
	wasRandPhrase = false
	return true
}*/
	if rActID==0 {
		// выдать случайное действие из уже известных ActionsImageArr
		var actualArr []int
		for k, _ := range ActionsImageArr {
			actualArr = append(actualArr, k)
		}
		rActID = lib.RandChooseIntArr(actualArr)
	}

	if rActID>0{
		// не повторяться
		if !lib.ExistsValInArr(oldRandActionsIDarr, rActID){
			return false
		}

			oldRandActionsIDarr=append(oldRandActionsIDarr,rActID)
			// далее - infoFunc7(): создать и запустить ментальный автоматизм запуска моторного автоматизма по действию
			infoCreateAndRunMentMotorAtmzmFromAction(rActID)
			return true
		}

	return false
}
///////////////////////////////////////////////////




//////////////////////////////////////////////////////////
/* случайный выбор ментальной функции, из тех, что еще не использовались в данном цикле (нет в functionsInThisCickles)
Непонятно, как отслеживать Эффект если нет памяти о том, какой именно случайный парамет был применен.
*/
func infoFindRundomMentalFunction()int{

	// return 13 // тестирование

	clinerMentalInfo()// чтобы было случайное clinerMentalInfo() для 14-й функции
	var infoArr=[]int{3,4,8,9,10,12,13,14,15,16}
	var actualArr []int
	for i := 0; i < len(infoArr); i++ {
		if !lib.ExistsValInArr(functionsInThisCickles, infoArr[i]){
			actualArr=append(actualArr,infoArr[i])
		}
	}
	if len(actualArr)>0{
		return lib.RandChooseIntArr(actualArr)
	}
	return 0
}
////////////////////////////////////////////////



/*
*************** Группировка зеркальных автоматизмов *********************
Были зафиксированы две цепочки диалога по 2 шага каждая (1 шаг - 1 пара вопрос + ответ):
а) Привет – привет, как дела – нормально
б) Приветствую – привет, Ты как – отлично
Из них сформировались 2 зеркальных автоматизма: привет - как дела, привет - ты как. Их можно сгруппировать и вывести в отдельном
массиве варианты ответов на один пускатель: привет - как дела, ты как.
При поиске ответа нужно искать в этом массиве и выбирать варианты например по счетчику успешности. Это будет намного быстрее,
чем перебирать весь массив автоматизмов.
Для такой группировки нужно при создании нового зеркального автоматизма дописывать
в этот массив новый вариант в нужной строке: находить в массиве пусковое слово и добавлять к нему вариант ответа.
По сути это групповой запрос, только выделенный в динамическую таблицу. Так БД-шники часто делают, если приходится ворочать объемные
данные под миллионы записей. Вместо тяжелых тормозных запросов строятся буферные таблицы и забиваются через хранимки при совершении
операций с данными.
*/