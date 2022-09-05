package roman_numerals

import (
	"golang.org/x/exp/slices"
	"strings"
)

type RomanArabicPair struct {
	Arabic int
	Roman  string
}

type RomanNumerals []RomanArabicPair

var romanArabicLookupTable = RomanNumerals{
	{1000, "M"},
	{900, "CM"},
	{500, "D"},
	{400, "CD"},
	{100, "C"},
	{90, "XC"},
	{50, "L"},
	{40, "XL"},
	{10, "X"},
	{9, "IX"},
	{5, "V"},
	{4, "IV"},
	{1, "I"},
}

var subtractiveSymbols = []uint8{'I', 'X', 'C'}

func isSubtractive(symbol uint8) bool {
	return slices.Contains(subtractiveSymbols, symbol)
}

func (r RomanNumerals) Exists(symbols ...byte) bool {
	symbol := string(symbols)
	for _, s := range r {
		if s.Roman == symbol {
			return true
		}
	}
	return false
}

func ConvertToRoman(arabic int) string {
	var result = strings.Builder{}

	for _, pair := range romanArabicLookupTable {
		for arabic >= pair.Arabic {
			result.WriteString(pair.Roman)
			arabic -= pair.Arabic
		}
	}

	return result.String()
}

func (r RomanNumerals) ValueOf(symbols ...byte) int {
	for _, s := range romanArabicLookupTable {
		if s.Roman == string(symbols) {
			return s.Arabic
		}
	}
	return 0
}

func ConvertToArabic(roman string) (total int) {
	for _, symbols := range windowedRoman(roman).Symbols() {
		total += romanArabicLookupTable.ValueOf(symbols...)
	}
	return
}

type windowedRoman string

func (w windowedRoman) Symbols() (symbols [][]byte) {
	for i := 0; i < len(w); i++ {
		symbol := w[i]
		notAtEnd := i+1 < len(w)

		if notAtEnd && isSubtractive(symbol) && romanArabicLookupTable.Exists(symbol, w[i+1]) {
			symbols = append(symbols, []byte{symbol, w[i+1]})
			i++
		} else {
			symbols = append(symbols, []byte{symbol})
		}
	}
	return
}
