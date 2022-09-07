/* БАЗОВЫЕ КОНТЕКСТЫ (основные стили поведения) - contexts.go

1	Пищевой	- Пищевое поведение, восполнение энергии, на что тратится время и тормозятся антагонистические стили поведения.
2	Поиск	- Поисковое поведение, любопытство. Обследование объекта внимания, поиск новых возможностей.
3	Игра	- Игровое поведение - отработка опыта в облегченных ситуациях или при обучении.
4	Гон	- Половое поведение. Тормозятся антагонистические стили
5	Защита	- Оборонительные поведение для явных признаков угрозы или плохом состоянии.
6	Лень	- Апатия в благополучном или безысходном состоянии.
7	Ступор	- Оцепенелость при непреодолимой опастbase_context_activnostности или когда нет мотивации при благополучии или отсуствии любых возможностей для активного поведения.
8	Страх	- Осторожность при признаках опасной ситуации.
9	Агрессия	- Агрессивное поведение для признаков легкой добычи или защиты (иногда - при плохом состоянии).
10	Злость	- Безжалостность в случае низкой оценки .
11	Доброта	- Альтруистическое поведение.
12	Сон - Состояние сна. Освобождение стрессового состояния. Реконструкция необработанной информации.

Число одновеременно активных контекстов - не более 3-х!!! Остаются только с наибольшим весом,
а лишние будут отсеиваться в порядке убывания весов контекстов.
Это имитирует распознаватель с активацией по частично-активному профилю на входе.
 */
package gomeostas

import (
	"BOT/lib"
	"sort"
	"strconv"
	"strings"
)


func init(){

}
/////////////////////////////////////
// название Базового контекста из его ID str:=gomeostas.GetBaseContextCondFromID(id)
func GetBaseContextCondFromID(id int)(string){
	var out=""
	switch id{
	case 1: out="Пищевой"
	case 2: out="Поиск"
	case 3: out="Игра"
	case 4: out="Гон"
	case 5: out="Защита"
	case 6: out="Лень"
	case 7: out="Ступор"
	case 8: out="Страх"
	case 9: out="Агрессия"
	case 10: out="Злость"
	case 11: out="Доброта"
	case 12: out="Сон"
	}
	return out
}
//////////////////////////////////////////////

/* масссив активностей базовых контекстов
активный - true, неактивный - false
 */
var BaseContextActive [15]bool // index 0 НЕ ИСПОЛЬЗУЕТСЯ! т.е. начинаетсмя с BaseContextActive[1]
var BaseContextWeight [15]int // index 0 НЕ ИСПОЛЬЗУЕТСЯ!

/* Прошивка несовместимых сочетаний контекстов
Для каждого основного контекста - антагонисты
 */
var antagonists =make(map[int][]int)

func initContextDetector(){
	antagonists =make(map[int][]int)
	path:=lib.GetMainPathExeFile()
	lines,_:=lib.ReadLines(path+"/memory_reflex/base_context_antagonists.txt")
	for i := 0; i < len(lines); i++ {
p := strings.Split(lines[i], "|")
id, _ := strconv.Atoi(p[0]) // ID параметра гомеостаза
a := strings.Split(p[1], ",")
for n := 0; n < len(a); n++ {
aID, _ := strconv.Atoi(a[n])
antagonists[id] =append(antagonists[id],aID)
}
}

	lines,_=lib.ReadLines(path+"/memory_reflex/base_context_weight.txt")
	for i := 0; i < len(lines); i++ {
		p:=strings.Split(lines[i], "|")
		id,_:=strconv.Atoi(p[0])
		val,_:=strconv.Atoi(p[1])
		BaseContextWeight[id]=val
	}

	for id, _ := range BaseContextWeight {
		BaseContextActive[id]=false
	}


/* // проверка ограничителя
	BaseContextActive[1]=true
	BaseContextActive[6]=true
	BaseContextActive[9]=true
	BaseContextActive[12]=true
	BaseContextActive[9]=true
	var activedC=make(map[int]int)
	for id, v := range BaseContextActive {
		if v{
			activedC[id]=BaseContextWeight[id]
		}
	}
	//карта только активных контекстов
	keys := make([]int, 0, len(activedC))
	for k := range activedC {
		keys = append(keys, k)
	}
	//СОРТИРОВКА ПО ЗНАЧЕНИЮ даже если значения повторяются
	sort.SliceStable(keys , func(i, j int) bool {
		return activedC[keys[i]] > activedC[keys[j]]
	})
	// ограничить только первыми тремя
	if len(keys)>3 {
		keys = keys[:3]
	}
	for id, _ := range BaseContextActive {
		BaseContextActive[id]=false
		for i := 0; i < len(keys); i++ {
			if id == keys[i]{
				BaseContextActive[id]=true
			}
		}
	}
 */



	return
}
////////////////
/* состояние базового контекста зависит
1) от выхода из нормы жизненных параметров
2) от безусловно прошитых признаков восприятия
Антагонисты конкурируют между собой со своими весами значимости.

!!! BaseContextActive[6]=true - Ступор регулируется как отсуствие реакций в данном контексте при опасности
 */
func baseContextUpdate(){
// сначала все контексты выключены:
	for id, _ := range BaseContextWeight {
		if id!=12 {// сон не гасить
			BaseContextActive[id] = false
		}
	}
	// если активен сон, все остальные неактивны
	if IsSlipping{
		//return    но во сне тоже должны быть некоторые рефлексы и эмоции
	}else {
		// При бодрствовании обязательно должен быть определен базовый контекст,
		// поэтому делаем по умолчанию Лень и гасим его в нужных случаях: BaseContextActive[6]=false
		// проверка в самом конце BaseContextActive[6] = true
	}

// Определяем текущее сочетания активных базовых контекстов:
	for i := 1; i < 9; i++ {
		if BadNormalWell[i]==2 { //НОРМА
			diapazoN:=getNormaDiapason(i)
			// правило для данного диапазона
			rule:=GomeostazActivnostArr[i][diapazoN+1] // +2 из-за 1 и 2 заняты под Хор. и Пл.
			// активируем или пассивируем контексты по заданному правилу в http://go/pages/gomeostaz.php
			activeOrPassiveContext(rule)
		}//НОРМА конец
// перекрывает предыдущее
		if BadNormalWell[i]==3 { //ХОРОШО
			rule:=GomeostazActivnostArr[i][1]
			// активируем или пассивируем контексты по заданному правилу в http://go/pages/gomeostaz.php
			activeOrPassiveContext(rule)
		}
// перекрывает предыдущее
		if BadNormalWell[i]==1 { //ПЛОХО
			rule:=GomeostazActivnostArr[i][0]
			// активируем или пассивируем контексты по заданному правилу в http://go/pages/gomeostaz.php
			activeOrPassiveContext(rule)
		}
	}

// Конкурентность антагонистов
	keys := make([]int, 0, len(BaseContextActive))
	for id, v := range BaseContextActive {
		if v {
			keys = append(keys, id)
		}
	}
	sort.Ints(keys)
	for _, id := range keys {
		for _, ida := range antagonists[id] {
			if BaseContextActive[ida] && BaseContextWeight[ida] > BaseContextWeight[id]{// активный антагонист значимее
			BaseContextActive[id]=false // погасить текущий контекст
			break // больше не нужно смотреть других антагонистов
				}else{
				BaseContextActive[ida]=false // погасить антогониста
			}
			}
	}

/*ограничение на число компонентов в образе б.контекстов: их число будет не более 3-х,
а лишние будут отсеиваться в порядке убывания весов контекстов.
Это неплохо имитирует распознаватель с активацией по частично-активному профилю на входе.
*/
	var activedC=make(map[int]int)
	for id, v := range BaseContextActive {
		if v{
			activedC[id]=BaseContextWeight[id]
		}
	}
	//карта только активных контекстов
	keys = make([]int, 0, len(activedC))
	for k := range activedC {
		keys = append(keys, k)
	}
	//СОРТИРОВКА ПО ЗНАЧЕНИЮ даже если значения повторяются
	sort.SliceStable(keys , func(i, j int) bool {
		return activedC[keys[i]] > activedC[keys[j]]
	})
	// ограничить только первыми тремя
	if len(keys)>3 {
		keys = keys[:3]
	}
	for id, _ := range BaseContextActive {
		BaseContextActive[id]=false
		for i := 0; i < len(keys); i++ {
			if id == keys[i]{
				BaseContextActive[id]=true
			}
		}
	}

}
//активируем или пассивируем контексты по заданному правилу в http://go/pages/gomeostaz.php
func activeOrPassiveContext(rule string){
	if len(rule)==0{
		return
	}
	// выделяем ID контекстов
	p:=strings.Split(rule, ",")
	for n := 0; n < len(p); n++ {
		p[n]=strings.TrimSpace(p[n])
		if len(p[n])==0{
			return
		}
		// активируем или пассивируем контексты
		cID,_:=strconv.Atoi(p[n])
		if cID>0{
			BaseContextActive[cID]=true
		}else{
			BaseContextActive[-cID]=false
		}
	}
}
/////////////////////////////////////////////////////


// для определения текущего сочетания ID Безовых контекстов  gomeostas.GetCurContextActiveIDarr()
func GetCurContextActiveIDarr()([]int) {
	var out []int
// concurrent map iteration and map write
	for id, v := range BaseContextActive {
		if(v){
			out=append(out,id)
		}
	}
	return out
}
/////////////////////////////////////////////////////////


// для Пульта
func GetCurContextActive()(string) {
	var out=""
	for id, v := range BaseContextActive {
		if(v){
			out+=strconv.Itoa(id)+";1|"
		}else{
			out+=strconv.Itoa(id)+";0|"
		}
	}
	return out
}
/////////////////////////////////////////////////////////

// контекст распознавания текущей фразы с Пульта для Vernike_detector.go
func GetActiveContextInfo()(map[int]int){
	var activeCW =make(map[int]int)

	keys := make([]int, 0, len(BaseContextActive))
	for id, v := range BaseContextActive {
		if v {
			keys = append(keys, id)
		}
	}
	sort.Ints(keys)
	for _, k := range keys {
		activeCW[k]=BaseContextWeight[k]
	}
	return activeCW
}
//////////////////////////////////////////////////////////



/* детектор изменения базового состояния и контекстов - проверка по каждому пульсу
 */
var IsNewConditions=false //  использовать из других пакетов: gomeostas.IsNewConditions

var oldBaseCondition=0
var oldActiveContextstStr=""// строка ID старого сочетания активных Базовых контекстов
func changingConditionsDetector(){

if oldBaseCondition!=CommonBadNormalWell{
	oldBaseCondition=CommonBadNormalWell
	IsNewConditions=true
	return
}

var activeContextstStr=""
	keys := make([]int, 0, len(BaseContextActive))
	for k := range BaseContextActive {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for k,v := range keys {
		if BaseContextActive[v] {
			activeContextstStr += strconv.Itoa(k)+"_" // "_" нужно разделять цифры, иначе будет ошибаться
		}
	}
if oldActiveContextstStr!=activeContextstStr{
	oldActiveContextstStr=activeContextstStr
	IsNewConditions=true
	return
}
	IsNewConditions=false
}
//////////////////////////////////////////


/* Есть ли хоть какая то активность контекстов */

func IsContextActive()bool{

  for _, v :=range BaseContextActive{

   if v{

    return true

   }

 }

 return false

}
////////////////////////////////////////////////////
