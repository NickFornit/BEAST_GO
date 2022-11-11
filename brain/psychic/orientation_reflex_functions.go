/*  Вспомогательные функции Ориентировочных рефлексоы

*/


package psychic

import (

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

	aArr:=pc.actionID.ActID

	sArr:=pc.actionID.PhraseID

	// тон и настроение
	t:=pc.actionID.ToneID
	m:=pc.actionID.MoodID-19

	ActionsImageID,_:=CreateNewlastActionsImageID(0,aArr,sArr,t,m,true)
	// создать автоматизм
	_,atzm:=CreateNewAutomatizm(BranchID,ActionsImageID)
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
func getImportantSigns()(bool){
	lenN:=len(CurrentAutomatizTreeEnd)
if lenN==0{
	return false
}
	var news []int // выделить новизну
	for i := 0; i < lenN; i++ {
		if CurrentAutomatizTreeEnd[i]>0{
			news=append(news,CurrentAutomatizTreeEnd[i])
		}
	}
	if news==nil{
		return false
	}
// выделить признаки и оценить важность
switch currentStepCount{
case 3: //остается после lev3 - ActivityID образ сочетания пусковых стимулов (только кнопок!)
/* пример CurrentAutomatizTreeEnd:
   lev4 90   GetToneMoodID(verb.ToneID, verb.MoodID)
   lev5 17   verb.SimbolID
   lev6 132  verb.PhraseID[0]
 */
	ton,mood:=GetToneMoodFromImg(90)
// 1-восклицательный,4 - Повышенный, настроение: 2-Плохое, 5-Агрессивное, 6-Защитное
	if ton==1 || ton==4 || mood==2 || mood==5 || mood==6{
		return true
	}
	if EvolushnStage >3 {
		// TODO опредять по словарному запасу в контексте дерева понимания
	}
case 4: //SimbolID
	/* пример CurrentAutomatizTreeEnd:
	   lev5 17   verb.SimbolID
	   lev6 132  verb.PhraseID[0]
	*/
	if EvolushnStage >3 {
		// TODO опредять по словарному запасу в контексте дерева понимания
	}

case 5: // ToneMoodID
	/* пример CurrentAutomatizTreeEnd:
	   lev6 132  verb.PhraseID[0]
	*/
	if EvolushnStage >3 {
		// TODO опредять по словарному запасу в контексте дерева понимания
	}

default: return false
}

	return false
}