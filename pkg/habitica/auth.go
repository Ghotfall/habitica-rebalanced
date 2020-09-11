package habitica

import (
	"github.com/valyala/fasthttp"
	"os"
)

func SetHeaders(req *fasthttp.Request) {
	req.Header.Set("x-api-user", os.Getenv("HABITICA_API_USER"))
	req.Header.Set("x-api-key", os.Getenv("HABITICA_API_KEY"))
	req.Header.Set("x-client", os.Getenv("HABITICA_CLIENT"))
}
