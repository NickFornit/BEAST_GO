/*  Выполнение действий. Выдать на Пульт акции Beast

вывести на Пульт действия Beast строкой todoAction("xcvxvxcv") с возможностью блокировки
Каждая акция - в формате: вид действия (1 - действие рефлекса, 2 - фраза) затем строка акции,
например: "1|Предлогает поиграть" или "2|Привет!"
Можно передавать неограниченную последовательность акций, разделяя их "||"
например: "1|Предлогает поиграть||2|Привет!"
*/

package reflexes

import (
	"BOT/brain/gomeostas"
	"BOT/brain/psychic"
	"BOT/brain/sleep"
	termineteAction "BOT/brain/terminete_action"
	"BOT/lib"
	"sort"
	"strconv"
)



/* Образы сочетаний действий Beast - при выполнении ТОЛЬКО рефлексов НО НЕ автоматизмов */

// последний образ сочетаний действий Beast
var ActivedTerminalImage []int // массив действий Beast
// предыдущий  образ сочетаний действий Beast
var oldActivedTerminalImage []int

// сохранить образ действий, обеспечивая сортировку ID в массиве
func UpdateActivedTerminalImage(actionsIdArr []int) {
// сохранить предыдущий образ
	oldActivedTerminalImage=nil
	//for _, v := range ActivedTerminalImage {
	for i := 0; i < len(ActivedTerminalImage); i++ {
		oldActivedTerminalImage=append(oldActivedTerminalImage,ActivedTerminalImage[i])
	}
// сохранить текущий образ, обеспечивая сортировку ID в массиве
	ActivedTerminalImage=nil
	sort.Ints(actionsIdArr)
	//for _, v := range actionsIdArr {
	for i := 0; i < len(actionsIdArr); i++ {
		ActivedTerminalImage=append(ActivedTerminalImage,actionsIdArr[i])
	}
	return
}

var WakeUpping=false // true - нужно проснуться 

/* выдать на Пульт подряд акции одного безусловного рефлекса
Рефлексы совершаются со средней "силой" ==5 т.к. нет мозжечкового механизма оптимизации
 */
func TerminateGeneticReflaxActions(id int){
	if id == 0 { return	}
	reflex:=GeneticReflexes[id]
	if reflex == nil { return }
	// сохранить образы сочетаний действий
	UpdateActivedTerminalImage(reflex.ActionIDarr)

	var out = "1|"
	for i := 0; i < len(reflex.ActionIDarr); i++ {
		//проснуться
		if reflex.ActionIDarr[i] == 20 { WakeUpping = true }
		if i > 0 { out += ", "}
		out += termineteAction.TerminalActonsNameFromID[reflex.ActionIDarr[i]]
	}
	todoAction(out)
	return
}

// выдать на Пульт подряд акции массива ID БЕЗУСЛОВНЫХ рефлексов
// reflexKind: 1-древний безусловный, 2-новый безусловный, 2-условный
func TerminateGeneticAllReflaxActions(reflexesIdArr []int, reflexKind int){
	if reflexesIdArr == nil { return }
	// очистить буфер передачи действий на пульт
	// lib.ActionsForPultStr = ""
	lastActivnostFromPult = ReflexPulsCount // новый период 10 секундного ослеживания

	// оставить только уникальные действия рефлексов и выдать такой список
	var unicumArr = make(map[int]int)
	// var rIdArr=make(map[int]int)// привязать действие к ID рефлекса
	for i := 0; i < len(reflexesIdArr); i++ {
		reflex := GeneticReflexes[reflexesIdArr[i]]
		if reflex != nil {
			for j := 0; j < len(reflex.ActionIDarr); j++ {
				unicumArr[reflex.ActionIDarr[j]] = reflexesIdArr[i]
			}
		}
	}
	if len(unicumArr) == 0 { return	}
	var out = "1|<b>БЕССМЫСЛЕННЫЙ безусловный рефлекс:</b><br>"
	var n = 0
	var aImage []int
	for k, rID := range unicumArr {
		// проснуться
		if k == 20{ WakeUpping=true	}
		if n > 0 { out +=", " }
		ta := termineteAction.TerminalActonsNameFromID[k]
		rIDstr := strconv.Itoa(rID)
		out += "<img src=\"/img/edit.png\" style=\"cursor:pointer;\" onClick=\"edit_b_reflex(" + rIDstr + ")\" title=\"Изменить действия рефлекса ID=" + rIDstr + "\">&nbsp;" + ta
		// при этом меняются гомео-параметры:
		ExpensesGomeostatParametersAfterAction(k)
		aImage=append(aImage, k)
		n++
	}
	// сохранить образы сочетаний действий
	UpdateActivedTerminalImage(aImage)
	// детектор нового вычленяет новые условия по сравненеию с условиями рефлекса и обрабатывает это:
	// здесь ВСЕГДА rank=0 т.к. тут обрабатываются только безусловные
	updateNewsConditions(0)
	todoAction(out)
	// очистить массивы рефлексов данного типа
	switch reflexKind{
	case 1: oldReflexesIdArr=nil
	case 2: geneticReflexesIdArr=nil
	// здесь не бывает	case 3: conditionReflexesIdArr=nil
	}

	return
}

// выдать на Пульт подряд акции массива ID УСЛОВНЫХ рефлексов
func TerminateConditionAllReflaxActions(reflexesIdArr []int) {
	if reflexesIdArr == nil || len(reflexesIdArr) == 0 { return }
	var out = "1|<b>БЕССМЫСЛЕННЫЙ условный рефлекс:</b><br>"

	// может быть только один член в reflexesIdArr
	reflex := ConditionReflexes[reflexesIdArr[0]]
	var n = 0
	var aImage []int
	// здесь не может быть неуникальных ID действий, так что просто:
	for i := 0; i < len(reflex.ActionIDarr); i++ {
		k:=reflex.ActionIDarr[i]
		//проснуться
		if k == 20 { WakeUpping = true }
		if n > 0 { out += ", " }
		out += termineteAction.TerminalActonsNameFromID[k]
		// при этом меняются гомео-параметры:
		ExpensesGomeostatParametersAfterAction(k)
		aImage=append(aImage,k)
		n++
	}
	// сохранить образы сочетаний действий
	UpdateActivedTerminalImage(aImage)
	// детектор нового вычленяет новые условия по сравненеию с условиями рефлекса и обрабатывает это:
	// для создания условных рефлексов более высого ранга на основе усл.рефлексов
	rank := ConditionReflexes[reflex.ID].rank
	updateNewsConditions(rank + 1)
	todoAction(out)
	conditionReflexesIdArr = nil
}

// выдать на Пульт подряд акции массива по каждому рефлексу отдельной строкой действий
/*
func TerminateGeneticEachReflaxActions(reflexesIdArr []int){
	if reflexesIdArr==nil{
		return
	}
	var out=""
	for i := 0; i < len(reflexesIdArr); i++ {
		if i>0{out+="||"} // пошел следующий рефлекс
		reflex:=GeneticReflexes[reflexesIdArr[i]]
		if reflex==nil{return}
		out+="1|"
		for j := 0; j < len(reflex.ActionIDarr); j++ {
			if j>0{out +=", "}
			out += TerminalActonsNameFromID[reflex.ActionIDarr[j]]
		}
	}
	todoAction(out)
}
 */

/* изменение гомео-параметров при действии
для рефлексов - нет соррекции по силе, она средняя =5
 */
func ExpensesGomeostatParametersAfterAction(actID int) {
	if gomeostas.NotAllowSetGomeostazParams { return }
	se := termineteAction.TerminalActionsExpensesFromID[actID]
	if se != nil {
		for j := 0; j < len(se); j++ {
			gomeostas.GomeostazParams[se[j].GomeoID] += se[j].Diff
			if gomeostas.GomeostazParams[se[j].GomeoID] > 100 {
				gomeostas.GomeostazParams[se[j].GomeoID] = 100
			}
			if gomeostas.GomeostazParams[se[j].GomeoID] < 0{
				gomeostas.GomeostazParams[se[j].GomeoID] = 0
			}
		}
	}
}

// совершение действий и возможность их блокировки
func todoAction(out string) {
	if !psychic.GetAllowReflexRuning() {
		return }
	// блокировка во время сна
	if IsBlockingMotorsAction() {
		return }

	lib.SentActionsForPult(out)
}

/* Блокировка рефлексотрных действий из Психики или во сне */
func IsBlockingMotorsAction() bool {
	notAllow1 := psychic.NotAllowReflexesAction()
	notAllow2 := sleep.NotAllowReflexesAction()
	if notAllow1 || notAllow2 { return true	}
return false
}

// запустить готовые к выполнению рефлексы
func toRunRefleses() {
	if !psychic.GetAllowReflexRuning() {
		lib.WritePultConsol("<span style='color:red'>Рефлекс <b>заблокирован</b></span>")
		return }
	// очистить буфер передачи действий на пульт
	// lib.ActionsForPultStr = ""
	lastActivnostFromPult=ReflexPulsCount // сбросить отчет времени бездействия
	if len(conditionReflexesIdArr) > 0 { // есть условные рефлексы
		// удалить более низкоуровневые рефлексы
		geneticReflexesIdArr = nil
		oldReflexesIdArr = nil
		// выдать на пульт действия
		conditionReflexesIdArr = lib.UniqueArr(conditionReflexesIdArr)
		TerminateConditionAllReflaxActions(conditionReflexesIdArr)
		// ОБНУЛЯЕТСЯ при активации дерева рефлексов, если вызвало какое-то действие
		oldActiveCurTriggerStimulsID = 0
		return
	}
	if len(geneticReflexesIdArr) > 0 { // есть новые безусловные рефлексы
		NoUnconditionRefles = ""
		// выдать на пульт действия
		geneticReflexesIdArr = lib.UniqueArr(geneticReflexesIdArr)
		TerminateGeneticAllReflaxActions(geneticReflexesIdArr,2)
		// В TerminateGeneticAllReflaxActions детектор нового updateNewsConditions() вычленяет новые условия по сравненеию с условиями рефлекса и обрабатывает это:
		// ОБНУЛЯЕТСЯ при активации дерева рефлексов, если вызвало какое-то действие
		oldActiveCurTriggerStimulsID = 0
		return
	}
	if len(oldReflexesIdArr) > 0 {// есть старые безусловные рефлексы
		// выдать на пульт действия
		oldReflexesIdArr = lib.UniqueArr(oldReflexesIdArr)
		TerminateGeneticAllReflaxActions(oldReflexesIdArr,1)
		// В TerminateGeneticAllReflaxActions детектор нового updateNewsConditions() вычленяет новые условия по сравненеию с условиями рефлекса и обрабатывает это:
		//ОБНУЛЯЕТСЯ при активации дерева рефлексов, если вызвало какое-то действие
		oldActiveCurTriggerStimulsID = 0
		return
	}
	/* если нет условного и безусловного рефлекса, то совершается самый простейший безусловный рефлекс
	по сочетаниям редактора http://go/pages/terminal_actions.php
	Данный редактор связывает действие с тем, какие гомео-параметры улучшает данное действие.
	*/
	findAndExecuteSimpeReflex()
}