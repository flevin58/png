//
// Based on: https://www.w3.org/TR/png/
//

package main

import "github.com/flevin58/png/cmds"

func main() {
	// Parse the args and copy the file
	subcmd, args := cmds.ParseCommandLine()
	switch subcmd {
	case "copy":
		cmds.DoCopy(args)
	case "dump":
		cmds.DoDump(args)
	}
}
