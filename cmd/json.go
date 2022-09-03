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
		replace, pretty bool
	)
	cmd := &cobra.Command{
		Use:   `json`,
		Short: `Convert jsonnet to json.`,
		Run: func(cmd *cobra.Command, args []string) {
			c, e := converter.New()
			if e != nil {
				log.Fatalln(e)
			}
			if pretty {
				c.Convert(converter.PrettyJsonMarshaler{}, test, move, copy, replace)
			} else {
				c.Convert(converter.JsonMarshaler{}, test, move, copy, replace)
			}
		},
	}
	flags := cmd.Flags()
	flags.BoolVarP(&test, `test`, `t`, false, `test and println json to stdout`)
	flags.BoolVarP(&move, `move`, `m`, false, `move json to target`)
	flags.BoolVarP(&copy, `copy`, `c`, false, `copy json to target`)
	flags.BoolVarP(&replace, `replace`, `r`, false, `replacement target does not need to compare whether it has changed`)
	flags.BoolVarP(&pretty, `pretty`, `p`, false, `output pretty json`)
	rootCmd.AddCommand(cmd)
}
