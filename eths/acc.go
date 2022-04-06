package eths

import (
	"Go-StandingbookServer/config"
	"Go-StandingbookServer/utils"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"os"
)
import "github.com/ethereum/go-ethereum/rpc"

func NewAcc(pass string, connstr string) (string, error) {
	client, err := rpc.Dial(connstr)
	if err != nil {
		fmt.Println("failed to connect to geth", err)
		return "", err
	}

	//创建账户
	var account string
	err = client.Call(&account, "personal_newAccount", pass)
	if err != nil {
		fmt.Println("failed to create acc", err)
		return "", err
	}
	fmt.Println("acc==", account)
	return account, nil
}

func UploadPic(from, pass string, tokenId int64) error {
	client, err := ethclient.Dial(config.ETHConnStr)
	if err != nil {
		fmt.Println("failed to connet to geth", err)
		return err
	}

	instance, err := NewStandingbook(common.HexToAddress(config.ContractAddr), client)
	if err != nil {
		fmt.Println("failed to interact with abi", err)
		return err
	}

	fileName, err := utils.GetFileName(string([]rune(from)[2:]), config.KeyStoreDie)
	if err != nil {
		fmt.Println("failed to find keystore file")
		return err
	}

	file, err := os.Open(config.KeyStoreDie + "/" + fileName)
	if err != nil {
		fmt.Println("failed to open file", err)
		return err
	}

	auth, err := bind.NewTransactor(file, pass)
	if err != nil {
		fmt.Println("failed to sign transactor", err)
		return err
	}

	// Mint to the minter...
	_, err = instance.UploadAndMint(auth, common.HexToAddress(from), big.NewInt(tokenId))
	if err != nil {
		fmt.Println("failed to Mint", err)
		return err
	}
	instance.IsOwner()

	return nil
}
