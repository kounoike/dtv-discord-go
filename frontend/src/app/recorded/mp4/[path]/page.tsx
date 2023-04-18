"use client"
import Image from "next/image"
import styles from "./page.module.css"
import { AppBar, IconButton, Toolbar, Typography } from "@mui/material"
import { VideoPlayer } from "@videojs-player/react"
import "video.js/dist/video-js.css"

export default function SearchPage({ params }: { params: { path: string } }) {
  return (
    <>
      {false ? (
        <main className={styles.loadingMain}>
          <Image
            src="/shichou-chan-logo.png"
            alt="視聴ちゃんロゴ"
            className={styles.logo}
            width={1280}
            height={313}
            priority
          />
        </main>
      ) : (
        <>
          <AppBar position="sticky" className={styles.appbar}>
            <Toolbar>
              <Typography variant="h6" className={styles.title}>
                視聴ちゃん 録画視聴
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
            <VideoPlayer
              src={"/encoded/" + decodeURIComponent(params.path)}
              controls={true}
              loop={false}
              volume={0.6}
              autoplay={true}
              muted={false}
              playsinline={true}
              width={1920}
              height={1080}
            />
          </main>
        </>
      )}
    </>
  )
}
