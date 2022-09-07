/*  параметры определены в http://go/pages/gomeostaz.php
Жизненные параметры гомеостаза и Базовые стили поведения (базовые контексты рефлексов)

1 - энергия 2 - стресс 3 - гон 4-потребность в общении 5-потребность в обучении
6- Поиск, 7- Самосохранение, 8 - Повреждения
файл сохранения состояния параметров: files/GomeostazParams.txt

Для параметров гомеостаза, напрямую не связанных с жизнеобеспечением (parMaxPulsCount: гон, потребность в общении,
потребность в обучении и любопытство) организована цикличность: при нарастании параметра до максимума,
он удерживается в течении 20 секунд, а потом сбрасывается.
Это позволяет создавать достаточные по времени периоды специфических контекстов реагирования.
*/
package gomeostas

//import word_sensor "BOT/brain/words_sensor"

import (
	"BOT/lib"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Текущие значения параметров гомеостаза заполняется из файла files/GomeostazParams.txt
var GomeostazParams =make(map[int]float64)
// старые значения, обновляемые при ухудшении ??????
var OldGomeostazParams=make(map[int]float64)

// скорость изменения параметров со временем
var GomeostazParamsSpeed =make(map[int]int)

// веса значимости для GomeostazSensor От 0 до 100 %
var GomeostazParamsWeight =make(map[int]int)

// правила активирования Базовых контекстов
var GomeostazActivnostArr =make(map[int][]string)

//нажатия кнопок действий Пульта дают гомео-эффект:
var GomeostazActionEffectArr =make(map[int]string)
/*нажатия кнопок действий Пульта дают mood-эффект (Плохо - Хорошо НЕ БАЗОВОЕ СОСТОЯНИЕ а НАСТРОЕНИЕ!!!):
На уровне рефлексов не оказывает действия - только в период ожидания результата действий Beast
- для (де)мотивации в automatizm_result.go
 */
var GomeostazActionCommonEffectArr =make(map[int]string)
/* Общее (де)мотивирующее действие - только в период ожидания результата действий Beast
- для (де)мотивации в automatizm_result.go
На уровне рефлексов не оказывает действия.
Значения: "" - ничего, "+" - позитив, "-" - негатив.
Держится 5 пульсов, потом становится "" - гасится в GomeostazPuls()
*/
var CommonMoodAfterAction=""
var CommonMoodAfterActionPulsCount=0

func SetGomeostazActionCommonEffectArr(actionID int){
	CommonMoodAfterAction=GomeostazActionCommonEffectArr[actionID]
	CommonMoodAfterActionPulsCount=PulsCount
}

var PeriodPulsCount=20 // 20 пульсов удерживается максимальное значение некритичных для жизни 4-х параметров
var parMaxPulsCount=make([]int, 4)// 0- гон, 1-потребность в общении, 2-потребность в обучении, 3- любопытство

/////////////////////////////////////////////////////////////


func init(){

	/*
	// подбор нужной крутизны экспоненты для BetterOrWorseNow() - множителем k
	var str="|"
	var val=0.0
	var k=0.17
	for i := 0; i < 70; i++ {
		if i>0 && i%10==0{str+="\r\n"}
		val=10.0 - 10.0/math.Exp(float64(i)*k)
		str+=strconv.Itoa(i)+":"+fmt.Sprintf("%10.2f", val)+"|"
	}
	*/

	path:=lib.GetMainPathExeFile()
	lines,_:=lib.ReadLines(path+"/memory_reflex/GomeostazParams.txt")
	if len(lines)<7{// испорчен файл, восстановить
		var def="1|10\r\n2|10\r\n3|0\r\n4|0\r\n5|0\r\n6|0\r\n7|0\r\n8|0"
		lib.WriteFileContent(lib.GetMainPathExeFile()+"/memory_reflex/GomeostazParams.txt",def)
		lines,_=lib.ReadLines(path+"/memory_reflex/GomeostazParams.txt")
	}
	for i := 0; i < len(lines); i++ {
		p:=strings.Split(lines[i], "|")
		id,_:=strconv.Atoi(p[0])
		val,_:=strconv.ParseFloat(p[1], 32)
		GomeostazParams[id]=val
		OldGomeostazParams[id]=val
	}
	//////////////////////////////////////////
	lines,_=lib.ReadLines(path+"/memory_reflex/GomeostasWeight.txt")
	for i := 0; i < len(lines); i++ {
		p:=strings.Split(lines[i], "|")
		id,_:=strconv.Atoi(p[0])
		weight,_:=strconv.Atoi(p[1])
		speed,_:=strconv.Atoi(p[2])
		GomeostazParamsWeight[id]=weight
		GomeostazParamsSpeed[id]=speed
	}
	//////////////////////////////////////////
	GomeostazActivnostArr =make(map[int][]string)
	lines,_=lib.ReadLines(path+"/memory_reflex/base_context_activnost.txt")
	for i := 0; i < len(lines); i++ {
		p:=strings.Split(lines[i], "|")
		id,_:=strconv.Atoi(p[0]) // ID параметра гомеостаза
		GomeostazActivnostArr[id]=make([]string, 7)
		GomeostazActivnostArr[id][0]=p[1] // плохо
		GomeostazActivnostArr[id][1]=p[2]     // хорошо
		GomeostazActivnostArr[id][2]=p[3]     //Норма 0-19%
		GomeostazActivnostArr[id][3]=p[4]     //Норма 20-39%
		GomeostazActivnostArr[id][4]=p[5]     //Норма 40-59%
		GomeostazActivnostArr[id][5]=p[6]     //Норма 60-79%
		GomeostazActivnostArr[id][6]=p[7]     //Норма 80-100%
	}
	//////////////////////////////////////////
	GomeostazActionEffectArr =make(map[int]string)
	GomeostazActionCommonEffectArr =make(map[int]string)
	lines,_=lib.ReadLines(path+"/memory_reflex/Gomeostaz_pult_actions.txt")
	for i := 0; i < len(lines); i++ {
		p := strings.Split(lines[i], "|")
		id, _ := strconv.Atoi(p[0]) // ID параметра гомеостаза
		GomeostazActionEffectArr[id]=p[1]
		GomeostazActionCommonEffectArr[id]=p[2]
	}
	//////////////////////////////////////////

	initBadDetector()
	initContextDetector()
return
}
/* приравнять OldGomeostazParams GomeostazParams
При каждом сравнении старого с новым приравнивать Для определения вектора изменения с каждым пульсом
 */
func copyToOldPsrams(){
	for id, _ := range GomeostazParams {
		OldGomeostazParams[id]=GomeostazParams[id]
	}
}
//////////////////////////////////////////////////////////



/////////////////////////////////////////////
/* ПУЛЬС ГОМЕОСТАЗА - обработка раз в секунду

 */
var PulsCount=0 // передача тика Пульса из brine.go
var LifeTime=0
var EvolushnStage=0 // стадия развития
var IsSlipping=false
// коррекция текущего состояния гомеостаза и базового контекста с каждым пульсом
func GomeostazPuls(evolushnStage int, lifeTime int,puls int,isSlipping bool){
	LifeTime=lifeTime
	EvolushnStage=evolushnStage
	PulsCount=puls // передача номера тика из более низкоуровневого пакета
	IsSlipping=isSlipping

	if CommonMoodAfterActionPulsCount>0 && CommonMoodAfterActionPulsCount+5 >PulsCount{
		CommonMoodAfterActionPulsCount=0
		CommonMoodAfterAction=""
	}

	if EvolushnStage>=3{//Период подражания
		IsLevelBeginParam5=true // уровень развития для Потребность в обучении достигнут
	}
	if EvolushnStage>=4{//Период преступной инициативы
		IsLevelBeginParam3=true // уровень развития для Гона достигнут
	}

	gomeostazUpdate()
	baseContextUpdate()
	badDetecting()
	// детектор изменения базового состояния и контекстов - проветка по каждому пульсу
	changingConditionsDetector()

	//При каждом сравнении старого с новым приравнивать Для определения вектора изменения с каждым пульсом
	copyToOldPsrams()
}
//////////////////////////////////////////////////////////////


// изменить параметр на величину
var NotAllowSetGomeostazParams=false
func ChangeGomeostazParametr(id int, diff float64){
	NotAllowSetGomeostazParams=true
	OldGomeostazParams[id]=GomeostazParams[id]
	GomeostazParams[id]+=diff
	if GomeostazParams[id] <0{GomeostazParams[id]=0}
	if GomeostazParams[id] >100{GomeostazParams[id]=100}
	NotAllowSetGomeostazParams=false
}
///////////////////////////////////////////




/////////////////////////////////////////////////////////////////////////////
////// выдать текущие значения жизненных параметров
func GetCurGomeoParams() (string) {
	if NotAllowSetGomeostazParams{
		return ""
	}
	var out="";
	for id, v := range GomeostazParams {
		out+=strconv.Itoa(int(id))+";"+strconv.Itoa(int(v))+"|";
	}
	return out
}
// установка параметров гомеостаза с Пульта:
func SetCurGomeoParams(parID int,parVal string){
	NotAllowSetGomeostazParams=true
	GomeostazParams[parID],_=strconv.ParseFloat(parVal, 64);
	SaveCurrentGomeoParams()
	NotAllowSetGomeostazParams=false

/*определить базовые контексты при новых пераметрах гомеостаза baseContextUpdate()
потом 	очистить сенсоры слов и стек слов и переактивировать Дерево понимания
 */
	baseContextUpdate()

	//fmt.Println("SET: ", p0Arr[0], p0Arr[1])
	return
}
////////////////////////////////////////////////////////////////////////////////

// до определенной стадии развития Гон и Потребности в обучении не влияют ни на что (нет у детей гона)
var IsLevelBeginParam3=false // true - уровень развития для Гона достигнут
var IsLevelBeginParam5=false  // true - уровень развития для Потребность в обучении достигнут

// сколько Beast уже съел и нужно превратить в энергию (сразу писать нельзя из-за паники: "concurrent map read and map write")
var FoodPortionEaten=0


/////////////////////////////////////////////////////////////
// состояние гомеостаза
func gomeostazUpdate() {

if NotAllowSetGomeostazParams{
	return
}
	if FoodPortionEaten>0{// Beast съел порцию, нужно переварить в энергию
		ChangeGomeostazParametr(1,float64(FoodPortionEaten))
		FoodPortionEaten=0
	}
	// изменение со временем
	changingParVal(1)
	changingParVal(2)
// Гон начинает изменяться с некоторого уровня развития Beast
if IsLevelBeginParam3 {
		changingParVal(3)
}
	changingParVal(4)
	// 	Потребность в обучении начинает изменяться с некоторого уровня развития Beast
	if IsLevelBeginParam5 {
		changingParVal(5)
	}
	changingParVal(6)
	changingParVal(7)

	//.............. эксклюзивные зависимости
// повреждение
	//При энергии <5% начинает увеличивается cкорость поврежедения
	var oldValGomeostazParamsSpeed8=GomeostazParamsSpeed[8]
	if GomeostazParams[1] < 5{
		// при GomeostazParams[1]==0 прибавка скорости будет 0, при при GomeostazParams[1]==0 прибавка скорости станет 100
		GomeostazParamsSpeed[8]+=100 - int(GomeostazParams[1]*20.0)
	}
	// при повышенном стрессе:
	if GomeostazParams[2] > 70{
		// при GomeostazParams[2]==70 прибавка скорости будет 0, при при GomeostazParams[2]==100 прибавка скорости станет 50
		GomeostazParamsSpeed[8]+=int((GomeostazParams[2] - 70.0)*1.7)
	}
	changingParVal(8)
	GomeostazParamsSpeed[8]=oldValGomeostazParamsSpeed8
	//ПРи 100% повреждений  - смерть. Фон пульта менятся от повреждений, начиная с критического порога.
	//Смерть - черный фор с траурной надписью, сброс памяти.
	////////////////////////////////


//самосохранение
		val1:=0.0
		if (float64(compareLimites[1])-GomeostazParams[1])>0{
			val1=(float64(compareLimites[1])-GomeostazParams[1])*2.5
		}
	val2:=0.0
	if (GomeostazParams[2]-float64(compareLimites[2]))>0{
		val2=(GomeostazParams[2]-float64(compareLimites[2]))/2
	}
	val3:=0.0
	if (GomeostazParams[8]-float64(compareLimites[8]))>0{
		val3=(GomeostazParams[8])-float64(compareLimites[8])
	}
	add:=val1+val2+val3
	// ограничить порцию изменения
	if add>1{
		add=1
	}
	GomeostazParams[7] += add
	if GomeostazParams[7]<0 {
		GomeostazParams[7]=0
	}
	if GomeostazParams[7]>100 {
		GomeostazParams[7]=100
	}
	//////////////////////////////////////////////


	if PulsCount>1 && PulsCount%10 == 0 { // записать в файл текущее состояние гомеостаза раз в 10 сек
		SaveCurrentGomeoParams()
	}
	baseContextUpdate()

	/*Для параметров гомеостаза, напрямую не связанных с жизнеобеспечением (parMaxPulsCount: гон, потребность в общении,
		потребность в обучении и любопытство) организована цикличность: при нарастании параметра до максимума,
		он удерживается в течении 20 секунд, а потом сбрасывается.
		Это позволяет создавать достаточные по времени периоды специфических контекстов реагирования.
	*/
	if GomeostazParams[3]>90 && parMaxPulsCount[0]==0 {
		parMaxPulsCount[0]=PulsCount
	}
	if GomeostazParams[4]>90 && parMaxPulsCount[1]==0 {
		parMaxPulsCount[1]=PulsCount
	}
	if GomeostazParams[5]>90 && parMaxPulsCount[2]==0 {
		parMaxPulsCount[2]=PulsCount
	}
	if GomeostazParams[6]>90 && parMaxPulsCount[3]==0 {
		parMaxPulsCount[3]=PulsCount
	}
	// сброс после периода удержания
	if GomeostazParams[3]>90 {
		if parMaxPulsCount[0]+PeriodPulsCount < PulsCount {
			parMaxPulsCount[0] = 0
			GomeostazParams[3] = 0
		}
	}else{
		parMaxPulsCount[0] = 0
	}

	if GomeostazParams[4]>90 {
		if parMaxPulsCount[1]+PeriodPulsCount < PulsCount {
			parMaxPulsCount[1] = 0
			GomeostazParams[4] = 0
		}
	}else{parMaxPulsCount[1] = 0}

	if GomeostazParams[5]>90 {
		if parMaxPulsCount[2]+PeriodPulsCount < PulsCount {
			parMaxPulsCount[2] = 0
			GomeostazParams[5] = 0
		}
	}else{parMaxPulsCount[2] = 0}

	if GomeostazParams[6]>90 {
		if parMaxPulsCount[3]+PeriodPulsCount < PulsCount {
			parMaxPulsCount[3] = 0
			GomeostazParams[6] = 0
		}
	}else{parMaxPulsCount[3] = 0}
}
/////////////////////////////////
func SaveCurrentGomeoParams(){
	if len(GomeostazParams)<7{
		return
	}
	var fStr=
		"1|"+strconv.Itoa(int(GomeostazParams[1]))+"\n" +
			"2|"+strconv.Itoa(int(GomeostazParams[2]))+"\n" +
			"3|"+strconv.Itoa(int(GomeostazParams[3]))+"\n" +
			"4|"+strconv.Itoa(int(GomeostazParams[4]))+"\n" +
			"5|"+strconv.Itoa(int(GomeostazParams[5]))+"\n" +
			"6|"+strconv.Itoa(int(GomeostazParams[6]))+"\n" +
			"7|"+strconv.Itoa(int(GomeostazParams[7]))+"\n" +
			"8|"+strconv.Itoa(int(GomeostazParams[8]))

	file, err := os.Create(lib.MainPathExeFile+"/memory_reflex/GomeostazParams.txt")
	if err != nil{
		fmt.Println("Unable to create file:", err)
		os.Exit(1)
	}
	defer file.Close()
	file.WriteString(fStr)
}
// шаг изменения парамктра со скоростью GomeostazParamsSpeed
func changingParVal(id int){
	if NotAllowSetGomeostazParams{
		return
	}
	step:=(float64(GomeostazParamsSpeed[id])/3600)
	if(id==1){
		ChangeGomeostazParametr(id,-step)
	}else{
		ChangeGomeostazParametr(id,step)
	}

}
///////////////////////////////////////////////

// true - смерть Beast при повреждении >99%
var IsBeastDeath=false // true - смерть Beast при повреждении >99%   gomeostas.IsBeastDeath
func CheckBeastDeath()(bool){
	if GomeostazParams[8] >99.0{
		IsBeastDeath=true
		return true
	}
	return false
}
/////////////////////////////////////

