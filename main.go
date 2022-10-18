package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

var tobeDeleted = []int{}
var punctuation = []string{".", ",", "!", "?", ":", ";"}
var vowels = []string{"a", "e", "i", "o", "u", "h"}
var upslowsCupsHexBinWnoNum = []string{"(cap)", "(low)", "(up)", "(hex)", "(bin)"}
var upsLowCupsWNum = []string{"(cap,", "(low,", "(up,"}
var isTrueq bool

func getFileData(s string) ([]string, error) {
	arr, err := os.ReadFile(s)
	return strings.Split(string(arr), " "), err
}
func isTrue(word string, is string) bool {
	return strings.TrimSpace(word) == is
}
func RemoveIndex(s []string, index int) []string {
	return append(s[:index], s[index+1:]...)
}

func getbinHext(word string, base int) (string, error) {
	num, err := strconv.ParseInt(word, base, 32)
	return strconv.Itoa(int(num)), err
}
func mutateUpdateSigleUpLowCapHexBin(index int, s []string, word string) {
	prevIndex := s[index-1]
	switch word {
	case "(cap)":
		s[index-1] = strings.Title(prevIndex)
		tobeDeleted = append(tobeDeleted, index)
	case "(up)":
		s[index-1] = strings.ToUpper(prevIndex)
		tobeDeleted = append(tobeDeleted, index)
	case "(low)":
		s[index-1] = strings.ToLower(prevIndex)
		tobeDeleted = append(tobeDeleted, index)
	case "(bin)":
		val, err := getbinHext(prevIndex, 2)
		if err == nil {
			s[index-1] = val
		}
		tobeDeleted = append(tobeDeleted, index)
	case "(hex)":
		val, err := getbinHext(prevIndex, 16)
		if err == nil {
			s[index-1] = val
		}
		tobeDeleted = append(tobeDeleted, index)
	}
}

func CuplowupWithnummuteSubArray(arr []string, t func(s string) string) {
	for i, cur := range arr {
		arr[i] = t(cur)
	}
}

func mutateCuplowupWithnum(word string, index int, s []string) {
	number, err := strconv.ParseInt(strings.ReplaceAll(word, ")", ""), 10, 32)
	arr := s[index-int(number)-1 : index-1]
	if err == nil {
		switch s[index-1] {
		case "(cap,":
			CuplowupWithnummuteSubArray(arr, strings.Title)
		case "(low,":
			CuplowupWithnummuteSubArray(arr, strings.ToLower)
		case "(up,":
			CuplowupWithnummuteSubArray(arr, strings.ToUpper)
		}
	}
	tobeDeleted = append(tobeDeleted, index-1)
	tobeDeleted = append(tobeDeleted, index)

}

func filter[T any](slice []T, f func(T) bool) []T {
	var n []T
	for _, e := range slice {
		if f(e) {
			n = append(n, e)
		}
	}
	return n
}

func Any[T any](slice []T, f func(T) bool) bool {

	for _, e := range slice {
		if f(e) {
			return true
		}
	}
	return false
}

func AnyConditionVal(word []string, t string) bool {
	return Any(word, func(val string) bool {
		if val == t {
			return true
		}
		return false
	})
}

func mutatePunctuation(word string, index int, s []string) {

	tBPlI := filter(strings.Split(word, ""), func(v string) bool {
		return (v == "." || v == "?" || v == "," || v == "!" || v == ":" || v == ";")
	})

	subS := strings.Join(tBPlI, "")

	repSubS := strings.Replace(strings.TrimSpace(word), subS, "", 1)
	im2g := index-2 >= 0
	deletS := strings.Trim(strings.Replace(fmt.Sprint(tobeDeleted), " ", ",", -1), "[]")

	switch {
	case index-3 >= 0 && im2g && AnyConditionVal(strings.Split(deletS, ","), strconv.Itoa(index-1)) && AnyConditionVal(strings.Split(deletS, ","), strconv.Itoa(index-2)):
		s[index-3] += subS
	case im2g && AnyConditionVal(strings.Split(deletS, ","), strconv.Itoa(index-1)) && !AnyConditionVal(strings.Split(deletS, ","), strconv.Itoa(index-2)):
		s[index-2] += subS
	default:
		s[index-1] += subS
	}
	if repSubS != "" {
		s[index] = repSubS
	} else {
		tobeDeleted = append(tobeDeleted, index)
	}
}

func sigleQ(word string, index int, s []string) {
	if !isTrueq {
		if len(s)-1 >= index+1 {
			// fmt.Println(s[index+1])
			s[index+1] = "'" + s[index+1]
			s[index] = ""
		}
		tobeDeleted = append(tobeDeleted, index)
		isTrueq = !isTrueq
	} else {
		deletS := strings.Trim(strings.Replace(fmt.Sprint(tobeDeleted), " ", ",", -1), "[]")

		fmt.Println(s[index-1])

		if AnyConditionVal(strings.Split(deletS, ","), strconv.Itoa(index-1)) {
			s[index-2] += s[index]
		} else {
			s[index-1] += s[index]
		}
		tobeDeleted = append(tobeDeleted, index)
	}

}

func MutateAll(word string, index int, s []string) {
	// fmt.Println(string(word[0]))
	isPrevExists := index-1 >= 0
	switch {
	case isPrevExists && AnyConditionVal(upsLowCupsWNum, s[index-1]) && strings.Contains(word, ")"):
		mutateCuplowupWithnum(word, index, s)
	case isPrevExists && AnyConditionVal(upslowsCupsHexBinWnoNum, strings.TrimSpace(word)):
		mutateUpdateSigleUpLowCapHexBin(index, s, word)
	case isPrevExists && word != "" && AnyConditionVal(punctuation, string(word[0])):
		mutatePunctuation(word, index, s)
	case isPrevExists && word == "'":
		sigleQ(word, index, s)
	case isPrevExists && isTrue(strings.ToLower(s[index-1]), "a") && AnyConditionVal(vowels, string(word[0])):
		s[index-1] += "n"
	}
}

func moDify(s []string) {
	for i, curr := range s {
		MutateAll(curr, i, s)
	}
}

func main() {
	arrgs := os.Args[1:]
	if len(arrgs) != 2 {
		return
	}
	str, err := getFileData(arrgs[0])
	if err != nil {
		fmt.Println(err)
	}
	moDify(str)
	fmt.Println(str)
	for i, cur := range tobeDeleted {
		str = RemoveIndex(str, cur-i)
	}

	fmt.Println(str)
	output := []byte(string(strings.Join(str, " ")))

	errorWriteFile := ioutil.WriteFile(arrgs[1], output, 0644)
	if err != nil {
		fmt.Println(errorWriteFile)
	}

}
