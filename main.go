/*
индексная страница http://localhost:8181/index
*/

package main

import (
	"BOT/brain"
	"BOT/brain/action_sensor"
	"BOT/brain/gomeostas"
	"BOT/brain/psychic"
	"BOT/brain/reflexes"
	"BOT/brain/update"
	"BOT/brain/words_sensor"
	"BOT/closer"
	"BOT/lib"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

/* В проекте много глобальных переменных, что привычно раздражает Свидетелей инкапсуляции и непорочного пространства имен,
но ТАК НУЖНО (спорно) для одганизации среды,
схожей с организацией линкующих указателей в мозге (т.е. связей с одного распознавателя к целому ансаблю - объекту).
Ну и есть немало других вещей, нарушающих Порядок и Традиции Golang.
Попытки использовать горутины оказались просто неуместными (спорно) и просто ненужными, учитывая вряд ли в чем-то могущий быть выигрыш.
Короче, код предоставляется на вольное растерзание и свободное экспериментирование, без претензий, сорри за возможный негатив.
Везде много пространных комментариев, которые запутывают даже меня, но они НУЖНЫ.
*/
var isGlobalStopAllActivnost = false // true - остановка всей активности для совершения критических глобальных операций

func receiveSend(resp http.ResponseWriter, r *http.Request) {
	resp.Header().Set("Access-Control-Allow-Origin", "*")
	resp.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// посылается постоянно раз в 1сек (из /common/linking.php) или с запросом или для подтверждения связи,
	// а так же для передачи по текущему пульсу информации от Beast, например WritePultConsol()
	if r.Method == "POST" {
		if gomeostas.CheckBeastDeath() {
			isGlobalStopAllActivnost = true
			_, _ = fmt.Fprint(resp, "!!!")
		}

		if !isGlobalStopAllActivnost {
			// 	текстовый блок для набивки дерева слов-фраз из http://go/pages/words.php
			text_block := r.FormValue("text_block")
			if len(text_block) > 0 {
				brain.IsPultActivnost = true
				res := word_sensor.SetNewTextBlock(text_block)
				brain.IsPultActivnost = false
				_, _ = fmt.Fprint(resp, res)
				return
			}

			text_dlg := r.FormValue("text_dlg")
			if len(text_dlg) > 0 {
				brain.IsPultActivnost = true
				is_input_rejim, _ := strconv.Atoi(r.FormValue("is_input_rejim"))
				if is_input_rejim == 0 { // наоборот
					reflexes.IsUnlimitedMode = 1
				} else {
					reflexes.IsUnlimitedMode = 0
				}
				toneID, _ := strconv.Atoi(r.FormValue("pult_tone"))
				pultMood := r.FormValue("pult_mood")
				moodID, _ := strconv.Atoi(pultMood)
				res := word_sensor.VerbalDetection(text_dlg, is_input_rejim, toneID, moodID)
				// если добавлены пусковые стимулы
				set_img_action := r.FormValue("set_img_action")
				if len(set_img_action) > 0 {
					//brain.IsPultActivnost = true
					enegry, _ := strconv.Atoi(r.FormValue("food_portion"))
					action_sensor.SetActionFromPult(set_img_action, enegry)
					/*
						//// активировать дерево действием
						reflexes.ActiveFromAction()
						brain.IsPultActivnost = false
					*/
				}

				reflexes.ActiveFromPhrase() // активировать дерево рефлексов фразой - только для условных рефлексов
				brain.IsPultActivnost = false
				_, _ = fmt.Fprint(resp, res)
				return
			}

			// отправить на пульт состояние гомеостаза Beast и его базовые контексты
			getParams := r.FormValue("get_params")
			if len(getParams) > 0 {
				brain.IsPultActivnost = true
				outStr := gomeostas.GetCurGomeoParams()

				outStr += "#|#" + gomeostas.GetCurGomeoStatus() + "#|#" + gomeostas.GetCurContextActive() +
					"#|# " + reflexes.GetCurrentConditionsStr() + //чтобы постоянно была инфа о сочетаниях контекстов
					"#|#" + strconv.Itoa(brain.LifeTime) +
					"#|#" + reflexes.NoUnconditionRefles
				brain.IsPultActivnost = false
				_, _ = fmt.Fprint(resp, outStr)
				return
			}

			// установка параметров гомеостаза с Пульта:
			// задать параметры гомеостаза Beast
			setParamsId := r.FormValue("set_params")
			if len(setParamsId) > 0 {
				brain.IsPultActivnost = true
				id, _ := strconv.Atoi(setParamsId)
				gomeostas.SetCurGomeoParams(id, r.FormValue("params_val"))
				brain.IsPultActivnost = false
				_, _ = fmt.Fprint(resp, "1")
				return
			}

			//  передача информации от Beast в Пульт различается идентификатором назначения информации перед самой информацией
			// информация для консоли ничинается с идентификатора назначения: "CONSOL:"
			if len(lib.WritePultConsolStr) > 0 {
				_, _ = fmt.Fprint(resp, "CONSOL:"+lib.WritePultConsolStr)
				lib.WritePultConsolStr = "" // очистка для новой порции
				return
			}

			// информация о действиях Beast ничинается с идентификатора назначения: "ACTION:"
			if len(lib.ActionsForPultStr) > 0 {
				_, _ = fmt.Fprint(resp, "ACTION:"+lib.ActionsForPultStr)
				lib.ActionsForPultStr = "" // очистка для новой порции
				return
			}

			// если ничего выше не было, то:
			// передача на Пульт сигнала готовности - когда нет других запросов - посылаетс сигнал на Пульт в function bot_answer(res)
			if word_sensor.IsReadyWordSensorLevel() {
				//идентификатор назначения информации: "READY"
				_, _ = fmt.Fprint(resp, "READY")
				return
			}
			_, _ = fmt.Fprint(resp, "POST")
		} else {
			brain.NotAllowAnyActions = true

			// Сформировать условные рефлексы на основе списка фраз-синонимов
			file_for_condition_reflexes := r.FormValue("file_for_condition_reflexes")
			if len(file_for_condition_reflexes) > 0 {
				reflexes.FormingConditionsRefleaxFromList(file_for_condition_reflexes)
				_, _ = fmt.Fprint(resp, "OK")
			}

		}
		//fmt.Println("EMPTY")
	}

	if r.Method == "GET" {
		// остановка любой активности Beast
		brain.IsPultActivnost = true
		stop_activnost := r.FormValue("stop_activnost")
		if stop_activnost == "1" {
			isGlobalStopAllActivnost = true
			_, _ = fmt.Fprint(resp, "stop")
			return
		}
		// восстановление активности Beast
		start_activnost := r.FormValue("start_activnost")
		if start_activnost == "1" {
			isGlobalStopAllActivnost = false
			brain.IsPultActivnost = false
			_, _ = fmt.Fprint(resp, "active")
			return
		}
		// ЗОНА ОСОБЫХ ДЕЙСТВИЙ в период остановленной активности Beast:

		// Сохранить текущее состояние Beast
		save_all_memory := r.FormValue("save_all_memory")
		if save_all_memory == "1" {
			brain.IsPultActivnost = true
			if brain.SaveAll() {
				_, _ = fmt.Fprint(resp, "yes")
				brain.IsPultActivnost = false
				return
			}
			_, _ = fmt.Fprint(resp, "no")
			brain.IsPultActivnost = false
			return
		}

		// корректное выключение Beast
		bot_closing := r.FormValue("bot_closing")
		if bot_closing == "1" {
			brain.IsPultActivnost = true
			cleanupFunc()
			closer.Close()
			return
		}

		if !isGlobalStopAllActivnost {
			setExpParam := r.FormValue("set_exp_param")
			if len(setExpParam) > 0 {
				brain.IsPultActivnost = true
				if setExpParam == "1" {
					IsExpTrue, expTxt := update.ExportFileUpdate([]int{1,2,3,4,5})
					if IsExpTrue == true {
						setExpParam = "yes|" + expTxt
					} else {
						setExpParam = "no|" + expTxt
					}
				}
				brain.IsPultActivnost = false
				_, _ = fmt.Fprint(resp, setExpParam)
				return
			}
			setImpParam := r.FormValue("set_imp_param")
			if len(setImpParam) > 0 {
				brain.IsPultActivnost = true
				if setImpParam == "1" {
					IsImpParam, impTxt := update.ImportFileUpdate([]int{1,2,3,4,5})
					if IsImpParam == true {
						setImpParam = "yes|" + impTxt
					} else {
						setImpParam = "no|" + impTxt
					}
				}
				brain.IsPultActivnost = false
				_, _ = fmt.Fprint(resp, setImpParam)
				return
			}

			get_reflex_tree := r.FormValue("get_reflex_tree")
			if get_reflex_tree == "1" {
				brain.IsPultActivnost = true
				tree := reflexes.GetReflexesTreeForPult()
				brain.IsPultActivnost = false
				if tree == "!!!" {
					return // запрет показа карты во время распознавания и записи
				}
				_, _ = fmt.Fprint(resp, tree)
				return
			}

			get_phrase_tree := r.FormValue("get_phrase_tree")
			if get_phrase_tree == "1" {
				brain.IsPultActivnost = true
				phTree := word_sensor.GetPhraseTreeForPult()
				brain.IsPultActivnost = false
				if phTree == "!!!" {
					return // запрет показа карты во время распознавания и записи
				}
				_, _ = fmt.Fprint(resp, phTree)
				return
			}

			get_word_tree := r.FormValue("get_word_tree")
			if get_word_tree == "1" {
				brain.IsPultActivnost = true
				phTree := word_sensor.GetWordTreeForPult()
				brain.IsPultActivnost = false
				if phTree == "!!!" {
					return // запрет показа карты во время распознавания и записи
				}
				_, _ = fmt.Fprint(resp, phTree)
				return
			}

			set_action := r.FormValue("set_action")
			if len(set_action) > 0 {
				brain.IsPultActivnost = true
				enegry, _ := strconv.Atoi(r.FormValue("food_portion"))
				action_sensor.SetActionFromPult(set_action, enegry)

				//// активировать дерево действием
				reflexes.ActiveFromAction()
				brain.IsPultActivnost = false
				_, _ = fmt.Fprint(resp, "!")
				return
			}

			get_condition_reflex_info := r.FormValue("get_condition_reflex_info")
			if len(get_condition_reflex_info) > 0 {
				ref := reflexes.GetConditionReflexInfo()
				_, _ = fmt.Fprint(resp, ref)
				return
			}

			get_self_perception_info := r.FormValue("get_self_perception_info")
			if len(get_self_perception_info) > 0 {
				ref := psychic.GetSelfPerceptionInfo()
				_, _ = fmt.Fprint(resp, ref)
				return
			}

			get_automatizm_tree := r.FormValue("get_automatizm_tree")
			if len(get_automatizm_tree) > 0 {
				ref := psychic.GetAutomatizmTreeForPult()
				_, _ = fmt.Fprint(resp, ref)
				return
			}
			_, _ = fmt.Fprint(resp, "GET")
		}
	}

}

// инициализация
func init() {
	lib.GetMainPathExeFile()
	log.SetFlags(log.Lshortfile | log.LstdFlags)


	// тестирование комбинаций. Если бы время работы было приемлемо,
	//то можно было бы запускать процесс из Пульта в меню Инструменты (шестеренка)
	// 	tools.MakeContextCombinations()
}

// старт
func main() {
	defer closer.Close()
	// для перехвата при завершении программы (использует пакет "BOT/closer" https://github.com/xlab/closer):
	closer.Bind(cleanupFunc)

	address := lib.ReadFileContent(lib.GetMainPathExeFile() + "/common/linking_address.txt")
	address = strings.TrimSpace(address)[7:]

	brain.RunInitialisation() // init.go
	brain.Puls()

	http.HandleFunc("/", receiveSend)
	_ = http.ListenAndServe(address, nil)
	fmt.Println("Сервер запущен...")

	// в самом конце
	closer.Hold()
}

// отключение Beast
func cleanupFunc() {
	lib.WritePultConsol("Beast вырубается.")
	fmt.Print("ПОСЛЕДНИЕ ДЕЙСТВИЯ ПЕРЕД ЗАКРЫВАНИЕМ ПРОГРАММЫ")
	// записать текущее состояние Дерева Моделей и Эпизодическую память
	//	brain.PrepareBeforCloseTreeModel()
	/* если внезапно откобчить мозг человека, то из памяти пропадет все, что было в последние полчаса
	так что просто записывать PrepareBeforCloseTreeModel() раз в 10 минут и при создании нового узла дерева
	*/
	brain.SaveAll()
}

// здесь могут быть функции для обеспечения связи между пакетами чтобы избегать цикличного импорта
