from base64 import encode
import hashlib 
import os

nama = input("Masukkan nama anda: ")
os.system("clear")
encoded = nama.encode()
nama_encoded = hashlib.sha3_256(encoded)
print("Nama sebelum di hash: ",nama)
print ("Nama setelah di hash: ",nama_encoded.hexdigest())