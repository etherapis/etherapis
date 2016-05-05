package dashboard

import (
	"fmt"

	"github.com/etherapis/etherapis/etherapis"
	"github.com/etherapis/etherapis/etherapis/geth/web3ext"
	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/jsre"
)

// newREPLConsole creates a new JavaScript interpreter console that can be used
// by the API to serve fronted console requests.
func newREPLConsole(eapis *etherapis.EtherAPIs) *jsre.JSRE {
	// Create a JavaScript interpreter with web3 injected
	console := jsre.New("")

	client, _ := eapis.Geth().Stack().Attach()

	jeth := utils.NewJeth(console, client)
	console.Set("jeth", struct{}{})
	t, _ := console.Get("jeth")
	jethObj := t.Object()

	jethObj.Set("send", jeth.Send)
	jethObj.Set("sendAsync", jeth.Send)

	console.Compile("bignumber.js", jsre.BigNumber_JS)
	console.Compile("web3.js", jsre.Web3_JS)
	console.Run("var Web3 = require('web3');")
	console.Run("var web3 = new Web3(jeth);")

	// Inject all of the APIs exposed by the inproc RPC client
	shortcut := "var eth = web3.eth; var personal = web3.personal; "
	apis, _ := client.SupportedModules()
	for api, _ := range apis {
		if api == "web3" || api == "rpc" {
			continue
		}
		if jsFile, ok := web3ext.Modules[api]; ok {
			console.Compile(fmt.Sprintf("%s.js", api), jsFile)
			shortcut += fmt.Sprintf("var %s = web3.%s; ", api, api)
		}
	}
	console.Run(shortcut)

	// Finally return the ready console
	return console
}
