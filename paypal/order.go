/**
 * @ClassName order
 * @Description //TODO 
 * @Author liwei
 * @Date 2021/7/7 18:37
 * @Version example V1.0
 **/

package paypal

import (
	"context"
	"fmt"
)

func (c *Client) CreateOrder(ctx context.Context, intent string, purchaseUnits []PurchaseUnitRequest, payer *CreateOrderPayer, appContext *ApplicationContext) (*Order, error) {
	return c.CreateOrderWithPaypalRequestID(ctx, intent, purchaseUnits, payer, appContext, "")
}


// CreateOrderWithPaypalRequestID - Use this call to create an order with idempotency
// Endpoint: POST /v2/checkout/orders
func (c *Client) CreateOrderWithPaypalRequestID(ctx context.Context,
	intent string,
	purchaseUnits []PurchaseUnitRequest,
	payer *CreateOrderPayer,
	appContext *ApplicationContext,
	requestID string,
) (*Order, error) {
	type createOrderRequest struct {
		Intent             string                `json:"intent"`
		Payer              *CreateOrderPayer     `json:"payer,omitempty"`
		PurchaseUnits      []PurchaseUnitRequest `json:"purchase_units"`
		ApplicationContext *ApplicationContext   `json:"application_context,omitempty"`
	}

	order := &Order{}

	req, err := c.NewRequest(ctx, "POST", fmt.Sprintf("%s%s", c.Domain, "/v2/checkout/orders"), createOrderRequest{Intent: intent, PurchaseUnits: purchaseUnits, Payer: payer, ApplicationContext: appContext})
	if err != nil {
		return order, err
	}

	if requestID != "" {
		req.Header.Set("PayPal-Request-Id", requestID)
	}

	if err = c.SendWithAuth(req, order); err != nil {
		return order, err
	}

	return order, nil
}