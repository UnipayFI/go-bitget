package request

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/UnipayFI/go-bitget/client"
	"github.com/UnipayFI/go-bitget/common"
	"github.com/go-json-experiment/json/jsontext"
)

// apiResponse is Bitget's uniform REST envelope. "code" is "00000" on success;
// "data" carries the endpoint-specific payload.
type apiResponse[T any] struct {
	Code        string `json:"code"`
	Msg         string `json:"msg"`
	RequestTime int64  `json:"requestTime"`
	Data        T      `json:"data"`
}

// Do executes the request and decodes the envelope's data field into *T. A
// non-"00000" code is returned as a *client.APIError.
func Do[T any](r *Request) (resp *T, err error) {
	if r.err != nil {
		return nil, r.err
	}
	if err = r.prepare(); err != nil {
		return nil, err
	}

	r.client.GetLogger().Debugf("request: %s %s body=%s", r.method, r.r.URL, r.bodyJSON)
	defer func() {
		if err != nil {
			r.client.GetLogger().Errorf("request %s %s failed: %s", r.method, r.r.URL, err)
		}
	}()

	response, err := r.r.Send()
	if err != nil {
		return nil, err
	}
	body := response.Body()
	r.client.GetLogger().Debugf("response: %s", common.BytesToString(body))

	var out apiResponse[T]
	if uerr := r.client.GetHttpClient().JSONUnmarshal(body, &out); uerr != nil {
		// The body was not a well-formed envelope (gateway error, HTML, ...).
		if apiErr := parseAPIError(r, body); apiErr != nil {
			return nil, apiErr
		}
		return nil, fmt.Errorf("request failed (status %d): %s", response.StatusCode(), common.BytesToString(body))
	}
	if out.Code != "00000" {
		return nil, &client.APIError{Code: out.Code, Message: out.Msg, RequestTime: out.RequestTime}
	}
	return &out.Data, nil
}

// DoRawData executes the request and returns the raw JSON bytes of the
// envelope's "data" field (after verifying the success code). Tests use it to
// diff the real response shape against the typed structs.
func DoRawData(r *Request) ([]byte, error) {
	if r.err != nil {
		return nil, r.err
	}
	if err := r.prepare(); err != nil {
		return nil, err
	}
	response, err := r.r.Send()
	if err != nil {
		return nil, err
	}
	body := response.Body()
	var env struct {
		Code        string         `json:"code"`
		Msg         string         `json:"msg"`
		RequestTime int64          `json:"requestTime"`
		Data        jsontext.Value `json:"data"`
	}
	if uerr := r.client.GetHttpClient().JSONUnmarshal(body, &env); uerr != nil {
		return nil, fmt.Errorf("request failed (status %d): %s", response.StatusCode(), common.BytesToString(body))
	}
	if env.Code != "00000" {
		return nil, &client.APIError{Code: env.Code, Message: env.Msg, RequestTime: env.RequestTime}
	}
	return env.Data, nil
}

// DoRaw executes the request and returns the raw, undecoded response body. Use
// it for the rare endpoints whose payload shape is non-uniform.
func DoRaw(r *Request) ([]byte, error) {
	if r.err != nil {
		return nil, r.err
	}
	if err := r.prepare(); err != nil {
		return nil, err
	}
	r.client.GetLogger().Debugf("request: %s %s body=%s", r.method, r.r.URL, r.bodyJSON)
	response, err := r.r.Send()
	if err != nil {
		return nil, err
	}
	body := response.Body()
	if apiErr := parseAPIError(r, body); apiErr != nil {
		return nil, apiErr
	}
	return body, nil
}

// prepare finalizes the URL, body and (when private) the ACCESS-* signing
// headers. The signed prehash is timestamp + method + requestPath + body,
// using the exact bytes that go on the wire.
func (r *Request) prepare() error {
	r.r.URL = r.fullURL()
	r.r.Method = r.method
	if r.bodyJSON != "" {
		r.r.SetHeader("Content-Type", "application/json")
		r.r.SetBody(r.bodyJSON)
	}
	if r.client.IsDemoTrading() {
		r.r.SetHeader("paptrading", "1")
	}
	if !r.needSign {
		return nil
	}

	apiKey := r.client.GetAPIKey()
	secret := r.client.GetAPISecret()
	passphrase := r.client.GetPassphrase()
	if apiKey == "" || secret == "" || passphrase == "" {
		return errors.New("missing credentials: configure client.WithAuth(apiKey, apiSecret, passphrase)")
	}

	ts := strconv.FormatInt(r.client.TimestampMs(), 10)
	prehash := ts + r.method + r.requestPath() + r.bodyJSON

	var (
		sign string
		err  error
	)
	if fn := r.client.GetSignFn(); fn != nil {
		sign, err = fn(secret, prehash)
	} else {
		sign, err = HMACSign(secret, prehash)
	}
	if err != nil {
		return err
	}

	r.r.SetHeader("Content-Type", "application/json")
	r.r.SetHeader("ACCESS-KEY", apiKey)
	r.r.SetHeader("ACCESS-SIGN", sign)
	r.r.SetHeader("ACCESS-TIMESTAMP", ts)
	r.r.SetHeader("ACCESS-PASSPHRASE", passphrase)
	return nil
}

// parseAPIError tries to decode body as a Bitget error envelope, returning nil
// when it is not an API-level error.
func parseAPIError(r *Request, body []byte) error {
	apiErr := &client.APIError{}
	if e := r.client.GetHttpClient().JSONUnmarshal(body, apiErr); e != nil {
		return nil
	}
	if !apiErr.IsValid() {
		return nil
	}
	return apiErr
}
