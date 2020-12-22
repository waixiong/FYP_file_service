package commons

import (
	"context"
	"fmt"

	"net/http"

	googleOauth2 "google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"

	// "google.golang.org/api/people/v1"
	"google.golang.org/grpc/metadata"

	firebase "firebase.google.com/go"
)

type urlParams struct {
	Key   string
	Value string
}

func (v urlParams) Get() (string, string) {
	return v.Key, v.Value
}

func initGoogleConfig() error {
	opt := option.WithCredentialsFile("./configs/key/imagechat-1fb81-firebase-adminsdk-jxsfx-277a90960d.json")
	a, err := firebase.NewApp(context.Background(), nil, opt)
	app = a
	if err != nil {
		return fmt.Errorf("error initializing app: %v", err)
	}
	return nil
}

var app *firebase.App
var httpClient = &http.Client{}

func VerifyGoogleIDToken(ctx context.Context, idToken string, accessToken string) (*googleOauth2.Tokeninfo, *googleOauth2.Userinfo, error) {
	// creds, err := ioutil.ReadFile("./configs/key/GBEC-Project-8614eec529e1.json")
	// // creds, err := ioutil.ReadFile("./configs/key/imagechat-1fb81-firebase-adminsdk-jxsfx-277a90960d.json")
	// if err != nil {
	// 	log.Fatalf("Unknown creds: %v", err)
	// }
	// // cfg, err := goauth.JWTConfigFromJSON(creds, SQLScope)
	// config, err := google.ConfigFromJSON(creds, "https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile")
	// if err != nil {
	// 	log.Fatalf("Unable to parse client secret file to config: %v", err)
	// }
	// authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	// fmt.Printf("Go to the following link in your browser then type the "+
	// 	"authorization code: \n%v\n", authURL)
	// var authCode string
	// if _, err := fmt.Scan(&authCode); err != nil {
	// 	log.Fatalf("Unable to read authorization code: %v", err)
	// }
	// tok, err := config.Exchange(context.TODO(), authCode)
	// if err != nil {
	// 	log.Fatalf("Unable to retrieve token from web: %v", err)
	// }
	// httpClient = config.Client(ctx, tok)

	oauth2Service, err := googleOauth2.New(httpClient)
	// options := []option.ClientOption{}
	// oauth2Service, err := googleOauth2.NewService(ctx,
	// 	option.WithCredentialsFile("./configs/key/imagechat-1fb81-firebase-adminsdk-jxsfx-277a90960d.json"),
	// 	// option.WithAudiences("996344380155-alqg0crulhacekd7s45eha4fbgr5s5lg.apps.googleusercontent.com"),
	// 	option.WithScopes("https://www.googleapis.com/auth/userinfo.email"),
	// 	option.WithScopes("https://www.googleapis.com/auth/userinfo.profile"),
	// )

	if err != nil {
		fmt.Println(err)
		return nil, nil, err
	}
	tokenInfoCall := oauth2Service.Tokeninfo()
	// tokenInfoCall.IdToken(idToken)
	tokenInfoCall.AccessToken(accessToken)
	tokenInfo, err := tokenInfoCall.Do()
	if err != nil {
		fmt.Println(err)
		return nil, nil, err
	}

	// userInfoGetCall := oauth2Service.Userinfo.V2.Me.Get()
	userInfoGetCall := oauth2Service.Userinfo.Get()
	userInfoGetCall.Header().Add("Authorization", accessToken)
	userInfoGetCall.Fields()
	userInfo, err := userInfoGetCall.Do(urlParams{"access_token", accessToken})
	if err != nil {
		fmt.Println(err)
		return nil, nil, err
	}

	return tokenInfo, userInfo, nil
}

// check access token from header for authorization
func VerifyGoogleAccessToken(ctx context.Context) (*googleOauth2.Tokeninfo, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	fmt.Printf("\tAuth Func Get metadata...\n")
	if !ok {
		return nil, ErrMissingMetadata
	}
	accessTokens := md.Get("Authorization")
	fmt.Println(accessTokens)

	oauth2Service, err := googleOauth2.New(httpClient)
	if err != nil {
		return nil, err
	}
	// userInfoGetCall := oauth2Service.Userinfo.Get()
	// userInfoGetCall.Header().Add("Authorization", accessTokens[0])
	// userInfoGetCall.Fields()
	// userInfo, err := userInfoGetCall.Do(urlParams{"access_token", accessTokens[0]})
	// if err != nil {
	// 	fmt.Println(err)
	// 	return nil, err
	// }
	// return userInfo, nil
	tokenInfoCall := oauth2Service.Tokeninfo()
	tokenInfoCall.AccessToken(accessTokens[0])
	tokenInfo, err := tokenInfoCall.Do()
	if err != nil {
		return nil, err
	}
	return tokenInfo, nil
}
