/*
Образы восприятия и образы действия
*/

package reflexes

import (
	"BOT/brain/action_sensor"
	wordSensor "BOT/brain/words_sensor"
	"BOT/lib"
	"strconv"
	"strings"
)

// сохранить образы сочетаний базовых стилей
func SaveBaseStyleArr() {
	var out = ""
	for k, v := range BaseStyleArr {
		out += strconv.Itoa(k) + "|"
		for i := 0; i < len(v.BSarr); i++ {
			out += strconv.Itoa(v.BSarr[i]) + ","
		}
		out += "\r\n"
	}
	lib.WriteFileContent(lib.GetMainPathExeFile() + "/memory_reflex/base_style_images.txt", out)
}

// загрузить образы сочетаний базовых стилей
func loadBaseStyleArr() {
	BaseStyleArr = make(map[int]*BaseStyle)
	strArr, _ := lib.ReadLines(lib.GetMainPathExeFile() + "/memory_reflex/base_style_images.txt")
	cunt := len(strArr)
	for n := 0; n < cunt; n++ {
		if len(strArr[n]) == 0 {
			continue
		}
		p := strings.Split(strArr[n], "|")
		id, _ := strconv.Atoi(p[0])
		s := strings.Split(p[1], ",")
		var BSarr []int
		for i := 0; i < len(s); i++ {
			if len(s[i]) == 0 {
				continue
			}
			si, _ := strconv.Atoi(s[i])
			BSarr = append(BSarr, si)
		}
		createNewBaseStyle(id, BSarr,false)
	}
	return
}

// базовый стиль - Образ сочетаний Базовых Контекстов гомеостаза
type BaseStyle struct {
	ID    int // идентификатор данного сочетания контекстов
	BSarr []int
}

var BaseStyleArr = make(map[int]*BaseStyle)

// создать образ сочетаний активных Базовых контекстов
var lastBaseStyleID = 0

func createNewBaseStyle(id int, BSarr []int,CheckUnicum bool) (int, *BaseStyle) {
	if CheckUnicum {
		oldID, oldVal := checkUnicumBaseStyle(BSarr)
		if oldVal != nil { return oldID, oldVal	}
	}

	if id == 0 {
		lastBaseStyleID++
		id = lastBaseStyleID
	} else {
		if lastBaseStyleID < id {
			lastBaseStyleID = id
		}
	}

	var node BaseStyle
	node.ID = id
	node.BSarr = BSarr
	BaseStyleArr[id] = &node

	return id, &node
}
func checkUnicumBaseStyle(bArr []int) (int, *BaseStyle) {
	for id, v := range BaseStyleArr {
		if lib.EqualArrs(bArr, v.BSarr) {
			return id, v
		}
	}
	return 0, nil
}

/* Образ текущего сочетания пусковых (Trigger) стимулов в восприятии
Сочетание пусковых стимулов, включая фразы и включая тон, настроение - должны быть уникальными
чтобы только для данного сочетания мог срабатывать данный рефлекс.

Это - образ много-модеального воздействия или просто образ восприятия.

Формат: ID|RSarr через ,|PhraseID через ,|ToneID|MoodID|

Идентична ActionsImage психики.
*/
type TriggerStimuls struct {
	ID    int   // идентификатор данного сочетания пусковых стимулов
	RSarr []int // массив действий с Пульта
	// для текущего сообщения с Пусльта:
	PhraseID []int // массив фразID (DetectedUnicumPhraseID) слова каждой фразы вытаскиваются wordSensor.GetWordArrFromPhraseID(PhraseID[0])
	ToneID   int   // тон сообщения с Пульта
	MoodID   int   // настроение оператора
}

// образы сочетаний пусковых стимулов
var TriggerStimulsArr = make(map[int]*TriggerStimuls)

// создать образ сочетаний пусковых стимулов
// В случае отсуствия пусковых стимулов создается ID такого отсутсвия, пример такой записи: 2|||0|0| - ID=2
var lastTriggerStimulsID = 0

// Создать новый образ сочетаний пусковых стимулов
func CreateNewlastTriggerStimulsID(id int, RSarr []int, PhraseID []int, ToneID int, MoodID int,CheckUnicum bool) (int, *TriggerStimuls) {
	if CheckUnicum {
		oldID, oldVal := checkUnicumTriggerStimuls(RSarr, PhraseID, ToneID, MoodID)
		if oldVal != nil {
			return oldID, oldVal
		}
	}

	if id == 0 {
		lastTriggerStimulsID++
		id = lastTriggerStimulsID
	} else {
		if lastTriggerStimulsID < id {
			lastTriggerStimulsID = id
		}
	}

	var node TriggerStimuls
	node.ID = id
	node.RSarr = RSarr
	node.PhraseID = PhraseID
	node.ToneID = ToneID
	node.MoodID = MoodID

	TriggerStimulsArr[id] = &node
	return id, &node
}

// проверка наличия образа сочетаний пусковых стимулов
func checkUnicumTriggerStimuls(bArr []int, PhraseID []int, ToneID int, MoodID int) (int, *TriggerStimuls) {
	for id, v := range TriggerStimulsArr {
		if !lib.EqualArrs(bArr, v.RSarr) { continue	}
		if !lib.EqualArrs(PhraseID, v.PhraseID) {	continue }
		if ToneID != v.ToneID || MoodID != v.MoodID {	continue }
		return id, v
	}

	return 0, nil
}

// создать новое сочетание пусковых стимулов если такого еще нет
func CreateNewTriggerStimulsImage() {
	PhraseID := wordSensor.CurrentPhrasesIDarr
	ToneID := wordSensor.DetectedTone
	MoodID := wordSensor.CurPultMood
	RSarr := action_sensor.CheckCurActions()

	setOldActiveCurTriggerStimulsVal(ActiveCurTriggerStimulsID)
	ActiveCurTriggerStimulsID, _ = CreateNewlastTriggerStimulsID(0, RSarr, PhraseID, ToneID, MoodID,true)

	// SaveTriggerStimulsArr()
}

// сохранить образы сочетаний пусковых стимулов
// В случае отсуствия пусковых стимулов создается ID такого отсутсвия, пример такой записи: 2|||0|0|
// ID|RSarr через ,|PhraseID через ,|ToneID|MoodID|
func SaveTriggerStimulsArr() {
	var out = ""
	for k, v := range TriggerStimulsArr {
		out += strconv.Itoa(k) + "|"
		for i := 0; i < len(v.RSarr); i++ {
			out += strconv.Itoa(v.RSarr[i]) + ","
		}
		out += "|"
		for i := 0; i < len(v.PhraseID); i++ {
			out += strconv.Itoa(v.PhraseID[i]) + ","
		}
		out += "|"
		out += strconv.Itoa(v.ToneID) + "|"
		out += strconv.Itoa(v.MoodID) + "|"
		out += "\r\n"
	}
	lib.WriteFileContent(lib.GetMainPathExeFile()+"/memory_reflex/trigger_stimuls_images.txt", out)
}

// загрузить образы сочетаний пусковых стимулов
// ID|RSarr через ,|PhraseID через ,|ToneID|MoodID|
func loadTriggerStimulsArr() {
	TriggerStimulsArr = make(map[int]*TriggerStimuls)
	strArr, _ := lib.ReadLines(lib.GetMainPathExeFile() + "/memory_reflex/trigger_stimuls_images.txt")
	cunt := len(strArr)
	for n := 0; n < cunt; n++ {
		if len(strArr[n]) == 0 {
			continue
		}
		p := strings.Split(strArr[n], "|")
		id, _ := strconv.Atoi(p[0])

		s := strings.Split(p[1], ",")
		var RSarr []int
		for i := 0; i < len(s); i++ {
			if len(s[i]) == 0 { continue }
			si, _ := strconv.Atoi(s[i])
			RSarr = append(RSarr, si)
		}

		s = strings.Split(p[2], ",")
		var PhraseID []int
		for i := 0; i < len(s); i++ {
			if len(s[i]) == 0 {	continue }
			si, _ := strconv.Atoi(s[i])
			PhraseID = append(PhraseID, si)
		}
		x, _ := strconv.Atoi(p[3])
		ToneID := x
		x, _ = strconv.Atoi(p[4])
		MoodID := x

		CreateNewlastTriggerStimulsID(id, RSarr, PhraseID, ToneID, MoodID,false)
	}
	return
}