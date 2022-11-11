/* Короткие цепочки вспомогательной ментальной памяти.

*/


package psychic

import (
	"BOT/lib"
	"strconv"
	"strings"
)

////////////////////////////////////////////

/* Для формирования ментального Правила -
временное сохранеиние цикла осмысления между объективным вызовом consciousness и
сразу после запуска ОТВЕТНОГО (на стимул) моторного автоматизма.
т.е. в конце saveFromNextIDAnswerCicle должен быть мент.автоматизм моторного запуска,
если только Стимул не возникнет, не дожидаясь ответа на предыдущий.
В таком случае ментальное Правило не формируется.

Может быть с разными базовыми звеньями из-за произвольной переактивации.
*/
var saveFromNextIDAnswerCicle []int // последовательность fromNextID
//////////////////////////////////////////////////////////////////////////


/* временное сохранеиние цикла осмысления между двумя объективными вызовами consciousness
Всегда повторяет последний фагмент Кратковременной памяти ShortTermMemory
Может быть с разными базовыми звеньями из-за произвольной переактивации.
*/
var saveFromNextIDcurretCicle []int // последовательность fromNextID
//////////////////////////////////////////////////////////////////////////

func initMentalMemories(){
	//savePorposeIDcurrentCicle=nil
	saveFromNextIDbaseLinksCicle=nil
	loadInterruptMemory()
}
//////////////////////////////////


/* стек для обобщений: 7 Базовых fromNextID
Сохраняет 7 Базовых fromNextID (звенья начала цепочки)
*/
var saveFromNextIDbaseLinksCicle []int // последовательность goNextFromUnderstandingNodeIDArr
func addMewBaseLinksMemory(BaseLinksID int){
	if len(InterruptMemory)>7{
		// удалить первый
		saveFromNextIDbaseLinksCicle=saveFromNextIDbaseLinksCicle[1:]
	}
	saveFromNextIDbaseLinksCicle=append(saveFromNextIDbaseLinksCicle, BaseLinksID)
}
//////////////////////////////////////////////////////////////////////////


/* стек до 7 прерванных задач
Запоминаются в файле и актуализируются при просыпании
 */
var InterruptMemory []int
func addMewInterruptMemory(stopNextID int){
	if len(InterruptMemory)>7{
		// удалить первый
		InterruptMemory=InterruptMemory[1:]
	}
	InterruptMemory=append(InterruptMemory, stopNextID)
}
// добавить в стек прерывание
func addInterruptMemory() {
	// Добавить в стек до 7 прерванных решений
	addMewInterruptMemory(currrentFromNextID)
	// что-то еще?
}
func loadInterruptMemory(){
	InterruptMemory=nil
	strArr,_:=lib.ReadLines(lib.GetMainPathExeFile()+"/memory_psy/interrupt_memory.txt")
	cunt:=len(strArr)
	for n := 0; n < cunt; n++ {
		if len(strArr[n]) == 0 {
			continue
		}
		p := strings.Split(strArr[n], "|")
		for i := 0; i < len(p); i++ {
			im,_:=strconv.Atoi(p[i])
			InterruptMemory = append(InterruptMemory, im)
		}
	}
}
func saveInterruptMemory(){
	var out=""
	for n := 0; n < len(InterruptMemory); n++ {
		if n>0{out+="|"}
		out+=strconv.Itoa(InterruptMemory[n])
	}
	out+="\r\n"
	lib.WriteFileContent(lib.GetMainPathExeFile()+"/memory_psy/interrupt_memory.txt",out)
}
//////////////////////////////////////////////////////////////////////////



/* стек текущих целей  НЕ ПРИГОДИЛОСЬ, непонятно, что это и зачем
реализация оценивается как эффект правил
При оценке успешности в afterWaitingPeriod() смотрится, какеи цели достигнуты:


var savePorposeIDcurrentCicle []int //
func addMewPorposeMemory(porposeID int){
	if len(InterruptMemory)>7{
		// удалить первый
		savePorposeIDcurrentCicle=savePorposeIDcurrentCicle[1:]
	}
	savePorposeIDcurrentCicle=append(savePorposeIDcurrentCicle, porposeID)
}
*/
//////////////////////////////////////////////////////////////////////////



