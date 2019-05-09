package cmds

import (
	"bufio"
	"crypto/rsa"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gomodules.xyz/cert"
)

func NewCmdGetPublicKey() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "pub-key",
		Short:             "Prints public key from PEM encoded RSA private key",
		DisableAutoGenTag: true,
		Run: func(cmd *cobra.Command, args []string) {
			reader := bufio.NewReader(os.Stdin)
			keyBytes, err := ioutil.ReadAll(reader)
			if err != nil {
				Fatal(errors.Wrap(err, "failed to read private key"))
			}
			key, err := cert.ParsePrivateKeyPEM(keyBytes)
			if err != nil {
				Fatal(errors.Wrap(err, "failed to parse private key"))
			}
			saKey, ok := key.(*rsa.PrivateKey)
			if !ok {
				Fatal(errors.Wrapf(err, "only supports rsa private key. Found %v", reflect.ValueOf(key).Kind()))
			}
			saPub, err := cert.EncodePublicKeyPEM(&saKey.PublicKey)
			if err != nil {
				Fatal(errors.Wrap(err, "failed to generate public key"))
			}
			fmt.Println(string(saPub))
			os.Exit(0)
		},
	}
	return cmd
}
