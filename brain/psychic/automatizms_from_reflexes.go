/* Создать автоматизмы на основе существующих рефлексов

Для тестирования возможно избежать долгий период воспитания с формированием автоматизмов и просто сгенерировать автоматизмы на основе существующих рефлексов (с приоритетом условных рефлексов).
При этом у автоматизмов будут установлены опции уже проверенного автоматизма с полезностью, равной 1 (вполне полезно). Это правомерно потому, что рефлексы создавались уже полезными для своих условий.
В дальнейшем такие автоматизмы будут проверяться в зависимости от реакции оператора и изменения состояния Beast, корректируясь настолько эффективно, насколько позволяет текущая стадия развития.
*/

package psychic

import (
	"BOT/lib"
	"strconv"
	"strings"
)

///////////////////////////////////////////////
type autNode struct { // узел дерева автоматизмов
	ID int
	baseID int
	EmotionID int
	ActivityID int
	ToneMoodID int
	SimbolID int
	VerbalID int
	Children []AutomatizmNode // дочерние узлы (ветвление) НЕ АДРЕСА, А РЕАЛЬНЫЕ ОБЪЕКТЫ
	ParentID int     // ID родителя
	ParentNode *autNode
}
var autTreeFromID=make(map[int]*autNode)



////////////////////////////////////////////
/* сканировать для всех условий у.рефлексы (и если нет  - смотреть б.рефлексы),
создавать ветку дерева автоматизма если такой еще нет,
создавать автоматизм, прикрепляя его к нужно ветке.
 */
func RunMakeAutomatizmsFromReflexes()(string){
var newCount=0

// считать дерево рефлексов
	strArr,_:=lib.ReadLines(lib.GetMainPathExeFile()+"/memory_psy/automatizm_tree.txt")
	cunt:=len(strArr)
	//просто проход по всем строкам файла подряд так что сначала идут дочки, потом - их родители
	for n := 0; n < cunt; n++ {
		if len(strArr[n])==0{
			continue
		}
		p:=strings.Split(strArr[n], "|")
		id,_:=strconv.Atoi(p[0])
		parentID,_:=strconv.Atoi(p[1])
		baseID,_:=strconv.Atoi(p[2])
		EmotionID,_:=strconv.Atoi(p[3])
		ActivityID,_:=strconv.Atoi(p[4])
		ToneMoodID,_:=strconv.Atoi(p[5])
		SimbolID,_:=strconv.Atoi(p[6])
		VerbalID,_:=strconv.Atoi(p[7])
		setNewAutNode(id,parentID,baseID,EmotionID,
			ActivityID,ToneMoodID,SimbolID,VerbalID)
	}

//............ считать рефлексы
	path := lib.GetMainPathExeFile()
	lines, _ := lib.ReadLines(path + "/memory_reflex/condition_reflexes.txt")
	for i := 0; i < len(lines); i++ {
		if len(lines[i]) < 4 {
			continue
		}
		p := strings.Split(lines[i], "|")
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
// есть ли такие условия в дереве автоматизмов
if lev1>0 && lev2!=nil && lev3>0{

}



	}


///////////// убрать мусор
	strArr=nil
	autTreeFromID=nil
return "Процесс нормально завершен, создано "+strconv.Itoa(newCount)+" новых автоматизмов."
}
/////////////////////////////////////////


//
var idCount=1
func setNewAutNode(id int,parentID int,baseID int,EmotionID int,
	ActivityID int,ToneMoodID int,SimbolID int,VerbalID int)(*autNode){
	var node autNode
	node.ID = id
	idCount++
	//node.ParentNode=parent
	node.ParentID=parentID
	node.EmotionID=EmotionID
	node.ActivityID=ActivityID
	node.ToneMoodID=ToneMoodID
	node.SimbolID=SimbolID
	node.VerbalID=VerbalID

	autTreeFromID[node.ID]=&node

	return &node
}
////////////////////////////////////////////////////

