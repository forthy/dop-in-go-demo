package main

import (
	"fmt"

	PS "dop-in-go/predicates/string"

	E "github.com/IBM/fp-go/either"
	F "github.com/IBM/fp-go/function"
	O "github.com/IBM/fp-go/option"
)

type FirstName struct {
	Name string
}

type LastName struct {
	Name string
}

type MiddleName struct {
	Name string
}

/*
data Email = VerifiedEmail | UnverifiedEmail
*/
type Email interface {
	Address() string
}

type VerifiedEmail struct {
	Email string
}

func (m VerifiedEmail) String() string {
	return fmt.Sprintf("Verified:[%s]", m.Email)
}

type UnverifiedEmail struct {
	Email string
}

func (m UnverifiedEmail) String() string {
	return fmt.Sprintf("Unverified:[%s]", m.Email)
}

func (m VerifiedEmail) Address() string {
	return m.Email
}

func (m UnverifiedEmail) Address() string {
	return m.Email
}

type Contact struct {
	FirstName  FirstName
	MiddleName O.Option[MiddleName]
	LastName   LastName
	Email      Email
}

// Predicate -> string -> Option FirstName
func OfFirstNameWithPredicate(predicate func(string) bool) func(string) O.Option[FirstName] {
	return func(name string) O.Option[FirstName] {
		if predicate(name) {
			return O.Some(FirstName{Name: name})
		}
		return O.None[FirstName]()
	}
}

// Predicate -> string -> Option LastName
func OfLastNameWithPredicate(predicate func(string) bool) func(string) O.Option[LastName] {
	return func(name string) O.Option[LastName] {
		if predicate(name) {
			return O.Some(LastName{Name: name})
		}
		return O.None[LastName]()
	}
}

// Predicate -> string -> Option MiddleName
func OfMiddleNameWithPredicate(predicate func(string) bool) func(string) O.Option[MiddleName] {
	return func(name string) O.Option[MiddleName] {
		if predicate(name) {
			return O.Some(MiddleName{Name: name})
		}
		return O.None[MiddleName]()
	}
}

// Predicate -> string -> Option Email
func OfEmailWithPredicate(predicate func(string) bool) func(string) O.Option[Email] {
	return func(address string) O.Option[Email] {
		if predicate(address) {
			return O.Some[Email](UnverifiedEmail{Email: address})
		}
		return O.None[Email]()
	}
}

// use maker functions to utilise the polymorphic nature of the domain
func contactWithMiddleNameOf(
	firstName FirstName,
	lastName LastName,
	middleName MiddleName,
	email Email,
) Contact {
	return Contact{
		FirstName:  firstName,
		MiddleName: O.Some(middleName),
		LastName:   lastName,
		Email:      email,
	}
}

// Test implementation of email verification
func VerifyEmail(m UnverifiedEmail) bool {
	return m.Email != ""
}

type EmailVerificationFailed struct {
	Reason string
}

func (e EmailVerificationFailed) Error() string {
	return fmt.Sprintf("Email verification failed: %s", e.Reason)
}

type EmailAlreadyVerified struct {
	Email string
}

func (e EmailAlreadyVerified) Error() string {
	return fmt.Sprintf("Email already verified: %s", e.Email)
}

// (UnverifiedEmail -> Bool) -> Contact -> Either error Contact
func VerifyContact(verifyEmail func(UnverifiedEmail) bool) func(Contact) E.Either[error, Contact] {
	return func(contact Contact) E.Either[error, Contact] {
		switch v := contact.Email.(type) {
		case UnverifiedEmail:
			if verifyEmail(v) {
				copy := contact
				copy.Email = VerifiedEmail(v)

				return E.Right[error](copy)
			} else {
				return E.Left[Contact, error](EmailVerificationFailed{"Test implementation"})
			}
		case VerifiedEmail:
			return E.Left[Contact, error](EmailAlreadyVerified(v))
		default:
			return E.Left[Contact, error](EmailVerificationFailed{"Unknown email type"})
		}
	}
}

type ContactInitFailed struct {
	FirstName  string
	LastName   string
	MiddleName O.Option[string]
	Email      string
}

func (e ContactInitFailed) Error() string {
	return fmt.Sprintf(
		"Contact initialization failed:\nfirst-name: %s, last-name: %s, %s%s",
		e.FirstName,
		e.LastName,
		O.GetOrElse(func() string {
			return ""
		})(O.Map(func(mn string) string {
			return fmt.Sprintf("middlename: %s,", mn)
		})(e.MiddleName)),
		e.Email,
	)
}

func main() {
	firstName := "Richard"
	lastName := "Chuo"
	email := "test@example.com"
	middleName := "Andrew"

	// FirstName -> LastName -> MiddleName -> Email -> Contact
	curryContactWithMiddleNameOfO := O.Of(F.Curry4(contactWithMiddleNameOf))

	firstNameO := OfFirstNameWithPredicate(PS.InBetween(1)(10))(firstName)
	lastNameO := OfLastNameWithPredicate(PS.InBetween(1)(15))(lastName)
	emailO := OfEmailWithPredicate(PS.ShouldBeEmail)(email)
	middleNameO := OfMiddleNameWithPredicate(PS.InBetween(1)(10))(middleName)

	contactO := F.Pipe4(
		curryContactWithMiddleNameOfO,
		O.Ap[func(LastName) func(MiddleName) func(Email) Contact](firstNameO),
		O.Ap[func(MiddleName) func(Email) Contact](lastNameO),
		O.Ap[func(Email) Contact](middleNameO),
		O.Ap[Contact](emailO),
	)

	fmt.Println(contactO)

	contactE := F.Flow2(
		E.FromOption[Contact](
			func() error {
				return ContactInitFailed{firstName, lastName, O.Of(middleName), email}
			},
		),
		E.Chain(VerifyContact(VerifyEmail)),
	)(contactO)

	fmt.Println(contactE)
}
