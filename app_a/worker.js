const express = require("express");
const app = express();

app.get("/up", (req, res) => res.status(200).send("OK"));

app.listen(3100, "0.0.0.0", () => {
  console.log("Worker health server running...");
});

console.log("APP-A WORKER STARTED");

setInterval(() => {
  console.log("Worker heartbeat:", new Date().toISOString());
}, 5000);
