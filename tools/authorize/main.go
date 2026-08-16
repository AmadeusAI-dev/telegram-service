package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/auth"
	"github.com/gotd/td/tg"
)

func main() {
	phone := os.Getenv("PHONE")
	password := os.Getenv("PASSWORD")
	ctx := context.Background()

	codePrompt := func(ctx context.Context, sentCode *tg.AuthSentCode) (string, error) {
		fmt.Print("Enter code: ")
		code, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			return "", err
		}
		return strings.TrimSpace(code), nil
	}

	client, err := telegram.ClientFromEnvironment(telegram.Options{})
	if err != nil {
		log.Fatalf("failed to create telegram client: %v", err)
	}

	flow := auth.NewFlow(
		auth.Constant(phone, password, auth.CodeAuthenticatorFunc(codePrompt)),
		auth.SendCodeOptions{},
	)

	err = client.Run(ctx, func(ctx context.Context) error {
		if err := flow.Run(ctx, client.Auth()); err != nil {
			log.Fatalf("failed to authorize: %v", err)
		}

		self, err := client.Self(ctx)
		if err != nil {
			log.Fatalf("failed to get authorized user: %v", err)
		}

		fmt.Printf("Logged as: %s \n", self.Username)

		return nil
	})

}
