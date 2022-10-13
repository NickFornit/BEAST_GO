/* Образ текущей ситуации для осмысления или Образ Модели Понимания
4- уровень Дерева Понимания (или дерева ментальных автоматизмов)
Здесь активная ID дерева автоматизмов (доступны моторные автоматизмы от дерева автоматизмов)

*/

package psychic

import (
	"BOT/lib"
	"strconv"
	"strings"
)

/* Образы текущей ситуации сохраняются в памяти
т.к. представляют собой 4- уровень Дерева Понимания (или дерева ментальных автоматизмов)
*/
type SituationImage struct {
	ID int
/* позволяет получить доступ к конечному узлу ветки дерева моторных автоматизмов и получить все инфу ее узлов
Здесь - вся конкретика ситуации.
*/
	autmzmTreeNodeID int

/* Смысловой контекст ситуации, то, на что в первую очередь обращается внимание (приоритет):
	1 - было ответное действие (смотреть в автоматизме ветки Usefulness int - (БЕС)ПОЛЕЗНОСТЬ: вред: -10 0 +10 +n польза diffPsyBaseMood )
	2 - был запуск автоматизма ветки.
	3 - ничего не делали, но нужно осмысление
	4 - все спокойно, можно экспериментивароть
	5 - оператор не прореагировал на действия в течение периода ожидания - игнорирует? нужно достучаться?

	10-17 - оператор выбрал настроение (14 - Учитель при отправке или нажал кнопку Поучить)
	20-37 - оператор нажал кнопку (17 - Игровое при отправке или нажал кнопку Поиграть)
	... и т.п.
*/
	SituationType int

}
var SituationImageFromIdArr=make(map[int]*SituationImage)
/////////////////////////////////


var lastSituationImageID=0
func createSituationImage(id int,autmzmTreeNodeID int,SituationType int,save bool)(int,*SituationImage){
	oldID,oldVal:=checkUnicumSituationImage(autmzmTreeNodeID,SituationType)
	if EvolushnStage < 4 { // только со стадии развития 4
		return 0,nil
	}
	if PulsCount<4{// не активировать пока все не устаканится
		return  0,nil
	}
	/////////////////////////////////////
	if oldVal!=nil{
		return oldID,oldVal
	}
	if id==0{
		lastSituationImageID++
		id=lastSituationImageID
	}else{
		//		newW.ID=id
		if lastSituationImageID<id{
			lastSituationImageID=id
		}
	}

	var node SituationImage
	node.ID = id
	node.autmzmTreeNodeID = autmzmTreeNodeID
	node.SituationType = SituationType

	SituationImageFromIdArr[id]=&node

	if save {
		SaveSituationImage()
	}

	return id,&node
}
/////////////////////////////////////
func checkUnicumSituationImage(autmzmTreeNodeID int,SituationType int)(int,*SituationImage){
	for id, v := range SituationImageFromIdArr {
		if autmzmTreeNodeID!=v.autmzmTreeNodeID && SituationType!=v.SituationType {
			continue
		}
		return id,v
	}

	return 0,nil
}
////////////////////////////////////


func SaveSituationImage(){
	var out=""
	for k, v := range SituationImageFromIdArr {
		out+=strconv.Itoa(k)+"|"
		out+=strconv.Itoa(v.autmzmTreeNodeID)+"|"
		out+=strconv.Itoa(v.SituationType)+"|"
		out+="\r\n"
	}
	lib.WriteFileContent(lib.GetMainPathExeFile()+"/memory_psy/situation_images.txt",out)
}
//////////////////////////////////////////


func loadSituationImage(){
	SituationImageFromIdArr=make(map[int]*SituationImage)
	strArr,_:=lib.ReadLines(lib.GetMainPathExeFile()+"/memory_psy/situation_images.txt")
	cunt:=len(strArr)
	for n := 0; n < cunt; n++ {
		if len(strArr[n])==0{
			continue
		}
		p:=strings.Split(strArr[n], "|")
		id,_:=strconv.Atoi(p[0])
		autmzmTreeNodeID,_:=strconv.Atoi(p[1])
		SituationType,_:=strconv.Atoi(p[2])

		createSituationImage(id,autmzmTreeNodeID,SituationType,false)
	}
	return

}
//////////////////////////////


//////////// ПРИОРИТЕТЫ СИТУАЦИЙ
/* приоритет в зависимости от ID параметра.
prioritetOfPultButtonActions имеет преимущество перед prioritetOfPultMoodActions
Это - наследственное пердопределение.
*/
// 0-7
func getPrioritetOfPultMoodActions(moodID int)(int){
switch moodID{
case 1: return 1 //Хорошее
case 6: return 2 //Защитное
case 7: return 3 //Протест
case 3: return 4 //Игровое
case 4: return 5 //Учитель
case 2: return 6 //Плохое
case 5: return 7 //Агрессивное
default: return 0
}
return 0
}
// 1 - 17
func getPrioritetOfPultButtonActions(actionID int)(int){
switch actionID{
case 6: return 1 //Успокоить
case 16: return 2 //Простить
case 9: return 3 //Игнорировать
case 11: return 4 //Сделать приятно
case 2: return 5 //Наказать
case 1: return 6  //Непонятно
case 17: return 7 //Вылечить
case 5: return 8 //Накормить
case 14: return 9 //Обрадоваться
case 15: return 10 //Испугаться
case 12: return 11 //Заплакать
case 7: return 12 //Предложить поиграть
case 8: return 13 //Предложить поучить
case 13: return 3 //Засмеяться
case 4: return 15 //Поощрить
case 10: return 16 //Сделать больно
case 3: return 17 //Наказать

default: return 0
}
return 0
}
/////////////////////////////////////////

/* определить ID ситуации: настроение при посылке сообщения, нажатые кнопки и т.п.
Может быть выбрана только одна из существующих ситуаций, поэтому выбор идет по приоритетам.
для определения узла SituationID дерева.
Это определяет контекст ситуации, при вызове активации дерева понимания.
В getCurSituationImageID() по-началу выбирается наугад (для первого приближения) более важные из существующих,
но потом дерево понимания может переактивироваться с произвольным заданием контекста.
От этого параметра зависит в каком направлении пойдет информационный поиск решений,
если не будет запущен штатный автоматизм ветки (ориентировочные реакции).
Инфа curActiveActions обновляется при каждой активации дерева моторных автоматизмов.

Функция возвращает предположительный ID смыслового контекста ситуации:
	1 - было ответное действие (смотреть в автоматизме ветки Usefulness int - (БЕС)ПОЛЕЗНОСТЬ: вред: -10 0 +10 +n польза diffPsyBaseMood )
	2 - был запуск автоматизма ветки.
	3 - ничего не делали, но нужно осмысление
	4 - все спокойно, можно экспериментивароть
	5 - оператор не прореагировал на действия в течение периода ожидания - игнорирует? нужно достучаться?

	10-17 - оператор выбрал настроение (14 - Учитель при отправке)
	20-37 - оператор нажал кнопку (17 - Игровое при отправке или нажал кнопку Поиграть)
	... и т.п.
*/
func getCurSituationImageID()(int){
	if detectedActiveLastNodID==0{
		return 0
	}
var sitID=0

if LastRunAutomatizmPulsCount > 0{// был и закончился ответом период ожидания на действия автоматизма
	// вышло время ожидания реакции
	if (LastRunAutomatizmPulsCount+WaitingPeriodForActionsVal) < PulsCount {
// оператор не прореагировал на действия в течение периода ожидания - игнорирует? нужно достучаться?
		id, _ := createSituationImage(0, detectedActiveLastNodID, 5,true)
		if id > 0 {
			return id
		}
	}
//было ответное действие (смотреть в автоматизме ветки Usefulness int - (БЕС)ПОЛЕЗНОСТЬ: вред: -10 0 +10 +n польза diffPsyBaseMood )
	id, _ := createSituationImage(0, detectedActiveLastNodID, 1,true)
	if id > 0 {
		return id
	}
}//if LastRunAutomatizmPulsCount > 0

	// ЕСТЬ ЛИ АВТОМАТИЗМ В моторной ВЕТКЕ и болеее ранних?
	if currentAutomatizmAfterTreeActivatedID > 0 {
		//был запуск автоматизма ветки
		id, _ := createSituationImage(0, detectedActiveLastNodID, 2,true)
		if id > 0 {
			return id
		}
	}


/////////////// ситуации действий с пульта
var max=0
// сначала настроение, чтобы оно перекрылось кнопками действий если они есть
	mood:=curActiveActions.moodID
	if mood != 0{// есть настроение
		prior:=getPrioritetOfPultMoodActions(mood)
			if max<prior{
				sitID=10+mood
				max=prior
			}
	}/////////////////////

// кнопки действий
aArr:=curActiveActions.actID
	if aArr != nil{// есть действия кнопок с Пульта
		max=0 // перекрываем
		for i := 0; i < len(aArr); i++ {
			prior:=getPrioritetOfPultButtonActions(aArr[i])
if max<prior{
	sitID=20+aArr[i]
	max=prior
}
		}
	}/////////////////////

if sitID>0 {
	id, _ := createSituationImage(0, detectedActiveLastNodID, sitID,true)
	if id > 0 {
		return id
	}
}//////////////////////////////

	if detectedActiveLastNodID == 0 {
		// ничего не делали, но нужно осмысление
		id, _ := createSituationImage(0, detectedActiveLastNodID, 3,true)
		if id > 0 {
			return id
		}
	}

	// все спокойно, можно экспериментивароть
	id, _ := createSituationImage(0, detectedActiveLastNodID, 4,true)
	if id > 0 {
		return id
	}

return 0
}
/////////////////////////////////////////////////////////////////////
