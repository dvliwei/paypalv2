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
	"crypto/tls"
	"fmt"
	"github.com/astaxie/beego/httplib"
	"time"
)

func (c *Client) CreateOrder(ctx context.Context, createOrder CreateOrder) (*Order, error) {
	var createOrderReponse Order
	postUrl:= c.Domain+"/v2/checkout/orders"
	req := httplib.Post(postUrl)
	req.Header("Accept","application/json")
	req.Header("Authorization","Bearer "+c.Token.Token)
	req.JSONBody(createOrder)
	req.SetTLSClientConfig(&tls.Config{InsecureSkipVerify:true})
	req.SetTimeout(100*time.Second, 30*time.Second).Response()
	str ,err:=req.String()
	if err!=nil{
		return &createOrderReponse,err
	}
	err =req.ToJSON(&createOrderReponse)
	fmt.Println(str)
	return &createOrderReponse,nil
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