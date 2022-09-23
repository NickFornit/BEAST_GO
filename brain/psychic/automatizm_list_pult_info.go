/* список автоматизмов для Пульта
для http://go/pages/automatizm_table.php
*/

package psychic

import (
	"BOT/brain/gomeostas"
	"sort"
	"strconv"
)

//////////////////////////////////////////


////////////////////////////////////////////
func GetAutomatizmInfo(limitBasicID int)(string){
	var out=""
	// сколько рефлексов есть
	uAutomatizmCount:=len(AutomatizmFromIdArr)
	// если больше 1000 то выдавать только по одному из 3-х базовыз состояний, иначе сильно тормозит
	if uAutomatizmCount > 1000{
		if limitBasicID==0{
			limitBasicID=1// начинать с Плохо
		}
	}
	// переключатель диапазона вывода
	if limitBasicID>0{
		out+="<br>Показывать: "
		out+="<span style='cursor:pointer;"
		if limitBasicID==1{out+="background-color:#FFFF9D;font-weight:bold;"}
		out+="' onClick='show_level(1)'>Плохо</span> "

		out+="<span style='cursor:pointer;"
		if limitBasicID==2{out+="background-color:#FFFF9D;font-weight:bold;"}
		out+="' onClick='show_level(2)'>Норма</span> "

		out+="<span style='cursor:pointer;"
		if limitBasicID==3{out+="background-color:#FFFF9D;font-weight:bold;"}
		out+="' onClick='show_level(3)'>Хорошо</span> "
	}

	header:="<tr><th width=70 class='table_header'>ID</th>" +
		"<th width=70 class='table_header'>ID <br><nobr>узла Дерева</nobr></th>" +
		"<th width=70 class='table_header'  title='ID базового состояния'>BaseID</th>" +
		"<th width='25%' class='table_header'  title='ID образа сочетания эмоций'>Эмоции</th>" +
		"<th width='30' class='table_header'  title='ID образа пускового стимула'>Пусковой стимул</th>" +
		"<th width='25%' class='table_header'  title='Действия автоматизма, их может быть много видов.'>Строка действий</th>" +
		"<th width='30' class='table_header' title='ID следующего автоматизма в цепочке'>NextID</th>" +
		"<th width='30' class='table_header' title='Энергичность'>Energy</th>" +
		"<th width='30' class='table_header' title='(БЕС)ПОЛЕЗНОСТЬ: -10 вред 0 +10 +n польза'>Полезность</th>" +
		"<th width='30' class='table_header' title='Уверенность'>Belief</th></tr>"


	out+="<table class='main_table'  cellpadding=0 cellspacing=0 border=1 width='1000px' style='font-size:14px;'>" +
		header

	if len(AutomatizmFromIdArr)==0{
		return out+"</table><br>Подождите пока не активируется психика (не более 4 секунд) и нажмите <a href='/pages/automatizm_table.php'>Обновить</a>"
	}



	keys := make([]int, 0, uAutomatizmCount)
	for k := range AutomatizmFromIdArr {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, k := range keys {
		v:=AutomatizmFromIdArr[k]

id := strconv.Itoa(k)
nodeID:=""
baseID:=""
emotionID:=""
emotionTitle:=""
activityID:=""
//activityTitle:=""

// узел дерева, к которому прикреплен автоматизм
nodeA:=AutomatizmTreeFromID[v.BranchID]
if nodeA!=nil {
	nodeID=strconv.Itoa(nodeA.ID)
	if limitBasicID > 0 && nodeA.BaseID!=limitBasicID{
		continue
	}

	baseID = strconv.Itoa(nodeA.BaseID)+" "+gomeostas.GetBaseCondFromID(nodeA.BaseID)

	/// эмоции
	emotionID=strconv.Itoa(nodeA.EmotionID)
	emo:=EmotionFromIdArr[nodeA.EmotionID]
	for i := 0; i < len(emo.BaseIDarr); i++ {
		if i > 0 {
			emotionTitle += ", "
		}
		emotionTitle += gomeostas.GetBaseContextCondFromID(emo.BaseIDarr[i])
	}
	/////////////////////////////

	// пусковой стимуд
	activityID=strconv.Itoa(nodeA.ActivityID)

/*  вызывает ЦИКЛИЧЕСКИЙ ИМПОРТ

	act := reflexes.TriggerStimulsArr[nodeA.ActivityID]
	if len(act.RSarr) > 0 {
		for i := 0; i < len(act.RSarr); i++ {
			if i > 0 {
				activityTitle += ", "
			}
			activityTitle += actionSensor.GetActionNameFromID(act.RSarr[i])
		}
	}
	if len(act.PhraseID) > 0 {
		if len(activityTitle) > 0 {
			activityTitle += "<br>"
		}
		for i := 0; i < len(act.PhraseID); i++ {
			if i > 0 {
				activityTitle += "; "
			}
			w := wordsSensor.GetPhraseStringsFromPhraseID(act.PhraseID[i])
			//w=strings.Trim(w,"")
			activityTitle += "\"" + w + "\""
		}
	}*/
}//if nodeA!=nil{


//////////////
out += "<tr >"
out += "<td class='table_cell' >"+id+"</td>";
out += "<td class='table_cell' >"+nodeID+"</td>";
out += "<td class='table_cell' ><nobr>"+baseID+"</nobr></td>";
out += "<td class='table_cell' title='"+emotionTitle+"'><nobr>"+emotionID+"</nobr></td>";
out += "<td class='table_cell' title='Информация по клику' onClick='show_trigger("+activityID+")' style='cursor:pointer;color:blue'>"+activityID+"</td>";
out += "<td class='table_cell' title='Информация по клику'  onClick='show_actions("+id+",`"+v.Sequence+"`)' style='cursor:pointer;color:blue'><nobr>"+v.Sequence+"</nobr></td>";
out += "<td class='table_cell' >"+strconv.Itoa(v.NextID)+"</td>";
out += "<td class='table_cell' >"+strconv.Itoa(v.Energy)+"</td>";
out += "<td class='table_cell' >"+strconv.Itoa(v.Usefulness)+"</td>";
out += "<td class='table_cell' >"+strconv.Itoa(v.Belief)+"</td>";
out += "</tr>"
	}
	out+="</table>"
	return out
}

