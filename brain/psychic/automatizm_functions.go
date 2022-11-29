/* функции подержки автоматизмов

*/

package psychic

//////////////////////////////////////////



/* Не раздумывая, а рефлексторно используя имеющуюся информацию,
ВЫБРАТЬ ЛУЧШИЙ АВТОМАТИЗМ для узла nodeID то более ранних, если нет у поздних.
а если нет, то учитывать общие автомтизмы, привязанные к действиям (виртуальная ветка ID от 1000000) и словам (>2000000)
*/
func getAutomatizmFromNodeID(nodeID int)(int){
	// список всех автоматизмов для ID узла Дерева
	aArr:=GetMotorsAutomatizmListFromTreeId(nodeID)
	var usefulness = -10 // полезность, выбрать наилучшую
	var autmtzm *Automatizm
	if aArr != nil {
		for i := 0; i < len(aArr); i++ {
var allowRun=false
			if aArr[i].Usefulness < 0 {
				/* Не блокировать сразу, а посмотреть в Правила,
				м.б. после плохого эффекта последует следующее Правило с хорошим эффектом
				и тогда можно допустить Usefulness<0 в расчете на последующий успех.
				Не в качестве волевого усилия, а чисто автоматически использовать такую информацию.
				*/
				isWellEffect := isNextWellEffectFromActonRules(aArr[i].ActionsImageID)
				if isWellEffect {
					allowRun=true
				}
			}

			if aArr[i].Belief == 2 && allowRun { // есть штатный, проверенный автоматизм
				return aArr[i].ID
			}
			if aArr[i].Usefulness > usefulness {
				usefulness = aArr[i].Usefulness
				autmtzm = aArr[i]
			}
		}
	}
		if autmtzm !=nil { // выбран самый полезный из всех
			/*формирование не привязанных к узлу автоматизмов при активации дерева
			- для всех фраз - и для всех действий на основе привязанного автоматизма,
			чтобы другие ветки могли пользоваться при разных условиях.
			*/
			createNodeUnattachedAutomatizm(nodeID, autmtzm.ID)
			return autmtzm.ID
	}
	// в данном узле нет привязанного к нему автоматизма
	// если это - узел действий или узел фразы, смотрим, если привязанные к таким объектам автоматизм
	node:=AutomatizmTreeFromID[nodeID] // должен быть обязательно, но...
	if node == nil{
		return 0}
	if node.VerbalID>0 { // это узел фразы
		atmzS:=GetAutomatizmBeliefFromPhraseId(node.VerbalID)
		if atmzS != nil{
			return atmzS.ID //это - штатный автоматизм
		}
	}
	/////////////
	if node.ActivityID>0 && node.ToneMoodID==0 { // это узел действий - конечный в активной ветке
		atmzA:=GetAutomatizmBeliefFromActionId(node.ActivityID)
		if atmzA != nil{
			return atmzA.ID //это - штатный автоматизм
		}
	}
	//////////// нет штатных автоматизмов, выбрать любой нештатный на пробу
	/* такого быть не должно, т.к. штатный должен быть всегда
	if node.VerbalID>0 { // это узел фразы
		aArr = AutomatizmIdFromPhraseId[node.VerbalID]
		if aArr != nil {
			return aArr[0].ID // первый попавшийся не штатный, раз уже не нашелся штатный
		}
	}
	if node.ActivityID>0 && node.ToneMoodID==0 {
		aArr = AutomatizmIdFromActionId[node.VerbalID]
		if aArr != nil {
			return aArr[0].ID // первый попавшийся не штатный
		}
	}
	*/
	/////////// нет никаких автоматизмов хоть как-то относящихся к данному узлу
	// найти у предыдущих узел действий
	for i := len(ActiveBranchNodeArr)-1; i >2 ; i-- {
		node=AutomatizmTreeFromID[ActiveBranchNodeArr[i]]
		if node.ActivityID>0{
			atmzA:=GetAutomatizmBeliefFromActionId(node.ActivityID)
			if atmzA != nil{
				return atmzA.ID //это - штатный автоматизм
			}
			// не штатные автоматизмы для данного образа действий не будем смотреть
		}
	}

	return 0
}

/////////////////////////////////////




/////////////////////////////////////////////////
/*задать тип автоматизма Belief.
Только один из автоматизмов, прикрепленных к ветке или образу, может иметь Belief=2 - проверенное собственное знание
Если задается Belief=2, остальные Belief=2 становится Belief=0.
ТАК ПРОСТО НЕЛЬЗЯ ЗАДАВАТЬ Belief=2: LastAutomatizmWeiting.Belief=2
 */
func SetAutomatizmBelief(atmzm *Automatizm,belief int){
	if atmzm==nil || atmzm.BranchID==0{
		return
	}
if belief==2{
	// привязанные к ID узла дерева
	if atmzm.BranchID<1000000 {// обнулить Belief у всех привязанных к узлу
		aArr := GetMotorsAutomatizmListFromTreeId(atmzm.BranchID)
		if len(aArr) > 1 {
			for i := 0; i < len(aArr); i++ {
				if aArr[i] != atmzm && aArr[i].Belief == 2 {
					aArr[i].Belief = 0
					AutomatizmBelief2FromTreeNodeId[aArr[i].BranchID] = nil
				}
			}
		}
	// внизу	atmzm.Belief=2
		AutomatizmBelief2FromTreeNodeId[atmzm.BranchID] = atmzm
	}
	// привязанные к ID образа действий с пульта ActivityID
	if atmzm.BranchID>1000000 && atmzm.BranchID<2000000{// обнулить Belief у всех привязанных к ActivityID
		imgID:=atmzm.BranchID-1000000
		for _, v := range AutomatizmIdFromActionId[imgID] {
			v.Belief = 0
		}
	}
	if atmzm.BranchID>2000000{// обнулить Belief у всех привязанных к VerbalID
		imgID:=atmzm.BranchID-2000000
		for _, v := range AutomatizmIdFromPhraseId[imgID] {
			v.Belief = 0
		}
	}
}//if belief==2{
	atmzm.Belief=belief
}
/////////////////////////////////////////////////////




// список всех автоматизмов для ID узла Дерева
func GetMotorsAutomatizmListFromTreeId(nodeID int) []*Automatizm {
	if nodeID == 0 { return nil	}
	var mArr[] *Automatizm
	for _, a := range AutomatizmFromIdArr {
		if a.BranchID < 1000000 && a.BranchID == nodeID{
			mArr = append(mArr, a)
		}
	}
	return mArr
}
// штатный, невредный автоматизм, привязанный к ветке
func GetBelief2AutomatizmListFromTreeId(nodeID int) *Automatizm {
	if nodeID == 0 {
		return nil
	}
	aArr:=AutomatizmBelief2FromTreeNodeId[nodeID]

	if aArr == nil {
		return nil
	}
	if aArr.Usefulness >= 0 { // есть штатный, невредный
			return aArr
	}
	return nil
}
//////////////////////////////////////////////////

// есть ли штатный автоматизм (с Belief==2), привязанные к узлу дерева
func ExistsAutomatizmForThisNodeID(nodeID int)(bool){
	aArr:=AutomatizmBelief2FromTreeNodeId[nodeID]
	if aArr!=nil {
		return true
	}
	return false
}
///////////////////////////////////////

/* если для прикрепленных к узлу дерева есть карта штатных AutomatizmBelief2FromTreeNodeId,
то для прикрепленных к образам нужны ФУНКЦИИ ПОЛУЧЕНИЯ ШТАТНОГО ДЛЯ ДАННОГО ОБРАЗА
 */
func GetAutomatizmBeliefFromActionId(activityID int)(*Automatizm){
	if AutomatizmIdFromActionId[activityID] == nil{
		return nil
	}
	for _, v := range AutomatizmIdFromActionId[activityID] {
		if v.Belief == 2{
			return v
		}
	}
	return nil
}
///////////////////////////////////////////////////
func GetAutomatizmBeliefFromPhraseId(verbalID int)(*Automatizm){
	if AutomatizmIdFromPhraseId[verbalID] == nil{
		return nil
	}
	for _, v := range AutomatizmIdFromPhraseId[verbalID] {
		if v.Belief == 2{
			return v
		}
	}
	return nil
}
///////////////////////////////////////////////////

/*формирование не привязанных к узлу автоматизмов при активации дерева
- для всех фраз - и для всех действий на основе привязанного автоматизма,
чтобы другие ветки могли пользоваться при разных условиях.
*/
func createNodeUnattachedAutomatizm(nodeID int,aID int){
	node:=AutomatizmTreeFromID[nodeID] // должен быть обязательно, но...
	if node == nil{
		return }
	autmzm0:=AutomatizmFromIdArr[aID] // должен быть обязательно, но...
	if autmzm0 == nil{
		return }

	if node.VerbalID>0 { // это узел фразы
		_,autmzm:= CreateAtutomatizmNoSaveFile(2000000+node.VerbalID,autmzm0.ActionsImageID)
		if autmzm!=nil{
			SetAutomatizmBelief(autmzm, 2)// сделать автоматизм штатным
			autmzm.Usefulness=1 // полезность
		}
	}
	/////////////
	if node.ActivityID>0 && node.ToneMoodID==0 { // это узел действий - конечный в активной ветке
		_,autmzm:= CreateAtutomatizmNoSaveFile(1000000+node.ActivityID,autmzm0.ActionsImageID)
		if autmzm!=nil{
			SetAutomatizmBelief(autmzm, 2)// сделать автоматизм штатным
			autmzm.Usefulness=1 // полезность
		}
	}

}
////////////////////////////////////////////////////



// разблоикровака автоматизма для http://go/pages/automatizm_table.php
func UnblockAutomatizmID(atmtzmID int)string{
	atmtzm:=AutomatizmFromIdArr[atmtzmID]
	if atmtzm==nil{
		return "0"
	}
	atmtzm.Usefulness=1
	return "1"
}
/////////////////////////////////////////////////////////////////////


// привязать общий автоматизм к активной ветке detectedActiveLastNodID
func linkCoomonAtmtzmToBrench(commonAutomatizm *Automatizm){
	if LastAutomatizmWeiting.BranchID<1000000 { // это НЕ общий - не должно такого быть
		return
	}
	atmtzm := GetBelief2AutomatizmListFromTreeId(detectedActiveLastNodID)
	if atmtzm!=nil || atmtzm.ID == commonAutomatizm.ID{
		return
	}
	CreateAtutomatizmNoSaveFile(detectedActiveLastNodID,commonAutomatizm.ActionsImageID)
}
/////////////////////////////////////////////////////////////////////


/*Если есть ли автоматизм с действием curStimulImageID, и если у него atmtzm.Usefulness<0 - снять блокировку
потому как это - новое авторитарное подтвержение полезности.
Для текущей ветки дерева автоматизмов.
 */
func checkForUnbolokingAutomatizm(actID int){
	for _, v := range AutomatizmFromIdArr {
		if v.BranchID==detectedActiveLastNodID && v.ActionsImageID==actID{
			v.Usefulness=1 //авторитарно
			SetAutomatizmBelief(v, 2)// сделать штатным
		}
	}
}
///////////////////////////////////////////////////////