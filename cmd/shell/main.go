package main

import (
	"fmt"
	gp "github.com/c-bata/go-prompt"
	"github.com/soracom/soracom-shell"
)

func main() {
	executor := shell.NewSoracomExecutor("/bin/sh")
	completer := shell.NewSoracomCompleter("/soracom-api.en.yaml")

	fmt.Print(`
              ..;;ttLLSSSSSSLLtt;;..
          ..11SSSSSSSSSSSSSSSSSSSSSS11..
        ::LLSSSSSSttii::,,::iittSSSSSSLL::
      ::SSSSSS11..              ..11SSSSSS::
    ::SSSSSSSSttii::..              ::LLSSSS::
  ..LLSSSSSSSSSSSSSSSSffii::..        ,,LLSSLL..
  11SSSS::,,;;ttLLSSSSSSSSSSSSff11::..  ::SSSS11
..SSSS11          ,,;;11LLSSSSSSSSSSSS..  11SSSS..
iiSSSS,,                  ..::11LLSSSS..  ,,SSSSii
ttSSff                          ;;SSSS..    ffSSff
LLSSii                          ;;SSSS..    iiSSLL
SSSS;;                        ,,11SSSS..    ;;SSSS
SSSS::                ,,iittLLSSSSSSSS..    ::SSSS
SSSS;;      ..::iittSSSSSSSSSSSSSSffii      ;;SSSS
LLSSii    ;;SSSSSSSSSSSSLLttii,,            iiSSLL
ttSSff    ..LLSSSStt;;,,          ::        ffSSff
iiSSSS,,    iiSSSS,,          ,,::tt,,..  ,,SSSSii
..SSSS11    ..LLSStt          ;;LLSStt..  11SSSS..
  11SSSS::    iiSSSS,,          LLff;;  ::SSSS11
  ..LLSSLL,,  ..LLSStt  ..tt11..,,  ::,,LLSSLL..
    ::SSSSLL::  iiSSSS::ffSSSS;;    ::LLSSSS::
      ::SSSSSS11,,LLSSSSSSSS11  ..11SSSSSS::
        ,,LLSSSSSSLLSSSSSSffiittSSSSSSLL::
          ..11LLSSSSSSSSSSSSSSSSSSLL11..
              ..;;ttLLSSSSSSLLtt;;..

                      Type exit or Ctrl-D to exit.
`)
	gp.New(
		executor.Execute,
		completer.Complete,
		gp.OptionTitle("SORACOM Shell"),
		gp.OptionPrefix("SORACOM> "),
		gp.OptionMaxSuggestion(10),
		gp.OptionPrefixTextColor(gp.Cyan),
		gp.OptionPreviewSuggestionTextColor(gp.Blue),
		gp.OptionSelectedSuggestionBGColor(gp.LightGray),
		gp.OptionSuggestionBGColor(gp.DarkGray),
	).Run()
}
