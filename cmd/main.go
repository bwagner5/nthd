package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"

	"github.com/bwagner5/nthd/pkg/imds"
	"github.com/bwagner5/nthd/pkg/machine"
	"github.com/jaypipes/envutil"
)

type Options struct {
	DryRun         bool
	MetadataIP     string
	MetadataIPMode string
}

func main() {
	options := MustParseFlags()
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	ctx := context.Background()
	machine, err := machine.New()
	if err != nil {
		log.Fatalln(err)
	}
	imdsAPI, err := imds.NewClient(ctx, options.MetadataIP, options.MetadataIPMode)
	if err != nil {
		log.Fatalln(err)
	}
	ticker := time.NewTicker(time.Second * 5)
	for {
		select {
		case <-sigs:
			return
		case <-ticker.C:
			termTime, ok, err := imdsAPI.GetSpotInterruptionNotification(ctx)
			if err != nil {
				log.Println(err)
				continue
			}
			if !ok {
				continue
			}
			log.Printf("Spot Termination Happening at %s (in %s)", termTime, time.Until(*termTime))
			if options.DryRun {
				log.Println("Would have shutdown but dry-run was set")
				return
			}
			if err := machine.Shutdown(); err != nil {
				log.Fatalln(err)
			}
		}
	}
}

func MustParseFlags() Options {
	options := Options{}
	root := flag.NewFlagSet(path.Base(os.Args[0]), flag.ExitOnError)
	root.BoolVar(&options.DryRun, "dry-run", envutil.WithDefaultBool("DRY_RUN", false), "Do not take action on events received")
	root.StringVar(&options.MetadataIP, "metadata-ip", envutil.WithDefault("METADATA_IP", "http://169.254.169.254"), "The IP address of the instance metadata service")
	root.StringVar(&options.MetadataIPMode, "metadata-ip-mode", envutil.WithDefault("METADATA_IP_MODE", "ipv4"), "IP mode (ipv4 or ipv6) to access the instance metadata service")
	root.Parse(os.Args[1:])
	return options
}
