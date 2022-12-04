/* Кратковременная память

Сохраняет подряд ID базовых ментальных цепочек goNext
Не сохраняется, но обрабатывается во сне или в спокойствии с выделением новых Правил.
А так же используется для быстрогого обобщения, ориентации и сознательного поиска.

Быстро растуший массив, ограниченный 1000 элементов
*/

package psychic

////////////////////////////////////////////

type shortTermMemory struct {
	aTreeNodID int // detectedActiveLastUnderstandingNodID
	uTreeNodID int // detectedActiveLastNodID
	GoNextID int // ID goNext
}
var termMemory []shortTermMemory // массив объектов (а не адресов)

/////////////////////////////////////////////

func addShortTermMemory(gnID int){
	if len(termMemory)>1000{
		// удалить первый
		termMemory=termMemory[1:]
	}
	var stm shortTermMemory
	stm.uTreeNodID=detectedActiveLastUnderstandingNodID
	stm.aTreeNodID=detectedActiveLastNodID
	stm.GoNextID=gnID

	termMemory=append(termMemory, stm)
}
///////////////////////////////////////////////



// обработка кратковременной памяти во сне или бездействии
func ShortTermMemoryProcessing(){


}
///////////////////////////////////////////



