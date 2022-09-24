/*
файловые функции

ReadLines
ReadIntArr
ReadFloate64Arr
WriteNewString
RewriteFileContent
WriteFileContent
ReadFileContent
*/

package lib

import (
	"bufio"
	"fmt"
	_ "fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

// путь к исполняемому файлу
var MainPathExeFile string

func GetMainPathExeFile() string {
	//	mainPathExeFile=os.Args[0] - путь с самим файлом
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	//	mainPathExeFile, _ = filepath.EvalSymlinks(ex)  - путь с самим файлом
	MainPathExeFile = filepath.Dir(ex)
	return MainPathExeFile
}

// read line by line into memory
// all file contents is stores in lines[]
func ReadLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	/*  нельзя, т.к. тогда будут срабатывать for n := 0; n < len(lines); n++ {
	if lines == nil{// всегда позволять считывать хотя бы пустое lines[0]
		lines = append(lines, "")
	}
	*/
	return lines, scanner.Err()
}

func ReadIntArr(path string) ([]int, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []int
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		val, _ := strconv.Atoi(scanner.Text())
		lines = append(lines, val)
	}
	return lines, scanner.Err()
}

func ReadFloate64Arr(path string) ([]float64, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []float64
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		val, _ := strconv.ParseFloat(scanner.Text(), 64)
		lines = append(lines, val)
	}
	return lines, scanner.Err()
}

func WriteNewString(file string, str string) {
	f, _ := os.OpenFile(file, os.O_APPEND, 0666)
	f.WriteString(str + "\n")
	//	fmt.Println("Запись: ", err);
	defer f.Close()
}

// перезаписать файл  (должен быть файл)
func RewriteFileContent(file string, content string) {
	f, _ := os.OpenFile(file, os.O_RDWR, 0666)
	f.WriteString(content)
	//	fmt.Println("Запись: ", err);
	defer f.Close()
}

// записать файл, если нет - создать
func WriteFileContent(file string, content string) {
	if len(content) == 0 {
		return
	}
	//	f, _ := os.OpenFile(file,os.O_CREATE, 0666)
	f, _ := os.Create(file)
	f.WriteString(content)
	//	fmt.Println("Запись: ", err);
	defer f.Close()
}

// считывание файла в строку
func ReadFileContent(file string) string {
	//f, _ := os.Open(file)
	data, _ := ioutil.ReadFile(file)
	//defer f.Close()
	return string(data)
}

// копировать файл
func CopyFile(sourceFile string, destinationFile string) {
	input, err := ioutil.ReadFile(sourceFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = ioutil.WriteFile(destinationFile, input, 0644)
	if err != nil {
		fmt.Println("Error creating", destinationFile)
		fmt.Println(err)
		return
	}
}
