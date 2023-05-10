# Task1-Chapter2-TEFA-Telkom

Helsa Nesta Dhaifullah / 5025201005/ Teknik Informatika ITS <br />

## API Routes
Program diatas adalah API untuk User dan Products. API user diperlukan untuk authorization. <br />
- API User : 
  - Untuk register User <br />
  ```POST "/users/signup"```
  - Untuk login User <br />
  ```POST "/users/login"```
- API Products : 
  - Untuk create Product <br />
  ```POST "/products"```
  - Untuk update Product <br />
  ```UPDATE "/products/:id"```
  - Untuk delete Product <br />
  ```DELETE "/products/:id"```
  - Untuk get all Products <br />
  ```GET "/products"```
  - Untuk get one Product <br />
  ```GET "/products/:id"```
  - Untuk search by name Products <br />
  ```GET "/products/search```

Sebagai tambahan, setelah login akan mendapatkan token, dimana token tersebut digunakan untuk authorization proses Create, Update, dan Delete. <br />
Untuk guest user, bisa akses get all products, get one product, dan search product.

## Dokumentasi Postman
1. Post Product <br />
![image](https://github.com/helsanesta/Task1-Chapter2-TEFA-Telkom/assets/70515589/88a8a737-7c3d-4009-9875-4dfaf81bacd5)

2. Fill authorization token for Create, Update, Delete <br>
![image](https://github.com/helsanesta/Task1-Chapter2-TEFA-Telkom/assets/70515589/cecb4710-7acb-4e81-80e7-d64f433e43a7)

3. Get One Product <br>
![image](https://github.com/helsanesta/Task1-Chapter2-TEFA-Telkom/assets/70515589/1f3e9885-ab3e-4da3-9e89-68da250019a9)

4. Get All Product <br>
![image](https://github.com/helsanesta/Task1-Chapter2-TEFA-Telkom/assets/70515589/e631eb95-0d23-4024-b756-7392e0393d64)

5. Search Product by Name <br>
![image](https://github.com/helsanesta/Task1-Chapter2-TEFA-Telkom/assets/70515589/dd1f297d-ce79-48c1-bd26-af604cdf2c1c)

6. Update Product <br>
![image](https://github.com/helsanesta/Task1-Chapter2-TEFA-Telkom/assets/70515589/ccf10a8b-6a4b-4a3f-9479-fce47f424d30)

9. Delete Product <br>
![image](https://github.com/helsanesta/Task1-Chapter2-TEFA-Telkom/assets/70515589/9ddded7c-9b5e-421c-9c87-bd0ca56a35f0)


## How to run this code
1. Jalankan `git clone https://github.com/helsanesta/Task1-Chapter2-TEFA-Telkom.git` di lokal komputer kalian.
2. Kemudian untuk import modul, bisa coba jalankan `go mod tidy` di terminal.
3. Buka env, konfigurasikan PORT dan url mongoDb kalian jika berbeda.
4. Buka mongoDb dan start koneksinya.
5. Jalankan command `go run main.go` pada terminal vscode.
6. Untuk tes api, ada di file .json dan bisa jalankan di postman

## Kendala
Kesulitan di bagian minio, karena minio masih belum familiar (baru mengenal). Sudah download di windows aplikasi minio dan berhasil untuk akses minio lewat url di web browser. Di source code juga sudah ada code untuk upload image ke minio saat Create maupun Update. Namun, file image masih belum berhasil terupload ke minio (namun sukses untuk masuk ke mongoDB).
