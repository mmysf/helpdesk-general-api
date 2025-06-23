package xenditrepo

import (
	"app/domain"
	"app/domain/model"
	"app/helpers"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/Yureka-Teknologi-Cipta/yureka/response"
	"github.com/sirupsen/logrus"
)

func (r *xenditRepo) GenereteSnapLink(ctx context.Context, pkg *model.HourPackage, pkgS *model.ServerPackage, order model.Order) (result response.Base, err error) {
	items := []map[string]interface{}{}

	if order.Type == model.HOUR_TYPE {
		if pkg == nil {
			result.Status = 400
			result.Message = "Package not found"
			return
		}
		items = append(items, map[string]interface{}{
			"name":     pkg.Name,
			"quantity": order.Amount,
			"price":    pkg.Price,
		})
	} else {
		if pkgS == nil {
			result.Status = 400
			result.Message = "Package not found"
			return
		}
		items = append(items, map[string]interface{}{
			"name":     pkgS.Name,
			"quantity": order.Amount,
			"price":    pkgS.Price,
		})
	}

	generateSnapLinkDataApi := struct {
		ExternalId      string                   `json:"external_id"`
		Amount          float64                  `json:"amount"`
		PayerEmail      string                   `json:"payer_email"`
		Currency        string                   `json:"currency"`
		Locale          string                   `json:"locale"`
		Description     string                   `json:"description"`
		InvoiceDuration int64                    `json:"invoice_duration"`
		Customer        map[string]interface{}   `json:"customer"`
		Fees            []map[string]interface{} `json:"fees"`
		Items           []map[string]interface{} `json:"items"`
		Metadata        map[string]interface{}   `json:"metadata"`
	}{
		ExternalId:      order.ID.Hex(),
		Amount:          order.GrandTotal,
		PayerEmail:      order.Customer.Email,
		Description:     r.metadataIssuer,
		InvoiceDuration: r.xenditInvoiceDuration,
		Currency:        "IDR",
		Locale:          "en",
		Items:           items,
		Customer: map[string]interface{}{
			"given_names": order.Customer.Name,
			"email":       order.Customer.Email,
			"phone":       "-",
		},
		Fees: []map[string]interface{}{
			{
				"type":  "admin fee",
				"value": order.AdminFee,
			},
			{
				"type":  "tax",
				"value": order.Tax,
			},
			{
				"type":  "discount",
				"value": order.Discount * -1,
			},
		},
		Metadata: map[string]interface{}{
			"issuer": r.metadataIssuer,
		},
	}

	jsonByte, err := json.Marshal(generateSnapLinkDataApi)

	payload := strings.NewReader(string(jsonByte))

	req, err := http.NewRequest("POST", r.generateSnapURL, payload)

	helpers.Dump(r.generateSnapURL)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Basic "+r.secretBasicAuth)

	res, err := r.Client.Do(req)

	if err != nil {
		logrus.Error("Generete Snap Link", err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if res.StatusCode != 200 {
		failed := domain.XenditResponseError{}
		err = json.Unmarshal(body, &failed)
		if err != nil {
			logrus.Error("Generete Snap Link Response Unmarshal", err)
			result.Status = 400
			result.Message = "Generete Snap Link Response Unmarshal"
			return
		}
		result.Status = res.StatusCode
		result.Message = failed.Message
		return
	}

	if err != nil {
		logrus.Error("Generete Snap Link Response ReadBody", err)
		result.Status = 400
		result.Message = err.Error()
		return
	}

	success := domain.XenditGenereteSnapLinkResponseSuccess{}
	err = json.Unmarshal(body, &success)
	if err != nil {
		logrus.Error("Generete Snap Link Response Unmarshal", err)
		result.Status = 400
		result.Message = err.Error()
		return
	}

	result.Status = res.StatusCode
	result.Message = "success"
	result.Data = success

	return
}
