/* Доминанта нерешенной проблемы
Это - атрибут 5-й стадии развития - творчества
*/

package psychic

///////////////////////////////////////
type Dominanta struct {
	ID int
	// проблема, которую нужно решить - цель доминанты к которой стремятся ментальные автоматизмы как своей успешности.
	// TODO problemsID int // структура проблемы, с параметром цели решения

	lastMentalID int // ментальный автоматизм, на котоором закончилось прошлое решение проблемы

	birthTime int
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