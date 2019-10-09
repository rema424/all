package lib

import (
	"log"

	"github.com/google/uuid"
)

func UUID() string {
	u, err := uuid.NewRandom()
	if err != nil {
		log.Println(err)
		return ""
	}

	return u.String()
}
