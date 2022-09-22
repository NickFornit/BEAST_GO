/*   Определение Цели в данной ситуации - ну уровне наследственных функций
исходя из текущей информационной среды CurrentInformationEnvironment:

Эти безусловные рефлексы психики прописываются в виде функций.

Генетические цели действий Beast ID гомео-параметров, которые призвано улучшить данное действие - по его ID
прописана в terminal_actons.txt (http://go/pages/terminal_actions.php)
var TerminalActionsTargetsFromID=make(map[int][]int)
*/

package psychic

///////////////////////////////////////////////

/* образ цели бессловестного действия Формируется временно и не сохранятся в файле
Объекты PurposeGeneticObject накапливаются в оперативке и удаляются во сне
 */
type PurposeGenetic struct {
	puls int // PulsCount
	veryActual bool // true - цель очень актуальна
	targetID []int //массив ID парамктров гомеостаза как цели для улучшения в данных условиях
	actionID *ActionImage //выбранный образ действия для данной цели
// для каждого actionID сила действий сначала принимается =5, а потом корректируется мозжечковыми рефлексами
}
var PurposeGeneticObject []*PurposeGenetic
// текущая цель сохраняется до перекрытия следующим orientation_N()
var CurrentPurposeGenetic PurposeGenetic
var OldPurposeGenetic PurposeGenetic  // OldPurposeGenetic=CurrentPurposeGenetic
///////////////////////////////////////


/* Определение Цели в данной ситуации - на уровне наследственных функций

 */
func getPurposeGenetic()(*PurposeGenetic){
	var pg PurposeGenetic
	pg.puls = PulsCount
	pg.veryActual=veryActualSituation
	pg.targetID=curTargetArrID

/*Сначала посмотреть подходит ли по условиям текущий безусловный или условный рефлекс и сделать автоматизм по его действиям
	чтобы проверить его в текущих условиях т.к. getPurposeGenetic() срабатывает по ориентировочному рефлексу.
		При этом уже не будет формироваться условный рефлекс при осознанном внимании
	(т.к. заблокируется выработанным пробным действием)
 */
	//есть ли подходящий по условиям безусловный или условный рефлекс и сделать автоматизм по его действиям
	if len(actualRelextActon)>0{
		_,aImg:=CreateNewActionImageImage(actualRelextActon,nil,0,0)
		pg.actionID = aImg
	}else {
/* этого практивески не может быть, потому, что если нет рефлексов,
		то дейсвтвия в GetActualReflexAction() берутся из эффектов действий
НО если вообще нет лействий, то остается только случайное действие или бездействие
Так что лучше не заморачиваться с этим!
 */
		// veryActualSituation: плохо для  1, 2, 7 и/или 8  параметров гомеостаза
		if veryActualSituation { // нужно хоть что-то сделать, ПАНИКА
			ActID:=[]int{21} // паника
			_,aImg:=CreateNewActionImageImage(ActID,nil,0,0)
			pg.actionID = aImg
		}
	}

	PurposeGeneticObject = append(PurposeGeneticObject, &pg)
	OldPurposeGenetic = CurrentPurposeGenetic
	CurrentPurposeGenetic = pg
	return &pg
}
/////////////////////////////////////////////////////////



/////////////////////////////////

