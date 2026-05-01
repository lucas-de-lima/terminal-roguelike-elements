package main

import (
	"fmt"
	"math/rand"
)

type Skill interface {
	Name() string
	Description() string
	Cast(caster *Character, target *Character) string
	Upgrade()
	Element() string
}

type BaseSkill struct {
	name    string
	level   int
	element string
}

func (b *BaseSkill) Name() string {
	if b.level > 1 {
		return fmt.Sprintf("%s (Nv.%d)", b.name, b.level)
	}
	return b.name
}
func (b *BaseSkill) Upgrade()        { b.level++ }
func (b *BaseSkill) Element() string { return b.element }

type HeavyStrike struct{ BaseSkill }

func (s *HeavyStrike) Cast(c, t *Character) string {
	dmg := c.Stats.Str * (1.0 + (0.2 * float64(s.level-1)))
	t.Stats.HP -= dmg
	return fmt.Sprintf("💥 Golpe: -%.0f HP", dmg)
}
func (s *HeavyStrike) Description() string { return "Dano físico massivo." }

type Fireball struct{ BaseSkill }

func (s *Fireball) Cast(c, t *Character) string {
	dmg := c.Stats.Int * (2.5 + (0.5 * float64(s.level-1)))
	t.Stats.HP -= dmg
	return fmt.Sprintf("🔥 Magia: -%.0f HP", dmg)
}
func (s *Fireball) Description() string { return "Alto dano elemental." }

func getRandomNewSkill() Skill {
	if rand.Intn(2) == 0 {
		return &Fireball{BaseSkill{name: "Chama Elemental", level: 1, element: "Fire"}}
	}
	return &HeavyStrike{BaseSkill{name: "Ataque Base", level: 1, element: "Neutral"}}
}
