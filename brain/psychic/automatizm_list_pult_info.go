/* список автоматизмов для Пульта
для http://go/pages/automatizm_table.php
*/

package psychic

import (
	"BOT/brain/gomeostas"
	termineteAction "BOT/brain/terminete_action"
	word_sensor "BOT/brain/words_sensor"
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
	if uAutomatizmCount > 100{
		if limitBasicID==0{
		//  теперь 0 - непривязанные	limitBasicID=1// начинать с Плохо
		}
	// переключатель диапазона вывода
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

		out+="<span style='cursor:pointer;"
		if limitBasicID==0{out+="background-color:#FFFF9D;font-weight:bold;"}
		out+="' onClick='show_level(0)' title='Автоматизмы, не привязанные к определнным условиям состояния Beast.'>Свободные</span> "
	}

	header:="<tr><th width=70 class='table_header'>ID ав-ма</th>" +
		"<th width=70 class='table_header' style='background-color:#CCC5FF;' title='Если <1000000, то - узел ветки дерева\nЕсли <2000000 - ID образа действий\nостальное - ID фразы'>ID <br>объхекта<br>Привязки</th>" +
		"<th width=70 class='table_header' style='background-color:#CCC5FF;'  title='ID базового состояния'>BaseID Дерева</th>" +
		"<th width='70' class='table_header' style='background-color:#CCC5FF;'  title='ID образа сочетания эмоций'>Эмоции Дерева</th>" +
		"<th width='10%' class='table_header' style='background-color:#CCC5FF;'  title='ID образа пускового стимула'>ДействияID-ФразаID-НастроениеID</th>" +
		"<th class='table_header'  title='Действия автоматизма, их может быть много видов.'>Строка действий</th>" +
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
trStr:=""
onclickID:=""
//activityTitle:=""
sStyle:="style='cursor:pointer;color:blue'"
emClick:=""
actClick:=""

// узел дерева, к которому прикреплен автоматизм
if v.BranchID<1000000{
nodeA:=AutomatizmTreeFromID[v.BranchID]
if nodeA!=nil {
	nodeID = strconv.Itoa(nodeA.ID)
	if limitBasicID > 0 {
		if nodeA.BaseID != limitBasicID || limitBasicID == 0 {
			continue
		}
	}
	baseID = strconv.Itoa(nodeA.BaseID) + " " + gomeostas.GetBaseCondFromID(nodeA.BaseID)

	/// эмоции
	emotionID = strconv.Itoa(nodeA.EmotionID)
	emo := EmotionFromIdArr[nodeA.EmotionID]
	if emo == nil { // так не должно быть! эмоция д.б. всегда
		//lib.WritePultConsol("В emo:=EmotionFromIdArr[nodeA.EmotionID] НЕТ ЭМОЦИИ!!!! Какой-то образ неверен, нужен программист!")
		// лучше бы выкинуть панику: panic("НЕТ ЭМОЦИИ!!!!")
		emotionTitle += "<span style='color:red'>НЕТ ЭМОЦИИ! нужно разобраться с EmotionFromIdArr[]</span>"
		continue
	}
	for i := 0; i < len(emo.BaseIDarr); i++ {
		if i > 0 {
			emotionTitle += ", "
		}
		emotionTitle += gomeostas.GetBaseContextCondFromID(emo.BaseIDarr[i])
	}
	/////////////////////////////

	// пусковой стимул
	trStr=strconv.Itoa(nodeA.ActivityID)+"-"+strconv.Itoa(nodeA.VerbalID)+"-"+strconv.Itoa(nodeA.ToneMoodID)

	emClick="onClick='show_emotion("+emotionID+")'"
	actClick="onClick='show_trigger("+nodeID+")'"

}//if nodeA!=nil{

//if v.BranchID<1000000{
}else{
	if limitBasicID!=0{
		continue
	}
	nodeID = strconv.Itoa(v.BranchID)
	onclickID="onClick='show_object("+strconv.Itoa(v.BranchID)+")' style='cursor:pointer;color:blue'"
	baseID="не"
	emotionID="привязан"
//	emotionTitle="привязан"
	trStr="к дереву"
	sStyle=""

}


//////////////
out += "<tr >"
out += "<td class='table_cell' >"+id+"</td>";
out += "<td class='table_cell' "+onclickID+">"+nodeID+"</td>";
out += "<td class='table_cell' ><nobr>"+baseID+"</nobr></td>";
out += "<td class='table_cell' title='"+emotionTitle+"' "+emClick+" "+sStyle+"><nobr>"+emotionID+"</nobr></td>";
out += "<td class='table_cell' title='Информация по клику' "+actClick+" "+sStyle+">"+trStr+"</td>";
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
////////////////////////////////////


// показать пусковой стимул, к которому привязан автоматизм
func GetStrnameFromobjectID(objectID int)(string){
	if objectID>1000000 && objectID<2000000{// это - действия кнопок с пульта
		imgID:=objectID-1000000
		act:=ActivityFromIdArr[imgID]
		if act == nil{
			return ""
		}
		out:=""
		for i := 0; i < len(act.ActID); i++ {
			if i>0{out+=", "}
			out+=termineteAction.TerminalActonsNameFromID[act.ActID[i]]
		}
		return "Пусковые действия:<br><br><b>"+out+"</b>"
	}
	if objectID>2000000{// это - фраза
		imgID:=objectID-2000000
		return "Пусковая фраза:<br><br><b>"+word_sensor.GetPhraseStringsFromPhraseID(imgID)+"</b>"
	}
	return ""
}
/////////////////////////////////////

