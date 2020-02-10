//go:generate sh -c "rm -rf vendor && mkdir -p vendor && git clone https://github.com/pojntfx/ipxe.git vendor/ipxe"

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	os.Chdir(filepath.Join("vendor", "ipxe", "src"))

	out, err := exec.Command("make").CombinedOutput()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(out))
}
