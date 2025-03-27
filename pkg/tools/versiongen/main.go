package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var (
	versionFile = flag.String("file", "", "Path to version file (format: MAJOR.MINOR)")
	mode        = flag.String("mode", "new", "Mode: 'new' (default) or 'last'")
)

func main() {
	flag.Parse()

	if *versionFile == "" {
		fmt.Fprintln(os.Stderr, "Missing -file argument")
		os.Exit(1)
	}

	// Читаем MAJOR.MINOR из файла
	data, err := os.ReadFile(*versionFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read version file: %v\n", err)
		os.Exit(1)
	}

	parts := strings.Split(strings.TrimSpace(string(data)), ".")
	if len(parts) != 2 {
		fmt.Fprintln(os.Stderr, "Version file must contain MAJOR.MINOR (e.g. 1.3)")
		os.Exit(1)
	}

	major := parts[0]
	minor := parts[1]

	// Получаем все git-теги
	cmd := exec.Command("git", "tag")
	output, err := cmd.Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get git tags: %v\n", err)
		os.Exit(1)
	}

	re := regexp.MustCompile(`^v` + major + `\.` + minor + `\.(\d+)$`)
	var patchVersions []int

	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := scanner.Text()
		matches := re.FindStringSubmatch(line)
		if len(matches) == 2 {
			patch, _ := strconv.Atoi(matches[1])
			patchVersions = append(patchVersions, patch)
		}
	}

	sort.Ints(patchVersions)
	lastPatch := 0
	if len(patchVersions) > 0 {
		lastPatch = patchVersions[len(patchVersions)-1]
	}

	if *mode == "last" {
		fmt.Printf("v%s.%s.%d\n", major, minor, lastPatch-1)
	} else {
		fmt.Printf("v%s.%s.%d\n", major, minor, lastPatch)
	}
}
