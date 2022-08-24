/*  Вспомогательные функции Ориентировочных рефлексоы

*/


package psychic

import (
	"strconv"
)

/////////////////////////////////////////////////////////


//////////////////////////////////////////////
/* ТОЛЬКО ДЛЯ orientation_1(), когда автоматизма нет у недоделанной ветки!
сформировать пробный автоматизм моторного действия и сразу запустить его в действие
   Зафиксироваь время действия
   10 пульсов следить за измнением жизненных параметров и ответными действиями - считать следствием действия
   оценить результат и скорректировать силу мозжечком в записи автоматизма.
*/
func createAutomatizm(pc *PurposeGenetic)(*Automatizm){

	BranchID:=detectedActiveLastNodID
	/*
		var existsAutomatizmID=0
		if AutomatizmTreeFromID[BranchID].ActivityID>0{
			existsAutomatizmID=AutomatizmTreeFromID[BranchID].ActivityID
		}
		// у ветки есть автматизм, но в случае  orientation_1() автоматизма нет у недоделанной ветки
		if existsAutomatizmID>0{

		}
	*/
	var Sequence=""
	aArr:=pc.actionID.ActID
	for i := 0; i < len(aArr); i++ {
		if i>0 {Sequence += "|"}
		Sequence += "Dnn:" + strconv.Itoa(aArr[i])
	}

	// создать автоматизм, даже если такой уже есть
	_,atzm:=CreateNewAutomatizm(BranchID,Sequence)
	if atzm!=nil {
		atzm.Energy = 5

		return atzm
	}

return nil
}
//////////////////////////////////////////////


/////////////////////////////////////
/*подобрать по тону и настроению хоть как-то ассоциирующуюся фразу из имеющихся
Tone int //Тон: 0 - обычный, 1 - восклицательный, 2- вопросительный, 3- вялый, 4 - Повышенный
Mood int // настроение при передаче фразы с Пульта: 20-Хорошее    21-Плохое    22-Игровое    23-Учитель    24-Агрессивное   25-Защитное    26-Протест
 */
func findSuitablePhrase()([]int){
	var ToneID=0
	var MoodID=0
	if PsyBaseMood==-1{// плохое настроение
		MoodID=21
		ToneID=4
		if CurrentInformationEnvironment.danger { // опасность состояния
			ToneID=1
			MoodID=25 // защитное
		}
	}
	if PsyBaseMood==0{// нормальное настроение
		MoodID=0
		ToneID=0
		if CurrentInformationEnvironment.danger { // опасность состояния
			ToneID=4
			MoodID=24 // защитное
		}
	}
	if PsyBaseMood==1{// хорошее настроение
		MoodID=20
		ToneID=4
	}
	for _, v := range VerbalFromIdArr {
		if v.ToneID==ToneID && v.MoodID==MoodID{
			return v.PhraseID
		}
	}

	return nil
}
///////////////////////////////////////////////


////////////////////////////////////////////////
/* найти важные (по опыту) признаки в новизне NoveltySituation
Это - чисто рефлексторный процесс поиска в опыте
 */
func getImportantSigns()([]int){
	// NoveltySituation

	return nil
}