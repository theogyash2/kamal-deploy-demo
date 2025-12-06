const express = require("express");
const app = express();
const port = process.env.PORT || 3000;

app.get("/", (req, res) => {
  res.send("APP-A WEB :: UPDATED VERSION :: " + new Date().toISOString());
});

// healthcheck endpoint for kamal-proxy
app.get("/up", (req, res) => res.status(200).send("OK"));

app.listen(port, () => {
  console.log(`APP-A WEB running on port ${port}`);
});
