# Food Delivery Application

The Food Delivery Application is a web-based platform that allows users to browse and order food from various suppliers. The application provides a user-friendly interface to explore different food items, view details, and place orders for home delivery.

## Table of Contents

- [Getting Started](#getting-started)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Features](#features)
- [Endpoints](#endpoints)
- [Technologies Used](#technologies-used)
- [Contributing](#contributing)
- [License](#license)

## Getting Started

### Prerequisites

Before running the application, make sure you have the following installed:

- Go (1.16 or higher)
- PostgreSQL
- Postman (for API testing)

### Installation

1. Clone the repository:

```bash
git clone https://github.com/amiralikh/food-delivery.git
```
Set up the PostgreSQL database and configure the connection details in config.yml file.

Build and run the application:
```bash
go build
./food-delivery
```

The application should now be running on http://localhost:8080.

### Features
User registration and login with JWT authentication.
Browse food items, suppliers, and categories.
Place orders for food items and track order status.
View order history and details.
Supplier and admin functions for managing food items, categories, and orders.
Endpoints


### Technologies Used
Go (Golang) - Backend programming language
PostgreSQL - Database management system
Gorilla Mux - Router and HTTP handler for Go
JWT (JSON Web Tokens) - User authentication
Postman - API testing and documentation
### Contributing
Contributions are welcome! If you find a bug or have an enhancement in mind, please open an issue or submit a pull request.

### License
This project is licensed under the MIT License - see the LICENSE file for details.
