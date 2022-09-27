/* Информационная среда - основа текущего самоощущения


*/


package psychic

import (
	"BOT/brain/action_sensor"
	"BOT/brain/gomeostas"
	word_sensor "BOT/brain/words_sensor"
	"BOT/lib"
	"strconv"
)

/////////////////////////////////////////////////////
/* Среда условий текущего состояния - интегративная информационная среда
Образ интегративной информационной среды Формируется временно и не сохранятся в файле.
Объекты InformationEnvironment накапливаются в оперативке, указалеи входят в  эпизодическую память и удаляются во сне.

Структуру информационной среды можно дополнять по мере необходимости, т.к. она не сохраняется в файле.
*/
type InformationEnvironment struct {
	IsSleep bool // true - организм спит (во сне контекст задает тоже InformationEnvironment)
/* базовое состояние доступно из
   gomeostas.CommonBadNormalWell - общее базовое состояние
   gomeostas.BadNormalWell[] - базовое состояние по каждому параметру гомеостаза
 */
	// текущая эмоция Emotion, может быть произвольно изменена
	PsyEmotionImg *Emotion

	// опасность состояния
	danger bool // получить из GetAttentionDanger
	// общая оценка гомео-настроения
	Mood int //сила Плохо -10 ... 0 ...+10 Хорошо
	PsyMood int //Субъективно ощущаемая оценка, текущее осознаваемое настроение, которое можно произвольно изменять.

	// текущий образ сочетания действий с Пульта Activity
	PsyActionImg *Activity    // хранится до следующего любого действия или фразы с Пульта но не более 100 пусльсов
	PsyActionImgPulsCount int // момент обновления
	// текущий образ фразы с Пульта Verbal
	PsyVerbImg *Verbal

	veryActualSituation bool // оценка опасности ситуации, необходиомсть срочных действий
	curTargetArrID []int //ID парамктров гомеостаза как цели для улучшения в данных условиях

}
var InformationEnvironmentObjects []*InformationEnvironment
var CurrentInformationEnvironment InformationEnvironment
func initCurrentInformationEnvironment(){

	CurrentInformationEnvironment.PsyActionImg=nil
	CurrentInformationEnvironment.PsyActionImgPulsCount=0

	CurrentInformationEnvironment.PsyVerbImg=nil
}

var OldInformationEnvironment InformationEnvironment
//////////////////////////////////////////////////////////////////////


///////////////////////////////////////////////////////
/*  получение текущего состояния информационной среды: отражение Базового состояния и Активных Базовых контекстов
только при ориентировчном рефлексе и осмыслении результатов - обновление самоощущения! и запись кадра эпизодической памяти
 */
func GetCurrentInformationEnvironment(){
	var ie InformationEnvironment
	InformationEnvironmentObjects=append(InformationEnvironmentObjects,&ie)
	OldInformationEnvironment=CurrentInformationEnvironment
	CurrentInformationEnvironment=ie
	initCurrentInformationEnvironment()

	ie.PsyActionImgPulsCount=PulsCount// момент обновления

	 CurrentInformationEnvironment.IsSleep=IsSlipping

// определение текущего сочетания ID Базовых контекстов - оно есть всегда, даже если ничего не сделано на Пульте - нулевое сочетание.
	bsIDarr:=gomeostas.GetCurContextActiveIDarr()
	// текущая эмоция Emotion, может быть произвольно изменена
_,CurrentInformationEnvironment.PsyEmotionImg=createNewBaseStyle(0,PsyBaseMood,bsIDarr)

	CurrentInformationEnvironment.veryActualSituation,CurrentInformationEnvironment.curTargetArrID=gomeostas.FindTargetGomeostazID()

	ActID:=action_sensor.CheckCurActionsContext()
	_,CurrentInformationEnvironment.PsyActionImg=createNewlastActivityID(0,ActID)// текущий образ сочетания действий с Пульта Activity

if len(word_sensor.CurrentPhrasesIDarr)>0{
		PhraseID := word_sensor.CurrentPhrasesIDarr
		FirstSimbolID:=word_sensor.GetFirstSymbolFromPraseID(PhraseID[0])
		ToneID := word_sensor.DetectedTone
		MoodID := word_sensor.CurPultMood
		_, CurrentInformationEnvironment.PsyVerbImg = CreateVerbalImage(FirstSimbolID,PhraseID, ToneID, MoodID)
	}

	CurrentInformationEnvironment.danger=GetAttentionDanger()
	CurrentInformationEnvironment.Mood=GetCurMood()
	CurrentInformationEnvironment.PsyMood=PsyMood

// запись эпизодической памяти
	saveEpisodicMenory()

	writeInformationEnvironmentMarker()
	return
}
///////////////////////////////////////////////////////




//////////////////////////////////////////////////////
// записать метку изменения information_environment при каждом обновлении
func writeInformationEnvironmentMarker(){
	strArr,_:=lib.ReadLines(lib.GetMainPathExeFile()+"/memory_psy/self_perception_count.txt")
	var old=0
	if strArr != nil{
	old, _ = strconv.Atoi(strArr[0])
	}
	lib.WriteFileContent(lib.GetMainPathExeFile()+"/memory_psy/self_perception_count.txt",strconv.Itoa(old+1))
}
//////////////////////////////////////////////////////
