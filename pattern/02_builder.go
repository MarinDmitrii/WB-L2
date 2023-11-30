package pattern

import "fmt"

/*
	Реализовать паттерн «строитель».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Builder_pattern
*/

/*
Применимость паттерна "строитель":
1. Позволяет избежать "телескопического конструктора", когда есть сложные объекты, состоящие из множества частей
2. Отделяет конструирование сложного объекта от его представления так, что в результате одного и того же процесса конструирования могут получаться разные представления

Плюсы паттерна "строитель":
1. Изолирует код для конструирования объекта от его представления
2. Гарантирует пошаговое построение сложных объектов
3. Один и тот же код для конструирования можно повторно использовать при создании различных представлений объекта

Минусы паттерна "строитель":
1. Увеличивает сложность кода из-за введения дополнительных классов
*/

// House представляет структуру для описания дома.
type House struct {
	Walls   string
	Windows string
	Doors   string
	Roof    string
}

// Builder определяет интерфейс для создания домов.
type Builder interface {
	SetWallsType()
	SetWindowsType()
	SetDoorsType()
	SetRoofType()
	GetHouse() *House
}

// getBuilder возвращает конкретного строителя в зависимости от типа.
func getBuilder(builderType string) Builder {
	if builderType == "wooden" {
		return NewWoodenBuilder()
	} else if builderType == "stone" {
		return NewStoneBuilder()
	} else {
		return nil
	}
}

// WoodenBuilder представляет строителя для деревянного дома.
type WoodenBuilder struct {
	WallsType   string
	WindowsType string
	DoorsType   string
	RoofType    string
}

// NewWoodenBuilder создает новый экземпляр WoodenBuilder.
func NewWoodenBuilder() *WoodenBuilder {
	return &WoodenBuilder{}
}

// SetWallsType устанавливает тип стен для деревянного дома.
func (wb *WoodenBuilder) SetWallsType() {
	wb.WallsType = "Wooden Walls"
}

// SetWindowsType устанавливает тип окон для деревянного дома.
func (wb *WoodenBuilder) SetWindowsType() {
	wb.WindowsType = "Wooden Windows"
}

// SetDoorsType устанавливает тип дверей для деревянного дома.
func (wb *WoodenBuilder) SetDoorsType() {
	wb.DoorsType = "Wooden Doors"
}

// SetRoofType устанавливает тип крыши для деревянного дома.
func (wb *WoodenBuilder) SetRoofType() {
	wb.RoofType = "Wooden Roof"
}

// GetHouse возвращает экземпляр дома, построенного деревянным строителем.
func (wb *WoodenBuilder) GetHouse() *House {
	return &House{
		Walls:   wb.WallsType,
		Windows: wb.WindowsType,
		Doors:   wb.DoorsType,
		Roof:    wb.RoofType,
	}
}

// StoneBuilder представляет строителя для каменного дома.
type StoneBuilder struct {
	WallsType   string
	WindowsType string
	DoorsType   string
	RoofType    string
}

// NewStoneBuilder создает новый экземпляр StoneBuilder.
func NewStoneBuilder() *StoneBuilder {
	return &StoneBuilder{}
}

// SetWallsType устанавливает тип стен для каменного дома.
func (sb *StoneBuilder) SetWallsType() {
	sb.WallsType = "Stone Walls"
}

// SetWindowsType устанавливает тип окон для каменного дома.
func (sb *StoneBuilder) SetWindowsType() {
	sb.WindowsType = "Stone Windows"
}

// SetDoorsType устанавливает тип дверей для каменного дома.
func (sb *StoneBuilder) SetDoorsType() {
	sb.DoorsType = "Stone Doors"
}

// SetRoofType устанавливает тип крыши для каменного дома.
func (sb *StoneBuilder) SetRoofType() {
	sb.RoofType = "Stone Roof"
}

// GetHouse возвращает экземпляр дома, построенного каменным строителем.
func (sb *StoneBuilder) GetHouse() *House {
	return &House{
		Walls:   sb.WallsType,
		Windows: sb.WindowsType,
		Doors:   sb.DoorsType,
		Roof:    sb.RoofType,
	}
}

// Director управляет процессом построения дома с использованием строителя.
type Director struct {
	builder Builder
}

// NewDirector создает новый экземпляр Director с указанным строителем.
func NewDirector(b Builder) *Director {
	return &Director{builder: b}
}

// SetBuilder устанавливает строителя для Director.
func (d *Director) SetBuilder(b Builder) {
	d.builder = b
}

// buildHouse создает дом, используя строителя, установленного в Director.
func (d *Director) buildHouse() *House {
	d.builder.SetWallsType()
	d.builder.SetWindowsType()
	d.builder.SetDoorsType()
	d.builder.SetRoofType()
	return d.builder.GetHouse()
}

// builderPattern демонстрирует использование паттерна Builder для построения различных домов.
func builderPattern() {
	stoneBuilder := getBuilder("stone")
	woodenBuilder := getBuilder("wooden")

	director := NewDirector(stoneBuilder)
	stoneHouse := director.buildHouse()
	fmt.Printf("Stone house consist of: %v, %v, %v, %v\n", stoneHouse.Walls, stoneHouse.Windows, stoneHouse.Doors, stoneHouse.Roof)

	director.SetBuilder(woodenBuilder)
	woodenHouse := director.buildHouse()
	fmt.Printf("Wooden house consist of: %v, %v, %v, %v\n", woodenHouse.Walls, woodenHouse.Windows, woodenHouse.Doors, woodenHouse.Roof)
}
