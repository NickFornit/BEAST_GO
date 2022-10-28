/* Ментальные (умственные) автоматизмы мышления.
Базовый мент.автоматизм может прикрепляться ТОЛЬКО к последнему узлу ветки - при полном понимании ситуации

id|BranchID|Usefulness||ActionsImageID||NextID|Count
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
/* К узлу UnderstandingNode привязывается только базовый автоматизм.
А к нему - цепочка по NextID у членов которой BranchID==0
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

// СОБСТВЕННОЙ ЭНЕРГИИ НЕТ	Energy int // от 1 до 10, по умолчанию = 5
// Belief int нет т.к. только один автоматизм может быть привязан к данной BranchID

	/* В случае, если в результате автоматизма его Usefulness изменит знак, то
	Count обнулится, а при таком же знаке - увеличивается на 1.
	*/
	Count int // надежность: число использований с подтверждением (бес)полезности Usefulness
}
var MentalAutomatizmsFromID=make(map[int]*MentalAutomatizm)
// Базовые автоматизмы, привязанные к узлу UnderstandingNode ветки дерева.
// Только один автоматизм может быть привязан к данной BranchID, поэтому тут нет Belief
var MentalAutomatizmsFromBranchID=make(map[int]*MentalAutomatizm)
///////////////////////////////////////////////////


/* создать новый автоматизм
checkLevel - глубина проверки на идентичность: 0 - нет проверки, 1 - поверхностная, 2 - полная
*/
var lastMentalAutomatizmID=0
var NoWarningMentalCreateShow=false
func createMentalAutomatizmID(id int,BranchID int,ActionsImageID int, checkLevel int)(int,*MentalAutomatizm) {
	/* Автоматизмы уникальны по сочетанию ActionsImageID.
	   	При попытке создать с таким же сочетанием возвращается уже созданный.
	    к одной вентке могут быть прикреплены неограниченное число автоматизмов
	*/
	if checkLevel>0 {
		oldID, oldVal := checkUnicumMentalAutomatizm(BranchID,ActionsImageID,checkLevel)
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
	node.ActionsImageID = ActionsImageID

	MentalAutomatizmsFromID[id] = &node
	MentalAutomatizmsFromBranchID[BranchID] = &node

	if !NoWarningMentalCreateShow {
		lib.WritePultConsol("Создан новый автоматизм.")
	}
	return id, &node
}
//////////////////////////////////////////
/* Автоматизмы уникальны по ActionsImageID.
Функцию можно использовать для выборки автоматизма с заданным ActionsImageID
checkLevel - глубина проверки на идентичность: 0 - нет проверки, 1 - поверхностная, 2 - полная
Полная проверка м.б. пригодиться для ментальных дел, в частности, нахождения автоматизма с заданным ActionsImageID
*/
func checkUnicumMentalAutomatizm(BranchID int,ActionsImageID int, checkLevel int)(int,*MentalAutomatizm){
	for id, v := range MentalAutomatizmsFromID {
		if BranchID != v.BranchID && ActionsImageID != v.ActionsImageID {
			continue
		}
		return id,v
	}
	return 0,nil
}
////////////////////////////////////////////


// СОХРАНИТЬ структура записи: id|BranchID|Usefulness||ActionsImageID||NextID|Count
// В случае отсуствия пусковых стимулов создается ID такого отсутсвия, пример такой записи: 2|||0|0|
func SaveMentalAutomatizm() {
	var out = ""

	for k, v := range MentalAutomatizmsFromID {
		out += strconv.Itoa(k) + "|"
		out += strconv.Itoa(v.BranchID) + "|"
		out += strconv.Itoa(v.Usefulness) + "|"
		out += strconv.Itoa(v.ActionsImageID) + "|"
		out += strconv.Itoa(v.NextID) + "|"
		out += strconv.Itoa(v.Count)
		out += "\r\n"
	}
	lib.WriteFileContent(lib.GetMainPathExeFile() + "/memory_psy/mental_automatizm_images.txt", out)
}

// ЗАГРУЗИТЬ структура записи: id|BranchID|Usefulness||ActionsImageID||NextID|Count
func loadMentalAutomatizm() {
	NoWarningMentalCreateShow = true
	MentalAutomatizmsFromID = make(map[int]*MentalAutomatizm)

	strArr, _ := lib.ReadLines(lib.GetMainPathExeFile() + "/memory_psy/mental_automatizm_images.txt")
	if strArr == nil { return	}
	for n := 0; n < len(strArr); n++ {
		if len(strArr[n]) == 0 { continue	}
		p := strings.Split(strArr[n], "|")
		id, _ := strconv.Atoi(p[0])
		BranchID, _ := strconv.Atoi(p[1])
		Usefulness, _ := strconv.Atoi(p[2])
		ActionsImageID, _ := strconv.Atoi(p[3])
		NextID, _ := strconv.Atoi(p[4])
		Count, _ := strconv.Atoi(p[5])
 var saveDoWritingFile= doWritingFile; doWritingFile =false
		_, a := createMentalAutomatizmID(id, BranchID, ActionsImageID,0)// без проверки на уникальность
		a.NextID = NextID
		a.Usefulness = Usefulness
		a.Count = Count
doWritingFile =saveDoWritingFile
	}
	NoWarningMentalCreateShow = false
	return
}
/////////////////////////////////////////////////////////////



/*


func createAndRunMentalAutomatizm(maImgID int)(int,*MentalAutomatizm){

	k, v := createMentalAutomatizmID(0,maImgID,1)
	return k, v
}*/
//////////////////////////////////////////////////////











