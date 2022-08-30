/*
восприятие действий и фраз с Пусльта

Все рефлексторные и автоматические активности начинаются отсюда.

_____________________________
Сначала активируестся Дерво рефлексов и собираются рефлексы на выполнение, но пока не выполняются,
потом активируестся Дерво автоматизмови собираются автоматизмы на выполнение, но пока не выполняются,
если возникает ориентировчный рефлекс, то
активируестся Дерво понимания (ментальных автоматизмов) и решается, что делать дальше.
если нет ориентировочного рефлекса, то
потом выполняютсяя автоматизмы, если их нет - то рефлексы.
________________________________

Создание образов различной иерархии контекстов восприятия:
BaseStyleArr - образ сочетаний активных Базовых контекстов
TriggerStimulsArr - образ сочетаний пусковых стимулов

Создание образов сочетаний ID действий:
*/

package reflexes

import (
	"BOT/brain/action_sensor"
	"BOT/brain/gomeostas"
	"BOT/brain/psychic"
	"BOT/brain/sleep"
)

///////////////////////////////////////
func loadImagesArrs() {
	//  загрузить образы сочетаний базовых стилей
	loadBaseStyleArr()
	//  загрузить образы сочетаний пусковых стимулов
	loadTriggerStimulsArr()

}

//  отслеживане времени с последнего изменения условий с Пульта в пульсах
var lastActivnostFromPult=0
/////////////////////////////////////////////
// ПУЛЬС рефлексов
var ReflexPulsCount = 0 // передача тика Пульса из brine.go
var LifeTime = 0
var EvolushnStage = 0 // стадия развития
var IsSlipping = false

// коррекция текущего состояния гомеостаза и базового контекста с каждым пульсом
func ReflexCountPuls(evolushnStage int, lifeTime int, puls int, isSlipping bool) {
	LifeTime = lifeTime
	EvolushnStage = evolushnStage
	ReflexPulsCount = puls // передача номера тика из более низкоуровневого пакета
	IsSlipping = isSlipping

	if activetedPulsCount != ReflexPulsCount { // защита от повторных срабатываний
		if gomeostas.IsNewConditions { // изменились условия
			ActiveFromConditionChange()
			lastActivnostFromPult=ReflexPulsCount
			return
		}
		/* если условия не меняются более 20 сек, то пусть срабатывает простейший инстинкт
		 только если Базоваое состояние Плохо или Хорошо
		 */
		if ReflexPulsCount - lastActivnostFromPult >20{
			bc:=gomeostas.CommonBadNormalWell
			if bc!=2 {
				// найти и выполнить простейший безусловный рефлекс
				findAndExecuteSimpeReflex()
			}
			lastActivnostFromPult=ReflexPulsCount // новый период 10 секундного ослеживания
		}
	}

	if WakeUpping {
		sleep.WakeUpping()
	}

	// обнулить причину возможного запуска рефлекса
	if oldActiveCurTriggerStimulsID > 0 && oldActiveCurTriggerStimulsPulsCount > (ReflexPulsCount+10) {
		oldActiveCurTriggerStimulsID = 0
	}
}

//////////////////////////////////////////////////////////////

//////  АКТИВАЦИЯ ДЕРЕВА РЕФЛЕКСОВ ПО изменению условий, действиям с Пульта или фразе с Пульта

/*  Вид активации дерева рефлексов:
1 - изменение сочетания базовых контекстов
2 - действия с Пульта
3 - фраза с Пульта
*/
var ActivationTypeSensor = 0

////////////////////////////////////
// текущее восприятие ID образов
//обновляющихся при каждом событии с Пульта или достаточно сильном изменении Базовых параметров
var ActiveCurBaseID = 0           // ID Базового состояния CommonBadNormalWell
var ActiveCurBaseStyleID = 0      // ID сочетания базовых контекстов BaseStyle
var ActiveCurTriggerStimulsID = 0 // ID теущего активного образа сочетаний пусковых стимулов TriggerStimuls
/* предыдущий образ сочетания пусковых стимулов
используется как ПРИЧИНА последующих событий - для формирования условных рефлексов
потому как ОБНУЛЯЕТСЯ при:
1) активации дерева рефлексов, если вызвало какое-то действие
2) через 10 пульсов после записи значения - типа причина устаревает
*/
var oldActiveCurTriggerStimulsID = 0

func setOldActiveCurTriggerStimulsVal(val int) {
	oldActiveCurTriggerStimulsID = val
	oldActiveCurTriggerStimulsPulsCount = ReflexPulsCount
}

// момент записи значения в тике Пульса
var oldActiveCurTriggerStimulsPulsCount = 0

////////////////////////
/* Дерево рефлексов активируется при любом изменении условий с проверкой по каждому пульсу.
 */
var activetedPulsCount = 0 // против многократных срабатываний
func ActiveFromConditionChange() {
	if activetedPulsCount > 0 { // ждет следующего пульса
		return
	}
	// очищать прежние акции с пульта при смене сочетания Базовых контекстов.
	action_sensor.DeactivationTriggers()

	activetedPulsCount = ReflexPulsCount
	ActivationTypeSensor = 1

	ActiveCurBaseID = gomeostas.CommonBadNormalWell
	// определение текущего сочетания ID Базовых контекстов
	bsIDarr := gomeostas.GetCurContextActiveIDarr()

	// создаем новый образ Базовых контекстов, если такого еще нет
	ActiveCurBaseStyleID, _ = createNewBaseStyle(0, bsIDarr)

	// активировать дерево рефлексов
	activeReflexTree()

	// активировать дерево автоматизмов
	res := psychic.SensorActivation(ActivationTypeSensor)
	if res { // блокировать выполнение рефлексов
		return
	}

	toRunRefleses()
}

////////////////////////
// активировать дерево действием reflexes.ActiveFromAction()
func ActiveFromAction() {
	ActivationTypeSensor = 2

	ActiveCurBaseID = gomeostas.CommonBadNormalWell
	// определение текущего сочетания ID Базовых контекстов
	bsIDarr := gomeostas.GetCurContextActiveIDarr()

	// создаем новый образ Базовых контекстов, если такого еще нет
	ActiveCurBaseStyleID, _ = createNewBaseStyle(0, bsIDarr)

	// создаем новый образ Пусковых стимулов, если такого еще нет
	CreateNewTriggerStimulsImage()

	// активировать дерево рефлексов
	activeReflexTree()

	// активировать дерево автоматизмов
	res := psychic.SensorActivation(ActivationTypeSensor)
	if res { // блокировать выполнение рефлексов
		return
	}

	toRunRefleses()

}

////////////////////////
// активировать дерево фразой  reflexes.ActiveFromPhrase() - только для условных рефлексов
// причем, с учетом недавного действия (или нескольких действий) т.к. контекст действий сохраняется 10 сек.
func ActiveFromPhrase() {
	ActivationTypeSensor = 3

	ActiveCurBaseID = gomeostas.CommonBadNormalWell
	// определение текущего сочетания ID Базовых контекстов
	bsIDarr := gomeostas.GetCurContextActiveIDarr()

	// создаем новый образ Базовых контекстов, если такого еще нет
	ActiveCurBaseStyleID, _ = createNewBaseStyle(0, bsIDarr)

	// создать новое сочетание пусковых стимулов если такого еще нет
	CreateNewTriggerStimulsImage()

	// активировать дерево рефлексов
	//BS:=gomeostas.GetCurContextActiveIDarr() // текущий стиль Бозовых контекстов
	activeReflexTree()

	// активировать дерево автоматизмов
	res := psychic.SensorActivation(ActivationTypeSensor)
	if res { // блокировать выполнение рефлексов
		return
	}

	toRunRefleses()
}

/////////////////////////////////////////////////////////

//////////////////////////////////////////////////////////////
/* создание иерархии образов контекстов условий и пусковых стимулов в виде ID образов в [3]int
 создать последовательность уровней условий в виде массива  ID последовательности ID уровней
В случае отсуствия пусковых стимулов создается ID такого отсутсвия, пример такой записи: 2|||0|0|
*/
func getConditionsArr(lev1ID int, lev2 []int, lev3 []int, PhraseID []int, ToneID int, MoodID int) []int {
	arr := make([]int, 3)
	arr[0] = lev1ID
	arr[1], _ = createNewBaseStyle(0, lev2)
	arr[2], _ = CreateNewlastTriggerStimulsID(0, lev3, PhraseID, ToneID, MoodID)
	return arr
}

////////////////////////////////////////////////////
