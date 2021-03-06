/**
 * @ClassName main
 * @Description //TODO 
 * @Author liwei
 * @Date 2021/7/7 18:40
 * @Version example V1.0
 **/

package main

import (
	"context"
	"example/paypalv2/paypal"
	"fmt"
)

func main()  {
	c,err:=paypal.PaypalClient("*****","***",paypal.APIBaseSandBox)
	if err!=nil{
		fmt.Println(err)
	}

	c.GetAccessToken(context.Background())

	var amout  paypal.PurchaseUnitAmount
	amout.Currency="USD"
	amout.Value="15.00"

	var purchaseUnit paypal.PurchaseUnit
	purchaseUnit.Amount =&amout
	purchaseUnit.ReferenceID = "asdasda13sadasdas5"

	var purchaseUnits []paypal.PurchaseUnit
	purchaseUnits = append(purchaseUnits,purchaseUnit)

	var createOrder  paypal.CreateOrder

	createOrder.Intent = "CAPTURE"
	createOrder.PurchaseUnits = purchaseUnits
	createOrder.ApplicationContext.ReturnURL="http://pages.ylwtd.com/paypal-result.html"
	order, err := c.CreateOrder(context.Background(),createOrder)
	if err!=nil{
		fmt.Println(err)
	}
	fmt.Println(order.Links)
}
