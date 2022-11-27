/* Автоматизмы, могут совершать внешние действия или внутренние произвольные действия.
К ветке дерева может быть прикреплено сколько угодно автоматизмов: GetMotorsAutomatizmListFromTreeId(branchID)
но только один из автоматизмов, прикрепленных к ветке, может иметь Belief=2 - проверенное собственное знание
Автоматизмы могут быть и не привязаны к конкретной ветке дерева, а быть привязаны к отдельным значениям AutomatizmNode:
- к ID образа действий с пульта ActivityID и тогда branchID начинается с 1000000,
сохраняются в карте AutomatizmIdFromActionId
- к ID фразы VerbalID  и тогда branchID начинается с 2000000,
сохраняются в карте AutomatizmIdFromPhraseId

Если задается Belief=2, остальные Belief=2 становится Belief=0.
!!! ПОЭТОМУ ВСЕГДА нужно задавать с помощью setAutomatizmBelief(atmzm *Automatizm,belief int))

Если для прикрепленных к узлу дерева есть карта штатных AutomatizmBelief2FromTreeNodeId,
то для прикрепленных к образам нужны ФУНКЦИИ ПОЛУЧЕНИЯ ШТАТНОГО ДЛЯ ДАННОГО ОБРАЗА:
func GetAutomatizmBeliefFromActionId(activityID int)(*Automatizm){
func GetAutomatizmBeliefFromPhraseId(verbalID int)(*Automatizm){

Структура записи: id|BranchID|Usefulness|ActionsImageID|NextID|Energy|Belief|Count
*/

package psychic

import (
	"BOT/lib"
	"strconv"
	"strings"
)

// инициализирующий блок - в порядке последовательности инициализаций
// вызывается из psychic.go
func automatizmInit(){
	loadAutomatizm()
	//res:=RumAutomatizm(AutomatizmFromIdArr[1])
	//if res{}
}

type Automatizm struct {
	ID	int
/* id объекта к кторому привязан автоматизм:
(он может быть  привязан к узлу дерева, к фразе (AutomatizmIdFromPhraseId) или действиям (AutomatizmIdFromActionId)
т.е. втоматизмы могут быть и не привязаны к конкретному узлу ветки дерева,
а быть привязаны к отдельным значениям AutomatizmNode:
   к ID образа действий с пульта ActivityID и тогда branchID начинается с 1000000,
сохраняются в карте AutomatizmIdFromActionId
   к ID фразы VerbalID  и тогда branchID начинается с 2000000,
сохраняются в карте AutomatizmIdFromPhraseId
*/
	BranchID   int 
	Usefulness int // (БЕС)ПОЛЕЗНОСТЬ: -10 вред 0 +10 +n польза

	// образ действий
	ActionsImageID int

	/* Следующий автоматизм в цепочке исполнения. Цепь может быть пройдена ментально, без выполнения автоматизмов, для этого не вызывается моторное выполнение а просто - проход цепочки с просмотром ее звеньев.
	или цепь может быть прервана осознанно
	или пройдена при ее пошаговом отслеживании StepByStepAutomatizm
	и пошаговом запуске: allowNextAutomatizm(automatizm.NextID):
	Бот смотрит, выполнить ли следующий шаг, добавить ли рефлекс мозжечка.
	*/
	NextID int
/* Энергичность действия или фразы.
   Но т.к. автоматизм может использоваться в разных случаях,
   лучше для этих конкретных случаев использования уточнять энергичность
   с помощью мозжечковых рефлексов.
 */
	Energy int // от 1 до 10, по умолчанию = 5
	/* Уверенность в авторитарном автоматизме высока в период авторитарного обучения
	и сильно падает в период собственной инициативы, когда нужно на себе проверить,
	а даст ли такой автоматизм в самом деле обещанное улучшение.
	Только один из автоматизмов, прикрепленных к ветке, может иметь Belief=2 - проверенное собственное знание
	Если задается Belief=2, остальные Belief=2 становится Belief=0.
	!!! ПОЭТОМУ ВСЕГДА нужно задавать с помощью setAutomatizmBelief(atmzm *Automatizm,belief int)
	*/
	Belief int // 0 - предположение, 1 - чужие сведения, 2 - проверенное собственное знание
	/* В случае, если в результате автоматизма его Usefulness изменит знак, то
	Count обнулится, а при таком же знаке - увеличивается на 1.
	*/
	Count int // надежность: число использований с подтверждением (бес)полезности Usefulness
	/* какие ID гомео-параметров улучшает это действие
	по аналогии и функциональности с http://go/pages/terminal_actions.php
	*/
	GomeoIdSuccesArr []int
}

// все, привязанные к узлу дерева или привязанные к id образа действия и к id фразы.
var AutomatizmFromIdArr = make(map[int]*Automatizm)

// ШТАТНЫЕ автоматизмы, прикрепленные к ID узла Дерева с Belief==2 т.е. ШТАТНЫЕ, выполняющиеся не раздумывая
// у узла может быть только один штатный автоматизм с Belief==2
var AutomatizmBelief2FromTreeNodeId = make(map[int]*Automatizm)
//привязанные к ID образа действий с пульта ActivityID и тогда их branchID начинается с 1000000
// среди привязанный к данному образуID может быть один штатный с Belief==2
var AutomatizmIdFromActionId = make(map[int] []*Automatizm)
//привязанные к ID фразы VerbalID и тогда их branchID начинается с 2000000
// среди привязанных к данной фразеID (неважны предыдущие условия) может быть один штатный с Belief==2
var AutomatizmIdFromPhraseId = make(map[int] []*Automatizm)

/* список удачных автоматизмов, относящихся к определенным условиям (привзяанных к определенной ветке Дерева)
В этом списке поле Usefulness >0
 */
var AutomatizmSuccessFromIdArr = make(map[int]*Automatizm)

// GetMotorsAutomatizmListFromTreeId(nodeID int) список всех автоматизмов для ID узла Дерева
// ExistsAutomatizmForThisNodeID(nodeID int) есть ли штатный автоматизм (с Belief==2), привязанные к узлу дерева
// GetBelief2AutomatizmListFromTreeId(nodeID int) штатный, невредный автоматизм, привязанный к ветке


var lastAutomatizmID = 0 // ID последнего созданного автоматизма
var NoWarningCreateShow = false // true - не выдавать сообщение о новом автоматизме

/* создать новый автоматизм
checkLevel - глубина проверки на идентичность: 0 - нет проверки, 1 - поверхностная, 2 - полная
 */
func createNewAutomatizmID(id int,BranchID int,ActionsImageID int,CheckUnicum bool)(int,*Automatizm) {
/* Автоматизмы уникальны по сочетанию BranchID и ActionsImageID.
	При попытке создать с таким же сочетанием возвращается уже созданный.
 к одной вентке могут быть прикреплены неограниченное число автоматизмов
 */
	if CheckUnicum {
		oldID, oldVal := checkUnicumMotorsAutomatizm(BranchID, ActionsImageID)
		if oldVal != nil {
			return oldID, oldVal
		}
	}

	if id == 0 {
		lastAutomatizmID++
		id=lastAutomatizmID
	} else {
		if lastAutomatizmID < id {
			lastAutomatizmID = id
		}
	}

	var node Automatizm
	node.ID = id
	node.BranchID = BranchID
	node.Energy = 5
	node.ActionsImageID = ActionsImageID

	AutomatizmFromIdArr[id] = &node
	if BranchID > 1000000 && BranchID < 2000000 {
		imgID := BranchID - 1000000
		AutomatizmIdFromActionId[imgID] = append(AutomatizmIdFromActionId[imgID], &node)
	}
	if BranchID > 2000000 {
		imgID := BranchID - 2000000
		AutomatizmIdFromPhraseId[imgID] = append(AutomatizmIdFromPhraseId[imgID], &node)
	}

	if !NoWarningCreateShow {
		lib.WritePultConsol("Создан новый автоматизм.")
	}
	return id, &node
}
//////////////////////////////////////////
/* Автоматизмы уникальны по сочетанию BranchID и ActionsImageID.
Функцию можно использовать для выборки автоматизма с заданными BranchID и ActionsImageID
*/
func checkUnicumMotorsAutomatizm(BranchID int,ActionsImageID int)(int,*Automatizm){
	for id, v := range AutomatizmFromIdArr {
		if BranchID != v.BranchID || ActionsImageID != v.ActionsImageID {
			continue
		}
		return id,v
	}
	return 0,nil
}
////////////////////////////////////////////

// создать новый автоматизм 
func CreateNewAutomatizm(BranchID int, ActionsImageID int)(int, *Automatizm) {
	// BranchID может быть ==0 для мозжечковых рефлексов
	if ActionsImageID == 0 { return 0, nil }

	id, verb := createNewAutomatizmID(0, BranchID, ActionsImageID,true)
	
	if doWritingFile {SaveAutomatizm() }

	return id, verb
}

// создать новый автоматизм без записи в файл
func CreateAtutomatizmNoSaveFile(BranchID int, ActionsImageID int)(int, *Automatizm) {
	// BranchID может быть ==0 для мозжечковых рефлексов
	if ActionsImageID == 0 {	return 0, nil }

	id, verb := createNewAutomatizmID(0, BranchID, ActionsImageID,true)

	return id,verb
}

// СОХРАНИТЬ структура записи: id|BranchID|Usefulness|ActionsImageID|NextID|Energy|Belief
// В случае отсуствия пусковых стимулов создается ID такого отсутсвия, пример такой записи: 2|||0|0|
func SaveAutomatizm() {
	var out = ""

	for k, v := range AutomatizmFromIdArr {
		out += strconv.Itoa(k) + "|"
		out += strconv.Itoa(v.BranchID) + "|"
		out += strconv.Itoa(v.Usefulness) + "|"
		out += strconv.Itoa(v.ActionsImageID) + "|"
		out += strconv.Itoa(v.NextID) + "|"
		out += strconv.Itoa(v.Energy) + "|"
		out += strconv.Itoa(v.Belief) + "|"
		out += strconv.Itoa(v.Count) + "|"
		for i := 0; i < len(v.GomeoIdSuccesArr); i++ {
			out += strconv.Itoa(v.GomeoIdSuccesArr[i]) + ","
		}
		out += "\r\n"
	}
	lib.WriteFileContent(lib.GetMainPathExeFile() + "/memory_psy/automatizm_images.txt", out)
}

// ЗАГРУЗИТЬ структура записи: id|BranchID|Usefulness||ActionsImageID||NextID|Energy|Belief
func loadAutomatizm() {
	NoWarningCreateShow = true
	AutomatizmFromIdArr = make(map[int]*Automatizm)

	strArr, _ := lib.ReadLines(lib.GetMainPathExeFile() + "/memory_psy/automatizm_images.txt")
	if strArr == nil { return	}
	for n := 0; n < len(strArr); n++ {
		if len(strArr[n]) == 0 { continue	}
		p := strings.Split(strArr[n], "|")
		id, _ := strconv.Atoi(p[0])
		BranchID, _ := strconv.Atoi(p[1])
		Usefulness, _ := strconv.Atoi(p[2])
		ActionsImageID, _ := strconv.Atoi(p[3])
		NextID, _ := strconv.Atoi(p[4])
		Energy, _ := strconv.Atoi(p[5])
		Belief, _ := strconv.Atoi(p[6])
		Count, _ := strconv.Atoi(p[7])
		s := strings.Split(p[4], ",")
		var GomeoIdSuccesArr[] int
		for i := 0; i < len(s); i++ {
			if len(s[i]) == 0 { continue }
			sp, _ := strconv.Atoi(s[i])
			GomeoIdSuccesArr = append(GomeoIdSuccesArr, sp)
		}
var saveDoWritingFile= doWritingFile; doWritingFile =false
		_, a := createNewAutomatizmID(id, BranchID, ActionsImageID,false)// без проверки на уникальность
		a.NextID = NextID
		a.Usefulness = Usefulness
		a.Energy = Energy
		a.Count = Count
		SetAutomatizmBelief(a, Belief)
doWritingFile =saveDoWritingFile
	}
	NoWarningCreateShow = false
	return
}
/////////////////////////////////////////////////////////////