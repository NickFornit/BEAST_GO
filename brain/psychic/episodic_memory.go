/* Эпизодическая память последовательности привлечения осозаннного (understandingSituation) внимания
Линейная, не ветсящаяся цепочка.
Запоминается в файле episod_memory.txt для каждого ор.рефлекса (до обработки с редуцированием operMemory)
т.е. это - самая многочисленная часть памяти, поэтому редуцируется во сне (в psychic_sleep_process.go).
*/

package psychic

import (
	"BOT/lib"
	"strconv"
	"strings"
)

/////////////////////////////////////////
func EpisodeMemoryPuls(){
	//запись каждые 100 пульсов
	if PulsCount>=100 && PulsCount%100==0{
		saveEpisodicMenory()
	}

}
///////////////////////////////////////
func EpisodeMemoryInit(){
	loadEpisodicMenory()
}

///////////////////////////////////////
type EpisodeMemory struct {
	NodeAID int // конечный узел активной ветки дерева моторных автоматизмов
	NodePID int // конечный узел активной ветки дерева ментальных автоматизмов
	DiffMood int //изменение ощущаемого настроения после активации веток
	Type int // 0 - объективный образ (внешний ор.рефлекс) или 1 - субъективный (произвольный ор.рефлекс)
}
///////////////////////////////////////
// Эпизодическая память добавляется подряд c возрастанием ID как простой массив,
//во сне редуцируется с перезаписью подряд оставшихся
var EpisodeMemoryObjects []*EpisodeMemory

//последний созданный образ эпизодической памяти
var EpisodeMemoryLastIDFrame int
///////////////////////////////////////////////////////


// добавить НОВЫЙ ЭПИЗОД ПАМЯТИ в understanding_tree_orientation.go
func newEpisodeMemory(){
var DiffMood=0
	if WasOperatorActiveted { // оператор отреагировал
		DiffMood=oldBetterOrWorse // текущее состояние на момент срабатывания автоматизма
	}else{// активация по изменению текущего состояния без действия оператора
		_,DiffMood,_ = wasChangingMoodCondition(1)
	}

	// новый эпизод памяти
	createEpisodeMemoryFrame(detectedActiveLastNodID,
		detectedActiveLastUnderstandingNodID,
		DiffMood,
		0) // TODO пока не вижу как определять Type т.к. нет произвольной активации
}


// создать новый образ сочетаний действий, если такого еще нет
func createEpisodeMemoryFrame(NodeAID int,NodePID int,DiffMood int,Type int)(*EpisodeMemory){

	var node EpisodeMemory
	node.NodeAID=NodeAID
	node.NodePID = NodePID
	node.DiffMood=DiffMood
	node.Type=Type

	EpisodeMemoryObjects=append(EpisodeMemoryObjects,&node)
	EpisodeMemoryLastIDFrame=len(EpisodeMemoryObjects)-1

	//saveEpisodicMenory() запись каждые 100 пульсов и при корректном выходе

	return &node
}
/////////////////////////////////////////



//////////////////////////////////////////
// запись эпизодической памяти
func saveEpisodicMenory(){
	var out=""
	for _, v := range EpisodeMemoryObjects {
		out+=strconv.Itoa(v.NodeAID)+"|"
		out+=strconv.Itoa(v.NodePID)+"|"
		out+=strconv.Itoa(v.DiffMood)+"|"
		out+=strconv.Itoa(v.Type)
		out+="\r\n"
	}
	lib.WriteFileContent(lib.GetMainPathExeFile()+"/memory_psy/episod_memory.txt",out)
}
///////////////////////////////////////////

// запись эпизодической памяти
func loadEpisodicMenory(){
	EpisodeMemoryObjects=nil
	strArr,_:=lib.ReadLines(lib.GetMainPathExeFile()+"/memory_psy/episod_memory.txt")
	cunt:=len(strArr)
	for n := 0; n < cunt; n++ {
		if len(strArr[n]) == 0 {
			continue
		}
		p := strings.Split(strArr[n], "|")
		NodeAID, _ := strconv.Atoi(p[0])
		NodePID, _ := strconv.Atoi(p[1])
		DiffMood, _ := strconv.Atoi(p[2])
		Type, _ := strconv.Atoi(p[3])

		createEpisodeMemoryFrame(NodeAID, NodePID, DiffMood, Type)
	}
	return
}
///////////////////////////////////////////