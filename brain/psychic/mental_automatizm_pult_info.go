/* список ментальных автоматизмов для Пульта
для http://go/pages/mental_automatizm_table.php
*/

package psychic

import (
	"sort"
	"strconv"
)

//////////////////////////////////////////


////////////////////////////////////////////
func GetMentalAutomatizmInfo(limitPage int)(string){
	var out=""
	var isLimited=false
	// сколько рефлексов есть
	uAutomatizmCount:=len(MentalAutomatizmsFromID)
	maxP:=1000
	// если больше maxP то выдавать только по одному из 3-х базовыз состояний, иначе сильно тормозит
	if uAutomatizmCount > maxP{
		isLimited=true
	// переключатель диапазона вывода
		var nPages=int(uAutomatizmCount/maxP)
		if uAutomatizmCount%maxP >0 {
			nPages++
		}
		out+="<br>Страницы: "
		for i := 0; i < nPages; i++ {
			out += "<span style='cursor:pointer;"
			if limitPage == i {
				out += "background-color:#FFFF9D;font-weight:bold;font-size:19px;"
			}
			out += "' onClick='show_page("+strconv.Itoa(i)+")'>"+strconv.Itoa(i+1)+"</span> "
		}

		}
		////////////////////////////////////////////////////

	header:="<tr><th width=70 class='table_header'>ID автоматизма</th>" +
		"<th width='10%' class='table_header' style='background-color:#CCC5FF;'  title='ID образа действия'>Образ действия</th>" +
		"<th width='30' class='table_header' title='(БЕС)ПОЛЕЗНОСТЬ: -10 вред 0 +10 +n польза'>Полезность</th>" +
		"<th width='30' class='table_header' title='надежность: число использований с подтверждением (бес)полезности Usefulness'>Число использований</th></tr>"


	out+="<table class='main_table'  cellpadding=0 cellspacing=0 border=1 width='1000px' style='font-size:14px;'>" +
		header

	if len(MentalAutomatizmsFromID)==0{
		return out+"</table><br>Подождите пока не активируется психика (не более 6 секунд) и нажмите <a href='/pages/mental_automatizm_table'>Обновить</a>"
	}



	keys := make([]int, 0, uAutomatizmCount)
	for k := range MentalAutomatizmsFromID {
		keys = append(keys, k)
	}
	sort.Ints(keys) // сортировка по ID автоматизма

var count=0
	//var showed=0
	for _, k := range keys {
		count++
		if isLimited {
			if count <= limitPage*maxP || count > (limitPage+1)*maxP{
				continue
			}
		}

		v:=MentalAutomatizmsFromID[k]
id := strconv.Itoa(k)

out += "<tr >"
out += "<td class='table_cell' >"+id+"</td>";
out += "<td class='table_cell' title='Информация по клику'  onClick='show_actions("+id+")' style='cursor:pointer;color:blue'><nobr>"+strconv.Itoa(v.ActionsImageID)+"</nobr></td>";
out += "<td class='table_cell' >"+strconv.Itoa(v.Usefulness)+"</td>";
out += "<td class='table_cell' >"+strconv.Itoa(v.Count)+"</td>";
out += "</tr>"

	}
	out+="</table>"
	return out
}
////////////////////////////////////




