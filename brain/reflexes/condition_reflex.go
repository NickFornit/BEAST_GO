/*  Условные рефлексы

Для образования условных рефлексов необходимо:
1. Совпадение во времени (сочетание) какого-либо индифферентного раздражителя (условного)
	с раздражителем, вызывающим соответствующий безусловный рефлекс (безусловный раздражитель).
2. Необходимо, чтобы действие условного раздражителя несколько предшествовало действию безусловного.
3. Условный раздражитель должен быть не вызывающим значительной самостоятельной реакции.
4. Рефлекс возникает только после нескольких повторений сочетаний 2 (news_detectior.go: if tempImg.Count > 2 - в func updateNewsConditions(rank int)),
это избавляет от случайных связей.
И возникаеющий образ рефлекса гасится, если долго не подтверждаются сочетания, за ночь - точно.
Условия затухания условного рефлекса:
1. Долгое отсутствие пускового стимула (узла ветки, с которого он запускается) – это легко реализуется добавлением в структуру у.рефлекса lastActivation int – в числе пульсов и времени протухания истекших рефлексов;
2. Действие конкурентных раздражителей – т.е. подавление конкурентными рефлексами и автоматизмами – т.е. если есть более значимый рефлекс или любой автоматизм на тот же пусковой стимул, то он блокирует у.рефлекс. В структуре у.рефлекса есть его ранг (число цепочки родителей), чем он выше, тем рефлекс приоритетнее среди других. Автоматизм приоритетнее любого рефлекса.
3. При отсутствии “подкрепления” после совершенного действия. Но безусловные рефлексы не угасают при этом, они безусловны и по отношению к тому, что происходит после действия и их “подкрепление” – обусловлено наследственной эволюцией: безусловные рефлексы постоянны, даны от рождения и не угасают на протяжении всей жизни.. У условных рефлексов точно та же функциональная задача, только с новыми стимулами, значит, им так же не нужно последующее подкреплением. В литературе часто путается отсутствие подкрепления с дезадаптация, а так же условные рефлексы и автоматизмы, образующиеся при осознании.
4. Разные условные рефлексы без подкрепления угасают с неодинаковой скоростью. Более "молодые" и непрочные условные рефлексы угасают быстрее, чем более "старые", прочные условно-рефлекторные связи (func conditionRexlexFound).

Условный рефлекс может образовывать на основе безусловного
или на основе имеющегося безусловного,
используя действия исходного рефлекса для новых условий.
Такие цепи рефлексов ничем не ограничены.

Формат записи - как у безусловного рефлекса, но,
в отличие от безусловного рефлексов, lev3 - только один ID образа пускового стимула типа TriggerStimulsID

РЕЗЮМЕ
1. Усл.рефлекс возникает там, где нет безусловного на основе нового стимула N,
привзяывая к нему действия того рефлекса (условного или безусловного),
которое вызывало реакцию ПОСЛЕ данного нового стимула M.
2. Но темерь если в восприятии появляется стимул N, то вызываемый им условный рефлекс перекрывает все рефлексы более низкого уровня,
в том числе условные меньшего уровня. Так что в структуре бесусловного рефлекса есть параметр: rank int,
который увеличивается, если реакция наследуется от условного рефлекса
и тогда рефлекс с рангом выше, перекрывает все рангом ниже.
По умолчаню у безусловных рефлексов rank равен 0.

Формат записи:
ID|lev1|lev2 через ,|lev3 типа TriggerStimulsID|ActionIDarr через ,|rank|lastActivation|birthTime

для набивки у.рефылексов Сначала (А НЕ ОДНОВРЕМЕННО!) стимул, потом рефлекс:
ПЕРЕД ПУСКОВЫМ СТИМУЛОМ НУЖНО ЗАПОСТИТЬ СЛОВО
*/

package reflexes

import (
	"BOT/lib"
	"strconv"
	"strings"
)


/* 1 - это режим - без ограничений числа повторов для формирования у.рефлекса
Он устанавливается галкой в Пульте "набивка рабочих фраз без отсеивания мусорных слов"
Работает в news_detectior.go: if tempImg.Count > 2 - в func updateNewsConditions(rank int)
 */
var IsUnlimitedMode=0

func initConditionReflex() {
	loadConditionReflexes()
	// psychic.PsychicInit()  после 4-го пульса!
}

type ConditionReflex struct {
	ID   int
	lev1 int
	lev2 []int
	// ID образа пускового стимула типа TriggerStimulsID, в отличие от безусловного рефлекса, а только один пусковой
	lev3 int
	ActionIDarr []int
	// ранг рефлекса (число цепочки родителей), чем он выше, тем рефлекс приоритетнее среди других условных
	rank int
	/* время последней активации в числе пульсов времени жизни
	- для отключения рефлекса при неиспользовании в течении 50 суток жизни 3600*24*10 = 864000 пульсов,
	но каждое использование укрепляют рефлекс на 10 дней жизни:
	conditionRexlexFound().
	*/
	lastActivation int
	// время рождения в LifeTime т.к. более "молодые" и непрочные условные рефлексы угасают быстрее, чем более "старые".
	birthTime int
}

// у-рефлексы
var ConditionReflexes = make(map[int]*ConditionReflex)
// у.рефлексы - по значению ConditionReflex.lev3 (ID пускового стимула )
var ConditionReflexesFrom3 = make(map[int][]*ConditionReflex)

// последний ID в массиве у-рефлексов
var lastConditionReflexID = 0
/* Детектор нового news_detectior.go выявляет новые условия причинного (предшествовавшего имеющемуся рефлесу) стимула,
пока не приводящего к рефлексу,
в дополнение к условиям активного рефлекса (безусловного или условного)
и обрабатывает это в updateNewsConditions(
Должно уже быть не менее 2 событий образования нововго условного рефлекса
 */
func CreateNewConditionReflex(id int, lev1 int, lev2 []int, lev3 int, ActionIDarr []int, rank int) (int, *ConditionReflex) {
	// посмотреть, если рефлекс с такими же условиями уже есть
	idOld, rOld := compareCRUnicum(lev1, lev2, lev3)
	if idOld > 0 {
		// если условия те же, но действия уже другие - подставить в существующий рефлекс новые действия
		if !lib.EqualArrs(rOld.ActionIDarr, ActionIDarr) {
			rOld.ActionIDarr = ActionIDarr
		}
		return idOld, rOld
	}

	if id == 0 {
		lastConditionReflexID++
		id = lastConditionReflexID
	} else {
		if lastConditionReflexID < id {
			lastConditionReflexID = id
		}
	}

	var newW ConditionReflex
	newW.ID = id
	newW.lev1 = lev1
	newW.lev2 = lev2
	newW.lev3 = lev3
	newW.ActionIDarr = ActionIDarr
	newW.rank = rank
	newW.lastActivation = int(LifeTime/(3600*24)) // последняя активация
	newW.birthTime = int(LifeTime/(3600*24)) // время рождения

	ConditionReflexes[id] = &newW
	ConditionReflexesFrom3[lev3] = append(ConditionReflexesFrom3[lev3],&newW)
	return id, &newW
}

// посмотреть, если условный рефлекс с такими же условиями уже есть
func compareCRUnicum(lev1 int, lev2 []int, lev3 int) (int, *ConditionReflex) {
	for k, v := range ConditionReflexes {
		if v.lev1 == lev1 && lib.EqualArrs(v.lev2, lev2) && v.lev3 == lev3 {
			// если это просроченный рефлекс, то установить его lastActivation в актуальное состояние
			v.lastActivation = int(LifeTime/(3600*24)) // последняя активация
			v.birthTime = int(LifeTime/(3600*24)) // время рождения
			return k, v
		}
	}
	return 0, nil
}

// сохранить имеющиеся условные рефлексы
/* формат такой же как у безусловных (ID|lev1|lev2_1,lev2_2,...|lev3|actin_1,actin_2,...)
но отличаетсмя для lev3, который есть - только один ID образа пускового стимула типа TriggerStimulsID
*/
func SaveConditionReflex() {
	var out = ""
	for k, v := range ConditionReflexes {
		out += ListConditionReflex(k, v) + "\r\n"
	}
	lib.WriteFileContent(lib.GetMainPathExeFile()+"/memory_reflex/condition_reflexes.txt", out)
}

// Строка условного рефлекса по ID и value
func ListConditionReflex(k int, v *ConditionReflex)string {
	var out = ""

	out += strconv.Itoa(k) + "|"
	out += strconv.Itoa(v.lev1) + "|"
	for i := 0; i < len(v.lev2); i++ {
		if i > 0 { out += ","	}
		out += strconv.Itoa(v.lev2[i])
	}
	out += "|"
	out += strconv.Itoa(v.lev3) + "|"
	for i := 0; i < len(v.ActionIDarr); i++ {
		if i > 0 { out += "," }
		out += strconv.Itoa(v.ActionIDarr[i])
	}
	out += "|"
	out += strconv.Itoa(v.rank) + "|"
	out += strconv.Itoa(v.lastActivation) + "|"
	out += strconv.Itoa(v.birthTime)

	return out
}

/*  загрузить  условные рефлексы из файла в формате
ID|lev1|lev2 через ,|lev3 типа TriggerStimulsID|ActionIDarr через ,|rank|lastActivation|birthTime
 в отличие от безусловного рефлекссв, а только один ID образа пускового стимула типа TriggerStimulsID
*/
func loadConditionReflexes() {
	path := lib.GetMainPathExeFile()
	lines, _ := lib.ReadLines(path + "/memory_reflex/condition_reflexes.txt")
	for i := 0; i < len(lines); i++ {
		if len(lines[i]) < 4 { continue	}
		p := strings.Split(lines[i], "|")
		id, _ := strconv.Atoi(p[0])
		lev1, _ := strconv.Atoi(p[1])
		// второй уровень
		pn := strings.Split(p[2], ",")
		var lev2 []int
		for i := 0; i < len(pn); i++ {
			b, _ := strconv.Atoi(pn[i])
			if b > 0 {
				lev2 = append(lev2, b)
			}
		}
		// третий уровень
		lev3, _ := strconv.Atoi(p[3])

		pn = strings.Split(p[4], ",")
		var ActionIDarr []int
		for i := 0; i < len(pn); i++ {
			b, _ := strconv.Atoi(pn[i])
			if b > 0 {
				ActionIDarr = append(ActionIDarr, b)
			}
		}
		rank, _ := strconv.Atoi(p[5])
		lastActivation, _ := strconv.Atoi(p[6])
		birthTime, _ := strconv.Atoi(p[7])

		_, r := CreateNewConditionReflex(id, lev1, lev2, lev3, ActionIDarr, rank)
		r.lastActivation = lastActivation
		r.birthTime = birthTime
	}
	return
}

/* Угас ли рефлекс или его можно использовать?
Вызывается:
1) при каждом срабатывании рефлекса
2) для проверки состояния рефлекса
В первом случае рефлекс продлевает время жизни,
во втором случае он может быть пассивирован, если время жизни превысило период его угасания
Возвращает true если рефлекс активен, - false - если рефлекс угас.
Принцип:
при каждой активации рефлекса (в conditionRexlexFound()) его время жизни продлевается
за счет перезаписывания времени рождения birthTime - уменьшая его вплоть до 0.

Если рефлекс пересоздается (его актуальность подтверждается новым сочетанием причины и следствия),
то его время жизни обновляется в func compareCRUnicum(
 */
func checkReflexLifeTime(reflex *ConditionReflex) bool {
	// рефлексы, только что созданные автоматически не проверять, они всегда новые:
	if reflex.lastActivation == 0  { // !!! только только что созданные || (reflex.lastActivation - reflex.birthTime)==0
		reflex.lastActivation = int(LifeTime / (3600 * 24)) // последняя активация
		return true
	}
	// время жизни рефлекса
	life:= reflex.lastActivation - reflex.birthTime
	if life != 0 && life > 50 { // рефлекс угас и не должен использоваться
		return false
	}
	// последняя активация
	reflex.lastActivation = int(LifeTime / (3600 * 24))
	// удлинить время жизни на 10 дней
	reflex.birthTime -= 10
	if reflex.birthTime < 0 {
		reflex.birthTime = 0
	}
	// SaveConditionReflex() reflex записывается при текущем сеансе сохранения памяти
	return true
}

// обновить время жизни всех рефлексов
func ClinerTimeConditionReflex() string {
	for _, v := range ConditionReflexes {
		v.lastActivation = int(LifeTime / (3600 * 24)) // последняя активация
		v.birthTime = int(LifeTime / (3600 * 24)) // время рождения
	}
	SaveConditionReflex()
	return "Обновлено время жизни всех рефлексов"
}