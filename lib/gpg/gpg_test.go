package gpg

import (
	"bytes"
	"encoding/base64"
	"io"
	"io/ioutil"
	"strings"
	"testing"
	"time"

	"github.com/juju/errors"
)

func TestSymmetricEncrypt_Decrypt(t *testing.T) {
	type TestCase struct {
		PlainText string
		Passphare []byte

		ExpectErr error
	}

	cases := []TestCase{
		{
			PlainText: "hello world",
			Passphare: []byte("abc"),
			ExpectErr: nil,
		},
	}

	for i, tc := range cases {
		ciphertext_buf := new(bytes.Buffer)
		plaintext, err := SymmetricEncrypt(ciphertext_buf, tc.Passphare)
		if errors.Cause(err) != tc.ExpectErr {
			t.Fatal("case", i,
				"expect", tc.ExpectErr,
				"got", errors.Cause(err),
				"errorstack:", errors.ErrorStack(err))
		}
		if tc.ExpectErr != nil {
			continue
		}

		go func() {
			defer plaintext.Close()
			io.WriteString(plaintext, tc.PlainText)
		}()

		// bytes.Buffer.Read will return err io.EOF if the buffer
		// has no data to return. So sleep some time to wait SymmetricEncrypt
		// products ciphertext into ciphertext_buf
		time.Sleep(time.Second / 20)
		plaintext_decrypted, err := SymmetricDecrypt(ciphertext_buf, tc.Passphare)
		if err != nil {
			t.Fatal("case", i, errors.ErrorStack(err))
		}

		data, err := ioutil.ReadAll(plaintext_decrypted)
		if err != nil {
			t.Fatal("case", i, errors.ErrorStack(err))
		}

		if string(data) != tc.PlainText {
			t.Fatal("expect", tc.PlainText, "got", string(data))
		}
	}
}

type AlwaysErrWriter struct {
	Err error
}

func (w *AlwaysErrWriter) Write(p []byte) (n int, err error) {
	return 0, w.Err
}

func TestSymmetricEncrypt(t *testing.T) {
	w := &AlwaysErrWriter{
		Err: errors.New("error for test"),
	}
	_, err := SymmetricEncrypt(w, []byte(nil))
	if err == nil {
		t.Fatal("expect err != nil")
	}
}

func TestSymmetricDecrypt(t *testing.T) {
	_, err := SymmetricDecrypt(strings.NewReader(""), nil)
	if err == nil {
		t.Fatal("expect err != nil")
	}

	_, err = SymmetricDecrypt(base64.NewDecoder(base64.StdEncoding, strings.NewReader("jA0ECQMCLqOWYuZvZvxg0kMBbyJqllazQptJt5SxKq98sQpmxAdJTFnralJWIh8hhjmo7WOSSDQP\nRshCt9uplY+dbU+j2D/NSEwtmu22zwVP0bM1")), []byte("23"))
	if errors.Cause(err) != ErrNotCurrectPassphrase {
		t.Fatal("expect ErrNotCurrectPassphrase got", errors.Cause(err), "errorstack:", errors.ErrorStack(err))
	}

	_, err = SymmetricDecrypt(base64.NewDecoder(base64.StdEncoding, strings.NewReader(`hQIOA4txhH5lDDKZEAgAhP+a8MX//9PLgt2OMLrPwbPcDzrBKKM52IAEP83wXJ3xUdnddk4HnNvBVGh/U1oVnir1vMBk7hte6KghlK92pEcjpgIWu0d06G01BchjxYDVWw/6WFlXgQeDTduEI2g1pKJvta6+grdaa9kSjf/TPJI5QV/HtdJqSOtDPl/APx/jpKrtW3Uf9wuR77/m8gjkUyYqgxH5hu3HUSAyeEdcvQYseZolq0mHUKaBF3IRv0rNDnXkU4xrQP5ZSiEc0Xs64b1+mNvTbvrSaZxBsov9c0B4bPthXQPBvei+6JZ6E9+a9aJTnkoGl/kLox0+VJd07n1VhfAlHEw/+Y2AYEcvNAf+KeoX9GputEWTmqCXiROkHgCeouIlhZsYPiRpZqol47ZAWi2xnuLT1sjN0EiIIIDFimYAXzsePt6iUEoQCheyGxk4DojGAXKiEQBVP47VvhfTWpvrcGOj8mXzQDOyKrxBJNZkD5acDnag6rqO5o8n/EsK2Bvrz9PHawFDcUcoXLGKiNM94Aw0hoxzCG6Lgqo6yvmWo8xh/tjpW0bwWM4VYoEO0mEDrk1RO4BgQMCCfOCrLf7EOferelhEaqhByQbKb4PZxpl9g3OZdD8ArIeDR2XBJqQs2D8xDf1UeaT/kEjtx8GN6Y4p6ImwVGPbF5+KcBAxAJ3yRQ2qIRPni5lAD9JJAWp7MLYGJqSaN7yyrV05J9PuMh4rdxkrq36ZNNwuQPaAfnskhcis41LEk3/zO0sQ2wEkHmErR9oChBvyikmnJ6U9wjKlqO/GpA==`)), []byte("23"))
	if err == nil {
		t.Fatal("expect err != nil")
	}
}
