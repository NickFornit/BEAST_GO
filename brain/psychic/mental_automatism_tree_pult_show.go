/*  Выдать на Пульт дерево ментальных автоматизмов
id|Usefulness|ActionsImageID|Count
*/


package psychic

import (
	"BOT/brain/gomeostas"
	"strconv"
)

////////////////////////////////////////////

/*запрет показа карты на пульте (func GetAutomatizmTreeForPult()) при обновлении
против паники типа "одновременная запись и считывание карты"
Использовать для всех операций записи узлов дерева
*/


// образ дерева автоматизмов для вывода
var automatizmMentalTreeModel=""
/////////////////////////////////////////////////
func GetMentalAutomatizmTreeForPult(limit int)(string){
	// против паники типа "одновременная запись и считывание карты"
	if notAllowScanInTreeThisTime{
		return "!Временно запрещена работа func GetAutomatizmTreeForPult() т.к. идет параллельная обработка."
	}
	if len(UnderstandingTree.Children)==0 { // еще нет никаких веток
		return "Еще нет Дерева понимания"
	}

//посмотреть число имеющихся узлов дерева
	//strArr,_:=lib.ReadLines(lib.GetMainPathExeFile()+"/memory_psy/automatizm_tree.txt")

	scanMentalAutomatizmNodes(-1, &UnderstandingTree)

	if len(automatizmMentalTreeModel)<10{
		return "Еще нет информационных веток дерева"
	}

	return automatizmMentalTreeModel

}
//////////////////////
// ID|ParentNode|Mood|EmotionID|SituationID|PurposeID
func scanMentalAutomatizmNodes(level int,node *UnderstandingNode) {

	if node.ID > 0 {
		automatizmMentalTreeModel +="<span style='color:#666666;'>"+strconv.Itoa(node.ID)+":</span> "

		automatizmMentalTreeModel += setOutShift(level)

		switch level {
		case 0:// Mood
			automatizmMentalTreeModel +=getMoodStr(node)
		case 1: // EmotionID
			automatizmMentalTreeModel +=getEmotionStr(node)
		case 2: // SituationID
			automatizmMentalTreeModel +=getSituationStr(node)
		case 3: // PurposeID
			automatizmMentalTreeModel +=getPurposeStr(node)
		}
		automatizmMentalTreeModel +="<br>\n"
	}
		level++
		for n := 0; n < len(node.Children); n++ {

			scanMentalAutomatizmNodes(level, &node.Children[n])
		}
}
//////////////////////////////////////////////////////////
// отступ
func setOutShift(level int)(string){
	var sh=""
	for n := 0; n < level; n++ {
		sh+="&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;"
	}
	return sh
}
////////////////////////////////////////////////////




/////////////////////////////////////////////////////
func GetMentalSituationsForNodeInfo(sID int)(string){

	st:=SituationImageFromIdArr[sID]
	if st==nil{
		return "Нет образа ситуации с ID="+strconv.Itoa(sID)
	}

	switch st.SituationType{
	case 1: return "Было ответное действие"
	case 2: return "Был запуск автоматизма ветки"
	case 3: return "Ничего не делали, но нужно осмысление"
	case 4: return "Все спокойно, можно экспериментивароть"
	case 5: return "Оператор не прореагировал на действия в течение периода ожидания"

	case 11: return "Оператор выбрал настроение Хорошее"
	case 12: return "Оператор выбрал настроение Плохое"
	case 13: return "Оператор выбрал настроение Игровое"
	case 14: return "Оператор выбрал настроение Учитель"
	case 15: return "Оператор выбрал настроение Агрессивное"
	case 16: return "Оператор выбрал настроение Защитное"
	case 17: return "Оператор выбрал настроение Протест"



	//case 20: return "Оператор нажал кнопку "
	case 21: return "Оператор нажал кнопку Непонятно"
	case 22: return "Оператор нажал кнопку Понятно"
	case 23: return "Оператор нажал кнопку Наказать"
	case 24: return "Оператор нажал кнопку Поощрить"
	case 25: return "Оператор нажал кнопку Накормить"
	case 26: return "Оператор нажал кнопку Успокоить"
	case 27: return "Оператор нажал кнопку Предложить поиграть"
	case 28: return "Оператор нажал кнопку Предложить поучить"
	case 29: return "Оператор нажал кнопку Игнорировать"
	case 30: return "Оператор нажал кнопку Сделать больно"
	case 31: return "Оператор нажал кнопку Сделать приятно"
	case 32: return "Оператор нажал кнопку Заплакать"
	case 33: return "Оператор нажал кнопку Засмеяться"
	case 34: return "Оператор нажал кнопку Обрадоваться"
	case 35: return "Оператор нажал кнопку Испугаться"
	case 36: return "Оператор нажал кнопку Простить"
	case 37: return "Оператор нажал кнопку Вылечить"
	}

	return ""
}
///////////////////////////////////////////////////////
func GetMentalPurposeForNodeInfo(pID int)(string){
	out:=""
	pp:=PurposeImageFromID[pID]
	if pp==nil{
		return "Нет образа цели с ID="+strconv.Itoa(pID)
	}
	if pp.veryActual{
		out+="Цель очень актуальна<br>"
	}
	if len(pp.targetID)>0{
		out+="Улучшение параметров гомеостаза:<br>"
		for i := 0; i < len(pp.targetID); i++ {
if i>0{out+=", "}
			out+=gomeostas.GetBaseContextCondFromID(pp.targetID[i])
		}
		out+="<br>"
	}
	if pp.actionID>0{
		out+="Выбранный образ действия для достижения данной цели:<br>"
		out+=GetActionsString(pp.actionID)
	}

	return out
}
///////////////////////////////////////////////////////
// Mood
func getMoodStr(node *UnderstandingNode)(string){
	moodS := ""
	switch node.Mood {
	case -1:
		moodS = "Плохо"
	case 0:
		moodS = "Норма"
	case 1:
		moodS = "Хорошо"
	}
	out := "Состояние: <b> " + moodS + "</b>"

	return out
}

// EmotionID
func getEmotionStr(node *UnderstandingNode)(string){
em := EmotionFromIdArr[node.EmotionID]
	out := "Эмоция: <b> " + getEmotonsComponentStr(em) + "</b>"
	return out
}

// SituationID
func getSituationStr(node *UnderstandingNode)(string){
pp := SituationImageFromIdArr[node.SituationID]
out:=""
if pp == nil {
	out += "<span style='color:red'>Нет образа цели с ID=" + strconv.Itoa(node.SituationID)+"</span>"
} else {
	out += "Ситуация: <b> <span style='cursor:pointer;color:blue' onClick='get_situation(" + strconv.Itoa(node.SituationID) + ")'>" + strconv.Itoa(node.SituationID) + "</span>" + "</b>"
}
	return out
}

// PurposeID
func getPurposeStr(node *UnderstandingNode)(string){
pp:=PurposeImageFromID[node.PurposeID]
out:=""
if pp==nil {
	out += "<span style='color:red'>Нет образа цели с ID=" + strconv.Itoa(node.PurposeID)+"</span>"
} else {
	out += "Цель: <b> <span style='cursor:pointer;color:blue' onClick='get_purpose(" + strconv.Itoa(node.PurposeID) + ")'>" + strconv.Itoa(node.PurposeID) + "</span>" + "</b>"
}
	return out
}
///////////////////////////////////////////////////////