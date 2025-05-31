const express = require("express");

const app = express();

app.use(express.json())
app.use(express.urlencoded({extended:true}));
app.get('/', (req,res)=>{
    res.header({"node-js-server" : "true"});
    res.send("hello from nodejs server")
});

app.listen(3000, ()=>{
    console.log("Server started on http://localhost:3000");
})