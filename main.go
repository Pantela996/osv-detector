package main

import (
	"flag"
	"fmt"
	"github.com/fatih/color"
	"os"
	"osv-detector/detector/database"
	"osv-detector/detector/parsers"
	"path"
)

func loadOSVDatabase(offline bool) database.OSVDatabase {
	db, err := database.NewDB(offline, database.GithubOSVDatabaseArchiveURL)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to load the OSV DB: %s\n", err)
		os.Exit(127)
	}

	fmt.Printf("Loaded %d vulnerabilities\n", len(db.Vulnerabilities(true)))

	return *db
}

func printEcosystems(db database.OSVDatabase) {
	ecosystems := db.ListEcosystems()

	fmt.Print("The loaded OSV has vulnerabilities for the following ecosystems:")

	for _, ecosystem := range ecosystems {
		fmt.Printf("  %s\n", ecosystem)
	}
}

func printVulnerabilities(
	db database.OSVDatabase,
	ecosystem database.Ecosystem,
	pkg parsers.EcosystemPackage,
) int {
	// fmt.Printf("%s: %s@%s\n", ecosystem, pkg.Name, pkg.Version)
	vulnerabilities := db.VulnerabilitiesAffectingPackage(ecosystem, pkg.Name, pkg.Version)

	if len(vulnerabilities) == 0 {
		return 0
	}

	fmt.Printf(
		"%s %s\n",
		color.YellowString("%s@%s", pkg.Name, pkg.Version),
		color.RedString("is affected by the following vulnerabilities:"),
	)

	for _, vulnerability := range vulnerabilities {
		fmt.Printf(
			"  %s %s\n",
			color.CyanString("%s:", vulnerability.ID),
			vulnerability.Summary,
		)
	}

	return len(vulnerabilities)
}

func main() {
	offline := flag.Bool("offline", false, "Update the OSV database")
	listEcosystems := flag.Bool("list-ecosystems", false, "List all the ecosystems present in the loaded OSV database")

	flag.Parse()
	pathToLockOrDirectory := flag.Arg(0)
	fmt.Println("Hello, world.")

	db := loadOSVDatabase(*offline)

	if *listEcosystems {
		printEcosystems(db)
		os.Exit(0)
	}

	out, err := parsers.ParseComposerLock(pathToLockOrDirectory)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing %s: %s\n", pathToLockOrDirectory, err)
		os.Exit(127)
	}

	file := path.Base(pathToLockOrDirectory)

	knownVulnerabilitiesCount := 0
	for _, pkg := range out.Packages {
		knownVulnerabilitiesCount += printVulnerabilities(db, database.Ecosystem(out.Ecosystem), pkg)
	}

	if knownVulnerabilitiesCount == 0 {
		fmt.Printf("%s\n", color.GreenString("%s has no known vulnerabilities!", file))
		os.Exit(0)
	}

	fmt.Printf("\n%s\n", color.RedString(
		"%s is affected by %d vulnerabilities!",
		file,
		knownVulnerabilitiesCount,
	))
}
