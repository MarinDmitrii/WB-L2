package pattern

import "fmt"

/*
	Реализовать паттерн «комманда».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Command_pattern
*/

/*
Применимость паттерна "комманда":
1. Организация очереди: Паттерн поддерживает создание и управление очередью запросов.
2. Отмена операций: Позволяет реализовать отмену и повторение операций.
3. Параметризация объектов: Команды могут быть параметризованы и передаваться между объектами, что делает их многоразовыми и гибкими.

Плюсы паттерна "комманда":
1. Принцип единой ответственности: Классы, вызывающие операции, можно отделить от классов, выполняющих эти операции.
2. Принцип открытости/закрытости: Можно вводить в приложение новые команды, не ломая существующий код.
3. Позволяет реализовать отмену и повторение операций.

Минусы паттерна "комманда":
1. Код может стать более сложным, поскольку вводится новый слой между отправителями и получателями.
*/

// Command - интерфейс команды с методом Execute, который будет вызывать выполнение команды.
type Command interface {
	Execute()
}

// Button - структура кнопки, которая содержит команду для выполнения.
type Button struct {
	command Command
}

// press - метод для нажатия кнопки, который вызывает выполнение связанной с кнопкой команды.
func (b *Button) press() {
	b.command.Execute()
}

// Device - интерфейс устройства с методами On, Off и Next.
type Device interface {
	On()
	Off()
	Next()
}

// OnCommand - команда включения устройства.
type OnCommand struct {
	device Device
}

// Execute - метод выполнения команды включения.
func (c *OnCommand) Execute() {
	c.device.On()
}

// OffCommand - команда выключения устройства.
type OffCommand struct {
	device Device
}

// Execute - метод выполнения команды выключения.
func (c *OffCommand) Execute() {
	c.device.Off()
}

// NextCommand - команда переключения на следующий канал устройства.
type NextCommand struct {
	device Device
}

// Execute - метод выполнения команды переключения на следующий канал.
func (c *NextCommand) Execute() {
	c.device.Next()
}

// TV - структура телевизора с полем isRunning, обозначающим состояние работы телевизора.
type TV struct {
	isRunning bool
}

// On - метод включения телевизора.
func (t *TV) On() {
	t.isRunning = true
	fmt.Println("Turning tv on")
}

// Off - метод выключения телевизора.
func (t *TV) Off() {
	t.isRunning = false
	fmt.Println("Turning tv off")
}

// Next - метод переключения на следующий канал телевизора, если он включен.
// В противном случае выводится сообщение о невозможности переключения, так как телевизор выключен.
func (t *TV) Next() {
	if t.isRunning {
		fmt.Println("Switching tv channel")
	} else {
		fmt.Println("Can't switch tv channel, tv off")
	}
}

func commandPattern() {
	tv := &TV{}

	onCommand := &OnCommand{device: tv}
	offCommand := &OffCommand{device: tv}
	nextCommand := &NextCommand{device: tv}

	onButton := &Button{command: onCommand}
	offButton := &Button{command: offCommand}
	nextButton := &Button{command: nextCommand}

	onButton.press()
	nextButton.press()
	offButton.press()
	nextButton.press()
}
