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

func (s *HeavyStrike) Cast(caster, target *Character) string {
	// 1. Calcula o multiplicador elemental
	elemMult := getElementalMultiplier(caster.Element, target.Element)

	// 2. Aplica o multiplicador no dano final
	dmgBase := caster.Stats.Str * (1.0 + (0.2 * float64(s.level-1)))
	dmgFinal := dmgBase * elemMult

	target.Stats.HP -= dmgFinal

	// 3. Adiciona um "Flavor Text" para o jogador saber que deu certo!
	feedback := ""
	if elemMult > 1.3 {
		feedback = " (SUPER EFETIVO! +30%)"
	} else if elemMult < 1.0 {
		feedback = " (Resistido... -30%)"
	}

	return fmt.Sprintf("💥 Golpe causou %.0f dano!%s", dmgFinal, feedback)
}
func (s *HeavyStrike) Description() string { return "Dano físico massivo." }

type Fireball struct{ BaseSkill }

func (s *Fireball) Cast(caster, target *Character) string {
	// 1. Calcula o multiplicador elemental
	elemMult := getElementalMultiplier(caster.Element, target.Element)

	// 2. Aplica o multiplicador no dano final
	dmgBase := caster.Stats.Int * (2.5 + (0.5 * float64(s.level-1)))
	dmgFinal := dmgBase * elemMult

	target.Stats.HP -= dmgFinal

	// 3. Adiciona um "Flavor Text" para o jogador saber que deu certo!
	feedback := ""
	if elemMult > 1.3 {
		feedback = " (SUPER EFETIVO! +30%)"
	} else if elemMult < 1.0 {
		feedback = " (Resistido... -30%)"
	}

	return fmt.Sprintf("🔥 Magia causou %.0f dano!%s", dmgFinal, feedback)
}
func (s *Fireball) Description() string { return "Alto dano elemental." }

func getRandomNewSkill() Skill {
	if rand.Intn(2) == 0 {
		return &Fireball{BaseSkill{name: "Chama Elemental", level: 1, element: "Fire"}}
	}
	return &HeavyStrike{BaseSkill{name: "Ataque Base", level: 1, element: "Neutral"}}
}
