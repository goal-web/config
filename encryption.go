package config

import (
	"fmt"
	"github.com/goal-web/contracts"
	"github.com/goal-web/supports/commands"
	"github.com/goal-web/supports/logs"
	"github.com/goal-web/supports/utils"
	"os"
)

func EncryptionCommand() (contracts.Command, contracts.CommandHandlerProvider) {
	return commands.Base("env {action:encrypt or decrypt} --{driver?} --{in?} --{out} --{key?} --{force?}", "env"),
		func(app contracts.Application) contracts.CommandHandler {
			return &encryptionCommand{app: app}
		}
}

type encryptionCommand struct {
	commands.Command
	app contracts.Application
}

func (cmd encryptionCommand) Handle() (result any) {
	key := cmd.StringOptional("key", os.Getenv("GOAL_ENV_ENCRYPTION_KEY"))

	if key == "" {
		key = utils.RandStr(32)
		logs.Default().WithField("key", key).Info("key generated.")
	}

	encryptDriver := cmd.app.Get("encryption").(contracts.EncryptManager).
		Driver(cmd.StringOptional("driver", "AES"))

	encryptor := encryptDriver(key)

	if cmd.GetString("action") == "encrypt" {
		input := cmd.StringOptional("in", ".env")
		output := cmd.StringOptional("out", ".env.encrypted")

		envBytes, err := os.ReadFile(input)

		if err != nil {
			logs.Default().WithError(err).Error("failed to read file.")
			return
		}

		if len(envBytes) == 0 {
			logs.Default().Warn("env file is empty.")
			return
		}

		if err = os.WriteFile(output, encryptor.Encrypt(envBytes), os.ModePerm); err != nil {
			logs.Default().WithError(err).Error("failed to write encrypted env.")
		} else {
			logs.Default().Info("The env encrypted file has been written to " + output)
		}
	} else {
		input := cmd.StringOptional("in", ".env.encrypted")
		output := cmd.StringOptional("out", ".env")
		envBytes, err := os.ReadFile(input)

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
			logs.Default().WithError(err).Error("failed to encrypt.")
			return
		}

		logs.Default().Info("env file is decrypted.")

		if !utils.ExistsPath(output) || cmd.GetBool("force") {
			if err = os.WriteFile(output, envBytes, os.ModePerm); err != nil {
				logs.Default().WithError(err).Error("failed to encrypt.")
			} else {
				logs.Default().Info("The env file has been written.")
			}
		} else {
			logs.Default().Warn(fmt.Sprintf("The file [%s] is exists.", output))
		}
	}

	return
}
