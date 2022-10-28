/* Дерево понимания или дерево ментальных автоматизмов

В конечных узлах Дерева накапливаются ментальные автоматизмы.

*/

package psychic

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
Узлы всех уровней могут произвольно меняться на другие для переактивации Дерева.
 */
type UnderstandingNode struct { // узел дерева автоматизмов
	ID int
	Mood int // ощущаемое настроение PsyBaseMood: -1 Плохое настроение, 0 Нормальное, 1 - хорошее настроение
	EmotionID int // эмоция, может произвольно меняться
/* SituationID определяет основной контекст ситуации, определяемый при вызове активации дерева понимания.
Если этот контекст не задан в understandingSituation(situationImageID
то в getCurSituationImageID() по-началу выбирается наугад (для первого приближения) более важные из существующих,
но потом дерево понимания может переактивироваться с произвольным заданием контекста.
От этого параметра зависит в каком направлении пойдет информационный поиск решений,
если не будет запущен штатный автоматизм ветки (ориентировочные реакции).
Более частный, целевой контекст ситуации определяется следующим параметром PurposeID
 */
	SituationID int // ID объекта структуры понимания SituationImage, может произвольно меняться

/* ID образа ЖЕЛАЕМОЙ при данных условиях цели - PurposeImage,
который по-началу наследует PurposeGenetic, но может произвольно меняться,
в том числе после подсказки оператором:
в результате осмысления ответа оператора и запуска мент.автоматизма корректироваки цели
с перезапуском дерева понимания.
Для достижения этой общей цели в цепочках мент.автоматизмов определяются промежуточные цели так,
что каждый мент.автоматизм оценивается успешным при появлении состояния, соотвествующему данной промежуточной цели,
а конечное звено цепи должно стремиться к соотвествию PurposeID.
 */
	PurposeID int

	Children []UnderstandingNode // дочерние узлы (ветвление) НЕ АДРЕСА, А РЕАЛЬНЫЕ ОБЪЕКТЫ
	ParentID int     // ID родителя
	ParentNode *UnderstandingNode  // адрес родителя
}
var UnderstandingTree UnderstandingNode
var UnderstandingNodeFromID=make(map[int]*UnderstandingNode)
// последовательность узлов активной ветки
var ActiveBranchUnderstandingArr []int
////////////////////////////////////////////////

// если в результате ментальных процессов было действие, то нужно заблокировать обработку активации дерева моторных автоматизмов
var MentalReasonBlocing=false
///////////////////////////////////////////////////////



////////////////////////////////////////////////////////////////////////////////////////
/* попытка активации дерева ментальных автоматизмов

*/
var detectedActiveLastUnderstandingNodID=0
// нераспознанный остаток - НОВИЗНА
var CurrentUnderstandingTreeEnd []int
var currentUnderstandingStepCount=0
// массив id узлов активной ветки
var currentUnderstandingNodeID[] int

var saveSituationImageID=0

var currentMentalAutomatizmID=0


/* Активация дерева ментальных автоматизмов происходит из:
func afterTreeActivation() - при каждой активации automatism_tree.go
и если было действия без ответа в течении 20 пульсов, то understandingSituation вызывается из
func noAutovatizmResult()
т.е. оба деревав работают совместно при EvolushnStage > 3 и по каждой активации UnderstandingTree
добавляется эпизд памяти newEpisodeMemory()

Аналогично дереву моторных автоматзмов, после активации могут быть ориентировочные рефлексы привлечения внимания.

При вызове может быть определен situationImageID или проставлен 0 и тогда образ ситуации определяется в самой функции.

Если были совершены действия, то нужно выставлять MotorTerminalBlocking=true !!!
 */
func understandingSituation()(bool){
	MentalReasonBlocing=false // освободить обработку дерева моторных автоматизмов

	if EvolushnStage < 4 { // только со стадии развития 4
		return false
	}
	if PulsCount<4{// не активировать пока все не устаканится
		return false
	}
	// определить ID ситуации: настроение при посылке сообщения, нажатые кнопки и т.п.
	situationImageID:=getCurSituationImageID()
	if situationImageID<0{// нет выбранной ситуации
			return false
		}
	curBaseStateImage.SituationID=situationImageID

	saveSituationImageID=situationImageID
	ps:=getPurposeGenetic() // - тут уже сохраняется savePurposeGenetic

	savePurposeGenetic=ps

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

	// результат поиска:
	if detectedActiveLastUnderstandingNodID>0{
		// есть ли неучтенные условия?
		conditionsCount:=getConditionsCount(condArr)
		CurrentUnderstandingTreeEnd=condArr[currentUnderstandingStepCount:] // НОВИЗНА
		if currentUnderstandingStepCount<conditionsCount { // не пройдено до конца имеющихся условий
			// нарастить недостающее в ветке дерева
			detectedActiveLastUnderstandingNodID = formingUnderstandingBranch(detectedActiveLastUnderstandingNodID, currentUnderstandingStepCount+1, condArr)
	//мент.автоматизм может прикрепляться ТОЛЬКО к последнему узлу ветки - при полном понимании ситуации

			res:=consciousness(1,0)
			return res
		}
	}else{// вообще нет совпадений для данных условий
		// нарастить недостающее в ветке дерева
		detectedActiveLastUnderstandingNodID = formingUnderstandingBranch(detectedActiveLastUnderstandingNodID, currentUnderstandingStepCount, condArr)

		CurrentUnderstandingTreeEnd=condArr // все - новизна

		res:=consciousness(1,0)
		return res
	}

	return false
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
			continue
		}

		level++
		currentUnderstandingStepCount=level
		conditionUnderstandingFound(level,ost, &node.Children[n])
		return // раз совпало, то другие ветки не смотреть
	}

	return
}
////////////////////////////////////////////////////////