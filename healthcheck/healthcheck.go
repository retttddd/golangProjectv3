package healthcheck

import (
	"fmt"
	"net/http"
)

func Check(port string) error {
	_, err := http.Get(fmt.Sprintf("http://127.0.0.1:%s/health", port))
	return err
}
