/*
индексная страница http://localhost:8181/index

В проекте много глобальных переменных, что привычно раздражает Свидетелей инкапсуляции и непорочного пространства имен,
но ТАК НУЖНО (спорно) для одганизации среды,
схожей с организацией линкующих указателей в мозге (т.е. связей с одного распознавателя к целому ансаблю - объекту).
Ну и есть немало других вещей, нарушающих Порядок и Традиции Golang.
Попытки использовать горутины оказались просто неуместными (спорно) и просто ненужными, учитывая вряд ли в чем-то могущий быть выигрыш.
Короче, код предоставляется на вольное растерзание и свободное экспериментирование, без претензий, сорри за возможный негатив.
Везде много пространных комментариев, которые запутывают даже меня, но они НУЖНЫ.
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

var xxxxxxx=0 // в дебаге иногда начинается циклический вызов, это - защелка

// true - остановка всей активности для совершения критических глобальных операций
var isGlobalStopAllActivnost = false

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
			// текстовый блок для набивки дерева слов-фраз из http://go/pages/words.php
			text_block := r.FormValue("text_block")
			if len(text_block) > 0 {
				brain.IsPultActivnost = true
				res := word_sensor.SetNewTextBlock(text_block)
				brain.IsPultActivnost = false
				_, _ = fmt.Fprint(resp, res)
				return
			}
			// текст из окна ввода с пульта
			text_dlg := r.FormValue("text_dlg")
			if len(text_dlg) > 0 {
				brain.IsPultActivnost = true
				is_input_rejim, _ := strconv.Atoi(r.FormValue("is_input_rejim")) // режим быстрого формирования у-рефлексов
				if is_input_rejim == 0 { // наоборот
					reflexes.IsUnlimitedMode = 1
				} else {
					reflexes.IsUnlimitedMode = 0
				}
				toneID, _ := strconv.Atoi(r.FormValue("pult_tone"))
				pultMood := r.FormValue("pult_mood") // тон сообщения
				moodID, _ := strconv.Atoi(pultMood) // настроение сообщения
				res := word_sensor.VerbalDetection(text_dlg, is_input_rejim, toneID, moodID)
				// если добавлены пусковые стимулы (нажаты кнопки на пульте)
				set_img_action := r.FormValue("set_img_action")
				if len(set_img_action) > 0 {
					// brain.IsPultActivnost = true
					enegry, _ := strconv.Atoi(r.FormValue("food_portion"))
					action_sensor.SetActionFromPult(set_img_action, enegry)
					/*
						// активировать дерево действием
						reflexes.ActiveFromAction()
						brain.IsPultActivnost = false
					*/
				}

				reflexes.ActiveFromPhrase() // активировать дерево рефлексов фразой - только для условных рефлексов
				brain.IsPultActivnost = false

				var answerStr=""
				if len(lib.ActionsForPultStr)>5{
					answerStr=lib.ActionsForPultStr
					lib.ActionsForPultStr = "" // очистка для новой порции
				}
				_, _ = fmt.Fprint(resp, res+"|&|"+answerStr)
				return
			}

			// отправить на пульт состояние гомеостаза Beast и его базовые контексты
			sincronism := r.FormValue("sincronism")
			if len(sincronism) > 0 {
				// выполнить цикл действий по пульсу перед отправкой результата на Пульт
				brain.SincroTic()
				_, _ = fmt.Fprint(resp, "sincronism")
				return
			}

			// отправить на пульт состояние гомеостаза Beast и его базовые контексты
			getParams := r.FormValue("get_params")
			if len(getParams) > 0 {
				brain.IsPultActivnost = true
				outStr := gomeostas.GetCurGomeoParams()

				var waitingPeriod = ""
				res,time := psychic.WaitingPeriodForActions()
				if res {
					waitingPeriod = strconv.Itoa(time)
				}

				var psichicReady = "0"
				if psychic.StartPsichicNow{
					psichicReady = "1"
				}
				if psychic.AllowConsciousnessProcess{
					psichicReady = "2"
				}

				outStr += "#|#" + gomeostas.GetCurGomeoStatus() + "#|#" + gomeostas.GetCurContextActive() +
					"#|# " + reflexes.GetCurrentConditionsStr() + // чтобы постоянно была инфа о сочетаниях контекстов
					"#|#" + strconv.Itoa(brain.LifeTime) +
					"#|#" + reflexes.NoUnconditionRefles +
					"#|#" + waitingPeriod +
					"#|#" + psichicReady
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
		// проверка активности Beast, аозвращает текущий brain.PulsCount
		brain.IsPultActivnost = true
		check_Beast_activnost := r.FormValue("check_Beast_activnost")
		if check_Beast_activnost == "1" {
			_, _ = fmt.Fprint(resp, brain.PulsCount)
			return
		}
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
		// Формирование зеркальных автоматизмов на основе списка ответов
		mirror_making_fool := r.FormValue("mirror_making_fool")
		if len(mirror_making_fool)>0 {
			res:=psychic.FormingMirrorAutomatizmFromList(mirror_making_fool)
			_, _ = fmt.Fprint(resp, res)
			return
		}
		// Формирование зеркальных автоматизмов на основе общего шаблона
		mirror_making_temp := r.FormValue("mirror_making_temp")
		if len(mirror_making_temp)>0 {
			res:=psychic.FormingMirrorAutomatizmFromTempList(mirror_making_temp)
			_, _ = fmt.Fprint(resp, res)
			return
		}

		if !isGlobalStopAllActivnost {
			setExpParam := r.FormValue("set_exp_param") // экспорт данных
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

			setImpParam := r.FormValue("set_imp_param") // импорт данных
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

			get_words_list := r.FormValue("get_words_list")
			if get_words_list == "1" {
				brain.IsPultActivnost = true
				phTree := word_sensor.GetWordsListForPult()
				brain.IsPultActivnost = false
				if phTree == "!!!" {
					return // запрет показа карты во время распознавания и записи
				}
				_, _ = fmt.Fprint(resp, phTree)
				return
			}

			deleting_word := r.FormValue("deleting_word")
			if deleting_word =="1"{
				delete_word := r.FormValue("delete_word")
				deleteWord,_:=strconv.Atoi(delete_word)
				brain.IsPultActivnost = true
				word_sensor.DeleteWord(deleteWord)
				brain.IsPultActivnost = false
				_, _ = fmt.Fprint(resp, "OK")
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

//

			set_action := r.FormValue("set_action")
			if len(set_action) > 0 && xxxxxxx==0 {
				xxxxxxx=1
				brain.IsPultActivnost = true
				//ActionsForPultOldStr = lib.ActionsForPultOldStr
				enegry, _ := strconv.Atoi(r.FormValue("food_portion"))
				action_sensor.SetActionFromPult(set_action, enegry)

				// активировать дерево действием
				reflexes.ActiveFromAction()
				brain.IsPultActivnost = false
				var answerStr=""
				if len(lib.ActionsForPultStr)>5{
					answerStr=lib.ActionsForPultStr
					lib.ActionsForPultStr = "" // очистка для новой порции
				}
				_, _ = fmt.Fprint(resp, answerStr)
				xxxxxxx=0
				return
			}

			cliner_time_condition_reflex := r.FormValue("cliner_time_condition_reflex")
			if len(cliner_time_condition_reflex) > 0 {
				isGlobalStopAllActivnost = true
				ret := reflexes.ClinerTimeConditionReflex()
				isGlobalStopAllActivnost = false
				_, _ = fmt.Fprint(resp, ret)
				return
			}

			get_condition_reflex_info := r.FormValue("get_condition_reflex_info")
			if len(get_condition_reflex_info) > 0 {
				base := r.FormValue("limitBasicID")
				limitBasicID,_:=strconv.Atoi(base)
				ref := reflexes.GetConditionReflexInfo(limitBasicID)
				_, _ = fmt.Fprint(resp, ref)
				return
			}

			get_automatizm_list_info := r.FormValue("get_automatizm_list_info")
			if len(get_automatizm_list_info) > 0 {
				lpage := r.FormValue("limitBasicID")
				limitPage,_:=strconv.Atoi(lpage)
				ref := psychic.GetAutomatizmInfo(limitPage)
				_, _ = fmt.Fprint(resp, ref)
				return
			}

			get_mental_automatizm_list_info := r.FormValue("get_mental_automatizm_list_info")
			if len(get_mental_automatizm_list_info) > 0 {
				base := r.FormValue("limitPage")
				limitPage,_:=strconv.Atoi(base)
				ref := psychic.GetMentalAutomatizmInfo(limitPage)
				_, _ = fmt.Fprint(resp, ref)
				return
			}

			get_trigger_info := r.FormValue("get_trigger_info")
			if len(get_trigger_info) > 0 {
				triggerID,_:=strconv.Atoi(r.FormValue("triggerID"))
				ref := reflexes.GetTreeAutomatizmTriggersInfo(triggerID)
				_, _ = fmt.Fprint(resp, ref)
				return
			}

			get_sequence_info := r.FormValue("get_sequence_info")
			if len(get_sequence_info) > 0 {
				autmzmID,_:=strconv.Atoi(r.FormValue("autmzmID"))
				ref := psychic.GetAutomatizmSequenceInfo(autmzmID, false)
				_, _ = fmt.Fprint(resp, ref)
				return
			}

			get_mental_actiom_img_info := r.FormValue("get_mental_actiom_img_info")
			if len(get_mental_actiom_img_info) > 0 {
				mautmzmID,_:=strconv.Atoi(r.FormValue("mautmzmID"))
				ref := psychic.GetMentalAutomotizmActionsString(mautmzmID, false)
				_, _ = fmt.Fprint(resp, ref)
				return
			}

			get_emotion_info := r.FormValue("get_emotion_info")
			if len(get_emotion_info) > 0 {
				emotionID,_:=strconv.Atoi(r.FormValue("emotionID"))
				ref := psychic.GetStrnameFromBaseImageID(emotionID)
				_, _ = fmt.Fprint(resp, ref)
				return
			}

			get_object_info := r.FormValue("get_object_info")
			if len(get_object_info) > 0 {
				objectID,_:=strconv.Atoi(r.FormValue("objectID"))
				ref := psychic.GetStrnameFromobjectID(objectID)
				_, _ = fmt.Fprint(resp, ref)
				return
			}

			cliner_atmtzm_block := r.FormValue("cliner_atmtzm_block")
			if len(cliner_atmtzm_block) > 0 {
				atmtzmID,_:=strconv.Atoi(r.FormValue("atmtzmID"))
				ref := psychic.UnblockAutomatizmID(atmtzmID)
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
				base := r.FormValue("limitBasicID")
				limitBasicID,_:=strconv.Atoi(base)
				ref := psychic.GetAutomatizmTreeForPult(limitBasicID)
				_, _ = fmt.Fprint(resp, ref)
				return
			}
			get_node_automatizms := r.FormValue("get_node_automatizms")
			if len(get_node_automatizms) > 0 {
				nodeID,_:=strconv.Atoi(r.FormValue("autNodeID"))
				ref := psychic.GetAutomatizmForNodeInfo(nodeID)
				_, _ = fmt.Fprint(resp, ref)
				return
			}


			get_mental_automatizm_tree := r.FormValue("get_mental_automatizm_tree")
			if len(get_mental_automatizm_tree) > 0 {
				base := r.FormValue("limit")
				limit,_:=strconv.Atoi(base)
				ref := psychic.GetMentalAutomatizmTreeForPult(limit)
				_, _ = fmt.Fprint(resp, ref)
				return
			}
			get_node_mental_situations := r.FormValue("get_node_mental_situations")
			if len(get_node_mental_situations) > 0 {
				nodeID,_:=strconv.Atoi(r.FormValue("autNodeID"))
				ref := psychic.GetMentalSituationsForNodeInfo(nodeID)
				_, _ = fmt.Fprint(resp, ref)
				return
			}
			get_node_mental_purpose := r.FormValue("get_node_mental_purpose")
			if len(get_node_mental_purpose) > 0 {
				nodeID,_:=strconv.Atoi(r.FormValue("autNodeID"))
				ref := psychic.GetMentalPurposeForNodeInfo(nodeID)
				_, _ = fmt.Fprint(resp, ref)
				return
			}


			get_rulles_list_info := r.FormValue("get_rulles_list_info")
			if get_rulles_list_info =="1" {
				rulles:=psychic.GetCur10lastRules()
				_, _ = fmt.Fprint(resp, rulles)
				return
			}

			get_mental_rulles_list_info := r.FormValue("get_mental_rulles_list_info")
			if get_mental_rulles_list_info =="1" {
				rulles:=psychic.GetCur10lastMentalRules()
				_, _ = fmt.Fprint(resp, rulles)
				return
			}

			get_mental_action_info := r.FormValue("get_mental_action_info")
			if len(get_mental_action_info)>0 {
				actID,_:=strconv.Atoi(get_mental_action_info)
				rulles:=psychic.GetMentalActionInfo(actID)
				_, _ = fmt.Fprint(resp, rulles)
				return
			}

			get_mental_rules_cickle_info := r.FormValue("get_mental_rules_cickle_info")
			if len(get_mental_rules_cickle_info)>0 {
				rulles:=psychic.GetMentalRulesCickleInfo(get_mental_rules_cickle_info)
				_, _ = fmt.Fprint(resp, rulles)
				return
			}

			get_mental_cickles_list_info := r.FormValue("get_mental_cickles_list_info")
			if get_mental_cickles_list_info =="1" {
				ref := psychic.GetCicklesToPult()
				_, _ = fmt.Fprint(resp, ref)
				return
			}
			get_atmzm_tree_info := r.FormValue("get_atmzm_tree_info")
			if get_atmzm_tree_info =="1" {
				objID,_:=strconv.Atoi(r.FormValue("objID"))
				ref := psychic.GetAtmzmTreeInfo(objID)
				_, _ = fmt.Fprint(resp, ref)
				return
			}
			get_undstg_tree_info := r.FormValue("get_undstg_tree_info")
			if get_undstg_tree_info =="1" {
				objID,_:=strconv.Atoi(r.FormValue("objID"))
				ref := psychic.GetUndstgTreeInfo(objID)
				_, _ = fmt.Fprint(resp, ref)
				return
			}
			get_ment_atmzm_info := r.FormValue("get_ment_atmzm_info")
			if get_ment_atmzm_info =="1" {
				objID,_:=strconv.Atoi(r.FormValue("objID"))
				ref := psychic.GetMentAtmzmInfo(objID)
				_, _ = fmt.Fprint(resp, ref)
				return
			}

			make_automatizms_from_reflexes := r.FormValue("make_automatizms_from_reflexes")
			if len(make_automatizms_from_reflexes) == 1 {
				ref := reflexes.RunMakeAutomatizmsFromReflexes()
				_, _ = fmt.Fprint(resp, ref)
				return
			}
			make_automatizms_from_genetic_reflexes := r.FormValue("make_automatizms_from_genetic_reflexes")
			if len(make_automatizms_from_genetic_reflexes) == 1 {
				ref := reflexes.RunMakeAutomatizmsFromGeneticReflexes()
				_, _ = fmt.Fprint(resp, ref)
				return
			}
// объект значимости
			get_mental_importance_list_info := r.FormValue("get_mental_importance_list_info")
			if get_mental_importance_list_info =="1" {
				ref := psychic.GetImportanceToPult()
				_, _ = fmt.Fprint(resp, ref)
				return
			}
			get_importance_object_info := r.FormValue("get_importance_object_info")
			if get_importance_object_info =="1" {
				objID,_:=strconv.Atoi(r.FormValue("objID"))
				objType,_:=strconv.Atoi(r.FormValue("objType"))
				ref := psychic.GetImportanceObjectInfo(objID,objType)
				_, _ = fmt.Fprint(resp, ref)
				return
			}

			//получить фразы, используемые в автоматизмах для иконки выбора Пульта
			conditions_words_basic := r.FormValue("conditions_words_basic")
			if conditions_words_basic =="1" {
				// против concurrent map iteration and map write
				brain.IsPultActivnost = true
				bID:=r.FormValue("basicID")
				bID=strings.Trim(bID," ")
				basicID,_:=strconv.Atoi(bID)
				contexts:=r.FormValue("contexts")
				ref := psychic.GetAutomatizmPraseList(basicID,contexts)
				brain.IsPultActivnost = false
				_, _ = fmt.Fprint(resp, ref)
				return
			}

			// обнулить параметры гомеостаза Beast
			cliner_gomeo_pars := r.FormValue("cliner_gomeo_pars")
			if len(cliner_gomeo_pars) > 0 {
				brain.IsPultActivnost = true
				value,_:=strconv.ParseFloat(cliner_gomeo_pars, 64)
				gomeostas.ClinerAllGomeoParams(value)
				brain.IsPultActivnost = false
				_, _ = fmt.Fprint(resp, "1")
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