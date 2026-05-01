package main

import (
	"fmt"
	"math/rand"
)

type Skill interface {
	Name() string
	Description() string
	ReqClass() string   // "Qualquer", "Guerreiro", "Mago Elementalista", "Ladino"
	ReqElement() string // "Qualquer", ElemFire, ElemWater, etc.
	Cast(caster *Character, target *Character) string
	Upgrade()
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

// --- HELPER DE DANO (Centraliza os cálculos) ---
func calculateDamage(caster, target *Character, baseDmg float64) string {
	elemMult := getElementalMultiplier(caster.Element, target.Element)

	isCrit := rand.Float64() < caster.Stats.CritChance
	critMult := 1.0
	critText := ""
	if isCrit {
		critMult = caster.Stats.CritMultiplier
		critText = " ⚡CRÍTICO!"
	}

	finalDmg := baseDmg * elemMult * critMult
	target.Stats.HP -= finalDmg

	elemText := ""
	if elemMult > 1.0 {
		elemText = " (Super Efetivo!)"
	}
	if elemMult < 1.0 {
		elemText = " (Resistido)"
	}

	return fmt.Sprintf("causou %.0f dano!%s%s", finalDmg, critText, elemText)
}

// --- EXEMPLOS DE HABILIDADES COM LORE ---

// 1. NEUTRA (Qualquer um pega)
type GolpeRapido struct{ BaseSkill }

func (s *GolpeRapido) ReqClass() string    { return "Qualquer" }
func (s *GolpeRapido) ReqElement() string  { return "Qualquer" }
func (s *GolpeRapido) Description() string { return "Ataque veloz neutro." }
func (s *GolpeRapido) Cast(c, t *Character) string {
	dmg := c.Stats.Str * (0.8 + (0.2 * float64(s.level)))
	return "💨 Golpe Rápido " + calculateDamage(c, t, dmg)
}

// 2. GUERREIRO DE FOGO
type LaminaVulcanica struct{ BaseSkill }

func (s *LaminaVulcanica) ReqClass() string    { return "Guerreiro" }
func (s *LaminaVulcanica) ReqElement() string  { return ElemFire }
func (s *LaminaVulcanica) Description() string { return "Corta com a fúria do vulcão." }
func (s *LaminaVulcanica) Cast(c, t *Character) string {
	dmg := c.Stats.Str * (1.5 + (0.3 * float64(s.level)))
	return "🌋 Lâmina Vulcânica " + calculateDamage(c, t, dmg)
}

// 3. MAGO DE ÁGUA
type TsunamiArcano struct{ BaseSkill }

func (s *TsunamiArcano) ReqClass() string    { return "Mago Elementalista" }
func (s *TsunamiArcano) ReqElement() string  { return ElemWater }
func (s *TsunamiArcano) Description() string { return "Invoca uma onda esmagadora." }
func (s *TsunamiArcano) Cast(c, t *Character) string {
	dmg := c.Stats.Int * (2.0 + (0.5 * float64(s.level)))
	return "🌊 Tsunami " + calculateDamage(c, t, dmg)
}

// 4. LADINO DE VENTO
type GolpeAssasino struct{ BaseSkill }

func (s *GolpeAssasino) ReqClass() string    { return "Ladino" }
func (s *GolpeAssasino) ReqElement() string  { return ElemWind }
func (s *GolpeAssasino) Description() string { return "Corte silencioso e mortal." }
func (s *GolpeAssasino) Cast(c, t *Character) string {
	dmg := c.Stats.Str * (1.8 + (0.4 * float64(s.level)))
	return "🗡️ Golpe Assassino " + calculateDamage(c, t, dmg)
}

// O GRANDE FILTRO DE LEVEL UP
var GlobalSkillPool = []Skill{
	&GolpeRapido{BaseSkill{name: "Golpe Rápido", level: 1, element: "Qualquer"}},
	&LaminaVulcanica{BaseSkill{name: "Lâmina Vulcânica", level: 1, element: ElemFire}},
	&TsunamiArcano{BaseSkill{name: "Tsunami Arcano", level: 1, element: ElemWater}},
	&GolpeAssasino{BaseSkill{name: "Golpe Assassino", level: 1, element: ElemWind}},
}

func getSkillOptions(player *Character) []Skill {
	var validSkills []Skill
	for _, skill := range GlobalSkillPool {
		classMatch := skill.ReqClass() == "Qualquer" || skill.ReqClass() == player.Name
		elemMatch := skill.ReqElement() == "Qualquer" || skill.ReqElement() == player.Element
		if classMatch && elemMatch {
			// Cria uma cópia limpa da skill
			newSkill := skill
			validSkills = append(validSkills, newSkill)
		}
	}

	// Retorna 3 aleatórias (embaralhar e pegar o slice)
	if len(validSkills) == 0 {
		// Se não houver skills específicas, oferece sempre o Golpe Rápido
		validSkills = append(validSkills, &GolpeRapido{BaseSkill{name: "Golpe Rápido", level: 1, element: "Qualquer"}})
	}

	rand.Shuffle(len(validSkills), func(i, j int) { validSkills[i], validSkills[j] = validSkills[j], validSkills[i] })
	if len(validSkills) > 3 {
		return validSkills[:3]
	}
	return validSkills
}

// Mantém compatibilidade com o getRandomNewSkill antigo
func getRandomNewSkill() Skill {
	skills := getSkillOptions(&Character{Name: "Guerreiro"}) // Default
	if len(skills) > 0 {
		return skills[0]
	}
	return &GolpeRapido{BaseSkill{name: "Golpe Rápido", level: 1, element: "Qualquer"}}
}
