/*
самый простейший безусловный рефлекс
	по сочетаниям редактора http://go/pages/terminal_actions.php
	Данный редактор связывает действие с тем, какие гомео-параметры улучшает данное действие.
*/


package reflexes

import (
	TerminalActions "BOT/brain/terminete_action"
	"BOT/lib"
)

/////////////////////////////////////

// найти и выполнить простейший безусловный рефлекс

func findAndExecuteSimpeReflex()(bool){
	_,actID,_:=TerminalActions.ChooseSimpleReflexexAction()
	if actID>0{// совершить это действие
		// очистить буфер передачи действий на пульт
		//lib.ActionsForPultStr = ""
		actStr:="0|"+TerminalActions.TerminalActonsNameFromID[actID]
		lib.SentActionsForPult(actStr)
	}

	return false
}
//////////////////////////////////////




