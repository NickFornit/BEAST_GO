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

// инициализирующий блок - в порядке последовательности инициализаций
// из psychic.go
func automatizmTreeInit(){

	loadAutomatizmTree()

}
/////////////////////////////////////////////////////////////

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

/*запрет показа карты при обновлении против паники типа "одновременная запись и считывание карты"
Использовать для всех операций записи узлов дерева
 */
var notAllowScanInTreeThisTime=false
////////////////////////////////////////////////








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



func automatizmTreeActivation()(int){
	if PulsCount<4{// не активировать пока все не устаканится
		return 0
	}
/* НУЖНО, просто новый ор.рефлекс будет ждать окончания периода isPeriodResultWaiting
	if isPeriodResultWaiting{// не активировать в период ожидания результатов действий!
		return 0
	}
 */
	detectedActiveLastNodID=0
	ActiveBranchNodeArr=nil
	CurrentAutomatizTreeEnd=nil
	currentStepCount=0

	// вытащить 3 уровня условий в виде ID их образов
	//Еще нет InformationEnvironment т.к. Дерево активруется ДО ор.рефлексов
	lev1:=gomeostas.CommonBadNormalWell

	bsIDarr:=gomeostas.GetCurContextActiveIDarr()
	lev2,_:=createNewBaseStyle(0,PsyBaseMood,bsIDarr)

ActID:=action_sensor.CheckCurActionsContext();//CheckCurActions()

	lev3,_:=createNewlastActivityID(0,ActID)// текущий образ сочетания действий с Пульта Activity

	var lev4=0
	var lev5=0
	var lev6=0
	if len(wordSensor.CurrentPhrasesIDarr)>0{
		PhraseID := wordSensor.CurrentPhrasesIDarr
		ToneID := wordSensor.DetectedTone
		MoodID := wordSensor.CurPultMood
		_,verb:= CreateVerbalImage(PhraseID, ToneID, MoodID)
		lev4= getToneMoodID(verb.ToneID, verb.MoodID)
		lev5=verb.SimbolID
		/* для дерева берется только первая фраза, остальные можно восстановить для сопоставлений из
		AutomatizmNode.VerbalID.PhraseID[]
		или из памяти о воспринятых фразах (Vernike_detector.go): var MemoryDetectedArr []MemoryDetected
		*/
		lev6=verb.PhraseID[0]
	}

	condArr:=getActiveConditionsArr(lev1, lev2, lev3, lev4, lev5, lev6)
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
	notAllowScanInTreeThisTime=true
// есть ли неучтенные условия?
		conditionsCount:=getConditionsCount(condArr)
		CurrentAutomatizTreeEnd=condArr[currentStepCount:] // НОВИЗНА
		if currentStepCount<conditionsCount { // не пройдено до конца имеющихся условий
		// нарастить недостающее в ветке дерева - всегда для orientation_1()
			detectedActiveLastNodID = formingBranch(detectedActiveLastNodID, currentStepCount, condArr)
				
			// автоматизма нет у недоделанной ветки
			automatizmID := orientation(0)
			return automatizmID // блокировка рефлексов, если automatizmID > 0

			// если нет неучтенных условий
		}else{// все условия пойдены, ветка существует,
			automatizmID := getAutomatizmFromNodeID(detectedActiveLastNodID)
			if automatizmID > 0 {//ориентировочный рефлекс 2
				// если нет срочности
				if !veryActualSituation {
					// проверить подходит ли автоматизм к текущим условиям, если нет, - режим нахождения альтернативы  - ориентировочный рефлекс 2
					automatizmID := orientation(automatizmID)
					return automatizmID // блокировка рефлексов
					
				}
				return automatizmID// выполнение штатного автоматизма
			}else{
				// автоматизма нет у нормальной ветки (условия не требовали срочного его создания)
				automatizmID := orientation(0)
				return automatizmID // блокировка рефлексов
			}
		}
	}else{// вообще нет совпадений для данных условий
// нарастить недостающее в ветке дерева - всегда для orientation_1()
		detectedActiveLastNodID = formingBranch(detectedActiveLastNodID, currentStepCount, condArr)
			
		// автоматизма нет у недоделанной ветки
		CurrentAutomatizTreeEnd=condArr // все - новизна
		automatizmID := orientation(0)
		return automatizmID // блокировка рефлексов

	}
notAllowScanInTreeThisTime=false
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
		var levID = 0
		switch level {
		case 0:
			levID = cld.BaseID
		case 1:
			levID = cld.EmotionID
		case 2:
			levID = cld.ActivityID
		case 3:
			levID = cld.SimbolID
		case 4:
			levID = cld.ToneMoodID
		case 5:
			levID = cld.VerbalID
		}
		if cond[0] == levID {
			detectedActiveLastNodID=cld.ID
			ActiveBranchNodeArr=append(ActiveBranchNodeArr,cld.ID)
		}

		level++
		currentStepCount=level
		conditionAutomatizmFound(level,ost, &node.Children[n])
		return // раз совпало, то другие ветки не смотреть
	}

	return
}
////////////////////////////////////////////////////////






