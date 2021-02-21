package gameworld


type Hero struct {
	HealthPoint int
	Location int
}

func (hero *Hero) isDead() bool {
	return hero.HealthPoint <= 0
}


func (hero *Hero) attack(monster *Monster) {
	monster.HealthPoint -= 100
}

func (hero *Hero) isReachable(monster *Monster) bool {
	return monster.Location <= hero.Location + 3 && monster.Location >= hero.Location - 3
}

type monsterStrategyInterface interface {
	AttackStrategy(hero *Hero)
	IsReachableStrategy(hero *Hero, monster *Monster) bool
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

// 不同策略模式的实现
type LevelOneMonsterStrategy struct {
}

func (l LevelOneMonsterStrategy) AttackStrategy(hero *Hero) {
	hero.HealthPoint -= 10
}

func (l LevelOneMonsterStrategy) IsReachableStrategy(hero *Hero, monster *Monster) bool {
	return hero.Location <= monster.Location + 1 && hero.Location >= monster.Location - 1
}

type LevelTwoMonsterStrategy struct {
}

func (l LevelTwoMonsterStrategy) AttackStrategy(hero *Hero) {
	hero.HealthPoint -= 20
}

func (l LevelTwoMonsterStrategy) IsReachableStrategy(hero *Hero, monster *Monster) bool {
	return hero.Location <= monster.Location + 2 && hero.Location >= monster.Location - 2
}

// 可以优化为工厂方法
func newLevelOneMonster() *Monster {
	// 策略可以换成单例
	return &Monster{
		HealthPoint: 50,
		Location:    3,
		strategy:    LevelOneMonsterStrategy{},
	}
}

func newLevelTwoMonster() *Monster {
	return &Monster{
		HealthPoint: 10,
		Location:    5,
		strategy:    LevelTwoMonsterStrategy{},
	}
}

// 战场实例化
type BattleField struct {
	Monsters []Monster
	NumOfRemainMonsters int
}

// 英雄进场，开始
func (b *BattleField) Start(hero *Hero) {
	for {

	}
}

func (b *BattleField) monsterAttack(hero *Hero) {

}

func (b *BattleField) heroAttack(hero *Hero) {

}

func newBattleField() *BattleField {

}