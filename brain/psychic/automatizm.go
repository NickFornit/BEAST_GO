/* Автоматизмы, могут совершать внешние действия или внутренние произвольные действия.
К ветке дерева может быть прикреплено сколько угодно автоматизмов: getAutomatizmFromTreeNodeIdArr(branchID)
но только один из автоматизмов, прикрепленных к ветке, может иметь Belief=2 - проверенное собственное знание
	Если задается Belief=2, остальные Belief=2 становится Belief=0.
!!! ПОЭТОМУ ВСЕГДА нужно задавать с помощью setAutomatizmBelief(atmzm *Automatizm,belief int))

Автоматизм обязательно привязан к определенной ветке дерева (условиях)
так что здесь нет каких-то универсальных действий, как в http://go/pages/terminal_actions.php

структура записи: id|BranchID|Usefulness||Sequence||NextID|Energy|Belief
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

}
/////////////////////////////////

// есть ли автоматизм для этой ветки
func getAutomatizmFromNodeID(nodeID int)(int){

	return 0
}
/////////////////////////////////////
/* для разделения строки Sequence автоматизма на составляющие
типы действий:
1 Snn - перечень ID фраз через запятую
2 Dnn - ID прогрмаммы действий, через запятую
3 Ann - последовательный запуск автоматизмов с id1,id2..
4 Mnn - внутренние произвольные действия с id1,id2...
*/
type ActsAutomatizm struct {
	Type int  // тип совершаемого действия
	Acts string // само действие
}
///////////////////////////////////////

type Automatizm struct {
	ID         int
	BranchID   int // id ветки дерева AutomatizmNode к кторой прикреплен автоматизм
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
	      фраза+действие, заготовленная фраза,уже имеющийся автоматизм:
	      Sequence="Snn:24243,1234,0,24234,11234|Fnn:4Dnn:24|Ann:3"

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

	// Энергичность действия или фразы
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


var AutomatizmFromIdArr=make(map[int]*Automatizm)

// автоматизмы, прикрепленные к ID узла Дерева получать по getAutomatizmFromTreeNodeIdArr(BranchID)
// автоматизмы, прикрепленные к ID узла Дерева с Belief==2 т.е. штатные, выполняющиеся не раздумывая
var AutomatizmBelief2FromTreeNodeId=make(map[int]*Automatizm)

/* список удачных автоматизмов, относящихся к определенным условиям (привзяанных к определенной ветке Дерева)
В этом списке поле Usefulness >0
 */
var AutomatizmSuccessFromIdArr=make(map[int]*Automatizm)

//////////////////////////////////////////////////////
// список всех автоматизмов для ID узла Дерева
func GetMotorsAutomatizmListFromTreeId(nodeID int)([]*Automatizm){
	var mArr[] *Automatizm
	for _, a := range AutomatizmFromIdArr {
		if a.BranchID == nodeID{
			mArr = append(mArr, a)
		}
	}
	return mArr
}
//////////////////////////////////////////////
/////////////////////////////////////
// создать новый автоматизм
var lastAutomatizmID=0
var noWarningCreateShow=false // true - не выдавать сообщение о новом автоматизме
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

	if !noWarningCreateShow {
		lib.WritePultConsol("Создан новый автоматизм.")
	}
	return id,&node
}
/////////////////////////////////////////
// создать новый автоматизм с записю
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
	noWarningCreateShow=true
	AutomatizmFromIdArr=make(map[int]*Automatizm)
	strArr,_:=lib.ReadLines(lib.GetMainPathExeFile()+"/memory_psy/automatizm_images.txt")
	if strArr == nil{
		return
	}
	for n := 0; n < len(strArr); n++ {
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
		a.Belief = Belief
		a.Count = Count
		if Belief==2 {//штатный, выполняющийся не раздумывая
			AutomatizmBelief2FromTreeNodeId[id]=a
		}
	}
	noWarningCreateShow=false
	return
}
///////////////////////////////////////////

/* разделить строку Sequence автоматизма на составляющие
типы действий:
1 Snn - перечень ID фраз через запятую
2 Dnn - ID прогрмаммы действий, через запятую
3 Ann - последовательный запуск автоматизмов с id1,id2..
4 Mnn - внутренние произвольные действия с id1,id2...
*/
func ParceAutomatizmSequence(Sequence string)([]ActsAutomatizm){
	var acts[] ActsAutomatizm

	sArr:=strings.Split(Sequence, "|")
	for i := 0; i < len(sArr); i++ {
		if len(sArr[i])==0{
			continue
		}
		var act ActsAutomatizm
		pArr:=strings.Split(sArr[i], ":")
		switch pArr[0]{
		case "Snn": act.Type=1
		//case "nn": act.Type=2
		case "Dnn": act.Type=2
		case "Ann": act.Type=3
		case "Mnn": act.Type=4
		}

		act.Acts = pArr[1] // строка действий (любого типа) через запятую
		acts = append(acts, act)
	}
	return acts
}
////////////////////////////////////////////////


/* получить массив wordID из Sequence автоматизма
 */
func GetWordArrFromSequence(sequence string)([]int){
	sArr:=strings.Split(sequence, "|")
	for i := 0; i < len(sArr); i++ {
		if len(sArr[i])==0{
			continue
		}
		pArr:=strings.Split(sArr[i], ":")
		if pArr[0]=="Snn"{
			sA:=strings.Split(pArr[1], ",")
			if len(sA)>0{
				var out[]int
				for a := 0; a < len(sA); a++ {
					wID,_:=strconv.Atoi(sA[a])
					out = append(out, wID)
				}
				return out
			}
		}
	}
	return nil
}

func GetAutomatizmSnn(ma *Automatizm)(string){
	sequence:=ma.Sequence // Sequence="Snn:24243,1234,0,24234,11234|Dnn:24,4,56"
	sArr:=strings.Split(sequence, "|")
	for i := 0; i < len(sArr); i++ {
		if len(sArr[i])==0{
			continue
		}
		pArr:=strings.Split(sArr[i], ":")
		if pArr[0]=="Snn"{
			if len(pArr[1])>0{
				return pArr[1]
			}
		}
	}
	return ""
}

func GetAutomatizmDnn(ma *Automatizm)(string){
	sequence:=ma.Sequence
	sArr:=strings.Split(sequence, "|")
	for i := 0; i < len(sArr); i++ {
		if len(sArr[i])==0{
			continue
		}
		pArr:=strings.Split(sArr[i], ":")
		if pArr[0]=="Dnn"{
			if len(pArr[1])>0{
				return pArr[1]
			}
		}
	}
	return ""
}
//////////////////////////////////////////////////


////////////////////////////////////////////////
//втоматизмы, прикрепленные к ID узла Дерева
func getAutomatizmFromTreeNodeIdArr(branchID int)([]*Automatizm){
	if branchID==0{
		return nil
	}
	var aArr []*Automatizm
	for _, v := range AutomatizmFromIdArr {
		if v.BranchID==branchID {
			aArr = append(aArr, v)
		}
	}
	return aArr
}
/////////////////////////////////////////////////
/*задать тип автоматизма Belief.
Только один из автоматизмов, прикрепленных к ветке, может иметь Belief=2 - проверенное собственное знание
Если задается Belief=2, остальные Belief=2 становится Belief=0.
ТАК ПРОСТО НЕЛЬЗЯ ЗАДАВАТЬ Belief=2: AutomatizmRunning.Belief=2
 */
func setAutomatizmBelief(atmzm *Automatizm,belief int){
	if atmzm==nil || atmzm.BranchID==0{
		return
	}
if belief==2{
	aArr:=getAutomatizmFromTreeNodeIdArr(atmzm.BranchID)
	if len(aArr)>1{
		for i := 0; i < len(aArr); i++ {
			if aArr[i]==atmzm && aArr[i].Belief==2{
				atmzm.Belief=0
				AutomatizmBelief2FromTreeNodeId[atmzm.BranchID]=nil
			}
		}
	}
	AutomatizmBelief2FromTreeNodeId[atmzm.BranchID]=atmzm
}
	atmzm.Belief=belief

}
/////////////////////////////////////////////////////









