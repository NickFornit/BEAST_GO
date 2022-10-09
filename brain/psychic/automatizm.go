/* Автоматизмы, могут совершать внешние действия или внутренние произвольные действия.
К ветке дерева может быть прикреплено сколько угодно автоматизмов: GetMotorsAutomatizmListFromTreeId(branchID)
но только один из автоматизмов, прикрепленных к ветке, может иметь Belief=2 - проверенное собственное знание
   Автоматизмы могут быть и не привязаны к конкретной ветке дерева,
а быть привязаны к отдельным значениям AutomatizmNode:
   к ID образа действий с пульта ActivityID и тогда branchID начинается с 1000000,
сохраняются в карте AutomatizmIdFromActionId
   к ID фразы VerbalID  и тогда branchID начинается с 2000000,
сохраняются в карте AutomatizmIdFromPhraseId

	Если задается Belief=2, остальные Belief=2 становится Belief=0.
!!! ПОЭТОМУ ВСЕГДА нужно задавать с помощью setAutomatizmBelief(atmzm *Automatizm,belief int))

Если для прикрепленных к узлу дерева есть карта штатных AutomatizmBelief2FromTreeNodeId,
то для прикрепленных к образам нужны ФУНКЦИИ ПОЛУЧЕНИЯ ШТАТНОГО ДЛЯ ДАННОГО ОБРАЗА:
func GetAutomatizmBeliefFromActionId(activityID int)(*Automatizm){
func GetAutomatizmBeliefFromPhraseId(verbalID int)(*Automatizm){

Структура записи: id|BranchID|Usefulness||Sequence||NextID|Energy|Belief|Count
*/

package psychic

import (
	"BOT/lib"
	"strconv"
	"strings"
)

/////////////////////////////////

// инициализирующий блок - в порядке последовательности инициализаций
// вызывается из psychic.go
func automatizmInit(){

	loadAutomatizm()

	//res:=RumAutomatizm(AutomatizmFromIdArr[1])
	//if res{}

}
/////////////////////////////////


/* для разделения строки Sequence автоматизма на составляющие
типы действий:
1 Snn - перечень ID фраз через запятую и к ней: Tnn:23 - образ тон-настроения
2 Dnn - ID прогрмаммы действий, через запятую
3 Ann - последовательный запуск автоматизмов с id1,id2..
4 Mnn - внутренние произвольные действия с id1,id2...
5 Tnn - образ тон-настроения одна цифра == образ тона-настроения (как в func GetToneMoodID(  и func GetToneMoodFromImg()
*/
type ActsAutomatizm struct {
	Type int  // тип совершаемого действия
	Acts string // само действие
}
///////////////////////////////////////

type Automatizm struct {
	ID         int
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

	/* Цепочка последовательности реагирования,
	   включающие элементарные реакции и ID уже имеющихся цепочек автоматизмов
	   Последовательность - строка, с разделителем | в которой виды действий обозначаются
	   символом типа и последующим числом == ID данного вида реагирования:
	   Snn- перечень ID фраз через запятую,
НЕТ nn - ID программы ответа (фраза),
	   Dnn - ID прогрмаммы действий, через запятую
	   Ann - последовательный запуск автоматизмов с id1,id2...

	   Эти типы разделяются с помощью символа "|".
	      НАПРИМЕР:
	      Sequence="Snn:24243,1234,0,24234,11234|Ann:4"
	      Sequence="Dnn:24|Ann:4"
	      Sequence="Dnn:24,4,56" - несколько действий одновременно
	      Sequence="Snn:24243,1234,0,24234,11234|Dnn:24,4,56"
	Sequence="Tnn:23 - образ тон-настроения
	      фраза+действие с тоном-настроением, заготовленная фраза,уже имеющийся автоматизм:
	      Sequence="Snn:24243,1234,0,24234,11234|Tnn:23|Dnn:24|Ann:3"

	   Сразу не используется ActsAutomatizm struct для удобства записи и считывания файла
	*/
	Sequence string
	/* цепь может быть пройдена ментально, без выполнения автоматизмов, для этого не вызывается моторное выполнение а просто - проход цепочки с просмотром ее звеньев.
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
	/*
	   Уверенность в авторитарном автоматизме высока в период авторитарного обучения
	и сильно падает в период собственной инициативы, когда нужно на себе проверить,
	а даст ли такой автоматизм в самом деле обещанное улучшение.
	Только один из автоматизмов, прикрепленных к ветке, может иметь Belief=2 - проверенное собственное знание
	Если задается Belief=2, остальные Belief=2 становится Belief=0.
!!! ПОЭТОМУ ВСЕГДА нужно задавать с помощью setAutomatizmBelief(atmzm *Automatizm,belief int)
	*/
	Belief int // 0 - предположение, 1 - чужие сведения, 2 - проверенное собственное знание, 3 - для мозжечкового рефлекса
	/* В случае, если в результате автоматизма его Usefulness изменит знак, то
	Count обнулится, а при таком же знаке - увеличивается на 1.
	 */
	Count int //надежность: число использований с подтверждением (бес)полезности Usefulness
	/* какие ID гомео-параметров улучшает это действие
	по аналогии и функциональности с http://go/pages/terminal_actions.php
	 */
	GomeoIdSuccesArr []int
}
/////////////////////////////////////

// все, привязанные к узлу дерева или привязанные к id образа действия и к id фразы.
var AutomatizmFromIdArr=make(map[int]*Automatizm)

// ШТАТНЫЕ автоматизмы, прикрепленные к ID узла Дерева с Belief==2 т.е. ШТАТНЫЕ, выполняющиеся не раздумывая
// у узла может быть только один штатный автоматизм с Belief==2
var AutomatizmBelief2FromTreeNodeId=make(map[int]*Automatizm)
//привязанные к ID образа действий с пульта ActivityID и тогда их branchID начинается с 1000000
// среди привязанный к данному образуID может быть один штатный с Belief==2
var AutomatizmIdFromActionId=make(map[int] []*Automatizm)
//привязанные к ID фразы VerbalID и тогда их branchID начинается с 2000000
// среди привязанный к данной фразеID может быть один штатный с Belief==2
var AutomatizmIdFromPhraseId=make(map[int] []*Automatizm)

/* список удачных автоматизмов, относящихся к определенным условиям (привзяанных к определенной ветке Дерева)
В этом списке поле Usefulness >0
 */
var AutomatizmSuccessFromIdArr=make(map[int]*Automatizm)
///////////////////////////////////////////////////////////////





///////////////////////////////////////


//////////////////////////////////////////////////////
// список всех автоматизмов для ID узла Дерева
func GetMotorsAutomatizmListFromTreeId(nodeID int)([]*Automatizm){
	if nodeID==0{
		return nil
	}
	var mArr[] *Automatizm
	for _, a := range AutomatizmFromIdArr {
		if a.BranchID < 1000000 && a.BranchID == nodeID{
			mArr = append(mArr, a)
		}
	}
	return mArr
}
//////////////////////////////////////////////

/////////////////////////////////////
// создать новый автоматизм
var lastAutomatizmID=0
var NoWarningCreateShow=false // true - не выдавать сообщение о новом автоматизме
func createNewAutomatizmID(id int,BranchID int,Sequence string)(int,*Automatizm){
// автоматизмы могут быть неуникальными, т.е. даже с тождественными Sequence, это не имеет значения.
// к одной вентке могут быть прикреплены неограниченное число автоматизмов
	if id==0{
		lastAutomatizmID++
		id=lastAutomatizmID
	}else{
		//		newW.ID=id
		if lastAutomatizmID<id{
			lastAutomatizmID=id
		}
	}

	var node Automatizm
	node.ID = id
	node.BranchID = BranchID
	node.Energy=5
	node.Sequence = Sequence

	AutomatizmFromIdArr[id]=&node
	if BranchID>1000000 && BranchID<2000000{
		imgID:=BranchID-1000000
		AutomatizmIdFromActionId[imgID]=append(AutomatizmIdFromActionId[imgID],&node)
	}
	if BranchID>2000000{
		imgID:=BranchID-2000000
		AutomatizmIdFromPhraseId[imgID]=append(AutomatizmIdFromPhraseId[imgID],&node)
	}

	if !NoWarningCreateShow {
		lib.WritePultConsol("Создан новый автоматизм.")
	}
	return id,&node
}
/////////////////////////////////////////

// создать новый автоматизм с записю в файл
func CreateNewAutomatizm(BranchID int,Sequence string)(int,*Automatizm){
// BranchID может быть ==0 для мозжечковых рефлексов
	if len(Sequence)==0{
		return 0,nil
	}

	id,verb:=createNewAutomatizmID(0,BranchID,Sequence)

	SaveAutomatizm()

	return id,verb
}
/////////////////////////////////////////
// создать новый автоматизм без записи в файл
func CreateAutomatizm(BranchID int,Sequence string)(int,*Automatizm){
	// BranchID может быть ==0 для мозжечковых рефлексов
	if len(Sequence)==0{
		return 0,nil
	}

	id,verb:=createNewAutomatizmID(0,BranchID,Sequence)

	return id,verb
}
/////////////////////////////////////////

// СОХРАНИТЬ структура записи: id|BranchID|Usefulness||Sequence||NextID|Energy|Belief
//В случае отсуствия пусковых стимулов создается ID такого отсутсвия, пример такой записи: 2|||0|0|
func SaveAutomatizm(){
	var out=""
	for k, v := range AutomatizmFromIdArr {
		out+=strconv.Itoa(k)+"|"
		out+=strconv.Itoa(v.BranchID)+"|"
		out+=strconv.Itoa(v.Usefulness)+"||"
		out+=v.Sequence+"||"
		out+=strconv.Itoa(v.Usefulness)+"|"
		out+=strconv.Itoa(v.Energy)+"|"
		out+=strconv.Itoa(v.Belief)+"|"
		out+=strconv.Itoa(v.Count)+"|"
		for i := 0; i < len(v.GomeoIdSuccesArr); i++ {
			out+=strconv.Itoa(v.GomeoIdSuccesArr[i])+","
		}
		out+="\r\n"
	}
	lib.WriteFileContent(lib.GetMainPathExeFile()+"/memory_psy/automatizm_images.txt",out)
}
// ЗАГРУЗИТЬ структура записи: id|BranchID|Usefulness||Sequence||NextID|Energy|Belief
func loadAutomatizm(){
	NoWarningCreateShow=true
	AutomatizmFromIdArr=make(map[int]*Automatizm)
	strArr,_:=lib.ReadLines(lib.GetMainPathExeFile()+"/memory_psy/automatizm_images.txt")
	if strArr == nil{
		return
	}
	for n := 0; n < len(strArr); n++ {
		if len(strArr[n])==0{
			continue
		}
		main := strings.Split(strArr[n], "||")
		p := strings.Split(main[0], "|")
		id, _ := strconv.Atoi(p[0])
		BranchID, _ := strconv.Atoi(p[1])
		Usefulness, _ := strconv.Atoi(p[2])

		Sequence := main[1]

		p = strings.Split(main[2], "|")
		NextID,_:= strconv.Atoi(p[0])
		Energy, _ := strconv.Atoi(p[1])
		Belief, _ := strconv.Atoi(p[2])
		Count, _ := strconv.Atoi(p[3])
		s:=strings.Split(p[4], ",")
		var GomeoIdSuccesArr[] int
		for i := 0; i < len(s); i++ {
			if len(s[i])==0{
				continue
			}
			sp,_:=strconv.Atoi(s[i])
			GomeoIdSuccesArr=append(GomeoIdSuccesArr,sp)
		}
		_, a := createNewAutomatizmID(id, BranchID, Sequence)
		a.NextID = NextID
		a.Usefulness = Usefulness
		a.Energy = Energy
		a.Count = Count
		SetAutomatizmBelief(a, Belief)
	}
	NoWarningCreateShow=false
	return
}
///////////////////////////////////////






