/*  Выполнение действий, список возможных акций Beast */

package TerminalActions

import (
	"BOT/brain/gomeostas"
	"BOT/lib"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

func init() {
	loadTerminalActons()
}

var EnergyDescrib = []string{
	"Едва (сила=1)",
	"Очень слабо (сила=2)",
	"Слабо (сила=3)",
	"Ощутимо (сила=4)",
	"Средне (сила=5)",
	"Повышенно (сила=6)",
	"Настойчиво (сила=7)",
	"Сильно (сила=8)",
	"Очень сильно (сила=9)",
	"Максимально (сила=10)",
}

/* Генетические цели действий Beast ID гомео-параметров, которые призвано улучшить данное действие - по его ID
Значение 10 означает - для всех параметров - улучшение.

Это – наследственно заданная цель действия, не осознаваемая при его совершении.
Но с опытом каждому действию в конкретных условиях (и к ним добавляются слова и фразы)
будет ассоциироваться смысл (осознаваемая значимость).
*/
var TerminalActionsTargetsFromID = make(map[int][]int)

// ПУЛЬС
var ReflexPulsCount = 0 // передача тика Пульса из brine.go
var LifeTime = 0
var EvolushnStage = 0 // стадия развития
var IsSlipping = false

// коррекция текущего состояния гомеостаза и базового контекста с каждым пульсом
func TermineteActionCountPuls(evolushnStage int, lifeTime int, puls int, isSlipping bool) {
	LifeTime = lifeTime
	EvolushnStage = evolushnStage
	ReflexPulsCount = puls // передача номера тика из более низкоуровневого пакета
	IsSlipping = isSlipping

	pulsSimpleReflexex()
}

/* затратные последствия ответного действия
т.е. то, как изменятся параметры гомеостаза при совершении данного действия
*/
type TerminalActionsExpenses struct {
	GomeoID int     // ID затрагиваемого параметра гоместаза
	Diff    float64 // на сколько изменится параметр
}

// имя ответного действия по его ID
var TerminalActonsNameFromID = make(map[int]string)

// затратные последствия ответного действия по его ID
var TerminalActionsExpensesFromID = make(map[int][]TerminalActionsExpenses)
var LastTerminalActons = 0 // максимальный номер действия

// Загрузить акции
func loadTerminalActons() {
	TerminalActonsNameFromID = make(map[int]string)
	lines, _ := lib.ReadLines(lib.GetMainPathExeFile() + "/memory_reflex/terminal_actons.txt")
	for i := 0; i < len(lines); i++ {
		p := strings.Split(lines[i], "|")
		id, _ := strconv.Atoi(p[0])
		if LastTerminalActons < id {
			LastTerminalActons = id
		}
		TerminalActonsNameFromID[id] = p[1]
		//это действия приводит к изменению гомео-параметров:
		var expenses []TerminalActionsExpenses
		se := strings.Split(p[2], ";")
		for j := 0; j < len(se); j++ {
			if len(se[j]) == 0 {
				continue
			}
			e := strings.Split(se[j], ">")
			e1, _ := strconv.Atoi(e[0])
			e2, _ := strconv.ParseFloat(e[1], 64)
			expenses = append(expenses, TerminalActionsExpenses{e1, e2})
		}
		TerminalActionsExpensesFromID[id] = expenses

		tg := strings.Split(p[3], ",")

		var tgArr []int
		for j := 0; j < len(tg); j++ {
			if len(tg[j]) == 0 {
				continue
			}
			gID, _ := strconv.Atoi(tg[j])
			tgArr = append(tgArr, gID)
		}
		TerminalActionsTargetsFromID[id] = tgArr
	}
	return
}

/* Сохранить массив действий в файл - НЕ ИСПОЛЬЗУЕТСЯ */
func SaveTerminalActons() {
	var out = ""

	// сохранение только в режиме личинки Larva
	if gomeostas.EvolushnStage > 0 {
		return
	}

	keys := make([]int, 0, len(TerminalActonsNameFromID))
	for k := range TerminalActonsNameFromID {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, id := range keys {
		out += strconv.Itoa(id) + "|" + TerminalActonsNameFromID[id] + "|" +
			GetListTerminalActionsExpenses(id) + "|" +
			strings.Join(lib.StrArrToIntArr(TerminalActionsTargetsFromID[id]), ",") + "\r\n"
	}
	lib.WriteFileContent(lib.GetMainPathExeFile()+"/memory_reflex/terminal_actons.txt", out)
}

/* Получить строку результатов ответных действий */
func GetListTerminalActionsExpenses(actID int) string {
	var out = ""

	se := TerminalActionsExpensesFromID[actID]
	if se != nil {
		for i := 0; i < len(se); i++ {
			out += strconv.Itoa(se[i].GomeoID) + ">" + fmt.Sprintf("%.1f", se[i].Diff) + ";"
		}
	}
	out = strings.TrimSuffix(out, ";")

	return out
}

/* Обновить список БП, которые улучшаются при действии */
func UpdateActionsTargetsFromID(id int, txt string) {
	var tArr []int

	sArr := strings.Split(txt, ",")
	for i := 0; i < len(sArr); i++ {
		if len(sArr[i]) == 0 {
			continue
		}
		gID, _ := strconv.Atoi(sArr[i])
		tArr = append(tArr, gID)
	}
	TerminalActionsTargetsFromID[id] = tArr
}

/* Обновить массив затратных действий */
func UpdateTerminalActionsExpenses(id int, txt string) {
	var expenses []TerminalActionsExpenses

	se := strings.Split(txt, ";")
	for i := 0; i < len(se); i++ {
		if len(se[i]) == 0 {
			continue
		}
		e := strings.Split(se[i], ">")
		e1, _ := strconv.Atoi(e[0])
		e2, _ := strconv.ParseFloat(e[1], 64)
		expenses = append(expenses, TerminalActionsExpenses{e1, e2})
	}
	TerminalActionsExpensesFromID[id] = expenses
}

// уже выполнявшиеся при бодствовании простейшие рефлексы
var usedSimpleReflexexsID []int

func pulsSimpleReflexex() {
	if IsSlipping && len(usedSimpleReflexexsID) > 0 {
		usedSimpleReflexexsID = nil
	}
}

/* выбрать подходящий простейший рефлекс и вернуть его действие
TerminalActions.ChooseSimpleReflexexAction
Eсли нет условного и безусловного рефлекса, то совершается самый простейший безусловный рефлекс
		по сочетаниям редактора http://go/pages/terminal_actions.php
		Данный редактор связывает действие с тем, какие гомео-параметры улучшает данное действие.
*/
func ChooseSimpleReflexexAction() (bool, int, []int) {
	// выявить ID парамктров гомеостаза как цели для улучшения в данных условиях
	veryActual, targetID := gomeostas.FindTargetGomeostazID()
	// подходящие действия
	var aArr []int
	// выбранные
	var fActsID []int
	// одно из выбранных
	var singleActID = 0
	// выдать массив возможных действий чтобы выбрать одно из них, пока еще не испытанное
	for id, gIDarr := range TerminalActionsTargetsFromID {
		if id > 0 {
			// выбрать подходящие ID параметров гомеостаза для данной цели
			aArr = lib.GetExistsIntArs(targetID, gIDarr)
			if aArr == nil {
				continue
			}
			for i := 0; i < len(aArr); i++ {
				// исключить те, что уже использовались
				if lib.ExistsValInArr(usedSimpleReflexexsID, aArr[i]) {
					continue
				}
				//fActsID = append(fActsID, aArr[i])
				fActsID = append(fActsID, id) // !!! id действия, а не параметр гомеостаза!
			}
		}
	}
	if len(fActsID) == 0 {
		if len(aArr) == 0 {
			return false, 0, nil
		}
		// снова обнулить usedSimpleReflexexsID
		singleActID = aArr[0]
		usedSimpleReflexexsID = nil
	} else {
		singleActID = fActsID[0]
	}
	usedSimpleReflexexsID = append(usedSimpleReflexexsID, singleActID)
	return veryActual, singleActID, targetID
}
//////////////////////////////////////////////////

/* выдать массив возможных действий по ID парамктров гомеостаза как цели для улучшения в данных условиях */
func GetSimpleActionForCurContitions() ([]int, []int) {
	// выявить ID парамктров гомеостаза как цели для улучшения в данных условиях
	_, targetArrID := gomeostas.FindTargetGomeostazID()
	if targetArrID == nil {return nil, nil} // если целей нет, незачем дальше проверять - будет nil

	// выявить ID параметров гомеостаза как цели для улучшения в данных условиях
	var fActsID []int
	// выдать массив возможных действий чтобы выбрать одно из них, пока еще не испытанное
	for id, gIDarr := range TerminalActionsTargetsFromID {
		if id > 0 {
			// выбрать подходящие ID параметров гомеостаза для данной цели
			aArr := lib.GetExistsIntArs(targetArrID, gIDarr)
			if aArr == nil {
				continue
			}
			for i := 0; i < len(aArr); i++ {
				fActsID = append(fActsID, id) // !!! id действия, а не параметр гомеостаза!
			}
		}
	}

	return targetArrID, fActsID
}