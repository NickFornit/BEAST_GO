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
	"BOT/lib"
	"BOT/brain/action_sensor"
	"BOT/brain/gomeostas"
	"BOT/brain/psychic"
	termineteAction "BOT/brain/terminete_action"
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

			findReflexesNodes(1,ost, &node,1)

			//break // только один из Базовых состояний
		}
	}

	// результат поиска:
	if detectedActiveLastNodID > 0 {// найден узел
	// единственная реакция 9. Игнорирует? - для диалога Пульта чтобы сформировать рефлекс
		isIgnor:=checkIgnorOnly(oldReflexesIdArr,geneticReflexesIdArr)
// нет старых или новых безусловных рефлексов для текущих условий и если игнорирует
		if (len(oldReflexesIdArr)==0 && len(geneticReflexesIdArr)==0) || isIgnor{
			//  сообщить на Пульт, что при данных условиях нет б.рефлекса.
			if EvolushnStage==0 { // только для стадии безусловных рефлексов
				if isIgnor{
					NoUnconditionRefles = "IGNORED" + getCurrentConditionsStr() //СТРОКА УСЛОВИЙ ДЛЯ РЕФЛЕКСА
				}else {
					NoUnconditionRefles = "NOREFLEX" + getCurrentConditionsStr() //СТРОКА УСЛОВИЙ ДЛЯ РЕФЛЕКСА
				}
				return
			}
		}
		// передать в психику информацию
		veryActual,targetArrID,acrArr:=GetActualReflexAction()
		psychic.GetReflexInformation(veryActual,targetArrID,acrArr)

		if EvolushnStage<2 {// сразу запустить имеющиеся рефлексы
			toRunRefleses()
		}// иначе сначала будут проверены автоматизмы в perception.go

	}else{// вообще еще нет такой ветки
		//  сообщить на Пульт, что при данных условиях нет б.рефлекса.
		if EvolushnStage==0 { // только для стадии безусловных рефлексов
			NoUnconditionRefles = "NOREFLEX" + getCurrentConditionsStr() //СТРОКА УСЛОВИЙ ДЛЯ РЕФЛЕКСА
		}
// ничего не делать, рефлексы образуют новые ветки условий.
	}
}
/* На каждом уровне допускаются ветвления его дочек - для нахождения не точного соотвествия, если не было найдено точное.
Например, если рефлексы имеет только один Базовый контекст, а текущее состояние Beast - сочетание нескольких,
то в результате должны сработать все рефлексы, для каждой цифры в сочетании образа текущих условий.
Но если в рефлексе заданы несколько условий и такой образ точно совпадент с текущим образом условий,
то именно этот рефлекс и сработает.

isExactly: 0 - сработал неточный рефлекс, не смотреть блок точного совпадения
 */

func findReflexesNodes(level int,cond []int,node *ReflexNode,isExactly int){
	if len(cond)==0{
		return
	}

// для проверки
	if node.ID==9{
		node.ID=9
	}

	////////////////////////////////////////////
	if conditionReflexesIdArr!=nil{// уже найден условный рефлекс для данных условий
		return
	}
	if level==2{// смотреть условные рефлексы
		if conditionRexlexFound(cond){
			return
		}
	}

	/////////////////////////////////////////////

	ost:=cond[1:]
	// найти наиболее подходящий следующий узел по точному соотвествию с образом - имеет преимущество над частичным совпадением.
	if isExactly==1 { // искать точное совпадение
		for n := 0; n < len(node.Children); n++ {
			cld := node.Children[n]
			var levID = 0
			switch level {
			case 0:
				levID = cld.baseID
			case 1:
				levID = cld.StyleID
			case 2:
				levID = cld.ActionID
			}
			if cond[0] == levID {
				switch level {
				case 1:
					// только если нет пусковых стимулов, позволено смотреть древние рефлексы
					if ActiveCurTriggerStimulsID==0 {
						if ReflexTreeFromID[cld.ID].GeneticReflexID > 0 {
							// древний рефлекс
							oldReflexesIdArr = append(oldReflexesIdArr, ReflexTreeFromID[cld.ID].GeneticReflexID)
						}
					}
				case 2:
					if ReflexTreeFromID[cld.ID].GeneticReflexID > 0 {
						geneticReflexesIdArr = append(geneticReflexesIdArr, ReflexTreeFromID[cld.ID].GeneticReflexID)
					}
				}
				detectedActiveLastNodID = cld.ID

				level++
				findReflexesNodes(level,ost, &cld, isExactly)
				return // раз совпало, то другие ветки не смотреть
			}
		}
	}//if isExactly==1{// искать точное совпадение



	/* если не найдено точное соотвествие образа, смотеть по совпадающей части образа
	как при частичном распознавании персептрона
	 */
	isExactly=0
/////////////////////////////
// Заплатка: если у узла есть рефлекс, то предварительно учесть его, пока не подставляя в geneticReflexesIdArr
/*
var oldReflex=0
var isPrioritet=false
//учесть рефлекс если еще нет других рефлексов в geneticReflexesIdArr
if node.GeneticReflexID>0 && node.ActionID==0 && oldReflexesIdArr==nil{
	// древний рефлекс
	oldReflex=node.GeneticReflexID
}
 */
/////////////////////////////

	for n := 0; n < len(node.Children); n++ {
		cld:=node.Children[n]
		if cld.ID==50{
			cld.ID=50
		}
		var rArr []int
		switch level{
		case 1:
			if BaseStyleArr[cld.StyleID]!=nil {
				rArr = BaseStyleArr[cld.StyleID].BSarr
			}
		case 2:
			if TriggerStimulsArr[cld.ActionID]!=nil {
				rArr = TriggerStimulsArr[cld.ActionID].RSarr
			}
		}
		if rArr!=nil {
			// есть ли частичное совпадение
			if compareCondImade(level, rArr) {
				switch level {
				case 1:
					// только если нет пусковых стимулов, позволено смотреть древние рефлексы
					/* Не должно вызывать древний рефлекс при неточном совпадении
					if ActiveCurTriggerStimulsID == 0 {
						if ReflexTreeFromID[cld.ID].GeneticReflexID > 0 {
							// древний рефлекс
							oldReflexesIdArr = append(oldReflexesIdArr, ReflexTreeFromID[cld.ID].GeneticReflexID)
							isPrioritet=true
						}
					}
					*/
				case 2:
					/*
					if ReflexTreeFromID[cld.ID].GeneticReflexID > 0 {
						geneticReflexesIdArr = append(geneticReflexesIdArr, ReflexTreeFromID[cld.ID].GeneticReflexID)
						isPrioritet=true
					}
					*/
				}
				detectedActiveLastNodID = cld.ID
				level++
				// иначе игнорирует дубли findReflexesNodes(level, ost, &cld, isExactly)
				//!!!return
			}
		}
		findReflexesNodes(level, ost, &cld, isExactly)
	}
	/*
	if !isPrioritet && oldReflex>0  && oldReflexesIdArr==nil{
		oldReflexesIdArr = append(oldReflexesIdArr, oldReflex)
	}
	*/
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


////СТРОКА УСЛОВИЙ ДЛЯ безусловного РЕФЛЕКСА
func getCurrentConditionsStr()(string){

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
	as:=action_sensor.CheckCurActions()
	for i := 0; i < len(as); i++ {
			if i>0{out+=","}
			out+=strconv.Itoa(as[i])
	}

return out
}
////////////////////////////////////////////////////////



func checkIgnorOnly(oldReflexesIdArr []int,geneticReflexesIdArr []int)(bool) {
	var isIgnor = false
	if GeneticReflexes != nil {
		if len(oldReflexesIdArr) == 1 {
			gr := GeneticReflexes[oldReflexesIdArr[0]]
			if gr.ActionIDarr[0] == 9 {
				isIgnor = true
			}
		} else {
			if len(geneticReflexesIdArr) == 1 {
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

	// выявить ID парамктров гомеостаза как цели для улучшения в данных условиях
	//здесь, чтобы сразу получить veryActual и targetArrID для возврата
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
		for i := 0; i < len(condArr); i++ {
			act:=GeneticReflexes[condArr[i]]
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
	// выдать массив возможных действий чтобы выбрать одно из них, пока еще не испытанное
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