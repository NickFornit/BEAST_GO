/* Эмоции.
Распознавание активности текущих Базовых контекстов в виде структур 
- осмысленной значимости сочетаний активных базовых контекстов.
Произвольно возможна активация имеющегося образа, создание образа новых сочетаний.
*/


package psychic
import (
	"BOT/brain/gomeostas"
	"BOT/lib"
	"strconv"
	"strings"
)
///////////////////////////////////


//////////////////////////////
func emotionsInit(){
	loadEmotionArr()
}
////////////////////////////////

/* Образ сочетания базовых контекстов и оценки успешности действий.
ПО результатам действий в дерево рефлексов подставляется результирующий Emotion
и возникает новый цикл активации Дерева, уже по внутренним причинам.
Это - простейший механизм саморегуляции по результату действий,
т.к. на уровне ор.рефлекса осмысливается ситуация и могут создаваться новые автоматизмы.
Основы произвольности - корректировка Emotion по PsyBaseMood
 */
type Emotion struct {
	ID int  // идентификатор данного сочетания контекстов
	BaseIDarr[] int // сочетание базовых контекстов
	/* упешность действий. В automatizm_result.go по результатам действий
	1) устанавливается настроение PsyBaseMood (-1,0,1),
	2) в GetCurrentInformationEnvironment возникает Emotion, принимающий значение PsyBaseMood
		а при совершенном дейсвии уже изменились Б.параметры и, значит и основа Emotion.
	4) переактивируется Дерево автоматизмов, уже по внутренним причинам
	5) при активации Дерева или тупо выполняется автоматизм или возникает новый ор.рефлекс.
	 */
	Success int // (-1,0,1)
}
////////////////////////////////

var EmotionFromIdArr=make(map[int]*Emotion)

/*  создать новую эмоцию, если такой еще нет
 */
var lastEmotionID=0
func createNewBaseStyle(id int,Success int,BaseIDarr []int)(int,*Emotion){
	oldID,oldVal:=checkUnicumEmotion(Success,BaseIDarr)
	if oldVal!=nil{
		return oldID,oldVal
	}
	if id==0{
		lastEmotionID++
		id=lastEmotionID
	}else{
		//		newW.ID=id
		if lastEmotionID<id{
			lastEmotionID=id
		}
	}

	var node Emotion
	node.ID = id
	node.BaseIDarr = BaseIDarr

	EmotionFromIdArr[id]=&node

	SaveEmotionArr()

	return id,&node
}
func checkUnicumEmotion(Success int,bArr []int)(int,*Emotion){
	for id, v := range EmotionFromIdArr {
		if Success==v.Success && lib.EqualArrs(bArr,v.BaseIDarr) {
			return id,v
		}
	}

	return 0,nil
}
////////////////////////////////////////


/////////////////  сохранить образы сочетаний базовых стилей
func SaveEmotionArr(){
	var out=""
	for k, v := range EmotionFromIdArr {
		out+=strconv.Itoa(k)+"|"
		for i := 0; i < len(v.BaseIDarr); i++ {
			out+=strconv.Itoa(v.BaseIDarr[i])+","
		}
		out+="|"
		out+=strconv.Itoa(v.Success)+"|"
		out+="\r\n"
	}
	lib.WriteFileContent(lib.GetMainPathExeFile()+"/memory_psy/emotion_images.txt",out)
}
//////////////////  загрузить образы сочетаний базовых стилей
func loadEmotionArr(){
	EmotionFromIdArr=make(map[int]*Emotion)
	strArr,_:=lib.ReadLines(lib.GetMainPathExeFile()+"/memory_psy/emotion_images.txt")
	cunt:=len(strArr)
	for n := 0; n < cunt; n++ {
		if len(strArr[n])==0{
			continue
		}
		p:=strings.Split(strArr[n], "|")
		id,_:=strconv.Atoi(p[0])
		s:=strings.Split(p[1], ",")
		var BaseIDarr[] int
		for i := 0; i < len(s); i++ {
			if len(s[i])==0{
				continue
			}
			bc,_:=strconv.Atoi(s[i])
			BaseIDarr=append(BaseIDarr,bc)
		}
		Success,_:=strconv.Atoi(p[2])
		createNewBaseStyle(id,Success,BaseIDarr)
	}
	return
}
/////////////////////////////////////////////////////////////////////



////////////////////////////////////////////////////////////////////
// Описать словами текущую эмоцию
func getEmotonsComponentStr(em *Emotion)(string){
var out=""
if em==nil{
	return "Еще не определились эмоции."
}
	for i := 0; i < len(em.BaseIDarr); i++ {
		if i>0{out+=", "}
		out+=gomeostas.GetBaseContextCondFromID(em.BaseIDarr[i])
	}
	switch em.Success{
	case -1: out+=" ( было НЕуспешное действие)"
	case 1:	out+=" ( было успешное действие)"
	}

return out
}
//////////////////////////////////////////
