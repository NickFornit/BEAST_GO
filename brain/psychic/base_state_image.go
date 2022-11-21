/* Образ текущего сосотояния для func understandingSituation
Возникает только при активации дерева от ActiveFromConditionChange()
т.е. только если это - не действия оператора Пульта: if !WasOperatorActiveted {

Аналогично образу действий BaseStateImage (оператора или Beast) 
этот образ участвует в формировании Правил.

Ранее использовалось для записи Правила по изменению состояния, а не ответу оператора, теперь НЕ ИСПОЛЬЗУЕТСЯ в Правилах.
*/


package psychic

import (
	"BOT/lib"
	"strconv"
	"strings"
)

/////////////////////////////////////////////////
type BaseStateImage struct {
	ID    int   // идентификатор данного сочетания пусковых стимулов
	Mood int // ощущаемое настроение PsyBaseMood: -1 Плохое настроение, 0 Нормальное, 1 - хорошее настроение
	EmotionID int // эмоция, может произвольно меняться
	SituationID int // ID объекта структуры понимания SituationImage, может произвольно меняться
}

var BaseStateImageArr=make(map[int]*BaseStateImage)

//////////////////////////////////////////

// вызывается из psychic.go
func BaseStateImageInit(){
	loadBaseStateImageArr()
}


////////////////////////////////////////////////
var lastBaseStateImageID=0
func CreateNewStatImageID(id int,Mood int,EmotionID int,SituationID int,CheckUnicum bool)(int,*BaseStateImage){
	if CheckUnicum {
		oldID,oldVal:=checkUnicumBaseStateImage(Mood,EmotionID,SituationID)
		if oldVal!=nil{
			return oldID,oldVal
		}
	}

	if id==0{
		lastBaseStateImageID++
		id=lastBaseStateImageID
	}else{
		//		newW.ID=id
		if lastBaseStateImageID<id{
			lastBaseStateImageID=id
		}
	}

	var node BaseStateImage
	node.ID = id
	node.Mood = Mood
	node.EmotionID=EmotionID
	node.SituationID=SituationID

	BaseStateImageArr[id]=&node

	if doWritingFile { SaveBaseStateImageArr() }

	return id,&node
}
func checkUnicumBaseStateImage(Mood int,EmotionID int,SituationID int)(int,*BaseStateImage){
	for id, v := range BaseStateImageArr {
		if Mood!=v.Mood {
			continue
		}
		if EmotionID!=v.EmotionID {
			continue
		}
		if SituationID!=v.SituationID {
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
func SaveBaseStateImageArr(){
	var out=""
	for k, v := range BaseStateImageArr {
		out+=strconv.Itoa(k)+"|"
		out+=strconv.Itoa(v.Mood)+"|"
		out+=strconv.Itoa(v.EmotionID)+"|"
		out+=strconv.Itoa(v.SituationID)
		out+="\r\n"
	}
	lib.WriteFileContent(lib.GetMainPathExeFile()+"/memory_psy/base_state_images.txt",out)

}
////////////////////  загрузить образы сочетаний ответных действий
// ID|ActID через ,|PhraseID через ,|ToneID|MoodID|
func loadBaseStateImageArr(){
	BaseStateImageArr=make(map[int]*BaseStateImage)
	strArr,_:=lib.ReadLines(lib.GetMainPathExeFile()+"/memory_psy/base_state_images.txt")
	cunt:=len(strArr)
	for n := 0; n < cunt; n++ {
		if len(strArr[n])==0{
			continue
		}
		p:=strings.Split(strArr[n], "|")
		id,_:=strconv.Atoi(p[0])

		Mood,_:=strconv.Atoi(p[1])
		EmotionID,_:=strconv.Atoi(p[2])
		SituationID,_:=strconv.Atoi(p[3])

		var saveDoWritingFile= doWritingFile; doWritingFile =false
		CreateNewStatImageID(id,Mood,EmotionID,SituationID,false)
		doWritingFile =saveDoWritingFile
	}
	return

}
///////////////////////////////////////////

func GetBaseStateImageString(bs int)(string){
	var out="не сделано"


	return out
}
/////////////////////////////////////////
