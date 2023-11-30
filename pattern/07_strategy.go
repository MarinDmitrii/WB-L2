package pattern

import "fmt"

/*
Реализовать паттерн «стратегия».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
https://en.wikipedia.org/wiki/Strategy_pattern
*/

/*
Применимость паттерна "стратегия":
1. Когда в одном и том же месте в зависимости от текущего состояния системы (или её окружения) должны использоваться различные алгоритмы
2. Когда в наличии много похожих классов, которые отличаются только способом выполнения некоторого поведения
3. Когда необходимо изолировать бизнес-логику класса от деталей реализации алгоритмов
4. Когда необходимо заменить массивный условный оператор, который переключается между различными вариантами одного и того же алгоритма

Плюсы паттерна "стратегия":
1. Отказ от использования переключателей и/или условных операторов
2. Принцип открытости/закрытости: Добавление новых стратегий не требует изменения кода контекста
3. Улучшение тестируемости: Каждая стратегия может быть протестирована независимо.

Минусы паттерна "стратегия":
1. Увеличение числа классов: Каждая стратегия требует свой класс, что может увеличить количество классов в системе.
2. Код может стать более сложным: В программе с небольшим количеством редко меняющихся алгоритмов, применение паттерна может привести к избыточности кода и усложнению его структуры
*/

// PaymentStrategy определяет интерфейс стратегии оплаты.
type PaymentStrategy interface {
	Pay(amount float64) string
}

// CreditCardPayment представляет стратегию оплаты кредитной картой.
type CreditCardPayment struct {
}

func (p *CreditCardPayment) Pay(amount float64) string {
	return fmt.Sprintf("Paid %.2f $ using credit card", amount)
}

// PayPalPayment представляет стратегию оплаты через PayPal.
type PayPalPayment struct {
}

func (p *PayPalPayment) Pay(amount float64) string {
	return fmt.Sprintf("Paid %.2f $ using PayPal", amount)
}

// PaymentContext представляет контекст, который использует выбранную стратегию для оплаты.
type PaymentContext struct {
	strategy PaymentStrategy
}

func (c *PaymentContext) SetStrategy(strategy PaymentStrategy) {
	c.strategy = strategy
}

func (c *PaymentContext) MakePayment(amount float64) string {
	return c.strategy.Pay(amount)
}

func strategyPattern() {
	creditCardStrategy := CreditCardPayment{}
	payPalStrategy := PayPalPayment{}

	context := PaymentContext{}
	context.SetStrategy(&creditCardStrategy)
	fmt.Println(context.MakePayment(100.48))

	context.SetStrategy(&payPalStrategy)
	fmt.Println(context.MakePayment(44.88))
}
