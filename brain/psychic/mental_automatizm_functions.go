/* фнукции ментальных автоматизмов

*/

package psychic

import "strings"

//////////////////////////////////////////

/*задать тип автоматизма Belief.
Только один из автоматизмов, прикрепленных к ветке или образу, может иметь Belief=2 - проверенное собственное знание
Если задается Belief=2, остальные Belief=2 становится Belief=0.
ТАК ПРОСТО НЕЛЬЗЯ ЗАДАВАТЬ Belief=2
*/
func SetMentalAutomatizmBelief(atmzm *MentalAutomatizm,belief int){
	if atmzm==nil || atmzm.BranchID==0{
		return
	}
	if belief==2{
		// привязанные к ID узла дерева
		aArr := GetMentalAutomatizmListFromTreeId(atmzm.BranchID)
		if len(aArr) > 1 {
			for i := 0; i < len(aArr); i++ {
				if aArr[i] != atmzm && aArr[i].Belief == 2 {
					atmzm.Belief = 0
					MentalAutomatizmBelief2FromTreeNodeId[atmzm.BranchID] = nil
				}
			}
		}
		MentalAutomatizmBelief2FromTreeNodeId[atmzm.BranchID] = atmzm
	}//if belief==2{
	atmzm.Belief=belief
}
//////////////////////////////////////////////////////////////////

// есть ли штатный автоматизм (с Belief==2), привязанные к узлу дерева
func ExistsMentalAutomatizmForThisNodeID(nodeID int)(bool){
	aArr:=MentalAutomatizmBelief2FromTreeNodeId[nodeID]
	if aArr!=nil {
		return true
	}
	return false
}
///////////////////////////////////////


/* список удачных автоматизмов, относящихся к определенным условиям (привзяанных к определенной ветке Дерева)
В этом списке поле Usefulness >0
*/
var MentalAutomatizmSuccessFromIdArr = make(map[int]*MentalAutomatizm)

// список всех автоматизмов для ID узла Дерева
func GetMentalAutomatizmListFromTreeId(nodeID int) []*MentalAutomatizm {
	if nodeID == 0 { return nil	}
	var mArr[] *MentalAutomatizm
	for _, a := range MentalAutomatizmsFromID {
		if a.BranchID == nodeID{
			mArr = append(mArr, a)
		}
	}
	return mArr
}
//////////////////////////////////////////////


/* выбрать лучший автоматизм для узла nodeID то более ранних, если нет у поздних.
 */
func getMentalAutomatizmFromNodeID(nodeID int)(int){
	// список всех автоматизмов для ID узла Дерева
	aArr:=GetMentalAutomatizmListFromTreeId(nodeID)
	var usefulness =-10 // полезность, выбрать наилучшую
	var usefulnessID=0
	for i := 0; i < len(aArr); i++ {
		if aArr[i].Belief==2{// есть штатный автоматизм
			return aArr[i].ID
		}
		if aArr[i].Usefulness > usefulness{
			usefulness=aArr[i].Usefulness
			usefulnessID=aArr[i].ID
		}
	}
	if usefulnessID >0{// выбран самый полезный из всех
		return usefulnessID
	}
	// нет никаких автоматизмов хоть как-то относящийся к данному узлу
	return 0
}
/////////////////////////////////////



/* разделить строку Sequence автоматизма на составляющие
типы действий:
1 Mnn - выполнить ментальную функцию с ID
2 Ann - выполнить моторный автоматизм с ID
*/
func ParceMentalAutomatizmSequence(Sequence string)([]ActsMentalAutomatizm){
	var acts[] ActsMentalAutomatizm

	sArr:=strings.Split(Sequence, "|")
	for i := 0; i < len(sArr); i++ {
		if len(sArr[i])==0{
			continue
		}
		var act ActsMentalAutomatizm
		pArr:=strings.Split(sArr[i], ":")
		switch pArr[0]{
		case "Mnn": act.Type=1
		case "Ann": act.Type=2
		}

		act.Acts = pArr[1] // строка действий (любого типа) через запятую
		acts = append(acts, act)
	}
	return acts
}
////////////////////////////////////////////////
