/* получение информации о пусковых стимулах Дерева автоматизмов психики 
Здесь - против цикличности импорта. */

package reflexes

import (
	"BOT/brain/action_sensor"
	"BOT/brain/psychic"
	word_sensor "BOT/brain/words_sensor"
)

func GetTreeAutomatizmTriggersInfo(treeNodeID int) string {
	var out = ""
	treeNode := psychic.AutomatizmTreeFromID[treeNodeID]
	if treeNode != nil && treeNode.ActivityID > 0 {
		out += "Действия кнопок: <b>" + GetAcivButtInfo(treeNode.ActivityID) + "</b><br>"
	}
	if treeNode.VerbalID > 0 {
		out += "Фраза: \"<b>" + word_sensor.GetPhraseStringsFromPhraseID(treeNode.VerbalID) + "\"</b><br>"
	}
	if treeNode.ToneMoodID > 0 {
		out += "" + psychic.GetToneMoodStrFromID(treeNode.ToneMoodID) + ""
	}

	return "<div style='font-weight:200;text-align:left'>" + out + "</div>"
}

// текстовое отображение сочетаний действия с пульта образа psychic.Activity
func GetAcivButtInfo(triggerID int) string {
	var out = ""
	//act := TriggerStimulsArr[triggerID]
	act := psychic.ActivityFromIdArr[triggerID]
	if act == nil { return ""	}

	if len(act.ActID) > 0 {
		for i := 0; i < len(act.ActID); i++ {
			if i > 0 { out += ", " }
			out += action_sensor.GetActionNameFromID(act.ActID[i])
		}
	} else {
		out += "нет"
	}
	return out
}