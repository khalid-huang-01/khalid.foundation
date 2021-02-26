package gameworld

import (
	"fmt"
	"testing"
	"time"
)

func TestWorld(t *testing.T) {
	//实例化英雄
	hero := &Hero{
		HealthPoint: 100,
		Location:    3,
	}
	battle := newBattleField()
	battle.Start(hero)
}


func TestTime(t *testing.T) {
	createTimeStamp := time.Now().Unix()
	createTime := time.Unix(createTimeStamp, 0)
	fmt.Println(createTime)
	fmt.Println(createTime.Add(time.Duration(10) * time.Second))
}