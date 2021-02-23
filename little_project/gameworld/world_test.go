package gameworld

import "testing"

func TestWorld(t *testing.T) {
	//实例化英雄
	hero := &Hero{
		HealthPoint: 100,
		Location:    3,
	}
	battle := newBattleField()
	battle.Start(hero)
}

