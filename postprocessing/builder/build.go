package builder

import (
	"fmt"
	"os"
	"os/exec"
)

func BuildLatex(rootDir, rootFileName string) error {
	fmt.Println("Build latex")

	cleanupCmd := getCmd(getCleanupStr())
	if err := cleanupCmd.Run(); err != nil {
		return err
	}

	runCmd := getCmd(getBuildStr(rootDir, rootFileName))
	runCmd.Run()
	return nil
}

func getCmd(text string) *exec.Cmd {
	cmd := exec.Command("bash", "-c", text)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	return cmd
}

func getCleanupStr() string {
	return "docker kill latex_container || true && docker rm latex_container || true"
}

func getBuildStr(buildDir, rootFileName string) string {
	return fmt.Sprintf(
		"docker run -it -d --name latex_container -v %s:/home -w /home sumdoc/texlive-2017 /bin/bash && \\"+
			"docker exec $(docker ps -aqf 'name=latex_container') pdflatex -shell-escape -syntex=1 -interaction=nonstopmode /home/%s",
		buildDir, rootFileName,
	)
}
