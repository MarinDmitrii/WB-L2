package pattern

import "fmt"

/*
Реализовать паттерн «фабричный метод».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
https://en.wikipedia.org/wiki/Factory_method_pattern
*/

/*
Применимость паттерна "фабричный метод":
1. Заранее не известны точные типы и зависимости объектов
2. Возможность расширения внутренних компонентов
3. Экономия системных ресурсов, повторно используя существующие объекты

Плюсы паттерна "фабричный метод":
1. позволяет сделать код создания объектов более универсальным, не привязываясь к конкретным классам
2. Принцип единой ответственности: Можно переместить код создания продукта в одно место программы, что упростит поддержку кода.
3. Принцип открытости/закрытости: Можно вводить в программу новые виды продуктов, не нарушая существующий код.

Минусы паттерна "фабричный метод":
1. Необходимость создавать наследника Creator для каждого нового типа продукта
2. Код может стать более сложным: В небольших программах добавление фабричных методов может привести к избыточности кода и усложнению его структуры
*/

// Невозможно реализовать классический шаблон фабричного метода в Go из-за отсутствия функций ООП,
// таких как классы и наследование.

// Интерфейс для оружия
type IGun interface {
	setName(name string)
	setPower(power int)
	getName() string
	getPower() int
}

// Реализация базового класса оружия
type Gun struct {
	name  string
	power int
}

// Методы для установки и получения имени и мощности оружия
func (g *Gun) setName(name string) {
	g.name = name
}

func (g *Gun) setPower(power int) {
	g.power = power
}

func (g *Gun) getName() string {
	return g.name
}

func (g *Gun) getPower() int {
	return g.power
}

// Структура для конкретного вида оружия - BFG
type BFG struct {
	Gun
}

// Фабричный метод для создания экземпляра конкретного вида оружия - BFG
func newBFG() IGun {
	return &BFG{
		Gun: Gun{
			name:  "Big Fucking Gun",
			power: 9,
		},
	}
}

// Структура для конкретного вида оружия - Buriza
type Buriza struct {
	Gun
}

// Фабричный метод для создания экземпляра конкретного вида оружия - Buriza
func newBuriza() IGun {
	return &Buriza{
		Gun: Gun{
			name:  "Buriza-Do Kyanon",
			power: 5,
		},
	}
}

// Функция для получения экземпляра оружия по типу
func getGun(gunType string) (IGun, error) {
	if gunType == "BFG" {
		return newBFG(), nil
	} else if gunType == "Buriza" {
		return newBuriza(), nil
	} else {
		return nil, fmt.Errorf("%v - is wrong gun type", gunType)
	}
}

func printDetails(g IGun) {
	fmt.Printf("Gun's name: %v\n", g.getName())
	fmt.Printf("Gun's power: %v\n", g.getPower())
}

func factoryMethodPattern() {
	// Создание экземпляра BFG
	bfg, err := getGun("BFG")
	if err != nil {
		fmt.Println(err)
	} else {
		printDetails(bfg)
	}

	// Создание экземпляра Buriza
	buriza, err := getGun("Buriza")
	if err != nil {
		fmt.Println(err)
	} else {
		printDetails(buriza)
	}

	// Попытка создания экземпляра оружия с неверным типом
	rocketLauncher, err := getGun("Rocket")
	if err != nil {
		fmt.Println(err)
	} else {
		printDetails(rocketLauncher)
	}
}
