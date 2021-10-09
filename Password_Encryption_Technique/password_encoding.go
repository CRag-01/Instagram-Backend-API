// The solution is to calculate a hashed password to deliberately increase the amount of resources and time it would take to crack it.
// We design a hash such that nobody could possibly have the resources required to compute the required rainbow table.
package main

import (
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func CheckPass() {
	storedUserPassword := "Password"
	hash, err := bcrypt.GenerateFromPassword([]byte(storedUserPassword), bcrypt.DefaultCost) //crating an encrypted form(hashing) of the given password
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Hash to store:", string(hash))
	newlyEnteredPwd := "Password"
	hashFromDatabase := hash

	// Comparing the hashed forms of both, the stored Password and newly entered one
	if err := bcrypt.CompareHashAndPassword(hashFromDatabase, []byte(newlyEnteredPwd)); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Password was correct!")
}

// This way the old password is not directly visible to the database administraors or who ever have access to it, and hence cannot be misused.
// Only the hashed form of the password is visible to who ever has access to the db and there in no possible way to reverse engineer this hashed for to retrieve the password.
