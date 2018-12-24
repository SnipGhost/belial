package tasks

import (
	"fmt"
	"strconv"
)

func bitsToStr(data, len uint64) string {
	return fmt.Sprintf("%0"+strconv.FormatUint(len, 10)+"b", data)
}

func getGenPoly(n, k uint64) uint64 {
	switch dif := n - k; dif {
	case 2:
		fallthrough
	case 3:
		fallthrough
	case 4:
		fallthrough
	case 6:
		return (1 << dif) + 3
	case 5:
		return (1 << dif) + 5
	case 7:
		return (1 << dif) + 9
	case 8:
		return (1 << dif) + 29
	default:
		return (1 << dif) + 1
	}
}

func divideCycleCode(r, n, k uint64) uint64 {
	g := getGenPoly(n, k)    // Получаем генераторный полином для n и k
	var x uint64             // Переменная для хранения "выращенного" генераторного полинома
	var c uint64             // Текущий сдвиг деления (с 1, т.к. первый рост полинома на k-1 бит)
	for c = 1; c <= k; c++ { // Пока возможно делить полиномы - делим
		if r >= (1 << (n - c)) { // Проверка на снос нуля, не прошли - пропускаем шаг
			x = g << (k - c) // Выращиваем генераторный полином заполняя нулями справа
			r = r ^ x        // Побитовый XOR, нули в генераторном не влияют на другие разряды
		}
	}
	return r
}

func encryptCycleCode(data, n, k uint64) (uint64, error) {
	// data - информационные биты (входной полином)
	// n - количество выходных бит
	// k - количество информационных бит
	if data >= (1 << k) { // Проверка на то, что степень входного полинома меньше k
		return 0, fmt.Errorf("Data is too big: %b", data)
	}
	result := data << (n - k)          // data * x^(n-k)
	r := divideCycleCode(result, n, k) // Остаток от деления полинома на генераторный
	return result + r, nil
}

func decryptCycleCode(code, n, k uint64) (uint64, uint64, error) {
	// code - закодированные биты (принятый полином)
	// n - количество принятых бит
	// k - количество информационных бит
	if code >= (1 << n) { // Проверка на то, что степень принятого полинома меньше n
		return 0, 0, fmt.Errorf("Code is too big: %b", code)
	}
	vec := divideCycleCode(code, n, k) // Остаток от деления полинома на генераторный
	r := vec                           // Сохраняем в отдельную переменную для экспериментов
	if r != 0 {                        // Остаток != 0, => есть ошибка, пробуем исправить
		var err uint64
		var errRem uint64
		for err = 1; err <= code; err = err << 1 {
			errRem = divideCycleCode(err, n, k)
			if errRem == r {
				code = code ^ err
				r = divideCycleCode(code, n, k)
				err = code + 1
			}
		}
	}
	code = code >> (n - k)
	return code, vec, nil
}

func encryptWrapper(dataStr string, codeLen, infLen, codeType uint64) (uint64, error) {
	if codeType != 0 {
		return 0, fmt.Errorf("Not implemented")
	}
	data, err := strconv.ParseUint(dataStr, 2, 64)
	if err != nil {
		return 0, fmt.Errorf("Can't parse dataStr: %s", err)
	}
	res, err := encryptCycleCode(data, codeLen, infLen)
	if err != nil {
		return 0, fmt.Errorf("Can't encrypt data: %s", err)
	}
	// fmt.Println("Source:", dataStr, "Encrypted:", bitsToStr(res, codeLen))
	return res, nil
}

func decryptWrapper(code, codeLen, infLen, codeType uint64) (uint64, uint64, error) {
	if codeType != 0 {
		return 0, 0, fmt.Errorf("Not implemented")
	}
	res, rem, err := decryptCycleCode(code, codeLen, infLen)
	if err != nil {
		return 0, 0, fmt.Errorf("Can't decrypt data: %s", err)
	}
	// if rem != 0 {
	// 	fmt.Println("Error detected")
	// }
	// fmt.Println(bitsToStr(res, infLen))
	return res, rem, nil
}

func weightBits(v uint64) uint64 {
	// Also you can use Brian Kernighan's algorithm that takes O(log N)
	v = v - ((v >> 1) & 0x5555555555555555)
	v = (v & 0x3333333333333333) + ((v >> 2) & 0x3333333333333333)
	c := (((v + (v >> 4)) & 0x0F0F0F0F0F0F0F0F) * 0x101010101010101) >> 56
	return c
}

func factorial(n uint64) uint64 {
	var factVal uint64 = 1
	var i uint64
	for i = 1; i <= n; i++ {
		factVal *= i
	}
	return factVal
}

func countCombinations(n, k uint64) uint64 {
	return factorial(n) / (factorial(k) * factorial(n-k))
}
