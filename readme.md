# Zero Trust Architecture Prototype for Sensitive Access Protection

**Proyek Skripsi / Tugas Akhir S1 Teknologi Informasi**
**Judul:** *Perancangan Arsitektur Zero Trust untuk Pengamanan Akses Sensitif pada Backend Sistem Audit Keuangan Internal Perusahaan*

---

## 📖 Latar Belakang & Tujuan

Model keamanan tradisional (*perimeter-based security*) memberikan *implicit trust* setelah proses login berhasil. Pendekatan ini berisiko tinggi ketika kredensial bocor atau sistem telah disusupi.

Proyek ini bertujuan untuk membangun **prototype backend** yang mengimplementasikan **Zero Trust Architecture (ZTA)** dengan prinsip utama:

> **Never Trust, Always Verify**

Setiap request dievaluasi secara **dinamis** berdasarkan identitas, konteks, dan kebijakan keamanan — bukan hanya status login atau session.

Fokus sistem adalah melindungi **resource sensitif berupa data audit keuangan internal perusahaan**.

---

## 🎯 Ruang Lingkup & Tujuan Teknis

* Membuktikan implementasi konsep **Zero Trust Architecture** sesuai NIST SP 800-207
* Memisahkan logika keamanan dan logika bisnis (*decoupled architecture*)
* Menerapkan **default deny policy** (*assume breach*)
* Menyediakan **audit trail terstruktur** untuk kebutuhan forensik dan SIEM

---

## 🛡️ Pilar Zero Trust yang Diimplementasikan

### 1. Identity Verification

* Validasi **role pengguna** pada **setiap request**
* Simulasi identitas menggunakan HTTP Header

### 2. Context-Aware Policy

Akses hanya diberikan jika memenuhi seluruh kondisi berikut:

* Role diizinkan (Manager / Auditor)
* Waktu akses berada pada jam kerja (08:00 – 17:00)
* Request berasal dari jaringan aman (localhost / VPN)

### 3. Assume Breach

* **Default policy: DENY**
* Akses hanya diizinkan jika secara eksplisit memenuhi kebijakan

### 4. Micro-Segmentation Logic

* Logika keamanan dipisahkan dari business logic
* Middleware berperan sebagai **Policy Enforcement Point (PEP)**

### 5. Visibility & Observability

* Audit log dicatat dalam format JSON
* Menggunakan pendekatan **5W1H (Who, What, When, Where, Why, How)**

---

## 🏗️ Arsitektur Sistem

Prototype ini mengimplementasikan komponen logis **NIST SP 800-207** menggunakan bahasa pemrograman **Golang**.

### Pemetaan Komponen NIST

| Komponen NIST                  | Implementasi          | Lokasi Kode                    |
| ------------------------------ | --------------------- | ------------------------------ |
| Subject                        | Pengguna / API Client | HTTP Request Header            |
| Enterprise Resource            | Data Audit Keuangan   | `internal/app/handlers`        |
| Policy Enforcement Point (PEP) | Middleware Gin        | `internal/security/middleware` |
| Policy Decision Point (PDP)    | Policy Engine         | `internal/security/policies`   |
| Trust Context                  | Waktu & IP Address    | `config/config.go`             |

---

## 🔄 Alur Request (Flow of Events)

1. **Request Masuk**
   Client mengakses endpoint `/api/v1/finance/reports`

2. **Intersepsi (PEP)**
   Middleware Zero Trust menahan request

3. **Evaluasi Kebijakan (PDP)**

   * Apakah role diizinkan?
   * Apakah waktu akses sesuai jam kerja?
   * Apakah IP berasal dari jaringan aman?

4. **Keputusan**

   * **ALLOW** → Request diteruskan ke handler
   * **DENY** → Response `403 Forbidden`

5. **Audit Logging**
   Semua keputusan dicatat ke log terstruktur

---

## 📂 Struktur Proyek

Struktur proyek mengikuti **Golang Standard Project Layout** untuk menjaga skalabilitas dan maintainability.

```
backend-zerotrust-skripsi/
├── cmd/
│   └── api/
│       └── main.go            # Entry point server
├── config/
│   └── config.go              # Konfigurasi kebijakan Zero Trust
├── internal/
│   ├── app/                   # Business Layer
│   │   ├── handlers/          # API controller
│   │   └── models/            # Domain models
│   └── security/              # Zero Trust Core
│       ├── audit/             # Audit logger (JSON)
│       ├── middleware/        # PEP
│       └── policies/          # PDP / Policy Engine
├── pkg/
│   └── response/              # Standar response API
├── go.mod
└── README.md
```

---

## 🚀 Cara Menjalankan Aplikasi

### Prasyarat

* Go **v1.18** atau lebih baru
* Terminal (Bash / PowerShell)
* Postman atau cURL

### Instalasi & Run

1. **Clone Repository**

```
git clone https://github.com/username-kamu/backend-zerotrust-skripsi.git
cd backend-zerotrust-skripsi
```

2. **Install Dependency**

```
go mod tidy
```

3. **Jalankan Server**

```
go run cmd/api/main.go
```

Output:

```
🛡️  Zero Trust Architecture Active on port :8080
📜  Policy: Work Hours 08:00 - 17:00 | Safe Net: 127.0.0.1
```

---

## 🧪 Skenario Pengujian

### 1. Akses Berhasil (Authorized)

* Role: `manager`
* Waktu: Jam kerja
* Jaringan: Localhost

```
curl -X GET http://localhost:8080/api/v1/finance/reports \
  -H "X-User-Role: manager" \
  -H "X-User-ID: 999"
```

**Response:** `200 OK`

---

### 2. Akses Ditolak – Role Tidak Diizinkan

```
curl -X GET http://localhost:8080/api/v1/finance/reports \
  -H "X-User-Role: staff" \
  -H "X-User-ID: 101"
```

**Response:** `403 Forbidden`

---

### 3. Akses Ditolak – Di Luar Jam Kerja

Walaupun role valid, request akan ditolak jika di luar jam operasional.

> Untuk simulasi, ubah nilai jam kerja di `config/config.go`

---

## 📊 Audit Logging (Observability)

Setiap request dicatat ke **stdout** dalam format JSON terstruktur.

Contoh log saat akses ditolak:

```json
{
  "timestamp": "2024-05-20T10:15:30+07:00",
  "subject_id": "101",
  "subject_role": "staff",
  "ip_address": "::1",
  "resource": "/api/v1/finance/reports",
  "action": "GET",
  "decision": "DENY",
  "reason": "Role tidak diizinkan mengakses data keuangan"
}
```

---

## 🛠️ Teknologi yang Digunakan

* **Language:** Golang
* **Web Framework:** Gin Gonic
* **Security Model:** Custom Zero Trust Policy Engine (RBAC + ABAC)
* **Logging:** `encoding/json`
* **CORS:** `gin-contrib/cors`

---

## 📝 Catatan & Batasan

* Proyek ini merupakan **Proof of Concept (PoC) akademik**
* Identitas masih disimulasikan via HTTP Header
* Data keuangan masih berupa mock data
* Kebijakan masih berbasis konfigurasi statis

### Rekomendasi Pengembangan Lanjutan

* Integrasi **JWT / OAuth2**
* External Policy Engine (OPA / Rego)
* Persist audit log ke SIEM / ELK Stack
* Integrasi database dan secrets management

---

**Author:** Mahasiswa S1 Teknologi Informasi
**Topik:** Zero Trust Architecture • Backend Security • Golang
