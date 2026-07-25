package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"maps"
	"os"
	"slices"
	"strings"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "gtr - translate text from the command line using Google Translate")
		fmt.Fprintln(os.Stderr, "\nUsage:")
		fmt.Fprintln(os.Stderr, "  gtr [flags] <text to translate>")
		fmt.Fprintln(os.Stderr, "\nExamples:")
		fmt.Fprintln(os.Stderr, "  gtr 猫")
		fmt.Fprintln(os.Stderr, "  gtr -t ja cat")
		fmt.Fprintln(os.Stderr, "  gtr -f ja -t en 猫")
		fmt.Fprintln(os.Stderr, "  gtr -codes")
		fmt.Fprintln(os.Stderr, "\nFlags:")
		flag.PrintDefaults()
	}

	from := flag.String("f", "auto", "source language code, or 'auto' to detect (see -codes)")
	to := flag.String("t", "en", "target language code (see -codes)")
	codes := flag.Bool("codes", false, "print all available language codes and exit")
	toJSON := flag.Bool("json", false, "enable json output")
	flag.Parse()

	log.SetFlags(0)
	log.SetPrefix("gtr: ")

	if *codes {
		for lang, code := range langCode {
			fmt.Println(lang+":", code)
		}
		os.Exit(0)
	}

	*from = strings.TrimSpace(strings.ToLower(*from))
	*to = strings.TrimSpace(strings.ToLower(*to))
	message := strings.TrimSpace(strings.Join(flag.Args(), " "))

	langCodes := slices.Collect(maps.Values(langCode))
	isValidCode := func(code string) bool { return slices.Contains(langCodes, code) }

	for k, v := range map[string]string{"from": *from, "to": *to, "<text to translate>": message} {
		if v == "" {
			log.Fatalf("'%s' cannot be empty\n", k)
		}
		if k == "from" || k == "to" {
			if !isValidCode(v) {
				log.Fatalf("invalid '%s' lang code. See '-codes'\n", k)
			}
		}
	}

	tr, err := Translate(*from, *to, message)
	if err != nil {
		log.Fatal("translate error:", err)
	}

	if *toJSON {
		if err := json.NewEncoder(os.Stdout).Encode(tr); err != nil {
			log.Fatal("json encode error:", err)
		}
		os.Exit(0)
	}

	fmt.Println(tr.Text)
	if tr.POS != "" {
		fmt.Println()
		fmt.Print(tr.POS)
	}
	if tr.Def != "" {
		fmt.Println()
		fmt.Print(tr.Def)
	}
}
