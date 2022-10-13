/* образ сочетания действия с Пульта (кроме словестных) для 3-го уровня дерева автоматизмов.
Смысл (осознанную значимость) образ приобретеает только в контексте Дерева Понимания (дерева мент.автоматизмов)

Особенность: т.к. Activity обязательно присутствует в ветке дерева, а могут быть слова без предваряющего действия,
то первым формируется образ бездействия в createInactionImg()
и из этого узла будет формироваться основное количество вербальных веток.

Область моторных действий (для слов схожая область Брока) ActivityFromIdArr=make(map[int]*Activity)
отвечает за смысл распознанных действий с Пульта,
за конструирование собственных последовательностей действий,
за моторное использование действий.
За все ответственная структура - образ осмысленных действий и их сочетаний.
*/

package psychic

import (
	"BOT/lib"
	"strconv"
	"strings"
)
//////////////////////////////////////////////



/////////////////////////////////////////////

func loadActivityInit(){

	loadActivityFromIdArr()

}
////////////////////////////////////////////

/*
для оптимизации поиска по дереву перед узлом Activity идет узел первого символа : var symbolsArr из word_tree.go

 */
type Activity struct {
	ID int
	ActID[] int // массив ID действий с Пульта 
}
var ActivityFromIdArr=make(map[int]*Activity)
//////////////////////////////////////////


// создать образ сочетаний пусковых стимулов
//В случае отсуствия пусковых стимулов создается ID такого отсутсвия, пример такой записи: 2|||0|0| - ID=2
var lastActivityID=0
func createNewlastActivityID(id int,ActID []int)(int,*Activity){
	oldID,oldVal:=checkUnicumActivity(ActID)
	if oldVal!=nil{
		return oldID,oldVal
	}
	if id==0{
		lastActivityID++
		id=lastActivityID
	}else{
		//		newW.ID=id
		if lastActivityID<id{
			lastActivityID=id
		}
	}

	var node Activity
	node.ID = id
	node.ActID = ActID

	ActivityFromIdArr[id]=&node
	return id,&node
}
func checkUnicumActivity(ActID []int)(int,*Activity){
	for id, v := range ActivityFromIdArr {
		if !lib.EqualArrs(ActID,v.ActID) {
			continue
		}
		return id,v
	}

	return 0,nil
}
/////////////////////////////////////////
// создать новый образ сочетаний действий, если такого еще нет
func CreateNewActivityImage(ActID []int)(int,*Activity){
	if ActID==nil{
		return 0,nil
	}

	id,verb:=createNewlastActivityID(0,ActID)

	SaveActivityFromIdArr()

	return id,verb
}

/////////////////////////////////////////

//////////////////// сохранить образы сочетаний пусковых стимулов
//В случае отсуствия пусковых стимулов создается ID такого отсутсвия, пример такой записи: 2|||0|0|
func SaveActivityFromIdArr(){
	var out=""
	for k, v := range ActivityFromIdArr {
		out+=strconv.Itoa(k)+"|"
		for i := 0; i < len(v.ActID); i++ {
			out+=strconv.Itoa(v.ActID[i])+","
		}
		out+="|"
		out+="\r\n"
	}
	lib.WriteFileContent(lib.GetMainPathExeFile()+"/memory_psy/activity_images.txt",out)

}
////////////////////  загрузить образы сочетаний пусковых стимулов
func loadActivityFromIdArr(){
	ActivityFromIdArr=make(map[int]*Activity)
	// сразу создать первым образ бездействий с Пульта - самое частое состояние образов действия
	createInactionImg()
	strArr,_:=lib.ReadLines(lib.GetMainPathExeFile()+"/memory_psy/activity_images.txt")
	cunt:=len(strArr)
	for n := 0; n < cunt; n++ {
		if len(strArr[n])==0{
			continue
		}
		p:=strings.Split(strArr[n], "|")
		id,_:=strconv.Atoi(p[0])
		s:=strings.Split(p[1], ",")
		var ActID[] int
		for i := 0; i < len(s); i++ {
			if len(s[i])==0{
				continue
			}
			si,_:=strconv.Atoi(s[i])
			ActID=append(ActID,si)
		}

		createNewlastActivityID(id,ActID)
	}
	return

}
//////////////////////////////
func createInactionImg(){
	_,CurrentInformationEnvironment.PsyActionImg=createNewlastActivityID(1,[]int{0})
}
///////////////////////////////////////////


