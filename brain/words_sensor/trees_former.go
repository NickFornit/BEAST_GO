/*  Формирователь Дерева слов и Дерева фраз при наборе текстов
из редактора http://go/pages/words.php
и при общении с Beast с Пульта http://go/pult.php

*/

package word_sensor

import (
	"BOT/lib"
	"sort"
	"strconv"
	"strings"
)
/////////////////////////////////////////////////////////////////////////

/* переносим в дерево слов достаточно повторяющиеся из tempArr
 limitWord - число повторений отдельных слов, с которых начинается пеернос.
limitFraze - число повторений отдельных фраз (с несколькими словами), с которых начинается пеернос.
Для автозаливки - (4,6),
для сообщения с Пульта - (2,4)
*/
func updateWordTreeFromTempArr(limitWord int,limitFraze int){
	//Удалять использованные строки накопительного массива
var newTempFileStr="" // это будет записано после процесса в /memory_reflex/words_temp_arr.txt

var str[]string // для выбранных повторяющихся слов - чтобы потом отсортирвать
	// чтобы фразы могли использовать слова (минимизация размера дерева)
		// слово или фраза?
	var existsWords=false
	// сначала прходим только слова, потом - только фразы,
	for k, v := range tempArr {
		sps:=strings.Split(k, " ")
		if len(sps)==1{// это слово
			if v >= limitWord {
				str = append(str, k)
			}else{
				newTempFileStr+=strconv.Itoa(v)+"|#|"+k+"\r\n"
			}
		}
	}
	sort.Strings(str)  // по алфавиту, чтобы максимально облегчить последовательное разделение слов
	for i := 0; i < len(str); i++ {
		cur:=str[i]
		SetNewWordTreeNode(cur)
		existsWords=true
//		SaveWordTree() // для пошагового контроля
//		if(i>1){}
	}
	if existsWords {
		SaveWordTree()
	}

	// проход для фраз
	var existsPhrase=false
	str=nil
	for k, v := range tempArr {
		sps:=strings.Split(k, " ")
		if len(sps)>1{// это фраза
			if v >= limitFraze {
				str = append(str, k)
			}else{
				newTempFileStr+=strconv.Itoa(v)+"|#|"+k+"\r\n"
			}
		}
	}
	sort.Strings(str)  // по алфавиту, чтобы максимально облегчить последовательное разделение слов
	for i := 0; i < len(str); i++ {
		cur:=str[i]
		SetNewPhraseTreeNode(cur)
		existsPhrase=true
		//		SaveWordTree() // для пошагового контроля
		//		if(i>1){}
	}
	if existsPhrase {
		SavePhraseTree()
	}

	// Удаление использованных строк накопительного массива: запись только незатронутых
	lib.WriteFileContent(lib.GetMainPathExeFile()+"/memory_reflex/words_temp_arr.txt",newTempFileStr)
}
/////////////////////////////////////////////////

