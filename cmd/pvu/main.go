package main

import (
	"fmt"
	"os"

	flag "github.com/spf13/pflag"

	"github.com/kha7iq/pvu/internal/app"
	"github.com/kha7iq/pvu/internal/models"
	"github.com/kha7iq/pvu/internal/version"
)

func main() {
	flag.CommandLine.SortFlags = false

	namespace := flag.StringP("namespace", "n", "", "Namespace (defaults to current context)")
	kubeconfig := flag.String("kubeconfig", "", "Path to kubeconfig file")
	showVersion := flag.BoolP("version", "v", false, "Print version information and exit")

	flag.Usage = func() {
		out := flag.CommandLine.Output()

		fmt.Fprintln(out, "pvu is a kubectl plugin for inspecting pod volumes.")
		fmt.Fprintln(out)
		fmt.Fprintln(out, "Usage:")
		fmt.Fprintln(out, "  kubectl pvu [pod-name] [flags]")
		fmt.Fprintln(out)
		fmt.Fprintln(out, "Behavior:")
		fmt.Fprintln(out, "  - With no pod name, starts the interactive TUI pod selector.")
		fmt.Fprintln(out, "  - With a pod name, prints the volume view for that pod.")
		fmt.Fprintln(out)
		fmt.Fprintln(out, "Examples:")
		fmt.Fprintln(out, "  kubectl pvu")
		fmt.Fprintln(out, "  kubectl pvu my-pod")
		fmt.Fprintln(out, "  kubectl pvu my-pod -n kube-system")
		fmt.Fprintln(out, "  kubectl pvu --kubeconfig ~/.kube/config")
		fmt.Fprintln(out)
		fmt.Fprintln(out, "Flags:")
		flag.PrintDefaults()
	}

	flag.Parse()

	if *showVersion {
		fmt.Println(version.Info())
		os.Exit(0)
	}

	if len(flag.Args()) > 1 {
		fmt.Fprintln(os.Stderr, "error: accepts at most one pod name")
		fmt.Fprintln(os.Stderr)
		flag.Usage()
		os.Exit(2)
	}

	var podName string
	if len(flag.Args()) == 1 {
		podName = flag.Args()[0]
	}

	opts := models.Options{
		PodName:    podName,
		Namespace:  *namespace,
		Kubeconfig: *kubeconfig,
	}

	if err := app.Run(opts); err != nil {
		fmt.Fprintf(os.Stderr, "pvu: %v\n", err)
		os.Exit(1)
	}
}
