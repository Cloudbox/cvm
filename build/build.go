package build

import "fmt"

var (
	Version   string
	Timestamp string
	GitCommit string

	UserAgent string
)

func init() {
	UserAgent = fmt.Sprintf("cvm/%s", Version)
}
