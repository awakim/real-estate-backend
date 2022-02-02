package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	mockcache "github.com/awakim/immoblock-backend/cache/mock"
	mockdb "github.com/awakim/immoblock-backend/db/mock"
	db "github.com/awakim/immoblock-backend/db/sqlc"
	"github.com/awakim/immoblock-backend/token"
	"github.com/awakim/immoblock-backend/util"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func randomProperty(t *testing.T) (property db.Property) {
	property = db.Property{
		ID:                  util.RandomPropertyID(),
		Name:                util.RandomString(6),
		Description:         util.RandomString(32),
		InitialBlockCount:   1000,
		RemainingBlockCount: 1000,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}
	return
}

func TestTransferAPI(t *testing.T) {
	amount := int64(10)

	user1, _ := randomUser(t)
	user2, _ := randomUser(t)
	user3, _ := randomUser(t)

	account1 := randomAccount(user1.ID)
	account2 := randomAccount(user2.ID)
	account3 := randomAccount(user3.ID)

	property1 := randomProperty(t)
	property2 := randomProperty(t)
	account1.PropertyID = property1.ID
	account2.PropertyID = property1.ID
	account3.PropertyID = property2.ID

	testCases := []struct {
		name          string
		body          gin.H
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore, cache *mockcache.MockCache)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"amount":          amount,
				"property_id":     property1.ID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user1.ID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore, cache *mockcache.MockCache) {
				cache.EXPECT().IsRevoked(gomock.Any(), gomock.Any()).Times(1).Return(false, nil)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account1.ID)).Times(1).Return(account1, nil)
				store.EXPECT().GetProperty(gomock.Any(), gomock.Eq(account1.PropertyID)).Times(1).Return(property1, nil)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account2.ID)).Times(1).Return(account2, nil)
				store.EXPECT().GetProperty(gomock.Any(), gomock.Eq(account2.PropertyID)).Times(1).Return(property1, nil)

				arg := db.TransferTxParams{
					FromAccountID: account1.ID,
					ToAccountID:   account2.ID,
					Amount:        amount,
				}
				store.EXPECT().TransferTx(gomock.Any(), gomock.Eq(arg)).Times(1)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "UnauthorizedUser",
			body: gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"amount":          amount,
				"property_id":     property1.ID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user2.ID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore, cache *mockcache.MockCache) {
				cache.EXPECT().IsRevoked(gomock.Any(), gomock.Any()).Times(1).Return(false, nil)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account1.ID)).Times(1).Return(account1, nil)
				store.EXPECT().GetProperty(gomock.Any(), gomock.Eq(account1.PropertyID)).Times(1).Return(property1, nil)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().GetProperty(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "NoAuthorization",
			body: gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"amount":          amount,
				"property_id":     property1.ID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore, cache *mockcache.MockCache) {
				cache.EXPECT().IsRevoked(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().GetProperty(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().GetProperty(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "FromAccountNotFound",
			body: gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"amount":          amount,
				"property_id":     property1.ID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user1.ID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore, cache *mockcache.MockCache) {
				cache.EXPECT().IsRevoked(gomock.Any(), gomock.Any()).Times(1).Return(false, nil)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account1.ID)).Times(1).Return(db.Account{}, sql.ErrNoRows)
				store.EXPECT().GetProperty(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().GetProperty(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "ToAccountNotFound",
			body: gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"amount":          amount,
				"property_id":     property1.ID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user1.ID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore, cache *mockcache.MockCache) {
				cache.EXPECT().IsRevoked(gomock.Any(), gomock.Any()).Times(1).Return(false, nil)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account1.ID)).Times(1).Return(account1, nil)
				store.EXPECT().GetProperty(gomock.Any(), gomock.Eq(account1.PropertyID)).Times(1).Return(property1, nil)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account2.ID)).Times(1).Return(db.Account{}, sql.ErrNoRows)
				store.EXPECT().GetProperty(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "FromAccountPropertyIDMismatch",
			body: gin.H{
				"from_account_id": account3.ID,
				"to_account_id":   account2.ID,
				"amount":          amount,
				"property_id":     property1.ID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user3.ID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore, cache *mockcache.MockCache) {
				cache.EXPECT().IsRevoked(gomock.Any(), gomock.Any()).Times(1).Return(false, nil)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account3.ID)).Times(1).Return(account3, nil)
				store.EXPECT().GetProperty(gomock.Any(), gomock.Eq(account1.PropertyID)).Times(1)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account2.ID)).Times(0)
				store.EXPECT().GetProperty(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		// {
		// 	name: "ToAccountCurrencyMismatch",
		// 	body: gin.H{
		// 		"from_account_id": account1.ID,
		// 		"to_account_id":   account3.ID,
		// 		"amount":          amount,
		// 		"property_id":     propertyID1,
		// 	},
		// 	setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
		// 		addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)
		// 	},
		// 	buildStubs: func(store *mockdb.MockStore, cache *mockcache.MockCache) {
		// 		store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account1.ID)).Times(1).Return(account1, nil)
		// 		store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account3.ID)).Times(1).Return(account3, nil)
		// 		store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(0)
		// 	},
		// 	checkResponse: func(recorder *httptest.ResponseRecorder) {
		// 		require.Equal(t, http.StatusBadRequest, recorder.Code)
		// 	},
		// },
		// {
		// 	name: "InvalidPropertyID",
		// 	body: gin.H{
		// 		"from_account_id": account1.ID,
		// 		"to_account_id":   account2.ID,
		// 		"amount":          amount,
		// 		"property_id":     "XYZ",
		// 	},
		// 	setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
		// 		addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)
		// 	},
		// 	buildStubs: func(store *mockdb.MockStore, cache *mockcache.MockCache) {
		// 		store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)
		// 		store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(0)
		// 	},
		// 	checkResponse: func(recorder *httptest.ResponseRecorder) {
		// 		require.Equal(t, http.StatusBadRequest, recorder.Code)
		// 	},
		// },
		// {
		// 	name: "NegativeAmount",
		// 	body: gin.H{
		// 		"from_account_id": account1.ID,
		// 		"to_account_id":   account2.ID,
		// 		"amount":          -amount,
		// 		"property_id":     propertyID1,
		// 	},
		// 	setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
		// 		addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)
		// 	},
		// 	buildStubs: func(store *mockdb.MockStore, cache *mockcache.MockCache) {
		// 		store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)
		// 		store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(0)
		// 	},
		// 	checkResponse: func(recorder *httptest.ResponseRecorder) {
		// 		require.Equal(t, http.StatusBadRequest, recorder.Code)
		// 	},
		// },
		// {
		// 	name: "GetAccountError",
		// 	body: gin.H{
		// 		"from_account_id": account1.ID,
		// 		"to_account_id":   account2.ID,
		// 		"amount":          amount,
		// 		"property_id":     propertyID1,
		// 	},
		// 	setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
		// 		addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)
		// 	},
		// 	buildStubs: func(store *mockdb.MockStore, cache *mockcache.MockCache) {
		// 		store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(1).Return(db.Account{}, sql.ErrConnDone)
		// 		store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(0)
		// 	},
		// 	checkResponse: func(recorder *httptest.ResponseRecorder) {
		// 		require.Equal(t, http.StatusInternalServerError, recorder.Code)
		// 	},
		// },
		// {
		// 	name: "TransferTxError",
		// 	body: gin.H{
		// 		"from_account_id": account1.ID,
		// 		"to_account_id":   account2.ID,
		// 		"amount":          amount,
		// 		"property_id":     propertyID1,
		// 	},
		// 	setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
		// 		addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)
		// 	},
		// 	buildStubs: func(store *mockdb.MockStore, cache *mockcache.MockCache) {
		// 		store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account1.ID)).Times(1).Return(account1, nil)
		// 		store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account2.ID)).Times(1).Return(account2, nil)
		// 		store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(1).Return(db.TransferTxResult{}, sql.ErrTxDone)
		// 	},
		// 	checkResponse: func(recorder *httptest.ResponseRecorder) {
		// 		require.Equal(t, http.StatusInternalServerError, recorder.Code)
		// 	},
		// },
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			cache := mockcache.NewMockCache(ctrl)
			tc.buildStubs(store, cache)

			server := newTestServer(t, store, cache)
			recorder := httptest.NewRecorder()

			// Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/transfers"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, request, server.TokenMaker)
			server.Router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}
