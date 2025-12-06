const express = require("express");
const app = express();
const port = process.env.PORT || 4000;

app.get("/", (req, res) => {
  res.send("APP-B WEB :: RUNNING :: " + new Date().toISOString());
});

// healthcheck endpoint
app.get("/up", (req, res) => res.status(200).send("OK"));

app.listen(port, () => {
  console.log(`APP-B running on port ${port}`);
});
