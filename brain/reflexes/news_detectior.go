/* Детектор нового условия для создания условного рефлекса
детектор нового выявляет новые условия причинного (предшествовавшего) стимула,
не приводящего к рефлексу, в дополнение к условиям активного рефлекса (безусловного или условного)
и обрабатывает это в updateNewsConditions()
*/

package reflexes

import (
	"BOT/brain/gomeostas"
	"BOT/lib"
)

/* массив нового сочетания пусковых (Trigger) стимулов в восприятии
Накапливает образы в течение бодрствования, т.е. не сохраняет в файле.
Идентичен TriggerStimuls но с Count int.
если Count>3 - создается условный рефлекс, а TriggerStimulsTemp записывается как новый образ TriggerStimuls,
который участвует в распознавании рефлексов в дереве рефлексов.

Массивы TriggerStimulsTemp не сохраняются в памяти т.к.
стимул должен набрать 3 повторения в течение ближайшего времени
- за время бодрствования.
*/
type TriggerStimulsTemp struct {
	ID int 								// идентификатор данного сочетания пусковых стимулов
	TriggerStimulsID int 	// ID образа пускового стимула типа TriggerStimulsID
	Count int 						// число повторений до оразования рефлекса
}

var TriggerStimulsTempArr = make(map[int]*TriggerStimulsTemp)

// создать образ сочетаний пусковых стимулов
// В случае отсуствия пусковых стимулов создается ID такого отсутсвия, пример такой записи: 2|||0|0| - ID=2
var lastTriggerStimulsTempID = 0

func createNewlastTriggerStimulsTempID(id int, ts *TriggerStimuls) (int, *TriggerStimulsTemp) {
	oldID, oldVal := checkUnicumTriggerStimulsTemp(ts.RSarr, ts.PhraseID, ts.ToneID, ts.MoodID)
	if oldVal != nil { // рефлекс с таким условием уже есть
		return oldID, oldVal
	}
	if id == 0 {
		lastTriggerStimulsTempID++
		id = lastTriggerStimulsTempID
	} else {
		if lastTriggerStimulsTempID < id {
			lastTriggerStimulsTempID = id
		}
	}

	var node TriggerStimulsTemp
	node.ID = id
	node.TriggerStimulsID = ts.ID
	TriggerStimulsTempArr[id] = &node
	return id, &node
}
func checkUnicumTriggerStimulsTemp(bArr []int, PhraseID []int, ToneID int, MoodID int) (int, *TriggerStimulsTemp) {
	for id, v := range TriggerStimulsTempArr {
		ts := TriggerStimulsArr[v.TriggerStimulsID]
		if !lib.EqualArrs(bArr, ts.RSarr) { continue }
		if !lib.EqualArrs(PhraseID, ts.PhraseID) { continue }
		if ToneID != ts.ToneID || MoodID != ts.MoodID { continue }
		return id, v
	}

	return 0, nil
}

/* Детектор нового news_detectior.go выявляет новые условия
причинного (предшествовавшего имеющемуся рефлесу) стимула, пока не приводящего к рефлексу,
в дополнение к условиям активного рефлекса (безусловного или условного).
Чтобы сделать у.рефлекс нужно чтобы в ксловиях имеющегося рефлекса появился новый признак
(пуск. стимул или слово). Т.е. перед пусковым стимулом нужно запостить слово.
 */
func updateNewsConditions(rank int) {
	// Не должен работать до рождения!
	if EvolushnStage < 1 { return }

	// была ли перед действиями Beast Причина == образ пусковых ситимулов, не приводящих к действиям
	if oldActiveCurTriggerStimulsID == 0 { return	} // нет значимых причин для поиска новизны
	if ActivedTerminalImage == nil { return	} // нет действий, которые могли быть порождены причиной
	TriggImage := TriggerStimulsArr[oldActiveCurTriggerStimulsID]
	if TriggImage == nil { // почему-то нет такого образа пусковых стимулов? а должен быть т.к. при любых событиях с Пульта они записываются
		return
	}
	// есть ли такой TriggerStimulsTemp?
	_, tempImg := checkUnicumTriggerStimulsTemp(TriggImage.RSarr, TriggImage.PhraseID, TriggImage.ToneID, TriggImage.MoodID)
	if tempImg == nil {
		_, tempImg = createNewlastTriggerStimulsTempID(0, TriggImage)
	}

	// должно уже быть не менее 2 событий образования рефлекса
	// или включен режим IsUnlimitedMode "набивка рабочих фраз без отсеивания мусорных слов"
	if tempImg.Count > 2 || IsUnlimitedMode == 1 {
		if IsUnlimitedMode==1 && tempImg.Count < 3 { // искусственно добавить в счетчик
			tempImg.Count = 3
		}
		/* нужно создать рефлекс или если условия те же, но действия уже другие
		   - подставить в существующий рефлекс новые действия (делается в createNewConditionReflex)
		    Иначе рефлекс просто остается прежним.
		*/
		// базовыми условиями рефлекс становится текущие условия:
		lev1 := gomeostas.CommonBadNormalWell
		lev2 := gomeostas.GetCurContextActiveIDarr()
		// ранг rank наследуется от условного рефлекса, от которого берутся действия:
		CreateNewConditionReflex(0, lev1, lev2, oldActiveCurTriggerStimulsID, oldActivedTerminalImage, rank)
		SaveConditionReflex()
	} else { // просто увеличить счетчик
		tempImg.Count++
	}
	/*
		// сравниваем только по пусковым стимулам ActiveCurTriggerStimulsID
		//потому как полные безусловные рефлексы всегда уже соотвествуют базовым условиям
		pCond := TriggerStimulsArr[ActiveCurTriggerStimulsID]
		reflex:=GeneticReflexes[reflexID]
		if len(pCond.RSarr)>0{
			diff:=lib.GetDifferentIntArs(pCond.RSarr, reflex.lev3)
	    if len(diff)>0{
				temp.RSarr=diff
	    }
		}
	*/
	// обнулить использованную причину
	oldActiveCurTriggerStimulsID = 0
	// дезактивировать использованный активный образа сочетаний пусковых стимулов
	ActiveCurTriggerStimulsID = 0
}