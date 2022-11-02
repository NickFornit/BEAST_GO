/* Образ ментального действия
ID|activateBaseID|activateEmotion|activateInfoFunc|activateMotorID|
*/

package psychic

import (
	"BOT/lib"
	"strconv"
	"strings"
)

///////////////////////////////////
type MentalActionsImages struct {
	ID    int   // идентификатор данного сочетания пусковых стимулов
	activateBaseID int // активация настроения Mood в дереве понимания 1 2 3 (отражает -1 0 1 UnderstandingNode.Mood)
	activateEmotion int // активация эмоции EmotionID в дереве понимания (UnderstandingNode.EmotionID)
	/* Разнообразие заготовленных инфо-функций дает больший потенциал
	   разных ментальных действий, поначалу случайных, но оптимизирующихся по эффекту Правила.
	*/
	activateInfoFunc int // вызов инфо функции
	activateMotorID int // запуск моторного автоматизма по результатам инфо-функции создания автоматизма
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
func CreateNewlastMentalActionsImagesID(id int,activateBaseID int,activateEmotion int,
	activateInfoFunc int,activateMotorID int)(int,*MentalActionsImages){

	oldID,oldVal:=checkUnicumMentalActionsImages(activateBaseID,activateEmotion,activateInfoFunc,activateMotorID)
	if oldVal!=nil{
		return oldID,oldVal
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
	node.activateBaseID = activateBaseID
	node.activateEmotion=activateEmotion
	node.activateInfoFunc=activateInfoFunc
	node.activateMotorID=activateMotorID

	MentalActionsImagesArr[id]=&node

	if doWritingFile { SaveMentalActionsImagesArr() }

	return id,&node
}
func checkUnicumMentalActionsImages(activateBaseID int,activateEmotion int,activateInfoFunc int,activateMotorID int)(int,*MentalActionsImages){
	for id, v := range MentalActionsImagesArr {
		if activateBaseID!=v.activateBaseID || 
			activateEmotion!=v.activateEmotion ||
			activateInfoFunc!=v.activateInfoFunc ||
			activateMotorID!=v.activateMotorID {
			continue
		}
		return id,v
	}

	return 0,nil
}
/////////////////////////////////////////






//////////////////// сохранить образы сочетаний ответных действий
// ID|activateBaseID|activateEmotion|activateInfoFunc|activateMotorID|
func SaveMentalActionsImagesArr(){
	var out=""
	for k, v := range MentalActionsImagesArr {
		out+=strconv.Itoa(k)+"|"
		out+=strconv.Itoa(v.activateBaseID)+"|"
		out+=strconv.Itoa(v.activateEmotion)+"|"
		out+=strconv.Itoa(v.activateInfoFunc)+"|"
		out+=strconv.Itoa(v.activateMotorID)
		out+="\r\n"
	}
	lib.WriteFileContent(lib.GetMainPathExeFile()+"/memory_psy/action_images_mental.txt",out)

}
////////////////////  загрузить образы сочетаний ответных действий
// ID|ActID через ,|PhraseID через ,|ToneID|MoodID|
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
		activateBaseID,_:=strconv.Atoi(p[1])
		activateEmotion,_:=strconv.Atoi(p[2])
		activateInfoFunc,_:=strconv.Atoi(p[3])
		activateMotorID,_:=strconv.Atoi(p[4])

var saveDoWritingFile= doWritingFile; doWritingFile =false
		CreateNewlastMentalActionsImagesID(id,activateBaseID,activateEmotion,activateInfoFunc,activateMotorID)
doWritingFile =saveDoWritingFile
	}
	return

}
///////////////////////////////////////////

func GetMentalActionsString(act int)(string){
	if MentalActionsImagesArr == nil || len(MentalActionsImagesArr)==0 || MentalActionsImagesArr[act]==nil{
		return ""
	}
	var out=""
	ai:=MentalActionsImagesArr[act]
	if ai.activateBaseID != 0 {
		out+=" активация настроения: "+getToneStrFromID(ai.activateBaseID)+" "
	}

	if ai.activateEmotion != 0 {
		out+=" активация эмоции: "+getToneStrFromID(ai.activateEmotion)+" "
	}

	if ai.activateInfoFunc != 0 {
		out+=" вызов инфо функции: "+getToneStrFromID(ai.activateInfoFunc)+" "
	}

	if ai.activateMotorID != 0 {
		out+=" запуск автоматизма: "+getToneStrFromID(ai.activateInfoFunc)+" "
	}

	return out
}
/////////////////////////////////////////