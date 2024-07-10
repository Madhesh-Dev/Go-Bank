package api

import (
	mockdb "bank/db/mock"
	db "bank/db/sqlc"
	"bank/util"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestGetAccountApi(t *testing.T) {
	account := randomAccount()

	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)

	store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(account, nil)

	// start test server and send request

	server := NewServer(store)
	recoder := httptest.NewRecorder()

	url := fmt.Sprintf("/accounts/%d", account.ID)
	request, err := http.NewRequest(http.MethodGet, url, nil)

	require.NoError(t, err)

	server.router.ServeHTTP(recoder, request)

	require.Equal(t, http.StatusOK, recoder.Code)
	requireBodyMatchAccount(t, recoder.Body, account)
}

func randomAccount() db.Account {
	return db.Account{
		ID:       util.RandInt(1, 1000),
		Owner:    util.RandomOnwer(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
}

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account db.Account) {
	data, err := io.ReadAll(body)

	require.NoError(t, err)

	var gotAccount db.Account
	err = json.Unmarshal(data, &gotAccount)

	require.NoError(t, err)
	require.Equal(t, account, gotAccount)
}
