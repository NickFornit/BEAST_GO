/* Передача инфы на Пульт */

package lib

import (
	"fmt"
)

// строка вывода на пульт - для func WritePultConsol
var WritePultConsolStr = ""
/* вывести на консоль Пульта
Сообщения накапливаются в WritePultConsolStr и откправлются с каждым пульсом
*/
func WritePultConsol(print string) {
	fmt.Println("НА ПУЛЬТ: ", print)
	WritePultConsolStr =	print + "<br>" + WritePultConsolStr
}

// функция вызова паники с информированием в логе Пульта
func TodoPanic(panicWarning string){
	WritePultConsol("<span style='color:red;font-size:19px;font-weight:bold;'>ПАНИКА: </span> "+panicWarning)
	panic(panicWarning)
}
//////////////////////////////////////////////////

var ActionsForPultStr = ""
/* вывести на Пульт действия Бота строкой lib.SentActionsForPult("xcvxvxcv")
Каждая акция - в формате: вид действия (1 - действие рефлекса, 2 - фраза) затем строка акции,
например: "1|Предлогает поиграть" или "2|Привет!"
Можно передавать неограниченную последовательность акций, разделяя их "||"
например: "1|Предлогает поиграть||2|Привет!"
*/
func SentActionsForPult(print string) {
	if len(ActionsForPultStr) > 0 { // еще не прочитана предыдущая инфа т.к. читается раз в пульс, а после действия может измениться услове и будет новое действие
		ActionsForPultStr += "||" + print
		return
	}
	ActionsForPultStr =	print
}