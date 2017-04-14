package cmd

import (
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"time"

	"github.com/immofon/databoard/lib/gpg"
	"github.com/juju/errors"
	"github.com/mholt/archiver"
	"github.com/spf13/cobra"
)

// packetCmd represents the packet command
var packetCmd = &cobra.Command{Use: "packet"}

func init() {
	RootCmd.AddCommand(packetCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// packetCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// packetCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// packet <file> [file [files...]]
	packetCmd.Run = func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			packetCmd.Usage()
			exit(1)
		}

		passphare, err := getPassphare()
		if err != nil {
			packetCmd.Usage()
			exit(2, "require passphare\n", errors.ErrorStack(errors.Trace(err)))
		}

		basedir := fmt.Sprintf("/tmp/databoard-%d", os.Getpid())
		defer os.RemoveAll(basedir)
		data_tar_gz_gpg, err := Packet(basedir, "", passphare, args)
		if err != nil {
			exit(3, errors.ErrorStack(errors.Trace(err)))
		}

		os.Link(data_tar_gz_gpg, "data.tar.gz.gpg")
	}
}

func Packet(basedir, fileprefix string, passphare []byte, files []string) (path_ string, err error) {
	for _, filename := range files {
		_, err := os.Stat(filename)
		if err != nil {
			return "", errors.Annotatef(err, "file %q is not exist", filename)
		}
	}

	os.Mkdir(basedir, 0700)

	fileprefix = strings.TrimSpace(fileprefix)
	if fileprefix == "" {
		fileprefix = fmt.Sprintf("data-%d", time.Now().UnixNano())
	}

	data_tar_gz := path.Join(basedir, fileprefix+".tar.gz")
	err = archiver.TarGz.Make(data_tar_gz, files)
	if err != nil {
		return "", errors.Trace(err)
	}

	data_tar_gz_gpg, err := GPGEncrypt(data_tar_gz, passphare)
	if err != nil {
		return "", errors.Trace(err)
	}

	return data_tar_gz_gpg, nil
}

func GPGEncrypt(file string, passphare []byte) (name string, err error) {
	_, err = os.Stat(file)
	if err != nil {
		return "", errors.Annotatef(err, "file %q is not exist", file)
	}

	in, err := os.Open(file)
	if err != nil {
		return "", errors.Trace(err)
	}
	defer in.Close()

	name = file + ".gpg"
	out, err := os.OpenFile(name, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return "", errors.Trace(err)
	}
	defer out.Close()

	plaintext, err := gpg.SymmetricEncrypt(out, passphare)
	if err != nil {
		return "", errors.Trace(err)
	}

	_, err = io.Copy(plaintext, in)
	plaintext.Close()
	return name, errors.Trace(err)
}
