/* функции для распознаваемых слов

*/


package word_sensor

import (
	"BOT/lib"
	"strconv"
	"strings"
)

///////////////////////////////////////////////////

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


/* удалить слово с ID из дерева слов
Пройти от данного ID по его родителям до ID предыдущего слова.
Так, при удалении "приветствую" должны удаляться узлы чтобы осталось "привет".
После этого удаляются это слово из Дерева Фраз.

!После удаления ряда слов нужно перезагрузить сервер чтобы обновились данные!
 */
func DeleteWord(wID int){
	node:=WordTreeFromID[wID]
	if node==nil{// уже удалена
		return
	}
	wordID:=wID // сохраняем для удаления в Дереве фраз
// проходим ветку дерева
	for node.ParentNode !=nil{
// это или вложенное слово или ветвление на другое слово
		if existsWordID(node.ParentNode.ID) || len(node.ParentNode.Children)>1 {// начинается уже другое, вложенное слово
			// удаляем узел прямо из файла word_tree.txt
			deleteNodeFromFile(node.ID)
			// удалить слово все всех упоминаниях в Дереве фраз
			deleteWordFromPhrase(wordID)
			return
		}
		next:=node.ParentNode
		// удаляем узел прямо из файла word_tree.txt
		deleteNodeFromFile(node.ID)

		node=next
	}
	return
}
// есть ли слово с данным ID
func existsWordID(ID int)(bool){
	for _, k := range WordIdFormWord {
		if k==ID{
			return true
		}
	}
	return false
}
// удалить строку из word_tree.txt
func deleteNodeFromFile(wID int){
	strArr,_:=lib.ReadLines(lib.GetMainPathExeFile()+"/memory_reflex/word_tree.txt")
var out=""
	for n := 0; n < len(strArr); n++ {
		p:=strings.Split(strArr[n], "|")
		id,_:=strconv.Atoi(p[0])
		if id==wID{
			continue
		}
		out+=strArr[n]+"\r\n"
	}
	lib.WriteFileContent(lib.GetMainPathExeFile()+"/memory_reflex/word_tree.txt",out)
}
///////////////////////////////////////////////////////////////////////