- Pertama buat terlebih dahulu go mod nya
- Installation gin 

``go get -u github.com/gin-gonic/gin``

- Installation gorm 

``go get -u gorm.io/gorm``

- lalu tentukan driver sql yang diinstal jika menggunakan mysql install menggunakan kode berikut

``go get -u gorm.io/driver/mysql``

- menggunakan plugin slug

``go get -u github.com/gosimple/slug``

https://medium.com/koding-kala-weekend/dependency-injection-dan-di-php-2c9d24a885cb
https://martinfowler.com/articles/injection.html
--------------
// Repository:
1. Create image/save data image ke dalam tabel campaign_images
2. Ubah is_primary true ke false
   // Service // Dengan kondisi memanggil point 2 pada repository, kemudian panggil repository pada poin 1
   // Save image campaign ke suatu folder
   // tangkap input dan ubah ke struct input
   // handler

#### deleted file in repository without delete in local system

``git rm --cached file.txt``

``git rm --cached ./dir/ -r``


#### menggunakan repo midtrans
``go get -u github.com/veritrans/go-midtrans``