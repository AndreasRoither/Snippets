
<h1 align="center">
  <!--<a name="logo" href=""><img src="" alt="Logo" width="200"></a>-->
  <br>
  Snippets
  <br>
  - Archived -

  [![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://opensource.org/licenses/MIT)
</h1>

# Foreword

This project is archived and meant for future reference. Any findings to CI/CD pipelines regarding the different cloud service provider can be found in the PDF `./Cloud_Computing_Documentation.pdf`. The project app and the state of the pipelines can not be guaranteed and will likely need some upgrading on your part if you chose to adopt them.

## Overview

This project is an exploration of different cloud service providers and how pipelines for 4 major cloud service provider (AWS, Google, Azure, Alibaba) can be implemented.
The project itself consists of two parts:

- Electron app
- Rest API

The main parts of this project besides the app are the files for the cloud build pipelines and container registry pushes. They can be found under the `.github/` folder. For Alibaba, the pipeline has to be implemented by yourself with Jenkins (for example) and can be found in the `./pipeline` folder.

## Motivation and Goal of the App

Sometimes when writing code you forget what the optimal solution for a particular task like creating a button listener is. Furthermore, a lot of code has to be remembered which can be quite challenging.

The goal of this project is to create a service where users can save code snippets. Users can log in and view their snippets that they previously added. These snippets can be easily copied and tweaked to the users liking.

# Project parts

## Electron app

The app has been built with electron and typescript. As a packaging solution we have used [electron-forge](https://www.electronforge.io/) with [webpack](https://webpack.js.org/). Webpack allows us to bundle modules dependencies to static assets. 
We didn't use the bundling option for modules but rather to include and pack the main application. The electron app connects to the heroku hosted rest api which will be explained later on.

Demo:  
![](https://i.imgur.com/vg38CJI.png)  
![](https://i.imgur.com/0ZHMFGS.png)  

Snippets can be added by pressing the + sign. To change the snippet name click on name above the editor. Changing the language can be achieved by pressing the language button below the code editor.

## RestAPI

The RestAPI is created using golang and jwt. The database that the api connects to is PostgreSQL. To ensure that the rest api is working correctly tests have been created that test every available functionality for any breaking changes. Since we use travis and heroku for deployment the connection string attributes are collected using environment variables. For CI/CD configuration through environment variables is the preferred way anyway.

### Docker & docker compose

Docker compose is used to start up three containers:  the rest api, database and pgAdmin to monitor and adjust the database.

### Tools Used

- [Electron](https://www.electronjs.org/) as GUI framework  
- [Golang](https://golang.org/) Rest API with [JWT](https://jwt.io/)  
- [Postgres](https://www.postgresql.org/) as database  
- [Docker & Docker Compose](https://www.docker.com/)

#### Development and extensions

- [Visual Studio Code](https://code.visualstudio.com/) for development  
  - [ESLint](https://marketplace.visualstudio.com/items?itemName=dbaeumer.vscode-eslint)  
  - [npm Inellisense](https://marketplace.visualstudio.com/items?itemName=christian-kohler.npm-intellisense)  
  - [Firefox debugger](https://marketplace.visualstudio.com/items?itemName=firefox-devtools.vscode-firefox-debug)  
  - [Prettier](https://marketplace.visualstudio.com/items?itemName=esbenp.prettier-vscode)  
