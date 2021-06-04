package post

import (
	"../../pkg"
	"crypto/md5"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"os/exec"
	"strings"
)

const (
	userkey = "user"
)

func Login(c *gin.Context) {
	session := sessions.Default(c)
	tableName := "user_passwd"
	username := c.PostForm("username")
	password := c.PostForm("password")
	passwordMd5Form := fmt.Sprintf("%x", md5.Sum([]byte(password)))

	db, err := sql.Open("mysql", "root:hushanglai@tcp(localhost:3306)/web")
	if err != nil {
		log.Printf("sql.Open: %v", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Printf("Fail to close db:%v", err)
		}
	}(db)

	// Query password for a specified username
	sqlString := fmt.Sprintf("select password from %s where username=\"%s\"",
		tableName, username)

	res, err := db.Query(sqlString)
	defer func(res *sql.Rows) {
		err := res.Close()
		if err != nil {
			log.Printf("res.Close: %v", err)
		}
	}(res)
	if err != nil {
		log.Printf("db.Query: %v", err)
	}

	var passwordMd5Mysql string
	for res.Next() {
		err = res.Scan(&passwordMd5Mysql)
		if err != nil {
			c.String(http.StatusUnauthorized, "username or password wrong!")
		}
	}
	if passwordMd5Mysql != passwordMd5Form {
		c.String(http.StatusUnauthorized, "username or password wrong!")
		return
	}
	c.String(http.StatusOK, "Login successful")

	// Save the username in the session
	session.Set(userkey, username) // In real world usage you'd set this to the users ID
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully authenticated user"})
}

// login is a handler that parses a form and checks for specific data
// func Login(c *gin.Context) {
// 	session := sessions.Default(c)
// 	username := c.PostForm("username")
// 	password := c.PostForm("password")
//
// 	// Validate form input
// 	if strings.Trim(username, " ") == "" || strings.Trim(password, " ") == "" {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Parameters can't be empty"})
// 		return
// 	}
//
// 	// Check for username and password match, usually from a database
// 	if username != "hello" || password != "itsme" {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
// 		return
// 	}
//
// 	// Save the username in the session
// 	session.Set(userkey, username) // In real world usage you'd set this to the users ID
// 	if err := session.Save(); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{"message": "Successfully authenticated user"})
// }
func Logout(c *gin.Context) {
}

func Create(c *gin.Context) {
	err := c.Request.ParseForm()
	if err != nil {
		log.Printf("c.Request.ParseForm: %v", err)
	}

	db, err := sql.Open("mysql", "root:hushanglai@tcp(localhost:3306)/web")
	if err != nil {
		log.Printf("sql.Open: %v", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Printf("Fail to close db:%v", err)
		}
	}(db)

	tableName := "user_passwd"
	field1 := "username"
	field2 := "password"
	username := c.PostForm("username")
	password := c.PostForm("password")
	confirmedPassword := c.PostForm("confirmed_password")

	if password != confirmedPassword {
		c.String(http.StatusBadRequest, "Two password not match")
		return
	}

	// Simple filter
	if username != "" && password != "" {
		if len(password) < 7 {
			c.String(http.StatusBadRequest, "Password less or equal than 7!")
			return
		}
		if !strings.ContainsAny(password, "./<>,()_+-=[]{}|") {
			c.String(http.StatusBadRequest, "Password must contains at least on special character!")
			return
		}

		/*
		 * Secure password with md5, need to add salt?
		 */
		password = fmt.Sprintf("%x", md5.Sum([]byte(password)))
		/*
		 * Store username and password to mysql
		 */
		// Sql injection?
		sqlString := fmt.Sprintf("insert into %s (%s, %s) values (\"%s\", \"%s\")",
			tableName, field1, field2, username, password)
		_, err = db.Exec(sqlString)
		if err != nil {
			log.Printf("create account in mysql: %v", err)
		}
	}

	_, err = fmt.Fprintf(c.Writer, "Created successful!\nBecareful with your password?")
	if err != nil {
		log.Print(err)
	}
}

func Upload(c *gin.Context) {
	file, _ := c.FormFile("file")
	log.Println(file.Filename)

	// Upload the file to specific dst.
	err := c.SaveUploadedFile(file, "upload/" + file.Filename)
	if err != nil {
		log.Printf("c.SaveUploadedFile: %v", err)
	}

	c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
}

type Cmd struct {
	cmd string
}

func WebShell(c *gin.Context) {
	var cmdName string
	var cmdArgs string

	str := Cmd {
		cmd:   c.PostForm("cmd"),
	}
	cmdName = strings.Split(str.cmd, " ")[0]
	for _, arg := range strings.Split(str.cmd, " ")[1:] {
		cmdArgs += arg + " "
	}

	// Strip useless character
	cmdName = strings.ReplaceAll(cmdName, "\"", "")
	cmdName = strings.TrimSpace(cmdName)
	cmdArgs = strings.ReplaceAll(cmdArgs, "\"", "")
	cmdArgs = strings.TrimSpace(cmdArgs)

	var out []byte
	var err error

	_, err = fmt.Fprintf(c.Writer, "Executing\nName:%v\tArgs:%v\n", cmdName, cmdArgs)
	if err != nil {
		log.Fatalf("fmt.Fprintf: %v", err)
	}

	// Execute with arguments
	if cmdArgs != "" {
		out, err = exec.Command(cmdName, cmdArgs).Output()
		if err != nil {
			_, err := fmt.Fprintf(c.Writer, "Error, %v", err)
			if err != nil {
				log.Fatalf("fmt.Fprintf: %v", err)
			}
		}
	} else {
		// Without arguments
		out, err = exec.Command(cmdName).Output()
		if err != nil {
			_, err := fmt.Fprintf(c.Writer, "Error, %v", err)
			if err != nil {
				log.Fatalf("fmt.Fprintf: %v", err)
			}
		}
	}

	// Print result to web browser
	_, err = fmt.Fprintf(c.Writer, string(out))
	if err != nil {
		log.Fatalf("fmt.Fprintf: %v", err)
	}
}

func Gif(c *gin.Context) {
	pkg.Lissajous(c.Writer)
}

func Playground(c *gin.Context) {
	games := c.PostForm("games")
	// c.Redirect(http.StatusMovedPermanently, "http://localhost:8888/htop_app")

	/*
	 * Build shell string
	 */
	// var cmdArgs string
	// port := "8888"

	// cmdStr := fmt.Sprintf("/usr/bin/shellinaboxd -t -b -p %s --no-beep " +
	// 	"-s '//:coolder:coolder:/:%s -d 10'", port, games)
	// cmdName := strings.Split(cmdStr, " ")[0]
	// for _, arg := range strings.Split(cmdStr, " ")[1:] {
	// 	cmdArgs += " " + arg + " "
	// }

	// Strip useless character
	// cmdName = strings.ReplaceAll(cmdName, "\"", "")
	// cmdName = strings.TrimSpace(cmdName)
	// cmdArgs = strings.ReplaceAll(cmdArgs, "\"", "")
	// cmdArgs = strings.TrimSpace(cmdArgs)

	// _, err := exec.Command(cmdName, cmdArgs).Output()
	// if err != nil {
	// 	log.Printf("exec.Command: %v", err)
	// }
	// argStr := fmt.Sprintf("'/%s/:coolder:coolder:/:%s -d 10'", games, games)
	// _, err := exec.Command("shellinaboxd","-t", "-b", "-p", "8888", argStr).Output()
	switch {
	case games == "sl":
		_, err := exec.Command("shellinaboxd", "-t", "-b", "-p", "8889", "--no-beep", "-s", "'/sl/:coolder:coolder:/:sl -d 10'").Output()

		if err != nil {
			log.Printf("exec.Command: %v", err)
		}
		c.Redirect(http.StatusMovedPermanently, "http://localhost:8889/sl")
	case games == "htop":
		_, err := exec.Command("shellinaboxd", "-t", "-b", "-p", "8890", "--no-beep", "-s", "'/htop/:coolder:coolder:/:htop -d 10'").Output()
		if err != nil {
			log.Printf("exec.Command: %v", err)
		}
		c.Redirect(http.StatusMovedPermanently, "http://localhost:8890/htop")
	case games == "asciiquarium":
		_, err := exec.Command("shellinaboxd", "-t", "-b", "-p", "8891", "--no-beep", "-s", "'/asciiquarium/:coolder:coolder:/:asciiquarium -d 10'").Output()
		if err != nil {
			log.Printf("exec.Command: %v", err)
		}
		c.Redirect(http.StatusMovedPermanently, "http://localhost:8891/asciiquarium")
	}

}
