/* Настроение: сила Плохо -10 ... 0 ...+10 Хорошо.
Ощущение силы Плохо и сил Хорошо и это информационно для оценок последствий действий.
Настроение в зависимости от гомеостатический параметров
а так же оченка опасноного состояния
*/


package psychic

import (
	actionSensor "BOT/brain/action_sensor"
	"BOT/brain/gomeostas"
	"BOT/lib"
)

////////////////////////////////////////
/* CurrentMood - ТЕКУЩАЯ ОЦЕНКА ИЗМЕНЕНИЯ гомео-СОСТОЯНИЯ (гомео-настроения) - сила ощущаемого настроения PsyBaseMood
-10 0 +10 изменение настроения
- основа для стремления и избегания.
*/
var CurrentMood=0 // Так же текущее значение - в CurrentInformationEnvironment.Mood
// предыдуще значение. Так же предыдущее значение - в OldInformationEnvironment.Mood
var PreviousMood=0

/* Определячет первый уровень Дерева Понимания (дерева ментальных автоматизмов): PsyBaseMood.
Субъективно ощущаемая оценка, текущее осознаваемое настроение, которое можно произвольно изменять.
Она стремится к нулю со временем.
Постоянное состояние Хорошо довольно быстро уходит, постоянное состояние Плохо уходит гораздо медленнее.
При резких изменениях возникает эффект "маятника настроения":
появление противоположного по знаку настроения, но меньшего значения.
Значение обновляется при значительных изменениях (CurrentMood - PreviousMood)
Предыдущее значение - в OldInformationEnvironment.PsyMood
 */
var PsyMood=0 // плохо -10...0...10 хорошо - сила ощущаемого настроения PsyBaseMood  корректируется с каждым пульсом moodePulse()
var PsyMoodPulse=0 // последнее изменение PsyMood в пульсах
var veryMachChanged=false // true - было сильное изменение настроения
/////////////////////////////////////////////////////////////////////////////////////
var PsyBaseMood=0 // -1 Плохое настроение, 0 Нормальное, 1 - хорошее настроение

/////////////////////////////////////////////////////////////////////////////////////

/* ПУЛЬС психики
var PulsCount=0 // передача тика Пульса из brine.go
var LifeTime=0
var EvolushnStage=0 // стадия развития
var IsSleeping=false
 */
func moodePulse(){
	updatePsyMood(false)
}


/* Обновляется при инициализации GetCurrentInformationEnvironment() до и после совершения действий
Это безусловно-рефлекторная оценка на освнове изменений
1) гомео параметров (улавливает их изменение >4%) и 2) основы базовых контекстов: 0- НОРМА, 1-ПЛОХО, 2-ХОРОШО
Хотя базовые контексты уже создаются на основе жизненных параметров,
но это - безусловный рефлекс на основе безусловного рефлекса активации базовых контекстов.
На этой основе происходит конкуренция Целей при их выборе
и, как следствие, произвольная оценка эмоционального состояния.
*/
func GetCurMood()(int){
	danger:=GetAttentionDanger()
	var mood =0
	GomeostazParams:=gomeostas.GomeostazParams
	OldGomeostazParams:=gomeostas.OldGomeostazParams
	CommonBadNormalWell:=gomeostas.CommonBadNormalWell
	CommonOldBadValue:=gomeostas.CommonOldBadValue
	for n := 0; n < len(GomeostazParams); n++ {
		if OldGomeostazParams[n]!=0{
			// процент изменения параметра от предыдущего значения позитивный или негативный:
			var proc=100*(GomeostazParams[n]-OldGomeostazParams[n])/OldGomeostazParams[n]

			switch n {
			case 0: // энергия   +/- 10
				if proc > 4{
					if danger{ mood+=6} else{mood+=3}
				}
				if proc < -4{
					if danger{ mood-=6} else{mood-=3}
				}
			case 1: // стресс  +/- 20
				if proc > 6{
					if danger{ mood-=4} else{mood-=2}
				}
				if proc < -6{
					if danger{ mood+=4} else{mood+=2}
				}
			case 2: // гон
				if proc > 20{
					if danger{ mood-=1} else{mood-=1}
				}
				if proc < -20{
					if danger{ mood+=1} else{mood+=1}
				}
			case 3: // потребность в общении +/- 20
				if proc > 10{
					if danger{ mood-=1} else{mood-=1}
				}
				if proc < -10{
					if danger{ mood+=1} else{mood+=1}
				}
			case 4: // потребность в обучении +/- 20
				if proc > 10{
					if danger{ mood-=1} else{mood-=1}
				}
				if proc < -10{
					if danger{ mood+=1} else{mood+=1}
				}
			case 5: // 	Поиск любопытство  +/- 20
				if proc > 10{
					if danger{ mood+=1} else{mood+=1}
				}
				if proc < -10{
					if danger{ mood-=1} else{mood-=1}
				}
			case 6: // Самосохранение (Жадность, эгоизм, самозащита, страх смерти. Зависит от ситуации, может сам уменьшаться при благополучии.)
				if proc > 20{
					if danger{ mood-=1} else{mood-=1}
				}
				if proc < -20{
					if danger{ mood+=1} else{mood+=1}
				}
			case 7: // Повреждения
				if proc > 10{
					if danger{ mood-=1}
				}else {
					if proc > 20 {
						if danger {
							mood -= 5
						}
					}
				}
				if proc > 40 {
				// при proc==40 mood=-6 а при proc==100 mood=-10
					mood=-int(10.0 - (100.0-proc)/15.0)
				}
			}//switch n
		}
		OldGomeostazParams[n]=GomeostazParams[n]
	}

	switch CommonBadNormalWell {
	case 2: // НОРМА
		if CommonOldBadValue==1 { // было ПЛОХО
			if danger{ mood+=5} else{mood+=2}
		}
		if CommonOldBadValue==3 { // было ХОРОШО
			if danger{ mood-=5} else{mood-=2}
		}
	case 1: // ПЛОХО
		if CommonOldBadValue==3 { // было ХОРОШО
			if danger{ mood-=5} else{mood-=2}
		}
		if CommonOldBadValue==2 { // было НОРМА
			if danger{ mood+=3} else{mood+=1}
		}
	case 3: // ХОРОШО
		if CommonOldBadValue==2 { // было НОРМА
			if danger{ mood+=3} else{mood+=1}
		}
		if CommonOldBadValue==1 { // было ПЛОХО
			if danger{ mood+=6} else{mood+=3}
		}
	}
	if mood > 10{mood=10}
	if mood < -10{mood=-10}

	PreviousMood=CurrentMood
	CurrentMood=mood

	updatePsyMood(true)
	updatePsyBaseMood()

	return mood
}
///////////////////////////////////////////////////////////////////



/* оценить опасность текущей ситуации: да-нет
 */
func GetAttentionDanger()(bool){
	// опасные жизненные параметры
	for k, pID := range gomeostas.BadNormalWell {
		if pID==1 { // плохо для данного параметра гомеостаза
			if k == 1 || k == 2 || k == 7 || k == 8 {
				return true
			}
		}
	}
	// по опасному действию с Пульта

	aArr:=actionSensor.CheckCurActionsContext()
	for i := 0; i < len(aArr); i++ {
		if aArr[i]==3 || aArr[i]==10 || aArr[i]==12 || aArr[i]==15{
			return true
		}
	}
	/*
	// по опасной фразе с Пульта
	if isDangerWordFromPult(){
		return true
	}
	*/
	return false
}
////////////////////////////////////////////////////////////////



/////////////////////////////////////////////
/* PsyMood стремится к нулю со временем.
Постоянное состояние Хорошо довольно быстро уходит, постоянное состояние Плохо уходит гораздо медленнее.
При резких изменениях возникает эффект "маятника настроения":
появление противоположного по знаку настроения, но меньшего значения.
 */
func updatePsyMood(newMood bool){
	if PreviousMood==0{// первое значение - просто копируем
		PsyMood=CurrentMood
		PsyMoodPulse=PulsCount
		return
	}
// было ли очень сильное изменение - для маятника настроения
if (lib.Abs(PsyMood - CurrentMood) > 6) || (lib.Abs(CurrentMood)>=3 && lib.IsDiffersOfSign(PsyMood,CurrentMood)) {
	veryMachChanged=true
}
// PsyMood изменяется при значительных измнениях CurrentMood
	if lib.Abs(PsyMood - CurrentMood) > 3{
		PsyMood=CurrentMood
		PsyMoodPulse=PulsCount
		return
	}
	// PsyMood изменяется при смене знака значения (одно положительное, другое отрицательное),
	// большего, чем 1 (чтобы не было постоянной смены около нуля)
	if lib.Abs(CurrentMood)>=1 && lib.IsDiffersOfSign(PsyMood,CurrentMood){
		PsyMood=CurrentMood
		PsyMoodPulse=PulsCount
		return
	}
	// маятник настроения - через 30 пульсов
	if veryMachChanged && (PulsCount-PsyMoodPulse>30){
		PsyMood=-(PsyMood-3)// смена настроение на противоположное, но уменьшенное значение
		PsyMoodPulse=PulsCount
		veryMachChanged=false
		return
	}
	////////// режим постепенного угасания
	if PsyMood>0 { // хорошее настроение угасает довольно быстро
		if PulsCount-PsyMoodPulse > 30 {
			if lib.Abs(CurrentMood) > 0 {
				if PsyMood > 0 {
					PsyMood--
				} else {
					PsyMood++
				}
				PsyMood = CurrentMood
				PsyMoodPulse = PulsCount
			}
		}
	}
	if PsyMood<0 { // плохое настроение более важно, угасает медленнее
		if PulsCount-PsyMoodPulse > 60 {
			if lib.Abs(CurrentMood) > 0 {
				if PsyMood > 0 {
					PsyMood--
				} else {
					PsyMood++
				}
				PsyMood = CurrentMood
				PsyMoodPulse = PulsCount
			}
		}
	}
}
/////////////////////////////////////////////

/* ощущаемое настроение (-1,0,1)
обновляется сразу после updatePsyMood

 */
func updatePsyBaseMood(){
	if PsyMood<=1{
		PsyBaseMood=-1
	}
	if PsyMood<2 && PsyMood>-2{
		PsyBaseMood=0
	}
	if PsyMood>=1{
		PsyBaseMood=1
	}
}
///////////////////////////////////////////////

