# monster-house

Thanks for looking at my submission! I called my image repository "Monster House" because it consumes images like
there's no tomorrow.

# Steps to Run Locally

1. In a terminal pane, start the webapp
   `cd webapp && yarn start`

2. run the persistent localstack instance running to mock aws
   `cd server && docker-compose up`

Now we have to spin up an ngrok instance to properly expose our locally-stored files to
the imagga API

3. Start ngrok in another terminal
   `./ngrok http 4566`

4. copy the ngrok url into the `ngrokURL` variable in `config.yaml`:
   The url should look like: `https://cd6017426bf0.ngrok.io`

5. In another terminal pane, start the go server
   `cd server && go run main.go`

# Testing

Sorry for not writing any tests. It's my last week of work and I'm trying to close up my tickets.
