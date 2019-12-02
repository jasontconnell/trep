package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type InMatchReplace struct {
	In  string
	Out string
}

func main() {
	dir := flag.String("d", ".", "directory")
	e := flag.String("e", "txt", "extension to include")
	reg := flag.String("r", "(.*)", "the regex")
	test := flag.String("t", ".*", "include files that match this regex")
	rep := flag.String("p", "%s[0]", "the replace, use go sprintf style tags")
	imr := flag.String("i", "", "in match replacements, comma separated (for each group matched. for instance, replace ' with '' if using in sql   syntax  ':'') ")
	help := flag.Bool("h", false, "help - print help and exit")
	flag.Parse()

	if *help {
		flag.PrintDefaults()
		return
	}

	var inMatchReplacements []InMatchReplace
	if *imr != "" {
		imrs := strings.Split(*imr, ",")

		for _, reps := range imrs {
			i := strings.Split(reps, ":")
			inMatchReplacements = append(inMatchReplacements, InMatchReplace{In: i[0], Out: i[1]})
		}
	}

	output := os.Stdout

	fullDir := *dir
	if !filepath.IsAbs(fullDir) {
		wd, _ := os.Getwd()
		fullDir = filepath.Join(wd, fullDir)
	}

	contents, err := run(*dir, *e, *reg, *test)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	regex := regexp.MustCompile(*reg)
	for _, c := range contents {
		matches := regex.FindAllStringSubmatch(c, -1)
		for _, m := range matches {
			repary := []string{}
			for i := 1; i < len(m); i++ {

				s := m[i]
				for _, imrep := range inMatchReplacements {
					s = strings.Replace(s, imrep.In, imrep.Out, -1)
				}

				repary = append(repary, s)
			}

			var prms []interface{}
			for _, p := range repary {
				prms = append(prms, p)
			}

			result := fmt.Sprintf(*rep, prms...)

			fmt.Fprintln(output, result)
		}
	}
}

func run(dir, e, reg, test string) ([]string, error) {
	if len(e) > 0 && e[0] != '.' {
		e = "." + e
	}

	fullContents := []string{}
	testreg, terr := regexp.Compile(test)
	if terr != nil {
		return nil, terr
	}

	err := filepath.Walk(dir, func(p string, file os.FileInfo, err error) error {
		if filepath.Ext(p) == e {
			rdr, rerr := os.Open(p)
			if rerr != nil {
				return rerr
			}
			defer rdr.Close()

			contents, cerr := ioutil.ReadAll(rdr)
			if cerr != nil {
				return cerr
			}

			if testreg.Match(contents) {
				fullContents = append(fullContents, string(contents))
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return fullContents, nil
}
