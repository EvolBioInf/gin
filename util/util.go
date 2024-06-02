// The package util collects utility functions for the package gin.
package util

import (
	"fmt"
	"github.com/evolbioinf/clio"
	"log"
	"os"
)

var version, date string
var name string

// Version prints program information and exits.
func Version() {
	author := "Bernhard Haubold"
	email := "haubold@evolbio.mpg.de"
	license := "Gnu General Public License, " +
		"https://www.gnu.org/licenses/gpl.html"
	clio.PrintInfo(name, version, date,
		author, email, license)
	os.Exit(0)
}

// Open opens a file with error checking.
func Open(file string) *os.File {
	f, err := os.Open(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "couldn't open %s\n",
			file)
		os.Exit(1)
	}
	return f
}

// Name sets the program name. It also customizes the error messages generated via the log package by prefixing them with the program name.
func Name(n string) {
	name = n
	m := fmt.Sprintf("%s: ", name)
	log.SetPrefix(m)
	log.SetFlags(0)
}

// Check checks an error and exits with message if the error isn't nil.
func Check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}
