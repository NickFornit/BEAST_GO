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
	PsyBaseMood int // ощущаемое настроение (-1,0,1) - база дерева понимания

	NodeID[] int // узлы активных веток Дерева Понимания в последовательности их активации

	Type int // 0 - объективный образ или 1 - субъективный
	PrevID int // ID предыдущего эпизода для просмотра прошлого
	NextID int // ID последующего, когда он появляется
}
///////////////////////////////////////
// Эпизодическая память добавляется подряд как простой массив,
//но во сне редуцируется, так что придется использовать карты
var EpisodeMemoryObjects=make(map[int]*EpisodeMemory)
//текущий активный образ эпизодической памяти
var EpisodeMemoryActiveFrame *EpisodeMemory
//последний созданный образ эпизодической памяти
var EpisodeMemoryLastFrame *EpisodeMemory
///////////////////////////////////////////////////////



// создать новый образ сочетаний действий, если такого еще нет
var lastEpisodeMemoryFrameID=0
func createEpisodeMemoryFrame(id int,LifeTime int,Mood int,PsyBaseMood int,NodeID[] int,Type int)(int,*EpisodeMemory){
	if id==0{
		lastEpisodeMemoryFrameID++
		id=lastEpisodeMemoryFrameID
	}else{
		//		newW.ID=id
		if lastEpisodeMemoryFrameID<id{
			lastEpisodeMemoryFrameID=id
		}
	}

	var node EpisodeMemory
	node.ID = id
	node.LifeTime=LifeTime
	node.Mood = Mood
	node.PsyBaseMood=PsyBaseMood
	node.NodeID=NodeID
	node.Type=Type
	if EpisodeMemoryLastFrame!=nil { // еще не было эп.памяти
		node.PrevID = EpisodeMemoryLastFrame.ID
		EpisodeMemoryObjects[EpisodeMemoryLastFrame.ID].NextID = id
	}
	EpisodeMemoryLastFrame=&node

	EpisodeMemoryObjects[id]=&node

	//saveEpisodicMenory() запись каждые 100 пульсов и при корректном выходе

	return id,&node
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
		out+=strconv.Itoa(v.Type)+"|"
		out+=strconv.Itoa(v.PrevID)+"|"
		out+=strconv.Itoa(v.NextID)
		out+="\r\n"
	}
	lib.WriteFileContent(lib.GetMainPathExeFile()+"/memory_psy/episod_memory.txt",out)
}
///////////////////////////////////////////

// запись эпизодической памяти
func loadEpisodicMenory(){
	EpisodeMemoryObjects=make(map[int]*EpisodeMemory)
	strArr,_:=lib.ReadLines(lib.GetMainPathExeFile()+"/memory_psy/episod_memory.txt")
	cunt:=len(strArr)
	var lastEM *EpisodeMemory
	for n := 0; n < cunt; n++ {
		if len(strArr[n])==0{
			continue
		}
		p:=strings.Split(strArr[n], "|")
		id,_:=strconv.Atoi(p[0])
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
		PrevID,_:=strconv.Atoi(p[6])
		NextID,_:=strconv.Atoi(p[7])

		_,em:=createEpisodeMemoryFrame(id,LifeTime,Mood,PsyBaseMood,NodeID,Type)
		em.PrevID=PrevID
		em.NextID=NextID
		lastEM=em
	}
	EpisodeMemoryLastFrame=lastEM
	return
}
///////////////////////////////////////////