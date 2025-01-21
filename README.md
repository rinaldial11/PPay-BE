
<a name="readme-top"></a>

[![Contributors][contributors-shield]][contributors-url]
[![Forks][forks-shield]][forks-url]

<br />
<div align="center">
  <a href="https://github.com/rinaldial11/PPay-BE">
    <img src="https://media.discordapp.net/attachments/990198975200104459/1330456891569995776/ppay_logo2.png?ex=678e0c09&is=678cba89&hm=7921e863f445e259cd9fb7b915068fee1e82ab4771f65e2bde2a9f00412ec517&=&format=webp&quality=lossless" alt="Logo" width="550" height="200">
  </a>

  <h3 align="center">P-pay E-wallet (RANER Team)</h3>

  <p align="center">
    <br />
    <a href="https://github.com/rinaldial11/PPay-BE"><strong>Explore the docs »</strong></a>
    <br />
  </p>
</div>



<!-- TABLE OF CONTENTS -->
<details>
  <summary>Table of Contents</summary>
  <ol>
    <li>
      <a href="#about-the-project">About The Project</a>
      <ul>
        <li><a href="#built-with">Built With</a></li>
      </ul>
    </li>
    <li>
      <a href="#getting-started">Getting Started</a>
      <ul>
        <li><a href="#prerequisites">Prerequisites</a></li>
        <li><a href="#installation">Installation</a></li>
      </ul>
    </li>
    <li><a href="#usage">Usage</a></li>
    <li><a href="#related-projects">Related Projects</a></li>
    <li><a href="#contact">Contact</a></li>
  </ol>
</details>



<!-- ABOUT THE PROJECT -->
## About The Project

This project aims to build a Rest-API built using gorm and other supporting modules with a focus on Golang.

RANER Team:
1. Ramadhan Rangkuti 
2. Rinaldi Ainurrahman Lengkana
3. Endra Dwi Jamaludin
4. Adivigo Khalishtama
5. Nanda Hasnan

<p align="right">(<a href="#readme-top">back to top</a>)</p>



### Built With

This project is based on the following packages:

* [![Golang][golang-shield]][golang-url]
* [![pg][pg-shield]][pg-url]
* [![gin-gonic][gin_gonic-shield]][gin_gonic-url]
* [![.env][.env-shield]][.env-url]

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- GETTING STARTED -->
## Getting Started

This project worker can follow the steps below:

### Prerequisites

* install golang

### Installation

1. Clone the repo
   ```sh
   git clone https://github.com/rinaldial11/PPay-BE
   ```
2. Install Golang packages
   ```sh
   go mod tidy
   ```
3. create a postgresql database and create a table and enter the data according to the files in the migration/sql folder
4. please configure in .env
5. Run
   ```sh
   go run main.go
   ```

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- USAGE EXAMPLES -->
## Usage

1. Install [postman](https://www.postman.com/)
2. Visit the following link to export P-pay postman workspace 
   ```sh
    https://speeding-astronaut-290968.postman.co/workspace/New-Team-Workspace~dc44c810-df4e-4ea3-8bc6-36f854a76dbf/collection/20936546-4a2277a9-6cd7-4e12-9d46-800514c66ca6?action=share&creator=20936546
   ```
3. Import the workspace that you already have in stage 2 into the postman application
4. Go to P-pay workspace -> auth -> Register -> Login.
5. Please try to do get data with the token. To insert a token, you can do it on the authorization tab and select Bearer Token

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- USAGE EXAMPLES -->
## Related Projects

[Frontend P-pay Wallet](https://github.com/endradwi/PPay-FE)

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- CONTACT -->
## Contact

Rinaldi Ainurrahman Lengkana - [@Rinaldi Ainurrahman engkana](https://www.linkedin.com/in/rinaldilengkana/)

Project Link: [https://github.com/rinaldial11/PPay-BE)

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->
[contributors-shield]: https://img.shields.io/badge/CONTRIBUTORS-%205%20-orange
[contributors-url]: https://github.com/rinaldial11/PPay-BE/graphs/contributors
[forks-shield]: https://img.shields.io/badge/FORKS-%204%20-blue
[forks-url]: https://github.com/rinaldial11/PPay-BE/fork
[golang-shield]: https://img.shields.io/badge/Go-00ADD8?logo=Go&logoColor=white&style=for-the-badge
[golang-url]: https://go.dev/
[pg-shield]: https://img.shields.io/badge/postgresql-4169e1?style=for-the-badge&logo=postgresql&logoColor=white
[pg-url]: https://www.postgresql.org/
[gin_gonic-shield]: https://img.shields.io/badge/gin%20gonic-grey?style=for-the-badge&logo=gin
[gin_gonic-url]: https://gin-gonic.com/
[.env-shield]: https://img.shields.io/badge/dot%20env-grey?style=for-the-badge&logo=dotenv
[.env-url]: https://www.npmjs.com/package/dotenv
