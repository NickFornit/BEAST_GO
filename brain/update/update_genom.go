/* Обмен базовыми данными происходит по схеме:
1. Заводится общий каталог, в который выводятся данные ботов в виде структурированных файлов,
которыми они хотят обменяться. Имя файла в формате: имя бота_имя файла. Например:
bot1_update_phrase_tree.txt.
2. В каталоге memory_save создается файл каталога обмена update_dir.txt, в котором указывается список файлов для
обновления:
1|bot2|update_phrase_tree|2022-07-12 09:40:48|1
1 - номер очередности загрузки. ВАЖНО: сначала файл действий, потом рефлексов.
bot2 - имя внешнего бота
update_phrase_tree - имя файла обмена (см. константы ниже)
2022-07-12 09:40:48 - дата/время, заполняется автоматом после успешного обновления
1 - статус блокировки записи: 0 - обмен заблокирован, 1 - обмен разрешен
3. Каждый бот следит за своим файлом, обновляя их. Чужие только читает. При обновлени делается проверка на
совместимость БП и БК, список действий.
Прочий обмен это обмен первичными сенсорами: чтобы бот правильно расположил их по дереву надо ему просто скормить
дерево фраз другого бота.
*/
package update

import (
	"BOT/brain/reflexes"
	termineteAction "BOT/brain/terminete_action"
	"BOT/brain/words_sensor"
	"BOT/lib"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

/* Шаблон дата-время
Может возникнуть справедливый вопрос: а что же такого магического есть в Mon Jan 2 15:04:05 MST 2006?
Давайте посмотрим на шаблон в другом порядке:
01/02 03:04:05PM 06-0700
Видно, что здесь нет двух одинаковых полей. А это означает, что для такой конкретной даты каждое поле будет
точно идентифицированным вне зависимости от форматирования. */
var Layout = "2006-01-02 15:04:05"

// Имя бота
var botName string

// каталог обмена
var pathUpdate string

// имя файла обновления фраз
const updatePhraseName = "update_phrase_tree"

// имя файла обновления рефлекторных действий
const updateActonsName = "update_terminal_actons"

// имя файла обновления рефлексов
const updateDnkReflexes = "update_dnk_reflexes"

// имя файла обновления пусковых стимулов У-рефлекса
const updateTriggerStimulsImages = "update_trigger_stimuls_images"

// имя файла обновления У-рефлексов
const updateConditionReflexes = "update_condition_reflexes"

func init() {
	address := lib.ReadFileContent(lib.GetMainPathExeFile() + "/common/linking_address.txt")
	botName = strings.TrimSpace(address)[7:]
	botName = strings.Split(botName, ":")[0]
	pathUpdate = strings.Replace(lib.GetMainPathExeFile(), "\\"+botName, "\\update", 1)
	LoadFileUpdate()
}

// структура каталога файлов обмена
type FileMod struct {
	ID       int       // ID записи
	FileName string    // имя файла
	LastMod  time.Time // дата-время последнего изменения
	Compat   int       // совместимость файлов обменна: 0 - не совместмый, 1 - совместимый
}

var FileUpdate = make(map[string]string)   // каталог обмена
var FileUpdateDir = make(map[int]*FileMod) // подкаталог обмена

/* Загрузить каталог обмена */
func LoadFileUpdate() {
	var sArr []string
	var tArr []string
	var i, id, MaxNum, Compat int
	var FileName, botNameBuf, idList string
	var LastMod time.Time

	sArr, _ = lib.ReadLines(lib.GetMainPathExeFile() + "/memory_save/update_dir.txt")
	MaxNum = len(sArr)
	if MaxNum == 0 {
		return
	}
	FileUpdateDir = make(map[int]*FileMod, MaxNum)

	for i = 0; i < MaxNum; i++ {
		tArr = strings.Split(sArr[i], "|")
		if botNameBuf != tArr[1] {
			if botNameBuf != "" {
				FileUpdate[botNameBuf] = strings.TrimSuffix(idList, "|")
			}
			botNameBuf = tArr[1]
			idList = ""
		}
		idList += tArr[0] + "|"
		if i == MaxNum-1 {
			idList = strings.TrimSuffix(idList, "|")
		}
		FileUpdate[botNameBuf] = idList
		FileName = tArr[2]
		id, _ = strconv.Atoi(tArr[0])
		LastMod, _ = time.Parse(Layout, tArr[3]) // строка в дату/время
		Compat, _ = strconv.Atoi(tArr[4])
		EditAddFileUpdate(id, FileName, LastMod, Compat)
	}
}

/* Обновить каталог обмена */
func EditAddFileUpdate(ID int, FileName string, LastMod time.Time, Compat int) *FileMod {
	var node FileMod

	node.ID = ID
	node.FileName = FileName
	node.LastMod = LastMod
	node.Compat = Compat
	FileUpdateDir[ID] = &node

	return &node
}

// сохранить изменения в каталоге обмена в файл
func SaveFileUpdateDir() {
	var sArr []string
	var i, id int
	var out, fileName, LastMod, compat string

	for BotName, idList := range FileUpdate {
		sArr = strings.Split(idList, "|")
		for i = 0; i < len(sArr); i++ {
			id, _ = strconv.Atoi(sArr[i])
			fileName = FileUpdateDir[id].FileName
			LastMod = FileUpdateDir[id].LastMod.Format(Layout)
			compat = strconv.Itoa(FileUpdateDir[id].Compat)
			out += strconv.Itoa(id) + "|" + BotName + "|" + fileName + "|" + LastMod + "|" + compat + "\r\n"
		}
	}
	out = strings.TrimSuffix(out, "\r\n")
	lib.WriteFileContent(lib.GetMainPathExeFile()+"/memory_save/update_dir.txt", out)
}

/* Импорт из файла обмена
Загружаем поочередно все типы файлов от ботов, которые прописаны в каталоге
*/
func ImportFileUpdate(NoCheckWordCount bool) bool {
	var sArr []string
	var tArr []string
	var IsMod, IsNoCompat, IsNoAllImport, IsCompat, FlgBreak bool
	var i, j, id, compat, aktIdCur int
	var FileName, updateName, FilePathUpdate, msgTxt, nameAkt string

	// смотрим каталог обмена и закачиваем согласно списку ботов их файлы
	for BotName, ListId := range FileUpdate {
		sArr = strings.Split(ListId, "|")
		for i = 0; i < len(sArr); i++ {
			id, _ = strconv.Atoi(sArr[i])
			compat = FileUpdateDir[id].Compat
			if compat == 0 {
				continue // не совместимый для обмена файл
			}
			IsCompat = true // хотя бы один файл в списке обновлений не блокирован по совместимости
			FileName = FileUpdateDir[id].FileName
			FilePathUpdate = pathUpdate + "/" + BotName + "_" + FileName + ".txt"
			// смотрим дату последнего изменения файла обмена
			file, err := os.Stat(FilePathUpdate)
			if err != nil || file.Size() == 0 {
				lib.WritePultConsol("Файл не найден или пустой: " + BotName + "_" + FileName)
				IsNoAllImport = true
				continue
			}
			// сравниваем дату/время изменения файла с зафиксированной при последнем обновлении
			ModTime := file.ModTime() // считывает с учетом часового пояса!
			// убираем пояс из времени
			ModtimeTxt := ModTime.Format(Layout)
			ModtimeLayout, _ := time.Parse(Layout, ModtimeTxt)
			LastMod := FileUpdateDir[id].LastMod
			dur := LastMod.Sub(ModtimeLayout) // сравнивает с учетом пояса, поэтому его надо убрать заранее!
			if dur != 0 {
				tArr, _ = lib.ReadLines(FilePathUpdate)
				if len(tArr) == 0 {
					continue
				}
				if FileName == updateActonsName {
					ActonsSincID = make(map[int]int)
				}
				if FileName == updateTriggerStimulsImages {
					TriggerSincID = make(map[int]int)
				}
				for j = 0; j < len(tArr); j++ {
					if tArr[j] != "" {
						switch FileName {
						case updatePhraseName: // дерево фраз
							updateName = updatePhraseName
							if NoCheckWordCount == true {
								// включаем флаг авторитарной записи новой фразы, без выдержки в tempArr
								word_sensor.NoCheckWordCount = true
							}
							word_sensor.SetNewPhraseTreeNode(tArr[j])
							word_sensor.NoCheckWordCount = false
						case updateActonsName: // список действий
							updateName = updateActonsName
							// делаем проверки на совместимость БП и БК, закачиваем только новые действия
							// уникальный идентификатор действия: имя + комбинация воздействий в любом порядке
							// при этом списка воздействия может не быть
							akt := strings.Split(tArr[j], "|")
							aktId, _ := strconv.Atoi(akt[0])
							if akt[1] == "" {
								continue // действие без имени игнорируется
							} else {
								if akt[2] == "" {
									ActonsSincID[aktId] = aktId
									nameAkt = akt[1]
								} else {
									aktIdCur, nameAkt = IsNameActionsToBot(akt[1], akt[2])
									// фиксируем связь ID текущего и внешнего ботов
									ActonsSincID[aktId] = aktIdCur
									if aktIdCur == -1 { // не совместимое действие по БП или БК
										continue
									}
								}
							}
							// есть такое действие по ключу: совпадают имя и действие (которое может быть пустым)
							if aktIdCur > 0 {
								continue
							}
							if akt[3] != "" {
								if IsListActionToBot(akt[3], 8) == false {
									continue
								}
							}
							// если дошло до этого места, значит все проверки прошли и такого действия нет в базе текущего бота
							// и его можно залить, зафиксировав при этом связь между ID действий ботов
							termineteAction.LastTerminalActons++
							aktIdNew := termineteAction.LastTerminalActons
							ActonsSincID[aktId] = aktIdNew // фиксируем связь ID текущего и внешнего ботов
							termineteAction.TerminalActonsNameFromID[aktIdNew] = nameAkt
							termineteAction.UpdateTerminalActionsExpenses(aktIdNew, akt[2])
							termineteAction.UpdateActionsTargetsFromID(aktIdNew, akt[3])
						case updateDnkReflexes: // список рефлексов
							updateName = updateDnkReflexes
							if len(ActonsSincID) == 0 { // если массив соответствий пустой, нет смысла проверять
								lib.WritePultConsol("Массив соотвествий ID действий пустой!")
								FlgBreak = true
								break
							}
							// делаем проверки на совместимость
							rf := strings.Split(tArr[j], "|")
							// совместимость базовых эмоций
							lev1, _ := strconv.Atoi(rf[1])
							if lib.ExistsValInArr([]int{1, 2, 3}, lev1) == false {
								continue
							}
							// совместимость базовых контекстов
							if IsCompareArrValue(rf[2], 12) == false {
								continue
							}
							// совместимость пусковых стимулов
							if IsCompareArrValue(rf[3], 8) == false {
								continue
							}
							// совместимость рефлекторных действий
							rf[4] = IsActionIdToBot(rf[4])
							if rf[4] == "-1" {
								continue
							} // не нашлось соответствий в массиве
							lev2 := lib.IntArrToStrArr(strings.Split(rf[2], ","))
							lev3 := lib.IntArrToStrArr(strings.Split(rf[3], ","))
							aktArr := strings.Split(rf[4], ",")
							reflexes.CreateNewGeneticReflex(0, lev1, lev2, lev3, lib.IntArrToStrArr(aktArr))
						case updateTriggerStimulsImages: // список пусковых стимулов У-рефлексов
							updateName = updateTriggerStimulsImages
							if len(ActonsSincID) == 0 { // если массив соответствий пустой, нет смысла проверять
								lib.WritePultConsol("Массив соотвествий ID действий пустой!")
								FlgBreak = true
								break
							}
							rt := strings.Split(tArr[j], "|")
							trId, _ := strconv.Atoi(rt[0])
							rt[1] = IsActionIdToBot(rt[1])
							if rt[1] == "-1" {
								continue
							} // не нашлось соответствий в массиве
							ton := 0
							if rt[3] != "" {
								ton, _ = strconv.Atoi(rt[3])
								if lib.ExistsValInArr([]int{0, 3, 4}, ton) == false {
									continue
								} // тон не совпадает
							}
							mod := 0
							if rt[4] != "" {
								mod, _ = strconv.Atoi(rt[4])
								if lib.ExistsValInArr([]int{0, 20, 21, 22, 23, 24, 25, 26}, mod) == false {
									continue
								} // настроение не совпадает
							}
							if rt[2] != "" {
								if NoCheckWordCount == true {
									word_sensor.NoCheckWordCount = true
								}
								phr := strings.Split(rt[2], "#")
								for p := 0; p < len(phr); p++ {
									word_sensor.SetNewPhraseTreeNode(phr[p])
								}
								word_sensor.NoCheckWordCount = false
								rsr := strings.Split(rt[1], ",")
								rsar := lib.IntArrToStrArr(rsr)
								TriggerSincID[trId], _ = reflexes.CreateNewlastTriggerStimulsID(0, rsar, word_sensor.CurrentPhrasesIDarr, ton, mod)
							}
						case updateConditionReflexes: // список У-рефлексов
							updateName = updateConditionReflexes
							if len(ActonsSincID) == 0 { // если массив соответствий пустой, нет смысла проверять
								lib.WritePultConsol("Массив соотвествий ID действий пустой!")
								FlgBreak = true
								break
							}
							if len(TriggerSincID) == 0 { // если массив соответствий пустой, нет смысла проверять
								lib.WritePultConsol("Массив соотвествий ID пусковых стимулов пустой!")
								FlgBreak = true
								break
							}
							// делаем проверки на совместимость
							rcf := strings.Split(tArr[j], "|")
							rfuId, _ := strconv.Atoi(rcf[0])
							// совместимость базовых эмоций
							lev1, _ := strconv.Atoi(rcf[1])
							if lib.ExistsValInArr([]int{1, 2, 3}, lev1) == false {
								continue
							}
							// совместимость базовых контекстов
							if IsCompareArrValue(rcf[2], 12) == false {
								continue
							}
							lev2 := lib.IntArrToStrArr(strings.Split(rcf[2], ","))
							// совместимость пусковых триггеров
							lev3 := 0
							if _, ok := TriggerSincID[rfuId]; ok {
								lev3 = TriggerSincID[rfuId]
							} else {
								continue
							}
							// совместимость рефлекторных действий
							rcf[4] = IsActionIdToBot(rcf[4])
							if rcf[4] == "-1" {
								continue
							} // не нашлось соответствий в массиве
							aktArr := strings.Split(rcf[4], ",")
							reflexes.CreateNewConditionReflex(0, lev1, lev2, lev3, lib.IntArrToStrArr(aktArr), 0)
						}
						if FlgBreak == true {
							break
						} else {
							IsMod = true      // прочитана хотя бы одна запись с файла
							IsNoCompat = true // все проверки на совместимость успешны хотя бы для одной записи
						}
					}
				}
				LastModeUpdate(BotName, updateName, ModTime)
				SaveFileUpdateDir()
			}
			// обновляем данные в каталоге обмена
			if IsMod == true {
				IsMod = false
				switch updateName {
				case updateActonsName:
					termineteAction.SaveTerminalActons()
				case updateDnkReflexes:
					reflexes.SaveGeneticReflexes()
				}
				msgTxt = "Успешный импорт: "
			} else {
				// не совместимый файл обмена
				if IsNoCompat == true {
					FileUpdateDir[id].Compat = 0
					IsNoCompat = false
					msgTxt = "Не совместимый для импорта файл обмена: "
				} else {
					msgTxt = "Нет новых данных в файле: "
				}
				IsNoAllImport = true
			}
			// очищаем массивы совметимостей
			if updateName == updateConditionReflexes {
				ActonsSincID = make(map[int]int)
				TriggerSincID = make(map[int]int)
			}
			lib.WritePultConsol(msgTxt + BotName + "_" + FileName)
		}
	}
	// все файлы заблокированы, нечего обновлять
	if IsCompat == false {
		lib.WritePultConsol("Все файлы в списке обмена не совместимы для импорта!")
		IsNoAllImport = true
	}

	return !IsNoAllImport
}

/* Экспорт в файл обмена
Выгружаем все типы файлов, которые встречаются в каталоге по одному разу */
func ExportFileUpdate() bool {
	var sArr []string
	var out, outBuf, FileName, FileNameList, msgTxt string
	var i, id int
	var flgExp bool

	_, err := os.Stat(pathUpdate)
	if os.IsNotExist(err) {
		lib.WritePultConsol("Каталог обмена не найден: " + pathUpdate)
		return false
	}

	// смотрим каталог обмена и выгружаем согласно списку ботов файлы для них
	for _, ListId := range FileUpdate {
		sArr = strings.Split(ListId, "|")
		for i = 0; i < len(sArr); i++ {
			id, _ = strconv.Atoi(sArr[i])
			FileName = FileUpdateDir[id].FileName
			if lib.ExistsValStrInList(FileNameList, FileName, "|") == true {
				continue
			}
			FileNameList += FileName + "|"
			switch FileName {
			case updatePhraseName: // дерево фраз
				cnt := len(word_sensor.PhraseTreeFromID)
				if cnt == 0 {
					continue
				}
				// добавляем самую длинную фразу ветки, конечный узел
				for n := 0; n < cnt; n++ {
					outBuf = word_sensor.GetPhraseStringsFromPhraseID(word_sensor.PhraseTreeFromID[n].ID)
					if word_sensor.PhraseTreeFromID[n].Children == nil {
						if outBuf != "" {
							out += outBuf + "\r\n"
							flgExp = true
						}
					}
					outBuf = ""
				}
			case updateActonsName: // список действий
				PathFileExport := lib.MainPathExeFile + "/memory_reflex/terminal_actons.txt"
				flgExp = CopyFileToExport(PathFileExport, FileName)
			case updateDnkReflexes: // список рефлексов
				PathFileExport := lib.MainPathExeFile + "/memory_reflex/dnk_reflexes.txt"
				flgExp = CopyFileToExport(PathFileExport, FileName)
			case updateTriggerStimulsImages: // список пусковых стимулов У-рефлексов
				cnt := len(word_sensor.PhraseTreeFromID)
				if cnt == 0 {
					continue
				}
				out = ""
				for k, v := range reflexes.TriggerStimulsArr {
					out += strconv.Itoa(k) + "|"
					for i := 0; i < len(v.RSarr); i++ {
						out += strconv.Itoa(v.RSarr[i]) + ","
					}
					out += "|"
					for i := 0; i < len(v.PhraseID); i++ {
						tmpArr := word_sensor.WordsArrFromPhraseID[v.PhraseID[i]]
						if tmpArr != nil {
							out += strings.TrimSpace(word_sensor.GetStrFromArrID(tmpArr)) + "#"
						}
					}
					out += "|"
					out += strconv.Itoa(v.ToneID) + "|"
					out += strconv.Itoa(v.MoodID)
					out += "\r\n"
					flgExp = true
				}
			case updateConditionReflexes: // список у-рефлексов
				PathFileExport := lib.MainPathExeFile + "/memory_reflex/condition_reflexes.txt"
				flgExp = CopyFileToExport(PathFileExport, FileName)
			}
			if out != "" {
				out = strings.TrimSuffix(out, "\r\n")
				lib.WriteFileContent(pathUpdate+"/"+botName+"_"+FileName+".txt", out)
				out = ""
			}
			if flgExp == true {
				msgTxt = "Успешный экспорт: "
				flgExp = false
			} else {
				msgTxt = "Не удачный экспорт: "
			}
			lib.WritePultConsol(msgTxt + botName + "_" + FileName)
		}
	}
	return true
}

/* Копировать файл данных в общий каталог */
func CopyFileToExport(PathFileExport string, updateName string) bool {
	file, err := os.Stat(PathFileExport)
	if err == nil {
		if file.Size() > 0 {
			lib.CopyFile(PathFileExport, pathUpdate+"/"+botName+"_"+updateName+".txt")
			return true
		} else {
			lib.WritePultConsol("Файл пустой: " + botName + "_" + updateName)
		}
	}
	return false
}

/* Обновить дату/время в каталоге обмена */
func LastModeUpdate(BotName string, FileName string, TimeUpdate time.Time) {
	var sArr []string
	var ListID string
	var i, id int

	ListID = FileUpdate[BotName]
	sArr = strings.Split(ListID, "|")
	for i = 0; i < len(sArr); i++ {
		id, _ = strconv.Atoi(sArr[i])
		if FileUpdateDir[id].FileName == FileName {
			// убираем часовой пояс из времени
			ModtimeTxt := TimeUpdate.Format(Layout)
			ModtimeLayout, _ := time.Parse(Layout, ModtimeTxt)
			FileUpdateDir[id].LastMod = ModtimeLayout
			break
		}
	}
}

// проверка совместимости списка значений с заданным рядом значений в массиве
func IsCompareArrValue(txt string, typeArr int) bool {
	var CntArr []int
	var i int

	if txt == "" {
		return true
	} // если список пустой, то нечего проверять

	for i = 1; i < typeArr+1; i++ {
		CntArr = append(CntArr, i)
	}

	sArr := strings.Split(txt, ",")
	for i = 0; i < len(sArr); i++ {
		id, _ := strconv.Atoi(sArr[i])
		if lib.ExistsValInArr(CntArr, id) == true {
			return true
		}
	}

	return false
}

/* Массив соответствия ID рефлекторных действий текущего и внешнего ботов
Заполняется при загрузке файла действий и очищается при завершении загрузки файла У-рефлексов */
var ActonsSincID = make(map[int]int)

/* Массив соответствия ID пускровых стимулов У-рефлекса текущего и внешнего ботов
Заполняется при загрузке файла триггеров и очищается при завершении загрузки файла У-рефлексов */
var TriggerSincID = make(map[int]int)

/* Поиск номера рефлекторного действия в базе текущего бота
соответствующий номеру внешнего бота. Для этого был заполнен массив соответствий
при загрузке файла действий. */
func IsActionIdToBot(txt string) string {
	var out string

	// если база пустая, то прекращаем проверять
	if len(termineteAction.TerminalActonsNameFromID) == 0 {
		return txt
	}
	if txt == "" {
		return ""
	}

	sArr := strings.Split(txt, ",")
	for i := 0; i < len(sArr); i++ {
		id, _ := strconv.Atoi(sArr[i])
		if as, ok := ActonsSincID[id]; ok {
			out += strconv.Itoa(as) + ","
		}
	}
	if out == "" {
		out = "-1"
	} else {
		out = strings.TrimSuffix(out, ",")
	}
	return out
}

/* Проверка наличия рефлекторного действия в базе текущего бота по ключу: имя + действие
Нужно проверять пары БП>Val не как целую строку, а в любом порядке вхождения в строке
*/
func IsNameActionsToBot(ActName string, ActList string) (int, string) {
	var tArr []string
	var id, i, j, n, num int
	var aktLst, NewName string
	var isFind bool
	var err error

	// если база пустая или действия не совместимы то прекращаем проверять
	if len(termineteAction.TerminalActonsNameFromID) == 0 {
		return 0, ActName
	}
	if len(termineteAction.TerminalActionsExpensesFromID) == 0 {
		return 0, ActName
	}

	// проверяем совместимость базовых параметров в строке затрат действия
	sArr := strings.Split(ActList, ";")
	for i = 0; i < len(sArr); i++ {
		aktLst += strings.Split(sArr[i], ">")[0] + ","
	}
	aktLst = strings.TrimSuffix(aktLst, ",")
	if IsListActionToBot(aktLst, 12) == false {
		return -1, ""
	}
	aktLst = ""

	for _, val := range termineteAction.TerminalActonsNameFromID {
		if val == ActName {
			isFind = true
			break
		}
	}
	if isFind == false {
		return 0, ActName
	} else {
		// создаем новое имя
		poz := strings.LastIndex(ActName, "_")
		if poz > 0 {
			if num, err = strconv.Atoi(ActName[poz+1:]); err != nil {
				num = 0
			}
			ActName = ActName[:poz]
		}
		NewName = ActName + "_" + strconv.Itoa(num+1)
	}

	// формируем сводку действий по имени ActName
	for AktId := range termineteAction.TerminalActionsExpensesFromID {
		se := termineteAction.TerminalActionsExpensesFromID[AktId]
		if termineteAction.TerminalActonsNameFromID[AktId] == ActName {
			for i = 0; i < len(se); i++ {
				aktLst += strconv.Itoa(se[i].GomeoID) + ">" + fmt.Sprintf("%.1f", se[i].Diff) + ";"
			}
			aktLst = strings.TrimSuffix(aktLst, ";")
			tArr = append(tArr, strconv.Itoa(AktId)+"|"+aktLst)
			aktLst = ""
		}
	}

	// ищем совпадения в любой последовательности с ActList
	sArr = strings.Split(ActList, ";")
	for i = 0; i < len(tArr); i++ {
		aktIdArr := strings.Split(tArr[i], "|")
		aktArr := strings.Split(aktIdArr[1], ";")
		for j = 0; j < len(sArr); j++ {
			for n = 0; n < len(aktArr); n++ {
				if aktArr[i] == sArr[j] {
					num++
				}
			}
		}
		// совпало имя и список действий по кол-во и вхождению
		if num == len(sArr) && num == len(aktArr) {
			id, _ = strconv.Atoi(aktIdArr[0])
			break
		}
	}

	return id, NewName
}

/* Проверка совместимости списка базовых параметров внешнего бота с базой текущего */
func IsListActionToBot(txt string, typeArr int) bool {
	sArr := strings.Split(txt, ",")
	for _, val := range sArr {
		if IsCompareArrValue(val, typeArr) == false {
			return false
		}
	}

	return true
}
