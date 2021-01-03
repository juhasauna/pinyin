package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"
)

func main() {
	// inputPath := "input.txt"
	inputPath := "C:\\Users\\FIJUSAU\\OneDrive - ABB\\cn\\xhzd\\xhzd5_rmTradHanzi.txt"

	input, err := os.Open(inputPath)
	if err != nil {
		log.Fatal(err)
	}
	defer input.Close()

	output, err := os.Create("output.txt")
	if err != nil {
		log.Fatal(err)
	}

	reps := pinyins()

	prev := ""

	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		s := scanner.Text()
		if s == prev {
			continue
		}
		prev = s

		w, p, r := splitText(s)
		p = strings.ReplaceAll(p, "－", " ")
		p2 := ""
		for _, v := range strings.Split(p, " ") {
			if len(v) == 0 {
				log.Fatalln("empty pinyin, ", s)
				continue
			}
			if end := v[len(v)-1]; end == '1' || end == '2' || end == '3' || end == '4' {
				p2 += " " + v
				continue // pinyin already in correct format
			}
			v += pinyinEndIndicator
			p2 += replaceIt(v, reps, pinyinEndsWithErHua(v))
		}
		if len(p2) == 0 {
			log.Fatalf("didn't find pinyin for: %s\n", s)
		}
		if string(p2[0]) == " " {
			p2 = p2[1:]
		}

		s = w + p2 + r

		_, err := output.WriteString(s + "\n")
		if err != nil {
			log.Fatal(err)
		}

	}
}

func pinyinEndsWithErHua(s string) bool {
	endPos := strings.Index(s, pinyinEndIndicator)
	if endPos == -1 {
		log.Fatalln("didn't find pinyinEndIndicator")
	}
	if endPos < 2 {
		return false
	}
	endPos -= 1

	if string(s[endPos]) == "r" {
		for _, v := range []string{"e", "è", "ě", "é", "ē"} {
			_, size := utf8.DecodeRuneInString(v)
			if string(s[endPos-size:endPos]) == v {
				return false
			}
		}
		return true
	}
	return false
}

const pinyinEndIndicator = "!#/)"

func splitText(s string) (string, string, string) {
	sepStr1 := "`2`"
	sep1 := strings.Index(s, sepStr1) + len(sepStr1)
	sep2 := strings.Index(s, "<")
	if sep2 < sep1 {
		return s[:sep1], s[sep1:], ""
	}
	word, pinyin, rest := s[:sep1], s[sep1:sep2], s[sep2:]
	return word, pinyin, rest
}

func replaceIt(txt string, reps []replacement, endsWithr bool) string {

	for _, v := range reps {
		if l := len(v.new); endsWithr && l > 0 {
			if l == 1 {
				v.new = v.new + "r"
			} else if _, err := strconv.Atoi(string(v.new[l-1])); err == nil {
				v.new = string(v.new[:l-1]) + "r" + string(v.new[l-1])
			} else {
				v.new = v.new + "r"
			}
			v.old = v.old + "r"
		}
		newTxt := strings.Replace(txt, v.old+pinyinEndIndicator, pinyinEndIndicator+" "+v.new, 1)
		if newTxt != txt {
			txt = replaceIt(newTxt, reps, pinyinEndsWithErHua(newTxt))
		}
	}
	return strings.Replace(txt, pinyinEndIndicator, "", 1)
}

func rev(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

type replacement struct {
	old string
	new string
}

func rmTradHanzi(s string) string {
	startString := "`2`"
	start := strings.Index(s, startString)

	if start < 1 {
		log.Fatalf("didn't find '%s'\n", startString)
	}
	// fmt.Printf("%d\t", start)
	start += len(startString)
	_, size := utf8.DecodeRuneInString(s[start:])
	// fmt.Println(r, string(r), size)
	if size < 3 { // The size of Hanzi is 3
		return s
	}
	end := strings.Index(s, "<br>") + 4

	return s[:start] + s[end:]
}
