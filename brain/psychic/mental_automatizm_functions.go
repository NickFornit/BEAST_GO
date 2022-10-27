/* фнукции ментальных автоматизмов

*/

package psychic

//////////////////////////////////////////

/*задать тип автоматизма Belief.
Только один из автоматизмов, прикрепленных к ветке или образу, может иметь Belief=2 - проверенное собственное знание
Если задается Belief=2, остальные Belief=2 становится Belief=0.
ТАК ПРОСТО НЕЛЬЗЯ ЗАДАВАТЬ Belief=2
*/
func SetMentalAutomatizmBelief(atmzm *MentalAutomatizm,belief int){
	if atmzm==nil {
		return
	}
	if belief==2{
		MentalAutomatizmBelief2FromTreeNodeId[atmzm.ID] = atmzm
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


/* список удачных автоматизмов, относящихся к определенным условиям )
В этом списке поле Usefulness >0
*/
var MentalAutomatizmSuccessFromIdArr = make(map[int]*MentalAutomatizm)


/////////////////////////////////////




