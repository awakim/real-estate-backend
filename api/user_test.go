package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	mockcache "github.com/awakim/immoblock-backend/cache/mock"
	mockdb "github.com/awakim/immoblock-backend/db/mock"
	db "github.com/awakim/immoblock-backend/db/sqlc"
	mockidentity "github.com/awakim/immoblock-backend/identity/mock"
	"github.com/awakim/immoblock-backend/util"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

type eqCreateUserParamsMatcher struct {
	arg      db.CreateUserParams
	password string
}

func (e eqCreateUserParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.CreateUserParams)
	if !ok {
		return false
	}

	err := util.CheckPassword(e.password, arg.HashedPassword)
	if err != nil {
		return false
	}

	e.arg.HashedPassword = arg.HashedPassword
	return reflect.DeepEqual(e.arg, arg)
}

func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

func EqCreateUserParams(arg db.CreateUserParams, password string) gomock.Matcher {
	return eqCreateUserParamsMatcher{arg, password}
}

func randomUser(t *testing.T) (user db.User, hashedPassword string) {
	password := util.RandomString(8)
	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)
	uid, err := uuid.NewRandom()
	require.NoError(t, err)

	user = db.User{
		ID:             uid,
		HashedPassword: hashedPassword,
		Nickname:       util.RandomString(6),
		Email:          util.RandomEmail(),
	}
	return
}

func requireBodyMatchUser(t *testing.T, body *bytes.Buffer, user db.User) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotUser db.User
	err = json.Unmarshal(data, &gotUser)

	require.NoError(t, err)
	require.Equal(t, user.Nickname, gotUser.Nickname)
	require.Equal(t, user.Email, gotUser.Email)
	require.Empty(t, gotUser.HashedPassword)
}

func TestCreateUserAPI(t *testing.T) {
	user, password := randomUser(t)

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore, cache *mockcache.MockCache, userManager *mockidentity.MockUserManagement)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"password": password,
				"nickname": user.Nickname,
				"email":    user.Email,
			},
			buildStubs: func(store *mockdb.MockStore, cache *mockcache.MockCache, userManager *mockidentity.MockUserManagement) {
				arg := db.CreateUserParams{
					Nickname: user.Nickname,
					Email:    user.Email,
				}
				store.EXPECT().CreateUser(gomock.Any(), EqCreateUserParams(arg, password)).Times(1).Return(user, nil)
				userManager.EXPECT().Create(gomock.Any(), gomock.Any()).Times(1).Return(nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchUser(t, recorder.Body, user)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"password": password,
				"nickname": user.Nickname,
				"email":    user.Email,
			},
			buildStubs: func(store *mockdb.MockStore, cache *mockcache.MockCache, userManager *mockidentity.MockUserManagement) {
				store.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(1).Return(db.User{}, sql.ErrConnDone)
				userManager.EXPECT().Create(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "DuplicateUsername",
			body: gin.H{
				"password": password,
				"nickname": user.Nickname,
				"email":    user.Email,
			},
			buildStubs: func(store *mockdb.MockStore, cache *mockcache.MockCache, userManager *mockidentity.MockUserManagement) {
				store.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(1).Return(db.User{}, &pq.Error{Code: "23505"})
				userManager.EXPECT().Create(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
		},
		{
			name: "InvalidEmail",
			body: gin.H{
				"password": password,
				"nickname": user.Nickname,
				"email":    "invalid-email",
			},
			buildStubs: func(store *mockdb.MockStore, cache *mockcache.MockCache, userManager *mockidentity.MockUserManagement) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
				userManager.EXPECT().Create(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "TooShortPassword",
			body: gin.H{
				"password": "123",
				"nickname": user.Nickname,
				"email":    user.Email,
			},
			buildStubs: func(store *mockdb.MockStore, cache *mockcache.MockCache, userManager *mockidentity.MockUserManagement) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
				userManager.EXPECT().Create(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
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

			// Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/users"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.Router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}
