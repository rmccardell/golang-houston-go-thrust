# golang-houston-go-thrust
example go-thrust application created for Golang Houston Meetup

How to use:

1. Navigate to https://secure.meetup.com/meetup_api/oauth_consumers/ and create your own OAuth consumer. Be sure to set the redirect uri to http://localhost:8080/auth.
2. go get https://github.com/rmccardell/golang-houston-go-thrust/
3. Inside of go-thrust-main.go edit the "meetUpAuth" variable to match your oauth2 consumer data.
4. Build go-thrust-main.go and run it.
