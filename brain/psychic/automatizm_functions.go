/* функции подержки автоматизмов

*/

package psychic

import (
	"strconv"
	"strings"
)

//////////////////////////////////////////



/* выбрать лучший автоматизм для узла nodeID то более ранних, если нет у поздних.
*/
func getAutomatizmFromNodeID(nodeID int)(int){
	// список всех автоматизмов для ID узла Дерева
	aArr:=GetMotorsAutomatizmListFromTreeId(nodeID)
	var usefulness =-10 // полезность, выбрать наилучшую
	var usefulnessID=0
	for i := 0; i < len(aArr); i++ {
		if aArr[i].Belief==2{// есть единственный проверенный автоматизм
			return aArr[i].ID
		}
		if aArr[i].Usefulness > usefulness{
			usefulness=aArr[i].Usefulness
			usefulnessID=aArr[i].ID
		}
	}
	if usefulnessID >0{// выбран самый полезный из всех
		/*формирование не привязанных к узлу автоматизмов при активации дерева
		- для всех фраз - и для всех действий на основе привязанного автоматизма,
		чтобы другие ветки могли пользоваться при разных условиях.
		*/
		createNodeUnattachedAutomatizm(nodeID, usefulnessID)
		return usefulnessID
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
	/////////// нет никаких автоматизмов хоть как-то относящийся к данному узлу
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


/* разделить строку Sequence автоматизма на составляющие
типы действий:
1 Snn - перечень ID фраз через запятую
2 Dnn - ID прогрмаммы действий, через запятую
3 Ann - последовательный запуск автоматизмов с id1,id2..
4 Mnn - внутренние произвольные действия с id1,id2...
5 Tnn - образ тон-настроения одна цифра == образ тона-настроения (как в func GetToneMoodID(  и func GetToneMoodFromImg()
*/
func ParceAutomatizmSequence(Sequence string)([]ActsAutomatizm){
	var acts[] ActsAutomatizm

	sArr:=strings.Split(Sequence, "|")
	for i := 0; i < len(sArr); i++ {
		if len(sArr[i])==0{
			continue
		}
		var act ActsAutomatizm
		pArr:=strings.Split(sArr[i], ":")
		switch pArr[0]{
		case "Snn": act.Type=1
		//case "nn": act.Type=2
		case "Dnn": act.Type=2
		case "Ann": act.Type=3
		case "Mnn": act.Type=4
		case "Tnn": act.Type=5
		}

		act.Acts = pArr[1] // строка действий (любого типа) через запятую
		acts = append(acts, act)
	}
	return acts
}
////////////////////////////////////////////////


/* получить массив wordID из Sequence автоматизма
 */
func GetWordArrFromSequence(sequence string)([]int){
	sArr:=strings.Split(sequence, "|")
	for i := 0; i < len(sArr); i++ {
		if len(sArr[i])==0{
			continue
		}
		pArr:=strings.Split(sArr[i], ":")
		if pArr[0]=="Snn"{
			sA:=strings.Split(pArr[1], ",")
			if len(sA)>0{
				var out[]int
				for a := 0; a < len(sA); a++ {
					wID,_:=strconv.Atoi(sA[a])
					out = append(out, wID)
				}
				return out
			}
		}
	}
	return nil
}

func GetAutomatizmSnn(ma *Automatizm)(string){
	sequence:=ma.Sequence // Sequence="Snn:24243,1234,0,24234,11234|Dnn:24,4,56"
	sArr:=strings.Split(sequence, "|")
	for i := 0; i < len(sArr); i++ {
		if len(sArr[i])==0{
			continue
		}
		pArr:=strings.Split(sArr[i], ":")
		if pArr[0]=="Snn"{
			if len(pArr[1])>0{
				return pArr[1]
			}
		}
	}
	return ""
}

func GetAutomatizmDnn(ma *Automatizm)(string){
	sequence:=ma.Sequence
	sArr:=strings.Split(sequence, "|")
	for i := 0; i < len(sArr); i++ {
		if len(sArr[i])==0{
			continue
		}
		pArr:=strings.Split(sArr[i], ":")
		if pArr[0]=="Dnn"{
			if len(pArr[1])>0{
				return pArr[1]
			}
		}
	}
	return ""
}


func GetAutomatizmTnn(ma *Automatizm)(string){
	sequence:=ma.Sequence
	sArr:=strings.Split(sequence, "|")
	for i := 0; i < len(sArr); i++ {
		if len(sArr[i])==0{
			continue
		}
		pArr:=strings.Split(sArr[i], ":")
		if pArr[0]=="Tnn"{
			if len(pArr[1])>0{
				return pArr[1]
			}
		}
	}
	return ""
}
//////////////////////////////////////////////////



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
		_,autmzm:=CreateAutomatizm(2000000+node.VerbalID,autmzm0.Sequence)
		if autmzm!=nil{
			SetAutomatizmBelief(autmzm, 2)// сделать автоматизм штатным
			autmzm.Usefulness=1 // полезность
		}
	}
	/////////////
	if node.ActivityID>0 && node.ToneMoodID==0 { // это узел действий - конечный в активной ветке
		_,autmzm:=CreateAutomatizm(1000000+node.ActivityID,autmzm0.Sequence)
		if autmzm!=nil{
			SetAutomatizmBelief(autmzm, 2)// сделать автоматизм штатным
			autmzm.Usefulness=1 // полезность
		}
	}

}
////////////////////////////////////////////////////




