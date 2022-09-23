/*  Функция сохранения всех текущих данных, записываемых в файлы

*/

package brain

import (
	"BOT/brain/psychic"
	"BOT/brain/reflexes"
	word_sensor "BOT/brain/words_sensor"
)


func init(){
	//res:=SaveAll()
	//if res{	}
}

////////////////////////////////////////
/* все данные сохраняются при нажатии на Пульте Выключить Beast,
при ручном сохранении файлов (шестеренка)
и в main.go при выходе
 !!! только после того, как все данные будут загружены:
 */
func SaveAll()(bool){
	if PulsCount <5{
		return true
	}
	var success=true
	defer func(){// ловим панику
		if err := recover(); err != nil {
			success=false
		}
	}()
	// сохранения всего
		//pppp()
	// !!!! saveLifeTime() записывает только в puls.go
	word_sensor.SaveWordTree()
	word_sensor.SavePhraseTree()

	// сохранение всех файлов по рефлексам
	reflexes.SaveReflexesAttributes()

	psychic.SaveAllPsihicMemory()


	return success
}
/////////////////////////////////////////

func pppp(){// паника для тектирования
	var n=1
	p:=12/(n-1)
	if p>0{	}
}