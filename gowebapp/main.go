package main

import (
	"context"
	"database/sql"
	"fmt"
	"html/template"
	"image"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type Adharno struct {
	AdharNumber int64
	//Adharreq    string
	Check     int
	Checkdesc string
	Otpcreate int
	a         []string
	b         string
	d         int
}
type Otps struct {
	Otp     int `json:"otp"`
	Otpdesc string
}
type variable struct {
	preb         float64
	checkbalance float64
}
type Wallet struct {
	// AdharId int64   `json:"adharno"`
	// Balance float64 `json:"balance"`
	AdharId int64
	Balance float64
	Maximum float64
	Minimum float64
	Active  int
	// Status         int     `json:"status"`
	// Statusdes      string  `json:"statudes"`
	// PresentBalance float64 `json:"presentbalance"`
	//bal            float64
	Status         int
	Statusdes      string
	PresentBalance float64
	reqbal         float64
	sno            int
	slno           int
	preb           float64
	curb           float64
	statuss        string
	types          string
	stupdate       string
	Otpcreate      int
	debitotp       int
}
type balanncecheck struct {
	addha int64
	k     []string
	l     string
	m     int
}

var vari variable
var data Adharno
var data1 Otps
var db *sql.DB
var err error
var data2 Wallet
var data3 balanncecheck

type MyImage struct {
	value *image.RGBA
}

//http://localhost:9090/adharcheck
func main() {
	fileServer := http.FileServer(http.Dir("./image"))
	http.Handle("/image/", http.StripPrefix("/image", fileServer))
	fmt.Println("my sql start")
	db, err = sql.Open("mysql", "root:Abyakta@123@tcp(localhost:3306)/isu")
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("connection successful")
	http.HandleFunc("/adharcheck", adharcheck)
	http.HandleFunc("/wrongadhar", wrongadhar)
	http.HandleFunc("/otpcheck", otpcheck)
	http.HandleFunc("/trans", trans)
	http.HandleFunc("/debit", debit)
	http.HandleFunc("/debitsuccess", debitsuccess)
	http.HandleFunc("/credit", credit)
	http.HandleFunc("/creditsuccess", creditsuccess)
	http.HandleFunc("/creditotp", creditotp)
	http.HandleFunc("/debitotp", debitotp)

	http.HandleFunc("/wrongotp", wrongotp)
	http.HandleFunc("/inactive", inactive)
	http.HandleFunc("/maxbalance", maxbalance)
	http.HandleFunc("/minbalance", minbalance)
	http.HandleFunc("/checkbalance", checkbalance)

	err := http.ListenAndServe(":8080", nil) // setting listening port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
	fmt.Println("MYSQL END")
	defer db.Close()
}
func checkbalance(w http.ResponseWriter, r *http.Request) {
	fmt.Print(r.Method)

	if r.Method == "GET" {
		t, _ := template.ParseFiles("checkbalance.html")

		t.Execute(w, nil)
	} else {
		// logic part of log in	t, _ := template.ParseFiles("balancec.html")
		r.ParseForm()
		fmt.Println("Adharno:", r.Form["Adharno"])

		ctx := context.Background()
		tx, err7 := db.BeginTx(ctx, nil)
		if err7 != nil {
			fmt.Println(err7.Error())
		}
		data3.k = r.Form["Adharno"]
		data3.l = strings.Join(data3.k, " ")

		data3.m, _ = strconv.Atoi(data3.l)
		data3.addha = int64(data3.m)
		fmt.Println("balance check addhar:", data3.addha)
		err = tx.QueryRow("SELECT balance FROM isu.wallet WHERE adharno=" + fmt.Sprint(data3.addha)).Scan(&vari.checkbalance)
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println("balance:::", vari.checkbalance)
		http.Redirect(w, r, "balanceshow", http.StatusFound)
	}
}

func adharcheck(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) //get request method

	if r.Method == "GET" {
		t, _ := template.ParseFiles("adharcheck.html")

		t.Execute(w, nil)
	} else {
		r.ParseForm()
		// logic part of log in
		fmt.Println("Adharno:", r.Form["Adharreq"])

		ctx := context.Background()
		tx, err7 := db.BeginTx(ctx, nil)
		if err7 != nil {
			fmt.Println(err7.Error())
		}
		data.a = r.Form["Adharreq"]
		data.b = strings.Join(data.a, " ")

		data.d, _ = strconv.Atoi(data.b)
		if err != nil {
			panic(err)
		}
		data.AdharNumber = int64(data.d)
		fmt.Print("thrthrthth:", data.AdharNumber)
		err = tx.QueryRow("SELECT balance FROM isu.wallet WHERE adharno=" + fmt.Sprint(data.AdharNumber)).Scan(&vari.preb)
		fmt.Println(err)
		if err == nil {
			fmt.Print(vari.preb)
			data.Otpcreate = rand.Intn(8999)
			data.Check = 1
			data.Checkdesc = "ADHAR NO VARIFIED SUCCESSFULLY"
			//http.Redirect(w, r, "otpcheck", http.StatusFound)
			fmt.Print(data.Checkdesc)
			fmt.Print(data.Otpcreate)
			//fmt.Fprintf(w, "<h1> ADHAR NO VARIFIED SUCCESSFULLY%s</h1>", r.URL.Path("/otpcheck"))
			// w.Header().Set("Content-Type", "application/json")
			// json.NewEncoder(w).Encode(data.Checkdesc)
			http.Redirect(w, r, "otpcheck", http.StatusFound)
			tx.Commit()
			return

		} else {
			tx.Rollback()
			tx.Commit()
			data.Check = -1
			data.Checkdesc = "ADHAR NUMBER NOT  LINKED OR ENTER WRONG ADHAR NUMBER"
			// w.Header().Set("Content-Type", "application/json")
			// json.NewEncoder(w).Encode(data.Checkdesc)
			fmt.Print(data.Checkdesc)
			http.Redirect(w, r, "wrongadhar", http.StatusFound)

		}
		//http.Redirect(w, r, "/Template/", http.StatusFound)

		//fmt.Println("password:", r.Form["password"])
		//http.Redirect(w, r, "chand", http.StatusFound)
	}
}

//var err error
func wrongotp(w http.ResponseWriter, r *http.Request) {
	//r.ParseForm()
	fmt.Print(r.Method)
	if r.Method == "GET" {
		t, _ := template.ParseFiles("wrongotp.html")
		t.Execute(w, nil)

	}

}
func minbalance(w http.ResponseWriter, r *http.Request) {
	//r.ParseForm()
	fmt.Print(r.Method)
	if r.Method == "GET" {
		t, _ := template.ParseFiles("minbalance.html")
		t.Execute(w, nil)

	}

}
func maxbalance(w http.ResponseWriter, r *http.Request) {
	//r.ParseForm()
	fmt.Print(r.Method)
	if r.Method == "GET" {
		t, _ := template.ParseFiles("maxbalance.html")
		t.Execute(w, nil)

	}

}
func inactive(w http.ResponseWriter, r *http.Request) {
	//r.ParseForm()
	fmt.Print(r.Method)
	if r.Method == "GET" {
		t, _ := template.ParseFiles("inactive.html")
		t.Execute(w, nil)

	}

}

func wrongadhar(w http.ResponseWriter, r *http.Request) {
	//r.ParseForm()
	fmt.Print(r.Method)
	if r.Method == "GET" {
		t, _ := template.ParseFiles("wrongadhar.html")
		t.Execute(w, nil)

	}

}

func otpcheck(w http.ResponseWriter, r *http.Request) {
	//r.ParseForm()

	if r.Method == "GET" {
		t, _ := template.ParseFiles("otpcheck.html")
		t.Execute(w, nil)

	} else {
		r.ParseForm()
		fmt.Println("Otpfrom webpge:", r.FormValue("Otp"))
		a := r.Form["Otp"]
		b := strings.Join(a, " ")
		d, _ := strconv.Atoi(b)

		if data.Otpcreate == d {
			fmt.Println("OTP varified")
			http.Redirect(w, r, "trans", http.StatusFound)
		} else {
			http.Redirect(w, r, "wrongotp", http.StatusFound)
			//fmt.Fprintf(w, "<h1> %s</h1>", r.URL.Path[len("/otpcheck"):])
		}

	}
}

func creditotp(w http.ResponseWriter, r *http.Request) {
	//r.ParseForm()

	if r.Method == "GET" {
		t, _ := template.ParseFiles("creditotp.html")
		t.Execute(w, nil)

	} else {
		r.ParseForm()
		fmt.Println("Otpfrom webpge:", r.FormValue("Otp"))
		a := r.Form["Otp"]
		b := strings.Join(a, " ")
		d, _ := strconv.Atoi(b)

		if data2.Otpcreate == d {
			fmt.Println("OTP varified")
			http.Redirect(w, r, "creditsuccess", http.StatusFound)
		} else {
			//fmt.Fprintf(w, "<h1>WRONG OTP %s</h1>", r.URL.Path[len("/adharheck"):])
			http.Redirect(w, r, "wrongotp", http.StatusFound)
		}

	}
}
func debitotp(w http.ResponseWriter, r *http.Request) {
	//r.ParseForm()

	if r.Method == "GET" {
		t, _ := template.ParseFiles("debitotp.html")
		t.Execute(w, nil)

	} else {
		r.ParseForm()
		fmt.Println("Otpfrom webpge:", r.FormValue("Otp"))
		a := r.Form["Otp"]
		b := strings.Join(a, " ")
		d, _ := strconv.Atoi(b)

		if data2.debitotp == d {
			fmt.Println("OTP varified")
			http.Redirect(w, r, "debitsuccess", http.StatusFound)
		} else {
			//fmt.Fprintf(w, "<h1>WRONG OTP %s</h1>", r.URL.Path[len("/adharheck"):])
			http.Redirect(w, r, "wrongotp", http.StatusFound)
		}

	}
}

func trans(w http.ResponseWriter, r *http.Request) {
	//r.ParseForm()
	fmt.Print(r.Method)
	if r.Method == "GET" {
		t, _ := template.ParseFiles("trans.html")
		t.Execute(w, nil)

	}

}
func debitsuccess(w http.ResponseWriter, r *http.Request) {
	//r.ParseForm()
	fmt.Print(r.Method)
	if r.Method == "GET" {

		t, _ := template.ParseFiles("debitsuccess.html")
		t.Execute(w, nil)

	}

}
func creditsuccess(w http.ResponseWriter, r *http.Request) {
	//r.ParseForm()
	fmt.Print(r.Method)
	if r.Method == "GET" {
		t, _ := template.ParseFiles("creditsuccess.html")
		t.Execute(w, nil)

	}

}

func debit(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		fmt.Println(r.Method)
		t, _ := template.ParseFiles("debit.html")
		t.Execute(w, nil)

	} else {

		r.ParseForm()
		fmt.Println(r.Method)
		//if r.Method == "GET" {

		fmt.Print("debit wbpge")
		a := r.Form["Adharreq"]
		b := strings.Join(a, " ")
		d, _ := strconv.Atoi(b)
		if err != nil {
			panic(err)
		}
		fmt.Print(" if prt")
		data2.AdharId = int64(d)

		k := r.Form["Balance"]
		l := strings.Join(k, " ")
		m, _ := strconv.ParseFloat(l, 64)
		data2.Balance = m
		fmt.Println(data2.AdharId, ":details:", data2.Balance)
		ctx := context.Background()
		tx, err7 := db.BeginTx(ctx, nil)
		if err7 != nil {
			fmt.Println(err7.Error())
		}
		err = tx.QueryRow("SELECT active FROM isu.wallet WHERE adharno=" + fmt.Sprint(data2.AdharId)).Scan(&data2.Active)
		if err == nil {
			if data2.Active == 1 {
				data2.slno = rand.Intn(10000)
				data2.statuss = "INITIATED"
				st := "INSERT INTO isu.walletledger (SLNO,amount,prevbls,currbls,type,status,adharno) VALUES(" + fmt.Sprint(data2.slno) + "," + fmt.Sprint(data2.Balance) + "," + fmt.Sprint(0) + "," + fmt.Sprint(0) + ",'" + data2.types + "','" + data2.statuss + "'," + fmt.Sprint(data2.AdharId) + ")"
				_, err6 := db.Exec(st)
				if err6 != nil {
					data2.Status = -1
					data2.Statusdes = err6.Error()

				} else {
					err2 := db.QueryRow("SELECT LAST_INSERT_ID()").Scan(&data2.sno)
					fmt.Println("ID", data2.slno)
					if err2 != nil {
						data2.sno = 1

					}
				}

				err3 := tx.QueryRow("SELECT balance,mimbls,maxbls FROM isu.wallet WHERE adharno="+fmt.Sprint(data2.AdharId)).Scan(&data2.preb, &data2.Minimum, &data2.Maximum)
				if err3 == nil {
					data2.curb = data2.preb - data2.Balance
					if data2.Minimum < data2.curb {

						_, err23 := tx.ExecContext(ctx, "UPDATE wallet SET wallet.balance = "+fmt.Sprint(data2.curb)+" WHERE wallet.adharno="+fmt.Sprint(data2.AdharId))
						if err23 == nil {
							fmt.Println("ENTER IN UPDTE WALLET")
							//fmt.Println("ID", total.sno)
							data2.Statusdes = "DEBITSUCCESS"
							_, err24 := tx.ExecContext(ctx, "UPDATE walletledger SET walletledger.prevbls="+fmt.Sprint(data2.preb)+","+"walletledger.currbls="+fmt.Sprint(data2.curb)+","+"walletledger.status='"+data2.statuss+"'\n"+"WHERE walletledger.SLNO="+fmt.Sprint(data2.slno))
							//fmt.Println(err24)
							if err24 == nil {
								data2.Status = 0
								data2.Statusdes = "DEBIT SUCCESS"
								data2.PresentBalance = data2.curb
								data2.debitotp = rand.Intn(10000)
								fmt.Println(data2.curb)

								fmt.Println(data2.curb)
								fmt.Println(data2.Statusdes)
								fmt.Println(data2.debitotp)
								tx.Commit()
								http.Redirect(w, r, "debitotp", http.StatusFound)

								//return
							} else {
								tx.Rollback()
								tx.Commit()
								data2.Status = -1
								data2.Statusdes = "UPDATE YOUR PASS PROBLEM PLEASE VISIT THE BANK"

								fmt.Fprintf(w, "<h1>UPDATE YOUR PASS PROBLEM PLEASE VISIT THE BANK%s</h1>", r.URL.Path[len("/trans"):])
							}
						} else {
							tx.Rollback()
							tx.Commit()
							data2.Status = -1
							data2.Statusdes = "PROBLEM IN UPDATE YOUR BANK ACCOUNT PLEASE VISIT YOUR BANK"

							fmt.Fprintf(w, "<h1>PROBLEM IN UPDATE YOUR BANK ACCOUNT PLEASE VISIT YOUR BANK%s</h1>", r.URL.Path[len("/trans"):])
						}
					} else {
						tx.Rollback()
						tx.Commit()
						data2.Status = -1
						data2.Statusdes = "MINIMUM BALANCE EXCEED"

						http.Redirect(w, r, "minbalance", http.StatusFound)
					}
				} else {
					tx.Rollback()
					tx.Commit()
					data2.Status = -1
					data2.Statusdes = " PROBLEM IN BALANCE FETCH "

					fmt.Println("")
					fmt.Fprintf(w, "<h1>PROBLEM IN BALANCE FETCH%s</h1>", r.URL.Path[len("/trans"):])

				}

			} else {
				tx.Rollback()
				tx.Commit()
				data2.Status = -1
				data2.Statusdes = "ACCOUNT IS NOT ACTIVE"

				http.Redirect(w, r, "inactive", http.StatusFound)

			}

		} else {
			tx.Rollback()
			tx.Commit()
			data2.Status = -1
			data2.Statusdes = "PLEASE ENTER YOUR CORRECT AADHAAR NUMBER"
			http.Redirect(w, r, "wrongadhar", http.StatusFound)

		}

	}

}
func credit(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		fmt.Println(r.Method)
		t, _ := template.ParseFiles("credit.html")
		t.Execute(w, nil)

	} else {

		r.ParseForm()
		fmt.Println(r.Method)
		//if r.Method == "GET" {

		fmt.Print("debit wbpge")
		a := r.Form["Adharreq"]
		b := strings.Join(a, " ")
		d, _ := strconv.Atoi(b)
		if err != nil {
			panic(err)
		}
		fmt.Print(" if prt")
		data2.AdharId = int64(d)

		k := r.Form["Balance"]
		l := strings.Join(k, " ")
		m, _ := strconv.ParseFloat(l, 64)
		data2.Balance = m
		fmt.Println(data2.AdharId, ":details:", data2.Balance)
		ctx := context.Background()
		tx, err7 := db.BeginTx(ctx, nil)
		if err7 != nil {
			fmt.Println(err7.Error())
		}
		err = tx.QueryRow("SELECT active FROM isu.wallet WHERE adharno=" + fmt.Sprint(data2.AdharId)).Scan(&data2.Active)
		if err == nil {
			if data2.Active == 1 {
				data2.slno = rand.Intn(10000)
				data2.statuss = "INITIATED"
				st := "INSERT INTO isu.walletledger (SLNO,amount,prevbls,currbls,type,status,adharno) VALUES(" + fmt.Sprint(data2.slno) + "," + fmt.Sprint(data2.Balance) + "," + fmt.Sprint(0) + "," + fmt.Sprint(0) + ",'" + data2.types + "','" + data2.statuss + "'," + fmt.Sprint(data2.AdharId) + ")"
				_, err6 := db.Exec(st)
				if err6 != nil {
					data2.Status = -1
					data2.Statusdes = err6.Error() // "problem in wallet ledger insertion"
					// w.Header().Set("Content-Type", "application/json")
					// json.NewEncoder(w).Encode(data2)
					// return

				} else {
					err2 := db.QueryRow("SELECT LAST_INSERT_ID()").Scan(&data2.sno)
					fmt.Println("ID", data2.slno)
					if err2 != nil {
						data2.sno = 1

					}
				}

				err3 := tx.QueryRow("SELECT balance,mimbls,maxbls FROM isu.wallet WHERE adharno="+fmt.Sprint(data2.AdharId)).Scan(&data2.preb, &data2.Minimum, &data2.Maximum)
				if err3 == nil {
					data2.curb = data2.preb + data2.Balance
					if data2.Maximum > data2.curb {

						_, err23 := tx.ExecContext(ctx, "UPDATE wallet SET wallet.balance = "+fmt.Sprint(data2.curb)+" WHERE wallet.adharno="+fmt.Sprint(data2.AdharId))
						if err23 == nil {
							fmt.Println("ENTER IN UPDTE WALLET")
							//fmt.Println("ID", total.sno)
							data2.Statusdes = "CREDITSUCCESS"
							_, err24 := tx.ExecContext(ctx, "UPDATE walletledger SET walletledger.prevbls="+fmt.Sprint(data2.preb)+","+"walletledger.currbls="+fmt.Sprint(data2.curb)+","+"walletledger.status='"+data2.Statusdes+"'\n"+"WHERE walletledger.SLNO="+fmt.Sprint(data2.slno))
							//fmt.Println(err24)
							if err24 == nil {
								data2.Otpcreate = rand.Intn(1000)
								data2.Status = 0
								data2.Statusdes = "CREDIT SUCCESS"
								data2.PresentBalance = data2.curb
								fmt.Println(data2.curb)
								// w.Header().Set("Content-Type", "application/json")
								// json.NewEncoder(w).Encode(data2.Statusdes)
								fmt.Println(data2.curb)
								fmt.Println(data2.Statusdes)
								fmt.Println(data2.Otpcreate)
								tx.Commit()
								//fmt.Fprint(w, "<h1>{{.data2.Statusdes}}%s</h1>", r.URL.Path[("/trans"):])
								http.Redirect(w, r, "creditotp", http.StatusFound)

								//return
							} else {
								tx.Rollback()
								tx.Commit()
								data2.Status = -1
								data2.Statusdes = "UPDATE YOUR PASS PROBLEM PLEASE VISIT THE BANK"
								// w.Header().Set("Content-Type", "application/json")
								// json.NewEncoder(w).Encode(data2.Statusdes)
								// return
								//http.Redirect(w, r, "trans", http.StatusFound)
								fmt.Fprintf(w, "<h1>UPDATE YOUR PASS PROBLEM PLEASE VISIT THE BANK%s</h1>", r.URL.Path[len("/trans"):])
							}
						} else {
							tx.Rollback()
							tx.Commit()
							data2.Status = -1
							data2.Statusdes = "PROBLEM IN UPDATE YOUR BANK ACCOUNT PLEASE VISIT YOUR BANK"
							// w.Header().Set("Content-Type", "application/json")
							// json.NewEncoder(w).Encode(data2.Statusdes)
							// return
							//http.Redirect(w, r, "trans", http.StatusFound)
							fmt.Fprintf(w, "<h1>PROBLEM IN UPDATE YOUR BANK ACCOUNT PLEASE VISIT YOUR BANK%s</h1>", r.URL.Path[len("/trans"):])
						}
					} else {
						tx.Rollback()
						tx.Commit()
						data2.Status = -1
						data2.Statusdes = "MAX BALANCE EXCEED"
						// w.Header().Set("Content-Type", "application/json")
						// json.NewEncoder(w).Encode(data2.Statusdes)
						// return
						//http.Redirect(w, r, "trans", http.StatusFound)
						//fmt.Fprintf(w, "<h1>MINIMUM BALANCE EXCEED%s</h1>", r.URL.Path[len("/trans"):])
						http.Redirect(w, r, "maxbalance", http.StatusFound)
					}
				} else {
					tx.Rollback()
					tx.Commit()
					data2.Status = -1
					data2.Statusdes = " PROBLEM IN BALANCE FETCH "
					// w.Header().Set("Content-Type", "application/json")
					// json.NewEncoder(w).Encode(data2.Statusdes)
					// return\
					//http.Redirect(w, r, "trans", http.StatusFound)
					fmt.Println("")
					fmt.Fprintf(w, "<h1>PROBLEM IN BALANCE FETCH%s</h1>", r.URL.Path[len("/trans"):])

				}

			} else {
				tx.Rollback()
				tx.Commit()
				data2.Status = -1
				data2.Statusdes = "ACCOUNT IS NOT ACTIVE"
				// w.Header().Set("Content-Type", "application/json")
				// json.NewEncoder(w).Encode(data2.Statusdes)
				// return
				http.Redirect(w, r, "inactive", http.StatusFound)
				//fmt.Fprintf(w, "<h1>ACCOUNT IS NOT ACTIVE%s</h1>", r.URL.Path[len("/trans"):])
				//http.Redirect(w, r, "trans", http.StatusFound)
			}

		} else {
			tx.Rollback()
			tx.Commit()
			data2.Status = -1
			data2.Statusdes = "PLEASE ENTER YOUR CORRECT AADHAAR NUMBER"
			// w.Header().Set("Content-Type", "application/json")
			// json.NewEncoder(w).Encode(data2.Statusdes)
			//return
			http.Redirect(w, r, "wrongadhar", http.StatusFound)
			//fmt.Fprintf(w, "<h1>PLEASE ENTER YOUR CORRECT AADHAAR NUMBER%s</h1>", r.URL.Path[len("/trans"):])
			//http.Redirect(w, r, "trans", http.StatusFound)
		}

		//if r.Method == "POST" {

		//http.Redirect(w, r, "trans", http.StatusFound)
		//}
		//t.Execute(w, nil)

	}

}
