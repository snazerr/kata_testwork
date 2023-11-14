package main

import (
	"fmt"
	"strings"
)

/* В функции main выводим предложение ввести математическое выражение.
Объявляем переменную input для хранения введенной строки.
Считываем строку из ввода пользователя. Если при этом возникает ошибка, она сохраняется в err.
Вызываем функцию calculate с введенной строкой. Результат и возможная ошибка сохраняются в result и err.
Выводит результат в консоль. Если при вычислении произошла ошибка, выводит ее.
*/

func main() {
	fmt.Println("Привет! Введи математическое выражение без использования пробелов (например, 1+2):")
	var input string
	_, err := fmt.Scanln(&input)
	if err != nil {
		fmt.Println("Ошибка ввода:", err)
		return
	}

	result, err := calculate(input)
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}

	fmt.Println("Результат:", result)
}

/* В функции calculate:
Удаляем пробелы из введенной строки, чтобы обработать выражения с пробелами и без.
Определяем допустимые операторы.
Ищем первый оператор в строке и сохраняет его индекс в operatorIndex.
Проверяем корректность формата выражения и определяются операнды и оператор.
Проверяем, используются ли римские числа.
Проверяем, используются ли арабские числа.
Если используются римские числа, они преобразуются в арабские.
Вызывает функцию performOperation для выполнения операции.
Если использовались римские числа, результат преобразуется обратно. Возвращается результат и возможная ошибка.
*/

func calculate(input string) (string, error) {
	// Убираем пробелы из входной строки
	input = strings.ReplaceAll(input, " ", "")

	// Операторы, которые поддерживаются (+, -, *, /)
	operators := "+-*/"
	operatorIndex := -1

	// Находим индекс оператора в строке
	for i, char := range input {
		if strings.ContainsRune(operators, char) {
			if operatorIndex != -1 {
				return "", fmt.Errorf("Ты ввел больше двух операндов...такое я пока что не умею считать")
			}
			operatorIndex = i
		}
	}

	// Проверяем корректность расположения оператора
	if operatorIndex == -1 || operatorIndex == 0 || operatorIndex == len(input)-1 {
		return "", fmt.Errorf("Некорректный формат выражения")
	}

	// Выделяем операнды и оператор
	operand1, operator, operand2 := input[:operatorIndex], string(input[operatorIndex]), input[operatorIndex+1:]

	// Проверяем, являются ли операнды римскими числами и числовыми значениями
	isRoman := isRomanNumeral(operand1) && isRomanNumeral(operand2)
	isNumeric := isNumeric(operand1) && isNumeric(operand2)

	// Проверяем, используется ли одна система счисления (римская или арабская)
	if !isRoman && !isNumeric {
		return "", fmt.Errorf("Ты по-моему перепутал...такое я не посчитаю")
	}

	var a, b int

	// Логика для работы с римскими и арабскими числами
	if isRoman {
		a, _ = romanToArabic(operand1)
		b, _ = romanToArabic(operand2)
	} else if isNumeric {

		// Сканируем операнды как целые числа
		fmt.Sscan(operand1, &a)
		fmt.Sscan(operand2, &b)

		// Проверка чисел на соответствие диапазону
		a, err := strconv.Atoi(operand1)
		if err != nil || a < 1 || a > 10 {
			return "", fmt.Errorf("Число должно быть в диапазоне от 1 до 10")
		}

		b, err := strconv.Atoi(operand2)
		if err != nil || b < 1 || b > 10 {
			return "", fmt.Errorf("Число должно быть в диапазоне от 1 до 10")
		}
	}

	result, err := performOperation(a, b, operator)
	if err != nil {
		return "", err
	}

	if isRoman {
		return arabicToRoman(result), nil
	}

	return fmt.Sprint(result), nil
}

/*В функции performeOperation:
Выполняем математическую операцию в зависимости от оператора.
Обрабатываем случай деления на ноль и неизвестного оператора, а также случай, где результат меньше или равен 0.
*/

func performOperation(a, b int, operator string) (int, error) {
	switch operator {
	case "+":
		return a + b, nil
	case "-":
		result := a - b
		if result <= 0 {
			return 0, fmt.Errorf("Результат меньше или равен нулю")
		}
		return result, nil
	case "*":
		return a * b, nil
	case "/":
		if b == 0 {
			return 0, fmt.Errorf("Деление на ноль")
		}
		return a / b, nil
	default:
		return 0, fmt.Errorf("Неизвестный оператор")
	}
}

/*В функции isRomanNumeral:
Задаем значения корректных римских чисел.
Проверяем, является ли строка корректным римским числом.
*/

func isRomanNumeral(s string) bool {
	validChars := "IVXLCDM"
	invalidPatterns := []string{"IIII", "VV", "XXXX", "LL", "CCCC", "DD", "MMMM"}

	for _, pattern := range invalidPatterns {
		if strings.Contains(s, pattern) {
			return false
		}
	}

	for _, char := range s {
		if !strings.ContainsRune(validChars, char) {
			return false
		}
	}

	return true
}

/*В функции arabicToRoman:
преобразуем целое число, представляющее арабскую систему счисления, в эквивалентное ему римское число
*/

func arabicToRoman(n int) string {
	if n <= 0 || n > 3999 {
		return ""
	}

	// Значения и символы для римских чисел
	vals := []int{1000, 900, 500, 400, 100, 90, 50, 40, 10, 9, 5, 4, 1}
	roman := []string{"M", "CM", "D", "CD", "C", "XC", "L", "XL", "X", "IX", "V", "IV", "I"}

	result := ""

	// Идем по каждому значению и добавляем соответствующие символы к результату
	for i := 0; n > 0; i++ {
		for n >= vals[i] {
			n -= vals[i]
			result += roman[i]
		}
	}
	return result
}

/*В функции romanToArabic:
преобразуем строку, представляющую римское число, в эквивалентное ему целое число в арабской системе счисления.
*/

func romanToArabic(s string) (int, error) {
	// Мапа, сопоставляющая римские символы их десятичным значениям
	romanNumerals := map[rune]int{
		'I': 1,
		'V': 5,
		'X': 10,
		'L': 50,
		'C': 100,
		'D': 500,
		'M': 1000,
	}

	total := 0     //представляет общее арабское значение числа
	prevValue := 0 //хранит значение предыдущего символа в римской записи

	// Идем по строке справа налево
	for i := len(s) - 1; i >= 0; i-- {
		value := romanNumerals[rune(s[i])]

		// Если значение текущего символа меньше предыдущего, вычитаем его
		if value < prevValue {
			total -= value
		} else {
			// Иначе прибавляем его к общей сумме
			total += value
		}
		prevValue = value
	}

	return total, nil
}

/*
Функция isNumeric выполняет проверку, является ли переданная строка числом.
Здесь она используется для определения, представляют ли операнды числовые значения в выражении.
*/

func isNumeric(s string) bool {
	_, err := fmt.Sscan(s, new(int))
	return err == nil
}
