//
// Based on: https://www.w3.org/TR/png/
//

package main

func main() {
	// Parse the args and copy the file
	subcmd, args := ParseCommandLine()
	switch subcmd {
	case "copy":
		doCopy(args)
	case "dump":
		doDump(args)
	}
}
