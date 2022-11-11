/*  Дерево автоматизмов

Все начинается с psychic.go (atomatizmID:=automatizmTreeActivation()) -> func automatizmTreeActivation()

Это дерево активируется при:
1. Всегда при любых событиях с Пульта – так же как дерево рефлексов, но если к ветке привязан автоматизм,
то он выполняется преимущественно, блокируя рефлексы потому,
что уже было произвольностью преодолено действие рефлекса при выработке автоматизма.
Такой автоматизм обладает меткой успешности ==1. Успешность ==0 означает предположительный вариант действий,
а успешность < 0 – заблокированный вариант действия.
Так что к ветке может быть прикреплено множество неудачных и предположительных автоматизмов
и только один удачный. Более удачный результат переводит ранее удачный автоматизм в предполагаемые.
2. При произвольной активации отдельных условий.
Отсуствие подходящей для данных условий ветки дерева вызывает
Ориентировочный рефлекс привлечения внимания к активной ветке с осмыслением ситуации
и рассмотрением альтернатив действиям (4 уровня глубины рассмотрения).
При формировании нового предположительного действия создается новая ветка дерева и к ней прикрепляется автоматизм.
Т.е. новые условия не создают новой ветки, а тольно - новый автоматизм,
а пока нет автоматизма будет ориентировочный рефлекс.

У дерева фиксированных 6 уровней:
0 нулевой - основание
1 Базовое состояние - 3 вида
2 Эмоция
3 Активность с Пульта - образ ActivityFromIdArr=make(map[int]*Activity)
4 Образ контекста сообщения: сочетание Tone и Mood из структуры vrbal
5 Первый символ фразы
6 Фраза - VerbalID
До 6-го уровня - полный аналог условным рефлексам, только вместо сочетаний контекстов - эмоция.

Для оптимизации поиска по дереву перед узлом Verbal идет узел первого символа : var symbolsArr из word_tree.go

Формат записи:
ID|ParentNode|BaseID|EmotionID|ActivityID|ToneMoodID|SimbolID|VerbalID


Самоадаптация уровня Дерева автоматизмов
В результате действия автоматизма могут измениться условия и, значит,
будут запущены дерево рефлексов и опять - Дерево автоматизмов.
Возникает новая итерация адаптивности, возможно, с новым ориентировочным рефлексом второго типа.
Такой процесс может продолжаться до прихода к устойчивому состоянию.

*/

package psychic

import (
	"BOT/brain/action_sensor"
	"BOT/brain/gomeostas"
	wordSensor "BOT/brain/words_sensor"
)

// психика инициализирована
var StartPsichicNow=false

// инициализирующий блок - в порядке последовательности инициализаций
// из psychic.go
func automatizmTreeInit(){

	loadAutomatizmTree()
	if len(AutomatizmTree.Children)==0{// еще нет никаких веток
		// создать первые три ветки базовых состояний
		createBasicAutomatizmTree()
	}
	StartPsichicNow=true
}
/////////////////////////////////////////////////////////////


////////////////////////////////////////////

////// ДЕРЕВО автоматизмов имеет фиксированных 6 уровней (кроме базового нулевого)
type AutomatizmNode struct { // узел дерева автоматизмов
	ID int
	BaseID int // базовое состояние, !это еще не произвольно меняющееся PsyBaseMood
/* эмоция (type Emotion struct) Эмоция может произвольно меняться, независимо от базовых контекстов
т.е., к примеру, при BaseID Плохо может быть позитивное EmotionID
 */
	EmotionID int
	ActivityID int // образ сочетания действия с Пульта
/* образ контекста сообщения: сочетание Tone и Mood из структуры vrbal из automatism_tree_verbal_img.go
т.е. просто toneID+moodID - в виде строки, например: "922" = "Обычный, Хорошее"
дешифруется func getToneMoodStrFromID(id string)(string) 
*/
	ToneMoodID int  
	SimbolID int
	VerbalID int

	Children []AutomatizmNode // дочерние узлы (ветвление) НЕ АДРЕСА, А РЕАЛЬНЫЕ ОБЪЕКТЫ
	ParentID int     // ID родителя
	ParentNode *AutomatizmNode  // адрес родителя
}
var AutomatizmTree AutomatizmNode
var AutomatizmTreeFromID=make(map[int]*AutomatizmNode)
// последовательность узлов активной ветки
var ActiveBranchNodeArr []int

////////////////////////////////////////////////


// создать первые три ветки базовых состояний
func createBasicAutomatizmTree(){
	notAllowScanInTreeThisTime=true // запрет показа карты при обновлении

	createNewAutomatizmNode(&AutomatizmTree,0,1,0,0,0,0,0,false)
	createNewAutomatizmNode(&AutomatizmTree,0,2,0,0,0,0,0,false)
	createNewAutomatizmNode(&AutomatizmTree,0,3,0,0,0,0,0,false)

	if doWritingFile {SaveAutomatizmTree() }
	//SaveAutomatizmTree()
	notAllowScanInTreeThisTime=false // запрет показа карты при обновлении
	return
}
/////////////////////////////////////////////////////

// структура действий оператора при активации дерева автоматизмов
var curActiveActions ActionsImage
// структура образа текущего сосотояния
var curBaseStateImage BaseStateImage

////////////////////////////////////////////////////////////////////////////////////////
/* попытка активации дерева автоматизмов, если неудачно - начать искать вариант действий
Используется активная текущая информационная среда из psychic.go:
var PsyBaseID=0 // текущее базовое состояние, может быть произвольно изменено
var PsyEmotionImg *Emotion // текущая эмоция Emotion, может быть произвольно изменена
var PsyActionImg *Activity // текущий образ сочетания действий с Пульта Activity
var PsyVerbImg *Verbal // текущий образ фразы с Пульта Verbal
*/
var detectedActiveLastNodID=0
// нераспознанный остаток - НОВИЗНА
var CurrentAutomatizTreeEnd []int
var currentStepCount=0
var currentAutomatizmAfterTreeActivatedID=0 //! это  не обязательно штатный автоматизм ветки, а выбранный мягким алгоритмом


func automatizmTreeActivation()(int){
	if PulsCount<4{// не активировать пока все не устаканится
		return 0
	}
/* НУЖНО, просто новый ор.рефлекс будет ждать окончания периода LastRunAutomatizmPulsCount
	if LastRunAutomatizmPulsCount >0{// не активировать в период ожидания результатов действий!
		return 0
	}
 */

/* ТЕПЕПЕРЬ ВСЕГДА АКТИВИРОВАТЬ потому как и по изменению состояния формируются Правила.
Но нужно блокировать ор.рефлексы!
// не активировать дерево по изменению гомеостатуса во время ожидания ответа оператора
//  LastRunAutomatizmPulsCount устанавливается в RumAutomatizm(
if LastRunAutomatizmPulsCount > 0{
if !WasOperatorActiveted {
	return 0
}
 */

	detectedActiveLastNodID=0
	ActiveBranchNodeArr=nil
	CurrentAutomatizTreeEnd=nil
	currentStepCount=0
	currentAutomatizmAfterTreeActivatedID=0

	// вытащить 3 уровня условий в виде ID их образов
	//Еще нет InformationEnvironment т.к. Дерево активруется ДО ор.рефлексов
	lev1:=gomeostas.CommonBadNormalWell

	bsIDarr:=gomeostas.GetCurContextActiveIDarr()
	lev2,_:=createNewBaseStyle(0,bsIDarr,true)

	curBaseStateImage.Mood=lev1
	curBaseStateImage.EmotionID=lev2
	curBaseStateImage.SituationID=0 // будет определн при активации дерева понимания, может и не быть выбранной ситуации

ActID:=action_sensor.CheckCurActionsContext();//CheckCurActions()

	lev3,_:=createNewlastActivityID(0,ActID,true)// текущий образ сочетания действий с Пульта Activity
	curActiveActions.ActID=ActID// сохраняем для отзеркаливания действий оператора
	curActiveActions.PhraseID=nil
	curActiveActions.ToneID=0
	curActiveActions.MoodID=0

	var lev4=0
	var lev5=0
	var lev6=0
	if len(wordSensor.CurrentPhrasesIDarr)>0{
		PhraseID := wordSensor.CurrentPhrasesIDarr
		FirstSimbolID:=wordSensor.GetFirstSymbolFromPraseID(PhraseID)
		ToneID := wordSensor.DetectedTone
		MoodID := wordSensor.CurPultMood
		_,verb:= CreateVerbalImage(FirstSimbolID,PhraseID, ToneID, MoodID)
		lev4= GetToneMoodID(verb.ToneID, verb.MoodID)
		lev5=verb.SimbolID
		/* для дерева берется только первая фраза, остальные можно восстановить для сопоставлений из
		AutomatizmNode.VerbalID.PhraseID[]
		или из памяти о воспринятых фразах (Vernike_detector.go): var MemoryDetectedArr []MemoryDetected
		*/
		lev6=verb.PhraseID[0]
		// сохраняем для отзеркаливания действий оператора
		curActiveActions.PhraseID=PhraseID
		curActiveActions.ToneID=ToneID
		curActiveActions.MoodID=MoodID
	}

	condArr:=getActiveConditionsArr(lev1, lev2, lev3, lev4, lev5, lev6)
	notAllowScanInTreeThisTime=true // защелка от повтора во время обработки
	// основа дерева
	cnt := len(AutomatizmTree.Children)
	for n := 0; n < cnt; n++ {
		node := AutomatizmTree.Children[n]
		lev1 := node.BaseID
		if condArr[0] == lev1 {
			detectedActiveLastNodID=node.ID
			ost:=condArr[1:]
			if len(ost)==0{

			}

			conditionAutomatizmFound(1,ost, &node)

			break // другие ветки не смотреть
		}
	}


	// результат активации Дерева:
	if detectedActiveLastNodID>0{
// есть ли еще неучтенные, нулевые условия? т.е. просто показаь число ненулевых значений condArr
		conditionsCount:=getConditionsCount(condArr)
		CurrentAutomatizTreeEnd=condArr[currentStepCount:] // НОВИЗНА
		if currentStepCount<conditionsCount {              // не пройдено до конца имеющихся условий
			// нарастить недостающее в ветке дерева - всегда для orientation_1()
			//oldDetectedActiveLastNodID:=detectedActiveLastNodID
			detectedActiveLastNodID = formingBranch(detectedActiveLastNodID, currentStepCount+1, condArr)
		}

	}else{// вообще нет совпадений для данных условий
// нарастить недостающее в ветке дерева - всегда для orientation_1()
		detectedActiveLastNodID = formingBranch(detectedActiveLastNodID, currentStepCount+1, condArr)
			
		// автоматизма нет у недоделанной ветки
		CurrentAutomatizTreeEnd=condArr // все - новизна

	}

	if afterTreeActivation(){
		notAllowScanInTreeThisTime=false // снять блокировку
		return 1
	}
	notAllowScanInTreeThisTime=false // снять блокировку
	return 0
}
//////////////////////////////////////////////////////////////////

func conditionAutomatizmFound(level int,cond []int,node *AutomatizmNode){
	if cond==nil || len(cond)==0{
		return
	}

	ost:=cond[1:]

	for n := 0; n < len(node.Children); n++ {
		cld := node.Children[n]
		var val = 0
		switch level {
		case 0:
			val = cld.BaseID
		case 1:
			val = cld.EmotionID
		case 2:
			val = cld.ActivityID
		case 3:
			val = cld.ToneMoodID
		case 4:
			val = cld.SimbolID
		case 5:
			val = cld.VerbalID
		}
		if cond[0] == val {
			detectedActiveLastNodID=cld.ID
			ActiveBranchNodeArr=append(ActiveBranchNodeArr,cld.ID)
		}else {
			currentStepCount=level-1
			continue
		}

		level++
		currentStepCount=level
		conditionAutomatizmFound(level,ost, &node.Children[n])
		return // раз совпало, то другие ветки не смотреть
	}

	return
}
////////////////////////////////////////////////////////



/////////////////////////////////////////////////////////
/* реакция после активации ветки дерева

если нет никаких действий, то возвращает false, инчае - true для блокировки более низкоуровневого
 */
var onliOnceWasConditionsActiveted=false // т.к. опять может продолжиться изменение состояния в период ожидания
func afterTreeActivation()(bool){
	/* Нельзя здесь определять currentAutomatizmAfterTreeActivatedID перед if LastRunAutomatizmPulsCount >0{
	// ЕСТЬ ЛИ АВТОМАТИЗМ В ВЕТКЕ и болеее ранних? выбрать лучший автоматизм для сформированной ветки nodeID
	currentAutomatizmAfterTreeActivatedID = getAutomatizmFromNodeID(detectedActiveLastNodID)
	 */

// Был запущен моторный автоматизм (в том числе и ментальным автоматизмом)
if LastRunAutomatizmPulsCount >0{// обработка периода ожидания ответа оператора
	effect:=0
		// 	Контроль за изменением состояния, возвращает:
		//	lastCommonDiffValue - насколько изменилось общее состояние
		//  	lastBetterOrWorse - стали лучше или хуже: величина измнения от -10 через 0 до 10
		//  	gomeoParIdSuccesArr - стали лучше следующие г.параметры []int гоменостаза
		if WasOperatorActiveted { // оператор отреагировал
			lastCommonDiffValue,lastBetterOrWorse,gomeoParIdSuccesArr := wasChangingMoodCondition(2)
			effect=lastCommonDiffValue
			// обработать изменение состояния
			calcAutomatizmResult(lastCommonDiffValue,lastBetterOrWorse, gomeoParIdSuccesArr)

// по результатам обработки, но до очистки 	LastRunAutomatizmPulsCount и LastAutomatizmWeiting
			if EvolushnStage > 3 {
// Активировать Дерево Понимания: или запустить ментальный автоматизм или - ориентировочная реакция для осмысления
				understandingSituation(1) // нельзя здесь делать прерывание! после обработки ожидаемой реакции Оператора - следует реакция Beast
				// return true
			}
// закончить период ожидания после реакции оператора
			clinerAutomatizmRunning()
			WasConditionsActiveted = false // иначе сразу сработает fixRulesBaseStateImage после изменения состояния при действияъ
		}

	if !onliOnceWasConditionsActiveted {// только один раз во время периода ожидания
		onliOnceWasConditionsActiveted = true
		if WasConditionsActiveted { // изменились условия (не действия оператора)
			WasConditionsActiveted = false
			if EvolushnStage > 3 {
				lastCommonDiffValue, _, _ := wasChangingMoodCondition(2)
				// обработать изменение состояния - записать Правило типа BaseStateImage
				fixRulesBaseStateImage(lastCommonDiffValue)// здесь корректируется успешность автоматизма - как в calcAutomatizmResult
				// Активировать Дерево Понимания: или запустить ментальный автоматизм или - ориентировочная реакция для осмысления
				understandingSituation(1)

				// НЕ заканчивать период ожидания после переактивации по изменившимся условиям, но не запускать ор.рефлекс:
				return true
			}
		}
	}
	/* после периода ожидания
	   Учесть последствия запуска мот.автоматизма
	   и если нужно обдумать их в новом ментальном consciousness(2,currrentFromNextID)
	*/
	afterWaitingPeriod(effect)


//  после обработки ожидаемой реакции Оператора - следует реакция Beast
//		return true  поэтому нельзя здесь делать прерывание!
}// конец обработки ожидания ответа оператора

	// ЕСТЬ ЛИ АВТОМАТИЗМ В ВЕТКЕ и болеее ранних? выбрать лучший автоматизм для сформированной ветки nodeID
	currentAutomatizmAfterTreeActivatedID = getAutomatizmFromNodeID(detectedActiveLastNodID)

// ОБРАБОТКА ВНЕ ПЕРИОДА ОЖИДАНИЯ ОТВЕТА
	// всегда сначала активировать дерево понимания, результаты которого могут заблокировать все внизу
	if EvolushnStage > 3 {
// Активировать Дерево Понимания: или запустить ментальный автоматизм или - ориентировочная реакция для осмысления
		understandingSituation(1)

		//в результате ментальных процессов было совершено действие (моторное или ментальное)
		if MentalReasonBlocing{
			return true
		}
	}
	//////////////////////
	//более примитивное реагирование, EvolushnStage < 4
	if currentAutomatizmAfterTreeActivatedID > 0 { //ориентировочный рефлекс 2
		// проверить подходит ли автоматизм к текущим условиям, если нет, - режим нахождения альтернативы  - ориентировочный рефлекс 2
		atzm:=orientation(currentAutomatizmAfterTreeActivatedID)
		// если автоматизм прошел проверку, то он уже был запущен
		if atzm>0{// блокировка рефлексов, если automatizmID > 0
			return true
		}
	}else {
		// автоматизма нет у недоделанной ветки
		atzm:=orientation(0)
		if atzm>0{// блокировка рефлексов, если automatizmID > 0
			return true
		}
	}

	return false
}
//////////////////////////////////////////////////////////











