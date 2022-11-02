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

После объективной активации (activationType==1) начинается рекурсивный субъективный вызов (activationType==2)
- цикл обдумывания (субъективный ориентировочный рефлекс), каждой шаг которого основывается на информации, даваемой предыдущим шагом
с целью найти подходящие действия для данной ситуации, что дает возможность снова сориентироваться.

Кроме гомеостатического инфо.окрыжения psi_information_environment.go
есть psi_information_environment_mental.go куда помещаются результаты информационных функций.
Использует субъективную часть эпизодов памяти - субъективный тип (для записи эпизодов цепочки мыслей).

Функция consciousness() проходит через 4 уровня решений - выделенных комментариями.
Если решение не найдено на данном уровне, то отрабатывает следующий уровень сложности обработки (что характеризует эволюционную последовательность их появления).

РАБОТА consciousness() каждого уровня описана по месту.

*/


package psychic

//////////////////////////////
var AllowConsciousnessProcess=false // при включении и просыпании - 1 раз

// сохранение значения уровня осмысления == стадии развития при произвольном изменении уровня
var saveEvolushnStage=0

// true - текущий режим активации consciousness - субъектиынй (activationType=2)
var isActivationType2=false

var currrentFromNextID=0 // текущий fromNextID в текущем запуске consciousness

// временное сохранеиние цикла осмысления между двумя объективными вызовами consciousness
var saveFromNextIDcurretCicle []int

////////////////////////////////////////////////////////////
/* Главная, постоянно активная с каждым пульсом функция поддержвания информационной среды и произвольности.
Изолированная от непосредственных воздействий и поэтому самостоятельная система оценки и корректировки состояния,
происходит с формированием опыта прозвольности выполнения Правил и их выбора для данных условий (psy_Experience.go),
в частности, Правил относительно самого себя (самосознание).

Начинает работать с if EvolushnStage > 3

По каждой активации дерева автоматизмов запускается система осознания consciousness.

Вид активации - при вызове функции осознания.
activationType == 0 - не бывает такого значея
activationType == 1 - активация ориентировочным рефлексом новой ситуации
activationType == 2 - активация "внутренним" (произвольным) ориентировочным рефлексом

В принципе здесь должны исправляться все лажи ответов предыдущих периодов...

fromNextID - ID MentalNext на котором была запущена переактивация consciousness при обдумывании
*/
func consciousness(activationType int,fromNextID int)(bool) {    return false
	if currrentFromNextID != fromNextID{//сохранеиние цикла осмысления между двумя объективными вызовами consciousness
		saveFromNextIDcurretCicle=append(saveFromNextIDcurretCicle,fromNextID)
	}
	currrentFromNextID=fromNextID

	if !AllowConsciousnessProcess {
		return false
	}
	var stopMentalWork=false
	if activationType == 1 && isActivationType2{// объективная активация
		// нужно прервать выполнение циклов субъективныъ активаций
		stopMentalWork=true
	}

	if activationType == 1 {
		isActivationType2 = false
		// посмотреть, есть ли прерванные цепочки осмысления, и если есть, выбрать, с какой продолжить осмысление.
	}
	if activationType == 2 {
		isActivationType2 = true
	}

	// вернуть стадию развития
	if saveEvolushnStage >0 {// иначе оно обнуляет EvolushnStage
		EvolushnStage = saveEvolushnStage // возвращаем уровень осмысления
	}
	///////////////////////////////////////////////
	if activationType == 1 {
		refreshCurrentInformationEnvironment()
		saveFromNextIDcurretCicle = nil
	}


	/* Период ожидания ответа LastRunAutomatizmPulsCount при поочередном Стимуле-Ответе есть всегда.
	   А здсь - поиск Ответа именно после каждого Стимула. Так что LastRunAutomatizmPulsCount в функции не учитываем.
	*/
if false {//для тестирования

	////////////////////////////// 1 уровень ////////////////////
	// ПЕРВЫЙ УРОВЕНЬ, самый примитивный уровень:
	// есть ли штатный мот.автоматизм и нужно ли его менять или задумываться
	if fromNextID == 0 {
		nArr := GetMotorsAutomatizmListFromTreeId(detectedActiveLastNodID)
		/* - только привязанные к ветке, а не мягкий алгоритм
		      currentAutomatizmAfterTreeActivatedID = getAutomatizmFromNodeID(detectedActiveLastNodID)
		   Если нет привязанных - находим решение - по сдедующим уровням привлечения осознания.
		*/
		if nArr != nil {
			for i := 0; i < len(nArr); i++ {
				if nArr[i].Belief == 2 && nArr[i].Usefulness > 0 { // есть нормальный, пусть выполняется
					//currentAutomatizmAfterTreeActivatedID=nArr[i].ID
					return false
				}
			}
		}
	}
	// нет штатного автоматизма

	// return false // для набора Объективных правил перед использованием 2-го уровня.

	/////////////////////////////////////////////////////////


	// if CurrentInformationEnvironment.veryActualSituation || CurrentInformationEnvironment.danger{

	//////////////////////////////// 2 уровень ////////////////////////////
	// ВТОРОЙ УРОВЕНЬ - попытка использования примитивных Правил
	if fromNextID == 0 {
		rules := getSuitableRules(activationType)
		if rules > 0 { // по правилу найти автоматизм и запустить его
			ta := TriggerAndActionArr[rules]
			purpose := getPurposeGenetic()
			// выбираем Ответное действие из Правила чтобы повторить его
			ai := ActionsImageArr[ta.Action]
			if ai != nil {
				purpose.actionID = ai
				atmzm := createAndRunAutomatizmFromPurpose(purpose)
				if atmzm != nil {
					return true // заблокирвать более низкоуровневое
				}
			}
		}
	}
	/////////////////////////////////////////////////////////
}// конец блокирования для тестирования

	//////////////////////////////// 3 уровень ////////////////////////////
	//if EvolushnStage > 3 { - уже обеспечено
	/* ТРЕТИЙ УРОВЕНЬ - попытка найти решение, используя всю текущую инфрмацию с учетом срочности.
	      Ментальные автоматизмы нужны только если нет мот.автоматизма или его нужно переделать.
	   Т.е. привязанный к ветке дерева понимания мент автоматизм должен срабатывать ЗДЕСЬ.
	   Он продолжается по NextID или даже ветвится в зависимости от ситуации.
	   А если его еще нет, то создать БАЗОВЫЙ.
	Работа третьего уровня.
	   Запуск ментального автоматизма сопровождается перезапуском consciousness() кроме запуска моторного автоматизма.

	У ветки UnderstandingNode всегда должен быть Базовый ментальный автоматизм, с которого начинается
	просмотр в функции consciousness() на ее Третьем уровне. От него может идти цепочка дочерних.

	Каждый запускаемый мент. автоматизм (кроме MentalActionsImages.activateMotorID) после отработки
	вызывает consciousness() прямо или косвенно. В течение одного пульса может быть множество
	перезапусков consciousness() с продолжением процесса мышления и добавления в цепь (.NextID) новых
	автоматизмов.

	Базовый автоматизм должен прикинуть, какой будет следующий – путем выбора MentalActionsImages,
	сделать его, запустить, а в следующем цикле consciousness() использовать инфу (и все окружение) для
	формирования моторного автоматизма (MentalActionsImages.activateMotorID)  и тогда запустить его с
	периодом ожидания.

	По результату записывается Правило. Если хорошо, то данная цепочка так и заканчивается запуском
	моторного автоматизма (MentalActionsImages.activateMotorID), если плохо – формируется следующая
	цепочка (.NextID) с выбором другого MentalActionsImages и т.д. Т.е. формирование следующего звена
	цепочки идет С УЧЕТОМ ОПЫТА (MentalActionsImages) ПРЕДЫДУЩИХ.

	Структура мент Правила MentalTriggerAndAction начинается или с мент.действия MentalActionsImage
	или с моторного ActionsImage, потом - Ответ MentalActionsImage и обычный Эффект.

	 После срабатывания инфо-функции (.activateInfoFunc) информация добавляется к текущему информационному окружения
	в виде глобальной структурц и задается значение глобальной переменной currentInfoStructId == ID инфо-функции),
	которые могут использоваться при запуске consciousness().
	*/
	if stopMentalWork{
		// запомнить текущую работу в момент ее прерывания, чтобы можно было вернуться к обдумыванию
		// TODO стек 7 прерванных дел
	}
	///////////////////////////////////////////



	if fromNextID==0 {

		// детекция ленивого состояния
		if isIdleness() {// ЛЕНЬ
			// пофиг все, можно лениво обрабатывать накопившиеся структуры, реагирование - на уровне - до EvolushnStage < 4
			saveEvolushnStage = EvolushnStage
			EvolushnStage = 3 // нагло и просто :) - произвольный откат уровня осознания

			processingFreeState() // обработка структур в свободном состоянии может быть долгой -

			EvolushnStage = saveEvolushnStage // возвращаем уровень осмысления, иначе зависнет на этой стадии
			return false                      // пусть выполняется все менее высокоуровневое
		}//if isIdleness()
		/////////////////////////  НЕТ ЛЕНИ



// НАЙТИ fromNextID = определяет ID звена цепи для продолжения осмысления до нахождения удачного моторного автоматизма.
		/* Есть ли ментальный автоматизм для ветки UnderstandingNode? Если нет, то сформировать базовый.
		   Если есть с .activateMotorID>0 - запустить моторный автоматизм activateMotorID,
		   если нет .activateMotorID==0 - создать следующий цепоцечный мент.автоматизм и прописать его ID в .NextID
		*/
		firstArr:=goNextFromUnderstandingNodeIDArr[detectedActiveLastUnderstandingNodID]
		if firstArr == nil {// нет веток у данного узла дерева - создать базовое звено цепочки

			// создание или использование ментального автоматизма инфо-функции №1
			fromNextID=createNexusFromNextID(fromNextID,1)
			if fromNextID>0 {
				// перезапуск осмысления
				return reloadConsciousness(stopMentalWork, fromNextID)
			}

			//if firstArr != nil
		}else{// есть ветки у данного узла дерева
			// найти конец подходящей цепочки
			fromNextID:=getLastAutomatizmFrom(detectedActiveLastUnderstandingNodID,detectedActiveLastNodID)
			if fromNextID>0 {
				// перезапуск осмысления
				return reloadConsciousness(stopMentalWork, fromNextID)
			}
		}// конец: есть ветки у данного узла дерева
	}/// Конец поиска fromNextID
	///////////////////////////////////



// НАЙДЕНО fromNextID
	if fromNextID >0 {// продолжение осмысления по цепочке goNext.ID == fromNextID
		// добавить в кратковреемнную память
		addShortTermMemory(fromNextID)

		// если нужно, учесть: ID последний запуск инфо-функции: switch currentInfoStructId{

		// fromNextID определяет звено цепи (другого не может быть после перезапуска с fromNextID)
		mImgID:=goNextFromIDArr[fromNextID].AutomatizmID
		if mImgID>0{
			// ментальный автоматизм звена
			matmzm:=MentalAutomatizmsFromID[mentalInfoStruct.mImgID]
			if matmzm != nil {
				// если акция - моторный автоматизм или переактивация состояния, то запустить
				mentAct := MentalActionsImagesArr[matmzm.ActionsImageID]
				if mentAct.activateInfoFunc == 0 {// если не инфо-функция
					// запуск моторного автоматизма  или объектиный перезапуск (через переактивированные деревья)
					RunMentalMentalAutomatizm(matmzm)
/* выход из рекурсивного цикла в understanding_tree.go в блок if detectedActiveLastUnderstandingNodID>0{
Если моторный окажется успешным, то он будет записан штатным для ветки detectedActiveLastUnderstandingNodID.
Тут же будут записаны не только моторные, но и ментальные Правила
TODO записать ментальные Правила
*/
					return true
				}// не моторный автоматизм и не объяктивный перезапуск, а ИНФО_ФУНКЦИЯ
				// запустить инфо-функцию
				RunMentalMentalAutomatizm(matmzm)
				// обработать результат

				// получить следующее звено с учетом результата
				fromNextID=calcNexusFromNextID(fromNextID)
				if fromNextID>0 {
					// перезапуск осмысления
					return reloadConsciousness(stopMentalWork, fromNextID)
				}
			}
		}
		//  не случилось ничего полезного и это конец цепочки - нужно доращивать

		// создание или использование ментального автоматизма инфо-функции №2
		fromNextID=createNexusFromNextID(fromNextID,2)
		if fromNextID>0 {
			// перезапуск осмысления с обновлением currrentFromNextID
			return reloadConsciousness(stopMentalWork, fromNextID)
		}
	///////////////////

}// конец фазы осмысления по цепоске goNext.ID == fromNextID
//////////////////////////////////////




//////////////////////////////// 4 уровень ////////////////////////////
if fromNextID==0 {
if EvolushnStage > 4 {
// ЧЕТВЕРТЫЙ УРОВЕНЬ - доминанта нерешенной проблемы - только если нет срочности
// и тут - Стек отложенных дел.
if CurrentInformationEnvironment.veryActualSituation || CurrentInformationEnvironment.danger {
// TODO
} else { // нет решения, паника, откатиться на прежний уровень регирования
// TODO аварийное решение проблемы
return false // пусть выполняется все менее высокоуровневое
}

} //if EvolushnStage > 4
}
/////////////////////////////////////////////////////////







/////////////////////////////////////////////////////////



return false // не блокировать последующий код ориентировочного рефлекса.
}
////////////////////////////////////////////


// перезапустить осмысление
func reloadConsciousness(stopMentalWork bool,fromNextID int)(bool){
	if stopMentalWork{
		// не запускать consciousness(2,fromNextID)
		// TODO сохранить fromNextID прерванного процесса в стеке 7 прерванных процессов

		return false
	}else {
		// перезапуск осмысления
		consciousness(2, fromNextID)
	}
	return true
}
///////////////////////////////////////////////////


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
	if CurrentInformationEnvironment.veryActualSituation || CurrentInformationEnvironment.danger {
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




///////////////////////////////////////////////////////
// создание или использование ментального автоматизма инфо-функции c ID= infoID
func createNexusFromNextID(fromNextID int,infoID int)(int){

imgID,_:=CreateNewlastMentalActionsImagesID(0,0,0,infoID,0)
if imgID>0 {
	aID, _ := createMentalAutomatizmID(0, imgID, 1)
	if aID > 0 {
		// создание звена цепочки (всегда уникальное звено)
		fromNextID0 := fromNextID
		// создание нового элемента цепочки
		fromNextID, _ = createNewNextFromUnderstandingNodeID(
			detectedActiveLastUnderstandingNodID,
			detectedActiveLastNodID,
			fromNextID0,
			aID)

		//addShortTermMemory(fromNextID)

		return fromNextID
	}
}

return 0
}
//////////////////////////////////////////////////


/* обработать результат инфо-фукнции и прикинуть, что дальше с fromNextID
 */
func calcNexusFromNextID(fromNextID int)(int){
	switch currentInfoStructId{
	// если это - самая первая ф-ция, то вызвать 2-ю ф-цию поиска продолжения
	case 1: return createNexusFromNextID(fromNextID,2)
	// анализ инфо стркутуры по currentInfoStructId и выдача решения
	case 2: return analisAndSintez(fromNextID)
	}

	return 0
}
//////////////////////////////////////////////////////



//////////////////////////////////////////////////////
/* после периода ожидания
Учесть последствия запуска мот.автоматизма,
из saveFromNextIDcurretCicle выявить Правила и записать в эпизод.пямять.
Если нужно - обдумать результат в новом ментальном consciousness(2,currrentFromNextID)
Сразу после обработки периода ожидания запускается дерево понимания и объявный запуск consciousness(1,0)
так что никаких действий совершать в afterWaitingPeriod() или в самой afterWaitingPeriod() не следует.

МОЖНО ДОСТАТОЧНО БЫСТРО СДЕЛАТЬ ПОСЛЕ ОБКАТКИ СИСТЕМЫ МЕНТАЛЬНЫХ ЦИКЛОВ
 */
func afterWaitingPeriod(){

	// TODO из saveFromNextIDcurretCicle выявить Правила и записать в эпизод.пямять и т.д.

//??  consciousness(2,currrentFromNextID)
}
///////////////////////////////////////////////////////


