# Wadood API

![Wadood Logo](https://img.icons8.com/color/96/pet-commands-train.png)

Wadood is a simple, robust, and extensible RESTful API for managing pets and user accounts, built with Go (Golang) and Gin. It is designed for pet lovers, shelters, and anyone who wants to build pet care or adoption platforms. Wadood provides secure authentication, OTP-based registration, and full CRUD operations for pets.

---

## 🚀 Features

### 🐾 Pet Management
- **Add Pet**: Register a new pet with name, type, gender, and age.
- **Edit Pet**: Update pet details by ID.
- **Delete Pet**: Remove a pet by ID.
- **List Pets**: Retrieve all pets owned by the authenticated user.

### 👤 User Authentication & Security
- **User Registration**: Register with username, email, password, phone number, and device IMEI.
- **Login/Logout**: Secure login and logout endpoints.
- **OTP Verification**: Phone-based OTP for registration and password reset.
- **Password Reset**: Request and verify OTP for password recovery.
- **Multi-language Support**: API messages in English and Arabic (customizable).

### 🛡️ Security
- **JWT Authentication**: All pet endpoints require a valid Bearer token.
- **Password Complexity**: Enforced strong password rules.
- **Input Validation**: Email, phone, and password validation.

### 🗄️ Database
- **MySQL**: Uses MySQL for persistent storage.
- **Auto-migration**: Tables for users and pets are created automatically.

### 🔔 SMS Integration
- **Twilio**: OTPs are sent via Twilio SMS (configurable in `utils/twilio.go`).

---

## 📚 API Endpoints

### Auth
| Method | Endpoint                        | Description                       |
|--------|----------------------------------|-----------------------------------|
| POST   | `/register`                     | Register a new user (OTP sent)    |
| POST   | `/login`                        | Login and receive JWT             |
| POST   | `/logout`                       | Logout user                       |
| POST   | `/send_otp`                     | Send OTP to phone                 |
| POST   | `/verify_otp`                   | Verify OTP for registration       |
| POST   | `/resend_otp`                   | Resend OTP                        |
| POST   | `/forget_password`              | Request password reset OTP        |
| POST   | `/verify_forget_password_otp`   | Verify OTP for password reset     |

### Pets (Require Bearer Token)
| Method | Endpoint         | Description           |
|--------|------------------|-----------------------|
| GET    | `/get_pets`      | List all user pets    |
| POST   | `/add_pet`       | Add a new pet         |
| PUT    | `/edit_pet/:id`  | Edit pet by ID        |
| DELETE | `/delete_pet/:id`| Delete pet by ID      |

---

## 🛠️ Tech Stack
- **Go** (Golang)
- **Gin** (Web Framework)
- **MySQL** (Database)
- **Twilio** (SMS OTP)
- **Air** (Live reload for Go)

---

## 🏁 Getting Started

1. **Clone the repository**
   ```bash
   git clone https://github.com/yourusername/wadood.git
   cd wadood
   ```
2. **Configure Environment**
   - Copy `.env.example` to `.env` and set your MySQL and Twilio credentials.
3. **Run MySQL**
   - Ensure MySQL is running and accessible.
4. **Run the API**
   - The Wadood API comes with [Air](https://github.com/cosmtrek/air) hot-reload support for rapid development.
   ```bash
   air
   ```
   The server will start on `localhost:8080` and automatically reload on code changes.

---

## 📝 License

This project is licensed under the Apache License 2.0. See [LICENSE](LICENSE) for details.

---

## 🤝 Contributing

Contributions, issues, and feature requests are welcome! Feel free to fork the repo and submit a pull request.

---

## 📬 Contact

For support or business inquiries, open an issue or contact the maintainer.

---

> Made with ❤️ for pets and their people.
