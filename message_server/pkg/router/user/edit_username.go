package user

import (
	"encoding/json"

	"context"
	"fmt"
	"net/http"
	"regexp"

	"go.mongodb.org/mongo-driver/bson"
	"wraith.me/message_server/pkg/db/mongoutil"
	"wraith.me/message_server/pkg/mw"
	"wraith.me/message_server/pkg/schema/user"
	"wraith.me/message_server/pkg/util"
)

// Handles incoming requests made to `PATCH /api/user/username`.
func ChangeUnameRoute(w http.ResponseWriter, r *http.Request) {
	var req struct {
		NewUsername string `json:"username"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.ErrResponse(http.StatusBadRequest, err).Respond(w) // Using ErrResponse from util
		return
	}

	// Validate the username
	if err := isValidUsername(req.NewUsername); err != nil {
		util.ErrResponse(http.StatusBadRequest, err).Respond(w)
		return
	}

	// Get the requestor's info
	user := r.Context().Value(mw.AuthCtxUserKey).(user.User)

	// Check for username uniqueness
	usernameExists, err := ensureUniqueUsername(uc, req.NewUsername, r.Context())
	if err != nil {
		util.ErrResponse(http.StatusInternalServerError, err).Respond(w)
		return
	}
	if usernameExists {
		util.ErrResponse(http.StatusConflict, fmt.Errorf("username already exists")).Respond(w)
		return
	}

	// Update the username in the database
	filter := bson.M{"_id": user.ID}
	update := bson.M{"$set": bson.M{
		"username":     req.NewUsername,
		"display_name": req.NewUsername,
	}}
	if err := uc.UpdateOne(r.Context(), filter, update); err != nil {
		util.ErrResponse(http.StatusInternalServerError, err).Respond(w)
		return
	}

	// Update the current Go User object
	user.Username = req.NewUsername
	user.DisplayName = req.NewUsername

	util.PayloadOkResponse("username changed successfully", user).Respond(w)
}

/*
Ensures that a user doesn't already exist in the database based on what
was given by the user. A `nil` error indicates that no matching records
were found. Checking collections for existent objects is expensive, so
not all records are checked if one fails.
*/
func ensureUniqueUsername(coll *user.UserCollection, username string, ctx context.Context) (bool, error) {
	// Create an aggregation pipeline to check for existing usernames
	agg := bson.A{
		bson.D{
			{Key: "$match",
				Value: bson.D{
					{Key: "username", Value: username},
				},
			},
		},
		bson.D{
			{Key: "$project",
				Value: bson.D{
					{Key: "_id", Value: 1},
				},
			},
		},
	}

	// Run the aggregation and check if any results are returned
	hits, err := mongoutil.Aggregate2IDArr(coll.Aggregate(ctx, agg))
	if err != nil {
		return false, err
	}

	// Return true if at least one user with the given username exists
	return len(hits) > 0, nil
}

// Function to validate that the username is alphanumeric with underscores allowed
func isValidUsername(username string) error {
	var validUsernamePattern = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)

	if len(username) < 4 || len(username) > 16 {
		return fmt.Errorf("username must be between 4 and 16 characters")
	}

	if !validUsernamePattern.MatchString(username) {
		return fmt.Errorf("username can only contain alphanumeric characters and underscores")
	}

	return nil
}
