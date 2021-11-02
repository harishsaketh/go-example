// First create a Config.
cfg := NewConfig("Example Application", "eu01xxdb465d7fbd762e6778e4a37c221355NRAL")

// Modify Config fields to control agent behavior.
cfg.Logger = NewDebugLogger(os.Stdout)

// Now use the Config the create an Application.
app, err := NewApplication(cfg)
if nil != err {
	fmt.Println(err)
	os.Exit(1)
}

// Now you can use the Application to collect data!  Create transactions
// to time inbound requests or background tasks. You can start and stop
// transactions directly using Application.StartTransaction and
// Transaction.End.
func() {
	txn := app.StartTransaction("myTask", nil, nil)
	defer txn.End()

	time.Sleep(time.Second)
}()

// WrapHandler and WrapHandleFunc make it easy to instrument inbound web
// requests handled by the http standard library without calling
// StartTransaction.  Popular framework instrumentation packages exist
// in the _integrations directory.
http.HandleFunc(WrapHandleFunc(app, "", func(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "this is the index page")
}))
helloHandler := func(w http.ResponseWriter, req *http.Request) {
	// WrapHandler and WrapHandleFunc add the transaction to the
	// inbound request's context.  Access the transaction using
	// FromContext to add attributes, create segments, and notice.
	// errors.
	txn := FromContext(req.Context())

	func() {
		// Segments help you understand where the time in your
		// transaction is being spent.  You can use them to time
		// functions or arbitrary blocks of code.
		defer StartSegment(txn, "helperFunction").End()
	}()

	io.WriteString(w, "hello world")
}
http.HandleFunc(WrapHandleFunc(app, "/hello", helloHandler))
http.ListenAndServe(":8000", nil)
