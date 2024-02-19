package service

import (
	"bot/internal/entities"
	"bot/internal/log"
	"bot/internal/repo"
	mockRepo "bot/internal/repo/mocks"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

// tests for user methods

// TestService_GetUser is a positive test for svc.GetUser
func TestService_GetUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mockRepo.NewMockUserRepo(ctrl)

	expected := &entities.User{
		ID:      123,
		Name:    "Cristiano Ronaldo",
		Phone:   "89251232323",
		Address: "madrid",
	}

	userRepo.EXPECT().GetUser(expected.ID).Return(expected, nil).Times(1)

	repos := repo.Repo{
		UserRepo: userRepo,
	}

	l := log.NewLogrus("debug")
	l.Named("service")

	svc := NewService(repos, l)

	got, err := svc.GetUser(expected.ID)

	require.NoError(t, err)
	require.Equal(t, expected, got)
}

// TestService_GetUser is a negative test for svc.GetUser
func TestService_GetUserError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mockRepo.NewMockUserRepo(ctrl)

	var userID int64 = 123

	repoErr := errors.New("random error in db")
	userRepo.EXPECT().GetUser(userID).Return(nil, repoErr).Times(1)

	repos := repo.Repo{
		UserRepo: userRepo,
	}

	l := log.NewLogrus("debug")
	l.Named("service")

	svc := NewService(repos, l)

	got, err := svc.GetUser(userID)

	require.Error(t, err)
	require.EqualError(t, err, repoErr.Error())
	require.Nil(t, got)
}

// tests for order methods

// TestService_GetAllCurrentOrders is a positive test for svc.GetAllCurrentOrders
func TestService_GetAllCurrentOrders(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	orderRepo := mockRepo.NewMockOrderRepo(ctrl)

	expected := []entities.CurrentOrder{
		{
			ID:     257,
			UserID: 123,
			Start:  time.Now(),
			Composition: []entities.Product{
				{
					ID:     1,
					Name:   "trousers",
					Size:   "XS",
					Color:  "black",
					Text:   "hello world",
					Img:    "",
					Amount: 2,
				},
				{
					ID:     2,
					Name:   "trousers",
					Size:   "S",
					Color:  "grey",
					Text:   "bye world",
					Img:    "",
					Amount: 1,
				},
			},
		},
		{
			ID:     258,
			UserID: 124,
			Start:  time.Now(),
			Composition: []entities.Product{
				{
					ID:     3,
					Name:   "hoodie",
					Size:   "L",
					Color:  "black",
					Text:   "AC Milan",
					Img:    "",
					Amount: 3,
				},
			},
		},
	}

	orderRepo.EXPECT().GetAllCurrentOrders().Return(expected, nil).Times(1)

	repos := repo.Repo{OrderRepo: orderRepo}

	l := log.NewLogrus("debug")
	l.Named("service")

	svc := NewService(repos, l)

	got, err := svc.GetAllCurrentOrders()

	require.NoError(t, err)
	require.Equal(t, expected, got)
}

// TestService_GetAllCurrentOrdersError is a negative test for svc.GetAllCurrentOrders
func TestService_GetAllCurrentOrdersError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	orderRepo := mockRepo.NewMockOrderRepo(ctrl)

	repoErr := errors.New("random error from db")
	orderRepo.EXPECT().GetAllCurrentOrders().Return(nil, repoErr).Times(1)

	repos := repo.Repo{OrderRepo: orderRepo}

	l := log.NewLogrus("debug")
	l.Named("service")

	svc := NewService(repos, l)

	got, err := svc.GetAllCurrentOrders()

	require.Error(t, err)
	require.EqualError(t, err, repoErr.Error())
	require.Nil(t, got)
}

func TestService_GetCartProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cartRepo := mockRepo.NewMockCartRepo(ctrl)

	var userID int64 = 11
	var idx int = 10

	expected := &entities.Product{
		ID:     15,
		Name:   "hoodie",
		Size:   "XS",
		Color:  "white",
		Text:   "",
		Img:    "link.jpg",
		Amount: 1,
	}

	cartRepo.EXPECT().GetCartProduct(userID, idx).Return(expected, nil).Times(1)

	repos := repo.Repo{CartRepo: cartRepo}

	l := log.NewLogrus("debug")
	l.Named("service")

	svc := NewService(repos, l)

	got, err := svc.GetCartProduct(userID, idx)

	require.NoError(t, err)
	require.Equal(t, expected, got)
}

func TestService_GetCartProductError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cartRepo := mockRepo.NewMockCartRepo(ctrl)

	var userID int64 = 11
	var idx int = 10

	repoErr := errors.New("random error from db")
	cartRepo.EXPECT().GetCartProduct(userID, idx).Return(nil, repoErr).Times(1)

	repos := repo.Repo{CartRepo: cartRepo}

	l := log.NewLogrus("debug")
	l.Named("service")

	svc := NewService(repos, l)

	got, err := svc.GetCartProduct(userID, idx)

	require.Error(t, err)
	require.EqualError(t, err, repoErr.Error())
	require.Nil(t, got)
}
