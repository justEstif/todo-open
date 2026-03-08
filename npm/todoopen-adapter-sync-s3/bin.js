#!/usr/bin/env node
"use strict";
const { spawnSync } = require("child_process");
const path = require("path");
const fs = require("fs");

const exe = process.platform === "win32" ? "todoopen-adapter-sync-s3.exe" : "todoopen-adapter-sync-s3";
const binPath = path.join(__dirname, "bin", exe);

if (!fs.existsSync(binPath)) {
  console.error(`todoopen-adapter-sync-s3: binary not found at ${binPath}`);
  console.error("Try reinstalling: npm install -g @justestif/todoopen-adapter-sync-s3");
  process.exit(1);
}

const result = spawnSync(binPath, process.argv.slice(2), { stdio: "inherit" });
process.exit(result.status ?? 1);
