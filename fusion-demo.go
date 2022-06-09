package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/big"

	"github.com/tyler-smith/go-bip39"
	"github.com/zenon-network/go-zenon/vm/constants"
	"github.com/zenon-network/go-zenon/wallet"
	"github.com/zenon-wiki/go-zdk/client"
	"github.com/zenon-wiki/go-zdk/utils"
	signer "github.com/zenon-wiki/go-zdk/wallet"
	"github.com/zenon-wiki/go-zdk/zdk"
)

func main() {
	mnemonic := "route become dream access impulse price inform obtain engage ski believe awful absent pig thing vibrant possible exotic flee pepper marble rural fire fancy"
	entropy, _ := bip39.EntropyFromMnemonic(mnemonic)
	ks := &wallet.KeyStore{
		Entropy:  entropy,
		Seed:     bip39.NewSeed(mnemonic, ""),
		Mnemonic: mnemonic,
	}
	_, kp, _ := ks.DeriveForIndexPath(0)
	ks.BaseAddress = kp.Address
	s := signer.NewSigner(kp)

	fmt.Println("Using example key:")
	fmt.Println(kp.Address)

	url := "ws://127.0.0.1:35998"
	rpc, err := client.NewClient(url, client.ChainIdentifier(321))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to", url)
	z := zdk.NewZdk(rpc)

	fmt.Println()
	fmt.Println("Let's see how many fusions we have right now")
	fmt.Println("Press Enter to continue")
	fmt.Scanln()

	fmt.Println("Getting current fusions for address (should be none)")
	fmt.Println("Fusions:")
	fusions, err := z.Embedded.Plasma.GetEntriesByAddress(kp.Address, 0, 50)
	for _, f := range fusions.Fusions {
		fmt.Println(f.Id, f.QsrAmount, f.Beneficiary, f.ExpirationHeight)
	}
	fmt.Println()
	fmt.Println("Let's fuse 5000 QSR (the maximum useful amount per address)")
	fmt.Println("Press Enter to continue")
	fmt.Scanln()

	fmt.Println("Making tx template to fuse 5000 QSR")
	tx, err := z.Embedded.Plasma.Fuse(kp.Address, big.NewInt(5000*constants.Decimals))
	if err != nil {
		log.Fatal(err)
	}
	j, _ := json.MarshalIndent(tx, "", "  ")
	fmt.Println(string(j))

	fmt.Println()
	fmt.Println("We now need fill in the template and send the tx. ")
	fmt.Println("Because we don't have any plasma yet, (we're fusing some right now), we will need to generate PoW.")
	fmt.Println("Press Enter to continue")
	fmt.Scanln()

	fmt.Println("Filling out tx. Generating PoW to send tx. This may take a while (using golang pow)")
	tx, err = utils.Send(z, tx, s, false)
	fmt.Println("Sent the following tx")
	if err != nil {
		log.Fatal(err)
	}
	j, _ = json.MarshalIndent(tx, "", "  ")
	fmt.Println(string(j))

	fmt.Println()
	fmt.Println("Let's see how many fusions we have right now")
	fmt.Println("Please wait ~20 seconds and then press Enter to continue")
	fmt.Scanln()

	fmt.Println("Fusions:")
	fusions, err = z.Embedded.Plasma.GetEntriesByAddress(kp.Address, 0, 50)
	for _, f := range fusions.Fusions {
		fmt.Println(f.Id, f.QsrAmount, f.Beneficiary, f.ExpirationHeight)
	}

	fmt.Println()
	fmt.Println("Let's fuse some more QSR. 2 more fuses of 5000 QSR. This will be quick since we have plasma now. Don't need PoW")
	fmt.Println("Press Enter to continue")
	fmt.Scanln()

	tx, err = z.Embedded.Plasma.Fuse(kp.Address, big.NewInt(5000*constants.Decimals))
	if err != nil {
		log.Fatal(err)
	}
	tx, err = utils.Send(z, tx, s, false)
	if err != nil {
		log.Fatal(err)
	}
	j, _ = json.MarshalIndent(tx, "", "  ")
	fmt.Println(string(j))

	fmt.Println()
	fmt.Println("That's one. Press Enter to continue.")
	fmt.Scanln()

	tx, err = z.Embedded.Plasma.Fuse(kp.Address, big.NewInt(5000*constants.Decimals))
	if err != nil {
		log.Fatal(err)
	}
	tx, err = utils.Send(z, tx, s, false)
	if err != nil {
		log.Fatal(err)
	}
	j, _ = json.MarshalIndent(tx, "", "  ")
	fmt.Println(string(j))

	fmt.Println()
	fmt.Println("That's two.")
	fmt.Println("Please wait ~20 seconds and press Enter to see how many fusions there are now.")
	fmt.Scanln()

	fmt.Println("Getting current fusions for address")
	fmt.Println("Fusions:")
	fusions, err = z.Embedded.Plasma.GetEntriesByAddress(kp.Address, 0, 50)
	for _, f := range fusions.Fusions {
		fmt.Println(f.Id, f.QsrAmount, f.Beneficiary, f.ExpirationHeight)
	}
	fmt.Println()
}
