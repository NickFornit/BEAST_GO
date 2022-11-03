/* Процессы осмысления: создание и использование ментальных автоматизмов
для Дерева понимания (или дерева ментальных автоматизмов)

*/

package psychic

import "strconv"

///////////////////////////////


//////////////////////////////////////////////////
/* Определение Цели в данной ситуации - на уровне дерева понимания
	Здесь выбирается действие пробного автоматизма из выполнившегося рефлекса actualRelextActon
	и запускается МЕНТАЛЬНЫЙ автоматизм.
	На стадии 4 - провоцировать оператора на ответы (почему, зачем, что такое?)
*/
func getPurposeUndestandingAndRunAutomatizm()(bool) {


	// мозжечковые рефлексы - самый первый уровень осознания - подгонка действий под заданную Цель.
	if EvolushnStage == 4 {
		/*  на стадии 4 - провоцировать оператора на ответы (почему, зачем, что такое?)


		 */
	}
/*
	// переосмыслить ситуацию - Активировать Дерево Понимания
	//understandingSituation()
	и затем создать новую цель understanding_purpose_image.go
*/


	return false
}
////////////////////////////////////////////////




///////////////////////////////////////////////////

// выдать текущую инфу для http://go/pages/mental_cicles.php
func GetCicklesToPult()(string){

//	saveFromNextIDcurretCicle=[]int{1,2,3} // тестирвоание
	out:=""
	cickl:=saveFromNextIDcurretCicle
	if cickl==nil || len(cickl)==0{
		out+="Нет текущего цикла осмысления.<hr><br><br>"
	}else{
		out+="<table cellpadding=0 cellspacing=0 border=1 class='main_table'>"
		out+="<tr><th class='table_header'>goNext ID</th>"
		out+="<th class='table_header'>ID дерева понимания</th>"
		out+="<th class='table_header'>ID дерева автоматизмов</th>"
		out+="<th class='table_header'>ID ментального автоматизма</th></tr>"
		for i := 0; i < len(cickl); i++ {
			gn:=goNextFromIDArr[cickl[i]]
			if gn!=nil{
				warn:=""
				if detectedActiveLastNodID != gn.MotorBranchID{
					warn="style='color:red' title='ID дерева автоматизмов не соотвествует goNext.MotorBranchID'"
				}
				style:="style='font-size:19px;font-weight:bold;cursor:pointer'"
				out+="<tr><td class='table_cell'><span style='color:#666666'>"+strconv.Itoa(gn.ID)+"</span>"
				out+="</td><td class='table_cell'><span "+style+" onClick='show_unde_tree("+strconv.Itoa(detectedActiveLastUnderstandingNodID)+")'>"+strconv.Itoa(detectedActiveLastUnderstandingNodID)+
					"</span></td><td class='table_cell'><span "+style+" onClick='show_atmzm_tree("+strconv.Itoa(detectedActiveLastNodID)+")'>"+strconv.Itoa(detectedActiveLastNodID)+" <span "+warn+">("+strconv.Itoa(gn.MotorBranchID)+")</span> "+
					"</span></td><td class='table_cell'><span "+style+" onClick='show_ment_atmzm("+strconv.Itoa(gn.AutomatizmID)+")'>"+strconv.Itoa(gn.AutomatizmID)+
					"</span></td></tr>"
			}
		}
		out+="</table>"
	}
	// последнгие 20 кадров кратковременной памяти
//	termMemory=[]shortTermMemory{{4800,5,1},{7188,4,2},{4800,3,3}} // тестирование
	if termMemory == nil{
		return "Еще нет кратковременной памяти."
	}
	var termMemoryFrag []shortTermMemory
	if len(termMemory)>20{
		termMemoryFrag = termMemory[:20]
	}else{
		termMemoryFrag=termMemory
	}
	out+="<br><b>Кратковременная память (последнгие 20 кадров):</b><br>"
	out+="<table cellpadding=0 cellspacing=0 border=1 class='main_table'>"
	out+="<tr><th class='table_header'>goNext ID</th>"
	out+="<th class='table_header'>ID дерева понимания</th>"
	out+="<th class='table_header'>ID дерева автоматизмов</th>"
	out+="<th class='table_header'>ID ментального автоматизма</th></tr>"
	for i := len(termMemoryFrag)-1; i >=0; i-- {
		sm:=termMemory[i]
		gn:=goNextFromIDArr[sm.GoNextID]

		style:="style='font-size:19px;font-weight:bold;cursor:pointer'"
		out+="<tr><td class='table_cell'><span style='color:#666666'>"+strconv.Itoa(gn.ID)+
			"</span></td><td class='table_cell'><span "+style+" onClick='show_unde_tree("+strconv.Itoa(sm.uTreeNodID)+")'>"+strconv.Itoa(sm.uTreeNodID)+
			"</span></td><td class='table_cell'><span "+style+" onClick='show_atmzm_tree("+strconv.Itoa(gn.MotorBranchID)+")'>"+strconv.Itoa(gn.MotorBranchID)+"</span> "+
			"</span></td><td class='table_cell'><span "+style+" onClick='show_ment_atmzm("+strconv.Itoa(gn.AutomatizmID)+")'>"+strconv.Itoa(gn.AutomatizmID)+
			"</span></td></tr>"
	}
	out+="</table>"

	return out
}
//////////////////////////////////////////////////////////////////
// для GetCicklesToPult() инфа о ветке дерева автоматизмов
func GetAtmzmTreeInfo(id int)(string){
	out:=""
	node:=AutomatizmTreeFromID[id]
	if node==nil{
		return "Нет такого узла дерева автоматизмов."
	}
	out+=getStrFromCond(0,node.BaseID)+"<br>"
	out+=getStrFromCond(1,node.EmotionID)+"<br>"
	out+=getStrFromCond(2,node.ActivityID)+"<br>"
	out+=getStrFromCond(3,node.ToneMoodID)+"<br>"
	out+=getStrFromCond(4,node.SimbolID)+"<br>"
	out+=getStrFromCond(5,node.VerbalID)+"<br>"
	return out
}
// для GetCicklesToPult() инфа о ветке дерева понимания
func GetUndstgTreeInfo(id int)(string){
	out:=""
	node:=UnderstandingNodeFromID[id]
	if node==nil{
		return "Нет такого узла дерева понимания."
	}
	out += "Состояние: <b> " + getMoodStr(node) + "</b><br>"
	out += "Эмоция: <b> " + getEmotionStr(node) + "</b><br>"
	out += "Ситуация: <b> " + getSituationStr(node) + "</b><br>"
	out += "Цель: <b> " + getPurposeStr(node) + "</b><br>"
	return out
}
// для GetCicklesToPult() инфа о ментальном автоматизме
func GetMentAtmzmInfo(id int)(string){
	out:=""
	atmz:=MentalAutomatizmsFromID[id]
	if atmz==nil{
		return "Нет такого узла дерева понимания."
	}
	out += "Действия: "+GetMentalAutomotizmActionsString(atmz.ID,false)+"<br>"
	out += "(Бес)Полезность: "+strconv.Itoa(atmz.Usefulness)+"<br>"
	out += "Число использований: "+strconv.Itoa(atmz.Count)+"<br>"

	return out
}
//////////////////////////////////////////////////////////////////






