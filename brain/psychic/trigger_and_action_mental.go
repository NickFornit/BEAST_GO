/* Образ ментального Правила.
Ментальное Правило используется как опыт использования ментальных автоматизмов для нахождения решения.

Начало фиксации Правила - ментальный запуск моторного автомтаизма (MentalAutomatizm.ActionsImageID ->  activateMotorID).
Эффект Правила отражает насколько мент.автоматизм привел к успешному решению,
т.е. оценивается в момент действий оператора в период ожидания.
После этого в эпизод память и MentalTriggerAndActionArr записывается новое Правило.
Так же при этом проставляется MentalAutomatizm.Usefulness автоматизму, запустившему действие (вообще MentalAutomatizm нивелируется Правилами).

При ментальном запуске моторного автомтаизма фиксируется последний фрагмент Кратковременной памяти (saveFromNextIDcurretCicle),
начиная с последней объективной активации consciousness т.е. только цепочка для данной активности дерева актоматизмов,
но могут быть разные активности дерева понимания из-за произвольной активации
(т.к. при переактивации деревьев могут измениться базовые циклы) и эти звенья оцениваются как успешнвые или нет.
Это Правило записывается в массив Правил rulesMentalArr, и в кадр эпиз.памяти с Type=1 .

*/
package psychic

import (
	"BOT/lib"
	"strconv"
	"strings"
)

/////////////////////////
/* Action - уже есь в реализованном последнем кадре ShortTermMemoryID
так что можно было бы не вводить переменную Action, но она явно не помешает.
Правило заключается в том, что если в прежнем опыте rulesMentalArr или кадрах эпизод.памяти (с Type=1)
встречается текущий набор ShortTermMemoryID с Effect>0, то можно смело ставить мент.автоматизм с действием Action.
А если встречается бОльшая часть такой последовательности, можно рискнуть поставить мент.автоматизм с действием Action.

Т.е. с кажой инерацией цикла осмысления можно смотреть, если ли уже похожий опыт.
В отличие от объективнх Правил здесь нет пошаговой сверки с реакцией Оператора (это внутренний цикл поиска решения),
так что нужно просто смотреть наиболее подходящее Правило.
Каждое Правило - один из усвоенных алгоритмов поиска решений.
С каждой итерацией поиска решения можно сверяться с таким опытом и при оптимистичном результате не искать дальше,
а применить действие.
 */
type MentalTriggerAndAction struct {
	ID int
	ShortTermMemoryID []int // последний фрагмент Кратковременной памяти из saveFromNextIDcurretCicle []int
	Action int // образ ответных действий - всегда MentalAutomatizm.ActionsImageID ->  activateMotorID
	Effect int // эффект от действий накапливается при каждой новой перезаписи и используется уже суммарное значение.
}
////////////////////////

var MentalTriggerAndActionArr=make(map[int]*MentalTriggerAndAction)

//////////////////////////////////////////

// вызывается из psychic.go
func MentalTriggerAndActionInit(){
	loadMentalTriggerAndActionArr()
}


////////////////////////////////////////////////
// создать новое сочетание ответных действий если такого еще нет
var lastMentalTriggerAndActionID=0
func createNewlastMentalTriggerAndActionID(id int,ShortTermMemoryID []int,Action int,Effect int)(int,*MentalTriggerAndAction){
	if Effect<0{Effect=-1}
	if Effect>0{Effect=1}

	oldID,oldVal:=checkUnicumMentalTriggerAndAction(ShortTermMemoryID,Action,Effect)
	if oldVal!=nil{
		return oldID,oldVal
	}
	if id==0{
		lastMentalTriggerAndActionID++
		id=lastMentalTriggerAndActionID
	}else{
		//		newW.ID=id
		if lastMentalTriggerAndActionID<id{
			lastMentalTriggerAndActionID=id
		}
	}

	var node MentalTriggerAndAction
	node.ID = id
	node.ShortTermMemoryID = ShortTermMemoryID
	node.Action=Action
	node.Effect=Effect

	MentalTriggerAndActionArr[id]=&node

	if doWritingFile{SaveMentalTriggerAndActionArr() }

	return id,&node
}
func checkUnicumMentalTriggerAndAction(ShortTermMemoryID []int,Action int,Effect int)(int,*MentalTriggerAndAction){
	for id, v := range MentalTriggerAndActionArr {
		if !lib.EqualArrs(ShortTermMemoryID,v.ShortTermMemoryID) {
			continue
		}
		if Action != v.Action {
			continue
		}
		if Effect != v.Effect {
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
func SaveMentalTriggerAndActionArr(){
	var out=""
	for k, v := range MentalTriggerAndActionArr {
		out+=strconv.Itoa(k)+"|"
		for i := 0; i < len(v.ShortTermMemoryID); i++ {
			out+=strconv.Itoa(v.ShortTermMemoryID[i])+","
		}
		out+=strconv.Itoa(v.Action)+"|"
		out+=strconv.Itoa(v.Effect)
		out+="\r\n"
	}
	lib.WriteFileContent(lib.GetMainPathExeFile()+"/memory_psy/trigger_and_actions_mental.txt",out)

}
////////////////////  загрузить образы стимула (действий оператора) - ответа Beast
// ID|ActID через ,|PhraseID через ,|ToneID|MoodID|
func loadMentalTriggerAndActionArr(){
	MentalTriggerAndActionArr=make(map[int]*MentalTriggerAndAction)
	strArr,_:=lib.ReadLines(lib.GetMainPathExeFile()+"/memory_psy/trigger_and_actions_mental.txt")
	cunt:=len(strArr)
	for n := 0; n < cunt; n++ {
		if len(strArr[n])==0{
			continue
		}
		p:=strings.Split(strArr[n], "|")
		id,_:=strconv.Atoi(p[0])

		s:=strings.Split(p[1], ",")
		var ShortTermMemoryID[] int
		for i := 0; i < len(s); i++ {
			si,_:=strconv.Atoi(s[i])
			ShortTermMemoryID=append(ShortTermMemoryID,si)
		}
		Action,_:=strconv.Atoi(p[1])
		Effect,_:=strconv.Atoi(p[2])
var saveDoWritingFile= doWritingFile; doWritingFile =false
		createNewlastMentalTriggerAndActionID(id,ShortTermMemoryID,Action,Effect)
doWritingFile =saveDoWritingFile
	}
	return

}
///////////////////////////////////////////

