package cmd

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	auroraclient "github.com/diamcircle/go/clients/auroraclient"
	hlog "github.com/diamcircle/go/support/log"
)

var DatabaseURL string
var Client *auroraclient.Client
var UseTestNet bool
var Logger = hlog.New()

var defaultDatabaseURL = getEnv("DB_URL", "postgres://localhost:5432/diamcircleticker01?sslmode=disable")

var rootCmd = &cobra.Command{
	Use:   "ticker",
	Short: "Diamcircle Development Foundation Ticker.",
	Long:  `A tool to provide Diamcircle Asset and Market data.`,
}

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(
		&DatabaseURL,
		"db-url",
		"d",
		defaultDatabaseURL,
		"database URL, such as: postgres://user:pass@localhost:5432/ticker",
	)
	rootCmd.PersistentFlags().BoolVar(
		&UseTestNet,
		"testnet",
		false,
		"use the Diamcircle Test Network, instead of the Diamcircle Public Network",
	)

	Logger.SetLevel(logrus.DebugLevel)
}

func initConfig() {
	if UseTestNet {
		Logger.Debug("Using Diamcircle Default Test Network")
		Client = auroraclient.DefaultTestNetClient
	} else {
		Logger.Debug("Using Diamcircle Default Public Network")
		Client = auroraclient.DefaultPublicNetClient
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
