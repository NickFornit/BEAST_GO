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

Цепочки Правил в Эпиз.памяти создабт карту решений в контексте одной темы:
карты местности - куда идти после очередного шага,
карту игры в шахматы: как ходить в данной позиции и на сколько шагов вперед обдумывать решения.

Поиск по Павилам ведется не по отдельным действиям (а это - не только слова, но и тон и/или кнопки действий),
а тупо по ID образов действий. Но, если крепко заматься и вспоминать детально,
то можно делать поиск по всем компонентам действий и делать его мягким.
Заготовлен прототип getSuitableMentalRulesСarefully() для использования в ментальных инфо-функциях.
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
	TAid []int // цепочка стимул-ответов ID TriggerAndAction - последовательность из эпизодов памяти подряд, сохраняющая последовательность общения ( дурак -> сам дурак!, маме скажу -> ябеда, щас в морду дам -> ну попробуй)
}
var rulesArr=make(map[int]*rules)

//////////////////////////////////////////


///////////////////////////////////////////
// вызывается из psychic.go
func rulesInit(){
	loadrulesArr()

//	getCur10lastRules()
}


////////////////////////////////////////////////
// создать новое сочетание ответных действий если такого еще нет
var lastrulesID=0
var isNotLoading=true
func createNewlastrulesID(id int,TAid []int)(int,*rules){

	if TAid == nil{
		return 0,nil
	}
	oldID,oldVal:=checkUnicumrules(TAid)
	if oldVal!=nil{
		return oldID,oldVal
	}
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
	node.TAid=TAid

	rulesArr[id]=&node

	if doWritingFile{
		SaveRulesArr()
	}
	if isNotLoading {
		if len(TAid) > 1 {
			lib.WritePultConsol("<span style='color:green'>Записано групповое <b>ПРАВИЛО № " + strconv.Itoa(id) + "</b></span>")
		} else {
			lib.WritePultConsol("<span style='color:green'>Записано <b>ПРАВИЛО № " + strconv.Itoa(id) + "</b></span>")
		}
	}

	return id,&node
}
func checkUnicumrules(TAid []int)(int,*rules){
	for id, v := range rulesArr {
		if !lib.EqualArrs(TAid,v.TAid) {
			continue
		}
		return id,v
	}

	return 0,nil
}
/////////////////////////////////////////






//////////////////// сохранить Образы rules
// ID|TAid через ,
func SaveRulesArr(){
	var out=""
	for k, v := range rulesArr {
		out+=strconv.Itoa(k)+"|"
		for i := 0; i < len(v.TAid); i++ {
			out+=strconv.Itoa(v.TAid[i])+","
		}
		out+="\r\n"
	}
	lib.WriteFileContent(lib.GetMainPathExeFile()+"/memory_psy/rules.txt",out)

}
////////////////////  загрузить образы rules
// ID|TAid через ,
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

		s:=strings.Split(p[1], ",")
		var TAid[] int
		for i := 0; i < len(s); i++ {
			if len(s[i])==0{
				continue
			}
			si,_:=strconv.Atoi(s[i])
			TAid=append(TAid,si)
		}
var saveDoWritingFile= doWritingFile; doWritingFile =false
		isNotLoading=false
		createNewlastrulesID(id,TAid)
		isNotLoading=true
doWritingFile =saveDoWritingFile
	}
	return

}
///////////////////////////////////////////


