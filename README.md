# Chatbot: Facebook integration

Building Locally
```
go get github.com/Masterminds/glide (dev version, if you don't already have glide)
glide update
glide install

go build
./chatbotfb 
```

### Heroku (for test deployment)

set the environmental variable `TOKEN` to be the API token below

```
EAATZAxfQTVYQBAD5RIvKCpLEK5BQ4TF7V2l6S4OYcWHZAxZAwQ1va2x5zGNZAgEke8ZC7Mik8CKOcwqPmSLZBrZB2PzBaXEeOvhvoxfHwjelZBMLZCGZCOvflQJ1cCSH2nPfdOVih79WoQK0F47I5BI6wetibxz0eTlsiWFv9gPbllZBgZDZD
```

Also make sure `MONGODBURL` is set to the url of your MongoDB database 