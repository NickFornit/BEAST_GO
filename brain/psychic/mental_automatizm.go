/* Ментальные (умственные) автоматизмы мышления.
мент.автоматизм может прикрепляться ТОЛЬКО к последнему узлу ветки - при полном понимании ситуации
К узлу может быть прикреплено сколько угодно ментальных автоматизмов, но только один из них - штатный (Belief==2)
Штатный устанавливается ТОЛЬКО функцией SetMentalAutomatizmBelief !
*/

package psychic

import (
	"BOT/lib"
	"strconv"
	"strings"
)
////////////////////////////////////////

func mentalAutomatizmInit(){
	loadMentalAutomatizm()
}
///////////////////////////////////////


///////////////////////////////////////
type MentalAutomatizm struct {
	ID int
	BranchID int
	Usefulness int // (БЕС)ПОЛЕЗНОСТЬ: -10 вред 0 +10 +n польза

	/* Цепочка последовательности реагирования,
Mnn:ID - выполнить ментальную функцию с ID
Ann:ID - выполнить моторный автоматизм с ID
	*/
	Sequence string

	/* Следующий автоматизм в цепочке исполнения. Цепь может быть пройдена ментально, без выполнения автоматизмов, для этого не вызывается моторное выполнение а просто - проход цепочки с просмотром ее звеньев.
	или цепь может быть прервана осознанно
	или пройдена при ее пошаговом отслеживании StepByStepAutomatizm
	и пошаговом запуске: allowNextAutomatizm(automatizm.NextID):
	Бот смотрит, выполнить ли следующий шаг, добавить ли рефлекс мозжечка.
	*/
	NextID int

// СОБСТВЕННОЙ ЭНЕРГИИ НЕТ	Energy int // от 1 до 10, по умолчанию = 5

	/* Уверенность в авторитарном автоматизме высока в период авторитарного обучения
	   	и сильно падает в период собственной инициативы, когда нужно на себе проверить,
	   	а даст ли такой автоматизм в самом деле обещанное улучшение.
Только один из автоматизмов, прикрепленных к ветке, может иметь Belief=2 - проверенное собственное знание
Если задается Belief=2, остальные Belief=2 становится Belief=0.
!!! ПОЭТОМУ ВСЕГДА нужно задавать с помощью SetMentalAutomatizmBelief(atmzm *Automatizm,belief int)
	*/
	Belief int // 0 - предположение, 1 - чужие сведения, 2 - проверенное собственное знание - ШТАТНЫЙ, умолчательный автоматизм
	/* В случае, если в результате автоматизма его Usefulness изменит знак, то
	Count обнулится, а при таком же знаке - увеличивается на 1.
	*/
	Count int // надежность: число использований с подтверждением (бес)полезности Usefulness
}
// привязанные к узлу дерева
var MentalAutomatizmsFromID=make(map[int]*MentalAutomatizm)
///////////////////////////////////////////////////
// ШТАТНЫЕ автоматизмы, прикрепленные к ID узла Дерева с Belief==2 т.е. ШТАТНЫЕ, выполняющиеся не раздумывая
// у узла может быть только один штатный автоматизм с Belief==2
var MentalAutomatizmBelief2FromTreeNodeId = make(map[int]*MentalAutomatizm)
/////////////////////////////////////////////////////////////////////


/* создать новый автоматизм
checkLevel - глубина проверки на идентичность: 0 - нет проверки, 1 - поверхностная, 2 - полная
*/
var lastMentalAutomatizmID=0
var NoWarningMentalCreateShow=false
func createMentalAutomatizmID(id int,BranchID int,Sequence string, checkLevel int)(int,*MentalAutomatizm) {
	/* Автоматизмы уникальны по сочетанию BranchID и Sequence.
	   	При попытке создать с таким же сочетанием возвращается уже созданный.
	    к одной вентке могут быть прикреплены неограниченное число автоматизмов
	*/
	if checkLevel>0 {
		oldID, oldVal := checkUnicumMentalAutomatizm(BranchID, Sequence,checkLevel)
		if oldVal != nil {
			return oldID, oldVal
		}
	}
	if id == 0 {
		lastMentalAutomatizmID++
		id=lastMentalAutomatizmID
	} else {
		if lastMentalAutomatizmID < id {
			lastMentalAutomatizmID = id
		}
	}

	var node MentalAutomatizm
	node.ID = id
	node.BranchID = BranchID
	node.Sequence = Sequence

	MentalAutomatizmsFromID[id] = &node

	if !NoWarningMentalCreateShow {
		lib.WritePultConsol("Создан новый автоматизм.")
	}
	return id, &node
}
//////////////////////////////////////////
/* Автоматизмы уникальны по сочетанию BranchID и Sequence.
Функцию можно использовать для выборки автоматизма с заданными BranchID и Sequence
checkLevel - глубина проверки на идентичность: 0 - нет проверки, 1 - поверхностная, 2 - полная
Полная проверка м.б. пригодиться для ментальных дел, в частности, нахождения автоматизма с заданными BranchID и Sequence
*/
func checkUnicumMentalAutomatizm(BranchID int,Sequence string, checkLevel int)(int,*MentalAutomatizm){
	for id, v := range MentalAutomatizmsFromID {
		if BranchID != v.BranchID || !compareMentalAutomatizmSequence(Sequence,v.Sequence,checkLevel) {
			continue
		}
		return id,v
	}
	return 0,nil
}
/* сравненеие на идентичности двух Sequence
Sequence="Snn:24243,1234,0,24234,11234|Tnn:23|Dnn:24,78"
тестирование - в func PsychicInit()
*/
func compareMentalAutomatizmSequence(Sequence1 string,Sequence2 string, checkLevel int)(bool){
	//a1Arr:=ParceMentalAutomatizmSequence(Sequence1)
	if Sequence1 == Sequence2{
		return true
	}
	if checkLevel==1{// на этом проверка завершается
		return false
	}
	// полная проверка
	sArr:=strings.Split(Sequence1, "|")
	for i := 0; i < len(sArr); i++ {
		if len(sArr[i]) == 0 {
			continue
		}
		pArr := strings.Split(sArr[i], ":")
		switch pArr[0] {
		case "Mnn": // есть ли такой у второго
			if strings.Contains(Sequence2, "Mnn"){
				if !compareMentalBlockContent(pArr[1],Sequence2,"Mnn"){
					return false}
			}else{return false}
		case "Ann":
			if strings.Contains(Sequence2, "Ann"){
				if !compareMentalBlockContent(pArr[1],Sequence2,"Ann"){
					return false}
			}else{return false}
		}
	}
	return true // все блоки действий совпадают
}
// сравнить содержимое блоков данного типа
func compareMentalBlockContent(block1 string,Sequence2 string,kind string)(bool){
	sArr:=strings.Split(Sequence2, "|")
	for i := 0; i < len(sArr); i++ {
		pArr := strings.Split(sArr[i], ":")
		// т.к. последовательности Dnn отсортированы по возрастанию, то можно проверять pArr[1]==block1
		if pArr[0]==kind && pArr[1]==block1{
			return true
		}
	}
	return false
}
////////////////////////////////////////////


// СОХРАНИТЬ структура записи: id|BranchID|Usefulness||Sequence||NextID|Energy|Belief
// В случае отсуствия пусковых стимулов создается ID такого отсутсвия, пример такой записи: 2|||0|0|
func SaveMentalAutomatizm() {
	var out = ""

	for k, v := range MentalAutomatizmsFromID {
		out += strconv.Itoa(k) + "|"
		out += strconv.Itoa(v.BranchID) + "|"
		out += strconv.Itoa(v.Usefulness) + "||"
		out += v.Sequence + "||"
		out += strconv.Itoa(v.NextID) + "|"
		out += strconv.Itoa(v.Belief) + "|"
		out += strconv.Itoa(v.Count)
		out += "\r\n"
	}
	lib.WriteFileContent(lib.GetMainPathExeFile() + "/memory_psy/mental_automatizm_images.txt", out)
}

// ЗАГРУЗИТЬ структура записи: id|BranchID|Usefulness||Sequence||NextID|Energy|Belief
func loadMentalAutomatizm() {
	NoWarningMentalCreateShow = true
	MentalAutomatizmsFromID = make(map[int]*MentalAutomatizm)

	strArr, _ := lib.ReadLines(lib.GetMainPathExeFile() + "/memory_psy/mental_automatizm_images.txt")
	if strArr == nil { return	}
	for n := 0; n < len(strArr); n++ {
		if len(strArr[n]) == 0 { continue	}
		main := strings.Split(strArr[n], "||")
		p := strings.Split(main[0], "|")
		id, _ := strconv.Atoi(p[0])
		BranchID, _ := strconv.Atoi(p[1])
		Usefulness, _ := strconv.Atoi(p[2])
		Sequence := main[1]

		p = strings.Split(main[2], "|")
		NextID, _ := strconv.Atoi(p[0])
		Belief, _ := strconv.Atoi(p[1])
		Count, _ := strconv.Atoi(p[2])
 var saveDoWritingFile= doWritingFile; doWritingFile =false
		_, a := createMentalAutomatizmID(id, BranchID, Sequence,0)// без проверки на уникальность
		a.NextID = NextID
		a.Usefulness = Usefulness
		a.Count = Count
		SetMentalAutomatizmBelief(a, Belief)
doWritingFile =saveDoWritingFile
	}
	NoWarningMentalCreateShow = false
	return
}
/////////////////////////////////////////////////////////////











