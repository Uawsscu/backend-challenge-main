# User Management & Lottery API

Backend API สำหรับจัดการข้อมูลผู้ใช้และค้นหาลอตเตอรี่ พัฒนาด้วย Go โดยใช้ Hexagonal Architecture

## 1. วิธีติดตั้ง (Installation)

### สิ่งที่ต้องเตรียม
- [Docker](https://www.docker.com/products/docker-desktop/) และ [Docker Compose](https://docs.docker.com/compose/install/)
- [Go 1.21+](https://go.dev/dl/) (กรณีต้องการรันแบบ Local)

### การติดตั้ง
```bash
git clone <repository-url>
cd backend-challenge-main
```

## 2. วิธีรันระบบ (Running the system)

รันระบบทั้งหมด (API, MongoDB, Redis, Mongo Express) ด้วยคำสั่งเดียว:

```bash
docker-compose up -d
```

- **API**: http://localhost:8080
- **Mongo Express (UI)**: http://localhost:8081 (User: `admin`, Pass: `admin123`)
- **Redis**: localhost:6379

## 3. คู่มือการใช้งาน JWT Token

ระบบใช้คู่ของ Access Token และ Refresh Token เพื่อความปลอดภัย

1. **การขอ Token**: เรียก API `/api/v1/auth/login` เพื่อรับ `accessToken` และ `refreshToken`
2. **การส่ง Token**: ใส่โทเค็นใน Header ของ Request ตามรูปแบบ:
   `Authorization: Bearer <accessToken>`
3. **การต่ออายุ Token**: เมื่อ Access Token หมดอายุ ให้ส่ง `refreshToken` ไปยัง `/api/v1/auth/refresh`


## 4. ตัวอย่าง API Request และ Response

คุณสามารถนำเข้าไฟล์ [Backend-Challenge.postman_collection.json](Backend-Challenge.postman_collection.json) เข้าสู่ Postman เพื่อดูตัวอย่าง API ทั้งหมดและทดสอบได้ทันที


## 5. Assumptions และ Design Decisions

- **Hexagonal Architecture**: แยก Logic ออกจากส่วนติดต่อภายนอก (HTTP, Database) เพื่อให้ง่ายต่อการทดสอบและเปลี่ยนเทคโนโลยีในอนาคต
- **MongoDB**: ใช้เก็บข้อมูลผู้ใช้และลอตเตอรี่ เนื่องจากขยายตัวได้ง่าย (Scalable)
- **Redis**: ใช้จัดการ Session ของ JWT (Token Blacklisting/Revocation) และช่วยในการจัดการลำดับ Ticket ลอตเตอรี่ (Atomic Allocation)
- **Security**: รหัสผ่านถูกเข้ารหัสด้วย `bcrypt` ก่อนเก็บลงฐานข้อมูล และใช้ JWT ในการยืนยันตัวตน (Stateless Auth)
- **Graceful Shutdown**: ระบบรองรับการปิดตัวอย่างปลอดภัยเพื่อจัดการงานที่ค้างอยู่
