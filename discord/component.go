package discord

import (
	"fmt"

	"github.com/DisgoOrg/disgo/json"
)

// ComponentType defines different Component(s)
type ComponentType int

// Supported ComponentType(s)
//goland:noinspection GoUnusedConst
const (
	ComponentTypeActionRow = iota + 1
	ComponentTypeButton
	ComponentTypeSelectMenu
)

type Component interface {
	json.Marshaler
	Type() ComponentType
}

type unmarshalComponent struct {
	Component
}

func (u *unmarshalComponent) UnmarshalJSON(data []byte) error {
	var cType struct {
		Type ComponentType `json:"type"`
	}

	if err := json.Unmarshal(data, &cType); err != nil {
		return err
	}

	var (
		component Component
		err       error
	)

	switch cType.Type {
	case ComponentTypeActionRow:
		v := ActionRowComponent{}
		err = json.Unmarshal(data, &v)
		component = v
	case ComponentTypeButton:
		v := ButtonComponent{}
		err = json.Unmarshal(data, &v)
		component = v
	case ComponentTypeSelectMenu:
		v := SelectMenuComponent{}
		err = json.Unmarshal(data, &v)
		component = v
	default:
		return fmt.Errorf("unkown component with type %d received", cType.Type)
	}
	if err != nil {
		return err
	}

	u.Component = component
	return nil
}

type ComponentEmoji struct {
	ID       Snowflake `json:"id,omitempty"`
	Name     string    `json:"name,omitempty"`
	Animated bool      `json:"animated,omitempty"`
}