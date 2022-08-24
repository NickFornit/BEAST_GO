/*
Распознаватель слов по символьно для заполнения дерева слов word_tree.go
и распознавания слов при вводе с пульта

*/

package word_sensor

import "strings"

/////////////
func wordRecognizerInit() {
	//	WordDetection("играть") // распознавание слова
}

////////////////////////////////////////
// текущий уникальный ID последней активной ветки дерева - результат детекции фразы - для дальнейшего использования
var DetectedUnicumID = 0

// при активации дерева слов
var FirstSimbolID = 0

// нераспознанный остаток
var CurrentVerbalEnd []rune

var detectedCurrentID = 0 // текущий найденный ID последней активной ветки дерева слов
var currentStepCount = 0  // текущеее число нераспознанных символов

//var lastFoundID=0 // последний ID символа при проходе дерева
//////////////////////////////////////////////

/////////////////////////////////////////////////////////
// проход одного слова - распознавание слова
// возвращает найденное ID слова или похожей альтернативы
func WordDetection(text string) int {
	text = strings.TrimSpace(text)
	if len(text) == 0 {
		return 0
	}
	CurrentVerbalEnd = []rune("")
	DetectedUnicumID = 0
	//var pultOut=""
	detectedCurrentID = 0
	currentStepCount = 0

	r := []rune(text)
	// основа дерева
	cnt := len(VernikeWordTree.Children)
	var curFirstLevelID = 0
	for n := 0; n < cnt; n++ {
		smblNode := VernikeWordTree.Children[n]
		rt := []rune(smblNode.Symbol)
		if r[0] == rt[0] {
			FirstSimbolID = GetSymbolIDfromRune(r[0])
			curFirstLevelID = smblNode.ID
			if len(r) == 1 { // это - символ, присвоить слову ID символа
				DetectedUnicumID = VernikeWordTree.Children[n].ID
				CurrentVerbalEnd = []rune("")
				return DetectedUnicumID
			}

			cldrn := VernikeWordTree.Children[n].Children
			cnt := len(cldrn)
			for n := 0; n < cnt; n++ {
				getWordTreeNode(r[1:], &cldrn[n])
			}
		}
	}
	//////////////// результат распознавания
	if detectedCurrentID > 0 {
		if currentStepCount == 0 { // полностью распознан
			DetectedUnicumID = detectedCurrentID
		} else {
			var nr = len(r) - currentStepCount
			CurrentVerbalEnd = r[nr:]
		}
	}
	/////////////////////////////////
	var needSave = false
	if DetectedUnicumID == 0 {

		// если обычный режим диалога (на ПУльте не стоит галка "набивка работчих фраз без отсеивания мусорных слов ")
		if !NoCheckWordCount {
			//отсеивать мусорных (редких - менее 4 повторов в tempArr) слов
			repeet := getWordTemparrCount(text)
			if repeet < 4 {
				// попробовать найти подходящее слово
				DetectedUnicumID = getAlternative(text)
				if DetectedUnicumID != 0 {
					return DetectedUnicumID
				}
			}
		}

		//  нераспознанный остаток
		if len(CurrentVerbalEnd) > 0 {
			r := CurrentVerbalEnd
			var tree *WordTree
			if detectedCurrentID > 0 {
				tree = WordTreeFromID[detectedCurrentID]
			} else {
				tree = &VernikeWordTree
			}
			// просто добавить новую ветку - из диалога это стоит делать за 1 раз т.к. слова уже известны
			node := createNewNodeWordTree(tree, 0, string(r[0]))
			tree = node
			id := createWordTreeNodes(r, WordTreeFromID[tree.ID])
			DetectedUnicumID = id
			//SaveWordTree()
			needSave = true
		}
	}

	// нет вообще такого, добавить все слово
	if DetectedUnicumID == 0 {
		tree := WordTreeFromID[curFirstLevelID]
		r = r[1:]
		// сразу создать первый узел
		if len(r) > 0 {
			node := createNewNodeWordTree(tree, 0, string(r[0]))
			tree = node
			if tree != nil {
				id := createWordTreeNodes(r, WordTreeFromID[tree.ID])
				DetectedUnicumID = id
				//SaveWordTree()
				needSave = true
			}
		}
	}
	if needSave {
		SaveWordTree()
	}

	return DetectedUnicumID //pultOut+"{"+strconv.Itoa(DetectedUnicumID)+")"
}

/////////////////////////////////////////////

// cканирование следует строго по нужной ветке
func getWordTreeNode(word []rune, wt *WordTree) {
	if len(word) == 0 {
		return
	}

	ost := word[1:]
	if string(word[0]) != wt.Symbol { // пошло не туда
		return
	}

	detectedCurrentID = wt.ID
	currentStepCount = len(ost)

	for n := 0; n < len(wt.Children); n++ {
		getWordTreeNode(ost, &wt.Children[n])
	}
}

//////////////////////////////////////////////////////////////////////

/* найти слово в tempArr и выдать его повторяемость, если слова нет в tempArr - добавить его.
 слово еще раньше добавляется в tempArr если его там нет
так что оно уже обязательно там будет
*/

func getWordTemparrCount(word string) int {
	for k, v := range tempArr {
		if k == word {
			return v
		}
	}
	return 0 // на всякий случай
}

//////////////////////////////////////////

/* попробовать найти подходящее слово с альтрнативным ID
Первые буквы должны совпадать, а остальные, кроме последней (разные окончания),
быть перемешаны, но в наличии >80%.
Сканирует дерево с начальной буквы строго по числу чимволов слова.
Это имитирует свойство персептронного распознавателя.
*/
var smilarlyArr []int // сбор схожих lastID для анализа
func getAlternative(word string) int {
	r := []rune(word)
	var wordSize = len(r)
	cnt := len(VernikeWordTree.Children)
	for n := 0; n < cnt; n++ {
		if VernikeWordTree.Children[n].Symbol == string(r[0]) {
			alphNode := &VernikeWordTree.Children[n]
			getWordiSmilarly(wordSize, r, alphNode)
			break
		}
	}
	if smilarlyArr == nil { // нет подходящих по размеру и начинающихся тождественно
		return 0
	}
	// смотрим схожесть у найденных
	for n := 0; n < len(smilarlyArr); n++ {
		str := GetWordFromWordID(smilarlyArr[n])
		// для сверки берем руны без первых в словах
		id := chooseSmilarly([]rune(str)[1:], r[1:], smilarlyArr[n])
		if id > 0 { // выбран подходящий аналог
			return id
		}
	}

	return 0
}
func getWordiSmilarly(wordSize int, word []rune, wt *WordTree) {
	if len(word) == 0 { //
		return
	}

	ost := word[1:]
	if wt.Children == nil {
		count := getSmilarlyCount(wt.ID)
		if count == wordSize { // совпадение по числу символов
			smilarlyArr = append(smilarlyArr, wt.ID)
		}
		return
	}
	for n := 0; n < len(wt.Children); n++ {
		getWordiSmilarly(wordSize, ost, &wt.Children[n])
	}
}
func getSmilarlyCount(lastID int) int {
	var count = 0
	for {
		node := WordTreeFromID[lastID]
		if node == nil || lastID == 0 {
			break
		}
		count++
		lastID = node.ParentID
	}
	return count
}

// подходит ли слово str как аналог word (без учета последних символов)
func chooseSmilarly(str []rune, word []rune, id int) int {
	var tCount = 0 // число совпадений
	for n := 0; n < len(word)-1; n++ {
		for m := 0; m < len(str)-1; m++ {
			if str[m] == word[n] {
				tCount++
				break
			}
		}
	}
	if (float64(tCount) / float64(len(word)-1)) > 0.8 { //более 80% совпадений
		return id
	}
	return 0
}

//////////////////////////////////////////////////
