package pattern

import "fmt"

/*
Реализовать паттерн «состояние».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
https://en.wikipedia.org/wiki/State_pattern
*/

/*
Применимость паттерна "состояние":
1. Когда есть объект, который ведёт себя по-разному в зависимости от текущего состояния
2. Когда код, управляющий состоянием, становится сложным и трудно поддерживаемым
3. Когда поведение объекта зависит от его состояния и должно изменяться во время выполнения

Плюсы паттерна "состояние":
1. Принцип единой ответственности: Код, относящийся к конкретным состояниям, можно выделить в отдельные классы.
2. Принцип открытости/закрытости: Можно вводить новые состояния, не меняя существующие классы состояний или контекст.

Минусы паттерна "состояние":
1. Код может стать более сложным: Если у объекта небольшое количество состояний или они редко меняются, применение паттерна может привести к избыточности кода и усложнению его структуры
*/

// DoorState определяет интерфейс для различных состояний двери.
type DoorState interface {
	Open()
	Close()
	Lock()
	Unlock()
}

// OpenState - состояние открытой двери
type OpenState struct{}

func (s *OpenState) Open() {
	fmt.Println("The door is already open")
}

func (s *OpenState) Close() {
	fmt.Println("Closing the door")
}

func (s *OpenState) Lock() {
	fmt.Println("Can't lock an opened door")
}

func (s *OpenState) Unlock() {
	fmt.Println("The door is already unlocked")
}

// CloseLockState - состояние закрытой и заблокированной двери
type CloseLockState struct{}

func (s *CloseLockState) Open() {
	fmt.Println("Cant' open an locked door")
}

func (s *CloseLockState) Close() {
	fmt.Println("The door is already closed")
}

func (s *CloseLockState) Lock() {
	fmt.Println("The door is already locked")
}

func (s *CloseLockState) Unlock() {
	fmt.Println("Unlocking the door")
}

// CloseUnlockState - состояние закрытой и незаблокированной двери
type CloseUnlockState struct{}

func (s *CloseUnlockState) Open() {
	fmt.Println("Opening the door")
}

func (s *CloseUnlockState) Close() {
	fmt.Println("The door is already closed")
}

func (s *CloseUnlockState) Lock() {
	fmt.Println("Locking the door")
}

func (s *CloseUnlockState) Unlock() {
	fmt.Println("The door is already unlocked")
}

// Door представляет собой контекст, который хранит текущее состояние двери.
type Door struct {
	state DoorState
}

func (d *Door) setState(state DoorState) {
	d.state = state
}

func (d *Door) Open() {
	d.state.Open()
	// Если состояние двери "закрыта, незаблокирована", то изменяем состояние на "открыта"
	if _, isCloseUnlock := d.state.(*CloseUnlockState); isCloseUnlock {
		d.setState(&OpenState{})
	}
}

func (d *Door) Close() {
	d.state.Close()
	// Если состояние двери "открыта", то изменяем состояние на "закрыта, незаблокирована"
	if _, isOpen := d.state.(*OpenState); isOpen {
		d.setState(&CloseUnlockState{})
	}
}

func (d *Door) Lock() {
	d.state.Lock()
	// Если состояние двери "закрыта,незаблокирована", то изменяем состояние двери на "закрыта, заблокирована"
	if _, isCloseUnlock := d.state.(*CloseUnlockState); isCloseUnlock {
		d.setState(&CloseLockState{})
	}
}

func (d *Door) Unlock() {
	d.state.Unlock()
	// Если состояние двери "закрыта, заблокирована", то изменяем состояние на "закрыта, незаблокирована"
	if _, isCloseLock := d.state.(*CloseLockState); isCloseLock {
		d.setState(&CloseUnlockState{})
	}
}

func statePattern() {
	door := &Door{state: &OpenState{}}

	// Открытое состояние двери
	door.Open()
	door.Lock()
	door.Unlock()
	door.Close()

	// Вызов door.Close() изменил состояние двери на "закрыта, незаблокирована"
	door.Close()
	door.Unlock()
	door.Lock()

	// Вызов door.Lock() изменил состояние двери на "закрыта, заблокирована"
	door.Open()
	door.Close()
	door.Lock()
	door.Unlock()

	// Вызов door.Unlock() изменил состояние двери на "закрыта, незаблокирована"
	door.Open()
	// Вызов door.Open() изменил состояние двери на "открытая"
}
