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
	ID int

	LifeTime int // момент жизни

	Mood int // сила ощущаемого настроения PsyMood (-10..0..10)
	PsyBaseMood int // ощущаемое настроение: -1 Плохое настроение, 0 Нормальное, 1 - хорошее настроение
// все ветки дерева понимания, включая текущую ситуацию:
	NodeID[] int // узлы активных веток Дерева Понимания в последовательности их активации

	Type int // 0 - объективный образ (внешний ор.рефлекс) или 1 - субъективный (произвольный ор.рефлекс)
}
///////////////////////////////////////
// Эпизодическая память добавляется подряд c возрастанием ID как простой массив,
//во сне редуцируется с перезаписью подряд оставшихся
var EpisodeMemoryObjects []*EpisodeMemory

//последний созданный образ эпизодической памяти
var EpisodeMemoryLastIDFrame int
///////////////////////////////////////////////////////


/////////// НОВЫЙ ЭПИЗОД ПАМЯТИ в реализации
func newEpisodeMemory(){
	// выдать массив ID узлов ветки по заданному ID узла
	currentUnderstandingNodeID:=getBrangeUnderstandingNodeIdArr(detectedActiveLastUnderstandingNodID)
	// новый эпизод памяти
	createEpisodeMemoryFrame(LifeTime,PsyMood,PsyBaseMood,currentUnderstandingNodeID,0)
}


// создать новый образ сочетаний действий, если такого еще нет
func createEpisodeMemoryFrame(LifeTime int,Mood int,PsyBaseMood int,NodeID[] int,Type int)(*EpisodeMemory){

	var node EpisodeMemory
	//node.ID = id
	node.LifeTime=LifeTime
	node.Mood = Mood
	node.PsyBaseMood=PsyBaseMood
	node.NodeID=NodeID
	node.Type=Type

	EpisodeMemoryObjects=append(EpisodeMemoryObjects,&node)
	EpisodeMemoryLastIDFrame=node.ID

	//saveEpisodicMenory() запись каждые 100 пульсов и при корректном выходе

	return &node
}
/////////////////////////////////////////



//////////////////////////////////////////
// запись эпизодической памяти
func saveEpisodicMenory(){
	var out=""
	for k, v := range EpisodeMemoryObjects {
		out+=strconv.Itoa(k)+"|"
		out+=strconv.Itoa(v.LifeTime)+"|"
		out+=strconv.Itoa(v.Mood)+"|"
		out+=strconv.Itoa(v.PsyBaseMood)+"|"

		for i := 0; i < len(v.NodeID); i++ {
			out+=strconv.Itoa(v.NodeID[i])+","
		}
		out+="|"
		out+=strconv.Itoa(v.Type) //+"|"
		//out+=strconv.Itoa(v.PrevID)+"|"
		//out+=strconv.Itoa(v.NextID)
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
	var lastEM int
	for n := 0; n < cunt; n++ {
		if len(strArr[n])==0{
			continue
		}
		p:=strings.Split(strArr[n], "|")
		//id,_:=strconv.Atoi(p[0])
		LifeTime,_:=strconv.Atoi(p[1])
		Mood,_:=strconv.Atoi(p[2])
		PsyBaseMood,_:=strconv.Atoi(p[3])
		s:=strings.Split(p[4], ",")
		var NodeID[] int
		for i := 0; i < len(s); i++ {
			if len(s[i])==0{
				continue
			}
			ni,_:=strconv.Atoi(s[i])
			NodeID=append(NodeID,ni)
		}
		Type,_:=strconv.Atoi(p[5])
	//	PrevID,_:=strconv.Atoi(p[6])
	//	NextID,_:=strconv.Atoi(p[7])

		em:=createEpisodeMemoryFrame(LifeTime,Mood,PsyBaseMood,NodeID,Type)
	//	em.PrevID=PrevID
	//	em.NextID=NextID
		lastEM=em.ID
	}
	EpisodeMemoryLastIDFrame=lastEM
	return
}
///////////////////////////////////////////