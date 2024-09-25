package service

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/cecchisandrone/smarthome-server/dto"
	"github.com/cecchisandrone/smarthome-server/model"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
)

type Rental struct {
	ConfigurationService *Configuration `inject:""`
}

func (r *Rental) Init() {
}

func (r *Rental) GenerateAccessLink(configuration model.Configuration, booking dto.Booking) (string, error) {
	link := configuration.Rental.Url
	key := os.Getenv("SMARTHOME_CRYPTO_KEY")
	log.Info("Generating access link for booking: ", booking)

	// Concatenate the arguments with | separator
	concatenatedArgs := fmt.Sprintf("%s|%s|%s", booking.StartDate, booking.EndDate, booking.Email)

	// Encrypt the concatenated arguments
	encryptedData, iv, err := encrypt(concatenatedArgs, key)
	if err != nil {
		fmt.Println("Error encrypting data:", err)
		return "", err
	}

	encryptedWithIv := fmt.Sprintf("%s:%s", encryptedData, iv)
	return fmt.Sprintf("%s?token=%s", link, encryptedWithIv), nil
}

// Encrypt function
func encrypt(text string, key string) (string, string, error) {

	// Define iv, key and text in byte format
	bIV := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, bIV); err != nil {
		return "", "", err
	}
	bKey := []byte(key)
	bText := PKCS5Padding([]byte(text), aes.BlockSize, len(text))

	block, err := aes.NewCipher(bKey)
	if err != nil {
		return "", "", err
	}

	ciphertext := make([]byte, len(bText))
	mode := cipher.NewCBCEncrypter(block, bIV)
	mode.CryptBlocks(ciphertext, bText)
	return hex.EncodeToString(ciphertext), hex.EncodeToString(bIV), nil
}

func PKCS5Padding(ciphertext []byte, blockSize int, after int) []byte {
	padding := (blockSize - len(ciphertext)%blockSize)
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}
