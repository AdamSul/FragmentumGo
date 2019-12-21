// main_test.go
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
)

var a App

func TestMain(m *testing.M) {
	a = App{}
	a.Initialize("/Users/adamsullivan/go/")
	ensureTableExists()
	code := m.Run()
	clearTable()
	os.Exit(code)
}
func ensureTableExists() {
	if _, err := a.DB.Exec(tableCreationQuery); err != nil {
		log.Fatal(err)
	}
}
func clearTable() {
	a.DB.Exec("DELETE FROM intraday_core_cash")
	a.DB.Exec("ALTER TABLE intraday_core_cash AUTO_INCREMENT = 1")
}

const tableCreationQuery = `
CREATE TABLE IF NOT EXISTS
    intraday_core_cash
    (
        source_ref VARCHAR(20) NOT NULL,
        account VARCHAR(10) NOT NULL,
        tran_date DATETIME,
        tran_date_time DATETIME,
        beneficiary VARCHAR(20),
        originator VARCHAR(30),
        tran_dir VARCHAR(10),
        tran_amount DECIMAL(50,15),
        aba VARCHAR(10),
        tran_path VARCHAR(20),
        counter_bank VARCHAR(30),
        tran_timestamp VARCHAR(30),
        INDEX icore_tran_timestamp (source_ref)
    )
    ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 DEFAULT COLLATE=utf8mb4_0900_ai_ci;
`

const tableLoadSQL = `

INSERT INTO intraday_core (source_ref, account, tran_date, tran_date_time, beneficiary, originator, tran_dir, tran_amount, ABA, tran_path, counter_bank, pre_tran_bal, post_tran_bal, tran_timestamp, is_affiliate, is_open) VALUES ('70303P00213000', 'acct24', '2017-03-03 00:00:00', '2017-03-03 06:01:17', 'Bene-202', 'New Originator - 15520', 'C', 46037.200000000000000, '', 'CHIPS Incoming', '#N/A', -932297903.590000030000000, -932251866.389999990000000, '03-MAR-17 06.03.04.510491 AM', 'SG Entity', 'Open');
INSERT INTO intraday_core (source_ref, account, tran_date, tran_date_time, beneficiary, originator, tran_dir, tran_amount, ABA, tran_path, counter_bank, pre_tran_bal, post_tran_bal, tran_timestamp, is_affiliate, is_open) VALUES ('70303P00213100', 'acct52', '2017-03-03 00:00:00', '2017-03-03 06:01:17', 'Bene-1363', 'New Originator - 11733', 'C', 42705.599999999999000, '', 'CHIPS Incoming', '#N/A', 10178419.960000001000000, 10221125.560000001000000, '03-MAR-17 06.03.04.970025 AM', 'SG Entity', 'Open');
INSERT INTO intraday_core (source_ref, account, tran_date, tran_date_time, beneficiary, originator, tran_dir, tran_amount, ABA, tran_path, counter_bank, pre_tran_bal, post_tran_bal, tran_timestamp, is_affiliate, is_open) VALUES ('70303P00213200', 'acct52', '2017-03-03 00:00:00', '2017-03-03 06:01:18', 'Bene-1363', 'New Originator - 11733', 'C', 32832.000000000000000, '', 'CHIPS Incoming', '#N/A', 10221125.560000001000000, 10253957.560000001000000, '03-MAR-17 06.03.05.384726 AM', 'SG Entity', 'Open');
INSERT INTO intraday_core (source_ref, account, tran_date, tran_date_time, beneficiary, originator, tran_dir, tran_amount, ABA, tran_path, counter_bank, pre_tran_bal, post_tran_bal, tran_timestamp, is_affiliate, is_open) VALUES ('70303P00213300', 'acct24', '2017-03-03 00:00:00', '2017-03-03 06:01:20', 'Bene-202', 'New Originator - 35039', 'C', 73771261.209999993000000, '', 'CHIPS Incoming', '#N/A', -932251866.389999990000000, -858480605.179999950000000, '03-MAR-17 06.03.05.823524 AM', 'SG Entity', 'Open');
INSERT INTO intraday_core (source_ref, account, tran_date, tran_date_time, beneficiary, originator, tran_dir, tran_amount, ABA, tran_path, counter_bank, pre_tran_bal, post_tran_bal, tran_timestamp, is_affiliate, is_open) VALUES ('70303P00213400', 'acct8', '2017-03-03 00:00:00', '2017-03-03 06:01:21', 'Bene-115', 'New Originator - 26539', 'C', 656723.550000000050000, '', 'CHIPS Incoming', 'Counterpary Bank-686', -754459363.179999950000000, -753802639.630000000000000, '03-MAR-17 06.15.03.138219 AM', 'SG Entity', 'Open');
INSERT INTO intraday_core (source_ref, account, tran_date, tran_date_time, beneficiary, originator, tran_dir, tran_amount, ABA, tran_path, counter_bank, pre_tran_bal, post_tran_bal, tran_timestamp, is_affiliate, is_open) VALUES ('70303P00213500', 'acct8', '2017-03-03 00:00:00', '2017-03-03 06:01:22', 'Bene-115', 'New Originator - 26539', 'C', 444839.740000000000000, '', 'CHIPS Incoming', 'Counterpary Bank-686', -753802639.630000000000000, -753357799.889999990000000, '03-MAR-17 06.15.03.639907 AM', 'SG Entity', 'Open');
INSERT INTO intraday_core (source_ref, account, tran_date, tran_date_time, beneficiary, originator, tran_dir, tran_amount, ABA, tran_path, counter_bank, pre_tran_bal, post_tran_bal, tran_timestamp, is_affiliate, is_open) VALUES ('70303P00213600', 'acct8', '2017-03-03 00:00:00', '2017-03-03 06:01:23', 'Bene-115', 'New Originator - 26539', 'C', 87351000.000000000000000, '', 'CHIPS Incoming', 'Counterpary Bank-686', -752807799.889999990000000, -665456799.889999990000000, '03-MAR-17 06.15.05.445543 AM', 'SG Entity', 'Open');
INSERT INTO intraday_core (source_ref, account, tran_date, tran_date_time, beneficiary, originator, tran_dir, tran_amount, ABA, tran_path, counter_bank, pre_tran_bal, post_tran_bal, tran_timestamp, is_affiliate, is_open) VALUES ('70303P00213700', 'acct24', '2017-03-03 00:00:00', '2017-03-03 06:01:37', 'Bene-202', 'New Originator - 26539', 'C', 28777.780000000000000, '', 'CHIPS Incoming', '#N/A', -861051045.120000000000000, -861022267.340000030000000, '03-MAR-17 06.03.06.669805 AM', 'SG Entity', 'Open');
INSERT INTO intraday_core (source_ref, account, tran_date, tran_date_time, beneficiary, originator, tran_dir, tran_amount, ABA, tran_path, counter_bank, pre_tran_bal, post_tran_bal, tran_timestamp, is_affiliate, is_open) VALUES ('70303P00213800', 'acct8', '2017-03-03 00:00:00', '2017-03-03 06:01:59', 'Bene-442', 'New Originator - 31802', 'C', 13975.000000000000000, '', 'CHIPS Incoming', '#N/A', -901118512.610000010000000, -901104537.610000010000000, '03-MAR-17 06.03.54.006428 AM', 'SG Entity', 'Open');
INSERT INTO intraday_core (source_ref, account, tran_date, tran_date_time, beneficiary, originator, tran_dir, tran_amount, ABA, tran_path, counter_bank, pre_tran_bal, post_tran_bal, tran_timestamp, is_affiliate, is_open) VALUES ('70303P00213900', 'acct24', '2017-03-03 00:00:00', '2017-03-03 06:02:13', 'Bene-202', 'New Originator - 24000', 'C', 25515.690000000000000, '', 'CHIPS Incoming', '#N/A', -861012528.679999950000000, -860987012.990000010000000, '03-MAR-17 06.03.54.650962 AM', 'SG Entity', 'Open');
INSERT INTO intraday_core (source_ref, account, tran_date, tran_date_time, beneficiary, originator, tran_dir, tran_amount, ABA, tran_path, counter_bank, pre_tran_bal, post_tran_bal, tran_timestamp, is_affiliate, is_open) VALUES ('70303P00214000', 'acct24', '2017-03-03 00:00:00', '2017-03-03 06:02:13', 'Bene-202', 'New Originator - 24000', 'C', 9738.660000000000000, '', 'CHIPS Incoming', '#N/A', -861022267.340000030000000, -861012528.679999950000000, '03-MAR-17 06.03.54.442334 AM', 'SG Entity', 'Open');
INSERT INTO intraday_core (source_ref, account, tran_date, tran_date_time, beneficiary, originator, tran_dir, tran_amount, ABA, tran_path, counter_bank, pre_tran_bal, post_tran_bal, tran_timestamp, is_affiliate, is_open) VALUES ('70303P00214100', 'acct24', '2017-03-03 00:00:00', '2017-03-03 06:02:14', 'Bene-202', 'New Originator - 24000', 'C', 73653.800000000000000, '', 'CHIPS Incoming', '#N/A', -860987012.990000010000000, -860913359.190000060000000, '03-MAR-17 06.03.54.864703 AM', 'SG Entity', 'Open');
INSERT INTO intraday_core (source_ref, account, tran_date, tran_date_time, beneficiary, originator, tran_dir, tran_amount, ABA, tran_path, counter_bank, pre_tran_bal, post_tran_bal, tran_timestamp, is_affiliate, is_open) VALUES ('70303P00214200', 'acct8', '2017-03-03 00:00:00', '2017-03-03 06:02:17', 'Bene-442', 'New Originator - 26539', 'C', 100000.000000000000000, '', 'CHIPS Incoming', '#N/A', 264083354.710000010000000, 264183354.710000010000000, '03-MAR-17 10.57.53.656959 AM', 'SG Entity', 'Open');
INSERT INTO intraday_core (source_ref, account, tran_date, tran_date_time, beneficiary, originator, tran_dir, tran_amount, ABA, tran_path, counter_bank, pre_tran_bal, post_tran_bal, tran_timestamp, is_affiliate, is_open) VALUES ('70303P00214300', 'acct8', '2017-03-03 00:00:00', '2017-03-03 06:02:25', 'Bene-442', 'New Originator - 15108', 'C', 695.920000000000000, '', 'CHIPS Incoming', '#N/A', -901104537.610000010000000, -901103841.690000060000000, '03-MAR-17 06.03.55.082522 AM', 'SG Entity', 'Open');
INSERT INTO intraday_core (source_ref, account, tran_date, tran_date_time, beneficiary, originator, tran_dir, tran_amount, ABA, tran_path, counter_bank, pre_tran_bal, post_tran_bal, tran_timestamp, is_affiliate, is_open) VALUES ('70303P00214500', 'acct24', '2017-03-03 00:00:00', '2017-03-03 06:02:44', 'Bene-202', 'New Originator - 14603', 'C', 47000.000000000000000, '', 'CHIPS Incoming', '#N/A', -860913359.190000060000000, -860866359.190000060000000, '03-MAR-17 06.04.54.355776 AM', 'SG Entity', 'Open');
INSERT INTO intraday_core (source_ref, account, tran_date, tran_date_time, beneficiary, originator, tran_dir, tran_amount, ABA, tran_path, counter_bank, pre_tran_bal, post_tran_bal, tran_timestamp, is_affiliate, is_open) VALUES ('70208S38352400', 'acct23', '2017-03-08 00:00:00', '2017-03-08 15:42:07', 'Bene-12946', 'New Originator - 25233', 'D', 709.100000000000000, 'aba-25', 'FED Outgoing', 'Counterpary Bank-2014', 2819973.440000000000000, 2819264.340000000000000, '08-MAR-17 03.42.54.055206 PM', 'SG Entity', 'Open');
INSERT INTO intraday_core (source_ref, account, tran_date, tran_date_time, beneficiary, originator, tran_dir, tran_amount, ABA, tran_path, counter_bank, pre_tran_bal, post_tran_bal, tran_timestamp, is_affiliate, is_open) VALUES ('70210S39756400', 'acct9', '2017-03-10 00:00:00', '2017-03-10 16:17:23', 'Bene-16988', 'New Originator - 10020', 'D', 6270.000000000000000, 'aba-25', 'CHIPS Outgoing', 'Counterpary Bank-1113', 56132620.939999998000000, 56126350.939999998000000, '10-MAR-17 04.18.53.549825 PM', 'SG Entity', 'Open');
INSERT INTO intraday_core (source_ref, account, tran_date, tran_date_time, beneficiary, originator, tran_dir, tran_amount, ABA, tran_path, counter_bank, pre_tran_bal, post_tran_bal, tran_timestamp, is_affiliate, is_open) VALUES ('70215S41766300', 'acct9', '2017-03-30 00:00:00', '2017-03-30 16:33:20', 'Bene-39958', 'New Originator - 11213', 'D', 11068.200000000000000, 'aba-162', 'CHIPS Outgoing', 'Counterpary Bank-2214', 19791196.000000000000000, 19780127.800000001000000, '30-MAR-17 04.33.54.095570 PM', 'SG Entity', 'Open');
INSERT INTO intraday_core (source_ref, account, tran_date, tran_date_time, beneficiary, originator, tran_dir, tran_amount, ABA, tran_path, counter_bank, pre_tran_bal, post_tran_bal, tran_timestamp, is_affiliate, is_open) VALUES ('70215U012892IN', 'acct1', '2017-03-01 00:00:00', '2017-03-01 00:18:44', 'Bene-1', 'New Originator - 35027', 'D', 96671.880000000000000, 'aba-762', 'FED Outgoing', 'Counterpary Bank-1727', -389683.680000000000000, -486355.560000000000000, '01-MAR-17 12.20.12.550406 AM', 'Not an SG entity', 'Open');
INSERT INTO intraday_core (source_ref, account, tran_date, tran_date_time, beneficiary, originator, tran_dir, tran_amount, ABA, tran_path, counter_bank, pre_tran_bal, post_tran_bal, tran_timestamp, is_affiliate, is_open) VALUES ('70215U012893IN', 'acct1', '2017-03-01 00:00:00', '2017-03-01 00:15:53', 'Bene-2', 'New Originator - 35027', 'D', 18791.000000000000000, 'aba-162', 'FED Outgoing', 'Counterpary Bank-2763', -309186.150000000020000, -327977.150000000020000, '01-MAR-17 12.17.11.488620 AM', 'Not an SG entity', 'Open');

`

/*
CREATE TABLE IF NOT EXISTS intraday_core_cash
(
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    age INT NOT NULL
) */

func TestEmptyTable(t *testing.T) {
	clearTable()
	req, _ := http.NewRequest("GET", "/cashtransactions", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}
func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func TestGetNonExistentCashTransaction(t *testing.T) {
	clearTable()
	req, _ := http.NewRequest("GET", "/cashtransaction/45", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Transaction not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Transaction not found'. Got '%s'", m["error"])
	}
}

func TestCreateCashTransaction(t *testing.T) {
	clearTable()
	payload := []byte(`{"name":"test user","age":30}`)
	req, _ := http.NewRequest("POST", "/cashtransaction", bytes.NewBuffer(payload))
	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)
	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["name"] != "test user" {
		t.Errorf("Expected user name to be 'test user'. Got '%v'", m["name"])
	}
	if m["age"] != 30.0 {
		t.Errorf("Expected user age to be '30'. Got '%v'", m["age"])
	}
	// the id is compared to 1.0 because JSON unmarshaling converts numbers to
	// floats, when the target is a map[string]interface{}
	if m["id"] != 1.0 {
		t.Errorf("Expected user ID to be '1'. Got '%v'", m["id"])
	}
}

func TestGetCashTransaction(t *testing.T) {
	clearTable()
	addCashTransactions(1)
	req, _ := http.NewRequest("GET", "/cashtransaction/1", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
}

func addCashTransactions(count int) {
	if count < 1 {
		count = 1
	}
	for i := 0; i < count; i++ {
		statement := fmt.Sprintf("INSERT INTO intraday_core_cash(name, age) VALUES('%s', %d)", ("User " + strconv.Itoa(i+1)), ((i + 1) * 10))
		a.DB.Exec(statement)
	}
}

func TestUpdateCashTransaction(t *testing.T) {
	clearTable()
	addCashTransactions(1)
	req, _ := http.NewRequest("GET", "/cashtransaction/1", nil)
	response := executeRequest(req)
	var originalUser map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &originalUser)
	payload := []byte(`{"name":"test user - updated name","age":21}`)
	req, _ = http.NewRequest("PUT", "/cashtransaction/1", bytes.NewBuffer(payload))
	response = executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["id"] != originalUser["id"] {
		t.Errorf("Expected the id to remain the same (%v). Got %v", originalUser["id"], m["id"])
	}
	if m["name"] == originalUser["name"] {
		t.Errorf("Expected the name to change from '%v' to '%v'. Got '%v'", originalUser["name"], m["name"], m["name"])
	}
	if m["age"] == originalUser["age"] {
		t.Errorf("Expected the age to change from '%v' to '%v'. Got '%v'", originalUser["age"], m["age"], m["age"])
	}
}
func TestDeleteCashTransactions(t *testing.T) {
	clearTable()
	addCashTransactions(1)
	req, _ := http.NewRequest("GET", "/cashtransaction/1", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
	req, _ = http.NewRequest("DELETE", "/cashtransaction/1", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
	req, _ = http.NewRequest("GET", "/cashtransaction/1", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
}
