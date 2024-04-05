// -------------------------------------------------------------------
// Copyright (c) Axis Communications AB, SWEDEN. All rights reserved.
// -------------------------------------------------------------------

package frontend

import (
	"fmt"
	"github.com/robinpersson/LoveLetter/internal/chat"
)

func getOpponents(userInfo []chat.UserInfo) string {
	var opp string
	for _, o := range userInfo {
		opp += fmt.Sprintf("F%d. %s\n", o.Order, o.Name)
	}

	return opp
}
