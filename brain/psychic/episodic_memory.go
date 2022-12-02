/* Эпизодическая память последовательности Правил (Стимул-Ответ-Эффект) прежнего опыта реагирования.
Фактически представляет собой дерево решений,
в котором каждое последующее решение тематически зависит от имеющихся цепочек последовательностей.

Поиск решений реализуется функцией getRulesFromEpisodicsSlice для только что пройденной последовательности Правил.
Подробнее - в описании данной функции.

Линейная, не ветвящаяся  цепочка.
Запоминается в файле episod_memory.txt для каждого ор.рефлекса (до обработки с редуцированием operMemory)
т.е. это - самая многочисленная часть памяти, поэтому оптимизируется, редуцируясь во сне (в psychic_sleep_process.go).
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
	NodeAID            int // конечный узел активной ветки дерева моторных автоматизмов
	NodePID            int // конечный узел активной ветки дерева ментальных автоматизмов
	TriggerAndActionID int // действия оператора, ответное действие и эффект

	LifeTime int
	Type int // 0 - объективный образ (внешний ор.рефлекс) или 1 - субъективный (произвольный ор.рефлекс)
}
///////////////////////////////////////
// Эпизодическая память добавляется подряд c возрастанием ID как простой массив,
//во сне редуцируется с перезаписью подряд оставшихся
var EpisodeMemoryObjects []*EpisodeMemory

//последний ID эпизода эпизодической памяти
var EpisodeMemoryLastIDFrameID=0

/* последний эпизод, который был осмыслен в лени или во сне
 */
var EpisodeMemoryLastCalcID=0
///////////////////////////////////////////////////////

/* период времени прерывания цепочки связанных правил
Для функции getLastRulesSequenceFromEpisodeMemory (вытащить из эпизод.памяти посленюю цепочку кадров)

Резон выбора значения EpisodeMemoryLastCalcID:
Время реакции оператора в период ожидания - максимум 20 пульсов.
Получается, что для группы в 5 Правил нужно значение EpisodeMemoryPause не менее 20*5=100 пульсов.
 */
var EpisodeMemoryPause=100 // в числе пульсов

/* кэш последних 7 эпизодов EpisodeMemory
Заполняется как регистр сдвига: последий эпизод делается первым, а все сдвигаются далее.
 */
//const EpisodeMemoryChechCount = 7
//var EpisodeMemoryChech=make([]*EpisodeMemory, EpisodeMemoryChechCount)

/* добавить НОВЫЙ ЭПИЗОД ПАМЯТИ в understanding_tree_orientation.go
вызывается только в func calcAutomatizmResult:
 */
func newEpisodeMemory(rulesID int,kind int)(*EpisodeMemory){
	// новый эпизод памяти - для условий ПРЕДЫДУЩЕЙ активации деревьев, - как и Правила
	em:=createEpisodeMemoryFrame(detectedActiveLastNodPrevID,
		detectedActiveLastUnderstandingNodPrevID,
		rulesID,
		LifeTime,
		kind)
	return em
}


// создать новый кадр памяти
func createEpisodeMemoryFrame(NodeAID int,NodePID int,rulesID int,lifeTime int,Type int)(*EpisodeMemory){

	var node EpisodeMemory
	node.NodeAID=NodeAID
	node.NodePID = NodePID
	node.TriggerAndActionID =rulesID
	node.LifeTime=lifeTime
	node.Type=Type

	EpisodeMemoryObjects=append(EpisodeMemoryObjects,&node)
	EpisodeMemoryLastIDFrameID=len(EpisodeMemoryObjects)-1

	if doWritingFile {saveEpisodicMenory() }

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
		out+=strconv.Itoa(v.TriggerAndActionID)+"|"
		out+=strconv.Itoa(v.LifeTime)+"|"
		out+=strconv.Itoa(v.Type)
		out+="\r\n"
	}
	lib.WriteFileContent(lib.GetMainPathExeFile()+"/memory_psy/episod_memory.txt",out)
}
///////////////////////////////////////////

// загрузка эпизодической памяти
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
		RulesID, _ := strconv.Atoi(p[2])
		lifeTime, _ := strconv.Atoi(p[3])
		Type, _ := strconv.Atoi(p[4])
var saveDoWritingFile= doWritingFile; doWritingFile =false
		createEpisodeMemoryFrame(NodeAID, NodePID, RulesID, lifeTime, Type)
doWritingFile =saveDoWritingFile
	}
	return
}
///////////////////////////////////////////


/* Добавление нового эпизода в кэш EpisodeMemoryChech
Заполняется как регистр сдвига: последий эпизод делается первым, а все сдвигаются далее.

func AddEpisodeMemoryChech(episode *EpisodeMemory){
	if episode == nil {return}
	// сдвинуть все имеющиеся
	for i := EpisodeMemoryChechCount-1; i >=0; i-- {
		EpisodeMemoryChech[i]=EpisodeMemoryChech[i-1]
	}
	EpisodeMemoryChech[0]=episode
}
// выдает, сколько эпизодов с последнего (первый эпизод) и до времени EpisodeMemoryPause
func GetEpisodeMemoryChechCount(kind int)(int){
	var count=0
	for i := 0; i < EpisodeMemoryChechCount; i++ {
		em:=EpisodeMemoryChech[i]
		if em == nil || em.Type != kind || (LifeTime - em.LifeTime) >EpisodeMemoryPause{
			return count
		}
		count++
	}
return count
}*/
////////////////////////////////////////////




/* Вытащить из эпизод.памяти посленюю цепочку кадров типа (kind=) 0 - объективная память
 */
func getLastRulesSequenceFromEpisodeMemory(kind int, limit int)([]int){
	if EpisodeMemoryLastIDFrameID==0{
		return nil
	}
	var beginID=0
	var preLifeTime=0
	for i := EpisodeMemoryLastIDFrameID; i >=0; i-- {
		em:=EpisodeMemoryObjects[i]
		if preLifeTime==0{
			preLifeTime=em.LifeTime
		}
		if em == nil ||
// объективные чредуются с ментальными во время диалога!!!!  em.Type != kind || // если самый последний эпизод уже не является em.Type == kind
(em.LifeTime - preLifeTime) >EpisodeMemoryPause ||// - главный критерий границы контекста цепочки э.памчти
			beginID >=limit{
			break // закончить выборку
		}
		if em.Type == kind {// игнорировать другой вид памяти
			beginID++
		}
	}
	if beginID == 0 {
		return nil
	}
	var rImg []int
	// перебор последнего фрагмента кадров эпиз.памяти
	//beginID+1 чтобы число проходов цикла было равно beginID и окончился на i <= EpisodeMemoryLastIDFrameID
	for i := EpisodeMemoryLastIDFrameID - beginID+1; i <= EpisodeMemoryLastIDFrameID; i++ {
		em := EpisodeMemoryObjects[i]
		if em.Type == kind {
			rImg = append(rImg, em.TriggerAndActionID)
		}
	}
	// групповое Правило
	if len(rImg)>1{
		lastFrame:=EpisodeMemoryObjects[len(rImg)-1]
		if lastFrame!=nil {
			createNewRules(0, lastFrame.NodeAID, lastFrame.NodePID, rImg, true) // записать (если еще нет такого) групповое правило
		}
		return rImg
	}

	return nil
}
/////////////////////////////////////////////////







