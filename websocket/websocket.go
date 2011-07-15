//target:github.com/droundy/gui/websocket

package websocket

import (
	"fmt"
	"http"
	"websocket"
	"bufio"
	"path"
	"os"
)

type Handler interface {
	AddSend(send func(string))
	Handle(evt string)
	Done(os.Error)
}

// We export a single function, which creates a page controlled by a
// single websocket.  It's quite primitive, and yet quite easy to use!
func Handle(url string, h Handler) {
	myh := func(ws *websocket.Conn) {
		h.AddSend(func (p string) { fmt.Fprintln(ws, p) })
		h.Handle("")
		r := bufio.NewReader(ws)
		for {
			x, err := r.ReadString('\n')
			if err == nil {
				h.Handle(x[:len(x)-1])
			} else {
				h.Done(os.NewError("Error from r.ReadString: " + err.String()))
				return
			}
		}
	}
	http.Handle(path.Join(url, "socket"), websocket.Handler(myh))

	skeleton := func(c http.ResponseWriter, req *http.Request) {
		c.Header().Set("Content-Type", "text/html")
		fmt.Fprintln(c, skeletonpage(req))
	}
	http.HandleFunc(url, skeleton)

}

// We export a single function, which creates a page controlled by a
// single websocket.  It's quite primitive, and yet quite easy to use!
func HandleSeparate(url string, hh func() Handler) {
	myh := func(ws *websocket.Conn) {
		h := hh()
		h.AddSend(func (p string) { fmt.Fprintln(ws, p) })
		fmt.Fprintln(ws, "start")
		r := bufio.NewReader(ws)
		for {
			x, err := r.ReadString('\n')
			if err == nil {
				h.Handle(x[:len(x)-1])
			} else {
				h.Done(os.NewError("Error from r.ReadString: " + err.String()))
				return
			}
		}
	}
	http.Handle(path.Join(url, "socket"), websocket.Handler(myh))

	skeleton := func(c http.ResponseWriter, req *http.Request) {
		c.Header().Set("Content-Type", "text/html")
		fmt.Fprintln(c, skeletonpage(req))
	}
	http.HandleFunc(url, skeleton)

}

// Run handles the case where you want each user should have the
// same session which will look identical.
func Run(url string, port int, handler Handler) os.Error {
	Handle("/", handler)
	return http.ListenAndServe(fmt.Sprint(":", port), nil);
}

// RunSeparate handles the case where you want each user who logs on
// to have a separate session with a separate handler.
func RunSeparate(url string, port int, handler func() Handler) os.Error {
	HandleSeparate("/", handler)
	return http.ListenAndServe(fmt.Sprint(":", port), nil);
}

func skeletonpage(req *http.Request) string {
	wsurl := *req.URL
	wsurl.Host = req.Host
	wsurl.Scheme = "ws" + req.URL.Scheme
	wsurl.Path = "/socket"
	return `<!DOCTYPE HTML>
<html>
<head>
<link href="/style.css" rel="stylesheet" type="text/css" />
<script type="text/javascript">

// Define helper cookie functions:
function createCookie(name,value,days) {
	if (days) {
		var date = new Date();
		date.setTime(date.getTime()+(days*24*60*60*1000));
		var expires = "; expires="+date.toGMTString();
	}
	else var expires = "";
	document.cookie = name+"="+value+expires+"; path=/";
}
function readCookie(name) {
	var nameEQ = name + "=";
	var ca = document.cookie.split(';');
	for(var i=0;i < ca.length;i++) {
		var c = ca[i];
		while (c.charAt(0)==' ') c = c.substring(1,c.length);
		if (c.indexOf(nameEQ) == 0) return c.substring(nameEQ.length,c.length);
	}
	return null;
}
function eraseCookie(name) {
	createCookie(name,"",-1);
}

// Set up the websocket
if (! "WebSocket" in window) {
 // The browser doesn't support WebSocket
 alert("WebSocket NOT supported by your Browser!");
}

// Let us open a web socket
var ws = new WebSocket("` + wsurl.String() + `");
function say(txt) {
   ws.send(txt + '\n')
};
ws.onclose = function() {
   // websocket is closed.
   alert("Connection is closed..."); 
};
window.onpopstate = function(event) {
  say('path:' + window.location.href)
}

ws.onmessage = function (evt) {
   if (evt.data.replace(/^\s+|\s+$/g,"") == 'read-cookie') {
       var cookie = readCookie('WebSocketCookie');
       if (cookie != null) {
         say('cookie is ' + readCookie('WebSocketCookie'));
       } else {
         say('cookie is unknown')
       }
       return
   }
   if (evt.data.substr(0,5) == 'start') {
     say('path:' + window.location.href)
   }
   if (evt.data.substr(0,8) == 'setpath ') {
     history.pushState('', evt.data.substr(8), evt.data.substr(8));
   }
   if (evt.data.substr(0,12) == 'write-cookie') {
      createCookie('WebSocketCookie', evt.data.substr(12), 365);
      say('got cookie');
      return
   }
   var everything = document.getElementById("everything")
   if (everything == null) {
     return
   }
   var received_msg = evt.data;
   //alert("Message is received: " + received_msg);
   everything.innerHTML=received_msg;
};

</script>
</head>
<body>
<div id="everything">

  Everything goes here.
</div>
</body>
</html>
`
}
