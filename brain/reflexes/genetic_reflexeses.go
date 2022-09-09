/*
безусловные рефлексы
*/

package reflexes

import (
	"BOT/lib"
	"strconv"
	"strings"
)

func init() {
	loadGeneticReflexes()
	loadImagesArrs()
	initReflexTree()
}

type GeneticReflex struct {
	ID          int
	lev1        int
	lev2        []int
	lev3        []int
	ActionIDarr []int
	// Result int - у безусловных рефлексов нет конкуренции, кроме того, что они подавляются более высокоуровневыми рефлексами и автоматизмами
}
var GeneticReflexes = make(map[int]*GeneticReflex)

// для быстрого поиска по совпадениям строк
type geneticReflexStr struct {
ID int
lev1 string
lev2 string
lev3 string
lev4 string
actions string
}

var geneticReflexesStr = make(map[int]*geneticReflexStr)
var lastGeneticReflexID = 0

func CreateNewGeneticReflex(id int, lev1 int, lev2 []int, lev3 []int, ActionIDarr []int) (int, *GeneticReflex) {
	// посмотреть, если рефлекс с такими же условиями уже есть
	idOld, rOld := compareUnicum(lev1, lev2, lev3)
	if idOld > 0 {
		return idOld, rOld
	}

	if id == 0 {
		lastGeneticReflexID++
		id = lastGeneticReflexID
	} else {
		//		newW.ID=id
		if lastGeneticReflexID < id {
			lastGeneticReflexID = id
		}
	}

	var newW GeneticReflex
	newW.ID = id
	newW.lev1 = lev1
	newW.lev2 = lev2
	newW.lev3 = lev3
	newW.ActionIDarr = ActionIDarr
	GeneticReflexes[id] = &newW
	return id, &newW
}

// посмотреть, если рефлекс с такими же условиями уже есть
func compareUnicum(lev1 int, lev2 []int, lev3 []int) (int, *GeneticReflex) {
	for k, v := range GeneticReflexes {
		if v.lev1 == lev1 && lib.EqualArrs(v.lev2, lev2) && lib.EqualArrs(v.lev3, lev3) {
			return k, v
		}
	}
	return 0, nil
}

// P.S. безусловные рефлексы создаются в редакторе и поэтому здесь нет функции их сохранения.
// а только загрузка имеющихся в формате ID|lev1|lev2_1,lev2_2,...|lev3_1,lev3_2,...|actin_1,actin_2,...:

func loadGeneticReflexes() {
	path := lib.GetMainPathExeFile()
	lines, _ := lib.ReadLines(path + "/memory_reflex/dnk_reflexes.txt")
	for i := 0; i < len(lines); i++ {
		if len(lines[i]) < 4 {
			continue
		}
		p := strings.Split(lines[i], "|")
		id, _ := strconv.Atoi(p[0])
		lev1, _ := strconv.Atoi(p[1])
		// второй уровень
		pn := strings.Split(p[2], ",")
		var lev2 []int
		for i := 0; i < len(pn); i++ {
			b, _ := strconv.Atoi(pn[i])
			if b > 0 {
				lev2 = append(lev2, b)
			}
		}
		// НЕ СОЗДАВАТЬ ЗАРАНЕЕ createNewBaseStyle(0, lev2)
		// третий уровень
		pn = strings.Split(p[3], ",")
		var lev3 []int
		for i := 0; i < len(pn); i++ {
			b, _ := strconv.Atoi(pn[i])
			if b > 0 {
				lev3 = append(lev3, b)
			}
		}
		// создать образ сочетаний пусковых стимулов
		//actID,_:=createNewlastTriggerStimulsID(0,lev3,nil,0,0)
		pn = strings.Split(p[4], ",")
		var ActionIDarr []int
		for i := 0; i < len(pn); i++ {
			b, _ := strconv.Atoi(pn[i])
			if b > 0 {
				ActionIDarr = append(ActionIDarr, b)
			}
		}
		CreateNewGeneticReflex(id, lev1, lev2, lev3, ActionIDarr)
		var newS geneticReflexStr
		newS.ID=id
		newS.lev1=p[1]
		newS.lev2=p[2]
		newS.lev3=p[3]
		newS.actions=p[4]
		geneticReflexesStr[id]=&newS
	}
	return
}

/* Сохранить в файл безусловные рефлексы */
func SaveGeneticReflexes() {
	var out string

	// сохранение только в режиме Larva
	if EvolushnStage > 0 {
		return
	}
	for i := 1; i < len(GeneticReflexes)+1; i++ {
		out += ListDnkReflex(i) + "\r\n"
	}
	lib.WriteFileContent(lib.GetMainPathExeFile()+"/memory_reflex/dnk_reflexes.txt", out)
}

// Получить строку ДНК-рефлекса по ID
func ListDnkReflex(ID int)string{
	return strconv.Itoa(GeneticReflexes[ID].ID) + "|" +
		strconv.Itoa(GeneticReflexes[ID].lev1) + "|" +
		strings.Join(lib.StrArrToIntArr(GeneticReflexes[ID].lev2), ",") + "|" +
		strings.Join(lib.StrArrToIntArr(GeneticReflexes[ID].lev3), ",") + "|" +
		strings.Join(lib.StrArrToIntArr(GeneticReflexes[ID].ActionIDarr), ",")
}