package cmd

import (
	"github.com/KBingsoo/entities/pkg/models"
	//"github.com/KBingsoo/mono/domain/mono"
	"gateways/web"
	"transactions/gateways/database"

	mono "../domain/mono"

	"github.com/joho/godotenv"
	"github.com/literalog/go-wise/wise"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start webserver",
	RunE: func(cmd *cobra.Command, args []string) error {

		err := godotenv.Load(".env")
		if err != nil {
			return err
		}

		col, err := database.GetCollection("cards")
		if err != nil {
			return err
		}

		cardsRepository, err := wise.NewMongoSimpleRepository[models.Card](col)
		if err != nil {
			return err
		}

		ordersRepository, err := wise.NewMongoSimpleRepository[models.Order](col)
		if err != nil {
			return err
		}

		participants := []mono.Participant{
			&mono.LoggingParticipant{},
		}

		service := mono.NewManager(cardsRepository, ordersRepository, participants)

		handler := mono.NewHandler(service)

		server := web.NewServer(handler)

		errCh := make(chan error)

		go func() {
			errCh <- server.Run(8080)
		}()

		return nil
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
