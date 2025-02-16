package porter

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/shopspring/decimal"
)

type Port struct {
	Name        string             `json:"name"        validate:"required"`
	City        string             `json:"city"        validate:"required"`
	Country     string             `json:"country"`
	Alias       []string           `json:"alias"`   // by assumption string type not required
	Regions     []string           `json:"regions"` // by assumption string type not required
	Coordinates [2]decimal.Decimal `json:"coordinates"`
	Province    string             `json:"province"`
	Timezone    string             `json:"timezone"`
	Unlocs      []string           `json:"unlocs"      validate:"required,dive,min=1"`
	Code        string             `json:"code"`
}

func (p *Port) Validate() error {
	validate := validator.New()

	return validate.Struct(p)
}

func (p *Port) Time() (time.Time, error) {
	loc, err := time.LoadLocation(p.Timezone)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to load port timezone: %w", err)
	}

	return time.Now().In(loc), nil
}
