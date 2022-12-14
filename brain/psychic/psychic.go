/* Общий файл психики

Психика активируется PsychicInit() после активации всех предшествующих структур.

По каждому событию с Пульта или изменению состояния организма активируется дерево автоматизмов.
Если есть автоматизм - перед его выполнением вызывается ориентировочный рефлекс привлечеия внимания orientation_2()
если нет автоматизма - ориентировчный рефлекс оценки ситуации orientation_1()
Привлечение осознанного внимания выявляет конечную цель - найти автоматизм или ничего не делать.

На уровне наивности (нет уверенно решающих проблему автоматизмов)
на первый план выходит отзеркаливание чужих действий и случайные пробы-ошибки.

Понимание смысла воспринимаемого и своих действий начинается с предопределенных
генетические целей действий Beast ID гомео-параметров, которые призвано улучшить данное действие - по его ID:
TerminalActionsTargetsFromID - это – наследственно заданная цель действия, не осознаваемая при его совершении.
Но с опытом каждому действию в конкретных условиях (и к ним добавляются слова и фразы)
будет ассоциироваться смысл (осознаваемая значимость).
Таким образом, отзеркаливая чужие "зеркальные" действия и совершая свои с оценкой результата,
будет пополняться МОДЕЛЬ ПОНИМАНИЯ при данных условиях.
Эта модель, фактически, будет составлена из наборов автоматизмов, привязанных к активной ветке дерева автоматизмов,
из которых один - текущий актуальный, остальные отбракованные и предположительные.
Автоматизмы с внутренними, ментальными действиями будут обеспечивать произвольность.


Это – наследственно заданная цель действия, не осознаваемая при его совершении.
Но с опытом каждому действию в конкретных условиях (и к ним добавляются слова и фразы)
будет ассоциироваться смысл (осознаваемая значимость).

Безусловные рефлексы психики прописываются в виде функций обработки
текущей инфрмационной среды CurrentInformationEnvironment.
В этой среде активируются текущие проблемы и доминанты нерешенной проблемы.

НЕ ЗАБЫВАТЬ для всех функций произвольной активации (по актуализации текущего самоощущения)
ставить блокировку по brain.NotAllowAnyActions - if brain.NotAllowAnyActions{ return }
*/
package psychic

import (
	"BOT/lib"
	"strconv"
)

///////////////////////////////

// true - ПРИ ТЕСТИРОВАНИИ СОХРАНЯТЬ В ФАЙЛАХ ВСЕ ЭЛЕМЕНТЫ
var doWritingFile =false // true  false

///////////////////////////////////////////////////////////////////////
// инициализирующий блок - в порядке последовательности инициализаций
// после condition_reflex.go
func PsychicInit(){

	if EvolushnStage<2{// еще нет психики
		return
	}
	automatizmTreeInit()
	loadActionsImageArr()
	automatizmInit()
	emotionsInit()
	loadActivityInit()
	verbalInit()
	cerebellumReflexInit()
	EpisodeMemoryInit()
	initCurrentInformationEnvironment()
	BaseStateImageInit()
	TriggerAndActionInit()
	loadMentalTriggerAndActionArr()
	rulesInit()
	ImportanceInit()
	loadSituationImage()
	goNextInit()
	rulesMentalInit()
	MentalActionsImagesInit()// -> loadMentalActionsImagesArr()
	mentalAutomatizmInit()
	UnderstandingTreeInit()

	// saveActionImageArr()// сохранить образы сочетаний ответных действий

	// просыпание - создание базового самоощущения CurrentInformationEnvironment
		wakingUp()

//	SensorActivation(1,1,[]int{1})
/*
	atmzm:=findAnySympleRandActions()
	if atmzm!=nil{	}
 */


//	FormingMirrorAutomatizmFromList("/mirror_reflexes_basic_phrases/1_2.txt")

//	FormingMirrorAutomatizmFromTempList("/lib/mirror_basic_phrases_common.txt")


}
/////////////////////////////////////////////////////////////

// ПУЛЬС психики
var PulsCount=0 // передача тика Пульса из brine.go
var LifeTime=0
var EvolushnStage=0 // стадия развития
var IsSleeping =false
func PsychicCountPuls(evolushnStage int,lifeTime int,puls int,isSleeping bool){

	if evolushnStage<2 { // недостаточная стадия развития
		return
	}

	LifeTime=lifeTime
	EvolushnStage=evolushnStage
	PulsCount=puls // передача номера тика из более низкоуровневого пакета
	IsSleeping =isSleeping

	// тики в automatizm_result.go для удобства
	orientarionPuls()
	automatizmActionsPuls()
	moodePulse()
	EpisodeMemoryPuls()

	if IsSleeping {
		sleepingProcess()
	}


	// осознание при включении и бодрствовании - один раз
	if evolushnStage > 3 && PulsCount >4 && !IsSleeping && !AllowConsciousnessProcess {
// начать мышление
		AllowConsciousnessProcess=true
		consciousness(1,0)
	}



	// просыпание - создание базового самоощущения CurrentInformationEnvironment
	//	if psychicPulsCount>3 {
	//		wakingUp()
	//	}


	if CurrentInformationEnvironment.PsyActionImgPulsCount>0 && CurrentInformationEnvironment.PsyActionImgPulsCount>(PulsCount+100){
		CurrentInformationEnvironment.PsyActionImg=ActivityFromIdArr[1]// образ бездействия
		CurrentInformationEnvironment.PsyActionImgPulsCount=0
	}

}
//////////////////////////////////////////////////////////////
var NotAllowAnyActions=false
func SetNotAllowAnyActions(notAllow bool){
	NotAllowAnyActions=notAllow
}
///////////////////////////////////////










// просыпание - создание базового самоощущения CurrentInformationEnvironment
func wakingUp(){

// осознание самоощущения
	SensorActivation(1)

	// очистить всякое при просырании
	usedActIdArr=nil
	usedPraseIdArr=nil
}
/////////////////////////////////////////////////////////////



/////////////////////////////////////////////////////////////////
/*  активация по событиям с Пульта - из perception.go
Для блокировки активации дерева рефлексов вернуть true
 */
var firstStadiesWarning=true // защелка от повторов
var ActivationTypeSensor=0//передача типа акетивации в психику из рефлексов
func SensorActivation(activationType int)(bool){
	if PulsCount<4{
		return false
	}
	ActivationTypeSensor=activationType

	if EvolushnStage<2{// недостаточная стадия развития
		if firstStadiesWarning {
			firstStadiesWarning=false
			lib.WritePultConsol("Стадия развития " + strconv.Itoa(EvolushnStage) + " НЕДОСТАТОЧНА ДЛЯ АВТОМАТИЗМОВ")
		}
			return false
	}

//  получение текущего состояния информационной среды: отражение Базового состояния и Активных Базовых контекстов
//!!!! GetCurrentInformationEnvironment() только при ориентировочном рефлексе - смена самоощущения !!!

atomatizmID:=automatizmTreeActivation()
if atomatizmID>0{

	return true
}


	return false
}
/////////////////////////////////////////////////////////////////




////////////////////////////////////
/* Блокировать выполнение рефлексов на время ожидания результата автоматизма
вызывается из reflex_action.go рефлексов
 */
func NotAllowReflexesAction()(bool){
	if MotorTerminalBlocking {
		return true
	}
	return false
}
////////////////////////////////////



///////////////////////////////////
func SaveAllPsihicMemory(){
	notAllowScanInTreeThisTime = true
	SaveEmotionArr()
	SaveActivityFromIdArr()
	SaveVerbalFromIdArr()
	SaveAutomatizmTree()
	SaveAutomatizm()
	SaveSituationImage()
	SaveActionsImageArr()
	SaveCerebellumReflex()
	SaveUnderstandingTree()
	saveEpisodicMenory()
	SaveBaseStateImageArr()
	SaveTriggerAndActionArr()
	SaveMentalTriggerAndActionArr()
	SavePurposeImageFromIdArr()
	SaveRulesArr()
	SaverulesMentalArr()
	Saveimportance()
	SavegoNext()
	saveInterruptMemory()
	SaveMentalActionsImagesArr()
	SaveMentalAutomatizm()
	SaveProblemDominenta()
	notAllowScanInTreeThisTime = false
}
//////////////////////////////////////


