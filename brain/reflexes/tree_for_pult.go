/* пока дерева рефлексов на Пульте */

package reflexes

import (
	actionSensor "BOT/brain/action_sensor"
	"BOT/brain/gomeostas"
	termineteAction "BOT/brain/terminete_action"
	wordsSensor "BOT/brain/words_sensor"
	"sort"
	"strconv"
)

// образ дерева рефлексов для вывода
var reflesTreeModel = ""

func GetReflexesTreeForPult() string {
	// против паники типа "одновременная запись и считывание карты"
	if notAllowScanInReflexesThisTime { return "!!!" }
	reflesTreeModel = ""
	scanReflexesNodes(-1, &ReflexTree)

	return reflesTreeModel
}

func scanReflexesNodes(level int, node *ReflexNode) {
	if node.ID == 69 {
		node.ID = 69
	}
	if node.ID > 0 {
		reflesTreeModel += setShift(level)
		switch level {
		case 0:
			reflesTreeModel += getStrFromCond(1, node.baseID) + "(nodeID=" + strconv.Itoa(node.ID) + ")"
		case 1:
			reflesTreeModel += getStrFromCond(2, node.StyleID) + "(nodeID=" + strconv.Itoa(node.ID) + ")"
		case 2:
			if node.ActionID == 0 { // древний рефлекс - без пусковыого стимула
				reflesTreeModel += getStrFromCond(2, node.StyleID) + "(nodeID=" + strconv.Itoa(node.ID) + ")"
			} else {
				reflesTreeModel += getStrFromCond(3, node.ActionID) + "(nodeID=" + strconv.Itoa(node.ID) + ")"
			}
		case 3:
			reflesTreeModel += getStrFromCond(3, node.ActionID) + "(nodeID=" + strconv.Itoa(node.ID) + ")"
		}

		// если есть рефлекс - показать действия
		if node.ConditionedReflex > 0 {
			//rID:=ConditioneReflexes[node.ConditionedReflex]
		} else if node.GeneticReflexID > 0 {
			reflex := GeneticReflexes[node.GeneticReflexID]
			if reflex.ID == 45 {
				reflex.ID = 45
			}
			reflesTreeModel += " <span style='color:blue'>Рефлекс(" + strconv.Itoa(reflex.ID) + "): "+getReflexActionsStrFromID(reflex) + "</span>"
		}
		reflesTreeModel += "<br>\n"
	}
	level++
	for n := 0; n < len(node.Children); n++ {
		scanReflexesNodes(level,&node.Children[n])
	}
}

// отступ
func setShift(level int) string {
	var sh = ""
	for n := 0; n < level; n++ {
		sh += "&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;"
	}
	return sh
}

// из ID образа получить составляющие в виде строк
func getStrFromCond(level int,imgID int) string {
	var out = ""

	switch level {
	case 1: if imgID > 0 && imgID < 4 { out=gomeostas.GetBaseCondFromID(imgID)
	} else {
		out += "<span style='color:red'>несуществующее Базовое состояние ID=" + strconv.Itoa(imgID) + "</span>"
	}
	case 2:	out = getStrnameFromBaseImageID(imgID)
	case 3: out = getStrnameFromStyleImageID(imgID)
	}

	return out
}

// названия базовых контекстов в их сочетании -из ID их образа
func getStrnameFromBaseImageID(id int) string {
	var out = ""

	img := BaseStyleArr[id]
	for i := 0; i < len(img.BSarr); i++ {
		if i > 0 { out += ", " }
		name := gomeostas.GetBaseContextCondFromID(img.BSarr[i]) + ""
		if len(name) == 0 {
			out += "<span style='color:red'>несуществующий Базовый контекст ID=" + strconv.Itoa(img.BSarr[i]) + "</span>"
		} else {
			out += name
		}
	}

	return out
}

// названия Пусковых стимулов в их сочетании -из ID их образа
func getStrnameFromStyleImageID(id int) string {
	var out = ""

	if id == 0 { return "" }
	img := TriggerStimulsArr[id]
	if img == nil { return "" }
	// Действия с Пульта?
	if len(img.RSarr) > 0 {
		for i := 0; i < len(img.RSarr); i++ {
			if i > 0 { out += ", " }
			name := actionSensor.GetActionNameFromID(img.RSarr[i])
			if len(name) == 0 {
				out += "<span style='color:red'>несуществующий стимул ID=" + strconv.Itoa(img.RSarr[i]) + "</span>"
			} else {
				out += name
			}
		}
	}
	// фразы с Пульта?
	if len(img.PhraseID) > 0 {
		out += "Фразы c ID: "
		for i := 0; i < len(img.PhraseID); i++ {
			if i > 0 { out += ", " }
			out += strconv.Itoa(img.PhraseID[i])
		}
	}

	return out
}

func getReflexActionsStrFromID(reflex *GeneticReflex) string {
	if reflex == nil { return "" }
	var out = ""
	for i := 0; i < len(reflex.ActionIDarr); i++ {
		if i > 0 { out += ", " }
		out += termineteAction.TerminalActonsNameFromID[reflex.ActionIDarr[i]]
	}

	return out
}

// выдать таблицу условных рефлексов для http://go/pages/condition_reflexes.php
func GetConditionReflexInfo(limitBasicID int) string {
	var out = ""

	// сколько рефлексов есть
	ureflexCount := len(ConditionReflexes)
	// если больше 1000 то выдавать только по одному из 3-х базовыз состояний, иначе сильно тормозит
	if ureflexCount > 1000 {
		if limitBasicID == 0 {
			limitBasicID = 1 // начинать с Плохо
		}
	}
	// переключатель диапазона вывода
	if limitBasicID > 0 {
		out += "<br>Показывать: "
		out += "<span style='cursor:pointer;"
		if limitBasicID == 1 { out += "background-color:#FFFF9D;font-weight:bold;" }
		out += "' onClick='show_level(1)'>Плохо</span> "

		out += "<span style='cursor:pointer;"
		if limitBasicID == 2 { out += "background-color:#FFFF9D;font-weight:bold;" }
		out += "' onClick='show_level(2)'>Норма</span> "

		out += "<span style='cursor:pointer;"
		if limitBasicID == 3 { out += "background-color:#FFFF9D;font-weight:bold;" }
		out += "' onClick='show_level(3)'>Хорошо</span> "
	}

	out += "<table class='main_table'  cellpadding=0 cellspacing=0 border=1 width='100%' style='font-size:14px;'>" +
		"<tr><th width=70 class='table_header'>ID<br>рефлекса</th>" +
		"<th width=70 class='table_header'>ID (1 уровень)<br><nobr>базового состояния</nobr></th>" +
		"<th width='25%' class='table_header'>ID (2) актуальных контекстов<br>через запятую</th>" +
		"<th width='25%' class='table_header'>ID (3) пусковых стимулов<br>через запятую</th>" +
		"<th width='25%' class='table_header'>ID действий<br>через запятую</th>" +
		"<th width='25%' class='table_header'>Родитель</th>" +
		"<th width='25%' class='table_header'>Возраст</th>" +
		"<th width='30' class='table_header' title='Удалить рефлекс'>Х</th></tr>"

	if len(ConditionReflexes) == 0 {
		return "<table id='main_table' class='main_table'  cellpadding=0 cellspacing=0 border=1 width='100%' style='font-size:14px;'>" +
		"<tr><th width=70 class='table_header'>ID<br>рефлекса</th>" +
		"<th width=70 class='table_header'>ID (1 уровень)<br><nobr>базового состояния</nobr></th>" +
		"<th width='25%' class='table_header'>ID (2) актуальных контекстов<br>через запятую</th>" +
		"<th width='25%' class='table_header'>ID (3) пусковых стимулов<br>через запятую</th>" +
		"<th width='25%' class='table_header'>ID действий<br>через запятую</th>" +
		"<th width='25%' class='table_header'>Родитель</th>" +
		"<th width='25%' class='table_header'>Возраст</th>" +
		"<th width='30' class='table_header' title='Удалить рефлекс'>Х</th></tr>" +
		"<tr class='highlighting' ><td colspan=7 >Еще нет условных рефлексов.</td></tr>"+
		"</table>"
	}

	keys := make([]int, 0, ureflexCount)
	for k := range ConditionReflexes {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, k := range keys {
		v := ConditionReflexes[k]
		if limitBasicID>0 && v.lev1!=limitBasicID {
			continue
		}
		id := strconv.Itoa(k)
		lev1 := gomeostas.GetBaseCondFromID(v.lev1)
		var lev2 = ""
		for i := 0; i < len(v.lev2); i++ {
			if i > 0 { lev2 += ", " }
			lev2 += gomeostas.GetBaseContextCondFromID(v.lev2[i])
		}
		act := TriggerStimulsArr[v.lev3]
		var lev3 = ""
		if len(act.RSarr) > 0 {
			for i := 0; i < len(act.RSarr); i++ {
				if i > 0 { lev3 += ", " }
				lev3 += actionSensor.GetActionNameFromID(act.RSarr[i])
			}
		}
		if len(act.PhraseID) > 0 {
			if len(lev3) > 0 {
				lev3 += "<br>"
			}
			for i := 0; i < len(act.PhraseID); i++ {
				if i > 0 { lev3 += "; " }
				w := wordsSensor.GetPhraseStringsFromPhraseID(act.PhraseID[i])
				//w=strings.Trim(w,"")
				lev3 += "\"" + w + "\""
			}
		}
		var tact=""
		for i := 0; i < len(v.ActionIDarr); i++ {
			if i > 0 { tact += ", " }
			tact += termineteAction.TerminalActonsNameFromID[v.ActionIDarr[i]]
		}
		var rank = ""
		if v.rank == 0 {
			rank = "Безусловный"
		} else {
			rank = "Условный"
		}

		t := v.birthTime
		y := int(t / (3600 * 24 * 365))
		d := int((t - y * 3600 * 24 * 365) / (3600 * 24))
		s := t - (y * 3600 * 24 * 365) - (d * 3600 * 24)
		old := strconv.Itoa(y) + "лет " + strconv.Itoa(d) + "дней " + strconv.Itoa(s) + "сек"

		out += "<tr >"
		out += "<td class='table_cell' >"+id+"</td>"
		out += "<td class='table_cell' >"+lev1+"</td>"
		out += "<td class='table_cell' >"+lev2+"</td>"
		out += "<td class='table_cell' >"+lev3+"</td>"
		out += "<td class='table_cell' >"+tact+"</td>"
		out += "<td class='table_cell' >"+rank+"</td>"
		out += "<td class='table_cell' >"+old+"</td>"

		out += "<td class='table_cell' align='center' title='Удалить рефлекс' style='cursor:pointer;' ><input type='checkbox' class='deleteCHBX' id='dchbx_" + id + "'></td>"
		out += "</tr>"
	}
	out += "</table>"
	return out
}