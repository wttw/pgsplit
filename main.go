package main

import "os"
import "fmt"
import "strings"
import "regexp"
import "io/ioutil"
import "bytes"

func main() {
	sql, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed: %s\n", err.Error())
		os.Exit(1)
	}

	split, err := splitSQL(string(sql))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed: %s\n", err.Error())
		os.Exit(1)
	}

	b := [][]byte{}
	for _, s := range split {
		b = append(b, []byte(s))
	}

	_, err = os.Stdout.Write(bytes.Join(b, []byte{0}))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed: %s\n", err.Error())
		os.Exit(1)
	}
}

func splitSQL(source string) ([]string, error) {
	re := regexp.MustCompile(`(?ims)(?:\$(?:[a-z][a-z0-9]*)?\$|'|;\s*$)`)
	ret := []string{}
	start := 0
	idx := 0
	for {
		// fmt.Printf("idx=%d start=%d len=%d\n", idx, start, len(source))
		if idx >= len(source) {
			ret = append(ret, source[start:])
			return ret, nil
		}
		match := re.FindStringIndex(source[idx:])
		if match == nil {
			// At the end of the script
			ret = append(ret, source[start:])
			return ret, nil
		}
		if source[idx+match[0]] == ';' {
			// End of statement
			ret = append(ret, source[start:idx+match[1]])
			idx = idx + match[1]
			start = idx
			continue
		}
		close := strings.Index(source[idx+match[1]:], source[idx+match[0]:idx+match[1]])
		// fmt.Printf("%s found at %d, close=%d idx=%d m[0]=%d m[1]=%d snip=%s\n", source[idx+match[0]:idx+match[1]], idx+match[0], close, idx, match[0], match[1], source[idx+match[1]:idx+match[1]+10])
		if close == -1 {
			return ret, fmt.Errorf("Unmatched quote %s at offset %d",
				source[idx+match[0]:idx+match[1]], idx+match[0])
		}
		idx = idx + match[1] + close + match[1] - match[0]
	}
}
