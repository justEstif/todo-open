#!/usr/bin/env node
"use strict";
const { spawnSync } = require("child_process");
const path = require("path");
const fs = require("fs");

const exe = process.platform === "win32" ? "todoopen-adapter-sync-git.exe" : "todoopen-adapter-sync-git";
const binPath = path.join(__dirname, "bin", exe);

if (!fs.existsSync(binPath)) {
  console.error(`todoopen-adapter-sync-git: binary not found at ${binPath}`);
  console.error("Try reinstalling: npm install -g @justestif/todoopen-adapter-sync-git");
  process.exit(1);
}

const result = spawnSync(binPath, process.argv.slice(2), { stdio: "inherit" });
process.exit(result.status ?? 1);
