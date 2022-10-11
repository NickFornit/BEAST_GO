/* Ментальные (умственные) автоматизмы мышления.

*/

package psychic

///////////////////////////////////////
type MentalAutomatizm struct {
	ID int
	/*
	   	   Уверенность в авторитарном автоматизме высока в период авторитарного обучения
	   	и сильно падает в период собственной инициативы, когда нужно на себе проверить,
	   	а даст ли такой автоматизм в самом деле обещанное улучшение.
Только один из автоматизмов, прикрепленных к ветке, может иметь Belief=2 - проверенное собственное знание
Если задается Belief=2, остальные Belief=2 становится Belief=0.
!!! ПОЭТОМУ ВСЕГДА нужно задавать с помощью SetAutomatizmBelief(atmzm *Automatizm,belief int)
	*/
	Belief int // 0 - предположение, 1 - чужие сведения, 2 - проверенное собственное знание, 3 - для мозжечкового рефлекса

}

var MentalAutomatizmsFromID=make(map[int]*MentalAutomatizm)
///////////////////////////////////////////////////


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