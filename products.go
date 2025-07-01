package bigcommerce

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// include_fields
var productFields = []string{"name", "sku", "custom_url", "is_visible", "price"}

// include (subresources, like variants images custom_fields bulk_pricing_rules primary_image modifiers options videos)
var productInclude []string

// extra arguments for product interface
var productArgs map[string]string

// Product is a BigCommerce product object
type Product struct {
	ID                      int64         `json:"id,omitempty"`
	Name                    string        `json:"name,omitempty"`
	Type                    string        `json:"type,omitempty"`
	Sku                     string        `json:"sku,omitempty"`
	Description             string        `json:"description,omitempty"`
	Weight                  float64       `json:"weight,omitempty"`
	Width                   float64       `json:"width,omitempty"`
	Depth                   float64       `json:"depth,omitempty"`
	Height                  float64       `json:"height,omitempty"`
	Price                   float64       `json:"price"`
	CostPrice               float64       `json:"cost_price"`
	RetailPrice             float64       `json:"retail_price,omitempty"`
	SalePrice               float64       `json:"sale_price,omitempty"`
	MapPrice                float64       `json:"map_price,omitempty"`
	TaxClassID              int64         `json:"tax_class_id"`
	ProductTaxCode          string        `json:"product_tax_code,omitempty"`
	CalculatedPrice         float64       `json:"calculated_price,omitempty"`
	Categories              []interface{} `json:"categories,omitempty"`
	BrandID                 int64         `json:"brand_id,omitempty"`
	OptionSetID             interface{}   `json:"option_set_id,omitempty"`
	OptionSetDisplay        string        `json:"option_set_display,omitempty"`
	InventoryLevel          int           `json:"inventory_level,omitempty"`
	InventoryWarningLevel   int           `json:"inventory_warning_level,omitempty"`
	InventoryTracking       string        `json:"inventory_tracking,omitempty"`
	ReviewsRatingSum        int           `json:"reviews_rating_sum,omitempty"`
	ReviewsCount            int           `json:"reviews_count,omitempty"`
	TotalSold               int           `json:"total_sold,omitempty"`
	FixedCostShippingPrice  float64       `json:"fixed_cost_shipping_price,omitempty"`
	IsFreeShipping          bool          `json:"is_free_shipping,omitempty"`
	IsVisible               bool          `json:"is_visible"`
	IsFeatured              bool          `json:"is_featured,omitempty"`
	RelatedProducts         []int         `json:"related_products,omitempty"`
	Warranty                string        `json:"warranty,omitempty"`
	BinPickingNumber        string        `json:"bin_picking_number,omitempty"`
	LayoutFile              string        `json:"layout_file,omitempty"`
	Upc                     string        `json:"upc,omitempty"`
	Mpn                     string        `json:"mpn,omitempty"`
	Gtin                    string        `json:"gtin,omitempty"`
	SearchKeywords          string        `json:"search_keywords,omitempty"`
	Availability            string        `json:"availability,omitempty"`
	AvailabilityDescription string        `json:"availability_description,omitempty"`
	GiftWrappingOptionsType string        `json:"gift_wrapping_options_type,omitempty"`
	GiftWrappingOptionsList []interface{} `json:"gift_wrapping_options_list,omitempty"`
	SortOrder               int           `json:"sort_order,omitempty"`
	Condition               string        `json:"condition,omitempty"`
	IsConditionShown        bool          `json:"is_condition_shown,omitempty"`
	OrderQuantityMinimum    int           `json:"order_quantity_minimum,omitempty"`
	OrderQuantityMaximum    int           `json:"order_quantity_maximum,omitempty"`
	PageTitle               string        `json:"page_title,omitempty"`
	MetaKeywords            []interface{} `json:"meta_keywords,omitempty"`
	MetaDescription         string        `json:"meta_description,omitempty"`
	DateCreated             time.Time     `json:"date_created,omitempty"`
	DateModified            time.Time     `json:"date_modified,omitempty"`
	ViewCount               int           `json:"view_count,omitempty"`
	PreorderReleaseDate     interface{}   `json:"preorder_release_date,omitempty"`
	PreorderMessage         string        `json:"preorder_message,omitempty"`
	IsPreorderOnly          bool          `json:"is_preorder_only,omitempty"`
	IsPriceHidden           bool          `json:"is_price_hidden,omitempty"`
	PriceHiddenLabel        string        `json:"price_hidden_label,omitempty"`
	CustomURL               struct {
		URL          string `json:"url,omitempty"`
		IsCustomized bool   `json:"is_customized,omitempty"`
	} `json:"custom_url,omitempty"`
	BaseVariantID               int64  `json:"base_variant_id,omitempty"`
	OpenGraphType               string `json:"open_graph_type,omitempty"`
	OpenGraphTitle              string `json:"open_graph_title,omitempty"`
	OpenGraphDescription        string `json:"open_graph_description,omitempty"`
	OpenGraphUseMetaDescription bool   `json:"open_graph_use_meta_description,omitempty"`
	OpenGraphUseProductName     bool   `json:"open_graph_use_product_name,omitempty"`
	OpenGraphUseImage           bool   `json:"open_graph_use_image,omitempty"`
	Variants                    []struct {
		ID                        int64         `json:"id,omitempty"`
		ProductID                 int64         `json:"product_id,omitempty"`
		Sku                       string        `json:"sku,omitempty"`
		SkuID                     interface{}   `json:"sku_id,omitempty"`
		Price                     float64       `json:"price,omitempty"`
		CalculatedPrice           float64       `json:"calculated_price,omitempty"`
		SalePrice                 float64       `json:"sale_price,omitempty"`
		RetailPrice               float64       `json:"retail_price,omitempty"`
		MapPrice                  float64       `json:"map_price,omitempty"`
		Weight                    float64       `json:"weight,omitempty"`
		Width                     int           `json:"width,omitempty"`
		Height                    int           `json:"height,omitempty"`
		Depth                     int           `json:"depth,omitempty"`
		IsFreeShipping            bool          `json:"is_free_shipping,omitempty"`
		FixedCostShippingPrice    float64       `json:"fixed_cost_shipping_price,omitempty"`
		CalculatedWeight          float64       `json:"calculated_weight,omitempty"`
		PurchasingDisabled        bool          `json:"purchasing_disabled,omitempty"`
		PurchasingDisabledMessage string        `json:"purchasing_disabled_message,omitempty"`
		ImageURL                  string        `json:"image_url,omitempty"`
		CostPrice                 float64       `json:"cost_price,omitempty"`
		Upc                       string        `json:"upc,omitempty"`
		Mpn                       string        `json:"mpn,omitempty"`
		Gtin                      string        `json:"gtin,omitempty"`
		InventoryLevel            int           `json:"inventory_level,omitempty"`
		InventoryWarningLevel     int           `json:"inventory_warning_level,omitempty"`
		BinPickingNumber          string        `json:"bin_picking_number,omitempty"`
		OptionValues              []interface{} `json:"option_values,omitempty"`
	} `json:"variants,omitempty"`
	Images       []Image       `json:"images,omitempty"`
	PrimaryImage interface{}   `json:"primary_image,omitempty"`
	Videos       []interface{} `json:"videos,omitempty"`
	CustomFields []struct {
		ID    int64  `json:"id,omitempty"`
		Name  string `json:"name,omitempty"`
		Value string `json:"value,omitempty"`
	} `json:"custom_fields,omitempty"`
	BulkPricingRules []interface{} `json:"bulk_pricing_rules,omitempty"`
	Options          []interface{} `json:"options,omitempty"`
	Modifiers        []interface{} `json:"modifiers,omitempty"`
}

type ProductInventory struct {
	ID                      int64  `json:"id,omitempty"`
	InventoryLevel          int    `json:"inventory_level"`
	InventoryWarningLevel   int    `json:"inventory_warning_level,omitempty"`
	InventoryTracking       string `json:"inventory_tracking"`
	IsVisible               bool   `json:"is_visible"`
	Availability            string `json:"availability,omitempty"`
	AvailabilityDescription string `json:"availability_description,omitempty"`
}

// create struct to update only sale prices
type ProductSalePrice struct {
	ID        int64   `json:"id,omitempty"`
	SalePrice float64 `json:"sale_price"`
}

type Variant struct {
	ID                        int64         `json:"id,omitempty"`
	ProductID                 int64         `json:"product_id,omitempty"`
	Sku                       string        `json:"sku,omitempty"`
	SkuID                     interface{}   `json:"sku_id,omitempty"`
	Price                     float64       `json:"price,omitempty"`
	CalculatedPrice           float64       `json:"calculated_price,omitempty"`
	SalePrice                 float64       `json:"sale_price"`
	RetailPrice               float64       `json:"retail_price,omitempty"`
	MapPrice                  float64       `json:"map_price,omitempty"`
	Weight                    float64       `json:"weight,omitempty"`
	Width                     int           `json:"width,omitempty"`
	Height                    int           `json:"height,omitempty"`
	Depth                     int           `json:"depth,omitempty"`
	IsFreeShipping            bool          `json:"is_free_shipping,omitempty"`
	FixedCostShippingPrice    float64       `json:"fixed_cost_shipping_price,omitempty"`
	CalculatedWeight          float64       `json:"calculated_weight,omitempty"`
	PurchasingDisabled        bool          `json:"purchasing_disabled,omitempty"`
	PurchasingDisabledMessage string        `json:"purchasing_disabled_message,omitempty"`
	ImageURL                  string        `json:"image_url,omitempty"`
	CostPrice                 float64       `json:"cost_price,omitempty"`
	Upc                       string        `json:"upc,omitempty"`
	Mpn                       string        `json:"mpn,omitempty"`
	Gtin                      string        `json:"gtin,omitempty"`
	InventoryLevel            int           `json:"inventory_level"`
	InventoryWarningLevel     int           `json:"inventory_warning_level,omitempty"`
	BinPickingNumber          string        `json:"bin_picking_number,omitempty"`
	OptionValues              []interface{} `json:"option_values,omitempty"`
}

type VariantInventory struct {
	ID                    int64 `json:"id,omitempty"`
	ProductID             int64 `json:"product_id,omitempty"`
	InventoryLevel        int   `json:"inventory_level"`
	InventoryWarningLevel int   `json:"inventory_warning_level,omitempty"`
}

type VariantSalePrice struct {
	ID        int64   `json:"id,omitempty"`
	ProductID int64   `json:"product_id,omitempty"`
	SalePrice float64 `json:"sale_price"`
}

// Metafield is a struct representing a BigCommerce product metafield
type Metafield struct {
	ID            int64     `json:"id,omitempty"`
	Key           string    `json:"key,omitempty"`
	Value         string    `json:"value,omitempty"`
	ResourceID    int64     `json:"resource_id,omitempty"`
	ResourceType  string    `json:"resource_type,omitempty"`
	Description   string    `json:"description,omitempty"`
	DateCreated   time.Time `json:"date_created,omitempty"`
	DateModified  time.Time `json:"date_modified,omitempty"`
	Namespace     string    `json:"namespace,omitempty"`
	PermissionSet string    `json:"permission_set,omitempty"`
}

type ChannelAssignment struct {
	ProductID int64 `json:"product_id,omitempty"`
	ChannelID int64 `json:"channel_id,omitempty"`
}

// GetAllProducts gets all products from BigCommerce
// args is a key-value map of additional arguments to pass to the API
func (bc *Client) GetAllProducts(args map[string]string) ([]Product, error) {
	ps := []Product{}
	var psp []Product
	page := 1
	more := true
	var err error
	retries := 0
	for more {
		psp, more, err = bc.GetProducts(args, page)
		// log.Printf("page %d entries %d", page, len(psp))
		if err != nil {
			retries++
			if retries > bc.MaxRetries {
				log.Println("Max retries reached")
				return ps, err
			}
			break
		}
		ps = append(ps, psp...)
		page++
	}
	return ps, err
}

// GetProducts gets a page of products from BigCommerce
// args is a key-value map of additional arguments to pass to the API
// page: the page number to download
func (bc *Client) GetProducts(args map[string]string, page int) ([]Product, bool, error) {
	fpart := ""
	for k, v := range args {
		fpart += "&" + k + "=" + v
	}
	url := "/v3/catalog/products?page=" + strconv.Itoa(page) + fpart
	// log.Printf("GET %s", url)

	req := bc.getAPIRequest(http.MethodGet, url, nil)
	res, err := bc.HTTPClient.Do(req)
	if err != nil {
		return nil, false, err
	}
	defer res.Body.Close()
	if res.StatusCode == http.StatusNoContent {
		return nil, false, ErrNoContent
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, false, err
	}
	var pp struct {
		Status int       `json:"status"`
		Title  string    `json:"title"`
		Data   []Product `json:"data"`
		Meta   struct {
			Pagination Pagination `json:"pagination"`
		} `json:"meta"`
	}
	err = json.Unmarshal(body, &pp)
	if err != nil {
		return nil, false, err
	}
	//	log.Printf("%d products (%+v)", len(pp.Data), pp.Meta.Pagination)

	if pp.Status != 0 {
		return nil, false, errors.New(pp.Title)
	}
	return pp.Data, pp.Meta.Pagination.CurrentPage < pp.Meta.Pagination.TotalPages, nil
}

func (bc *Client) GetVariants(args map[string]string, page int) ([]Variant, bool, error) {
	fpart := ""
	for k, v := range args {
		fpart += "&" + k + "=" + v
	}
	url := "/v3/catalog/variants?page=" + strconv.Itoa(page) + fpart

	req := bc.getAPIRequest(http.MethodGet, url, nil)
	res, err := bc.HTTPClient.Do(req)
	if err != nil {
		return nil, false, err
	}
	defer res.Body.Close()
	if res.StatusCode == http.StatusNoContent {
		return nil, false, ErrNoContent
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, false, err
	}
	var pp struct {
		Status int       `json:"status"`
		Title  string    `json:"title"`
		Data   []Variant `json:"data"`
		Meta   struct {
			Pagination Pagination `json:"pagination"`
		} `json:"meta"`
	}
	err = json.Unmarshal(body, &pp)
	if err != nil {
		return nil, false, err
	}
	//	log.Printf("%d products (%+v)", len(pp.Data), pp.Meta.Pagination)

	if pp.Status != 0 {
		return nil, false, errors.New(pp.Title)
	}
	return pp.Data, pp.Meta.Pagination.CurrentPage < pp.Meta.Pagination.TotalPages, nil
}

// GetProductByID gets a product from BigCommerce by ID
// productID: BigCommerce product ID to get
func (bc *Client) GetProductByID(productID int64) (*Product, error) {
	url := "/v3/catalog/products/" + strconv.FormatInt(productID, 10) + "?include=variants,images,custom_fields,bulk_pricing_rules,primary_image,modifiers,options,videos"
	req := bc.getAPIRequest(http.MethodGet, url, nil)
	res, err := bc.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := processBody(res)
	if err != nil {
		return nil, err
	}

	var productResponse struct {
		Data Product `json:"data"`
	}
	err = json.Unmarshal(body, &productResponse)
	if err != nil {
		return nil, err
	}
	return &productResponse.Data, nil
}

// GetProductMetafields gets metafields values for a product
// productID: BigCommerce product ID to get metafields for
func (bc *Client) GetProductMetafields(productID int64) (map[string]Metafield, error) {
	url := "/v3/catalog/products/" + strconv.FormatInt(productID, 10) + "/metafields"
	req := bc.getAPIRequest(http.MethodGet, url, nil)
	res, err := bc.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := processBody(res)
	if err != nil {
		return nil, err
	}

	var metafieldsResponse struct {
		Metafields []Metafield `json:"data,omitempty"`
	}
	err = json.Unmarshal(body, &metafieldsResponse)
	if err != nil {
		return nil, err
	}
	ret := map[string]Metafield{}
	for _, mf := range metafieldsResponse.Metafields {
		ret[mf.Key] = mf
	}
	return ret, nil
}

func (bc *Client) CreateProduct(payload *Product) (*Product, error) {
	var b []byte
	prod := &payload
	b, _ = json.Marshal(prod)
	req := bc.getAPIRequest(http.MethodPost, "/v3/catalog/products", bytes.NewBuffer(b))
	res, err := bc.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := processBody(res)
	if err != nil {
		if res.StatusCode == http.StatusUnprocessableEntity {
			var errResp ErrorResult
			err = json.Unmarshal(body, &errResp)
			if err != nil {
				log.Printf("Error: %s\nResult: %s", err, string(body))
				return nil, err
			}
			if len(errResp.Errors) > 0 {
				errors := []string{}
				for _, e := range errResp.Errors {
					errors = append(errors, e)
				}
				return nil, fmt.Errorf("%s", strings.Join(errors, ", "))
			}
			return nil, errors.New("unknown error")
		}
		log.Printf("Error: %s\nResult: %s", err, string(body))
		return nil, err
	}
	var productResponse struct {
		Data Product `json:"data"`
	}
	err = json.Unmarshal(body, &productResponse)
	if err != nil {
		return nil, err
	}
	return &productResponse.Data, nil
}

// Update a product based on SKU and not on ID
func (bc *Client) UpdateProductBySku(payload *Product) (*Product, error) {
	var b []byte
	prod := &payload

	// retrieve product ID from BigCommerce
	searchCriteria := map[string]string{
		"sku": payload.Sku,
	}
	productsArray, err := bc.GetAllProducts(searchCriteria)
	// product doens't exists
	if err != nil {
		return nil, err
	}

	if len(productsArray) == 0 {
		return nil, errors.New("Empty response back on getting product by sku: " + payload.Sku)
	}
	prodId := strconv.Itoa(int(productsArray[0].ID))

	b, _ = json.Marshal(prod)

	req := bc.getAPIRequest(http.MethodPut, "/v3/catalog/products/"+prodId, bytes.NewBuffer(b))
	res, err := bc.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := processBody(res)
	if err != nil {
		if res.StatusCode == http.StatusUnprocessableEntity {
			var errResp ErrorResult
			err = json.Unmarshal(body, &errResp)
			if err != nil {
				log.Printf("Error: %s\nResult: %s", err, string(body))
				return nil, err
			}
			if len(errResp.Errors) > 0 {
				errors := []string{}
				for _, e := range errResp.Errors {
					errors = append(errors, e)
				}
				return nil, fmt.Errorf("%s", strings.Join(errors, ", "))
			}
			return nil, errors.New("unknown error")
		}
		log.Printf("Error: %s\nResult: %s", err, string(body))
		return nil, err
	}
	var productResponse struct {
		Data Product `json:"data"`
	}
	err = json.Unmarshal(body, &productResponse)
	if err != nil {
		return nil, err
	}
	return &productResponse.Data, nil
}

// Update a product inventory data based on ID
func (bc *Client) UpdateProductInventory(prodInventoryPayload *ProductInventory) (*Product, error) {
	var b []byte

	// make API request path
	path := fmt.Sprintf("/v3/catalog/products/%d", prodInventoryPayload.ID)

	b, _ = json.Marshal(prodInventoryPayload)
	// make the API request
	req := bc.getAPIRequest(http.MethodPut, path, bytes.NewBuffer(b))
	res, err := bc.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := processBody(res)
	if err != nil {
		if res.StatusCode == http.StatusUnprocessableEntity {
			var errResp ErrorResult
			err = json.Unmarshal(body, &errResp)
			if err != nil {
				log.Printf("Error: %s\nResult: %s", err, string(body))
				return nil, err
			}
			if len(errResp.Errors) > 0 {
				errors := []string{}
				for _, e := range errResp.Errors {
					errors = append(errors, e)
				}
				return nil, fmt.Errorf("%s", strings.Join(errors, ", "))
			}
			return nil, errors.New("unknown error")
		}
		log.Printf("Error: %s\nResult: %s", err, string(body))
		return nil, err
	}
	var productResponse struct {
		Data Product `json:"data"`
	}
	err = json.Unmarshal(body, &productResponse)
	if err != nil {
		return nil, err
	}
	return &productResponse.Data, nil
}

// Update a product sale price based on ID
func (bc *Client) UpdateProductSalePrice(prodSalePricePayload *ProductSalePrice) (*Product, error) {
	var b []byte

	// make API request path
	path := fmt.Sprintf("/v3/catalog/products/%d", prodSalePricePayload.ID)

	b, _ = json.Marshal(prodSalePricePayload)
	// make the API request
	req := bc.getAPIRequest(http.MethodPut, path, bytes.NewBuffer(b))
	res, err := bc.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := processBody(res)
	if err != nil {
		if res.StatusCode == http.StatusUnprocessableEntity {
			var errResp ErrorResult
			err = json.Unmarshal(body, &errResp)
			if err != nil {
				log.Printf("Error: %s\nResult: %s", err, string(body))
				return nil, err
			}
			if len(errResp.Errors) > 0 {
				errors := []string{}
				for _, e := range errResp.Errors {
					errors = append(errors, e)
				}
				return nil, fmt.Errorf("%s", strings.Join(errors, ", "))
			}
			return nil, errors.New("unknown error")
		}
		log.Printf("Error: %s\nResult: %s", err, string(body))
		return nil, err
	}
	var productResponse struct {
		Data Product `json:"data"`
	}
	err = json.Unmarshal(body, &productResponse)
	if err != nil {
		return nil, err
	}
	return &productResponse.Data, nil
}

// Update a product based on SKU and not on ID
func (bc *Client) UpdateVariantBySku(payload *Variant) (*Variant, error) {
	var b []byte
	prod := &payload

	// retrieve product ID from BigCommerce
	searchCriteria := map[string]string{
		"sku": payload.Sku,
	}
	variantsArray, _, err := bc.GetVariants(searchCriteria, 1)
	// product doens't exists
	if err != nil {
		return nil, err
	}

	if len(variantsArray) == 0 {
		return nil, errors.New("Empty response back on getting variant by sku: " + payload.Sku)
	}
	variantId := strconv.Itoa(int(variantsArray[0].ID))
	parentProductId := strconv.Itoa(int(variantsArray[0].ProductID))

	b, _ = json.Marshal(prod)
	log.Println("payload for variant:")
	log.Println(string(b))
	req := bc.getAPIRequest(http.MethodPut, "/v3/catalog/products/"+parentProductId+"/variants/"+variantId, bytes.NewBuffer(b))
	res, err := bc.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := processBody(res)
	if err != nil {
		if res.StatusCode == http.StatusUnprocessableEntity {
			var errResp ErrorResult
			err = json.Unmarshal(body, &errResp)
			if err != nil {
				log.Printf("Error: %s\nResult: %s", err, string(body))
				return nil, err
			}
			if len(errResp.Errors) > 0 {
				errors := []string{}
				for _, e := range errResp.Errors {
					errors = append(errors, e)
				}
				return nil, fmt.Errorf("%s", strings.Join(errors, ", "))
			}
			return nil, errors.New("unknown error")
		}
		log.Printf("Error: %s\nResult: %s", err, string(body))
		return nil, err
	}
	var variantResponse struct {
		Data Variant `json:"data"`
	}
	err = json.Unmarshal(body, &variantResponse)
	if err != nil {
		return nil, err
	}
	return &variantResponse.Data, nil
}

// Update a product inventory data based on ID
func (bc *Client) UpdateVariantInventory(variantPayload *VariantInventory) (*Variant, error) {
	var b []byte

	variantId := strconv.Itoa(int(variantPayload.ID))
	parentProductId := strconv.Itoa(int(variantPayload.ProductID))

	b, _ = json.Marshal(variantPayload)
	req := bc.getAPIRequest(http.MethodPut, "/v3/catalog/products/"+parentProductId+"/variants/"+variantId, bytes.NewBuffer(b))
	res, err := bc.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := processBody(res)
	if err != nil {
		if res.StatusCode == http.StatusUnprocessableEntity {
			var errResp ErrorResult
			err = json.Unmarshal(body, &errResp)
			if err != nil {
				log.Printf("Error: %s\nResult: %s", err, string(body))
				return nil, err
			}
			if len(errResp.Errors) > 0 {
				errors := []string{}
				for _, e := range errResp.Errors {
					errors = append(errors, e)
				}
				return nil, fmt.Errorf("%s", strings.Join(errors, ", "))
			}
			return nil, errors.New("unknown error")
		}
		log.Printf("Error: %s\nResult: %s", err, string(body))
		return nil, err
	}
	var variantResponse struct {
		Data Variant `json:"data"`
	}
	err = json.Unmarshal(body, &variantResponse)
	if err != nil {
		return nil, err
	}
	return &variantResponse.Data, nil
}

func (bc *Client) UpdateVariantSalePrice(variantSalePricePayload *VariantSalePrice) (*Variant, error) {
	var b []byte

	variantId := strconv.Itoa(int(variantSalePricePayload.ID))
	parentProductId := strconv.Itoa(int(variantSalePricePayload.ProductID))

	b, _ = json.Marshal(variantSalePricePayload)
	req := bc.getAPIRequest(http.MethodPut, "/v3/catalog/products/"+parentProductId+"/variants/"+variantId, bytes.NewBuffer(b))
	res, err := bc.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := processBody(res)
	if err != nil {
		if res.StatusCode == http.StatusUnprocessableEntity {
			var errResp ErrorResult
			err = json.Unmarshal(body, &errResp)
			if err != nil {
				log.Printf("Error: %s\nResult: %s", err, string(body))
				return nil, err
			}
			if len(errResp.Errors) > 0 {
				errors := []string{}
				for _, e := range errResp.Errors {
					errors = append(errors, e)
				}
				return nil, fmt.Errorf("%s", strings.Join(errors, ", "))
			}
			return nil, errors.New("unknown error")
		}
		log.Printf("Error: %s\nResult: %s", err, string(body))
		return nil, err
	}
	var variantResponse struct {
		Data Variant `json:"data"`
	}
	err = json.Unmarshal(body, &variantResponse)
	if err != nil {
		return nil, err
	}
	return &variantResponse.Data, nil
}

// deletes a product from a specific channel
func (bc *Client) DeleteProductFromChannel(productId int64, channelId int64) (bool, error) {
	req := bc.getAPIRequest(http.MethodDelete, fmt.Sprintf("/v3/catalog/products/channel-assignments?product_id:in=%d&channel_id:in=%d", productId, channelId), nil)
	res, err := bc.HTTPClient.Do(req)
	if err != nil {
		return false, err
	}
	defer res.Body.Close()

	if res.StatusCode > 200 && res.StatusCode < 300 {
		return true, nil
	}

	return false, nil
}

// add a product to a specific channel
func (bc *Client) AddProductToChannel(productId int64, channelId int64) (bool, error) {
	var payload []ChannelAssignment
	var channelAssignment ChannelAssignment
	channelAssignment.ProductID = productId
	channelAssignment.ChannelID = channelId

	payload = append(payload, channelAssignment)
	b, _ := json.Marshal(&payload)

	req := bc.getAPIRequest(http.MethodPut, "/v3/catalog/products/channel-assignments", bytes.NewBuffer(b))
	res, err := bc.HTTPClient.Do(req)
	if err != nil {
		return false, err
	}
	defer res.Body.Close()

	if res.StatusCode >= 200 && res.StatusCode < 300 {
		return true, nil
	}

	return false, nil
}
