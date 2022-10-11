/*
самый простейший безусловный рефлекс
по сочетаниям редактора http://go/pages/terminal_actions.php
Данный редактор связывает действие с тем, какие гомео-параметры улучшает данное действие.
Это внутренний рефлекс самоадаптации, выполняемый, если нет пусковых стимулов для более высокоуровневых рефлексов
Или рефлекс по умолчанию, что означает - реагирование будет всегда хотя бы на таком уровне
*/

package reflexes

import (
	TerminalActions "BOT/brain/terminete_action"
	"BOT/lib"
)

// найти и выполнить простейший безусловный рефлекс
func findAndExecuteSimpeReflex() {
	// только если нет другого рефлекса
	if len(lib.ActionsForPultStr) > 0 { return }
	_, actID,_ := TerminalActions.ChooseSimpleReflexexAction()
	if actID > 0 { // совершить это действие
		// очистить буфер передачи действий на пульт
		// lib.ActionsForPultStr = ""
		actStr := "0|" + TerminalActions.TerminalActonsNameFromID[actID]
		lib.SentActionsForPult(actStr)
	}
}