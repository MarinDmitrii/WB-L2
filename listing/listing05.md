Что выведет программа? Объяснить вывод программы.

```go
package main

type customError struct {
	msg string
}

func (e *customError) Error() string {
	return e.msg
}

func test() *customError {
	{
		// do something
	}
	return nil
}

func main() {
	var err error
	err = test()
	if err != nil {
		println("error")
		return
	}
	println("ok")
}
```

Ответ:
```
error

объявили переменную err типа error
error - встроенный интерфейс, который определяет один метод Error(), возвращающий строку.
Переменной err присваем значение функции test(), которая возвращает указатель на стуктуру customError, равный nil.
У структуры customError есть метод Error(), который возвращает string. Значит структуру реализует интерфейс error.
Так как у переменной err определён тип, то она не является пустым интерфейсом. Поэтому проверка 'err != nil будет пройдена и на печать будет выведено "error"

```
