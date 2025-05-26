package auth

import (
	"crypto/rand"
	"encoding/hex"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	db "swagtask/internal/db/generated"
	"swagtask/internal/utils"

	"github.com/google/uuid"
)

// TODO: write this well
func HandleSignup(queries *db.Queries, w http.ResponseWriter, r *http.Request) {
	const MAX_UPLOAD_SIZE = 10 << 20 // 10MB
	username := r.FormValue("username")
	password := r.FormValue("password")

	r.Body = http.MaxBytesReader(w, r.Body, MAX_UPLOAD_SIZE)
	err := r.ParseMultipartForm(MAX_UPLOAD_SIZE)
	if err != nil {
		if err.Error() == "http: request body too large" {
			http.Error(w, "Uploaded file is too big. Max 10MB allowed.", http.StatusRequestEntityTooLarge)
		} else {
			http.Error(w, "Failed to parse form", http.StatusInternalServerError)
		}
		utils.LogError("Error parsing multipart form: ", err)
		return
	}

	file, handler, err := r.FormFile("img")
	if err != nil {
		if err == http.ErrMissingFile {
			log.Println("No profile picture uploaded.")

			// no pfp case
			hash, err := HashPassword(password)
			if err != nil {
				http.Error(w, "Server error", 500)
				return
			}

			err = queries.SignUpAndCreateDefaultVaultNoPfp(r.Context(), db.SignUpAndCreateDefaultVaultNoPfpParams{
				Username:     username,
				PasswordHash: hash,
			})
			if err != nil {
				utils.LogError("", err)
				http.Error(w, "User exists", 400)
				return
			}
			http.Redirect(w, r, "/login/", http.StatusSeeOther)

			return
		} else {
			http.Error(w, "Error retrieving file from form", http.StatusInternalServerError)
			log.Printf("Error retrieving 'img' file: %v", err)
			return
		}
	}
	defer file.Close() // IMPORTANT: Close the file when done!

	filename := handler.Filename
	fileSize := handler.Size
	fileHeader := handler.Header
	validateImage(filename, fileSize, fileHeader)

	// create valid path
	fileExtension := filepath.Ext(filename)
	newFilename := uuid.New().String() + fileExtension
	savePath := filepath.Join("./web/pfps/", newFilename)

	// create file
	dst, err := os.Create(savePath)
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		log.Printf("Error creating destination file: %v", err)
		return
	}
	defer dst.Close()

	// copy file
	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, "Failed to save file content", http.StatusInternalServerError)
		log.Printf("Error copying file content: %v", err)
		// Clean up the partially created file if copying failed
		os.Remove(savePath)
		return
	}
	log.Printf("File saved to: %s", savePath)

	hash, err := HashPassword(password)
	if err != nil {
		http.Error(w, "Server error", 500)
		return
	}

	pathToPfp := "/pfps/" + newFilename
	err = queries.SignUpAndCreateDefaultVault(r.Context(), db.SignUpAndCreateDefaultVaultParams{
		Username:     username,
		PathToPfp:    pathToPfp,
		PasswordHash: hash,
	})
	if err != nil {
		http.Error(w, "User exists", 400)
		return
	}
	http.Redirect(w, r, "/login/", http.StatusSeeOther)
}

func generateSessionID() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}
func HandleLogin(queries *db.Queries, w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	credentials, err := queries.GetUserCredentials(r.Context(), username)
	if err != nil || !CheckPasswordHash(password, credentials.PasswordHash) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	sessionID := generateSessionID()
	errSesh := queries.CreateSession(r.Context(), db.CreateSessionParams{ID: sessionID, UserID: credentials.ID})
	if errSesh != nil {
		http.Error(w, "Error creating session", http.StatusInternalServerError)
		return
	}
	cookie := http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Secure:   true,
		HttpOnly: true,
		Path:     "/",
	}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/tasks/", http.StatusSeeOther)
}

func HandleLogout(queries *db.Queries, w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err == nil {
		errDeleteCookie := queries.DeleteSession(r.Context(), cookie.Value)
		if errDeleteCookie != nil {
			http.Error(w, "Error deleting cookie, try logging out again", http.StatusInternalServerError)
			return
		}
	}

	http.SetCookie(w, &http.Cookie{
		Name:   "session_id",
		Secure: true,
		MaxAge: -1,
		Path:   "/",
	})
	http.Redirect(w, r, "/login/", http.StatusSeeOther)
}
