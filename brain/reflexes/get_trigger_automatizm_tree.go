/* получение информации о пусковых стимулах Дерева автоматизмов психики 
Здесь - против цикличности импорта.

*/

package reflexes

import (
	"BOT/brain/action_sensor"
	"BOT/brain/psychic"
	wordSensor "BOT/brain/words_sensor"
)

///////////////////////////////////////////////////////////
func GetTreeAutomatizmTriggersInfo(treeNodeID int)(string){
var out=""
	treeNode:=psychic.AutomatizmTreeFromID[treeNodeID]
	if treeNode.ActivityID >0{
		out+="Действия кнопок: <b>"+GetAcivButtInfo(treeNode.ActivityID)+"</b><br>"
	}
	if treeNode.VerbalID >0{
		out+="Фраза: \"<b>"+wordSensor.GetWordFromPraseNodeID(treeNode.VerbalID)+"\"</b><br>"
	}
	if treeNode.ToneMoodID >0{
		out+=""+psychic.GetToneMoodStrFromID(treeNode.ToneMoodID)+""
	}

	return "<div style='font-weight:200;text-align:left'>"+out+"</div>"
}
//////////////////////////////////////////////////////////


func GetAcivButtInfo(triggerID int)(string){
	var out=""
	act := TriggerStimulsArr[triggerID]
	if act == nil {
		return ""
	}

	if len(act.RSarr) > 0 {
		for i := 0; i < len(act.RSarr); i++ {
			if i > 0 {
				out += ", "
			}
			out += action_sensor.GetActionNameFromID(act.RSarr[i])
		}
	}else{
		out +="нет"
	}
	/*
	if len(act.PhraseID) > 0 {
		if len(out) > 0 {
			out += "<br>"
		}
		for i := 0; i < len(act.PhraseID); i++ {
			if i > 0 {
				out += "; "
			}
			w := wordSensor.GetPhraseStringsFromPhraseID(act.PhraseID[i])
			//w=strings.Trim(w,"")
			out += "\"" + w + "\""
		}
	}
	 */
	return out
}
////////////////////////////////////////////////////