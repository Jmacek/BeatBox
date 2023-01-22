# Beatbox frontend

A Node.js application written in TypeScript to act as a music streaming service. (mostly just a display for my go server, so not very UX rich)

## Disclaimer

Please note that this application is for educational and demonstration purposes only and should not be used in production. It is intended to be used as a tool for learning and growing your knowledge in Node.js, TypeScript and related technologies. It may contain bugs, security vulnerabilities or other issues that could compromise your system's stability. Use it at your own risk.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

- Node.js version 19.3.0 or higher
- npm version 9.2.0 or higher
- Go server (see beatbox-backend)

## Configuration

The application uses an API to communicate with the backend server. Before running the application, make sure to change the `apiUrl` in the following file to match the URL of the Go server you're running:

`beatbox-frontend/src/configs/constants.ts`


This file contains a constant named `apiUrl` that is used throughout the application to make API calls. By default, it is set to `http://localhost:8080`, but you will need to change it to match the URL of the Go server you're running.

For example, if you're running the Go server on your local machine on port 8080, you don't need to change anything. However, if you're running the Go server on a different machine or port, you will need to update the `apiUrl` accordingly.

### Installing

1. Clone the repository

2. Install the dependencies

`npm install`

3. Build the application

`npm run build`

### Running the Application

To start the application, run the following command:

`npm run start`

This will start the application on the default port (3000).

### Running the tests

SURPRISE! There are no tests because that wasn't my focus for this project :) maybe i'll add some later

## Built With

- [Node.js](https://nodejs.org/)
- [TypeScript](https://www.typescriptlang.org/)

## Authors

- **Jake Macek** - *Initial work* - [macekj](https://github.com/macekj)

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details
