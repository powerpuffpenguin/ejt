package cmd

import (
	"log"

	"github.com/powerpuffpenguin/ejt/internal/converter"
	"github.com/spf13/cobra"
)

func init() {
	var (
		test bool
		move bool
		copy bool
	)
	cmd := &cobra.Command{
		Use:   `yaml`,
		Short: `Convert jsonnet to yaml.`,
		Run: func(cmd *cobra.Command, args []string) {
			c, e := converter.New()
			if e != nil {
				log.Fatalln(e)
			}
			c.Yaml(test, move, copy)
		},
	}
	flags := cmd.Flags()
	flags.BoolVarP(&test, `test`, `t`, false, `test and println yaml to stdout`)
	flags.BoolVarP(&move, `move`, `m`, false, `move yaml to target`)
	flags.BoolVarP(&copy, `copy`, `c`, false, `copy yaml to target`)
	rootCmd.AddCommand(cmd)
}
