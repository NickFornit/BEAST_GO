/* распознаватели отклонения жизненных параметров GomeostazParams[id] от нормы

*/

package gomeostas

import (
	"BOT/lib"
	"math"
	"strconv"
	"strings"
)

/* детектор общего состояния:
1 - Похо, 2 - Норма, 3 - Хорошо (==выход из состояния Плохо)

// название Базового состояния из его ID: str:=reflexes.getBaseCondFromID(id)
*/
var CommonBadNormalWell=2  // из других пакетов: gomeostas.CommonBadNormalWell
// название Базового состояния из его ID str:=gomeostas.getBaseCondFromID(id)
func GetBaseCondFromID(id int)(string){
	var out=""
	switch id{
	case 1: out="Плохo"
	case 2: out="Норма"
	case 3: out="Хорошо"
	}
	return out
}


// величина - насколько плохо
var CommonBadValue=0
// момент появления Хорошо (в тиках пульса)
var CommonWellValueStart=0
// предыдущее значение Плохо (CommonBadValue)
var CommonOldBadValue=0
//////////////////////////////

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
var BadNormalWell =make(map[int]int) // 1 - Похо, 2 - Норма, 3 - Хорошо (==выход из состояния Плохо)
var BadValue =make(map[int]int) // насколько плохо, значение от  0 до 10(максимально Плохо)
var wellValueStart =make(map[int]int)// время возникновения состояния Хорошо
var oldBadValue =make(map[int]int) // предыдущие значения BadValue



// пороги начала выхода параметров из нормы
var compareLimites=make(map[int]int)
//////////////////////////////

func initBadDetector(){
	for id, _ := range GomeostazParams {
		BadNormalWell[id]=2 // норма
		BadValue[id]=0
		wellValueStart[id]=0
		oldBadValue[id]=0
	}

	// порогb выхода из нормы жиз.параметров:
	path:=lib.GetMainPathExeFile()
	lines,_:=lib.ReadLines(path+"/memory_reflex/GomeostazLimits.txt")
	for i := 0; i < len(lines); i++ {
		p := strings.Split(lines[i], "|")
		id, _ := strconv.Atoi(p[0])
		limit, _ := strconv.Atoi(p[1])
		compareLimites[id]=limit
	}
	/*
// прошивка порогов выхода из нормы жиз.параметров:
	compareLimites[1]=30
	compareLimites[2]=60
	compareLimites[3]=80
	compareLimites[4]=80
	compareLimites[5]=80
	compareLimites[6]=80
	compareLimites[7]=80
	compareLimites[8]=30
	 */
	return
}
//////////////////////////////



// время удерхания состояния Хорошо для возврата в Плохо или Норму
var dinamicTimeFromBad=50 // секунд
// эти параметры стали лучше
var curGomeoParIdSuccesArr []int // временный массив, на время измерения, чтобы никогда не было обнуления
var GomeoParIdSuccesArr []int

/* после Хорошо по PulsCount определяется когда оно истечет и станет Норма
if wellValueStart[id]+dinamicTimeFromBad < PulsCount {
				BadNormalWell[id]=2 // норма
			}
 */
func badDetecting(){
	if NotAllowSetGomeostazParams{
		return
	}
	curGomeoParIdSuccesArr=nil

	for id, _ := range GomeostazParams {
		cur := int(GomeostazParams[id])
		//var isNotNorma = false

//////////
	if id == 1{//// только для энергии:
		var isBad=false
		if cur < compareLimites[id] { // вышел из нормы - для остальных
			isBad = true
		}
		detector(id,cur,isBad)

		//  изменения состояния Плохо (в норме не нужны)
		if PulsCount>2 {// чтобы было уже установленные значения PulsCount!!!
			dif := oldBadValue[id] - BadValue[id]
			if (id == 1 && dif < -5) || (id > 0 && dif > 5) {
				BadNormalWell[id] = 3 // хорошо
				wellValueStart[id] = PulsCount
				curGomeoParIdSuccesArr = append(curGomeoParIdSuccesArr, id)
			}
		}
		if oldBadValue[id] > BadValue[id]{// стало хуже
			BadNormalWell[id]=1 // плохо
			wellValueStart[id]=0
		}
	}
	////////////////////////////////////////////////
	if id > 1{//// для остальных:
		var isBad=false
		if cur > compareLimites[id] { // вышел из нормы - для остальных
				 isBad = true
		}
		detector(id,cur,isBad)

		//  изменения состояния Плохо (в норме не нужны)
		if (oldBadValue[id] > BadValue[id]) || (oldBadValue[id]>0 && BadValue[id]==0) {// стало лучше
			BadNormalWell[id]=3 // хорошо
			wellValueStart[id]=PulsCount
			curGomeoParIdSuccesArr=append(curGomeoParIdSuccesArr,id)
		}
		if oldBadValue[id] < BadValue[id]{// стало хуже
			BadNormalWell[id]=1 // плохо
			wellValueStart[id]=0
		}
	}

		///////////////////
		oldBadValue[id]=BadValue[id]
	}//for id, _ := range GomeostazParams {


	commonBadDetecting()
}
/* BadValue[id] - все нули кроме Плохих параметров, у которых Плохо от 1 до 10
 */
func detector(id int,cur int,isBad bool){
	if NotAllowSetGomeostazParams{
		return
	}
	if isBad { // вышел из нормы
//BadValue насколько Плохо - как процент, только умноженный не на 100, а на 10 и округленной дробной части:
		diff := compareLimites[id] - cur
		BadValue[id] = lib.Round(float64((diff * 10) / compareLimites[id]))
		if id>1{// кроме энергии - поменять знак!
			BadValue[id]*=-1
		}
		if BadNormalWell[id]==3 {// было Хорошо - удерживать Хорошо dinamicTimeFromBad сек
			//  состояние Хорошо протухло, нужно менять на...
			if wellValueStart[id]+dinamicTimeFromBad < PulsCount {
				BadNormalWell[id]=1 // плохо
			}
		}else {
			BadNormalWell[id] = 1 // плохо
		}
		//isNotNorma = true
	} else { // если было Хорошо, то удерживать его на время dinamicTimeFromBad
		BadValue[id] = 0      // обнулить насколько Плохо
		if BadNormalWell[id]==3 {// удерживать Хорошо dinamicTimeFromBad сек
			//  состояние Хорошо протухло, нужно менять на...
			if wellValueStart[id]+dinamicTimeFromBad < PulsCount {
				BadNormalWell[id]=2 // норма
			}
		}else {
			BadNormalWell[id] = 2 // норма
		}
	}
	return
}


/* распознавание CommonBadNormalWell
- пороговый (compareLevel) сумматор значений состояний плохо
 */
var compareLevel=100 // порог начала состояния Плохо
func commonBadDetecting(){
	if NotAllowSetGomeostazParams{
		return
	}

	CommonBadValue=0
	commonPerception=0

	for id, v := range GomeostazParams {
// насколько плохо умножаем на вес значимости параметра
		CommonBadValue+=BadValue[id]*GomeostazParamsWeight[id]

		if id==1{// у энергии - чем больше значение параметра - тем лучше
			commonPerception += int(v) * GomeostazParamsWeight[id]
		}else {// чем больше значение параметра - тем хуже
			commonPerception += int(100-v) * GomeostazParamsWeight[id]
		}
	}

	if CommonBadValue > compareLevel{
		if CommonBadNormalWell==3 {// удерживать Хорошо dinamicTimeFromBad сек
			//  состояние Хорошо протухло, нужно менять на...
			if CommonWellValueStart+dinamicTimeFromBad < PulsCount {
				CommonBadNormalWell=2 // норма
			}
		}else {
			CommonBadNormalWell = 1 // Плохо
		}
	}else{
		if CommonBadNormalWell==3 {// удерживать Хорошо dinamicTimeFromBad сек
			//  состояние Хорошо протухло, нужно менять на...
			if CommonWellValueStart+dinamicTimeFromBad < PulsCount {
				CommonBadNormalWell=2 // норма
			}
		}else {
			CommonBadNormalWell = 2 // норма
		}
	}

	//  изменения состояния Плохо (в норме не нужны)
	if (CommonOldBadValue > CommonBadValue) || (CommonOldBadValue>0 && CommonBadValue==0) {// стало лучше
		CommonBadNormalWell=3 // хорошо
		CommonWellValueStart=PulsCount
	}
	if CommonOldBadValue < CommonBadValue{// стало хуже
		CommonBadNormalWell=1 // плохо
		CommonWellValueStart=0
	}

}
//////////////

///////////////////////////////////////////////



// для Пульта
func GetCurGomeoStatus()(string) {
	var out="0;"+strconv.Itoa(CommonBadNormalWell)+"|"
	for id, v := range BadNormalWell {
		out+=strconv.Itoa(id)+";"+strconv.Itoa(v)+"|"
	}
	return out
}
////////////////////////////////////////////////


/* для Психики: стало хуже или лучше теперь
ОСОБЕННОСТЬ: срабатывание только при изменении параметра (любого) от нормы в плохо или наоборот,
просто изменение нормы не влияет, т.е. реально индицирует, что был выход Б.параметра из нормы или возврат в норму
и то, на какую условную величину были изменения, влияющие на осознание настроения PsyBaseMood
Возвращает величину измнения от -10 через 0 до 10
Меняться может не чаще раз в 1 пульс
Значение экспоненциально стремиться к пределам -10 и 10
Это имитирует ограничение природных распознавателей на число дискретов распознавания.

Сканируется с каждым пульсом постоянно
 */
var lastBetterOrWorse=0
var curBetterOrWorsePulsCount=0
func prepBetterOrWorseNow(){
	diff:= int((CommonOldBadValue - CommonBadValue)/10)
	if diff!=0 && curBetterOrWorsePulsCount!=PulsCount{//
		//Значение экспоненциально стремится к пределам -10 и 10
		lastBetterOrWorse=int(10.0 - 10.0/math.Exp(float64(lib.Abs(diff))*0.17))
		if diff<0{
			lastBetterOrWorse*=-1
		}
		if lastBetterOrWorse>0{// если стало лучше, то и показывать GomeoParIdSuccesArr
			GomeoParIdSuccesArr=curGomeoParIdSuccesArr
		}else{
			GomeoParIdSuccesArr=nil
		}
		curBetterOrWorsePulsCount=PulsCount
	}

	return
}
/////////////////////////////////////
// текущее общее ощущение: сумма произведений параметров на их веса
var commonPerception =0 // постоянно обновляемое значение
// предыдущее общее ощущение
var commonOldPerception =0 // меняется только по запросам психики функции BetterOrWorseNow()
// насколько изменилось общее состояние, значение от  -10(максимально Плохо) через 0 до 10(максимально Хорошо)
var commonDiffValue=0
var curcommonOldPerceptionPulsCount=0

func commonPerceptionNow(){
	// !!! чтобы точно успело учесть изменения параметров
	commonBadDetecting()

	if commonOldPerception==0{
		for id, _ := range GomeostazParams {
			commonOldPerception += 50 * GomeostazParamsWeight[id]
		}
	}
	diff:= int((commonPerception - commonOldPerception)/10)
	if  curcommonOldPerceptionPulsCount!=PulsCount{//
		//Значение экспоненциально стремится к пределам -10 и 10
		commonDiffValue=int(10.0 - 10.0/math.Exp(float64(lib.Abs(diff))*0.17))
		if diff<0{
			commonDiffValue*=-1
		}

		curcommonOldPerceptionPulsCount=PulsCount
	}

	return
}
/////////////////////////////////////
/* вызывается из психики res:=gomeostas.BetterOrWorseNow()
Сканируется с каждым пульсом в func automatizmActionsPuls() во время ожидания
CommonMoodAfterAction - "+" действия оператора привели к позитиву, "-" к негативу
ВОЗВРАЩАЕТ:
commonDiffValue - насколько изменилось общее состояние, значение от  -10(максимально Плохо) через 0 до 10(максимально Хорошо)
lastBetterOrWorse - стали лучше или хуже: величина измнения от -10 через 0 до 10
GomeoParIdSuccesArr - стали лучше следующие г.параметры []int гоменостаза

Если было очень плохо, а стало не очень плохо, то commonDiffValue станет позитивным.
 */
func BetterOrWorseNow()(int,int,[]int){
	// насколько изменилось общее состояние
	commonPerceptionNow()
	commonOldPerception = commonPerception

	//перед перекрытием старого значения: стало хуже или лучше теперь
	prepBetterOrWorseNow()
	CommonOldBadValue=CommonBadValue

	return commonDiffValue,lastBetterOrWorse,GomeoParIdSuccesArr
	}
//////////////////////////////////////////