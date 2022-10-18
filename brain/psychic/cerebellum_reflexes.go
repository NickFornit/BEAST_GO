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

/*  по результатам выполнения автоматизма выбираются дополнительные действия
или изменяется сила действия автоматизма.
Это - средство не переписывать автоматизм, а оптимизировать его.
В качестве дополнительных действий используются имеющиеся автоматизмы на основе которых создаются
мозжечковые рефлексы
 */
type cerebellumReflex struct {
	id int
	typeAut int 									// тип корректруемого автоматизма: 0 - это ID моторного, 1 - ID ментального
	sourceAutomatizmID int 				// корректируемый моторный или ментальный автоматизм
	addEnergy int 								// добавление (-убавление) силы действия Energy. Может быть отрицательное число, чтобы уменьшить энергию автоматизма
	// кроме корректирования самого автоматизма по силе действия, могут быть запущены дополнительные автоматизмы:
	additionalAutomatizmID []int 	// массив ID дополнительных моторных автоматизмов
	additionalMentalAutID []int 	// массив ID дополнительных ментальных автоматизмов
}
// общий массив рефлексов мозжечка
var cerebellumReflexFromID = make(map[int]*cerebellumReflex)
// рефлексы мозжечка по моторному ID
var cerebellumReflexFromMotorsID = make(map[int]*cerebellumReflex)
// рефлексы мозжечка по ментальному ID
var cerebellumReflexFromMentalsID = make(map[int]*cerebellumReflex)

func cerebellumReflexInit() {
	loadCerebellumReflex()
}

// последний ID рефлексов мохжечка
var lastCRid = 0

// создать новый автоматизм
// В случае отсуствия автоматизма создается ID такого отсутсвия, пример такой записи: 2|||0|0| - ID=2
func createNewCerebellumReflex(id int, typeAut int, sourceAutomatizmID int)(int, *cerebellumReflex){
	oldID, oldVal := checkUnicumCerebellumReflex(typeAut, sourceAutomatizmID)
	if oldVal != nil { return oldID, oldVal }
	if id == 0 {
		lastCRid++
		id = lastCRid
	} else {
		if lastCRid < id {
			lastCRid = id
		}
	}

	var node cerebellumReflex
	node.id = id
	node.typeAut = typeAut
	node.sourceAutomatizmID = sourceAutomatizmID
	node.addEnergy = 5 // сразу придаст максимум, сложившись с энергией автоматизма :)
	// node.additionalAutomatizmID = additionalAutomatizmID
	cerebellumReflexFromID[id] = &node
	if typeAut == 0 {
		cerebellumReflexFromMotorsID[sourceAutomatizmID] = &node
	} else {
		cerebellumReflexFromMentalsID[sourceAutomatizmID] = &node
	}
	SaveCerebellumReflex()
	return id, &node
}

// поиск рефлекса мозжечка
func checkUnicumCerebellumReflex(typeAut int, sourceAutomatizmID int)(int, *cerebellumReflex) {
	for id, v := range cerebellumReflexFromID {
		if typeAut == v.typeAut && sourceAutomatizmID == v.sourceAutomatizmID { return id, v }
	}
	return 0, nil
}

// Сохранить рефлексы мозжечка
// структура записи: id|typeAut|sourceAutomatizmID||addEnergy|additionalAutomatizmID|additionalMentalAutID
// В случае отсуствия пусковых стимулов создается ID такого отсутсвия, пример такой записи: 2|||0|0|
func SaveCerebellumReflex() {
	var out = ""
	for k, v := range cerebellumReflexFromID {
		out += strconv.Itoa(k)+"|"
		out += strconv.Itoa(v.typeAut)+"|"
		out += strconv.Itoa(v.sourceAutomatizmID)+"|"
		out += strconv.Itoa(v.addEnergy)+"|"
		for i := 0; i < len(v.additionalAutomatizmID); i++ {
			out += strconv.Itoa(v.additionalAutomatizmID[i]) + ","
		}
		out += "|"
		for i := 0; i < len(v.additionalMentalAutID); i++ {
			out += strconv.Itoa(v.additionalMentalAutID[i]) + ","
		}
		out += "\r\n"
	}
	lib.WriteFileContent(lib.GetMainPathExeFile() + "/memory_psy/cerebellum_reflex.txt", out)
}

// Загрузить рефлексы мозжечка
// структура записи: id|typeAut|sourceAutomatizmID||addEnergy|additionalAutomatizmID|additionalMentalAutID
func loadCerebellumReflex() {
	cerebellumReflexFromID = make(map[int]*cerebellumReflex)
	strArr, _ := lib.ReadLines(lib.GetMainPathExeFile() + "/memory_psy/cerebellum_reflex.txt")
	for n := 0; n < len(strArr); n++ {
		p := strings.Split(strArr[n], "|")
		if len(p) < 5 { return }
		id, _ := strconv.Atoi(p[0])
		typeAut, _ := strconv.Atoi(p[1])
		sourceAutomatizmID, _ := strconv.Atoi(p[2])
		addEnergy, _ := strconv.Atoi(p[3])
		a := strings.Split(p[4], "|")
		var additionalAutomatizmID []int
		for i := 0; i < len(a); i++ {
			aid, _ := strconv.Atoi(a[i])
			additionalAutomatizmID = append(additionalAutomatizmID, aid)
		}
		a = strings.Split(p[5], "|")
		var additionalMentalAutID []int
		for i := 0; i < len(a); i++ {
			aid, _ := strconv.Atoi(a[i])
			additionalMentalAutID = append(additionalMentalAutID,aid)
		}
		_, ca := createNewCerebellumReflex(id, typeAut, sourceAutomatizmID)
		ca.addEnergy = addEnergy
		ca.additionalAutomatizmID = additionalAutomatizmID
		ca.additionalMentalAutID = additionalMentalAutID
	}
	return
}

// вернуть скорректированную силу действия
func getCerebellumReflexAddEnergy(kind int, automatizmID int) int {
	var e *cerebellumReflex

	if kind == 0 {
		e = cerebellumReflexFromMotorsID[automatizmID]
	} else {
		e = cerebellumReflexFromMentalsID[automatizmID]
	}

	if e == nil { return 0 }
	return e.addEnergy
}

// выполнить дополнительные мозжечковые автоматизмы сразу после выполняющегося автоматизма
func runCerebellumAdditionalAutomatizm(kind int, automatizmID int) {
	var cr *cerebellumReflex

	if kind == 0 {
		cr = cerebellumReflexFromMotorsID[automatizmID]
	} else {
		cr = cerebellumReflexFromMentalsID[automatizmID]
	}
	if cr == nil { return	}
	if kind == 0 {
		aArr := cr.additionalAutomatizmID
		for i := 0; i < len(aArr); i++ {
			RumAutomatizmID(aArr[i])
		}
	} else {
		aArr := cr.additionalMentalAutID
		for i := 0; i < len(aArr); i++ {
			RunMentalMentalAutomatizmsID(aArr[i])
		}
	}
}