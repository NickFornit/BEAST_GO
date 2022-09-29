/* Словестные образы (область Брока)  для 4, 5 и 6-го уровня дерева автоматизмов.
Смысл (осознанную значимость) образ приобретеает только в контексте Дерева Понимания (дерева мент.автоматизмов)

Детекторы зоны Вернике распознают слова и словосочетания, 
а область Брока отвечает за смысл распознанных слов и словосочетений,
за конструирование собственных словосочетаний,
за моторное использование слов и словосочетаний.
За все ответственная структура - образ осмысленных слов и сочетаний.

! Нужно иметь в виду, что в Vernike_detector.go есть массив памяти фраз, накапливается в течении дня
var MemoryDetectedArr []MemoryDetected - структур фразы с контекстным окружением
и Verbal.PhraseID[] - можно найти в этом массиве для ориентировки что бы ло раньше и позже.
MemoryDetectedArr - как бы оперативная память фраз для сопоставлений.
*/

package psychic

import (
	"BOT/lib"
	"strconv"
	"strings"
)

/////////////////////////////////////////////

func verbalInit(){
	loadVerbalFromIdArr()

/*
var tm=922// "Обычный, Хорошее"
	str:=getToneMoodStrFromID(tm)
	if len(str)>0{	}
 */

}
////////////////////////////////////////////

/* для оптимизации поиска по дереву перед узлом Verbal идет узел первого символа : var symbolsArr из word_tree.go
Смысл (осознанную значимость) образ приобретеает только в контексте Дерева Понимания (дерева мент.автоматизмов)
 */
type Verbal struct {
	ID int
	// для текущего сообщения с Пусльта:
	SimbolID int // id первого символа первой фразы PhraseID: var symbolsArr из word_tree.go
	PhraseID[] int // массив фразID (DetectedUnicumPhraseID) слова каждой фразы вытаскиваются wordSensor.GetWordArrFromPhraseID(PhraseID[0])
//0 - обычный, 1 - восклицательный, 2- вопросительный, 3- вялый, 4 - Повышенный	
	ToneID int // тон сообщения с Пульта
//0 - обычный, 1 (из 20)-Хорошее    2 (21)-Плохое    3(22)-Игровое    4(23)-Учитель
//5(24)-Агрессивное   6(25)-Защитное    7(26)-Протест
	MoodID int // настроение оператора
}
var VerbalFromIdArr=make(map[int]*Verbal)
//////////////////////////////////////////
// для поиска по ID фразы
//var VerbalFromPhraseIdArr=make(map[int]*Verbal)
///////////////////////////////////////////

// создать образ сочетаний пусковых стимулов
//В случае отсуствия пусковых стимулов создается ID такого отсутсвия, пример такой записи: 2|||0|0| - ID=2
var lastVerbalID=0
func createNewlastVerbalID(id int,SimbolID int,PhraseID []int,ToneID int,MoodID int)(int,*Verbal){
	oldID,oldVal:=checkUnicumVerbal(PhraseID,ToneID,MoodID)
	if oldVal!=nil{
		return oldID,oldVal
	}
	if id==0{
		lastVerbalID++
		id=lastVerbalID
	}else{
		//		newW.ID=id
		if lastVerbalID<id{
			lastVerbalID=id
		}
	}

	var node Verbal
	node.ID = id
	node.SimbolID=SimbolID
	node.PhraseID = PhraseID
	node.ToneID=ToneID
	node.MoodID=MoodID
	if MoodID>19 {
		MoodID= MoodID - 19
	}else{
		MoodID=0
	}

	VerbalFromIdArr[id]=&node
	return id,&node
}
func checkUnicumVerbal(PhraseID []int,ToneID int,MoodID int)(int,*Verbal){
	for id, v := range VerbalFromIdArr {
		if !lib.EqualArrs(PhraseID,v.PhraseID) {
			continue
		}
		if ToneID!=v.ToneID || MoodID!=v.MoodID {
			continue
		}
		return id,v
	}

	return 0,nil
}
/////////////////////////////////////////
// создать новый вербальный образ, если такого еще нет
func CreateVerbalImage(FirstSimbolID int,PhraseID []int,ToneID int,MoodID int)(int,*Verbal){
	if PhraseID==nil{
		return 0,nil
	}
	// достаем первый символ первой фразы
	// получить последователньость wordID из уникального идентификатора первой фразы
/*
	wordIDarr:=wordSensor.GetWordArrFromPhraseID(PhraseID[0])
	// первое слово в виде строки
	if wordIDarr==nil || len(wordIDarr)==0{
		return 0,nil
	}
	word:=wordSensor.GetWordFromWordID(wordIDarr[0])
	//rw:=[]rune(word)
	//SimbolID:=wordSensor.GetSymbolIDfromRune(rw[0])
//	word:=wordSensor.GetPhraseStringsFromPhraseID(PhraseID[0])
	//SimbolID:=wordSensor.GetSymbolIDfromString(rw[0])
*/
	id,verb:=createNewlastVerbalID(0,FirstSimbolID,PhraseID,ToneID,MoodID)

	SaveVerbalFromIdArr()

	return id,verb
}

/////////////////////////////////////////

//////////////////// сохранить образы сочетаний пусковых стимулов
//В случае отсуствия пусковых стимулов создается ID такого отсутсвия, пример такой записи: 2|||0|0|
func SaveVerbalFromIdArr(){
	var out=""
	for k, v := range VerbalFromIdArr {
		out+=strconv.Itoa(k)+"|"
		out+=strconv.Itoa(v.SimbolID)+"|"
		for i := 0; i < len(v.PhraseID); i++ {
			out+=strconv.Itoa(v.PhraseID[i])+","
		}
		out+="|"
		out+=strconv.Itoa(v.ToneID)+"|"
		out+=strconv.Itoa(v.MoodID)+"|"
		out+="\r\n"
	}
	lib.WriteFileContent(lib.GetMainPathExeFile()+"/memory_psy/verbal_images.txt",out)

}
////////////////////  загрузить образы сочетаний пусковых стимулов
func loadVerbalFromIdArr(){
	VerbalFromIdArr=make(map[int]*Verbal)
	strArr,_:=lib.ReadLines(lib.GetMainPathExeFile()+"/memory_psy/verbal_images.txt")
	cunt:=len(strArr)
	for n := 0; n < cunt; n++ {
		if len(strArr[n])==0{
			continue
		}
		p:=strings.Split(strArr[n], "|")
		id,_:=strconv.Atoi(p[0])
		SimbolID,_:=strconv.Atoi(p[1])
		s:=strings.Split(p[2], ",")
		var PhraseID[] int
		for i := 0; i < len(s); i++ {
			if len(s[i])==0{
				continue
			}
			si,_:=strconv.Atoi(s[i])
			PhraseID=append(PhraseID,si)
		}
		ToneID,_:=strconv.Atoi(p[3])
		MoodID,_:=strconv.Atoi(p[4])

		createNewlastVerbalID(id,SimbolID,PhraseID,ToneID,MoodID)
	}
	return

}
//////////////////////////////


/* получить уникальное сочетание в виде int из двух компонентов int
На входе int2 -  виде ID настроения (jn 20 до 26) преобразуется в диапазон от 1 до 7
простым вычитанием int2-19
 */
func GetToneMoodID(int1 int,int2 int)(int){
	// вмето первой 0 (для "обычный") ставим 9 !!!
	if int1==0{
		int1=9
	}
	s:=strconv.Itoa(int1)
	if int2>19 {
		s += strconv.Itoa((int2 - 19))
	}else{
		s += "0"
	}
	ToneMoodID,_:=strconv.Atoi(s)
	return ToneMoodID
}
//////////////////////////////////
// получить тон и настроение из уникального сочетания
func getToneMoodFromImg(img int)(int,int){
	tonmoode:=strconv.Itoa(img)
	var t=0
	ton:=tonmoode[:1]
	if ton=="9"{
		t=0
	}else{
		t,_=strconv.Atoi(ton)
	}
	m,_:=strconv.Atoi(tonmoode[1:])
	return t,m
}
// расшифровка в виде строки
func getToneMoodStrFromID(img int)(string){
	t,m:=getToneMoodFromImg(img)
	out:="Тон: "+getToneStrFromID(t)
	out+=" Настроение: "+getMoodStrFromID(m)

	return out
}

//////////////////////////////////////////////////


///////////////////////////////
func getToneStrFromID(id int)(string){
var ret=""
//0 - обычный, 1 - восклицательный, 2- вопросительный, 3- вялый, 4 - Повышенный
switch id{
case 0: ret="обычный"
case 1: ret="восклицательный"
case 2: ret="вопросительный"
case 3: ret="вялый"
case 4: ret="Повышенный"

}
return ret
}
////////////////////////////////
func getMoodStrFromID(id int)(string){
var ret=""
// из 20-Хорошее    21-Плохое    22-Игровое    23-Учитель    24-Агрессивное   25-Защитное    26-Протест
// отнимаем 19
switch id{
case 0: ret="Нормальное"
case 1: ret="Хорошее"
case 2: ret="Плохое"
case 3: ret="Игровое"
case 4: ret="Учитель"
case 5: ret="Агрессивное"
case 6: ret="Защитное"
case 7: ret="Протест"

}
return ret
}
////////////////////////////////

