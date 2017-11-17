package builder

import (
	"fmt"
	"os/exec"
	"os"
)

func BuildLatex(rootDir, rootFileName string) error {
	fmt.Println("Build latex")
	var filePath = rootDir + "/" + rootFileName
	var cmd = exec.Command(
		"bash",
		"-c",
		fmt.Sprintf("cd %s; pdflatex -synctex=1 -interaction=nonstopmode %s", rootDir, filePath),
	)
	cmd.Stdout = os.Stdout
	return cmd.Run()
}
