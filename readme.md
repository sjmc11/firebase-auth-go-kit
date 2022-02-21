<!-- PROJECT LOGO -->
<br />
<div align="center">

<h3 align="center">Firebase Auth Go Starter Kit</h3>

  <p align="center">
    Firebase ready starter API project written in Go
    <br />
    <a href="https://documenter.getpostman.com/view/5420516/UVkmQGdv"><strong>Explore the API docs »</strong></a>
  </p>
</div>


<!-- ABOUT THE PROJECT -->
## About The Project

This is a starter project for creating an App that utilises Googles Firebase Authentication system.

It allows you to get started with verifying tokens sent from your front-end application and focus on developing your application.

Here's why:
* Extendable minimalist project
* Graceful API authentication middleware
* Covers basic app requirements in order to focus on meaningful development
  * Bootstrapping, Routing & Server, DB, Caching, Auth etc..



### Built With

This project utilizes the following packages. Many of these can be replaced and interchanged with other packages of your choosing.

* Router - [Gin](https://example.com)
* Database driver - [PGX](https://example.com)
* Query Builder - [SQLB](https://example.com)
* Caching - [go-cache](https://github.com/patrickmn/go-cache)
* Commands - [Cobra](https://github.com/spf13/cobra)
* Firebase - [Admin SDK](https://github.com/firebase/firebase-admin-go)

<p align="right">(<a href="#top">back to top</a>)</p>



<!-- GETTING STARTED -->
## Getting Started

To get a local copy up and running clone the repo and follow these steps.

1. In the project directory, install dependencies & packages with `go get`
2. Create an `.env` file - link to your DB & key file, use .env.example for reference.
3. Create a Postgres database for the application
4. Run the table migration using `go run main.go migrate`
   1. This will auto-generate a test user allowing you to pass the auth token "development" for testing in postman
5. Serve the application using `go run main.go serve`

### Prerequisites

Before getting started you should have Firebase configured and a private key file generated.
For information on generating your private key file visit [the firebase docs.](https://firebase.google.com/docs/cloud-messaging/auth-server)

Ideally you should have a front-end application that at the least implements Googles Firebase UI. Alternatively, you can test all the functionality using the token "development" through postman.

<!-- USAGE -->
## Usage


When a user logs in to your front-end application using the Firebase UI, pass the users accessToken as the Authorization header on requests to this API.

All requests go through token check middleware. Once verified by the Firebase Admin SDK, the user is verified in the systems db and additional checks are performed.

The `user/sync` endpoint acts to verify a user and return their application specific user information to use in your front-end.

This example Go App is written as a private application, this means users must be pre-registered in the db (invited) in order to pass authentication checks.

Admin users can create accounts for other users.

Manually create an account through the command

`go run main.go user:create`

When a registered user logs in for the first time, their UUID is captured and recorded.

### Available commands

- Create user & optional project `go run main.go user:create`
- Run migrations & create test user `go run main.go migrate`
- Serve application `go run main.go serve`

<p align="right">(<a href="#top">back to top</a>)</p>


<!-- ROADMAP -->
## Roadmap

- [x] ReadMe
- [x] Postman collection
- [ ] Testing
- [ ] Caching
- [ ] Improved readme & documentation
- [ ] Email notification on user registrations

<p align="right">(<a href="#top">back to top</a>)</p>



<!-- CONTRIBUTING -->
## Contributing

Contributions are what make the open source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

If you have a suggestion that would make this better, please fork the repo and create a pull request. You can also simply open an issue with the tag "enhancement".
Don't forget to give the project a star! Thanks again!

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

<p align="right">(<a href="#top">back to top</a>)</p>



<!-- CONTACT -->
## Contact

Jack Callow - [@jvckcallow](https://twitter.com/jvckcallow) - [Linkedin](https://www.linkedin.com/in/jack-callow-11002b8a/) - sjmc11@gmail.com

Project Link: [https://github.com/sjmc11/firebase-auth-go-kit](https://github.com/sjmc11/firebase-auth-go-kit)

<p align="right">(<a href="#top">back to top</a>)</p>
