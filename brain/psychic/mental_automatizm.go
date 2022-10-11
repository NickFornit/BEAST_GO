/* Ментальные (умственные) автоматизмы мышления.
мент.автоматизм может прикрепляться ТОЛЬКО к последнему узлу ветки - при полном понимании ситуации
*/

package psychic

///////////////////////////////////////
type MentalAutomatizm struct {
	ID int
	BranchID int
	Usefulness int // (БЕС)ПОЛЕЗНОСТЬ: -10 вред 0 +10 +n польза

	/* Цепочка последовательности реагирования,
Mnn:ID - выполнить ментальную функцию с ID
Ann:ID - выполнить моторный автоматизм с ID
	*/
	Sequence string

	/* Следующий автоматизм в цепочке исполнения. Цепь может быть пройдена ментально, без выполнения автоматизмов, для этого не вызывается моторное выполнение а просто - проход цепочки с просмотром ее звеньев.
	или цепь может быть прервана осознанно
	или пройдена при ее пошаговом отслеживании StepByStepAutomatizm
	и пошаговом запуске: allowNextAutomatizm(automatizm.NextID):
	Бот смотрит, выполнить ли следующий шаг, добавить ли рефлекс мозжечка.
	*/
	NextID int

// СОБСТВЕННОЙ ЭНЕРГИИ НЕТ	Energy int // от 1 до 10, по умолчанию = 5

	/* Уверенность в авторитарном автоматизме высока в период авторитарного обучения
	   	и сильно падает в период собственной инициативы, когда нужно на себе проверить,
	   	а даст ли такой автоматизм в самом деле обещанное улучшение.
Только один из автоматизмов, прикрепленных к ветке, может иметь Belief=2 - проверенное собственное знание
Если задается Belief=2, остальные Belief=2 становится Belief=0.
!!! ПОЭТОМУ ВСЕГДА нужно задавать с помощью SetMentalAutomatizmBelief(atmzm *Automatizm,belief int)
	*/
	Belief int // 0 - предположение, 1 - чужие сведения, 2 - проверенное собственное знание - ШТАТНЫЙ, умолчательный автоматизм
	/* В случае, если в результате автоматизма его Usefulness изменит знак, то
	Count обнулится, а при таком же знаке - увеличивается на 1.
	*/
	Count int // надежность: число использований с подтверждением (бес)полезности Usefulness
}
// привязанные к узлу дерева
var MentalAutomatizmsFromID=make(map[int]*MentalAutomatizm)
///////////////////////////////////////////////////
// ШТАТНЫЕ автоматизмы, прикрепленные к ID узла Дерева с Belief==2 т.е. ШТАТНЫЕ, выполняющиеся не раздумывая
// у узла может быть только один штатный автоматизм с Belief==2
var MentalAutomatizmBelief2FromTreeNodeId = make(map[int]*MentalAutomatizm)

/* список удачных автоматизмов, относящихся к определенным условиям (привзяанных к определенной ветке Дерева)
В этом списке поле Usefulness >0
*/
var MentalAutomatizmSuccessFromIdArr = make(map[int]*MentalAutomatizm)

// список всех автоматизмов для ID узла Дерева
func GetMentalAutomatizmListFromTreeId(nodeID int) []*MentalAutomatizm {
	if nodeID == 0 { return nil	}
	var mArr[] *MentalAutomatizm
	for _, a := range MentalAutomatizmsFromID {
		if a.BranchID < 1000000 && a.BranchID == nodeID{
			mArr = append(mArr, a)
		}
	}
	return mArr
}


func SaveMentalAutomatizms(){

}

func RunMentalAutomatizmsID(mentID int){
	am:=MentalAutomatizmsFromID[mentID]

	// вернуть скорректированную силу действия
	//addE:=0
	//if am.Belief!=3 {// не рефлекс мозжечка
		addE := getCerebellumReflexAddEnergy(1,am.ID)
	//}
	if addE>0{

	}


	//  выполнить дополнительные мозжечковые автоматизмы сразу после выполняющегося автоматизма
	runCerebellumAdditionalAutomatizm(1,mentID)
}
////////////////////////////////////////////////////////////////


