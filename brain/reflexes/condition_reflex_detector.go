/* распознаватель условного рефлекса

*/

package reflexes



//////////////////////////////////////////

func conditionRexlexFound(cond []int)(bool){
	if cond==nil || len(cond)==0{
		return false
	}
// var ConditionReflexesFrom3=make(map[int]*ConditionReflex)
	reflex:=ConditionReflexesFrom3[cond[0]]
	if reflex==nil {
		return false
	}
		// время жизни рефлекса минус 10 дней
		//life:=LifeTime - reflex.activationTime -10
		life:= reflex.lastActivation - reflex.activationTime -3600*24*10
		if life<0{life=0}
		// коэффициент влияния времени жизни: каждые 10 дней укрепляют рефлекс в 2 раза
		k:=1+ (2*life)/(3600*24*10)
// определить время просрочки рефлекса при неиспользовании > 10 дней с учетом времни жизни
		delay:=(3600*24*10)*k
//Рефлкс НЕ активен, если его lastActivation меньше, чем activationTime-delay
		if reflex.lastActivation > reflex.activationTime-delay {
			// обновить актуальность рефлекса
			reflex.lastActivation=reflex.activationTime
			conditionReflexesIdArr = append(conditionReflexesIdArr, reflex.ID)
			return true
		}else{// рефлекс просрочен и должен быть дезактивирован
			reflex.lastActivation=0
		}

	return false
}
//////////////////////////////////////////////