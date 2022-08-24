/*  общие глобальные дела

нейросеть
 */
package brain

import (
	"BOT/brain/action_sensor"
	"BOT/brain/gomeostas"
	"BOT/brain/psychic"
	"BOT/brain/reflexes"
	"BOT/brain/sleep"
	TerminalActions "BOT/brain/terminete_action"
)

///////////////////////////////////////////////////////////////////////////
var IsStressing=0;// не успевает за 1 сек отсканировать
var isWorking=0



/* Главный цикл (пульс) опроса состояния Beast и его ответов
активируется каждую секунду из neurons.go
 */

func Аctivation(lifeTime int){
	// true - смерть Beast при повреждении >99%
	if gomeostas.IsBeastDeath{// прекращение раздачи пульса
		return
	}
	if isWorking==1 {
		IsStressing=1
		return
	}else {
		IsStressing=0
	}
	// начало обработки нейросети
	isWorking=1

	isSlipping:=sleep.GetSleepCondition()

	// текущее состояние гомеостаза и базового контекста с каждым пульсом
	gomeostas.GomeostazPuls(EvolushnStage,lifeTime,PulsCount,isSlipping)
	action_sensor.ActionSensorPuls(EvolushnStage,lifeTime,PulsCount,isSlipping)
	sleep.SleepPuls(EvolushnStage,lifeTime,PulsCount)
	reflexes.ReflexCountPuls(EvolushnStage,lifeTime,PulsCount,isSlipping)
	TerminalActions.TermineteActionCountPuls(EvolushnStage,lifeTime,PulsCount,isSlipping)
	psychic.PsychicCountPuls(EvolushnStage,lifeTime,PulsCount,isSlipping)
	psychic.SetNotAllowAnyActions(NotAllowAnyActions)

	// конец обработки нейросети
	isWorking=0
}

