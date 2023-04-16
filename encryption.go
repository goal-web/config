package config

import (
	"github.com/goal-web/contracts"
	"github.com/goal-web/encryption"
	"github.com/goal-web/supports/commands"
	"github.com/goal-web/supports/logs"
	"os"
)

func EncryptionCommand(app contracts.Application) contracts.Command {
	return &encryptionCommand{
		Command: commands.Base("env {action:encrypt or decrypt} --{driver?} --{env?} --{key?}", "env"),
	}
}

type encryptionCommand struct {
	commands.Command
}

func (cmd encryptionCommand) Handle() (result any) {
	key := cmd.StringOptional("key", os.Getenv("GOAL_ENV_ENCRYPTION_KEY"))

	encryptor := encryption.DefaultDrivers()[cmd.StringOptional("driver", "AES")](key)

	if cmd.GetString("action") == "encrypt" {

		envPath := cmd.StringOptional("env", ".env")

		envBytes, err := os.ReadFile(envPath)

		if err != nil {
			panic(err)
		}

		if len(envBytes) == 0 {
			logs.Default().Warn("env file is empty.")
			return
		}

		err = os.WriteFile(".env.encrypted", encryptor.Encrypt(envBytes), os.ModePerm)

		if err != nil {
			panic(err)
		}

		logs.Default().Info("The env encrypted file has been written to .env.encrypted")

	} else {
		envBytes, err := os.ReadFile(".env.encrypted")

		if err != nil {
			logs.Default().Error(err.Error())
			return
		}

		if len(envBytes) == 0 {
			logs.Default().Warn("env file is empty.")
			return
		}

		envBytes, err = encryptor.Decrypt(envBytes)
		if err != nil {
			panic(err)
		}

		err = os.WriteFile(cmd.StringOptional("env", ".env"), envBytes, os.ModePerm)

		if err != nil {
			panic(err)
		}

		logs.Default().Info("env file is decrypted.")
	}

	return
}
