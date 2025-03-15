package general

import (
	"errors"

	"gorm.io/gorm"
)

type Phone struct {
	gorm.Model
	Name        string
	Description string
	Brand       string
	CPU         string
	ScreenSize  string
	Camera      string
	Battery     string
	Storage     string
	Price       float64
	IsUsed      bool
	Issues      string
	Condition
	Image
}

type Condition string

const (
	Excellent Condition = "Excellent"
	Good      Condition = "Good"
	Fair      Condition = "Fair"
	Poor      Condition = "Poor"
)

func ValidateCondition(condition Condition) error {
	switch condition {
	case Excellent, Good, Fair, Poor:
		return nil
	default:
		return errors.New("invalid condition")
	}
}
