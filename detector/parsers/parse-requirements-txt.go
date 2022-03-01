package parsers

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

const PipEcosystem Ecosystem = "PyPI"

// todo: expand this to support more things, e.g.
//   https://pip.pypa.io/en/stable/reference/requirements-file-format/#example
func parseLine(line string) (PackageDetails, error) {
	var constraint string
	name := line

	version := "0.0.0"

	if strings.Contains(line, "==") {
		constraint = "=="
	}

	if strings.Contains(line, ">=") {
		constraint = ">="
	}

	if strings.Contains(line, "~=") {
		constraint = "~="
	}

	if strings.Contains(line, "!=") {
		constraint = "!="
	}

	if constraint != "" {
		splitted := strings.Split(line, constraint)

		name = strings.TrimSpace(splitted[0])

		if constraint != "!=" {
			version = strings.TrimSpace(splitted[1])
		}
	}

	return PackageDetails{
		Name:      name,
		Version:   version,
		Ecosystem: PipEcosystem,
	}, nil
}

func removeComments(line string) string {
	var re = regexp.MustCompile(`(^|\s+)#.*$`)

	return strings.TrimSpace(re.ReplaceAllString(line, ""))
}

func isNotRequirementLine(line string) bool {
	return line == "" ||
		// flags are not supported
		strings.HasPrefix(line, "-") ||
		// file urls
		strings.HasPrefix(line, "https://") ||
		strings.HasPrefix(line, "http://") ||
		// file paths are not supported (relative or absolute)
		strings.HasPrefix(line, ".") ||
		strings.HasPrefix(line, "/")
}

func ParseRequirementsTxt(pathToLockfile string) ([]PackageDetails, error) {
	var packages []PackageDetails

	file, err := os.Open(pathToLockfile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := removeComments(scanner.Text())

		if isNotRequirementLine(line) {
			continue
		}

		detail, err := parseLine(line)

		if err != nil {
			fmt.Printf("Was unable to parse line '%s'", line)

			continue
		}

		packages = append(packages, detail)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return packages, nil
}