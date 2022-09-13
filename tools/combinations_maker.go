/*
Генерация сочетаний для редактора б.рефлексов.


*/

package tools

import (
	"BOT/lib"
	"strconv"
)

func init(){


}
/////////////////////////////////////////////////////
// найти все сочетания ряда чисел от 0 до максимального подряд (без перемешивания
// limit - максимальное число элементов в сочетании
func GetAllCombinationsOfSeriesNumbers(max int,limit int)([][]int){
	var out [][]int
	for i := 0; i < max; i++ {
		var arr []int
		end:=limit+i
		if i > max-limit{
			end=max
		}
		for j := i; j < end; j++ {
			arr=append(arr,j)
		}
		out=append(out,arr)
		// подсочетания
		var arr0 []int
		for n := 2; n < len(arr); n++ {
			for m := 0; m < n; m++ {
				arr0=append(arr0,arr[m])
			}
			out=append(out,arr0)
		}
	}
	return out
}
////////////////////////////////////////////////////


// найти все сочетания ряда чисел от 0 до максимального без перемешивания
// limit - максимальное число элементов в сочетании
func GetAllCombinationsNumbers(max int,limit int)([][]int){
	var out [][]int
	var series []string
	for i := 0; i < max; i++ {
		series=append(series,strconv.Itoa(i))
	}
	comb:=Combinations(series,limit)
	for _, v := range comb {
		var arr []int
		for n := 0; n < len(v); n++ {
			nun,_:=strconv.Atoi(v[n])
			arr=append(arr,nun)
		}
		out=append(out,arr)
	}
	return out
}
////////////////////////////////////////////////////



func MakeContextCombinations() {
	var cellArr []string
	for i := 1; i < 18; i++ {
			cellArr=append(cellArr,strconv.Itoa(i))
	}
	arr:=Combinations(cellArr,1)
	//arr:=All(cellArr)

	out:=""

	for _, v := range arr {
		for i := 0; i < len(v); i++ {
			out+=v[i]+","
		}
		out+="\r\n"
	}

	lib.WriteFileContent(lib.GetMainPathExeFile()+"/combination.txt", out)
	if len(out)>0{

	}
	return
}

var nid []int
var outArr []string
var cols int
var gCArr=make(map[int][]string)
// сейчас вызов из main.init()
/*
func MakeContextCombinations() {

	gCArr = gomeostas.GomeostazActivnostArr

	cols = len(gCArr[1])
//	rows := len(gCArr)
//	nid = make([]int, rows)
	arr1 := gCArr[1];
	arr2 := gCArr[2];
	arr3 := gCArr[3];
	arr4 := gCArr[4];
	arr5 := gCArr[5];
	arr6 := gCArr[6];
	arr7 := gCArr[7];
	arr8 := gCArr[8];
	var out = "";
	for n1:= 0; n1 < cols; n1++ {
		for n2:= 0; n2 < cols; n2++ {
			for n3:= 0; n3 < cols; n3++ {
				for n4:= 0; n4 < cols; n4++ {
					for n5:= 0; n5 < cols; n5++ {
						for n6:= 0; n6 < cols; n6++ {
							for n7:= 0; n7 < cols; n7++ {
								for n8:= 0; n8 < cols; n8++ {
out +=arr1[n1]+"|"+arr2[n2]+"|"+arr3[n3]+"|"+arr4[n4]+"|"+arr5[n5]+"|"+arr6[n6]+"|"+arr7[n7]+"|"+arr8[n8]
								}
							}
						}
					}
				}
			}
		}
	}

	// удаляем повторы из каждого сочетания
	var comb = make([][]int, len(outArr))
	for k, v := range outArr {
		p := strings.Split(v, ",")
		for i := 0; i < len(p); i++ {
			id, _ := strconv.Atoi(p[i])
			comb[k] = append(comb[k], id)
		}
		comb[k] = lib.RemoveDuplicateValues(comb[k])
		sort.Ints(comb[k])
	}


	return
}
 */


/////////////////////////////////////////////////////

/* // перебор с использованием проф. библиотеки
var cellArr []string
cols:=len(gCArr[1])
rows:=len(gCArr)
	for i := 1; i < rows; i++ {
		for j := 0; j < cols; j++ {
			cellArr=append(cellArr,gCArr[i][j])
			//cellArr=append(cellArr,strconv.Itoa(i*j))
		}
	}

	// Время выполнения для перебора 49  - очень большое...
	out:=All(cellArr)
	//out:=Combinations(cellArr,3)
	if len(out)>0{

	}

return
}
 */