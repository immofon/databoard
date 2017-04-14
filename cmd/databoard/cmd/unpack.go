package cmd

import (
	"io"
	"os"
	"strings"

	"github.com/immofon/databoard/lib/gpg"
	"github.com/juju/errors"
	"github.com/mholt/archiver"
	"github.com/spf13/cobra"
)

// unpackCmd represents the unpack command
var unpackCmd = &cobra.Command{
	Use: "unpack",
}

func init() {
	RootCmd.AddCommand(unpackCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// unpackCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// unpackCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	var (
		src string // file name
		dst string // dir  name
	)

	unpackCmd.Flags().StringVarP(&src, "src", "s", "", "src file(must end with .tar.gz.gpg)")
	unpackCmd.Flags().StringVarP(&dst, "dst", "d", "", "dst dir")

	unpackCmd.Run = func(cmd *cobra.Command, args []string) {
		passphare, err := getPassphare()
		if err != nil {
			unpackCmd.Usage()
			exit(1, "require passphare\n", errors.ErrorStack(errors.Trace(err)))
		}
		err := Unpack(src, passphare, dst)
		if err != nil {
			exit(2, errors.ErrorStack(errors.Annotatef(err, "Unpack %q to %q", src, dst)))
		}
	}
}

func Unpack(srcfile string, passphare []byte, dstDir string) error {
	if !strings.HasSuffix(srcfile, ".tar.gz.gpg") {
		return errors.Errorf("srcfile hasn't suffix .tar.gz.gpg <srcfile: %q>", srcfile)
	}
	plaintfile, err := GPGDecrypt(srcfile, passphare)
	if err != nil {
		return errors.Annotate(err, "GPG symmetric decrypt")
	}
	defer os.Remove(plaintfile)

	err = archiver.TarGz.Open(plaintfile, dstDir)
	if err != nil {
		return errors.Annotate(err, "archiver.TarGz.Open")
	}

	return nil
}

func GPGDecrypt(cipherfile string, passphare []byte) (plainfile string, err error) {
	if !strings.HasSuffix(cipherfile, ".gpg") {
		return "", errors.Errorf("expect filename %q has suffix .gpg", cipherfile)
	}

	in, err := os.Open(cipherfile)
	if err != nil {
		return "", errors.Trace(err)
	}

	plainfile = cipherfile[:len(cipherfile)-len(".gpg")]
	out, err := os.OpenFile(plainfile, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return "", errors.Trace(err)
	}

	plaintext, err := gpg.SymmetricDecrypt(in, passphare)
	if err != nil {
		return "", errors.Trace(err)
	}

	_, err = io.Copy(out, plaintext)
	return plainfile, err
}
