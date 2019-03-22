package cmd

import (
	"fmt"
	"os"

	"github.com/ncr-devops-platform/nagiosfoundation/lib/app/nagiosfoundation"
	"github.com/spf13/cobra"
)

// Execute runs the root command
func Execute() {
	var greaterThan bool
	var warning, critical float64
	var pollingAttempts, pollingDelay int
	var metricName, counterName string

	var rootCmd = &cobra.Command{
		Use:   "check_performance_counter",
		Short: "Retrieve and compare values on a performance counter.",
		Long: `The performance counter check is Windows only. It retrieves a Windows Performance Counter
(--counter_name) and compares it to --critical and --warning then outputs an appropriate response
based on the check. Many flags make this check quite configurable.

The defaults for this check have the --critical and --warning values set to 0, and the counter value
retrieved is compared to be lesser than those values. Generally a counter value will be > 0, causing
this check to emit an OK response when using these defaults.`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.ParseFlags(os.Args)

			msg, retval := nagiosfoundation.CheckPerformanceCounter(warning, critical, greaterThan, pollingAttempts,
				pollingDelay, metricName, counterName)

			fmt.Println(msg)
			os.Exit(retval)
		},
	}

	nagiosfoundation.AddVersionCommand(rootCmd)

	const counterNameFlag = "counter_name"
	rootCmd.Flags().StringVarP(&counterName, counterNameFlag, "n", "", "the name of the performance counter to check")
	rootCmd.MarkFlagRequired(counterNameFlag)
	rootCmd.Flags().StringVarP(&metricName, "metric_name", "m", "", "the name of the metric generated by this check")
	rootCmd.Flags().Float64VarP(&warning, "warning", "w", 0, "the threshold to issue a warning alert")
	rootCmd.Flags().Float64VarP(&critical, "critical", "c", 0, "the threshold to issue a critical alert")
	rootCmd.Flags().BoolVarP(&greaterThan, "greater_than", "g", false, "issue warnings if the metric is greater than the expected thresholds (default false)")
	rootCmd.Flags().IntVarP(&pollingAttempts, "polling_attempts", "a", 2, "the number of times to fetch and average the performance counter")
	rootCmd.Flags().IntVarP(&pollingDelay, "polling_delay", "d", 1, "the number of seconds to delay between polling attempts")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
