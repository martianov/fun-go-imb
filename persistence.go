package main

import (
	"fmt"
	"os"
	"os/signal"
	"errors"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID          bson.ObjectId `bson:"_id" json:"-"`
	Email       string        `bson:"email" json:"email"`
	Password    string        `bson:"password,omitempty" json:"-"`
	DisplayName string        `bson:"displayName,omitempty" json:"displayName,omitempty"`
}


func DBConnect(address string) *mgo.Session {
	session, err := mgo.Dial(address)
	if err != nil {
		panic(err)
	}

	session.SetMode(mgo.Monotonic, true)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			fmt.Printf("%v captured - Closing database connection", sig)
			session.Close()
			os.Exit(1)
		}
	}()

	return session
}

func DBEnsureIndicesAndDefaults(s *mgo.Session, dbName string) error {
	rootUser := User{}
	rootUser.Email="e@e.e"
	rootUser.Password="eee"
	rootUser.DisplayName="Root"

	CreateUser(s.DB(dbName), &rootUser)

	i := mgo.Index{
		Key:        []string{"email"},
		Unique:     true,
		Background: true,
		Name:       "email",
	}

	return s.DB(dbName).C("users").EnsureIndex(i)
}

func (u *User) Save(db *mgo.Database) error {
	uC := db.C("users")
	_, err := uC.UpsertId(u.ID, bson.M{"$set": u})
	return err
}

func NewUser() (u *User) {
	u = &User{}
	u.ID = bson.NewObjectId()
	return
}

func CreateUser(db *mgo.Database, u *User) *Error {
	uC := db.C("users")
	pwHash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return &Error{Reason: errors.New("Couldn't hash password"), Internal: true}
	}
	u.Password = string(pwHash)
	u.ID = bson.NewObjectId()
	err = uC.Insert(u)
	if mgo.IsDup(err) {
		return &Error{Reason: errors.New("User already exists"), Internal: false}
	}
	return nil
}

func AuthUser(db *mgo.Database, email, password string) (*User, *Error) {
	uC := db.C("users")
	user := &User{}
	err := uC.Find(bson.M{"email": email}).One(user)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, &Error{Reason: errors.New("User wasn't found on our servers"), Internal: false}
		}
		return nil, &Error{Reason: err, Internal: true}
	}
	if user.ID == "" {
		return nil, &Error{Reason: errors.New("No user found"), Internal: false}
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, &Error{Reason: errors.New("Incorrect password"), Internal: false}
	}
	return user, nil
}

func FindUserById(db *mgo.Database, id bson.ObjectId) (*User, *Error) {
	uC := db.C("users")
	user := &User{}
	err := uC.FindId(id).One(user)
	if err != nil {
		return nil, &Error{Reason: err, Internal: true}
	} else if user.ID == "" {
		return nil, &Error{Reason: errors.New("No user found"), Internal: false}
	}
	return user, nil
}