package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

//go:generate go run generate.go -w data/generatedhelp.json

var out = map[string]interface{}{}
var outputFile = flag.String("w", "generatedhelp.json", "Output filename")

func main() {
	flag.Parse()

	process(out, []string{}, nil /* no parent */, "" /* no short help */)

	fl, err := os.Create(*outputFile)
	if err != nil {
		log.Fatalln("Err opening output file:", err)
	}
	defer fl.Close()

	enc := json.NewEncoder(fl)
	enc.SetIndent("", "  ")
	enc.Encode(out)
}

var reShort = regexp.MustCompile(``)
var reUsage = regexp.MustCompile(`(?m)Usage:\n\s+(.*)\n\nAvail`)
var reLong = regexp.MustCompile(`(?ms)(.*)\n\nUsage:`)
var reFlags = regexp.MustCompile(`(?msU)Flags:\n(.*)\n\n`)
var reHelpFlag = regexp.MustCompile(`\n?(\s+-h, --help.*\n?)`)
var reCommands = regexp.MustCompile(`(?msU)Available Commands:\n(.*)\n\n`)

func process(out map[string]interface{}, args []string, parent []string, shortHelp string) {
	fmt.Println("Processing eosc", args, "--help")
	c := exec.Command("eosc", append(args, "--help")...)

	entry := map[string]interface{}{}
	if shortHelp != "" {
		entry["short"] = shortHelp
	}

	if parent != nil {
		entry["parent"] = dashName(parent)
	}

	cnt, err := c.Output()
	if err != nil {
		log.Fatalln("Error running eosc", args, ":", err)
	}
	lines := string(cnt)

	if match := reUsage.FindStringSubmatch(lines); match != nil {
		entry["usage"] = match[1]
	}

	if match := reLong.FindStringSubmatch(lines); match != nil {
		if len(match[1]) != 0 {
			entry["long"] = match[1]
		}
	}

	if match := reFlags.FindStringSubmatch(lines); match != nil {
		var flags = match[1]
		helpFlag := reHelpFlag.FindStringSubmatch(flags)
		if helpFlag != nil {
			flags = strings.Replace(flags, helpFlag[1], "", 1)
		}

		if strings.TrimSpace(flags) != "" {
			entry["flags"] = flags
		}
	}

	match := reCommands.FindStringSubmatch(lines)
	if match != nil {
		for _, cmd := range strings.Split(match[1], "\n") {
			chunks := strings.SplitN(strings.TrimSpace(cmd), " ", 2)
			var cmdName, shortHelp string
			if len(chunks) == 2 {
				shortHelp = strings.TrimSpace(chunks[1])
			}
			cmdName = chunks[0]
			process(out, append(args, cmdName), args, shortHelp)
		}
	}

	out[dashName(args)] = entry
}

func dashName(in []string) string {
	if len(in) == 0 {
		return "eosc"
	}
	return "eosc_" + strings.Join(in, "_")
}
