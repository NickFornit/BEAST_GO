/* Образ ментального действия

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
	activateBaseID int // активация настроения
	activateEmotion int // активация эмоции
	/* Разнообразие заготовленных инфо-функций дает больший потенциал
	   разных ментальных действий, по-началу случайных, но оптимизирующихся по эффекту Правила.
	*/
	activateInfoFunc int // вызов инфо функции
	activateConsciousness bool // ментальная активация func consciousness
}

var MentalActionsImagesArr=make(map[int]*MentalActionsImages)

//////////////////////////////////////////

// вызывается из psychic.go
func MentalActionsImagesInit(){
	loadMentalActionsImagesArr()
}


////////////////////////////////////////////////
// создать новое сочетание ответных действий если такого еще нет 
var lastMentalActionsImagesID=0
func CreateNewlastMentalActionsImagesID(id int,activateBaseID int,activateEmotion int,activateInfoFunc int,activateConsciousness bool)(int,*MentalActionsImages){
	oldID,oldVal:=checkUnicumMentalActionsImages(activateBaseID,activateEmotion,activateInfoFunc,activateConsciousness)
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
	node.activateConsciousness=activateConsciousness

	MentalActionsImagesArr[id]=&node

	if doWritingFile { SaveMentalActionsImagesArr() }

	return id,&node
}
func checkUnicumMentalActionsImages(activateBaseID int,activateEmotion int,activateInfoFunc int,activateConsciousness bool)(int,*MentalActionsImages){
	for id, v := range MentalActionsImagesArr {
		if activateBaseID!=v.activateBaseID || 
			activateEmotion!=v.activateEmotion ||
			activateInfoFunc!=v.activateInfoFunc ||
			activateConsciousness!=v.activateConsciousness{
			continue
		}
		return id,v
	}

	return 0,nil
}
/////////////////////////////////////////






//////////////////// сохранить образы сочетаний ответных действий
//В случае отсуствия ответных действий создается ID такого отсутсвия, пример такой записи: 2|||0|0|
// ID|ActID через ,|PhraseID через ,|ToneID|MoodID|
func SaveMentalActionsImagesArr(){
	var out=""
	for k, v := range MentalActionsImagesArr {
		out+=strconv.Itoa(k)+"|"
		out+=strconv.Itoa(v.activateBaseID)+"|"
		out+=strconv.Itoa(v.activateEmotion)+"|"
		out+=strconv.Itoa(v.activateInfoFunc)+"|"
		out+=strconv.FormatBool(v.activateConsciousness)+"|"
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
		activateBaseID,_:=strconv.Atoi(p[3])
		activateEmotion,_:=strconv.Atoi(p[3])
		activateInfoFunc,_:=strconv.Atoi(p[3])
		activateConsciousness,_:=strconv.ParseBool(p[4])

var saveDoWritingFile= doWritingFile; doWritingFile =false
		CreateNewlastMentalActionsImagesID(id,activateBaseID,activateEmotion,activateInfoFunc,activateConsciousness)
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

	if ai.activateConsciousness {
		out+=" ментальная активация consciousness: "+strconv.FormatBool(ai.activateConsciousness)+"<br>"
	}

	return out
}
/////////////////////////////////////////