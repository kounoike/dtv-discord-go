"use client"
import Image from "next/image"
import { Inter } from "next/font/google"
import styles from "./page.module.css"
import { AppBar, Icon, IconButton, Toolbar, Typography } from "@mui/material"

const inter = Inter({ subsets: ["latin"] })

export default function Home() {
  return (
    <>
      <AppBar position="sticky">
        <Toolbar>
          <Typography variant="h6" className={styles.title}>
            視聴ちゃん
          </Typography>

          <div className={styles.grow} />
          <IconButton>
            <Image
              alt="視聴ちゃんアイコン"
              src="/web-icon.png"
              width={50}
              height={50}
            ></Image>
          </IconButton>
        </Toolbar>
      </AppBar>
      <main className={styles.main}>
        <Image
          src="/shichou-chan-logo.png"
          alt="視聴ちゃんロゴ"
          className={styles.logo}
          width={1280}
          height={313}
          priority
        />
      </main>
    </>
  )
}
