package client

import (
	"context"
	"fmt"

	"github.com/gotd/td/telegram/peers"
)

type User struct {
	ID       int
	Username string
}

type UserRepo struct {
	peerManager *peers.Manager
}

func NewUserRepo(peerManager *peers.Manager) *UserRepo {
	return &UserRepo{peerManager: peerManager}
}

func (u *UserRepo) Get(ctx context.Context, userID int) (User, error) {
	userPeer, err := u.peerManager.ResolveUserID(ctx, int64(userID))
	if err != nil {
		return User{}, fmt.Errorf("failed to get user via peerManager: %w", err)
	}

	return User{
		ID:       int(userPeer.Raw().ID),
		Username: userPeer.Raw().Username,
	}, nil
}
