# monster-house

Thanks for looking at my submission! I called my image repository "Monster House" because it consumes images like
there's no tomorrow.

## Features

- Upload multiple photos to the server
- Tag each photo upon upload
- Each photo will also receive auto-generated tags
- A search box to search for photos with a specific tag

## Installation

1. `cd webapp && npm install`
2. Run the migrations by hand in `server/migrations`
3. After starting the localstack instance (which you'll do in the instructions below), create the `monsterhouse`
   s3bucket with `cd server && make create-bucket`

## Steps to Run Locally

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

## Technology

- React Frontend
- Golang Backend
- AWS localstack hosting s3 from docker
- postgres hosted from docker
- Ngrok to expose images to the internet
- The Imagga API to auto-tag images

## Standards:

### Webapp

- `TSDoc` for docstrings above functions
- Standard TypeScript `Prettier` with `ESLint`

### Server

- `gofmt`

## Testing

Sorry for not writing any tests. It's my last week of work and I'm trying to close up my tickets.
