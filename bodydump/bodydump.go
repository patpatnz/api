package bodydump

type bodyDump struct {
	Time            time.Time    `json:"time"`
	RequestID       string       `json:"id"`
	RequestBody     string       `json:"request_body"`
	ResponseBody    string       `json:"response_body"`
	RequestHeaders  *http.Header `json:"response_headers"`
	ResponseHeaders *http.Header `json:"request_headers"`
}

func BodyDump(c echo.Context, reqBody, resBody []byte) {
	status := c.Response().Status
	if viper.GetBool("debug") || ((status < 199 || status > 299) && status != 404) {
		resHeaders := c.Response().Header()
		v := bodyDump{
			Time:            time.Now(),
			RequestID:       c.Response().Header().Get(echo.HeaderXRequestID),
			RequestBody:     string(reqBody),
			ResponseBody:    string(resBody),
			RequestHeaders:  &c.Request().Header,
			ResponseHeaders: &resHeaders,
		}
		buf, _ := json.Marshal(v)
		fmt.Println(string(buf))
	}
}

