package main

import (
	"fmt"

	"github.com/achwanyusuf/carrent-lib/pkg/hash"
)

func main() {
	// str := []string{
	// 	"JzHNlyNIaqpQEBHuUwaWNyUkfsmSVPha",
	// 	"CZaVUuMYKMHOzfgNFUTZByHyoIxQSHIz",
	// 	"mMuVPnFcLzdZCyuuPfuDoEzhvzNAwYqR",
	// 	"oIkhcMLVnMoJGpNvMuYOuRmwfMjTIFfp",
	// 	"dFMMZMpITtxUnsDdzWLirxrkvXYQdpWm",
	// 	"DGdfYiFnJXIkylBZITvcnihCsOumTLqY",
	// 	"VwLmsbXEeyszGXiijSgUFEYNxjymCLUH",
	// 	"LOFIqJWJOODYcsDuzPASkuJHzetbjKux",
	// 	"jGZWYzfsLEdMllxOvNoLppPAVLGVNkHV",
	// 	"KBNnzGWmORKNQRYAsRvmpTealBwGKiDO",
	// }

	// for _, a := range str {
	// 	key := "62157hasjhjas"
	// 	k, _ := hash.EncAES(a, key)
	// 	fmt.Println(k)
	// }

	d, _ := hash.Hash("12345678")
	fmt.Println(d)
}
