package interfaces

import "food-app/utils/mock"

var (
	userApp    mock.UserAppInterface
	foodApp    mock.FoodAppInterface
	fakeUpload mock.UploadFileInterface
	fakeAuth   mock.AuthInterface
	fakeToken  mock.TokenInterface

	s  = NewUsers(&userApp, &fakeAuth, &fakeToken)                       //We use all mocked data here
	f  = NewFood(&foodApp, &userApp, &fakeUpload, &fakeAuth, &fakeToken) //We use all mocked data here
	au = NewAuthenticate(&userApp, &fakeAuth, &fakeToken)                //We use all mocked data here

)
