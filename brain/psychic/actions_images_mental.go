/* Образ ментального действия
ID|activateBaseID|activateEmotion|activatePurpose|activateInfoFunc|activateMotorID|
*/

package psychic

import (
	"BOT/lib"
	"strconv"
	"strings"
)

///////////////////////////////////
/* виды ментальных действий для MentalActionsImages
Можно наращивать при необходимости новые действия, добавляя в:
func GetMentalActionsString
func RunMentalAutomatizm
 */
var MentalActionsType=[]int{
	1,//активация настроения Mood в дереве понимания 1 2 3 (отражает -1 0 1 UnderstandingNode.Mood)
	2,//активация эмоции EmotionID в дереве понимания (UnderstandingNode.EmotionID)
	3,//активация PurposeImage в дереве понимания (UnderstandingNode.PurposeID)
	4,//запуск инфо-функции
	5,//запуск моторного автоматизма
	6,//запуск Доминанты
	7,// создание новой Доминанты
	8,// xxxxx
}



type MentalActionsImages struct {
	ID    int   // идентификатор данного сочетания пусковых стимулов
	typeID int // вид действия - MentalActionsType
	valID int // ID для действия
}
var MentalActionsImagesArr=make(map[int]*MentalActionsImages)

//////////////////////////////////////////

// вызывается из psychic.go
func MentalActionsImagesInit(){
	loadMentalActionsImagesArr()
}


////////////////////////////////////////////////
/* создать новое сочетание ответных действий если такого еще нет
!!! activateBaseID в виде 1,2,3, а не -1,0,1 !!!!
при activateBaseID==0 не задается переадресация
 */
var lastMentalActionsImagesID=0
func CreateNewlastMentalActionsImagesID(id int,typeID int,valID int,CheckUnicum bool)(int,*MentalActionsImages){
	if CheckUnicum {
		oldID,oldVal:=checkUnicumMentalActionsImages(typeID,valID)
		if oldVal!=nil{
			return oldID,oldVal
		}
	}

	if id==0{
		lastMentalActionsImagesID++
		id=lastMentalActionsImagesID
	}else{
		//		newW.ID=id
		if lastMentalActionsImagesID<id{
			lastMentalActionsImagesID=id
		}
	}

	var node MentalActionsImages
	node.ID = id
	node.typeID=typeID
	node.valID=valID

	MentalActionsImagesArr[id]=&node

	if doWritingFile { SaveMentalActionsImagesArr() }

	return id,&node
}
func checkUnicumMentalActionsImages(typeID int,valID int)(int,*MentalActionsImages){
	for id, v := range MentalActionsImagesArr {
		if typeID!=v.typeID ||
			valID!=v.valID  {
			continue
		}
		return id,v
	}

	return 0,nil
}
/////////////////////////////////////////






//////////////////// сохранить образы сочетаний ментальных действий
// ID|activateBaseID|activateEmotion|activatePurpose|activateInfoFunc|activateMotorID|
func SaveMentalActionsImagesArr(){
	var out=""
	for k, v := range MentalActionsImagesArr {
		out+=strconv.Itoa(k)+"|"
		out+=strconv.Itoa(v.typeID)+"|"
		out+=strconv.Itoa(v.valID)
		out+="\r\n"
	}
	lib.WriteFileContent(lib.GetMainPathExeFile()+"/memory_psy/action_images_mental.txt",out)

}
////////////////////  загрузить образы сочетаний ментальных действий
func loadMentalActionsImagesArr(){
	MentalActionsImagesArr=make(map[int]*MentalActionsImages)
	strArr,_:=lib.ReadLines(lib.GetMainPathExeFile()+"/memory_psy/action_images_mental.txt")
	cunt:=len(strArr)
	for n := 0; n < cunt; n++ {
		if len(strArr[n])==0{
			continue
		}
		p:=strings.Split(strArr[n], "|")
		id,_:=strconv.Atoi(p[0])
		typeID,_:=strconv.Atoi(p[1])
		valID,_:=strconv.Atoi(p[2])

var saveDoWritingFile= doWritingFile; doWritingFile =false
		CreateNewlastMentalActionsImagesID(id,typeID,valID,false)
doWritingFile =saveDoWritingFile
	}
	return

}
///////////////////////////////////////////

//  для Пульта
func GetMentalActionsString(act int)(string){
	if act==0{
		return "Нулевой ID ментального действия."
	}
	ai:=MentalActionsImagesArr[act]
	if ai==nil{
		return "Несуществующий ID ментального действия: "+strconv.Itoa(act)
	}
	switch ai.typeID{
	case 1: return "Активация настроения Mood <span style='color:blue;cursor:pointer' onClick='show_mental_actions("+strconv.Itoa(ai.ID)+")'>"+strconv.Itoa(ai.ID)+ "</span>в дереве понимания"
	case 2: return "Активация эмоции  <span style='color:blue;cursor:pointer' onClick='show_mental_actions("+strconv.Itoa(ai.ID)+")'>"+strconv.Itoa(ai.ID)+ "</span>в дереве понимания"
	case 3: return "Активация цели PurposeImage <span style='color:blue;cursor:pointer' onClick='show_mental_actions("+strconv.Itoa(ai.ID)+")'>"+strconv.Itoa(ai.ID)+ "</span>в дереве понимания"
	case 4: return "Запуск инфо-функции № <span style='color:blue;cursor:pointer' onClick='show_mental_actions("+strconv.Itoa(ai.ID)+")'>"+strconv.Itoa(ai.ID)+"</span>"
	case 5: return "Запуск моторного автоматизма № <span style='color:blue;cursor:pointer' onClick='show_mental_actions("+strconv.Itoa(ai.ID)+")'>"+strconv.Itoa(ai.ID)+"</span>"
	case 6: return "Запуск Доминанты <span style='color:blue;cursor:pointer' onClick='show_mental_actions("+strconv.Itoa(ai.ID)+")'>"+strconv.Itoa(ai.ID)+"</span>"
	case 7: return "Создание новой Доминанты"
	case 8: return "xxxxxxx"
	}

	return "Нет такого типа действия с ID="+strconv.Itoa(ai.typeID)
}
/////////////////////////////////////////
// раскрыть значение ai.valID стврокой для Пульта
func GetMentalActionInfo(actID int)string{
	ai:=MentalActionsImagesArr[actID]
	if ai==nil{
		return "Несуществующий ID ментального действия: "+strconv.Itoa(actID)
	}
	switch ai.typeID { // еще есть getSituationDetaileString(
	case 1: return getMoodString(ai.valID)
	case 2: return getEmotionString(ai.valID)
	case 3: return getPurposeDetaileString(ai.valID)
	case 4: return getMentalFunctionString(ai.valID)
	case 5: return GetAutomotizmIDString(ai.valID)
	case 6: return GetDominantaIDString(ai.valID)
	case 7: // нет значения

	}
	return "Несуществующий ID ментального действия: "+strconv.Itoa(actID)
}
///////////////////////////////////////////