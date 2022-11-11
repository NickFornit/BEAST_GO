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
}
var mentalInfoStruct mentalInfo

func clinerMentalInfo(){ // НУЖНА ЛИ?..
	mentalInfoStruct.mImgID=0
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
	case 1: infoFunc1()//Подобрать MentalActionsImages для базового звена цепочки
	case 2: infoFunc2()//Подобрать MentalActionsImages для последующего звена цепочки
	case 3: infoFunc3()//айти подходящий мент.автоматизм по опыту ментальных Правил
	case 4: infoFunc4()//нализ инфо стркутуры и др. информации по currentInfoStructId и выдача решения
	case 5: infoFunc5()//создать и запустить ментальный автоматизм по акции
	case 6: infoFunc6()//ПОДВЕРГНУТЬ СОМНЕНИЮ автоматизм, если нет опасности (не нужно реагировать аффектно) и ситуация важна
	case 7: infoFunc7()//создать и запустить ментальный автоматизм запуска моторного автоматизма по действию ActionsImageID
	}
}

func getMentalFunctionString(id int)string{
	switch id {
	case 1: return "Подобрать MentalActionsImages для базового звена цепочки"
	case 2: return "Подобрать MentalActionsImages для последующего звена цепочки"
	case 3: return "Найти подходящий мент.автоматизм по опыту ментальных Правил"
	case 4: return "Анализ инфо стркутуры и др. информации по currentInfoStructId и выдача решения"
	case 5: return "Создать и запустить ментальный автоматизм по акции"
	case 6: return "Подвергнуть сомнению автоматизм, если нет опасности (не нужно реагировать аффектно) и ситуация важна"
	case 7: return "Создать и запустить ментальный автоматизм запуска моторного автоматизма по действию ActionsImageID"
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




/* №0 Подобрать MentalActionsImages для базового звена цепочки
c вызовом activateInfoFunc для начальной информированности,
случайно или по заготовке редактора с Пульта
*/
func infoFunc1(){
	iID:=0// ID MentalActionsImages

	typeID:=0
	valID:=0

	// TODO подобрать
	/*
	1. Поиск MentalActionsImages для следующего .NextID начинается по ментальным Правилам.
	2. Если нет правил, посмотреть, есть ли дургие ветки у данного узла и использовать их начальную инфу.
	3. Если нет дргих веток, выбрать какой-то провоцирующий.
	 */
	// поиск в Правилах
	//action:=infoFindRightMentalRukes()

	// создать
	iID,_=CreateNewlastMentalActionsImagesID(0,typeID,valID)

	mentalInfoStruct.mImgID = iID // передача инфы в структуру
	currentInfoStructId=1 // определение актуального поля mentalInfo
}
/////////////////////////////////////////////////////////


/* №2 Подобрать MentalActionsImages для последующего звена цепочки, если не найдено в опыте
с вызовом различных activateInfoFunc или
 c вызовом activateMotorID моторнорнного автоматизма (а значит, с запуском моторного с периодом ожидания),
т.е. раз нет решения, пробовать подобрать моторные действия, просмотрев цеаочку с начала.
Если в базовом звене есть другие MotorBranchID то - подсмотреть как они продолжаются с насовпажающего звена цепи!
*/
func infoFunc2(){
	iID:=0// ID MentalActionsImages

	typeID:=0
	valID:=0

	// TODO подобрать
	/*
		1. Поиск MentalActionsImages для следующего .NextID начинается по ментальным Правилам.
		2. Если нет правил, посмотреть, есть ли дургие ветки у данного узла и использовать их начальную инфу.
		3. Если нет дргих веток, выбрать какой-то провоцирующий.
	*/
	// поиск в Правилах
	//action:=infoFindRightMentalRukes()


	// создать
	iID,_=CreateNewlastMentalActionsImagesID(0,typeID,valID)

	mentalInfoStruct.mImgID = iID // передача инфы в структуру
	currentInfoStructId=2 // определение актуального поля mentalInfo
}
//////////////////////////////////////////////////////////





/* №3 найти подходящий мент.автоматизм по опыту ментальных Правил
*/
func infoFunc3() {
	infoFindRightMentalRukes()
	// получили mentalInfoStruct.mImgID
	currentInfoStructId=3 // определение актуального поля mentalInfo
}
func infoFindRightMentalRukes()(int){
	mrules:=getSuitableMentalRules()
	if mrules > 0 { // по правилу найти автоматизм и запустить его
		mta := MentalTriggerAndActionArr[mrules]
		// выбираем Ответное действие из Правила
		if mta != nil {
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
	clinerMentalInfo()
	// сначала стандартно:
	/*
		1. Поиск MentalActionsImages для следующего .NextID начинается по ментальным Правилам.
		2. Если нет правил, посмотреть, есть ли дургие ветки у данного узла и использовать их начальную инфу.
		3. Если нет дргих веток, выбрать какой-то провоцирующий.
	*/
	// поиск в Правилах
	//action:=infoFindRightMentalRukes()

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
		infoCreateAndRunMentAtmzmFromAction(mentalInfoStruct.mImgID)
	}
	// получили mentalInfoStruct.mImgID
	currentInfoStructId=5 // определение актуального поля mentalInfo
}
func infoCreateAndRunMentAtmzmFromAction(actImgID int){
	if actImgID ==0 {
		return
	}
	id, matmzm := createMentalAutomatizmID(0, actImgID, 1)
	if id >0 {
		// запустить мент.автоматизм
		RunMentalMentalAutomatizm(matmzm)
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
		infoCreateAndRunNewActionMentAtmzmFromAction(mentalInfoStruct.motorAtmzmID)
	}
	// получили mentalInfoStruct.mImgID
	currentInfoStructId=6 // определение актуального поля mentalInfo
}
func infoCreateAndRunNewActionMentAtmzmFromAction(actImgID int){
	if actImgID ==0 {// ID действия мот.автоматизма, рвущегося на выполнение
		return
	}

// сопоставление
newActionID:=todoComparison(actImgID)
if newActionID==0{// не получилось
	return
}

	id, matmzm := createMentalAutomatizmID(0, actImgID, 1)
	if id >0 {
		// запустить мент.автоматизм
		RunMentalMentalAutomatizm(matmzm)
	}
}
/* сопоставить действия actImgID с другими из Правил и т.п. инфы и если компоненты найденных действий
дают более желательный результат, то Создать новый образ действий из таких компонентов.
Или же, если действия actImgID в данных условиях прогнозируют плохой эффект, то
попытаться сгенерировать новое действие по имеющейся информации.

Нужно иметь в виду, что в Vernike_detector.go есть массив памяти фраз, накапливается в течении дня
var MemoryDetectedArr []MemoryDetected - структур фразы с контекстным окружением
и Verbal.PhraseID[] - можно найти в этом массиве для ориентировки что бы ло раньше и позже.
MemoryDetectedArr - как бы оперативная память фраз для сопоставлений.
При активности дерева автоматизмов возникает Правило: Verbal оператора -> Verbal бота -> эффект
*/
func todoComparison(actImgID int)(int){

	// TODO сопоставление
	return 0
}
//////////////////////////////////////////////////////////



/* №7 создать и запустить ментальный автоматизм запуска моторного автоматизма по действию ActionsImageID -
ВСЕГДА ПОСЛЕ ПОЛУЧЕНИЯ ОБРАЗА ДЕЙСТВИЯ mentalInfoStruct.ActionsImageID
Создается моторный автоматизм (если такого еще нет), привязанный к ветке текущей активности дерева автоматизмов.
*/
var prevMotorAtmzmID=0
func infoFunc7() {
	if mentalInfoStruct.ActionsImageID >0 {
		infoCreateAndRunMentMotorAtmzmFromAction(mentalInfoStruct.ActionsImageID)
	}
	// получили mentalInfoStruct.mImgID
	currentInfoStructId=7 // определение актуального поля mentalInfo
}
func infoCreateAndRunMentMotorAtmzmFromAction(ActionsImageID int){
	if ActionsImageID ==0 {
		return
	}
	motorID,motorAtmzm:=createNewAutomatizmID(0,detectedActiveLastNodID,mentalInfoStruct.ActionsImageID)
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
	actImgID,_:=CreateNewlastMentalActionsImagesID(0,5,motorID)
if actImgID==0{
	return
}
	id, matmzm := createMentalAutomatizmID(0, actImgID, 1)
	if id >0 {
		mentalInfoStruct.mentalAtmzmID=id
		// запустить мент.автоматизм
		RunMentalMentalAutomatizm(matmzm)
	}
}
//////////////////////////////////////////////////////////











//////////////////////////////////////////////////////////
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