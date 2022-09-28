/*  Выдать на Пульт дерево автоматизмов

*/


package psychic

import (
	actionSensor "BOT/brain/action_sensor"
	"BOT/brain/gomeostas"
	termineteAction "BOT/brain/terminete_action"
	word_sensor "BOT/brain/words_sensor"
	"strconv"
	"strings"
)

////////////////////////////////////////////

/*запрет показа карты на пульте (func GetAutomatizmTreeForPult()) при обновлении
против паники типа "одновременная запись и считывание карты"
Использовать для всех операций записи узлов дерева
*/
var notAllowScanInTreeThisTime=false


// образ дерева автоматизмов для вывода
var automatizmTreeModel=""
/////////////////////////////////////////////////
func GetAutomatizmTreeForPult()(string){
	// против паники типа "одновременная запись и считывание карты"
	if notAllowScanInTreeThisTime{
		return "!Временно запрещена работа func GetAutomatizmTreeForPult() т.к. идет параллельная обработка."
	}
	if len(AutomatizmTree.Children)==0 { // еще нет никаких веток
		return "Еще нет Дерева автоматизмов"
	}
	automatizmTreeModel=""
	scanAutomatizmNodes(-1,&AutomatizmTree)
	if len(automatizmTreeModel)<10{
		return "Еще нет информационных веток дерева"
	}

	return automatizmTreeModel
}
//////////////////////

func scanAutomatizmNodes(level int,node *AutomatizmNode){
	if node.ID==69{
		node.ID=69
	}
	if node.ID>0 {
		automatizmTreeModel += setShift(level)
		switch level{
		case 0:
			automatizmTreeModel += getStrFromCond(level,node.BaseID) + "(nodeID=" + strconv.Itoa(node.ID) + ")"
		case 1:
			automatizmTreeModel += getStrFromCond(level,node.EmotionID) + "(nodeID=" + strconv.Itoa(node.ID) + ")"
		case 2:
			automatizmTreeModel += getStrFromCond(level,node.ActivityID) + "(nodeID=" + strconv.Itoa(node.ID) + ")"
		case 3:
			automatizmTreeModel += getStrFromCond(level, node.ToneMoodID) + "(nodeID=" + strconv.Itoa(node.ID) + ")"
		case 4:
			automatizmTreeModel += getStrFromCond(level, node.SimbolID) + "(nodeID=" + strconv.Itoa(node.ID) + ")"
		case 5:
			automatizmTreeModel += getStrFromCond(level, node.VerbalID) + "(nodeID=" + strconv.Itoa(node.ID) + ")"

		}

		// если есть штатный автоматизм - показать действия
		atmzm:=AutomatizmBelief2FromTreeNodeId[node.ID]
		if atmzm!=nil{
			automatizmTreeModel += " <span style='color:blue'>АВТОМАТИЗМ(" + strconv.Itoa(atmzm.ID) + "): "+
				TranslateAutomatizmSequence(atmzm) + "</span>"
		}
		automatizmTreeModel +="<br>\n"
	}
	level++
	for n := 0; n < len(node.Children); n++ {
		scanAutomatizmNodes(level,&node.Children[n])
	}
}
// отступ
func setShift(level int)(string){
	var sh=""
	for n := 0; n < level; n++ {
		sh+="&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;"
	}
	return sh
}
////////////////////////////////////////////////////




/////////////////////////////////////////////////////
// из ID образа получить составляющие в виде строк
func getStrFromCond(level int,imgID int)(string){
	var out=""
	switch level{
	case 0:
		if imgID>0 && imgID<4{ out="Состояние: <b>"+gomeostas.GetBaseCondFromID(imgID)+"</b>"
	}else{
		out += "<span style='color:red'>несуществующее Базовое состояние ID="+strconv.Itoa(imgID)+"</span>"
	}
	case 1: // эмоция
		out="Эмоция ("+strconv.Itoa(imgID)+"): <b>"+getStrnameFromBaseImageID(imgID)+"</b>"
	case 2: // действия
		out=getStrnameFromStyleImageID(imgID)
		if len(out)==0{
			return "Нет действий с Пульта "
		}else{
			out="Действия с Пульта: <b>"+out+"</b>"
		}
	case 3: // тон-настроение фразы
		out=getToneStrFromID(imgID)
		if len(out)==0{
			return "Нормальное настроение"
		}else{
			out="Тон-Настроение: <b>"+out+"</b>"
		}
	case 4:// первый символ
		out=word_sensor.GetSynbolFromID(imgID)
		if len(out)==0{
			return "Нет первого символа фразы"
		}else{
			out="Первый символ: <b>&quot;"+out+"&quot;</b>"
		}
	case 5:// фраза
		//vrbal:=VerbalFromIdArr[imgID]
		//if vrbal != nil {
			//out = word_sensor.GetPhraseStringsFromPhraseID(vrbal.PhraseID[0])
		out = word_sensor.GetPhraseStringsFromPhraseID(imgID)
		//out = word_sensor.GetWordFromPraseNodeID(imgID)
			if len(out) == 0 {
				return "Нет фразы"
			} else {
				out = "Фраза: <b>&quot;" + out + "&quot;</b>"
			}
		//}
	}
	return out
}
// названия базовых контекстов в их сочетании -из ID эмоции
func getStrnameFromBaseImageID(id int)(string){
	var out=""
	if EmotionFromIdArr[id]==nil{
		return "Нет эмоций"
	}
	img:=EmotionFromIdArr[id].BaseIDarr
	for i := 0; i < len(img); i++ {
		if i>0{out +=", "}
		name:=gomeostas.GetBaseContextCondFromID(img[i])+""
		out +=name
	}
	if len(out)==0{
		return "Нет эмоций"
	}
	return out
}
// названия Пусковых стимулов в их сочетании -из ID их образа
func getStrnameFromStyleImageID(id int)(string){
	if ActivityFromIdArr[id]==nil{
		return "Нет действий с Пульта"
	}
	var out=""
	img:=ActivityFromIdArr[id].ActID
	for i := 0; i < len(img); i++ {
		if i>0{out +=", "}
		name:=actionSensor.GetActionNameFromID(img[i])+""
			out +=name
	}
	if len(out)==0{
		return "Нет действий с Пульта"
	}
	return out
}
/////////////////////////////////////////////




/////////////////////////////////////////////
/*расшифровать действия автоматизма для инфы пульта: Snn:21812,27777,0,1478,13388,0,27303,24882Dnn:4
Сделано на основе запуска автоматизма на выполнение: func RumAutomatizmID(id int) из automatizm_actions.go
 */
func TranslateAutomatizmSequence(am *Automatizm)(string){
	if am==nil{
		return ""
	}
	if len(am.Sequence)==0{
		return ""
	}

	out:=GetAutomatizmSequenceInfo(am.ID,am.Sequence)

	return out
}
////////////////////////////////////////


// действия - в виде строки
func GetAutomatizmSequenceInfo(idA int,sequence string)(string){

	am:=AutomatizmFromIdArr[idA]
	if am == nil{
		return ""
	}
	var out=""
	actArr:=ParceAutomatizmSequence(sequence)

	for i := 0; i < len(actArr); i++ {
		// строка действий (любого типа) через запятую
		aArr := strings.Split(actArr[i].Acts, ",")
		var idArr []int

		switch actArr[i].Type {
		case 1: // Snn- перечень ID сенсора слов через запятую,
			for n := 0; n < len(aArr); n++ {
				aID, _ := strconv.Atoi(aArr[n])
				idArr = append(idArr, aID)
			}
			addE := 0
			if am.Belief != 3 { // не рефлекс мозжечка
				addE = getCerebellumReflexAddEnergy(am.ID)
			}
			for i := 0; i < len(idArr); i++ {
				if i==0{out+="Фраза Beast c энергией "+strconv.Itoa(am.Energy+addE)+":</b>"}else{out+=", "}
				out+= "<b>"+word_sensor.GetPhraseStringsFromPhraseID(idArr[i])+"</b>"
			}
			//TerminatePraseAutomatizmActions(idArr, am.Energy+addE)
		case 2: //Dnn - ID прогрмаммы действий, через запятую

			for n := 0; n < len(aArr); n++ {
				aID, _ := strconv.Atoi(aArr[n])
				idArr = append(idArr, aID)
			}
			addE := 0
			if am.Belief != 3 { // не рефлекс мозжечка
				addE = getCerebellumReflexAddEnergy(am.ID)
			}
			for i := 0; i < len(idArr); i++ {
				if i==0{out+="Действие Beast c энергией "+strconv.Itoa(am.Energy+addE)+": "}else{out+=", "}
				out+= "<b>"+termineteAction.TerminalActonsNameFromID[idArr[i]]+"</b>"
			}

		case 3: //Ann - последовательный запуск автоматизмов с id1,id2..
			// НО нужно как-то дожидаться выплнения предыдущего до запуска следующего !!!!!!
			// последний просто перекроет все. Лучше выполнять следующий просто по следующему двойному тику пульса??
			out+="<b>Цепочка автоматизмов:</b>"
			for n := 0; n < len(aArr); n++ {
				if i>0{out+=", "}
				out+=aArr[i]
			}
			///////////////////////////////////////
		}
	}


	return out
}