/* Доминанта нерешенной проблемы
Это - атрибут 5-й стадии развития - творчества
*/

package psychic

///////////////////////////////////////
type Dominanta struct {
	ID int
	weight int // вес значимости проблемы
	emotionID int // эмоциональный контекст проблемы

	/* проблема, которую нужно решить - цель доминанты к которой стремятся ментальные автоматизмы как своей успешности.
	// TODO problemsID int // структура проблемы, с параметром цели решения
!!! Эта проблема должна быть совместима с ментальным использованием: mentalInfoStruct.mentalPurposeID=
	 */

	lastGoNextID int // ID структуры goNext, на котоором закончилось прошлое решение проблемы

	// TODO: Обновить время рождения Доминант  /lib/set_life_time.php при установки нового времени жизни
	birthTime int // время рождения проблемы в днях жизни (как в условных рефлексах)
	/* // степень решения проблемы:
	0 - ничего нет,
	1 - начато решение,
	2 - частичное решение,
	3 - успешное решение - гельштат закрыт, можно не думать, а пользолваться опытом
	*/
	isSuccess int
}
var DominantaProblem=make(map[int]*Dominanta)

var isCurrentProblemDominanta *Dominanta


/////////////////////////////////  if doWritingFile{SaveProblemDominenta() }
func SaveProblemDominenta(){

}
/////////////////////////////////////////


// степень решенной Доминанты (зактыть гештальт isSuccess=3 )
func solutionDominanta(dID int,isSuccess int){
	DominantaProblem[dID].isSuccess=isSuccess
}
/////////////////////////////////////////////


///////////////////////////////////
// наиболее важная доминанта в заданном эмоциональном контексте
func getMainDominanta(emotionID int) int {
	var dominantaID=0
	// TODO

		switch emotionID {
		case 1: //Пищевой	- Пищевое поведение, восполнение энергии, на что тратится время и тормозятся антагонистические стили поведения.

		case 2: //Поиск	- Поисковое поведение, любопытство. Обследование объекта внимания, поиск новых возможностей.

		case 3: //Игра	- Игровое поведение - отработка опыта в облегченных ситуациях или при обучении.

		case 4: //Гон	- Половое поведение. Тормозятся антагонистические стили

		case 5: //Защита	- Оборонительные поведение для явных признаков угрозы или плохом состоянии.

		case 6: //Лень	- Апатия в благополучном или безысходном состоянии.

		case 7: //Ступор	- Оцепенелость при непреодолимой опастbase_context_activnostности или когда нет мотивации при благополучии или отсуствии любых возможностей для активного поведения.

		case 8: //Страх	- Осторожность при признаках опасной ситуации.

		case 9: //Агрессия	- Агрессивное поведение для признаков легкой добычи или защиты (иногда - при плохом состоянии).

		case 10: //Злость	- Безжалостность в случае низкой оценки .

		case 11: //Доброта	- Альтруистическое поведение.

		case 12: //Сон - Состояние сна. Освобождение стрессового состояния. Реконструкция необработанной информации.

		}

return dominantaID
}
////////////////////////////////////////////////////


// для пульта
func GetDominantaIDString(id int)string{

	return "Еще не сделан вывод Доминант"
}
////////////////////////////////////////////