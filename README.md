# E-Commerce Web Application

## Overview
This is a **multi-page e-commerce web application** built as part of an internship assignment.  
The application allows users to:

- **Sign up and login** securely (JWT authentication)
- **Browse products** with category and price filters
- **Add items to cart** and **view/remove items**
- **Persist cart** even after logging out and logging back in

The application has **separate pages** for different functionalities:

- `index.html` – Home / Product listing
- `login.html` – Login page
- `signup.html` – Signup page
- `cart.html` – Cart page
- `seller.html` – Seller dashboard (add/delete products)

---

## Tech Stack

### Backend
- Language: **Go (Golang)**
- JWT Authentication for secure login/signup
- REST APIs for:
  - CRUD operations on items
  - Filtering items by category and price
  - Adding/removing items to/from cart
- Database: PostgreSQL (or any relational database)
- Router: Gorilla Mux
- ORM: GORM

### Frontend
- HTML, CSS, JavaScript
- Multi-page structure:
  - `index.html` – Products listing
  - `login.html` – Login
  - `signup.html` – Signup
  - `cart.html` – Cart
  - `seller.html` – Seller product management

---

## Project Structure

project-root/
│
├─ frontend/
│ ├─ index.html
│ ├─ login.html
│ ├─ signup.html
│ ├─ cart.html
│ ├─ seller.html
│ └─ app.js
│
├─ backend/
│ ├─ main.go
│ ├─ routes/
│ ├─ controllers/
│ ├─ models/
│ └─ config/
│
├─ .env.example
└─ README.md
