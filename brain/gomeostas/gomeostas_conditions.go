/* функции использования текущих условий гомесотаза
*/

package gomeostas

import "sort"

/* выявить ID параметров гомеостаза как цели для улучшения в данных условиях
Возвращает PurposeGenetic.veryActual gomeostas.FindTargetGomeostazID
сортировка по уменьшению важности
*/
func FindTargetGomeostazID()(bool,[]int) {
	var veryActual = false
	var idArr []int
	// BadNormalWell - состояние каждого параметра гомеостаза: 1 - Похо, 2 - Норма, 3 - Хорошо
	// отсортировать по убыли важности GomeostazParamsWeight
	badNormalWellImp := sortingForImpotents()

	for k, pID := range badNormalWellImp {
		if pID == 1 { // плохо для данного параметра гомеостаза
			idArr = append(idArr, k)
			if k == 1 || k == 2 || k == 7 || k == 8 {
				veryActual = true
			}
		}
	}
	return veryActual, idArr
}

// Сортировка ID параметров гомеостаза по убыванию значимости GomeostazParamsWeight
func sortingForImpotents() map[int]int {
	var impC = make(map[int]int)
	for id, _ := range BadNormalWell {
		impC[GomeostazParamsWeight[id]] = id
	}

	vals := make([]int, 0, len(impC))
	for k := range impC {
		vals = append(vals, k)
	}
	// СОРТИРОВКА ПО УБЫВАНИЮ
	sort.Slice(vals, func(i, j int) bool {
		return vals[i] > vals[j]
	})

	var arr = make(map[int]int)
	for _,v := range vals {
		arr[impC[v]] = BadNormalWell[impC[v]]
	}

	return arr
}

/* в каком из 5 диапазоне Плохо находится Жизненный параметр
1 Плохо 0-19%
2 Плохо 20-39%
3 Плохо 40-59%
4 Плохо 60-79%
5 Плохо 80-100%
 */
func getBadDiapazon(pID int) int {
	gp := int(GomeostazParams[pID])
	limit := compareLimites[pID] // порог начала критического выхода параметров из нормы
	if pID == 1 && gp >= limit { return 0 } // для энергии - наоборот
	if pID > 1 && gp <= limit { return 0 }
	// для Плохо
	var bad = 0
	if pID == 1 {
		bad = limit // остаток параметра вне Норма
	} else {
		bad = 100 - limit // остаток параметра вне Норма
		gp = gp - limit // убираем Норма
	}
	// какой процент составляет gp от Плохо
	proc := int((gp * 100) / bad)
	if proc < 20 { return 1 }
	if proc < 40 { return 2 }
	if proc < 60 { return 3 }
	if proc < 80 { return 4 }

	return 5
}