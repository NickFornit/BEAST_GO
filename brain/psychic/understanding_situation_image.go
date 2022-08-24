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
	ReflexTreeNodeID int
	/* Тип ситуации:
	0 - ничего не делали, но нужно осмысление
	1 - было действие автоматизма (смотреть в автоматизме ветки Usefulness int - (БЕС)ПОЛЕЗНОСТЬ: вред: -10 0 +10 +n польза diffPsyBaseMood )
	2 - был автоматический запуск автоматизма без ориентировочного рефлекса.
	3 - был плохой автоматизм, нужно найти лучше
	4 - все спокойно, можно экспериментивароть
	5 - есть важные (по опыту) признаки в новизне NoveltySituation - осмыслить их
	 */
	SituationType int

}
var SituationImageFromIdArr=make(map[int]*SituationImage)
/////////////////////////////////


var lastSituationImageID=0
func createSituationImage(id int,reflexTreeNodeID int,SituationType int)(int,*SituationImage){
	oldID,oldVal:=checkUnicumSituationImage(reflexTreeNodeID,SituationType)
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

	var node SituationImage
	node.ID = id
	node.ReflexTreeNodeID = reflexTreeNodeID
	node.SituationType = SituationType

	SituationImageFromIdArr[id]=&node

	SaveSituationImage()

	return id,&node
}
/////////////////////////////////////
func checkUnicumSituationImage(reflexTreeNodeID int,SituationType int)(int,*SituationImage){
	for id, v := range SituationImageFromIdArr {
		if reflexTreeNodeID!=v.ReflexTreeNodeID && SituationType!=v.SituationType {
			continue
		}
		return id,v
	}

	return 0,nil
}
////////////////////////////////////


func SaveSituationImage(){
	var out=""
	for k, v := range ActivityFromIdArr {
		out+=strconv.Itoa(k)+"|"
		for i := 0; i < len(v.ActID); i++ {
			out+=strconv.Itoa(v.ActID[i])+","
		}
		out+="|"
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
		reflexTreeNodeID,_:=strconv.Atoi(p[1])
		SituationType,_:=strconv.Atoi(p[2])

		createSituationImage(id,reflexTreeNodeID,SituationType)
	}
	return

}
//////////////////////////////