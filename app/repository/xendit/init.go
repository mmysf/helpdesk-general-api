package xenditrepo

import (
	"app/domain/model"
	"context"
	"encoding/base64"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/Yureka-Teknologi-Cipta/yureka/response"
)

type xenditRepo struct {
	Client                *http.Client
	baseURL               *url.URL
	generateQrisURL       string
	generateSnapURL       string
	secret                string
	secretBasicAuth       string
	metadataIssuer        string
	xenditInvoiceDuration int64
}

func NewXenditRepo() XenditRepo {
	client := &http.Client{}
	baseURL, _ := url.Parse(os.Getenv("XENDIT_BASE_URL"))
	secret := os.Getenv("XENDIT_API_KEY")
	metadataIssuer := os.Getenv("XENDIT_METADATA_ISSUER")
	xenditInvoiceDuration, _ := strconv.Atoi(os.Getenv("XENDIT_INVOICE_DURATION"))

	return &xenditRepo{
		Client:                client,
		baseURL:               baseURL,
		secret:                secret,
		metadataIssuer:        metadataIssuer,
		secretBasicAuth:       base64.StdEncoding.EncodeToString([]byte(secret + ":")),
		generateQrisURL:       baseURL.ResolveReference(&url.URL{Path: "/qr_codes"}).String(),
		generateSnapURL:       baseURL.ResolveReference(&url.URL{Path: "/v2/invoices"}).String(),
		xenditInvoiceDuration: int64(xenditInvoiceDuration),
	}
}

type XenditRepo interface {
	GenereteSnapLink(ctx context.Context, pkg *model.HourPackage, pkgS *model.ServerPackage, order model.Order) (response.Base, error)
}
