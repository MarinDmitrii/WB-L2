package main

import (
	"fmt"
	"sort"
	"strings"
)

/*
=== Поиск анаграмм по словарю ===

Напишите функцию поиска всех множеств анаграмм по словарю.
Например:
'пятак', 'пятка' и 'тяпка' - принадлежат одному множеству,
'листок', 'слиток' и 'столик' - другому.

Входные данные для функции: ссылка на массив - каждый элемент которого - слово на русском языке в кодировке utf8.
Выходные данные: Ссылка на мапу множеств анаграмм.
Ключ - первое встретившееся в словаре слово из множества
Значение - ссылка на массив, каждый элемент которого, слово из множества. Массив должен быть отсортирован по возрастанию.
Множества из одного элемента не должны попасть в результат.
Все слова должны быть приведены к нижнему регистру.
В результате каждое слово должно встречаться только один раз.

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

func sortString(s string) string {
	r := []rune(s)
	compare := func(i, j int) bool {
		r1 := r[i]
		r2 := r[j]

		result := strings.Compare(string(r1), string(r2))

		return result < 0
	}

	sort.SliceStable(r, compare)

	return string(r)
}

func getAnagram(words []string) map[string][]string {
	myMap := make(map[string][]string)
	result := make(map[string][]string)

	for _, word := range words {
		// Приведем слово к нижнему регистру и отсортируем буквы
		word = strings.ToLower(word)
		sortedWord := sortString(word)

		if _, ok := myMap[sortedWord]; ok {
			// Если ключ уже существует, то добавим слово в массив:
			myMap[sortedWord] = append(myMap[sortedWord], word)
		} else {
			// Если ключ не существует, то создадим новый массив с текущим словом в качестве первого элемента
			myMap[sortedWord] = []string{word}
		}
	}

	for _, v := range myMap {
		if len(v) > 1 {
			result[v[0]] = v[1:]
		}
	}

	for _, v := range result {
		sort.Strings(v)
	}

	return result
}

func main() {
	words := []string{"ток", "пятак", "тяпка", "кот", "столик", "листок", "пятка", "слиток", "молоток"}

	result := getAnagram(words)

	for key, v := range result {
		fmt.Printf("Key: %v, value: %v\n", key, v)
	}
}
