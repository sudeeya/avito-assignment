package integration

import (
	"bytes"
	"encoding/json"
	"io"
	"math/rand/v2"
	"net/http"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/sudeeya/avito-assignment/internal/model"
)

type IntegrationSuite struct {
	suite.Suite

	url    string
	bearer string
	client *http.Client
}

func (s *IntegrationSuite) SetupSuite() {
	s.url = "http://localhost:8080/api/v1"
	s.client = &http.Client{}

	req, err := http.NewRequest(http.MethodPost, s.url+"/dummyLogin", bytes.NewReader(
		[]byte(`{"role": "moderator"}`),
	))
	s.Require().NoError(err, "Failed to create request")

	resp, err := s.client.Do(req)
	s.Require().NoError(err, "Failed to do request")

	token, err := io.ReadAll(resp.Body)
	s.Require().NoError(err, "Failed to read token")

	s.bearer = "Bearer " + string(token)
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, &IntegrationSuite{})
}

func (s *IntegrationSuite) TestCreateAndCloseReception() {
	// Create PVZ
	req, err := http.NewRequest(http.MethodPost, s.url+"/pvz", bytes.NewReader(
		[]byte(`{"city":"Москва"}`),
	))
	s.Require().NoError(err, "Failed to create request")

	s.addToken(req)

	resp, err := s.client.Do(req)
	s.Require().NoError(err, "Failed to do request")

	var pvz model.PVZ
	err = json.NewDecoder(resp.Body).Decode(&pvz)
	s.Require().NoError(err, "Failed to read PVZ")

	s.Require().NotEmpty(pvz, "PVZ is empty")
	s.Require().Equal("Москва", pvz.City, "Another city was returned")

	// Create reception
	req, err = http.NewRequest(http.MethodPost, s.url+"/receptions", bytes.NewReader(
		[]byte(`{"pvz_id":"`+pvz.ID.String()+`"}`),
	))
	s.Require().NoError(err, "Failed to create request")

	s.addToken(req)

	resp, err = s.client.Do(req)
	s.Require().NoError(err, "Failed to do request")

	var reception model.Reception
	err = json.NewDecoder(resp.Body).Decode(&reception)
	s.Require().NoError(err, "Failed to read reception")

	s.Require().NotEmpty(reception, "Reception is empty")
	s.Require().NotEmpty(reception, "Reception is empty")
	s.Require().Equal("in_progress", reception.Status, "Not in progress reception was returned")
	s.Require().Equal(pvz.ID, reception.PVZID, "Not in progress reception was returned")

	// Add products
	productTypes := []string{"обувь", "одежда", "электроника"}
	for range 50 {
		productType := productTypes[rand.IntN(len(productTypes))]
		req, err = http.NewRequest(http.MethodPost, s.url+"/products", bytes.NewReader(
			[]byte(`{"pvz_id":"`+pvz.ID.String()+`", "type":"`+productType+`"}`),
		))
		s.Require().NoError(err, "Failed to create request")

		s.addToken(req)

		resp, err = s.client.Do(req)
		s.Require().NoError(err, "Failed to do request")

		var product model.Product
		err = json.NewDecoder(resp.Body).Decode(&product)
		s.Require().NoError(err, "Failed to read product")

		s.Require().NotEmpty(reception, "Reception is empty")
		s.Require().Equal(productType, product.Type, "Another product type was returned")
		s.Require().Equal(reception.ID, product.ReceptionID, "Another reception ID was returned")
	}

	// Close reception
	req, err = http.NewRequest(http.MethodPost, s.url+"/pvz/"+pvz.ID.String()+"/close_last_reception", nil)
	s.Require().NoError(err, "Failed to create request")

	s.addToken(req)

	resp, err = s.client.Do(req)
	s.Require().NoError(err, "Failed to do request")

	var receptionOnClose model.Reception
	err = json.NewDecoder(resp.Body).Decode(&receptionOnClose)
	s.Require().NoError(err, "Failed to read reception")

	s.Require().NotEmpty(receptionOnClose, "Reception is empty")
	s.Require().Equal(reception.ID, receptionOnClose.ID, "Another reception ID was returned")
	s.Require().Equal(reception.PVZID, receptionOnClose.PVZID, "Another PVZ ID was returned")
}

func (s *IntegrationSuite) addToken(req *http.Request) {
	req.Header.Set("Authorization", s.bearer)
}
