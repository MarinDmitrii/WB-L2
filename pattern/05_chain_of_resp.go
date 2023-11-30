package pattern

import "fmt"

/*
Реализовать паттерн «цепочка вызовов».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
https://en.wikipedia.org/wiki/Chain-of-responsibility_pattern
*/

/*
Применимость паттерна "цепочка вызовов":
1. Обработка различных типов запросов различными способами, когда не известны типы запросов и их последовательность
2. Выполнение нескольких обработчиков в определённом порядке

Плюсы паттерна "цепочка вызовов":
1. Принцип единой ответственности: Можете отделить классы, вызывающие операции, от классов, выполняющих операции.
2. Принцип открытости/закрытости: Можно вводить в приложение новые обработчики, не ломая существующий код.
3. Контроль порядка обработки запросов

Минусы паттерна "цепочка вызовов":
1. Некоторые запросы могут быть не обработаны
*/

// Информация о клиенте
type Patient struct {
	name             string
	registrationDone bool
	doctorChekUpDone bool
	medicineDone     bool
	paymentDone      bool
}

// Интерфейс для выполнения операций в разных отделах
type Departament interface {
	execute(*Patient)
	setNext(Departament)
}

// Регистратура
type Reception struct {
	next Departament
}

func (r *Reception) execute(p *Patient) {
	if p.registrationDone {
		fmt.Println("Patient registration already done")
		r.next.execute(p)
		return
	} else {
		fmt.Println("Reception registering patient")
		p.registrationDone = true
		r.next.execute(p)
	}
}

func (r *Reception) setNext(next Departament) {
	r.next = next
}

// Врач
type Doctor struct {
	next Departament
}

func (d *Doctor) execute(p *Patient) {
	if p.doctorChekUpDone {
		fmt.Println("Doctor checkup already done")
		d.next.execute(p)
		return
	} else {
		fmt.Println("Doctor checking patient")
		p.doctorChekUpDone = true
		d.next.execute(p)
	}
}

func (d *Doctor) setNext(next Departament) {
	d.next = next
}

// Медицинский отдел
type Medical struct {
	next Departament
}

func (m *Medical) execute(p *Patient) {
	if p.medicineDone {
		fmt.Println("Medicine already given to patient")
		m.next.execute(p)
		return
	} else {
		fmt.Println("Medical giving medicine to patient")
		p.medicineDone = true
		m.next.execute(p)
	}
}

func (m *Medical) setNext(next Departament) {
	m.next = next
}

// Касса
type Cashier struct {
	next Departament
}

func (c *Cashier) execute(p *Patient) {
	if p.paymentDone {
		fmt.Println("Payment Done")
		return
	} else {
		fmt.Println("Patient paying to the cashier")
		p.paymentDone = true
	}
}

func (c *Cashier) setNext(next Departament) {
	c.next = next
}

func chainOfReponsibilityPattern() {
	cashier := &Cashier{}

	// Указываем следующий отдел (касса)
	medical := &Medical{}
	medical.setNext(cashier)

	// Указываем следующий отдел (медицинский)
	doctor := &Doctor{}
	doctor.setNext(medical)

	// Указываем следующий отдел (доктор)
	reception := &Reception{}
	reception.setNext(doctor)

	patient := &Patient{name: "Dave"}
	// Пациент пришёл в регистратуру
	reception.execute(patient)
}
