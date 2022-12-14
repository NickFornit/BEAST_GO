/* показывать состояние психики на Пульте

*/

package psychic

import (
	word_sensor "BOT/brain/words_sensor"
	"sort"
	"strconv"
)

/////////////////////////////////////////////

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
		return "Еще нет кратковременной памяти. Beast только что проснулся и еще ни о чем не подумал."
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
		if sm.GoNextID==0{
			return "Нулевой образ звена цикла в GetCicklesToPult()."
		}
		gn:=goNextFromIDArr[sm.GoNextID]
		if gn==nil{
			out+="<tr><td colspan=10>Нет образа звена цикла с ID = "+strconv.Itoa(sm.GoNextID)+"</td></tr>"
		}
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
////////////////////////////////////

//для http://go/pages/mental_rules.php инфа о ID goNext
func GetGoNextInfo(id int)(string){
	if id==0{
		return "Нулевой образ звена цикла в GetGoNextInfo."
	}
		gn:=goNextFromIDArr[id]
		if gn==nil{
			return "Нет образа звена цикла с ID = "+strconv.Itoa(id)
		}
		style:="style='font-size:19px;font-weight:bold;cursor:pointer'"
		out:="<tr><td class='table_cell'><span style='color:#666666'>"+strconv.Itoa(id)+
			"</span></td><td class='table_cell'><span "+style+" onClick='show_atmzm_tree("+strconv.Itoa(gn.MotorBranchID)+")'>"+strconv.Itoa(gn.MotorBranchID)+"</span> "+
			"</span></td><td class='table_cell'><span "+style+" onClick='show_ment_atmzm("+strconv.Itoa(gn.AutomatizmID)+")'>"+strconv.Itoa(gn.AutomatizmID)+
			"</span></td></tr>"

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
	//if node.VerbalID>0 { нафиг масикировать лажи! первый символ не должен быть в ветке, если нет фразы!
		out += getStrFromCond(4, node.SimbolID) + "<br>"
	//}
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



// выдать текущую инфу Oбъекты значимости для http://go/pages/mental_importance.php
func GetImportanceToPult()(string){

	//	saveFromNextIDcurretCicle=[]int{1,2,3} // тестирвоание
	out:=""
	iArr:=importanceFromID
	if iArr==nil || len(iArr)==0{
		out+="Еще нет объектов значимости.<hr><br><br>"
	}else{
		out+="<table cellpadding=0 cellspacing=0 border=1 class='main_table'>"
		out+="<tr><th class='table_header'>объект ID</th>"
		out+="<th class='table_header'>ID<br>дерева<br>автоматизмов</th>"
		out+="<th class='table_header'>ID<br>дерева<br>понимания</th></th>"
		out+="<th class='table_header' width=350>Тип объекта</th>"
		out+="<th class='table_header'>Значимость</th></tr>"
		keys := make([]int, 0, len(iArr))
		for k := range iArr {
			keys = append(keys, k)
		}
		sort.Ints(keys)
		for _, k := range keys {
			oi:=importanceFromID[k]

				style:="style='font-size:19px;font-weight:bold;cursor:pointer'"
				out+="<tr><td class='table_cell'><span style='color:#666666'>"+strconv.Itoa(oi.ID)+"</span>"+

"</td><td class='table_cell' "+style+" onClick='show_atmzm_tree("+strconv.Itoa(oi.NodeAID)+")'>"+strconv.Itoa(oi.NodeAID)+
"<td class='table_cell' "+style+" onClick='show_unde_tree("+strconv.Itoa(oi.NodePID)+")'>"+strconv.Itoa(oi.NodePID)

out+="</td><td class='table_cell' onClick='show_object("+strconv.Itoa(oi.Type)+","+strconv.Itoa(oi.ObjectID)+")'>id="+strconv.Itoa(oi.ObjectID)+" (type="+strconv.Itoa(oi.Type)+"): <b>"+importanceTypeName[oi.Type]+"</b>"+
"</td><td class='table_cell'>"+strconv.Itoa(oi.Value)+
"</td></tr>"
			}
		}
		out+="</table>"

	return out
}
// информация об объекте значимости
func GetImportanceObjectInfo(objID int,objType int)(string){
out:=""

switch objType {
case 1:     // ID ActionsImage
	out=GetActionsString(objID)
case 2: // ID MentalActionsImages
	out=GetMentalAutomotizmActionsString(objID, false) // GetMentalActionInfo(objID) ?
case 3: // ID несловестного действия ActionsImage.ActID[n]
	//out=GetAutomatizmSequenceInfo(objID,true)
case 4: // ID Verbal - при активации дерева автоматизмов
	out=GetStringsFromVerbalID(objID)
case 5: // ID отдельной фразы Verbal.PhraseID[n]
	out=word_sensor.GetPhraseStringsFromPhraseID(objID)
case 6: // ID отдельного слова  из Verbal.PhraseID[n]
	out=word_sensor.GetWordFromWordID(objID)
case 7: // ID тон сообщения с Пульта  Verbal.ToneID
	out=getToneStrFromID(objID)
case 8: // ID настроение оператора  Verbal.MoodID
	out=getMoodStrFromID(objID-19)
}

	return out
}
/////////////////////////////////////////////////
