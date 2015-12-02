// result is a helper package to give consistent result formats
package result

import (
	"github.com/labstack/echo"
)

type Type int

const (
	JSON Type = iota
	String
	XML
)

type APIResult struct {
	Code        int         `json:"code"`
	TxnID       string      `json:"transaction_id"`
	ErrorReason string      `json:"error_reason,omitempty"`
	Result      interface{} `json:"result,omitempty"`
}

func Error(ctx *echo.Context, code int, err error) error {
	txn := ctx.Get("txn").(string)
	ctx.JSON(code, &APIResult{Code: code, TxnID: txn, ErrorReason: err.Error()})
	return err
}

func Send(ctx *echo.Context, code int, res interface{}) error {
	txn := ctx.Get("txn").(string)
	return ctx.JSON(code, &APIResult{Code: code, TxnID: txn, Result: res})
}
