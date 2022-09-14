package sugar

import "flag"

var (
	listen  *string = flag.String("listen", ":8080", "sugar listen")
	gListen *string = flag.String("glisten", ":8081", "sugar regist listen")
)
