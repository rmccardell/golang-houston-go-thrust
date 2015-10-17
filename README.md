# golang-houston-go-thrust
example go-thrust application created for Golang Houston Meetup

How to use:

1. Navigate to https://secure.meetup.com/meetup_api/oauth_consumers/ and create your own OAuth consumer. be sure to set the redirect uri to http://localhost:8080/auth.
2. Inside of go-thrust-main.go edit the "meetupAuth" variable to match your oauth2 consumer data.
3. Build go-thrust-main.go and run it.
