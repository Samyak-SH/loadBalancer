const express = require("express")

const app1 = express();
const app2 = express();
const app3 = express();
const app4 = express();

app1.listen(3000, ()=>{console.log("started on 3000")})
app2.listen(3001, ()=>{console.log("started on 3001")})
app3.listen(3002, ()=>{console.log("started on 3002")})
app4.listen(3003, ()=>{console.log("started on 3003")})

app1.get("/",(req,res)=>{
    res.send("hello from server 1");
})
app2.get("/",(req,res)=>{
    res.send("hello from server 2");
})
app3.get("/",(req,res)=>{
    res.send("hello from server 3");
})
app4.get("/",(req,res)=>{
    res.send("hello from server 4");
})
