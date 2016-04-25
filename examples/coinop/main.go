package main

import (
	"fmt"
	// "github.com/dappstore/dapp/dapp"
)

// specifying the flag "-dapp.version" causes the binary to output it's version.
// specifying the flag "-dapp.update='stable'" will update the binary to the latest published stable version fo the binary.
// by default the application will output to stderr when there is a new version available, will output a warning when the version has known issues and will exit when the version has been nuked.
// developers use the flag "-dapp.dev" to disable auto update, or perhaps they have a config file specifies their

func main() {
	// dapp.Register("SD427TEBFKFYJOOMFLA723WWKSY7HXSZPG62LX5CL5UA52CVGNVE7AGN")

	// regular program code goes here
	fmt.Println("Hello world!")
}
