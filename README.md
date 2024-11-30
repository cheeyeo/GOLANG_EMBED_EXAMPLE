### GOLANG Websocket application

An example golang chat application built using Websockets and Gin router with ReactJS frontend UI and the gorilla websocket library.

This is based on the following tutorial:

https://www.jetbrains.com/guide/go/tutorials/webapp_go_react_part_one/


And also referenced in the blog post here:

https://www.cheeyeo.dev/golang/embed/reactjs/websocket/2024/11/07/golang-embed/


To really understand how this all works, please refer to the original jetbrains article and the gorilla/websocket broadcast example.


### Building

```
cd chatui

npm install

npm run build
```

The JS must be inside the `chatui/prod/` directory before proceeding to next step.

To build the binary:
```
make build
```

This bundles the JS from before into a single binary in `bin/chatapp`

To run:
```
bin/chatapp
```

NOTE:

This application is for experimental purposes only and not fit for production usage as there's no security involved. Use at your own risk.

### Local Development

The golang app can run in a terminal on port 8080 via:
```
go run main.go
```

The UI is a ReactJS application which uses `http-proxy-middleware` and has a setupProxy.js file setup pointing to localhost:8080

Run the UI in a separate terminal:
```
npm start
```