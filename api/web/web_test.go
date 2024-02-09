package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"github.com/mimirsoft/mimirledger/api/cfg"
	"github.com/mimirsoft/mimirledger/api/datastore"
	"github.com/onsi/gomega"
	"github.com/onsi/gomega/format"
	"github.com/onsi/gomega/types"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"
)

var TestPostgresConfig datastore.PostgresConfig
var TestPostgresClient *sqlx.DB
var TestRouter *chi.Mux
var TestDataStore *datastore.Datastores

func TestMain(m *testing.M) {
	cfg.LoadEnv()
	myConfig := datastore.LoadPostgresConfigFromEnv()
	TestPostgresConfig = myConfig
	myClient, err := datastore.NewClient(&TestPostgresConfig)
	if err != nil {
		panic(err)
	}
	ds := datastore.NewDatastores(myClient)
	TestDataStore = ds
	TestPostgresClient = myClient
	TestRouter = NewRouter(ds, nil)
	result := m.Run()
	os.Exit(result)
}

// setupDatastores wipes postgres for tests
func setupDatastores(ds *datastore.Datastores) {
	if err := TeardownTestAccounts(ds.PGClient()); err != nil {
		log.Panicln(err)
	}
}

// TeardownTestAccounts truncates the transactions_accounts table
func TeardownTestAccounts(client *sqlx.DB) (err error) {
	_, err = client.Exec("TRUNCATE TABLE transaction_accounts CASCADE;")
	return
}

// TableTest represents the methods required to run table tests.
type TableTest interface {
	Exec()
}

// RouterTableTest is a table test of router tests that tracks the current index.
type RouterTableTest struct {
	tableTests []RouterTest
}

// NewRouterTableTest creates a new ControllerTableTestV2.
func NewRouterTableTest(tableTests []RouterTest) *RouterTableTest {
	return &RouterTableTest{tableTests}
}

// Exec executes table tests and optionally asserts status code (equal) and
// response body (contains string).
func (c *RouterTableTest) Exec() {
	c.invokeRequests(func(test *RouterTest) {
		if test.RespBody != "" {
			test.GomegaWithT.Expect(test.actualRespBody).NotTo(gomega.BeNil())
			test.GomegaWithT.Expect(test.actualRespBody.String()).To(test.bodyMatcher(test.RespBody))
		}
	})
}

// exec runs the table tests, sets the currentIndex, execs the individual
// tests and invokes a callback.
func (c *RouterTableTest) invokeRequests(fn func(test *RouterTest)) {
	for i := range c.tableTests {
		v := &c.tableTests[i]
		v.currentIndex = i
		v.invokeRequestAndCheckRespCode()
		fn(v)
	}
}

// Request is used to test controller actions
type Request struct {
	Method     string
	RequestURL string
	Payload    interface{}
	Router     *chi.Mux
}

// Invoke handles the setup and invocation of controller action tests
func (r Request) Invoke() (response *httptest.ResponseRecorder) {
	response = httptest.NewRecorder()
	request, _ := http.NewRequest(r.Method, r.RequestURL, nil)

	if r.Payload != nil && r.Payload != http.NoBody {
		b, err := json.Marshal(r.Payload)
		if err != nil {
			fmt.Printf("payload marshal error:%s \n", err)
			return
		}
		request, _ = http.NewRequest(r.Method, r.RequestURL, strings.NewReader(string(b)))
	}
	if r.Payload == http.NoBody {
		request, _ = http.NewRequest(r.Method, r.RequestURL, http.NoBody)
	}
	r.Router.ServeHTTP(response, request)
	return
}

// RouterTest represents one API request test.
type RouterTest struct {
	Request
	Code        int
	RespBody    string
	GomegaWithT *gomega.WithT
	// actualRespBody should not be set by tests, use ActualRespBody get method
	actualRespBody *bytes.Buffer
	currentIndex   int // used to track the index of the table test slice
	response       *http.Response
}

// Exec executes a single ControllerTableTestV2 and matches the status code and
// if provided, contains string on the body.
func (c *RouterTest) Exec() {
	c.invokeRequestAndCheckRespCode()
	if c.RespBody != "" {
		c.GomegaWithT.Expect(c.actualRespBody.String()).To(c.bodyMatcher(c.RespBody))
	}
}

// ExecWithUnmarshal does the same thing as Exec except it unmarshals the response
// and does not do a string contains check.
func (c *RouterTest) ExecWithUnmarshal(dest interface{}) {
	c.invokeRequestAndCheckRespCode()
	c.GomegaWithT.Expect(c.actualRespBody).NotTo(gomega.BeNil())
	c.GomegaWithT.Expect(json.Unmarshal(c.actualRespBody.Bytes(), &dest)).To(gomega.Succeed())
}

// invokeRequestAndCheckRespCode invokes the request, matches the code, and sets
// actualRespBody for a ControllerTestV2.
func (c *RouterTest) invokeRequestAndCheckRespCode() {
	rw := c.Request.Invoke()
	if c.Code > 0 {
		if rw.Code != c.Code {
			fmt.Printf("%s \n", rw.Body)
		}
		c.GomegaWithT.Expect(rw.Code).To(c.codeMatcher(c.Code))
	}
	c.response = rw.Result() //nolint:bodyclose
	c.actualRespBody = rw.Body
}

// codeMatcher returns a new requestCodeMatcher.
func (r RouterTest) codeMatcher(expected interface{}) types.GomegaMatcher {
	return &requestCodeMatcher{expected, r.currentIndex}
}

// bodyMatcher returns a new requestBodyMatcher.
func (r RouterTest) bodyMatcher(expected interface{}) types.GomegaMatcher {
	return &requestBodyMatcher{expected, r.currentIndex}
}

// requestCodeMatcher fulfills the Gomega Matcher interface and prints table test
// index with failure messages
type requestCodeMatcher struct {
	expected interface{}
	idx      int
}

func (matcher *requestCodeMatcher) Match(actual interface{}) (success bool, err error) {
	code, ok := actual.(int)
	if !ok {
		return false, fmt.Errorf("[%d]: CodeMatcher matcher expects an int %v",
			matcher.idx, actual)
	}
	return reflect.DeepEqual(code, matcher.expected), nil
}

func (matcher *requestCodeMatcher) FailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("[%d]: Expected\n\t%#v\nto equal status code \n\t%#v",
		matcher.idx, actual, matcher.expected)
}

func (matcher *requestCodeMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("[%d]: Expected\n\t%#v\nnot to equal status code \n\t%#v",
		matcher.idx, actual, matcher.expected)
}

// requestBodyMatcher fulfills the Gomega Matcher interface and prints table test
// index with failure messages
type requestBodyMatcher struct {
	expected interface{}
	idx      int
}

func (matcher *requestBodyMatcher) Match(actual interface{}) (success bool, err error) {
	actualString, ok := actual.(string)
	if !ok {
		return false, fmt.Errorf("[%d]: BodyMatcher matcher requires a string. Got:\n%s",
			matcher.idx, format.Object(actual, 1))
	}
	return strings.Contains(actualString, matcher.expected.(string)), nil
}

func (matcher *requestBodyMatcher) FailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("[%d]: Expected\n\t%#v\nto contain substring \n\t%#v",
		matcher.idx, actual, matcher.expected)
}

func (matcher *requestBodyMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("[%d]: Expected\n\t%#v\nnot to contain substring \n\t%#v",
		matcher.idx, actual, matcher.expected)
}
