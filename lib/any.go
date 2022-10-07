package lib

import (
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

/* разные вспомогательные функции */

// удалить повторяющиеся значения
func RemoveDuplicateValues(intSlice []int) []int {
	keys := make(map[int]bool)
	list := []int{}

	// If the key(values of the slice) is not equal
	// to the already present value in new slice (list)
	// then we append it. else we jump on another element.
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

// случайное из диапазона чисел
func RandChooseDiapazonQuest(first int, last int) int {
	rand.Seed(time.Now().UnixNano())
	n := first + rand.Intn(last-first+1)
	return n
}

// Round возвращает ближайшее целочисленное значение.
func Round(x float64) int {
	t := math.Trunc(x)
	if math.Abs(x-t) >= 0.5 {
		return int(t + math.Copysign(1, x))
	}
	return int(t)
}

// случайное из нескольких чисел
func RandChooseQuest(nums ...int) int {
	var count = 0
	var num []int
	for _, n := range nums {
		num = append(num, n)
		count++
	}
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(count + 1)
	return num[n]
}

// случайное из массива
func RandChooseIntArr(n []int) int {
	var count = 0
	var num []int
	for _, n := range n {
		num = append(num, n)
		count++
	}
	rand.Seed(time.Now().UnixNano())
	v := rand.Intn(count)
	return num[v]
}

// сравнение массивовв int на идентичность
func EqualArrs(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

// вернуть наиболее повторяющееся по числу значение массива и среди них первое из равных
func GetMaxCountVal(arr []int) int {
	// группируем значения в a
	var a = make(map[int]int)
	for i := 0; i < len(arr); i++ {
		a[arr[i]]++
	}
	// находим максимальное или первое из равных
	var max = 0
	var ind = 0
	for k, v := range a {
		if max < v {
			max = v
			ind = k
		}
	}
	return ind
}

// индекс значения в массиве
func IndexValInArr(arr []int, val int) (bool, int) {
	for k, n := range arr {
		if n == val {
			return true, k
		}
	}
	return false, 0
}

// есть такое значение в массиве
func ExistsValInArr(arr []int, val int) bool {
	for _, n := range arr {
		if n == val {
			return true
		}
	}
	return false
}

// есть такое значение в массиве с учетом его сортировки
func ExistsValInArrSort(arr []int, val int) bool {
	for i:=0; i<len(arr); i++{
		if arr[i] == val { return true }
	}
	return false
}

// есть полное строковое значение в строке значений
func ExistsValStrInList(List string, val string, razd string) bool {
	var sArr []string

	sArr = strings.Split(List, razd)
	for _, n := range sArr {
		if n == val {
			return true
		}
	}
	return false
}

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// числа имеют разные знаки (одно положительное, другое отрицательное)
func IsDiffersOfSign(n1 int, n2 int) bool {
	if (n1 > 0 && n2 < 0) || (n1 < 0 && n2 > 0) {
		return true
	}
	return false
}

func AbsFloate(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

func Max(x, y int) int {
	if x < y {
		return y
	}
	return x
}

func Min(x, y int) int {
	if x > y {
		return y
	}
	return x
}

// сохранить массив int в переменной
func SaveArrToVar(arr []int, to []int) []int {
	for a := 0; a < len(arr); a++ {
		to = append(to, arr[a])
	}
	return to
}

// текстовый массив в числовой
func IntArrToStrArr(sArr []string) []int {
	var out []int
	var id int

	for i := 0; i < len(sArr); i++ {
		id, _ = strconv.Atoi(sArr[i])
		if id != 0 {
			out = append(out, id)
		}
	}

	return out
}

// числовой массив в текстовый
func StrArrToIntArr(sArr []int) []string {
	var out []string
	var id string

	for i := 0; i < len(sArr); i++ {
		id = strconv.Itoa(sArr[i])
		out = append(out, id)
	}

	return out
}

// объединить 2 массива в один
func SumArr(arr1 []int, arr2 []int) []int {
	var out = append(arr1, arr2...)
	return UniqueArr(out)
}

// убрать дублеры в массиве
func UniqueArr(intSlice []int) []int {
	keys := make(map[int]bool)
	list := []int{}
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

/* найти какие значения m1 есть в m2
т.е. m2 должен быть частью m1 или полностью совпадать с ним
*/
func GetExistsIntArs(m1 []int, m2 []int) []int {
	var found []int
	for i := 0; i < len(m1); i++ {
		for j := 0; j < len(m2); j++ {
			if m1[i] == m2[j] {
				found = append(found, m1[i])
			}
		}
	}
	return found
}

/* найти каких значений m1 нет в m2
т.е. m2 должен быть частью m1 или полностью совпадать с ним
*/
func GetDifferentIntArs(m1 []int, m2 []int) []int {
	var diff []int
	var exists = false
	for i := 0; i < len(m1); i++ {
		exists = false
		for j := 0; j < len(m2); j++ {
			if m1[i] == m2[j] {
				exists = true
				break
			}
			if !exists {
				diff = append(diff, m1[i])
			}
		}
	}
	return diff
}
