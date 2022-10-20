/*  Моторные дейсвтвия автоматизма

Для каждого действия brain\reflexes\terminete_action.go задается "сила" действия в градации от 1 до 10, которая передается наПульт словами:
Максимально (сила=10), wwww (сила=8)", "Очень сильно (сила=9), ... Едва (сила=1).
При этом пропорционально расходуется энергия и могут происходить другие изменения гоместаза.
Такой результат сопоставляется с допустимым сразу при действии и корректируется установкой рефлекса мозжечка.

Две области моторного терминала уровня психики:
Область Брока VerbalFromIdArr=make(map[int]*Verbal)
отвечает за смысл распознанных слов и словосочетаний,
за конструирование собственных словосочетаний,
за моторное использование сло и словосочетаний.
За все ответственная структура - образ осмысленных слов и сочетаний.

Область моторных действий ActivityFromIdArr=make(map[int]*Activity)
отвечает за смысл распознанных действий с Пульта,
за конструирование собственных последовательностей действий,
за моторное использование действий.
За все ответственная структура - образ осмысленных действий и их сочетаний.

*/

package psychic

import (
	"BOT/brain/gomeostas"
	_ "BOT/brain/gomeostas"
	termineteAction "BOT/brain/terminete_action"
	word_sensor "BOT/brain/words_sensor"
	"BOT/lib"
)

// при срабатывании автоматизма - блокируются все рефлексторные действия
var MotorTerminalBlocking=false

////////////////////////////////////////////////


// момент запуска автоматизма в числе пульсов
var LastRunAutomatizmPulsCount =0 //сбрасывать ожидание результата автоматизма если прошло 20 пульсов
// ожидается результат запущенного MotAutomatizm
var LastAutomatizmWeiting *Automatizm
// активный узел дерева в момент запуска автоматизма
var LastDetectedActiveLastNodID=0


/////////////////////////////////////////
/* запуск автоматизма на выполнение
возвращает true при успехе
 */
func RumAutomatizmID(id int)(bool){
	a:=AutomatizmFromIdArr[id]
	if a==nil{
		return false
	}
	return RumAutomatizm(a)
}
////////////////////
// todo = true - выполнить полюбому,
func RumAutomatizm(am *Automatizm)(bool){
	if am==nil{
		return false
	}
	if MotorTerminalBlocking { //блокировка моторных терминалов во сне или произвольно
		return false
	}

// NotAllowAnyActions ставится тогда, когда сохранение памяти должно выполняться в тишине, в бездействии
	if  NotAllowAnyActions{
		return false
	}
	if am.ActionsImageID==0{
		return false
	}

// блокировка выполнения плохого автоматизма, если только не применена "ИЛА ВОЛИ"
if am.Usefulness<0{

	return false
}

	var out="3|"

	out+=GetAutomotizmActionsString(am)

	lib.SentActionsForPult(out)

	//выполнить мозжечковый рефлекс сразу после выполняющегося автоматизма
	runCerebellumAdditionalAutomatizm(0,am.ID)

	MotorTerminalBlocking=true // блокировка рефлексов
	LastRunAutomatizmPulsCount =PulsCount // активность мот.автоматизма в чисде пульсов
	LastAutomatizmWeiting=am
	LastDetectedActiveLastNodID=detectedActiveLastNodID
	/* Блокировать выполнение рефлексов на время ожидания результата автоматизма
	вызывается из reflex_action.go рефлексов
	*/
	//isReflexesActionBloking=true // отмена в automatizm_result.go или просто isReflexesActionBloking=false

	return true
}
//////////////////////////////////////////


func GetAutomotizmActionsString(am *Automatizm)(string){
	var out=""
	ai:=ActionsImageArr[am.ActionsImageID]
	if ai.ActID != nil {
		// учесть рефлекс мозжечка
		addE := getCerebellumReflexAddEnergy(0,am.ID)
		sumEnergy:=am.Energy+addE
		if sumEnergy>10{
			sumEnergy=10
		}
		if sumEnergy<1{
			sumEnergy=1
		}
		out+=TerminateMotorAutomatizmActions(ai.ActID,sumEnergy)
	}

	if ai.PhraseID != nil {
		addE := getCerebellumReflexAddEnergy(0,am.ID)
		out+=TerminatePraseAutomatizmActions(ai.PhraseID,am.Energy+addE)
	}

	if ai.ToneID != 0 {
		out+="<br>"+getToneStrFromID(ai.ToneID)+"<br>"
	}

	if ai.MoodID != 0 {
		out+="<br>"+getMoodStrFromID(ai.MoodID)+"<br>"
	}
	return out
}
/////////////////////////////////////////




/* совершить МОТОРНОЕ (http://go/pages/terminal_actions.php) действие  - Dnn-часть автоматизма (не фраза)
cила действия сначала задается =5, а потот корректируется мозжечковыми рефлексами
Использование: 	TerminateMotorAutomatizmActions(actIDarr,energy)
 */
func TerminateMotorAutomatizmActions(actIDarr []int,energy int)string{
	// energy=1
	// название силы:
	var enegrName=""
	if energy < len(termineteAction.EnergyDescrib) {
		enegrName = termineteAction.EnergyDescrib[energy]
	}
	var out=""
	var isAct=false
	for i := 0; i < len(actIDarr); i++ {
		if len(out) > 0{
			out += ", "
		}
		// при моторном действии  меняются гомео-параметры:
		expensesGomeostatParametersAfterAction(actIDarr[i],energy)
		// выдать на Пульт:
		actName:= termineteAction.TerminalActonsNameFromID[actIDarr[i]]
		// ЭНЕРГИЧНОСТЬ
		switch energy{
		case 1:
			out +="<span style=\"font-size:10px;\">"+actName+"</span>"
		case 2:
			out +="<span style=\"font-size:11px;\">"+actName+"</span>"
		case 3:
			out +="<span style=\"font-size:12px;\">"+actName+"</span>"
		case 4:
			out +="<span style=\"font-size:13px;\">"+actName+"</span>"
		case 5:
			out +="<span style=\"font-size:14px;\">"+actName+"</span>"
		case 6:
			out +="<span style=\"font-size:14px;\"><b>"+actName+"<b></span>"
		case 7:
			out +="<span style=\"font-size:17px;color:#927ACC\"><b>"+actName+"<b></span>"
		case 8:
			out +="<span style=\"font-size:19px;color:#E8A7A7\"><b>"+actName+"<b></span>"
		case 9:
			out +="<span style=\"font-size:21px;color:#E86966\"><b>"+actName+"<b></span>"
		case 10:
			out +="<span style=\"font-size:25px;color:#FF0000\"><b>"+actName+"<b></span>"
		}
		isAct=true
	}
	if isAct {
		out = "Действие: <b>"+out+"</b><br><span style=\"font-size:14px;\">Энергичность: <b>" + enegrName+"</b></span><br>"
		return out
	}
	return ""
}

/* совершить МОТОРНОЕ (ВЫДАТЬ ФРАЗУ) действие - Snn-часть автоматизма
cила действия сначала задается = 5, а потот корректируется мозжечковыми рефлексами
*/
func TerminatePraseAutomatizmActions(IDarr []int, energy int)string{
	// при моторном действии  меняются гомео-параметры:
	// expensesGomeostatParametersAfterAction(aI) болтать можно без устали?

	// выдать на ПУльт
	var out = "Фраза Beast: "
	for i := 0; i < len(IDarr); i++ {
		prase := word_sensor.GetPhraseStringsFromPhraseID(IDarr[i])
		out += "<b>"+prase+"</b>"
	}
	// название силы:
	if energy < len(termineteAction.EnergyDescrib) {
		out += " " + termineteAction.EnergyDescrib[energy] + "</b>"
	}
	return out
}

/* изменение гомео-параметров при действии
сила действия корректирует воздействие на параметр гомеостаза
*/
func expensesGomeostatParametersAfterAction(actID int,energy int){
	se :=termineteAction.TerminalActionsExpensesFromID[actID]
	if se != nil {
		for j := 0; j < len(se); j++ {
			// (2*aI.Energy/10) при силе==5 коэффициент будет 1, при силе==10 воздействие увеличиться в 2 раза
			if !gomeostas.NotAllowSetGomeostazParams{
				k := float64(2 * energy / 10)
				gomeostas.GomeostazParams[se[j].GomeoID] += se[j].Diff * k
				if gomeostas.GomeostazParams[se[j].GomeoID] > 100{
					gomeostas.GomeostazParams[se[j].GomeoID] = 100
				}
				if gomeostas.GomeostazParams[se[j].GomeoID] < 0{
					gomeostas.GomeostazParams[se[j].GomeoID] = 0
				}
			}
		}
	}
}
/////////////////////////////////////////////////////////////////////


