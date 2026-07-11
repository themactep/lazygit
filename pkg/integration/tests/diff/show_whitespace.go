package diff

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ShowWhitespace = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Toggle visual whitespace indicators in the diff",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFile("myfile", "two  spaces\n")
		shell.GitAddAll()
		shell.Commit("initial commit")
		shell.UpdateFile("myfile", "two  spaces\n\twithtab\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Main().ContainsLines(
			Contains("two  spaces"),
		)

		t.Views().Files().
			IsFocused().
			Press(keys.Universal.ToggleShowWhitespace)

		t.ExpectToast(Equals("(showing whitespace)"))

		t.Views().Main().ContainsLines(
			Contains("·two··spaces"),
			Contains("+→--withtab"),
		)
	},
})
