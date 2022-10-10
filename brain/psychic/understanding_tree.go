/* Дерево понимания или дерево ментальных автоматизмов

В конечных узлах Дерева накапливаются ментальные автоматизмы.

*/

package psychic

import (
	"BOT/lib"
)
///////////////////////////////

// инициализирующий блок - в порядке последовательности инициализаций
// из psychic.go
func UnderstandingTreeInit(){

	loadPurposeImageFromIdArr()
	loadUnderstandingTree()
	if len(UnderstandingTree.Children)==0{// еще нет никаких веток
		// создать первые три ветки базовых состояний
		createBasicUnderstandingTree()
	}
}
/////////////////////////////////////////////////////////////

/* ДЕРЕВО понимания или Дерево ментальных автоматизмов.
Имеет фиксированных 4 уровней (кроме базового нулевого)
формат записи: ID|Mood|EmotionID|SituationID|PurposeI
Узлы всех уровней могут произвольно меняться на другие с переактивацией Дерева.
 */
type UnderstandingNode struct { // узел дерева автоматизмов
	ID int
	Mood int // ощущаемое настроение PsyBaseMood, может произвольно меняться
	EmotionID int // эмоция, может произвольно меняться
	SituationID int // ID объекта структуры понимания SituationImage, может произвольно меняться
	PurposeID int // ID образа PurposeImage  который наследует PurposeGenetic, может произвольно меняться

	Children []UnderstandingNode // дочерние узлы (ветвление) НЕ АДРЕСА, А РЕАЛЬНЫЕ ОБЪЕКТЫ
	ParentID int     // ID родителя
	ParentNode *UnderstandingNode  // адрес родителя
}
var UnderstandingTree UnderstandingNode
var UnderstandingNodeFromID=make(map[int]*UnderstandingNode)
// последовательность узлов активной ветки
var ActiveBranchUnderstandingArr []int
////////////////////////////////////////////////



///////////////////////////////////////////////////////



////////////////////////////////////////////////////////////////////////////////////////
/* попытка активации дерева автоматизмов, если неудачно - начать искать вариант действий
Используется активная текущая информационная среда из psychic.go:
var PsyBaseID=0 // текущее базовое состояние, может быть произвольно изменено
var PsyEmotionImg *Emotion // текущая эмоция Emotion, может быть произвольно изменена
var PsyActionImg *Activity // текущий образ сочетания действий с Пульта Activity
var PsyVerbImg *Verbal // текущий образ фразы с Пульта Verbal
*/
var detectedActiveLastUnderstandingNodID=0
// нераспознанный остаток - НОВИЗНА
var CurrentUnderstandingTreeEnd []int
var currentUnderstandingStepCount=0
// массив id узлов активной ветки
var currentUnderstandingNodeID[] int


/* вызывается из:
automatizm_result.go - в calcAutomatizmResult(
на 4-й стадии и позже - при любых действиях с пульта, не сопровождающихся срабатыванием автоматизма

Если были совершены действия, то нужно выставлять isReflexesActionBloking=true !!!
 */
func understandingSituation(situationImageID int,ps *PurposeGenetic)(bool){
	if EvolushnStage < 4 { // только со стадии развития 4
		return false
	}
	if PulsCount<4{// не активировать пока все не устаканится
		return false
	}
	/* может думать в это время!
	if LastRunAutomatizmPulsCount >0{// не активировать в период ожидания результатов действий!
		return false
	}
	*/
	detectedActiveLastUnderstandingNodID=0
	ActiveBranchNodeArr=nil
	CurrentUnderstandingTreeEnd=nil
	currentUnderstandingStepCount=0
	currentUnderstandingNodeID=nil

	// вытащить 3 уровня условий в виде ID их образов
	var lev1=PsyBaseMood
	var lev2=0
	if CurrentInformationEnvironment.PsyActionImg!=nil{
		lev2=CurrentInformationEnvironment.PsyEmotionImg.ID
	}

	var lev3=situationImageID
	var lev4,_=createPurposeImageID(0,ps.veryActual,ps.targetID,ps.actionID.ID)


	condArr:=getUnderstandingActiveConditionsArr(lev1, lev2, lev3, lev4)
	// основа дерева
	cnt := len(UnderstandingTree.Children)
	for n := 0; n < cnt; n++ {
		node := UnderstandingTree.Children[n]
		lev1 := node.Mood
		if condArr[0] == lev1 {
			detectedActiveLastUnderstandingNodID=node.ID
			ost:=condArr[1:]
			if len(ost)==0{

			}

			conditionUnderstandingFound(1,ost, &node)

			break // другие ветки не смотреть
		}
	}

	lib.WritePultConsol("Активировалось Дерево понимания.")

	// результат поиска:
	if detectedActiveLastUnderstandingNodID>0{
		// есть ли неучтенные условия?
		conditionsCount:=getConditionsCount(condArr)
		CurrentUnderstandingTreeEnd=condArr[currentUnderstandingStepCount:] // НОВИЗНА
		if currentUnderstandingStepCount<conditionsCount { // не пройдено до конца имеющихся условий
			// нарастить недостающее в ветке дерева
			detectedActiveLastUnderstandingNodID = formingUnderstandingBranch(detectedActiveLastUnderstandingNodID, currentUnderstandingStepCount+1, condArr)

			// Ориентировочный рефлекс осознания ситуации - частичная новизна условий
			res:=orientationConsciousness(1)
			/*
			// автоматизма нет у недоделанной ветки
			automatizm := orientation_1()
			if automatizm !=nil  {
					automatizm.BranchID = detectedActiveLastUnderstandingNodID
				// сразу запустить МОТОРНОЕ действие - в  orientation_1()

			}
			 */
			newEpisodeMemory()
			return res
		}else{// все условия пойдены,. ветка существует,
			//МЕНТ.АВТОМАТИЗМ может и не БЫТЬ
			// Ориентировочный рефлекс осознания ситуации - только новизна ситуации
			res:=orientationConsciousness(2)
			/*
			automatizmID := getAutomatizmFromNodeID(detectedActiveLastUnderstandingNodID)
			if automatizmID > 0 {//ориентировочный рефлекс 2
				// проверить подходит ли автоматизм к текущим условиям, если нет, - режим нахождения альтернативы  - ориентировочный рефлекс 2
				automatizm := orientation_2(automatizmID)
				if automatizm !=nil {
						automatizm.BranchID=nodeID
					// сразу попробовать выполнить этот автоматизм, т.к. если вернуло автоматизм, значит хочет попробовать
					RumAutomatizmID(automatizm.ID)
				}

			}
			*/
			newEpisodeMemory()
			return res
		}
	}else{// вообще нет совпадений для данных условий
		// Ориентировочный рефлекс осознания ситуации полная новизна условий
		res:=orientationConsciousness(0)
		//нарастить недостающую ветку
		// нарастить недостающее в ветке дерева
		detectedActiveLastUnderstandingNodID = formingUnderstandingBranch(detectedActiveLastUnderstandingNodID, currentUnderstandingStepCount, condArr)

		// автоматизма нет у недоделанной ветки
		CurrentUnderstandingTreeEnd=condArr // все - новизна
		automatizm := orientation_1()
		if automatizm !=nil {
			automatizm.BranchID = detectedActiveLastUnderstandingNodID
			// сразу запустить МОТОРНОЕ действие - в  orientation_1()
			return true
		}
		newEpisodeMemory()
		return res
	}

	return false
}
/////////// НОВЫЙ ЭПИЗОД ПАМЯТИ
func newEpisodeMemory(){
	// выдать массив ID узлов ветки по заданному ID узла
	currentUnderstandingNodeID:=getBrangeUnderstandingNodeIdArr(detectedActiveLastUnderstandingNodID)
	// новый эпизод памяти
	createEpisodeMemoryFrame(0,LifeTime,PsyMood,PsyBaseMood,currentUnderstandingNodeID,0)
}
//////////////////////////////////////////////////////////////////

func conditionUnderstandingFound(level int,cond []int,node *UnderstandingNode){
	if cond==nil || len(cond)==0{
		return
	}

	ost:=cond[1:]

	for n := 0; n < len(node.Children); n++ {
		cld := node.Children[n]
		var levID = 0
		switch level {
		case 0:
			levID = cld.Mood
		case 1:
			levID = cld.EmotionID
		case 2:
			levID = cld.SituationID
		case 3:
			levID = cld.PurposeID
		}
		if cond[0] == levID {
			detectedActiveLastUnderstandingNodID=cld.ID
			ActiveBranchNodeArr=append(ActiveBranchNodeArr,cld.ID)
		}else {
			currentUnderstandingStepCount=level-1
			return
		}

		level++
		currentUnderstandingStepCount=level
		conditionUnderstandingFound(level,ost, &node.Children[n])
		return // раз совпало, то другие ветки не смотреть
	}

	return
}
////////////////////////////////////////////////////////