# Beatbox

A simple Go application for streaming music from a server.

## Disclaimer

Please note that this application is for educational and demonstration purposes only and should not be used in production. It is intended to be used as a tool for learning and growing your knowledge in Go programming language and related technologies.
It may contain bugs, security vulnerabilities or other issues that could compromise your system's stability. Use it at your own risk.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

- Go version 1.19 or higher
- A database for storing the generated song metadata (dynamoDB)
- An S3 bucket for storing mp3 files

### Installing

1. Clone the repository

2. Build the application

`go build`

### Setting up the AWS environment

Before you can run the terraform script, you will need to have the following installed on your local machine:
- [Terraform](https://www.terraform.io/downloads.html)
- [AWS CLI](https://aws.amazon.com/cli/)
- [AWS Access Key and Secret Key](https://docs.aws.amazon.com/general/latest/gr/aws-sec-cred-types.html#access-keys-and-secret-access-keys)

Once you have the above installed and configured, follow these steps to set up the AWS environment:

1. Initialize Terraform: `terraform init`

2. Run the terraform script: `terraform apply`

3. Review the changes that will be made, and type `yes` to confirm

4. Wait for Terraform to finish creating the AWS resources

5. Verify that the resources have been created in your AWS console

Please note that there may be additional steps or configurations required depending on your specific use case.

```
To utilize CloudFront as the caching mechanism for your files, you may uncomment the appropriate configuration within the 'aws-terraform.tf' file.
Please be advised that this action will necessitate the implementation of code to retrieve the files through CloudFront rather than S3.
```

### Generating Data

The application includes a `dataUtils.go` file in the `beatbox-backend/scripts` directory for generating data.
To use it, uncomment relevant functions in the file to generate JSON files and insert data into your database (default is DynamoDB)

Then run

`go run beatbox-backend/scripts/dataUtils.go`

This will generate records of data as defined by the constants at the top. To change number of data generated, modify `NumData`. Default is 43.

### Starting the Server

To start the server, run the following command:

`go run main.go`

### Command Line Flags

The following flags are available when starting the server:

- `port`: The port to start the server on if running locally (not useable under ngrok mode) (default: "8080")
- `ngrok`: Enable running on ngrok. Default random domain (default: false)
- `hostname`: The hostname to use if running ngrok (only useable if you already own the domain and are registered with ngrok) (default: random)

### Examples

- Start the server locally on port 8080:

`go run main.go -port=8080`

- Start the server using ngrok:

`go run main.go -ngrok`

- Start the server using ngrok with a custom hostname:

`go run main.go -ngrok -hostname=example.com`

Please note that if you want to use a custom hostname with ngrok, you need to have already registered the domain with ngrok and have the proper authentication token.


## Built With

- [Go](https://golang.org/)

## Contributing

Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct, and the process for submitting pull requests to us.

## Authors

- **Jake Macek** - *Initial work* - [macekj](https://github.com/macekj)

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details

## Credits

The music used in this application was obtained from [Pixabay](https://pixabay.com/), a website that offers a wide variety of royalty-free music and sound effects. We would like to extend our thanks to Pixabay for providing this resource.

I apologize that the data generated doesn't reflect real artists/titles, since the files didn't (at time of writing) contain any metadata and pixbay prevents webcrawling, I was unable to automate the creation of real artist/song titles.
Please support the artists at [Pixabay](https://pixabay.com/)