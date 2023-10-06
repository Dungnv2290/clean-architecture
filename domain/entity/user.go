package entity

import (
	"context"
	"errors"
	"time"

	"github.com/dungnguyen/clean-architecture/domain/vo"
)

var (
	ErrUserInsufficientBalance = errors.New("user does not have sufficient balance")

	ErrNotFoundUser = errors.New("not found user")

	ErrUpdateUserWallet = errors.New("error update the value of the wallet")

	ErrCreateUser = errors.New("error creating user")

	ErrFindUserByID = errors.New("error fetching user by ID")
)

type (
	// UserRepositoryCreator defines the operation of creating User entity
	UserRepositoryCreator interface {
		Create(context.Context, User) (User, error)
	}

	// UserRepositoryFinder defines the search operation for a user entity
	UserRepositoryFinder interface {
		FindByID(context.Context, vo.Uuid) (User, error)
	}

	// UserRepositoryUpdated defines the update operation of a user entity wallet
	UserRepositoryUpdater interface {
		UpdateWallet(context.Context, vo.Uuid, vo.Money) error
	}

	// User define the user entity
	User struct {
		id        vo.Uuid
		fullName  vo.FullName
		email     vo.Email
		password  vo.Password
		document  vo.Document
		wallet    *vo.Wallet
		typeUser  vo.TypeUser
		roles     vo.Roles
		createdAt time.Time
	}
)

// NewUser is a factory for user
func NewUser(
	ID vo.Uuid,
	fullName vo.FullName,
	email vo.Email,
	password vo.Password,
	document vo.Document,
	wallet *vo.Wallet,
	typeUser vo.TypeUser,
	createdAt time.Time,
) (User, error) {
	switch typeUser.ToUpper() {
	case vo.COMMON:
		return NewCommonUser(
			ID,
			fullName,
			email,
			password,
			document,
			wallet,
			createdAt,
		), nil
	case vo.MERCHANT:
		return NewMerchantUser(
			ID,
			fullName,
			email,
			password,
			document,
			wallet,
			createdAt,
		), nil
	}

	return User{}, vo.ErrInvalidTypeUser
}

// NewCommonUser create new common user
func NewCommonUser(
	ID vo.Uuid,
	fullName vo.FullName,
	email vo.Email,
	password vo.Password,
	document vo.Document,
	wallet *vo.Wallet,
	createdAt time.Time,
) User {
	return User{
		id:       ID,
		fullName: fullName,
		email:    email,
		password: password,
		document: document,
		wallet:   wallet,
		roles: vo.Roles{
			CanTransfer: true,
		},
		typeUser:  vo.COMMON,
		createdAt: createdAt,
	}
}

// NewMerchantUser create new merchant user
func NewMerchantUser(
	ID vo.Uuid,
	fullName vo.FullName,
	email vo.Email,
	password vo.Password,
	document vo.Document,
	wallet *vo.Wallet,
	createdAt time.Time,
) User {
	return User{
		id:       ID,
		fullName: fullName,
		email:    email,
		password: password,
		document: document,
		wallet:   wallet,
		roles: vo.Roles{
			CanTransfer: false,
		},
		typeUser:  vo.MERCHANT,
		createdAt: createdAt,
	}
}

// Withdraw remove value of money of wallet
func (u *User) Withdraw(money vo.Money) error {
	if u.Wallet().Money().Amount().Value() < money.Amount().Value() {
		return ErrUserInsufficientBalance
	}

	u.Wallet().Sub(money.Amount())

	return nil
}

// Deposit add value of money of wallet
func (u *User) Deposit(money vo.Money) {
	u.Wallet().Add(money.Amount())
}

// CanTransfer returns whether it is possible to transfer
func (u User) CanTransfer() error {
	if u.Roles().CanTransfer {
		return nil
	}

	return vo.ErrNotAllowedTypeUser
}

// ID return the id property
func (u User) ID() vo.Uuid {
	return u.id
}

// FullName return the fullName property
func (u User) FullName() vo.FullName {
	return u.fullName
}

// Password return the password property
func (u User) Password() vo.Password {
	return u.password
}

// Email return the email property
func (u User) Email() vo.Email {
	return u.email
}

// Roles return the roles property
func (u User) Roles() vo.Roles {
	return u.roles
}

// TypeUser return the typeUser property
func (u User) TypeUser() vo.TypeUser {
	return u.typeUser
}

// Wallet return the wallet property
func (u User) Wallet() *vo.Wallet {
	return u.wallet
}

// Document return the document property
func (u User) Document() vo.Document {
	return u.document
}

// CreatedAt return the createdAt property
func (u User) CreatedAt() time.Time {
	return u.createdAt
}
