/*   функции использования текущих условий гомесотаза


*/


package gomeostas

import "sort"

/////////////////////////////////////////////


/* выявить ID парамктров гомеостаза как цели для улучшения в данных условиях
Возвращает PurposeGenetic.veryActual  gomeostas.FindTargetGomeostazID
сортировка по уменьшению важности
*/
func FindTargetGomeostazID()(bool,[]int){
	var veryActual=false
	var idArr []int
// BadNormalWell - состояние каждого параметра гомеостаза: 1 - Похо, 2 - Норма, 3 - Хорошо
// отсортировать по убыли важности GomeostazParamsWeight
	badNormalWellImp:=sortingForImpotents()

	for k, pID := range badNormalWellImp {
		if pID==1 { // плохо для данного параметра гомеостаза
			idArr = append(idArr, k)
			if k == 1 || k == 2 || k == 7 || k == 8 {
				veryActual = true
			}
		}
	}
	return veryActual,idArr
}
/////////////////////////////////
func sortingForImpotents()(map[int]int){
	var impC=make(map[int]int)
	for id, _ := range BadNormalWell {
		impC[GomeostazParamsWeight[id]]=id
	}

	vals := make([]int, 0, len(impC))
	for k := range impC {
		vals = append(vals, k)
	}
	//СОРТИРОВКА ПО УБЫВАНИЮ
	sort.Slice(vals , func(i, j int) bool {
		return vals[i] > vals[j]
	})

	var arr=make(map[int]int)
	for _,v := range vals {
			arr[impC[v]]=BadNormalWell[impC[v]]
		}

return arr
}
//////////////////////////////////////////////


/////////////////////////////////////////////
/* в каком из 5 диапазоне нормы находится Базовый параметр
0 - это не норма
1 Норма 0-19%
2 Норма 20-39%
3 Норма 40-59%
4 Норма 60-79%
5 Норма 80-100%
 */
func getNormaDiapason(pID int)(int){
	gp:=int(GomeostazParams[pID])
limit:=compareLimites[pID]// порог начала критического выхода параметров из нормы
if pID==1 && gp <= limit{ 	return 0 }// для энергии - наоборот
if pID>1 && gp >= limit{ 	return 0 }
// для нормы
var norm=0
	if pID==1 {
		norm = 100 - limit      // остаток параметра вне критического
		gp = gp-limit // убираем критическую часть
	}else{
		norm = limit            // остаток параметра вне критического
	}

// какой процент составляет gp от norm
proc:=int((gp*100)/norm)
	if proc <20 { return 1 }
	if proc <40 { return 2 }
	if proc <60 { return 3 }
	if proc <80 { return 4 }
return 5
}
//////////////////////////////////////////



