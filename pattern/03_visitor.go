package pattern

import (
	"fmt"
	"math"
)

/*
	Реализовать паттерн «посетитель».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Visitor_pattern
*/

/*
Применимость паттерна "посетитель":
1. Необходимо выполнить некоторую (одну и ту же) операцию для ряда объектов, без добавления этой операции в сам класс каждого объекта.
2. Определить новую операцию над объектами без изменения их структуры.

Плюсы паттерна "посетитель":
1. Принцип открытости/закрытости. Можно ввести новое поведение, которое может работать с объектами различных классов, не изменяя эти классы.
2. Принцип единой ответственности. Можно перенести несколько версий одного и того же поведения в один класс.

Минусы паттерна "посетитель":
1. Если структура элементов часто меняется, приходится изменять все посетителей
2. Усложнение кода из-за введения дополнительных классов и интерфейсов
*/

// Shape - интерфейс для всех геометрических форм
type Shape interface {
	accept(Visitor)
}

// Rectangle - структура, представляющая прямоугольник
type Rectangle struct {
	a float64 // длина
	b float64 // ширина
}

// accept - метод, реализующий интерфейс Shape для прямоугольника
func (r *Rectangle) accept(v Visitor) {
	v.visitForRectangle(r)
}

// Circle - структура, представляющая круг
type Circle struct {
	radius float64 // радиус
}

// accept - метод, реализующий интерфейс Shape для круга
func (c *Circle) accept(v Visitor) {
	v.visitForCircle(c)
}

// Visitor - интерфейс для посетителя геометрических форм
type Visitor interface {
	visitForRectangle(*Rectangle)
	visitForCircle(*Circle)
}

// AreaCalculator - структура, реализующая интерфейс Visitor для вычисления площади
type AreaCalculator struct {
}

// visitForRectangle - метод для вычисления площади прямоугольника
func (ac *AreaCalculator) visitForRectangle(r *Rectangle) {
	area := r.a * r.b
	fmt.Printf("Area of rectangle is - %v", area)
}

// visitForCircle - метод для вычисления площади круга
func (ac *AreaCalculator) visitForCircle(c *Circle) {
	area := math.Pi * c.radius * c.radius
	fmt.Printf("Area of circle with radius %v is - %.2f", c.radius, area)
}

func visitorPattern() {
	rectangle := &Rectangle{3, 4}
	circle := &Circle{4}

	ac := &AreaCalculator{}
	rectangle.accept(ac)
	circle.accept(ac)
}
