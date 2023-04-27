"use client"
import {
  AppBar,
  IconButton,
  MenuItem,
  Toolbar,
  Typography,
} from "@mui/material"
import styles from "./page.module.css"
import Image from "next/image"

export default function Page() {
  return (
    <>
      <AppBar position="sticky" className={styles.appbar}>
        <Toolbar>
          <Typography
            variant="h6"
            className={styles.title}
            sx={{ flexGrow: 1 }}
          >
            視聴ちゃん
          </Typography>
          <MenuItem onClick={() => (window.location.href = "/program/search")}>
            <Typography>→番組検索</Typography>
          </MenuItem>
          <MenuItem onClick={() => (window.location.href = "/recorded/search")}>
            <Typography>→録画検索</Typography>
          </MenuItem>
        </Toolbar>
      </AppBar>
      <main className={styles.main}>
        <Image
          src="/shichou-chan-logo.png"
          width={1280}
          height={313}
          alt="視聴ちゃんロゴ"
        ></Image>
      </main>
    </>
  )
}
