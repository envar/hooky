# hooky

## Installation

- Create a personal access token which the app will use to access the api
- Configure a push webhook for the repo. You can use ngrok to expose
  the app when developing locally. The Payload URL should end with
  '/gh'. You should also create a secret here.
- build and run hooky with the required environment variables below

## Environment variables

~~~~
GITHUB_OWNER= # owner of repo
GITHUB_REPO= # repo name
GITHUB_SECRET= # secret set in push webhook for repo
GITHUB_API_TOKEN= # developer api token with access to repo
~~~~

## Usage

- GETs to url/path/to/file will fetch and return the contents of
  the file in the repo 
- POSTs/PUTs to url/path/to/file with the contents of the file in
  the body of the request will create/update a file in the repo
- DELETEs to url/path/to/file will delete the file

For each request above as well as directly pushing to the repo, the
webhook will be triggered and print the affected files.
