/* Эпизодическая память последовательности привлечения осозаннного (understandingSituation) внимания
Срабатывает по каждому действию оператора orientationConsciousness
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
	// последствия действий оператора
	NodeAID int // конечный узел активной ветки дерева моторных автоматизмов
	NodePID int // конечный узел активной ветки дерева ментальных автоматизмов
	RulesID int // действия оператора, ответное действие и эффект

	LifeTime int
	Type int // 0 - объективный образ (внешний ор.рефлекс) или 1 - субъективный (произвольный ор.рефлекс)
}
///////////////////////////////////////
// Эпизодическая память добавляется подряд c возрастанием ID как простой массив,
//во сне редуцируется с перезаписью подряд оставшихся
var EpisodeMemoryObjects []*EpisodeMemory

//последний созданный образ эпизодической памяти
var EpisodeMemoryLastIDFrame=0

// последний эпизод, который был осмыслен в лени или во сне
var EpisodeMemoryLastCalcID=0
///////////////////////////////////////////////////////


/* добавить НОВЫЙ ЭПИЗОД ПАМЯТИ в understanding_tree_orientation.go
вызывается только в func calcAutomatizmResult:
 */
func newEpisodeMemory(rulesID int){
	// новый эпизод памяти
	createEpisodeMemoryFrame(detectedActiveLastNodID,
		rulesID,
		LifeTime,
		0) // TODO пока не вижу как определять Type т.к. нет произвольной активации
}


// создать новый образ сочетаний действий, если такого еще нет
func createEpisodeMemoryFrame(NodePID int,rulesID int,lifeTime int,Type int)(*EpisodeMemory){

	var node EpisodeMemory
	node.NodeAID=NodePID
	node.NodePID = NodePID
	node.RulesID=rulesID
	node.LifeTime=lifeTime
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
		out+=strconv.Itoa(v.RulesID)+"|"
		out+=strconv.Itoa(v.LifeTime)+"|"
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
		NodePID, _ := strconv.Atoi(p[0])
		RulesID, _ := strconv.Atoi(p[1])
		lifeTime, _ := strconv.Atoi(p[2])
		Type, _ := strconv.Atoi(p[3])

		createEpisodeMemoryFrame(NodePID, RulesID, lifeTime, Type)
	}
	return
}
///////////////////////////////////////////