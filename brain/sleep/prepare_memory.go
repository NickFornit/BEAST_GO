/*  информационная обработка различных видов памяти

Образ цели бессловестного действия Формируется временно и не сохранятся в файле
Объекты PurposeGeneticObject накапливаются в оперативке и удаляются во сне

*/

package sleep

import (
	word_sensor "BOT/brain/words_sensor"
	"BOT/brain/psychic"
)



// обработка накопившегося массива распознанных фраз
func prepareWordArr(){
		wCount:=len(word_sensor.MemoryDetectedArr)
		if wCount==0{

		}


	// обработка кратковременной памяти во сне или бездействии
	psychic.ShortTermMemoryProcessing()
}




