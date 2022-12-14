/* Cон, его стадии и циклы

В каждом пакете есть флаг  и там для сна выполняется
	if IsSleeping {
		sleepingProcess()
	}
в психике для этого есть psychic_sleep_process.go
*/


package sleep

import (
	"BOT/brain/gomeostas"
	"BOT/brain/psychic"
)

//////////////////////////////////////////
var SlipPulsCount=0 // передача тика Пульса из brine.go
var LifeTime=0
var EvolushnStage=0  // стадия развития
// коррекция текущего состояния с каждым пульсом
func SleepPuls(evolushnStage int,lifeTime int,puls int){
	LifeTime=lifeTime
	EvolushnStage=evolushnStage
	SlipPulsCount=puls // передача номера тика из более низкоуровневого пакета

	// разбудить при включении. Сторожевая функция - в Пульсе рефлексов вызывает sleep.WakeUpping()
	if SlipPulsCount>3{
		WakeUpping() // разбудить при включении
	}

if SlipPulsCount>5{
	prepareWordArr()// обработка накопившегося массива распознанных фраз


}

// понижение повреждений во сне каждый пульс или в час 0.01*3600=36
if !gomeostas.NotAllowSetGomeostazParams{
	gomeostas.GomeostazParams[8]-=0.01
	if gomeostas.GomeostazParams[8] <0{gomeostas.GomeostazParams[8]=0}
}

}
//////////////////////////////////////////////////////


////////////////////////////////////
/* Блокировать выполнение рефлексов на время сна sleep.AllowReflexesAction()
вызывается из reflex_action.go рефлексов
 */
//var isReflexesActionBloking=false
func NotAllowReflexesAction()(bool){
	if psychic.MotorTerminalBlocking{
		return true
	}
	return false
}
////////////////////////////////////