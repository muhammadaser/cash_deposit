# Cash Deposit

Aplikasi Cash Deposit untuk T.N Assignment.

## Route untuk Frontend

Berikut Route untuk Frontend :

- [Root or Home](#root-or-home)
- [Detail Account](#detail-account)
- [New Account](#new-account)
- [Deposit History](#deposit-history)
- [Cash Deposit](#cash-deposit-nasabah)

### Root or Home (`/`)

Root atau home merupakan halaman awal dari aplikasi. Pada halaman ini terdapat form pencarian _account number_. Pencarian akan menampilkan daftar _account number_ yang ada didalam database.

Jika data yang dicari belum terdaftar, maka pengguna dapat membuat account baru dengan mengklik tombol **New Account**.

```sh
 ---------------------------         ------
|Account Number ...         |       | cari |
 ---------------------------         ------

 -----------
|New Acount |
 -----------

| Account Number | Full Name | Action |
| -------------- |:---------:|-------:|
| 03292827       |Jhone Doe  | Detail |
| 03292828       |Jhone Doe 1| Detail |
| 03292829       |Jhone Doe 2| Detail |
| 03292830       |Jhone Doe 3| Detail |
| 03292831       |Jhone Doe 4| Detail |
```

### Detail Account

```url
/{acount_id}/detail
```

Detail Account merupakan halaman ketika detail di klik. Halaman ini menampilkan detail dari acount yang di pilih. Selain menampilkan detail account seperti nama, alamat, email, dan phone no, halaman ini juga menampilkan _total balance_ dari pemilik acoount.

Saat berada didalam halaman ini, pengguna bisa membuat _cash deposit_ baru atau melihat _deposit history_.

### New Account

```url
/new-account
```

New Account merupakan halaman untuk membuat account baru bagi nasabah yang belum memiliki account didalam sistem. Halaman ini memungkinkan pengguna untuk memasukan data-data seperti First Name, Last Name, email, phone_no, address etc. Setelah pembuatan account baru, pengguna akan di arahkan ke halaman [detail account](#detail-account).

### Deposit History

```url
/{account_id}/deposit-history
```

Halaman Deposit History merupakah halaman yang menampilkan history dari transaksi *cash deposit*.

### Cash Deposit Nasabah

```url
/{account_id}/cash-deposit
```

Halaman Cash Deposit merupakan halaman untuk transaksi deposit uang ke bank. Halaman ini menampilkan form inputan untuk deposit uang ke account nasabah.

- get total balance
- get list deposit
- post deposit 