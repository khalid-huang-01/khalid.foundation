package gameworld

import (
	"fmt"
	"time"
)

type Hero struct {
	HealthPoint int
	Location int
}

func (hero *Hero) isDead() bool {
	return hero.HealthPoint <= 0
}


func (hero *Hero) attack(monster *Monster) {
	monster.HealthPoint -= 10
}

func (hero *Hero) isReachable(monster *Monster) bool {
	return monster.Location <= hero.Location + 3 && monster.Location >= hero.Location - 3
}

func (hero *Hero) String() string {
	return fmt.Sprintf("hero healthPoint is: %v", hero.HealthPoint)
}

type monsterStrategyInterface interface {
	AttackStrategy(hero *Hero)
	IsReachableStrategy(hero *Hero, monster *Monster) bool
	Name() string
}

// 稳定不应该依赖于不稳定的具体实现，而要去依赖稳定的抽象接口
type Monster struct {
	HealthPoint int
	Location int
	strategy monsterStrategyInterface
}

func (m *Monster) isDead() bool {
	return m.HealthPoint <= 0
}

func (m *Monster) attack(hero *Hero) {
	m.strategy.AttackStrategy(hero)
}

func (m *Monster) isReachable(hero *Hero) bool {
	return m.strategy.IsReachableStrategy(hero, m)
}

func (m *Monster) String() string {
	return fmt.Sprintf("%v health point is %v", m.strategy.Name(), m.HealthPoint)
}

// 不同策略模式的实现
type LevelOneMonsterStrategy struct {
}

func (l LevelOneMonsterStrategy) AttackStrategy(hero *Hero) {
	hero.HealthPoint -= 10
}

func (l LevelOneMonsterStrategy) IsReachableStrategy(hero *Hero, monster *Monster) bool {
	return hero.Location <= monster.Location + 1 && hero.Location >= monster.Location - 1
}

func (l LevelOneMonsterStrategy) Name() string {
	return "Level One Monster"
}

type LevelTwoMonsterStrategy struct {
}

func (l LevelTwoMonsterStrategy) AttackStrategy(hero *Hero) {
	hero.HealthPoint -= 20
}

func (l LevelTwoMonsterStrategy) IsReachableStrategy(hero *Hero, monster *Monster) bool {
	return hero.Location <= monster.Location + 2 && hero.Location >= monster.Location - 2
}

func (l LevelTwoMonsterStrategy) Name() string {
	return "Level Two Monster"
}

// 可以优化为工厂方法
func newLevelOneMonster() *Monster {
	// 策略可以换成单例
	return &Monster{
		HealthPoint: 30,
		Location:    2,
		strategy:    LevelOneMonsterStrategy{},
	}
}

func newLevelTwoMonster() *Monster {
	return &Monster{
		HealthPoint: 20,
		Location:    4,
		strategy:    LevelTwoMonsterStrategy{},
	}
}

// 战场实例化
type BattleField struct {
	Monsters            []*Monster
	NumOfRemainMonsters int
	stopCh              chan struct{}
}


// 英雄进场，开始
func (b *BattleField) Start(hero *Hero) {
	t := time.NewTicker(2*time.Second)
	for {
		select {
		case <-t.C:
			//轮番进行攻击
			 go func() {
				 fmt.Println("hero attack monster")
				 b.heroAttack(hero)
				 printSituation(hero, b.Monsters)
				 time.Sleep(1*time.Second) // 通过sleep的方式让出执行权限，让主线程有时间check b.stopCh是否结束
				 fmt.Println("monsters attack hero")
				 b.monsterAttack(hero)
				 printSituation(hero, b.Monsters)
			 }()
		case <-b.stopCh:
			if hero.isDead() {
				fmt.Println("hero dead, fail the battle")
			} else {
				fmt.Println("hero win, pass the battle")
			}
			return
		}
	}

}

func (b *BattleField) monsterAttack(hero *Hero) {
	for _, m := range b.Monsters {
		if m.isDead() {
			continue
		}
		if m.isReachable(hero) {
			m.attack(hero)
		}
		if hero.isDead() {
			close(b.stopCh)
		}
	}
}

func (b *BattleField) heroAttack(hero *Hero) {
	for _, m := range b.Monsters {
		if m.isDead() {
			continue
		}
		if hero.isReachable(m) {
			hero.attack(m)
		}
		if m.isDead() {
			b.NumOfRemainMonsters++
			if b.NumOfRemainMonsters == len(b.Monsters) {
				close(b.stopCh)
			}
		}
	}
}

func newBattleField() *BattleField {
	b := &BattleField{}
	b.stopCh = make(chan struct{})
	b.Monsters = make([]*Monster, 0)
	b.Monsters = append(b.Monsters, newLevelOneMonster())
	b.Monsters = append(b.Monsters, newLevelOneMonster())
	b.Monsters = append(b.Monsters, newLevelTwoMonster())
	b.NumOfRemainMonsters = 0
	return b
}


func printSituation(hero *Hero, monsters []*Monster) {
	fmt.Println(hero)
	for _, m := range monsters {
		fmt.Println(m)
	}
	fmt.Println("")
}
