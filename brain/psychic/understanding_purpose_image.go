/* образ желаемойцели для Дерева понимания (дерева ментальных автоматизмов)
Сначала - на основе жизненной цели PurposeGenetic и опасности, но потом может задаваться произвольно
с переактивацией дерева понимания.
*/

package psychic

import (
	"BOT/lib"
	"strconv"
	"strings"
)

///////////////////////////////


/* образ цели действия

*/
type PurposeImage struct {
	ID int
	veryActual bool // true - цель очень актуальна
	targetID []int //массив ID парамктров гомеостаза как цели для улучшения в данных условиях
	actionID int //выбранный образ действия для достижения данной цели (ActionsImage)
	// для каждого actionID сила действий сначала принимается =5, а потом корректируется мозжечковыми рефлексами
}
var PurposeImageFromID=make(map[int]*PurposeImage)
///////////////////////////////////////////////////////////////////


// создать новый образ желаемой цели, если такого еще нет
var lastPurposeImagePurposeID=0
func createPurposeImageID(id int,veryActual bool,targetID []int,actionID int)(int,*PurposeImage){
	oldID,oldVal:=checkUnicumPurposeImage(veryActual,targetID,actionID)
	if oldVal!=nil{
		return oldID,oldVal
	}
	if id==0{
		lastPurposeImagePurposeID++
		id=lastPurposeImagePurposeID
	}else{
		//		newW.ID=id
		if lastPurposeImagePurposeID<id{
			lastPurposeImagePurposeID=id
		}
	}

	var node PurposeImage
	node.ID = id
	node.veryActual=veryActual
	node.targetID = targetID
	node.actionID=actionID

	PurposeImageFromID[id]=&node

	if doWritingFile {SavePurposeImageFromIdArr() }

	return id,&node
}
func checkUnicumPurposeImage(veryActual bool,targetID []int,actionID int)(int,*PurposeImage){
	for id, v := range PurposeImageFromID {
		if !lib.EqualArrs(targetID,v.targetID) || veryActual!=v.veryActual || actionID != v.actionID{
			continue
		}
		return id,v
	}
	return 0,nil
}
/////////////////////////////////////////



/////////////////////////////////////////
// сохранить образы
func SavePurposeImageFromIdArr(){
	var out=""
	for k, v := range PurposeImageFromID {
		out+=strconv.Itoa(k)+"|"
		out+=strconv.FormatBool(v.veryActual)+"|"
		for i := 0; i < len(v.targetID); i++ {
			out+=strconv.Itoa(v.targetID[i])+","
		}
		out+="|"
		out+=strconv.Itoa(v.actionID)
		out+="\r\n"
	}
	lib.WriteFileContent(lib.GetMainPathExeFile()+"/memory_psy/purpose_images.txt",out)

}
////////////////////  загрузить образы 
func loadPurposeImageFromIdArr(){
	PurposeImageFromID=make(map[int]*PurposeImage)
	strArr,_:=lib.ReadLines(lib.GetMainPathExeFile()+"/memory_psy/purpose_images.txt")
	cunt:=len(strArr)
	for n := 0; n < cunt; n++ {
		if len(strArr[n])==0{
			continue
		}
		p:=strings.Split(strArr[n], "|")
		id,_:=strconv.Atoi(p[0])
		veryActual,_:=strconv.ParseBool(p[1])
		s:=strings.Split(p[2], ",")
		var targetID[] int
		for i := 0; i < len(s); i++ {
			if len(s[i])==0{
				continue
			}
			si,_:=strconv.Atoi(s[i])
			targetID=append(targetID,si)
		}
		actionID,_:=strconv.Atoi(p[3])

var saveDoWritingFile= doWritingFile; doWritingFile =false
		createPurposeImageID(id,veryActual,targetID,actionID)
doWritingFile =saveDoWritingFile
	}
	return

}
//////////////////////////////