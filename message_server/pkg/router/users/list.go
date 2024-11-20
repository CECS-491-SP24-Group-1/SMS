package users

import (
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"wraith.me/message_server/pkg/db/qpage"
	"wraith.me/message_server/pkg/http_types/response"
	"wraith.me/message_server/pkg/schema/user"
	"wraith.me/message_server/pkg/util"
)

// Handles incoming requests made to `GET /api/users/list`.
func UserListRoute(w http.ResponseWriter, r *http.Request) {
	//Get the paging params from the URL query params
	pagingParams := qpage.ParseQuery(r)

	//Construct the pager object
	pager, err := qpage.NewQPage(uc.Collection)
	if err != nil {
		util.ErrResponse(http.StatusInternalServerError, err).Respond(w)
		return
	}

	//Construct the aggregate query and perform the paging query
	users := make([]response.UInfo, 0)
	query := bson.A{user.PublicQuery}
	pagination, err := pager.Aggregate(&users, r.Context(), query, pagingParams)
	if err != nil {
		util.ErrResponse(http.StatusInternalServerError, err).Respond(w)
		return
	}

	//Wrap the users and pagination and return the pagination data
	out := response.NewPaginatedData(users, *pagination)
	util.PayloadOkResponse(out.Desc(), out).Respond(w)
}
