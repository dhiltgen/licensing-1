package model

import (
	"fmt"
	"time"

	validation "github.com/docker/licensing/lib/go-validation"
)

// PricingComponents represents a collection of pricing components
type PricingComponents []*SubscriptionPricingComponent

func (comps PricingComponents) Len() int      { return len(comps) }
func (comps PricingComponents) Swap(i, j int) { comps[i], comps[j] = comps[j], comps[i] }

// always sorting by name
func (comps PricingComponents) Less(i, j int) bool { return comps[i].Name < comps[j].Name }

// SubscriptionDetail presents Subscription information to billing service clients.
type SubscriptionDetail struct {
	Name               string     `json:"name"`
	ID                 string     `json:"subscription_id"`
	DockerID           string     `json:"docker_id"`
	ProductID          string     `json:"product_id"`
	CreatedByID        string     `json:"created_by_docker_id"`
	ProductRatePlan    string     `json:"product_rate_plan"`
	ProductRatePlanID  string     `json:"product_rate_plan_id"`
	InitialPeriodStart time.Time  `json:"initial_period_start"`
	CurrentPeriodStart time.Time  `json:"current_period_start"`
	CurrentPeriodEnd   *time.Time `json:"current_period_end,omitempty"`

	State       string     `json:"state"`
	Eusa        *EusaState `json:"eusa,omitempty"`
	Origin      string     `json:"origin,omitempty"`
	OrderID     string     `json:"order_id,omitempty"`
	OrderItemID string     `json:"order_item_id,omitempty"`

	PricingComponents PricingComponents `json:"pricing_components"`

	// If true, the product for this subscription uses product keys. To
	// obtain the keys, the frontend or billing client will need to
	// make additional calls to the fulfillment service.
	UsesProductKeys bool `json:"uses_product_keys,omitempty"`

	// If non-empty, this is a managed subscription, and this identifier is
	// known to the fulfillment service as a means to uniquely identify the
	// partner that manages this subscription.
	//
	// Different permissions checking will be used to authorize changes and
	// cancellation; the entity entitled to this subscription (represented
	// by DockerID) may not change or cancel it directly.
	ManagingPartnerID string `json:"managing_partner_id,omitempty"`

	// If non-empty, this is a managed subscription, and this ID belongs to the
	// account of a user within a partner's account system.
	PartnerAccountID string `json:"partner_account_id,omitempty"`

	// Marketing opt-in for the subscription. This means customer agrees to receive additional marketing emails
	MarketingOptIn bool `json:"marketing_opt_in"`
}

// SubscriptionPricingComponent captures pricing component values that have been selected by the user.
type SubscriptionPricingComponent struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

// SubscriptionCreationRequest represents a subscription creation request
type SubscriptionCreationRequest struct {
	Name     string `json:"name"`
	DockerID string `json:"docker_id"`

	ProductID       string     `json:"product_id"`
	ProductRatePlan string     `json:"product_rate_plan"`
	Eusa            *EusaState `json:"eusa,omitempty"`
	Origin          string     `json:"origin,omitempty"`
	OrderID         string     `json:"order_id,omitempty"`
	OrderItemID     string     `json:"order_item_id,omitempty"`

	End   *time.Time `json:"end,omitempty"`
	Start *time.Time `json:"start,omitempty"`

	CouponCodes []string `json:"coupon_codes"`

	PricingComponents PricingComponents `json:"pricing_components"`

	// If true, the product for this subscription uses product keys. To
	// obtain the keys, the frontend or billing client will need to
	// make additional calls to the fulfillment service.
	UsesProductKeys bool `json:"uses_product_keys,omitempty"`

	// Should be non-empty only if creating a managed subscription that will
	// be controlled by a partner or publisher. This identifier matches
	// whatever the fulfillment service uses as guid's for partners.
	ManagingPartnerID string `json:"managing_partner_id,omitempty"`

	// Should be non-empty only if creating a managed subscription on behalf
	// of a partner, and this ID represent's a partner's user's account id.
	PartnerAccountID string `json:"partner_account_id,omitempty"`

	// Marketing opt-in for the subscription. This means customer agrees to receive additional marketing emails
	MarketingOptIn bool `json:"marketing_opt_in"`
}

// Validate returns true if the subscription request is valid, false otherwise.
// If invalid, one or more validation Errors will be returned.
func (s *SubscriptionCreationRequest) Validate() (bool, validation.Errors) {
	var errs validation.Errors

	if validation.IsEmpty(s.Name) {
		errs = append(errs, validation.InvalidEmpty("name"))
	}

	if validation.IsEmpty(s.DockerID) {
		errs = append(errs, validation.InvalidEmpty("docker_id"))
	}

	if validation.IsEmpty(s.ProductID) {
		errs = append(errs, validation.InvalidEmpty("product_id"))
	}

	if validation.IsEmpty(s.ProductRatePlan) {
		errs = append(errs, validation.InvalidEmpty("product_rate_plan"))
	}

	for i, component := range s.PricingComponents {
		if validation.IsEmpty(component.Name) {
			name := fmt.Sprintf("pricing_component[%v]/name", i)
			errs = append(errs, validation.InvalidEmpty(name))
		}
	}

	valid := len(errs) == 0
	return valid, errs
}

// EusaState encodes whether the subscription's EUSA has been accepted,
// and if so, by whom and when.
// See json marshal & unmarshal below.
type EusaState struct {
	Accepted   bool   `json:"accepted"`
	AcceptedBy string `json:"accepted_by,omitempty"`
	AcceptedOn string `json:"accepted_on,omitempty"`
}
