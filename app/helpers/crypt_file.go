package helpers

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"os"

	"github.com/goravel/framework/facades"
)

// generateKeyFromPassphraseLaravel membuat kunci 32-byte dari passphrase menggunakan SHA-256.
func generateKeyFromPassphraseLaravel(passphrase string) ([]byte, error) {
	if passphrase == "" {
		return nil, errors.New("FILE_ENCRYPTION_KEY tidak boleh kosong")
	}
	hash := sha256.Sum256([]byte(passphrase))
	return hash[:], nil
}

// EncryptFileLaravelCompat mengenkripsi file dan otomatis menghapus file sumber setelah berhasil.
func EncryptFileLaravelCompat(sourcePath string, destPath string) (string, error) {
	passphrase := facades.Config().GetString("FILE_ENCRYPTION_KEY")
	key, err := generateKeyFromPassphraseLaravel(passphrase)
	if err != nil {
		return "", err
	}

	plaintext, err := os.ReadFile(sourcePath)
	if err != nil {
		return "", fmt.Errorf("gagal membaca file asli: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	iv := make([]byte, 16)
	if _, err := rand.Read(iv); err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCMWithNonceSize(block, 16)
	if err != nil {
		return "", err
	}

	sealedData := gcm.Seal(nil, iv, plaintext, nil)
	tagSize := 16
	ciphertext := sealedData[:len(sealedData)-tagSize]
	tag := sealedData[len(sealedData)-tagSize:]

	var buffer bytes.Buffer
	buffer.Write(iv)
	buffer.Write(tag)
	buffer.Write(ciphertext)

	err = os.WriteFile(destPath, buffer.Bytes(), 0644)
	if err != nil {
		return "", fmt.Errorf("gagal menulis file terenkripsi: %w", err)
	}

	// --- PENAMBAHAN ---
	// Hapus file sumber (file sementara) setelah enkripsi berhasil.
	if err := os.Remove(sourcePath); err != nil {
		// Jika gagal, catat sebagai peringatan tapi jangan gagalkan operasi utama.
		facades.Log().Warningf("gagal menghapus file sumber sementara '%s': %v", sourcePath, err)
	}

	return destPath, nil
}

// DecryptFileToMemory mendekripsi file dan mengembalikan hasilnya sebagai byte slice di memori.
func DecryptFileToMemory(sourcePath string) ([]byte, error) {
	// 1. Ambil passphrase & buat kunci
	passphrase := facades.Config().GetString("FILE_ENCRYPTION_KEY")
	key, err := generateKeyFromPassphraseLaravel(passphrase)
	if err != nil {
		return nil, err
	}

	// 2. Baca file terenkripsi
	encryptedData, err := os.ReadFile(sourcePath)
	if err != nil {
		return nil, fmt.Errorf("gagal membaca file terenkripsi: %w", err)
	}

	// 3. Pisahkan file: [16 IV] + [16 Tag] + [Ciphertext]
	ivSize := 16
	tagSize := 16
	if len(encryptedData) < ivSize+tagSize {
		return nil, errors.New("file terenkripsi tidak valid atau rusak")
	}
	iv := encryptedData[:ivSize]
	tag := encryptedData[ivSize : ivSize+tagSize]
	ciphertext := encryptedData[ivSize+tagSize:]

	// 4. Lakukan dekripsi
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCMWithNonceSize(block, 16)
	if err != nil {
		return nil, err
	}
	ciphertextWithTag := append(ciphertext, tag...)
	plaintext, err := gcm.Open(nil, iv, ciphertextWithTag, nil)
	if err != nil {
		return nil, fmt.Errorf("gagal dekripsi: %w", err)
	}

	// 5. Kembalikan data mentah, bukan menyimpannya
	return plaintext, nil
}
