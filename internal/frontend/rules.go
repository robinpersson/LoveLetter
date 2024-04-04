// -------------------------------------------------------------------
// Copyright (c) Axis Communications AB, SWEDEN. All rights reserved.
// -------------------------------------------------------------------

package frontend

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"github.com/robinpersson/LoveLetter/internal/card"
)

func (ui *UI) RulesView() error {
	maxX, maxY := ui.Size()
	//name := fmt.Sprintf("v%v", 0)
	//0, 0, maxX-25, maxY-30
	v, err := ui.SetView("Rules", 0, 0+3, maxX-1, maxY-1)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Wrap = true

		fmt.Fprintln(v, getRules())
	}
	if _, err := ui.SetCurrentView("Rules"); err != nil {
		return err
	}
	return nil
}

func getRules() string {
	return card.GetPrincessRule() +
		card.GetCountessRule() +
		card.GetkingRule() +
		card.GetChancellorRule() +
		card.GetPrinceRule() +
		card.GetHandmaidRule() +
		card.GetBaronRule() +
		card.GetPriestRule() +
		card.GetGuardRule() +
		card.GetSpyRule()
}
