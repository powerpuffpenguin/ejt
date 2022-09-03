package cmd

import (
	"log"
	"os"

	"github.com/powerpuffpenguin/ejt/version"
	"github.com/spf13/cobra"
)

func init() {
	cmd := &cobra.Command{
		Use:   `init`,
		Short: `Initializes a ejt project and creates a ejt.jsonnet file.`,
		Run: func(cmd *cobra.Command, args []string) {
			f, e := os.OpenFile(`ejt.jsonnet`, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
			if e != nil {
				log.Fatalln(e)
			}
			defer f.Close()
			_, e = f.WriteString(`{
  version: '` + version.Version + `',
  endpoints: [
    {
      output: './dst',  // redirect output structure to the directory.
      target: './envoy',  // target root directory.
      source: './src',  // source root directory.
      resources: [
        'envoy.jsonnet',
      ],
    },
  ],
}`)
			if e != nil {
				log.Fatalln(e)
			}
			log.Println(`Successfully created a ejt.jsonnet file.`)

		},
	}
	rootCmd.AddCommand(cmd)
}
