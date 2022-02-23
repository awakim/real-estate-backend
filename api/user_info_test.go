package api

import (
	"bytes"
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	mockcache "github.com/awakim/immoblock-backend/cache/mock"
	mockdb "github.com/awakim/immoblock-backend/db/mock"
	db "github.com/awakim/immoblock-backend/db/sqlc"
	mockidentity "github.com/awakim/immoblock-backend/identity/mock"
	"github.com/awakim/immoblock-backend/token"
	"github.com/awakim/immoblock-backend/util"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func randomUserInfo(userID uuid.UUID) db.UserInformation {
	return db.UserInformation{
		UserID:           userID,
		Firstname:        util.RandomString(6),
		Lastname:         util.RandomString(6),
		PhoneNumber:      util.RandomString(6),
		Nationality:      util.RandomString(6),
		Address:          util.RandomString(32),
		PostalCode:       "75005",
		City:             "Paris",
		Country:          "France",
		VerificationStep: 1,
	}
}

func TestGetUserInfoAPI(t *testing.T) {
	user1, _ := randomUser(t)
	user2, _ := randomUser(t)

	userInfo1 := randomUserInfo(user1.ID)
	userInfo2 := db.UserInformation{}

	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore, cache *mockcache.MockCache, userManager *mockidentity.MockUserManagement)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",

			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user1.ID, user1.IsAdmin, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore, cache *mockcache.MockCache, userManager *mockidentity.MockUserManagement) {
				cache.EXPECT().IsRevoked(gomock.Any(), gomock.Any()).Times(1).Return(false, nil)
				store.EXPECT().GetUserInfo(gomock.Any(), gomock.Eq(user1.ID)).Times(1).Return(userInfo1, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "NotFound",

			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user2.ID, user2.IsAdmin, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore, cache *mockcache.MockCache, userManager *mockidentity.MockUserManagement) {
				cache.EXPECT().IsRevoked(gomock.Any(), gomock.Any()).Times(1).Return(false, nil)
				store.EXPECT().GetUserInfo(gomock.Any(), gomock.Eq(user2.ID)).Times(1).Return(userInfo2, sql.ErrNoRows)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "InternalError",

			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user2.ID, user2.IsAdmin, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore, cache *mockcache.MockCache, userManager *mockidentity.MockUserManagement) {
				cache.EXPECT().IsRevoked(gomock.Any(), gomock.Any()).Times(1).Return(false, nil)
				store.EXPECT().GetUserInfo(gomock.Any(), gomock.Eq(user2.ID)).Times(1).Return(userInfo2, errors.New("internal server error"))
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "Unauthorized",

			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				// addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user2.ID, user2.IsAdmin, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore, cache *mockcache.MockCache, userManager *mockidentity.MockUserManagement) {
				cache.EXPECT().IsRevoked(gomock.Any(), gomock.Any()).Times(0).Return(false, nil)
				// store.EXPECT().GetUserInfo(gomock.Any(), gomock.Eq(user2.ID)).Times(1).Return(userInfo2, errors.New("internal server error"))
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "UserRevoked",

			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user2.ID, user2.IsAdmin, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore, cache *mockcache.MockCache, userManager *mockidentity.MockUserManagement) {
				cache.EXPECT().IsRevoked(gomock.Any(), gomock.Any()).Times(1).Return(true, nil)
				// store.EXPECT().GetUserInfo(gomock.Any(), gomock.Eq(user2.ID)).Times(1).Return(userInfo2, errors.New("internal server error"))
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			cache := mockcache.NewMockCache(ctrl)
			userManager := mockidentity.NewMockUserManagement(ctrl)
			tc.buildStubs(store, cache, userManager)

			server := newTestServer(t, store, cache, userManager)
			recorder := httptest.NewRecorder()

			url := "/users/info"
			request, err := http.NewRequest(http.MethodGet, url, bytes.NewReader(nil))
			require.NoError(t, err)

			tc.setupAuth(t, request, server.TokenMaker)
			server.Router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

// func TestCreateUserInfoAPI(t *testing.T) {
// 	user1, _ := randomUser(t)
// 	// user2, _ := randomUser(t)
// 	body1 := gin.H{
// 		"firstname":    util.RandomString(6),
// 		"lastname":     util.RandomString(6),
// 		"phone_number": util.RandomPhoneNumber(),
// 		"nationality":  "France",
// 		"address":      util.RandomString(32),
// 		"postal_code":  "75009",
// 		"city":         "Paris",
// 		"country":      "France",
// 	}
// 	userInfo1 := randomUserInfo(user1.ID)
// 	arg1 := db.CreateUserInfoParams{
// 		UserID:           user1.ID,
// 		Firstname:        body1["firstname"].(string),
// 		Lastname:         body1["firstname"].(string),
// 		PhoneNumber:      body1["phone_number"].(string),
// 		Nationality:      body1["nationality"].(string),
// 		Address:          body1[""],
// 		PostalCode:       "",
// 		City:             "",
// 		Country:          "",
// 		VerificationStep: 0,
// 	}
// 	// userInfo2 := db.UserInformation{}

// 	testCases := []struct {
// 		name          string
// 		body          gin.H
// 		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
// 		buildStubs    func(store *mockdb.MockStore, cache *mockcache.MockCache, userManager *mockidentity.MockUserManagement)
// 		checkResponse func(recoder *httptest.ResponseRecorder)
// 	}{
// 		{
// 			name: "OK",
// 			body: body,

// 			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
// 				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user1.ID, user1.IsAdmin, time.Minute)
// 			},
// 			buildStubs: func(store *mockdb.MockStore, cache *mockcache.MockCache, userManager *mockidentity.MockUserManagement) {
// 				cache.EXPECT().IsRevoked(gomock.Any(), gomock.Any()).Times(1).Return(false, nil)
// 				store.EXPECT().ExistsUserInfo(gomock.Any(), gomock.Eq(user1.ID)).Times(1).Return(false, nil)
// 				store.EXPECT().CreateUserInfo(gomock.Any(), gomock.Eq(arg1)).Times(1).Return(userInfo1, nil)
// 			},
// 			checkResponse: func(recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusOK, recorder.Code)
// 			},
// 		},
// 	}

// 	for i := range testCases {
// 		tc := testCases[i]

// 		t.Run(tc.name, func(t *testing.T) {
// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()

// 			store := mockdb.NewMockStore(ctrl)
// 			cache := mockcache.NewMockCache(ctrl)
// 			tc.buildStubs(store, cache, userManager)

// 			server := newTestServer(t, store, cache)
// 			recorder := httptest.NewRecorder()

// 			// Marshal body data to JSON
// 			data, err := json.Marshal(tc.body)
// 			require.NoError(t, err)

// 			url := "/users/info"
// 			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
// 			require.NoError(t, err)

// 			tc.setupAuth(t, request, server.TokenMaker)
// 			server.Router.ServeHTTP(recorder, request)
// 			tc.checkResponse(recorder)
// 		})
// 	}
// }
