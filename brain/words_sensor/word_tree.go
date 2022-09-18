/*  дерево символов, слов
вербальная иерархия распознавателей

*/

package word_sensor

import (
	"BOT/lib"
	"strconv"
	"strings"
)

/////////////////////////////////////////////
// Первый уровень дерева - символы. Заглавных не может быть в воспринимаемом
var symbolsArr = []string{" ","а","б","в","г","д","е","ё","ж","з","и","й","к","л","м","н","о","п","р","с","т","у","ф","х","ц","ч","ш","щ","ъ","ы","ь","э","ю","я","a","b","c","d","e","f","g","h","i","j","k","l","m","n","o","p","q","r","s","t","u","v","w","x","y","z","!","?","@","#","$","%","^","&","*","(",")","+","=","-","1","2","3","4","5","6","7","8","9","0","[","]","{","}","<",">",".",",","/"}
var symbolsRuneArr=make([]rune, len(symbolsArr)) // []rune

func GetSynbolFromID(id int)(string){
	if len(symbolsArr)<id{
		return ""
	}
	return symbolsArr[id]
}

// для получения ID symbolsArr по букве (по string символа)
var SymboIDfromSymbl=make(map[string]int)
// для получения ID symbolsArr по rune символа)
var SymboIDfromRune=make(map[rune]int)

//////////////////////////////////////////////
func IsSymbol(w string)(bool){
//	s:=[]rune(w)
	for i := 0; i < len(symbolsArr); i++ {
		if w == symbolsArr[i]{
			return true
		}
	}
	return false
}
func IsАlphabeticSymbol(w string)(bool){
	var alphabeticEnd=60 // до начала неалфавитных символов
	for i := 1; i < alphabeticEnd; i++ {
		if w == symbolsArr[i]{
			return true
		}
	}
	return false
}
// word_sensor.GetSymbolID(s string)
func GetSymbolID(s string)(int){
	/*
	for i := 0; i < len(symbolsArr); i++ {
		if s == symbolsArr[i]{
			return i
		}
	}
	return 0
	 */
	return SymboIDfromSymbl[s]
}
// для получения ID symbolsArr по rune символа  не работает
func GetSymbolIDfromRune(r rune)(int){
/*
	for i := 0; i < len(symbolsArr); i++ {
		sr:=rune(symbolsArr[i][0])
		if  sr== r{
			return i
		}
	}
	return 0
	*/
	return SymboIDfromRune[r]
}

/////////////////////////////////////////




////////////////////////////////////////////////////////////////////

// дерево слов, разбитых на символы
/*у алфавитных узлов (каждый узел - 1 символ) проявился недостаток:
в дереве слов невозможно вычленить реальные слова,
например, слово "приветствую" перекрывает веткой слово "привет".
Но в дереве фраз все слова имеют ID узла дерева слов, так что список старых (сохраненных) слов
формируется в виде WordIdFormWord=make(map[string]int) проходом в func getWordIdFormWord()
*/
type WordTree struct {
	ID int // id узла символа
	Symbol string // один символ
	Children []WordTree // дочерние узлы (ветвление) НЕ АДРЕСА, А РЕАЛЬНЫЕ ОБЪЕКТЫ
	ParentID int     // ID родителя
	ParentNode *WordTree  // адрес родителя
}

var VernikeWordTree WordTree // дерево слов
var WordTreeFromID=make(map[int]*WordTree) // узел дерева от его ID
var WordTreeFromStr=make(map[string][]*WordTree)// массив узлов с такой SymbolID
//var WordFromID=make(map[int]string) лучше не пытаться получать это при загрузке дерева, используем GetWordFromWordID

// по слову найти его ID - быстрая проверки уже имеющихся слов
var WordIdFormWord=make(map[string]int)
/*
/////////////// для обеспечения уникальности узлов:
type WordUnicum struct {
	ID int
	word string
}
var WordUnicumIdStr=make(map[WordUnicum]int)// для каждого сочетания  выдается ID узла
*/

// счктчик при создании узлов
var lastIDwordTree=0
///////////////////////////////////////////////////////////////////////

// подошла очередь инициализации
func afterLoadTempArr(){
	loadWordTree()
	/*
	   //%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
	   	SetNewWordTreeNode("и") // а я нехочу такое дело
	   //%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
	   	SaveWordTree()
	*/
	//SetNewWordTreeNode("привет")
	//str:=GetWordFromWordID(144)
	//if len(str)>0{}

//	updateWordTreeFromTempArr(2,2)


afterLoadPhraseArr()
	
}

// создать первый, символьный уровень дерева
func initWordTree(vt *WordTree){
	for i := 0; i < len(symbolsArr); i++ {
		createNewNodeWordTree(vt,0,symbolsArr[i])
		SymboIDfromSymbl[symbolsArr[i]]=i
		sr:=[]rune(symbolsArr[i])
		symbolsRuneArr[i]=sr[0]
		SymboIDfromRune[sr[0]]=i

		if symbolsArr[i]=="п"{
			symbolsArr[i]="п"

		}
	}
	//updateWordTreeFromID()
	return
}
/////////////////////////////////

func createNewNodeWordTree(parent *WordTree,id int,word string)(*WordTree){
	if parent==nil{
		return nil
	}
	
	if word==""{
		return nil
	}

//	notAllowScanInThisTime=true // запрет показа карты при обновлении
	if id==0{
		lastIDwordTree++;
		id=lastIDwordTree
	}else{
		//		newW.ID=id
		if lastIDwordTree<id{
			lastIDwordTree=id
		}
	}
/*   лишнее
// не позволять дублировать узлы
	unic:=WordUnicumIdStr[WordUnicum{id,word}]
	if unic!=0 {// уже есть такой узел
		old:=WordTreeFromID[unic] // WordTreeFromID еще актуальный для всех уже созданных, раз не создано нового
		return old //nil
	}
*/
	var newW WordTree
	newW.ID = id
	newW.ParentID = parent.ID
	newW.ParentNode = parent
	newW.Symbol = word

	parent.Children = append(parent.Children, newW)
	// четко находим новый вставленный член (а не &parent.Children[count-1])
	var new *WordTree
	for i := 0; i < len(parent.Children); i++ {
		if parent.Children[i].ID == newW.ID {
			new = &parent.Children[i]
		}
	}

	//WordTreeFromID[new.ID]=new
	// т.к. append меняет длину массива, перетусовывая адреса, то нужно:
	updateWordTreeFromID(parent)// здесь потому, что при загрузке из файла нужно на лету получать адреса
//	notAllowScanInThisTime=false  НЕТ, иначе оно прерывает более общий запрет и возникают гонки

	return new
}
// корректируем адреса всех узлов
func updateWordTreeFromID(parent *WordTree){
	updatingWordTreeFromID(parent)
}
// проход всего дерева
func updatingWordTreeFromID(wt *WordTree){
	if wt.ID>0 {
		wt.ParentNode=WordTreeFromID[wt.ParentID] // wt.ParentNode адрес меняется из=за corretsParent(,
		WordTreeFromID[wt.ID] = wt
		WordTreeFromStr[wt.Symbol] = append(WordTreeFromStr[wt.Symbol],wt)
		//WordUnicumIdStr[WordUnicum{wt.ParentID,wt.Symbol}] = wt.ID
	}
	if wt.Children == nil{// конец ветки
		return
	}
	for i := 0; i < len(wt.Children); i++ {
		updatingWordTreeFromID(&wt.Children[i])
	}
}
///////////////////////////////////////////////////////////


func loadWordTree(){

	initWordTree(&VernikeWordTree)
	strArr,_:=lib.ReadLines(lib.GetMainPathExeFile()+"/memory_reflex/word_tree.txt")
	cunt:=len(strArr)
	//просто проход по всем строкам файла подряд так что сначала идут дочки, потом - их родители
	for n := 0; n < cunt; n++ {
		if len(strArr[n])<2{
			panic("Сбой загрузки дерева слов: ["+strconv.Itoa(n) + "] " + strArr[n])
			return
		}
		p:=strings.Split(strArr[n], "|#|")
		word:=p[1]
		idP:=strings.Split(p[0], "|")
		id,_:=strconv.Atoi(idP[0])
		parentID,_:=strconv.Atoi(idP[1])
		// новый узел с каждой строкой из файла
		createNewNodeWordTree(WordTreeFromID[parentID],id,word)
	}

	return
}
////////////////////////////////////////////
func SaveWordTree(){
	var out=""
	// записываем, начиная со второго уровня (уровень символов не пишем)
	cnt:=len(VernikeWordTree.Children)
	for n := 0; n < cnt; n++ {
		out+=getWtreeNode(&VernikeWordTree.Children[n])
	}
	lib.WriteFileContent(lib.GetMainPathExeFile()+"/memory_reflex/word_tree.txt",out)
	return
}
func getWtreeNode(wt *WordTree)(string){
	var out=""
	if wt.ParentID>0 {
		out += strconv.Itoa(wt.ID) + "|"
		out += strconv.Itoa(wt.ParentID) + "|#|"
		out += wt.Symbol + "\r\n"
	}
	if(wt.Children==nil){// конец
		return out
	}
	for n := 0; n < len(wt.Children); n++ {
		out+=getWtreeNode(&wt.Children[n])
	}
	return out
}
//////////////////////////////////////////////////////////////////////////






// вставить новое слово в дерево:
//найти подходящий узел и если еще нет - вставить новый.
func SetNewWordTreeNode(word string)(int){
		if len(word) == 0 {
			return 0
		}
	notAllowScanInThisTime=true // запрет показа карты при обновлении

		// распознаем что можем
	if len(word)>0 {
		WordDetection(word)

	// Запись недостающего в дерево слов происходит в WordDetection
	}


	notAllowScanInThisTime=false
	return DetectedUnicumID
}
///////////////////
// создание ветки символов, начиная с заданного узла
func createWordTreeNodes(word []rune,wt *WordTree)(int){
	ost:=word[1:]
	if len(ost)==0 {
		return wt.ID
	}

	node:=createNewNodeWordTree(WordTreeFromID[wt.ID], 0, string(ost[0]))
	id:=createWordTreeNodes(ost, WordTreeFromID[node.ID])

	return id
}
///////////////////////////////////////



/////////////////////////////////////////////////
// получить слово из ID слова (из ID конечного узла дерева слов)
func GetWordFromWordID(lastID int)(string){

	var idArr []string
	for {
		node:=WordTreeFromID[lastID]
		if node==nil || lastID==0{
			break
		}
		idArr=append(idArr,node.Symbol)
		lastID=node.ParentID
	}

	var str=""
	for i := len(idArr)-1; i >=0; i-- {
		str+=idArr[i]
	}

	return str
}
////////////////////////////////////////////







//////////////////////////////////////////////////////////////////
/* выдать массив wordsArr[]int из фразы (а не абзаца или текста!)
использовать ТОЛЬКО ДЛЯ func PhraseSeparator !!!
 */
func GetWordIDfromPhrase(phrase string)([]int){
	var out []int
	/*  Делим фразу на слова (в строке нет других разделительных символов,
	т.к. они уже сработали при разделении на фразы).
	*/
	wArr := strings.Split(phrase, " ")
	for n := 0; n < len(wArr); n++ { // перебор отдельных слов
		curWord := strings.TrimSpace(wArr[n])

	wID:= SetNewWordTreeNode(curWord)
		out=append(out,wID)
	}
	return out
}


/* для старых слов получить WordIdFormWord - для распознавания неточно введенных слов и т.п.
Проход всего дерева фраз - там выделены известные слова.
Запускачется при нициализации Дерева фраз,
а при вставки новых слов в Дерево слов - сразу заполняется WordIdFormWord
 */
func getWordIdFormWord(){
	for _, ph := range PhraseTreeFromID {
		word := GetWordFromWordID(ph.WordID)
		WordIdFormWord[word]=ph.WordID
	}
	return
}
///////////////////////////////////////

/* Cлово проверяется на наличие в списке старых (до сохранения) слов WordIdFormWord=make(map[string]int)
и если оно есть, то возвращается его ID.
Если слова нет и оно имеет более 4-х символов, то делается предположение об описке внутренних символов
(в природном распознавателе слово узнается если точно совпали первая и последняя буквы,
а внутренние буквы могут быть как угодно перемешаны)
Если слово распознается, то возвращается ID слова.
*/
func tryWordRecognize(word string)(int){
	id:=WordIdFormWord[word]
	if id != 0{
		return id
	}
	/////////////////////////
	id = getAlternative(word)
	if id != 0 {
		return id
	}

	return 0
}
/////////////////////////////////////////////////////////////////////