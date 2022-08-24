/* Рефлексы мозжечка для корректировки автоматизмов 

Допонение автоматизма другими корректирующими действиеями 
или корректировка самого автоматизма.
*/

package psychic

import (
	"BOT/lib"
	"strconv"
	"strings"
)

/////////////////////////////////////////////////

/*  по результатам выполнения автоматизма выбираются дополнительные действи
или изменяется сила действия автоматизма.
Это - средство не переписывать автоматизм, а оптимизировать его.
В качестве дополнительных действий используются имеющиеся автоматизмы на основе которых сохдаются
автоматизмы мозжечковые, с Belief==3, наследующие ID ветки
Или создаются совершенно новые автоматизмы с Belief==3, и НЕ наследующие ID ветки с BranchID==0

 */
type cerebellumReflex struct {
	id int
	typeAut int // тип корректруемого автоматизма: 0 - это ID моторного, 1 - ID ментального
	sourceAutomatizmID int // корректируемый моторный автоматизм
	addEnergy int // добавление (-убавление) силы действия Energy
	additionalAutomatizmID []int // массив ID дополнительных моторных автоматизмов
	additionalMentalAutID []int // массив ID дополнительных ментальных автоматизмов
}
var cerebellumReflexFromID=make(map[int]*cerebellumReflex)

var cerebellumReflexFromAutomatizmID=make(map[int]*cerebellumReflex)
//////////////////////////////////////////////////

func cerebellumReflexInit(){
	loadCerebellumReflex()
}

////////////////////////////////////////////////////////////////////////////
// создать новый автоматизм
//В случае отсуствия автоматизма создается ID такого отсутсвия, пример такой записи: 2|||0|0| - ID=2
var lastCRid=0
func createNewCerebellumReflex(id int,typeAut int,sourceAutomatizmID int)(int,*cerebellumReflex){
	oldID,oldVal:=checkUnicumCerebellumReflex(typeAut,sourceAutomatizmID)
	if oldVal!=nil{
		return oldID,oldVal
	}
	if id==0{
		lastCRid++
		id=lastCRid
	}else{
		//		newW.ID=id
		if lastCRid<id{
			lastCRid=id
		}
	}

	var node cerebellumReflex
	node.id = id
	node.typeAut = typeAut
	node.sourceAutomatizmID = sourceAutomatizmID
	node.addEnergy = 5
	//node.additionalAutomatizmID = additionalAutomatizmID
	cerebellumReflexFromID[id]=&node
	cerebellumReflexFromAutomatizmID[sourceAutomatizmID]=&node
	SaveCerebellumReflex()
	return id,&node
}
/////////////////////////////////////////
func checkUnicumCerebellumReflex(typeAut int,sourceAutomatizmID int)(int,*cerebellumReflex){
	for id, v := range cerebellumReflexFromID {
		if typeAut==v.typeAut && sourceAutomatizmID==v.sourceAutomatizmID{
			return id,v
		}
	}
	return 0,nil
}
/////////////////////////////////////////

// СОХРАНИТЬ структура записи: id|BranchID|Usefulness||Sequence||NextID|Energy|Belief
//В случае отсуствия пусковых стимулов создается ID такого отсутсвия, пример такой записи: 2|||0|0|
func SaveCerebellumReflex(){
	var out=""
	for k, v := range cerebellumReflexFromID {
		out+=strconv.Itoa(k)+"|"
		out+=strconv.Itoa(v.typeAut)+"|"
		out+=strconv.Itoa(v.sourceAutomatizmID)+"|"
		out+=strconv.Itoa(v.addEnergy)+"|"
		for i := 0; i < len(v.additionalAutomatizmID); i++ {
			out+=strconv.Itoa(v.additionalAutomatizmID[i])+","
		}
		out+="|"
		for i := 0; i < len(v.additionalMentalAutID); i++ {
			out+=strconv.Itoa(v.additionalMentalAutID[i])+","
		}
		out+="\r\n"
	}
	lib.WriteFileContent(lib.GetMainPathExeFile()+"/memory_psy/cerebellum_reflex.txt",out)
}
// ЗАГРУЗИТЬ структура записи: id|BranchID|Usefulness||Sequence||NextID|Energy|Belief
func loadCerebellumReflex(){
	cerebellumReflexFromID=make(map[int]*cerebellumReflex)
	strArr,_:=lib.ReadLines(lib.GetMainPathExeFile()+"/memory_psy/cerebellum_reflex.txt")
	for n := 0; n < len(strArr); n++ {
		p := strings.Split(strArr[n], "|")
		id, _ := strconv.Atoi(p[0])
		typeAut, _ := strconv.Atoi(p[1])
		sourceAutomatizmID, _ := strconv.Atoi(p[2])
		addEnergy, _ := strconv.Atoi(p[3])
		a := strings.Split(p[4], "|")
		var additionalAutomatizmID []int
		for i := 0; i < len(a); i++ {
			aid, _ := strconv.Atoi(a[i])
			additionalAutomatizmID=append(additionalAutomatizmID,aid)
		}
		a = strings.Split(p[5], "|")
		var additionalMentalAutID []int
		for i := 0; i < len(a); i++ {
			aid, _ := strconv.Atoi(a[i])
			additionalMentalAutID=append(additionalMentalAutID,aid)
		}
		_,ca:=createNewCerebellumReflex(id, typeAut, sourceAutomatizmID)
		ca.addEnergy=addEnergy
		ca.additionalAutomatizmID=additionalAutomatizmID
		ca.additionalMentalAutID=additionalMentalAutID
	}
	return
}
///////////////////////////////////////////



// вернуть скорректированную силу действия
func getCerebellumReflexAddEnergy(automatizmID int)(int){
	e:=cerebellumReflexFromID[automatizmID]
	if e==nil{
		return 0
	}
	return e.addEnergy
}
//////////////////////////////////////////////////////////



/*  выполнить текущий мозжечковый рефлекс сразу после выполняющегося автоматизма

 */
func runCerebellumReflex(automatizmID int){
	cr:=cerebellumReflexFromID[automatizmID]
	if cr==nil{
		return
	}

	aArr:=cr.additionalAutomatizmID
	for i := 0; i < len(aArr); i++ {
		RumAutomatizmID(aArr[i])
	}

}
/////////////////////////////////////////////////////