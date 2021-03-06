package main

import (
	"fmt"
	. "github.com/dappstore/go-dapp/app"
	"log"
)

// specifying the flag "-dapp.version" causes the binary to output it's version.
// specifying the flag "-dapp.update='stable'" will update the binary to the latest published stable version fo the binary.
// by default the application will output to stderr when there is a new version available, will output a warning when the version has known issues and will exit when the version has been nuked.
// developers use the flag "-dapp.dev" to disable auto update, or perhaps they have a config file specifies their

func main() {

	app, err := New("GDGIXJPUTJIYHHJ2TYWO2HJMFNT7M767ZB33SFGTD77JUE3YZ6YZBUD4",
		Defaults,
		Developer("GA6AJ6WPO6BDFUKUJKPDW3SILWSXLP62O72JTY3JDUJVR2EMIOBMJDLM"),
		Name("coinop"),
		Description("an example stellar integration"),
	)
	if err != nil {
		log.Fatal(err)
	}

	p, err := app.WaitForPayment()
	if err != nil {
		log.Fatal(err)
	}

	// regular program code goes here
	fmt.Println("Hello world!")
	fmt.Println(p)
}
