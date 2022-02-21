<!-- PROJECT SHIELDS -->
[![Contributors][contributors-shield]][contributors-url]
[![Forks][forks-shield]][forks-url]
[![Stargazers][stars-shield]][stars-url]
[![Issues][issues-shield]][issues-url]
[![MIT License][license-shield]][license-url]
[![LinkedIn][linkedin-shield]][linkedin-url]



<!-- PROJECT LOGO -->
<br />
<div align="center">
  <a href="https://github.com/sjmc11/firebase-auth-go-kit">
    <img src="/logo.svg" alt="Logo" width="80" height="80">
  </a>

<h3 align="center">Firebase Auth Go Project</h3>

  <p align="center">
    Firebase ready starter API project written in Go
    <br />
    <a href="https://documenter.getpostman.com/view/5420516/UVkmQGdv"><strong>Explore the API docs »</strong></a>
  </p>
</div>


<!-- ABOUT THE PROJECT -->
## About The Project

This is a demonstration starter kit for creating an App that utilises Googles Firebase Authentication system.

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



<!-- LICENSE -->
## License

Distributed under the MIT License. See `LICENSE.txt` for more information.

<p align="right">(<a href="#top">back to top</a>)</p>



<!-- CONTACT -->
## Contact

Your Name - [@jvckcallow](https://twitter.com/jvckcallow) - sjmc11@gmail.com

Project Link: [https://github.com/sjmc11/firebase-auth-go-kit](https://github.com/sjmc11/firebase-auth-go-kit)

<p align="right">(<a href="#top">back to top</a>)</p>



<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->
[contributors-shield]: https://img.shields.io/github/contributors/othneildrew/Best-README-Template.svg?style=for-the-badge
[contributors-url]: https://github.com/othneildrew/Best-README-Template/graphs/contributors
[forks-shield]: https://img.shields.io/github/forks/othneildrew/Best-README-Template.svg?style=for-the-badge
[forks-url]: https://github.com/othneildrew/Best-README-Template/network/members
[stars-shield]: https://img.shields.io/github/stars/othneildrew/Best-README-Template.svg?style=for-the-badge
[stars-url]: https://github.com/othneildrew/Best-README-Template/stargazers
[issues-shield]: https://img.shields.io/github/issues/othneildrew/Best-README-Template.svg?style=for-the-badge
[issues-url]: https://github.com/othneildrew/Best-README-Template/issues
[license-shield]: https://img.shields.io/github/license/othneildrew/Best-README-Template.svg?style=for-the-badge
[license-url]: https://github.com/othneildrew/Best-README-Template/blob/master/LICENSE.txt
[linkedin-shield]: https://img.shields.io/badge/-LinkedIn-black.svg?style=for-the-badge&logo=linkedin&colorB=555
[linkedin-url]: https://linkedin.com/in/othneildrew
[product-screenshot]: images/screenshot.png
