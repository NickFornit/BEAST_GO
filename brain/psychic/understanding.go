/* Базовая система самоощущения (актуализация доступной информации),
функция поддержвания информационной среды в зависимости от текущй ситуации -
обдщая картина понимания ситуации,
с функцией лени - думать или не думать: мотивированность думать зависит от ситуации, ее опасности. 
Решение что-делать или просто игнорировать ситуацию - важнейший параметр индивидуальной адаптивности. 
Не просто искать способ выйти из опасного состояния, а искать как улучшить состояние, 
каким бы оно ни было вплоть до полной неудовлетворенности существующим, когда поис ведется в любой ситуации.
Нужна функция, определяющая лень с индивидуализированными параметрами 
(м.б. зависящими от более базовых индивид.параментров).
Функция поддерживается системой самоощущения, 
которая работает независимо от активации дерева понимания по каждому пульсу. 

Активация func consciousness - после включения или пробуждения - начало цепочки мыслей,
а так же - по активности дерева (ор.рефлекс) - срабатывает функция делать-неделать.

Кроме гомеостатического инфо.окрыжения psi_information_environment.go
есть psi_information_environment_mental.go куда помещаются результаты информационных функций.
Использует субъективную часть эпизодов памяти - субъективный тип (для записи эпизодов цепочки мыслей).
*/


package psychic

//////////////////////////////
var allowConsciousnessProcess=false // при включении и просыпании - 1 раз


var saveEvolushnStage=0 // сохранение значения уровня осмысления == стадии развития при произвольном изменении уровня


////////////////////////////////////////////////////////////
/* Главная, постоянно активная с каждым пульсом функция поддержвания информационной среды и произвольности.
Изолированная от непосредственных воздействий и поэтому самостоятельная система оценки и корректировки состояния,
происходит с формированием опыта прозвольности выполнения Правил и их выбора для данных условий (psy_Experience.go),
в частности, Правил относительно самого себя (самосознание).

По каждому пульсу система осознания активируется по ор.рефлексу 
(из ор.рефлекса - вызов функции активации самосознания и по результатам - продолжить), 
а если в этот такт нет ор.рефлекса - то из PsychicCountPuls(.
Вид активации - при вызове функции осознания.
activation_type == 0 - не бывает
activation_type == 1 - активация ориентировочным рефлексом новой ситуации
activation_type == 2 - активация "внутренним" (произвольным) ориентировочным рефлексом
*/
func consciousness(activation_type int)(bool){

	if !allowConsciousnessProcess{
		return false
	}

// TODO если период ожидания ответа LastRunAutomatizmPulsCount ??  - УЧЕСТЬ

// ЧЕТВЕРТЫЙ УРОВЕНЬ, самый примитивный уровень:
	// есть ли штатный мот.автоматизм и нужно ли его менять или задумываться
	if currentAutomatizmAfterTreeActivatedID > 0 {
		am:=AutomatizmSuccessFromIdArr[currentAutomatizmAfterTreeActivatedID]
		if am != nil && am.Belief==2 && am.Usefulness>0{// нормальный, пусть выполняется
			return false
		}
	}
	/////////////////////////////////////////////////////////


// ВТОРОЙ УРОВЕНЬ - попытка использования примитивных Правил


	/////////////////////////////////////////////////////////


// ТРЕТИЙ УРОВЕНЬ - попытка найти решение, используя всю текущую инфрмацию с учетом срочности
	// исходные параметры:
	refreshCurrentInformationEnvironment()
	// детекция ленивого состояния
	if isIdleness(){
		// пофиг все, можно лениво обрабатывать накопившиеся структуры, реагирование - на уровне - до EvolushnStage < 4
		saveEvolushnStage=EvolushnStage
		EvolushnStage = 3 // нагло и просто :) - произвольный откат уровня осознания

		processingFreeState()// обработка структур в свободном состоянии может быть долгой -

		EvolushnStage=saveEvolushnStage // возвращаем уровень осмысления, иначе зависнет на этой стадии
		return false
	}


	/////////////////////////////////////////////////////////


// ЧЕТВЕРТЫЙ УРОВЕНЬ - доминанта нерешенной проблемы - только если нет срочности



	/////////////////////////////////////////////////////////



//	ВЫБРАНА АКТИВНОСТЬ поика решения
if saveEvolushnStage >0 {// иначе оно обнуляет EvolushnStage
	EvolushnStage = saveEvolushnStage // возвращаем уровень осмысления
}

	/* приоритетное направление мотивации задается UnderstandingNode.SituationID
	    непроизвольно идет поиск решения в этом направлении.
	Не просто искать способ выйти из опасного состояния, а искать как улучшить состояние,
	каким бы оно ни было вплоть до полной неудовлетворенности существующим, когда поис ведется в любой ситуации.

	Во вторую очередь тему общения и как оно было с прогнозом получаем из эпизод.памяти
	 */







return false // не блокировать последующий код ориентировочного рефлекса.
}
////////////////////////////////////////////


///////////////////////////////
// обновление состояния информационной среды
func refreshCurrentInformationEnvironment(){
	///////// Информационная среда осознания ситуации
	// Нужно собрать всю информацию, которая может повлиять на решение.
	//  получение текущего состояния информационной среды: отражение Базового состояния и Активных Базовых контекстов
	GetCurrentInformationEnvironment()

	// оценка опасности ситуации, необходиомсть срочных действий
	veryActualSituation=CurrentInformationEnvironment.veryActualSituation
	// выявить ID парамктров гомеостаза как цели для улучшения в данных условиях
	curTargetArrID=CurrentInformationEnvironment.curTargetArrID

	/* Еще информация:
	жизненный опыт  psy_Experience.go
	доминанта psy_problem_dominanta.go
	субъектиная оценка ситуации для применения произвольности
	*/

	// актуальной инфой являются узлы активной ветки дерева понимания, особенно контекст SituationID
}
///////////////////////////////////////////////


// детекция ленивого состояния
func isIdleness()(bool){
	if veryActualSituation {
		return false
	}

	if isCurrentProblemDominanta != nil{
		return false
	}

	return true
}
///////////////////////////////////////////////////////
// обработка структур в свободном состоянии, в первую очередь - эпизодической памяти
func processingFreeState(){
	// TODO переработка происходившего, в первую очередь - эпизодической памяти
	//EpisodeMemoryLastCalcID - последний эпизод, который был осмыслен в лени или во сне
}
//////////////////////////////////////////////////////
