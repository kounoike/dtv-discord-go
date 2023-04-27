"use client"
import Image from "next/image"
import styles from "./page.module.css"
import {
  AppBar,
  IconButton,
  MenuItem,
  Toolbar,
  Typography,
} from "@mui/material"
import { VideoPlayer } from "@videojs-player/react"
import "video.js/dist/video-js.css"
import { useEffect, useState } from "react"

export default function SearchPage() {
  const [mp4Path, setMp4Path] = useState<string | undefined>(undefined)
  useEffect(() => {
    const searchParams = new URLSearchParams(window.location.search)
    const f = searchParams.get("f")
    if (f !== null) {
      setMp4Path(f)
    }
  }, [])

  return mp4Path === undefined ? (
    <></>
  ) : (
    <>
      <AppBar position="sticky" className={styles.appbar}>
        <Toolbar>
          <Typography
            variant="h6"
            className={styles.title}
            sx={{ flexGrow: 1 }}
          >
            視聴ちゃん 録画視聴 {decodeURIComponent(mp4Path)}
          </Typography>
          <MenuItem onClick={() => (window.location.href = "/program/search")}>
            <Typography>→番組検索</Typography>
          </MenuItem>
          <MenuItem onClick={() => (window.location.href = "/recorded/search")}>
            <Typography>→録画検索</Typography>
          </MenuItem>
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
          src={"/encoded/" + mp4Path}
          controls={true}
          loop={false}
          volume={0.6}
          autoplay={true}
          muted={false}
          playsinline={true}
          // width={1920}
          // height={1080}
          aspectRatio="16:9"
          responsive={true}
        />
      </main>
    </>
  )
}
