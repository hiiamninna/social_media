package config

import (
	"fmt"
	"social_media/controller"
	"social_media/library"
	"social_media/repository"
	"time"
)

const TIME_FORMAT = "2006-01-02 15:01:02 "

type Context struct {
	CFG Config
	CTL controller.Controller
	JWT library.JWT
	S3  library.S3
}

func NewContext() (Context, error) {
	// read config
	config, err := NewConfiguration()
	if err != nil {
		fmt.Println(time.Now().Format(TIME_FORMAT), "new config : "+err.Error())
	}

	// set up jwt
	jwt := library.NewJWT(config.App.JWTSecret)

	// set up db
	db, err := library.NewDatabaseConnection(config.DB)
	if err != nil {
		fmt.Println(time.Now().Format(TIME_FORMAT), "new db : "+err.Error())
	}

	s3, err := library.NewS3(config.S3Config)
	if err != nil {
		fmt.Println(time.Now().Format(TIME_FORMAT), "new s3 : "+err.Error())
	}

	// set up repo and controller
	repo := repository.NewRepository(db)
	ctl := controller.NewController(repo, jwt, config.App.BcryptSalt, s3)

	return Context{
		CFG: config,
		JWT: jwt,
		CTL: ctl,
		S3:  s3,
	}, nil
}
