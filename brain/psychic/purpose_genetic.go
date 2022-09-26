/*   Определение Цели в данной ситуации - ну уровне наследственных функций
исходя из текущей информационной среды CurrentInformationEnvironment:

Эти безусловные рефлексы психики прописываются в виде функций.

Генетические цели действий Beast ID гомео-параметров, которые призвано улучшить данное действие - по его ID
прописана в terminal_actons.txt (http://go/pages/terminal_actions.php)
var TerminalActionsTargetsFromID=make(map[int][]int)
*/

package psychic

import (
	"BOT/brain/gomeostas"
	TerminalActions "BOT/brain/terminete_action"
	word_sensor "BOT/brain/words_sensor"
	"BOT/lib"
	"strconv"
)

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



//////////////////////////////////////////////////
/* Определение Цели в данной ситуации - ну уровне наследственных функций
Здесь выбирается действие пробного автоматизма из выполнившегося рефлекса actualRelextActon.
*/
func getPurposeGeneticAndRunAutomatizm()(*Automatizm) {

	purpose:=getPurposeGenetic()// выбираются наиболее подходящие действия для автоматизмаы
	// мозжечковые рефлексы - самый первый уровень осознания - подгонка действий под заданную Цель.

	// нужно ли вообще шевелиться?
	// veryActualSituation: плохо для  1, 2, 7 и/или 8  параметров гомеостаза
	if purpose.veryActual { // нужно создавать автоматизм и тут же запускать его
		if purpose.actionID.ID > 0 {
			/* сформировать пробный автоматизм моторного действия и сразу запустить его в действие
			      Зафиксироваь время действия
			      20 пульсов следить за измнением состояния и ответными действиями - считать следствием действия
			      оценить результат и скорректировать силу мозжечком в записи автоматизма.
			   Выбрать любое действие, т.к. послед создания автоматизма в данной ветке detectedActiveLastNodID
			   он уже не вызовет orientation_1(), а будет orientation_2()
			*/
			atmzm :=createAndRunAutomatizmFromPurpose(purpose)
			// в automatizm_result.go после оценки результата будет осмысление с активацией Дерева Понимания
			return atmzm
		}
		///////////////////////////////////////
		// игнорировать новое внимание на время ожидания результата автоматизма
		if isPeriodResultWaiting {
			return AutomatizmRunning
		}

		/* нет действий (практически невозможная ситуация потому, что если нет рефлексов,
			то дейсвтвия в GetActualReflexAction() берутся из эффектов действий)
		 */
		if purpose.actionID == nil{ // AutomatizmSuccessFromIdArr=make(map[int]*Automatizm)
			// ранее найденные удачные автоматизмы
			// выбрать из ранее удачного автоматизма, перекрыть цель новой и запустить новый автоматизм
			atmzm :=chooseAutomatizmSuccessAndRun(purpose)
			return atmzm
		}
		if purpose.actionID == nil {
			// нет действий, попробовать бессмысленно выдать фразы имеющиеся Вериниковские раз нужно что-то срочно делать
			if purpose.veryActual {
				// подобрать хоть как-то ассоциирующуюся фразу из имеющизся
				phraseID := findSuitablePhrase()
				if len(phraseID) > 0 {
					purpose.actionID.PhraseID = phraseID
					atmzm := createAutomatizm(purpose)
					// запустить автоматизм
					if RumAutomatizm(atmzm) {
						// отслеживать последствия в automatizm_result.go
						setAutomatizmRunning(atmzm, purpose)
					}
					// в automatizm_result.go после оценки результата будет осмысление с активацией Дерева Понимания
					return atmzm
				}

				// return autmzm
			}
		}
		//  ЗДЕСЬ активировать Дерево Понимания НЕ НУЖНО, действие уже запущено, омысление будет по результату.
	}else { // нет атаса, можно спокойно поэкспериментивроать, если есть любопытсво

		if gomeostas.BaseContextActive[2] || gomeostas.BaseContextActive[3] { // активен Поиск или Игра
			// тупо метод тыка
				// Тупо поэкспериментировать для пополнения опыта (не)удачных автоматизмов
			// TODO !не проверено!
// в отличии от createAndRunAutomatizmFromPurpose(purpose) не использовать текущие рефлексы, а пробовать всякое
// Выдавая это на стадии 3, тварь получает реакцию оператора, которую отзеркаливает
				atmzm :=findAnySympleRandActions()
				return atmzm

			///////////////////////////////////////////////////
		}else{// НЕ ИГРА  И НЕ ПОИСК, пониженная мотивация что-то делать если нет актуальности
			if EvolushnStage == 2 {
				// нет действий, попробовать использовать AutomatizmSuccessFromIdArr.GomeoIdSuccesArr
				// выбрать из ранее удачного автоматизма, перекрыть цель новой и запустить новый автоматизм
				atmzm :=chooseAutomatizmSuccessAndRun(purpose)
				return atmzm
			}
			if EvolushnStage == 3 {
				/*  ???? если пониженная мотивация что-то делать ничего не делать - лень
				 */
			}
		}
	}
	// ЛЕНЬ
	return nil
}
/////////////////////////////////////////////////////////////

// atmzm :=createAndRunAutomatizmFromPurpose(purpose)
func createAndRunAutomatizmFromPurpose(purpose *PurposeGenetic)(*Automatizm){
	atmzm := createAutomatizm(purpose)
	// запустить автоматизм
	if RumAutomatizm(atmzm) {
		// отслеживать последствия в automatizm_result.go
		setAutomatizmRunning(atmzm, purpose)
	}
	// в automatizm_result.go после оценки результата будет осмысление с активацией Дерева Понимания
	return atmzm
}
//////////////////////////////////////////////////////////

// выбрать из ранее удачного автоматизма, перекрыть цель новой и запустить новый автоматизм
func chooseAutomatizmSuccessAndRun(purpose *PurposeGenetic)(*Automatizm) {
	// ранее найденные удачные автоматизмы
	//  AutomatizmSuccessFromIdArr[n].GomeoIdSuccesArr[] - какие ID гомео-параметров улучшает это действие
	for _, v := range AutomatizmSuccessFromIdArr {
		targID := v.GomeoIdSuccesArr
		for i := 0; i < len(targID); i++ {
			if lib.ExistsValInArr(purpose.targetID, targID[i]) {
				// первый попавшися
				// создать новый автоматизм на основе успешного, но для данных условий и запустить его
				// TODO !не проверено!
				purpose.targetID = nil
				purpose.targetID = append(purpose.targetID, targID[i])
				// вытащить действия автоматизма
				trigID := CreateNewActionImageFromAutomatizm(v.Sequence)
				purpose.actionID = trigID
				atmzm := createAndRunAutomatizmFromPurpose(purpose)
				return atmzm
			}
		}
	}
	return nil
}
///////////////////////////////////////////////////////


/*  пробовать всякие случайныее простые действия, не повторяясь
Выдавая это на стадии 3, тварь получает реакцию оператора, которую отзеркаливает
 */
var usedActIdArr []int // какие деййствия уже были испробованы, погасить во сне wakingUp()
var usedPraseIdArr []int
func findAnySympleRandActions()(*Automatizm){

		// выдать массив возможных действий по ID парамктров гомеостаза как цели для улучшения в данных условиях
		targID,actID:=TerminalActions.GetSimpleActionForCurContitions()
   	// удалить уже использованное
		var tmp []int
		for i := 0; i < len(actID); i++ {
			if !lib.ExistsValInArr(usedActIdArr, actID[i]){
				tmp=append(tmp,actID[i])
			}
		}
		actID=tmp

	if len(actID)>0 {
		var actArrId []int
		if len(actID) > 2 { // сделать случайное сочетание
			actArrId1 := lib.RandChooseIntArr(actID)
			actArrId = append(actArrId, actArrId1)
			actArrId2 := lib.RandChooseIntArr(actID)
			actArrId = append(actArrId, actArrId2)
			actArrId = lib.UniqueArr(actArrId)

		} else {
			actArrId1 := lib.RandChooseIntArr(actID)
			actArrId = append(actArrId, actArrId1)
		}
		// чтобы не повторяться
		for i := 0; i < len(actArrId); i++ {
			usedActIdArr=append(usedActIdArr,actArrId[i])
		}
		var sequence = "Dnn:"
		for i := 0; i < len(actArrId); i++ {
			if i > 0 {
				sequence += ","
			}
			sequence += strconv.Itoa(actArrId[i])
		}
		var purpose PurposeGenetic
		purpose.targetID = targID
		// вытащить действия автоматизма
		trigID := CreateNewActionImageFromAutomatizm(sequence)
		purpose.actionID= trigID
		atmzm := createAndRunAutomatizmFromPurpose(&purpose)
		return atmzm
	}

	// если кончились действия, начали порверять фразы  ДЛЯ ПОПОЛНЕНИЯ type Verbal struct (verbalFromIdArr[])
// Выдавая это на стадии 3, тварь получает реакцию оператора, которую отзеркаливает
	for k, _ := range word_sensor.PhraseTreeFromID {
			if lib.ExistsValInArr(usedPraseIdArr, k){
				continue
			}
		usedPraseIdArr=append(usedPraseIdArr,k)
		sequence := "Snn:"+strconv.Itoa(k)
		var purpose PurposeGenetic
		purpose.targetID = targID
		// вытащить действия автоматизма
		trigID := CreateNewActionImageFromAutomatizm(sequence)
		purpose.actionID = trigID
		atmzm := createAndRunAutomatizmFromPurpose(&purpose)
		return atmzm
	}

	return nil
}
//////////////////////////////////////////////////////