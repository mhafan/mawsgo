package mawsgo

import (
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// ---------------------------------------------------------------------------
//
func LambdaTesting() bool {
	//
	return true
}

// ---------------------------------------------------------------------------
// Vystup z Lambdy
type MAWSLambdaResponse struct {
	// Text pripadneho erroru
	Message string
	// HTPP Status
	Status int
	// JSON navratove hodnoty dobrovolne
	Body string
}

// ---------------------------------------------------------------------------
//
func LocalExec(cmdString string) error {
	//
	cmd := exec.Command("sh", "-cex", strings.TrimSpace(cmdString))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// -----------------------------------------------------------------------
	// spusteni. navratova chyba
	if err_run := cmd.Run(); err_run != nil {
		//
		return err_run
	}

	return nil
}

// ---------------------------------------------------------------------------
//
func (m *MAWSLambdaResponse) JSON() string {
	//
	v, _ := json.Marshal(m)

	//
	return string(v)
}

// ---------------------------------------------------------------------------
// Invokace Lambdy (na strane AWS). Velmi opatrne ;)
// Pozn.: V aplikaci by jedna Lambda funkce snad radeji nemela volat jinou ;)
func (maws *MAWS) MAWSCallLambda(funName string, async bool, args interface{}) ([]byte, error) {
	// Handle na volani AWS/Lambda
	/*
		client := lambda.New(maws.AWS)

		// zakodovani vstupu
		_enco, _encoError := json.Marshal(args)

		// ...
		if _encoError != nil {
			//
			return []byte{}, _encoError
		}

		//
		_args := &lambda.InvokeInput{
			FunctionName: aws.String(funName),
			Payload:      _enco,
		}

		//
		if async {
			//
			_args.InvocationType = aws.String("Event")
		}

		// invokace: jmeno funkce + payload
		_resp, _err := client.Invoke(_args)

		//
		return _resp.Payload, _err*/
	return []byte{}, errors.New("not impl")
}

// ---------------------------------------------------------------------------
// Invokace se vstupnim a vystupnim kodovanim
func (maws *MAWS) MAWSCallLambdaResponse(funName string, args, resp interface{}) error {
	// volani s parametry
	_calling, _cerr := maws.MAWSCallLambda(funName, false, args)

	// ...
	if _cerr != nil {
		//
		return _cerr
	}

	// AWS z nejakeho duvodu transformuje na ESC sekvenci -> \"
	// davam ty \ pryc
	_cunqo, _ := strconv.Unquote(string(_calling))

	// rozbaleni JSON odpovedi
	return json.Unmarshal([]byte(_cunqo), resp)
}

// ---------------------------------------------------------------------------
// Invokace se vstupnim a vystupnim kodovanim
func (maws *MAWS) MAWSCallLambdaStandard(funName string, args interface{}) (*MAWSLambdaResponse, error) {
	//
	var _out MAWSLambdaResponse

	//
	_err := maws.MAWSCallLambdaResponse(funName, args, &_out)

	//
	return &_out, _err
}
