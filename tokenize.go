package prose

import (
	"html"
	"regexp"
	"strings"
	"unicode"

	"github.com/mingrammer/commonregex"
	"github.com/willf/pad"
)

// iterTokenizer splits a sentence into words.
type iterTokenizer struct {
}

// newIterTokenizer is a iterTokenizer constructor.
func newIterTokenizer() *iterTokenizer {
	return new(iterTokenizer)
}

var emoticonRE = regexp.MustCompile(`:>|\._\.|\[-:|:X|\(-_-\)|\(\^_\^\)|:-\}|ಠ_ಠ|¯\\\\\(ツ\)/¯|;\)|O\.o|:-\(|\(╯°□°）╯︵┻━┻|:\*|\(-8|\^__\^|8-D|O\.O|\(-;|:-D|=D|v\.v|:o\)|=/|-__-|;D|@_@|8\)|:o|</3|:-x|O_O|:-\)\)|\(-:|\(\._\.\)|V\.V|:-\||0\.o|:->|:p|8-\)|:-0|xDD|>\.>|:\(\)|:1|<33|\)-:|:-p|0\.0|<3|><\(\(\(\*>|\[:|:-\*|-_-|=\)|;_;|:\(\(|ಠ︵ಠ|:P|:\(|>:\(|o\.o|xD|\):|\(=|:\}|:3|;-D|\(¬_¬\)|:-\(\(\(|\(ಠ_ಠ\)|:\)|:0|:-\(\(|v_v|:-\)|o_o|:\)\)|\(:|0_0|:\)\)\)|0_o|o_O|o\.O|:\(\(\(|\(\*_\*\)|O_o|:\]|;-\)|\^___\^|\(>_<\)|\(o:|:-P|:-\)\)\)|:D|o_0|<333|XDD|=3|:-o|:-3|=\(|:O|o\.0|:-X|:\||:-\]|>:o|V_V|\(;|8D|XD|:-/|\^_\^|:-O|<\.<|:/|>\.<|:x|=\|`)
var terminator = regexp.MustCompile(`\b[a-z]{2,}\.$`)
var contractRE = regexp.MustCompile(`([^' ])'([sS]|[mM]|[dD]|ll|LL|re|RE|ve|VE)`)
var endtracrRE = regexp.MustCompile(`([^' ]+)(n't|N'T)`)
var sanitizer = strings.NewReplacer(
	"\u201c", `"`,
	"\u201d", `"`,
	"\u2018", "'",
	"\u2019", "'",
	"&rsquo;", "'",
	"\u2013", "-",
	"\u2014", "-",
	"&mdash;", "-",
	"&ndash;", "-",
	"\r\n", "\n",
	"\r", "\n")

func replace(pat *regexp.Regexp, text string) string {
	replacements := []string{}
	for _, found := range pat.FindAllString(text, -1) {
		key := pad.Right("", len(found), "x")
		if !stringInSlice(found, replacements) {
			text = strings.Replace(text, found, key, -1)
		}
		replacements = append(replacements, found)
	}
	return text
}

func preTokenize(text string) string {
	text = replace(commonregex.LinkRegex, text)
	text = replace(emoticonRE, text)
	return html.UnescapeString(sanitizer.Replace(text))
}

func postTokenize(spans []span, text string) []Token {
	processed := []Token{}
	last := -1
	for _, s := range spans {
		substr := strings.TrimSpace(text[s.begin:s.end])
		if last >= 0 && substr != "-" && substr != "" {
			processed = append(processed, Token{Text: text[last:s.begin]})
			last = -1
		}
		if terminator.MatchString(substr) {
			processed = append(processed, Token{Text: text[s.begin : s.end-1]})
			processed = append(processed, Token{Text: text[s.end-1 : s.end]})
		} else if substr == "-" && last < 0 {
			last = s.begin
		} else if substr != "" && substr != "-" {
			processed = append(processed, Token{Text: text[s.begin:s.end]})
		}

	}
	return processed
}

// tokenize splits a sentence into a slice of words.
func (t iterTokenizer) tokenize(text string) []Token {
	spans := []span{}
	i, j := 0, 0

	for _, uc := range preTokenize(text) {
		space := unicode.IsSpace(uc)
		punct := unicode.IsPunct(uc) && uc != '.' && uc != '\''
		if space || punct {
			substr := text[i:j]
			if contractRE.MatchString(substr) {
				spans = append(spans, span{begin: i, end: i + 1})
				spans = append(spans, span{begin: i + 1, end: j})
			} else if endtracrRE.MatchString(substr) {
				split := i + strings.Index(substr, "'")
				spans = append(spans, span{begin: i, end: split - 1})
				spans = append(spans, span{begin: split - 1, end: j})
			} else {
				spans = append(spans, span{begin: i, end: j})
				if punct {
					spans = append(spans, span{begin: j, end: j + 1})
				}
			}
			i = j + 1
		}
		j++
	}
	spans = append(spans, span{begin: i, end: j})
	return postTokenize(spans, text)
}

var emoticons = []string{
	`:>`, `._.`, `[-:`, `:X`, `(-_-)`, `(^_^)`, `:-}`, `ಠ_ಠ`, `¯\\(ツ)/¯`, `;)`,
	`O.o`, `:-(`, `(╯°□°）╯︵┻━┻`, `:*`, `(-8`, `^__^`, `8-D`, `O.O`, `(-;`,
	`:-D`, `=D`, `v.v`, `:o)`, `=/`, `-__-`, `;D`, `@_@`, `8)`, `:o`, `</3`,
	`:-x`, `O_O`, `:-))`, `(-:`, `(._.)`, `V.V`, `:-|`, `0.o`, `:->`, `:p`,
	`8-)`, `:-0`, `xDD`, `>.>`, `:()`, `:1`, `<33`, `)-:`, `:-p`, `0.0`,
	":`-(", `<3`, `><(((*>`, `[:`, `:-*`, `-_-`, `=)`, `;_;`, `:((`, `ಠ︵ಠ`,
	`:P`, `:(`, `>:(`, `o.o`, `xD`, `):`, `(=`, `:}`, `:3`, `;-D`, `(¬_¬)`,
	`:-(((`, `(ಠ_ಠ)`, `:)`, `:0`, `:-((`, `v_v`, `:-)`, `o_o`, `:))`, `(:`,
	`0_0`, `:)))`, `0_o`, `o_O`, `o.O`, `:(((`, `(*_*)`, `O_o`, ":`-)", `:]`,
	`;-)`, `^___^`, `(>_<)`, `(o:`, `:-P`, `:-)))`, `:D`, `o_0`, `<333`, `XDD`,
	`=3`, `:-o`, `:-3`, `=(`, `:O`, ":`)", `o.0`, `:-X`, `:|`, `:-]`, `>:o`,
	`V_V`, `(;`, `8D`, `XD`, `:-/`, `^_^`, ":`(", `:-O`, `<.<`, `:/`, `>.<`,
	`:x`, `=|`}
