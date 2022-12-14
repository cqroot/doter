package cmd

import (
	"os"
	"path"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"

	"github.com/cqroot/dotm/pkg/dotmanager"
)

var (
	Tag     string
	Verbose bool
	rootCmd = &cobra.Command{
		Use:   "dm",
		Short: "DotM - Manage dotfiles more easily.",
		Long:  "DotM - Manage dotfiles more easily.",
		Run:   runRootCmd,
	}
)

func init() {
	rootCmd.PersistentFlags().StringVarP(
		&Tag, "tag", "t", "", "use dotfiles with specified tags")
	rootCmd.PersistentFlags().BoolVarP(
		&Verbose, "verbose", "v", false, "")
}

func Execute() {
	err := rootCmd.Execute()
	cobra.CheckErr(err)
}

func runRootCmd(cmd *cobra.Command, args []string) {
	baseDir, err := getBaseDir()
	cobra.CheckErr(err)

	dm, err := dotmanager.New(baseDir, path.Join(baseDir, "dotm.toml"), Tag)
	cobra.CheckErr(err)

	t := newTable()
	if Verbose {
		t.AppendHeader(table.Row{
			"#", "Name", "Type", "Source Path", "Target Path", "Status",
		})
	} else {
		t.AppendHeader(table.Row{
			"#", "Name", "Type", "Status",
		})
	}

	for idx, dot := range dm.Dots {
		err := dot.Check()
		msg := ""

		switch err {
		case nil:
			msg = text.FgGreen.Sprint("OK")
		default:
			msg = text.FgRed.Sprintf("ERROR: %s", err.Error())
		}

		if Verbose {
			t.AppendRow([]interface{}{
				idx, dot.Name(), dot.Type(), dot.Source(), dot.Target(), msg,
			})
		} else {
			t.AppendRow([]interface{}{
				idx, dot.Name(), dot.Type(), msg,
			})
		}
	}

	t.Render()
}

func newTable() table.Writer {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleRounded)
	t.Style().Format.Header = text.FormatDefault

	return t
}

func getBaseDir() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	return path.Join(configDir, "dotfiles"), nil
}
