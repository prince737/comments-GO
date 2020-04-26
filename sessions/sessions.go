package sessions

import "github.com/gorilla/sessions"

//Store is the session store
var Store = sessions.NewCookieStore([]byte("t0p-s3cr3t"))
