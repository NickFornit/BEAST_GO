/* обобщение примитивных правил на основе эпизодической памяти

В детстве опыт ответов на то, чего пока не знаешь набирается или пробно
или отзеркаливаются чужие ответы. Это становится шаблоном ответа в данной ситуации.
Шаблон усложняется после ответа на ответ и растет цепочка понимания как можно отвечать.
каждый может вспомнить, как учился отвечать на колкости.
Если тебе сказали - "ты дурак", и раньше никогда так не было, очень важно как другие детки на такое отвечали,
ты просто делашь точно так же, отвечаешь "Сам дурак". А тебе: "От дурака слышу!",
ты опять в ступоре, но постепенно набираются цепочки: на такою предъяву - такой-то ответ.
И, как в обучении игры в шахматы развиваются последовательности действий от исходной комбинации.

Вся детская лексика - практически только такие цепочки.

Я очень ясно помню как в детстве искал ответы на значимые реплики,
без чего оказываелся в проигрыше в ловесных перепалках.
Так однажды придумал "мне до лампочки", в другой раз "ты тупой как автобус" -
это были вполне удачные эксперименты, на которые оппонетн затыкался или начинал корчить рожи, т.к. нет ответа.
Но это - уже процесс творчества...


Заполняется при активации дерева (один эпизод)
и при обобщении эпизодической памяти (последовательность эпизодов).
*/


package psychic

import (
	"BOT/lib"
	"strconv"
	"strings"
)

///////////////////////////////////

/* Правила примитивного опыта, обобщающие стимулы->ответы->эффект для таких цепочек в диалогах
На основе этих правил становятся возможны более системные обобщения.
 */
type rules struct {
	ID int
	BranchID int// правило оптимизировано для ветки дерева понимания
	TAid []int // цепочка стимул-ответов ID TriggerAndAction - последовательность из эпизодов памяти подряд, сохраняющая последовательность общения ( дурак -> сам дурак!, маме скажу -> ябеда, щас в морду дам -> ну попробуй)
}


var rulesArr=make(map[int]*rules)

//////////////////////////////////////////

// вызывается из psychic.go
func rulesInit(){
	loadrulesArr()
}


////////////////////////////////////////////////
// создать новое сочетание ответных действий если такого еще нет с ЗАПОМИНАНИЕМ
func CreateNewrulesImage(BranchID int,TAid []int)(int,*rules){
	oldID,oldVal:=checkUnicumrules(BranchID,TAid)
	if oldVal!=nil{
		return oldID,oldVal
	}
	aImgID,aImg:=createNewlastrulesID(0,BranchID,TAid)

	SaveRulesArr()

	return aImgID,aImg
}
/////////////////////////////////////////
// создать образ сочетаний ответных действий
var lastrulesID=0
func createNewlastrulesID(id int,BranchID int,TAid []int)(int,*rules){

	if id==0{
		lastrulesID++
		id=lastrulesID
	}else{
		//		newW.ID=id
		if lastrulesID<id{
			lastrulesID=id
		}
	}

	var node rules
	node.ID = id
	node.BranchID = BranchID
	node.TAid=TAid

	rulesArr[id]=&node
	return id,&node
}
func checkUnicumrules(BranchID int,TAid []int)(int,*rules){
	for id, v := range rulesArr {
		if BranchID != v.BranchID {
			continue
		}
		if !lib.EqualArrs(TAid,v.TAid) {
			continue
		}
		return id,v
	}

	return 0,nil
}
/////////////////////////////////////////






//////////////////// сохранить Образы стимула (действий оператора) - ответа Beast
//В случае отсуствия ответных действий создается ID такого отсутсвия, пример такой записи: 2|||0|0|
// ID|ActID через ,|PhraseID через ,|ToneID|MoodID|
func SaveRulesArr(){
	var out=""
	for k, v := range rulesArr {
		out+=strconv.Itoa(k)+"|"
		out+=strconv.Itoa(v.BranchID)+"|"
		for i := 0; i < len(v.TAid); i++ {
			out+=strconv.Itoa(v.TAid[i])+","
		}
		out+="\r\n"
	}
	lib.WriteFileContent(lib.GetMainPathExeFile()+"/memory_psy/rules.txt",out)

}
////////////////////  загрузить образы стимула (действий оператора) - ответа Beast
// ID|ActID через ,|PhraseID через ,|ToneID|MoodID|
func loadrulesArr(){
	rulesArr=make(map[int]*rules)
	strArr,_:=lib.ReadLines(lib.GetMainPathExeFile()+"/memory_psy/rules.txt")
	cunt:=len(strArr)
	for n := 0; n < cunt; n++ {
		if len(strArr[n])==0{
			continue
		}
		p:=strings.Split(strArr[n], "|")
		id,_:=strconv.Atoi(p[0])

		BranchID,_:=strconv.Atoi(p[1])
		s:=strings.Split(p[2], ",")
		var TAid[] int
		for i := 0; i < len(s); i++ {
			if len(s[i])==0{
				continue
			}
			si,_:=strconv.Atoi(s[i])
			TAid=append(TAid,si)
		}

		createNewlastrulesID(id,BranchID,TAid)
	}
	return

}
///////////////////////////////////////////