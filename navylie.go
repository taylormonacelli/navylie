package navylie

import (
	"fmt"
	"log/slog"
	"os/exec"
	"path/filepath"

	"github.com/taylormonacelli/ashpalm"
	"github.com/taylormonacelli/coalfoot"
)

func Main(userProjectDir string) int {
	userProjectDirAbs, err := filepath.Abs(userProjectDir)
	if err != nil {
		slog.Error("filepath.abs", "path", userProjectDirAbs, "error", err.Error())
	}
	slog.Debug("user project dir", "dir", userProjectDirAbs)

	templateData := TemplateData{
		ModuleName:          filepath.Base(userProjectDirAbs),
		GithubUsername:      "taylormonacelli",
		BazelGazelleVersion: "0.33.0",
		BazelRulesGo:        "0.42.0",
	}

	tpl := coalfoot.NewTxtarTemplate()
	tpl.FetchFromRemoteIfOld()

	renderTemplate(tpl, userProjectDir, templateData)

	slog.Debug("extracting txtar", "path", tpl.LocalPathRendered, "target dir", userProjectDirAbs)
	err = tpl.Extract(userProjectDirAbs)
	if err != nil {
		slog.Error("extracting", "error", err.Error())
		return 1
	}

	slog.Debug("running func", "func", "runGoModTidy")
	err = runGoModTidy(userProjectDirAbs)
	if err != nil {
		return 1
	}

	return 0
}

func runGoModTidy(projectDir string) error {
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = projectDir

	code, outStr, errStr := ashpalm.RunCmd(cmd)
	if code != 0 {
		slog.Error("runcmd", "cmd", cmd.String(), "stdout", outStr, "stderr", errStr, "code", code)
		return fmt.Errorf("%s had unexpected exit code %d with error %s", cmd.String(), code, errStr)
	}

	slog.Debug("runcmd", "cmd", cmd.String(), "stdout", outStr, "stderr", errStr, "code", code)

	return nil
}
