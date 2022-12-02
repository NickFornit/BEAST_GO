package brain

import (
	"BOT/lib"
	_ "fmt"
	"strconv"
	"time"
)

/* Пульс - сердце Beast :)
запускает  каждую секунду Аctivatin() в net.go, поддерживая возможность самостоятельных действий Beast
В общем, это обеспечивает возможность параллельной работы функций Beast, отводя им выполнение с каждым пульсом
*/

var PulsCount = 0 		// счетчик пульса
var noWorkung = false // не активировать очередной цикл если true - тишина
var startWait = 0  		// начало ожидания
var secWaiting = 1 		// время ожидания
var LifeTime int 			// время жизни в числе пульсов
var EvolushnStage int // стадия развития

func init() {
// Puls() Запускается только в одном месте (main.go) чтобы не было двух потоков пульса!!!
}

// сохранить время жизни в файл
func saveLifeTime() {
	if LifeTime>10 {// иногда life_time.txt обнулялся...
		lib.WriteFileContent(lib.GetMainPathExeFile()+"/memory_reflex/life_time.txt", strconv.Itoa(LifeTime))
	}
}

/* запуск пульса по сигналам из getParams := r.FormValue("get_params")

 */
var blockinfP=0
func SincroTic(){
	blockinfP=1
	Puls()
}
////////////////////////////////


//водитель ритма пульса
var IsPultActivnost = false	// начало активности с Пульта  brain.IsPultActivnost=true brain.IsPultActivnost=false
func Puls() {
	if blockinfP>0{// не запускать генеретор пульса в этот раз - синхронизация из Пульта
		pulsActions()// действия дял этого такта пульса
		blockinfP=0
		return
	}
	time.AfterFunc(1 * time.Second, func() {

		// действия дял этого такта пульса
		pulsActions()

		Puls()
	})
}
// действия, совершаемые по каждому пульсу
func pulsActions(){
	// сканирование состояния (пульс):
	//now := time.Now()
	//curTime := now.Second()//.Millisecond()
	//log.Println("пульс",curTime)
	if PulsCount == 2 {
		lib.WritePultConsol("Beast активируется.")
	}
	/* ЗАЧЕМ ОСТАНАВЛИВАТЬ АКТИВНОСТЬ ПРИ СОБЫТИИ С ПУЛЬТА??
		 NotAllowAnyActions=true должно использоваться ТОЛЬКО для остановки активности дял ЗАПИСИ ПАМЯТИ!
	if  IsPultActivnost{// остановка любой активности Beast
		NotAllowAnyActions=true
	}else{
		NotAllowAnyActions=false
	}
	*/
	if !noWorkung {
		Аctivation(LifeTime) // Главный цикл (пульс) опроса состояния Beast и его ответов в brain.go
	}
	if startWait > 0 && !IsPultActivnost { // ждем следующего цикла и выполняем inaction
		if PulsCount > (startWait + 1) {
			inaction()
			noWorkung = false // конец ожидания
		}
	}
	LifeTime++
	//if (PulsCount + 1) % 10 == 0 {
	//	saveLifeTime()
	//}
	//ac := time.Date(0, time.January, 1, 0, 0, 0, 0, time.UTC)
	//lastPulsTime=curTime
	PulsCount++
}
///////////////////////////////////////////////////

var NotAllowAnyActions = false // - запрет любой активности  brain.NotAllowAnyActions
// то, что должно выполнять в тишине, в бездействии
func inaction() {
	/* Память сохраняется при корректном выключении, попытки автоматизировать - очень геморные и ненадежные
	// сохранять все каждые 100 сек если нет действий
	if (PulsCount+1)%100==0{
		NotAllowAnyActions=true
		SaveAll()
		NotAllowAnyActions=false
	}
	*/
}

/* если нужно остановить все на время waitSec для выполнения какой-то функции, то
Пример для функции созранения всей памяти:
func WaitSaveAllPryMemory(){
	OwnFuncForRun=SaveAllPryMemory
	StopAll() // выполнить  saveAllPryMemory() в тишине, когда ничто не работает
}
*/
func StopAll() {
	noWorkung = true	// останавливает puls()
	startWait = PulsCount
}

func StopRunAll(stop bool) {
	if stop {// останавливает puls()
		noWorkung = true
	} else {
		noWorkung = false
	}
}