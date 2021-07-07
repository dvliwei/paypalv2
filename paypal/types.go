/**
 * @ClassName types
 * @Description //TODO 
 * @Author liwei
 * @Date 2021/7/7 17:53
 * @Version example V1.0
 **/

package paypal

import (
	"io"
	"net/http"
	"sync"
	"time"
)

const(
	APIBaseSandBox="https://api.sandbox.paypal.com"
	// APIBaseLive points to the live version of the API
	APIBaseLive = "https://api.paypal.com"

	// RequestNewTokenBeforeExpiresIn is used by SendWithAuth and try to get new Token when it's about to expire
	RequestNewTokenBeforeExpiresIn = time.Duration(60) * time.Second
)

const (
	EventCheckoutOrderApproved         string = "CHECKOUT.ORDER.APPROVED"
	EventPaymentCaptureCompleted       string = "PAYMENT.CAPTURE.COMPLETED"
	EventPaymentCaptureDenied          string = "PAYMENT.CAPTURE.DENIED"
	EventPaymentCaptureRefunded        string = "PAYMENT.CAPTURE.REFUNDED"
	EventMerchantOnboardingCompleted   string = "MERCHANT.ONBOARDING.COMPLETED"
	EventMerchantPartnerConsentRevoked string = "MERCHANT.PARTNER-CONSENT.REVOKED"
)

const (
	OrderIntentCapture   string = "CAPTURE"
	OrderIntentAuthorize string = "AUTHORIZE"
)


// Amount struct
type (

	JSONTime time.Time

	Amount struct {
			Currency string  `json:"currency"`
			Total    string  `json:"total"`
			Details  Details `json:"details,omitempty"`
	}

	// AmountPayout struct
	AmountPayout struct {
		Currency string `json:"currency"`
		Value    string `json:"value"`
	}

	// ApplicationContext struct
	//Doc: https://developer.paypal.com/docs/api/subscriptions/v1/#definition-application_context
	ApplicationContext struct {
		BrandName          string             `json:"brand_name,omitempty"`
		Locale             string             `json:"locale,omitempty"`
		//LandingPage        string `json:"landing_page,omitempty"` // not found in documentation
		ReturnURL string `json:"return_url,omitempty"`
		CancelURL string `json:"cancel_url,omitempty"`
	}

	// Authorization struct
	Authorization struct {
		ID               string                `json:"id,omitempty"`
		CustomID         string                `json:"custom_id,omitempty"`
		InvoiceID        string                `json:"invoice_id,omitempty"`
		Status           string                `json:"status,omitempty"`
		Amount           *PurchaseUnitAmount   `json:"amount,omitempty"`
		CreateTime       *time.Time            `json:"create_time,omitempty"`
		UpdateTime       *time.Time            `json:"update_time,omitempty"`
		ExpirationTime   *time.Time            `json:"expiration_time,omitempty"`
		Links            []Link                `json:"links,omitempty"`
	}

	// https://developer.paypal.com/docs/api/payments/v2/#definition-platform_fee
	PlatformFee struct {
		Amount *Money          `json:"amount,omitempty"`
		Payee  *PayeeForOrders `json:"payee,omitempty"`
	}

	//https://developer.paypal.com/docs/api/payments/v2/#captures_get
	CaptureDetailsResponse struct {
		Status                    string                     `json:"status,omitempty"`
		ID                        string                     `json:"id,omitempty"`
		Amount                    *Money                     `json:"amount,omitempty"`
		InvoiceID                 string                     `json:"invoice_id,omitempty"`
		CustomID                  string                     `json:"custom_id,omitempty"`
		FinalCapture              bool                       `json:"final_capture,omitempty"`
		DisbursementMode          string                     `json:"disbursement_mode,omitempty"`
		Links                     []Link                     `json:"links,omitempty"`
		UpdateTime                *time.Time                 `json:"update_time,omitempty"`
		CreateTime                *time.Time                 `json:"create_time,omitempty"`
	}
	PaymentSource struct {
		Card  *PaymentSourceCard  `json:"card"`
		Token *PaymentSourceToken `json:"token"`
	}
	PaymentSourceCard struct {
		ID             string              `json:"id"`
		Name           string              `json:"name"`
		Number         string              `json:"number"`
		Expiry         string              `json:"expiry"`
		SecurityCode   string              `json:"security_code"`
		LastDigits     string              `json:"last_digits"`
		CardType       string              `json:"card_type"`
	}
	PaymentSourceToken struct {
		ID   string `json:"id"`
		Type string `json:"type"`
	}
	// CaptureOrderRequest - https://developer.paypal.com/docs/api/orders/v2/#orders_capture
	CaptureOrderRequest struct {
		PaymentSource *PaymentSource `json:"payment_source"`
	}

	// SenderBatchHeader struct
	SenderBatchHeader struct {
		EmailSubject  string `json:"email_subject"`
		EmailMessage  string `json:"email_message"`
		SenderBatchID string `json:"sender_batch_id,omitempty"`
	}
	// BatchHeader struct
	BatchHeader struct {
		Amount            *AmountPayout      `json:"amount,omitempty"`
		Fees              *AmountPayout      `json:"fees,omitempty"`
		PayoutBatchID     string             `json:"payout_batch_id,omitempty"`
		BatchStatus       string             `json:"batch_status,omitempty"`
		TimeCreated       *time.Time         `json:"time_created,omitempty"`
		TimeCompleted     *time.Time         `json:"time_completed,omitempty"`
		SenderBatchHeader *SenderBatchHeader `json:"sender_batch_header,omitempty"`
	}

	// Capture struct
	Capture struct {
		ID             string     `json:"id,omitempty"`
		Amount         *Amount    `json:"amount,omitempty"`
		State          string     `json:"state,omitempty"`
		ParentPayment  string     `json:"parent_payment,omitempty"`
		TransactionFee string     `json:"transaction_fee,omitempty"`
		IsFinalCapture bool       `json:"is_final_capture"`
		CreateTime     *time.Time `json:"create_time,omitempty"`
		UpdateTime     *time.Time `json:"update_time,omitempty"`
		Links          []Link     `json:"links,omitempty"`
	}

	// Client represents a Paypal REST API Client
	Client struct {
		sync.Mutex
		Client               *http.Client
		ClientID             string
		Secret               string
		Domain              string
		Log                  io.Writer // If user set log file name all requests will be logged there
		Token                *TokenResponse
		tokenExpiresAt       time.Time
		returnRepresentation bool
	}

	// Currency struct
	Currency struct {
		Currency string `json:"currency,omitempty"`
		Value    string `json:"value,omitempty"`
	}

	// LastPayment struct
	LastPayment struct {
		Amount Money     `json:"amount,omitempty"`
		Time   time.Time `json:"time,omitempty"`
	}

	// Details structure used in Amount structures as optional value
	Details struct {
		Subtotal         string `json:"subtotal,omitempty"`
		Shipping         string `json:"shipping,omitempty"`
		Tax              string `json:"tax,omitempty"`
		HandlingFee      string `json:"handling_fee,omitempty"`
		ShippingDiscount string `json:"shipping_discount,omitempty"`
		Insurance        string `json:"insurance,omitempty"`
		GiftWrap         string `json:"gift_wrap,omitempty"`
	}


	// Item struct
	Item struct {
		Name        string `json:"name"`
		UnitAmount  *Money `json:"unit_amount,omitempty"`
		Tax         *Money `json:"tax,omitempty"`
		Quantity    string `json:"quantity"`
		Description string `json:"description,omitempty"`
		SKU         string `json:"sku,omitempty"`
		Category    string `json:"category,omitempty"`
	}

	// ItemList struct
	ItemList struct {
		Items           []Item           `json:"items,omitempty"`
		ShippingAddress *ShippingAddress `json:"shipping_address,omitempty"`
	}

	// Link struct
	Link struct {
		Href        string `json:"href"`
		Rel         string `json:"rel,omitempty"`
		Method      string `json:"method,omitempty"`
		Description string `json:"description,omitempty"`
		Enctype     string `json:"enctype,omitempty"`
	}

	// PurchaseUnitAmount struct
	PurchaseUnitAmount struct {
		Currency  string                       `json:"currency_code"`
		Value     string                       `json:"value"`
		Breakdown *PurchaseUnitAmountBreakdown `json:"breakdown,omitempty"`
	}

	// PurchaseUnitAmountBreakdown struct
	PurchaseUnitAmountBreakdown struct {
		ItemTotal        *Money `json:"item_total,omitempty"`
		Shipping         *Money `json:"shipping,omitempty"`
		Handling         *Money `json:"handling,omitempty"`
		TaxTotal         *Money `json:"tax_total,omitempty"`
		Insurance        *Money `json:"insurance,omitempty"`
		ShippingDiscount *Money `json:"shipping_discount,omitempty"`
		Discount         *Money `json:"discount,omitempty"`
	}

	// Money struct
	//
	// https://developer.paypal.com/docs/api/orders/v2/#definition-money
	Money struct {
		Currency string `json:"currency_code"`
		Value    string `json:"value"`
	}

	// TaxInfo used for orders.
	TaxInfo struct {
		TaxID     string `json:"tax_id,omitempty"`
		TaxIDType string `json:"tax_id_type,omitempty"`
	}

	// CreateOrderPayerName create order payer name
	CreateOrderPayerName struct {
		GivenName string `json:"given_name,omitempty"`
		Surname   string `json:"surname,omitempty"`
	}

	// CreateOrderPayer used with create order requests
	CreateOrderPayer struct {
		Name         *CreateOrderPayerName          `json:"name,omitempty"`
		EmailAddress string                         `json:"email_address,omitempty"`
		PayerID      string                         `json:"payer_id,omitempty"`
		BirthDate    string                         `json:"birth_date,omitempty"`
		TaxInfo      *TaxInfo                       `json:"tax_info,omitempty"`
		Address      *ShippingDetailAddressPortable `json:"address,omitempty"`
	}
	PayerWithNameAndPhone struct {
		Name         *CreateOrderPayerName `json:"name,omitempty"`
		EmailAddress string                `json:"email_address,omitempty"`
		PayerID      string                `json:"payer_id,omitempty"`
		Address      Address               `json:"address,omitempty"`
	}
	PurchaseUnit struct {
		ReferenceID        string              `json:"reference_id"`
		Amount             *PurchaseUnitAmount `json:"amount,omitempty"`
		Payments           *CapturedPayments   `json:"payments,omitempty"`
	}
	// Order struct
	Order struct {
		ID            string                 `json:"id,omitempty"`
		Status        string                 `json:"status,omitempty"`
		Intent        string                 `json:"intent,omitempty"`
		Payer         *PayerWithNameAndPhone `json:"payer,omitempty"`
		PurchaseUnits []PurchaseUnit         `json:"purchase_units,omitempty"`
		Links         []Link                 `json:"links,omitempty"`
		CreateTime    *time.Time             `json:"create_time,omitempty"`
		UpdateTime    *time.Time             `json:"update_time,omitempty"`
	}

	// CaptureAmount struct
	CaptureAmount struct {
		ID                        string                     `json:"id,omitempty"`
		CustomID                  string                     `json:"custom_id,omitempty"`
		Amount                    *PurchaseUnitAmount        `json:"amount,omitempty"`
	}

	// CapturedPayments has the amounts for a captured order
	CapturedPayments struct {
		Captures []CaptureAmount `json:"captures,omitempty"`
	}

	// CapturedPurchaseItem are items for a captured order
	CapturedPurchaseItem struct {
		Quantity    string `json:"quantity"`
		Name        string `json:"name"`
		SKU         string `json:"sku,omitempty"`
		Description string `json:"description,omitempty"`
	}

	// RedirectURLs struct
	RedirectURLs struct {
		ReturnURL string `json:"return_url,omitempty"`
		CancelURL string `json:"cancel_url,omitempty"`
	}

	// Sale struct
	Sale struct {
		ID                        string     `json:"id,omitempty"`
		Amount                    *Amount    `json:"amount,omitempty"`
		TransactionFee            *Currency  `json:"transaction_fee,omitempty"`
		Description               string     `json:"description,omitempty"`
		CreateTime                *time.Time `json:"create_time,omitempty"`
		State                     string     `json:"state,omitempty"`
		ParentPayment             string     `json:"parent_payment,omitempty"`
		UpdateTime                *time.Time `json:"update_time,omitempty"`
		PaymentMode               string     `json:"payment_mode,omitempty"`
		PendingReason             string     `json:"pending_reason,omitempty"`
		ReasonCode                string     `json:"reason_code,omitempty"`
		ClearingTime              string     `json:"clearing_time,omitempty"`
		ProtectionEligibility     string     `json:"protection_eligibility,omitempty"`
		ProtectionEligibilityType string     `json:"protection_eligibility_type,omitempty"`
		Links                     []Link     `json:"links,omitempty"`
	}

	//ShippingAmount struct
	ShippingAmount struct {
		Money
	}

	// ShippingAddress struct
	ShippingAddress struct {
		RecipientName string `json:"recipient_name,omitempty"`
		Type          string `json:"type,omitempty"`
		Line1         string `json:"line1"`
		Line2         string `json:"line2,omitempty"`
		City          string `json:"city"`
		CountryCode   string `json:"country_code"`
		PostalCode    string `json:"postal_code,omitempty"`
		State         string `json:"state,omitempty"`
		Phone         string `json:"phone,omitempty"`
	}

	// ShippingDetailAddressPortable used with create orders
	ShippingDetailAddressPortable struct {
		AddressLine1 string `json:"address_line_1,omitempty"`
		AddressLine2 string `json:"address_line_2,omitempty"`
		AdminArea1   string `json:"admin_area_1,omitempty"`
		AdminArea2   string `json:"admin_area_2,omitempty"`
		PostalCode   string `json:"postal_code,omitempty"`
		CountryCode  string `json:"country_code,omitempty"`
	}

	// Name struct
	//Doc: https://developer.paypal.com/docs/api/subscriptions/v1/#definition-name
	Name struct {
		FullName   string `json:"full_name,omitempty"`
		Suffix     string `json:"suffix,omitempty"`
		Prefix     string `json:"prefix,omitempty"`
		GivenName  string `json:"given_name,omitempty"`
		Surname    string `json:"surname,omitempty"`
		MiddleName string `json:"middle_name,omitempty"`
	}

	expirationTime int64

	// TokenResponse is for API response for the /oauth2/token endpoint
	TokenResponse struct {
		RefreshToken string         `json:"refresh_token"`
		Token        string         `json:"access_token"`
		Type         string         `json:"token_type"`
		ExpiresIn    expirationTime `json:"expires_in"`
	}
	PaymentOptions struct {
		AllowedPaymentMethod string `json:"allowed_payment_method,omitempty"`
	}
	Related struct {
		Sale          *Sale          `json:"sale,omitempty"`
		Authorization *Authorization `json:"authorization,omitempty"`
		Order         *Order         `json:"order,omitempty"`
		Capture       *Capture       `json:"capture,omitempty"`
	}
	// Transaction struct
	Transaction struct {
		Amount           *Amount         `json:"amount"`
		Description      string          `json:"description,omitempty"`
		ItemList         *ItemList       `json:"item_list,omitempty"`
		InvoiceNumber    string          `json:"invoice_number,omitempty"`
		Custom           string          `json:"custom,omitempty"`
		SoftDescriptor   string          `json:"soft_descriptor,omitempty"`
		RelatedResources []Related       `json:"related_resources,omitempty"`
		PaymentOptions   *PaymentOptions `json:"payment_options,omitempty"`
		NotifyURL        string          `json:"notify_url,omitempty"`
		OrderURL         string          `json:"order_url,omitempty"`
		Payee            *Payee          `json:"payee,omitempty"`
	}

	//Payee struct
	Payee struct {
		Email string `json:"email"`
	}

	// PayeeForOrders struct
	PayeeForOrders struct {
		EmailAddress string `json:"email_address,omitempty"`
		MerchantID   string `json:"merchant_id,omitempty"`
	}

	// Webhook struct
	Webhook struct {
		ID         string             `json:"id"`
		URL        string             `json:"url"`
		EventTypes []WebhookEventType `json:"event_types"`
		Links      []Link             `json:"links"`
	}

	// Event struct.
	//
	// The basic webhook event data type. This struct is intended to be
	// embedded into resource type specific event structs.
	Event struct {
		ID              string    `json:"id"`
		CreateTime      time.Time `json:"create_time"`
		ResourceType    string    `json:"resource_type"`
		EventType       string    `json:"event_type"`
		Summary         string    `json:"summary,omitempty"`
		Links           []Link    `json:"links"`
		EventVersion    string    `json:"event_version,omitempty"`
		ResourceVersion string    `json:"resource_version,omitempty"`
	}

	// WebhookEventType struct
	WebhookEventType struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Status      string `json:"status,omitempty"`
	}
	SearchPayerName struct {
		GivenName string `json:"given_name"`
		Surname   string `json:"surname"`
	}

	SearchPayerInfo struct {
		AccountID     string               `json:"account_id"`
		EmailAddress  string               `json:"email_address"`
		AddressStatus string               `json:"address_status"`
		PayerStatus   string               `json:"payer_status"`
		PayerName     SearchPayerName      `json:"payer_name"`
		CountryCode   string               `json:"country_code"`
	}

	SearchTaxAmount struct {
		TaxAmount Money `json:"tax_amount"`
	}

	SearchTransactionInfo struct {
		PayPalAccountID           string   `json:"paypal_account_id"`
		TransactionID             string   `json:"transaction_id"`
		PayPalReferenceID         string   `json:"paypal_reference_id"`
		PayPalReferenceIDType     string   `json:"paypal_reference_id_type"`
		TransactionEventCode      string   `json:"transaction_event_code"`
		TransactionInitiationDate JSONTime `json:"transaction_initiation_date"`
		TransactionUpdatedDate    JSONTime `json:"transaction_updated_date"`
		TransactionAmount         Money    `json:"transaction_amount"`
		FeeAmount                 *Money   `json:"fee_amount"`
		InsuranceAmount           *Money   `json:"insurance_amount"`
		ShippingAmount            *Money   `json:"shipping_amount"`
		ShippingDiscountAmount    *Money   `json:"shipping_discount_amount"`
		ShippingTaxAmount         *Money   `json:"shipping_tax_amount"`
		OtherAmount               *Money   `json:"other_amount"`
		TipAmount                 *Money   `json:"tip_amount"`
		TransactionStatus         string   `json:"transaction_status"`
		TransactionSubject        string   `json:"transaction_subject"`
		PaymentTrackingID         string   `json:"payment_tracking_id"`
		BankReferenceID           string   `json:"bank_reference_id"`
		TransactionNote           string   `json:"transaction_note"`
		EndingBalance             *Money   `json:"ending_balance"`
		AvailableBalance          *Money   `json:"available_balance"`
		InvoiceID                 string   `json:"invoice_id"`
		CustomField               string   `json:"custom_field"`
		ProtectionEligibility     string   `json:"protection_eligibility"`
		CreditTerm                string   `json:"credit_term"`
		CreditTransactionalFee    *Money   `json:"credit_transactional_fee"`
		CreditPromotionalFee      *Money   `json:"credit_promotional_fee"`
		AnnualPercentageRate      string   `json:"annual_percentage_rate"`
		PaymentMethodType         string   `json:"payment_method_type"`
	}


	SearchItemDetails struct {
		ItemCode            string                 `json:"item_code"`
		ItemName            string                 `json:"item_name"`
		ItemDescription     string                 `json:"item_description"`
		ItemOptions         string                 `json:"item_options"`
		ItemQuantity        string                 `json:"item_quantity"`
		ItemUnitPrice       Money                  `json:"item_unit_price"`
		ItemAmount          Money                  `json:"item_amount"`
		DiscountAmount      *Money                 `json:"discount_amount"`
		AdjustmentAmount    *Money                 `json:"adjustment_amount"`
		GiftWrapAmount      *Money                 `json:"gift_wrap_amount"`
		TaxPercentage       string                 `json:"tax_percentage"`
		TaxAmounts          []SearchTaxAmount      `json:"tax_amounts"`
		BasicShippingAmount *Money                 `json:"basic_shipping_amount"`
		ExtraShippingAmount *Money                 `json:"extra_shipping_amount"`
		HandlingAmount      *Money                 `json:"handling_amount"`
		InsuranceAmount     *Money                 `json:"insurance_amount"`
		TotalItemAmount     Money                  `json:"total_item_amount"`
		InvoiceNumber       string                 `json:"invoice_number"`
	}

	SearchCartInfo struct {
		ItemDetails     []SearchItemDetails `json:"item_details"`
		TaxInclusive    *bool               `json:"tax_inclusive"`
		PayPalInvoiceID string              `json:"paypal_invoice_id"`
	}
	Address struct {
		Line1       string `json:"line1,omitempty"`
		Line2       string `json:"line2,omitempty"`
		City        string `json:"city,omitempty"`
		CountryCode string `json:"country_code,omitempty"`
		PostalCode  string `json:"postal_code,omitempty"`
		State       string `json:"state,omitempty"`
		Phone       string `json:"phone,omitempty"`
	}
	SearchShippingInfo struct {
		Name                     string   `json:"name"`
		Method                   string   `json:"method"`
		Address                  Address  `json:"address"`
		SecondaryShippingAddress *Address `json:"secondary_shipping_address"`
	}

	SearchTransactionDetails struct {
		TransactionInfo SearchTransactionInfo `json:"transaction_info"`
		PayerInfo       *SearchPayerInfo      `json:"payer_info"`
		ShippingInfo    *SearchShippingInfo   `json:"shipping_info"`
		CartInfo        *SearchCartInfo       `json:"cart_info"`
	}

	SharedResponse struct {
		CreateTime string `json:"create_time"`
		UpdateTime string `json:"update_time"`
		Links      []Link `json:"links"`
	}

	ListParams struct {
		Page          string `json:"page,omitempty"`           //Default: 0.
		PageSize      string `json:"page_size,omitempty"`      //Default: 10.
		TotalRequired string `json:"total_required,omitempty"` //Default: no.
	}

	SharedListResponse struct {
		TotalItems int    `json:"total_items,omitempty"`
		TotalPages int    `json:"total_pages,omitempty"`
		Links      []Link `json:"links,omitempty"`
	}

	// https://developer.paypal.com/docs/api/payments/v2/#definition-payment_instruction
	PaymentInstruction struct {
		PlatformFees     []PlatformFee `json:"platform_fees,omitempty"`
		DisbursementMode string        `json:"disbursement_mode,omitempty"`
	}

	// PurchaseUnitRequest struct
	PurchaseUnitRequest struct {
		ReferenceID        string              `json:"reference_id,omitempty"`
		Amount             *PurchaseUnitAmount `json:"amount"`
		Payee              *PayeeForOrders     `json:"payee,omitempty"`
		Description        string              `json:"description,omitempty"`
		CustomID           string              `json:"custom_id,omitempty"`
		InvoiceID          string              `json:"invoice_id,omitempty"`
		SoftDescriptor     string              `json:"soft_descriptor,omitempty"`
		Items              []Item              `json:"items,omitempty"`
		Shipping           *ShippingDetail     `json:"shipping,omitempty"`
		PaymentInstruction *PaymentInstruction `json:"payment_instruction,omitempty"`
	}

	// ShippingDetail struct
	ShippingDetail struct {
		Name    *Name                          `json:"name,omitempty"`
		Address *ShippingDetailAddressPortable `json:"address,omitempty"`
	}


	// ErrorResponse https://developer.paypal.com/docs/api/errors/
	ErrorResponse struct {
		Response        *http.Response        `json:"-"`
		Name            string                `json:"name"`
		DebugID         string                `json:"debug_id"`
		Message         string                `json:"message"`
		InformationLink string                `json:"information_link"`
		Details         []ErrorResponseDetail `json:"details"`
	}

	// ErrorResponseDetail struct
	ErrorResponseDetail struct {
		Field string `json:"field"`
		Issue string `json:"issue"`
		Links []Link `json:"link"`
	}
)


