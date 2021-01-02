# monster-house

# Run locally

1. run the persistent localstack instance running to mock aws
   `cd server && docker-compose up`
2. In another terminal pane, start the go server
   `cd server && go run main.go`
3. In another terminal pane, start the webapp
   `cd webapp && yarn start`

planning:

- I want to host localstack on docker
  - I'll use s3 to store the images
  - I have to build an s3 client

I want to test that uploading and everything works

- upload to localstack and see if permissions work and everything

The question: command line or webui

if webui, it'll be a react app communicating with the backend

- I can setup graphql but I'm lazy
- Just rest is enough

docker pull localstack/localstack
