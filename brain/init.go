/*
инициализация при запуске
*/

package brain

import (
	"BOT/lib"
	"strconv"
	"strings"
)

// самая первая инициализация (из main.go) после всех приготовлений, до запуска пульса в puls.go
func RunInitialisation(){

	str,_:=lib.ReadLines(lib.GetMainPathExeFile()+"/memory_reflex/life_time.txt")
	if len(str)>0 {
		LifeTime, _ = strconv.Atoi(strings.TrimSpace(str[0]))
	}else{
		panic("ОБНУЛИЛОСЬ ВРЕМЯ ЯЖИЗНИ LifeTime")
	}

	str,_=lib.ReadLines(lib.GetMainPathExeFile()+"/memory_reflex/stages.txt")
	EvolushnStage,_=strconv.Atoi(strings.TrimSpace(str[0]))

}