package user

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"go.uber.org/zap"

	// "mado/internal"
	"mado/internal/auth"
	"mado/internal/auth/model"
	"mado/models"
)

// Repository is a user repository.
type Repository interface {
	Create(ctx context.Context, dto *User) (*User, error)
	CheckIfUserExistsByIIN(ctx context.Context, iin string) (int, bool, error)
	GetUsersSignature(ctx context.Context, userId int) (string, error)
	GetUser(ctx context.Context, userId int) (models.User, error)
}

// Service is a user service interface.
type Service struct {
	userRepository Repository
	logger         *zap.Logger
}

const (
	block   = "Блок данных на подпись"
	baseURL = "https://sigex.kz"
)

// NewService creates a new user service.
func NewService(userRepository Repository, logger *zap.Logger) Service {
	return Service{
		userRepository: userRepository,
		logger:         logger,
	}
}

func (s Service) GetUsersSignature(ctx context.Context, userId int) (string, error) {
	return s.userRepository.GetUsersSignature(ctx, userId)
}

func (s Service) Login(requirements model.LoginRequirements) (*User, error) {
	// ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	// defer cancel()

	ctx := context.Background()

	signature := auth.GetNonceSignature(requirements.QrSigner)
	fmt.Println("LEN OF sign", len(*signature))
	fmt.Println("SIGNATURE111: ", signature)
	req := model.AuthRequest{
		Nonce:     requirements.Nonce,
		Signature: signature,
		External:  true,
	}

	response, err := authentification(req)
	if err != nil {
		fmt.Println("Authentication error:", err)
	}
	fmt.Println(response)
	iin := (response.UserID)[3:]
	user := &User{Username: getName(response.Subject), IIN: &iin, Email: &response.Email, BIN: &response.BusinessID, Is_manager: requirements.Is_manager, Signature: signature}
	id, exist, err := s.userRepository.CheckIfUserExistsByIIN(ctx, *user.IIN)
	if err != nil {
		return nil, err
	}

	if !exist {
		s.userRepository.Create(ctx, user)
	} else {
		user.ID = id
	}

	return user, nil
}

func authentification(request model.AuthRequest) (*model.AuthResponse, error) {
	requestData, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	response, err := http.Post(baseURL+"/api/auth", "application/json", bytes.NewReader(requestData))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Server returned status '%d: %s", response.StatusCode, response.Status)
	}

	var responseJSON model.AuthResponse
	err = json.NewDecoder(response.Body).Decode(&responseJSON)
	if err != nil {
		return nil, err
	}

	return &responseJSON, nil
}

func getName(input string) *string {

	// Define regular expressions to match CN and GIVENNAME values
	cnRegex := regexp.MustCompile(`CN=([^,]+)`)
	// givenNameRegex := regexp.MustCompile(`GIVENNAME=([^,]+)`)

	// Find CN and GIVENNAME values using regular expressions
	cnMatch := cnRegex.FindStringSubmatch(input)
	// givenNameMatch := givenNameRegex.FindStringSubmatch(input)

	// Check if both CN and GIVENNAME values were found
	if len(cnMatch) > 1 /* && len(givenNameMatch) > 1 */ {
		cnValue := cnMatch[1]
		// givenNameValue := givenNameMatch[1]

		// Print the extracted values
		result := cnValue
		return &result
	} else {
		fmt.Println("CN and/or GIVENNAME not found in the input string.")
	}
	return nil
}

func (s Service) GetUser(ctx context.Context, userId int) (models.User, error) {
	return s.userRepository.GetUser(ctx, userId)
}
