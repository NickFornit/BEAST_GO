/*  активация дерева рефлексов изменением:
текущих условий, 
действиями с Пульта 
и фразой с Пульта
функциями из perception.go
ActiveFromConditionChange()
ActiveFromAction()
ActiveFromPhrase()

Здесь отрабатываются три вида рефлексов:
1. Древние безусловные - у которых в условиях не прописаны пусковые стимулы. Собираются в oldReflexesIdArr
2. Новые безусловные - с прописанными пусковыми стимулами. Собираются в geneticReflexesIdArr
3. Условные рефлексы - на основе предыдущих безусловных (но не древних, а полных) или условных - связанных с новыми стимулами. Собираются в conditionReflexesIdArr
Каждый последующий вид рефлексов имеет приоритет над предыдущими, т.е. те не выполняются, если есть приоритетный.

*/

package reflexes

import (
	"BOT/brain/action_sensor"
	"BOT/brain/gomeostas"
	"BOT/brain/psychic"
	termineteAction "BOT/brain/terminete_action"
	"BOT/lib"
	"strconv"
)

/*
запрет показа карты при обновлении против паники типа "одновременная запись и считывание карты"
Использовать для всех операций записи узлов дерева
func xxxxxx(){
notAllowScanInReflexesThisTime=true
ОПЕРАЦИИ ...
notAllowScanInReflexesThisTime=false
}
*/


//////////////////////////////////
func readyForRecognitionRflexes() { // init() для дерева распознавания рефлексов

/*
		ActiveCurBaseID =1
		ActiveCurBaseStyleID=22 //14
		ActiveCurTriggerStimulsID =3
			activeReflexTree()
	}
 */
	initConditionReflex()


//	FormingConditionsRefleaxFromList("")
}
//////////////////////////////////////

//для частичного распознавания нужен массив текущих активных Базовых контекстов
var curBaseCondArr []int
// массив текущих пусковых стимулов
var curPultActionsArr []int
//////////////////////////////////
var detectedActiveLastNodID=0
/* распознавание идет строго по совпадающим веткам, без ветвлений.
Но на каждом уровне смотрятся ветвления дочек - для нахождения не точного соотвествия.
Например, если рефлексы имеют только один Базовый контекст, а текущее состояние Beast - сочетание нескольких,
то в результате должны сработать все рефлексы, для каждой цифры в сочетании образа текущих условий.
Но если в рефлексе заданы несколько условий и такой образ точно совпадент с текущим образом условий,
то именно этот рефлекс и сработает.
*/
var detectedActiveLevel=0 // уровень условий, до которого дошло распознавание в дереве

//собираются рефлексы, подходящие для текущих оразов услвий:
var oldReflexesIdArr []int //собираются Древние безусловные - у которых в условиях не прописаны пусковые стимулы.
var geneticReflexesIdArr []int //собираются Новые безусловные - с прописанными пусковыми стимулами.
var conditionReflexesIdArr []int //собираются Условные рефлексы - на основе предыдущих безусловных или условных - связанных с новыми стимулами.
/////////////////////////////////////////////

//  сообщить на Пульт, что при данных условиях нет б.рефлекса.
var NoUnconditionRefles=""

//////////////////////////////////////////////////////////////////
// распознавание рефлексов
func activeReflexTree(){
	detectedActiveLastNodID=0
	detectedActiveLevel=0

	oldReflexesIdArr=nil
	geneticReflexesIdArr=nil
	conditionReflexesIdArr=nil

	//для частичного распознавания нужен массив текущих активных Базовых контекстов
	curBaseCondArr=gomeostas.GetCurContextActiveIDarr()
	if curBaseCondArr==nil{// еще не определились актичные Базовые контексты
		return
	}
	// массив текущих пусковых стимулов
	curPultActionsArr=action_sensor.CheckCurActions()
	// для условныз реф-в
	//=GetActiveContextInfo()

	// ЗАПАРИВАЕТ lib.WritePultConsol("Активация дерева рефлексов.")
	// вытащить 3 уровня условий в виде ID их образов
	condArr:=getActiveConditionsArr(ActiveCurBaseID, ActiveCurBaseStyleID,ActiveCurTriggerStimulsID)
//	ActiveCurTriggerStimulsID=1; condArr:=[]int{1, 2, ActiveCurTriggerStimulsID}; // проверка подстановкой произвольных сочетаний


if ActiveCurBaseStyleID==22{
	ActiveCurBaseStyleID=22
}

	// основа дерева
	cnt := len(ReflexTree.Children)
	for n := 0; n < cnt; n++ {
		node := ReflexTree.Children[n]
		lev1 := node.baseID
		if condArr[0] == lev1 {
			detectedActiveLastNodID=node.ID
			ost:=condArr[1:]
			detectedActiveLevel=1
			findReflexesNodes(detectedActiveLevel,ost, &node,1)
			//findReflexesNodes(1,condArr, &node,1)

			break // только один из Базовых состояний
		}
	}

	// результат поиска:
	if detectedActiveLastNodID > 0 {// найден узел
// если есть фраза или условный рефлекс, то погасить б.рефлексы
		if ActivationTypeSensor==3 || len(conditionReflexesIdArr)>0 {
			// удалить более низкоуровневые рефлексы
			geneticReflexesIdArr = nil
			oldReflexesIdArr = nil
		}else {
			// единственная реакция 9. Игнорирует? - для диалога Пульта чтобы сформировать рефлекс
			isIgnor := checkIgnorOnly(oldReflexesIdArr, geneticReflexesIdArr)
			// нет старых или новых безусловных рефлексов для текущих условий и если игнорирует
			if (len(oldReflexesIdArr) == 0 && len(geneticReflexesIdArr) == 0) || isIgnor {
				/* если в GeneticReflexes (список всех dnk_reflexes.txt) есть совпадающее условие,
				то создать узел дерева
				*/
				addGeneticReflexesToTree(detectedActiveLastNodID, condArr)

				//  сообщить на Пульт, что при данных условиях нет б.рефлекса.
					if EvolushnStage == 0 { // только для стадии безусловных рефлексов
						if isIgnor {
							NoUnconditionRefles = "IGNORED" + GetCurrentConditionsStr() //СТРОКА УСЛОВИЙ ДЛЯ РЕФЛЕКСА
						} else {
							NoUnconditionRefles = "NOREFLEX" + GetCurrentConditionsStr() //СТРОКА УСЛОВИЙ ДЛЯ РЕФЛЕКСА
						}
						return
					}
			}
			// в консоль:
			consol := "<br>__________ РЕФЛЕКС: "
			for c := 0; c < len(oldReflexesIdArr); c++ {
				consol += "ID=" + strconv.Itoa(oldReflexesIdArr[c]) + "; "
			}
			for c := 0; c < len(geneticReflexesIdArr); c++ {
				consol += "ID=" + strconv.Itoa(geneticReflexesIdArr[c]) + "; "
			}
			//consol+="<br>"
			if (len(oldReflexesIdArr)> 0 || len(geneticReflexesIdArr) > 0) {
				lib.WritePultConsol(consol)
			}else{// не прописано никаких реакций
				lib.WritePultConsol("Не определен рефлекс")
				NoUnconditionRefles = "NOREFLEX" + GetCurrentConditionsStr() //СТРОКА УСЛОВИЙ ДЛЯ РЕФЛЕКСА
			}

			//////// ДЛЯ ПСИХИКИ:
			veryActual, targetArrID, acrArr := GetActualReflexAction()
			if len(acrArr)>6{// может быть и 64 если
// просто ограничить 3 предположительными акциями, т.к. есть сортировка по убыванию значимости
				acrArr=acrArr[:3]
				// на стадии рефлексов был сигнал NOREFLEX и диалог на пульте
			}
			// передать в психику информацию
			psychic.GetReflexInformation(veryActual, targetArrID, acrArr)
			//////////////////////

			if EvolushnStage < 2 { // сразу запустить имеющиеся рефлексы
				toRunRefleses()
			} // иначе сначала будут проверены автоматизмы в perception.go
		}

	}else{// вообще еще нет такого случая :) т.к. всегда есть нулевая
			//  сообщить на Пульт, что при данных условиях нет б.рефлекса.
	//	if EvolushnStage==0 { // только для стадии безусловных рефлексов
	//		NoUnconditionRefles = "NOREFLEX" + GetCurrentConditionsStr() //СТРОКА УСЛОВИЙ ДЛЯ РЕФЛЕКСА
	//	}
// ничего не делать
	}
}
/* На каждом уровне допускаются ветвления его дочек - для нахождения не точного соотвествия, если не было найдено точное.
Например, если рефлексы имеет только один Базовый контекст, а текущее состояние Beast - сочетание нескольких,
то в результате должны сработать все рефлексы, для каждой цифры в сочетании образа текущих условий.
Но если в рефлексе заданы несколько условий и такой образ точно совпадент с текущим образом условий,
то именно этот рефлекс и сработает.

isExactly: 0 - сработал неточный рефлекс, не смотреть блок точного совпадения
 */
// БЕЗ РЕКУРСИИ т.к. всего 2 уровня проверяется
func findReflexesNodes(level int,cond []int,node *ReflexNode,isExactly int){
	detectedActiveLevel=level
	if len(cond)==0{
		return
	}

	// если последний уровень Дерева
	if level==2{// смотреть условные рефлексы
		if conditionRexlexFound(cond){
			return
		}
	}

	for n := 0; n < len(node.Children); n++ {
		cld := node.Children[n]
		if cld.StyleID==cond[0] {
			// только если нет пусковых стимулов, позволено смотреть древние рефлексы
			if ActiveCurTriggerStimulsID==0 {
				if ReflexTreeFromID[cld.ID].GeneticReflexID > 0 {
					// древний рефлекс
					oldReflexesIdArr = append(oldReflexesIdArr, ReflexTreeFromID[cld.ID].GeneticReflexID)
					detectedActiveLevel=level+1
					detectedActiveLastNodID=cld.ID
					return
				}
			}
			// есть ли условный рефлекс?
			if conditionRexlexFound(cond[1:]){
				return
			}
			if cld.ActionID==cond[1]{
				if ReflexTreeFromID[cld.ID].GeneticReflexID > 0{
					geneticReflexesIdArr = append(geneticReflexesIdArr, ReflexTreeFromID[cld.ID].GeneticReflexID)
					detectedActiveLastNodID=cld.ID
					detectedActiveLevel=level+1
					return
				}
			}
		}
	}

	return
}
////////////////////////////////////////////////////
func compareCondImade(level int,rArr []int)(bool){
	switch level{
	case 1:
		for i := 0; i < len(curBaseCondArr); i++ {
			// есть такое значение в массиве
			for j := 0; j < len(rArr); j++ {
				if curBaseCondArr[i] == rArr[j] {
					return true
				}
			}
		}
	case 2:
		for i := 0; i < len(curPultActionsArr); i++ {
			for j := 0; j < len(rArr); j++ {
				if curPultActionsArr[i] == rArr[j] {
					return true
				}
			}
		}
	}

	return false
}
////////////////////////////




//////////////////////////////////////////////////////////////
/* создание иерархии АКТИВНЫХ образов контекстов условий и пусковых стимулов в виде ID образов в [3]int
 создать последовательность уровней условий в виде массива  ID последовательности ID уровней
*/
func getActiveConditionsArr(lev1 int, lev2 int, lev3 int)([]int){
	arr:=make([]int, 3)
	arr[0]=lev1
	arr[1]=lev2
	arr[2]=lev3
	return arr
}
////////////////////////////////////////////////////


/*СТРОКА УСЛОВИЙ ДЛЯ безусловного РЕФЛЕКСА типа   1|2,5,8|11
Базовое состояние
Активные контексты
Пусковые стимулы
 */
func GetCurrentConditionsStr()(string){

//ID базового состояния (1 уровень)
	var out=strconv.Itoa(gomeostas.CommonBadNormalWell)+"|"

//ID (2) актуальных контекстов через запятую
	bs:=gomeostas.GetCurContextActiveIDarr()
	for i := 0; i < len(bs); i++ {
		if i>0{out+=","}
		out+=strconv.Itoa(bs[i])
	}
	out+="|"

//ID (3) пусковых стимулов через запятую
	as:=action_sensor.CheckCurActionsContext()
	for i := 0; i < len(as); i++ {
			if i>0{out+=","}
			out+=strconv.Itoa(as[i])
	}

return out
}
////////////////////////////////////////////////////////



func checkIgnorOnly(oldReflexesIdArr []int,geneticReflexesIdArr []int)(bool) {
	var isIgnor = false
	if len(GeneticReflexes) >0 {
		if len(oldReflexesIdArr) == 1 && GeneticReflexes[oldReflexesIdArr[0]]!=nil {
			gr := GeneticReflexes[oldReflexesIdArr[0]]
			if gr.ActionIDarr[0] == 9 {
				isIgnor = true
			}
		} else {
			if len(geneticReflexesIdArr) == 1 && GeneticReflexes[geneticReflexesIdArr[0]]!=nil{
				gr := GeneticReflexes[geneticReflexesIdArr[0]]
				if gr.ActionIDarr[0] == 9 {
					isIgnor = true
				}
			}
		}
	}
	return isIgnor
}
/////////////////////////////////////////////



//////////////////////////////////////////////
/* сразу после активации дерева передать инфу для Психики
		// передать в психику информацию
		acrArr:=GetActualReflexAction()
		psychic.GetReflexInformation(acrArr)
*/
// вернуть 1)veryActual 2)targetArrID 3)actArtr
func GetActualReflexAction()(bool,[]int,[]int){
	var actArtr []int

	/* выявить ID парамктров гомеостаза как цели для улучшения в данных условиях
	здесь, чтобы сразу получить veryActual и targetArrID для возврата
	targetArrID - отсортирован по убыванию значимости
	 */
	veryActual,targetArrID:=gomeostas.FindTargetGomeostazID()

	//есть ли подходящий по условиям безусловный или условный рефлекс и сделать автоматизм по его действиям
	/* возвращает текущие массивы найденных при активации видов рефлексов
	   1-conditionReflexesIdArr []int - Условные рефлексы
	   2-geneticReflexesIdArr []int - Новые безусловные
	   3-oldReflexesIdArr []int - Древние безусловные
	    condArr,geneticArr,OldArr:=GetActualReflex()
	*/
	condArr,geneticArr,OldArr:=GetActualReflex()
	if condArr!=nil && len(condArr)>0{
		for i := 0; i < len(condArr); i++ {
			act:=ConditionReflexes[condArr[i]]
			for j := 0; j < len(act.ActionIDarr); j++ {
				actArtr = append(actArtr, act.ActionIDarr[j])
			}
		}
		return veryActual,targetArrID,actArtr
	}
	if geneticArr!=nil && len(geneticArr)>0{
		for i := 0; i < len(geneticArr); i++ {
			act:=GeneticReflexes[geneticArr[i]]
			for j := 0; j < len(act.ActionIDarr); j++ {
				actArtr = append(actArtr, act.ActionIDarr[j])
			}
		}
		return veryActual,targetArrID,actArtr
	}
	if OldArr!=nil && len(OldArr)>0{
		for i := 0; i < len(condArr); i++ {
			act:=GeneticReflexes[condArr[i]]
			for j := 0; j < len(act.ActionIDarr); j++ {
				actArtr = append(actArtr, act.ActionIDarr[j])
			}
		}
		return veryActual,targetArrID,actArtr
	}

	// остались самые древние реации - действия по их цели (Какие ID гомео-параметров
	//улучшает действие из http://go/pages/terminal_actions.php)
	// выдать массив возможных действий для данных условий чтобы выбрать одно из них, пока еще не испытанное
	for aID, gIDarr := range termineteAction.TerminalActionsTargetsFromID {
		// выбрать подхлдящие ID параметров гомеостаза для данной цели
		aArr:=lib.GetExistsIntArs(targetArrID,gIDarr)
		if aArr==nil{
			continue
		}
		actArtr = append(actArtr, aID)
	}

	return veryActual,targetArrID,actArtr
}
//////////////////////////////////////////////////


/* GetActualReflex() возвращает текущие массивы найденных при активации видов рефлексов
1-conditionReflexesIdArr []int - Условные рефлексы
2-geneticReflexesIdArr []int - Новые безусловные
3-oldReflexesIdArr []int - Древние безусловные
 condArr,geneticArr,OldArr:=reflexes.GetActualReflex()
*/
func GetActualReflex()([]int,[]int,[]int){

	return conditionReflexesIdArr,geneticReflexesIdArr,oldReflexesIdArr
}
/////////////////////////////////////////////



//////////////////////////////////////////////////////

