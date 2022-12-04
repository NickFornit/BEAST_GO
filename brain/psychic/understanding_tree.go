/* Дерево понимания или дерево ментальных автоматизмов

В конечных узлах Дерева накапливаются ментальные автоматизмы.
формат записи: ID|ParentNode|Mood|EmotionID|SituationID|PurposeID
*/

package psychic

import "BOT/brain/gomeostas"

///////////////////////////////

// инициализирующий блок - в порядке последовательности инициализаций
// из psychic.go
func UnderstandingTreeInit(){
	if EvolushnStage < 4 { // только со стадии развития 4
		return
	}
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
формат записи: ID|ParentNode|Mood|EmotionID|SituationID|PurposeID
Узлы всех уровней могут произвольно меняться на другие для переактивации Дерева.

Дерево может переактивароваться при срабатывании мент.автоматизмов с действиями
MentalActionsImages.activateBaseID и MentalActionsImages.activateEmotion
в mental_automatizm_actions.go RunMentalAutomatizm(
 */
type UnderstandingNode struct { // узел дерева понимания
	ID int
	//Mood = PsyBaseMood: -1 Плохое настроение, 0 Нормальное, 1 - хорошее настроение
	Mood int // ощущаемое настроение МОЖЕТ ПРОИЗВОЛЬНО МЕНЯТЬСЯ
	EmotionID int // эмоция, МОЖЕТ ПРОИЗВОЛЬНО МЕНЯТЬСЯ
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
//var ActiveBranchUnderstandingArr []int
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
// массив узлов активной ветки   currentUnderstandingActivedNodes=getcurrentUnderstandingActivedNodes(lastID)
var currentUnderstandingActivedNodes[]*UnderstandingNode // вначале конечный узел

var saveSituationImageID=0

var currentMentalAutomatizmID=0


//текущие образы  гомеостатической этиологии, колторые могут быть произвольно перекрыты ментальными образами
var newMoodID=0
var newEmotionID=0
var newPurposeID=0
//сохраненные образы гомеостатической этиологии, колторые могут быть произвольно перекрыты ментальными образами
var preMoodID=0
var preEmotionID=0
var prePurposeID=0


/* Активация дерева ментальных автоматизмов происходит из:
func afterTreeActivation() - при каждой активации automatism_tree.go
и если было действия без ответа в течении 20 пульсов, то understandingSituation вызывается из
func noAutovatizmResult()
т.е. оба деревав работают совместно при EvolushnStage > 3 и по каждой активации UnderstandingTree
добавляется эпизд памяти newEpisodeMemory()

Аналогично дереву моторных автоматзмов, после активации могут быть ориентировочные рефлексы привлечения внимания.

При вызове может быть определен situationImageID или проставлен 0 и тогда образ ситуации определяется в самой функции.

Если были совершены действия, то нужно выставлять MotorTerminalBlocking=true !!!

activationType =1 - объективная активация
activationType =2 - произвольная переактивация
 */
func understandingSituation(activationType int)(bool){
	MentalReasonBlocing=false // освободить обработку дерева моторных автоматизмов

	if EvolushnStage < 4 { // только со стадии развития 4
		return false
	}
	if PulsCount<4{// не активировать пока все не устаканится
		return false
	}
	// определить ID ситуации: настроение при посылке сообщения, нажатые кнопки и т.п.
	situationImageID:=getCurSituationImageID()
	if situationImageID == 0{// нет выбранной ситуации
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
	currentUnderstandingActivedNodes=nil

	// текушие гомео-зависимые параметы
	newMoodID=PsyBaseMood

	bsIDarr:=gomeostas.GetCurContextActiveIDarr()
	newEmotionID,_=createNewBaseStyle(0,bsIDarr,true)

	newPurposeID,_ = createPurposeImageID(0, ps.veryActual, ps.targetID, ps.actionID.ID,true)
	////////////////////////////////


	/* не просрочены ли произвольно активированные параметры
	Держатся на время, пока не изменятся генетически определенные соотвествующие параметры или
	если активация была в данном пульсе
	 */
	if mentalMoodVolitionPulsCount!=PulsCount && preMoodID != newMoodID{
		mentalMoodVolitionID=0
	}
	if mentalEmotionVolitionPulsCount!=PulsCount && preEmotionID != newEmotionID{
		mentalEmotionVolitionID=0
	}
	if mentalPurposeImagePulsCount!=PulsCount && prePurposeID != newPurposeID{
		mentalPurposeImageID=0
	}

	// сохранять только изменившиеся значения
	if preMoodID != newMoodID {
		preMoodID = newMoodID
	}
	if preEmotionID != newEmotionID {
		preEmotionID = newEmotionID
	}
	if prePurposeID != newPurposeID {
		prePurposeID = newPurposeID
	}
	////////////////////////////////

	// 3 уровня условий в виде ID их образов

	var lev1=newMoodID // PsyBaseMood: -1 Плохое настроение, 0 Нормальное, 1 - хорошее настроение
	if mentalMoodVolitionID >0 {
		lev1=mentalMoodVolitionID
	}

	var lev2=newEmotionID
	if mentalEmotionVolitionID >0 {
		lev2=mentalEmotionVolitionID
	}

	var lev3=situationImageID

	var lev4=newPurposeID
	if mentalPurposeImageID >0 {
		lev4=mentalPurposeImageID
	}

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
		}
	}else{// вообще нет совпадений для данных условий
		// нарастить недостающее в ветке дерева
		detectedActiveLastUnderstandingNodID = formingUnderstandingBranch(detectedActiveLastUnderstandingNodID, currentUnderstandingStepCount, condArr)
		CurrentUnderstandingTreeEnd=condArr // все - новизна
	}
	// все узлы активной ветки, в начале - конечный узел
	currentUnderstandingActivedNodes=getcurrentUnderstandingActivedNodes(detectedActiveLastUnderstandingNodID)

	/* объективный запуск consciousness - по кативации дерева автоматизмов
	ментальный запуск в случае произвольной переактивации дерева понимания или цикле осмысления (5-я ступень)
	 */
	res:=consciousness(activationType,0)
	if isActivationType2{// был ментальный перезапус
		return true// заблокировать все низкоуровневое
	}
/* вернуть res consciousness после конца рекурсивного цикла ментального осмысления:
	если был запущен моторный автоматизм, то true
 */
	return res
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


