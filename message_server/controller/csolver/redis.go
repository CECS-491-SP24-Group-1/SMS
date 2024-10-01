package csolver

import (
	"context"
	"fmt"
	"time"

	"wraith.me/message_server/globals"
	"wraith.me/message_server/obj/challenge"
)

// Rejects signed tokens that were already submitted to prevent replay attacks.
func CheckRedis(token *challenge.CToken, ctx context.Context) error {
	//Check if token ID exists in Redis
	tokenID := token.ID.String()
	exists, err := globals.Rcl.Exists(ctx, tokenID).Result()
	if err != nil {
		return fmt.Errorf("error checking token in Redis: %w", err)
	}

	//If the token ID exists in Redis, reject it (replay attack)
	if exists > 0 {
		return fmt.Errorf("token already used")
	}

	//Store the token ID in Redis with an expiration time
	expiration := time.Until(token.Expiry)
	err = globals.Rcl.Set(ctx, tokenID, "used", expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to store token ID in Redis: %w", err)
	}

	//No error so return nil
	return nil
}
