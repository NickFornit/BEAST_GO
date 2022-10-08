/* распознаватели отклонения жизненных параметров GomeostazParams[id] от нормы

*/

package gomeostas

import (
	"BOT/lib"
	"math"
	"strconv"
	"strings"
)

/* детектор общего базового состояния:
1 - Похо, 2 - Норма, 3 - Хорошо (==выход из состояния Плохо)
*/
var CommonBadNormalWell = 2  // из других пакетов: gomeostas.CommonBadNormalWell
// название Базового состояния из его ID str:=gomeostas.getBaseCondFromID(id)
func GetBaseCondFromID(id int)(string){
	var out = ""
	switch id{
	case 1: out = "Плохo"
	case 2: out = "Норма"
	case 3: out = "Хорошо"
	}
	return out
}

// величина - насколько плохо
var CommonBadValue = 0
// момент появления Хорошо (в тиках пульса)
var CommonWellValueStart = 0
// предыдущее значение Плохо (CommonBadValue)
var CommonOldBadValue = 0

/* детекторы по каждому из жизненных параметров
имеют те же значения, что и в CommonBadNormalWell
Алгоритм на примере энергии.
Если энергия истощилась, то BadValue[id] будет тем отрицательнее, чем сильнее истощилась.
Если началось восполнение энергии, то BadNormalWell[id]=3 (хорошо) и BadValue[id] уменьшается по мере насыщения.
Но если насыщение остановилось (не меняется в течение ),
то через время BadValue[id] снова становится ==1 (плохо), если параметр еще не восстановлен: остался еще голоден.
В природе чем сильнее голод, тем сильнее Хорошо с началом его удовлетворения и это Хорошо уменьшается с насыщением.
Но если еды было мало, то тварь довольно скоро опять почувствует голод, но уже не такой большой.
*/

var BadNormalWell = make(map[int]int) 	// 1 - Похо, 2 - Норма, 3 - Хорошо (==выход из состояния Плохо)
var BadValue = make(map[int]int) 				// насколько плохо, значение от  0 до 10(максимально Плохо)
var wellValueStart = make(map[int]int)	// время возникновения состояния Хорошо
var oldBadValue = make(map[int]int)			// предыдущие значения BadValue

// пороги начала выхода параметров из нормы
var compareLimites=make(map[int]int)

func initBadDetector() {
	for id, _ := range GomeostazParams {
		BadNormalWell[id] = 2 // норма
		BadValue[id] = 0
		wellValueStart[id] = 0
		oldBadValue[id] = 0
	}
	// порог выхода из нормы жиз.параметров:
	path := lib.GetMainPathExeFile()
	lines,_ := lib.ReadLines(path + "/memory_reflex/GomeostazLimits.txt")
	for i := 0; i < len(lines); i++ {
		p := strings.Split(lines[i], "|")
		id, _ := strconv.Atoi(p[0])
		limit, _ := strconv.Atoi(p[1])
		compareLimites[id] = limit
	}
	return
}

// время в секундах удержания состояния Хорошо для возврата в Плохо или Норму
var dinamicTimeFromBad = 50
// эти параметры стали лучше
var curGomeoParIdSuccesArr []int // временный массив, на время измерения, чтобы никогда не было обнуления
var GomeoParIdSuccesArr []int

/* Определение текущего состояния по каждому базовому параметру
затем общее интегральное состояние */
func badDetecting() {
	if NotAllowSetGomeostazParams {	return }
	curGomeoParIdSuccesArr = nil

	for id, _ := range GomeostazParams {
		cur := int(GomeostazParams[id])
		// только для энергии
		if id == 1 {
			var isBad=false
			if cur < compareLimites[id] { // вышел из нормы
				isBad = true
			}
			detector(id, cur, isBad) // насколько Плохо
			// изменения состояния Плохо (в норме не нужны)
			if PulsCount > 2 { // чтобы было уже установленные значения PulsCount!!!
				dif := oldBadValue[id] - BadValue[id]
				if (id == 1 && dif < -5) || (id > 0 && dif > 5) {
					BadNormalWell[id] = 3 // хорошо
					wellValueStart[id] = PulsCount
					curGomeoParIdSuccesArr = append(curGomeoParIdSuccesArr, id)
				}
			}
			if oldBadValue[id] > BadValue[id]{ // стало хуже
				BadNormalWell[id] = 1 // плохо
				wellValueStart[id] = 0
			}
		}
		// для остальных базовых параметров
		if id > 1 {
			var isBad = false
			if cur > compareLimites[id] { // вышел из нормы - для остальных
				isBad = true
			}
			detector(id, cur, isBad) // насколько Плохо
			//  изменения состояния Плохо (в норме не нужны)
			if (oldBadValue[id] > BadValue[id]) || (oldBadValue[id]>0 && BadValue[id]==0) {// стало лучше
				BadNormalWell[id] = 3 // хорошо
				wellValueStart[id] = PulsCount
				curGomeoParIdSuccesArr = append(curGomeoParIdSuccesArr, id)
			}
			if oldBadValue[id] < BadValue[id]{// стало хуже
				BadNormalWell[id] = 1 // плохо
				wellValueStart[id] = 0
			}
		}
		oldBadValue[id] = BadValue[id]
	}
	commonBadDetecting()
}

/* Определение насколько Плохо как % умноженный	на 10
BadValue[id] - все нули кроме Плохих параметров, у которых Плохо от 1 до 10
 */
func detector(id int, cur int, isBad bool) {
	if NotAllowSetGomeostazParams {	return }
	if isBad { // вышел из нормы
		// BadValue насколько Плохо - как процент, только умноженный не на 100, а на 10 и округленной дробной части:
		diff := compareLimites[id] - cur
		BadValue[id] = lib.Round(float64((diff * 10) / compareLimites[id]))
		if id > 1 { // кроме энергии - поменять знак!
			BadValue[id] *= -1
		}
		if BadNormalWell[id] == 3 { // было Хорошо - удерживать Хорошо dinamicTimeFromBad сек
			// состояние Хорошо протухло, нужно менять на...
			if wellValueStart[id]+dinamicTimeFromBad < PulsCount {
				BadNormalWell[id] = 1 // плохо
			}
		} else {
			BadNormalWell[id] = 1 // плохо
		}
	} else { // если было Хорошо, то удерживать его на время dinamicTimeFromBad
		BadValue[id] = 0	// обнулить насколько Плохо
		if BadNormalWell[id] == 3 { // удерживать Хорошо dinamicTimeFromBad сек
			// состояние Хорошо протухло, нужно менять на...
			if wellValueStart[id] + dinamicTimeFromBad < PulsCount {
				BadNormalWell[id] = 2 // норма
			}
		} else {
			BadNormalWell[id] = 2 // норма
		}
	}
	return
}

// порог начала состояния Плохо
var compareLevel = 100
/* Распознавание CommonBadNormalWell
пороговый (compareLevel) сумматор значений состояний Плохо
Логика работы:
1. Если суммарное значение Плохо выше порога - базовое состояние Плохо
2. Если суммарное значение Плохо ниже порога и:
	2.1. если предыдущее базовое состояние было Плохо - базовое состояние Хорошо
	2.2. если предыдущее базовое состояние было Норма или Хорошо - базовое состояние Норма
3. Если базовое состояние Хорошо держится больше dinamicTimeFromBad (50 сек) - базовое состояние Норма
*/
func commonBadDetecting() {
	if NotAllowSetGomeostazParams {	return }
	CommonBadValue=0
	commonPerception=0

	for id, v := range GomeostazParams {
		// насколько плохо умножаем на вес значимости параметра
		CommonBadValue += BadValue[id] * GomeostazParamsWeight[id]
		if id == 1{ // у энергии - чем больше значение параметра - тем лучше
			commonPerception += int(v) * GomeostazParamsWeight[id]
		} else { // чем больше значение параметра - тем хуже
			commonPerception += int(100 - v) * GomeostazParamsWeight[id]
		}
	}
	if CommonBadValue > compareLevel {
		if CommonBadNormalWell == 3 { // удерживать Хорошо dinamicTimeFromBad сек
			// состояние Хорошо протухло, нужно менять на...
			if CommonWellValueStart + dinamicTimeFromBad < PulsCount {
				CommonBadNormalWell = 2 // норма
			}
		} else {
			CommonBadNormalWell = 1 // Плохо
		}
	} else {
		if CommonBadNormalWell == 3 { // удерживать Хорошо dinamicTimeFromBad сек
			// состояние Хорошо протухло, нужно менять на...
			if CommonWellValueStart + dinamicTimeFromBad < PulsCount {
				CommonBadNormalWell = 2 // норма
			}
		} else {
			CommonBadNormalWell = 2 // норма
		}
	}
	// изменения состояния Плохо (в норме не нужны)
	if (CommonOldBadValue > CommonBadValue) || (CommonOldBadValue>0 && CommonBadValue == 0) { // стало лучше
		CommonBadNormalWell = 3 // хорошо
		CommonWellValueStart = PulsCount
	}
	if CommonOldBadValue < CommonBadValue{ // стало хуже
		CommonBadNormalWell = 1 // плохо
		CommonWellValueStart = 0
	}
	return
}

// для Пульта
func GetCurGomeoStatus()string {
	var out = "0;" + strconv.Itoa(CommonBadNormalWell) + "|"
	for id, v := range BadNormalWell {
		out += strconv.Itoa(id) + ";" + strconv.Itoa(v) + "|"
	}
	return out
}

var lastBetterOrWorse = 0 // стали лучше или хуже: величина измнения от -10 через 0 до 10
var curBetterOrWorsePulsCount = 0 // значение пульса для фиксации изменения стало лучше/хуже
/* для Психики: стало хуже или лучше теперь
ОСОБЕННОСТЬ: срабатывание только при изменении параметра (любого) от нормы в Плохо или наоборот,
просто изменение Нормы не влияет, т.е. реально индицирует, что был выход Б.параметра из Нормы или возврат в Норму
и то, на какую условную величину были изменения, влияющие на осознание настроения PsyBaseMood
Возвращает величину измнения от -10 через 0 до 10
Меняться может не чаще раз в 1 пульс
Значение экспоненциально стремиться к пределам -10 и 10
Это имитирует ограничение природных распознавателей на число дискретов распознавания.

Сканируется с каждым пульсом постоянно
*/
func prepBetterOrWorseNow() {
	diff := int((CommonOldBadValue - CommonBadValue) / 10)
	if diff != 0 && curBetterOrWorsePulsCount != PulsCount {
		// Значение экспоненциально стремится к пределам -10 и 10
		lastBetterOrWorse = int(10.0 - 10.0 / math.Exp(float64(lib.Abs(diff)) * 0.17))
		if diff < 0{
			lastBetterOrWorse *= -1
		}
		if lastBetterOrWorse > 0 { // если стало лучше, то и показывать GomeoParIdSuccesArr
			GomeoParIdSuccesArr = curGomeoParIdSuccesArr
		} else {
			GomeoParIdSuccesArr = nil
		}
		curBetterOrWorsePulsCount = PulsCount
	}
	return
}

// текущее общее ощущение: сумма произведений параметров на их веса
var commonPerception = 0 // постоянно обновляемое значение
// предыдущее общее ощущение
var commonOldPerception = 0 // меняется только по запросам психики функции BetterOrWorseNow()
// насколько изменилось общее состояние, значение от  -10(максимально Плохо) через 0 до 10(максимально Хорошо)
var commonDiffValue = 0
var curcommonOldPerceptionPulsCount=0

/* Экспоненциальная оценка изменения общего состояния */
func commonPerceptionNow() {
	// !!! чтобы точно успело учесть изменения параметров
	commonBadDetecting()
	if commonOldPerception == 0{
		/* В самый первый момент старое значение просто должно равняться текущему вот почему.
		Всегда делается два измерения: 1 - в момент срабатывания автоматизма и потом проверяется изменения от действий.*/
		commonOldPerception = commonPerception
	}
	diff := int((commonPerception - commonOldPerception) / 10)
	if  curcommonOldPerceptionPulsCount != PulsCount{
		// Значение экспоненциально стремится к пределам -10 и 10
		commonDiffValue = int(10.0 - 10.0 / math.Exp(float64(lib.Abs(diff)) * 0.17))
		if diff < 0 {
			commonDiffValue *= -1
		}
		curcommonOldPerceptionPulsCount = PulsCount
	}
	return
}

/* вызывается из психики res:=gomeostas.BetterOrWorseNow()
Сканируется с каждым пульсом в func automatizmActionsPuls() во время ожидания
CommonMoodAfterAction - "+" действия оператора привели к позитиву, "-" к негативу
ВОЗВРАЩАЕТ:
commonDiffValue - насколько изменилось общее состояние, значение от -10(максимально Плохо) через 0 до 10(максимально Хорошо)
lastBetterOrWorse - стали лучше или хуже: величина измнения от -10 через 0 до 10
GomeoParIdSuccesArr - стали лучше следующие г.параметры []int гоменостаза

Если было очень плохо, а стало не очень плохо, то commonDiffValue станет позитивным.
 */
func BetterOrWorseNow(kind int)(int, int, []int) {
	// насколько изменилось общее состояние
	commonPerceptionNow()
	commonOldPerception = commonPerception
	// перед перекрытием старого значения: стало хуже или лучше теперь
	prepBetterOrWorseNow()
	CommonOldBadValue = CommonBadValue

	// применить "эффект" кнопок с Пульта (в таблице он забивался в виде "+" "-")
	if kind == 2 { // второй вызов при измерении эффекта реакции
		if commonDiffValue == 0 {
			if CommonMoodAfterAction == "+" {
				commonDiffValue = 1
			}
			if CommonMoodAfterAction == "-" {
				commonDiffValue = -1
			}
		}
	}

	return commonDiffValue, lastBetterOrWorse, GomeoParIdSuccesArr
}