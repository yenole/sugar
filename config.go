package sugar

import "flag"

var Listen *string = flag.String("listen", ":8080", "sugar listen")

func init() {
	flag.Parsed()
}
