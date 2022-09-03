package cmd

import (
	"log"

	"github.com/powerpuffpenguin/ejt/internal/converter"
	"github.com/spf13/cobra"
)

func init() {
	var (
		test,
		move, copy,
		replace bool
	)
	cmd := &cobra.Command{
		Use:   `yaml`,
		Short: `Convert jsonnet to yaml.`,
		Run: func(cmd *cobra.Command, args []string) {
			c, e := converter.New()
			if e != nil {
				log.Fatalln(e)
			}
			c.Convert(converter.YamlMarshaler{}, test, move, copy, replace)
		},
	}
	flags := cmd.Flags()
	flags.BoolVarP(&test, `test`, `t`, false, `test and println yaml to stdout`)
	flags.BoolVarP(&move, `move`, `m`, false, `move yaml to target`)
	flags.BoolVarP(&copy, `copy`, `c`, false, `copy yaml to target`)
	flags.BoolVarP(&replace, `replace`, `r`, false, `replacement target does not need to compare whether it has changed`)
	rootCmd.AddCommand(cmd)
}
