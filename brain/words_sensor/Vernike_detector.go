/*  распознаватель слов и фраз по типу зоны Вернике в мозге.
Память о воспринятых фразах в текущем активном контексте (Vernike_detector.go): var MemoryDetectedArr []MemoryDetected

Распознавание фраз начинается в main.go с word_sensor.VerbalDetection(text_dlg, is_input_rejim, moodID)
С ПУльта приходит текст, который в VerbalDetection() разбирается на фрзацы (\r\n):

абзацы в PhraseSeparator() разбираются на фразы по разделителям (знаки препинания)
фразы в WordDetection() разбиваются на слова.
Распознанные (и нераспознанные) последовательности сохраняются в оперативной памяти Beast MemoryDetectedArr.
где распознанный текст представлен в виде уникального laslID фразы

ОПИСКИ при вводе слова. Если слово не распознается и оно имеет более 3-х символов,
то делается предположение об описке внутренних символов
(в природном распознавателе слово узнается если точно совпали первая и последняя буквы,
а внутренние буквы могут быть как угодно перемешаны)
Если слово распознается, то подставляется ID слова.

Нераспознанной фразы НЕ БЫВАЕТ т.к. она тут же создается

Тон фразы можно задать 1) с помощью знаков ! и ? в конце фразы
или задать преимущественно - выбрав Тон под окном ввода фразы.
*/

package word_sensor

import (
	"BOT/brain/gomeostas"
	_ "BOT/lib"
	"regexp"
	"strconv"
	"strings"
)

// запрет показа карт WordTreeFromID и PhraseTreeFromID во время распознавания и записи
// против паники типа "одновременная запись и считывание карты"
var notAllowScanInThisTime = false

// подошла очередь инициализации
func afetrInitPhraseTree() {
	wordRecognizerInit() //
	// VerbalDetection("привет новая абзаца",1) // текст с пульта
	// PhraseSeparator("привет") // распознавание фразы
	// WordDetection("привет") // распознавание слова
	isReadyWordSensorLevel = 1
	initWordPult()
	initPrasePult()
}

// индикация, что дерево загружено, можно вводить тексты
var isReadyWordSensorLevel = 0

func IsReadyWordSensorLevel() bool {
	if isReadyWordSensorLevel > 0 { return true	} // связь с корнями проекта
	return false
}

// word_sensor.VerbalDetectin("активностный")
// для использования в SetNewWordTreeNode

// здесь набирается массив lastID распознанных фраз
var CurrentPhrasesIDarr []int
// тон сообщения с Пульта при передаче фразы
var CurPultTone = 0
// настроение с Пульта при передаче фразы
var CurPultMood = 0
// текущий тон фразы: 0-обычный, 1-восклицательный, 2-вопросительный
var DetectedTone = 0

/* Память о воспринятых фразах в текущем активном контексте:
В принципе это не нужно, т.к. имеется эпизодическая память
(правда только для Правил, но это как раз и всязывает фразы с их полезностью)
 */
type MemoryDetected struct {
	//распознанный текст в виде lastID выделенных фраз
	PhrasesID []int // массив структур распознанных фраз
	Tone      int   // Тон: 0-обычный, 1-восклицательный, 2-вопросительный, 3-вялый, 4-Повышенный
	Mood      int   // настроение при передаче фразы с Пульта: 20-Хорошее 21-Плохое 22-Игровое 23-Учитель 24-Агрессивное 25-Защитное 26-Протест
	// индекс==ID активного базового контекста, значение - вес этого контекста
	ActveContextWeight map[int]int
}
// массив памяти накапливается в течении дня, обрабатывается и очищается во сне
var MemoryDetectedArr []MemoryDetected

// добавить строку в массив памяти о воспринятых фразах в текущем активном контексте
func addNewMemoryDetected() {
	var newM MemoryDetected
	newM.PhrasesID = CurrentPhrasesIDarr
	// тон может указываться 1) в виде ! или ? во фразе - DetectedTone 	И/ИЛИ 2) в виде радиокнопки Тон с Пульта - CurPultTone
	var tone = 0
	if DetectedTone > 0 { // преимущество - у задатчика тона 2)
		tone = DetectedTone
	} else { // есть ! или ? во фразе
		tone = CurPultTone
	}
	newM.Tone = tone
	newM.Mood = CurPultMood // настроение с Пульта: повышенный нормальный вялый Хорошее Плохое Игровое Учитель Агрессивное Защитное Протест
	newM.ActveContextWeight = gomeostas.GetActiveContextInfo()
	MemoryDetectedArr = append(MemoryDetectedArr, newM)
}

/*  распознавание фразы с Пульта - бывает только в нижнем регистре
 */
var wlev = 0
var pultOut = ""

//// если обычный режим диалога (на ПУльте не стоит галка "набивка рабочих фраз без отсеивания мусорных слов ")
var NoCheckWordCount = false // is_input_rejim - набивка рабочих фраз с отсеиванием мусорных слов

// вызывается фразой с Пульта
func VerbalDetection(text string, isDialog int, toneID int, moodID int) string {
	notAllowScanInThisTime = true // запрет показа карты при обновлении
	NoCheckWordCount = false
	CurrentPhrasesIDarr = nil
	if isDialog == 0 { // это набивка рабочих фраз без отсеивания мусорных слов
		// игнорировать getWordTemparrCount и всегда распознавать слова
		NoCheckWordCount = true
	}
	CurPultTone = toneID
	CurPultMood = moodID

	pultOut = ""
	// стандартно разделить текст на короткие фразы, отправить на накопление
	// разделяем на фразы
	strArr := strings.Split(text, "|#") // а не |#| - чтобы оставлять разделитель "|"
	for i := 0; i < len(strArr); i++ {
		if i > 0 { pultOut += "<br>" }
		str := addNewtempArr(strArr[i])
		for n := 0; n < len(str); n++ {
			if n > 0 { pultOut += " " }
			// проход фразы с распознаванием
			pultOut += PhraseSeparator(str[n])
			if DetectedUnicumPhraseID > 0 { // распознанная фраза
				CurrentPhrasesIDarr = append(CurrentPhrasesIDarr, DetectedUnicumPhraseID)
			} else {
				// нераспознанной фразы НЕ БЫВАЕТ т.к. она тут же создается
			}
		}
	}
	// добавить в стек памяти распознанных
	addNewMemoryDetected()
	// reflexes.ActiveFromPhrase() // активировать дерево рефлексов фразой - только для условных рефлексов
	// ответ на Пульт:
	notAllowScanInThisTime = false
	return pultOut
}

// проход одной фразы (т.е. по разделителям в предложении, а не пр \r\n)
func PhraseSeparator(text string) string {
	var pultOut = ""
	// чистим лишние пробелы
	rp := regexp.MustCompile("s+")
	text = rp.ReplaceAllString(text, " ")
	text = strings.TrimSpace(text)

	wordsArr := GetWordIDfromPhrase(text) // распознаватель слов
	str := PhraseDetection(wordsArr) // распознаватель фразы
	pultOut += str + "(" + strconv.Itoa(DetectedUnicumPhraseID) + ")"

	// тон сообщения
	DetectedTone = 0
	if strings.Contains(text, "!") {
		DetectedTone = 1
	}
	if strings.Contains(text, "?") {
		DetectedTone = 2
	}

	return pultOut
}

/* получить последователньость wordID из уникального идентификатора фразы CurrentPhrasesIDarr[i]
начиная с любого узла дерева Фраз (не обязательно конечного!) - к первому узлу ветки
 */
func GetWordArrFromPhraseID(PhraseNodeID int) []int {
	var wArr []int
	// пройти фразу от последнего слова до первого
	w:=PhraseTreeFromID[PhraseNodeID]
	if w==nil{
		return wArr
	}
	wArr=append(wArr,w.ID)
	for w.ParentID >0{
		w=PhraseTreeFromID[w.ParentID]
		if w!=nil {
			wArr = append(wArr, w.ID)
		}
	}
	// восстановить порядок слов
	var wordIDarr []int
	for i := len(wArr)-1; i >=0 ; i-- {
		wordIDarr=append(wordIDarr,wArr[i])
	}
	return wordIDarr
}
////////////////////////////////////////






