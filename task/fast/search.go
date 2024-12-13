package fast

import (
	"bufio"
	"fmt"
	easyjson "github.com/mailru/easyjson"
	"io"
	"os"
	"strings"
)

type User struct {
	Browsers []string
	Email    string
	Name     string
}

func Search(filePath string, out io.Writer) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	seenBrowsers := []string{}
	uniqueBrowsers := 0

	fmt.Fprintln(out, "found users:")
	scanner := bufio.NewScanner(file)
	for i := 0; scanner.Scan(); i++ {
		line := scanner.Bytes()
		user := &User{}
		// fmt.Printf("%v %v\n", err, line)

		// easyjson -all -snake_case fast/search.go
		// https://github.com/mailru/easyjson/issues/10
		err := easyjson.Unmarshal(line, user)

		if err != nil {
			panic(err)
		}

		isAndroid := false
		isMSIE := false

		browsers := user.Browsers

		for _, browser := range browsers {
			if strings.Contains(browser, "Android") {
				isAndroid = true
				notSeenBefore := true
				for _, item := range seenBrowsers {
					if item == browser {
						notSeenBefore = false
					}
				}
				if notSeenBefore {
					// log.Printf("SLOW New browser: %s, first seen: %s", browser, user["name"])
					seenBrowsers = append(seenBrowsers, browser)
					uniqueBrowsers++
				}
			}

			if strings.Contains(browser, "MSIE") {
				isMSIE = true
				notSeenBefore := true
				for _, item := range seenBrowsers {
					if item == browser {
						notSeenBefore = false
					}
				}
				if notSeenBefore {
					// log.Printf("SLOW New browser: %s, first seen: %s", browser, user["name"])
					seenBrowsers = append(seenBrowsers, browser)
					uniqueBrowsers++
				}
			}
		}

		if !(isAndroid && isMSIE) {
			continue
		}

		// log.Println("Android and MSIE user:", user["name"], user["email"])
		email := strings.Replace(user.Email, "@", " [at] ", 1)

		fmt.Fprintln(out, fmt.Sprintf("[%d] %s <%s>", i, user.Name, email))
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	fmt.Fprintln(out)
	fmt.Fprintln(out, "Total unique browsers", len(seenBrowsers))
}
