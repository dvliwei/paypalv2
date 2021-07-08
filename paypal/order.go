/**
 * @ClassName order
 * @Description //TODO 
 * @Author liwei
 * @Date 2021/7/7 18:37
 * @Version example V1.0
 **/

package paypal

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

func (c *Client) CreateOrder(ctx context.Context, createOrder CreateOrder) (*Order, error) {
	var createOrderReponse Order
	postUrl:= c.Domain+"/v2/checkout/orders"
	var buf io.Reader
	jsonStr, err := json.Marshal(&createOrder)
	if err!=nil{
		return nil,err
	}
	buf = bytes.NewBuffer(jsonStr)
	tr := &http.Transport{
		MaxIdleConns:       30,
		IdleConnTimeout:    100 * time.Second,
		DisableCompression: true,
		TLSClientConfig:&tls.Config{InsecureSkipVerify:true},
	}
	client := &http.Client{
		Transport:tr,
	}
	req,err:=http.NewRequest("POST",postUrl,ioutil.NopCloser(buf))
	if err!=nil{
		return nil,err
	}
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	req.Header.Add("Authorization","Bearer "+c.Token.Token)

	resp, err := client.Do(req)
	if err!=nil{
		return nil,err
	}
	defer resp.Body.Close()
	//
	//fmt.Println("response Status:", resp.Status)
	//fmt.Println("response Headers:", resp.Header)
	str,err:=ioutil.ReadAll(resp.Body)
	if err!=nil{
		return nil,err
	}
	err1:=json.Unmarshal(str,&createOrderReponse)
	if err1!=nil{
		return nil,err
	}
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