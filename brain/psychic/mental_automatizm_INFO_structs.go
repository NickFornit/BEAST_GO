/* Структуры информации для инфо-функций

Прототипы:
type infoStruct4 struct {
	par int
}
var info4 infoStruct4
func infoFunc4(){
	res:=0

	info4.par = res // передача инфы в структуру
	currentInfoStructId=4 // определение актуальной инфо-структуры
}
*/

package psychic

//////////////////////////////////////////
/* Подобрать MentalActionsImages
случайно или по заготовке редактора с Пульта
 */
type infoStruct0 struct {
	mImgID int // ID MentalActionsImages
}
var info0 infoStruct0
